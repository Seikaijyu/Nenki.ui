// SPDX-License-Identifier: Unlicense OR MIT

package app

import (
	"errors"
	"fmt"
	"image"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	syscall "golang.org/x/sys/windows"

	"gioui.org/app/internal/windows"
	"gioui.org/unit"
	gowindows "golang.org/x/sys/windows"

	"gioui.org/f32"
	"gioui.org/io/clipboard"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
)

// ViewEvent 结构体包含了一个 HWND，这是一个窗口的句柄
type ViewEvent struct {
	HWND uintptr
}

// window 结构体定义了一个窗口的各种属性
type window struct {
	hwnd        syscall.Handle  // 窗口的句柄
	hdc         syscall.Handle  // 设备上下文的句柄
	w           *callbacks      // 回调函数的集合
	stage       system.Stage    // 系统的阶段
	pointerBtns pointer.Buttons // 指针按钮的状态

	// cursorIn 标记鼠标光标是否在窗口内，根据最近的 WM_SETCURSOR 消息来判断
	cursorIn bool
	cursor   syscall.Handle // 光标的句柄

	// placement 在全屏模式下保存上一次窗口的位置
	placement *windows.WindowPlacement

	animating bool // 标记窗口是否在动画中
	focused   bool // 标记窗口是否被聚焦

	borderSize image.Point // 窗口边框的大小
	config     Config      // 窗口的配置信息
}

// _WM_WAKEUP 是一个自定义的 Windows 消息，用于唤醒窗口
const _WM_WAKEUP = windows.WM_USER + iota

// gpuAPI 结构体定义了一个 GPU API，包含了优先级和初始化函数
type gpuAPI struct {
	priority    int                              // 优先级
	initializer func(w *window) (context, error) // 初始化函数
}

// drivers 是一个存放所有可能的 Context 实现的列表
var drivers []gpuAPI

// winMap 是一个映射，将 win32 的 HWND 映射到 *windows
var winMap sync.Map

// iconID 是资源文件中图标的 ID
const iconID = 1

// resources 结构体包含了一些资源，如模块句柄、窗口类和光标资源
var resources struct {
	once   sync.Once      // 用于只执行一次的同步标记
	handle syscall.Handle // 从 GetModuleHandle 获取的模块句柄
	class  uint16         // 从 RegisterClassEx 注册的 Gio 窗口类
	cursor syscall.Handle // 光标资源的句柄
}

// osMain 函数是操作系统主函数，这里是一个空实现
func osMain() {
	select {}
}

// newWindow 函数用于创建一个新的窗口
func newWindow(window *callbacks, options []Option) error {
	// 创建一个错误通道
	cerr := make(chan error)
	// 在一个新的协程中创建窗口
	go func() {
		// GetMessage 和 PeekMessage 可以基于窗口的 HWND 进行过滤，但是
		// 这样会忽略掉特定于线程的消息，如 WM_QUIT。
		// 因此，我们锁定线程，让窗口消息通过未经过滤的 GetMessage 调用到达。
		runtime.LockOSThread()
		// 创建一个原生窗口
		w, err := createNativeWindow()
		// 如果创建窗口时出错，将错误发送到错误通道并返回
		if err != nil {
			cerr <- err
			return
		}
		// 将 nil 错误发送到错误通道
		cerr <- nil
		// 将创建的窗口存储到 winMap 中
		winMap.Store(w.hwnd, w)
		// 在函数返回时从 winMap 中删除窗口
		defer winMap.Delete(w.hwnd)
		// 设置窗口的回调函数
		w.w = window
		// 设置窗口的驱动程序
		w.w.SetDriver(w)
		// 发送一个 ViewEvent 事件
		w.w.Event(ViewEvent{HWND: uintptr(w.hwnd)})
		// 配置窗口
		w.Configure(options)
		// 将窗口设置为前台窗口
		windows.SetForegroundWindow(w.hwnd)
		// 设置窗口的焦点
		windows.SetFocus(w.hwnd)
		// 由于光标的窗口类是空的，
		// 所以在这里设置它以显示光标。
		w.SetCursor(pointer.CursorDefault)
		// 进入窗口的消息循环
		if err := w.loop(); err != nil {
			// 如果消息循环出错，抛出 panic
			panic(err)
		}
	}()
	// 返回错误通道中的错误
	return <-cerr
}

// initResources 函数用于初始化全局的 resources。
func initResources() error {
	// 设置当前进程为 DPI Aware，使得窗口在高 DPI 设置下能够正确显示
	windows.SetProcessDPIAware()
	// 获取当前模块的句柄
	hInst, err := windows.GetModuleHandle()
	if err != nil {
		// 如果获取失败，返回错误
		return err
	}
	// 将获取到的模块句柄赋值给 resources 的 handle 字段
	resources.handle = hInst
	// 加载系统预定义的光标
	c, err := windows.LoadCursor(windows.IDC_ARROW)
	if err != nil {
		// 如果加载失败，返回错误
		return err
	}
	// 将加载到的光标赋值给 resources 的 cursor 字段
	resources.cursor = c
	// 加载图标资源
	icon, _ := windows.LoadImage(hInst, iconID, windows.IMAGE_ICON, 0, 0, windows.LR_DEFAULTSIZE|windows.LR_SHARED)
	// 定义窗口类
	wcls := windows.WndClassEx{
		CbSize:        uint32(unsafe.Sizeof(windows.WndClassEx{})),                // 结构体的大小
		Style:         windows.CS_HREDRAW | windows.CS_VREDRAW | windows.CS_OWNDC, // 窗口样式
		LpfnWndProc:   syscall.NewCallback(windowProc),                            // 窗口过程函数
		HInstance:     hInst,                                                      // 模块句柄
		HIcon:         icon,                                                       // 图标
		LpszClassName: syscall.StringToUTF16Ptr("GioWindow"),                      // 窗口类名
	}
	// 注册窗口类
	cls, err := windows.RegisterClassEx(&wcls)
	if err != nil {
		// 如果注册失败，返回错误
		return err
	}
	// 将注册到的窗口类赋值给 resources 的 class 字段
	resources.class = cls
	// 如果所有操作都成功，返回 nil 表示没有错误
	return nil
}

// 定义窗口的扩展样式
const dwExStyle = windows.WS_EX_APPWINDOW | windows.WS_EX_WINDOWEDGE

// 窗口ID，默认为 0
var HWND syscall.Handle = 0

// createNativeWindow 函数用于创建一个本地窗口
func createNativeWindow() (*window, error) {
	var resErr error
	// 使用 sync.Once 确保全局的 resources 只被初始化一次
	resources.once.Do(func() {
		resErr = initResources()
	})
	// 如果初始化 resources 时出错，返回错误
	if resErr != nil {
		return nil, resErr
	}
	// 定义窗口的样式
	const dwStyle = windows.WS_OVERLAPPEDWINDOW

	// 调用 CreateWindowEx 函数创建窗口
	hwnd, err := windows.CreateWindowEx(
		dwExStyle,       // 窗口的扩展样式
		resources.class, // 窗口类
		"",              // 窗口标题
		dwStyle|windows.WS_CLIPSIBLINGS|windows.WS_CLIPCHILDREN, // 窗口样式
		windows.CW_USEDEFAULT, windows.CW_USEDEFAULT, // 窗口的初始位置
		windows.CW_USEDEFAULT, windows.CW_USEDEFAULT, // 窗口的初始大小
		0,                // 父窗口的句柄
		0,                // 菜单的句柄
		resources.handle, // 应用程序实例的句柄
		0)                // 创建窗口的参数
	// 如果创建窗口时出错，返回错误
	if err != nil {
		return nil, err
	}
	// 创建 window 结构体实例
	w := &window{
		hwnd: hwnd,
	}

	HWND = hwnd
	// 获取窗口的设备上下文
	w.hdc, err = windows.GetDC(hwnd)
	// 如果获取设备上下文时出错，返回错误
	if err != nil {
		return nil, err
	}
	// 如果所有操作都成功，返回创建的窗口
	return w, nil
}

// 拖动处理函数
var dragHandler func([]string) = func(s []string) {}

// 窗口是否接受文件拖放
func DragAcceptFiles(hwnd syscall.Handle, accept bool) {
	windows.DragAcceptFiles(hwnd, accept)
}

// 自定义处理文件拖放事件
func CustomDragHandler(fn func([]string)) {
	dragHandler = fn
}

// update() 函数用于处理用户所做的更改，并更新配置。
// 它会读取窗口的样式以及大小/位置，并更新 w.config。
// 如果有任何更改，它会发出一个 ConfigEvent 来通知应用程序。
func (w *window) update() {
	// 获取窗口的客户区大小
	cr := windows.GetClientRect(w.hwnd)
	// 更新窗口的大小
	w.config.Size = image.Point{
		// 宽度为客户区的右边界减去左边界
		X: int(cr.Right - cr.Left),
		// 高度为客户区的下边界减去上边界
		Y: int(cr.Bottom - cr.Top),
	}

	// 获取窗口边框的大小
	w.borderSize = image.Pt(
		// 边框的宽度为系统的 SM_CXSIZEFRAME 参数
		windows.GetSystemMetrics(windows.SM_CXSIZEFRAME),
		// 边框的高度为系统的 SM_CYSIZEFRAME 参数
		windows.GetSystemMetrics(windows.SM_CYSIZEFRAME),
	)
	// 发出 ConfigEvent 事件，通知应用程序窗口的配置已经更改
	w.w.Event(ConfigEvent{Config: w.config})
}

// windowProc 函数是窗口过程函数，用于处理窗口接收到的消息。
// 它接收四个参数：窗口句柄、消息、以及两个消息参数。
func windowProc(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	// 从 winMap 中获取窗口句柄对应的窗口
	// winMap 是一个 map，键为窗口句柄，值为窗口对象
	win, exists := winMap.Load(hwnd)
	// 如果 winMap 中不存在该窗口句柄，说明这是一个未注册的窗口
	// 对于未注册的窗口，我们直接调用 DefWindowProc 函数进行默认处理
	if !exists {
		return windows.DefWindowProc(hwnd, msg, wParam, lParam)
	}

	// 如果 winMap 中存在该窗口句柄，我们将其转换为 *window 类型
	// 并赋值给 w，以便后续处理
	w := win.(*window)
	// 根据消息的类型进行不同的处理
	switch msg {
	case windows.WM_DROPFILES:
		files := windows.DragFiles(uintptr(wParam))
		if files != nil {
			// 如果接收到的是 WM_DROPFILES 消息，处理文件拖放事件
			dragHandler(files)
		}
		// 释放拖放操作
		windows.DragFinish(wParam)
		// 消息已经被处理
		return windows.TRUE
	case windows.WM_UNICHAR:
		// 如果接收到的是 WM_UNICHAR 消息
		if wParam == windows.UNICODE_NOCHAR {
			// 如果参数表示没有字符，那么告诉系统我们接受 WM_UNICHAR 消息
			return windows.TRUE
		}
		fallthrough
	case windows.WM_CHAR:
		// 如果接收到的是 WM_CHAR 消息
		if r := rune(wParam); unicode.IsPrint(r) {
			// 如果参数是可打印的字符，那么在编辑器中插入该字符
			w.w.EditorInsert(string(r))
		}
		// 消息已经被处理
		return windows.TRUE
	case windows.WM_DPICHANGED:
		// 如果接收到的是 WM_DPICHANGED 消息，告诉 Windows 我们已经准备好进行运行时 DPI 的改变
		return windows.TRUE
	case windows.WM_ERASEBKGND:
		// 如果接收到的是 WM_ERASEBKGND 消息，为了避免 GPU 内容和背景颜色之间的闪烁，返回 TRUE
		return windows.TRUE
	case windows.WM_KEYDOWN, windows.WM_KEYUP, windows.WM_SYSKEYDOWN, windows.WM_SYSKEYUP:
		// 如果接收到的是键盘按下或释放的消息
		if n, ok := convertKeyCode(wParam); ok {
			// 如果参数是有效的键码，那么创建一个键盘事件
			e := key.Event{
				Name:      n,
				Modifiers: getModifiers(),
				State:     key.Press,
			}
			// 如果消息是键盘释放，那么修改事件的状态
			if msg == windows.WM_KEYUP || msg == windows.WM_SYSKEYUP {
				e.State = key.Release
			}

			// 发出键盘事件
			w.w.Event(e)

			if (wParam == windows.VK_F10) && (msg == windows.WM_SYSKEYDOWN || msg == windows.WM_SYSKEYUP) {
				// 如果按下的是 F10 键，那么保留它，不让它打开系统菜单。其他 Windows 程序
				// 如 cmd.exe 和图形调试器也会保留 F10 键。
				return 0
			}
		}
	case windows.WM_LBUTTONDOWN:
		// 如果接收到的是鼠标左键按下的消息，处理鼠标按键事件
		w.pointerButton(pointer.ButtonPrimary, true, lParam, getModifiers())
	case windows.WM_LBUTTONUP:
		// 如果接收到的是鼠标左键释放的消息，处理鼠标按键事件
		w.pointerButton(pointer.ButtonPrimary, false, lParam, getModifiers())
	case windows.WM_RBUTTONDOWN:
		// 如果接收到的是鼠标右键按下的消息，处理鼠标按键事件
		w.pointerButton(pointer.ButtonSecondary, true, lParam, getModifiers())
	case windows.WM_RBUTTONUP:
		// 如果接收到的是鼠标右键释放的消息，处理鼠标按键事件
		w.pointerButton(pointer.ButtonSecondary, false, lParam, getModifiers())
	case windows.WM_MBUTTONDOWN:
		// 如果接收到的是鼠标中键按下的消息，处理鼠标按键事件
		w.pointerButton(pointer.ButtonTertiary, true, lParam, getModifiers())
	case windows.WM_MBUTTONUP:
		// 如果接收到的是鼠标中键释放的消息，处理鼠标按键事件
		w.pointerButton(pointer.ButtonTertiary, false, lParam, getModifiers())
	case windows.WM_CANCELMODE:
		// 如果接收到的是 WM_CANCELMODE 消息，发出一个取消的鼠标事件
		w.w.Event(pointer.Event{
			Kind: pointer.Cancel,
		})
	case windows.WM_SETFOCUS:
		// 如果接收到的是 WM_SETFOCUS 消息，设置窗口为获取焦点状态，并发出一个焦点事件
		w.focused = true
		w.w.Event(key.FocusEvent{Focus: true})
	case windows.WM_KILLFOCUS:
		// 如果接收到的是 WM_KILLFOCUS 消息，设置窗口为失去焦点状态，并发出一个焦点事件
		w.focused = false
		w.w.Event(key.FocusEvent{Focus: false})
	case windows.WM_NCACTIVATE:
		// 如果接收到的是 WM_NCACTIVATE 消息，根据 wParam 的值改变窗口的状态
		if w.stage >= system.StageInactive {
			if wParam == windows.TRUE {
				w.setStage(system.StageRunning)
			} else {
				w.setStage(system.StageInactive)
			}
		}
	case windows.WM_NCHITTEST:
		// 如果接收到的是 WM_NCHITTEST 消息，如果窗口是装饰的，让系统处理它
		if w.config.Decorated {
			// 让系统处理它
			break
		}
		// 否则，将 lParam 转换为坐标，并进行命中测试
		x, y := coordsFromlParam(lParam)
		np := windows.Point{X: int32(x), Y: int32(y)}
		windows.ScreenToClient(w.hwnd, &np)
		return w.hitTest(int(np.X), int(np.Y))
	case windows.WM_MOUSEMOVE:
		// 如果接收到的是 WM_MOUSEMOVE 消息，将 lParam 转换为坐标，并发出一个鼠标移动事件
		x, y := coordsFromlParam(lParam)
		p := f32.Point{X: float32(x), Y: float32(y)}
		w.w.Event(pointer.Event{
			Kind:      pointer.Move,
			Source:    pointer.Mouse,
			Position:  p,
			Buttons:   w.pointerBtns,
			Time:      windows.GetMessageTime(),
			Modifiers: getModifiers(),
		})
	case windows.WM_MOUSEWHEEL:
		// 如果接收到的是 WM_MOUSEWHEEL 消息，处理鼠标滚轮事件
		w.scrollEvent(wParam, lParam, false, getModifiers())
	case windows.WM_MOUSEHWHEEL:
		// 如果接收到的是 WM_MOUSEHWHEEL 消息，处理鼠标水平滚轮事件
		w.scrollEvent(wParam, lParam, true, getModifiers())
	case windows.WM_DESTROY:
		// 如果接收到的是 WM_DESTROY 消息，发出一个视图事件和一个销毁事件
		w.w.Event(ViewEvent{})
		w.w.Event(system.DestroyEvent{})
		// 如果设备上下文句柄不为 0，释放它
		if w.hdc != 0 {
			windows.ReleaseDC(w.hdc)
			w.hdc = 0
		}
		// 系统会为我们销毁窗口句柄
		w.hwnd = 0
		// 发送一个退出消息
		windows.PostQuitMessage(0)
	case windows.WM_NCCALCSIZE:
		// 如果接收到的是 WM_NCCALCSIZE 消息，如果窗口是装饰的，让 Windows 处理装饰
		if w.config.Decorated {
			// 让 Windows 处理装饰
			break
		}
		// 没有客户区域；我们自己绘制装饰
		if wParam != 1 {
			return 0
		}
		// lParam 包含一个 NCCALCSIZE_PARAMS，我们可以调整它
		place := windows.GetWindowPlacement(w.hwnd)
		if !place.IsMaximized() {
			// 没有需要调整的
			return 0
		}
		// 调整窗口位置以避免在最大化状态下的额外填充
		// 参见 https://devblogs.microsoft.com/oldnewthing/20150304-00/?p=44543
		// 注意，试图在 WM_GETMINMAXINFO 中进行调整会被 Windows 忽略
		szp := (*windows.NCCalcSizeParams)(unsafe.Pointer(uintptr(lParam)))
		mi := windows.GetMonitorInfo(w.hwnd)
		szp.Rgrc[0] = mi.WorkArea
		return 0
	case windows.WM_PAINT:
		// 如果接收到的是 WM_PAINT 消息，执行绘制操作
		w.draw(true)
	case windows.WM_SIZE:
		// 如果接收到的是 WM_SIZE 消息，更新窗口大小
		w.update()
		switch wParam {
		case windows.SIZE_MINIMIZED:
			// 如果窗口被最小化，设置窗口模式为最小化，并将窗口状态设置为暂停
			w.config.Mode = Minimized
			w.setStage(system.StagePaused)
		case windows.SIZE_MAXIMIZED:
			// 如果窗口被最大化，设置窗口模式为最大化，并将窗口状态设置为运行
			w.config.Mode = Maximized
			w.setStage(system.StageRunning)
		case windows.SIZE_RESTORED:
			// 如果窗口被恢复到正常大小，设置窗口模式为窗口化（除非当前模式是全屏），并将窗口状态设置为运行
			if w.config.Mode != Fullscreen {
				w.config.Mode = Windowed
			}
			w.setStage(system.StageRunning)
		}
	case windows.WM_GETMINMAXINFO:
		// 如果接收到的是 WM_GETMINMAXINFO 消息，获取窗口的最小和最大尺寸信息
		mm := (*windows.MinMaxInfo)(unsafe.Pointer(uintptr(lParam)))
		var bw, bh int32
		if w.config.Decorated {
			// 如果窗口是装饰的，获取窗口和客户区的尺寸，计算边框的宽度和高度
			r := windows.GetWindowRect(w.hwnd)
			cr := windows.GetClientRect(w.hwnd)
			bw = r.Right - r.Left - (cr.Right - cr.Left)
			bh = r.Bottom - r.Top - (cr.Bottom - cr.Top)
		}
		if p := w.config.MinSize; p.X > 0 || p.Y > 0 {
			// 如果设置了窗口的最小尺寸，将其加上边框的尺寸，设置为窗口的最小跟踪尺寸
			mm.PtMinTrackSize = windows.Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		if p := w.config.MaxSize; p.X > 0 || p.Y > 0 {
			// 如果设置了窗口的最大尺寸，将其加上边框的尺寸，设置为窗口的最大跟踪尺寸
			mm.PtMaxTrackSize = windows.Point{
				X: int32(p.X) + bw,
				Y: int32(p.Y) + bh,
			}
		}
		return 0
	case windows.WM_SETCURSOR:
		// 如果接收到的是 WM_SETCURSOR 消息，检查光标是否在客户区内
		w.cursorIn = (lParam & 0xffff) == windows.HTCLIENT
		if w.cursorIn {
			// 如果光标在客户区内，设置光标
			windows.SetCursor(w.cursor)
			return windows.TRUE
		}
	case _WM_WAKEUP:
		// 如果接收到的是 _WM_WAKEUP 消息，触发唤醒事件
		w.w.Event(wakeupEvent{})
	case windows.WM_IME_STARTCOMPOSITION:
		// 如果接收到的是 WM_IME_STARTCOMPOSITION 消息，开始输入法编辑
		imc := windows.ImmGetContext(w.hwnd)
		if imc == 0 {
			// 如果无法获取输入法上下文，返回 TRUE
			return windows.TRUE
		}
		defer windows.ImmReleaseContext(w.hwnd, imc)
		// 获取编辑器的选择状态
		sel := w.w.EditorState().Selection
		// 转换选择的光标位置
		caret := sel.Transform.Transform(sel.Caret.Pos.Add(f32.Pt(0, sel.Caret.Descent)))
		icaret := image.Pt(int(caret.X+.5), int(caret.Y+.5))
		// 设置输入法的组合窗口和候选窗口位置
		windows.ImmSetCompositionWindow(imc, icaret.X, icaret.Y)
		windows.ImmSetCandidateWindow(imc, icaret.X, icaret.Y)
	case windows.WM_IME_COMPOSITION:
		// 如果接收到的是 WM_IME_COMPOSITION 消息，进行输入法编辑
		imc := windows.ImmGetContext(w.hwnd)
		if imc == 0 {
			// 如果无法获取输入法上下文，返回 TRUE
			return windows.TRUE
		}
		defer windows.ImmReleaseContext(w.hwnd, imc)
		// 获取编辑器状态
		state := w.w.EditorState()
		rng := state.compose
		if rng.Start == -1 {
			rng = state.Selection.Range
		}
		if rng.Start > rng.End {
			rng.Start, rng.End = rng.End, rng.Start
		}
		var replacement string
		switch {
		case lParam&windows.GCS_RESULTSTR != 0:
			// 如果 lParam 的 GCS_RESULTSTR 位被设置，获取输入法的结果字符串
			replacement = windows.ImmGetCompositionString(imc, windows.GCS_RESULTSTR)
		case lParam&windows.GCS_COMPSTR != 0:
			// 如果 lParam 的 GCS_COMPSTR 位被设置，获取输入法的组合字符串
			replacement = windows.ImmGetCompositionString(imc, windows.GCS_COMPSTR)
		}
		// 替换编辑器的内容
		end := rng.Start + utf8.RuneCountInString(replacement)
		w.w.EditorReplace(rng, replacement)
		state = w.w.EditorState()
		comp := key.Range{
			Start: rng.Start,
			End:   end,
		}
		if lParam&windows.GCS_DELTASTART != 0 {
			start := windows.ImmGetCompositionValue(imc, windows.GCS_DELTASTART)
			comp.Start = state.RunesIndex(state.UTF16Index(comp.Start) + start)
		}
		// 设置编辑器的组合区域
		w.w.SetComposingRegion(comp)
		pos := end
		if lParam&windows.GCS_CURSORPOS != 0 {
			rel := windows.ImmGetCompositionValue(imc, windows.GCS_CURSORPOS)
			pos = state.RunesIndex(state.UTF16Index(rng.Start) + rel)
		}
		// 设置编辑器的选择区域
		w.w.SetEditorSelection(key.Range{Start: pos, End: pos})
		return windows.TRUE
	case windows.WM_IME_ENDCOMPOSITION:
		// 如果接收到的是 WM_IME_ENDCOMPOSITION 消息，结束输入法编辑
		w.w.SetComposingRegion(key.Range{Start: -1, End: -1})
		return windows.TRUE
	}

	// 如果没有匹配的消息处理，调用默认的窗口处理函数处理消息
	return windows.DefWindowProc(hwnd, msg, wParam, lParam)
}

// getModifiers 函数用于获取当前按下的修饰键（如Ctrl、Alt等）的状态
func getModifiers() key.Modifiers {
	var kmods key.Modifiers
	// 检查左右 Win 键是否被按下
	if windows.GetKeyState(windows.VK_LWIN)&0x1000 != 0 || windows.GetKeyState(windows.VK_RWIN)&0x1000 != 0 {
		kmods |= key.ModSuper
	}
	// 检查 Alt 键是否被按下
	if windows.GetKeyState(windows.VK_MENU)&0x1000 != 0 {
		kmods |= key.ModAlt
	}
	// 检查 Ctrl 键是否被按下
	if windows.GetKeyState(windows.VK_CONTROL)&0x1000 != 0 {
		kmods |= key.ModCtrl
	}
	// 检查 Shift 键是否被按下
	if windows.GetKeyState(windows.VK_SHIFT)&0x1000 != 0 {
		kmods |= key.ModShift
	}
	// 返回按键状态
	return kmods
}

// hitTest 函数用于检测鼠标点击的位置是否在非客户区，
// 该函数主要用于处理 WM_NCHITTEST 消息。
func (w *window) hitTest(x, y int) uintptr {
	// 如果窗口处于全屏模式，则所有点击都视为在客户区内
	if w.config.Mode == Fullscreen {
		return windows.HTCLIENT
	}
	// 如果窗口不处于窗口模式，则不允许调整窗口大小
	if w.config.Mode != Windowed {
		return windows.HTCLIENT
	}
	// 检查鼠标是否在窗口的边缘，用于调整窗口大小
	top := y <= w.borderSize.Y
	bottom := y >= w.config.Size.Y-w.borderSize.Y
	left := x <= w.borderSize.X
	right := x >= w.config.Size.X-w.borderSize.X
	switch {
	case top && left:
		return windows.HTTOPLEFT
	case top && right:
		return windows.HTTOPRIGHT
	case bottom && left:
		return windows.HTBOTTOMLEFT
	case bottom && right:
		return windows.HTBOTTOMRIGHT
	case top:
		return windows.HTTOP
	case bottom:
		return windows.HTBOTTOM
	case left:
		return windows.HTLEFT
	case right:
		return windows.HTRIGHT
	}
	// 检查鼠标是否在窗口的移动区域
	p := f32.Pt(float32(x), float32(y))
	if a, ok := w.w.ActionAt(p); ok && a == system.ActionMove {
		return windows.HTCAPTION
	}
	// 其他情况，视为在客户区内
	return windows.HTCLIENT
}

// pointerButton 函数处理鼠标按钮的按下和释放事件
func (w *window) pointerButton(btn pointer.Buttons, press bool, lParam uintptr, kmods key.Modifiers) {
	// 如果窗口没有焦点，设置焦点到该窗口
	if !w.focused {
		windows.SetFocus(w.hwnd)
	}

	var kind pointer.Kind
	// 如果是按下事件
	if press {
		kind = pointer.Press
		// 如果没有其他按钮被按下，获取鼠标捕获
		if w.pointerBtns == 0 {
			windows.SetCapture(w.hwnd)
		}
		// 更新按下按钮的状态
		w.pointerBtns |= btn
	} else {
		// 如果是释放事件
		kind = pointer.Release
		// 更新按下按钮的状态
		w.pointerBtns &^= btn
		// 如果所有按钮都已释放，释放鼠标捕获
		if w.pointerBtns == 0 {
			windows.ReleaseCapture()
		}
	}
	// 从 lParam 中获取鼠标的坐标
	x, y := coordsFromlParam(lParam)
	p := f32.Point{X: float32(x), Y: float32(y)}
	// 发送鼠标事件
	w.w.Event(pointer.Event{
		Kind:      kind,
		Source:    pointer.Mouse,
		Position:  p,
		Buttons:   w.pointerBtns,
		Time:      windows.GetMessageTime(),
		Modifiers: kmods,
	})
}

// coordsFromlParam 函数从 lParam 中解析出鼠标的坐标
func coordsFromlParam(lParam uintptr) (int, int) {
	x := int(int16(lParam & 0xffff))
	y := int(int16((lParam >> 16) & 0xffff))
	return x, y
}

// scrollEvent 函数处理鼠标滚轮事件
func (w *window) scrollEvent(wParam, lParam uintptr, horizontal bool, kmods key.Modifiers) {
	// 从 lParam 中获取鼠标的屏幕坐标
	x, y := coordsFromlParam(lParam)
	// 将屏幕坐标转换为客户区坐标
	np := windows.Point{X: int32(x), Y: int32(y)}
	windows.ScreenToClient(w.hwnd, &np)
	p := f32.Point{X: float32(np.X), Y: float32(np.Y)}
	// 获取滚动的距离
	dist := float32(int16(wParam >> 16))
	var sp f32.Point
	// 如果是水平滚动
	if horizontal {
		sp.X = dist
	} else {
		// 如果按下 Shift 键，支持水平滚动
		if kmods == key.ModShift {
			sp.X = -dist
		} else {
			// 否则为垂直滚动
			sp.Y = -dist
		}
	}
	// 发送鼠标滚轮事件
	w.w.Event(pointer.Event{
		Kind:      pointer.Scroll,
		Source:    pointer.Mouse,
		Position:  p,
		Buttons:   w.pointerBtns,
		Scroll:    sp,
		Modifiers: kmods,
		Time:      windows.GetMessageTime(),
	})
}

// loop 函数是窗口的消息循环
// 该函数参考了 https://blogs.msdn.microsoft.com/oldnewthing/20060126-00/?p=32513/
func (w *window) loop() error {
	msg := new(windows.Msg)
loop:
	for {
		anim := w.animating
		// 如果窗口正在动画中，并且没有待处理的消息，绘制窗口
		if anim && !windows.PeekMessage(msg, 0, 0, 0, windows.PM_NOREMOVE) {
			w.draw(false)
			continue
		}
		// 获取消息
		switch ret := windows.GetMessage(msg, 0, 0, 0); ret {
		case -1:
			// 如果 GetMessage 返回 -1，表示出错
			return errors.New("GetMessage failed")
		case 0:
			// 如果 GetMessage 返回 0，表示接收到 WM_QUIT 消息，退出消息循环
			break loop
		}
		// 转换和分发消息
		windows.TranslateMessage(msg)
		windows.DispatchMessage(msg)
	}
	return nil
}

// EditorStateChanged 方法用于处理编辑器状态变化
// 当编辑器的选区或者代码片段发生变化时，会取消当前的输入法编辑
func (w *window) EditorStateChanged(old, new editorState) {
	// 获取当前窗口的输入法上下文
	imc := windows.ImmGetContext(w.hwnd)
	if imc == 0 {
		return
	}
	// 确保输入法上下文在函数返回时被释放
	defer windows.ImmReleaseContext(w.hwnd, imc)
	// 如果选区或者代码片段发生了变化
	if old.Selection.Range != new.Selection.Range || old.Snippet != new.Snippet {
		// 取消当前的输入法编辑
		windows.ImmNotifyIME(imc, windows.NI_COMPOSITIONSTR, windows.CPS_CANCEL, 0)
	}
}

// SetAnimating 方法用于设置窗口是否处于动画状态
func (w *window) SetAnimating(anim bool) {
	w.animating = anim
}

// Wakeup 方法用于唤醒窗口
// 它会向窗口发送一个 _WM_WAKEUP 消息
func (w *window) Wakeup() {
	if err := windows.PostMessage(w.hwnd, _WM_WAKEUP, 0, 0); err != nil {
		panic(err)
	}
}

// setStage 方法用于设置窗口的阶段
// 如果阶段发生了变化，它会发送一个 StageEvent 事件
func (w *window) setStage(s system.Stage) {
	if s != w.stage {
		w.stage = s
		w.w.Event(system.StageEvent{Stage: s})
	}
}

// draw 方法用于绘制窗口
// 如果窗口的大小为 0，它会直接返回
// 否则，它会根据窗口的 DPI 创建一个配置，并发送一个 frameEvent 事件
func (w *window) draw(sync bool) {
	if w.config.Size.X == 0 || w.config.Size.Y == 0 {
		return
	}
	dpi := windows.GetWindowDPI(w.hwnd)
	cfg := configForDPI(dpi)
	w.w.Event(frameEvent{
		FrameEvent: system.FrameEvent{
			Now:    time.Now(),
			Size:   w.config.Size,
			Metric: cfg,
		},
		Sync: sync,
	})
}

// NewContext 方法用于创建一个新的上下文
// 它会按照优先级顺序尝试所有的驱动程序，直到成功创建一个上下文
// 如果所有的驱动程序都无法创建上下文，它会返回一个错误
func (w *window) NewContext() (context, error) {
	// 按照优先级对驱动程序进行排序
	sort.Slice(drivers, func(i, j int) bool {
		return drivers[i].priority < drivers[j].priority
	})
	// 用于保存每个驱动程序的错误信息
	var errs []string
	// 遍历所有的驱动程序
	for _, b := range drivers {
		// 尝试使用当前驱动程序创建上下文
		ctx, err := b.initializer(w)
		// 如果创建成功，返回创建的上下文
		if err == nil {
			return ctx, nil
		}
		// 如果创建失败，保存错误信息
		errs = append(errs, err.Error())
	}
	// 如果所有的驱动程序都无法创建上下文，返回一个错误
	if len(errs) > 0 {
		return nil, fmt.Errorf("NewContext: failed to create a GPU device, tried: %s", strings.Join(errs, ", "))
	}
	return nil, errors.New("NewContext: no available GPU drivers")
}

// ReadClipboard 方法用于读取剪贴板的内容
func (w *window) ReadClipboard() {
	w.readClipboard()
}

// readClipboard 方法用于读取剪贴板的内容
// 它会打开剪贴板，获取剪贴板中的数据，然后发送一个剪贴板事件
func (w *window) readClipboard() error {
	// 打开剪贴板
	if err := windows.OpenClipboard(w.hwnd); err != nil {
		return err
	}
	// 确保在函数返回时关闭剪贴板
	defer windows.CloseClipboard()
	// 获取剪贴板中的数据
	mem, err := windows.GetClipboardData(windows.CF_UNICODETEXT)
	if err != nil {
		return err
	}
	// 锁定内存，获取数据的指针
	ptr, err := windows.GlobalLock(mem)
	if err != nil {
		return err
	}
	// 确保在函数返回时解锁内存
	defer windows.GlobalUnlock(mem)
	// 将数据转换为字符串
	content := gowindows.UTF16PtrToString((*uint16)(unsafe.Pointer(ptr)))
	// 发送剪贴板事件
	w.w.Event(clipboard.Event{Text: content})
	return nil
}

// Configure 方法用于配置窗口
// 它会根据给定的选项来设置窗口的各项参数
func (w *window) Configure(options []Option) {
	// 获取系统的 DPI
	dpi := windows.GetSystemDPI()
	// 根据 DPI 创建一个配置
	metric := configForDPI(dpi)
	// 应用配置
	w.config.apply(metric, options)
	// 设置窗口的标题
	windows.SetWindowText(w.hwnd, w.config.Title)

	// 获取窗口的样式
	style := windows.GetWindowLong(w.hwnd, windows.GWL_STYLE)
	var showMode int32
	var x, y, width, height int32
	swpStyle := uintptr(windows.SWP_NOZORDER | windows.SWP_FRAMECHANGED)
	winStyle := uintptr(windows.WS_OVERLAPPEDWINDOW)
	style &^= winStyle
	// 根据窗口的模式来设置窗口的样式和显示模式
	switch w.config.Mode {
	case Minimized:
		style |= winStyle
		swpStyle |= windows.SWP_NOMOVE | windows.SWP_NOSIZE
		showMode = windows.SW_SHOWMINIMIZED

	case Maximized:
		style |= winStyle
		swpStyle |= windows.SWP_NOMOVE | windows.SWP_NOSIZE
		showMode = windows.SW_SHOWMAXIMIZED

	case Windowed:
		style |= winStyle
		showMode = windows.SW_SHOWNORMAL
		// 获取目标的客户区大小
		width = int32(w.config.Size.X)
		height = int32(w.config.Size.Y)
		// 获取当前窗口的大小和位置
		wr := windows.GetWindowRect(w.hwnd)
		x = wr.Left
		y = wr.Top
		if w.config.Decorated {
			// 计算客户区的大小和位置。注意，当我们控制装饰时，客户区的大小等于窗口的大小
			r := windows.Rect{
				Right:  width,
				Bottom: height,
			}
			windows.AdjustWindowRectEx(&r, uint32(style), 0, dwExStyle)
			width = r.Right - r.Left
			height = r.Bottom - r.Top
		}
		if !w.config.Decorated {
			// 当我们绘制装饰时，启用阴影效果
			windows.DwmExtendFrameIntoClientArea(w.hwnd, windows.Margins{-1, -1, -1, -1})
		}

	case Fullscreen:
		swpStyle |= windows.SWP_NOMOVE | windows.SWP_NOSIZE
		mi := windows.GetMonitorInfo(w.hwnd)
		x, y = mi.Monitor.Left, mi.Monitor.Top
		width = mi.Monitor.Right - mi.Monitor.Left
		height = mi.Monitor.Bottom - mi.Monitor.Top
		showMode = windows.SW_SHOWMAXIMIZED
	}
	// 设置窗口的样式
	windows.SetWindowLong(w.hwnd, windows.GWL_STYLE, style)
	// 设置窗口的位置和大小
	windows.SetWindowPos(w.hwnd, 0, x, y, width, height, swpStyle)
	// 显示窗口
	windows.ShowWindow(w.hwnd, showMode)

	// 更新窗口
	w.update()
}

// WriteClipboard 方法将指定的字符串写入剪贴板
func (w *window) WriteClipboard(s string) {
	w.writeClipboard(s)
}

// writeClipboard 方法将指定的字符串写入剪贴板，如果出现错误则返回错误
func (w *window) writeClipboard(s string) error {
	// 打开剪贴板
	if err := windows.OpenClipboard(w.hwnd); err != nil {
		return err
	}
	// 确保剪贴板在函数结束后关闭
	defer windows.CloseClipboard()
	// 清空剪贴板
	if err := windows.EmptyClipboard(); err != nil {
		return err
	}
	// 将字符串转换为 UTF16 编码
	u16, err := gowindows.UTF16FromString(s)
	if err != nil {
		return err
	}
	// 分配全局内存
	n := len(u16) * int(unsafe.Sizeof(u16[0]))
	mem, err := windows.GlobalAlloc(n)
	if err != nil {
		return err
	}
	// 锁定全局内存
	ptr, err := windows.GlobalLock(mem)
	if err != nil {
		windows.GlobalFree(mem)
		return err
	}
	// 将字符串复制到全局内存
	u16v := unsafe.Slice((*uint16)(ptr), len(u16))
	copy(u16v, u16)
	// 解锁全局内存
	windows.GlobalUnlock(mem)
	// 将全局内存设置为剪贴板的数据
	if err := windows.SetClipboardData(windows.CF_UNICODETEXT, mem); err != nil {
		windows.GlobalFree(mem)
		return err
	}
	// 返回无错误
	return nil
}

// SetCursor 方法设置窗口的光标
func (w *window) SetCursor(cursor pointer.Cursor) {
	// 加载光标
	c, err := loadCursor(cursor)
	if err != nil {
		c = resources.cursor
	}
	// 设置光标
	w.cursor = c
	// 如果光标在窗口内，设置窗口的光标
	if w.cursorIn {
		windows.SetCursor(w.cursor)
	}
}

// windowsCursor 包含从 pointer.Cursor 到 IDC 的映射
var windowsCursor = [...]uint16{
	pointer.CursorDefault:                  windows.IDC_ARROW,       // 默认光标，对应 Windows 的箭头光标
	pointer.CursorNone:                     0,                       // 无光标
	pointer.CursorText:                     windows.IDC_IBEAM,       // 文本光标，对应 Windows 的 I 形光标
	pointer.CursorVerticalText:             windows.IDC_IBEAM,       // 垂直文本光标，对应 Windows 的 I 形光标
	pointer.CursorPointer:                  windows.IDC_HAND,        // 指针光标，对应 Windows 的手形光标
	pointer.CursorCrosshair:                windows.IDC_CROSS,       // 十字光标，对应 Windows 的十字形光标
	pointer.CursorAllScroll:                windows.IDC_SIZEALL,     // 全方向滚动光标，对应 Windows 的四向箭头光标
	pointer.CursorColResize:                windows.IDC_SIZEWE,      // 列调整光标，对应 Windows 的双向箭头（左右）光标
	pointer.CursorRowResize:                windows.IDC_SIZENS,      // 行调整光标，对应 Windows 的双向箭头（上下）光标
	pointer.CursorGrab:                     windows.IDC_SIZEALL,     // 抓取光标，对应 Windows 的四向箭头光标
	pointer.CursorGrabbing:                 windows.IDC_SIZEALL,     // 正在抓取光标，对应 Windows 的四向箭头光标
	pointer.CursorNotAllowed:               windows.IDC_NO,          // 不允许光标，对应 Windows 的禁止符号光标
	pointer.CursorWait:                     windows.IDC_WAIT,        // 等待光标，对应 Windows 的等待符号光标
	pointer.CursorProgress:                 windows.IDC_APPSTARTING, // 进度光标，对应 Windows 的应用启动光标
	pointer.CursorNorthWestResize:          windows.IDC_SIZENWSE,    // 西北调整光标，对应 Windows 的双向箭头（左上右下）光标
	pointer.CursorNorthEastResize:          windows.IDC_SIZENESW,    // 东北调整光标，对应 Windows 的双向箭头（右上左下）光标
	pointer.CursorSouthWestResize:          windows.IDC_SIZENESW,    // 西南调整光标，对应 Windows 的双向箭头（右上左下）光标
	pointer.CursorSouthEastResize:          windows.IDC_SIZENWSE,    // 东南调整光标，对应 Windows 的双向箭头（左上右下）光标
	pointer.CursorNorthSouthResize:         windows.IDC_SIZENS,      // 南北调整光标，对应 Windows 的双向箭头（上下）光标
	pointer.CursorEastWestResize:           windows.IDC_SIZEWE,      // 东西调整光标，对应 Windows 的双向箭头（左右）光标
	pointer.CursorWestResize:               windows.IDC_SIZEWE,      // 西调整光标，对应 Windows 的双向箭头（左右）光标
	pointer.CursorEastResize:               windows.IDC_SIZEWE,      // 东调整光标，对应 Windows 的双向箭头（左右）光标
	pointer.CursorNorthResize:              windows.IDC_SIZENS,      // 北调整光标，对应 Windows 的双向箭头（上下）光标
	pointer.CursorSouthResize:              windows.IDC_SIZENS,      // 南调整光标，对应 Windows 的双向箭头（上下）光标
	pointer.CursorNorthEastSouthWestResize: windows.IDC_SIZENESW,    // 东北西南调整光标，对应 Windows 的双向箭头（右上左下）光标
	pointer.CursorNorthWestSouthEastResize: windows.IDC_SIZENWSE,    // 西北东南调整光标，对应 Windows 的双向箭头（左上右下）光标
}

// loadCursor 函数根据指定的光标类型加载光标，如果出现错误则返回错误
func loadCursor(cursor pointer.Cursor) (syscall.Handle, error) {
	switch cursor {
	case pointer.CursorDefault:
		return resources.cursor, nil // 默认光标，直接返回预设的光标
	case pointer.CursorNone:
		return 0, nil // 无光标，返回0
	default:
		return windows.LoadCursor(windowsCursor[cursor]) // 其他类型的光标，通过 windows.LoadCursor 加载
	}
}

// ShowTextInput 方法用于显示或隐藏文本输入，此处为空实现
func (w *window) ShowTextInput(show bool) {}

// SetInputHint 方法用于设置输入提示，此处为空实现
func (w *window) SetInputHint(_ key.InputHint) {}

// HDC 方法返回窗口的设备上下文句柄
func (w *window) HDC() syscall.Handle {
	return w.hdc
}

// HWND 方法返回窗口的句柄以及窗口的宽度和高度
func (w *window) HWND() (syscall.Handle, int, int) {
	return w.hwnd, w.config.Size.X, w.config.Size.Y
}

// Perform 方法用于执行一系列的系统动作
func (w *window) Perform(acts system.Action) {
	walkActions(acts, func(a system.Action) {
		switch a {
		case system.ActionCenter: // 窗口居中动作
			if w.config.Mode != Windowed {
				break
			}
			r := windows.GetWindowRect(w.hwnd) // 获取窗口的矩形区域
			dx := r.Right - r.Left             // 计算窗口的宽度
			dy := r.Bottom - r.Top             // 计算窗口的高度
			// Calculate center position on current monitor.
			mi := windows.GetMonitorInfo(w.hwnd).Monitor // 获取当前显示器的信息
			x := (mi.Right - mi.Left - dx) / 2           // 计算窗口在显示器上居中的横坐标
			y := (mi.Bottom - mi.Top - dy) / 2           // 计算窗口在显示器上居中的纵坐标
			// 设置窗口的位置和大小，使其居中
			windows.SetWindowPos(w.hwnd, 0, x, y, dx, dy, windows.SWP_NOZORDER|windows.SWP_FRAMECHANGED)
		case system.ActionRaise: // 窗口置顶动作
			w.raise()
		case system.ActionClose: // 关闭窗口动作
			windows.PostMessage(w.hwnd, windows.WM_CLOSE, 0, 0)
		}
	})
}

// raise 方法用于将窗口置顶
func (w *window) raise() {
	windows.SetForegroundWindow(w.hwnd) // 将窗口设置为前台窗口
	// 将窗口置顶，但不改变其位置和大小
	windows.SetWindowPos(w.hwnd, windows.HWND_TOPMOST, 0, 0, 0, 0,
		windows.SWP_NOMOVE|windows.SWP_NOSIZE|windows.SWP_SHOWWINDOW)
}

// convertKeyCode 函数用于将虚拟键码转换为字符串
func convertKeyCode(code uintptr) (string, bool) {
	if '0' <= code && code <= '9' || 'A' <= code && code <= 'Z' {
		return string(rune(code)), true // 如果虚拟键码是数字或字母，则直接转换为字符串
	}
	var r string
	// 根据虚拟键码的值，将其转换为对应的键名
	switch code {
	case windows.VK_ESCAPE:
		r = key.NameEscape // Escape 键
	case windows.VK_LEFT:
		r = key.NameLeftArrow // 左箭头键
	case windows.VK_RIGHT:
		r = key.NameRightArrow // 右箭头键
	case windows.VK_RETURN:
		r = key.NameReturn // 回车键
	case windows.VK_UP:
		r = key.NameUpArrow // 上箭头键
	case windows.VK_DOWN:
		r = key.NameDownArrow // 下箭头键
	case windows.VK_HOME:
		r = key.NameHome // Home 键
	case windows.VK_END:
		r = key.NameEnd // End 键
	case windows.VK_BACK:
		r = key.NameDeleteBackward // Backspace 键
	case windows.VK_DELETE:
		r = key.NameDeleteForward // Delete 键
	case windows.VK_PRIOR:
		r = key.NamePageUp // Page Up 键
	case windows.VK_NEXT:
		r = key.NamePageDown // Page Down 键
	case windows.VK_F1:
		r = key.NameF1 // F1 键
	case windows.VK_F2:
		r = key.NameF2 // F2 键
	case windows.VK_F3:
		r = key.NameF3 // F3 键
	case windows.VK_F4:
		r = key.NameF4 // F4 键
	case windows.VK_F5:
		r = key.NameF5 // F5 键
	case windows.VK_F6:
		r = key.NameF6 // F6 键
	case windows.VK_F7:
		r = key.NameF7 // F7 键
	case windows.VK_F8:
		r = key.NameF8 // F8 键
	case windows.VK_F9:
		r = key.NameF9 // F9 键
	case windows.VK_F10:
		r = key.NameF10 // F10 键
	case windows.VK_F11:
		r = key.NameF11 // F11 键
	case windows.VK_F12:
		r = key.NameF12 // F12 键
	case windows.VK_TAB:
		r = key.NameTab // Tab 键
	case windows.VK_SPACE:
		r = key.NameSpace // 空格键
	case windows.VK_OEM_1:
		r = ";" // 分号键
	case windows.VK_OEM_PLUS:
		r = "+" // 加号键
	case windows.VK_OEM_COMMA:
		r = "," // 逗号键
	case windows.VK_OEM_MINUS:
		r = "-" // 减号键
	case windows.VK_OEM_PERIOD:
		r = "." // 句号键
	case windows.VK_OEM_2:
		r = "/" // 斜杠键
	case windows.VK_OEM_3:
		r = "`" // 反引号键
	case windows.VK_OEM_4:
		r = "[" // 左方括号键
	case windows.VK_OEM_5, windows.VK_OEM_102:
		r = "\\" // 反斜杠键
	case windows.VK_OEM_6:
		r = "]" // 右方括号键
	case windows.VK_OEM_7:
		r = "'" // 单引号键
	case windows.VK_CONTROL:
		r = key.NameCtrl // Ctrl 键
	case windows.VK_SHIFT:
		r = key.NameShift // Shift 键
	case windows.VK_MENU:
		r = key.NameAlt // Alt 键
	case windows.VK_LWIN, windows.VK_RWIN:
		r = key.NameSuper // Win 键
	default:
		return "", false // 如果没有匹配的键名，返回空字符串和 false
	}
	return r, true // 返回键名和 true
}

// configForDPI 函数根据给定的 DPI（每英寸点数）值生成一个 unit.Metric 对象
// 这个对象包含了每个设备独立像素(DP)和比例独立像素(SP)的像素数
func configForDPI(dpi int) unit.Metric {
	// 每个设备独立像素(DP)的英寸数，这里设定为 1/96，这是 Android 系统的标准
	const inchPrDp = 1.0 / 96.0
	// 计算每个设备独立像素(DP)和比例独立像素(SP)的像素数
	ppdp := float32(dpi) * inchPrDp
	// 返回一个 unit.Metric 对象，其中 PxPerDp 和 PxPerSp 都设置为计算得到的像素数
	return unit.Metric{
		PxPerDp: ppdp,
		PxPerSp: ppdp,
	}
}

// ImplementsEvent 方法是一个空实现，用于满足 Event 接口的要求
// 这里的 ViewEvent 是一个空结构体，它实现了 Event 接口，但并没有添加任何额外的方法或属性
func (_ ViewEvent) ImplementsEvent() {}

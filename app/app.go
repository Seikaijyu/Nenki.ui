// 应用程序窗口管理器
// 应用程序窗口管理器不实现MonadInterface
package app

import (
	"image/color"
	"runtime"

	gapp "gioui.org/app"
	glayout "gioui.org/layout"
	"gioui.org/op"
	gunit "gioui.org/unit"
	"nenki.ui/context"
)

// Orientation是应用程序的方向
//
// 仅支持Android和JS
type Orientation uint8

const (
	// 允许窗口自由定向
	AnyOrientation Orientation = iota
	// 将窗口限制为横向
	LandscapeOrientation
	// 将窗口限制为纵向
	PortraitOrientation
)

// 窗口模式
//
// 模式可以通过编程方式或用户点击窗口标题栏上的最小化/最大化按钮来更改
type WindowMode = gapp.WindowMode

const (
	// Windowed是带有特定于操作系统的窗口装饰的正常窗口模式。
	Windowed WindowMode = iota
	// 全屏窗口模式
	Fullscreen
	// 最小化
	Minimized
	// 最大化
	Maximized
)

// 程序
type App struct {
	window    *gapp.Window   // 窗口
	uiContext *context.AppUI // UI上下文管理器
}

// 主动更新UI
func (p *App) update(gtx glayout.Context) *App {
	op.InvalidateOp{}.Add(gtx.Ops)
	return p
}

// 此函数会在每次UI循环时调用，用于更新UI
//
// 此函数执行后会根据返回值判断是否需要完全更新UI，减少UI更新次数以提高性能
func (p *App) Loop(fn func(self *App, root *context.Root)) *App {
	p.uiContext.OnUILoop(func(gtx glayout.Context) {
		fn(p, p.uiContext.GetUIWidget().(*context.Root))
		p.update(gtx)
	})
	return p
}

// 此函数仅在UI循环中执行一次，用于初始化UI或者修改UI
//
// 此函数执行后会根据返回值判断是否需要完全更新UI，减少UI更新次数以提高性能
func (p *App) Then(fn func(self *App, root *context.Root)) *App {
	p.uiContext.AppendSingleUIHandler(func(gtx glayout.Context) {
		fn(p, p.uiContext.GetUIWidget().(*context.Root))
		p.update(gtx)
	})
	return p
}

// 设置标题
func (p *App) Title(title string) *App {
	p.window.Option(gapp.Title(title))
	return p
}

// 设置窗口尺寸
func (p *App) Size(width, height float32) *App {
	p.window.Option(gapp.Size(gunit.Dp(width), gunit.Dp(height)))
	return p
}

// 设置窗口最小尺寸
func (p *App) MinSize(width, height float32) *App {
	p.window.Option(gapp.MinSize(gunit.Dp(width), gunit.Dp(height)))
	return p
}

// 设置窗口最大尺寸
func (p *App) MaxSize(width, height float32) *App {
	p.window.Option(gapp.MaxSize(gunit.Dp(width), gunit.Dp(height)))
	return p
}

// 设置Android导航栏或者浏览器地址栏的颜色
//
// 仅支持Android和JS
func (p *App) NavigationColor(r, g, b, a uint8) *App {
	p.window.Option(gapp.NavigationColor(color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}))
	return p
}

// 控制窗口是否自绘装饰边框，false表示应用程序将不绘制自己的装饰边框
func (p *App) Decorated(visible bool) *App {
	p.window.Option(gapp.Decorated(visible))
	return p
}

// 用于设置 Android 状态栏的颜色
func (p *App) StatusColor(r, g, b, a uint8) *App {
	p.window.Option(gapp.StatusColor(color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}))
	return p
}

// 设置文件拖拽到窗口
func (p *App) DragFiles(enabled bool) *App {
	if runtime.GOOS == "windows" {
		gapp.DragAcceptFiles(gapp.HWND, enabled)

	}
	return p
}

// 设置窗口布局方向
//
// 仅支持Android和JS
func (p *App) Orientation(orientation Orientation) *App {
	p.window.Option(gapp.Orientation(orientation).Option())
	return p
}

// 设置背景颜色
func (p *App) Background(r, g, b, a uint8) *App {
	p.uiContext.Background(r, g, b, a)
	return p
}

// 设置窗口模式
func (p *App) WindowMode(mode WindowMode) *App {
	p.window.Option(mode.Option())
	return p
}

// 设置文件拖拽到窗口处理函数
//
// 仅支持Windows
func (p *App) OnDropFiles(fn func(files []string)) *App {
	gapp.CustomDragHandler(func(files []string) {
		p.Then(func(self *App, root *context.Root) {
			fn(files)
		})
	})
	return p
}

// 自定义UI上下文错误处理函数
func (p *App) OnUIContextError(fn func(err error)) *App {
	p.Then(func(self *App, root *context.Root) {
		p.uiContext.OnUIContextError(fn)
	})
	return p
}

// 创建应用，需要传入窗口名字
func NewApp(title string) *App {
	window := gapp.NewWindow()

	var application = &App{
		window:    window,
		uiContext: context.NewAppUI(window),
	}
	application.Title(title) // 设置标题
	// 提前调用一次以获取HWND和更快的加载
	window.NextEvent()
	return application
}

// 阻塞以进行UI循环
func Run() {
	select {}
}

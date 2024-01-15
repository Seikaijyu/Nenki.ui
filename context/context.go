// UI上下文管理器，是拼接UI和App的中间件
// 因为是中层（用户不可见），所以不实现MonadInterface
package context

import (
	"fmt"
	"image/color"
	"os"

	"github.com/Seikaijyu/nenki.ui/widget"

	gio "github.com/Seikaijyu/gio/app"
	"github.com/Seikaijyu/gio/io/system"
	glayout "github.com/Seikaijyu/gio/layout"
	"github.com/Seikaijyu/gio/op"
	"github.com/Seikaijyu/gio/op/clip"
	"github.com/Seikaijyu/gio/op/paint"
)

// 根组件
type Root = widget.ContainerLayout

// UI上下文配置
type contextConfig struct {
	// 背景颜色
	background *color.NRGBA
}

// UI上下文管理器
type AppUI struct {
	fatalHandler        func(err error)               // 错误处理
	window              *gio.Window                   // 基础窗口
	uiWidget            widget.WidgetInterface        // UI组件
	data                chan map[string]any           // 数据
	config              *contextConfig                // 配置
	updateHandler       func(glayout.Context)         // UI每次更新时执行的函数
	singleUpdateHandler *Queue[func(glayout.Context)] // 单次执行的UI函数
	graphContext        glayout.Context               // 渲染上下文
}

// UI循环
func (p *AppUI) loop() error {
	var ops op.Ops
	for {

		select {
		case data := <-p.data:
			fmt.Println(data)
			p.window.Invalidate()
		default:
			switch e := p.window.NextEvent().(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:

				p.graphContext = glayout.NewContext(&ops, e)
				var stack = clip.Rect{Max: e.Size}.Push(p.graphContext.Ops)
				// 设置背景颜色
				if p.config.background != nil {
					paint.ColorOp{Color: *p.config.background}.Add(p.graphContext.Ops)
					paint.PaintOp{}.Add(p.graphContext.Ops)
				}
				p.updateHandler(p.graphContext)
				// 执行队列中的UI函数，如果有的话
				if fn, ok := p.singleUpdateHandler.Dequeue(); ok {
					fn(p.graphContext)
					p.window.Invalidate()
				}
				// 渲染UI
				p.uiWidget.Layout(p.graphContext)
				stack.Pop()
				e.Frame(p.graphContext.Ops)

			}
		}
	}
}

// 设置背景颜色
func (p *AppUI) Background(r, g, b, a uint8) *AppUI {
	p.config.background = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 自定义UI循环
func (p *AppUI) OnUILoop(fn func(glayout.Context)) {
	p.updateHandler = fn
}

// 自定义错误处理函数
func (p *AppUI) OnUIContextError(fn func(err error)) *AppUI {
	p.fatalHandler = fn
	return p
}

func (p *AppUI) AppendSingleUIHandler(fn func(glayout.Context)) {
	p.singleUpdateHandler.Enqueue(fn)
	p.window.Invalidate()
}

// 获取UI组件
func (p *AppUI) GetUIWidget() widget.WidgetInterface {
	return p.uiWidget
}

// 获取渲染上下文
func (p *AppUI) GetGraphContext() glayout.Context {
	return p.graphContext
}

// 创建一个UI上下文管理器
func NewAppUI(window *gio.Window) *AppUI {
	uiContext := &AppUI{
		fatalHandler: func(err error) {
			panic(err)
		},
		window:              window,
		uiWidget:            widget.NewContainerLayout(),
		updateHandler:       func(glayout.Context) {},
		singleUpdateHandler: &Queue[func(glayout.Context)]{},
		config:              &contextConfig{},
	}
	uiContext.uiWidget.OnDestroy(func() {})
	go func() {
		// 进行UI循环
		if err := uiContext.loop(); err != nil {
			// 进行错误处理
			uiContext.fatalHandler(err)
		}
		os.Exit(0) // 退出程序
	}()
	return uiContext
}
// UI上下文管理器，是拼接UI和App的中间件
// 因为是中层（用户不可见），所以不实现MonadInterface
package context

import (
	"fmt"
	"image/color"
	"os"

	"github.com/Seikaijyu/nenki.ui/widget"

	gio "gioui.org/app"
	"gioui.org/io/system"
	glayout "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// 根组件
type Root = widget.ContainerLayout

// UI上下文配置
type contextConfig struct {
	// 背景颜色
	background *color.NRGBA
}

// UI上下文管理器
type AppUI struct {
	fatalHandler        func(err error)               // 错误处理
	window              *gio.Window                   // 基础窗口
	uiWidget            widget.WidgetInterface        // UI组件
	data                chan map[string]any           // 数据
	config              *contextConfig                // 配置
	updateHandler       func(glayout.Context)         // UI每次更新时执行的函数
	singleUpdateHandler *Queue[func(glayout.Context)] // 单次执行的UI函数
	graphContext        glayout.Context               // 渲染上下文
}

// UI循环
func (p *AppUI) loop() error {
	var ops op.Ops
	for {

		select {
		case data := <-p.data:
			fmt.Println(data)
			p.window.Invalidate()
		default:
			switch e := p.window.NextEvent().(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:

				p.graphContext = glayout.NewContext(&ops, e)
				var stack = clip.Rect{Max: e.Size}.Push(p.graphContext.Ops)
				// 设置背景颜色
				if p.config.background != nil {
					paint.ColorOp{Color: *p.config.background}.Add(p.graphContext.Ops)
					paint.PaintOp{}.Add(p.graphContext.Ops)
				}
				p.updateHandler(p.graphContext)
				// 执行队列中的UI函数，如果有的话
				if fn, ok := p.singleUpdateHandler.Dequeue(); ok {
					fn(p.graphContext)
					p.window.Invalidate()
				}
				// 渲染UI
				p.uiWidget.Layout(p.graphContext)
				stack.Pop()
				e.Frame(p.graphContext.Ops)

			}
		}
	}
}

// 设置背景颜色
func (p *AppUI) Background(r, g, b, a uint8) *AppUI {
	p.config.background = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 自定义UI循环
func (p *AppUI) OnUILoop(fn func(glayout.Context)) {
	p.updateHandler = fn
}

// 自定义错误处理函数
func (p *AppUI) OnUIContextError(fn func(err error)) *AppUI {
	p.fatalHandler = fn
	return p
}

func (p *AppUI) AppendSingleUIHandler(fn func(glayout.Context)) {
	p.singleUpdateHandler.Enqueue(fn)
	p.window.Invalidate()
}

// 获取UI组件
func (p *AppUI) GetUIWidget() widget.WidgetInterface {
	return p.uiWidget
}

// 获取渲染上下文
func (p *AppUI) GetGraphContext() glayout.Context {
	return p.graphContext
}

// 创建一个UI上下文管理器
func NewAppUI(window *gio.Window) *AppUI {
	uiContext := &AppUI{
		fatalHandler: func(err error) {
			panic(err)
		},
		window:              window,
		uiWidget:            widget.NewContainerLayout(),
		updateHandler:       func(glayout.Context) {},
		singleUpdateHandler: &Queue[func(glayout.Context)]{},
		config:              &contextConfig{},
	}
	uiContext.uiWidget.OnDestroy(func() {})
	go func() {
		// 进行UI循环
		if err := uiContext.loop(); err != nil {
			// 进行错误处理
			uiContext.fatalHandler(err)
		}
		os.Exit(0) // 退出程序
	}()
	return uiContext
}

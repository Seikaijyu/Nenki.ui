// UI上下文管理器，是拼接UI和App的中间件
// 因为是中层（用户不可见），所以不实现MonadInterface
package context

import (
	"fmt"
	"image/color"
	"os"

	gio "gioui.org/app"
	"gioui.org/io/system"
	glayout "gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"nenki.ui/widget"
)

type ContextConfig struct {
	// 背景颜色
	Background *color.NRGBA
}

// UI上下文管理器
type AppUI struct {
	fatal           func(err error)        // 错误处理
	window          *gio.Window            // 基础窗口
	uiWidget        widget.WidgetInterface // UI组件
	data            chan map[string]any    // 数据
	config          *ContextConfig         // 配置
	uiHandler       func(glayout.Context)
	singleUIHandler *Queue[func(glayout.Context)] // 单次执行的UI函数
	gtx             glayout.Context
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
				p.gtx = glayout.NewContext(&ops, e)
				var stack = clip.Rect{Max: e.Size}.Push(p.gtx.Ops)
				// 设置背景颜色
				if p.config.Background != nil {
					paint.ColorOp{Color: color.NRGBA{R: p.config.Background.R, G: p.config.Background.G, B: p.config.Background.B, A: p.config.Background.A}}.Add(p.gtx.Ops)
					paint.PaintOp{}.Add(p.gtx.Ops)
				}
				p.uiHandler(p.gtx)
				// 执行队列中的UI函数，如果有的话
				if fn, ok := p.singleUIHandler.Dequeue(); ok {
					fn(p.gtx)
					p.window.Invalidate()
				}
				// 渲染UI
				p.uiWidget.Layout(p.gtx)
				stack.Pop()
				e.Frame(p.gtx.Ops)

			}
		}
	}
}

// 设置背景颜色
func (p *AppUI) SetBackground(r, g, b, a uint8) *AppUI {
	p.config.Background = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 获取背景颜色
func (p *AppUI) Background() (uint8, uint8, uint8, uint8) {
	return p.config.Background.R, p.config.Background.G, p.config.Background.B, p.config.Background.A
}

// 自定义UI循环
func (p *AppUI) CustomUIHandler(fn func(glayout.Context)) {
	p.uiHandler = fn
}

// 自定义错误处理函数
func (p *AppUI) CustomFatalHandler(fn func(err error)) *AppUI {
	p.fatal = fn
	return p
}

func (p *AppUI) AppendSingleUIHandler(fn func(glayout.Context)) {
	p.singleUIHandler.Enqueue(fn)
	p.window.Invalidate()
}

// 获取UI组件
func (p *AppUI) GetUIWidget() widget.WidgetInterface {
	return p.uiWidget
}

// 获取渲染上下文
func (p *AppUI) GetGraphContext() glayout.Context {
	return p.gtx
}

// 创建一个UI上下文管理器
func NewAppUI(window *gio.Window) *AppUI {
	uiContext := &AppUI{
		fatal: func(err error) {
			panic(err)
		},
		window:          window,
		uiWidget:        widget.NewAnchorLayout(widget.Center).SetDirection(widget.TopLeft),
		uiHandler:       func(glayout.Context) {},
		singleUIHandler: &Queue[func(glayout.Context)]{},
		config:          &ContextConfig{},
	}
	go func() {
		// 进行UI循环
		if err := uiContext.loop(); err != nil {
			// 进行错误处理
			uiContext.fatal(err)
		}
		os.Exit(0) // 退出程序
	}()
	return uiContext
}

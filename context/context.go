// UI上下文管理器，是拼接UI和App的中间件
// 因为是中层（用户不可见），所以不实现MonadInterface
package context

import (
	"fmt"
	"os"

	gio "gioui.org/app"
	"gioui.org/io/system"
	glayout "gioui.org/layout"
	"gioui.org/op"
	"nenki.ui/widget"
)

// UI上下文管理器
type AppUI struct {
	fatal     func(err error)        // 错误处理
	window    *gio.Window            // 基础窗口
	uiWidget  widget.WidgetInterface // UI组件
	data      chan map[string]any    // 数据
	uiHandler func(glayout.Context, *widget.AnchorLayout)
}

// UI循环
func (p *AppUI) loop() error {
	var ops op.Ops
	for {
		select {
		case data := <-p.data:
			fmt.Println(data)
			// 更新数据
			p.window.Invalidate()
		default:
			switch e := p.window.NextEvent().(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := glayout.NewContext(&ops, e)
				p.uiHandler(gtx, p.uiWidget.(*widget.AnchorLayout))
				// 渲染UI
				p.uiWidget.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

// 自定义UI组件
func (p *AppUI) CustomUIHandler(fn func(glayout.Context, *widget.AnchorLayout)) {
	p.uiHandler = fn
}

// 自定义错误处理函数
func (p *AppUI) CustomFatalHandler(fn func(err error)) *AppUI {
	p.fatal = fn
	return p
}

// 创建一个UI上下文管理器
func NewAppUI(window *gio.Window) *AppUI {
	uiContext := &AppUI{
		fatal: func(err error) {
			panic(err)
		},
		window:    window,
		uiWidget:  widget.NewAnchorLayout(widget.Center).SetDirection(widget.TopLeft),
		uiHandler: func(glayout.Context, *widget.AnchorLayout) {},
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

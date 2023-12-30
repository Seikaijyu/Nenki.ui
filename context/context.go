// UI上下文管理器，是拼接UI和App的中间件
package context

import (
	"os"

	gio "gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget/material"
)

// UI上下文管理器
type AppUI struct {
	globalTheme *material.Theme // 全局主题
	fatal       func(err error) // 错误处理
	window      *gio.Window     // 基础窗口
}

// UI循环
func (p *AppUI) loop() error {
	var ops op.Ops
	for {
		switch e := p.window.NextEvent().(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)

			e.Frame(gtx.Ops)
		}
	}
}

// 自定义错误处理函数
func (p *AppUI) CustomFatalHandler(fn func(err error)) *AppUI {
	p.fatal = fn
	return p
}

// 创建一个UI上下文管理器
func NewAppUI(window *gio.Window) *AppUI {
	uiContext := &AppUI{
		globalTheme: material.NewTheme(),
		fatal: func(err error) {
			panic(err)
		},
		window: window,
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

package app

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"nenki.ui/monad"
)

var _ monad.MonadInterface[*App, string] = &App{}

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
type WindowMode uint8

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

// 应用程序配置
type AppConfig struct {
	// 窗口尺寸
	width  float32
	height float32
	// 最小尺寸
	minWidth  float32
	minHeight float32
	// 最大尺寸
	maxWidth  float32
	maxHeight float32
	// 窗口名字
	title string
	// Android导航栏或者浏览器地址栏的颜色
	navigationColor color.NRGBA
	// 是否显示窗口边框
	decoratedVisible bool
	// Android状态栏颜色
	statusColor color.NRGBA
	// 窗口布局方向
	orientation Orientation
	// 窗口模式
	windowMode WindowMode
}

// 未经过包装的原始数据
type AppBaseData struct {
	window      *app.Window     // 窗口
	globalTheme *material.Theme // 全局主题
	fatal       func(err error) // 错误处理
}

// UI循环
func (p *AppBaseData) loop(app *App) error {
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

// 程序
type App struct {
	config *AppConfig   // 配置项
	base   *AppBaseData // 原始数据
}

// 绑定函数
func (p *App) Bind(fn func(*App) *App) *App {
	return fn(p)
}

// 包装一个值，相当于调用NewApp
func (p *App) Unit(title string) *App {
	return NewApp(title)
}

// 设置标题
func (p *App) SetTitle(title string) *App {
	p.config.title = title
	p.base.window.Option(app.Title(title))
	return p
}

// 获取标题
func (p *App) Title() string {
	return p.config.title
}

// 设置窗口尺寸
func (p *App) SetSize(width, height float32) *App {
	p.config.width = width
	p.config.height = height
	p.base.window.Option(app.Size(unit.Dp(width), unit.Dp(height)))
	return p
}

// 获取窗口尺寸
func (p *App) Size(width, height float32) (float32, float32) {
	return p.config.width, p.config.height
}

// 设置窗口最小尺寸
func (p *App) SetMinSize(width, height float32) *App {
	p.config.minWidth = width
	p.config.minHeight = height
	p.base.window.Option(app.MinSize(unit.Dp(width), unit.Dp(height)))
	return p
}

// 获取窗口最小尺寸
func (p *App) MinSize(width, height float32) (float32, float32) {
	return p.config.minWidth, p.config.minHeight
}

// 设置窗口最大尺寸
func (p *App) SetMaxSize(width, height float32) *App {
	p.config.maxWidth = width
	p.config.maxHeight = height
	p.base.window.Option(app.MaxSize(unit.Dp(width), unit.Dp(height)))
	return p
}

// 获取窗口最大尺寸
func (p *App) MaxSize(width, height float32) (float32, float32) {
	return p.config.maxWidth, p.config.maxHeight
}

// 设置Android导航栏或者浏览器地址栏的颜色
//
// 仅支持Android和JS
func (p *App) SetNavigationColor(r, g, b, a uint8) *App {
	p.config.navigationColor = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	p.base.window.Option(app.NavigationColor(p.config.navigationColor))
	return p
}

// 获取Android导航栏或者浏览器地址栏的颜色
//
// 仅支持Android和JS
func (p *App) NavigationColor() (uint8, uint8, uint8, uint8) {
	return p.config.navigationColor.R, p.config.navigationColor.G, p.config.navigationColor.B, p.config.navigationColor.A
}

// 控制窗口是否自绘装饰边框，false表示应用程序将不绘制自己的装饰边框
func (p *App) SetDecorated(visible bool) *App {
	p.config.decoratedVisible = visible
	p.base.window.Option(app.Decorated(visible))
	return p
}

// 获取控制窗口是否自绘装饰边框
func (p *App) Decorated() bool {
	return p.config.decoratedVisible
}

// 用于设置 Android 状态栏的颜色
func (p *App) SetStatusColor(r, g, b, a uint8) *App {
	p.config.statusColor = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}

	p.base.window.Option(app.StatusColor(p.config.statusColor))
	return p
}

// 设置窗口布局方向
//
// 仅支持Android和JS
func (p *App) SetOrientation(orientation Orientation) *App {
	p.config.orientation = orientation
	p.base.window.Option(app.Orientation(orientation).Option(), app.Fullscreen.Option())
	return p
}

// 获取窗口布局方向
//
// 仅支持Android和JS
func (p *App) Orientation() Orientation {
	return p.config.orientation
}

// 获取 Android 状态栏的颜色
func (p *App) StatusColor() (uint8, uint8, uint8, uint8) {
	return p.config.statusColor.R, p.config.statusColor.G, p.config.statusColor.B, p.config.statusColor.A
}

// 设置窗口模式
func (p *App) SetWindowMode(mode WindowMode) *App {
	p.config.windowMode = mode
	p.base.window.Option(app.WindowMode(mode).Option())
	return p
}

// 获取窗口模式
func (p *App) WindowMode() WindowMode {
	return p.config.windowMode
}

// 自定义错误处理函数
func (p *App) CustomFatalHandler(fn func(err error)) *App {
	p.base.fatal = fn
	return p
}

// 创建应用，需要传入窗口名字
func NewApp(title string) *App {
	window := app.NewWindow()
	var application = &App{
		config: &AppConfig{},
		base: &AppBaseData{
			window:      window,
			globalTheme: material.NewTheme(),
			fatal: func(err error) {
				panic(err)
			},
		},
	}
	application.SetTitle(title) // 设置标题
	go func() {
		// 进行UI循环
		if err := application.base.loop(application); err != nil {
			// 进行错误处理
			application.base.fatal(err)
		}
		os.Exit(0) // 退出程序
	}()
	return application
}

// 阻塞以进行UI循环
func Run() {
	select {}
}

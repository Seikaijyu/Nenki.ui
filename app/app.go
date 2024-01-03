// 应用程序窗口管理器
// 应用程序窗口管理器不实现MonadInterface
package app

import (
	"image/color"

	"gioui.org/app"
	glayout "gioui.org/layout"
	"gioui.org/unit"
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
	Width  float32
	Height float32
	// 最小尺寸
	MinWidth  float32
	MinHeight float32
	// 最大尺寸
	MaxWidth  float32
	MaxHeight float32
	// 窗口名字
	Title string
	// Android导航栏或者浏览器地址栏的颜色
	NavigationColor *color.NRGBA
	// 是否显示窗口边框
	DecoratedVisible bool
	// Android状态栏颜色
	StatusColor *color.NRGBA
	// 窗口布局方向
	Orientation Orientation
	// 窗口模式
	WindowMode WindowMode
}

// 程序
type App struct {
	window    *app.Window    // 窗口
	config    *AppConfig     // 配置项
	uiContext *context.AppUI // UI上下文管理器
}

func (p *App) update(gtx glayout.Context) *App {
	p.uiContext.GetUIWidget().(*context.Root).Layout(gtx)
	return p
}

// 此函数会在每次UI循环时调用，用于更新UI
//
// 此函数执行后会根据返回值判断是否需要完全更新UI，减少UI更新次数以提高性能
func (p *App) Loop(fn func(self *App, root *context.Root)) *App {
	p.uiContext.CustomUIHandler(func(gtx glayout.Context) {
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
func (p *App) SetTitle(title string) *App {
	p.config.Title = title
	p.window.Option(app.Title(title))
	return p
}

// 获取标题
func (p *App) Title() string {
	return p.config.Title
}

// 设置窗口尺寸
func (p *App) SetSize(width, height float32) *App {
	p.config.Width = width
	p.config.Height = height
	p.window.Option(app.Size(unit.Dp(width), unit.Dp(height)))
	return p
}

// 获取窗口尺寸
func (p *App) Size(width, height float32) (float32, float32) {
	return p.config.Width, p.config.Height
}

// 设置窗口最小尺寸
func (p *App) SetMinSize(width, height float32) *App {
	p.config.MinWidth = width
	p.config.MinHeight = height
	p.window.Option(app.MinSize(unit.Dp(width), unit.Dp(height)))
	return p
}

// 获取窗口最小尺寸
func (p *App) MinSize(width, height float32) (float32, float32) {
	return p.config.MinWidth, p.config.MinHeight
}

// 设置窗口最大尺寸
func (p *App) SetMaxSize(width, height float32) *App {
	p.config.MaxWidth = width
	p.config.MaxHeight = height
	p.window.Option(app.MaxSize(unit.Dp(width), unit.Dp(height)))
	return p
}

// 获取窗口最大尺寸
func (p *App) MaxSize(width, height float32) (float32, float32) {
	return p.config.MaxWidth, p.config.MaxHeight
}

// 设置Android导航栏或者浏览器地址栏的颜色
//
// 仅支持Android和JS
func (p *App) SetNavigationColor(r, g, b, a uint8) *App {
	p.config.NavigationColor = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	p.window.Option(app.NavigationColor(*p.config.NavigationColor))
	return p
}

// 获取Android导航栏或者浏览器地址栏的颜色
//
// 仅支持Android和JS
func (p *App) NavigationColor() (uint8, uint8, uint8, uint8) {
	return p.config.NavigationColor.R, p.config.NavigationColor.G, p.config.NavigationColor.B, p.config.NavigationColor.A
}

// 控制窗口是否自绘装饰边框，false表示应用程序将不绘制自己的装饰边框
func (p *App) SetDecorated(visible bool) *App {
	p.config.DecoratedVisible = visible
	p.window.Option(app.Decorated(visible))
	return p
}

// 获取控制窗口是否自绘装饰边框
func (p *App) Decorated() bool {
	return p.config.DecoratedVisible
}

// 用于设置 Android 状态栏的颜色
func (p *App) SetStatusColor(r, g, b, a uint8) *App {
	p.config.StatusColor = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}

	p.window.Option(app.StatusColor(*p.config.StatusColor))
	return p
}

// 设置窗口布局方向
//
// 仅支持Android和JS
func (p *App) SetOrientation(orientation Orientation) *App {
	p.config.Orientation = orientation
	p.window.Option(app.Orientation(orientation).Option())
	return p
}

// 获取窗口布局方向
//
// 仅支持Android和JS
func (p *App) Orientation() Orientation {
	return p.config.Orientation
}

// 设置背景颜色
func (p *App) SetBackground(r, g, b, a uint8) *App {
	p.uiContext.SetBackground(r, g, b, a)
	return p
}

// 获取背景颜色
func (p *App) Background() (uint8, uint8, uint8, uint8) {
	return p.uiContext.Background()
}

// 获取 Android 状态栏的颜色
func (p *App) StatusColor() (uint8, uint8, uint8, uint8) {
	return p.config.StatusColor.R, p.config.StatusColor.G, p.config.StatusColor.B, p.config.StatusColor.A
}

// 设置窗口模式
func (p *App) SetWindowMode(mode WindowMode) *App {
	p.config.WindowMode = mode
	p.window.Option(app.WindowMode(mode).Option())
	return p
}

// 获取窗口模式
func (p *App) WindowMode() WindowMode {
	return p.config.WindowMode
}

// 自定义错误处理函数
func (p *App) CustomFatalHandler(fn func(err error)) *App {
	p.uiContext.CustomFatalHandler(fn)
	return p
}

// 创建应用，需要传入窗口名字
func NewApp(title string) *App {
	window := app.NewWindow()
	var application = &App{
		window:    window,
		config:    &AppConfig{},
		uiContext: context.NewAppUI(window),
	}
	application.SetTitle(title) // 设置标题

	return application
}

// 阻塞以进行UI循环
func Run() {
	select {}
}

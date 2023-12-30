// 应用程序窗口管理器
// 应用程序窗口管理器不实现MonadInterface
package app

import (
	"image/color"

	"gioui.org/app"
	glayout "gioui.org/layout"
	"gioui.org/unit"
	"nenki.ui/context"
	"nenki.ui/widget"
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

// 程序
type App struct {
	window    *app.Window    // 窗口
	config    *AppConfig     // 配置项
	uiContext *context.AppUI // UI上下文管理器
}

// 绑定函数
func (p *App) Then(fn func(*App, glayout.Context, *widget.AnchorLayout)) {
	p.uiContext.CustomUIHandler(func(gtx glayout.Context, al *widget.AnchorLayout) {
		fn(p, gtx, al)
	})
}

// 设置标题
func (p *App) SetTitle(title string) *App {
	p.config.title = title
	p.window.Option(app.Title(title))
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
	p.window.Option(app.Size(unit.Dp(width), unit.Dp(height)))
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
	p.window.Option(app.MinSize(unit.Dp(width), unit.Dp(height)))
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
	p.window.Option(app.MaxSize(unit.Dp(width), unit.Dp(height)))
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
	p.window.Option(app.NavigationColor(p.config.navigationColor))
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
	p.window.Option(app.Decorated(visible))
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

	p.window.Option(app.StatusColor(p.config.statusColor))
	return p
}

// 设置窗口布局方向
//
// 仅支持Android和JS
func (p *App) SetOrientation(orientation Orientation) *App {
	p.config.orientation = orientation
	p.window.Option(app.Orientation(orientation).Option(), app.Fullscreen.Option())
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
	p.window.Option(app.WindowMode(mode).Option())
	return p
}

// 获取窗口模式
func (p *App) WindowMode() WindowMode {
	return p.config.windowMode
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

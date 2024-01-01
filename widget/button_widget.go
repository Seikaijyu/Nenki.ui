package widget

import (
	"image"
	"image/color"

	glayout "gioui.org/layout"
	"gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
)

// 按钮配置
type ButtonConfig struct {
	// 文本
	Text string
	// 字体大小
	FontSize int
	// 字体颜色
	FontColor *color.NRGBA
	// 圆角
	CornerRadius float32
	// 背景颜色
	Background *color.NRGBA
	//
	Inset image.Point
	// 点击事件
	OnClick func(*Button)
}
type Button struct {
	// 配置
	config *ButtonConfig
	// 主题
	button gmaterial.ButtonStyle
	// 组件是否被删除
	isRemove bool
}

// 校验接口是否实现
var _ WidgetInterface = &Button{}

// 绑定函数
func (p *Button) Then(fn func(*Button)) *Button {
	fn(p)
	return p
}

// 按钮事件
func (p *Button) OnClick(fn func(*Button)) *Button {
	p.config.OnClick = fn
	return p
}

// 是否被删除
func (p *Button) IsDestroy() bool {
	return p.isRemove
}

// 注销自身，清理所有引用
func (p *Button) Destroy() {
	p.isRemove = true
}

// 设置字体大小
func (p *Button) SetFontSize(size int) *Button {
	p.config.FontSize = size
	p.button.TextSize = unit.Sp(size)
	return p
}

// 获取字体大小
func (p *Button) FontSize() int {
	return p.config.FontSize
}

// 设置圆角
func (p *Button) SetCornerRadius(radius float32) *Button {
	p.config.CornerRadius = radius
	p.button.CornerRadius = unit.Dp(radius)
	return p
}

// 获取圆角
func (p *Button) CornerRadius() float32 {
	return p.config.CornerRadius
}

// 设置背景颜色
func (p *Button) SetBackground(r, g, b, a uint8) *Button {
	p.config.Background = &color.NRGBA{R: r, G: g, B: b, A: a}
	p.button.Background = *p.config.Background
	return p
}

// 获取背景颜色
func (p *Button) Background() (uint8, uint8, uint8, uint8) {
	return p.config.Background.R, p.config.Background.G, p.config.Background.B, p.config.Background.A
}

// 设置文本
func (p *Button) SetText(text string) *Button {
	p.config.Text = text
	p.button.Text = text
	return p
}

// 获取文本
func (p *Button) Text() string {
	return p.config.Text
}

// 布局
func (p *Button) Layout(gtx glayout.Context) glayout.Dimensions {
	// 点击事件
	if p.config.OnClick != nil && p.button.Button.Clicked(gtx) {
		p.config.OnClick(p)
	}

	return p.button.Layout(gtx)
}

// 创建按钮
func NewButton(text string) *Button {
	widget := &Button{
		config: &ButtonConfig{
			Text: text,
		},
		button: gmaterial.Button(gmaterial.NewTheme(), &gwidget.Clickable{}, text),
	}
	return widget

}

// 从ID创建按钮
func NewButtonWithID(id string, text string) *Button {
	widget := &Button{
		config: &ButtonConfig{
			Text: text,
		},
		button: gmaterial.Button(gmaterial.NewTheme(), &gwidget.Clickable{}, text),
	}

	return widget
}

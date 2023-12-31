package widget

import (
	glayout "gioui.org/layout"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
)

type ButtonConfig struct {
	// 宽度
	Width int
	// 高度
	Height int
	// 文本
	Text string
	// 字体大小
	FontSize int
	// 字体颜色
	FontColor string
	// 点击事件
	OnClick func(*Button)
}
type Button struct {
	// ID
	id string
	// 索引
	index int
	// 配置
	config *ButtonConfig
	// 主题
	theme *gmaterial.Theme
	// 按钮事件
	clickable *gwidget.Clickable
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
	pool.RemoveAtIndex(p.id, p.index)
}

// 布局
func (p *Button) Layout(gtx glayout.Context) glayout.Dimensions {
	// 点击事件
	if p.config.OnClick != nil && p.clickable.Clicked(gtx) {
		p.config.OnClick(p)
	}

	return gmaterial.Button(p.theme, p.clickable, p.config.Text).Layout(gtx)
}

// 创建按钮
func NewButton(text string) *Button {
	widget := &Button{
		id: "",
		config: &ButtonConfig{
			Text: text,
		},
		theme:     gmaterial.NewTheme(),
		clickable: &gwidget.Clickable{},
	}
	widget.index = pool.AddWidget("", widget)
	return widget
}

// 从ID创建按钮
func NewButtonWithID(id string, text string) *Button {
	widget := &Button{
		id: "#" + id,
		config: &ButtonConfig{
			Text: text,
		},
		theme:     gmaterial.NewTheme(),
		clickable: &gwidget.Clickable{},
	}
	widget.index = pool.AddWidget(id, widget)
	return widget
}

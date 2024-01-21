package widget

import (
	"image/color"

	glayout "github.com/Seikaijyu/gio/layout"
	gunit "github.com/Seikaijyu/gio/unit"
	gmaterial "github.com/Seikaijyu/gio/widget/material"
	"github.com/Seikaijyu/nenki.ui/widget/text"
	"github.com/Seikaijyu/nenki.ui/widget/theme"
)

// 校验接口是否实现
var _ WidgetInterface = &Label{}

type labelConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

type Label struct {
	// 配置
	config *labelConfig
	// 外边距
	margin *glayout.Inset
	// 组件
	labelWidget *gmaterial.LabelStyle
}

// 绑定函数
func (p *Label) Then(fn func(self *Label)) *Label {
	fn(p)
	return p
}

// 注销自身，清理所有引用
func (p *Label) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
	}
	p.config._destroy = nil
}

// 注册删除事件
func (p *Label) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *Label) Update(update bool) {
	p.config.update = update
}

// 外边距
func (p *Label) Margin(Top, Left, Bottom, Right float32) *Label {
	p.margin.Top = gunit.Dp(Top)
	p.margin.Left = gunit.Dp(Left)
	p.margin.Bottom = gunit.Dp(Bottom)
	p.margin.Right = gunit.Dp(Right)
	return p
}

// 布局
func (p *Label) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		p.config.update = false
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.labelWidget.Layout(gtx)
	})
}

// 设置文本内容
func (p *Label) Text(text string) *Label {
	p.labelWidget.Text = text
	return p
}

// 设置文本大小
func (p *Label) FontSize(size float32) *Label {
	p.labelWidget.TextSize = gunit.Sp(size)
	return p
}

// 设置文字颜色
func (p *Label) FontColor(r, g, b, a uint8) *Label {
	p.labelWidget.Color = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置文本对齐方式
func (p *Label) Alignment(alignment text.Alignment) *Label {
	p.labelWidget.Alignment = alignment
	return p
}

// 设置文本行高
func (p *Label) LineHeight(lineHeight float32) *Label {
	p.labelWidget.LineHeight = gunit.Sp(lineHeight)
	return p
}

// 设置文本行高缩放
func (p *Label) LineHeightScale(lineHeightScale float32) *Label {
	p.labelWidget.LineHeightScale = lineHeightScale
	return p
}

// 设置文本最大行
func (p *Label) MaxLines(maxLines int) *Label {
	p.labelWidget.State.MaxLines = maxLines
	return p
}

// 设置文本超出最大行显示的文本
func (p *Label) Truncator(truncator string) *Label {
	p.labelWidget.Truncator = truncator
	return p
}

// 设置如何显示文本换行
func (p *Label) WrapPolicy(wrapPolicy text.WrapPolicy) *Label {
	p.labelWidget.WrapPolicy = wrapPolicy
	return p
}

// 设置文本
func NewLabel(text string) *Label {
	label := gmaterial.Label(theme.NewTheme(), gunit.Sp(18), text)
	return &Label{
		labelWidget: &label,
		margin:      &glayout.Inset{},
		config:      &labelConfig{},
	}
}

package widget

import (
	"github.com/Seikaijyu/nenki.ui/widget/anchor"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
)

// 校验接口是否实现
var _ WidgetInterface = &AnchorLayout{}
var _ SingleChildLayoutInterface[*AnchorLayout] = &AnchorLayout{}

// 锚点布局配置
type anchorLayoutConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

// 锚定布局
type AnchorLayout struct {
	// 配置
	config *anchorLayoutConfig
	// 内边距
	padding *glayout.Inset
	// 外边距
	margin *glayout.Inset
	// 子节点，可以为任意组件
	childWidget WidgetInterface
	// 配置
	direction anchor.Direction
}

// 绑定函数
func (p *AnchorLayout) Then(fn func(self *AnchorLayout)) *AnchorLayout {
	fn(p)
	return p
}

// 注册删除事件
func (p *AnchorLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *AnchorLayout) Update(update bool) {
	p.config.update = update
}

// 重新设置父节点
func (p *AnchorLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild()
	})
}

// 设置子节点
func (p *AnchorLayout) AppendChild(child WidgetInterface) *AnchorLayout {
	p.childWidget = child
	return p
}

// 获取子节点
func (p *AnchorLayout) GetChild() WidgetInterface {
	return p.childWidget
}

// 删除子节点
func (p *AnchorLayout) RemoveChild() *AnchorLayout {
	p.childWidget = nil
	return p
}

// 删除自身
func (p *AnchorLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.childWidget.Destroy()
	}
	p.config._destroy = nil
}

// 设置外边距
func (p *AnchorLayout) Margin(Top, Left, Bottom, Right float32) *AnchorLayout {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 设置锚定方向
func (p *AnchorLayout) Direction(direc anchor.Direction) *AnchorLayout {
	p.direction = direc
	return p
}

// 渲染
func (p *AnchorLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {
	if !p.config.update || p.childWidget == nil {
		return glayout.Dimensions{}
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.direction.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
			return p.childWidget.Layout(gtx)
		})
	})
}

// 创建锚点布局
func NewAnchorLayout(direction anchor.Direction) *AnchorLayout {
	widget := &AnchorLayout{
		// 无子节点
		childWidget: nil,
		direction:   direction,
		padding:     &glayout.Inset{},
		margin:      &glayout.Inset{},
		config:      &anchorLayoutConfig{update: true},
	}
	return widget
}

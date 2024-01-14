package widget

import (
	"image/color"

	glayout "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	gunit "gioui.org/unit"
)

// 校验接口是否实现
var _ WidgetInterface = &ContainerLayout{}
var _ SingleChildLayoutInterface[*ContainerLayout] = &ContainerLayout{}

type containerConfig struct {
	// 背景颜色
	background *color.NRGBA
}

// 容器布局，只用于包裹组件
type ContainerLayout struct {
	// 配置
	config *containerConfig
	// 内边距
	padding *glayout.Inset
	// 外边距
	margin *glayout.Inset
	// 子节点，可以为任意组件
	childWidget WidgetInterface
	// 组件是否被删除
	isRemove bool
}

// 绑定函数
func (p *ContainerLayout) Then(fn func(self *ContainerLayout)) *ContainerLayout {
	fn(p)
	return p
}

// 设置子节点
func (p *ContainerLayout) AppendChild(child WidgetInterface) *ContainerLayout {
	p.childWidget = child
	return p
}

// 获取子节点
func (p *ContainerLayout) GetChild() WidgetInterface {
	return p.childWidget
}

// 删除子节点
func (p *ContainerLayout) RemoveChild() *ContainerLayout {
	p.childWidget = nil
	return p
}

// 是否被删除
func (p *ContainerLayout) IsDestroy() bool {
	return p.isRemove
}

// 删除自身
func (p *ContainerLayout) Destroy() {
	// 如果有子节点
	if p.childWidget != nil {
		// 注销子节点
		p.childWidget.Destroy()
		// 断开子节点
		p.RemoveChild()
	}
	p.isRemove = true
}

// 设置外边距
func (p *ContainerLayout) Margin(Top, Left, Bottom, Right float32) *ContainerLayout {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 背景颜色
func (p *ContainerLayout) Background(r, g, b, a uint8) *ContainerLayout {
	p.config.background = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 渲染
func (p *ContainerLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		var stack = clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
		defer stack.Pop()
		// 设置背景颜色
		if p.config.background != nil {
			paint.ColorOp{Color: *p.config.background}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		// 如果有子节点
		if p.childWidget != nil {
			// 如果子节点被删除
			if p.childWidget.IsDestroy() {
				// 断开子节点
				p.RemoveChild()
			} else {
				return p.childWidget.Layout(gtx)
			}
		}
		return glayout.Dimensions{Size: gtx.Constraints.Max}
	})
}

// 创建锚点布局
func NewContainerLayout() *ContainerLayout {
	widget := &ContainerLayout{
		// 无子节点
		childWidget: nil,
		padding:     &glayout.Inset{},
		margin:      &glayout.Inset{},
		config:      &containerConfig{},
	}
	return widget
}

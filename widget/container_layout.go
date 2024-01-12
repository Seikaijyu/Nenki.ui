package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
)

// 校验接口是否实现
var _ WidgetInterface = &ContainerLayout{}
var _ SingleChildLayoutInterface[*ContainerLayout] = &ContainerLayout{}

// 容器布局，只用于包裹组件
type ContainerLayout struct {
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
func (p *ContainerLayout) Then(fn func(*ContainerLayout)) *ContainerLayout {
	fn(p)
	return p
}

// 设置子节点
func (p *ContainerLayout) AppendChild(child WidgetInterface) *ContainerLayout {
	p.childWidget = child
	return p
}

// 获取子节点
func (p *ContainerLayout) Child() WidgetInterface {
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

// 渲染
func (p *ContainerLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
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
	}
	return widget
}

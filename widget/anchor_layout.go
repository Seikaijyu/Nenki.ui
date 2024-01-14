package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	"nenki.ui/widget/anchor"
)

// 校验接口是否实现
var _ WidgetInterface = &AnchorLayout{}
var _ SingleChildLayoutInterface[*AnchorLayout] = &AnchorLayout{}

// 锚定布局
type AnchorLayout struct {
	// 内边距
	padding *glayout.Inset
	// 外边距
	margin *glayout.Inset
	// 子节点，可以为任意组件
	childWidget WidgetInterface
	// 配置
	direction anchor.Direction
	// 组件是否被删除
	isRemove bool
}

// 绑定函数
func (p *AnchorLayout) Then(fn func(self *AnchorLayout)) *AnchorLayout {
	fn(p)
	return p
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

// 是否被删除
func (p *AnchorLayout) IsDestroy() bool {
	return p.isRemove
}

// 删除自身
func (p *AnchorLayout) Destroy() {
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
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.direction.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
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
	}
	return widget
}

package widget

import (
	"image/color"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
)

// 校验接口是否实现
var _ WidgetInterface = &Border{}
var _ SingleChildLayoutInterface[*Border] = &Border{}

type Border struct {
	// 外边距
	margin *glayout.Inset
	// 内边距
	padding *glayout.Inset
	// 边框
	border *gwidget.Border
	// 间隔，用于占位
	spacer *glayout.Inset
	// 包裹的组件
	childWidget WidgetInterface
	// 组件是否被删除
	isRemove bool
}

// 绑定函数
func (p *Border) Then(fn func(*Border)) *Border {
	fn(p)
	return p
}

// 是否被删除
func (p *Border) IsDestroy() bool {
	return p.isRemove
}

// 注销自身，清理所有引用
func (p *Border) Destroy() {
	p.isRemove = true
}

// 设置子节点
func (p *Border) AppendChild(child WidgetInterface) *Border {
	p.childWidget = child
	return p
}

// 获取子节点
func (p *Border) GetChild() WidgetInterface {
	return p.childWidget
}

// 删除子节点
func (p *Border) RemoveChild() *Border {
	p.childWidget = nil
	return p
}

// 设置边框颜色
func (p *Border) Color(r, g, b, a uint8) *Border {
	p.border.Color = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 设置边框宽度
func (p *Border) Width(width float32) *Border {
	p.border.Width = gunit.Dp(width)
	p.spacer.Top = gunit.Dp(width)
	p.spacer.Left = gunit.Dp(width)
	p.spacer.Bottom = gunit.Dp(width)
	p.spacer.Right = gunit.Dp(width)

	return p
}

// 设置边框圆角
func (p *Border) CornerRadius(radius float32) *Border {
	p.border.CornerRadius = gunit.Dp(radius)
	return p
}

// 设置外边距
func (p *Border) Margin(Top, Left, Bottom, Right float32) *Border {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 设置内边距
func (p *Border) Padding(Top, Left, Bottom, Right float32) *Border {
	p.padding = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 渲染UI
func (p *Border) Layout(gtx glayout.Context) glayout.Dimensions {
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.border.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
			return p.padding.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
				return p.spacer.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
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
					return glayout.Dimensions{}
				})
			})
		})
	})

}

// 创建一个边框
func NewBorder(widget WidgetInterface) *Border {
	border := &Border{
		childWidget: widget,
		border: &gwidget.Border{
			Color: color.NRGBA{
				R: 0x00,
				G: 0x00,
				B: 0x00,
				A: 0xff,
			},
			Width: gunit.Dp(1),
		},
		margin:  &glayout.Inset{},
		padding: &glayout.Inset{},
		spacer:  &glayout.Inset{},
	}
	border.AppendChild(widget)
	return border
}

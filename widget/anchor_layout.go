package widget

import (
	glayout "gioui.org/layout"
)

var _ WidgetInterface = &AnchorLayout{}

// 锚点方向
type Direction uint8

const (
	// 顶部左边
	TopLeft Direction = iota
	// 顶部
	Top
	// 顶部右边
	TopRight
	// 右边
	Right
	// 底部右边
	BottomRight
	// 底部
	Bottom
	// 底部左边
	BottomLeft
	// 左边
	Left
	// 居中
	Center
)

// 锚定布局
type AnchorLayout struct {
	// 父节点布局
	parent WidgetInterface
	// 子节点
	child WidgetInterface
	// 锚定方向
	direction Direction
}

// 绑定函数
func (p *AnchorLayout) AndThen(fn func(*AnchorLayout) *AnchorLayout) *AnchorLayout {
	return fn(p)
}

// 设置子节点
func (p *AnchorLayout) SetChild(child WidgetInterface) *AnchorLayout {
	p.child = child
	return p
}

// 设置锚定方向
func (p *AnchorLayout) SetDirection(direc Direction) *AnchorLayout {
	p.direction = direc
	return p
}

// 获取锚定方向
func (p *AnchorLayout) Direction() Direction {
	return p.direction
}

// 渲染
func (p *AnchorLayout) Layout(gtx glayout.Context) glayout.Dimensions {

	return glayout.Direction(p.direction).Layout(gtx,
		func(gtx glayout.Context) glayout.Dimensions {
			if p.child != nil {
				return p.child.Layout(gtx)
			}
			return glayout.Dimensions{}
		},
	)
}

// 创建锚点布局
func NewAnchorLayout(direction Direction) *AnchorLayout {
	return &AnchorLayout{
		parent:    nil,
		child:     nil,
		direction: direction,
	}
}

package widget

import (
	glayout "gioui.org/layout"
)

// 校验接口是否实现
var _ WidgetInterface = &AnchorLayout{}
var _ SingleChildLayoutInterface[*AnchorLayout] = &AnchorLayout{}

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
	// ID
	id string
	// 子节点，可以为任意组件
	child WidgetInterface
	// 锚定方向
	direction Direction
	// 组件池所在的索引位置
	index int
	// 组件是否被删除
	isRemove bool
}

// 绑定函数
func (p *AnchorLayout) Then(fn func(*AnchorLayout)) *AnchorLayout {
	fn(p)
	return p
}

// 设置子节点
func (p *AnchorLayout) AppendChild(child WidgetInterface) *AnchorLayout {
	p.child = child
	return p
}

// 获取子节点
func (p *AnchorLayout) Child() WidgetInterface {
	return p.child
}

// 设置锚定方向
func (p *AnchorLayout) SetDirection(direc Direction) *AnchorLayout {
	p.direction = direc
	return p
}

// 删除子节点
func (p *AnchorLayout) RemoveChild() *AnchorLayout {
	p.child = nil
	return p
}

// 是否被删除
func (p *AnchorLayout) IsDestroy() bool {
	return p.isRemove
}

// 删除自身
func (p *AnchorLayout) Destroy() {
	// 如果有子节点
	if p.child != nil {
		// 注销子节点
		p.child.Destroy()
		// 断开子节点
		p.RemoveChild()
	}
	p.isRemove = true
	pool.RemoveAtIndex(p.id, p.index)
}

// 获取锚定方向
func (p *AnchorLayout) Direction() Direction {
	return p.direction
}

// 渲染
func (p *AnchorLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {

	return glayout.Direction(p.direction).Layout(gtx,
		func(gtx glayout.Context) glayout.Dimensions {
			// 如果有子节点
			if p.child != nil {
				// 如果子节点被删除
				if p.child.IsDestroy() {
					// 断开子节点
					p.RemoveChild()
				} else {
					return p.child.Layout(gtx)
				}
			}
			return glayout.Dimensions{}
		},
	)
}

// 创建锚点布局
func NewAnchorLayout(direction Direction) *AnchorLayout {
	widget := &AnchorLayout{
		id:        "",
		child:     nil,
		direction: direction,
	}
	widget.index = pool.AddWidget("", widget)
	return widget
}

// 从ID创建锚点布局
func NewAnchorLayoutWithID(id string, direction Direction) *AnchorLayout {
	widget := &AnchorLayout{
		id:        "#" + id,
		child:     nil,
		direction: direction,
	}
	widget.index = pool.AddWidget(id, widget)
	return widget
}

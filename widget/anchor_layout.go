package widget

import (
	glayout "gioui.org/layout"
)

// 校验接口是否实现
var _ WidgetInterface = &AnchorLayout{}
var _ SingleChildLayoutInterface[*AnchorLayout] = &AnchorLayout{}

// 锚点方向
type AnchorDirection uint8

const (
	// 顶部左边
	TopLeft AnchorDirection = iota
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

// 锚点布局配置
type AnchorLayoutConfig struct {
	// 锚定方向
	Direction AnchorDirection
}

// 锚定布局
type AnchorLayout struct {
	// 子节点，可以为任意组件
	childWidget WidgetInterface
	// 配置
	config *AnchorLayoutConfig
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
	p.childWidget = child
	return p
}

// 获取子节点
func (p *AnchorLayout) Child() WidgetInterface {
	return p.childWidget
}

// 设置锚定方向
func (p *AnchorLayout) SetDirection(direc AnchorDirection) *AnchorLayout {
	p.config.Direction = direc
	return p
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

// 获取锚定方向
func (p *AnchorLayout) Direction() AnchorDirection {
	return p.config.Direction
}

// 渲染
func (p *AnchorLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {
	return glayout.Direction(p.config.Direction).Layout(gtx,
		func(gtx glayout.Context) glayout.Dimensions {

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
		},
	)
}

// 创建锚点布局
func NewAnchorLayout(direction AnchorDirection) *AnchorLayout {
	widget := &AnchorLayout{
		childWidget: nil,
		config: &AnchorLayoutConfig{
			Direction: direction,
		},
	}
	return widget
}

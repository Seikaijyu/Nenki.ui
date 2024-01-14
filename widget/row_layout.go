package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	"nenki.ui/widget/anchor"
)

type RowLayout struct {
	// 是否被删除
	isRemove         bool
	margin           *glayout.Inset
	childWidgets     []WidgetInterface
	flexChilds       []glayout.FlexChild
	HorizontalWidget *glayout.Flex
}

// 校验接口是否实现
var _ WidgetInterface = &RowLayout{}
var _ MultiChildLayoutInterface[*RowLayout] = &RowLayout{}

// 绑定函数
func (p *RowLayout) Then(fn func(self *RowLayout)) *RowLayout {
	fn(p)
	return p
}

// 是否被删除
func (p *RowLayout) IsDestroy() bool {
	return p.isRemove
}

// 注销自身，清理所有引用
func (p *RowLayout) Destroy() {
	p.isRemove = true
}

// 添加子节点，并且可以定义权重
func (p *RowLayout) AppendFlexChild(weight float32, child WidgetInterface) *RowLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, func(gtx glayout.Context) glayout.Dimensions {
		return child.Layout(gtx)
	}))
	return p
}

// 添加子节点，可以根据方向进行布局堆叠
func (p *RowLayout) AppendFlexAnchorChild(weight float32, direction anchor.Direction, child WidgetInterface) *RowLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, func(gtx glayout.Context) glayout.Dimensions {
		return direction.Layout(gtx, child.Layout)
	}))
	return p
}

// 添加子节点，组件得到基于子组件的固定空间
func (p *RowLayout) AppendRigidChild(child WidgetInterface) *RowLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(func(gtx glayout.Context) glayout.Dimensions {
		return child.Layout(gtx)
	}))
	return p
}

// 从指定索引删除子节点
func (p *RowLayout) RemoveChildAt(index int) *RowLayout {
	// 现在进行删除操作
	if index >= 0 && index < len(p.childWidgets) {
		p.childWidgets = append(p.childWidgets[:index], p.childWidgets[index+1:]...)
		p.flexChilds = append(p.flexChilds[:index], p.flexChilds[index+1:]...)
	}
	return p
}

// 删除所有子节点
func (p *RowLayout) RemoveChildAll() *RowLayout {
	p.childWidgets = []WidgetInterface{}
	p.flexChilds = []glayout.FlexChild{}
	return p
}

// 获取所有子节点
func (p *RowLayout) GetChildAll() []WidgetInterface {
	return p.childWidgets
}

// 获取指定索引的子节点
func (p *RowLayout) GetChildAt(index int) WidgetInterface {
	if index >= 0 && index < len(p.childWidgets) {
		return p.childWidgets[index]
	}
	return nil
}

// 获取子节点数量
func (p *RowLayout) GetChildCount() int {
	return len(p.childWidgets)
}

// 设置外边距
func (p *RowLayout) Margin(Top, Left, Bottom, Right float32) *RowLayout {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 渲染UI
func (p *RowLayout) Layout(gtx glayout.Context) glayout.Dimensions {
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.HorizontalWidget.Layout(gtx, p.flexChilds...)
	})
}

// 创建一个水平布局
func NewRowLayout() *RowLayout {
	return &RowLayout{
		childWidgets:     []WidgetInterface{},
		margin:           &glayout.Inset{},
		flexChilds:       []glayout.FlexChild{},
		HorizontalWidget: &glayout.Flex{Axis: glayout.Horizontal},
	}
}

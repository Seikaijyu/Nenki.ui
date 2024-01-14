package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	"nenki.ui/widget/anchor"
)

type ColumnLayout struct {
	// 是否被删除
	isRemove       bool
	margin         *glayout.Inset
	childWidgets   []WidgetInterface
	flexChilds     []glayout.FlexChild
	verticalWidget *glayout.Flex
}

// 校验接口是否实现
var _ WidgetInterface = &ColumnLayout{}
var _ MultiChildLayoutInterface[*ColumnLayout] = &ColumnLayout{}

// 绑定函数
func (p *ColumnLayout) Then(fn func(*ColumnLayout)) *ColumnLayout {
	fn(p)
	return p
}

// 是否被删除
func (p *ColumnLayout) IsDestroy() bool {
	return p.isRemove
}

// 注销自身，清理所有引用
func (p *ColumnLayout) Destroy() {
	p.isRemove = true
}

// 添加多个子节点，并且可以定义权重
func (p *ColumnLayout) AppendFlexChild(weight float32, child WidgetInterface) *ColumnLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, func(gtx glayout.Context) glayout.Dimensions {
		return child.Layout(gtx)
	}))
	return p
}

// 添加多个子节点，可以根据方向进行布局堆叠
func (p *ColumnLayout) AppendStackedAnchorChild(direction anchor.Direction, child WidgetInterface) *ColumnLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(func(gtx glayout.Context) glayout.Dimensions {
		return glayout.Stack{Alignment: direction}.Layout(gtx, glayout.Stacked(func(gtx glayout.Context) glayout.Dimensions {
			return child.Layout(gtx)
		}))
	}))
	return p
}

// 添加多个子节点，可以进行堆叠
func (p *ColumnLayout) AppendStackedChild(child WidgetInterface) *ColumnLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(func(gtx glayout.Context) glayout.Dimensions {
		return glayout.Stack{}.Layout(gtx, glayout.Stacked(func(gtx glayout.Context) glayout.Dimensions {
			return child.Layout(gtx)
		}))
	}))
	return p
}

// 添加多个子节点，组件得到基于子组件的固定空间
func (p *ColumnLayout) AppendRigidChild(child WidgetInterface) *ColumnLayout {
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(func(gtx glayout.Context) glayout.Dimensions {
		return child.Layout(gtx)
	}))
	return p
}

// 从指定索引删除子节点
func (p *ColumnLayout) RemoveChildAt(index int) *ColumnLayout {
	// 现在进行删除操作
	if index >= 0 && index < len(p.childWidgets) {
		p.childWidgets = append(p.childWidgets[:index], p.childWidgets[index+1:]...)
		p.flexChilds = append(p.flexChilds[:index], p.flexChilds[index+1:]...)
	}
	return p
}

// 删除所有子节点
func (p *ColumnLayout) RemoveChildAll() *ColumnLayout {
	p.childWidgets = []WidgetInterface{}
	p.flexChilds = []glayout.FlexChild{}
	return p
}

// 获取所有子节点
func (p *ColumnLayout) GetChildAll() []WidgetInterface {
	return p.childWidgets
}

// 获取指定索引的子节点
func (p *ColumnLayout) GetChildAt(index int) WidgetInterface {
	if index >= 0 && index < len(p.childWidgets) {
		return p.childWidgets[index]
	}
	return nil
}

// 获取子节点数量
func (p *ColumnLayout) GetChildCount() int {
	return len(p.childWidgets)
}

// 设置外边距
func (p *ColumnLayout) Margin(Top, Left, Bottom, Right float32) *ColumnLayout {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 渲染UI
func (p *ColumnLayout) Layout(gtx glayout.Context) glayout.Dimensions {
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.verticalWidget.Layout(gtx, p.flexChilds...)
	})
}

// 创建一个垂直布局
func NewColumnLayout() *ColumnLayout {
	return &ColumnLayout{
		childWidgets:   []WidgetInterface{},
		margin:         &glayout.Inset{},
		flexChilds:     []glayout.FlexChild{},
		verticalWidget: &glayout.Flex{Axis: glayout.Vertical},
	}
}

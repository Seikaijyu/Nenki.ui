package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	"nenki.ui/widget/anchor"
)

type rowLayoutConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

type RowLayout struct {
	config           *rowLayoutConfig
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

// 注册删除事件
func (p *RowLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 注销自身，清理所有引用
func (p *RowLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.RemoveChildAll()
	}
	p.config._destroy = nil
}

// 是否可见
func (p *RowLayout) Update(update bool) {
	p.config.update = update
}

// 重新设置父节点
func (p *RowLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild(child)
	})
}

// 添加子节点，并且可以定义权重
func (p *RowLayout) AppendFlexChild(weight float32, child WidgetInterface) *RowLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, child.Layout))
	return p
}

// 添加子节点，可以根据方向进行布局堆叠
func (p *RowLayout) AppendFlexAnchorChild(weight float32, direction anchor.Direction, child WidgetInterface) *RowLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, func(gtx glayout.Context) glayout.Dimensions {
		return direction.Layout(gtx, child.Layout)
	}))
	return p
}

// 添加子节点，组件得到基于子组件的固定空间
func (p *RowLayout) AppendRigidChild(child WidgetInterface) *RowLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(child.Layout))
	return p
}

// 从组件删除子节点
func (p *RowLayout) RemoveChild(child WidgetInterface) *RowLayout {
	// 为了加速删除，这里使用异步删除
	go func() {
		for i, _child := range p.childWidgets {
			if _child == child {
				p.RemoveChildAt(i)
				break
			}
		}
	}()
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
	if !p.config.update {
		return glayout.Dimensions{}
	}
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
		config:           &rowLayoutConfig{update: true},
	}
}

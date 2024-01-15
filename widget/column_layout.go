package widget

import (
	"github.com/Seikaijyu/nenki.ui/widget/anchor"

	glayout "github.com/Seikaijyu/gio/layout"
	gunit "github.com/Seikaijyu/gio/unit"
)

type columnLayoutConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

type ColumnLayout struct {
	config         *columnLayoutConfig
	margin         *glayout.Inset
	childWidgets   []WidgetInterface
	flexChilds     []glayout.FlexChild
	verticalWidget *glayout.Flex
}

// 校验接口是否实现
var _ WidgetInterface = &ColumnLayout{}
var _ MultiChildLayoutInterface[*ColumnLayout] = &ColumnLayout{}

// 绑定函数
func (p *ColumnLayout) Then(fn func(self *ColumnLayout)) *ColumnLayout {
	fn(p)
	return p
}

// 注册删除事件
func (p *ColumnLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *ColumnLayout) Update(update bool) {
	p.config.update = update
}

// 注销自身，清理所有引用
func (p *ColumnLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.RemoveChildAll()
	}
	p.config._destroy = nil
}

// 重新设置父节点
func (p *ColumnLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild(child)
	})
}

// 添加子节点，并且可以定义权重
func (p *ColumnLayout) AppendFlexChild(weight float32, child WidgetInterface) *ColumnLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, child.Layout))
	return p
}

// 添加子节点，可以根据方向进行布局堆叠
func (p *ColumnLayout) AppendFlexAnchorChild(weight float32, direction anchor.Direction, child WidgetInterface) *ColumnLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, func(gtx glayout.Context) glayout.Dimensions {
		return direction.Layout(gtx, child.Layout)
	}))
	return p
}

// 添加子节点，组件得到基于子组件的固定空间
func (p *ColumnLayout) AppendRigidChild(child WidgetInterface) *ColumnLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(child.Layout))
	return p
}

// 从组件删除子节点
func (p *ColumnLayout) RemoveChild(child WidgetInterface) *ColumnLayout {
	go func() {
		for index, value := range p.childWidgets {
			if value == child {
				p.RemoveChildAt(index)
				break
			}
		}
	}()
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
	if !p.config.update {
		return glayout.Dimensions{}
	}
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
		config:         &columnLayoutConfig{update: true},
	}
}
package widget

import (
	"github.com/Seikaijyu/nenki.ui/widget/anchor"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
)

type columnLayoutConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

type ColumnLayout struct {
	config         *columnLayoutConfig
	margin         *glayout.Inset
	childWidgets   []WidgetInterface
	flexChilds     []glayout.FlexChild
	verticalWidget *glayout.Flex
}

// 校验接口是否实现
var _ WidgetInterface = &ColumnLayout{}
var _ MultiChildLayoutInterface[*ColumnLayout] = &ColumnLayout{}

// 绑定函数
func (p *ColumnLayout) Then(fn func(self *ColumnLayout)) *ColumnLayout {
	fn(p)
	return p
}

// 注册删除事件
func (p *ColumnLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *ColumnLayout) Update(update bool) {
	p.config.update = update
}

// 注销自身，清理所有引用
func (p *ColumnLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.RemoveChildAll()
	}
	p.config._destroy = nil
}

// 重新设置父节点
func (p *ColumnLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild(child)
	})
}

// 添加子节点，并且可以定义权重
func (p *ColumnLayout) AppendFlexChild(weight float32, child WidgetInterface) *ColumnLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, child.Layout))
	return p
}

// 添加子节点，可以根据方向进行布局堆叠
func (p *ColumnLayout) AppendFlexAnchorChild(weight float32, direction anchor.Direction, child WidgetInterface) *ColumnLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Flexed(weight, func(gtx glayout.Context) glayout.Dimensions {
		return direction.Layout(gtx, child.Layout)
	}))
	return p
}

// 添加子节点，组件得到基于子组件的固定空间
func (p *ColumnLayout) AppendRigidChild(child WidgetInterface) *ColumnLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	p.flexChilds = append(p.flexChilds, glayout.Rigid(child.Layout))
	return p
}

// 从组件删除子节点
func (p *ColumnLayout) RemoveChild(child WidgetInterface) *ColumnLayout {
	go func() {
		for index, value := range p.childWidgets {
			if value == child {
				p.RemoveChildAt(index)
				break
			}
		}
	}()
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
	if !p.config.update {
		return glayout.Dimensions{}
	}
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
		config:         &columnLayoutConfig{update: true},
	}
}

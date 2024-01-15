package widget

import (
	"image/color"

	glayout "github.com/Seikaijyu/gio/layout"
	"github.com/Seikaijyu/gio/op/clip"
	"github.com/Seikaijyu/gio/op/paint"
	gunit "github.com/Seikaijyu/gio/unit"
)

// 校验接口是否实现
var _ WidgetInterface = &ContainerLayout{}
var _ SingleChildLayoutInterface[*ContainerLayout] = &ContainerLayout{}

type containerConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 背景颜色
	background *color.NRGBA
}

// 容器布局，只用于包裹组件
type ContainerLayout struct {
	// 配置
	config *containerConfig
	// 外边距
	margin *glayout.Inset
	// 子节点，可以为任意组件
	childWidget WidgetInterface
}

// 绑定函数
func (p *ContainerLayout) Then(fn func(self *ContainerLayout)) *ContainerLayout {
	fn(p)
	return p
}

// 注册删除事件
func (p *ContainerLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *ContainerLayout) Update(update bool) {
	p.config.update = update
}

// 重新设置父节点
func (p *ContainerLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild()
	})
}

// 设置子节点
func (p *ContainerLayout) AppendChild(child WidgetInterface) *ContainerLayout {
	p.ResetParent(child)
	p.childWidget = child
	return p
}

// 获取子节点
func (p *ContainerLayout) GetChild() WidgetInterface {
	return p.childWidget
}

// 删除子节点
func (p *ContainerLayout) RemoveChild() *ContainerLayout {
	p.childWidget = nil
	return p
}

// 删除自身
func (p *ContainerLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.childWidget.Destroy()
	}
	p.config._destroy = nil
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

// 背景颜色
func (p *ContainerLayout) Background(r, g, b, a uint8) *ContainerLayout {
	p.config.background = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 渲染
func (p *ContainerLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {
	if !p.config.update || p.childWidget == nil {
		return glayout.Dimensions{Size: gtx.Constraints.Max}
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		var stack = clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
		defer stack.Pop()
		// 设置背景颜色
		if p.config.background != nil {
			paint.ColorOp{Color: *p.config.background}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		return p.childWidget.Layout(gtx)

	})
}

// 创建锚点布局
func NewContainerLayout() *ContainerLayout {
	widget := &ContainerLayout{
		// 无子节点
		childWidget: nil,
		margin:      &glayout.Inset{},
		config:      &containerConfig{update: true},
	}
	return widget
}
package widget

import (
	"image/color"

	glayout "gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	gunit "gioui.org/unit"
)

// 校验接口是否实现
var _ WidgetInterface = &ContainerLayout{}
var _ SingleChildLayoutInterface[*ContainerLayout] = &ContainerLayout{}

type containerConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 背景颜色
	background *color.NRGBA
}

// 容器布局，只用于包裹组件
type ContainerLayout struct {
	// 配置
	config *containerConfig
	// 外边距
	margin *glayout.Inset
	// 子节点，可以为任意组件
	childWidget WidgetInterface
}

// 绑定函数
func (p *ContainerLayout) Then(fn func(self *ContainerLayout)) *ContainerLayout {
	fn(p)
	return p
}

// 注册删除事件
func (p *ContainerLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *ContainerLayout) Update(update bool) {
	p.config.update = update
}

// 重新设置父节点
func (p *ContainerLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild()
	})
}

// 设置子节点
func (p *ContainerLayout) AppendChild(child WidgetInterface) *ContainerLayout {
	p.ResetParent(child)
	p.childWidget = child
	return p
}

// 获取子节点
func (p *ContainerLayout) GetChild() WidgetInterface {
	return p.childWidget
}

// 删除子节点
func (p *ContainerLayout) RemoveChild() *ContainerLayout {
	p.childWidget = nil
	return p
}

// 删除自身
func (p *ContainerLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.childWidget.Destroy()
	}
	p.config._destroy = nil
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

// 背景颜色
func (p *ContainerLayout) Background(r, g, b, a uint8) *ContainerLayout {
	p.config.background = &color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 渲染
func (p *ContainerLayout) Layout(gtx glayout.Context) (dimensions glayout.Dimensions) {
	if !p.config.update || p.childWidget == nil {
		return glayout.Dimensions{Size: gtx.Constraints.Max}
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		var stack = clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops)
		defer stack.Pop()
		// 设置背景颜色
		if p.config.background != nil {
			paint.ColorOp{Color: *p.config.background}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
		}
		return p.childWidget.Layout(gtx)

	})
}

// 创建锚点布局
func NewContainerLayout() *ContainerLayout {
	widget := &ContainerLayout{
		// 无子节点
		childWidget: nil,
		margin:      &glayout.Inset{},
		config:      &containerConfig{update: true},
	}
	return widget
}

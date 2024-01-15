package widget

import (
	"image/color"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
	"nenki.ui/widget/axis"
)

// 列表配置
type listLayoutConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

// 列表布局
type ListLayout struct {
	// 配置
	config       *listLayoutConfig
	margin       *glayout.Inset
	childWidgets []WidgetInterface
	listWidget   *glayout.List
	listMaterial *gmaterial.ListStyle
}

// 校验接口是否实现
var _ WidgetInterface = &ListLayout{}
var _ MultiChildLayoutInterface[*ListLayout] = &ListLayout{}

// 绑定函数
func (p *ListLayout) Then(fn func(self *ListLayout)) *ListLayout {
	fn(p)
	return p
}

// 注册删除事件
func (p *ListLayout) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 删除组件
func (p *ListLayout) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.RemoveChildAll()
	}

	p.config._destroy = nil
}

// 是否更新组件
func (p *ListLayout) Update(update bool) {
	p.config.update = update
}

// 重新设置父节点
func (p *ListLayout) ResetParent(child WidgetInterface) {
	child.Destroy()
	child.Update(true)
	child.OnDestroy(func() {
		child.Update(false)
		p.RemoveChild(child)
	})
}

// 添加子节点
func (p *ListLayout) AppendChild(child WidgetInterface) *ListLayout {
	p.ResetParent(child)
	p.childWidgets = append(p.childWidgets, child)
	return p
}

// 从组件中删除子节点
func (p *ListLayout) RemoveChild(child WidgetInterface) *ListLayout {
	for index, childWidget := range p.childWidgets {
		if childWidget == child {
			p.RemoveChildAt(index)
			break
		}
	}
	return p
}

// 从指定索引删除子节点
func (p *ListLayout) RemoveChildAt(index int) *ListLayout {
	// 现在进行删除操作
	if index >= 0 && index < len(p.childWidgets) {
		p.childWidgets = append(p.childWidgets[:index], p.childWidgets[index+1:]...)
	}
	return p
}

// 删除所有子节点
func (p *ListLayout) RemoveChildAll() *ListLayout {
	p.childWidgets = []WidgetInterface{}
	return p
}

// 获取所有子节点
func (p *ListLayout) GetChildAll() []WidgetInterface {
	return p.childWidgets
}

// 获取指定索引的子节点
func (p *ListLayout) GetChildAt(index int) WidgetInterface {
	if index >= 0 && index < len(p.childWidgets) {
		return p.childWidgets[index]
	}
	return nil
}

// 获取子节点数量
func (p *ListLayout) GetChildCount() int {
	return len(p.childWidgets)
}

// 当设置为true，会让列表在更新item后保持滚动到最后更新位置
func (p *ListLayout) ScrollToEnd(scrollToEnd bool) *ListLayout {
	p.listWidget.ScrollToEnd = scrollToEnd
	return p
}

// 滚动到列表的指定位置
func (p *ListLayout) ScrollBy(offset float32) *ListLayout {
	p.listWidget.ScrollBy(offset)
	return p
}

// 滚动到列表中存在的指定项的位置
func (p *ListLayout) ScrollToItem(index int) *ListLayout {
	p.listWidget.ScrollTo(index)
	return p
}

// 滚动条背景颜色
func (p *ListLayout) ScrollBgColor(r, g, b, a uint8) *ListLayout {
	p.listMaterial.Track.Color = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 滚动条的间距宽高
func (p *ListLayout) ScrollPadding(width, height float32) *ListLayout {
	p.listMaterial.Track.MinorPadding = gunit.Dp(width)
	p.listMaterial.Track.MajorPadding = gunit.Dp(height)
	return p
}

// 滚动条的宽度
func (p *ListLayout) ScrollWidth(width float32) *ListLayout {
	p.listMaterial.Indicator.MinorWidth = gunit.Dp(width)
	return p
}

// 滚动条最小长度
func (p *ListLayout) ScrollMinLen(minLen float32) *ListLayout {
	p.listMaterial.Indicator.MajorMinLen = gunit.Dp(minLen)
	return p
}

// 滚动条的圆角
func (p *ListLayout) ScrollCornerRadius(radius float32) *ListLayout {
	p.listMaterial.Indicator.CornerRadius = gunit.Dp(radius)
	return p
}

// 滚动条鼠标默认颜色
func (p *ListLayout) ScrollColor(r, g, b, a uint8) *ListLayout {
	p.listMaterial.Indicator.Color = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 滚动条鼠标悬浮颜色
func (p *ListLayout) ScrollHoverColor(r, g, b, a uint8) *ListLayout {
	p.listMaterial.Indicator.HoverColor = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 是否启用滚动条间距，如果禁用，滚动条会覆着在列表上
func (p *ListLayout) ScrollPaddingEnable(enable bool) *ListLayout {
	if enable {
		p.listMaterial.AnchorStrategy = gmaterial.Occupy
	} else {
		p.listMaterial.AnchorStrategy = gmaterial.Overlay
	}
	return p
}

// 设置外边距
func (p *ListLayout) Margin(Top, Left, Bottom, Right float32) *ListLayout {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(Top),
		Left:   gunit.Dp(Left),
		Bottom: gunit.Dp(Bottom),
		Right:  gunit.Dp(Right),
	}
	return p
}

// 渲染UI
func (p *ListLayout) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		return glayout.Dimensions{}
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.listMaterial.Layout(gtx, len(p.childWidgets), func(gtx glayout.Context, index int) glayout.Dimensions {
			return p.childWidgets[index].Layout(gtx)
		})
	})
}

// 设置方向
func (p *ListLayout) Axis(axis axis.Axis) *ListLayout {
	p.listWidget.Axis = axis
	return p
}

// 创建一个指定方向的列表布局
func NewListLayout(axis axis.Axis) *ListLayout {
	listWidget := gwidget.List{}
	listMaterial := gmaterial.List(&gmaterial.Theme{}, &listWidget)
	listWidget.Axis = axis
	return &ListLayout{
		childWidgets: []WidgetInterface{},
		margin:       &glayout.Inset{},
		listWidget:   &listWidget.List,
		listMaterial: &listMaterial,
		config:       &listLayoutConfig{update: true},
	}
}

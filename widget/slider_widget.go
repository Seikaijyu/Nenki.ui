package widget

import (
	"image/color"

	glayout "github.com/Seikaijyu/gio/layout"
	gunit "github.com/Seikaijyu/gio/unit"
	"github.com/Seikaijyu/gio/widget"
	gmaterial "github.com/Seikaijyu/gio/widget/material"
	"github.com/Seikaijyu/nenki.ui/widget/axis"
	"github.com/Seikaijyu/nenki.ui/widget/theme"
)

// 校验接口是否实现
var _ WidgetInterface = &Slider{}

type sliderConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 拖动事件
	_dragging func(*Slider, float32)
}

type Slider struct {
	// 配置
	config *sliderConfig
	// 外边距
	margin *glayout.Inset
	// 滑块组件
	slider *gmaterial.SliderStyle
}

// 绑定函数
func (p *Slider) Then(fn func(self *Slider)) *Slider {
	fn(p)
	return p
}

// 注销自身，清理所有引用
func (p *Slider) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
	}
	p.config._destroy = nil
}

// 注册删除事件
func (p *Slider) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 注册拖动事件
func (p *Slider) OnDragging(fn func(p *Slider, value float32)) *Slider {
	p.config._dragging = fn
	return p
}

// 是否更新组件
func (p *Slider) Update(update bool) {
	p.config.update = update
}

// 外边距
func (p *Slider) Margin(Top, Left, Bottom, Right float32) *Slider {
	p.margin.Top = gunit.Dp(Top)
	p.margin.Left = gunit.Dp(Left)
	p.margin.Bottom = gunit.Dp(Bottom)
	p.margin.Right = gunit.Dp(Right)
	return p
}

// 滑块的颜色
func (p *Slider) Color(r, g, b, a uint8) *Slider {
	p.slider.Color = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 滑块可点击的范围
func (p *Slider) FingerSize(size float32) *Slider {
	p.slider.FingerSize = gunit.Dp(size)
	return p
}

// 获取滑块的值
func (p *Slider) GetValue() float32 {
	return p.slider.Float.Value
}

// 渲染组件
func (p *Slider) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		p.config.update = false
	}
	if p.config._dragging != nil && p.slider.Float.Dragging() {
		p.config._dragging(p, p.slider.Float.Value)
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.slider.Layout(gtx)
	})
}

// 创建滑块组件，值为0-1
func NewSlider(axis axis.Axis) *Slider {
	slider := gmaterial.Slider(theme.NewTheme(), &widget.Float{})
	slider.Axis = axis
	return &Slider{
		slider: &slider,
		margin: &glayout.Inset{},
		config: &sliderConfig{},
	}
}

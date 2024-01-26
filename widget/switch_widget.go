package widget

import (
	"image/color"

	glayout "github.com/Seikaijyu/gio/layout"
	gunit "github.com/Seikaijyu/gio/unit"
	gwidget "github.com/Seikaijyu/gio/widget"
	gmaterial "github.com/Seikaijyu/gio/widget/material"
	"github.com/Seikaijyu/nenki.ui/widget/theme"
)

// 校验接口是否实现
var _ WidgetInterface = &Switch{}

type switchConfig struct {
	prevValue bool
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 选择事件
	_change func(*Switch, bool)
}

type Switch struct {
	// 配置
	config *switchConfig
	// 外边距
	margin *glayout.Inset
	// 组件
	switchWidget *gmaterial.SwitchStyle
}

// 绑定函数
func (p *Switch) Then(fn func(self *Switch)) *Switch {
	fn(p)
	return p
}

// 注销自身，清理所有引用
func (p *Switch) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
	}
	p.config._destroy = nil
}

// 注册删除事件
func (p *Switch) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *Switch) Update(update bool) {
	p.config.update = update
}

// 外边距
func (p *Switch) Margin(Top, Left, Bottom, Right float32) *Switch {
	p.margin.Top = gunit.Dp(Top)
	p.margin.Left = gunit.Dp(Left)
	p.margin.Bottom = gunit.Dp(Bottom)
	p.margin.Right = gunit.Dp(Right)
	return p
}

// 选择事件
func (p *Switch) OnChange(fn func(p *Switch, value bool)) *Switch {
	p.config._change = fn
	return p
}

// 启用颜色
func (p *Switch) EnabledColor(r, g, b, a uint8) *Switch {
	p.switchWidget.Color.Enabled = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 禁用颜色
func (p *Switch) DisabledColor(r, g, b, a uint8) *Switch {
	p.switchWidget.Color.Disabled = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 轨道颜色
func (p *Switch) TrackColor(r, g, b, a uint8) *Switch {
	p.switchWidget.Color.Track = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 启用开关
func (p *Switch) Enabled(value bool) *Switch {
	p.switchWidget.Switch.Value = value
	return p
}

// 布局
func (p *Switch) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		p.config.update = false
	}
	if p.config._change != nil && p.config.prevValue != p.switchWidget.Switch.Value {
		p.config._change(p, p.switchWidget.Switch.Value)
	}

	p.config.prevValue = p.switchWidget.Switch.Value
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.switchWidget.Layout(gtx)
	})
}

func NewSwitch() *Switch {
	switchWidget := gmaterial.Switch(theme.NewTheme(), &gwidget.Bool{}, "")
	return &Switch{
		switchWidget: &switchWidget,
		margin:       &glayout.Inset{},
		config:       &switchConfig{},
	}
}

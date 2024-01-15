package widget

import (
	"image/color"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
	"nenki.ui/widget/axis"
	"nenki.ui/widget/theme"
)

// 校验接口是否实现
var _ WidgetInterface = &RadioButtons{}

type radioButtonsConfig struct {
	// 文字颜色
	color color.NRGBA
	// 勾选框左边标记的颜色
	iconColor color.NRGBA
	// 记录选择的key
	selectKey string
	// 尺寸
	size float32
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 鼠标悬浮事件
	_hovered func(*RadioButtons, string)
	// 焦点事件
	_focused func(*RadioButtons, string)
	// 选择事件
	_selected func(*RadioButtons, string)
}

// 复选框
type RadioButtons struct {
	// 固定布局
	flexChilds []glayout.FlexChild
	// 单选组件组
	radioButtonWidgets []*gmaterial.RadioButtonStyle
	// 单选组件
	radioEnum *gwidget.Enum
	// 主题
	radioTheme *gmaterial.Theme
	// 布局
	flexWidget *glayout.Flex
	// 配置
	config *radioButtonsConfig
	// 外边距
	margin *glayout.Inset
}

// 绑定函数
func (p *RadioButtons) Then(fn func(self *RadioButtons)) *RadioButtons {
	fn(p)
	return p
}

// 注册删除事件
func (p *RadioButtons) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *RadioButtons) Update(update bool) {
	p.config.update = update
}

// 添加子节点
func (p *RadioButtons) AppendRadioButton(key, text string) *RadioButtons {
	radio := gmaterial.RadioButton(p.radioTheme, p.radioEnum, key, text)
	radio.Color = p.config.color
	radio.IconColor = p.config.iconColor
	radio.TextSize = gunit.Sp(p.config.size)
	radio.Size = gunit.Dp(p.config.size + 1 + p.config.size*0.3)
	pradio := &radio
	p.radioButtonWidgets = append(p.radioButtonWidgets, pradio)
	p.flexChilds = append(p.flexChilds,
		glayout.Rigid(func(gtx glayout.Context) glayout.Dimensions {
			return (*pradio).Layout(gtx)
		}),
	)
	return p
}

// 鼠标悬浮事件
func (p *RadioButtons) OnHovered(fn func(p *RadioButtons, key string)) *RadioButtons {
	p.config._hovered = fn
	return p
}

// 焦点事件
func (p *RadioButtons) OnFocused(fn func(p *RadioButtons, key string)) *RadioButtons {
	p.config._focused = fn
	return p
}

// 选择事件
func (p *RadioButtons) OnSelected(fn func(p *RadioButtons, key string)) *RadioButtons {
	p.config._selected = fn
	return p
}

// 注销自身，清理所有引用
func (p *RadioButtons) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.flexChilds = []glayout.FlexChild{}
		p.radioButtonWidgets = []*gmaterial.RadioButtonStyle{}
	}
	p.config._destroy = nil
}

// 外边距
func (p *RadioButtons) Margin(Top, Left, Bottom, Right float32) *RadioButtons {
	p.margin.Top = gunit.Dp(Top)
	p.margin.Left = gunit.Dp(Left)
	p.margin.Bottom = gunit.Dp(Bottom)
	p.margin.Right = gunit.Dp(Right)
	return p
}

// 重新设置父节点
func (p *RadioButtons) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		p.config.update = false
	}
	if p.config._hovered != nil {
		if value, ok := p.radioEnum.Hovered(); ok {
			p.config._hovered(p, value)
		}
	}
	if value, ok := p.radioEnum.Focused(); ok {
		if p.config._focused != nil {
			p.config._focused(p, value)
		}
		if p.config._selected != nil && p.config.selectKey != value {
			p.config._selected(p, value)
		}
		p.config.selectKey = value
	}
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.flexWidget.Layout(gtx, p.flexChilds...)
	})
}

// 设置字体大小
func (p *RadioButtons) Size(size float32) *RadioButtons {
	if size < 16 {
		size = 16
	} else if size > 30 {
		size = 30
	}
	for _, v := range p.radioButtonWidgets {
		v.TextSize = gunit.Sp(size)
		v.Size = gunit.Dp(size + 1 + size*0.3)
	}

	p.config.size = size

	return p
}

// 设置文字颜色
func (p *RadioButtons) FontColor(r, g, b, a uint8) *RadioButtons {
	for _, v := range p.radioButtonWidgets {
		v.Color = color.NRGBA{R: r, G: g, B: b, A: a}
	}
	p.config.color = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置单选框左边标记的颜色
func (p *RadioButtons) RadioMarkColor(r, g, b, a uint8) *RadioButtons {
	for _, v := range p.radioButtonWidgets {
		v.IconColor = color.NRGBA{R: r, G: g, B: b, A: a}
	}
	p.config.iconColor = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置方向
func (p *RadioButtons) Axis(axis axis.Axis) *RadioButtons {
	p.flexWidget.Axis = axis
	return p
}

// 创建一个指定方向的复选框
func NewRadioButtons(axis axis.Axis) *RadioButtons {
	radioWidget := &RadioButtons{
		radioEnum:  &gwidget.Enum{},
		radioTheme: theme.NewTheme(),
		flexChilds: []glayout.FlexChild{},
		config: &radioButtonsConfig{
			update: true, size: 16,
			color:     color.NRGBA{A: 255},
			iconColor: color.NRGBA{R: 63, G: 81, B: 181, A: 255},
		},
		margin:     &glayout.Inset{},
		flexWidget: &glayout.Flex{Axis: axis},
	}
	return radioWidget
}

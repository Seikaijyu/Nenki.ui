package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
	"gioui.org/widget/material"
	gmaterial "gioui.org/widget/material"
	"nenki.ui/widget/axis"
)

// 校验接口是否实现
var _ WidgetInterface = &RadioButtons{}

type radioButtonsConfig struct {
	// 尺寸
	size float32
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 鼠标悬浮事件
	_hovered func(*RadioButtons)
	// 点击事件
	_pressed func(*RadioButtons)
}

// 复选框
type RadioButtons struct {
	flexChilds         []glayout.FlexChild
	radioButtonWidgets []gmaterial.RadioButtonStyle
	radioEnum          *gwidget.Enum
	radioTheme         *gmaterial.Theme
	flexWidget         *glayout.Flex
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
	radio.TextSize = gunit.Sp(p.config.size)
	radio.Size = gunit.Dp(p.config.size + 1 + p.config.size*0.3)
	p.radioButtonWidgets = append(p.radioButtonWidgets, radio)
	p.flexChilds = append(p.flexChilds,
		glayout.Rigid(radio.Layout),
	)
	return p
}

// 设置字体大小
func (p *RadioButtons) Size(size float32) *RadioButtons {
	if size < 16 {
		size = 16
	} else if size > 30 {
		size = 30
	}
	p.config.size = size

	return p
}

// 鼠标悬浮事件
func (p *RadioButtons) OnHovered(fn func(*RadioButtons)) *RadioButtons {
	p.config._hovered = fn
	return p
}

// 鼠标按下事件
func (p *RadioButtons) OnPressed(fn func(*RadioButtons)) *RadioButtons {
	p.config._pressed = fn
	return p
}

// 获取选中的按钮的Key
func (p *RadioButtons) GetSelectedKey() (string, bool) {
	return p.radioEnum.Focused()
}

// 获取选中的按钮的值
func (p *RadioButtons) GetSelectedValue() string {
	return p.radioEnum.Value
}

// 注销自身，清理所有引用
func (p *RadioButtons) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
		p.flexChilds = []glayout.FlexChild{}
		p.radioButtonWidgets = []gmaterial.RadioButtonStyle{}
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
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.flexWidget.Layout(gtx, p.flexChilds...)
	})
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
		radioTheme: material.NewTheme(),
		flexChilds: []glayout.FlexChild{},
		config:     &radioButtonsConfig{update: true, size: 20},
		margin:     &glayout.Inset{},
		flexWidget: &glayout.Flex{Axis: axis},
	}
	return radioWidget
}

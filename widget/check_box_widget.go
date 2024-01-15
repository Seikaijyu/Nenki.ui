package widget

import (
	"image/color"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
)

// 校验接口是否实现
var _ WidgetInterface = &CheckBox{}

type checkBoxConfig struct {
	// 记录选择的bool
	checkedBool bool
	// 焦点事件
	_focused func(*CheckBox, bool)
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 鼠标悬浮事件
	_hovered func(*CheckBox, bool)
	// 选择事件
	_checked func(*CheckBox, bool)
}

// 复选框
type CheckBox struct {
	checkBoxWidget *gmaterial.CheckBoxStyle
	checkBool      *gwidget.Bool
	// 配置
	config *checkBoxConfig
	// 外边距
	margin *glayout.Inset
}

// 绑定函数
func (p *CheckBox) Then(fn func(self *CheckBox)) *CheckBox {
	fn(p)
	return p
}

// 注册删除事件
func (p *CheckBox) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *CheckBox) Update(update bool) {
	p.config.update = update
}

// 注销自身，清理所有引用
func (p *CheckBox) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
	}
	p.config._destroy = nil
}

// 外边距
func (p *CheckBox) Margin(Top, Left, Bottom, Right float32) *CheckBox {
	p.margin.Top = gunit.Dp(Top)
	p.margin.Left = gunit.Dp(Left)
	p.margin.Bottom = gunit.Dp(Bottom)
	p.margin.Right = gunit.Dp(Right)
	return p
}

// 重新设置父节点
func (p *CheckBox) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		p.config.update = false
	}
	if p.config._hovered != nil && p.checkBoxWidget.CheckBox.Hovered() {
		p.config._hovered(p, p.checkBool.Value)
	}

	if p.config._focused != nil && p.checkBoxWidget.CheckBox.Focused() {
		p.config._focused(p, p.checkBool.Value)
	}

	p.config.checkedBool = p.checkBool.Value
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.checkBoxWidget.Layout(gtx)
	})
}

// 设置尺寸
func (p *CheckBox) Size(size float32) *CheckBox {
	if size < 16 {
		size = 16
	} else if size > 30 {
		size = 30
	}
	p.checkBoxWidget.TextSize = gunit.Sp(size)
	p.checkBoxWidget.Size = gunit.Dp(size + 1 + size*0.3)
	return p
}

// 鼠标悬浮事件
func (p *CheckBox) OnHovered(fn func(p *CheckBox, check bool)) *CheckBox {
	p.config._hovered = fn
	return p
}

// 选择事件
func (p *CheckBox) OnChecked(fn func(p *CheckBox, check bool)) *CheckBox {
	p.config._checked = fn
	p.checkBool.OnChecked(func(b bool) {
		if p.config._checked != nil {
			p.config._checked(p, b)
		}
	})

	return p
}

// 设置文字颜色
func (p *CheckBox) FontColor(r, g, b, a uint8) *CheckBox {
	p.checkBoxWidget.Color = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置勾选框左边标记的颜色
func (p *CheckBox) CheckMarkColor(r, g, b, a uint8) *CheckBox {
	p.checkBoxWidget.IconColor = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置文字
func (p *CheckBox) Text(text string) *CheckBox {
	p.checkBoxWidget.Label = text
	return p
}

// 设置是否勾选
func (p *CheckBox) Checked(checked bool) *CheckBox {
	p.checkBool.Value = checked
	return p
}

// 获取是否勾选
func (p *CheckBox) GetChecked() bool {
	return p.checkBool.Value
}

// 焦点事件
func (p *CheckBox) OnFocused(fn func(*CheckBox, bool)) *CheckBox {
	p.config._focused = fn
	return p
}

// 创建复选框
func NewCheckBox(text string) *CheckBox {
	checkBool := gwidget.Bool{}
	widget := gmaterial.CheckBox(gmaterial.NewTheme(), &checkBool, text)
	checkBox := &CheckBox{
		checkBool:      &checkBool,
		checkBoxWidget: &widget,
		config:         &checkBoxConfig{update: true, checkedBool: false},
		margin:         &glayout.Inset{},
	}
	return checkBox.Size(16)
}

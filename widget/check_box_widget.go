package widget

import (
	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
)

// 校验接口是否实现
var _ WidgetInterface = &CheckBox{}

type checkBoxConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
}

// 复选框
type CheckBox struct {
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
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return glayout.Dimensions{}
	})
}

func NewCheckBox() *CheckBox {
	return &CheckBox{
		config: &checkBoxConfig{update: true},
		margin: &glayout.Inset{},
	}
}

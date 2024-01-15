package widget

import (
	"image/color"
	"time"

	"github.com/Seikaijyu/nenki.ui/widget/theme"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
)

// 按钮配置
type buttonConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 鼠标点击事件
	_clicked func(*Button)
	// 模拟点击
	_click func()
	// 鼠标悬浮事件
	_hovered func(*Button)
	// 鼠标按下事件
	_pressed func(*Button)
	// 长点击
	_longClicked func(*Button)
	// 焦点事件
	_focused func(*Button)
}
type Button struct {
	margin *glayout.Inset
	// 配置
	config *buttonConfig
	// 主题
	button gmaterial.ButtonStyle
	// 记录点击次数
	clickCount int
	// 上一次点击时间
	lastClickTime time.Time
}

// 校验接口是否实现
var _ WidgetInterface = &Button{}

// 绑定函数
func (p *Button) Then(fn func(self *Button)) *Button {
	fn(p)
	return p
}

// 注册删除事件
func (p *Button) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 删除组件
func (p *Button) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
	}
	p.config._destroy = nil
}

// 是否更新组件
func (p *Button) Update(update bool) {
	p.config.update = update
}

// 设置焦点为按钮
func (p *Button) Focus() *Button {
	p.button.Button.Focus()
	return p
}

// 设置字体大小
func (p *Button) FontSize(size int) *Button {
	p.button.TextSize = gunit.Sp(size)
	return p
}

// 设置圆角
func (p *Button) CornerRadius(radius float32) *Button {
	p.button.CornerRadius = gunit.Dp(radius)
	return p
}

// 设置文字颜色
func (p *Button) FontColor(r, g, b, a uint8) *Button {
	p.button.Color = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置背景颜色
func (p *Button) Background(r, g, b, a uint8) *Button {
	p.button.Background = color.NRGBA{R: r, G: g, B: b, A: a}
	return p
}

// 设置内边距
func (p *Button) Padding(top, left, bottom, right float32) *Button {
	p.button.Inset = glayout.Inset{
		Top:    gunit.Dp(top),
		Left:   gunit.Dp(left),
		Bottom: gunit.Dp(bottom),
		Right:  gunit.Dp(right),
	}
	return p
}

// 设置外边距
func (p *Button) Margin(top, left, bottom, right float32) *Button {
	p.margin = &glayout.Inset{
		Top:    gunit.Dp(top),
		Left:   gunit.Dp(left),
		Bottom: gunit.Dp(bottom),
		Right:  gunit.Dp(right),
	}
	return p
}

// 设置文本
func (p *Button) Text(text string) *Button {
	p.button.Text = text
	return p
}

// 获取焦点状态
func (p *Button) OnFocused(fn func(*Button)) *Button {
	p.config._focused = fn
	return p
}

// 模拟点击
func (p *Button) Click() *Button {
	if p.config._click != nil {
		// 模拟点击，触发点击事件
		p.button.Button.Click()
	}
	return p
}

// 点击事件
func (p *Button) OnClicked(fn func(*Button)) *Button {
	p.config._clicked = fn
	return p
}

// 鼠标悬浮事件
func (p *Button) OnHovered(fn func(*Button)) *Button {
	p.config._hovered = fn
	return p
}

// 鼠标按下事件
func (p *Button) OnPressed(fn func(*Button)) *Button {
	p.config._pressed = fn
	return p
}

// 多次点击事件
func (p *Button) OnLongClicked(fn func(*Button)) *Button {
	p.config._longClicked = fn
	return p
}

// 布局
func (p *Button) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		return glayout.Dimensions{}
	}
	// 悬浮事件
	if p.config._hovered != nil && p.button.Button.Hovered() {
		p.config._hovered(p)
	}
	// 按下事件
	if p.config._pressed != nil && p.button.Button.Pressed() {
		p.config._pressed(p)
	}
	// 焦点事件
	if p.config._focused != nil && p.button.Button.Focused() {
		p.config._focused(p)
	}
	// 双击事件，其中实现了点击事件
	if p.config._longClicked != nil {
		// 拿到点击事件的具体状态
		if history := p.button.Button.History(); len(history) > 0 {
			// 拿到最后一次点击事件
			press := history[len(history)-1]
			// 如果最后一次点击事件的结束时间为 0，说明还没点击完
			if press.End.Second() != 0 {
				// 计算点击事件的时间间隔
				clickInterval := press.End.Sub(press.Start).Milliseconds()
				if p.lastClickTime.Second() == 0 || clickInterval < 200 {
					if clickInterval < 200 {
						intervalSinceLastClick := press.End.Sub(p.lastClickTime).Milliseconds()
						if intervalSinceLastClick != 0 {
							p.clickCount++
							if intervalSinceLastClick < 200 {
								if p.clickCount >= 2 {
									if p.config._longClicked != nil {
										// 作为长点击事件触发
										p.config._longClicked(p)
									}
									// 重置点击次数
									p.clickCount = 0
								}
							} else {
								if p.config._clicked != nil {
									// 作为点击事件触发
									p.config._clicked(p)
								}
							}
						}
					} else {
						// 重置点击次数
						p.clickCount = 0
					}
					p.lastClickTime = press.End

				}

			}

		} else {
			// 重置点击次数
			p.clickCount = 0
		}
		// 为了确保单击事件捕获性能，仅在有双击事件的时候才会使用以上方式捕获点击事件（使用以上方式捕获点击事件会减缓响应时间至少200ms）
	} else if p.config._clicked != nil && p.button.Button.Clicked(gtx) {
		// 点击事件
		p.config._clicked(p)
	}
	// 外边距
	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		// 按钮
		return p.button.Layout(gtx)
	})
}

// 创建按钮
func NewButton(text string) *Button {
	widget := &Button{
		clickCount:    0,
		lastClickTime: time.Time{},
		margin:        &glayout.Inset{},
		config:        &buttonConfig{update: true},
		button:        gmaterial.Button(theme.NewTheme(), &gwidget.Clickable{}, text),
	}
	return widget
}

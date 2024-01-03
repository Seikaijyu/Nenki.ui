package widget

import (
	"image/color"
	"time"

	glayout "gioui.org/layout"
	"gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
	"nenki.ui/widget/edge"
)

// 按钮配置
type ButtonConfig struct {
	// 文本
	Text string
	// 字体大小
	FontSize int
	// 圆角
	CornerRadius float32
	// 文字颜色
	FontColor *color.NRGBA
	// 背景颜色
	Background *color.NRGBA
	// 内边距
	Padding edge.Padding
	// 鼠标点击事件
	Clicked func(*Button)
	// 模拟点击
	Click func()
	// 鼠标悬浮事件
	Hovered func(*Button)
	// 鼠标按下事件
	Pressed func(*Button)
	// 长点击
	LongClicked func(*Button)
}
type Button struct {
	// 配置
	config *ButtonConfig
	// 主题
	button gmaterial.ButtonStyle
	// 组件是否被删除
	isRemove bool
	// 记录点击次数
	clickCount int
	// 上一次点击时间
	lastClickTime time.Time
}

// 校验接口是否实现
var _ WidgetInterface = &Button{}

// 绑定函数
func (p *Button) Then(fn func(*Button)) *Button {
	fn(p)
	return p
}

// 点击事件
func (p *Button) OnClicked(fn func(*Button)) *Button {
	p.config.Clicked = fn
	return p
}

// 模拟点击
func (p *Button) Click() *Button {
	if p.config.Click != nil {
		// 模拟点击，触发点击事件
		p.button.Button.Click()
	}
	return p
}

// 鼠标悬浮事件
func (p *Button) OnHovered(fn func(*Button)) *Button {
	p.config.Hovered = fn
	return p
}

// 鼠标按下事件
func (p *Button) OnPressed(fn func(*Button)) *Button {
	p.config.Pressed = fn
	return p
}

// 多次点击事件
func (p *Button) OnLongClicked(fn func(*Button)) *Button {
	p.config.LongClicked = fn
	return p
}

// 设置焦点为按钮
func (p *Button) SetFocus() *Button {
	p.button.Button.Focus()
	return p
}

// 设置焦点状态
func (p *Button) Focused() bool {
	return p.button.Button.Focused()
}

// 是否被删除
func (p *Button) IsDestroy() bool {
	return p.isRemove
}

// 注销自身，清理所有引用
func (p *Button) Destroy() {
	p.isRemove = true
}

// 设置字体大小
func (p *Button) SetFontSize(size int) *Button {
	p.config.FontSize = size
	p.button.TextSize = unit.Sp(size)
	return p
}

// 获取字体大小
func (p *Button) FontSize() int {
	return p.config.FontSize
}

// 设置圆角
func (p *Button) SetCornerRadius(radius float32) *Button {
	p.config.CornerRadius = radius
	p.button.CornerRadius = unit.Dp(radius)
	return p
}

// 获取圆角
func (p *Button) CornerRadius() float32 {
	return p.config.CornerRadius
}

// 设置文字颜色
func (p *Button) SetFontColor(r, g, b, a uint8) *Button {
	p.config.FontColor = &color.NRGBA{R: r, G: g, B: b, A: a}
	p.button.Color = *p.config.FontColor
	return p
}

// 获取文字颜色
func (p *Button) FontColor() (uint8, uint8, uint8, uint8) {
	return p.config.FontColor.R, p.config.FontColor.G, p.config.FontColor.B, p.config.FontColor.A
}

// 设置背景颜色
func (p *Button) SetBackground(r, g, b, a uint8) *Button {
	p.config.Background = &color.NRGBA{R: r, G: g, B: b, A: a}
	p.button.Background = *p.config.Background
	return p
}

// 获取背景颜色
func (p *Button) Background() (uint8, uint8, uint8, uint8) {
	return p.config.Background.R, p.config.Background.G, p.config.Background.B, p.config.Background.A
}

// 设置内边距
func (p *Button) SetPadding(top, left, bottom, right float32) *Button {
	p.config.Padding = edge.Padding{
		Top:    top,
		Left:   left,
		Bottom: bottom,
		Right:  right,
	}
	p.button.Inset = glayout.Inset{
		Top:    unit.Dp(top),
		Left:   unit.Dp(left),
		Bottom: unit.Dp(bottom),
		Right:  unit.Dp(right),
	}
	return p
}

// 获取内边距
func (p *Button) Padding() (float32, float32, float32, float32) {
	return p.config.Padding.Top, p.config.Padding.Left, p.config.Padding.Bottom, p.config.Padding.Right
}

// 设置文本
func (p *Button) SetText(text string) *Button {
	p.config.Text = text
	p.button.Text = text
	return p
}

// 获取文本
func (p *Button) Text() string {
	return p.config.Text
}

// 布局
func (p *Button) Layout(gtx glayout.Context) glayout.Dimensions {

	// 悬浮事件
	if p.config.Hovered != nil && p.button.Button.Hovered() {
		p.config.Hovered(p)
	}
	// 按下事件
	if p.config.Pressed != nil && p.button.Button.Pressed() {
		p.config.Pressed(p)
	}
	// 双击事件，其中实现了点击事件
	if p.config.LongClicked != nil {
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
									if p.config.LongClicked != nil {
										// 作为长点击事件触发
										p.config.LongClicked(p)
									}
									// 重置点击次数
									p.clickCount = 0
								}
							} else {
								if p.config.Clicked != nil {
									// 作为点击事件触发
									p.config.Clicked(p)
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
		// 为了确保性能，仅在有双击事件的时候才会使用以上方式捕获点击事件
	} else if p.config.Clicked != nil && p.button.Button.Clicked(gtx) {
		// 点击事件
		p.config.Clicked(p)
	}
	return p.button.Layout(gtx)
}

// 创建按钮
func NewButton(text string) *Button {
	widget := &Button{
		clickCount:    0,
		lastClickTime: time.Time{},
		config: &ButtonConfig{
			Text: text,
		},
		button: gmaterial.Button(gmaterial.NewTheme(), &gwidget.Clickable{}, text),
	}
	return widget

}

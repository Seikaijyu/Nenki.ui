package widget

import (
	"image/color"

	"github.com/Seikaijyu/nenki.ui/widget/text"
	"github.com/Seikaijyu/nenki.ui/widget/theme"

	glayout "gioui.org/layout"
	gunit "gioui.org/unit"
	gwidget "gioui.org/widget"
	gmaterial "gioui.org/widget/material"
)

// 校验接口是否实现
var _ WidgetInterface = &Editor{}

// 编辑框配置
type editorConfig struct {
	// 是否更新组件
	update bool
	// 删除事件
	_destroy func()
	// 鼠焦点事件
	_focused func(*Editor, string)
	// 回车事件，仅在单行编辑时有效
	_submit func(*Editor, string)
	// 选择事件
	_select func(*Editor, string)
	// 文本改变事件
	_change func(*Editor, string)
}

// 编辑框
type Editor struct {
	// 配置
	config *editorConfig
	// 外边距
	margin *glayout.Inset
	// editor组件
	editorMaterial *gmaterial.EditorStyle
}

// 绑定函数
func (p *Editor) Then(fn func(self *Editor)) *Editor {
	fn(p)
	return p
}

// 注册删除事件
func (p *Editor) OnDestroy(fn func()) {
	p.config._destroy = fn
}

// 是否更新组件
func (p *Editor) Update(update bool) {
	p.config.update = update
}

// 注销自身，清理所有引用
func (p *Editor) Destroy() {
	p.config.update = false
	if p.config._destroy != nil {
		p.config._destroy()
	}
	p.config._destroy = nil
}

// 外边距
func (p *Editor) Margin(Top, Left, Bottom, Right float32) *Editor {
	p.margin.Top = gunit.Dp(Top)
	p.margin.Left = gunit.Dp(Left)
	p.margin.Bottom = gunit.Dp(Bottom)
	p.margin.Right = gunit.Dp(Right)
	return p
}

// 渲染
func (p *Editor) Layout(gtx glayout.Context) glayout.Dimensions {
	if !p.config.update {
		return glayout.Dimensions{}
	}
	for _, item := range p.editorMaterial.Editor.Events() {
		switch item.(type) {
		case gwidget.ChangeEvent:
			if p.config._change != nil {
				p.config._change(p, p.GetText())
			}
		case gwidget.SubmitEvent:
			if p.config._submit != nil {
				p.config._submit(p, p.GetText())
			}
		case gwidget.SelectEvent:
			if p.config._select != nil {
				p.config._select(p, p.GetSelectedText())
			}
		}
	}
	if p.config._focused != nil {
		if p.editorMaterial.Editor.Focused() {
			p.config._focused(p, p.GetText())
		}
	}

	return p.margin.Layout(gtx, func(gtx glayout.Context) glayout.Dimensions {
		return p.editorMaterial.Layout(gtx)
	})
}

// 设置只读
func (p *Editor) ReadOnly(readOnly bool) *Editor {
	p.editorMaterial.Editor.ReadOnly = readOnly
	return p
}

// 设置文字
func (p *Editor) Text(text string) *Editor {
	p.editorMaterial.Editor.SetText(text)
	return p
}

// 设置单行编辑
func (p *Editor) SingleLine(singleLine bool) *Editor {
	p.editorMaterial.Editor.SingleLine = singleLine
	return p
}

// 文本对齐方向
func (p *Editor) Alignment(alignment text.Alignment) *Editor {
	p.editorMaterial.Editor.Alignment = alignment
	return p
}

// 仅允许输入指定字符
func (p *Editor) AllowOnly(filter string) *Editor {
	p.editorMaterial.Editor.Filter = filter
	return p
}

// 设置提示文字
func (p *Editor) Hint(hint string) *Editor {
	p.editorMaterial.Hint = hint
	return p
}

// 设置提示文字颜色
func (p *Editor) HintColor(r, g, b, a uint8) *Editor {
	p.editorMaterial.HintColor = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 设置文字行高
func (p *Editor) LineHeight(lineHeight float32) *Editor {
	p.editorMaterial.Editor.LineHeight = gunit.Sp(lineHeight)
	return p
}

// 设置文字行高缩放
func (p *Editor) LineHeightScale(lineHeightScale float32) *Editor {
	p.editorMaterial.Editor.LineHeightScale = lineHeightScale
	return p
}

// 设置文字掩盖符号，一般用于密码输入
func (p *Editor) Mask(mask rune) *Editor {
	p.editorMaterial.Editor.Mask = mask
	return p
}

// 将回车键作为提交按钮，启用时文本框自动变为单行文本框
func (p *Editor) Submit(submit bool) *Editor {
	p.editorMaterial.Editor.Submit = submit
	return p
}

// 设置最可输入长度
func (p *Editor) MaxLength(maxLines int) *Editor {
	p.editorMaterial.Editor.MaxLen = maxLines
	return p
}

// 设置如何显示文本换行
func (p *Editor) WrapPolicy(wrapPolicy text.WrapPolicy) *Editor {
	p.editorMaterial.Editor.WrapPolicy = wrapPolicy
	return p
}

// 移动光标到 start 的位置，并将选择区域的结束位置设定到 end。start 和 end 是以字符计数的，代表在编辑器文本中的偏移量。
func (p *Editor) SetCaret(start, end int) *Editor {
	p.editorMaterial.Editor.SetCaret(start, end)
	return p
}

// 获取光标位置的字符串
func (p *Editor) GetSelectedText() string {
	return p.editorMaterial.Editor.SelectedText()
}

// 返回选择区域的开始和结束位置，开始的位置可以大于结束的位置
func (p *Editor) GetSelection() (start, end int) {
	return p.editorMaterial.Editor.Selection()
}

// 返回选择区域的长度
func (p *Editor) GetSelectionLen() int {
	return p.editorMaterial.Editor.SelectionLen()
}

// 清除选择区域，通过将选择区域的结束位置设定为开始位置
func (p *Editor) ClearSelection() *Editor {
	p.editorMaterial.Editor.ClearSelection()
	return p
}

// 方法会返回光标所在的行和列的编号
func (p *Editor) CaretPos() (line, col int) {
	return p.editorMaterial.Editor.CaretPos()
}

// 方法会从光标位置删除字符。参数的符号指定删除的方向：正数表示向前删除，负数表示向后删除。
//
// 如果有选定的文字，那么会被删除，这会算作删除了一个字符串。
func (p *Editor) Delete(count int) *Editor {
	p.editorMaterial.Editor.Delete(count)
	return p
}

// 方法会请求为编辑器获取输入焦点
func (p *Editor) Focus() *Editor {
	p.editorMaterial.Editor.Focus()
	return p
}

// 方法返回编辑器是否处于焦点状态。
func (p *Editor) OnFocused(fn func(p *Editor, text string)) *Editor {
	p.config._focused = fn
	return p
}

// 提交事件
func (p *Editor) OnSubmit(fn func(p *Editor, text string)) *Editor {
	p.config._submit = fn
	return p
}

// 选择事件
func (p *Editor) OnSelect(fn func(p *Editor, text string)) *Editor {
	p.config._select = fn
	return p
}

// 文本改变事件
func (p *Editor) OnChange(fn func(p *Editor, text string)) *Editor {
	p.config._change = fn
	return p
}

// 方法在选中处插入一段字符串
func (p *Editor) Insert(text string) *Editor {
	p.editorMaterial.Editor.Insert(text)
	return p
}

// 方法移动光标（也就是选择开始）和选择结束，相对于它们当前的位置。正距离向前移动，负距离向后移动。距离是以字形簇为单位，即使字符由多个代码点组成，这也非常接近用户所认为的"字符"。
func (p *Editor) MoveCaret(startDelta, endDelta int) *Editor {
	p.editorMaterial.Editor.MoveCaret(startDelta, endDelta)
	return p
}

// 获取文本框的文字
func (p *Editor) GetText() string {

	return p.editorMaterial.Editor.Text()
}

// 获取文本框文本长度
func (p *Editor) GetTextLen() int {
	return p.editorMaterial.Editor.Len()
}

// 设置文本框选中部分的背景颜色
func (p *Editor) SelectionColor(r, g, b, a uint8) *Editor {
	p.editorMaterial.SelectionColor = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 设置文本框字体大小
func (p *Editor) FontSize(textSize float32) *Editor {
	p.editorMaterial.TextSize = gunit.Sp(textSize)
	return p
}

// 设置文本框字体粗细
func (p *Editor) TextWeight(textWeight text.Weight) *Editor {
	p.editorMaterial.Font.Weight = textWeight
	return p
}

// 设置文本框文本颜色
func (p *Editor) TextColor(r, g, b, a uint8) *Editor {
	p.editorMaterial.Color = color.NRGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
	return p
}

// 设置屏幕键盘类型
func (p *Editor) KeyboardType(keyboardType text.InputHint) *Editor {
	p.editorMaterial.Editor.InputHint = keyboardType
	return p
}

// 创建编辑框
func NewEditor(hint string) *Editor {
	editorMaterial := gmaterial.Editor(theme.NewTheme(), &gwidget.Editor{}, hint)
	return &Editor{
		config:         &editorConfig{update: true},
		margin:         &glayout.Inset{},
		editorMaterial: &editorMaterial,
	}
}

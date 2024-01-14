package text

import gfont "gioui.org/font"

// Weight 类型定义了字体的粗细程度，其底层类型是 gfont.Weight
type Weight = gfont.Weight

// 这里定义了一些常用的字体粗细等级常量
const (
	// Thin 代表极细字体，值为-300
	Thin Weight = -300
	// ExtraLight 代表超轻字体，值为-200
	ExtraLight Weight = -200
	// Light 代表轻字体，值为-100
	Light Weight = -100
	// Normal 代表正常字体，值为0
	Normal Weight = 0
	// Medium 代表中等字体，值为100
	Medium Weight = 100
	// SemiBold 代表半粗字体，值为200
	SemiBold Weight = 200
	// Bold 代表粗字体，值为300
	Bold Weight = 300
	// ExtraBold 代表超粗字体，值为400
	ExtraBold Weight = 400
	// Black 代表最粗的字体，值为500
	Black Weight = 500
)

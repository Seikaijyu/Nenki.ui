package anchor

import glayout "gioui.org/layout"

// 锚点方向
type Direction = glayout.Direction

const (
	// 顶部左边
	TopLeft Direction = iota
	// 顶部
	Top
	// 顶部右边
	TopRight
	// 右边
	Right
	// 底部右边
	BottomRight
	// 底部
	Bottom
	// 底部左边
	BottomLeft
	// 左边
	Left
	// 居中
	Center
)

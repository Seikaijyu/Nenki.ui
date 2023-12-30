package widget

import (
	glayout "gioui.org/layout"
)

// 定义一个通用的布局接口
type WidgetInterface interface {
	Layout(gtx glayout.Context) glayout.Dimensions
}

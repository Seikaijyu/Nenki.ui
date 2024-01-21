package text

import gtext "github.com/Seikaijyu/gio/text"

// 设置文本对齐
type Alignment = gtext.Alignment

const (
	// 开始方向
	Start Alignment = iota
	// 结束方向
	End
	// 中间方向
	Middle
)

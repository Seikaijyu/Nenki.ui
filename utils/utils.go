package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// 根据16进制字符串将其转换为RGBA值
// 如果字符串有任何的16进制错误则返回255,255,255,255
// 否则根据情况隐式处理，即：如果字符串不足8位则补全，填充完全使用f，如果字符串超过8位则截断
func HexToRGBA(hex string) (uint8, uint8, uint8, uint8) {
	// 默认值，如果有任何16进制错误
	r, g, b, a := uint8(255), uint8(255), uint8(255), uint8(255)

	// 补全位数
	if len(hex) < 7 {
		hex = fmt.Sprintf("#%s", strings.TrimPrefix(hex, "#"))
		hex = strings.TrimRight(hex, "#")
		hex = fmt.Sprintf("%-8s", hex)
		hex = strings.Replace(hex, " ", "f", -1)
	}
	// 将16进制字符串转换为RGBA值
	if rVal, err := strconv.ParseUint(hex[1:3], 16, 8); err == nil {
		r = uint8(rVal)
	}
	if gVal, err := strconv.ParseUint(hex[3:5], 16, 8); err == nil {
		g = uint8(gVal)
	}
	if bVal, err := strconv.ParseUint(hex[5:7], 16, 8); err == nil {
		b = uint8(bVal)
	}

	return r, g, b, a
}

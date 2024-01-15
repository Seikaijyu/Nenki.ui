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
	r, g, b, a := uint8(0), uint8(0), uint8(0), uint8(255)
	hex = strings.TrimPrefix(hex, "#")
	length := len(hex)

	// 补全位数
	if length == 3 {
		hex = fmt.Sprintf("%s%s%s%s%s%s", string(hex[0]), string(hex[0]), string(hex[1]), string(hex[1]), string(hex[2]), string(hex[2]))
	} else if length < 8 {
		hex = fmt.Sprintf("%-8s", hex)
		hex = strings.Replace(hex, " ", "f", -1)
	} else if length > 8 {
		hex = hex[:8]
	}

	// 将16进制字符串转换为RGBA值
	if rVal, err := strconv.ParseUint(hex[0:2], 16, 8); err == nil {
		r = uint8(rVal)
	}
	if gVal, err := strconv.ParseUint(hex[2:4], 16, 8); err == nil {
		g = uint8(gVal)
	}
	if bVal, err := strconv.ParseUint(hex[4:6], 16, 8); err == nil {
		b = uint8(bVal)
	}
	if aVal, err := strconv.ParseUint(hex[6:8], 16, 8); err == nil {
		a = uint8(aVal)
	}
	return r, g, b, a
}

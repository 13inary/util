package util

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

// StrTerminalLen 获取字符串在终端中输出的长度
// 例子：
// 中文 => 4
// en   => 2
func StrTerminalLen(str string) int {
	byteLen := len(str)
	charLen := utf8.RuneCountInString(str)
	chLen := (byteLen - charLen) / 2 // 1个中文 == 3个字节
	enLen := charLen - chLen
	return chLen*2 + enLen
}

func Float2String(f float64, suffix string) string {
	b := make([]byte, 0, 16+len(suffix))      // 预先分配足够的空间
	b = strconv.AppendFloat(b, f, 'f', 2, 64) // 2位小数
	b = append(b, suffix...)
	return string(b)
}

func Int2String(i int, suffix string) string {
	b := make([]byte, 0, 16+len(suffix)) // 预先分配足够的空间
	b = strconv.AppendInt(b, int64(i), 10)
	b = append(b, suffix...)
	return string(b)
}

func HasPrefixs(str string, prefixs []string) bool {
	if str == "" || len(prefixs) == 0 {
		return false
	}
	for _, prefix := range prefixs {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

func HasSuffixs(str string, suffixs []string) bool {
	if str == "" || len(suffixs) == 0 {
		return false
	}
	for _, suffix := range suffixs {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

package util

import "strconv"

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

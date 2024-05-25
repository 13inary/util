package util

import (
	"bytes"
	"strings"
)

// TableStr 获取可以展示表格的字符串
// 当 colMaxWidth 为nil，花费时间来自动计算
func TableStr(lines [][]string, colMaxWidth []int) string {
	// 获取列最小长度
	if colMaxWidth == nil {
		colMaxWidth = make([]int, 0, 16)
		for _, line := range lines {
			// 保证 colMaxWidth 数组空间足够
			if len(line) > len(colMaxWidth) {
				extended := make([]int, len(line))
				copy(extended, colMaxWidth)
				colMaxWidth = extended
			}
			// 从当前行获取各个列的长度最大值
			for k, col := range line {
				realLen := StrTerminalLen(col)
				if realLen > colMaxWidth[k] {
					colMaxWidth[k] = realLen
				}
			}
		}
	}

	// 生成表格
	colInterval := 2    // 列间隔的空格数
	colDecollator := "" // 列之间的分割符号
	var buffer bytes.Buffer
	buffer.Grow(len(lines) * 100)
	for _, line := range lines {
		for k, col := range line {
			var newCol string
			realLen := StrTerminalLen(col)
			if realLen > colMaxWidth[k] {
				newCol = col + colDecollator
			} else {
				newCol = col + strings.Repeat(" ", colMaxWidth[k]-realLen+colInterval) + colDecollator
			}
			buffer.WriteString(newCol)
		}
		buffer.WriteString("\n")
	}
	return buffer.String()
}

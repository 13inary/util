package util

import (
	"fmt"
	"os"
)

const (
	// ANSI 转义序列
	reset         = "\033[0m"
	reverse       = "\033[7m" // 反转显示效果（反色）
	bold          = "\033[1m"
	moveUp        = "\033[%dA"  // 向上移动光标
	moveDown      = "\033[%dB"  // 向下移动光标
	moveRight     = "\033[%dC"  // 向右移动光标
	moveLeft      = "\033[%dD"  // 向左移动光标
	saveCursor    = "\033[s"    // 保存光标位置
	restoreCursor = "\033[u"    // 恢复光标位置
	clearLine     = "\033[2K"   // 清除整行
	clearScreen   = "\033[2J"   // 清屏
	moveToHome    = "\033[H"    // 移动光标到左上角
	moveToCol     = "\033[%dG"  // 移动到指定列
	hideCursor    = "\033[?25l" // 隐藏光标
	showCursor    = "\033[?25h" // 显示光标
)

type printDiff struct {
	lines     []string
	prevLines []string
}

func NewPrintCharPositionDiff() *printDiff {
	return &printDiff{
		lines:     make([]string, 0),
		prevLines: make([]string, 0),
	}
}

// Start 初始化输出环境（隐藏光标、清屏）
func (lu *printDiff) Start() {
	// 隐藏光标
	fmt.Print(hideCursor)
	// 清屏
	fmt.Print(clearScreen)
	fmt.Print(moveToHome)
}

// Close 清理输出环境（显示光标、换行）
func (lu *printDiff) Close() {
	// 显示光标
	fmt.Print(showCursor)
	// 换行
	fmt.Print("\n")
}

// 设置所有行
func (lu *printDiff) SetLines(lines []string) {
	// 保存旧的行数（通过计算 prevLines 的长度得到）
	oldLineCount := len(lu.prevLines)

	// 移动到第一行（如果之前有内容，向上移动；如果是第一次调用，也移动到第一行）
	n := len(lu.prevLines)
	if n == 0 {
		n = len(lines)
	}
	if n > 0 {
		fmt.Printf(moveUp, n)
	}

	// 更新每一行
	for i, line := range lines {
		lu.updateLine(i, line)
		if i < len(lines)-1 {
			fmt.Print("\n")
		}
	}

	// 如果行数减少，清除多余的行
	if len(lines) < oldLineCount {
		for i := len(lines); i < oldLineCount; i++ {
			fmt.Print("\n")
			fmt.Print(clearLine)
		}
		// 移回最后一行
		fmt.Printf(moveUp, oldLineCount-len(lines))
	}

	// 刷新输出
	_ = os.Stdout.Sync()

	// 保存当前状态
	lu.lines = make([]string, len(lines))
	copy(lu.lines, lines)
	lu.prevLines = make([]string, len(lines))
	copy(lu.prevLines, lines)
}

// 更新单行，只刷新变化的部分
func (lu *printDiff) updateLine(lineNum int, newContent string) {
	if lineNum >= len(lu.prevLines) {
		// 新行，直接输出（正常显示，不加反转效果，因为这是初始内容）
		fmt.Print(newContent)
		return
	}

	oldContent := lu.prevLines[lineNum]
	if oldContent == newContent {
		// 内容没有变化，不更新
		return
	}

	// 计算段列表（同步比较rune，连续相同/不同作为segment）
	segments := lu.calculateSegments(oldContent, newContent)

	// 移动到该行的开始位置
	fmt.Printf(moveToCol, 1)
	fmt.Print(reset)

	// 遍历所有segment，重新输出整行
	for _, seg := range segments {
		if seg.isFixed {
			// 固定段：正常显示
			fmt.Print(reset + seg.content)
		} else {
			// 变化段：带反转效果
			fmt.Print(reverse + seg.content + reset)
		}
	}

	// 如果旧内容比新内容长，清除末尾多余的字符
	oldWidth := lu.displayWidth(oldContent)
	newWidth := lu.displayWidth(newContent)
	if oldWidth > newWidth {
		for i := newWidth; i < oldWidth; i++ {
			fmt.Print(" ")
		}
	}
}

// 计算字符的显示宽度（中文字符占2个位置，英文字符占1个位置）
func (lu *printDiff) displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r >= 0x1100 && (r <= 0x115F || r >= 0x2E80 && r <= 0x9FFF ||
			r >= 0xAC00 && r <= 0xD7AF || r >= 0xF900 && r <= 0xFAFF ||
			r >= 0xFE30 && r <= 0xFE4F || r >= 0x1F200 && r <= 0x1F2FF) {
			width += 2 // 中文字符或宽字符
		} else {
			width += 1 // 普通字符
		}
	}
	return width
}

// 段类型：固定段或变化段
type segment struct {
	isFixed bool   // true表示固定段，false表示变化段
	content string // 段内容
}

// 计算两个字符串的差异（同步比较rune，连续相同/不同作为segment）
func (lu *printDiff) calculateSegments(old, new string) []segment {
	oldRunes := []rune(old)
	newRunes := []rune(new)
	var segments []segment
	idx := 0

	for idx < len(oldRunes) || idx < len(newRunes) {
		// 找到连续相同的部分（固定段）
		start := idx
		for idx < len(oldRunes) && idx < len(newRunes) && oldRunes[idx] == newRunes[idx] {
			idx++
		}
		if idx > start {
			segments = append(segments, segment{isFixed: true, content: string(newRunes[start:idx])})
		}

		// 找到连续不同的部分（变化段）
		start = idx
		for idx < len(oldRunes) || idx < len(newRunes) {
			if idx < len(oldRunes) && idx < len(newRunes) && oldRunes[idx] == newRunes[idx] {
				break
			}
			idx++
		}
		if idx > start {
			segments = append(segments, segment{isFixed: false, content: string(newRunes[start:idx])})
		}
	}

	return segments
}

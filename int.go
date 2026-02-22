package util

// AbsInt 绝对值，返回true表示是负数
func AbsInt(n int) (bool, int) {
	if n < 0 {
		return true, -n
	}
	return false, n
}

package util

import "os"

// FileExists 检查指定路径的文件是否存在
// 参数:
//
//	path - 要检查的文件路径
//
// 返回值:
//
//	bool - 文件存在返回true，否则false
//	error - 底层系统调用返回的错误（当错误不是文件不存在时）
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

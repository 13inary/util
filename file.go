package util

import (
	"encoding/gob"
	"os"
)

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

// 异步写文件。写完文件并关闭后，若在操作系统同步文件前断电，会丢失文件。
func AtomicWriteSmallFile(file string, content []byte, perm os.FileMode) error {
	tmpName := file + ".tmp"
	if err := os.WriteFile(tmpName, content, perm); err != nil { // 覆盖式写入
		return err
	}
	return os.Rename(tmpName, file) // 对于目标已经存在，linux和mac直接覆盖，windows需要自己先删除后Rename
}

// 同步写文件。写完文件后，让操作系统同步文件，然后再关闭文件，避免文件的丢失。
// 同步操作会损耗性能
func SyncAtomicWriteSmallFile(file string, content []byte, perm os.FileMode) error {
	tmpName := file + ".tmp"

	f, err := os.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(content); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	return os.Rename(tmpName, file) // 对于目标已经存在，linux和mac直接覆盖，windows需要自己先删除后Rename
}

func SaveToCache[T any](data T, cacheFile string) error {
	f, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer f.Close()

	return gob.NewEncoder(f).Encode(data)
}

func LoadFromCache[T any](cacheFile string) (T, error) {
	var data T

	f, err := os.Open(cacheFile)
	if err != nil {
		return data, err
	}
	defer f.Close()

	return data, gob.NewDecoder(f).Decode(&data)
}

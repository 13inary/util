package util

import (
	"bytes"
	"io"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// GBKToUTF8 将GBK编码字节切片转换为UTF-8编码字节切片
func GBKToUTF8(gbk []byte) ([]byte, error) {
	// 预分配缓冲区避免扩容。预分配2倍原始大小缓冲区（中文UTF-8通常比GBK长）
	buf := bytes.NewBuffer(make([]byte, 0, len(gbk)*2))
	// 使用transform链式调用减少嵌套
	reader := transform.NewReader(
		bytes.NewReader(gbk),
		simplifiedchinese.GBK.NewDecoder(),
	)

	if _, err := io.Copy(buf, reader); err != nil { // io.ReadAll()适合小数据和代码简洁。io.Copy()适合大数据和高性能
		return nil, err
	}
	return buf.Bytes(), nil
}

// UTF8ToGBK 将UTF-8编码字节切片转换为GBK编码字节切片
func UTF8ToGBK(utf8 []byte) ([]byte, error) {
	// 预分配缓冲区（按GBK最大压缩比预估）
	buf := bytes.NewBuffer(make([]byte, 0, len(utf8)/2*3))
	writer := transform.NewWriter(
		buf,
		simplifiedchinese.GBK.NewEncoder(),
	)

	// 使用defer确保writer关闭
	defer writer.Close()

	if _, err := writer.Write(utf8); err != nil { // io.ReadAll()适合小数据和代码简洁。io.Copy()适合大数据和高性能
		return nil, err
	}
	return buf.Bytes(), nil
}

// TrimZeroBytes 移除字节切片首尾的0x00值
func TrimZeroBytes(b []byte) []byte {
	if len(b) == 0 {
		return b
	}
	start, end := 0, len(b)-1

	// 处理前面
	for ; start < end && b[start] == 0x00; start++ {
	}

	// 全零情况快速返回
	if start > end {
		return b[:0]
	}

	// 处理后面
	for ; end >= start && b[end] == 0x00; end-- {
	}

	return b[start : end+1]
}

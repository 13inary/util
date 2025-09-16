package util

import (
	"encoding/binary"
	"math"
)

func Float32ToBytes(f float32, order binary.ByteOrder) []byte {
	bits := math.Float32bits(f) // 转换为 uint32
	buf := make([]byte, 4)
	order.PutUint32(buf, bits)
	return buf
}

func BytesToFloat32(b []byte, order binary.ByteOrder) float32 {
	bits := order.Uint32(b)
	return math.Float32frombits(bits)
}

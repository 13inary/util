package util

import (
	"bytes"
	"encoding/json"
)

// FormatJSON 格式化json字符串
func FormatJSON(input []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, input, "", "    ")
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

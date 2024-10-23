package util

import (
	"encoding/json"
	"fmt"
)

// ConvertModel 源模型的字段会覆盖目标模型的字段。注意转换多个源到目的模型时多个源的字段相同问题。
func ConvertModel(srcModel any, dstModelPointer any) error {
	srcBytes, err := json.Marshal(srcModel)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal %+v, %s", srcModel, err.Error())
	}
	if err := json.Unmarshal(srcBytes, dstModelPointer); err != nil {
		return fmt.Errorf("failed to json.Unmarshal from %+v to %+v, %s", string(srcBytes), dstModelPointer, err.Error())
	}
	return nil
}

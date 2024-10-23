package util

import (
	"encoding/json"
	"fmt"
)

func ConvertModel(srcModel any, srcModelPointer any) error {
	srcBytes, err := json.Marshal(srcModel)
	if err != nil {
		return fmt.Errorf("failed to json.Marshal %+v, %s", srcModel, err.Error())
	}
	if err := json.Unmarshal(srcBytes, srcModelPointer); err != nil {
		return fmt.Errorf("failed to json.Unmarshal from %+v to %+v, %s", string(srcBytes), srcModelPointer, err.Error())
	}
	return nil
}

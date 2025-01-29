package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// 对于转化，存在的处理方法： .()   go:generate   copier   reflect   json   mapstructure

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

// Map2Struct 对于结构体字段名，若map中有就拿过来覆盖。
// 对于嵌套结构体的问题：调用方传入嵌套的字段。
func Map2Struct(srcMap map[string]any, dstStructPointer interface{}) error {
	// 获取结构体类型和字段
	val := reflect.ValueOf(dstStructPointer).Elem()
	typ := val.Type()

	// 遍历结构体的每个字段
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// 检查 map 中是否有对应的键
		if value, ok := srcMap[fieldName]; ok {
			// 设置结构体字段的值
			fieldVal := val.Field(i)
			if fieldVal.CanSet() {
				// 先将value转换为字段的类型
				convertedValue := reflect.ValueOf(value).Convert(fieldVal.Type())
				// fieldVal.Set(reflect.ValueOf(value))
				fieldVal.Set(convertedValue)
			}
		}
	}
	return nil
}

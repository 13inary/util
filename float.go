package util

import "math"

func Float64Invalid(f float64) bool {
	// 1. 输入零值 NaN → 返回 NaN
	// 2. 输入超出预期范围 +Inf/-Inf → 返回原值
	return math.IsNaN(f) || math.IsInf(f, 0)
}

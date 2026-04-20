package util

import "sort"

// GroupSort 对切片进行分组排序处理，没必要排序的元素会保持原顺序
// 参数:
//
//	src - 需要处理的源切片
//	getGroupKey - 获取元素分组键的函数，返回字符串类型的分组标识
//	sortGroup - 组间排序函数，参数是两个组的元素切片，返回true表示第一个组应排在第二个组前，为nil时不排序
//	sortMember - 组内排序函数，用于对同一组的元素进行排序，返回true表示第一个元素应排在第二个元素前，为nil时不排序
//
// 返回值:
//
//	[]T - 返回按分组排序后的新切片
//
// 注意:
//
//	排序函数返回false都会触发元素更换位置，也就是大小比较用好等于号
//	sortGroup有2个slice传进来，它们的长度>=1
func GroupSort[T any](src []T, getGroupKey func(T) string, sortGroup func([]T, []T) bool, sortMember func(T, T) bool) []T {
	groupKeys := make([]string, 0, len(src))

	// 分类成组
	groupMap := make(map[string][]T, len(src))
	for _, one := range src {
		key := getGroupKey(one)
		group, ok := groupMap[key]
		if !ok {
			groupKeys = append(groupKeys, key) // 记录下所有组的key，写在这里是为了去重
			group = make([]T, 0, len(src))
		}
		group = append(group, one)
		groupMap[key] = group
	}

	// 组内排序
	if sortMember != nil {
		for _, v := range groupMap {
			sort.Slice(v, func(i, j int) bool {
				return sortMember(v[i], v[j])
			})
		}
	}

	// 组间排序
	if sortGroup != nil {
		sort.Slice(groupKeys, func(i, j int) bool {
			return sortGroup(groupMap[groupKeys[i]], groupMap[groupKeys[j]])
		})
	}

	// 返回结果
	dest := make([]T, 0, len(src))
	for _, v := range groupKeys {
		dest = append(dest, groupMap[v]...)
	}
	return dest
}

// GetReverseIndex 获取当前数组下标对应的元素在数据反转后的新下标
//
// 1.单数数组场景：
// 元素 a b c
// 下标 0 1 2
// 反转 c b a
// 位置 2 1 0
// 2.双数数组场景：
// 元素 a b c d
// 下标 0 1 2 3
// 反转 d c b a
// 位置 3 2 1 0
func GetReverseIndex(arrayLength int, currentIndex int) int {
	return arrayLength - 1 - currentIndex
}

// 将slice切成多个部分，每个部分限制最大数量为limit
func SplitSliceByMaxLength[T any](srcSlice []T, maxLength int) [][]T {
	sliceLen := len(srcSlice)
	if sliceLen == 0 {
		return [][]T{}
	}
	ints := make([][]T, 0, sliceLen/maxLength+1)
	for start := 0; start < sliceLen; start += maxLength {
		// 找到part的右边界
		end := start + maxLength // 获取的到达end-1为止
		end = min(end, sliceLen) // 避免越界
		// 保存part数据
		ints = append(ints, srcSlice[start:end])
	}
	return ints
}

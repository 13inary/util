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

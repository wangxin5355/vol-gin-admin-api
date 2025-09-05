package utils

import "strconv"

// StringSliceToIntSliceFilter 过滤掉无法转换的字符串
func StringSliceToIntSliceFilter(strings []string) []int {
	var result []int
	for _, s := range strings {
		if num, err := strconv.Atoi(s); err == nil {
			result = append(result, num)
		}
	}
	return result
}

func FindCommonElementsSimple(a, b []string) []string {
	var common []string
	for _, itemA := range a {
		for _, itemB := range b {
			if itemA == itemB {
				common = append(common, itemA)
				break // 找到后跳出内层循环
			}
		}
	}
	return common
}

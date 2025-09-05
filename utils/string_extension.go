package utils

import (
	"strconv"
	"strings"
)

func GenerateRandomNumber(length int) string {
	var sb strings.Builder
	k := len(str)
	for i := 0; i < length; i++ {
		sb.WriteByte(str[RandomInt(0, k)])
	}
	return sb.String()
}

// 定义一个固定的数据，里边有数字和字母
var str = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 将字符串按某种char 分隔
func ToListStrString(liststr []string, splitChar string) string {
	return strings.Join(liststr, splitChar)
}

// 将int数组。按某种char 分隔，转string
func ToListIntString(listint []int, splitChar string) string {
	if len(listint) == 0 {
		return ""
	}
	// 先将int转换为string切片
	strSlice := make([]string, len(listint))
	for i, num := range listint {
		strSlice[i] = strconv.Itoa(num)
	}
	// 使用strings.Join连接
	return strings.Join(strSlice, splitChar)
}

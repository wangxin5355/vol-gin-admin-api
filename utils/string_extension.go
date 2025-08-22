package utils

import (
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

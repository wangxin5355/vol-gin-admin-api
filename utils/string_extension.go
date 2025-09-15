package utils

import (
	"strconv"
	"strings"
	"unicode"
)

// GenerateRandomNumber 生成指定长度的随机字符串
func GenerateRandomNumber(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	k := len(str)
	for i := 0; i < length; i++ {
		sb.WriteByte(str[RandomInt(0, k)])
	}
	return sb.String()
}

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

func FirstLetterUpper(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// FirstLetterLower 首字母小写
func FirstLetterLower(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

// CamelCase 驼峰命名法
func CamelCase(s string) string {
	if strings.Contains(s, "_") {
		parts := strings.Split(s, "_")
		for i := range parts {
			if len(parts[i]) > 0 {
				parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
			}
		}
		return strings.Join(parts, "")
	}
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CamelCaseLower 首字母小写的驼峰命名法
func CamelCaseLower(str string) string {
	parts := strings.Split(str, "_")
	for i := range parts {
		if i == 0 {
			parts[i] = FirstLetterLower(parts[i])
		} else {
			parts[i] = FirstLetterUpper(parts[i])
		}
	}
	return strings.Join(parts, "")
}

// GoTypeWithNull 根据字段类型和是否可空返回Go类型
func GoTypeWithNull(columnType string, isNull int) string {
	switch columnType {
	case "int32":
		if isNull == 1 {
			return "*int32"
		}
		return "int32"
	case "int64":
		if isNull == 1 {
			return "*int64"
		}
		return "int64"
	case "float64":
		if isNull == 1 {
			return "*float64"
		}
		return "float64"
	case "bool":
		if isNull == 1 {
			return "*bool"
		}
		return "bool"
	case "string":
		if isNull == 1 {
			return "*string"
		}
		return "string"
	case "time.Time":
		if isNull == 1 {
			return "*time.Time"
		}
		return "time.Time"
	default:
		return columnType
	}
}

// 检查字符串是否为空或空白字符
func IsNull(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}

package common

import (
	"fmt"
	"time"
)

// JsonTime 用于 API 输出时格式化时间（去掉 T）
type JsonTime time.Time

// MarshalJSON —— 控制序列化为 "yyyy-MM-dd HH:mm:ss"
func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

// UnmarshalJSON —— 支持反序列化前端传入的字符串
func (t *JsonTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if str == "null" {
		return nil
	}
	str = str[1 : len(str)-1] // 去掉引号
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		return err
	}
	*t = JsonTime(parsed)
	return nil
}

// ToTime —— 转回原生 time.Time
func (t JsonTime) ToTime() time.Time {
	return time.Time(t)
}

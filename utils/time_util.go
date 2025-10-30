package utils

import (
	"fmt"
	"time"
)

func FormatTimeMillis(t time.Time) string {
	return t.Format("20060102150405") + fmt.Sprintf("%03d", t.Nanosecond()/1e6)
}

type JsonTime time.Time

func (t JsonTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(formatted), nil
}

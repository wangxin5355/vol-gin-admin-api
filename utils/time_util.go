package utils

import (
	"fmt"
	"time"
)

func FormatTimeMillis(t time.Time) string {
	return t.Format("20060102150405") + fmt.Sprintf("%03d", t.Nanosecond()/1e6)
}

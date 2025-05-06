package utils

import (
	"fmt"
	"time"
)

const timeLayout = "15:04:05.000"

func ParseTime(value string) (time.Time, error) {
	return time.Parse(timeLayout, value)
}

func FormatTime(t time.Time) string {
	return t.Format(timeLayout)
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	milliseconds := int(d.Milliseconds()) % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, milliseconds)
}

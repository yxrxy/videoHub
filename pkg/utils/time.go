package utils

import "time"

// FormatTimestamp 将时间戳格式化为 ISO8601 字符串
func FormatTimestamp(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(timestamp, 0).Format(time.RFC3339)
}

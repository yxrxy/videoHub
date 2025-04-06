package constants

import "time"

const (
	// 文件权限
	DirPermission  = 0o755
	FilePermission = 0o644

	// WebSocket 相关
	WebSocketPingRatio = 9.0 / 10.0

	// 缓存时间
	DayInHours   = 24
	MonthInDays  = 30
	TokenExpiry  = time.Hour * DayInHours
	TokenExpiry2 = time.Hour * DayInHours * MonthInDays

	// 默认分页
	DefaultPageSize = 10

	// Kafka 相关
	DefaultKafkaPartitions = 50

	// 视频评分权重
	VideoScoreWeight = 1.5
)

package mq

import (
	"sync/atomic"

	"github.com/yxrxy/videoHub/app/video/domain/repository"
	"github.com/yxrxy/videoHub/pkg/kafka"
)

type VideoMQ struct {
	client *kafka.Kafka
	done   atomic.Bool
}

func NewVideoMQ(client *kafka.Kafka) repository.VideoMQ {
	return &VideoMQ{client: client}
}

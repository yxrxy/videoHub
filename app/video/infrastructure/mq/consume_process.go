package mq

import (
	"context"

	"github.com/yxrxy/videoHub/pkg/kafka"
)

const (
	VideoProcessConsumerNum = 10
	VideoProcessGroupID     = "video_process"
	DefaultConsumerChanCap  = 100
)

func (c *VideoMQ) ConsumeProcessVideo(ctx context.Context) <-chan *kafka.Message {
	msgCh := c.client.Consume(ctx,
		VideoTopic,
		VideoProcessConsumerNum,
		VideoProcessGroupID,
		DefaultConsumerChanCap)
	return msgCh
}

package mq

import (
	"context"
	"fmt"
	"strconv"

	"github.com/bytedance/sonic"
	"github.com/yxrxy/videoHub/app/video/domain/model"
	"github.com/yxrxy/videoHub/pkg/kafka"
)

func (c *VideoMQ) SendProcessVideo(ctx context.Context, videoID int64, videoPath string) error {
	msgv := &model.ProcessVideoMsg{
		VideoID:   videoID,
		VideoPath: videoPath,
	}

	v, err := sonic.Marshal(msgv)
	if err != nil {
		return fmt.Errorf("sonic.Marshal: %w", err)
	}
	msg := []*kafka.Message{
		{
			K: []byte(strconv.FormatInt(videoID%VideoProcessConsumerNum, 10)),
			V: v,
		},
	}
	if err = c.send(ctx, msg); err != nil {
		return fmt.Errorf("mq.SendAddGoods: marshal msg failed, err: %w", err)
	}
	return nil
}

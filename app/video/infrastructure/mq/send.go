package mq

import (
	"context"
	"fmt"
	"strings"

	"github.com/yxrxy/videoHub/pkg/kafka"
)

var (
	VideoTopic = "video"
)

func (c *VideoMQ) send(ctx context.Context, msg []*kafka.Message) (err error) {
	if !c.done.Load() {
		err = c.client.SetWriter(VideoTopic, true)
		if err != nil {
			return err
		}
		c.done.Swap(true)
	}
	errs := c.client.Send(ctx, VideoTopic, msg)
	if len(errs) != 0 {
		var errMsg string
		for _, e := range errs {
			errMsg = strings.Join([]string{errMsg, e.Error(), ";"}, "")
		}
		err = fmt.Errorf("mq.Send: send msg failed, errs: %v", errMsg)
		return err
	}
	return nil
}

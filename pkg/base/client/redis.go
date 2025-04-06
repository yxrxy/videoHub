/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package client

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/pkg/errno"
)

func NewRedisClient(db int) (*redis.Client, error) {
	client := config.GetRedisClient()
	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		log.Printf("redis ping failed: %v", err)
		return nil, errno.NewErrNo(errno.InternalRedisErrorCode, fmt.Sprintf("redis ping failed: %v", err))
	}
	return client, nil
}

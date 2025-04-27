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

package context

import (
	"context"
	"strconv"

	"github.com/yxrxy/videoHub/pkg/errno"
)

const UserIDKey = "user_id"

// WithUserID 将UserID加入到context中，通过metainfo传递到RPC server
func WithUserID(ctx context.Context, uid int64) context.Context {
	return newContext(ctx, UserIDKey, strconv.FormatInt(uid, 10))
}

// GetUserID 从context中取出UserID
func GetUserID(ctx context.Context) (int64, error) {
	user, ok := fromContext(ctx, UserIDKey)
	if !ok {
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get user_id from context")
	}

	value, err := strconv.ParseInt(user, 10, 64)
	if err != nil {
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to parse user_id from context")
	}
	return value, nil
}

// GetStreamUserID 流式传输传递ctx, 获取user_id
func GetStreamUserID(ctx context.Context) (int64, error) {
	uid, success := streamFromContext(ctx, UserIDKey)
	if !success {
		return -1, errno.NewErrNo(errno.ParamMissingErrorCode, "Failed to get user_id from context")
	}

	value, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return -1, errno.NewErrNo(errno.InternalServiceErrorCode, "Failed to parse user_id from context")
	}
	return value, nil
}

// SetStreamUserID 流式传输传递ctx, 设置ctx值
func SetStreamUserID(ctx context.Context, uid int64) context.Context {
	value := strconv.FormatInt(uid, 10)
	return streamAppendContext(ctx, UserIDKey, value)
}

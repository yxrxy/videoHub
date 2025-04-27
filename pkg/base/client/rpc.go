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
	"errors"
	"fmt"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/yxrxy/videoHub/config"
	"github.com/yxrxy/videoHub/kitex_gen/interaction/interactionservice"
	"github.com/yxrxy/videoHub/kitex_gen/social/socialservice"
	"github.com/yxrxy/videoHub/kitex_gen/user/userservice"
	"github.com/yxrxy/videoHub/kitex_gen/video/videoservice"
)

const (
	MuxConnection                 = 1
	KitexClientEndpointInfoFormat = "videohub.%s.client"
	UserServiceName               = "UserService"
	VideoServiceName              = "VideoService"
	InteractionServiceName        = "InteractionService"
	SocialServiceName             = "SocialService"
)

// 通用的RPC客户端初始化函数
func initRPCClient[T any](serviceName string, newClientFunc func(string, ...client.Option) (T, error)) (*T, error) {
	if config.Etcd.Addr == "" {
		return nil, errors.New("config.Etcd.Addr is empty")
	}

	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		return nil, fmt.Errorf("initRPCClient etcd.NewEtcdResolver failed: %w", err)
	}

	client, err := newClientFunc(serviceName,
		client.WithResolver(r),
		client.WithMuxConnection(MuxConnection),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: fmt.Sprintf(KitexClientEndpointInfoFormat, serviceName)}),
	)
	if err != nil {
		return nil, fmt.Errorf("initRPCClient NewClient failed: %w", err)
	}
	return &client, nil
}

func InitUserRPC() (*userservice.Client, error) {
	return initRPCClient(UserServiceName, userservice.NewClient)
}

func InitVideoRPC() (*videoservice.Client, error) {
	return initRPCClient(VideoServiceName, videoservice.NewClient)
}

func InitInteractionRPC() (*interactionservice.Client, error) {
	return initRPCClient(InteractionServiceName, interactionservice.NewClient)
}

func InitSocialRPC() (*socialservice.Client, error) {
	return initRPCClient(SocialServiceName, socialservice.NewClient)
}

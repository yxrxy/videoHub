// Code generated by Kitex v0.13.1. DO NOT EDIT.
package videoservice

import (
	server "github.com/cloudwego/kitex/server"
	video "github.com/yxrxy/videoHub/kitex_gen/video"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler video.VideoService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)
	options = append(options, server.WithCompatibleMiddlewareForUnary())

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler video.VideoService, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}

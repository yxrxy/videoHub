package main

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/yxrrxy/videoHub/app/user/router"
	"github.com/yxrrxy/videoHub/config"
	"github.com/yxrrxy/videoHub/kitex_gen/user/userservice"
)

func main() {
	config.Init()

	h := server.Default(server.WithHostPorts(config.User.HTTPAddr))

	c, err := userservice.NewClient(
		"user",
		client.WithHostPorts("127.0.0.1"+config.User.RPCAddr),
	)
	if err != nil {
		panic(err)
	}

	h.Static("/static/uploads", config.Upload.Avatar.UploadDir)

	router.Register(h, c)

	h.Spin()
}

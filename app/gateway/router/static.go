package router

import (
	"log"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func RegisterStaticRoutes(h *server.Hertz) {
	dirs := []string{
		"src/storage/videos",
		"src/storage/avatars",
		"src/storage/chat",
		"src/storage/covers",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("创建目录失败 %s: %v", dir, err)
		}
	}

	h.StaticFS("/", &app.FS{
		Root: "src/storage",
	})
}

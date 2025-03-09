package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/yxrrxy/videoHub/config"
)

func main() {
	config.Init()

	if len(os.Args) < 3 {
		fmt.Println("Usage: config_tool [service] [field]")
		os.Exit(1)
	}

	service := os.Args[1]
	field := os.Args[2]

	switch service {
	case "mysql":
		switch strings.ToLower(field) {
		case "password":
			fmt.Print(config.MySQL.Password)
		case "database":
			fmt.Print(config.MySQL.Database)
		case "port":
			fmt.Print(config.MySQL.Port)
		}
	case "user":
		switch strings.ToLower(field) {
		case "port":
			fmt.Print(strings.TrimPrefix(config.User.HTTPAddr, ":"))
		case "name":
			fmt.Print(config.User.Name)
		}
	}
}

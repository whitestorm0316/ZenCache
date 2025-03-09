package main

import (
	"log"
	"zencache/internal/config"
	"zencache/internal/transport/http"
)

func main() {
	// 加载配置文件
	conf, err := config.LoadConfig("config.json")
	if err != nil {
		log.Printf("加载配置文件失败，使用默认配置: %v", err)
		conf = &config.DefaultConfig
	}

	// 使用配置初始化服务器
	s := http.NewWithConfig(conf)
	s.Run()
}

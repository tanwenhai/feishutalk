package main

import (
	"github.com/tanwenhai/feishutalk/config"
	"github.com/tanwenhai/feishutalk/logger"
	"github.com/tanwenhai/feishutalk/web"
)

func main() {
	logger.Init()
	config.Init()
	web.Run()
}

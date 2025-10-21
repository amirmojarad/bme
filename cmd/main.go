package main

import (
	"bme/cmd/server"
	"bme/conf"
	"bme/pkg/logger"
)

func main() {
	cfg, err := conf.New()
	if err != nil {
		logger.GetLogger().Fatal(err)
	}

	if err = server.Run(cfg); err != nil {
		logger.GetLogger().Fatal(err)
	}
}

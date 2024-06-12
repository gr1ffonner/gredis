package main

import (
	"fmt"
	"gredis/internal/app"
	"gredis/internal/cache"
	"gredis/internal/config"
	"gredis/internal/db"
	"gredis/pkg/logging"
)

func main() {
	// Setup logging
	logger := logging.GetLogger("trace")
	logger.Info("logger is working")
	// Read config
	cfg := config.GetConfig()
	logger.Info("config is working")
	// Setup db
	db, err := db.NewDB(*cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	rc, err := cache.NewRedicClient(logger)
	if err != nil {
		logger.Fatal(err)
	}
	rc.Set("1", "1")
	v, err := rc.Get("1")
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(v)

	app.StartApp(*cfg, logger, db, rc)
}

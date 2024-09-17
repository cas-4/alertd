package main

import (
	"github.com/charmbracelet/log"

	"github.com/cas-4/alertd"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.AllowAll())

	var err error

	// Read environment variables and stops execution if any errors occur
	err = alertd.LoadConfig()
	if err != nil {
		log.Printf("failed to load config. err %v", err)

		return
	}

	// Ignore error because if it failed on loading, it should raised an error above.
	cfg, _ := alertd.GetConfig()

	if !cfg.Bool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	if err := alertd.RedisConnect(cfg.String("redis")); err != nil {
		log.Fatalf("Can't connect to Redis with `%s`", cfg.String("redis"))
	}

	router.POST("/alerts/", alertd.NewAlert)

	router.Run(cfg.String("address"))
}

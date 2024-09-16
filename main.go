package main

import (
	"log"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.AllowAll())

	var err error

	// Read environment variables and stops execution if any errors occur
	err = LoadConfig()
	if err != nil {
		log.Printf("failed to load config. err %v", err)

		return
	}

	// Ignore error because if it failed on loading, it should raised an error above.
	cfg, _ := GetConfig()

	if !cfg.Bool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	router.POST("/alerts/", NewAlert)

	router.Run(cfg.String("address"))
}

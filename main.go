package main

import (
	"log"

	"github.com/cas-4/alertd/config"
	"github.com/cas-4/alertd/handlers"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	router := gin.Default()

	router.Use(cors.AllowAll())

	var err error

	// Read environment variables and stops execution if any errors occur
	err = config.LoadConfig()
	if err != nil {
		log.Printf("failed to load config. err %v", err)

		return
	}

	// Ignore error because if it failed on loading, it should raised an error above.
	conf, _ := config.GetConfig()

	if !conf.Bool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	router.POST("/alerts/", handlers.NewAlert)

	router.Run(conf.String("address"))
}

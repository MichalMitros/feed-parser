package main

import (
	"time"

	"github.com/MichalMitros/feed-parser/controllers"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	// Create gin server
	r := gin.New()

	// Set logger
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Add routes and controllers
	r.POST("/parse-feed", controllers.PostParseFeed)

	// Run server
	r.Run()
}

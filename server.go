package main

import (
	"os"
	"strings"
	"time"

	"github.com/MichalMitros/feed-parser/controllers"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	envMode := os.Getenv("ENV")

	// Set logger
	var logger *zap.Logger

	// Set mode of the application (logger and gin server)
	if strings.ToLower(envMode) == "developmnent" {
		gin.SetMode(gin.DebugMode)
		logger, _ = zap.NewProduction()
		logger.Info("Running in DEVELOPMENT mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		logger, _ = zap.NewDevelopment()
		logger.Info("Running in PRODUCTION mode")
	}
	defer logger.Sync()

	// Create gin server
	r := gin.New()

	// Use logger in gin server
	zap.ReplaceGlobals(logger)
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Add routes and controllers
	r.POST("/parse-feed", controllers.PostParseFeed)

	logger.Info("Listening at :8080")

	// Run server
	err := r.Run()
	if err != nil {
		logger.Panic(
			"Couldn't start the server",
			zap.Error(err),
		)
	}
}

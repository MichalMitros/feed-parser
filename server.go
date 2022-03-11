package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MichalMitros/feed-parser/controllers"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	envMode := os.Getenv("ENV")

	// Set logger
	var logger *zap.Logger

	// Set mode of the application (logger and gin server)
	if strings.ToLower(envMode) == "development" {
		gin.SetMode(gin.DebugMode)
		logger, _ = zap.NewDevelopment()
		logger.Info("Running in DEVELOPMENT mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		logger, _ = zap.NewProduction()
		logger.Info("Running in PRODUCTION mode")
	}
	defer logger.Sync()

	serverAddress, isServerAddrSet := os.LookupEnv("SERVER_ADDRESS")
	if !isServerAddrSet {
		logger.Error(
			"'SERVER_ADDRESS' variable not set, starting on default ':8080'",
		)
		serverAddress = ":8080"
	}

	// Create gin server
	r := gin.New()

	// Use logger in gin server
	zap.ReplaceGlobals(logger)
	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger, true))

	// Add routes and controllers
	r.POST("/parse-feed", controllers.PostParseFeed)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run server
	logger.Info(
		fmt.Sprintf("Listening and serving HTTP on %s", serverAddress),
	)
	err := r.Run(serverAddress)
	if err != nil {
		logger.Panic(
			"Couldn't start the server",
			zap.Error(err),
		)
	}
}

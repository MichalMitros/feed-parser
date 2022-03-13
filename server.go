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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feedparser_requests_total",
		Help: "The total number of received requests",
	})
)

func main() {
	envMode := os.Getenv("ENV")

	// Set logger
	var logger *zap.Logger

	// Set mode of the application (logger and gin server)
	if strings.ToLower(envMode) == "development" {
		gin.SetMode(gin.DebugMode)
		logger, _ = zap.NewDevelopment()
		zap.L().Info("Running in DEVELOPMENT mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
		logger, _ = zap.NewProduction()
		zap.L().Info("Running in PRODUCTION mode")
	}
	defer logger.Sync()

	serverAddress, isServerAddrSet := os.LookupEnv("SERVER_ADDRESS")
	if !isServerAddrSet {
		zap.L().Error(
			"'SERVER_ADDRESS' variable not set, starting on default ':8080'",
		)
		serverAddress = ":8080"
	}

	// Create gin server
	r := gin.New()

	// Use zap logger in gin server
	zap.ReplaceGlobals(logger)
	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		TimeFormat: time.RFC3339,
		UTC:        true,
		SkipPaths:  []string{"/metrics"},
	}))
	r.Use(ginzap.RecoveryWithZap(zap.L(), true))

	// Use prometheus middleware
	r.Use(promMiddleware)

	// Add routes and controllers
	r.POST("/parse-feed", controllers.PostParseFeed)
	r.POST("/parse-feed-async", controllers.PostParseFeedAsync)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Run server
	zap.L().Info(
		fmt.Sprintf("Listening and serving HTTP on %s", serverAddress),
	)
	err := r.Run(serverAddress)
	if err != nil {
		zap.L().Panic(
			"Couldn't start the server",
			zap.Error(err),
		)
	}
}

func promMiddleware(c *gin.Context) {
	if !strings.HasSuffix(c.Request.URL.Path, "/metrics") {
		// Increment prom total requests counter
		opsProcessed.Inc()
	}

	// Pass on to the next-in-chain
	c.Next()
}

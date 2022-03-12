package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MichalMitros/feed-parser/controllers/contracts"
	"github.com/MichalMitros/feed-parser/feedparser"
	"github.com/MichalMitros/feed-parser/filefetcher/httpfilefetcher"
	"github.com/MichalMitros/feed-parser/fileparser/xmlparser"
	"github.com/MichalMitros/feed-parser/queuewriter/rabbitwriter"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

// Feed parser instance
var feedParser *feedparser.FeedParser

// Initialize feedParser
func init() {
	defer zap.L().Sync()

	fetcher := httpfilefetcher.DefaultHttpFileFetcher()

	queueWriter, err := rabbitwriter.NewRabbitWriter(rabbitwriter.RabbitWriterOptions{
		Hostname: getEnvVarOrPanic("RABBITMQ_HOST"),
		Username: getEnvVarOrPanic("RABBITMQ_USER"),
		Password: getEnvVarOrPanic("RABBITMQ_PASSWORD"),
	})
	if err != nil {
		zap.L().Fatal(
			"Cannot establish connection to RabbitMQ",
			zap.Error(err),
		)
	}

	fileParser := xmlparser.NewXmlFeedParser()

	// Create FeedParser instance for controllers usage
	feedParser = feedparser.NewFeedParser(fetcher, fileParser, queueWriter)
}

func PostParseFeed(c *gin.Context) {
	defer zap.L().Sync()

	// Parse request json to object
	var request contracts.ParseFeedRequest
	if err := c.BindJSON(&request); err != nil || len(request.FeedUrls) == 0 {
		zap.L().Warn("POST /parse-feed Bad Request", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": "Request should contain field 'feedUrls' with not empty list of urls",
		})
		return
	}

	// Parse all feeds from the request
	feedParser.ParseFeedsAsync(request.FeedUrls)

	// Send response
	c.IndentedJSON(http.StatusAccepted, gin.H{
		"status": "ACCEPTED",
	})
}

// Get environment variable or panic when variable is not set
func getEnvVarOrPanic(key string) string {
	defer zap.L().Sync()

	envVar, isEnvSet := os.LookupEnv(key)
	if !isEnvSet {
		zap.L().Panic(
			fmt.Sprintf("Required '%s' environment variable is not set", key),
		)
	}
	return envVar
}

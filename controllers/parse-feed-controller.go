package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MichalMitros/feed-parser/controllers/contracts"
	"github.com/MichalMitros/feed-parser/feedparser"
	"github.com/MichalMitros/feed-parser/filefetcher"
	"github.com/MichalMitros/feed-parser/fileparser/xmlparser"
	"github.com/MichalMitros/feed-parser/rabbitwriter"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

var feedParser *feedparser.FeedParser

func init() {

	fetcher := filefetcher.NewHttpFileFetcher(
		http.DefaultClient,
	)

	queueWriter, err := rabbitwriter.NewRabbitWriter(
		getEnvVarOrPanic("RABBITMQ_USER"),
		getEnvVarOrPanic("RABBITMQ_PASSWORD"),
		getEnvVarOrPanic("RABBITMQ_HOST"),
	)

	if err != nil {
		zap.L().Error(
			"Cannot establish connection to RabbitMQ",
			zap.Error(err),
		)
	}

	fileParser := xmlparser.NewXmlFeedParser()
	feedParser = feedparser.NewFeedParser(fetcher, fileParser, queueWriter)
	feedParser.Run()
}

func PostParseFeed(c *gin.Context) {
	var request contracts.ParseFeedRequest

	if err := c.BindJSON(&request); err != nil || len(request.FeedUrls) == 0 {
		zap.L().Warn("POST /parse-feed Bad Request", zap.Error(err))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status":  "BAD_REQUEST",
			"message": "Request should contain field 'feedUrls' with not empty list of urls",
		})
		return
	}

	feedUrls := feedParser.GetFeedUrlsChannel()

	for _, url := range request.FeedUrls {
		feedUrls <- url
	}

	// feedParser.ParseFeeds(request.FeedUrls)

	c.IndentedJSON(http.StatusAccepted, gin.H{
		"status": "ACCEPTED",
	})
}

func getEnvVarOrPanic(key string) string {
	envVar, isEnvSet := os.LookupEnv(key)
	if !isEnvSet {
		zap.L().Fatal(
			fmt.Sprintf("Required '%s' environment variable is not set", key),
		)
	}
	return envVar
}

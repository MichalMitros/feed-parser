package controllers

import (
	"net/http"

	"github.com/MichalMitros/feed-parser/controllers/contracts"
	"github.com/MichalMitros/feed-parser/feedparser"
	"github.com/MichalMitros/feed-parser/filefetcher"
	"github.com/MichalMitros/feed-parser/fileparser/xmlparser"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var feedParser *feedparser.FeedParser

func init() {
	fetcher := filefetcher.NewHttpFileFetcher(
		http.DefaultClient,
	)
	fileParser := xmlparser.NewXmlFeedParser()
	feedParser = feedparser.NewFeedParser(fetcher, fileParser)
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

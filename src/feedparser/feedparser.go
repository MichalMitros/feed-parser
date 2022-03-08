package feedparser

import (
	"fmt"
	"sync"

	"github.com/MichalMitros/feed-parser/filefetcher"
	"github.com/MichalMitros/feed-parser/fileparser"
	"go.uber.org/zap"
)

type FeedParser struct {
	feedUrls   chan string
	fetcher    filefetcher.FileFetcherInterface
	fileParser fileparser.FeedFileParserInterface
}

// Creates new FeedParser instance
func NewFeedParser(
	fetcher filefetcher.FileFetcherInterface,
	fileParser fileparser.FeedFileParserInterface,
) *FeedParser {
	feedUrls := make(chan string)
	return &FeedParser{
		feedUrls:   feedUrls,
		fetcher:    fetcher,
		fileParser: fileParser,
	}
}

// Runs routine listening for files to fetch
func (p *FeedParser) Run() {
	go func() {
		for {
			url := <-p.feedUrls
			go p.ParseFeed(url)
		}
	}()
}

func (p *FeedParser) ParseFeeds(feedUrls []string) {
	var wg sync.WaitGroup
	wg.Add(len(feedUrls))

	for _, url := range feedUrls {
		go func(url string) {
			p.ParseFeed(url)
			wg.Done()
		}(url)
	}

	wg.Wait()
}

func (p FeedParser) GetFeedUrlsChannel() chan string {
	return p.feedUrls
}

func (p *FeedParser) ParseFeed(feedUrl string) error {
	defer zap.L().Sync()

	zap.L().Info(
		fmt.Sprintf("Started parsing feed from %s", feedUrl),
		zap.String("feedUrl", feedUrl),
	)

	// Fetch feed file from url
	zap.L().Info("Fetching feed file", zap.String("feedUrl", feedUrl))
	feedFile, lastModified, err := p.fetcher.FetchFile(feedUrl)
	if err != nil {
		zap.L().Error(
			"Error while fetching feed file",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return err
	}
	if len(lastModified) == 0 {
		zap.L().Warn(`Feed file has no "Last-Modified" header`, zap.String("feedUrl", feedUrl))
	}
	zap.L().Info("Feed file fetched", zap.String("feedUrl", feedUrl))

	// Parse xml to object
	shop, err := p.fileParser.ParseFile(feedFile)
	if err != nil {
		zap.L().Error(
			"Error while parsing xml file",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return err
	}

	// Clear unused file to save memory
	feedFile = nil

	fmt.Println(len(shop.ShopItems))

	return nil
}

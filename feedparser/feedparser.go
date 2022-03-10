package feedparser

import (
	"fmt"
	"sync"

	"github.com/MichalMitros/feed-parser/filefetcher"
	"github.com/MichalMitros/feed-parser/fileparser"
	"github.com/MichalMitros/feed-parser/models"
	"github.com/MichalMitros/feed-parser/rabbitwriter"
	"go.uber.org/zap"
)

type FeedParser struct {
	feedUrls    chan string
	fetcher     filefetcher.FileFetcherInterface
	fileParser  fileparser.FeedFileParserInterface
	queueWriter rabbitwriter.RabbitWriterItnerface
}

// Creates new FeedParser instance
func NewFeedParser(
	fetcher filefetcher.FileFetcherInterface,
	fileParser fileparser.FeedFileParserInterface,
	queueWriter rabbitwriter.RabbitWriterItnerface,
) *FeedParser {
	feedUrls := make(chan string)
	return &FeedParser{
		feedUrls:    feedUrls,
		fetcher:     fetcher,
		fileParser:  fileParser,
		queueWriter: queueWriter,
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
	parsedShopItems := make(chan models.ShopItem)
	go p.fileParser.ParseFile(feedFile, parsedShopItems)

	allItems := make(chan models.ShopItem, 100)
	biddingItems := make(chan models.ShopItem, 100)

	go filterItems(
		parsedShopItems,
		allItems,
		biddingItems,
	)

	p.queueWriter.WriteToQueue("shop_items", allItems)
	p.queueWriter.WriteToQueue("shop_items_bidding", biddingItems)

	return nil
}

func filterItems(
	input chan models.ShopItem,
	allItemsOutput chan models.ShopItem,
	biddingItemsOutput chan models.ShopItem,
) {
	for item := range input {
		if len(item.HeurekaCPC) > 0 {
			biddingItemsOutput <- item
		}
		allItemsOutput <- item
	}
	close(allItemsOutput)
	close(biddingItemsOutput)
}

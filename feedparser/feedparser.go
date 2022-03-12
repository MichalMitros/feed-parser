package feedparser

import (
	"fmt"

	"github.com/MichalMitros/feed-parser/filefetcher"
	"github.com/MichalMitros/feed-parser/fileparser"
	"github.com/MichalMitros/feed-parser/models"
	"github.com/MichalMitros/feed-parser/queuewriter"
	"go.uber.org/zap"
)

type FeedParser struct {
	fetcher     filefetcher.FileFetcherInterface
	fileParser  fileparser.FeedFileParserInterface
	queueWriter queuewriter.QueueWriterInterface
}

// Creates new FeedParser instance
func NewFeedParser(
	fetcher filefetcher.FileFetcherInterface,
	fileParser fileparser.FeedFileParserInterface,
	queueWriter queuewriter.QueueWriterInterface,
) *FeedParser {
	return &FeedParser{
		fetcher:     fetcher,
		fileParser:  fileParser,
		queueWriter: queueWriter,
	}
}

func (p *FeedParser) ParseFeedsAsync(feedUrls []string) {
	for _, url := range feedUrls {
		go func(url string) {
			p.ParseFeed(url)
		}(url)
	}
}

func (p *FeedParser) ParseFeed(feedUrl string) error {
	defer zap.L().Sync()

	zap.L().Info(
		fmt.Sprintf("Started parsing feed from %s", feedUrl),
		zap.String("feedUrl", feedUrl),
	)

	// Fetch feed file from url
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

	// Parse xml to object
	defer zap.L().Info("Parsing feed file", zap.String("feedUrl", feedUrl))
	parsedShopItems := make(chan models.ShopItem)
	go p.fileParser.ParseFile(feedFile, parsedShopItems)

	// crewate channels for filtered shop items
	allItems := make(chan models.ShopItem, 100)
	biddingItems := make(chan models.ShopItem, 100)

	// Filter items
	defer zap.L().Info("Filtering shop items", zap.String("feedUrl", feedUrl))
	go filterItems(
		parsedShopItems,
		allItems,
		biddingItems,
	)

	defer zap.L().Info("Publishing shop items", zap.String("feedUrl", feedUrl))
	p.queueWriter.WriteToQueue("shop_items", allItems)
	p.queueWriter.WriteToQueue("shop_items_bidding", biddingItems)

	return nil
}

func filterItems(
	input chan models.ShopItem,
	allItemsOutput chan models.ShopItem,
	biddingItemsOutput chan models.ShopItem,
) {
	// Close channels after filtering
	defer close(allItemsOutput)
	defer close(biddingItemsOutput)

	for item := range input {
		// Send items with HeurekaCPC to biddingItemsOutput
		if len(item.HeurekaCPC) > 0 {
			biddingItemsOutput <- item
		}
		// Send all items to allItemsOutput
		allItemsOutput <- item
	}
}

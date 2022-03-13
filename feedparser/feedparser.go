package feedparser

import (
	"fmt"

	"github.com/MichalMitros/feed-parser/errorscollector"
	"github.com/MichalMitros/feed-parser/filefetcher"
	"github.com/MichalMitros/feed-parser/fileparser"
	"github.com/MichalMitros/feed-parser/models"
	"github.com/MichalMitros/feed-parser/queuewriter"
	"go.uber.org/zap"
)

type FeedParser struct {
	fetcher         filefetcher.FileFetcherInterface
	fileParser      fileparser.FeedFileParserInterface
	queueWriter     queuewriter.QueueWriterInterface
	errorsCollector errorscollector.ErrorsCollectorInterface
}

// Creates new FeedParser instance
func NewFeedParser(
	fetcher filefetcher.FileFetcherInterface,
	fileParser fileparser.FeedFileParserInterface,
	queueWriter queuewriter.QueueWriterInterface,
	errorsCollector errorscollector.ErrorsCollectorInterface,
) *FeedParser {
	return &FeedParser{
		fetcher:         fetcher,
		fileParser:      fileParser,
		queueWriter:     queueWriter,
		errorsCollector: errorsCollector,
	}
}

func (p *FeedParser) ParseFeedsAsync(feedUrls []string) {
	for _, url := range feedUrls {
		go func(url string) {
			p.ParseFeed(url)
		}(url)
	}
}

func (p *FeedParser) ParseFeed(feedUrl string) {
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
		return
	}
	// Check if feed has last modified value
	if len(lastModified) == 0 {
		zap.L().Warn(`Feed file has no "Last-Modified" header`, zap.String("feedUrl", feedUrl))
	} else {
		zap.L().Info(
			fmt.Sprintf(`Feed file %s last modification: %s`, feedUrl, lastModified),
			zap.String("feedUrl", feedUrl),
		)
	}

	// Parse xml to object
	zap.L().Info("Parsing feed file", zap.String("feedUrl", feedUrl))
	parsedShopItems := make(chan models.ShopItem)
	parseErrors, err := p.errorsCollector.HandleErrors(feedUrl, "file_parsing")
	if err != nil {
		zap.L().Error(
			"Error while starting errors collector",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return
	}
	go p.fileParser.ParseFile(feedFile, parsedShopItems, parseErrors)

	// Create channels for filtered shop items
	allItems := make(chan models.ShopItem, 100)
	biddingItems := make(chan models.ShopItem, 100)

	// Filter items
	zap.L().Info("Filtering shop items", zap.String("feedUrl", feedUrl))
	go filterItems(
		parsedShopItems,
		allItems,
		biddingItems,
	)

	// Publishing shop item to the queue
	zap.L().Info("Publishing shop items", zap.String("feedUrl", feedUrl))
	allItemsPublishErrors, err := p.errorsCollector.HandleErrors(feedUrl, "shop_items_publishing")
	if err != nil {
		zap.L().Error(
			"Error while starting errors collector",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return
	}
	biddingItemsPublishErrors, err := p.errorsCollector.HandleErrors(feedUrl, "bidding_shop_items_publishing")
	if err != nil {
		zap.L().Error(
			"Error while starting errors collector",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return
	}
	p.queueWriter.WriteToQueue("shop_items", allItems, allItemsPublishErrors)
	p.queueWriter.WriteToQueue("shop_items_bidding", biddingItems, biddingItemsPublishErrors)
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

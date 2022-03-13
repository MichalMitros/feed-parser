package feedparser

import (
	"fmt"
	"io"
	"sync"

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

	var wg sync.WaitGroup

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
	logFeedLastModification(feedUrl, lastModified)

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
	wg.Add(1)
	p.parseFeedFileAsync(
		feedFile,
		parsedShopItems,
		parseErrors,
		&wg,
	)

	// Create channels for filtered shop items
	allItems := make(chan models.ShopItem)
	biddingItems := make(chan models.ShopItem)

	// Filter items
	zap.L().Info("Filtering shop items", zap.String("feedUrl", feedUrl))
	wg.Add(1)
	p.filterItemsAsync(
		parsedShopItems,
		allItems,
		biddingItems,
		&wg,
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
	wg.Add(2)
	p.writeItemsToQueueAsync("shop_items", allItems, allItemsPublishErrors, &wg)
	p.writeItemsToQueueAsync("shop_items_bidding", biddingItems, biddingItemsPublishErrors, &wg)

	wg.Wait()
}

func (p *FeedParser) filterItemsAsync(
	input chan models.ShopItem,
	allItemsOutput chan models.ShopItem,
	biddingItemsOutput chan models.ShopItem,
	wg *sync.WaitGroup,
) {
	go func(
		input chan models.ShopItem,
		allItemsOutput chan models.ShopItem,
		biddingItemsOutput chan models.ShopItem,
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		p.filterItems(input, allItemsOutput, biddingItemsOutput)
	}(input, allItemsOutput, biddingItemsOutput, wg)
}

func (p FeedParser) filterItems(
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

func (p *FeedParser) parseFeedFileAsync(
	feedFile io.ReadCloser,
	parsedShopItems chan models.ShopItem,
	parseErrors chan error,
	wg *sync.WaitGroup,
) {
	go func(
		fileParser fileparser.FeedFileParserInterface,
		feedFile io.ReadCloser,
		parsedShopItems chan models.ShopItem,
		parseErrors chan error,
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		fileParser.ParseFile(feedFile, parsedShopItems, parseErrors)
	}(p.fileParser, feedFile, parsedShopItems, parseErrors, wg)
}

func (p *FeedParser) writeItemsToQueueAsync(
	queueName string,
	shopItemsInput chan models.ShopItem,
	writingErrorsInput chan error,
	wg *sync.WaitGroup,
) {
	go func(
		queueWriter queuewriter.QueueWriterInterface,
		queueName string,
		shopItemsInput chan models.ShopItem,
		writingErrorsInput chan error,
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		queueWriter.WriteToQueue(queueName, shopItemsInput, writingErrorsInput)
	}(p.queueWriter, queueName, shopItemsInput, writingErrorsInput, wg)
}

func logFeedLastModification(feedUrl string, lastModified string) {
	// Check if feed has last modified value
	if len(lastModified) == 0 {
		zap.L().Warn(`Feed file has no "Last-Modified" header`, zap.String("feedUrl", feedUrl))
	} else {
		zap.L().Info(
			fmt.Sprintf(`Feed file %s last modification: %s`, feedUrl, lastModified),
			zap.String("feedUrl", feedUrl),
		)
	}
}

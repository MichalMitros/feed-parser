package feedparser

import (
	"fmt"
	"io"
	"sync"

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

func (p *FeedParser) ParseFeeds(feedUrls []string) {
	var wg sync.WaitGroup
	parsingStatuses := []models.FeedParsingResult{}
	for _, url := range feedUrls {
		wg.Add(1)
		go func(url string, parsingStatus []models.FeedParsingResult, wg *sync.WaitGroup) {
			defer wg.Done()
			parsingResult := models.FeedParsingResult{
				FeedUrl: url,
				Status:  models.ParsingErrors,
			}
			err := p.ParseFeed(url)
			if err == nil {
				parsingResult.Status = models.ParsedSuccessfully
			}
			parsingStatuses = append(parsingStatuses, parsingResult)
		}(url, parsingStatuses, &wg)
	}
	wg.Wait()

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

	var wg sync.WaitGroup
	asyncJobsErrors := []error{}

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
	// Check if feed has last modified value
	logFeedLastModification(feedUrl, lastModified)

	// Parse xml to object
	zap.L().Info("Parsing feed file", zap.String("feedUrl", feedUrl))
	parsedShopItems := make(chan models.ShopItem)
	if err != nil {
		zap.L().Error(
			"Error while starting errors collector",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return err
	}
	parsingError := p.parseFeedFileAsync(feedFile, parsedShopItems, &wg)
	asyncJobsErrors = append(asyncJobsErrors, parsingError)

	// Create channels for filtered shop items
	allItems := make(chan models.ShopItem)
	biddingItems := make(chan models.ShopItem)

	// Filter items
	zap.L().Info("Filtering shop items", zap.String("feedUrl", feedUrl))
	p.filterItemsAsync(
		parsedShopItems,
		allItems,
		biddingItems,
		&wg,
	)

	// Publishing shop item to the queue
	zap.L().Info("Publishing shop items", zap.String("feedUrl", feedUrl))
	allItemsQueueError := p.writeItemsToQueueAsync("shop_items", allItems, &wg)
	biddingItemsQueueError := p.writeItemsToQueueAsync("shop_items_bidding", biddingItems, &wg)
	asyncJobsErrors = append(asyncJobsErrors, allItemsQueueError)
	asyncJobsErrors = append(asyncJobsErrors, biddingItemsQueueError)

	wg.Wait()

	for _, jobError := range asyncJobsErrors {
		if err != nil {
			zap.L().Error(
				"Feed parsing error",
				zap.String("feedUrl", feedUrl),
				zap.Error(err),
			)
			return jobError
		}
	}

	zap.L().Info(
		fmt.Sprintf("Successfully finished parsing feed from %s", feedUrl),
		zap.String("feedUrl", feedUrl),
	)

	return nil
}

func (p *FeedParser) filterItemsAsync(
	input chan models.ShopItem,
	allItemsOutput chan models.ShopItem,
	biddingItemsOutput chan models.ShopItem,
	wg *sync.WaitGroup,
) {
	wg.Add(1)
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
	wg *sync.WaitGroup,
) error {
	wg.Add(1)
	var err error
	err = nil
	go func(
		fileParser fileparser.FeedFileParserInterface,
		feedFile io.ReadCloser,
		parsedShopItems chan models.ShopItem,
		err error,
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		err = fileParser.ParseFile(feedFile, parsedShopItems)
	}(p.fileParser, feedFile, parsedShopItems, err, wg)

	return err
}

func (p *FeedParser) writeItemsToQueueAsync(
	queueName string,
	shopItemsInput chan models.ShopItem,
	wg *sync.WaitGroup,
) error {
	wg.Add(1)
	var err error
	err = nil
	go func(
		queueWriter queuewriter.QueueWriterInterface,
		queueName string,
		shopItemsInput chan models.ShopItem,
		err error,
		wg *sync.WaitGroup,
	) {
		defer wg.Done()
		err = queueWriter.WriteToQueue(queueName, shopItemsInput)
	}(p.queueWriter, queueName, shopItemsInput, err, wg)
	return err
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

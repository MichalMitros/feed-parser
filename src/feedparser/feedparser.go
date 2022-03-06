package feedparser

import (
	"fmt"
	"sync"

	"github.com/MichalMitros/feed-parser/filefetcher"
	"go.uber.org/zap"
)

type FeedParser struct {
	feedUrls chan string
	fetcher  filefetcher.FileFetcher
}

func NewFeedParser(
	fetcher filefetcher.FileFetcher,
) *FeedParser {
	feedUrls := make(chan string)
	return &FeedParser{
		feedUrls: feedUrls,
		fetcher:  fetcher,
	}
}

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

	feedFile, err := p.fetcher.FetchFile(feedUrl)

	if err != nil {
		zap.L().Error(
			"Error while fetching feed file",
			zap.String("feedUrl", feedUrl),
			zap.Error(err),
		)
		return err
	}

	fmt.Println(feedFile)

	return nil
}

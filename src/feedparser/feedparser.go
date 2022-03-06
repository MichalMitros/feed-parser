package feedparser

type FeedParser struct {
	feedUrls chan string
}

func NewFeedParser() *FeedParser {
	feedUrls := make(chan string)
	return &FeedParser{
		feedUrls: feedUrls,
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

func (p FeedParser) GetFeedUrlsChannel() chan string {
	return p.feedUrls
}

func (p *FeedParser) ParseFeed(feedUrl string) {

}

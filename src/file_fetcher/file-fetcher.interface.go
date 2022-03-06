package file_fetcher

// Interface of file fetcher
type FileFetcher interface {
	FetchFiles(urls []string) ([]string, error)
}

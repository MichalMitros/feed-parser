package filefetcher

// Interface of file fetcher
type FileFetcherInterface interface {
	FetchFile(url string) ([]byte, error)
}

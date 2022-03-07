package filefetcher

// Interface of file fetcher
type FileFetcher interface {
	FetchFile(url string) ([]byte, error)
}

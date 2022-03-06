package file_fetcher

// File fetcher for fetching file from http asynchronously
// Implements FileFetcher interface
type AsyncFileFetcher struct {
	httpClient *HttpClientInterface
}

// Creates new AsyncFileFetcher instance
func NewAsyncFileFetcher(
	httpClient HttpClientInterface,
) *AsyncFileFetcher {
	return &AsyncFileFetcher{
		httpClient: &httpClient,
	}
}

func (f *AsyncFileFetcher) FetchFiles(
	urls []string,
) ([]string, error) {
	return []string{}, nil
}

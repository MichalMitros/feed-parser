package filefetcher

import "io"

// Interface of file fetcher
type FileFetcherInterface interface {
	FetchFile(url string) (io.ReadCloser, string, error)
}

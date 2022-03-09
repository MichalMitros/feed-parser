package filefetcher

import (
	"io"

	"go.uber.org/zap"
)

// File fetcher for fetching file from http asynchronously
// Implements FileFetcher interface
type HttpFileFetcher struct {
	httpClient HttpClientInterface
}

// Creates new AsyncFileFetcher instance
func NewHttpFileFetcher(
	httpClient HttpClientInterface,
) *HttpFileFetcher {
	return &HttpFileFetcher{
		httpClient: httpClient,
	}
}

func (f *HttpFileFetcher) FetchFile(
	url string,
) (io.ReadCloser, string, error) {

	resp, err := f.httpClient.Get(url)
	if err != nil {
		return nil, "", err
	}
	zap.L().Debug("Feed file HTTP headers", zap.Any("responseHeaders", resp.Header))

	lastModified := resp.Header.Get("Last-Modified")

	return resp.Body, lastModified, nil
}

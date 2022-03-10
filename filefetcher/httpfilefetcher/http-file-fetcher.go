package httpfilefetcher

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

// Fetcher for getting files from http urls
// Implements FileFetcher interface
type HttpFileFetcher struct {
	httpClient HttpClientInterface
}

// Creates new FileFetcher instance
func NewHttpFileFetcher(
	httpClient HttpClientInterface,
) *HttpFileFetcher {
	return &HttpFileFetcher{
		httpClient: httpClient,
	}
}

// Creates new FileFetcher instance with default httpClient
func DefaultHttpFileFetcher() *HttpFileFetcher {
	return &HttpFileFetcher{
		httpClient: http.DefaultClient,
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

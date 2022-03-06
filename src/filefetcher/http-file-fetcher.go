package filefetcher

import (
	"io"
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
) (string, error) {

	resp, err := f.httpClient.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

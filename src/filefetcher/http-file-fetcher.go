package filefetcher

import (
	"fmt"
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
) ([]byte, error) {

	fmt.Println("Fetching...")
	resp, err := f.httpClient.Get(url)
	fmt.Println("Fetched")
	if err != nil {
		return nil, err
	}

	fmt.Println("Reading Body...")
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	fmt.Println("Body reading complete...")

	if err != nil {
		return nil, err
	}

	return body, nil
}

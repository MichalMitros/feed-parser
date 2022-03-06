package filefetcher

import (
	"net/http"
	"testing"
)

// Mocked http.Client as struct implementing FileFetcher interface
type MockedHttpClient struct{}

func (c MockedHttpClient) Get(url string) (*http.Response, error) {
	return &http.Response{
		Status:     "SOME STATUS",
		StatusCode: 200,
		Body:       MockedHttpResponseBody{},
	}, nil
}

type MockedHttpResponseBody struct{}

func (b MockedHttpResponseBody) Read(p []byte) (int, error) {

	return 0, nil
}

func (b MockedHttpResponseBody) Close() error {
	return nil
}

func TestFetchFile(t *testing.T) {
	client := MockedHttpClient{}
	filesFetcher := NewHttpFileFetcher(
		client,
	)
	files, err := filesFetcher.FetchFile("url")

	if err != nil {
		t.Fatalf(`FetchFiles([]string{}) = %q, %v, want []string{""}, nil`, files, err)
	}

}

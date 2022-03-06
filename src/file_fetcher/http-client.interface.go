package file_fetcher

import "net/http"

// Interface for http.Client struct made for easier
// mocking and testing as there is not built-in native interface
//
// http.Client docs: https://pkg.go.dev/net/http#Client
type HttpClientInterface interface {
	Get(url string) (*http.Response, error)
}

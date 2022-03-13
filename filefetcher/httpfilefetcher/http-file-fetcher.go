package httpfilefetcher

import (
	"io"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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

// Fetch file and returns response body as io.ReadCloser,
// "Last-Modified" header as string and potentially an error
func (f *HttpFileFetcher) FetchFile(
	url string,
) (*io.ReadCloser, string, error) {
	defer zap.L().Sync()

	resp, err := f.httpClient.Get(url)
	if err != nil {
		filesFetchedFailures.Inc()
		return nil, "", err
	}
	filesFetched.Inc()
	zap.L().Debug("Feed file HTTP headers", zap.Any("responseHeaders", resp.Header))

	lastModified := resp.Header.Get("Last-Modified")

	return &resp.Body, lastModified, nil
}

// Prometheus fetched xml files counter
var (
	filesFetched = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feedparser_fetched_xml_files_total",
		Help: "The total number of fetched XML files",
	})
	filesFetchedFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feedparser_fetched_xml_files_failures_total",
		Help: "The total number of failures in fetching XML files",
	})
)

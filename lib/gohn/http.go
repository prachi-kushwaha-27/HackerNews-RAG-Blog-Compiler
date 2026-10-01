package gohn

import (
	"net/http"
	"time"
)

const (
	_workQueueMaxSize = 10_000
)

type rateControlledHTTPClient interface {
	Get(url string) <-chan httpResponse
}

func newRateControlledHTTPClient(
	timeout time.Duration,
	rps int,
) rateControlledHTTPClient {
	client := &rcHTTPClient{
		http: http.Client{
			Timeout: timeout,
		},
		rps:       rps,
		workItems: make(chan _httpQueueItem, _workQueueMaxSize),
	}
	go client._run()

	return client
}

type httpResponse struct {
	resp *http.Response
	err  error
}

type _httpQueueItem struct {
	url string
	out chan<- httpResponse
}

type rcHTTPClient struct {
	http      http.Client
	rps       int
	workItems chan _httpQueueItem
}

func (c *rcHTTPClient) Get(url string) <-chan httpResponse {
	out := make(chan httpResponse, 1)
	c.workItems <- _httpQueueItem{
		url: url,
		out: out,
	}
	return out
}

func (c *rcHTTPClient) _run() {
	for item := range c.workItems {
		resp, err := c.http.Get(item.url)
		item.out <- httpResponse{
			resp: resp,
			err:  err,
		}
		time.Sleep(time.Duration((1.0 / float64(c.rps)) * float64(time.Second)))
	}
}

func (c *rcHTTPClient) Terminate() {
	close(c.workItems)
}

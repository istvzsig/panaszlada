package internal

import (
	"context"
	"errors"
	"net/http"

	"github.com/istvzsig/retryx"
)

type Client struct {
	http    *http.Client
	wrapper retryx.Wrapper
}

func New(httpClient *http.Client, wrapper retryx.Wrapper) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		http:    httpClient,
		wrapper: wrapper,
	}
}

func (c *Client) Get(ctx context.Context, url string) (*http.Response, error) {
	var resp *http.Response

	err := c.wrapper.Do(ctx, func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			url,
			nil,
		)
		if err != nil {
			return err
		}

		r, err := c.http.Do(req)
		if err != nil {
			return err
		}

		// retryable HTTP statuses
		if r.StatusCode >= 500 || r.StatusCode == 429 {
			r.Body.Close()
			return errors.New("retryable http status")
		}

		resp = r
		return nil
	})

	return resp, err
}

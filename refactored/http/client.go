package http

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	c http.Client
}

func NewClient(c http.Client) *Client {
	return &Client{c: c}
}

func (c *Client) Get(u *url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO: Handle different status codes

	return resp.Body, nil
}

func (c *Client) DownloadDocument(u url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", u.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	// TODO: Handle different status codes

	return resp.Body, nil
}

func (c *Client) DownloadTo(ctx context.Context, w io.Writer) error {
	req, err := http.NewRequest("GET", "", http.NoBody)
	if err != nil {
		return err
	}

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

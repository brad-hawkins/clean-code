package types

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type HTTPManager interface {
	HTTPClient() (*http.Client, error)
	NewHTTPRequest(context.Context, string, url.URL, io.Reader) (*http.Request, error)
	Headers() map[string]string
}

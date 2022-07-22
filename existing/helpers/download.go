package helpers

import (
	"context"
	types2 "github.com/brad-hawkins/clean-code/existing/types"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

/*
	DownloadFile - given a URI and a full path for the downloaded file. If the file already exists it
	will be overwritten.
*/
func DownloadFile(ctx context.Context, httpClient types2.HTTPManager, uri url.URL, w io.Writer) error {

	req, err := httpClient.NewHTTPRequest(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return types2.WrapError(err, "unable to create http request")
	}

	c, err := httpClient.HTTPClient()
	if err != nil {
		return types2.WrapError(err, "unable to get http client")
	}

	resp, err := c.Do(req)
	if err != nil {
		return types2.WrapError(err, "unable to make request to get file")
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		e := types2.NewError("invalid response code received", types2.WithTag("statusCode", resp.StatusCode), types2.WithTag("url", uri.String()))
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return e
		}
		e.AddTag("body", string(body))
		return e
	}

	_, err = io.Copy(w, resp.Body)

	if err != nil {
		return types2.WrapError(err, "unable to copy response to temp file")
	}

	return nil
}

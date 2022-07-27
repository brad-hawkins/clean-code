package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/brad-hawkins/clean-code/refactored/controller"
	"github.com/brad-hawkins/clean-code/refactored/document"
	"github.com/brad-hawkins/clean-code/refactored/download"
	httpish "github.com/brad-hawkins/clean-code/refactored/http"
	"github.com/brad-hawkins/clean-code/refactored/integration"
)

func main() {

	ctx := context.Background()

	httpClient := httpish.NewClient(http.Client{
		Timeout: 30 * time.Second,
	})

	downloadFactory := download.NewDownloaderFactory(httpClient, &url.URL{}, &url.URL{})

	fileDownloader := download.NewFileDownloader(downloadFactory)
	integrationManager := integration.NewManager()

	c := controller.NewController(integrationManager, fileDownloader)

	err := c.ImportDocument(ctx, &document.Document{
		DocumentID: "123",
		CategoryID: "asdf",
		CustomerID: "98734",
	})
	if err != nil {
		log.Fatal(err)
	}
}

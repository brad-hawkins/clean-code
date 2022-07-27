package controller

import (
	"context"
	"github.com/brad-hawkins/clean-code/refactored/document"
	"github.com/brad-hawkins/clean-code/refactored/download"
	"github.com/brad-hawkins/clean-code/refactored/integration"
)

type Controller struct {
	integrationManager *integration.Manager
	fileDownloader     *download.File
}

func NewController(integrationManager *integration.Manager, fileDownloader *download.File) *Controller {
	return &Controller{
		integrationManager: integrationManager,
		fileDownloader:     fileDownloader,
	}
}

func (c *Controller) ImportDocument(ctx context.Context, doc *document.Document) error {

	documentImporter, err := c.integrationManager.GetDocumentImporter()
	if err != nil {
		return err
	}

	tempFile, err := c.fileDownloader.DownloadToTempFile(ctx, doc)
	if err != nil {
		return err
	}

	defer tempFile.Close()

	err = documentImporter.ImportDocument(ctx, doc, tempFile)
	if err != nil {
		return err
	}

	return nil

}

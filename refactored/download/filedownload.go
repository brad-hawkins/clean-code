package download

import (
	"context"
	"github.com/brad-hawkins/clean-code/refactored/document"
	"github.com/brad-hawkins/clean-code/refactored/filestore"
	"io"
)

type File struct {
	downloadFactory *DownloaderFactory
}

func NewFileDownloader(downloadFactory *DownloaderFactory) *File {
	return &File{downloadFactory: downloadFactory}
}

// DownloadToTempFile downloads the document provided to a temp file and returns an io.ReadCloser so it can be read
// and cleaned up
func (f *File) DownloadToTempFile(ctx context.Context, doc *document.Document) (io.ReadCloser, error) {
	downloader := f.downloadFactory.NewDownloader(doc)

	tempFile, err := filestore.NewTempFile("document", doc.FileName(), "pdf")
	if err != nil {
		return nil, err
	}

	err = downloader.DownloadTo(ctx, tempFile)
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}

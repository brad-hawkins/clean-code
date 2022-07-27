package download

import (
	"context"
	types2 "github.com/brad-hawkins/clean-code/existing/types"
	"github.com/brad-hawkins/clean-code/refactored/document"
	"io"
	"net/url"
)

type DownloaderTo interface {
	DownloadTo(ctx context.Context, w io.Writer) error
}

type getter interface {
	Get(u *url.URL) (io.ReadCloser, error)
}

type DownloaderFactory struct {
	getter               getter
	formsDocumentURL     *url.URL
	insuranceDocumentURL *url.URL
}

func NewDownloaderFactory(getter getter, formsDocumentURL *url.URL, insuranceDocumentURL *url.URL) *DownloaderFactory {
	return &DownloaderFactory{
		getter:               getter,
		formsDocumentURL:     formsDocumentURL,
		insuranceDocumentURL: insuranceDocumentURL,
	}
}

// NewDownloader returns a downloader that can be used to download to an io.Writer. The downloader returned depends on
// the document provided. Currently, Forms and Insurance documents are supported
func (f *DownloaderFactory) NewDownloader(doc *document.Document) DownloaderTo {
	u := f.formsDocumentURL
	if doc.CategoryID == types2.PatientInsuranceCategory {
		u = f.insuranceDocumentURL
	}
	return &Downloader{
		getter:      f.getter,
		downloadURL: u,
	}
}

type Downloader struct {
	getter      getter
	downloadURL *url.URL
}

// DownloadTo will call the configured url and write the response to the provided writer
func (d *Downloader) DownloadTo(ctx context.Context, w io.Writer) error {
	body, err := d.getter.Get(d.downloadURL)
	if err != nil {
		return err
	}

	defer body.Close()

	_, err = io.Copy(w, body)
	if err != nil {
		return err
	}

	return nil
}

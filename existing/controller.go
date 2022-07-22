package existing

import (
	"context"
	"fmt"
	helpers "github.com/brad-hawkins/clean-code/existing/helpers"
	types2 "github.com/brad-hawkins/clean-code/existing/types"
	"github.com/prometheus/common/log"
	"net/url"
	"os"
	"path/filepath"
	"weavelab.xyz/monorail/shared/wlib/werror"
)

const (
	Prod1Hostname = "sync-app.weaveconnect.com"
)

type Controller struct {
	integrationManager IntegrationManager
	http               types2.HTTPManager
	settings           types2.Settings
	hostname           string
}

type IntegrationManager struct {
	Current types2.Integration
}

func (c *Controller) ImportDocument(doc *types2.Document) error {

	if c.integrationManager.Current == nil {
		return types2.UninitializedError.Here("integration has not been initialized")
	}

	importDoc, ok := c.integrationManager.Current.(types2.DocumentImporter)
	if !ok {
		return types2.NotImplementedError.Here("integration does not support importing a document")
	}

	if c.isWritebackDisabled() {
		return types2.FunctionalityDisabled.Here("writebacks are disabled for location")
	}

	// get absolute path for the application
	aPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return types2.WrapError(err, "cannot get absolute path for the sync-app", types2.WithFileSystemError())
	}

	parentDir := filepath.Dir(aPath)

	filename := fmt.Sprintf("%s-", doc.CustomerID)

	destDir := filepath.Join(parentDir, "document")

	err = os.MkdirAll(destDir, 0600)
	if err != nil {
		return types2.WrapError(err, "unable to create document directory", types2.WithFileSystemError())
	}

	tmpFile, err := TempFile(destDir, filename, ".pdf")
	if err != nil {
		return types2.WrapError(err, "unable to create file", types2.WithFileSystemError(), types2.WithTag("path", destDir), types2.WithTag("name", filename))
	}

	// ensure document is removed so nothing is left behind
	defer func() error {
		if err := tmpFile.Close(); err != nil {
			log.Infoln(werror.Wrap(err, "unable to close document"))
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			return types2.WrapError(err, "unable to remove imported document", types2.WithFileSystemError(), types2.WithTag("fileName", tmpFile.Name()))
		}
		return nil
	}()

	// Download document from sync-app-api
	u := c.getDocumentURL(doc.DocumentID, doc.CategoryID)
	err = helpers.DownloadFile(context.Background(), c.http, u, tmpFile)
	if err != nil {
		return err
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		return types2.WrapError(err, "unable to seek to beginning of temp file", types2.WithFileSystemError(), types2.WithTag("fileName", tmpFile.Name()))
	}

	err = importDoc.ImportDocument(doc, tmpFile)
	if err != nil {
		return werror.Wrap(err, "unable to import a document").Add("documentID", doc.DocumentID).Add("personID", doc.CustomerID)
	}

	return nil
}

func (c *Controller) isWritebackDisabled() bool {
	return c.settings.GetRaw(types2.SettingWritebackDisabled) == "true"
}

func (c *Controller) getDocumentURL(id string, category string) url.URL {
	path := "sync-app-api/document/%s"
	if category == types2.PatientInsuranceCategory {
		path = "sync-app-api/insurance-document/%s"
	}
	return url.URL{
		Scheme: "https",
		Host:   c.GetWeaveHostname(),
		Path:   fmt.Sprintf(path, id),
	}
}

func (c *Controller) GetWeaveHostname() string {
	if c.hostname == "" {
		return Prod1Hostname //If configs are not loaded for some reason, default to production
	}
	return c.hostname
}

package integration

import (
	"context"
	types2 "github.com/brad-hawkins/clean-code/existing/types"
	"github.com/brad-hawkins/clean-code/refactored/document"
	"io"
)

type Manager struct {
	integration types2.Integration
	settings    types2.Settings
}

type DocumentImporter interface {
	ImportDocument(ctx context.Context, doc *document.Document, reader io.Reader) error
}

func NewManager() *Manager {
	return &Manager{
		integration: nil,
		settings:    nil,
	}
}

// GetDocumentImporter checks if the configured integration supports importing documents and is enabled.
// If it is, it will return a DocumentImporter that allows importing of documents to the PMS.
func (i *Manager) GetDocumentImporter() (DocumentImporter, error) {
	if i.isWritebackDisabled() {
		return nil, types2.FunctionalityDisabled.Here("writebacks are disabled for location")
	}

	if i.integration == nil {
		return nil, types2.UninitializedError.Here("integration has not been initialized")
	}

	importDoc, ok := i.integration.(DocumentImporter)
	if !ok {
		return nil, types2.NotImplementedError.Here("integration does not support importing a document")
	}

	return importDoc, nil

}

func (i *Manager) isWritebackDisabled() bool {
	return i.settings.GetRaw(types2.SettingWritebackDisabled) == "true"
}

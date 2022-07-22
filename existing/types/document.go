package types

import "io"

type Document struct {
	DocumentID string `json:"documentID"`
	CategoryID string `json:"categoryID"`
	CustomerID string `json:"customerID"`
}

type DocumentImporter interface {
	ImportDocument(doc *Document, file io.Reader) error
}

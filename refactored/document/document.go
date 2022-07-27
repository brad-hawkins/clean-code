package document

type Document struct {
	DocumentID string
	CategoryID string
	CustomerID string
}

func (d *Document) FileName() string {
	return d.CustomerID
}

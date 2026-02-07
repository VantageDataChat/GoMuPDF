//go:build !cgo || nomupdf

package gomupdf

// This file provides stub types and functions when CGO is not available.
// The full implementation requires MuPDF C library via CGO.

// Document represents a document (PDF, XPS, EPUB, etc.).
type Document struct {
	name     string
	isClosed bool
}

// Page represents a document page.
type Page struct {
	doc    *Document
	number int
}

// Pixmap represents a pixel map (raster image).
type Pixmap struct{}

// TextPage represents extracted text from a page.
type TextPage struct{}

// Annot represents a PDF annotation.
type Annot struct{}

// Widget represents a PDF form field widget.
type Widget struct{}

// Open opens a document (stub - requires CGO with MuPDF).
func Open(filename string) (*Document, error) {
	return nil, ErrInitFailed
}

// OpenFromMemory opens a document from memory (stub - requires CGO).
func OpenFromMemory(data []byte, magic string) (*Document, error) {
	return nil, ErrInitFailed
}

// NewPDF creates a new empty PDF (stub - requires CGO).
func NewPDF() (*Document, error) {
	return nil, ErrInitFailed
}

// NewPixmap creates a new pixmap (stub - requires CGO).
func NewPixmap(colorspace int, width, height int, alpha bool) (*Pixmap, error) {
	return nil, ErrInitFailed
}

// NewPixmapFromImage creates a pixmap from a PDF image (stub - requires CGO).
func NewPixmapFromImage(doc *Document, xref int) (*Pixmap, error) {
	return nil, ErrInitFailed
}

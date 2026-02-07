//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Document represents a document (PDF, XPS, EPUB, etc.).
type Document struct {
	ctx      *context
	doc      *C.fz_document
	pdf      *C.pdf_document
	name     string
	isClosed bool
}

// Open opens a document from a file path.
func Open(filename string) (*Document, error) {
	ctx, err := newContext()
	if err != nil {
		return nil, err
	}
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var errcode C.int
	doc := C.gomupdf_open_document(ctx.ctx, cFilename, &errcode)
	if errcode != 0 || doc == nil {
		ctx.close()
		return nil, fmt.Errorf("%w: %s", ErrOpenFailed, filename)
	}
	d := &Document{ctx: ctx, doc: doc, name: filename}
	d.pdf = C.gomupdf_pdf_document(ctx.ctx, doc)
	return d, nil
}

// OpenFromMemory opens a document from a byte slice.
func OpenFromMemory(data []byte, magic string) (*Document, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("%w: empty data", ErrOpenFailed)
	}
	ctx, err := newContext()
	if err != nil {
		return nil, err
	}
	cMagic := C.CString(magic)
	defer C.free(unsafe.Pointer(cMagic))

	var errcode C.int
	doc := C.gomupdf_open_document_from_memory(ctx.ctx, cMagic,
		(*C.uchar)(unsafe.Pointer(&data[0])), C.int(len(data)), &errcode)
	if errcode != 0 || doc == nil {
		ctx.close()
		return nil, ErrOpenFailed
	}
	d := &Document{ctx: ctx, doc: doc, name: magic}
	d.pdf = C.gomupdf_pdf_document(ctx.ctx, doc)
	return d, nil
}

// NewPDF creates a new empty PDF document.
func NewPDF() (*Document, error) {
	ctx, err := newContext()
	if err != nil {
		return nil, err
	}
	pdfDoc := C.pdf_create_document(ctx.ctx)
	if pdfDoc == nil {
		ctx.close()
		return nil, ErrInitFailed
	}
	doc := C.gomupdf_pdf_to_fz_document(pdfDoc)
	d := &Document{ctx: ctx, doc: doc, pdf: pdfDoc, name: "application/pdf"}
	return d, nil
}

func (d *Document) Close() {
	if d.isClosed {
		return
	}
	d.isClosed = true
	if d.doc != nil {
		C.gomupdf_drop_document(d.ctx.ctx, d.doc)
		d.doc = nil
		d.pdf = nil
	}
	if d.ctx != nil {
		d.ctx.close()
		d.ctx = nil
	}
}

func (d *Document) IsClosed() bool { return d.isClosed }
func (d *Document) Name() string   { return d.name }
func (d *Document) IsPDF() bool    { return d.pdf != nil }

func (d *Document) PageCount() int {
	if d.isClosed {
		return 0
	}
	return int(C.gomupdf_page_count(d.ctx.ctx, d.doc))
}

func (d *Document) NeedsPass() bool {
	if d.isClosed {
		return false
	}
	return C.gomupdf_needs_password(d.ctx.ctx, d.doc) != 0
}

func (d *Document) Authenticate(password string) (int, error) {
	if d.isClosed {
		return 0, ErrClosed
	}
	cPw := C.CString(password)
	defer C.free(unsafe.Pointer(cPw))
	result := int(C.gomupdf_authenticate_password(d.ctx.ctx, d.doc, cPw))
	if result == 0 {
		return 0, ErrAuthFailed
	}
	return result, nil
}

func (d *Document) IsReflowable() bool {
	if d.isClosed {
		return false
	}
	return C.gomupdf_is_document_reflowable(d.ctx.ctx, d.doc) != 0
}

func (d *Document) Layout(width, height, fontsize float64) {
	if d.isClosed || !d.IsReflowable() {
		return
	}
	C.gomupdf_layout_document(d.ctx.ctx, d.doc, C.float(width), C.float(height), C.float(fontsize))
}

func (d *Document) LoadPage(pageNum int) (*Page, error) {
	if d.isClosed {
		return nil, ErrClosed
	}
	count := d.PageCount()
	if pageNum < 0 {
		pageNum += count
	}
	if pageNum < 0 || pageNum >= count {
		return nil, fmt.Errorf("%w: %d (document has %d pages)", ErrPageNotFound, pageNum, count)
	}
	var errcode C.int
	page := C.gomupdf_load_page(d.ctx.ctx, d.doc, C.int(pageNum), &errcode)
	if errcode != 0 || page == nil {
		return nil, fmt.Errorf("%w: page %d", ErrPageNotFound, pageNum)
	}
	return &Page{ctx: d.ctx, page: page, doc: d, number: pageNum}, nil
}

func (d *Document) Pages(args ...int) ([]*Page, error) {
	if d.isClosed {
		return nil, ErrClosed
	}
	count := d.PageCount()
	start, stop, step := 0, count, 1
	switch len(args) {
	case 3:
		step = args[2]
		fallthrough
	case 2:
		stop = args[1]
		fallthrough
	case 1:
		start = args[0]
	}
	if start < 0 {
		start += count
	}
	var pages []*Page
	for i := start; (step > 0 && i < stop) || (step < 0 && i > stop); i += step {
		pno := i
		if pno < 0 {
			pno += count
		}
		pno = pno % count
		if pno < 0 {
			pno += count
		}
		page, err := d.LoadPage(pno)
		if err != nil {
			return pages, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (d *Document) Metadata() map[string]string {
	if d.isClosed {
		return nil
	}
	meta := make(map[string]string)
	keys := map[string]string{
		"format":            "format",
		"encryption":        "encryption",
		"info:Title":        "title",
		"info:Author":       "author",
		"info:Subject":      "subject",
		"info:Keywords":     "keywords",
		"info:Creator":      "creator",
		"info:Producer":     "producer",
		"info:CreationDate": "creationDate",
		"info:ModDate":      "modDate",
	}
	for cKey, goKey := range keys {
		val := d.lookupMetadata(cKey)
		if val != "" {
			meta[goKey] = val
		}
	}
	return meta
}

func (d *Document) lookupMetadata(key string) string {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var errcode C.int
	cVal := C.gomupdf_lookup_metadata(d.ctx.ctx, d.doc, cKey, &errcode)
	if errcode != 0 || cVal == nil {
		return ""
	}
	defer d.ctx.freeString(cVal)
	return C.GoString(cVal)
}

func (d *Document) GetTOC(simple bool) ([]TOCItem, error) {
	if d.isClosed {
		return nil, ErrClosed
	}
	var errcode C.int
	outline := C.gomupdf_load_outline(d.ctx.ctx, d.doc, &errcode)
	if errcode != 0 {
		return nil, ErrOutline
	}
	if outline == nil {
		return nil, nil
	}
	defer C.gomupdf_drop_outline(d.ctx.ctx, outline)
	var items []TOCItem
	d.walkOutline(outline, 1, &items)
	return items, nil
}

func (d *Document) walkOutline(outline *C.fz_outline, level int, items *[]TOCItem) {
	for ol := outline; ol != nil; ol = ol.next {
		item := TOCItem{
			Level: level,
			Title: C.GoString(ol.title),
			Page:  int(ol.page.page) + 1,
		}
		if ol.uri != nil {
			uri := C.GoString(ol.uri)
			if item.Dest == nil {
				item.Dest = &LinkDest{}
			}
			item.Dest.URI = uri
		}
		*items = append(*items, item)
		if ol.down != nil {
			d.walkOutline(ol.down, level+1, items)
		}
	}
}

// Convenience methods that delegate to Page.

func (d *Document) GetPageText(pno int, output string) (string, error) {
	page, err := d.LoadPage(pno)
	if err != nil {
		return "", err
	}
	defer page.Close()
	return page.GetText(output)
}

func (d *Document) GetPagePixmap(pno int, opts ...PixmapOption) (*Pixmap, error) {
	page, err := d.LoadPage(pno)
	if err != nil {
		return nil, err
	}
	defer page.Close()
	return page.GetPixmap(opts...)
}

func (d *Document) SearchPageFor(pno int, needle string, quads bool) ([]Quad, error) {
	page, err := d.LoadPage(pno)
	if err != nil {
		return nil, err
	}
	defer page.Close()
	return page.SearchFor(needle, quads)
}

func (d *Document) GetPageFonts(pno int) ([]FontInfo, error) {
	if !d.IsPDF() {
		return nil, ErrNotPDF
	}
	page, err := d.LoadPage(pno)
	if err != nil {
		return nil, err
	}
	defer page.Close()
	return page.GetFonts()
}

func (d *Document) GetPageImages(pno int) ([]ImageInfo, error) {
	if !d.IsPDF() {
		return nil, ErrNotPDF
	}
	page, err := d.LoadPage(pno)
	if err != nil {
		return nil, err
	}
	defer page.Close()
	return page.GetImages()
}

func (d *Document) ConvertToPDF(fromPage, toPage, rotate int) ([]byte, error) {
	if d.isClosed {
		return nil, ErrClosed
	}
	var outlen, errcode C.int
	data := C.gomupdf_convert_to_pdf(d.ctx.ctx, d.doc,
		C.int(fromPage), C.int(toPage), C.int(rotate), &outlen, &errcode)
	if errcode != 0 || data == nil {
		return nil, ErrConvert
	}
	defer d.ctx.freeBytes(data)
	return C.GoBytes(unsafe.Pointer(data), outlen), nil
}

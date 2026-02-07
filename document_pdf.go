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

func (d *Document) Save(filename string, opts ...SaveOptions) error {
	if d.isClosed {
		return ErrClosed
	}
	if !d.IsPDF() {
		return ErrNotPDF
	}
	opt := DefaultSaveOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	var cOwnerPW, cUserPW *C.char
	if opt.OwnerPW != "" {
		cOwnerPW = C.CString(opt.OwnerPW)
		defer C.free(unsafe.Pointer(cOwnerPW))
	}
	if opt.UserPW != "" {
		cUserPW = C.CString(opt.UserPW)
		defer C.free(unsafe.Pointer(cUserPW))
	}

	boolToInt := func(b bool) int { if b { return 1 }; return 0 }

	errcode := C.gomupdf_pdf_save(d.ctx.ctx, d.pdf, cFilename,
		C.int(opt.Garbage), C.int(boolToInt(opt.Deflate)), C.int(boolToInt(opt.Linear)),
		C.int(boolToInt(opt.Clean)), C.int(boolToInt(opt.ASCII)),
		C.int(boolToInt(opt.Incremental)), C.int(boolToInt(opt.Pretty)),
		C.int(opt.Encryption), cOwnerPW, cUserPW, C.int(opt.Permissions))
	if errcode != 0 {
		return fmt.Errorf("%w: %s", ErrSave, filename)
	}
	return nil
}

func (d *Document) EzSave(filename string) error {
	return d.Save(filename, EzSaveOptions())
}

func (d *Document) ToBytes(opts ...SaveOptions) ([]byte, error) {
	if d.isClosed {
		return nil, ErrClosed
	}
	if !d.IsPDF() {
		return nil, ErrNotPDF
	}
	opt := DefaultSaveOptions()
	if len(opts) > 0 {
		opt = opts[0]
	}
	boolToInt := func(b bool) int { if b { return 1 }; return 0 }

	var outlen, errcode C.int
	data := C.gomupdf_pdf_tobytes(d.ctx.ctx, d.pdf,
		C.int(opt.Garbage), C.int(boolToInt(opt.Deflate)),
		C.int(boolToInt(opt.Clean)), C.int(boolToInt(opt.ASCII)),
		C.int(boolToInt(opt.Pretty)), &outlen, &errcode)
	if errcode != 0 || data == nil {
		return nil, ErrSave
	}
	defer d.ctx.freeBytes(data)
	return C.GoBytes(unsafe.Pointer(data), outlen), nil
}

func (d *Document) NewPage(pno int, width, height float64) (*Page, error) {
	if d.isClosed {
		return nil, ErrClosed
	}
	if !d.IsPDF() {
		return nil, ErrNotPDF
	}
	if width <= 0 {
		width = 595
	}
	if height <= 0 {
		height = 842
	}
	if pno < 0 {
		pno = d.PageCount()
	}
	errcode := C.gomupdf_insert_page(d.ctx.ctx, d.pdf, C.int(pno), C.float(width), C.float(height))
	if errcode != 0 {
		return nil, fmt.Errorf("%w: insert page at %d", ErrSave, pno)
	}
	return d.LoadPage(pno)
}

func (d *Document) DeletePage(pno int) error {
	if d.isClosed {
		return ErrClosed
	}
	if !d.IsPDF() {
		return ErrNotPDF
	}
	count := d.PageCount()
	if pno < 0 {
		pno += count
	}
	if pno < 0 || pno >= count {
		return fmt.Errorf("%w: %d", ErrPageNotFound, pno)
	}
	errcode := C.gomupdf_delete_page(d.ctx.ctx, d.pdf, C.int(pno))
	if errcode != 0 {
		return fmt.Errorf("%w: delete page %d", ErrSave, pno)
	}
	return nil
}

func (d *Document) DeletePages(pages ...int) error {
	if d.isClosed {
		return ErrClosed
	}
	if !d.IsPDF() {
		return ErrNotPDF
	}
	sorted := make([]int, len(pages))
	copy(sorted, pages)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] > sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	for _, pno := range sorted {
		if err := d.DeletePage(pno); err != nil {
			return err
		}
	}
	return nil
}

func (d *Document) Select(pages []int) error {
	if d.isClosed {
		return ErrClosed
	}
	if !d.IsPDF() {
		return ErrNotPDF
	}
	if len(pages) == 0 {
		return ErrInvalidArg
	}
	count := d.PageCount()
	for _, p := range pages {
		if p < 0 || p >= count {
			return fmt.Errorf("%w: page %d out of range", ErrInvalidArg, p)
		}
	}
	cPages := make([]C.int, len(pages))
	for i, p := range pages {
		cPages[i] = C.int(p)
	}
	C.gomupdf_rearrange_pages(d.ctx.ctx, d.pdf, C.int(len(pages)), &cPages[0])
	return nil
}

func (d *Document) XrefLength() int {
	if d.isClosed || !d.IsPDF() {
		return 0
	}
	return int(C.gomupdf_xref_len(d.ctx.ctx, d.pdf))
}

func (d *Document) XrefObject(xref int, compressed bool) (string, error) {
	if d.isClosed {
		return "", ErrClosed
	}
	if !d.IsPDF() {
		return "", ErrNotPDF
	}
	comp := 0
	if compressed {
		comp = 1
	}
	var errcode C.int
	result := C.gomupdf_xref_object_str(d.ctx.ctx, d.pdf, C.int(xref), C.int(comp), &errcode)
	if errcode != 0 || result == nil {
		return "", fmt.Errorf("%w: xref %d", ErrXref, xref)
	}
	defer d.ctx.freeString(result)
	return C.GoString(result), nil
}

func (d *Document) PDFCatalog() int {
	if d.isClosed || !d.IsPDF() {
		return 0
	}
	return int(C.gomupdf_pdf_catalog_xref(d.ctx.ctx, d.pdf))
}

func (d *Document) SetMetadata(meta map[string]string) error {
	if d.isClosed {
		return ErrClosed
	}
	if !d.IsPDF() {
		return ErrNotPDF
	}
	keyMap := map[string]string{
		"title":        "info:Title",
		"author":       "info:Author",
		"subject":      "info:Subject",
		"keywords":     "info:Keywords",
		"creator":      "info:Creator",
		"producer":     "info:Producer",
		"creationDate": "info:CreationDate",
		"modDate":      "info:ModDate",
	}
	for goKey, val := range meta {
		pdfKey, ok := keyMap[goKey]
		if !ok {
			continue
		}
		cKey := C.CString(pdfKey)
		cVal := C.CString(val)
		C.gomupdf_set_metadata(d.ctx.ctx, d.doc, cKey, cVal)
		C.free(unsafe.Pointer(cKey))
		C.free(unsafe.Pointer(cVal))
	}
	return nil
}

func (d *Document) InsertPDF(src *Document, opts ...InsertPDFOptions) error {
	if d.isClosed || src.isClosed {
		return ErrClosed
	}
	if !d.IsPDF() || !src.IsPDF() {
		return ErrNotPDF
	}
	opt := InsertPDFOptions{FromPage: -1, ToPage: -1, StartAt: -1, Rotate: -1, Links: true, Annots: true}
	if len(opts) > 0 {
		opt = opts[0]
	}
	srcCount := src.PageCount()
	if opt.FromPage < 0 {
		opt.FromPage = 0
	}
	if opt.ToPage < 0 || opt.ToPage >= srcCount {
		opt.ToPage = srcCount - 1
	}
	if opt.StartAt < 0 {
		opt.StartAt = d.PageCount()
	}

	graftMap := C.pdf_new_graft_map(d.ctx.ctx, d.pdf)
	defer C.pdf_drop_graft_map(d.ctx.ctx, graftMap)

	for i := opt.FromPage; i <= opt.ToPage; i++ {
		pageTo := opt.StartAt + i - opt.FromPage
		errcode := C.gomupdf_graft_page(d.ctx.ctx, d.pdf, src.pdf, graftMap, C.int(pageTo), C.int(i))
		if errcode != 0 {
			return fmt.Errorf("%w: graft page %d", ErrSave, i)
		}
	}
	return nil
}

func (d *Document) CanSaveIncrementally() bool {
	if d.isClosed || !d.IsPDF() {
		return false
	}
	return C.gomupdf_can_save_incrementally(d.ctx.ctx, d.pdf) != 0
}

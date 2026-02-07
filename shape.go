//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import "unsafe"

// InsertText inserts text at the given position on the page. PDF only.
func (p *Page) InsertText(pos Point, text string, opts ...TextInsertOption) (int, error) {
	if !p.doc.IsPDF() {
		return 0, ErrNotPDF
	}
	cfg := textInsertConfig{fontname: "Helvetica", fontsize: 11, color: ColorBlack}
	for _, opt := range opts {
		opt(&cfg)
	}

	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	cFont := C.CString(cfg.fontname)
	defer C.free(unsafe.Pointer(cFont))

	errcode := C.gomupdf_insert_text(p.ctx.ctx, p.doc.pdf, C.int(p.number),
		C.float(pos.X), C.float(pos.Y), cText, cFont, C.float(cfg.fontsize),
		C.float(cfg.color.R), C.float(cfg.color.G), C.float(cfg.color.B))
	if errcode != 0 {
		return 0, ErrSave
	}
	return 1, nil
}

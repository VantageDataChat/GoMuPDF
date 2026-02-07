//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import "unsafe"

// InsertHTMLBox inserts HTML content into a rectangle on the page.
// Uses MuPDF's Story API to layout styled HTML/CSS with automatic line breaks,
// font handling (including CJK), and support for HTML tags like <b>, <i>,
// <span style="...">, <img>, <table>, etc.
//
// This is the Go equivalent of PyMuPDF's Page.insert_htmlbox().
//
// The rect uses top-left origin coordinates (same as other GoMuPDF APIs).
// Returns HTMLBoxResult with spare_height (remaining space) and scale factor.
//
// Example:
//
//	result, err := page.InsertHTMLBox(
//	    gomupdf.Rect{X0: 50, Y0: 50, X1: 500, Y1: 400},
//	    `<p style="color:red; font-size:16px;">Hello <b>World</b></p>`,
//	)
func (p *Page) InsertHTMLBox(rect Rect, html string, opts ...HTMLBoxOptions) (HTMLBoxResult, error) {
	if !p.doc.IsPDF() {
		return HTMLBoxResult{}, ErrNotPDF
	}
	if html == "" {
		return HTMLBoxResult{Scale: 1.0}, nil
	}

	opt := HTMLBoxOptions{ScaleLow: 0, Overlay: true}
	if len(opts) > 0 {
		opt = opts[0]
	}

	cHTML := C.CString(html)
	defer C.free(unsafe.Pointer(cHTML))

	var cCSS *C.char
	if opt.CSS != "" {
		cCSS = C.CString(opt.CSS)
		defer C.free(unsafe.Pointer(cCSS))
	}

	overlay := C.int(0)
	if opt.Overlay {
		overlay = 1
	}

	var spareHeight, scaleUsed C.float
	errcode := C.gomupdf_insert_htmlbox(p.ctx.ctx, p.doc.pdf, C.int(p.number),
		C.float(rect.X0), C.float(rect.Y0), C.float(rect.X1), C.float(rect.Y1),
		cHTML, cCSS, C.float(opt.ScaleLow), overlay,
		&spareHeight, &scaleUsed)

	if errcode == 1 {
		return HTMLBoxResult{}, ErrSave
	}

	result := HTMLBoxResult{
		SpareHeight: float64(spareHeight),
		Scale:       float64(scaleUsed),
	}

	if errcode == 2 {
		// Content overflow — didn't fit, but not a hard error
		return result, ErrOverflow
	}

	return result, nil
}

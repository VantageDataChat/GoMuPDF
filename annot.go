//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import "unsafe"

// Annot represents a PDF annotation.
type Annot struct {
	ctx   *context
	annot *C.pdf_annot
	page  *Page
}

func (a *Annot) Type() int { return int(C.gomupdf_annot_type(a.ctx.ctx, a.annot)) }

func (a *Annot) TypeString() string {
	names := map[int]string{
		AnnotText: "Text", AnnotLink: "Link", AnnotFreeText: "FreeText",
		AnnotLine: "Line", AnnotSquare: "Square", AnnotCircle: "Circle",
		AnnotPolygon: "Polygon", AnnotPolyLine: "PolyLine",
		AnnotHighlight: "Highlight", AnnotUnderline: "Underline",
		AnnotSquiggly: "Squiggly", AnnotStrikeOut: "StrikeOut",
		AnnotRedact: "Redact", AnnotStamp: "Stamp", AnnotCaret: "Caret",
		AnnotInk: "Ink", AnnotPopup: "Popup", AnnotFileAttachment: "FileAttachment",
	}
	if name, ok := names[a.Type()]; ok {
		return name
	}
	return "Unknown"
}

func (a *Annot) Rect() Rect {
	r := C.gomupdf_annot_rect(a.ctx.ctx, a.annot)
	return Rect{X0: float64(r.x0), Y0: float64(r.y0), X1: float64(r.x1), Y1: float64(r.y1)}
}

func (a *Annot) Contents() string {
	s := C.gomupdf_annot_contents(a.ctx.ctx, a.annot)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

func (a *Annot) SetContents(text string) {
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	C.gomupdf_set_annot_contents(a.ctx.ctx, a.annot, cText)
}

func (a *Annot) Xref() int { return int(C.gomupdf_annot_xref(a.ctx.ctx, a.annot)) }

func (p *Page) GetAnnots() []*Annot {
	if !p.doc.IsPDF() {
		return nil
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return nil
	}
	var annots []*Annot
	for annot := C.gomupdf_first_annot(p.ctx.ctx, pdfPage); annot != nil; annot = C.gomupdf_next_annot(p.ctx.ctx, annot) {
		annots = append(annots, &Annot{ctx: p.ctx, annot: annot, page: p})
	}
	return annots
}

func (p *Page) AddTextAnnot(pos Point, text string) (*Annot, error) {
	if !p.doc.IsPDF() {
		return nil, ErrNotPDF
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return nil, ErrNotPDF
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	annot := C.gomupdf_add_text_annot(p.ctx.ctx, pdfPage, C.float(pos.X), C.float(pos.Y), cText)
	if annot == nil {
		return nil, ErrInvalidArg
	}
	return &Annot{ctx: p.ctx, annot: annot, page: p}, nil
}

func (p *Page) AddHighlightAnnot(quads []Quad) (*Annot, error) {
	if !p.doc.IsPDF() {
		return nil, ErrNotPDF
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return nil, ErrNotPDF
	}
	if len(quads) == 0 {
		return nil, ErrInvalidArg
	}
	cQuads := make([]C.fz_quad, len(quads))
	for i, q := range quads {
		cQuads[i] = C.fz_quad{
			ul: C.fz_point{x: C.float(q.UL.X), y: C.float(q.UL.Y)},
			ur: C.fz_point{x: C.float(q.UR.X), y: C.float(q.UR.Y)},
			ll: C.fz_point{x: C.float(q.LL.X), y: C.float(q.LL.Y)},
			lr: C.fz_point{x: C.float(q.LR.X), y: C.float(q.LR.Y)},
		}
	}
	annot := C.gomupdf_add_highlight_annot(p.ctx.ctx, pdfPage, &cQuads[0], C.int(len(cQuads)))
	if annot == nil {
		return nil, ErrInvalidArg
	}
	return &Annot{ctx: p.ctx, annot: annot, page: p}, nil
}

func (p *Page) AddFreetextAnnot(rect Rect, text string, fontsize float64) (*Annot, error) {
	if !p.doc.IsPDF() {
		return nil, ErrNotPDF
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return nil, ErrNotPDF
	}
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cText))
	annot := C.gomupdf_add_freetext_annot(p.ctx.ctx, pdfPage,
		C.float(rect.X0), C.float(rect.Y0), C.float(rect.X1), C.float(rect.Y1),
		cText, C.float(fontsize))
	if annot == nil {
		return nil, ErrInvalidArg
	}
	return &Annot{ctx: p.ctx, annot: annot, page: p}, nil
}

func (p *Page) DeleteAnnot(annot *Annot) error {
	if !p.doc.IsPDF() {
		return ErrNotPDF
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return ErrNotPDF
	}
	C.gomupdf_delete_annot(p.ctx.ctx, pdfPage, annot.annot)
	return nil
}

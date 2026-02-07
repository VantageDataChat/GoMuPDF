//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import "unsafe"

// InsertImage inserts an image into the page at the given rectangle.
func (p *Page) InsertImage(rect Rect, imageData []byte, opts ...InsertImageOptions) error {
	if !p.doc.IsPDF() {
		return ErrNotPDF
	}
	if len(imageData) == 0 {
		return ErrInvalidArg
	}
	opt := InsertImageOptions{KeepProportion: true, Overlay: true}
	if len(opts) > 0 {
		opt = opts[0]
	}
	keepProp := 0
	if opt.KeepProportion {
		keepProp = 1
	}
	overlay := 0
	if opt.Overlay {
		overlay = 1
	}
	errcode := C.gomupdf_insert_image(p.ctx.ctx, p.doc.pdf, C.int(p.number),
		C.float(rect.X0), C.float(rect.Y0), C.float(rect.X1), C.float(rect.Y1),
		(*C.uchar)(unsafe.Pointer(&imageData[0])), C.int(len(imageData)),
		C.int(keepProp), C.int(overlay))
	if errcode != 0 {
		return ErrPixmap
	}
	return nil
}

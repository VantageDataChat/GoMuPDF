//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import (
	"fmt"
	"image"
	"image/color"
	"unsafe"
)

// Pixmap represents a pixel map (raster image).
type Pixmap struct {
	ctx *context
	pix *C.fz_pixmap
}

func NewPixmap(colorspace int, width, height int, alpha bool) (*Pixmap, error) {
	ctx, err := newContext()
	if err != nil {
		return nil, err
	}
	a := 0
	if alpha {
		a = 1
	}
	pix := C.gomupdf_new_pixmap(ctx.ctx, C.int(colorspace), C.int(width), C.int(height), C.int(a))
	if pix == nil {
		ctx.close()
		return nil, ErrPixmap
	}
	return &Pixmap{ctx: ctx, pix: pix}, nil
}

func NewPixmapFromImage(doc *Document, xref int) (*Pixmap, error) {
	if !doc.IsPDF() {
		return nil, ErrNotPDF
	}
	var errcode C.int
	pix := C.gomupdf_pixmap_from_image(doc.ctx.ctx, doc.doc, C.int(xref), &errcode)
	if errcode != 0 || pix == nil {
		return nil, ErrPixmap
	}
	return &Pixmap{ctx: doc.ctx, pix: pix}, nil
}

func (px *Pixmap) Close() {
	if px.pix != nil {
		C.gomupdf_drop_pixmap(px.ctx.ctx, px.pix)
		px.pix = nil
	}
}

func (px *Pixmap) Width() int  { return int(C.gomupdf_pixmap_width(px.pix)) }
func (px *Pixmap) Height() int { return int(C.gomupdf_pixmap_height(px.pix)) }
func (px *Pixmap) N() int      { return int(C.gomupdf_pixmap_n(px.pix)) }
func (px *Pixmap) Alpha() int  { return int(C.gomupdf_pixmap_alpha(px.pix)) }
func (px *Pixmap) Stride() int { return int(C.gomupdf_pixmap_stride(px.pix)) }
func (px *Pixmap) X() int      { return int(C.gomupdf_pixmap_x(px.pix)) }
func (px *Pixmap) Y() int      { return int(C.gomupdf_pixmap_y(px.pix)) }

func (px *Pixmap) IRect() IRect {
	return NewIRect(px.X(), px.Y(), px.X()+px.Width(), px.Y()+px.Height())
}

func (px *Pixmap) Samples() []byte {
	ptr := C.gomupdf_pixmap_samples(px.pix)
	length := C.gomupdf_pixmap_samples_len(px.pix)
	return C.GoBytes(unsafe.Pointer(ptr), length)
}

func (px *Pixmap) ToBytes() ([]byte, error) {
	var outlen, errcode C.int
	data := C.gomupdf_pixmap_to_png(px.ctx.ctx, px.pix, &outlen, &errcode)
	if errcode != 0 || data == nil {
		return nil, ErrPixmap
	}
	defer px.ctx.freeBytes(data)
	return C.GoBytes(unsafe.Pointer(data), outlen), nil
}

func (px *Pixmap) Save(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	errcode := C.gomupdf_pixmap_save_png(px.ctx.ctx, px.pix, cFilename)
	if errcode != 0 {
		return fmt.Errorf("%w: %s", ErrPixmap, filename)
	}
	return nil
}

func (px *Pixmap) SavePNG(filename string) error { return px.Save(filename) }

func (px *Pixmap) SavePNM(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))
	errcode := C.gomupdf_pixmap_save_pnm(px.ctx.ctx, px.pix, cFilename)
	if errcode != 0 {
		return fmt.Errorf("%w: %s", ErrPixmap, filename)
	}
	return nil
}

func (px *Pixmap) SetPixel(x, y int, c []byte) {
	if x < 0 || x >= px.Width() || y < 0 || y >= px.Height() {
		return
	}
	n := px.N()
	if len(c) < n {
		return
	}
	C.gomupdf_pixmap_set_pixel(px.pix, C.int(x), C.int(y), (*C.uchar)(unsafe.Pointer(&c[0])), C.int(n))
}

func (px *Pixmap) GetPixel(x, y int) []byte {
	n := px.N()
	c := make([]byte, n)
	if x < 0 || x >= px.Width() || y < 0 || y >= px.Height() {
		return c
	}
	C.gomupdf_pixmap_get_pixel(px.pix, C.int(x), C.int(y), (*C.uchar)(unsafe.Pointer(&c[0])), C.int(n))
	return c
}

func (px *Pixmap) Clear(value int) {
	C.gomupdf_pixmap_clear(px.ctx.ctx, px.pix, C.int(value))
}

func (px *Pixmap) Invert() { C.gomupdf_pixmap_invert(px.ctx.ctx, px.pix) }

func (px *Pixmap) Gamma(gamma float64) {
	C.gomupdf_pixmap_gamma(px.ctx.ctx, px.pix, C.float(gamma))
}

func (px *Pixmap) Tint(black, white int) {
	C.gomupdf_pixmap_tint(px.ctx.ctx, px.pix, C.int(black), C.int(white))
}

func (px *Pixmap) Convert(colorspace int) (*Pixmap, error) {
	var errcode C.int
	newPix := C.gomupdf_pixmap_convert(px.ctx.ctx, px.pix, C.int(colorspace), &errcode)
	if errcode != 0 || newPix == nil {
		return nil, ErrPixmap
	}
	return &Pixmap{ctx: px.ctx, pix: newPix}, nil
}

func (px *Pixmap) ToImage() image.Image {
	w, h, n, alpha := px.Width(), px.Height(), px.N(), px.Alpha()
	samples := px.Samples()
	stride := px.Stride()

	switch {
	case n-alpha == 1 && alpha == 0:
		img := image.NewGray(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				img.SetGray(x, y, color.Gray{Y: samples[y*stride+x]})
			}
		}
		return img
	case n-alpha == 1 && alpha == 1:
		img := image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*stride + x*2
				g := samples[idx]
				img.SetNRGBA(x, y, color.NRGBA{R: g, G: g, B: g, A: samples[idx+1]})
			}
		}
		return img
	case n-alpha == 3 && alpha == 0:
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*stride + x*3
				img.SetRGBA(x, y, color.RGBA{R: samples[idx], G: samples[idx+1], B: samples[idx+2], A: 255})
			}
		}
		return img
	case n-alpha == 3 && alpha == 1:
		img := image.NewNRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*stride + x*4
				img.SetNRGBA(x, y, color.NRGBA{R: samples[idx], G: samples[idx+1], B: samples[idx+2], A: samples[idx+3]})
			}
		}
		return img
	default:
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				idx := y*stride + x*n
				if n >= 4 && alpha == 0 {
					r := 255 - min8(255, uint8(int(samples[idx])+int(samples[idx+3])))
					g := 255 - min8(255, uint8(int(samples[idx+1])+int(samples[idx+3])))
					b := 255 - min8(255, uint8(int(samples[idx+2])+int(samples[idx+3])))
					img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
				}
			}
		}
		return img
	}
}

func min8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

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

// Page represents a document page.
type Page struct {
	ctx    *context
	page   *C.fz_page
	doc    *Document
	number int
}

func (p *Page) Close() {
	if p.page != nil {
		C.gomupdf_drop_page(p.ctx.ctx, p.page)
		p.page = nil
	}
}

func (p *Page) Number() int { return p.number }

func (p *Page) Rect() Rect {
	r := C.gomupdf_page_bound(p.ctx.ctx, p.page)
	return Rect{X0: float64(r.x0), Y0: float64(r.y0), X1: float64(r.x1), Y1: float64(r.y1)}
}

func (p *Page) Width() float64  { return p.Rect().Width() }
func (p *Page) Height() float64 { return p.Rect().Height() }

func (p *Page) Rotation() int {
	if !p.doc.IsPDF() {
		return 0
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return 0
	}
	return int(C.gomupdf_page_rotation(p.ctx.ctx, pdfPage))
}

func (p *Page) SetRotation(rotation int) error {
	if !p.doc.IsPDF() {
		return ErrNotPDF
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return ErrNotPDF
	}
	C.gomupdf_set_page_rotation(p.ctx.ctx, pdfPage, C.int(rotation))
	return nil
}

func (p *Page) GetText(output string, flags ...int) (string, error) {
	if output == "" {
		output = "text"
	}
	flag := TextFlagsDefault
	if len(flags) > 0 {
		flag = flags[0]
	}
	var errcode C.int
	tp := C.gomupdf_new_stext_page(p.ctx.ctx, p.page, C.int(flag), &errcode)
	if errcode != 0 || tp == nil {
		return "", ErrTextExtract
	}
	defer C.gomupdf_drop_stext_page(p.ctx.ctx, tp)

	cText := C.gomupdf_stext_page_as_text(p.ctx.ctx, tp, &errcode)
	if errcode != 0 || cText == nil {
		return "", ErrTextExtract
	}
	defer p.ctx.freeString(cText)
	return C.GoString(cText), nil
}

func (p *Page) GetTextWords(flags ...int) ([]TextWord, error) {
	flag := TextFlagsDefault
	if len(flags) > 0 {
		flag = flags[0]
	}
	var errcode C.int
	tp := C.gomupdf_new_stext_page(p.ctx.ctx, p.page, C.int(flag), &errcode)
	if errcode != 0 || tp == nil {
		return nil, ErrTextExtract
	}
	defer C.gomupdf_drop_stext_page(p.ctx.ctx, tp)

	var words []TextWord
	for block := tp.first_block; block != nil; block = block.next {
		if block._type != C.FZ_STEXT_BLOCK_TEXT {
			continue
		}
		for line := C.gomupdf_stext_block_first_line(block); line != nil; line = line.next {
			word := TextWord{BlockNo: len(words)}
			wordText := ""
			for ch := line.first_char; ch != nil; ch = ch.next {
				c := rune(ch.c)
				if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
					if wordText != "" {
						word.Text = wordText
						word.Rect.Y0 = float64(line.bbox.y0)
						word.Rect.X1 = float64(ch.origin.x)
						word.Rect.Y1 = float64(line.bbox.y1)
						words = append(words, word)
						word = TextWord{BlockNo: len(words)}
						wordText = ""
					}
					continue
				}
				if wordText == "" {
					word.Rect.X0 = float64(ch.origin.x)
				}
				wordText += string(c)
			}
			if wordText != "" {
				word.Text = wordText
				word.Rect.Y0 = float64(line.bbox.y0)
				word.Rect.Y1 = float64(line.bbox.y1)
				words = append(words, word)
			}
		}
	}
	return words, nil
}

func (p *Page) GetTextBlocks(flags ...int) ([]TextBlock, error) {
	flag := TextFlagsDefault
	if len(flags) > 0 {
		flag = flags[0]
	}
	var errcode C.int
	tp := C.gomupdf_new_stext_page(p.ctx.ctx, p.page, C.int(flag), &errcode)
	if errcode != 0 || tp == nil {
		return nil, ErrTextExtract
	}
	defer C.gomupdf_drop_stext_page(p.ctx.ctx, tp)

	var blocks []TextBlock
	blockNo := 0
	for block := tp.first_block; block != nil; block = block.next {
		tb := TextBlock{
			BlockNo: blockNo,
			Rect: Rect{
				X0: float64(block.bbox.x0), Y0: float64(block.bbox.y0),
				X1: float64(block.bbox.x1), Y1: float64(block.bbox.y1),
			},
		}
		if block._type == C.FZ_STEXT_BLOCK_TEXT {
			tb.Type = "text"
			text := ""
			for line := C.gomupdf_stext_block_first_line(block); line != nil; line = line.next {
				for ch := line.first_char; ch != nil; ch = ch.next {
					text += string(rune(ch.c))
				}
				text += "\n"
			}
			tb.Text = text
		} else {
			tb.Type = "image"
		}
		blocks = append(blocks, tb)
		blockNo++
	}
	return blocks, nil
}

func (p *Page) GetPixmap(opts ...PixmapOption) (*Pixmap, error) {
	cfg := pixmapConfig{matrix: Identity, colorspace: CsRGB, alpha: false, annots: true}
	for _, opt := range opts {
		opt(&cfg)
	}
	if cfg.dpi > 0 {
		scale := float64(cfg.dpi) / 72.0
		cfg.matrix = ScaleMatrix(scale, scale).Concat(cfg.matrix)
	}
	alpha := 0
	if cfg.alpha {
		alpha = 1
	}
	var errcode C.int
	var pix *C.fz_pixmap
	if cfg.clip != nil {
		pix = C.gomupdf_page_to_pixmap_clipped(p.ctx.ctx, p.page,
			C.float(cfg.matrix.A), C.float(cfg.matrix.B),
			C.float(cfg.matrix.C), C.float(cfg.matrix.D),
			C.float(cfg.matrix.E), C.float(cfg.matrix.F),
			C.int(cfg.colorspace), C.int(alpha),
			C.float(cfg.clip.X0), C.float(cfg.clip.Y0),
			C.float(cfg.clip.X1), C.float(cfg.clip.Y1), &errcode)
	} else {
		pix = C.gomupdf_page_to_pixmap(p.ctx.ctx, p.page,
			C.float(cfg.matrix.A), C.float(cfg.matrix.B),
			C.float(cfg.matrix.C), C.float(cfg.matrix.D),
			C.float(cfg.matrix.E), C.float(cfg.matrix.F),
			C.int(cfg.colorspace), C.int(alpha), &errcode)
	}
	if errcode != 0 || pix == nil {
		return nil, ErrPixmap
	}
	return &Pixmap{ctx: p.ctx, pix: pix}, nil
}

func (p *Page) SearchFor(needle string, quads bool) ([]Quad, error) {
	cNeedle := C.CString(needle)
	defer C.free(unsafe.Pointer(cNeedle))
	maxQuads := 500
	cQuads := make([]C.fz_quad, maxQuads)
	var errcode C.int
	count := int(C.gomupdf_search_page(p.ctx.ctx, p.page, cNeedle, &cQuads[0], C.int(maxQuads), &errcode))
	if errcode != 0 {
		return nil, ErrSearch
	}
	result := make([]Quad, count)
	for i := 0; i < count; i++ {
		q := cQuads[i]
		result[i] = Quad{
			UL: Point{float64(q.ul.x), float64(q.ul.y)},
			UR: Point{float64(q.ur.x), float64(q.ur.y)},
			LL: Point{float64(q.ll.x), float64(q.ll.y)},
			LR: Point{float64(q.lr.x), float64(q.lr.y)},
		}
	}
	return result, nil
}

func (p *Page) GetLinks() ([]Link, error) {
	var errcode C.int
	fzLinks := C.gomupdf_load_links(p.ctx.ctx, p.page, &errcode)
	if errcode != 0 {
		return nil, fmt.Errorf("failed to load links")
	}
	if fzLinks == nil {
		return nil, nil
	}
	defer C.gomupdf_drop_link(p.ctx.ctx, fzLinks)
	var links []Link
	for link := fzLinks; link != nil; link = link.next {
		l := Link{
			Rect: Rect{
				X0: float64(link.rect.x0), Y0: float64(link.rect.y0),
				X1: float64(link.rect.x1), Y1: float64(link.rect.y1),
			},
		}
		if link.uri != nil {
			l.URI = C.GoString(link.uri)
		}
		links = append(links, l)
	}
	return links, nil
}

func (p *Page) GetFonts() ([]FontInfo, error)   { return nil, nil }
func (p *Page) GetImages() ([]ImageInfo, error) { return nil, nil }

func (p *Page) GetLabel() string {
	var errcode C.int
	label := C.gomupdf_page_label(p.ctx.ctx, p.page, &errcode)
	if errcode != 0 || label == nil {
		return ""
	}
	defer p.ctx.freeString(label)
	return C.GoString(label)
}

func (p *Page) MediaBox() Rect { return p.Rect() }

func (p *Page) CropBox() Rect {
	if !p.doc.IsPDF() {
		return p.Rect()
	}
	return p.Rect()
}

func (p *Page) TransformationMatrix() Matrix {
	r := p.Rect()
	switch p.Rotation() {
	case 90:
		return NewMatrix(0, -1, 1, 0, 0, r.Height())
	case 180:
		return NewMatrix(-1, 0, 0, -1, r.Width(), r.Height())
	case 270:
		return NewMatrix(0, 1, -1, 0, r.Width(), 0)
	default:
		return Identity
	}
}

func (p *Page) DerotationMatrix() Matrix {
	m := p.TransformationMatrix()
	inv, ok := m.Invert()
	if !ok {
		return Identity
	}
	return inv
}

//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"

// TextPage represents extracted text and images from a page.
type TextPage struct {
	ctx *context
	tp  *C.fz_stext_page
}

func (t *TextPage) Close() {
	if t.tp != nil {
		C.gomupdf_drop_stext_page(t.ctx.ctx, t.tp)
		t.tp = nil
	}
}

func (t *TextPage) ExtractText() (string, error) {
	var errcode C.int
	cText := C.gomupdf_stext_page_as_text(t.ctx.ctx, t.tp, &errcode)
	if errcode != 0 || cText == nil {
		return "", ErrTextExtract
	}
	defer t.ctx.freeString(cText)
	return C.GoString(cText), nil
}

func (t *TextPage) Blocks() []STextBlock {
	var blocks []STextBlock
	for block := t.tp.first_block; block != nil; block = block.next {
		b := STextBlock{
			Rect: Rect{
				X0: float64(block.bbox.x0), Y0: float64(block.bbox.y0),
				X1: float64(block.bbox.x1), Y1: float64(block.bbox.y1),
			},
		}
		if block._type == C.FZ_STEXT_BLOCK_TEXT {
			b.Type = STextBlockText
			b.Lines = t.extractLines(block)
		} else {
			b.Type = STextBlockImage
		}
		blocks = append(blocks, b)
	}
	return blocks
}

func (t *TextPage) extractLines(block *C.fz_stext_block) []STextLine {
	var lines []STextLine
	for line := C.gomupdf_stext_block_first_line(block); line != nil; line = line.next {
		l := STextLine{
			Rect: Rect{
				X0: float64(line.bbox.x0), Y0: float64(line.bbox.y0),
				X1: float64(line.bbox.x1), Y1: float64(line.bbox.y1),
			},
			Dir: Point{X: float64(line.dir.x), Y: float64(line.dir.y)},
		}
		for ch := line.first_char; ch != nil; ch = ch.next {
			l.Chars = append(l.Chars, STextChar{
				C:      rune(ch.c),
				Origin: Point{X: float64(ch.origin.x), Y: float64(ch.origin.y)},
				Rect: Rect{
					X0: float64(ch.quad.ul.x), Y0: float64(ch.quad.ul.y),
					X1: float64(ch.quad.lr.x), Y1: float64(ch.quad.lr.y),
				},
			})
		}
		lines = append(lines, l)
	}
	return lines
}

// GetTextPage creates a TextPage from a Page for detailed text analysis.
func (p *Page) GetTextPage(flags ...int) (*TextPage, error) {
	flag := TextFlagsDefault
	if len(flags) > 0 {
		flag = flags[0]
	}
	var errcode C.int
	tp := C.gomupdf_new_stext_page(p.ctx.ctx, p.page, C.int(flag), &errcode)
	if errcode != 0 || tp == nil {
		return nil, ErrTextExtract
	}
	return &TextPage{ctx: p.ctx, tp: tp}, nil
}

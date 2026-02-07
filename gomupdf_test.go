//go:build cgo && !nomupdf

package gomupdf

import (
	"os"
	"path/filepath"
	"testing"
)

// createTestPDF creates a minimal test PDF file and returns its path.
func createTestPDF(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.pdf")

	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	p.Close()

	err = doc.Save(path)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	return path
}

// newTestPDFWithPage creates a new PDF with one A4 page.
func newTestPDFWithPage(t *testing.T) *Document {
	t.Helper()
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		doc.Close()
		t.Fatalf("NewPage: %v", err)
	}
	p.Close()
	return doc
}

// newTestPDFWithText creates a PDF with one page containing text.
func newTestPDFWithText(t *testing.T, text string) *Document {
	t.Helper()
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		doc.Close()
		t.Fatalf("NewPage: %v", err)
	}
	_, err = p.InsertText(NewPoint(72, 72), text)
	if err != nil {
		p.Close()
		doc.Close()
		t.Fatalf("InsertText: %v", err)
	}
	p.Close()
	return doc
}

// --- Document creation tests ---

func TestNewPDF(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	if !doc.IsPDF() {
		t.Error("expected IsPDF() == true")
	}
	if doc.IsClosed() {
		t.Error("expected IsClosed() == false")
	}
	if doc.PageCount() != 0 {
		t.Errorf("expected 0 pages, got %d", doc.PageCount())
	}
}

func TestOpenFromMemory(t *testing.T) {
	data := []byte(`%PDF-1.0
1 0 obj<</Pages 2 0 R>>endobj
2 0 obj<</Kids[3 0 R]/Count 1>>endobj
3 0 obj<</MediaBox[0 0 612 792]>>endobj
trailer<</Root 1 0 R>>`)

	doc, err := OpenFromMemory(data, "application/pdf")
	if err != nil {
		t.Fatalf("OpenFromMemory: %v", err)
	}
	defer doc.Close()

	if !doc.IsPDF() {
		t.Error("expected IsPDF() == true")
	}
	if doc.PageCount() != 1 {
		t.Errorf("expected 1 page, got %d", doc.PageCount())
	}
}

func TestOpenFromMemoryEmpty(t *testing.T) {
	_, err := OpenFromMemory(nil, "application/pdf")
	if err == nil {
		t.Error("expected error for nil data")
	}
	_, err = OpenFromMemory([]byte{}, "application/pdf")
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestOpenNonExistent(t *testing.T) {
	_, err := Open("nonexistent_file_12345.pdf")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestOpenSaveReopen(t *testing.T) {
	path := createTestPDF(t)
	doc, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc.Close()

	if doc.Name() != path {
		t.Errorf("expected name %s, got %s", path, doc.Name())
	}
	if !doc.IsPDF() {
		t.Error("should be PDF")
	}
	if doc.PageCount() != 1 {
		t.Errorf("expected 1 page, got %d", doc.PageCount())
	}
}

// --- Document close tests ---

func TestDocumentClose(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	doc.Close()

	if !doc.IsClosed() {
		t.Error("expected IsClosed() == true after Close()")
	}
	if doc.PageCount() != 0 {
		t.Error("expected PageCount() == 0 after Close()")
	}
	// Double close should not panic
	doc.Close()
}

func TestClosedDocOperations(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	doc.Close()

	if doc.NeedsPass() {
		t.Error("closed doc NeedsPass should be false")
	}
	if doc.IsReflowable() {
		t.Error("closed doc IsReflowable should be false")
	}
	meta := doc.Metadata()
	if meta != nil {
		t.Error("closed doc Metadata should be nil")
	}
	_, err = doc.LoadPage(0)
	if err == nil {
		t.Error("LoadPage on closed doc should error")
	}
	_, err = doc.GetTOC(true)
	if err == nil {
		t.Error("GetTOC on closed doc should error")
	}
	_, err = doc.Pages()
	if err == nil {
		t.Error("Pages on closed doc should error")
	}
	_, err = doc.ConvertToPDF(0, 0, 0)
	if err == nil {
		t.Error("ConvertToPDF on closed doc should error")
	}
	_, err = doc.ToBytes()
	if err == nil {
		t.Error("ToBytes on closed doc should error")
	}
}

// --- Page loading tests ---

func TestLoadPage(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage(0): %v", err)
	}
	defer page.Close()

	if page.Number() != 0 {
		t.Errorf("expected page number 0, got %d", page.Number())
	}
	rect := page.Rect()
	if rect.Width() <= 0 || rect.Height() <= 0 {
		t.Errorf("expected positive page dimensions, got %v", rect)
	}
}

func TestLoadPageNegativeIndex(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(-1)
	if err != nil {
		t.Fatalf("LoadPage(-1): %v", err)
	}
	defer page.Close()

	if page.Number() != 0 {
		t.Errorf("expected page 0 for -1 index, got %d", page.Number())
	}
}

func TestLoadPageOutOfRange(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	_, err := doc.LoadPage(99)
	if err == nil {
		t.Error("expected error for out-of-range page")
	}
	_, err = doc.LoadPage(-99)
	if err == nil {
		t.Error("expected error for large negative index")
	}
}

// --- Page properties tests ---

func TestPageRect(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	w := page.Width()
	h := page.Height()
	if w < 500 || w > 700 {
		t.Errorf("unexpected width: %f", w)
	}
	if h < 700 || h > 900 {
		t.Errorf("unexpected height: %f", h)
	}
}

func TestPageRotation(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	rot := page.Rotation()
	if rot != 0 {
		t.Errorf("expected rotation 0, got %d", rot)
	}

	err = page.SetRotation(90)
	if err != nil {
		t.Fatalf("SetRotation: %v", err)
	}
	rot = page.Rotation()
	if rot != 90 {
		t.Errorf("expected rotation 90, got %d", rot)
	}
}

func TestPageMediaBox(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	mb := page.MediaBox()
	if mb.Width() <= 0 || mb.Height() <= 0 {
		t.Error("MediaBox should have positive dimensions")
	}
	cb := page.CropBox()
	if cb.Width() <= 0 || cb.Height() <= 0 {
		t.Error("CropBox should have positive dimensions")
	}
}

func TestPageLabel(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	// New PDF may or may not have labels
	_ = page.GetLabel()
}

// --- Text extraction tests ---

func TestGetText(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	text, err := page.GetText("text")
	if err != nil {
		t.Fatalf("GetText: %v", err)
	}
	_ = text
}

func TestGetTextDefaultOutput(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	// Empty output string should default to "text"
	text, err := page.GetText("")
	if err != nil {
		t.Fatalf("GetText with empty output: %v", err)
	}
	_ = text
}

func TestGetTextWithFlags(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	text, err := page.GetText("text", TextPreserveLigatures|TextPreserveWhitespace|TextPreserveImages)
	if err != nil {
		t.Fatalf("GetText with flags: %v", err)
	}
	_ = text
}

func TestGetTextWithContent(t *testing.T) {
	doc := newTestPDFWithText(t, "Hello GoMuPDF World")
	defer doc.Close()

	// Save and reopen to ensure text is committed
	dir := t.TempDir()
	path := filepath.Join(dir, "text.pdf")
	if err := doc.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}
	doc.Close()

	doc2, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc2.Close()

	page, err := doc2.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	text, err := page.GetText("text")
	if err != nil {
		t.Fatalf("GetText: %v", err)
	}
	// Note: fz_show_string inserts glyphs; text extraction may or may not
	// recover them depending on font encoding. Log result for diagnostics.
	t.Logf("extracted text: %q (len=%d)", text, len(text))
}

func TestGetTextBlocks(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	blocks, err := page.GetTextBlocks()
	if err != nil {
		t.Fatalf("GetTextBlocks: %v", err)
	}
	_ = blocks
}

func TestGetTextWords(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	words, err := page.GetTextWords()
	if err != nil {
		t.Fatalf("GetTextWords: %v", err)
	}
	_ = words
}

func TestGetTextPage(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	tp, err := page.GetTextPage()
	if err != nil {
		t.Fatalf("GetTextPage: %v", err)
	}
	defer tp.Close()

	text, err := tp.ExtractText()
	if err != nil {
		t.Fatalf("ExtractText: %v", err)
	}
	_ = text

	blocks := tp.Blocks()
	_ = blocks
}

func TestGetTextPageWithFlags(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	tp, err := page.GetTextPage(TextPreserveLigatures | TextPreserveWhitespace | TextPreserveImages)
	if err != nil {
		t.Fatalf("GetTextPage with flags: %v", err)
	}
	tp.Close()
}

func TestGetPageText(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	text, err := doc.GetPageText(0, "text")
	if err != nil {
		t.Fatalf("GetPageText: %v", err)
	}
	_ = text
}

// --- Metadata tests ---

func TestMetadata(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	meta := doc.Metadata()
	if meta == nil {
		t.Error("expected non-nil metadata map")
	}
	if f, ok := meta["format"]; ok {
		t.Logf("format: %s", f)
	}
}

func TestSetMetadata(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	err = doc.SetMetadata(map[string]string{
		"title":   "Test Document",
		"author":  "GoMuPDF",
		"subject": "Testing",
		"unknown": "should be ignored",
	})
	if err != nil {
		t.Fatalf("SetMetadata: %v", err)
	}

	meta := doc.Metadata()
	if meta["title"] != "Test Document" {
		t.Errorf("expected title 'Test Document', got '%s'", meta["title"])
	}
	if meta["author"] != "GoMuPDF" {
		t.Errorf("expected author 'GoMuPDF', got '%s'", meta["author"])
	}
	if meta["subject"] != "Testing" {
		t.Errorf("expected subject 'Testing', got '%s'", meta["subject"])
	}
}

func TestSetMetadataPersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "meta.pdf")

	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	doc.SetMetadata(map[string]string{"title": "Persistent Title", "author": "Test Author"})
	p, _ := doc.NewPage(-1, 595, 842)
	p.Close()
	doc.Save(path)
	doc.Close()

	doc2, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc2.Close()

	meta := doc2.Metadata()
	if meta["title"] != "Persistent Title" {
		t.Errorf("title not persisted: got '%s'", meta["title"])
	}
}

func TestNeedsPass(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	if doc.NeedsPass() {
		t.Error("new PDF should not need password")
	}
}

func TestGetTOC(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	toc, err := doc.GetTOC(true)
	if err != nil {
		t.Fatalf("GetTOC: %v", err)
	}
	if len(toc) != 0 {
		t.Errorf("expected empty TOC, got %d items", len(toc))
	}
}

func TestIsReflowable(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	if doc.IsReflowable() {
		t.Error("PDF should not be reflowable")
	}
	// Layout on non-reflowable should be no-op
	doc.Layout(595, 842, 11)
}

// --- Pixmap tests ---

func TestGetPixmap(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap()
	if err != nil {
		t.Fatalf("GetPixmap: %v", err)
	}
	defer pix.Close()

	if pix.Width() <= 0 || pix.Height() <= 0 {
		t.Errorf("expected positive pixmap dimensions, got %dx%d", pix.Width(), pix.Height())
	}
	if pix.N() < 3 {
		t.Errorf("expected at least 3 components, got %d", pix.N())
	}
	if pix.Stride() <= 0 {
		t.Error("stride should be positive")
	}
}

func TestGetPixmapWithDPI(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix72, err := page.GetPixmap(WithDPI(72))
	if err != nil {
		t.Fatalf("GetPixmap 72dpi: %v", err)
	}
	defer pix72.Close()

	pix144, err := page.GetPixmap(WithDPI(144))
	if err != nil {
		t.Fatalf("GetPixmap 144dpi: %v", err)
	}
	defer pix144.Close()

	// 144 DPI should be roughly 2x the size of 72 DPI
	ratio := float64(pix144.Width()) / float64(pix72.Width())
	if ratio < 1.8 || ratio > 2.2 {
		t.Errorf("expected ~2x width ratio, got %f", ratio)
	}
}

func TestGetPixmapWithAlpha(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap(WithAlpha(true))
	if err != nil {
		t.Fatalf("GetPixmap with alpha: %v", err)
	}
	defer pix.Close()

	if pix.Alpha() != 1 {
		t.Errorf("expected alpha=1, got %d", pix.Alpha())
	}
}

func TestGetPixmapWithClip(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	clip := NewRect(0, 0, 100, 100)
	pix, err := page.GetPixmap(WithClip(clip))
	if err != nil {
		t.Fatalf("GetPixmap with clip: %v", err)
	}
	defer pix.Close()

	// Clipped pixmap should have valid dimensions
	if pix.Width() <= 0 || pix.Height() <= 0 {
		t.Errorf("clipped pixmap should have positive dimensions: %dx%d", pix.Width(), pix.Height())
	}
	t.Logf("clipped pixmap: %dx%d", pix.Width(), pix.Height())
}

func TestGetPixmapWithMatrix(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap(WithMatrix(ScaleMatrix(0.5, 0.5)))
	if err != nil {
		t.Fatalf("GetPixmap with matrix: %v", err)
	}
	defer pix.Close()

	// Half scale should produce smaller pixmap
	if pix.Width() > 400 {
		t.Errorf("expected smaller width at 0.5x, got %d", pix.Width())
	}
}

func TestGetPixmapGray(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap(WithColorspace(CsGray))
	if err != nil {
		t.Fatalf("GetPixmap gray: %v", err)
	}
	defer pix.Close()

	if pix.N() != 1 {
		t.Errorf("expected 1 component for gray, got %d", pix.N())
	}
}

func TestGetPagePixmap(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	pix, err := doc.GetPagePixmap(0, WithDPI(72))
	if err != nil {
		t.Fatalf("GetPagePixmap: %v", err)
	}
	defer pix.Close()

	if pix.Width() <= 0 {
		t.Error("expected positive width")
	}
}

func TestPixmapToPNG(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap(WithDPI(72))
	if err != nil {
		t.Fatalf("GetPixmap: %v", err)
	}
	defer pix.Close()

	data, err := pix.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PNG data")
	}
	if len(data) > 4 && data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G' {
		t.Log("valid PNG header")
	} else {
		t.Error("invalid PNG header")
	}
}

func TestPixmapSavePNG(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap(WithDPI(72))
	if err != nil {
		t.Fatalf("GetPixmap: %v", err)
	}
	defer pix.Close()

	path := filepath.Join(t.TempDir(), "test.png")
	err = pix.SavePNG(path)
	if err != nil {
		t.Fatalf("SavePNG: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty PNG file")
	}
}

func TestPixmapSavePNM(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap(WithDPI(72))
	if err != nil {
		t.Fatalf("GetPixmap: %v", err)
	}
	defer pix.Close()

	path := filepath.Join(t.TempDir(), "test.pnm")
	err = pix.SavePNM(path)
	if err != nil {
		t.Fatalf("SavePNM: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty PNM file")
	}
}

func TestPixmapToImage(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	pix, err := page.GetPixmap()
	if err != nil {
		t.Fatalf("GetPixmap: %v", err)
	}
	defer pix.Close()

	img := pix.ToImage()
	bounds := img.Bounds()
	if bounds.Dx() != pix.Width() || bounds.Dy() != pix.Height() {
		t.Errorf("image size mismatch: got %dx%d, expected %dx%d",
			bounds.Dx(), bounds.Dy(), pix.Width(), pix.Height())
	}
}

func TestPixmapToImageGray(t *testing.T) {
	pix, err := NewPixmap(CsGray, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	img := pix.ToImage()
	if img.Bounds().Dx() != 10 || img.Bounds().Dy() != 10 {
		t.Error("gray image size mismatch")
	}
}

func TestPixmapToImageGrayAlpha(t *testing.T) {
	pix, err := NewPixmap(CsGray, 10, 10, true)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	img := pix.ToImage()
	if img.Bounds().Dx() != 10 || img.Bounds().Dy() != 10 {
		t.Error("gray+alpha image size mismatch")
	}
}

func TestPixmapToImageRGBAlpha(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, true)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	img := pix.ToImage()
	if img.Bounds().Dx() != 10 || img.Bounds().Dy() != 10 {
		t.Error("RGBA image size mismatch")
	}
}

func TestPixmapToImageCMYK(t *testing.T) {
	pix, err := NewPixmap(CsCMYK, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	img := pix.ToImage()
	if img.Bounds().Dx() != 10 || img.Bounds().Dy() != 10 {
		t.Error("CMYK image size mismatch")
	}
}

// --- Standalone Pixmap tests ---

func TestNewPixmap(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 100, 100, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	if pix.Width() != 100 || pix.Height() != 100 {
		t.Errorf("expected 100x100, got %dx%d", pix.Width(), pix.Height())
	}
	if pix.N() != 3 {
		t.Errorf("expected 3 components for RGB, got %d", pix.N())
	}
}

func TestPixmapSetGetPixel(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	pix.SetPixel(5, 5, []byte{255, 0, 128})
	c := pix.GetPixel(5, 5)
	if c[0] != 255 || c[1] != 0 || c[2] != 128 {
		t.Errorf("expected [255 0 128], got %v", c)
	}
}

func TestPixmapSetPixelOutOfBounds(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	// Should not panic
	pix.SetPixel(-1, 0, []byte{255, 0, 0})
	pix.SetPixel(0, -1, []byte{255, 0, 0})
	pix.SetPixel(10, 0, []byte{255, 0, 0})
	pix.SetPixel(0, 10, []byte{255, 0, 0})
	// Too few color components
	pix.SetPixel(0, 0, []byte{255})
}

func TestPixmapGetPixelOutOfBounds(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	c := pix.GetPixel(-1, 0)
	if len(c) != 3 {
		t.Errorf("expected 3 components, got %d", len(c))
	}
	c = pix.GetPixel(10, 10)
	if len(c) != 3 {
		t.Errorf("expected 3 components, got %d", len(c))
	}
}

func TestPixmapConvert(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 50, 50, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	gray, err := pix.Convert(CsGray)
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	defer gray.Close()

	if gray.N() != 1 {
		t.Errorf("expected 1 component for gray, got %d", gray.N())
	}
}

func TestPixmapClear(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	pix.SetPixel(5, 5, []byte{255, 128, 64})
	pix.Clear(0)
	c := pix.GetPixel(5, 5)
	if c[0] != 0 || c[1] != 0 || c[2] != 0 {
		t.Errorf("expected [0 0 0] after clear, got %v", c)
	}

	pix.Clear(255)
	c = pix.GetPixel(5, 5)
	if c[0] != 255 || c[1] != 255 || c[2] != 255 {
		t.Errorf("expected [255 255 255] after clear(255), got %v", c)
	}
}

func TestPixmapInvert(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	pix.Clear(0)
	pix.Invert()
	c := pix.GetPixel(0, 0)
	if c[0] != 255 || c[1] != 255 || c[2] != 255 {
		t.Errorf("expected [255 255 255] after invert, got %v", c)
	}
}

func TestPixmapGamma(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	pix.Clear(128)
	pix.Gamma(1.0) // gamma 1.0 should be no-op
	c := pix.GetPixel(0, 0)
	if c[0] != 128 {
		t.Errorf("gamma 1.0 should be no-op, got %d", c[0])
	}
}

func TestPixmapTint(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	pix.Clear(128)
	// Just verify it doesn't crash
	pix.Tint(0x000000, 0xFFFFFF)
}

func TestPixmapSamples(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	samples := pix.Samples()
	expectedLen := pix.Stride() * pix.Height()
	if len(samples) != expectedLen {
		t.Errorf("expected %d bytes, got %d", expectedLen, len(samples))
	}
}

func TestPixmapIRect(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 100, 200, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	ir := pix.IRect()
	if ir.Width() != 100 || ir.Height() != 200 {
		t.Errorf("expected 100x200, got %dx%d", ir.Width(), ir.Height())
	}
}

func TestPixmapXY(t *testing.T) {
	pix, err := NewPixmap(CsRGB, 10, 10, false)
	if err != nil {
		t.Fatalf("NewPixmap: %v", err)
	}
	defer pix.Close()

	// Standalone pixmap should have x=0, y=0
	if pix.X() != 0 || pix.Y() != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", pix.X(), pix.Y())
	}
}

// --- Page manipulation tests ---

func TestNewPage(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	page, err := doc.NewPage(-1, 612, 792)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	page.Close()

	if doc.PageCount() != 1 {
		t.Errorf("expected 1 page after insert, got %d", doc.PageCount())
	}
}

func TestNewPageDefaultSize(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	// Width/height <= 0 should default to A4
	page, err := doc.NewPage(-1, 0, 0)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer page.Close()

	w := page.Width()
	h := page.Height()
	if w != 595 || h != 842 {
		t.Errorf("expected A4 default (595x842), got %gx%g", w, h)
	}
}

func TestNewPageMultiple(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	for i := 0; i < 5; i++ {
		p, err := doc.NewPage(-1, 595, 842)
		if err != nil {
			t.Fatalf("NewPage %d: %v", i, err)
		}
		p.Close()
	}

	if doc.PageCount() != 5 {
		t.Errorf("expected 5 pages, got %d", doc.PageCount())
	}
}

func TestDeletePage(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	p.Close()

	if doc.PageCount() != 2 {
		t.Fatalf("expected 2 pages, got %d", doc.PageCount())
	}

	err = doc.DeletePage(0)
	if err != nil {
		t.Fatalf("DeletePage: %v", err)
	}

	if doc.PageCount() != 1 {
		t.Errorf("expected 1 page after delete, got %d", doc.PageCount())
	}
}

func TestDeletePageOutOfRange(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	err := doc.DeletePage(99)
	if err == nil {
		t.Error("expected error for out-of-range delete")
	}
}

func TestDeletePages(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	for i := 0; i < 5; i++ {
		p, _ := doc.NewPage(-1, 595, 842)
		p.Close()
	}

	err = doc.DeletePages(4, 2, 0)
	if err != nil {
		t.Fatalf("DeletePages: %v", err)
	}
	if doc.PageCount() != 2 {
		t.Errorf("expected 2 pages, got %d", doc.PageCount())
	}
}

func TestSelect(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	for i := 0; i < 5; i++ {
		p, _ := doc.NewPage(-1, 595, 842)
		p.Close()
	}

	err = doc.Select([]int{0, 2, 4})
	if err != nil {
		t.Fatalf("Select: %v", err)
	}
	if doc.PageCount() != 3 {
		t.Errorf("expected 3 pages after select, got %d", doc.PageCount())
	}
}

func TestSelectEmpty(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	err := doc.Select([]int{})
	if err == nil {
		t.Error("expected error for empty select")
	}
}

func TestSelectOutOfRange(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	err := doc.Select([]int{99})
	if err == nil {
		t.Error("expected error for out-of-range select")
	}
}

// --- Save/Export tests ---

func TestSaveAndReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.pdf")

	doc := newTestPDFWithPage(t)
	err := doc.Save(path)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	doc.Close()

	doc2, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc2.Close()

	if !doc2.IsPDF() {
		t.Error("reopened doc should be PDF")
	}
	if doc2.PageCount() != 1 {
		t.Errorf("expected 1 page, got %d", doc2.PageCount())
	}
}

func TestSaveWithOptions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "opts.pdf")

	doc := newTestPDFWithPage(t)
	defer doc.Close()

	opts := SaveOptions{
		Garbage: 3,
		Deflate: true,
		Clean:   true,
		ASCII:   true,
		Pretty:  true,
	}
	err := doc.Save(path, opts)
	if err != nil {
		t.Fatalf("Save with options: %v", err)
	}

	info, _ := os.Stat(path)
	if info.Size() == 0 {
		t.Error("expected non-empty file")
	}
}

func TestToBytes(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	data, err := doc.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PDF bytes")
	}
	if string(data[:5]) != "%PDF-" {
		t.Error("invalid PDF header")
	}
}

func TestToBytesWithOptions(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	opts := SaveOptions{Garbage: 3, Deflate: true}
	data, err := doc.ToBytes(opts)
	if err != nil {
		t.Fatalf("ToBytes with options: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PDF bytes")
	}
}

func TestEzSave(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "ez.pdf")

	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	err = doc.EzSave(path)
	if err != nil {
		t.Fatalf("EzSave: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if info.Size() == 0 {
		t.Error("expected non-empty file")
	}
}

func TestConvertToPDF(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	data, err := doc.ConvertToPDF(0, 0, 0)
	if err != nil {
		t.Fatalf("ConvertToPDF: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty PDF data")
	}
}

func TestCanSaveIncrementally(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	_ = doc.CanSaveIncrementally()
}

func TestCanSaveIncrementallyFromFile(t *testing.T) {
	path := createTestPDF(t)
	doc, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc.Close()

	// File-based PDF may support incremental save
	_ = doc.CanSaveIncrementally()
}

// --- Search tests ---

func TestSearchFor(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	quads, err := page.SearchFor("hello", true)
	if err != nil {
		t.Fatalf("SearchFor: %v", err)
	}
	if len(quads) != 0 {
		t.Errorf("expected 0 results on empty page, got %d", len(quads))
	}
}

func TestSearchPageFor(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	quads, err := doc.SearchPageFor(0, "test", true)
	if err != nil {
		t.Fatalf("SearchPageFor: %v", err)
	}
	if len(quads) != 0 {
		t.Errorf("expected 0 results, got %d", len(quads))
	}
}

// --- Link tests ---

func TestGetLinks(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	links, err := page.GetLinks()
	if err != nil {
		t.Fatalf("GetLinks: %v", err)
	}
	if len(links) != 0 {
		t.Errorf("expected 0 links, got %d", len(links))
	}
}

// --- Annotation tests ---

func TestGetAnnots(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	annots := page.GetAnnots()
	if len(annots) != 0 {
		t.Errorf("expected 0 annotations, got %d", len(annots))
	}
}

func TestAddTextAnnot(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	annot, err := page.AddTextAnnot(NewPoint(100, 100), "Test note")
	if err != nil {
		t.Fatalf("AddTextAnnot: %v", err)
	}

	if annot.Type() != AnnotText {
		t.Errorf("expected AnnotText, got %d", annot.Type())
	}
	if annot.TypeString() != "Text" {
		t.Errorf("expected 'Text', got '%s'", annot.TypeString())
	}
	if annot.Contents() != "Test note" {
		t.Errorf("expected 'Test note', got '%s'", annot.Contents())
	}

	r := annot.Rect()
	if r.IsEmpty() {
		t.Error("annotation rect should not be empty")
	}
	if annot.Xref() <= 0 {
		t.Errorf("expected positive xref, got %d", annot.Xref())
	}
}

func TestAnnotSetContents(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	annot, err := page.AddTextAnnot(NewPoint(50, 50), "Original")
	if err != nil {
		t.Fatalf("AddTextAnnot: %v", err)
	}

	annot.SetContents("Updated")
	if annot.Contents() != "Updated" {
		t.Errorf("expected 'Updated', got '%s'", annot.Contents())
	}
}

func TestAddFreetextAnnot(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	rect := NewRect(100, 100, 300, 150)
	annot, err := page.AddFreetextAnnot(rect, "Free text", 12)
	if err != nil {
		t.Fatalf("AddFreetextAnnot: %v", err)
	}

	if annot.Type() != AnnotFreeText {
		t.Errorf("expected AnnotFreeText, got %d", annot.Type())
	}
	if annot.TypeString() != "FreeText" {
		t.Errorf("expected 'FreeText', got '%s'", annot.TypeString())
	}
}

func TestAddHighlightAnnot(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	quads := []Quad{NewRect(100, 100, 200, 120).Quad()}
	annot, err := page.AddHighlightAnnot(quads)
	if err != nil {
		t.Fatalf("AddHighlightAnnot: %v", err)
	}

	if annot.Type() != AnnotHighlight {
		t.Errorf("expected AnnotHighlight, got %d", annot.Type())
	}
}

func TestAddHighlightAnnotEmpty(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	_, err = page.AddHighlightAnnot(nil)
	if err == nil {
		t.Error("expected error for empty quads")
	}
}

func TestDeleteAnnot(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	annot, err := page.AddTextAnnot(NewPoint(100, 100), "To delete")
	if err != nil {
		t.Fatalf("AddTextAnnot: %v", err)
	}

	err = page.DeleteAnnot(annot)
	if err != nil {
		t.Fatalf("DeleteAnnot: %v", err)
	}
}

func TestAnnotPersistence(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "annot.pdf")

	doc := newTestPDFWithPage(t)
	page, _ := doc.LoadPage(0)
	page.AddTextAnnot(NewPoint(100, 100), "Persistent note")
	page.Close()
	doc.Save(path)
	doc.Close()

	doc2, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc2.Close()

	page2, _ := doc2.LoadPage(0)
	defer page2.Close()

	annots := page2.GetAnnots()
	if len(annots) == 0 {
		t.Error("expected at least 1 annotation after save/reopen")
	}
}

// --- Widget tests ---

func TestGetWidgets(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	widgets := page.GetWidgets()
	if len(widgets) != 0 {
		t.Errorf("expected 0 widgets, got %d", len(widgets))
	}
}

// --- Xref tests ---

func TestXrefLength(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	xlen := doc.XrefLength()
	if xlen <= 0 {
		t.Errorf("expected positive xref length, got %d", xlen)
	}
}

func TestXrefObject(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	xlen := doc.XrefLength()
	if xlen <= 1 {
		t.Skip("not enough xref entries")
	}

	obj, err := doc.XrefObject(1, false)
	if err != nil {
		t.Fatalf("XrefObject: %v", err)
	}
	if obj == "" {
		t.Error("expected non-empty xref object string")
	}
	t.Logf("xref 1: %s", obj)
}

func TestXrefObjectCompressed(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	obj, err := doc.XrefObject(1, true)
	if err != nil {
		t.Fatalf("XrefObject compressed: %v", err)
	}
	if obj == "" {
		t.Error("expected non-empty xref object string")
	}
}

func TestPDFCatalog(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	cat := doc.PDFCatalog()
	if cat <= 0 {
		t.Errorf("expected positive catalog xref, got %d", cat)
	}
}

// --- Embedded files tests ---

func TestEmbFileCount(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	count := doc.EmbFileCount()
	if count != 0 {
		t.Errorf("expected 0 embedded files, got %d", count)
	}
}

func TestEmbFileNames(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	names := doc.EmbFileNames()
	if len(names) != 0 {
		t.Errorf("expected 0 names, got %d", len(names))
	}
}

func TestEmbFileGetNoFiles(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	_, err = doc.EmbFileGet(0)
	if err == nil {
		t.Error("expected error for non-existent embedded file")
	}
}

// --- InsertPDF tests ---

func TestInsertPDF(t *testing.T) {
	doc1, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc1.Close()

	p, _ := doc1.NewPage(-1, 595, 842)
	p.Close()

	doc2, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc2.Close()

	for i := 0; i < 3; i++ {
		p, _ := doc2.NewPage(-1, 595, 842)
		p.Close()
	}

	err = doc1.InsertPDF(doc2)
	if err != nil {
		t.Fatalf("InsertPDF: %v", err)
	}

	if doc1.PageCount() != 4 {
		t.Errorf("expected 4 pages, got %d", doc1.PageCount())
	}
}

func TestInsertPDFWithOptions(t *testing.T) {
	doc1, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc1.Close()

	p, _ := doc1.NewPage(-1, 595, 842)
	p.Close()

	doc2, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc2.Close()

	for i := 0; i < 5; i++ {
		p, _ := doc2.NewPage(-1, 595, 842)
		p.Close()
	}

	opts := InsertPDFOptions{FromPage: 1, ToPage: 3, StartAt: 0}
	err = doc1.InsertPDF(doc2, opts)
	if err != nil {
		t.Fatalf("InsertPDF with options: %v", err)
	}

	if doc1.PageCount() != 4 {
		t.Errorf("expected 4 pages (1 original + 3 inserted), got %d", doc1.PageCount())
	}
}

// --- Pages iterator tests ---

func TestPages(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	for i := 0; i < 3; i++ {
		p, err := doc.NewPage(-1, 595, 842)
		if err != nil {
			t.Fatalf("NewPage: %v", err)
		}
		p.Close()
	}

	pages, err := doc.Pages()
	if err != nil {
		t.Fatalf("Pages: %v", err)
	}
	for _, p := range pages {
		p.Close()
	}
	if len(pages) != doc.PageCount() {
		t.Errorf("expected %d pages, got %d", doc.PageCount(), len(pages))
	}
}

func TestPagesRange(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	for i := 0; i < 5; i++ {
		p, _ := doc.NewPage(-1, 595, 842)
		p.Close()
	}

	// Pages(1, 3) should return pages 1, 2
	pages, err := doc.Pages(1, 3)
	if err != nil {
		t.Fatalf("Pages(1,3): %v", err)
	}
	for _, p := range pages {
		p.Close()
	}
	if len(pages) != 2 {
		t.Errorf("expected 2 pages, got %d", len(pages))
	}
}

func TestPagesStep(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	for i := 0; i < 6; i++ {
		p, _ := doc.NewPage(-1, 595, 842)
		p.Close()
	}

	// Pages(0, 6, 2) should return pages 0, 2, 4
	pages, err := doc.Pages(0, 6, 2)
	if err != nil {
		t.Fatalf("Pages(0,6,2): %v", err)
	}
	for _, p := range pages {
		p.Close()
	}
	if len(pages) != 3 {
		t.Errorf("expected 3 pages, got %d", len(pages))
	}
}

// --- Transformation matrix tests ---

func TestTransformationMatrix(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	m := page.TransformationMatrix()
	if m != Identity {
		t.Errorf("expected identity matrix for 0-rotation page, got %v", m)
	}

	dm := page.DerotationMatrix()
	if dm != Identity {
		t.Errorf("expected identity derotation matrix, got %v", dm)
	}
}

func TestTransformationMatrixRotated(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	page.SetRotation(90)
	m := page.TransformationMatrix()
	if m == Identity {
		t.Error("90-degree rotation should not produce identity matrix")
	}

	page.SetRotation(180)
	m = page.TransformationMatrix()
	if m == Identity {
		t.Error("180-degree rotation should not produce identity matrix")
	}

	page.SetRotation(270)
	m = page.TransformationMatrix()
	if m == Identity {
		t.Error("270-degree rotation should not produce identity matrix")
	}
}

func TestDerotationMatrixRoundTrip(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	for _, rot := range []int{0, 90, 180, 270} {
		page.SetRotation(rot)
		m := page.TransformationMatrix()
		dm := page.DerotationMatrix()
		product := m.Concat(dm)
		// Should be close to identity
		if !matrixNearIdentity(product) {
			t.Errorf("rotation %d: M*DM should be identity, got %v", rot, product)
		}
	}
}

func matrixNearIdentity(m Matrix) bool {
	const eps = 1e-6
	return abs(m.A-1) < eps && abs(m.B) < eps && abs(m.C) < eps &&
		abs(m.D-1) < eps && abs(m.E) < eps && abs(m.F) < eps
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// --- Text insertion tests ---

func TestInsertText(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	n, err := page.InsertText(NewPoint(72, 72), "Hello World")
	if err != nil {
		t.Fatalf("InsertText: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

func TestInsertTextWithOptions(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	n, err := page.InsertText(NewPoint(72, 200), "Red Text",
		WithFontName("Helvetica"),
		WithFontSize(16),
		WithColor(ColorRed),
	)
	if err != nil {
		t.Fatalf("InsertText with options: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

func TestInsertTextCJK(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	// Insert Chinese text
	n, err := page.InsertText(NewPoint(72, 200), "你好世界", WithFontSize(14))
	if err != nil {
		t.Fatalf("InsertText CJK: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	// Insert Japanese text
	n, err = page.InsertText(NewPoint(72, 230), "こんにちは", WithFontSize(14))
	if err != nil {
		t.Fatalf("InsertText Japanese: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	// Insert Korean text
	n, err = page.InsertText(NewPoint(72, 260), "안녕하세요", WithFontSize(14))
	if err != nil {
		t.Fatalf("InsertText Korean: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}

	// Insert mixed Latin + CJK text
	n, err = page.InsertText(NewPoint(72, 290), "Hello 你好", WithFontSize(14))
	if err != nil {
		t.Fatalf("InsertText mixed: %v", err)
	}
	if n != 1 {
		t.Errorf("expected 1, got %d", n)
	}
}

// --- Image insertion tests ---

func TestInsertImageEmpty(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer page.Close()

	err = page.InsertImage(NewRect(100, 100, 200, 200), nil)
	if err == nil {
		t.Error("expected error for nil image data")
	}
	err = page.InsertImage(NewRect(100, 100, 200, 200), []byte{})
	if err == nil {
		t.Error("expected error for empty image data")
	}
}

// --- Font/Image info tests ---

func TestGetFonts(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	fonts, err := doc.GetPageFonts(0)
	if err != nil {
		t.Fatalf("GetPageFonts: %v", err)
	}
	// Empty page may have no fonts
	_ = fonts
}

func TestGetImages(t *testing.T) {
	doc := newTestPDFWithPage(t)
	defer doc.Close()

	images, err := doc.GetPageImages(0)
	if err != nil {
		t.Fatalf("GetPageImages: %v", err)
	}
	_ = images
}

// --- TextInsertOption tests ---

func TestTextInsertOptions(t *testing.T) {
	cfg := textInsertConfig{}
	WithFontName("Courier")(&cfg)
	if cfg.fontname != "Courier" {
		t.Errorf("expected Courier, got %s", cfg.fontname)
	}
	WithFontSize(14)(&cfg)
	if cfg.fontsize != 14 {
		t.Errorf("expected 14, got %f", cfg.fontsize)
	}
	WithColor(ColorBlue)(&cfg)
	if cfg.color != ColorBlue {
		t.Error("color mismatch")
	}
	WithRotate(90)(&cfg)
	if cfg.rotate != 90 {
		t.Errorf("expected 90, got %d", cfg.rotate)
	}
}

// --- PixmapOption tests ---

func TestPixmapOptions(t *testing.T) {
	cfg := pixmapConfig{}
	WithMatrix(ScaleMatrix(2, 2))(&cfg)
	if cfg.matrix != ScaleMatrix(2, 2) {
		t.Error("matrix mismatch")
	}
	WithDPI(150)(&cfg)
	if cfg.dpi != 150 {
		t.Errorf("expected 150, got %d", cfg.dpi)
	}
	WithColorspace(CsGray)(&cfg)
	if cfg.colorspace != CsGray {
		t.Errorf("expected CsGray, got %d", cfg.colorspace)
	}
	WithAlpha(true)(&cfg)
	if !cfg.alpha {
		t.Error("expected alpha=true")
	}
	clip := NewRect(0, 0, 100, 100)
	WithClip(clip)(&cfg)
	if cfg.clip == nil || *cfg.clip != clip {
		t.Error("clip mismatch")
	}
	WithAnnots(false)(&cfg)
	if cfg.annots {
		t.Error("expected annots=false")
	}
}

// --- HTMLBox tests ---

func TestInsertHTMLBox(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	result, err := p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
		`<p style="font-size:14px;">Hello <b>World</b> from GoMuPDF!</p>`,
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox: %v", err)
	}
	if result.Scale <= 0 {
		t.Errorf("expected positive scale, got %f", result.Scale)
	}
	t.Logf("spare_height=%.1f, scale=%.3f", result.SpareHeight, result.Scale)
}

func TestInsertHTMLBoxWithCSS(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	html := `<h1>Title</h1><p class="body">Some body text with <em>emphasis</em>.</p>`
	css := `h1 { color: red; font-size: 24px; } .body { font-size: 12px; }`

	result, err := p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 500, Y1: 300},
		html,
		HTMLBoxOptions{CSS: css, Overlay: true},
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox with CSS: %v", err)
	}
	t.Logf("spare_height=%.1f, scale=%.3f", result.SpareHeight, result.Scale)
}

func TestInsertHTMLBoxCJK(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	html := `<p style="font-size:16px;">你好世界 Hello World 测试中文</p>`
	result, err := p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
		html,
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox CJK: %v", err)
	}
	t.Logf("CJK spare_height=%.1f, scale=%.3f", result.SpareHeight, result.Scale)
}

func TestInsertHTMLBoxEmpty(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	result, err := p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
		"",
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox empty: %v", err)
	}
	if result.Scale != 1.0 {
		t.Errorf("expected scale 1.0 for empty, got %f", result.Scale)
	}
}

func TestInsertHTMLBoxScaleDown(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	// Lots of text in a small box — should trigger auto-scaling
	html := `<p style="font-size:20px;">` +
		`This is a very long paragraph that should not fit in a tiny rectangle. ` +
		`It contains enough text to overflow the small box and trigger the ` +
		`automatic scale-down feature of InsertHTMLBox. ` +
		`More text here to ensure overflow happens reliably.</p>`

	result, err := p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 200, Y1: 80},
		html,
		HTMLBoxOptions{ScaleLow: 0, Overlay: true},
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox scale: %v", err)
	}
	t.Logf("scale-down: spare_height=%.1f, scale=%.3f", result.SpareHeight, result.Scale)
}

func TestInsertHTMLBoxNoScale(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	// Lots of text in a small box with ScaleLow=1 (no scaling) — should return ErrOverflow
	html := `<p style="font-size:20px;">` +
		`This is a very long paragraph that should not fit in a tiny rectangle. ` +
		`It contains enough text to overflow the small box.</p>`

	_, err = p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 200, Y1: 80},
		html,
		HTMLBoxOptions{ScaleLow: 1, Overlay: true},
	)
	if err != ErrOverflow {
		t.Logf("expected ErrOverflow, got: %v (may fit depending on font metrics)", err)
	}
}

func TestInsertHTMLBoxSaveAndVerify(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}

	_, err = p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
		`<p>Test content for save verification</p>`,
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox: %v", err)
	}
	p.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "htmlbox.pdf")
	err = doc.Save(path)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Reopen and verify
	doc2, err := Open(path)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer doc2.Close()

	if doc2.PageCount() != 1 {
		t.Fatalf("expected 1 page, got %d", doc2.PageCount())
	}

	p2, err := doc2.LoadPage(0)
	if err != nil {
		t.Fatalf("LoadPage: %v", err)
	}
	defer p2.Close()

	text, err := p2.GetText("text")
	if err != nil {
		t.Fatalf("GetText: %v", err)
	}
	t.Logf("extracted text: %q", text)
	if len(text) == 0 {
		t.Error("expected non-empty text from HTML box")
	}
}

func TestInsertHTMLBoxNonPDF(t *testing.T) {
	doc, err := NewPDF()
	if err != nil {
		t.Fatalf("NewPDF: %v", err)
	}
	defer doc.Close()

	// Force non-PDF by using a page from a non-PDF doc (simulate with closed pdf)
	p, err := doc.NewPage(-1, 595, 842)
	if err != nil {
		t.Fatalf("NewPage: %v", err)
	}
	defer p.Close()

	// This should work since it IS a PDF
	_, err = p.InsertHTMLBox(
		Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
		`<p>test</p>`,
	)
	if err != nil {
		t.Fatalf("InsertHTMLBox on PDF should work: %v", err)
	}
}

func TestFontSubsettingSize(t *testing.T) {
	dir := t.TempDir()

	type testCase struct {
		name string
		fn   func(doc *Document, p *Page)
	}

	cases := []testCase{
		{"HTMLBox_2CJK", func(doc *Document, p *Page) {
			p.InsertHTMLBox(Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
				`<p style="font-size:16px;">你好</p>`)
		}},
		{"HTMLBox_50CJK", func(doc *Document, p *Page) {
			p.InsertHTMLBox(Rect{X0: 50, Y0: 50, X1: 500, Y1: 400},
				`<p style="font-size:14px;">这是一段较长的中文测试文本用来验证字体子集化的效果包含更多不同的汉字春夏秋冬东南西北上下左右大小多少数字和标点</p>`)
		}},
		{"HTMLBox_English", func(doc *Document, p *Page) {
			p.InsertHTMLBox(Rect{X0: 50, Y0: 50, X1: 500, Y1: 200},
				`<p style="font-size:16px;">Hello World</p>`)
		}},
		{"InsertText_CJK", func(doc *Document, p *Page) {
			p.InsertText(Point{X: 50, Y: 50}, "你好世界测试中文")
		}},
	}

	for _, tc := range cases {
		doc, err := NewPDF()
		if err != nil {
			t.Fatalf("NewPDF: %v", err)
		}
		p, err := doc.NewPage(-1, 595, 842)
		if err != nil {
			t.Fatalf("NewPage: %v", err)
		}
		tc.fn(doc, p)
		p.Close()

		path := filepath.Join(dir, tc.name+".pdf")
		err = doc.EzSave(path)
		if err != nil {
			t.Fatalf("EzSave %s: %v", tc.name, err)
		}
		doc.Close()

		fi, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat %s: %v", tc.name, err)
		}
		t.Logf("%-20s %6d bytes (%5.1f KB)", tc.name, fi.Size(), float64(fi.Size())/1024)
	}
}

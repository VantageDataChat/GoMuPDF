//go:build cgo && !nomupdf

package gomupdf_test

import (
	"fmt"

	"github.com/nicejuice/gomupdf"
)

// Example_openDocument demonstrates opening a PDF and extracting text.
func Example_openDocument() {
	doc, err := gomupdf.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	fmt.Printf("Pages: %d\n", doc.PageCount())
	fmt.Printf("Is PDF: %v\n", doc.IsPDF())

	// Get metadata
	meta := doc.Metadata()
	fmt.Printf("Title: %s\n", meta["title"])
	fmt.Printf("Author: %s\n", meta["author"])

	// Extract text from each page
	for i := 0; i < doc.PageCount(); i++ {
		page, err := doc.LoadPage(i)
		if err != nil {
			continue
		}
		text, _ := page.GetText("text")
		fmt.Printf("Page %d text length: %d\n", i, len(text))
		page.Close()
	}
}

// Example_renderPage demonstrates rendering a page to an image.
func Example_renderPage() {
	doc, err := gomupdf.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		panic(err)
	}
	defer page.Close()

	// Render at 2x zoom
	pix, err := page.GetPixmap(
		gomupdf.WithMatrix(gomupdf.ScaleMatrix(2, 2)),
		gomupdf.WithAlpha(false),
	)
	if err != nil {
		panic(err)
	}
	defer pix.Close()

	// Save as PNG
	err = pix.SavePNG("page0.png")
	if err != nil {
		panic(err)
	}

	// Or get as Go image
	img := pix.ToImage()
	fmt.Printf("Image size: %dx%d\n", img.Bounds().Dx(), img.Bounds().Dy())
}

// Example_searchText demonstrates searching for text on a page.
func Example_searchText() {
	doc, err := gomupdf.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		panic(err)
	}
	defer page.Close()

	quads, err := page.SearchFor("hello", true)
	if err != nil {
		panic(err)
	}

	for _, q := range quads {
		fmt.Printf("Found at: %s\n", q.Rect())
	}
}

// Example_tableOfContents demonstrates extracting the TOC.
func Example_tableOfContents() {
	doc, err := gomupdf.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	toc, err := doc.GetTOC(true)
	if err != nil {
		panic(err)
	}

	for _, item := range toc {
		indent := ""
		for i := 1; i < item.Level; i++ {
			indent += "  "
		}
		fmt.Printf("%s%s (page %d)\n", indent, item.Title, item.Page)
	}
}

// Example_manipulatePDF demonstrates PDF manipulation.
func Example_manipulatePDF() {
	doc, err := gomupdf.Open("input.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	// Delete pages 5-10
	for i := 10; i >= 5; i-- {
		doc.DeletePage(i)
	}

	// Add a new blank page
	_, err = doc.NewPage(-1, 595, 842)
	if err != nil {
		panic(err)
	}

	// Save with compression
	err = doc.EzSave("output.pdf")
	if err != nil {
		panic(err)
	}
}

// Example_mergePDFs demonstrates merging two PDFs.
func Example_mergePDFs() {
	doc1, _ := gomupdf.Open("file1.pdf")
	defer doc1.Close()

	doc2, _ := gomupdf.Open("file2.pdf")
	defer doc2.Close()

	// Insert all pages of doc2 at the end of doc1
	doc1.InsertPDF(doc2)

	doc1.EzSave("merged.pdf")
}

// Example_annotations demonstrates working with annotations.
func Example_annotations() {
	doc, err := gomupdf.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		panic(err)
	}
	defer page.Close()

	// Add a text annotation
	_, err = page.AddTextAnnot(
		gomupdf.NewPoint(100, 100),
		"This is a note",
	)
	if err != nil {
		panic(err)
	}

	// List existing annotations
	annots := page.GetAnnots()
	for _, a := range annots {
		fmt.Printf("Annot: %s at %s\n", a.TypeString(), a.Rect())
	}

	doc.Save("annotated.pdf")
}

// Example_textExtraction demonstrates detailed text extraction.
func Example_textExtraction() {
	doc, err := gomupdf.Open("example.pdf")
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	page, err := doc.LoadPage(0)
	if err != nil {
		panic(err)
	}
	defer page.Close()

	// Get text blocks
	blocks, _ := page.GetTextBlocks()
	for _, b := range blocks {
		fmt.Printf("Block at %s: %s...\n", b.Rect, b.Text[:min(50, len(b.Text))])
	}

	// Get text words
	words, _ := page.GetTextWords()
	for _, w := range words {
		fmt.Printf("Word '%s' at %s\n", w.Text, w.Rect)
	}

	// Detailed text page analysis
	tp, _ := page.GetTextPage()
	defer tp.Close()
	for _, block := range tp.Blocks() {
		for _, line := range block.Lines {
			fmt.Printf("Line: %s\n", line.Text())
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

[中文](API_CN.md) | English

# GoMuPDF API Reference

Package `gomupdf` provides Go bindings for MuPDF.

All coordinates use a top-left origin system (same as PyMuPDF). Units are PDF points (1 point = 1/72 inch).

---

## Package Functions

### Document Creation & Opening

```go
func Open(filename string) (*Document, error)
```
Opens a document from a file path. Supports PDF, XPS, EPUB, CBZ, FB2, and image formats.

```go
func OpenFromMemory(data []byte, magic string) (*Document, error)
```
Opens a document from a byte slice. `magic` is a MIME type or file extension hint (e.g. `"application/pdf"`, `".pdf"`).

```go
func NewPDF() (*Document, error)
```
Creates a new empty PDF document.

### Pixmap Creation

```go
func NewPixmap(colorspace int, width, height int, alpha bool) (*Pixmap, error)
```
Creates a new blank pixmap. `colorspace`: `CsGray`, `CsRGB`, or `CsCMYK`.

```go
func NewPixmapFromImage(doc *Document, xref int) (*Pixmap, error)
```
Creates a pixmap from an image object in the document by xref number.

### Geometry Constructors

```go
func NewPoint(x, y float64) Point
func NewRect(x0, y0, x1, y1 float64) Rect
func RectFromPoints(topLeft, bottomRight Point) Rect
func NewIRect(x0, y0, x1, y1 int) IRect
func NewQuad(ul, ur, ll, lr Point) Quad
func QuadFromRect(r Rect) Quad
```

### Matrix Constructors

```go
func NewMatrix(a, b, c, d, e, f float64) Matrix
func ScaleMatrix(sx, sy float64) Matrix
func TranslateMatrix(tx, ty float64) Matrix
func RotateMatrix(deg float64) Matrix
func ShearMatrix(sx, sy float64) Matrix
```

### Utilities

```go
func Version() string          // Returns GoMuPDF version ("0.1.0")
func MuPDFVersion() string     // Returns MuPDF version ("1.24.9")
func PaperSize(name string) Rect  // Returns Rect for named paper size ("a4", "letter", "a3-landscape", etc.)
func GetPDFStr(s string) string   // Converts string to PDF string format (handles Unicode)
```

---

## Document

Represents an opened document (PDF, XPS, EPUB, etc.).

### Properties

```go
func (d *Document) PageCount() int
func (d *Document) Name() string
func (d *Document) IsPDF() bool
func (d *Document) IsClosed() bool
func (d *Document) NeedsPass() bool
func (d *Document) IsReflowable() bool
```

### Lifecycle

```go
func (d *Document) Close()
func (d *Document) Authenticate(password string) (int, error)
```

### Page Access

```go
func (d *Document) LoadPage(pageNum int) (*Page, error)    // 0-based, supports negative index
func (d *Document) Pages(args ...int) ([]*Page, error)     // args: start, stop, step
```

### Metadata

```go
func (d *Document) Metadata() map[string]string
func (d *Document) SetMetadata(meta map[string]string) error
```
Keys: `"title"`, `"author"`, `"subject"`, `"keywords"`, `"creator"`, `"producer"`, `"creationDate"`, `"modDate"`.

### Table of Contents

```go
func (d *Document) GetTOC(simple bool) ([]TOCItem, error)
```

### Save & Export

```go
func (d *Document) Save(filename string, opts ...SaveOptions) error
func (d *Document) EzSave(filename string) error              // Garbage=3, Deflate=true
func (d *Document) ToBytes(opts ...SaveOptions) ([]byte, error)
func (d *Document) ConvertToPDF(fromPage, toPage, rotate int) ([]byte, error)
func (d *Document) CanSaveIncrementally() bool
```

Font subsetting is performed automatically on save — only glyphs actually used are embedded.

### Page Manipulation

```go
func (d *Document) NewPage(pno int, width, height float64) (*Page, error)  // pno=-1 appends
func (d *Document) DeletePage(pno int) error
func (d *Document) DeletePages(pages ...int) error
func (d *Document) Select(pages []int) error
func (d *Document) InsertPDF(src *Document, opts ...InsertPDFOptions) error
```

### Layout (Reflowable Documents)

```go
func (d *Document) Layout(width, height, fontsize float64)
```

### Xref Access

```go
func (d *Document) XrefLength() int
func (d *Document) XrefObject(xref int, compressed bool) (string, error)
func (d *Document) PDFCatalog() int
```

### Embedded Files

```go
func (d *Document) EmbFileCount() int
func (d *Document) EmbFileNames() []string
func (d *Document) EmbFileGet(index int) ([]byte, error)
```

### Convenience Methods

```go
func (d *Document) GetPageText(pno int, output string) (string, error)
func (d *Document) GetPagePixmap(pno int, opts ...PixmapOption) (*Pixmap, error)
func (d *Document) SearchPageFor(pno int, needle string, quads bool) ([]Quad, error)
func (d *Document) GetPageFonts(pno int) ([]FontInfo, error)
func (d *Document) GetPageImages(pno int) ([]ImageInfo, error)
```

---

## Page

Represents a single page in a document.

### Properties

```go
func (p *Page) Number() int
func (p *Page) Rect() Rect
func (p *Page) Width() float64
func (p *Page) Height() float64
func (p *Page) MediaBox() Rect
func (p *Page) CropBox() Rect
func (p *Page) Rotation() int
func (p *Page) SetRotation(rotation int) error
func (p *Page) GetLabel() string
```

### Lifecycle

```go
func (p *Page) Close()
```

### Text Extraction

```go
func (p *Page) GetText(output string, flags ...int) (string, error)
```
`output`: `"text"` (plain text). `flags`: combination of `TextPreserveLigatures`, `TextPreserveWhitespace`, etc.

```go
func (p *Page) GetTextWords(flags ...int) ([]TextWord, error)
func (p *Page) GetTextBlocks(flags ...int) ([]TextBlock, error)
func (p *Page) GetTextPage(flags ...int) (*TextPage, error)
```

### Search

```go
func (p *Page) SearchFor(needle string, quads bool) ([]Quad, error)
```

### Rendering

```go
func (p *Page) GetPixmap(opts ...PixmapOption) (*Pixmap, error)
```

Pixmap options (functional options pattern):
```go
WithDPI(dpi int)
WithMatrix(m Matrix)
WithColorspace(cs int)       // CsGray, CsRGB, CsCMYK
WithAlpha(alpha bool)
WithClip(clip Rect)
WithAnnots(annots bool)
```

### Content Insertion

```go
func (p *Page) InsertText(pos Point, text string, opts ...TextInsertOption) (int, error)
```
Inserts text at position. CJK text is auto-detected and uses non-embedded CID fonts.

Text options:
```go
WithFontName(name string)    // Default: "Helvetica"
WithFontSize(size float64)   // Default: 11
WithColor(color Color)       // Default: ColorBlack
WithRotate(angle int)
```

```go
func (p *Page) InsertImage(rect Rect, imageData []byte, opts ...InsertImageOptions) error
```

```go
func (p *Page) InsertHTMLBox(rect Rect, html string, opts ...HTMLBoxOptions) (HTMLBoxResult, error)
```
Inserts styled HTML/CSS content using MuPDF's Story API. Supports CJK, auto scale-down, font subsetting.

### Links & Annotations

```go
func (p *Page) GetLinks() ([]Link, error)
func (p *Page) GetAnnots() []*Annot
func (p *Page) AddTextAnnot(pos Point, text string) (*Annot, error)
func (p *Page) AddFreetextAnnot(rect Rect, text string, fontsize float64) (*Annot, error)
func (p *Page) AddHighlightAnnot(quads []Quad) (*Annot, error)
func (p *Page) DeleteAnnot(annot *Annot) error
```

### Widgets (Form Fields)

```go
func (p *Page) GetWidgets() []*Widget
```

### Fonts & Images Info

```go
func (p *Page) GetFonts() ([]FontInfo, error)
func (p *Page) GetImages() ([]ImageInfo, error)
```

### Transformation

```go
func (p *Page) TransformationMatrix() Matrix
func (p *Page) DerotationMatrix() Matrix
```

---

## Pixmap

Represents a raster image (pixel map).

### Properties

```go
func (px *Pixmap) Width() int
func (px *Pixmap) Height() int
func (px *Pixmap) N() int          // Number of components (including alpha)
func (px *Pixmap) Alpha() int      // 1 if has alpha, 0 otherwise
func (px *Pixmap) Stride() int
func (px *Pixmap) X() int
func (px *Pixmap) Y() int
func (px *Pixmap) IRect() IRect
func (px *Pixmap) Samples() []byte
```

### Lifecycle

```go
func (px *Pixmap) Close()
```

### Export

```go
func (px *Pixmap) ToBytes() ([]byte, error)       // PNG bytes
func (px *Pixmap) Save(filename string) error      // Save as PNG
func (px *Pixmap) SavePNG(filename string) error   // Alias for Save
func (px *Pixmap) SavePNM(filename string) error
func (px *Pixmap) ToImage() image.Image            // Convert to Go image.Image
```

### Pixel Operations

```go
func (px *Pixmap) SetPixel(x, y int, c []byte)
func (px *Pixmap) GetPixel(x, y int) []byte
func (px *Pixmap) Clear(value int)
func (px *Pixmap) Invert()
func (px *Pixmap) Gamma(gamma float64)
func (px *Pixmap) Tint(black, white int)
func (px *Pixmap) Convert(colorspace int) (*Pixmap, error)
```

---

## TextPage

Provides detailed structured text analysis.

```go
func (t *TextPage) Close()
func (t *TextPage) ExtractText() (string, error)
func (t *TextPage) Blocks() []STextBlock
```

---

## Annot

Represents a PDF annotation.

```go
func (a *Annot) Type() int
func (a *Annot) TypeString() string
func (a *Annot) Rect() Rect
func (a *Annot) Contents() string
func (a *Annot) SetContents(text string)
func (a *Annot) Xref() int
```

---

## Widget

Represents a PDF form field.

```go
func (w *Widget) FieldType() int
func (w *Widget) FieldTypeString() string
func (w *Widget) FieldName() string
func (w *Widget) FieldValue() string
func (w *Widget) SetFieldValue(value string) error
func (w *Widget) Rect() Rect
func (w *Widget) Xref() int
```

---

## Geometry Types

### Point

```go
type Point struct { X, Y float64 }

func (p Point) Add(other Point) Point
func (p Point) Sub(other Point) Point
func (p Point) Mul(factor float64) Point
func (p Point) Abs() float64
func (p Point) Transform(m Matrix) Point
func (p Point) IsZero() bool
```

### Rect

```go
type Rect struct { X0, Y0, X1, Y1 float64 }

func (r Rect) Width() float64
func (r Rect) Height() float64
func (r Rect) IsEmpty() bool
func (r Rect) Contains(p Point) bool
func (r Rect) ContainsRect(other Rect) bool
func (r Rect) Intersects(other Rect) bool
func (r Rect) Intersect(other Rect) Rect
func (r Rect) Union(other Rect) Rect
func (r Rect) IncludePoint(p Point) Rect
func (r Rect) Transform(m Matrix) Rect
func (r Rect) Normalize() Rect
func (r Rect) TopLeft() Point
func (r Rect) TopRight() Point
func (r Rect) BottomLeft() Point
func (r Rect) BottomRight() Point
func (r Rect) Quad() Quad
func (r Rect) IRect() IRect
```

### IRect

```go
type IRect struct { X0, Y0, X1, Y1 int }

func (r IRect) Width() int
func (r IRect) Height() int
func (r IRect) IsEmpty() bool
func (r IRect) Rect() Rect
```

### Quad

```go
type Quad struct { UL, UR, LL, LR Point }

func (q Quad) Rect() Rect
func (q Quad) IsEmpty() bool
func (q Quad) IsRectangular() bool
func (q Quad) IsConvex() bool
func (q Quad) Transform(m Matrix) Quad
```

### Matrix

```go
type Matrix struct { A, B, C, D, E, F float64 }

var Identity Matrix

func (m Matrix) Concat(other Matrix) Matrix
func (m Matrix) PreScale(sx, sy float64) Matrix
func (m Matrix) PreTranslate(tx, ty float64) Matrix
func (m Matrix) PreRotate(deg float64) Matrix
func (m Matrix) Invert() (Matrix, bool)
func (m Matrix) IsRectilinear() bool
```

---

## Option Types

### SaveOptions

```go
type SaveOptions struct {
    Garbage     int    // 0-4, garbage collection level
    Deflate     bool   // compress streams
    Clean       bool   // clean content streams
    ASCII       bool   // ASCII hex encode binary data
    Linear      bool   // linearize (web-optimized)
    Pretty      bool   // pretty-print objects
    Incremental bool   // incremental save
    NoNewID     bool   // don't generate new file ID
    Encryption  int    // EncryptNone, EncryptAESV3, etc.
    Permissions int    // PDF permission flags
    OwnerPW     string
    UserPW      string
}

func DefaultSaveOptions() SaveOptions
func EzSaveOptions() SaveOptions    // Garbage=3, Deflate=true
```

### HTMLBoxOptions

```go
type HTMLBoxOptions struct {
    CSS      string   // additional CSS
    ScaleLow float64  // min scale (0=auto-shrink, 1=no scaling)
    Overlay  bool     // true=foreground, false=background
}
```

### HTMLBoxResult

```go
type HTMLBoxResult struct {
    SpareHeight float64  // remaining height after content
    Scale       float64  // actual scale factor used
}
```

### InsertPDFOptions

```go
type InsertPDFOptions struct {
    FromPage int   // source start page (default: 0)
    ToPage   int   // source end page (default: last)
    StartAt  int   // destination position (default: end)
    Rotate   int
    Links    bool
    Annots   bool
}
```

### InsertImageOptions

```go
type InsertImageOptions struct {
    KeepProportion bool  // maintain aspect ratio
    Overlay        bool  // foreground or background
}
```

---

## Constants

### Colorspaces

```go
CsGray = 0
CsRGB  = 1
CsCMYK = 2
```

### Annotation Types

`AnnotText`, `AnnotLink`, `AnnotFreeText`, `AnnotLine`, `AnnotSquare`, `AnnotCircle`, `AnnotHighlight`, `AnnotUnderline`, `AnnotStrikeOut`, `AnnotRedact`, `AnnotStamp`, `AnnotInk`, etc.

### Widget Types

`WidgetTypeButton`, `WidgetTypeCheckbox`, `WidgetTypeCombobox`, `WidgetTypeListbox`, `WidgetTypeRadioButton`, `WidgetTypeSignature`, `WidgetTypeText`.

### Encryption

`EncryptNone`, `EncryptRC4V1` (40-bit), `EncryptRC4V2` (128-bit), `EncryptAESV2` (128-bit), `EncryptAESV3` (256-bit).

### Paper Sizes

`PaperA4`, `PaperA3`, `PaperA5`, `PaperLetter`, `PaperLegal` — predefined `Rect` values.

### Predefined Colors

`ColorBlack`, `ColorWhite`, `ColorRed`, `ColorGreen`, `ColorBlue`, `ColorYellow`, `ColorMagenta`, `ColorCyan`.

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrInitFailed` | MuPDF context initialization failed |
| `ErrOpenFailed` | Document cannot be opened |
| `ErrPageNotFound` | Page number out of range |
| `ErrNotPDF` | PDF-only operation on non-PDF |
| `ErrEncrypted` | Document encrypted and not authenticated |
| `ErrAuthFailed` | Password authentication failed |
| `ErrClosed` | Operation on closed document |
| `ErrTextExtract` | Text extraction failed |
| `ErrPixmap` | Pixmap operation failed |
| `ErrSave` | Save operation failed |
| `ErrInvalidArg` | Invalid argument |
| `ErrOutline` | TOC/outline operation failed |
| `ErrSearch` | Search failed |
| `ErrConvert` | Document conversion failed |
| `ErrEmbeddedFile` | Embedded file operation failed |
| `ErrXref` | Xref operation failed |
| `ErrOverflow` | Content does not fit in target rectangle |

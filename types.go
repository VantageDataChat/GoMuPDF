package gomupdf

// TOCItem represents a table of contents entry.
type TOCItem struct {
	Level int
	Title string
	Page  int // 1-based page number, -1 if no destination
	Dest  *LinkDest
}

// LinkDest contains link destination details.
type LinkDest struct {
	Kind      int
	Page      int
	To        Point
	Zoom      float64
	URI       string
	File      string
	NamedDest string
}

// Link represents a hyperlink on a page.
type Link struct {
	Rect Rect
	URI  string
	Kind int
	Page int
	To   Point
}

// FontInfo contains information about a font referenced by a page.
type FontInfo struct {
	Xref     int
	Ext      string
	Type     string
	BaseName string
	Name     string
	Encoding string
}

// ImageInfo contains information about an image referenced by a page.
type ImageInfo struct {
	Xref       int
	SMask      int
	Width      int
	Height     int
	BPC        int
	Colorspace string
	AltCS      string
	Name       string
	Filter     string
}

// TextWord represents a word with its bounding box.
type TextWord struct {
	Rect    Rect
	Text    string
	BlockNo int
	LineNo  int
	WordNo  int
}

// TextBlock represents a text block with its bounding box.
type TextBlock struct {
	Rect    Rect
	Text    string
	Type    string // "text" or "image"
	BlockNo int
}

// Color represents an RGB color with values 0.0-1.0.
type Color struct {
	R, G, B float64
}

// EmbFileInfo contains metadata about an embedded file.
type EmbFileInfo struct {
	Name         string
	Filename     string
	UFilename    string
	Description  string
	Size         int
	Length       int
	CreationDate string
	ModDate      string
}

// SaveOptions configures how a PDF is saved.
type SaveOptions struct {
	Garbage     int
	Deflate     bool
	Clean       bool
	ASCII       bool
	Linear      bool
	Pretty      bool
	Incremental bool
	NoNewID     bool
	Encryption  int
	Permissions int
	OwnerPW     string
	UserPW      string
}

// InsertPDFOptions configures page insertion from another PDF.
type InsertPDFOptions struct {
	FromPage int
	ToPage   int
	StartAt  int
	Rotate   int
	Links    bool
	Annots   bool
}

// DefaultSaveOptions returns default save options.
func DefaultSaveOptions() SaveOptions {
	return SaveOptions{
		Permissions: -1,
		Encryption:  EncryptNone,
	}
}

// EzSaveOptions returns save options optimized for small file size.
// Equivalent to PyMuPDF's ez_save().
func EzSaveOptions() SaveOptions {
	return SaveOptions{
		Garbage:     3,
		Deflate:     true,
		Permissions: -1,
		Encryption:  EncryptNone,
	}
}

// InsertImageOptions configures image insertion.
type InsertImageOptions struct {
	KeepProportion bool
	Overlay        bool
}

// STextBlockType indicates the type of a structured text block.
type STextBlockType int

const (
	STextBlockText  STextBlockType = 0
	STextBlockImage STextBlockType = 1
)

// STextBlock represents a block of text or an image.
type STextBlock struct {
	Type  STextBlockType
	Rect  Rect
	Lines []STextLine
}

// STextLine represents a line of text.
type STextLine struct {
	Rect  Rect
	Dir   Point
	Chars []STextChar
}

// Text returns the text content of the line.
func (l STextLine) Text() string {
	s := ""
	for _, ch := range l.Chars {
		s += string(ch.C)
	}
	return s
}

// STextChar represents a single character.
type STextChar struct {
	C      rune
	Origin Point
	Rect   Rect
	Size   float64
	Font   string
}

// PixmapOption configures pixmap creation.
type PixmapOption func(*pixmapConfig)

type pixmapConfig struct {
	matrix     Matrix
	dpi        int
	colorspace int
	alpha      bool
	clip       *Rect
	annots     bool
}

// WithMatrix sets the transformation matrix for pixmap creation.
func WithMatrix(m Matrix) PixmapOption {
	return func(c *pixmapConfig) { c.matrix = m }
}

// WithDPI sets the DPI for pixmap creation.
func WithDPI(dpi int) PixmapOption {
	return func(c *pixmapConfig) { c.dpi = dpi }
}

// WithColorspace sets the colorspace (CsGray, CsRGB, CsCMYK).
func WithColorspace(cs int) PixmapOption {
	return func(c *pixmapConfig) { c.colorspace = cs }
}

// WithAlpha enables/disables alpha channel.
func WithAlpha(alpha bool) PixmapOption {
	return func(c *pixmapConfig) { c.alpha = alpha }
}

// WithClip sets a clip rectangle for pixmap creation.
func WithClip(clip Rect) PixmapOption {
	return func(c *pixmapConfig) { c.clip = &clip }
}

// WithAnnots includes/excludes annotations in the pixmap.
func WithAnnots(annots bool) PixmapOption {
	return func(c *pixmapConfig) { c.annots = annots }
}

// TextInsertOption configures text insertion.
type TextInsertOption func(*textInsertConfig)

type textInsertConfig struct {
	fontname string
	fontsize float64
	color    Color
	rotate   int
}

// WithFontName sets the font name for text insertion.
func WithFontName(name string) TextInsertOption {
	return func(c *textInsertConfig) { c.fontname = name }
}

// WithFontSize sets the font size for text insertion.
func WithFontSize(size float64) TextInsertOption {
	return func(c *textInsertConfig) { c.fontsize = size }
}

// WithColor sets the text color.
func WithColor(color Color) TextInsertOption {
	return func(c *textInsertConfig) { c.color = color }
}

// WithRotate sets the text rotation angle.
func WithRotate(angle int) TextInsertOption {
	return func(c *textInsertConfig) { c.rotate = angle }
}

// HTMLBoxOptions configures HTML box insertion.
type HTMLBoxOptions struct {
	CSS      string  // additional CSS styling
	ScaleLow float64 // minimum scale factor (0 = auto-shrink to fit, 1 = no scaling)
	Overlay  bool    // true = foreground, false = background
}

// HTMLBoxResult contains the result of an HTML box insertion.
type HTMLBoxResult struct {
	SpareHeight float64 // remaining height in the rect after content
	Scale       float64 // actual scale factor used (1.0 = no scaling)
}

// Predefined colors.
var (
	ColorBlack   = Color{0, 0, 0}
	ColorWhite   = Color{1, 1, 1}
	ColorRed     = Color{1, 0, 0}
	ColorGreen   = Color{0, 1, 0}
	ColorBlue    = Color{0, 0, 1}
	ColorYellow  = Color{1, 1, 0}
	ColorMagenta = Color{1, 0, 1}
	ColorCyan    = Color{0, 1, 1}
)

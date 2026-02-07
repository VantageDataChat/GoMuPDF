package gomupdf

// Link destination kinds — corresponds to PyMuPDF link kinds.
const (
	LinkNone  = 0 // No destination
	LinkGoto  = 1 // Go to a page in this document
	LinkURI   = 2 // Open a URI
	LinkGoToR = 3 // Go to a page in another document
	LinkNamed = 4 // Named action
	LinkLaunch = 5 // Launch an application
)

// Text extraction flags — corresponds to PyMuPDF TEXT_* flags.
const (
	TextPreserveLigatures  = 1 << 0
	TextPreserveWhitespace = 1 << 1
	TextPreserveImages     = 1 << 2
	TextInhibitSpaces      = 1 << 3
	TextDeHyphenate        = 1 << 4
	TextPreserveSpans      = 1 << 5
	TextMediaboxClip       = 1 << 6
	TextCIDForUnknownUnicode = 1 << 7
)

// Default text flags (matches PyMuPDF default flags=3).
const TextFlagsDefault = TextPreserveLigatures | TextPreserveWhitespace

// Annotation types — corresponds to PyMuPDF PDF_ANNOT_* constants.
const (
	AnnotText           = 0
	AnnotLink           = 1
	AnnotFreeText       = 2
	AnnotLine           = 3
	AnnotSquare         = 4
	AnnotCircle         = 5
	AnnotPolygon        = 6
	AnnotPolyLine       = 7
	AnnotHighlight      = 8
	AnnotUnderline      = 9
	AnnotSquiggly       = 10
	AnnotStrikeOut      = 11
	AnnotRedact         = 12
	AnnotStamp          = 13
	AnnotCaret          = 14
	AnnotInk            = 15
	AnnotPopup          = 16
	AnnotFileAttachment = 17
	AnnotSound          = 18
	AnnotMovie          = 19
	AnnotRichMedia      = 20
	AnnotWidget         = 21
	AnnotScreen         = 22
	AnnotPrinterMark    = 23
	AnnotTrapNet        = 24
	AnnotWatermark      = 25
	Annot3D             = 26
	AnnotProjection     = 27
	AnnotUnknown        = -1
)

// Widget field types — corresponds to PyMuPDF PDF_WIDGET_TYPE_* constants.
const (
	WidgetTypeUnknown    = 0
	WidgetTypeButton     = 1
	WidgetTypeCheckbox   = 2
	WidgetTypeCombobox   = 3
	WidgetTypeListbox    = 4
	WidgetTypeRadioButton = 5
	WidgetTypeSignature  = 6
	WidgetTypeText       = 7
)

// PDF permissions — corresponds to PyMuPDF PDF_PERM_* constants.
const (
	PermPrint           = 1 << 2  // Print the document
	PermModify          = 1 << 3  // Modify document contents
	PermCopy            = 1 << 4  // Copy or extract text/graphics
	PermAnnotate        = 1 << 5  // Add or modify annotations
	PermForm            = 1 << 8  // Fill in forms
	PermAccessibility   = 1 << 9  // Extract for accessibility
	PermAssemble        = 1 << 10 // Assemble (insert, rotate, delete pages)
	PermPrintHQ         = 1 << 11 // High quality print
)

// PDF encryption methods.
const (
	EncryptNone    = 0
	EncryptRC4V1   = 1 // RC4, 40-bit
	EncryptRC4V2   = 2 // RC4, 128-bit
	EncryptAESV2   = 3 // AES, 128-bit
	EncryptAESV3   = 4 // AES, 256-bit
	EncryptKeep    = -1 // Keep existing encryption
)

// Colorspace identifiers.
const (
	CsGray = iota
	CsRGB
	CsCMYK
)

// Standard page sizes in points (1 point = 1/72 inch).
var (
	PaperA4     = Rect{X0: 0, Y0: 0, X1: 595, Y1: 842}
	PaperA3     = Rect{X0: 0, Y0: 0, X1: 842, Y1: 1191}
	PaperA5     = Rect{X0: 0, Y0: 0, X1: 420, Y1: 595}
	PaperLetter = Rect{X0: 0, Y0: 0, X1: 612, Y1: 792}
	PaperLegal  = Rect{X0: 0, Y0: 0, X1: 612, Y1: 1008}
)

// OC (Optional Content) visibility actions.
const (
	OCPDF_OC_ON     = 0
	OCPDF_OC_TOGGLE = 1
	OCPDF_OC_OFF    = 2
)

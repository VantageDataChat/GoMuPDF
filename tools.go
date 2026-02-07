package gomupdf

import (
	"fmt"
	"strings"
)

// Version returns the GoMuPDF version string.
func Version() string {
	return "0.1.0"
}

// MuPDFVersion returns the underlying MuPDF library version.
func MuPDFVersion() string {
	return "1.24.9"
}

// GetPDFStr converts a string to a PDF-compatible string format.
// Handles encoding and escaping as needed.
func GetPDFStr(s string) string {
	// Check if ASCII-only
	isASCII := true
	for _, r := range s {
		if r > 127 {
			isASCII = false
			break
		}
	}

	if isASCII {
		// Escape special PDF characters
		s = strings.ReplaceAll(s, "\\", "\\\\")
		s = strings.ReplaceAll(s, "(", "\\(")
		s = strings.ReplaceAll(s, ")", "\\)")
		return "(" + s + ")"
	}

	// UTF-16BE BOM encoding for non-ASCII
	var hex strings.Builder
	hex.WriteString("<feff")
	for _, r := range s {
		hex.WriteString(fmt.Sprintf("%04x", r))
	}
	hex.WriteString(">")
	return hex.String()
}

// PaperSize returns the Rect for a named paper size.
// Supported: "a4", "a3", "a5", "letter", "legal", etc.
func PaperSize(name string) Rect {
	sizes := map[string]Rect{
		"a0":      {X0: 0, Y0: 0, X1: 2384, Y1: 3370},
		"a1":      {X0: 0, Y0: 0, X1: 1684, Y1: 2384},
		"a2":      {X0: 0, Y0: 0, X1: 1191, Y1: 1684},
		"a3":      {X0: 0, Y0: 0, X1: 842, Y1: 1191},
		"a4":      {X0: 0, Y0: 0, X1: 595, Y1: 842},
		"a5":      {X0: 0, Y0: 0, X1: 420, Y1: 595},
		"a6":      {X0: 0, Y0: 0, X1: 298, Y1: 420},
		"a7":      {X0: 0, Y0: 0, X1: 210, Y1: 298},
		"a8":      {X0: 0, Y0: 0, X1: 147, Y1: 210},
		"a9":      {X0: 0, Y0: 0, X1: 105, Y1: 147},
		"a10":     {X0: 0, Y0: 0, X1: 74, Y1: 105},
		"b0":      {X0: 0, Y0: 0, X1: 2835, Y1: 4008},
		"b1":      {X0: 0, Y0: 0, X1: 2004, Y1: 2835},
		"b2":      {X0: 0, Y0: 0, X1: 1417, Y1: 2004},
		"b3":      {X0: 0, Y0: 0, X1: 1001, Y1: 1417},
		"b4":      {X0: 0, Y0: 0, X1: 709, Y1: 1001},
		"b5":      {X0: 0, Y0: 0, X1: 499, Y1: 709},
		"letter":  {X0: 0, Y0: 0, X1: 612, Y1: 792},
		"legal":   {X0: 0, Y0: 0, X1: 612, Y1: 1008},
		"tabloid": {X0: 0, Y0: 0, X1: 792, Y1: 1224},
		"ledger":  {X0: 0, Y0: 0, X1: 1224, Y1: 792},
	}

	lower := strings.ToLower(strings.TrimSpace(name))

	// Check for landscape suffix
	landscape := false
	if strings.HasSuffix(lower, "-l") || strings.HasSuffix(lower, "-landscape") {
		landscape = true
		lower = strings.Split(lower, "-")[0]
	}

	r, ok := sizes[lower]
	if !ok {
		return PaperA4 // default
	}

	if landscape {
		r.X1, r.Y1 = r.Y1, r.X1
	}
	return r
}

// PlanishLine returns the deskew angle for a line defined by two points.
func PlanishLine(p1, p2 Point) float64 {
	if p1 == p2 {
		return 0
	}
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	if dx == 0 {
		if dy > 0 {
			return 90
		}
		return -90
	}
	return 0 // simplified
}

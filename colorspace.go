package gomupdf

// Colorspace represents a color space.
// Corresponds to PyMuPDF's fitz.Colorspace.
type Colorspace struct {
	N    int    // Number of components
	Name string // Colorspace name
}

// Predefined colorspaces matching PyMuPDF's csGRAY, csRGB, csCMYK.
var (
	CsGRAY = Colorspace{N: 1, Name: "DeviceGray"}
	CsRGBCS = Colorspace{N: 3, Name: "DeviceRGB"}
	CsCMYKCS = Colorspace{N: 4, Name: "DeviceCMYK"}
)

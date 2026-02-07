# GoMuPDF

Go bindings for MuPDF — a high-performance library for PDF (and other document) data extraction, analysis, conversion & manipulation.

This is a Go port of [PyMuPDF](https://github.com/pymupdf/PyMuPDF), providing an idiomatic Go API that mirrors PyMuPDF's functionality.

## Supported Platforms

| OS      | amd64 | arm64 |
|---------|-------|-------|
| Windows | ✅    | ✅    |
| Linux   | ✅    | ✅    |
| macOS   | ✅    | ✅    |

A `nomupdf` build tag is available for compiling without CGO/MuPDF (stub functions return `ErrInitFailed`).

## Supported Document Formats

PDF, XPS, OpenXPS, CBZ, CBR, FB2, EPUB, and image formats (JPEG, PNG, BMP, GIF, TIFF, etc.)

## Prerequisites

- Go 1.21+
- C compiler: GCC (MinGW-W64 on Windows) or Clang
- MuPDF 1.24.9 static libraries (see Building below)

## Building MuPDF Libraries

GoMuPDF links against platform-specific static libraries stored in the `libs/` directory. The naming convention is:

```
libs/libmupdf_<os>_<arch>.a
libs/libmupdfthird_<os>_<arch>.a
```

For example: `libmupdf_linux_amd64.a`, `libmupdf_darwin_arm64.a`, `libmupdf_windows_amd64.a`.

### Automated Build

```bash
# Linux / macOS
chmod +x build_libs.sh
./build_libs.sh

# Or specify an existing MuPDF source directory:
./build_libs.sh /path/to/mupdf-src

# Windows (MinGW)
build_libs.bat
build_libs.bat D:\mupdf-src
```

### Manual Build

```bash
git clone --depth 1 --branch 1.24.9 --recurse-submodules \
    https://github.com/ArtifexSoftware/mupdf.git mupdf-src
cd mupdf-src
make -j$(nproc) HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no \
    USE_SYSTEM_LIBS=no XCFLAGS="-fPIC" libs
cd ..

# Copy to libs/ with the correct naming
mkdir -p libs
cp mupdf-src/build/release/libmupdf.a libs/libmupdf_linux_amd64.a
cp mupdf-src/build/release/libmupdf-third.a libs/libmupdfthird_linux_amd64.a
```

### Using Pre-built Libraries

You can also extract the static libraries from [go-fitz](https://github.com/gen2brain/go-fitz) v1.24.15 and rename them to match the naming convention above.

## Installation

```bash
go get github.com/nicejuice/gomupdf
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/nicejuice/gomupdf"
)

func main() {
    doc, err := gomupdf.Open("example.pdf")
    if err != nil {
        panic(err)
    }
    defer doc.Close()

    fmt.Printf("Pages: %d\n", doc.PageCount())

    for i := 0; i < doc.PageCount(); i++ {
        page, err := doc.LoadPage(i)
        if err != nil {
            continue
        }
        text, _ := page.GetText("text")
        fmt.Printf("Page %d: %s\n", i, text)
        page.Close()
    }
}
```

## Testing

```bash
# Full tests (requires MuPDF libs for current platform)
go test -v -count=1 .

# Pure Go tests only (no CGO/libs required, works on any platform)
go test -v -count=1 -tags nomupdf .
```

## Core Types

| Go Type | PyMuPDF Equivalent | Description |
|---------|-------------------|-------------|
| `Document` | `fitz.Document` | Represents a document |
| `Page` | `fitz.Page` | Represents a page |
| `Pixmap` | `fitz.Pixmap` | Pixel map (raster image) |
| `Rect` | `fitz.Rect` | Rectangle |
| `Matrix` | `fitz.Matrix` | Transformation matrix |
| `Point` | `fitz.Point` | 2D point |
| `Quad` | `fitz.Quad` | Quadrilateral |
| `TextPage` | `fitz.TextPage` | Text extraction results |
| `Annot` | `fitz.Annot` | Annotation |
| `Widget` | `fitz.Widget` | Form field widget |

## Build Tags

- Default: CGO enabled, links against MuPDF static libraries
- `nomupdf`: Disables CGO, all MuPDF functions return `ErrInitFailed`. Useful for CI or environments without MuPDF.

## License

AGPL-3.0 (same as MuPDF/PyMuPDF)

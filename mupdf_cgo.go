//go:build cgo && !nomupdf

package gomupdf

/*
#cgo CFLAGS: -Iinclude
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/libs -lmupdf_windows_amd64 -lmupdfthird_windows_amd64 -lm -lgdi32 -lcomdlg32 -luser32 -ladvapi32 -lshell32
#cgo windows,arm64 LDFLAGS: -L${SRCDIR}/libs -lmupdf_windows_arm64 -lmupdfthird_windows_arm64 -lm -lgdi32 -lcomdlg32 -luser32 -ladvapi32 -lshell32
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/libs -lmupdf_linux_amd64 -lmupdfthird_linux_amd64 -lm -lstdc++ -lpthread -ldl
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/libs -lmupdf_linux_arm64 -lmupdfthird_linux_arm64 -lm -lstdc++ -lpthread -ldl
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/libs -lmupdf_darwin_amd64 -lmupdfthird_darwin_amd64 -lm -lc++ -framework CoreFoundation -framework Security
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/libs -lmupdf_darwin_arm64 -lmupdfthird_darwin_arm64 -lm -lc++ -framework CoreFoundation -framework Security

#include "gomupdf.h"
*/
import "C"
import (
	"unsafe"
)

// context wraps the MuPDF fz_context.
type context struct {
	ctx *C.fz_context
}

// newContext creates a new MuPDF context.
func newContext() (*context, error) {
	ctx := C.gomupdf_new_context()
	if ctx == nil {
		return nil, ErrInitFailed
	}
	C.fz_register_document_handlers(ctx)
	return &context{ctx: ctx}, nil
}

// close releases the context.
func (c *context) close() {
	if c.ctx != nil {
		C.gomupdf_drop_context(c.ctx)
		c.ctx = nil
	}
}

// clone creates a thread-safe clone of the context.
func (c *context) clone() *context {
	return &context{ctx: C.gomupdf_clone_context(c.ctx)}
}

// freeString frees a C string allocated by MuPDF.
func (c *context) freeString(s *C.char) {
	if s != nil {
		C.gomupdf_free(c.ctx, unsafe.Pointer(s))
	}
}

// freeBytes frees a byte buffer allocated by MuPDF.
func (c *context) freeBytes(p *C.uchar) {
	if p != nil {
		C.gomupdf_free(c.ctx, unsafe.Pointer(p))
	}
}

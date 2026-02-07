//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"

// Embedded files via PDF Names tree (no portfolio API in 1.24.9)
static int gomupdf_embfile_count(fz_context *ctx, pdf_document *pdf) {
    int count = 0;
    fz_try(ctx) {
        pdf_obj *root = pdf_dict_get(ctx, pdf_trailer(ctx, pdf), PDF_NAME(Root));
        pdf_obj *names = pdf_dict_get(ctx, root, PDF_NAME(Names));
        if (!names) return 0;
        pdf_obj *efs = pdf_dict_get(ctx, names, PDF_NAME(EmbeddedFiles));
        if (!efs) return 0;
        pdf_obj *namesArr = pdf_dict_get(ctx, efs, PDF_NAME(Names));
        if (!namesArr) return 0;
        count = pdf_array_len(ctx, namesArr) / 2;
    }
    fz_catch(ctx) { count = 0; }
    return count;
}

static const char* gomupdf_embfile_name(fz_context *ctx, pdf_document *pdf, int idx) {
    const char *name = NULL;
    fz_try(ctx) {
        pdf_obj *root = pdf_dict_get(ctx, pdf_trailer(ctx, pdf), PDF_NAME(Root));
        pdf_obj *names = pdf_dict_get(ctx, root, PDF_NAME(Names));
        pdf_obj *efs = pdf_dict_get(ctx, names, PDF_NAME(EmbeddedFiles));
        pdf_obj *namesArr = pdf_dict_get(ctx, efs, PDF_NAME(Names));
        name = pdf_to_text_string(ctx, pdf_array_get(ctx, namesArr, idx * 2));
    }
    fz_catch(ctx) { name = NULL; }
    return name;
}

static unsigned char* gomupdf_embfile_get(fz_context *ctx, pdf_document *pdf,
    int idx, int *outlen, int *errcode) {
    unsigned char *data = NULL;
    fz_try(ctx) {
        pdf_obj *root = pdf_dict_get(ctx, pdf_trailer(ctx, pdf), PDF_NAME(Root));
        pdf_obj *names = pdf_dict_get(ctx, root, PDF_NAME(Names));
        pdf_obj *efs = pdf_dict_get(ctx, names, PDF_NAME(EmbeddedFiles));
        pdf_obj *namesArr = pdf_dict_get(ctx, efs, PDF_NAME(Names));
        pdf_obj *filespec = pdf_array_get(ctx, namesArr, idx * 2 + 1);
        pdf_obj *ef = pdf_dict_get(ctx, filespec, PDF_NAME(EF));
        pdf_obj *stream = pdf_dict_get(ctx, ef, PDF_NAME(F));
        fz_buffer *buf = pdf_load_stream(ctx, stream);
        unsigned char *bufdata;
        size_t len = fz_buffer_storage(ctx, buf, &bufdata);
        data = (unsigned char*)fz_malloc(ctx, len);
        memcpy(data, bufdata, len);
        *outlen = (int)len;
        fz_drop_buffer(ctx, buf);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; data = NULL; *outlen = 0; }
    return data;
}
*/
import "C"
import "unsafe"

func (d *Document) EmbFileCount() int {
	if d.isClosed || !d.IsPDF() {
		return 0
	}
	return int(C.gomupdf_embfile_count(d.ctx.ctx, d.pdf))
}

func (d *Document) EmbFileNames() []string {
	if d.isClosed || !d.IsPDF() {
		return nil
	}
	count := d.EmbFileCount()
	names := make([]string, 0, count)
	for i := 0; i < count; i++ {
		name := C.gomupdf_embfile_name(d.ctx.ctx, d.pdf, C.int(i))
		if name != nil {
			names = append(names, C.GoString(name))
		}
	}
	return names
}

func (d *Document) EmbFileGet(index int) ([]byte, error) {
	if d.isClosed || !d.IsPDF() {
		return nil, ErrNotPDF
	}
	var outlen, errcode C.int
	data := C.gomupdf_embfile_get(d.ctx.ctx, d.pdf, C.int(index), &outlen, &errcode)
	if errcode != 0 || data == nil {
		return nil, ErrEmbeddedFile
	}
	defer d.ctx.freeBytes(data)
	return C.GoBytes(unsafe.Pointer(data), outlen), nil
}

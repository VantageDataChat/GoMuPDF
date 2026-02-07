// gomupdf.h - C wrapper functions for MuPDF 1.24.9
// This header is included by all CGO Go files via mupdf_cgo.go.

#ifndef GOMUPDF_H
#define GOMUPDF_H

#include <mupdf/fitz.h>
#include <mupdf/pdf.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

// ============================================================
// Context management
// ============================================================

static fz_context* gomupdf_new_context(void) {
    return fz_new_context(NULL, NULL, FZ_STORE_DEFAULT);
}

static void gomupdf_drop_context(fz_context *ctx) {
    fz_drop_context(ctx);
}

static fz_context* gomupdf_clone_context(fz_context *ctx) {
    return fz_clone_context(ctx);
}

// Free memory allocated by MuPDF
static void gomupdf_free(fz_context *ctx, void *ptr) {
    fz_free(ctx, ptr);
}

// ============================================================
// Document operations
// ============================================================

static fz_document* gomupdf_open_document(fz_context *ctx, const char *filename, int *errcode) {
    fz_document *doc = NULL;
    fz_try(ctx) { doc = fz_open_document(ctx, filename); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; doc = NULL; }
    return doc;
}

static fz_document* gomupdf_open_document_from_memory(fz_context *ctx,
    const char *magic, const unsigned char *data, int datalen, int *errcode) {
    fz_document *doc = NULL;
    fz_try(ctx) {
        fz_buffer *buf = fz_new_buffer_from_copied_data(ctx, data, datalen);
        fz_stream *stream = fz_open_buffer(ctx, buf);
        fz_drop_buffer(ctx, buf);
        doc = fz_open_document_with_stream(ctx, magic, stream);
        fz_drop_stream(ctx, stream);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; doc = NULL; }
    return doc;
}

static void gomupdf_drop_document(fz_context *ctx, fz_document *doc) {
    fz_drop_document(ctx, doc);
}

static int gomupdf_page_count(fz_context *ctx, fz_document *doc) {
    int count = 0;
    fz_try(ctx) { count = fz_count_pages(ctx, doc); }
    fz_catch(ctx) { count = 0; }
    return count;
}

static int gomupdf_needs_password(fz_context *ctx, fz_document *doc) {
    return fz_needs_password(ctx, doc);
}

static int gomupdf_authenticate_password(fz_context *ctx, fz_document *doc, const char *password) {
    return fz_authenticate_password(ctx, doc, password);
}

static pdf_document* gomupdf_pdf_document(fz_context *ctx, fz_document *doc) {
    return pdf_document_from_fz_document(ctx, doc);
}

static int gomupdf_is_document_reflowable(fz_context *ctx, fz_document *doc) {
    return fz_is_document_reflowable(ctx, doc);
}

static void gomupdf_layout_document(fz_context *ctx, fz_document *doc, float w, float h, float em) {
    fz_try(ctx) { fz_layout_document(ctx, doc, w, h, em); }
    fz_catch(ctx) { }
}

// ============================================================
// Metadata
// ============================================================

static char* gomupdf_lookup_metadata(fz_context *ctx, fz_document *doc, const char *key, int *errcode) {
    char buf[512];
    int n;
    fz_try(ctx) { n = fz_lookup_metadata(ctx, doc, key, buf, sizeof(buf)); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; return NULL; }
    if (n == -1) return NULL;
    return fz_strdup(ctx, buf);
}

static void gomupdf_set_metadata(fz_context *ctx, fz_document *doc, const char *key, const char *value) {
    fz_try(ctx) { fz_set_metadata(ctx, doc, key, value); }
    fz_catch(ctx) { }
}

// ============================================================
// Outline / TOC
// ============================================================

static fz_outline* gomupdf_load_outline(fz_context *ctx, fz_document *doc, int *errcode) {
    fz_outline *outline = NULL;
    fz_try(ctx) { outline = fz_load_outline(ctx, doc); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; outline = NULL; }
    return outline;
}

static void gomupdf_drop_outline(fz_context *ctx, fz_outline *outline) {
    fz_drop_outline(ctx, outline);
}

// ============================================================
// Page operations
// ============================================================

static fz_page* gomupdf_load_page(fz_context *ctx, fz_document *doc, int number, int *errcode) {
    fz_page *page = NULL;
    fz_try(ctx) { page = fz_load_page(ctx, doc, number); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; page = NULL; }
    return page;
}

static void gomupdf_drop_page(fz_context *ctx, fz_page *page) {
    fz_drop_page(ctx, page);
}

static fz_rect gomupdf_page_bound(fz_context *ctx, fz_page *page) {
    return fz_bound_page(ctx, page);
}

static int gomupdf_page_rotation(fz_context *ctx, pdf_page *page) {
    return pdf_to_int(ctx, pdf_dict_get_inheritable(ctx, page->obj, PDF_NAME(Rotate)));
}

static void gomupdf_set_page_rotation(fz_context *ctx, pdf_page *page, int rotation) {
    pdf_dict_put_int(ctx, page->obj, PDF_NAME(Rotate), rotation);
}

static char* gomupdf_page_label(fz_context *ctx, fz_page *page, int *errcode) {
    char buf[256];
    fz_try(ctx) { fz_page_label(ctx, page, buf, sizeof(buf)); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; return NULL; }
    return fz_strdup(ctx, buf);
}

// ============================================================
// Text extraction
// ============================================================

static fz_stext_page* gomupdf_new_stext_page(fz_context *ctx, fz_page *page, int flags, int *errcode) {
    fz_stext_page *tp = NULL;
    fz_stext_options opts;
    memset(&opts, 0, sizeof(opts));
    opts.flags = flags;
    fz_try(ctx) { tp = fz_new_stext_page_from_page(ctx, page, &opts); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; tp = NULL; }
    return tp;
}

static void gomupdf_drop_stext_page(fz_context *ctx, fz_stext_page *tp) {
    fz_drop_stext_page(ctx, tp);
}

static char* gomupdf_stext_page_as_text(fz_context *ctx, fz_stext_page *tp, int *errcode) {
    char *text = NULL;
    fz_try(ctx) {
        fz_buffer *buf = fz_new_buffer_from_stext_page(ctx, tp);
        text = fz_strdup(ctx, fz_string_from_buffer(ctx, buf));
        fz_drop_buffer(ctx, buf);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; text = NULL; }
    return text;
}

// ============================================================
// Search
// ============================================================

static int gomupdf_search_page(fz_context *ctx, fz_page *page, const char *needle,
    fz_quad *quads, int max_quads, int *errcode) {
    int count = 0;
    fz_try(ctx) {
        count = fz_search_page(ctx, page, needle, NULL, quads, max_quads);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; count = 0; }
    return count;
}

// ============================================================
// Links
// ============================================================

static fz_link* gomupdf_load_links(fz_context *ctx, fz_page *page, int *errcode) {
    fz_link *links = NULL;
    fz_try(ctx) { links = fz_load_links(ctx, page); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; links = NULL; }
    return links;
}

static void gomupdf_drop_link(fz_context *ctx, fz_link *link) {
    fz_drop_link(ctx, link);
}

// ============================================================
// Pixmap operations
// ============================================================

static fz_pixmap* gomupdf_page_to_pixmap(fz_context *ctx, fz_page *page,
    float a, float b, float c, float d, float e, float f,
    int colorspace, int alpha, int *errcode) {
    fz_pixmap *pix = NULL;
    fz_matrix ctm = {a, b, c, d, e, f};
    fz_colorspace *cs;
    switch(colorspace) {
        case 0: cs = fz_device_gray(ctx); break;
        case 2: cs = fz_device_cmyk(ctx); break;
        default: cs = fz_device_rgb(ctx); break;
    }
    fz_try(ctx) { pix = fz_new_pixmap_from_page(ctx, page, ctm, cs, alpha); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; pix = NULL; }
    return pix;
}

static fz_pixmap* gomupdf_page_to_pixmap_clipped(fz_context *ctx, fz_page *page,
    float a, float b, float c, float d, float e, float f,
    int colorspace, int alpha,
    float cx0, float cy0, float cx1, float cy1, int *errcode) {
    fz_pixmap *pix = NULL;
    fz_matrix ctm = {a, b, c, d, e, f};
    fz_rect clip = {cx0, cy0, cx1, cy1};
    fz_irect iclip = fz_round_rect(clip);
    fz_colorspace *cs;
    switch(colorspace) {
        case 0: cs = fz_device_gray(ctx); break;
        case 2: cs = fz_device_cmyk(ctx); break;
        default: cs = fz_device_rgb(ctx); break;
    }
    fz_try(ctx) {
        pix = fz_new_pixmap_from_page_contents(ctx, page, ctm, cs, alpha);
        // TODO: clip the pixmap to iclip if needed
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; pix = NULL; }
    (void)iclip;
    return pix;
}

static void gomupdf_drop_pixmap(fz_context *ctx, fz_pixmap *pix) {
    fz_drop_pixmap(ctx, pix);
}

static int gomupdf_pixmap_width(fz_pixmap *pix) { return pix->w; }
static int gomupdf_pixmap_height(fz_pixmap *pix) { return pix->h; }
static int gomupdf_pixmap_n(fz_pixmap *pix) { return pix->n; }
static int gomupdf_pixmap_alpha(fz_pixmap *pix) { return pix->alpha; }
static int gomupdf_pixmap_stride(fz_pixmap *pix) { return pix->stride; }
static int gomupdf_pixmap_x(fz_pixmap *pix) { return pix->x; }
static int gomupdf_pixmap_y(fz_pixmap *pix) { return pix->y; }
static unsigned char* gomupdf_pixmap_samples(fz_pixmap *pix) { return pix->samples; }
static int gomupdf_pixmap_samples_len(fz_pixmap *pix) { return pix->h * pix->stride; }

static unsigned char* gomupdf_pixmap_to_png(fz_context *ctx, fz_pixmap *pix, int *outlen, int *errcode) {
    unsigned char *data = NULL;
    fz_try(ctx) {
        fz_buffer *buf = fz_new_buffer_from_pixmap_as_png(ctx, pix, fz_default_color_params);
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

static fz_pixmap* gomupdf_new_pixmap(fz_context *ctx, int colorspace, int w, int h, int alpha) {
    fz_colorspace *cs;
    switch(colorspace) {
        case 0: cs = fz_device_gray(ctx); break;
        case 2: cs = fz_device_cmyk(ctx); break;
        default: cs = fz_device_rgb(ctx); break;
    }
    fz_pixmap *pix = fz_new_pixmap(ctx, cs, w, h, NULL, alpha);
    fz_clear_pixmap(ctx, pix);
    return pix;
}

static fz_pixmap* gomupdf_pixmap_from_image(fz_context *ctx, fz_document *doc,
    int xref, int *errcode) {
    fz_pixmap *pix = NULL;
    fz_try(ctx) {
        pdf_document *pdf = pdf_document_from_fz_document(ctx, doc);
        pdf_obj *ref = pdf_new_indirect(ctx, pdf, xref, 0);
        fz_image *img = pdf_load_image(ctx, pdf, ref);
        pix = fz_get_pixmap_from_image(ctx, img, NULL, NULL, NULL, NULL);
        fz_drop_image(ctx, img);
        pdf_drop_obj(ctx, ref);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; pix = NULL; }
    return pix;
}

static void gomupdf_pixmap_set_pixel(fz_pixmap *pix, int x, int y, unsigned char *color, int n) {
    unsigned char *s = pix->samples + y * pix->stride + x * pix->n;
    int i;
    for (i = 0; i < n && i < pix->n; i++) s[i] = color[i];
}

static void gomupdf_pixmap_get_pixel(fz_pixmap *pix, int x, int y, unsigned char *color, int n) {
    unsigned char *s = pix->samples + y * pix->stride + x * pix->n;
    int i;
    for (i = 0; i < n && i < pix->n; i++) color[i] = s[i];
}

static void gomupdf_pixmap_clear(fz_context *ctx, fz_pixmap *pix, int value) {
    if (value < 0) fz_clear_pixmap(ctx, pix);
    else fz_clear_pixmap_with_value(ctx, pix, value);
}

static void gomupdf_pixmap_invert(fz_context *ctx, fz_pixmap *pix) {
    fz_invert_pixmap(ctx, pix);
}

static void gomupdf_pixmap_gamma(fz_context *ctx, fz_pixmap *pix, float gamma) {
    fz_gamma_pixmap(ctx, pix, gamma);
}

static void gomupdf_pixmap_tint(fz_context *ctx, fz_pixmap *pix, int black, int white) {
    fz_tint_pixmap(ctx, pix, black, white);
}

static int gomupdf_pixmap_save_png(fz_context *ctx, fz_pixmap *pix, const char *filename) {
    int errcode = 0;
    fz_try(ctx) { fz_save_pixmap_as_png(ctx, pix, filename); }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

static int gomupdf_pixmap_save_pnm(fz_context *ctx, fz_pixmap *pix, const char *filename) {
    int errcode = 0;
    fz_try(ctx) { fz_save_pixmap_as_pnm(ctx, pix, filename); }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

static fz_pixmap* gomupdf_pixmap_convert(fz_context *ctx, fz_pixmap *pix, int colorspace, int *errcode) {
    fz_colorspace *cs;
    switch(colorspace) {
        case 0: cs = fz_device_gray(ctx); break;
        case 2: cs = fz_device_cmyk(ctx); break;
        default: cs = fz_device_rgb(ctx); break;
    }
    fz_pixmap *result = NULL;
    fz_try(ctx) { result = fz_convert_pixmap(ctx, pix, cs, NULL, NULL, fz_default_color_params, 1); *errcode = 0; }
    fz_catch(ctx) { *errcode = 1; result = NULL; }
    return result;
}

// ============================================================
// PDF save / write
// ============================================================

static int gomupdf_pdf_save(fz_context *ctx, pdf_document *pdf, const char *filename,
    int garbage, int deflate, int linear, int clean, int ascii,
    int incremental, int pretty, int encryption,
    const char *owner_pw, const char *user_pw, int permissions) {
    int errcode = 0;
    pdf_write_options opts;
    memset(&opts, 0, sizeof(opts));
    opts.do_garbage = garbage;
    opts.do_compress = deflate;
    opts.do_linear = linear;
    opts.do_clean = clean;
    opts.do_ascii = ascii;
    opts.do_incremental = incremental;
    opts.do_pretty = pretty;
    opts.do_encrypt = encryption;
    opts.permissions = permissions;
    if (owner_pw) strncpy(opts.opwd_utf8, owner_pw, sizeof(opts.opwd_utf8)-1);
    if (user_pw) strncpy(opts.upwd_utf8, user_pw, sizeof(opts.upwd_utf8)-1);
    fz_try(ctx) { pdf_save_document(ctx, pdf, filename, &opts); }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

static unsigned char* gomupdf_pdf_tobytes(fz_context *ctx, pdf_document *pdf,
    int garbage, int deflate, int clean, int ascii, int pretty,
    int *outlen, int *errcode) {
    unsigned char *data = NULL;
    fz_try(ctx) {
        fz_buffer *buf = fz_new_buffer(ctx, 8192);
        fz_output *out = fz_new_output_with_buffer(ctx, buf);
        pdf_write_options opts;
        memset(&opts, 0, sizeof(opts));
        opts.do_garbage = garbage;
        opts.do_compress = deflate;
        opts.do_clean = clean;
        opts.do_ascii = ascii;
        opts.do_pretty = pretty;
        pdf_write_document(ctx, pdf, out, &opts);
        fz_close_output(ctx, out);
        fz_drop_output(ctx, out);
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

// ============================================================
// PDF page manipulation
// ============================================================

static int gomupdf_insert_page(fz_context *ctx, pdf_document *pdf, int pno,
    float width, float height) {
    int errcode = 0;
    fz_try(ctx) {
        fz_rect mediabox = {0, 0, width, height};
        pdf_obj *resources = NULL;
        fz_buffer *contents = NULL;
        fz_device *dev = pdf_page_write(ctx, pdf, mediabox, &resources, &contents);
        fz_close_device(ctx, dev);
        fz_drop_device(ctx, dev);
        pdf_obj *page_obj = pdf_add_page(ctx, pdf, mediabox, 0, resources, contents);
        pdf_insert_page(ctx, pdf, pno, page_obj);
        pdf_drop_obj(ctx, page_obj);
        fz_drop_buffer(ctx, contents);
        pdf_drop_obj(ctx, resources);
    }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

static int gomupdf_delete_page(fz_context *ctx, pdf_document *pdf, int pno) {
    int errcode = 0;
    fz_try(ctx) { pdf_delete_page(ctx, pdf, pno); }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

static void gomupdf_rearrange_pages(fz_context *ctx, pdf_document *pdf, int count, const int *pages) {
    fz_try(ctx) { pdf_rearrange_pages(ctx, pdf, count, pages); }
    fz_catch(ctx) { }
}

// ============================================================
// PDF xref operations
// ============================================================

static int gomupdf_xref_len(fz_context *ctx, pdf_document *pdf) {
    return pdf_xref_len(ctx, pdf);
}

static char* gomupdf_xref_object_str(fz_context *ctx, pdf_document *pdf, int xref,
    int compressed, int *errcode) {
    char *result = NULL;
    fz_try(ctx) {
        pdf_obj *obj = pdf_load_object(ctx, pdf, xref);
        fz_buffer *buf = fz_new_buffer(ctx, 512);
        fz_output *out = fz_new_output_with_buffer(ctx, buf);
        pdf_print_obj(ctx, out, obj, compressed ? 1 : 0, 0);
        fz_close_output(ctx, out);
        fz_drop_output(ctx, out);
        result = fz_strdup(ctx, fz_string_from_buffer(ctx, buf));
        fz_drop_buffer(ctx, buf);
        pdf_drop_obj(ctx, obj);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; result = NULL; }
    return result;
}

// ============================================================
// PDF InsertPDF (graft pages)
// ============================================================

static int gomupdf_graft_page(fz_context *ctx, pdf_document *dst, pdf_document *src,
    pdf_graft_map *map, int page_to, int page_from) {
    int errcode = 0;
    fz_try(ctx) {
        pdf_graft_mapped_page(ctx, map, page_to, src, page_from);
    }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

// ============================================================
// Convert document to PDF
// ============================================================

static unsigned char* gomupdf_convert_to_pdf(fz_context *ctx, fz_document *doc,
    int from_page, int to_page, int rotate, int *outlen, int *errcode) {
    unsigned char *data = NULL;
    fz_try(ctx) {
        pdf_document *pdfout = pdf_create_document(ctx);
        int page_count = fz_count_pages(ctx, doc);
        if (from_page < 0) from_page = 0;
        if (to_page < 0 || to_page >= page_count) to_page = page_count - 1;

        for (int i = from_page; i <= to_page; i++) {
            fz_page *page = fz_load_page(ctx, doc, i);
            fz_rect mediabox = fz_bound_page(ctx, page);
            pdf_obj *resources = NULL;
            fz_buffer *contents = NULL;
            fz_device *dev = pdf_page_write(ctx, pdfout, mediabox, &resources, &contents);
            fz_run_page(ctx, page, dev, fz_identity, NULL);
            fz_close_device(ctx, dev);
            fz_drop_device(ctx, dev);
            pdf_obj *page_obj = pdf_add_page(ctx, pdfout, mediabox, rotate, resources, contents);
            pdf_insert_page(ctx, pdfout, -1, page_obj);
            pdf_drop_obj(ctx, page_obj);
            fz_drop_buffer(ctx, contents);
            pdf_drop_obj(ctx, resources);
            fz_drop_page(ctx, page);
        }

        fz_buffer *buf = fz_new_buffer(ctx, 8192);
        fz_output *out = fz_new_output_with_buffer(ctx, buf);
        pdf_write_options opts;
        memset(&opts, 0, sizeof(opts));
        opts.do_garbage = 4;
        opts.do_compress = 1;
        pdf_write_document(ctx, pdfout, out, &opts);
        fz_close_output(ctx, out);
        fz_drop_output(ctx, out);
        unsigned char *bufdata;
        size_t len = fz_buffer_storage(ctx, buf, &bufdata);
        data = (unsigned char*)fz_malloc(ctx, len);
        memcpy(data, bufdata, len);
        *outlen = (int)len;
        fz_drop_buffer(ctx, buf);
        pdf_drop_document(ctx, pdfout);
        *errcode = 0;
    }
    fz_catch(ctx) { *errcode = 1; data = NULL; *outlen = 0; }
    return data;
}

// ============================================================
// Annotations
// ============================================================

static pdf_annot* gomupdf_first_annot(fz_context *ctx, pdf_page *page) {
    return pdf_first_annot(ctx, page);
}

static pdf_annot* gomupdf_next_annot(fz_context *ctx, pdf_annot *annot) {
    return pdf_next_annot(ctx, annot);
}

static int gomupdf_annot_type(fz_context *ctx, pdf_annot *annot) {
    return (int)pdf_annot_type(ctx, annot);
}

static fz_rect gomupdf_annot_rect(fz_context *ctx, pdf_annot *annot) {
    return pdf_bound_annot(ctx, annot);
}

static const char* gomupdf_annot_contents(fz_context *ctx, pdf_annot *annot) {
    return pdf_annot_contents(ctx, annot);
}

static void gomupdf_set_annot_contents(fz_context *ctx, pdf_annot *annot, const char *text) {
    pdf_set_annot_contents(ctx, annot, text);
}

static void gomupdf_delete_annot(fz_context *ctx, pdf_page *page, pdf_annot *annot) {
    pdf_delete_annot(ctx, page, annot);
}

static int gomupdf_annot_xref(fz_context *ctx, pdf_annot *annot) {
    return pdf_to_num(ctx, pdf_annot_obj(ctx, annot));
}

static pdf_annot* gomupdf_add_text_annot(fz_context *ctx, pdf_page *page,
    float x, float y, const char *text) {
    pdf_annot *annot = NULL;
    fz_try(ctx) {
        annot = pdf_create_annot(ctx, page, PDF_ANNOT_TEXT);
        pdf_set_annot_contents(ctx, annot, text);
        fz_rect r = {x, y, x + 20, y + 20};
        pdf_set_annot_rect(ctx, annot, r);
        pdf_update_annot(ctx, annot);
    }
    fz_catch(ctx) { annot = NULL; }
    return annot;
}

static pdf_annot* gomupdf_add_highlight_annot(fz_context *ctx, pdf_page *page,
    fz_quad *quads, int nquads) {
    pdf_annot *annot = NULL;
    fz_try(ctx) {
        annot = pdf_create_annot(ctx, page, PDF_ANNOT_HIGHLIGHT);
        pdf_set_annot_quad_points(ctx, annot, nquads, quads);
        pdf_update_annot(ctx, annot);
    }
    fz_catch(ctx) { annot = NULL; }
    return annot;
}

static pdf_annot* gomupdf_add_freetext_annot(fz_context *ctx, pdf_page *page,
    float x0, float y0, float x1, float y1, const char *text, float fontsize) {
    pdf_annot *annot = NULL;
    fz_try(ctx) {
        annot = pdf_create_annot(ctx, page, PDF_ANNOT_FREE_TEXT);
        fz_rect r = {x0, y0, x1, y1};
        pdf_set_annot_rect(ctx, annot, r);
        pdf_set_annot_contents(ctx, annot, text);
        pdf_update_annot(ctx, annot);
    }
    fz_catch(ctx) { annot = NULL; }
    return annot;
}

// ============================================================
// Widgets (form fields)
// ============================================================

static pdf_annot* gomupdf_first_widget(fz_context *ctx, pdf_page *page) {
    return pdf_first_widget(ctx, page);
}

static pdf_annot* gomupdf_next_widget(fz_context *ctx, pdf_annot *widget) {
    return pdf_next_widget(ctx, widget);
}

static int gomupdf_widget_type(fz_context *ctx, pdf_annot *widget) {
    return (int)pdf_widget_type(ctx, widget);
}

static const char* gomupdf_widget_name(fz_context *ctx, pdf_annot *widget) {
    return pdf_annot_field_label(ctx, widget);
}

static const char* gomupdf_widget_value(fz_context *ctx, pdf_annot *widget) {
    return pdf_annot_field_value(ctx, widget);
}

static int gomupdf_set_widget_value(fz_context *ctx, pdf_document *doc, pdf_annot *widget, const char *value) {
    int errcode = 0;
    fz_try(ctx) {
        pdf_set_annot_field_value(ctx, doc, widget, value, 0);
        pdf_update_annot(ctx, widget);
    }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

// ============================================================
// Text insertion (Shape)
// ============================================================

/* Detect whether a UTF-8 string contains any non-ASCII characters (CJK, etc.).
   Returns 1 if CJK/non-Latin content is found, 0 if pure ASCII/Latin. */
static int gomupdf_text_needs_cjk(const char *text) {
    for (const unsigned char *p = (const unsigned char *)text; *p; p++) {
        if (*p > 0x7F) return 1;
    }
    return 0;
}

/* Decode one UTF-8 codepoint from *src, advance *src, return the codepoint.
   Returns 0xFFFD on invalid sequences. */
static unsigned int gomupdf_utf8_decode(const unsigned char **src) {
    const unsigned char *s = *src;
    unsigned int cp;
    int extra;
    if (s[0] < 0x80)       { cp = s[0]; extra = 0; }
    else if (s[0] < 0xC0)  { cp = 0xFFFD; extra = 0; }
    else if (s[0] < 0xE0)  { cp = s[0] & 0x1F; extra = 1; }
    else if (s[0] < 0xF0)  { cp = s[0] & 0x0F; extra = 2; }
    else if (s[0] < 0xF8)  { cp = s[0] & 0x07; extra = 3; }
    else                    { cp = 0xFFFD; extra = 0; }
    s++;
    for (int i = 0; i < extra; i++) {
        if ((*s & 0xC0) != 0x80) { cp = 0xFFFD; break; }
        cp = (cp << 6) | (*s & 0x3F);
        s++;
    }
    *src = s;
    return cp;
}

/* Append a UTF-8 string as a PDF hex string (<FEFF...>) encoded in UTF-16BE.
   This is the standard encoding for CID fonts with Identity-H CMap. */
static void gomupdf_append_utf16be_hex(fz_context *ctx, fz_buffer *buf, const char *text) {
    static const char hex[] = "0123456789ABCDEF";
    fz_append_byte(ctx, buf, '<');
    /* BOM: FEFF */
    fz_append_string(ctx, buf, "FEFF");
    const unsigned char *s = (const unsigned char *)text;
    while (*s) {
        unsigned int cp = gomupdf_utf8_decode(&s);
        if (cp <= 0xFFFF) {
            /* BMP character: single UTF-16 code unit */
            fz_append_byte(ctx, buf, hex[(cp >> 12) & 0xF]);
            fz_append_byte(ctx, buf, hex[(cp >>  8) & 0xF]);
            fz_append_byte(ctx, buf, hex[(cp >>  4) & 0xF]);
            fz_append_byte(ctx, buf, hex[ cp        & 0xF]);
        } else if (cp <= 0x10FFFF) {
            /* Supplementary character: surrogate pair */
            cp -= 0x10000;
            unsigned int hi = 0xD800 | (cp >> 10);
            unsigned int lo = 0xDC00 | (cp & 0x3FF);
            fz_append_byte(ctx, buf, hex[(hi >> 12) & 0xF]);
            fz_append_byte(ctx, buf, hex[(hi >>  8) & 0xF]);
            fz_append_byte(ctx, buf, hex[(hi >>  4) & 0xF]);
            fz_append_byte(ctx, buf, hex[ hi        & 0xF]);
            fz_append_byte(ctx, buf, hex[(lo >> 12) & 0xF]);
            fz_append_byte(ctx, buf, hex[(lo >>  8) & 0xF]);
            fz_append_byte(ctx, buf, hex[(lo >>  4) & 0xF]);
            fz_append_byte(ctx, buf, hex[ lo        & 0xF]);
        }
    }
    fz_append_byte(ctx, buf, '>');
}

/* Detect the best CJK ordering for a UTF-8 string by scanning codepoint ranges.
   Returns FZ_ADOBE_GB (Simplified Chinese) as default for any CJK text.
   ordering: 0=CNS (Traditional Chinese), 1=GB (Simplified Chinese),
             2=Japan, 3=Korea */
static int gomupdf_detect_cjk_ordering(const char *text) {
    int has_jp = 0, has_kr = 0, has_tc = 0;
    const unsigned char *s = (const unsigned char *)text;
    while (*s) {
        unsigned int cp = gomupdf_utf8_decode(&s);
        /* Hiragana / Katakana => Japanese */
        if ((cp >= 0x3040 && cp <= 0x309F) || (cp >= 0x30A0 && cp <= 0x30FF))
            has_jp = 1;
        /* Hangul => Korean */
        else if ((cp >= 0xAC00 && cp <= 0xD7AF) || (cp >= 0x1100 && cp <= 0x11FF))
            has_kr = 1;
        /* Bopomofo => Traditional Chinese */
        else if (cp >= 0x3100 && cp <= 0x312F)
            has_tc = 1;
    }
    if (has_jp) return FZ_ADOBE_JAPAN;
    if (has_kr) return FZ_ADOBE_KOREA;
    if (has_tc) return FZ_ADOBE_CNS;
    return FZ_ADOBE_GB; /* default: Simplified Chinese / generic CJK */
}

/* Return the CMap name for a given CJK ordering (UTF-16 horizontal). */
static const char* gomupdf_cjk_cmap_name(int ordering) {
    switch (ordering) {
        case FZ_ADOBE_CNS:   return "UniCNS-UTF16-H";
        case FZ_ADOBE_GB:    return "UniGB-UTF16-H";
        case FZ_ADOBE_JAPAN: return "UniJIS-UTF16-H";
        case FZ_ADOBE_KOREA: return "UniKS-UTF16-H";
        default:             return "UniGB-UTF16-H";
    }
}

/* Return the Adobe ordering string for a given CJK ordering. */
static const char* gomupdf_cjk_ordering_name(int ordering) {
    switch (ordering) {
        case FZ_ADOBE_CNS:   return "CNS1";
        case FZ_ADOBE_GB:    return "GB1";
        case FZ_ADOBE_JAPAN: return "Japan1";
        case FZ_ADOBE_KOREA: return "Korea1";
        default:             return "GB1";
    }
}

/* Return the supplement number for a given CJK ordering. */
static int gomupdf_cjk_supplement(int ordering) {
    switch (ordering) {
        case FZ_ADOBE_CNS:   return 7;
        case FZ_ADOBE_GB:    return 5;
        case FZ_ADOBE_JAPAN: return 7;
        case FZ_ADOBE_KOREA: return 2;
        default:             return 5;
    }
}

/* Create a non-embedded CID font (Type0) for CJK text.
   This builds the PDF font dictionary manually so it works even when
   MuPDF is compiled without built-in CJK font data.
   The PDF viewer will substitute an appropriate system CJK font. */
static pdf_obj* gomupdf_create_cjk_font(fz_context *ctx, pdf_document *doc, int ordering) {
    const char *cmap = gomupdf_cjk_cmap_name(ordering);
    const char *ord_name = gomupdf_cjk_ordering_name(ordering);
    int supplement = gomupdf_cjk_supplement(ordering);

    /* CIDSystemInfo dictionary */
    pdf_obj *sysinfo = pdf_new_dict(ctx, doc, 3);
    pdf_dict_put_text_string(ctx, sysinfo, PDF_NAME(Registry), "Adobe");
    pdf_dict_put_text_string(ctx, sysinfo, PDF_NAME(Ordering), ord_name);
    pdf_dict_put_int(ctx, sysinfo, PDF_NAME(Supplement), supplement);

    /* CIDFont (Type 0 descendant) dictionary */
    pdf_obj *cidfont = pdf_new_dict(ctx, doc, 5);
    pdf_dict_put(ctx, cidfont, PDF_NAME(Type), PDF_NAME(Font));
    pdf_dict_put_name(ctx, cidfont, PDF_NAME(Subtype), "CIDFontType0");
    pdf_dict_put_text_string(ctx, cidfont, PDF_NAME(BaseFont), "Adobe-Identity");
    pdf_dict_put(ctx, cidfont, PDF_NAME(CIDSystemInfo), sysinfo);
    /* DW (default width) = 1000 (standard for CJK fonts) */
    pdf_dict_put_int(ctx, cidfont, PDF_NAME(DW), 1000);

    pdf_obj *cidfont_ref = pdf_add_object(ctx, doc, cidfont);
    pdf_drop_obj(ctx, cidfont);
    pdf_drop_obj(ctx, sysinfo);

    /* Descendants array */
    pdf_obj *descendants = pdf_new_array(ctx, doc, 1);
    pdf_array_push(ctx, descendants, cidfont_ref);
    pdf_drop_obj(ctx, cidfont_ref);

    /* Type0 (composite) font dictionary */
    pdf_obj *fontdict = pdf_new_dict(ctx, doc, 5);
    pdf_dict_put(ctx, fontdict, PDF_NAME(Type), PDF_NAME(Font));
    pdf_dict_put_name(ctx, fontdict, PDF_NAME(Subtype), "Type0");
    pdf_dict_put_text_string(ctx, fontdict, PDF_NAME(BaseFont), "Adobe-Identity");
    pdf_dict_put_name(ctx, fontdict, PDF_NAME(Encoding), cmap);
    pdf_dict_put(ctx, fontdict, PDF_NAME(DescendantFonts), descendants);
    pdf_drop_obj(ctx, descendants);

    pdf_obj *fontdict_ref = pdf_add_object(ctx, doc, fontdict);
    pdf_drop_obj(ctx, fontdict);

    return fontdict_ref;
}

static int gomupdf_insert_text(fz_context *ctx, pdf_document *doc, int pno,
    float x, float y, const char *text, const char *fontname, float fontsize,
    float r, float g, float b) {
    int errcode = 0;
    fz_try(ctx) {
        int use_cjk = gomupdf_text_needs_cjk(text);

        /* Look up the existing page object */
        pdf_obj *page_obj = pdf_lookup_page_obj(ctx, doc, pno);

        /* Get or create the Resources dict and its Font sub-dict */
        pdf_obj *resources = pdf_dict_get(ctx, page_obj, PDF_NAME(Resources));
        if (!resources)
            resources = pdf_dict_put_dict(ctx, page_obj, PDF_NAME(Resources), 2);
        pdf_obj *fonts = pdf_dict_get(ctx, resources, PDF_NAME(Font));
        if (!fonts)
            fonts = pdf_dict_put_dict(ctx, resources, PDF_NAME(Font), 4);

        /* Create a unique font resource name */
        char fname[32];
        snprintf(fname, sizeof(fname), "F%d", pdf_create_object(ctx, doc));

        if (use_cjk) {
            /* CJK path: create non-embedded CID font with standard CMap */
            int ordering = gomupdf_detect_cjk_ordering(text);
            pdf_obj *font_obj = gomupdf_create_cjk_font(ctx, doc, ordering);
            pdf_dict_puts(ctx, fonts, fname, font_obj);
            pdf_drop_obj(ctx, font_obj);
        } else {
            /* Latin path: use Base14 simple font */
            fz_font *font = fz_new_base14_font(ctx, fontname);
            pdf_obj *font_obj = pdf_add_simple_font(ctx, doc, font, PDF_SIMPLE_ENCODING_LATIN);
            pdf_dict_puts(ctx, fonts, fname, font_obj);
            pdf_drop_obj(ctx, font_obj);
            fz_drop_font(ctx, font);
        }

        /* Build the content stream */
        fz_buffer *content = fz_new_buffer(ctx, 256);
        fz_append_printf(ctx, content, "q BT\n");
        fz_append_printf(ctx, content, "%g %g %g rg\n", r, g, b);
        fz_append_printf(ctx, content, "/%s %g Tf\n", fname, fontsize);
        fz_append_printf(ctx, content, "%g %g Td\n", x, y);

        if (use_cjk) {
            /* CJK: emit UTF-16BE hex string for CID font */
            gomupdf_append_utf16be_hex(ctx, content, text);
            fz_append_string(ctx, content, " Tj\n");
        } else {
            /* Latin: emit escaped PDF literal string */
            fz_append_byte(ctx, content, '(');
            for (const char *p = text; *p; p++) {
                if (*p == '(' || *p == ')' || *p == '\\')
                    fz_append_byte(ctx, content, '\\');
                fz_append_byte(ctx, content, (unsigned char)*p);
            }
            fz_append_string(ctx, content, ") Tj\n");
        }
        fz_append_string(ctx, content, "ET Q\n");

        /* Append the new content stream to the page's Contents array */
        pdf_obj *existing = pdf_dict_get(ctx, page_obj, PDF_NAME(Contents));
        if (pdf_is_array(ctx, existing)) {
            pdf_obj *newstream = pdf_add_stream(ctx, doc, content, NULL, 0);
            pdf_array_push(ctx, existing, newstream);
            pdf_drop_obj(ctx, newstream);
        } else {
            pdf_obj *arr = pdf_new_array(ctx, doc, 2);
            if (existing) pdf_array_push(ctx, arr, existing);
            pdf_obj *newstream = pdf_add_stream(ctx, doc, content, NULL, 0);
            pdf_array_push(ctx, arr, newstream);
            pdf_drop_obj(ctx, newstream);
            pdf_dict_put(ctx, page_obj, PDF_NAME(Contents), arr);
            pdf_drop_obj(ctx, arr);
        }

        fz_drop_buffer(ctx, content);
    }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

// ============================================================
// Image insertion
// ============================================================

static int gomupdf_insert_image(fz_context *ctx, pdf_document *doc, int pno,
    float x0, float y0, float x1, float y1,
    const unsigned char *imgdata, int imglen,
    int keep_proportion, int overlay) {
    int errcode = 0;
    fz_try(ctx) {
        fz_buffer *buf = fz_new_buffer_from_copied_data(ctx, imgdata, imglen);
        fz_image *img = fz_new_image_from_buffer(ctx, buf);
        fz_drop_buffer(ctx, buf);

        fz_rect rect = {x0, y0, x1, y1};
        if (keep_proportion) {
            float img_w = (float)img->w;
            float img_h = (float)img->h;
            float rect_w = x1 - x0;
            float rect_h = y1 - y0;
            float scale_w = rect_w / img_w;
            float scale_h = rect_h / img_h;
            float scale = scale_w < scale_h ? scale_w : scale_h;
            float new_w = img_w * scale;
            float new_h = img_h * scale;
            rect.x0 = x0 + (rect_w - new_w) / 2;
            rect.y0 = y0 + (rect_h - new_h) / 2;
            rect.x1 = rect.x0 + new_w;
            rect.y1 = rect.y0 + new_h;
        }

        pdf_obj *page_obj = pdf_lookup_page_obj(ctx, doc, pno);
        pdf_obj *resources = pdf_dict_get(ctx, page_obj, PDF_NAME(Resources));
        if (!resources)
            resources = pdf_dict_put_dict(ctx, page_obj, PDF_NAME(Resources), 2);

        pdf_obj *xobjects = pdf_dict_get(ctx, resources, PDF_NAME(XObject));
        if (!xobjects)
            xobjects = pdf_dict_put_dict(ctx, resources, PDF_NAME(XObject), 4);

        char name[32];
        snprintf(name, sizeof(name), "Img%d", pdf_create_object(ctx, doc));
        pdf_obj *imgref = pdf_add_image(ctx, doc, img);
        pdf_dict_puts(ctx, xobjects, name, imgref);

        fz_buffer *content = fz_new_buffer(ctx, 256);
        fz_append_printf(ctx, content, "q %g 0 0 %g %g %g cm /%s Do Q\n",
            rect.x1 - rect.x0, rect.y1 - rect.y0, rect.x0, rect.y0, name);

        pdf_obj *existing = pdf_dict_get(ctx, page_obj, PDF_NAME(Contents));
        if (pdf_is_array(ctx, existing)) {
            pdf_obj *newstream = pdf_add_stream(ctx, doc, content, NULL, 0);
            pdf_array_push(ctx, existing, newstream);
            pdf_drop_obj(ctx, newstream);
        } else {
            pdf_obj *arr = pdf_new_array(ctx, doc, 2);
            if (existing) pdf_array_push(ctx, arr, existing);
            pdf_obj *newstream = pdf_add_stream(ctx, doc, content, NULL, 0);
            pdf_array_push(ctx, arr, newstream);
            pdf_drop_obj(ctx, newstream);
            pdf_dict_put(ctx, page_obj, PDF_NAME(Contents), arr);
            pdf_drop_obj(ctx, arr);
        }

        fz_drop_buffer(ctx, content);
        fz_drop_image(ctx, img);
        pdf_drop_obj(ctx, imgref);
    }
    fz_catch(ctx) { errcode = 1; }
    return errcode;
}

// ============================================================
// PDF creation
// ============================================================

static fz_document* gomupdf_pdf_to_fz_document(pdf_document *pdf) {
    return &pdf->super;
}

// ============================================================
// Structured text block helpers (union access)
// ============================================================

static fz_stext_line* gomupdf_stext_block_first_line(fz_stext_block *block) {
    if (block->type != FZ_STEXT_BLOCK_TEXT) return NULL;
    return block->u.t.first_line;
}

// ============================================================
// PDF catalog helper
// ============================================================

static int gomupdf_pdf_catalog_xref(fz_context *ctx, pdf_document *pdf) {
    pdf_obj *trailer = pdf_trailer(ctx, pdf);
    pdf_obj *root = pdf_dict_get(ctx, trailer, PDF_NAME(Root));
    return pdf_to_num(ctx, root);
}

// ============================================================
// PDF incremental save check
// ============================================================

static int gomupdf_can_save_incrementally(fz_context *ctx, pdf_document *pdf) {
    return pdf_can_be_saved_incrementally(ctx, pdf);
}

#endif /* GOMUPDF_H */

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
// System font loading (for Story/HTML CJK support)
// ============================================================

/* CJK 字体加载回调 — 从系统字体目录加载 */
static fz_font* gomupdf_load_system_cjk_font(fz_context *ctx,
    const char *name, int ordering, int serif) {
    (void)name; (void)serif;
    const char *fontpath = NULL;
#ifdef _WIN32
    switch (ordering) {
        case FZ_ADOBE_GB:    fontpath = "C:\\Windows\\Fonts\\simsun.ttc"; break;
        case FZ_ADOBE_CNS:   fontpath = "C:\\Windows\\Fonts\\simsun.ttc"; break;
        case FZ_ADOBE_JAPAN: fontpath = "C:\\Windows\\Fonts\\msgothic.ttc"; break;
        case FZ_ADOBE_KOREA: fontpath = "C:\\Windows\\Fonts\\malgun.ttf"; break;
        default:             fontpath = "C:\\Windows\\Fonts\\simsun.ttc"; break;
    }
#elif defined(__APPLE__)
    (void)ordering;
    fontpath = "/System/Library/Fonts/PingFang.ttc";
#else
    (void)ordering;
    fontpath = "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc";
#endif
    if (!fontpath) return NULL;
    fz_font *font = NULL;
    fz_try(ctx) {
        font = fz_new_font_from_file(ctx, NULL, fontpath, 0, 0);
    }
    fz_catch(ctx) { font = NULL; }
    return font;
}

/* 普通字体加载回调 — 返回 NULL 让 MuPDF 用内置 Base14 */
static fz_font* gomupdf_load_system_font(fz_context *ctx,
    const char *name, int bold, int italic, int needs_exact_metrics) {
    (void)name; (void)bold; (void)italic; (void)needs_exact_metrics;
    return NULL;
}

/* 回退字体加载回调 — 用于找不到主字体时的 fallback */
static fz_font* gomupdf_load_system_fallback_font(fz_context *ctx,
    int script, int language, int serif, int bold, int italic) {
    (void)script; (void)language; (void)serif; (void)bold; (void)italic;
#ifdef _WIN32
    const char *fontpath = "C:\\Windows\\Fonts\\simsun.ttc";
#elif defined(__APPLE__)
    const char *fontpath = "/System/Library/Fonts/PingFang.ttc";
#else
    const char *fontpath = "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc";
#endif
    fz_font *font = NULL;
    fz_try(ctx) {
        font = fz_new_font_from_file(ctx, NULL, fontpath, 0, 0);
    }
    fz_catch(ctx) { font = NULL; }
    return font;
}

// ============================================================
// Context management
// ============================================================

static fz_context* gomupdf_new_context(void) {
    fz_context *ctx = fz_new_context(NULL, NULL, FZ_STORE_DEFAULT);
    if (ctx) {
        fz_install_load_system_font_funcs(ctx,
            gomupdf_load_system_font,
            gomupdf_load_system_cjk_font,
            gomupdf_load_system_fallback_font);
    }
    return ctx;
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
    fz_try(ctx) {
        pdf_subset_fonts(ctx, pdf, 0, NULL);
        pdf_save_document(ctx, pdf, filename, &opts);
    }
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
        pdf_subset_fonts(ctx, pdf, 0, NULL);
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
        // 不使用 pdf_page_write —— 它会注入 Y-flip CTM 到初始内容流，
        // 导致后续追加的 InsertText/InsertImage 内容流全部镜像翻转。
        // 直接创建空白页，让所有内容流在标准 PDF 坐标系下工作。
        pdf_obj *resources = pdf_new_dict(ctx, pdf, 2);
        fz_buffer *contents = fz_new_buffer(ctx, 1);
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
// Text insertion (Shape) — CJK helpers
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

/* Append a UTF-8 string as a PDF hex string with 2-byte CIDs for Identity-H.
   Identity-H encoding: CID = Unicode code point directly, no BOM needed. */
static void gomupdf_append_cid_hex(fz_context *ctx, fz_buffer *buf, const char *text) {
    static const char hex[] = "0123456789ABCDEF";
    fz_append_byte(ctx, buf, '<');
    // Identity-H 编码下不需要 BOM，每个 2 字节直接是 CID = Unicode 码点
    const unsigned char *s = (const unsigned char *)text;
    while (*s) {
        unsigned int cp = gomupdf_utf8_decode(&s);
        if (cp > 0xFFFF) cp = 0xFFFD;
        fz_append_byte(ctx, buf, hex[(cp >> 12) & 0xF]);
        fz_append_byte(ctx, buf, hex[(cp >>  8) & 0xF]);
        fz_append_byte(ctx, buf, hex[(cp >>  4) & 0xF]);
        fz_append_byte(ctx, buf, hex[ cp        & 0xF]);
    }
    fz_append_byte(ctx, buf, '>');
}

/* Detect the best CJK ordering for a UTF-8 string by scanning codepoint ranges.
   Returns FZ_ADOBE_GB (Simplified Chinese) as default for any CJK text. */
static int gomupdf_detect_cjk_ordering(const char *text) {
    int has_jp = 0, has_kr = 0, has_tc = 0;
    const unsigned char *s = (const unsigned char *)text;
    while (*s) {
        unsigned int cp = gomupdf_utf8_decode(&s);
        if ((cp >= 0x3040 && cp <= 0x309F) || (cp >= 0x30A0 && cp <= 0x30FF))
            has_jp = 1;
        else if ((cp >= 0xAC00 && cp <= 0xD7AF) || (cp >= 0x1100 && cp <= 0x11FF))
            has_kr = 1;
        else if (cp >= 0x3100 && cp <= 0x312F)
            has_tc = 1;
    }
    if (has_jp) return FZ_ADOBE_JAPAN;
    if (has_kr) return FZ_ADOBE_KOREA;
    if (has_tc) return FZ_ADOBE_CNS;
    return FZ_ADOBE_GB;
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

/* 返回各 CJK ordering 对应的标准 PDF 字体名，PDF 阅读器可据此做系统字体替换 */
static const char* gomupdf_cjk_font_name(int ordering) {
    switch (ordering) {
        case FZ_ADOBE_CNS:   return "MHei-Medium";
        case FZ_ADOBE_GB:    return "STSong-Light";
        case FZ_ADOBE_JAPAN: return "KozMinPr6N-Regular";
        case FZ_ADOBE_KOREA: return "HYSMyeongJoStd-Medium";
        default:             return "STSong-Light";
    }
}

/* Create a non-embedded CID font for CJK text.
   Uses standard CJK font names (STSong-Light etc.) that PDF readers can substitute.
   Identity-H encoding: CID = Unicode code point directly. */
static pdf_obj* gomupdf_create_cjk_font(fz_context *ctx, pdf_document *doc, int ordering) {
    const char *font_name = gomupdf_cjk_font_name(ordering);
    const char *ord_name = gomupdf_cjk_ordering_name(ordering);
    int supplement = gomupdf_cjk_supplement(ordering);

    /* CIDSystemInfo - 使用 Identity ordering，CID 直接等于 Unicode 码点 */
    pdf_obj *sysinfo = pdf_new_dict(ctx, doc, 3);
    pdf_dict_put_text_string(ctx, sysinfo, PDF_NAME(Registry), "Adobe");
    pdf_dict_put_text_string(ctx, sysinfo, PDF_NAME(Ordering), "Identity");
    pdf_dict_put_int(ctx, sysinfo, PDF_NAME(Supplement), 0);

    /* Type0 层级的 CIDSystemInfo（带正确的 ordering，帮助阅读器做字体替换） */
    pdf_obj *sysinfo2 = pdf_new_dict(ctx, doc, 3);
    pdf_dict_put_text_string(ctx, sysinfo2, PDF_NAME(Registry), "Adobe");
    pdf_dict_put_text_string(ctx, sysinfo2, PDF_NAME(Ordering), ord_name);
    pdf_dict_put_int(ctx, sysinfo2, PDF_NAME(Supplement), supplement);

    /* CIDFont 字典 - 使用 CIDFontType2 (TrueType-based) */
    pdf_obj *cidfont = pdf_new_dict(ctx, doc, 6);
    pdf_dict_put(ctx, cidfont, PDF_NAME(Type), PDF_NAME(Font));
    pdf_dict_put_name(ctx, cidfont, PDF_NAME(Subtype), "CIDFontType2");
    pdf_dict_put_name(ctx, cidfont, PDF_NAME(BaseFont), font_name);
    pdf_dict_put(ctx, cidfont, PDF_NAME(CIDSystemInfo), sysinfo);
    pdf_dict_put_int(ctx, cidfont, PDF_NAME(DW), 1000);
    pdf_dict_put(ctx, cidfont, PDF_NAME(CIDToGIDMap), PDF_NAME(Identity));

    /* FontDescriptor */
    pdf_obj *fd = pdf_new_dict(ctx, doc, 10);
    pdf_dict_put(ctx, fd, PDF_NAME(Type), PDF_NAME(FontDescriptor));
    pdf_dict_put_name(ctx, fd, PDF_NAME(FontName), font_name);
    pdf_dict_put_int(ctx, fd, PDF_NAME(Flags), 6);
    pdf_obj *bbox = pdf_new_array(ctx, doc, 4);
    pdf_array_push_int(ctx, bbox, -200);
    pdf_array_push_int(ctx, bbox, -200);
    pdf_array_push_int(ctx, bbox, 1200);
    pdf_array_push_int(ctx, bbox, 1000);
    pdf_dict_put(ctx, fd, PDF_NAME(FontBBox), bbox);
    pdf_drop_obj(ctx, bbox);
    pdf_dict_put_int(ctx, fd, PDF_NAME(ItalicAngle), 0);
    pdf_dict_put_int(ctx, fd, PDF_NAME(Ascent), 880);
    pdf_dict_put_int(ctx, fd, PDF_NAME(Descent), -120);
    pdf_dict_put_int(ctx, fd, PDF_NAME(StemV), 80);
    pdf_dict_put_int(ctx, fd, PDF_NAME(CapHeight), 700);
    pdf_obj *fd_ref = pdf_add_object(ctx, doc, fd);
    pdf_drop_obj(ctx, fd);
    pdf_dict_put(ctx, cidfont, PDF_NAME(FontDescriptor), fd_ref);
    pdf_drop_obj(ctx, fd_ref);

    pdf_obj *cidfont_ref = pdf_add_object(ctx, doc, cidfont);
    pdf_drop_obj(ctx, cidfont);
    pdf_drop_obj(ctx, sysinfo);

    pdf_obj *descendants = pdf_new_array(ctx, doc, 1);
    pdf_array_push(ctx, descendants, cidfont_ref);
    pdf_drop_obj(ctx, cidfont_ref);

    /* ToUnicode CMap - Identity-H 下 CID = Unicode，映射是 1:1 的 */
    const char *tounicode_str =
        "/CIDInit /ProcSet findresource begin\n"
        "12 dict begin\n"
        "begincmap\n"
        "/CIDSystemInfo\n"
        "<< /Registry (Adobe) /Ordering (UCS) /Supplement 0 >> def\n"
        "/CMapName /Adobe-Identity-UCS def\n"
        "/CMapType 2 def\n"
        "1 begincodespacerange\n"
        "<0000> <FFFF>\n"
        "endcodespacerange\n"
        "1 beginbfrange\n"
        "<0000> <FFFF> <0000>\n"
        "endbfrange\n"
        "endcmap\n"
        "CMapName currentdict /CMap defineresource pop\n"
        "end\n"
        "end\n";
    fz_buffer *tounicode_buf = fz_new_buffer_from_copied_data(ctx,
        (const unsigned char *)tounicode_str, strlen(tounicode_str));
    pdf_obj *tounicode_ref = pdf_add_stream(ctx, doc, tounicode_buf, NULL, 0);
    fz_drop_buffer(ctx, tounicode_buf);

    /* Type0 复合字体字典 */
    pdf_obj *fontdict = pdf_new_dict(ctx, doc, 6);
    pdf_dict_put(ctx, fontdict, PDF_NAME(Type), PDF_NAME(Font));
    pdf_dict_put_name(ctx, fontdict, PDF_NAME(Subtype), "Type0");
    pdf_dict_put_name(ctx, fontdict, PDF_NAME(BaseFont), font_name);
    pdf_dict_put_name(ctx, fontdict, PDF_NAME(Encoding), "Identity-H");
    pdf_dict_put(ctx, fontdict, PDF_NAME(DescendantFonts), descendants);
    pdf_drop_obj(ctx, descendants);
    pdf_dict_put(ctx, fontdict, PDF_NAME(ToUnicode), tounicode_ref);
    pdf_drop_obj(ctx, tounicode_ref);
    pdf_dict_put(ctx, fontdict, PDF_NAME(CIDSystemInfo), sysinfo2);
    pdf_drop_obj(ctx, sysinfo2);

    pdf_obj *fontdict_ref = pdf_add_object(ctx, doc, fontdict);
    pdf_drop_obj(ctx, fontdict);
    return fontdict_ref;
}

// ============================================================
// Text insertion
// ============================================================

static int gomupdf_insert_text(fz_context *ctx, pdf_document *doc, int pno,
    float x, float y, const char *text, const char *fontname, float fontsize,
    float r, float g, float b) {
    int errcode = 0;
    fz_try(ctx) {
        int use_cjk = gomupdf_text_needs_cjk(text);
        pdf_obj *page_obj = pdf_lookup_page_obj(ctx, doc, pno);

        pdf_obj *resources = pdf_dict_get(ctx, page_obj, PDF_NAME(Resources));
        if (!resources)
            resources = pdf_dict_put_dict(ctx, page_obj, PDF_NAME(Resources), 2);
        pdf_obj *fonts = pdf_dict_get(ctx, resources, PDF_NAME(Font));
        if (!fonts)
            fonts = pdf_dict_put_dict(ctx, resources, PDF_NAME(Font), 4);

        char fname[32];
        snprintf(fname, sizeof(fname), "F%d", pdf_create_object(ctx, doc));

        if (use_cjk) {
            int ordering = gomupdf_detect_cjk_ordering(text);
            pdf_obj *font_obj = gomupdf_create_cjk_font(ctx, doc, ordering);
            pdf_dict_puts(ctx, fonts, fname, font_obj);
            pdf_drop_obj(ctx, font_obj);
        } else {
            fz_font *font = fz_new_base14_font(ctx, fontname);
            pdf_obj *font_obj = pdf_add_simple_font(ctx, doc, font, PDF_SIMPLE_ENCODING_LATIN);
            pdf_dict_puts(ctx, fonts, fname, font_obj);
            pdf_drop_obj(ctx, font_obj);
            fz_drop_font(ctx, font);
        }

        /* Get page mediabox height for coordinate conversion */
        fz_rect mediabox;
        pdf_obj *mb = pdf_dict_get(ctx, page_obj, PDF_NAME(MediaBox));
        if (mb) {
            mediabox = pdf_to_rect(ctx, mb);
        } else {
            mediabox.x0 = 0; mediabox.y0 = 0;
            mediabox.x1 = 612; mediabox.y1 = 792;
        }
        float page_height = mediabox.y1 - mediabox.y0;
        float pdf_y = page_height - y;  /* Go API 左上角原点 → PDF 左下角原点 */

        /* Build content stream — 不使用 cm 变换，直接在原生 PDF 坐标系下定位 */
        fz_buffer *content = fz_new_buffer(ctx, 256);
        fz_append_printf(ctx, content, "q\n");
        fz_append_printf(ctx, content, "BT\n");
        fz_append_printf(ctx, content, "%g %g %g rg\n", r, g, b);
        fz_append_printf(ctx, content, "/%s %g Tf\n", fname, fontsize);
        fz_append_printf(ctx, content, "%g %g Td\n", x, pdf_y);

        if (use_cjk) {
            gomupdf_append_cid_hex(ctx, content, text);
            fz_append_string(ctx, content, " Tj\n");
        } else {
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

        /* Get page height for Y coordinate conversion */
        fz_rect mediabox;
        fz_matrix page_ctm;
        pdf_page_obj_transform(ctx, page_obj, &mediabox, &page_ctm);
        float page_height = mediabox.y1 - mediabox.y0;

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

        /* Y coordinate conversion: 图片顶边在 PDF 坐标系中的位置
           高度取负值使图片正向显示（PDF cm 矩阵中负高度 = 向下绘制） */
        float img_w = rect.x1 - rect.x0;
        float img_h = rect.y1 - rect.y0;
        float pdf_x = rect.x0;
        float pdf_y = page_height - rect.y0;

        fz_buffer *content = fz_new_buffer(ctx, 256);
        fz_append_printf(ctx, content, "q %g 0 0 %g %g %g cm /%s Do Q\n",
            img_w, -img_h, pdf_x, pdf_y, name);

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
// HTML box insertion (Story-based)
// ============================================================

/* Insert HTML content into a rectangle on an existing PDF page.
   Uses MuPDF's Story API to layout styled HTML/CSS into the given rect.
   Returns: 0 = all content fitted, 1 = error, 2 = content overflow (didn't fit).
   spare_height: if non-NULL, receives the remaining height in the rect.
   scale_used: if non-NULL, receives the actual scale factor used. */
static int gomupdf_insert_htmlbox(fz_context *ctx, pdf_document *doc, int pno,
    float x0, float y0, float x1, float y1,
    const char *html, const char *css, float scale_low, int overlay,
    float *spare_height, float *scale_used) {
    int errcode = 0;
    fz_try(ctx) {
        pdf_obj *page_obj = pdf_lookup_page_obj(ctx, doc, pno);

        /* Get page mediabox */
        fz_rect mediabox;
        fz_matrix page_ctm;
        pdf_page_obj_transform(ctx, page_obj, &mediabox, &page_ctm);

        /* 不做 Y 坐标转换 — pdf_page_write 的 device 已经自带 Y-flip CTM，
           所以 Go API 的左上角原点坐标可以直接传给 fz_place_story。 */
        fz_rect where = {x0, y0, x1, y1};

        float rect_w = where.x1 - where.x0;
        float rect_h = where.y1 - where.y0;
        if (rect_w <= 0 || rect_h <= 0) {
            if (spare_height) *spare_height = 0;
            if (scale_used) *scale_used = 1.0f;
            fz_throw(ctx, FZ_ERROR_ARGUMENT, "invalid rectangle");
        }

        /* Create the story from HTML + CSS */
        fz_buffer *html_buf = fz_new_buffer_from_copied_data(ctx,
            (const unsigned char *)html, strlen(html));
        fz_story *story = fz_new_story(ctx, html_buf, css, 12.0f, NULL);
        fz_drop_buffer(ctx, html_buf);

        /* Try to place the story; if it doesn't fit, optionally scale down */
        float scale = 1.0f;
        int more = 0;
        fz_rect filled = fz_empty_rect;

        more = fz_place_story(ctx, story, where, &filled);

        if (more && scale_low < 1.0f) {
            /* Content didn't fit — try scaling down */
            float lo = (scale_low > 0) ? scale_low : 0.05f;
            float hi = 1.0f;
            /* Binary search for the right scale */
            for (int iter = 0; iter < 20; iter++) {
                float mid = (lo + hi) / 2.0f;
                fz_reset_story(ctx, story);
                fz_rect scaled_where = {
                    where.x0, where.y0,
                    where.x0 + rect_w / mid,
                    where.y0 + rect_h / mid
                };
                fz_rect test_filled = fz_empty_rect;
                int test_more = fz_place_story(ctx, story, scaled_where, &test_filled);
                if (test_more)
                    hi = mid;
                else {
                    lo = mid;
                    if (hi - lo < 0.005f) break;
                }
            }
            scale = lo;
            /* Final placement at the found scale */
            fz_reset_story(ctx, story);
            fz_rect scaled_where = {
                where.x0, where.y0,
                where.x0 + rect_w / scale,
                where.y0 + rect_h / scale
            };
            filled = fz_empty_rect;
            more = fz_place_story(ctx, story, scaled_where, &filled);
        }

        /* Draw the story to a pdf_page_write device to capture content stream */
        pdf_obj *resources = NULL;
        fz_buffer *contents = NULL;
        fz_device *dev = pdf_page_write(ctx, doc, mediabox, &resources, &contents);

        fz_matrix draw_ctm = fz_identity;
        if (scale < 1.0f) {
            /* Scale around the top-left corner of the target rect */
            draw_ctm = fz_concat(
                fz_translate(-where.x0, -where.y0),
                fz_concat(fz_scale(scale, scale), fz_translate(where.x0, where.y0))
            );
        }

        fz_draw_story(ctx, story, dev, draw_ctm);
        fz_close_device(ctx, dev);
        fz_drop_device(ctx, dev);

        /* Merge resources from the story into the page's resources */
        pdf_obj *page_resources = pdf_dict_get(ctx, page_obj, PDF_NAME(Resources));
        if (!page_resources)
            page_resources = pdf_dict_put_dict(ctx, page_obj, PDF_NAME(Resources), 4);

        /* Merge each resource category (Font, XObject, ExtGState, etc.) */
        static const pdf_obj *res_keys[] = {
            PDF_NAME(Font), PDF_NAME(XObject), PDF_NAME(ExtGState),
            PDF_NAME(ColorSpace), PDF_NAME(Pattern), PDF_NAME(Shading),
            PDF_NAME(Properties), NULL
        };
        for (int ki = 0; res_keys[ki] != NULL; ki++) {
            pdf_obj *key = (pdf_obj *)res_keys[ki];
            pdf_obj *src_dict = pdf_dict_get(ctx, resources, key);
            if (!src_dict) continue;
            pdf_obj *dst_dict = pdf_dict_get(ctx, page_resources, key);
            if (!dst_dict)
                dst_dict = pdf_dict_put_dict(ctx, page_resources, key, 4);
            int n = pdf_dict_len(ctx, src_dict);
            for (int i = 0; i < n; i++) {
                pdf_obj *k = pdf_dict_get_key(ctx, src_dict, i);
                pdf_obj *v = pdf_dict_get_val(ctx, src_dict, i);
                pdf_dict_put(ctx, dst_dict, k, v);
            }
        }

        /* Append the content stream to the page, wrapped in q/Q to isolate
           pdf_page_write 注入的 Y-flip CTM，防止泄漏到同一页面的其他内容流。 */
        fz_buffer *wrapped = fz_new_buffer(ctx, 2 + (int)fz_buffer_storage(ctx, contents, NULL) + 3);
        fz_append_string(ctx, wrapped, "q\n");
        unsigned char *cdata;
        size_t clen = fz_buffer_storage(ctx, contents, &cdata);
        fz_append_data(ctx, wrapped, cdata, clen);
        fz_append_string(ctx, wrapped, "\nQ\n");

        pdf_obj *existing = pdf_dict_get(ctx, page_obj, PDF_NAME(Contents));
        pdf_obj *newstream = pdf_add_stream(ctx, doc, wrapped, NULL, 0);
        if (pdf_is_array(ctx, existing)) {
            if (overlay)
                pdf_array_push(ctx, existing, newstream);
            else
                pdf_array_insert(ctx, existing, newstream, 0);
        } else {
            pdf_obj *arr = pdf_new_array(ctx, doc, 2);
            if (existing) {
                if (overlay) {
                    pdf_array_push(ctx, arr, existing);
                    pdf_array_push(ctx, arr, newstream);
                } else {
                    pdf_array_push(ctx, arr, newstream);
                    pdf_array_push(ctx, arr, existing);
                }
            } else {
                pdf_array_push(ctx, arr, newstream);
            }
            pdf_dict_put(ctx, page_obj, PDF_NAME(Contents), arr);
            pdf_drop_obj(ctx, arr);
        }
        pdf_drop_obj(ctx, newstream);

        /* Calculate spare height */
        if (spare_height) {
            float used_h = filled.y1 - filled.y0;
            *spare_height = rect_h - used_h * scale;
            if (*spare_height < 0) *spare_height = 0;
        }
        if (scale_used) *scale_used = scale;

        fz_drop_buffer(ctx, wrapped);
        fz_drop_buffer(ctx, contents);
        pdf_drop_obj(ctx, resources);
        fz_drop_story(ctx, story);

        if (more && scale_low >= 1.0f) {
            /* Content didn't fit and no scaling allowed — report overflow */
            errcode = 2;
        }
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

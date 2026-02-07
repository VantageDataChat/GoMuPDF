# GoMuPDF `gomupdf.h` 修改建议

## 概述

修复 `InsertText` / `InsertImage` 产生的两类显示问题：

1. **文字和图片镜像翻转 + 上下颠倒** — 坐标系冲突
2. **CJK 中文文字乱码** — CID 字体字典编码错误

共涉及 **4 处修改**，影响 4 个函数 + 1 个辅助函数。

---

## 问题一：坐标系冲突导致内容镜像翻转

### 根因分析

`gomupdf_insert_page` 使用 `pdf_page_write()` 创建新页面。该 MuPDF API 会在初始内容流中注入一个 Y 轴翻转矩阵：

```
1 0 0 -1 0 842 cm
```

这将坐标系从 PDF 原生的"左下角原点，Y 向上"变为"左上角原点，Y 向下"。

**问题在于**：PDF 规范中，多个内容流共享同一图形状态。这个 CTM 变换在初始内容流结束后仍然生效，后续通过 `InsertText` / `InsertImage` 追加的内容流都在这个已翻转的坐标系下运行，但它们并不知道翻转的存在，导致所有内容水平镜像 + 上下颠倒。

### 修改 1：`gomupdf_insert_page` — 不使用 `pdf_page_write`

**原代码：**

```c
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
```

**修改后：**

```c
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
```

### 修改 2：`gomupdf_insert_text` — 坐标转换改为直接计算

去掉内容流中的 `cm` 翻转变换，改为在 C 代码中直接计算 PDF 原生 Y 坐标。

**原代码（内容流构建部分）：**

```c
float page_height = mediabox.y1 - mediabox.y0;

fz_buffer *content = fz_new_buffer(ctx, 256);
fz_append_printf(ctx, content, "q\n");
fz_append_printf(ctx, content, "1 0 0 -1 0 %g cm\n", page_height);
fz_append_printf(ctx, content, "BT\n");
fz_append_printf(ctx, content, "%g %g %g rg\n", r, g, b);
fz_append_printf(ctx, content, "/%s %g Tf\n", fname, fontsize);
fz_append_printf(ctx, content, "%g %g Td\n", x, y);
```

**修改后：**

```c
float page_height = mediabox.y1 - mediabox.y0;
float pdf_y = page_height - y;  // Go API 左上角原点 → PDF 左下角原点

fz_buffer *content = fz_new_buffer(ctx, 256);
fz_append_printf(ctx, content, "q\n");
// 不再使用 cm 变换，直接在原生 PDF 坐标系下定位
fz_append_printf(ctx, content, "BT\n");
fz_append_printf(ctx, content, "%g %g %g rg\n", r, g, b);
fz_append_printf(ctx, content, "/%s %g Tf\n", fname, fontsize);
fz_append_printf(ctx, content, "%g %g Td\n", x, pdf_y);
```

### 修改 3：`gomupdf_insert_image` — 图片坐标同样转换

**原代码：**

```c
fz_buffer *content = fz_new_buffer(ctx, 256);
fz_append_printf(ctx, content, "q %g 0 0 %g %g %g cm /%s Do Q\n",
    rect.x1 - rect.x0, rect.y1 - rect.y0, rect.x0, rect.y0, name);
```

**修改后：**

```c
// 需要先获取 page_height（在函数前面添加）
fz_rect mediabox;
fz_matrix page_ctm;
pdf_page_obj_transform(ctx, page_obj, &mediabox, &page_ctm);
float page_height = mediabox.y1 - mediabox.y0;

// ...（资源字典处理代码不变）...

float img_w = rect.x1 - rect.x0;
float img_h = rect.y1 - rect.y0;
float pdf_x = rect.x0;
float pdf_y = page_height - rect.y0;  // 图片顶边在 PDF 坐标系中的位置

fz_buffer *content = fz_new_buffer(ctx, 256);
// 高度取负值使图片正向显示（PDF cm 矩阵中负高度 = 向下绘制）
fz_append_printf(ctx, content, "q %g 0 0 %g %g %g cm /%s Do Q\n",
    img_w, -img_h, pdf_x, pdf_y, name);
```

---

## 问题二：CJK 中文文字乱码

### 根因分析

原实现的 CID 字体字典结构：

```
Type0 Font
├── BaseFont: Adobe-Identity        ← 不是有效的 PDF 字体名
├── Encoding: UniGB-UTF16-H         ← CMap 期望 Adobe-GB1 CID，不是原始 Unicode
└── DescendantFonts:
    └── CIDFontType0                ← PostScript-based CID，兼容性差
        ├── BaseFont: Adobe-Identity
        └── CIDSystemInfo: Adobe-GB1-5
```

问题：
- `Adobe-Identity` 不是 PDF 阅读器能识别的字体名，无法做字体替换
- `UniGB-UTF16-H` CMap 将 UTF-16 码点映射到 Adobe-GB1 字符集的 CID，但字体声明的是 `Adobe-Identity`，CID 映射不匹配
- `CIDFontType0`（PostScript-based）在非嵌入场景下兼容性不如 `CIDFontType2`
- hex 字符串开头的 BOM (`FEFF`) 会被当作 CID 0xFEFF 渲染

### 修改 4：`gomupdf_create_cjk_font` + `gomupdf_append_utf16be_hex`

#### 4a. 新增辅助函数 `gomupdf_cjk_font_name`

```c
/* 返回各 CJK ordering 对应的标准 PDF 字体名，PDF 阅读器可据此做系统字体替换 */
static const char* gomupdf_cjk_font_name(int ordering) {
    switch (ordering) {
        case FZ_ADOBE_CNS:   return "MHei-Medium";           /* 繁体中文 */
        case FZ_ADOBE_GB:    return "STSong-Light";           /* 简体中文 */
        case FZ_ADOBE_JAPAN: return "KozMinPr6N-Regular";     /* 日文 */
        case FZ_ADOBE_KOREA: return "HYSMyeongJoStd-Medium";  /* 韩文 */
        default:             return "STSong-Light";
    }
}
```

#### 4b. 重写 `gomupdf_create_cjk_font`

新的字体字典结构：

```
Type0 Font
├── BaseFont: STSong-Light          ← 标准 CJK 字体名，阅读器可替换
├── Encoding: Identity-H            ← 直接映射，CID = Unicode 码点
├── ToUnicode: CMap stream          ← 支持文本提取和搜索
└── DescendantFonts:
    └── CIDFontType2                ← TrueType-based CID，兼容性好
        ├── BaseFont: STSong-Light
        ├── CIDSystemInfo: Adobe-Identity-0
        ├── CIDToGIDMap: Identity
        ├── DW: 1000
        └── FontDescriptor: (完整描述符)
```

**完整代码：**

```c
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

    /* ToUnicode CMap - Identity-H 下 CID = Unicode，所以映射是 1:1 的 */
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
```

#### 4c. `gomupdf_append_utf16be_hex` — 去掉 BOM

```c
// 原代码：
fz_append_byte(ctx, buf, '<');
fz_append_string(ctx, buf, "FEFF");  // ← 删掉这行
const unsigned char *s = ...

// 改为：
fz_append_byte(ctx, buf, '<');
// Identity-H 编码下不需要 BOM，每个 2 字节直接是 CID = Unicode 码点
const unsigned char *s = ...
```

---

## 修改汇总

| # | 函数 | 改动 | 原因 |
|---|------|------|------|
| 1 | `gomupdf_insert_page` | `pdf_page_write` → `pdf_new_dict` + `fz_new_buffer` | `pdf_page_write` 注入 Y-flip CTM 污染后续内容流 |
| 2 | `gomupdf_insert_text` | 删除 `1 0 0 -1 0 h cm`，改用 `page_height - y` | 在原生 PDF 坐标系下直接定位 |
| 3 | `gomupdf_insert_image` | 添加 Y 坐标转换，高度取负值 | 同上，图片需要负高度才能正向显示 |
| 4 | `gomupdf_create_cjk_font` | 重写字体字典（见上方详细代码） | 修复 CJK 乱码 |
| 4+ | `gomupdf_append_utf16be_hex` | 删除 BOM `FEFF` | Identity-H 下 BOM 被当作无效 CID |
| 4+ | 新增 `gomupdf_cjk_font_name` | 返回标准 CJK 字体名 | 供 `gomupdf_create_cjk_font` 使用 |

修改 1-3 是一组关联改动（坐标系统一），修改 4 是独立的 CJK 编码修复。两组可以分开提交。

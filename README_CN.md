[English](README.md) | 中文

# GoMuPDF

MuPDF 的 Go 语言绑定 — 高性能 PDF（及其他文档格式）数据提取、分析、转换和操作库。

本项目是 [PyMuPDF](https://github.com/pymupdf/PyMuPDF) 的 Go 语言移植版，提供符合 Go 惯用风格的 API，功能对标 PyMuPDF。

## 支持平台

| 操作系统 | amd64 | arm64 |
|----------|-------|-------|
| Windows  | ✅    | ✅    |
| Linux    | ✅    | ✅    |
| macOS    | ✅    | ✅    |

提供 `nomupdf` 构建标签，可在无 CGO/MuPDF 环境下编译（桩函数返回 `ErrInitFailed`）。

## 支持的文档格式

PDF、XPS、OpenXPS、CBZ、CBR、FB2、EPUB 及图片格式（JPEG、PNG、BMP、GIF、TIFF 等）

## 前置要求

- Go 1.21+
- C 编译器：GCC（Windows 下使用 MinGW-W64）或 Clang
- MuPDF 1.24.9 静态库（参见下方构建说明）

## 构建 MuPDF 静态库

GoMuPDF 链接 `libs/` 目录下的平台特定静态库，命名规则：

```
libs/libmupdf_<os>_<arch>.a
libs/libmupdfthird_<os>_<arch>.a
```

例如：`libmupdf_linux_amd64.a`、`libmupdf_darwin_arm64.a`、`libmupdf_windows_amd64.a`。

### 自动构建

```bash
# Linux / macOS
chmod +x build_libs.sh
./build_libs.sh

# 指定已有的 MuPDF 源码目录：
./build_libs.sh /path/to/mupdf-src

# Windows (MinGW)
build_libs.bat
build_libs.bat D:\mupdf-src
```

### 手动构建

```bash
git clone --depth 1 --branch 1.24.9 --recurse-submodules \
    https://github.com/ArtifexSoftware/mupdf.git mupdf-src
cd mupdf-src
make -j$(nproc) HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no \
    USE_SYSTEM_LIBS=no XCFLAGS="-fPIC" libs
cd ..

# 复制到 libs/ 并按命名规则重命名
mkdir -p libs
cp mupdf-src/build/release/libmupdf.a libs/libmupdf_linux_amd64.a
cp mupdf-src/build/release/libmupdf-third.a libs/libmupdfthird_linux_amd64.a
```

### 使用预编译库

也可以从 [go-fitz](https://github.com/gen2brain/go-fitz) v1.24.15 提取静态库，按上述命名规则重命名即可。

## 安装

```bash
go get github.com/VantageDataChat/GoMuPDF
```

## 快速开始

```go
package main

import (
    "fmt"
    "github.com/VantageDataChat/GoMuPDF"
)

func main() {
    doc, err := gomupdf.Open("example.pdf")
    if err != nil {
        panic(err)
    }
    defer doc.Close()

    fmt.Printf("页数: %d\n", doc.PageCount())

    for i := 0; i < doc.PageCount(); i++ {
        page, err := doc.LoadPage(i)
        if err != nil {
            continue
        }
        text, _ := page.GetText("text")
        fmt.Printf("第 %d 页: %s\n", i, text)
        page.Close()
    }
}
```

## 插入 HTML 内容

```go
page, _ := doc.NewPage(-1, 595, 842)
result, err := page.InsertHTMLBox(
    gomupdf.Rect{X0: 50, Y0: 50, X1: 500, Y1: 400},
    `<p style="color:red; font-size:16px;">你好 <b>世界</b></p>`,
)
// result.SpareHeight: 剩余高度
// result.Scale: 实际缩放比例
```

支持 HTML/CSS 样式、CJK 中日韩文字、自动缩放适配，保存时自动进行字体子集化以减小文件体积。

## 插入文本

```go
// CJK 文本（非嵌入字体，PDF 阅读器替换显示）
page.InsertText(gomupdf.Point{X: 72, Y: 72}, "你好世界")

// 自定义字体和颜色
page.InsertText(gomupdf.Point{X: 72, Y: 100}, "Hello",
    gomupdf.WithFontName("helv"),
    gomupdf.WithFontSize(14),
    gomupdf.WithColor(gomupdf.ColorRed),
)
```

## 测试

```bash
# 完整测试（需要当前平台的 MuPDF 静态库）
go test -v -count=1 .

# 纯 Go 测试（无需 CGO/静态库，任意平台可用）
go test -v -count=1 -tags nomupdf .
```

## 核心类型

| Go 类型 | PyMuPDF 对应类型 | 说明 |
|---------|-----------------|------|
| `Document` | `fitz.Document` | 文档对象 |
| `Page` | `fitz.Page` | 页面对象 |
| `Pixmap` | `fitz.Pixmap` | 像素图（光栅图像） |
| `Rect` | `fitz.Rect` | 矩形 |
| `Matrix` | `fitz.Matrix` | 变换矩阵 |
| `Point` | `fitz.Point` | 二维点 |
| `Quad` | `fitz.Quad` | 四边形 |
| `TextPage` | `fitz.TextPage` | 文本提取结果 |
| `Annot` | `fitz.Annot` | 注释 |
| `Widget` | `fitz.Widget` | 表单控件 |

## 构建标签

- 默认：启用 CGO，链接 MuPDF 静态库
- `nomupdf`：禁用 CGO，所有 MuPDF 函数返回 `ErrInitFailed`，适用于 CI 或无 MuPDF 的环境

## 许可证

AGPL-3.0（与 MuPDF/PyMuPDF 一致）

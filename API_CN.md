中文 | [English](API_EN.md)

# GoMuPDF API 参考文档

包 `gomupdf` 提供 MuPDF 的 Go 语言绑定。

所有坐标使用左上角原点坐标系（与 PyMuPDF 一致）。单位为 PDF 点（1 点 = 1/72 英寸）。

---

## 包级函数

### 文档创建与打开

```go
func Open(filename string) (*Document, error)
```
从文件路径打开文档。支持 PDF、XPS、EPUB、CBZ、FB2 及图片格式。

```go
func OpenFromMemory(data []byte, magic string) (*Document, error)
```
从字节切片打开文档。`magic` 为 MIME 类型或扩展名提示（如 `"application/pdf"`、`".pdf"`）。

```go
func NewPDF() (*Document, error)
```
创建新的空白 PDF 文档。

### 像素图创建

```go
func NewPixmap(colorspace int, width, height int, alpha bool) (*Pixmap, error)
```
创建空白像素图。`colorspace`：`CsGray`、`CsRGB` 或 `CsCMYK`。

```go
func NewPixmapFromImage(doc *Document, xref int) (*Pixmap, error)
```
从文档中的图像对象（按 xref 编号）创建像素图。

### 几何构造函数

```go
func NewPoint(x, y float64) Point
func NewRect(x0, y0, x1, y1 float64) Rect
func RectFromPoints(topLeft, bottomRight Point) Rect
func NewIRect(x0, y0, x1, y1 int) IRect
func NewQuad(ul, ur, ll, lr Point) Quad
func QuadFromRect(r Rect) Quad
```

### 矩阵构造函数

```go
func NewMatrix(a, b, c, d, e, f float64) Matrix
func ScaleMatrix(sx, sy float64) Matrix
func TranslateMatrix(tx, ty float64) Matrix
func RotateMatrix(deg float64) Matrix
func ShearMatrix(sx, sy float64) Matrix
```

### 工具函数

```go
func Version() string          // 返回 GoMuPDF 版本号 ("0.1.0")
func MuPDFVersion() string     // 返回 MuPDF 版本号 ("1.24.9")
func PaperSize(name string) Rect  // 返回指定纸张尺寸的 Rect（"a4"、"letter"、"a3-landscape" 等）
func GetPDFStr(s string) string   // 将字符串转换为 PDF 字符串格式（处理 Unicode）
```

---

## Document（文档）

表示已打开的文档（PDF、XPS、EPUB 等）。

### 属性

```go
func (d *Document) PageCount() int      // 页数
func (d *Document) Name() string        // 文件名
func (d *Document) IsPDF() bool         // 是否为 PDF
func (d *Document) IsClosed() bool      // 是否已关闭
func (d *Document) NeedsPass() bool     // 是否需要密码
func (d *Document) IsReflowable() bool  // 是否可重排（EPUB 等）
```

### 生命周期

```go
func (d *Document) Close()
func (d *Document) Authenticate(password string) (int, error)
```

### 页面访问

```go
func (d *Document) LoadPage(pageNum int) (*Page, error)    // 0 起始，支持负数索引
func (d *Document) Pages(args ...int) ([]*Page, error)     // 参数：start, stop, step
```

### 元数据

```go
func (d *Document) Metadata() map[string]string
func (d *Document) SetMetadata(meta map[string]string) error
```
支持的键：`"title"`、`"author"`、`"subject"`、`"keywords"`、`"creator"`、`"producer"`、`"creationDate"`、`"modDate"`。

### 目录

```go
func (d *Document) GetTOC(simple bool) ([]TOCItem, error)
```

### 保存与导出

```go
func (d *Document) Save(filename string, opts ...SaveOptions) error
func (d *Document) EzSave(filename string) error              // Garbage=3, Deflate=true，最小体积
func (d *Document) ToBytes(opts ...SaveOptions) ([]byte, error)
func (d *Document) ConvertToPDF(fromPage, toPage, rotate int) ([]byte, error)
func (d *Document) CanSaveIncrementally() bool
```

保存时自动执行字体子集化 — 仅嵌入实际使用的字形，大幅减小文件体积。

### 页面操作

```go
func (d *Document) NewPage(pno int, width, height float64) (*Page, error)  // pno=-1 追加到末尾
func (d *Document) DeletePage(pno int) error
func (d *Document) DeletePages(pages ...int) error
func (d *Document) Select(pages []int) error
func (d *Document) InsertPDF(src *Document, opts ...InsertPDFOptions) error
```

### 布局（可重排文档）

```go
func (d *Document) Layout(width, height, fontsize float64)
```

### Xref 访问

```go
func (d *Document) XrefLength() int
func (d *Document) XrefObject(xref int, compressed bool) (string, error)
func (d *Document) PDFCatalog() int
```

### 嵌入文件

```go
func (d *Document) EmbFileCount() int
func (d *Document) EmbFileNames() []string
func (d *Document) EmbFileGet(index int) ([]byte, error)
```

### 便捷方法

```go
func (d *Document) GetPageText(pno int, output string) (string, error)
func (d *Document) GetPagePixmap(pno int, opts ...PixmapOption) (*Pixmap, error)
func (d *Document) SearchPageFor(pno int, needle string, quads bool) ([]Quad, error)
func (d *Document) GetPageFonts(pno int) ([]FontInfo, error)
func (d *Document) GetPageImages(pno int) ([]ImageInfo, error)
```

---

## Page（页面）

表示文档中的单个页面。

### 属性

```go
func (p *Page) Number() int            // 页码（0 起始）
func (p *Page) Rect() Rect            // 页面矩形
func (p *Page) Width() float64        // 宽度
func (p *Page) Height() float64       // 高度
func (p *Page) MediaBox() Rect        // 媒体框
func (p *Page) CropBox() Rect         // 裁剪框
func (p *Page) Rotation() int         // 旋转角度
func (p *Page) SetRotation(rotation int) error
func (p *Page) GetLabel() string      // 页面标签
```

### 生命周期

```go
func (p *Page) Close()
```

### 文本提取

```go
func (p *Page) GetText(output string, flags ...int) (string, error)
```
`output`：`"text"`（纯文本）。`flags`：`TextPreserveLigatures`、`TextPreserveWhitespace` 等的组合。

```go
func (p *Page) GetTextWords(flags ...int) ([]TextWord, error)    // 按单词提取
func (p *Page) GetTextBlocks(flags ...int) ([]TextBlock, error)  // 按文本块提取
func (p *Page) GetTextPage(flags ...int) (*TextPage, error)      // 获取结构化文本页
```

### 搜索

```go
func (p *Page) SearchFor(needle string, quads bool) ([]Quad, error)
```

### 渲染

```go
func (p *Page) GetPixmap(opts ...PixmapOption) (*Pixmap, error)
```

像素图选项（函数式选项模式）：
```go
WithDPI(dpi int)                 // 设置 DPI
WithMatrix(m Matrix)             // 设置变换矩阵
WithColorspace(cs int)           // CsGray, CsRGB, CsCMYK
WithAlpha(alpha bool)            // 启用/禁用透明通道
WithClip(clip Rect)              // 设置裁剪区域
WithAnnots(annots bool)          // 包含/排除注释
```

### 内容插入

```go
func (p *Page) InsertText(pos Point, text string, opts ...TextInsertOption) (int, error)
```
在指定位置插入文本。自动检测 CJK 文字并使用非嵌入 CID 字体。

文本选项：
```go
WithFontName(name string)    // 默认："Helvetica"
WithFontSize(size float64)   // 默认：11
WithColor(color Color)       // 默认：ColorBlack
WithRotate(angle int)        // 旋转角度
```

```go
func (p *Page) InsertImage(rect Rect, imageData []byte, opts ...InsertImageOptions) error
```
在指定矩形区域插入图片。

```go
func (p *Page) InsertHTMLBox(rect Rect, html string, opts ...HTMLBoxOptions) (HTMLBoxResult, error)
```
使用 MuPDF Story API 插入 HTML/CSS 样式内容。支持 CJK、自动缩放、字体子集化。

### 链接与注释

```go
func (p *Page) GetLinks() ([]Link, error)
func (p *Page) GetAnnots() []*Annot
func (p *Page) AddTextAnnot(pos Point, text string) (*Annot, error)
func (p *Page) AddFreetextAnnot(rect Rect, text string, fontsize float64) (*Annot, error)
func (p *Page) AddHighlightAnnot(quads []Quad) (*Annot, error)
func (p *Page) DeleteAnnot(annot *Annot) error
```

### 表单控件

```go
func (p *Page) GetWidgets() []*Widget
```

### 字体与图片信息

```go
func (p *Page) GetFonts() ([]FontInfo, error)
func (p *Page) GetImages() ([]ImageInfo, error)
```

### 变换矩阵

```go
func (p *Page) TransformationMatrix() Matrix   // 页面变换矩阵
func (p *Page) DerotationMatrix() Matrix       // 反旋转矩阵
```

---

## Pixmap（像素图）

表示光栅图像。

### 属性

```go
func (px *Pixmap) Width() int       // 宽度（像素）
func (px *Pixmap) Height() int      // 高度（像素）
func (px *Pixmap) N() int           // 分量数（含 alpha）
func (px *Pixmap) Alpha() int       // 1=有 alpha，0=无
func (px *Pixmap) Stride() int      // 行字节数
func (px *Pixmap) X() int
func (px *Pixmap) Y() int
func (px *Pixmap) IRect() IRect
func (px *Pixmap) Samples() []byte  // 原始像素数据
```

### 生命周期

```go
func (px *Pixmap) Close()
```

### 导出

```go
func (px *Pixmap) ToBytes() ([]byte, error)       // PNG 字节数据
func (px *Pixmap) Save(filename string) error      // 保存为 PNG
func (px *Pixmap) SavePNG(filename string) error   // Save 的别名
func (px *Pixmap) SavePNM(filename string) error   // 保存为 PNM
func (px *Pixmap) ToImage() image.Image            // 转换为 Go image.Image
```

### 像素操作

```go
func (px *Pixmap) SetPixel(x, y int, c []byte)
func (px *Pixmap) GetPixel(x, y int) []byte
func (px *Pixmap) Clear(value int)                     // 清除为指定值
func (px *Pixmap) Invert()                             // 反色
func (px *Pixmap) Gamma(gamma float64)                 // 伽马校正
func (px *Pixmap) Tint(black, white int)               // 着色
func (px *Pixmap) Convert(colorspace int) (*Pixmap, error)  // 色彩空间转换
```

---

## TextPage（结构化文本页）

提供详细的结构化文本分析。

```go
func (t *TextPage) Close()
func (t *TextPage) ExtractText() (string, error)
func (t *TextPage) Blocks() []STextBlock
```

---

## Annot（注释）

表示 PDF 注释。

```go
func (a *Annot) Type() int              // 注释类型编号
func (a *Annot) TypeString() string     // 注释类型名称
func (a *Annot) Rect() Rect            // 注释矩形
func (a *Annot) Contents() string       // 注释内容
func (a *Annot) SetContents(text string)
func (a *Annot) Xref() int
```

---

## Widget（表单控件）

表示 PDF 表单字段。

```go
func (w *Widget) FieldType() int           // 字段类型编号
func (w *Widget) FieldTypeString() string  // 字段类型名称
func (w *Widget) FieldName() string        // 字段名
func (w *Widget) FieldValue() string       // 字段值
func (w *Widget) SetFieldValue(value string) error
func (w *Widget) Rect() Rect
func (w *Widget) Xref() int
```

---

## 几何类型

### Point（点）

```go
type Point struct { X, Y float64 }

func (p Point) Add(other Point) Point       // 加法
func (p Point) Sub(other Point) Point       // 减法
func (p Point) Mul(factor float64) Point    // 缩放
func (p Point) Abs() float64               // 到原点距离
func (p Point) Transform(m Matrix) Point   // 矩阵变换
func (p Point) IsZero() bool               // 是否为原点
```

### Rect（矩形）

```go
type Rect struct { X0, Y0, X1, Y1 float64 }

func (r Rect) Width() float64                      // 宽度
func (r Rect) Height() float64                     // 高度
func (r Rect) IsEmpty() bool                       // 是否为空
func (r Rect) Contains(p Point) bool               // 是否包含点
func (r Rect) ContainsRect(other Rect) bool        // 是否包含矩形
func (r Rect) Intersects(other Rect) bool          // 是否相交
func (r Rect) Intersect(other Rect) Rect           // 交集
func (r Rect) Union(other Rect) Rect               // 并集
func (r Rect) IncludePoint(p Point) Rect           // 扩展以包含点
func (r Rect) Transform(m Matrix) Rect             // 矩阵变换
func (r Rect) Normalize() Rect                     // 规范化
func (r Rect) TopLeft() Point                      // 左上角
func (r Rect) TopRight() Point                     // 右上角
func (r Rect) BottomLeft() Point                   // 左下角
func (r Rect) BottomRight() Point                  // 右下角
func (r Rect) Quad() Quad                          // 转为四边形
func (r Rect) IRect() IRect                        // 转为整数矩形
```

### IRect（整数矩形）

```go
type IRect struct { X0, Y0, X1, Y1 int }

func (r IRect) Width() int
func (r IRect) Height() int
func (r IRect) IsEmpty() bool
func (r IRect) Rect() Rect     // 转为浮点矩形
```

### Quad（四边形）

```go
type Quad struct { UL, UR, LL, LR Point }  // 左上、右上、左下、右下

func (q Quad) Rect() Rect                  // 最小外接矩形
func (q Quad) IsEmpty() bool
func (q Quad) IsRectangular() bool         // 是否为矩形
func (q Quad) IsConvex() bool              // 是否为凸四边形
func (q Quad) Transform(m Matrix) Quad     // 矩阵变换
```

### Matrix（变换矩阵）

```go
type Matrix struct { A, B, C, D, E, F float64 }

var Identity Matrix                                // 单位矩阵

func (m Matrix) Concat(other Matrix) Matrix        // 矩阵乘法
func (m Matrix) PreScale(sx, sy float64) Matrix    // 前置缩放
func (m Matrix) PreTranslate(tx, ty float64) Matrix // 前置平移
func (m Matrix) PreRotate(deg float64) Matrix      // 前置旋转
func (m Matrix) Invert() (Matrix, bool)            // 求逆
func (m Matrix) IsRectilinear() bool               // 是否保持矩形
```

---

## 选项类型

### SaveOptions（保存选项）

```go
type SaveOptions struct {
    Garbage     int    // 0-4，垃圾回收级别
    Deflate     bool   // 压缩流
    Clean       bool   // 清理内容流
    ASCII       bool   // ASCII 十六进制编码二进制数据
    Linear      bool   // 线性化（Web 优化）
    Pretty      bool   // 美化输出
    Incremental bool   // 增量保存
    NoNewID     bool   // 不生成新文件 ID
    Encryption  int    // EncryptNone, EncryptAESV3 等
    Permissions int    // PDF 权限标志
    OwnerPW     string // 所有者密码
    UserPW      string // 用户密码
}

func DefaultSaveOptions() SaveOptions
func EzSaveOptions() SaveOptions    // Garbage=3, Deflate=true
```

### HTMLBoxOptions

```go
type HTMLBoxOptions struct {
    CSS      string   // 附加 CSS 样式
    ScaleLow float64  // 最小缩放比例（0=自动缩小适配，1=不缩放）
    Overlay  bool     // true=前景，false=背景
}
```

### HTMLBoxResult

```go
type HTMLBoxResult struct {
    SpareHeight float64  // 内容后剩余高度
    Scale       float64  // 实际使用的缩放比例
}
```

### InsertPDFOptions（PDF 插入选项）

```go
type InsertPDFOptions struct {
    FromPage int   // 源起始页（默认：0）
    ToPage   int   // 源结束页（默认：最后一页）
    StartAt  int   // 目标插入位置（默认：末尾）
    Rotate   int
    Links    bool
    Annots   bool
}
```

### InsertImageOptions（图片插入选项）

```go
type InsertImageOptions struct {
    KeepProportion bool  // 保持宽高比
    Overlay        bool  // 前景或背景
}
```

---

## 常量

### 色彩空间

```go
CsGray = 0   // 灰度
CsRGB  = 1   // RGB
CsCMYK = 2   // CMYK
```

### 注释类型

`AnnotText`、`AnnotLink`、`AnnotFreeText`、`AnnotLine`、`AnnotSquare`、`AnnotCircle`、`AnnotHighlight`、`AnnotUnderline`、`AnnotStrikeOut`、`AnnotRedact`、`AnnotStamp`、`AnnotInk` 等。

### 表单控件类型

`WidgetTypeButton`、`WidgetTypeCheckbox`、`WidgetTypeCombobox`、`WidgetTypeListbox`、`WidgetTypeRadioButton`、`WidgetTypeSignature`、`WidgetTypeText`。

### 加密方式

`EncryptNone`、`EncryptRC4V1`（40 位）、`EncryptRC4V2`（128 位）、`EncryptAESV2`（128 位）、`EncryptAESV3`（256 位）。

### 标准纸张尺寸

`PaperA4`、`PaperA3`、`PaperA5`、`PaperLetter`、`PaperLegal` — 预定义 `Rect` 值。

### 预定义颜色

`ColorBlack`、`ColorWhite`、`ColorRed`、`ColorGreen`、`ColorBlue`、`ColorYellow`、`ColorMagenta`、`ColorCyan`。

---

## 错误

| 错误 | 说明 |
|------|------|
| `ErrInitFailed` | MuPDF 上下文初始化失败 |
| `ErrOpenFailed` | 无法打开文档 |
| `ErrPageNotFound` | 页码超出范围 |
| `ErrNotPDF` | 对非 PDF 文档执行了 PDF 专用操作 |
| `ErrEncrypted` | 文档已加密且未认证 |
| `ErrAuthFailed` | 密码认证失败 |
| `ErrClosed` | 对已关闭的文档执行操作 |
| `ErrTextExtract` | 文本提取失败 |
| `ErrPixmap` | 像素图操作失败 |
| `ErrSave` | 保存失败 |
| `ErrInvalidArg` | 无效参数 |
| `ErrOutline` | 目录/大纲操作失败 |
| `ErrSearch` | 搜索失败 |
| `ErrConvert` | 文档转换失败 |
| `ErrEmbeddedFile` | 嵌入文件操作失败 |
| `ErrXref` | Xref 操作失败 |
| `ErrOverflow` | 内容超出目标矩形范围 |

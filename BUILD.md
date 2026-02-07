# GoMuPDF 编译文档

本文档详细说明如何在 Windows、Linux、macOS 三个平台上编译和构建 GoMuPDF。

## 目录

- [项目结构](#项目结构)
- [环境要求](#环境要求)
- [编译 MuPDF 静态库](#编译-mupdf-静态库)
- [构建 GoMuPDF](#构建-gomupdf)
- [测试](#测试)
- [构建标签说明](#构建标签说明)
- [CGO 链接参数详解](#cgo-链接参数详解)
- [CI/CD 持续集成](#cicd-持续集成)
- [常见问题排查](#常见问题排查)

---

## 项目结构

```
GoMuPDF/
├── include/mupdf/          # MuPDF 1.24.9 C 头文件
│   ├── fitz/                # fitz 核心头文件 (52个)
│   ├── pdf/                 # PDF 专用头文件
│   └── helpers/             # 辅助头文件
├── libs/                    # 平台特定的 MuPDF 静态库
│   ├── libmupdf_<os>_<arch>.a
│   └── libmupdfthird_<os>_<arch>.a
├── gomupdf.h                # C 包装函数头文件 (CGO 桥接层)
├── mupdf_cgo.go             # CGO 绑定入口 (含平台 LDFLAGS)
├── stubs_nocgo.go           # 无 CGO 时的桩实现
├── document.go              # Document 类型实现
├── document_pdf.go          # PDF 专用 Document 方法
├── page.go                  # Page 类型实现
├── page_insert.go           # 页面插入操作
├── pixmap.go                # Pixmap 像素图实现
├── textpage.go              # 文本提取实现
├── annot.go                 # 注释实现
├── widget.go                # 表单控件实现
├── shape.go                 # 绘图形状
├── embfile.go               # 嵌入文件
├── outline.go               # 大纲/书签
├── geometry.go              # 几何类型 (Rect, Point)
├── matrix.go                # 变换矩阵
├── quad.go                  # 四边形
├── colorspace.go            # 色彩空间
├── types.go                 # 公共类型定义
├── constants.go             # 常量定义
├── errors.go                # 错误类型
├── tools.go                 # 工具函数
├── gomupdf_test.go          # CGO 集成测试 (107个)
├── geometry_test.go         # 纯 Go 测试 (75个)
├── example_test.go          # 示例测试
├── build_libs.sh            # Linux/macOS 编译脚本
├── build_libs.bat           # Windows 编译脚本
├── Makefile                 # 构建目标
├── .github/workflows/ci.yml # GitHub Actions CI
├── go.mod                   # Go 模块定义
└── README.md                # 项目说明
```

### 静态库命名规则

所有 MuPDF 静态库存放在 `libs/` 目录下，命名格式为：

```
libmupdf_<os>_<arch>.a        # MuPDF 核心库
libmupdfthird_<os>_<arch>.a   # MuPDF 第三方依赖库
```

支持的组合：

| 文件名 | 平台 |
|--------|------|
| `libmupdf_windows_amd64.a` | Windows x86_64 |
| `libmupdf_windows_arm64.a` | Windows ARM64 |
| `libmupdf_linux_amd64.a` | Linux x86_64 |
| `libmupdf_linux_arm64.a` | Linux ARM64 |
| `libmupdf_darwin_amd64.a` | macOS Intel |
| `libmupdf_darwin_arm64.a` | macOS Apple Silicon |

---

## 环境要求

### 通用要求

| 工具 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.21+ (推荐 1.23) | 需要 CGO 支持 |
| Git | 任意 | 用于克隆 MuPDF 源码 |
| Make | GNU Make | 编译 MuPDF |

### Windows

| 工具 | 说明 |
|------|------|
| MinGW-W64 | 提供 GCC、ar、make 等工具 |
| CGO_ENABLED=1 | Go 环境变量，需确保开启 |

推荐安装方式：
```cmd
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

安装后确保 MinGW 的 `bin` 目录在 `PATH` 中，验证：
```cmd
gcc --version
make --version
```

> **注意**：Windows 上 `CGO_ENABLED` 默认为 1（当检测到 GCC 时），但如果 Go 找不到 GCC，需要手动设置：
> ```cmd
> set CGO_ENABLED=1
> set CC=gcc
> ```

### Linux

```bash
# Debian/Ubuntu
sudo apt-get update
sudo apt-get install -y build-essential git

# Fedora/RHEL
sudo dnf groupinstall "Development Tools"
sudo dnf install git

# Arch Linux
sudo pacman -S base-devel git
```

需要的系统库（运行时链接）：
- `libstdc++` — C++ 标准库
- `libpthread` — POSIX 线程
- `libdl` — 动态加载
- `libm` — 数学库

这些在大多数 Linux 发行版上已预装。

### macOS

```bash
# 安装 Xcode 命令行工具 (提供 Clang 和 make)
xcode-select --install

# 或通过 Homebrew 安装 GCC
brew install gcc git
```

需要的系统框架：
- `CoreFoundation`
- `Security`
- `libc++` — C++ 标准库

---

## 编译 MuPDF 静态库

GoMuPDF 需要 MuPDF **1.24.9** 版本的静态库。头文件已包含在 `include/` 目录中，**版本必须匹配**。

### 方法一：使用自动化脚本（推荐）

#### Linux / macOS

```bash
chmod +x build_libs.sh

# 自动克隆 MuPDF 1.24.9 并编译
./build_libs.sh

# 或指定已有的 MuPDF 源码目录
./build_libs.sh /path/to/mupdf-src
```

脚本会自动：
1. 检测当前操作系统和 CPU 架构
2. 如未指定源码目录，克隆 MuPDF 1.24.9 到 `/tmp/mupdf-1.24.9`
3. 使用 `make` 并行编译（自动检测 CPU 核心数）
4. 将编译产物复制到 `libs/` 并按命名规则重命名

#### Windows

```cmd
REM 自动克隆并编译
build_libs.bat

REM 或指定已有源码目录
build_libs.bat D:\mupdf-src
```

> **Windows 注意事项**：需要在 MinGW/MSYS2 环境中运行 `make`，或确保 GNU Make 在 PATH 中。

#### 使用 Makefile

```bash
# 自动编译（等同于运行 build_libs.sh）
make libs

# 指定 MuPDF 源码路径
make libs MUPDF_SRC=/path/to/mupdf-src
```

### 方法二：手动编译

#### 1. 获取 MuPDF 源码

```bash
git clone --depth 1 --branch 1.24.9 \
    --recurse-submodules \
    https://github.com/ArtifexSoftware/mupdf.git mupdf-src
```

> `--recurse-submodules` 是必须的，MuPDF 依赖多个子模块（freetype、harfbuzz、libjpeg 等）。

#### 2. 编译

```bash
cd mupdf-src
make -j$(nproc) \
    HAVE_X11=no \
    HAVE_GLUT=no \
    HAVE_CURL=no \
    USE_SYSTEM_LIBS=no \
    XCFLAGS="-fPIC" \
    libs
```

编译参数说明：

| 参数 | 说明 |
|------|------|
| `HAVE_X11=no` | 不编译 X11 GUI 支持 |
| `HAVE_GLUT=no` | 不编译 OpenGL 查看器 |
| `HAVE_CURL=no` | 不编译 HTTP 支持 |
| `USE_SYSTEM_LIBS=no` | 使用 MuPDF 自带的第三方库（避免版本冲突） |
| `XCFLAGS="-fPIC"` | 生成位置无关代码（Linux 必需） |
| `libs` | 只编译库文件，不编译可执行程序 |

#### 3. 复制库文件

编译完成后，产物位于 `mupdf-src/build/release/`：

```bash
mkdir -p libs

# Linux amd64 示例
cp mupdf-src/build/release/libmupdf.a     libs/libmupdf_linux_amd64.a
cp mupdf-src/build/release/libmupdf-third.a libs/libmupdfthird_linux_amd64.a

# macOS arm64 示例
cp mupdf-src/build/release/libmupdf.a     libs/libmupdf_darwin_arm64.a
cp mupdf-src/build/release/libmupdf-third.a libs/libmupdfthird_darwin_arm64.a

# Windows amd64 示例
copy mupdf-src\build\release\libmupdf.a     libs\libmupdf_windows_amd64.a
copy mupdf-src\build\release\libmupdf-third.a libs\libmupdfthird_windows_amd64.a
```

### 方法三：从 go-fitz 提取

如果你已安装 [go-fitz](https://github.com/gen2brain/go-fitz) v1.24.15，可以直接提取其静态库：

```bash
# 查找 go-fitz 模块缓存路径
go env GOMODCACHE
# 通常在: ~/go/pkg/mod/github.com/gen2brain/go-fitz@v1.24.15/

# 复制并重命名库文件到 libs/
```

> go-fitz v1.24.15 使用的 MuPDF 版本与 1.24.9 头文件兼容。

---

## 构建 GoMuPDF

### 标准构建（需要 MuPDF 静态库）

```bash
# 确保 libs/ 目录下有当前平台的静态库
go build .
```

### 验证构建

```bash
go vet .
```

### 使用 Makefile

```bash
make build    # go build .
make vet      # go vet .
```

### 无 CGO 构建（桩模式）

如果没有 MuPDF 静态库，可以使用 `nomupdf` 标签编译：

```bash
go build -tags nomupdf .
```

此模式下所有 MuPDF 相关函数返回 `ErrInitFailed`，但纯 Go 类型（Rect、Matrix、Point 等）正常工作。

---

## 测试

### 完整测试（需要 MuPDF 静态库）

```bash
go test -v -count=1 .
```

运行 107 个 CGO 集成测试，覆盖：
- 文档打开/创建/保存
- 页面操作（加载、渲染、文本提取）
- PDF 专用功能（注释、书签、嵌入文件、表单控件）
- 像素图操作
- 内存管理

### 纯 Go 测试（无需 MuPDF 库）

```bash
go test -v -count=1 -tags nomupdf .
```

运行 75 个纯 Go 测试，覆盖：
- 几何类型（Rect、Point、IRect）
- 变换矩阵（Matrix）
- 四边形（Quad）
- 工具函数
- 大纲/书签
- 类型定义
- 色彩空间常量
- 错误类型
- 桩函数行为

### 使用 Makefile

```bash
make test       # 完整测试
make test-pure  # 纯 Go 测试
```

---

## 构建标签说明

GoMuPDF 使用 Go 构建标签控制编译模式：

### 默认模式（CGO 启用）

```
//go:build cgo && !nomupdf
```

- 文件：`mupdf_cgo.go`、`document.go`、`page.go` 等所有功能文件
- 条件：CGO 可用且未设置 `nomupdf` 标签
- 行为：通过 CGO 调用 MuPDF C 库，提供完整功能

### 桩模式（无 CGO）

```
//go:build !cgo || nomupdf
```

- 文件：`stubs_nocgo.go`
- 条件：CGO 不可用，或显式设置 `nomupdf` 标签
- 行为：所有 MuPDF 函数返回 `ErrInitFailed`
- 用途：CI 环境、交叉编译、无 MuPDF 库的平台

### 测试文件标签

| 文件 | 构建标签 | 说明 |
|------|---------|------|
| `gomupdf_test.go` | `cgo && !nomupdf` | CGO 集成测试 |
| `geometry_test.go` | `!cgo \|\| nomupdf` | 纯 Go 测试 |
| `example_test.go` | `cgo && !nomupdf` | 示例测试 |

---

## CGO 链接参数详解

`mupdf_cgo.go` 中定义了各平台的 CGO 链接参数：

### 编译参数（所有平台通用）

```
#cgo CFLAGS: -Iinclude
```

将 `include/` 目录加入 C 头文件搜索路径。

### Windows (amd64 / arm64)

```
-L${SRCDIR}/libs
-lmupdf_windows_<arch>
-lmupdfthird_windows_<arch>
-lm -lgdi32 -lcomdlg32 -luser32 -ladvapi32 -lshell32
```

| 库 | 说明 |
|----|------|
| `-lm` | 数学库 |
| `-lgdi32` | Windows GDI 图形接口 |
| `-lcomdlg32` | 通用对话框 |
| `-luser32` | 用户界面 |
| `-ladvapi32` | 高级 Windows API（注册表、安全等） |
| `-lshell32` | Shell 功能 |

### Linux (amd64 / arm64)

```
-L${SRCDIR}/libs
-lmupdf_linux_<arch>
-lmupdfthird_linux_<arch>
-lm -lstdc++ -lpthread -ldl
```

| 库 | 说明 |
|----|------|
| `-lm` | 数学库 |
| `-lstdc++` | C++ 标准库（MuPDF 第三方依赖需要） |
| `-lpthread` | POSIX 线程 |
| `-ldl` | 动态链接加载器 |

### macOS (amd64 / arm64)

```
-L${SRCDIR}/libs
-lmupdf_darwin_<arch>
-lmupdfthird_darwin_<arch>
-lm -lc++ -framework CoreFoundation -framework Security
```

| 库 | 说明 |
|----|------|
| `-lm` | 数学库 |
| `-lc++` | C++ 标准库（macOS 使用 libc++ 而非 libstdc++） |
| `-framework CoreFoundation` | macOS 核心基础框架 |
| `-framework Security` | macOS 安全框架（证书、加密等） |

### `${SRCDIR}` 说明

`${SRCDIR}` 是 CGO 特殊变量，自动替换为包含 Go 源文件的目录路径。这确保无论从哪里运行 `go build`，都能正确找到 `libs/` 目录。

---

## CI/CD 持续集成

项目使用 GitHub Actions，配置文件位于 `.github/workflows/ci.yml`。

### 工作流概览

CI 包含两个 Job：

#### 1. `test-pure` — 纯 Go 测试

- 平台：Ubuntu、macOS、Windows
- 无需编译 MuPDF
- 运行 `go test -tags nomupdf` 和 `go vet -tags nomupdf`
- 验证桩实现和纯 Go 代码在所有平台上正常工作

#### 2. `test-cgo` — 完整 CGO 测试

- 平台：Linux amd64、macOS arm64、Windows amd64
- 从源码编译 MuPDF 1.24.9
- 运行 `go build`、`go vet`、`go test`
- 验证完整功能在各平台上正常工作

### 触发条件

- 推送到 `main` 分支
- 向 `main` 分支发起 Pull Request

---

## 常见问题排查

### 1. `undefined reference to ...` 链接错误

**原因**：`libs/` 目录下缺少当前平台的静态库。

**解决**：
```bash
# 检查 libs/ 目录
ls libs/

# 编译当前平台的库
./build_libs.sh
```

### 2. `cannot find -lmupdf_<os>_<arch>` 错误

**原因**：库文件命名不正确。

**解决**：确保文件名严格遵循 `libmupdf_<os>_<arch>.a` 格式。例如 Linux amd64 必须是 `libmupdf_linux_amd64.a`，不能是 `libmupdf.a`。

### 3. `fatal error: mupdf/fitz.h: No such file or directory`

**原因**：`include/` 目录缺失或头文件不完整。

**解决**：确保项目根目录下有完整的 `include/mupdf/` 头文件目录。头文件版本必须是 1.24.9。

### 4. Windows 上 `exec: "gcc": executable file not found in %PATH%`

**原因**：MinGW-W64 未安装或未加入 PATH。

**解决**：
```cmd
REM 安装 MinGW-W64
winget install BrechtSanders.WinLibs.POSIX.UCRT

REM 将 MinGW bin 目录加入 PATH
set PATH=%PATH%;C:\mingw64\bin

REM 验证
gcc --version
```

### 5. macOS 上 `ld: framework not found CoreFoundation`

**原因**：Xcode 命令行工具未安装。

**解决**：
```bash
xcode-select --install
```

### 6. Linux 上 `-lstdc++` 找不到

**原因**：缺少 C++ 标准库开发包。

**解决**：
```bash
# Debian/Ubuntu
sudo apt-get install -y libstdc++-dev g++

# Fedora/RHEL
sudo dnf install gcc-c++ libstdc++-devel
```

### 7. 头文件版本不匹配

**症状**：编译通过但运行时崩溃，或出现结构体大小不匹配的错误。

**原因**：`include/` 中的头文件版本与 `libs/` 中的静态库版本不一致。

**解决**：确保两者都是 MuPDF 1.24.9 版本。可以从同一份 MuPDF 源码中同时获取头文件和编译静态库。

### 8. 交叉编译

GoMuPDF 不支持直接交叉编译，因为需要目标平台的 MuPDF 静态库。每个平台需要在对应的操作系统上编译静态库，然后将 `.a` 文件复制到 `libs/` 目录。

CI 流程展示了如何在各平台上自动完成此过程。

### 9. 想跳过 CGO 编译

如果只需要使用纯 Go 类型（Rect、Matrix、Point 等），或在没有 C 编译器的环境中工作：

```bash
go build -tags nomupdf .
go test -tags nomupdf .
```

此模式下所有需要 MuPDF 的函数会返回 `ErrInitFailed`。

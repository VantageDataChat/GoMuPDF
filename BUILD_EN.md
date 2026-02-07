# GoMuPDF Build Documentation

This document provides detailed instructions for building GoMuPDF on Windows, Linux, and macOS.

## Table of Contents

- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Building MuPDF Static Libraries](#building-mupdf-static-libraries)
- [Building GoMuPDF](#building-gomupdf)
- [Testing](#testing)
- [Build Tags](#build-tags)
- [CGO Linker Flags Explained](#cgo-linker-flags-explained)
- [CI/CD](#cicd)
- [Troubleshooting](#troubleshooting)

---

## Project Structure

```
GoMuPDF/
├── include/mupdf/          # MuPDF 1.24.9 C headers
│   ├── fitz/                # Core fitz headers (52 files)
│   ├── pdf/                 # PDF-specific headers
│   └── helpers/             # Helper headers
├── libs/                    # Platform-specific MuPDF static libraries
│   ├── libmupdf_<os>_<arch>.a
│   └── libmupdfthird_<os>_<arch>.a
├── gomupdf.h                # C wrapper functions (CGO bridge layer)
├── mupdf_cgo.go             # CGO entry point (platform-specific LDFLAGS)
├── stubs_nocgo.go           # Stub implementation for non-CGO builds
├── document.go              # Document type implementation
├── document_pdf.go          # PDF-specific Document methods
├── page.go                  # Page type implementation
├── page_insert.go           # Page insertion operations
├── pixmap.go                # Pixmap (raster image) implementation
├── textpage.go              # Text extraction implementation
├── annot.go                 # Annotation implementation
├── widget.go                # Form widget implementation
├── shape.go                 # Drawing shapes
├── embfile.go               # Embedded files
├── outline.go               # Outlines / bookmarks
├── geometry.go              # Geometry types (Rect, Point)
├── matrix.go                # Transformation matrix
├── quad.go                  # Quadrilateral
├── colorspace.go            # Color spaces
├── types.go                 # Common type definitions
├── constants.go             # Constants
├── errors.go                # Error types
├── tools.go                 # Utility functions
├── gomupdf_test.go          # CGO integration tests (107)
├── geometry_test.go         # Pure Go tests (75)
├── example_test.go          # Example tests
├── build_libs.sh            # Linux/macOS build script
├── build_libs.bat           # Windows build script
├── Makefile                 # Build targets
├── .github/workflows/ci.yml # GitHub Actions CI
├── go.mod                   # Go module definition
└── README.md                # Project overview
```

### Static Library Naming Convention

All MuPDF static libraries are stored in the `libs/` directory with the following naming format:

```
libmupdf_<os>_<arch>.a        # MuPDF core library
libmupdfthird_<os>_<arch>.a   # MuPDF third-party dependencies
```

Supported combinations:

| Filename | Platform |
|----------|----------|
| `libmupdf_windows_amd64.a` | Windows x86_64 |
| `libmupdf_windows_arm64.a` | Windows ARM64 |
| `libmupdf_linux_amd64.a` | Linux x86_64 |
| `libmupdf_linux_arm64.a` | Linux ARM64 |
| `libmupdf_darwin_amd64.a` | macOS Intel |
| `libmupdf_darwin_arm64.a` | macOS Apple Silicon |

---

## Prerequisites

### Common Requirements

| Tool | Version | Notes |
|------|---------|-------|
| Go | 1.21+ (1.23 recommended) | CGO support required |
| Git | Any | For cloning MuPDF source |
| Make | GNU Make | For building MuPDF |

### Windows

| Tool | Notes |
|------|-------|
| MinGW-W64 | Provides GCC, ar, make |
| CGO_ENABLED=1 | Go environment variable, must be enabled |

Recommended installation:
```cmd
winget install BrechtSanders.WinLibs.POSIX.UCRT
```

After installation, ensure the MinGW `bin` directory is in your `PATH`:
```cmd
gcc --version
make --version
```

> **Note**: On Windows, `CGO_ENABLED` defaults to 1 when GCC is detected. If Go cannot find GCC, set it manually:
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

Required system libraries (linked at runtime):
- `libstdc++` — C++ standard library
- `libpthread` — POSIX threads
- `libdl` — Dynamic loading
- `libm` — Math library

These are pre-installed on most Linux distributions.

### macOS

```bash
# Install Xcode Command Line Tools (provides Clang and make)
xcode-select --install

# Or install GCC via Homebrew
brew install gcc git
```

Required system frameworks:
- `CoreFoundation`
- `Security`
- `libc++` — C++ standard library

---

## Building MuPDF Static Libraries

GoMuPDF requires MuPDF **1.24.9** static libraries. The headers are already included in the `include/` directory — **versions must match**.

### Option 1: Automated Scripts (Recommended)

#### Linux / macOS

```bash
chmod +x build_libs.sh

# Auto-clone MuPDF 1.24.9 and build
./build_libs.sh

# Or specify an existing MuPDF source directory
./build_libs.sh /path/to/mupdf-src
```

The script automatically:
1. Detects the current OS and CPU architecture
2. Clones MuPDF 1.24.9 to `/tmp/mupdf-1.24.9` if no source directory is specified
3. Builds in parallel using `make` (auto-detects CPU core count)
4. Copies build artifacts to `libs/` with the correct naming convention

#### Windows

```cmd
REM Auto-clone and build
build_libs.bat

REM Or specify an existing source directory
build_libs.bat D:\mupdf-src
```

> **Windows note**: Requires `make` (GNU Make) available in PATH, typically from MinGW/MSYS2.

#### Using the Makefile

```bash
# Auto-build (equivalent to running build_libs.sh)
make libs

# Specify MuPDF source path
make libs MUPDF_SRC=/path/to/mupdf-src
```

### Option 2: Manual Build

#### 1. Get MuPDF Source

```bash
git clone --depth 1 --branch 1.24.9 \
    --recurse-submodules \
    https://github.com/ArtifexSoftware/mupdf.git mupdf-src
```

> `--recurse-submodules` is required — MuPDF depends on several submodules (freetype, harfbuzz, libjpeg, etc.).

#### 2. Build

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

Build parameter reference:

| Parameter | Description |
|-----------|-------------|
| `HAVE_X11=no` | Disable X11 GUI support |
| `HAVE_GLUT=no` | Disable OpenGL viewer |
| `HAVE_CURL=no` | Disable HTTP support |
| `USE_SYSTEM_LIBS=no` | Use MuPDF's bundled third-party libraries (avoids version conflicts) |
| `XCFLAGS="-fPIC"` | Generate position-independent code (required on Linux) |
| `libs` | Build libraries only, skip executables |

#### 3. Copy Library Files

Build artifacts are located in `mupdf-src/build/release/`:

```bash
mkdir -p libs

# Linux amd64 example
cp mupdf-src/build/release/libmupdf.a     libs/libmupdf_linux_amd64.a
cp mupdf-src/build/release/libmupdf-third.a libs/libmupdfthird_linux_amd64.a

# macOS arm64 example
cp mupdf-src/build/release/libmupdf.a     libs/libmupdf_darwin_arm64.a
cp mupdf-src/build/release/libmupdf-third.a libs/libmupdfthird_darwin_arm64.a

# Windows amd64 example
copy mupdf-src\build\release\libmupdf.a     libs\libmupdf_windows_amd64.a
copy mupdf-src\build\release\libmupdf-third.a libs\libmupdfthird_windows_amd64.a
```

### Option 3: Extract from go-fitz

If you already have [go-fitz](https://github.com/gen2brain/go-fitz) v1.24.15 installed, you can extract its static libraries:

```bash
# Find the go-fitz module cache path
go env GOMODCACHE
# Typically at: ~/go/pkg/mod/github.com/gen2brain/go-fitz@v1.24.15/

# Copy and rename library files to libs/
```

> go-fitz v1.24.15 uses a MuPDF version compatible with the 1.24.9 headers.

---

## Building GoMuPDF

### Standard Build (requires MuPDF static libraries)

```bash
# Ensure libs/ contains static libraries for the current platform
go build .
```

### Verify Build

```bash
go vet .
```

### Using the Makefile

```bash
make build    # go build .
make vet      # go vet .
```

### Non-CGO Build (Stub Mode)

If MuPDF static libraries are not available, build with the `nomupdf` tag:

```bash
go build -tags nomupdf .
```

In this mode, all MuPDF-related functions return `ErrInitFailed`, but pure Go types (Rect, Matrix, Point, etc.) work normally.

---

## Testing

### Full Tests (requires MuPDF static libraries)

```bash
go test -v -count=1 .
```

Runs 107 CGO integration tests covering:
- Document open/create/save
- Page operations (load, render, text extraction)
- PDF-specific features (annotations, bookmarks, embedded files, form widgets)
- Pixmap operations
- Memory management

### Pure Go Tests (no MuPDF libraries needed)

```bash
go test -v -count=1 -tags nomupdf .
```

Runs 75 pure Go tests covering:
- Geometry types (Rect, Point, IRect)
- Transformation matrices (Matrix)
- Quadrilaterals (Quad)
- Utility functions
- Outlines / bookmarks
- Type definitions
- Color space constants
- Error types
- Stub function behavior

### Using the Makefile

```bash
make test       # Full tests
make test-pure  # Pure Go tests only
```

---

## Build Tags

GoMuPDF uses Go build tags to control compilation mode:

### Default Mode (CGO Enabled)

```
//go:build cgo && !nomupdf
```

- Files: `mupdf_cgo.go`, `document.go`, `page.go`, and all feature files
- Condition: CGO is available and `nomupdf` tag is not set
- Behavior: Calls MuPDF C library via CGO, providing full functionality

### Stub Mode (No CGO)

```
//go:build !cgo || nomupdf
```

- File: `stubs_nocgo.go`
- Condition: CGO is unavailable, or `nomupdf` tag is explicitly set
- Behavior: All MuPDF functions return `ErrInitFailed`
- Use cases: CI environments, cross-compilation, platforms without MuPDF libraries

### Test File Tags

| File | Build Tag | Description |
|------|-----------|-------------|
| `gomupdf_test.go` | `cgo && !nomupdf` | CGO integration tests |
| `geometry_test.go` | `!cgo \|\| nomupdf` | Pure Go tests |
| `example_test.go` | `cgo && !nomupdf` | Example tests |

---

## CGO Linker Flags Explained

Platform-specific CGO linker flags are defined in `mupdf_cgo.go`:

### Compiler Flags (all platforms)

```
#cgo CFLAGS: -Iinclude
```

Adds the `include/` directory to the C header search path.

### Windows (amd64 / arm64)

```
-L${SRCDIR}/libs
-lmupdf_windows_<arch>
-lmupdfthird_windows_<arch>
-lm -lgdi32 -lcomdlg32 -luser32 -ladvapi32 -lshell32
```

| Library | Description |
|---------|-------------|
| `-lm` | Math library |
| `-lgdi32` | Windows GDI graphics interface |
| `-lcomdlg32` | Common dialog boxes |
| `-luser32` | User interface |
| `-ladvapi32` | Advanced Windows API (registry, security, etc.) |
| `-lshell32` | Shell functions |

### Linux (amd64 / arm64)

```
-L${SRCDIR}/libs
-lmupdf_linux_<arch>
-lmupdfthird_linux_<arch>
-lm -lstdc++ -lpthread -ldl
```

| Library | Description |
|---------|-------------|
| `-lm` | Math library |
| `-lstdc++` | C++ standard library (required by MuPDF third-party deps) |
| `-lpthread` | POSIX threads |
| `-ldl` | Dynamic linker/loader |

### macOS (amd64 / arm64)

```
-L${SRCDIR}/libs
-lmupdf_darwin_<arch>
-lmupdfthird_darwin_<arch>
-lm -lc++ -framework CoreFoundation -framework Security
```

| Library | Description |
|---------|-------------|
| `-lm` | Math library |
| `-lc++` | C++ standard library (macOS uses libc++ instead of libstdc++) |
| `-framework CoreFoundation` | macOS Core Foundation framework |
| `-framework Security` | macOS Security framework (certificates, encryption, etc.) |

### About `${SRCDIR}`

`${SRCDIR}` is a special CGO variable that is automatically replaced with the directory path containing the Go source files. This ensures that `go build` can always locate the `libs/` directory regardless of the working directory.

---

## CI/CD

The project uses GitHub Actions. The configuration is at `.github/workflows/ci.yml`.

### Workflow Overview

The CI pipeline consists of two jobs:

#### 1. `test-pure` — Pure Go Tests

- Platforms: Ubuntu, macOS, Windows
- No MuPDF compilation required
- Runs `go test -tags nomupdf` and `go vet -tags nomupdf`
- Verifies that stub implementations and pure Go code work correctly on all platforms

#### 2. `test-cgo` — Full CGO Tests

- Platforms: Linux amd64, macOS arm64, Windows amd64
- Builds MuPDF 1.24.9 from source
- Runs `go build`, `go vet`, `go test`
- Verifies full functionality on each platform

### Triggers

- Push to `main` branch
- Pull request targeting `main` branch

---

## Troubleshooting

### 1. `undefined reference to ...` linker errors

**Cause**: Missing static libraries for the current platform in `libs/`.

**Solution**:
```bash
# Check libs/ directory
ls libs/

# Build libraries for the current platform
./build_libs.sh
```

### 2. `cannot find -lmupdf_<os>_<arch>` error

**Cause**: Library files are incorrectly named.

**Solution**: Ensure filenames strictly follow the `libmupdf_<os>_<arch>.a` format. For example, Linux amd64 must be `libmupdf_linux_amd64.a`, not `libmupdf.a`.

### 3. `fatal error: mupdf/fitz.h: No such file or directory`

**Cause**: The `include/` directory is missing or incomplete.

**Solution**: Ensure the project root contains a complete `include/mupdf/` header directory. Headers must be version 1.24.9.

### 4. Windows: `exec: "gcc": executable file not found in %PATH%`

**Cause**: MinGW-W64 is not installed or not in PATH.

**Solution**:
```cmd
REM Install MinGW-W64
winget install BrechtSanders.WinLibs.POSIX.UCRT

REM Add MinGW bin directory to PATH
set PATH=%PATH%;C:\mingw64\bin

REM Verify
gcc --version
```

### 5. macOS: `ld: framework not found CoreFoundation`

**Cause**: Xcode Command Line Tools are not installed.

**Solution**:
```bash
xcode-select --install
```

### 6. Linux: `-lstdc++` not found

**Cause**: Missing C++ standard library development package.

**Solution**:
```bash
# Debian/Ubuntu
sudo apt-get install -y libstdc++-dev g++

# Fedora/RHEL
sudo dnf install gcc-c++ libstdc++-devel
```

### 7. Header Version Mismatch

**Symptoms**: Build succeeds but crashes at runtime, or struct size mismatch errors.

**Cause**: Headers in `include/` and static libraries in `libs/` are from different MuPDF versions.

**Solution**: Ensure both are from MuPDF 1.24.9. Ideally, obtain headers and compile static libraries from the same MuPDF source tree.

### 8. Cross-Compilation

GoMuPDF does not support direct cross-compilation because it requires MuPDF static libraries for the target platform. Each platform's static libraries must be compiled on the corresponding OS, then the `.a` files copied to the `libs/` directory.

The CI workflow demonstrates how to automate this process on each platform.

### 9. Skipping CGO Compilation

If you only need pure Go types (Rect, Matrix, Point, etc.) or are working in an environment without a C compiler:

```bash
go build -tags nomupdf .
go test -tags nomupdf .
```

In this mode, all functions requiring MuPDF will return `ErrInitFailed`.

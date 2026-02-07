#!/bin/bash
# build_libs.sh - Build MuPDF static libraries for the current platform.
#
# Usage:
#   ./build_libs.sh [mupdf_source_dir]
#
# If mupdf_source_dir is not provided, it will clone MuPDF 1.24.9 into /tmp/mupdf-src.
#
# Prerequisites:
#   - GCC or Clang toolchain
#   - make
#   - git (if cloning)
#
# Output:
#   libs/libmupdf_<os>_<arch>.a
#   libs/libmupdfthird_<os>_<arch>.a

set -e

MUPDF_VERSION="1.24.9"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
LIBS_DIR="${SCRIPT_DIR}/libs"

# Detect OS
case "$(uname -s)" in
    Linux*)   OS=linux;;
    Darwin*)  OS=darwin;;
    MINGW*|MSYS*|CYGWIN*) OS=windows;;
    *)        echo "Unsupported OS: $(uname -s)"; exit 1;;
esac

# Detect architecture
case "$(uname -m)" in
    x86_64|amd64)  ARCH=amd64;;
    aarch64|arm64) ARCH=arm64;;
    *)             echo "Unsupported arch: $(uname -m)"; exit 1;;
esac

echo "Building MuPDF ${MUPDF_VERSION} for ${OS}/${ARCH}"

# Source directory
MUPDF_SRC="${1:-}"
if [ -z "${MUPDF_SRC}" ]; then
    MUPDF_SRC="/tmp/mupdf-${MUPDF_VERSION}"
    if [ ! -d "${MUPDF_SRC}" ]; then
        echo "Cloning MuPDF ${MUPDF_VERSION}..."
        git clone --depth 1 --branch "${MUPDF_VERSION}" \
            --recurse-submodules \
            https://github.com/ArtifexSoftware/mupdf.git "${MUPDF_SRC}"
    fi
fi

if [ ! -f "${MUPDF_SRC}/Makefile" ]; then
    echo "Error: MuPDF source not found at ${MUPDF_SRC}"
    exit 1
fi

echo "Building from ${MUPDF_SRC}..."

# Build MuPDF
NPROC=4
if command -v nproc &>/dev/null; then
    NPROC=$(nproc)
elif command -v sysctl &>/dev/null; then
    NPROC=$(sysctl -n hw.ncpu 2>/dev/null || echo 4)
fi

(
    cd "${MUPDF_SRC}"
    make -j"${NPROC}" \
        HAVE_X11=no \
        HAVE_GLUT=no \
        HAVE_CURL=no \
        USE_SYSTEM_LIBS=no \
        XCFLAGS="-fPIC" \
        libs
)

# Copy libraries
mkdir -p "${LIBS_DIR}"

MUPDF_LIB="${MUPDF_SRC}/build/release/libmupdf.a"
THIRD_LIB="${MUPDF_SRC}/build/release/libmupdf-third.a"

if [ ! -f "${MUPDF_LIB}" ]; then
    echo "Error: ${MUPDF_LIB} not found. Build may have failed."
    exit 1
fi

cp "${MUPDF_LIB}" "${LIBS_DIR}/libmupdf_${OS}_${ARCH}.a"

if [ -f "${THIRD_LIB}" ]; then
    cp "${THIRD_LIB}" "${LIBS_DIR}/libmupdfthird_${OS}_${ARCH}.a"
else
    echo "Warning: libmupdf-third.a not found, creating empty archive"
    ar rcs "${LIBS_DIR}/libmupdfthird_${OS}_${ARCH}.a"
fi

echo ""
echo "Done! Libraries installed:"
ls -lh "${LIBS_DIR}/libmupdf_${OS}_${ARCH}.a" "${LIBS_DIR}/libmupdfthird_${OS}_${ARCH}.a"
echo ""
echo "You can now build gomupdf with: go build ."

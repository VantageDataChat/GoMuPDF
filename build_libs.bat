@echo off
REM build_libs.bat - Build MuPDF static libraries on Windows.
REM
REM Usage:
REM   build_libs.bat [mupdf_source_dir]
REM
REM Prerequisites:
REM   - MinGW-W64 (gcc, ar, make) in PATH
REM   - git (if cloning)
REM
REM Output:
REM   libs\libmupdf_windows_amd64.a
REM   libs\libmupdfthird_windows_amd64.a

setlocal enabledelayedexpansion

set MUPDF_VERSION=1.24.9
set SCRIPT_DIR=%~dp0
set LIBS_DIR=%SCRIPT_DIR%libs

REM Detect architecture
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set ARCH=arm64
) else (
    echo Unsupported architecture: %PROCESSOR_ARCHITECTURE%
    exit /b 1
)

echo Building MuPDF %MUPDF_VERSION% for windows/%ARCH%

REM Source directory
set MUPDF_SRC=%1
if "%MUPDF_SRC%"=="" (
    set MUPDF_SRC=%TEMP%\mupdf-%MUPDF_VERSION%
    if not exist "!MUPDF_SRC!\Makefile" (
        echo Cloning MuPDF %MUPDF_VERSION%...
        git clone --depth 1 --branch %MUPDF_VERSION% --recurse-submodules ^
            https://github.com/ArtifexSoftware/mupdf.git "!MUPDF_SRC!"
    )
)

if not exist "%MUPDF_SRC%\Makefile" (
    echo Error: MuPDF source not found at %MUPDF_SRC%
    exit /b 1
)

echo Building from %MUPDF_SRC%...

pushd "%MUPDF_SRC%"
make -j4 HAVE_X11=no HAVE_GLUT=no HAVE_CURL=no USE_SYSTEM_LIBS=no XCFLAGS="-fPIC" libs
popd

if not exist "%LIBS_DIR%" mkdir "%LIBS_DIR%"

set MUPDF_LIB=%MUPDF_SRC%\build\release\libmupdf.a
set THIRD_LIB=%MUPDF_SRC%\build\release\libmupdf-third.a

if not exist "%MUPDF_LIB%" (
    echo Error: %MUPDF_LIB% not found. Build may have failed.
    exit /b 1
)

copy /Y "%MUPDF_LIB%" "%LIBS_DIR%\libmupdf_windows_%ARCH%.a"

if exist "%THIRD_LIB%" (
    copy /Y "%THIRD_LIB%" "%LIBS_DIR%\libmupdfthird_windows_%ARCH%.a"
) else (
    echo Warning: libmupdf-third.a not found
)

echo.
echo Done! Libraries installed in %LIBS_DIR%
dir "%LIBS_DIR%\libmupdf_windows_%ARCH%.a" "%LIBS_DIR%\libmupdfthird_windows_%ARCH%.a"
echo.
echo You can now build gomupdf with: go build .

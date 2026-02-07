# Makefile for gomupdf
#
# Targets:
#   make libs          - Build MuPDF static libraries for current platform
#   make build         - Build the Go package
#   make test          - Run all tests (requires MuPDF libs)
#   make test-pure     - Run pure Go tests only (no CGO required)
#   make vet           - Run go vet
#   make clean         - Clean build artifacts

MUPDF_SRC ?=

.PHONY: libs build test test-pure vet clean

libs:
	@bash build_libs.sh $(MUPDF_SRC)

build:
	go build .

test:
	go test -v -count=1 .

test-pure:
	go test -v -count=1 -tags nomupdf .

vet:
	go vet .

clean:
	go clean .

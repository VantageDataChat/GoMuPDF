package gomupdf

import "errors"

var (
	// ErrInitFailed is returned when MuPDF context initialization fails.
	ErrInitFailed = errors.New("gomupdf: failed to initialize MuPDF context")

	// ErrOpenFailed is returned when a document cannot be opened.
	ErrOpenFailed = errors.New("gomupdf: failed to open document")

	// ErrPageNotFound is returned when a page number is out of range.
	ErrPageNotFound = errors.New("gomupdf: page not found")

	// ErrNotPDF is returned when a PDF-only operation is attempted on a non-PDF.
	ErrNotPDF = errors.New("gomupdf: not a PDF document")

	// ErrEncrypted is returned when the document is encrypted and not authenticated.
	ErrEncrypted = errors.New("gomupdf: document is encrypted")

	// ErrAuthFailed is returned when password authentication fails.
	ErrAuthFailed = errors.New("gomupdf: authentication failed")

	// ErrClosed is returned when operating on a closed document.
	ErrClosed = errors.New("gomupdf: document is closed")

	// ErrTextExtract is returned when text extraction fails.
	ErrTextExtract = errors.New("gomupdf: text extraction failed")

	// ErrPixmap is returned when pixmap creation fails.
	ErrPixmap = errors.New("gomupdf: pixmap operation failed")

	// ErrSave is returned when saving fails.
	ErrSave = errors.New("gomupdf: save failed")

	// ErrInvalidArg is returned for invalid arguments.
	ErrInvalidArg = errors.New("gomupdf: invalid argument")

	// ErrOutline is returned when outline/TOC operations fail.
	ErrOutline = errors.New("gomupdf: outline operation failed")

	// ErrSearch is returned when search fails.
	ErrSearch = errors.New("gomupdf: search failed")

	// ErrConvert is returned when document conversion fails.
	ErrConvert = errors.New("gomupdf: conversion failed")

	// ErrEmbeddedFile is returned when embedded file operations fail.
	ErrEmbeddedFile = errors.New("gomupdf: embedded file operation failed")

	// ErrXref is returned when xref operations fail.
	ErrXref = errors.New("gomupdf: xref operation failed")
)

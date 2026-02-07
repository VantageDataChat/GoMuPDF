//go:build cgo && !nomupdf

package gomupdf

/*
#include "gomupdf.h"
*/
import "C"
import "unsafe"

// Widget represents a PDF form field widget.
type Widget struct {
	ctx    *context
	widget *C.pdf_annot
	page   *Page
}

func (w *Widget) FieldType() int { return int(C.gomupdf_widget_type(w.ctx.ctx, w.widget)) }

func (w *Widget) FieldTypeString() string {
	names := map[int]string{
		WidgetTypeButton: "Button", WidgetTypeCheckbox: "CheckBox",
		WidgetTypeCombobox: "ComboBox", WidgetTypeListbox: "ListBox",
		WidgetTypeRadioButton: "RadioButton", WidgetTypeSignature: "Signature",
		WidgetTypeText: "Text",
	}
	if name, ok := names[w.FieldType()]; ok {
		return name
	}
	return "Unknown"
}

func (w *Widget) FieldName() string {
	s := C.gomupdf_widget_name(w.ctx.ctx, w.widget)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

func (w *Widget) FieldValue() string {
	s := C.gomupdf_widget_value(w.ctx.ctx, w.widget)
	if s == nil {
		return ""
	}
	return C.GoString(s)
}

func (w *Widget) SetFieldValue(value string) error {
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	errcode := C.gomupdf_set_widget_value(w.ctx.ctx, w.page.doc.pdf, w.widget, cValue)
	if errcode != 0 {
		return ErrInvalidArg
	}
	return nil
}

func (w *Widget) Rect() Rect {
	r := C.gomupdf_annot_rect(w.ctx.ctx, w.widget)
	return Rect{X0: float64(r.x0), Y0: float64(r.y0), X1: float64(r.x1), Y1: float64(r.y1)}
}

func (w *Widget) Xref() int { return int(C.gomupdf_annot_xref(w.ctx.ctx, w.widget)) }

func (p *Page) GetWidgets() []*Widget {
	if !p.doc.IsPDF() {
		return nil
	}
	pdfPage := C.pdf_page_from_fz_page(p.ctx.ctx, p.page)
	if pdfPage == nil {
		return nil
	}
	var widgets []*Widget
	for w := C.gomupdf_first_widget(p.ctx.ctx, pdfPage); w != nil; w = C.gomupdf_next_widget(p.ctx.ctx, w) {
		widgets = append(widgets, &Widget{ctx: p.ctx, widget: w, page: p})
	}
	return widgets
}

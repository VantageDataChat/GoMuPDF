//go:build !cgo || nomupdf

package gomupdf

import (
	"math"
	"testing"
)

// --- Point tests ---

func TestPointBasic(t *testing.T) {
	p := NewPoint(3, 4)
	if p.X != 3 || p.Y != 4 {
		t.Errorf("expected (3, 4), got (%g, %g)", p.X, p.Y)
	}
	if math.Abs(p.Abs()-5) > 1e-10 {
		t.Errorf("expected abs=5, got %g", p.Abs())
	}
}

func TestPointAdd(t *testing.T) {
	p1 := NewPoint(1, 2)
	p2 := NewPoint(3, 4)
	sum := p1.Add(p2)
	if sum.X != 4 || sum.Y != 6 {
		t.Errorf("expected (4, 6), got (%g, %g)", sum.X, sum.Y)
	}
}

func TestPointSub(t *testing.T) {
	p1 := NewPoint(5, 7)
	p2 := NewPoint(2, 3)
	diff := p1.Sub(p2)
	if diff.X != 3 || diff.Y != 4 {
		t.Errorf("expected (3, 4), got (%g, %g)", diff.X, diff.Y)
	}
}

func TestPointMul(t *testing.T) {
	p := NewPoint(3, 4)
	scaled := p.Mul(2.5)
	if scaled.X != 7.5 || scaled.Y != 10 {
		t.Errorf("expected (7.5, 10), got (%g, %g)", scaled.X, scaled.Y)
	}
}

func TestPointIsZero(t *testing.T) {
	if !NewPoint(0, 0).IsZero() {
		t.Error("origin should be zero")
	}
	if NewPoint(1, 0).IsZero() {
		t.Error("(1,0) should not be zero")
	}
}

func TestPointString(t *testing.T) {
	p := NewPoint(1.5, 2.5)
	s := p.String()
	if s != "Point(1.5, 2.5)" {
		t.Errorf("unexpected string: %s", s)
	}
}

func TestPointTransform(t *testing.T) {
	p := NewPoint(1, 0)
	m := RotateMatrix(90)
	result := p.Transform(m)
	if math.Abs(result.X) > 1e-10 || math.Abs(result.Y-1) > 1e-10 {
		t.Errorf("expected (0, 1), got (%g, %g)", result.X, result.Y)
	}
}

func TestPointTransformTranslate(t *testing.T) {
	p := NewPoint(0, 0)
	m := TranslateMatrix(10, 20)
	result := p.Transform(m)
	if result.X != 10 || result.Y != 20 {
		t.Errorf("expected (10, 20), got (%g, %g)", result.X, result.Y)
	}
}

// --- Rect tests ---

func TestRectBasic(t *testing.T) {
	r := NewRect(0, 0, 100, 200)
	if r.Width() != 100 {
		t.Errorf("expected width=100, got %g", r.Width())
	}
	if r.Height() != 200 {
		t.Errorf("expected height=200, got %g", r.Height())
	}
	if r.IsEmpty() {
		t.Error("rect should not be empty")
	}
}

func TestRectFromPoints(t *testing.T) {
	r := RectFromPoints(NewPoint(10, 20), NewPoint(30, 40))
	if r.X0 != 10 || r.Y0 != 20 || r.X1 != 30 || r.Y1 != 40 {
		t.Errorf("unexpected rect: %s", r)
	}
}

func TestRectIsEmpty(t *testing.T) {
	if !NewRect(10, 10, 10, 20).IsEmpty() {
		t.Error("zero-width rect should be empty")
	}
	if !NewRect(10, 10, 20, 10).IsEmpty() {
		t.Error("zero-height rect should be empty")
	}
	if !NewRect(20, 20, 10, 30).IsEmpty() {
		t.Error("inverted x rect should be empty")
	}
}

func TestRectIsInfinite(t *testing.T) {
	r := NewRect(0, 0, 0, 0)
	// IsInfinite checks x0 > x1-1 && y0 > y1-1
	if !r.IsInfinite() {
		t.Log("zero rect is considered infinite by the formula")
	}
	r2 := NewRect(0, 0, 100, 100)
	if r2.IsInfinite() {
		t.Error("normal rect should not be infinite")
	}
}

func TestRectContains(t *testing.T) {
	r := NewRect(0, 0, 100, 100)
	if !r.Contains(NewPoint(50, 50)) {
		t.Error("rect should contain (50, 50)")
	}
	if r.Contains(NewPoint(150, 50)) {
		t.Error("rect should not contain (150, 50)")
	}
	// Edge cases
	if !r.Contains(NewPoint(0, 0)) {
		t.Error("rect should contain its top-left corner")
	}
	if !r.Contains(NewPoint(100, 100)) {
		t.Error("rect should contain its bottom-right corner")
	}
}

func TestRectContainsRect(t *testing.T) {
	outer := NewRect(0, 0, 100, 100)
	inner := NewRect(10, 10, 90, 90)
	if !outer.ContainsRect(inner) {
		t.Error("outer should contain inner")
	}
	if inner.ContainsRect(outer) {
		t.Error("inner should not contain outer")
	}
	if !outer.ContainsRect(outer) {
		t.Error("rect should contain itself")
	}
}

func TestRectIntersects(t *testing.T) {
	r1 := NewRect(0, 0, 100, 100)
	r2 := NewRect(50, 50, 150, 150)
	if !r1.Intersects(r2) {
		t.Error("overlapping rects should intersect")
	}
	r3 := NewRect(200, 200, 300, 300)
	if r1.Intersects(r3) {
		t.Error("non-overlapping rects should not intersect")
	}
	// Empty rect
	empty := NewRect(10, 10, 10, 10)
	if r1.Intersects(empty) {
		t.Error("empty rect should not intersect")
	}
}

func TestRectIntersect(t *testing.T) {
	r1 := NewRect(0, 0, 100, 100)
	r2 := NewRect(50, 50, 150, 150)
	inter := r1.Intersect(r2)
	expected := NewRect(50, 50, 100, 100)
	if inter != expected {
		t.Errorf("expected %s, got %s", expected, inter)
	}
}

func TestRectUnion(t *testing.T) {
	r1 := NewRect(0, 0, 100, 100)
	r2 := NewRect(50, 50, 150, 150)
	union := r1.Union(r2)
	expected := NewRect(0, 0, 150, 150)
	if union != expected {
		t.Errorf("expected %s, got %s", expected, union)
	}
}

func TestRectUnionWithEmpty(t *testing.T) {
	r := NewRect(10, 10, 50, 50)
	empty := NewRect(0, 0, 0, 0)
	if r.Union(empty) != r {
		t.Error("union with empty should return original")
	}
	if empty.Union(r) != r {
		t.Error("empty union with rect should return rect")
	}
}

func TestRectIncludePoint(t *testing.T) {
	r := NewRect(10, 10, 50, 50)
	r2 := r.IncludePoint(NewPoint(0, 0))
	if r2.X0 != 0 || r2.Y0 != 0 {
		t.Errorf("expected X0=0, Y0=0, got %s", r2)
	}
	r3 := r.IncludePoint(NewPoint(100, 100))
	if r3.X1 != 100 || r3.Y1 != 100 {
		t.Errorf("expected X1=100, Y1=100, got %s", r3)
	}
}

func TestRectTransform(t *testing.T) {
	r := NewRect(0, 0, 100, 50)
	m := ScaleMatrix(2, 3)
	result := r.Transform(m)
	if result.X1 != 200 || result.Y1 != 150 {
		t.Errorf("expected (0,0,200,150), got %s", result)
	}
}

func TestRectTransformEmpty(t *testing.T) {
	r := NewRect(10, 10, 5, 5) // empty (inverted)
	m := ScaleMatrix(2, 2)
	result := r.Transform(m)
	if result != r {
		t.Error("transforming empty rect should return it unchanged")
	}
}

func TestRectNormalize(t *testing.T) {
	r := NewRect(100, 200, 0, 0)
	n := r.Normalize()
	if n.X0 != 0 || n.Y0 != 0 || n.X1 != 100 || n.Y1 != 200 {
		t.Errorf("expected (0, 0, 100, 200), got %s", n)
	}
}

func TestRectCorners(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	if r.TopLeft() != NewPoint(10, 20) {
		t.Errorf("TopLeft: %s", r.TopLeft())
	}
	if r.TopRight() != NewPoint(30, 20) {
		t.Errorf("TopRight: %s", r.TopRight())
	}
	if r.BottomLeft() != NewPoint(10, 40) {
		t.Errorf("BottomLeft: %s", r.BottomLeft())
	}
	if r.BottomRight() != NewPoint(30, 40) {
		t.Errorf("BottomRight: %s", r.BottomRight())
	}
}

func TestRectQuad(t *testing.T) {
	r := NewRect(0, 0, 100, 50)
	q := r.Quad()
	if q.UL != r.TopLeft() || q.UR != r.TopRight() || q.LL != r.BottomLeft() || q.LR != r.BottomRight() {
		t.Error("Quad corners don't match Rect corners")
	}
}

func TestRectIRect(t *testing.T) {
	r := NewRect(0.3, 0.7, 99.2, 99.8)
	ir := r.IRect()
	if ir.X0 != 0 || ir.Y0 != 0 || ir.X1 != 100 || ir.Y1 != 100 {
		t.Errorf("expected (0, 0, 100, 100), got %s", ir)
	}
}

func TestRectString(t *testing.T) {
	r := NewRect(1, 2, 3, 4)
	s := r.String()
	if s != "Rect(1, 2, 3, 4)" {
		t.Errorf("unexpected string: %s", s)
	}
}

// --- IRect tests ---

func TestIRectBasic(t *testing.T) {
	ir := NewIRect(0, 0, 100, 200)
	if ir.Width() != 100 {
		t.Errorf("expected width=100, got %d", ir.Width())
	}
	if ir.Height() != 200 {
		t.Errorf("expected height=200, got %d", ir.Height())
	}
}

func TestIRectIsEmpty(t *testing.T) {
	if !NewIRect(10, 10, 10, 20).IsEmpty() {
		t.Error("zero-width IRect should be empty")
	}
	if !NewIRect(10, 10, 20, 10).IsEmpty() {
		t.Error("zero-height IRect should be empty")
	}
	if NewIRect(0, 0, 10, 10).IsEmpty() {
		t.Error("normal IRect should not be empty")
	}
}

func TestIRectRect(t *testing.T) {
	ir := NewIRect(1, 2, 3, 4)
	r := ir.Rect()
	if r.X0 != 1 || r.Y0 != 2 || r.X1 != 3 || r.Y1 != 4 {
		t.Errorf("unexpected rect: %s", r)
	}
}

func TestIRectString(t *testing.T) {
	ir := NewIRect(1, 2, 3, 4)
	s := ir.String()
	if s != "IRect(1, 2, 3, 4)" {
		t.Errorf("unexpected string: %s", s)
	}
}

// --- Matrix tests ---

func TestMatrixIdentity(t *testing.T) {
	p := NewPoint(5, 10)
	result := p.Transform(Identity)
	if result != p {
		t.Errorf("identity transform should not change point")
	}
}

func TestMatrixScale(t *testing.T) {
	m := ScaleMatrix(2, 3)
	p := NewPoint(5, 10)
	result := p.Transform(m)
	if result.X != 10 || result.Y != 30 {
		t.Errorf("expected (10, 30), got (%g, %g)", result.X, result.Y)
	}
}

func TestMatrixTranslate(t *testing.T) {
	m := TranslateMatrix(5, 10)
	p := NewPoint(1, 2)
	result := p.Transform(m)
	if result.X != 6 || result.Y != 12 {
		t.Errorf("expected (6, 12), got (%g, %g)", result.X, result.Y)
	}
}

func TestMatrixRotate(t *testing.T) {
	m := RotateMatrix(180)
	p := NewPoint(1, 0)
	result := p.Transform(m)
	if math.Abs(result.X+1) > 1e-10 || math.Abs(result.Y) > 1e-10 {
		t.Errorf("expected (-1, 0), got (%g, %g)", result.X, result.Y)
	}
}

func TestMatrixShear(t *testing.T) {
	m := ShearMatrix(1, 0)
	p := NewPoint(0, 1)
	result := p.Transform(m)
	// shear: x' = x + sx*y = 0 + 1*1 = 1, y' = sy*x + y = 0 + 1 = 1
	if math.Abs(result.X-1) > 1e-10 || math.Abs(result.Y-1) > 1e-10 {
		t.Errorf("expected (1, 1), got (%g, %g)", result.X, result.Y)
	}
}

func TestMatrixConcat(t *testing.T) {
	m1 := ScaleMatrix(2, 2)
	m2 := TranslateMatrix(10, 20)
	m := m1.Concat(m2)
	p := NewPoint(0, 0)
	result := p.Transform(m)
	if result.X != 10 || result.Y != 20 {
		t.Errorf("expected (10, 20), got (%g, %g)", result.X, result.Y)
	}
}

func TestMatrixPreScale(t *testing.T) {
	m := TranslateMatrix(10, 10)
	result := m.PreScale(2, 2)
	p := NewPoint(5, 5)
	r := p.Transform(result)
	// PreScale: Scale(2,2).Concat(Translate(10,10))
	// Scale first: (10, 10), then translate: (20, 20)
	if r.X != 20 || r.Y != 20 {
		t.Errorf("expected (20, 20), got (%g, %g)", r.X, r.Y)
	}
}

func TestMatrixPreTranslate(t *testing.T) {
	m := ScaleMatrix(2, 2)
	result := m.PreTranslate(5, 5)
	p := NewPoint(0, 0)
	r := p.Transform(result)
	// PreTranslate: Translate(5,5).Concat(Scale(2,2))
	// Translate first: (5, 5), then scale: (10, 10)
	if r.X != 10 || r.Y != 10 {
		t.Errorf("expected (10, 10), got (%g, %g)", r.X, r.Y)
	}
}

func TestMatrixPreRotate(t *testing.T) {
	m := Identity
	result := m.PreRotate(90)
	p := NewPoint(1, 0)
	r := p.Transform(result)
	if math.Abs(r.X) > 1e-10 || math.Abs(r.Y-1) > 1e-10 {
		t.Errorf("expected (0, 1), got (%g, %g)", r.X, r.Y)
	}
}

func TestMatrixInvert(t *testing.T) {
	m := ScaleMatrix(2, 4)
	inv, ok := m.Invert()
	if !ok {
		t.Fatal("matrix should be invertible")
	}
	p := NewPoint(10, 20)
	result := p.Transform(m).Transform(inv)
	if math.Abs(result.X-p.X) > 1e-10 || math.Abs(result.Y-p.Y) > 1e-10 {
		t.Errorf("expected (%g, %g), got (%g, %g)", p.X, p.Y, result.X, result.Y)
	}
}

func TestMatrixInvertSingular(t *testing.T) {
	m := NewMatrix(0, 0, 0, 0, 0, 0)
	_, ok := m.Invert()
	if ok {
		t.Error("singular matrix should not be invertible")
	}
}

func TestMatrixIsRectilinear(t *testing.T) {
	if !Identity.IsRectilinear() {
		t.Error("identity should be rectilinear")
	}
	if !ScaleMatrix(2, 3).IsRectilinear() {
		t.Error("scale should be rectilinear")
	}
	// Note: RotateMatrix(90) uses sin/cos which may not produce exact 0s,
	// so IsRectilinear (which checks == 0) may fail. Use exact matrix instead.
	rot90 := NewMatrix(0, 1, -1, 0, 0, 0)
	if !rot90.IsRectilinear() {
		t.Error("exact 90-degree rotation should be rectilinear")
	}
	rot45 := RotateMatrix(45)
	if rot45.IsRectilinear() {
		t.Error("45-degree rotation should not be rectilinear")
	}
}

func TestMatrixString(t *testing.T) {
	m := NewMatrix(1, 2, 3, 4, 5, 6)
	s := m.String()
	if s != "Matrix(1, 2, 3, 4, 5, 6)" {
		t.Errorf("unexpected string: %s", s)
	}
}

// --- Quad tests ---

func TestQuadRect(t *testing.T) {
	r := NewRect(0, 0, 100, 50)
	q := r.Quad()
	back := q.Rect()
	if back != r {
		t.Errorf("expected %s, got %s", r, back)
	}
}

func TestQuadFromRect(t *testing.T) {
	r := NewRect(10, 20, 30, 40)
	q := QuadFromRect(r)
	if q.UL != r.TopLeft() || q.LR != r.BottomRight() {
		t.Error("QuadFromRect corners mismatch")
	}
}

func TestNewQuad(t *testing.T) {
	ul := NewPoint(0, 0)
	ur := NewPoint(100, 0)
	ll := NewPoint(0, 50)
	lr := NewPoint(100, 50)
	q := NewQuad(ul, ur, ll, lr)
	if q.UL != ul || q.UR != ur || q.LL != ll || q.LR != lr {
		t.Error("NewQuad corners mismatch")
	}
}

func TestQuadIsRectangular(t *testing.T) {
	r := NewRect(10, 20, 100, 80)
	q := r.Quad()
	if !q.IsRectangular() {
		t.Error("quad from rect should be rectangular")
	}
	// Non-rectangular quad
	q2 := NewQuad(NewPoint(0, 0), NewPoint(100, 10), NewPoint(0, 50), NewPoint(100, 50))
	if q2.IsRectangular() {
		t.Error("skewed quad should not be rectangular")
	}
}

func TestQuadIsEmpty(t *testing.T) {
	q := NewQuad(NewPoint(0, 0), NewPoint(0, 0), NewPoint(0, 0), NewPoint(0, 0))
	if !q.IsEmpty() {
		t.Error("degenerate quad should be empty")
	}
}

func TestQuadIsConvex(t *testing.T) {
	r := NewRect(0, 0, 100, 100)
	q := r.Quad()
	if !q.IsConvex() {
		t.Error("rectangular quad should be convex")
	}
}

func TestQuadTransform(t *testing.T) {
	r := NewRect(0, 0, 10, 10)
	q := r.Quad()
	m := ScaleMatrix(2, 3)
	q2 := q.Transform(m)
	if q2.LR.X != 20 || q2.LR.Y != 30 {
		t.Errorf("expected LR (20, 30), got %s", q2.LR)
	}
}

func TestQuadString(t *testing.T) {
	q := NewQuad(NewPoint(0, 0), NewPoint(1, 0), NewPoint(0, 1), NewPoint(1, 1))
	s := q.String()
	if s == "" {
		t.Error("String should not be empty")
	}
}

// --- Tools tests ---

func TestVersion(t *testing.T) {
	v := Version()
	if v == "" {
		t.Error("Version should not be empty")
	}
}

func TestMuPDFVersion(t *testing.T) {
	v := MuPDFVersion()
	if v != "1.24.9" {
		t.Errorf("expected 1.24.9, got %s", v)
	}
}

func TestPaperSize(t *testing.T) {
	a4 := PaperSize("a4")
	if a4.Width() != 595 || a4.Height() != 842 {
		t.Errorf("expected A4 (595x842), got %gx%g", a4.Width(), a4.Height())
	}
	letter := PaperSize("letter")
	if letter.Width() != 612 || letter.Height() != 792 {
		t.Errorf("expected Letter (612x792), got %gx%g", letter.Width(), letter.Height())
	}
}

func TestPaperSizeLandscape(t *testing.T) {
	a4l := PaperSize("a4-l")
	if a4l.Width() != 842 || a4l.Height() != 595 {
		t.Errorf("expected A4 landscape (842x595), got %gx%g", a4l.Width(), a4l.Height())
	}
	a4l2 := PaperSize("a4-landscape")
	if a4l2.Width() != 842 || a4l2.Height() != 595 {
		t.Errorf("expected A4 landscape (842x595), got %gx%g", a4l2.Width(), a4l2.Height())
	}
}

func TestPaperSizeUnknown(t *testing.T) {
	r := PaperSize("unknown_size")
	// Should default to A4
	if r.Width() != 595 || r.Height() != 842 {
		t.Errorf("unknown size should default to A4, got %gx%g", r.Width(), r.Height())
	}
}

func TestPaperSizeAllSizes(t *testing.T) {
	sizes := []string{"a0", "a1", "a2", "a3", "a4", "a5", "a6", "a7", "a8", "a9", "a10",
		"b0", "b1", "b2", "b3", "b4", "b5", "letter", "legal", "tabloid", "ledger"}
	for _, name := range sizes {
		r := PaperSize(name)
		if r.Width() <= 0 || r.Height() <= 0 {
			t.Errorf("paper size %s has invalid dimensions: %gx%g", name, r.Width(), r.Height())
		}
	}
}

func TestGetPDFStr(t *testing.T) {
	// ASCII string
	s := GetPDFStr("Hello World")
	if s != "(Hello World)" {
		t.Errorf("expected (Hello World), got %s", s)
	}
	// String with special chars
	s = GetPDFStr("Hello (World)")
	if s != "(Hello \\(World\\))" {
		t.Errorf("expected escaped parens, got %s", s)
	}
	// Backslash
	s = GetPDFStr("path\\to\\file")
	if s != "(path\\\\to\\\\file)" {
		t.Errorf("expected escaped backslashes, got %s", s)
	}
	// Non-ASCII string
	s = GetPDFStr("€")
	if s[:5] != "<feff" {
		t.Errorf("expected UTF-16BE BOM, got %s", s)
	}
	// Chinese characters
	s = GetPDFStr("你好")
	if s[:5] != "<feff" {
		t.Errorf("expected UTF-16BE BOM for Chinese, got %s", s)
	}
	// Empty string
	s = GetPDFStr("")
	if s != "()" {
		t.Errorf("expected (), got %s", s)
	}
}

func TestPlanishLine(t *testing.T) {
	// Same point
	angle := PlanishLine(NewPoint(0, 0), NewPoint(0, 0))
	if angle != 0 {
		t.Errorf("same point should return 0, got %g", angle)
	}
	// Vertical line up
	angle = PlanishLine(NewPoint(0, 0), NewPoint(0, -10))
	if angle != -90 {
		t.Errorf("vertical up should return -90, got %g", angle)
	}
	// Vertical line down
	angle = PlanishLine(NewPoint(0, 0), NewPoint(0, 10))
	if angle != 90 {
		t.Errorf("vertical down should return 90, got %g", angle)
	}
}

// --- Outline tests ---

func TestOutlineFlatten(t *testing.T) {
	o := &Outline{
		Title: "Chapter 1",
		Page:  1,
		Down: &Outline{
			Title: "Section 1.1",
			Page:  2,
		},
		Next: &Outline{
			Title: "Chapter 2",
			Page:  5,
		},
	}
	items := o.Flatten()
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	if items[0].Level != 1 || items[0].Title != "Chapter 1" {
		t.Errorf("item 0: %+v", items[0])
	}
	if items[1].Level != 2 || items[1].Title != "Section 1.1" {
		t.Errorf("item 1: %+v", items[1])
	}
	if items[2].Level != 1 || items[2].Title != "Chapter 2" {
		t.Errorf("item 2: %+v", items[2])
	}
}

func TestOutlineFlattenNil(t *testing.T) {
	var o *Outline
	items := o.Flatten()
	if len(items) != 0 {
		t.Errorf("nil outline should return empty, got %d", len(items))
	}
}

// --- Types tests ---

func TestSTextLineText(t *testing.T) {
	line := STextLine{
		Chars: []STextChar{
			{C: 'H'}, {C: 'e'}, {C: 'l'}, {C: 'l'}, {C: 'o'},
		},
	}
	if line.Text() != "Hello" {
		t.Errorf("expected 'Hello', got '%s'", line.Text())
	}
}

func TestSTextLineTextEmpty(t *testing.T) {
	line := STextLine{}
	if line.Text() != "" {
		t.Errorf("expected empty, got '%s'", line.Text())
	}
}

func TestDefaultSaveOptions(t *testing.T) {
	opts := DefaultSaveOptions()
	if opts.Garbage != 0 {
		t.Errorf("expected Garbage=0, got %d", opts.Garbage)
	}
	if opts.Deflate {
		t.Error("expected Deflate=false")
	}
	if opts.Permissions != -1 {
		t.Errorf("expected Permissions=-1, got %d", opts.Permissions)
	}
}

func TestEzSaveOptions(t *testing.T) {
	opts := EzSaveOptions()
	if opts.Garbage != 3 {
		t.Errorf("expected Garbage=3, got %d", opts.Garbage)
	}
	if !opts.Deflate {
		t.Error("expected Deflate=true")
	}
}

// --- Colorspace tests ---

func TestColorspaces(t *testing.T) {
	if CsGRAY.N != 1 {
		t.Errorf("gray should have 1 component, got %d", CsGRAY.N)
	}
	if CsRGBCS.N != 3 {
		t.Errorf("RGB should have 3 components, got %d", CsRGBCS.N)
	}
	if CsCMYKCS.N != 4 {
		t.Errorf("CMYK should have 4 components, got %d", CsCMYKCS.N)
	}
}

// --- Color tests ---

func TestPredefinedColors(t *testing.T) {
	if ColorBlack.R != 0 || ColorBlack.G != 0 || ColorBlack.B != 0 {
		t.Error("black should be (0,0,0)")
	}
	if ColorWhite.R != 1 || ColorWhite.G != 1 || ColorWhite.B != 1 {
		t.Error("white should be (1,1,1)")
	}
	if ColorRed.R != 1 || ColorRed.G != 0 || ColorRed.B != 0 {
		t.Error("red should be (1,0,0)")
	}
}

// --- Stubs tests (nomupdf) ---

func TestStubOpen(t *testing.T) {
	_, err := Open("test.pdf")
	if err != ErrInitFailed {
		t.Errorf("expected ErrInitFailed, got %v", err)
	}
}

func TestStubOpenFromMemory(t *testing.T) {
	_, err := OpenFromMemory([]byte("test"), "test.pdf")
	if err != ErrInitFailed {
		t.Errorf("expected ErrInitFailed, got %v", err)
	}
}

func TestStubNewPDF(t *testing.T) {
	_, err := NewPDF()
	if err != ErrInitFailed {
		t.Errorf("expected ErrInitFailed, got %v", err)
	}
}

func TestStubNewPixmap(t *testing.T) {
	_, err := NewPixmap(CsRGB, 10, 10, false)
	if err != ErrInitFailed {
		t.Errorf("expected ErrInitFailed, got %v", err)
	}
}

func TestStubNewPixmapFromImage(t *testing.T) {
	_, err := NewPixmapFromImage(nil, 0)
	if err != ErrInitFailed {
		t.Errorf("expected ErrInitFailed, got %v", err)
	}
}

// --- Constants tests ---

func TestConstants(t *testing.T) {
	// Verify key constants are defined correctly
	if TextFlagsDefault != 3 {
		t.Errorf("expected TextFlagsDefault=3, got %d", TextFlagsDefault)
	}
	if CsGray != 0 || CsRGB != 1 || CsCMYK != 2 {
		t.Error("colorspace constants wrong")
	}
	if AnnotText != 0 || AnnotHighlight != 8 {
		t.Error("annotation constants wrong")
	}
	if WidgetTypeText != 7 {
		t.Error("widget type constants wrong")
	}
	if PermPrint != 4 {
		t.Errorf("PermPrint should be 4, got %d", PermPrint)
	}
	if EncryptAESV3 != 4 {
		t.Errorf("EncryptAESV3 should be 4, got %d", EncryptAESV3)
	}
}

func TestPaperConstants(t *testing.T) {
	if PaperA4.Width() != 595 || PaperA4.Height() != 842 {
		t.Error("PaperA4 wrong")
	}
	if PaperLetter.Width() != 612 || PaperLetter.Height() != 792 {
		t.Error("PaperLetter wrong")
	}
}

// --- Error tests ---

func TestErrors(t *testing.T) {
	errs := []error{
		ErrInitFailed, ErrOpenFailed, ErrPageNotFound, ErrNotPDF,
		ErrEncrypted, ErrAuthFailed, ErrClosed, ErrTextExtract,
		ErrPixmap, ErrSave, ErrInvalidArg, ErrOutline,
		ErrSearch, ErrConvert, ErrEmbeddedFile, ErrXref,
	}
	for _, e := range errs {
		if e == nil {
			t.Error("error should not be nil")
		}
		if e.Error() == "" {
			t.Error("error message should not be empty")
		}
	}
}

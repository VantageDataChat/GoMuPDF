package gomupdf

import (
	"fmt"
	"math"
)

// Point represents a 2D point (x, y).
// Corresponds to PyMuPDF's fitz.Point.
type Point struct {
	X, Y float64
}

// NewPoint creates a new Point.
func NewPoint(x, y float64) Point {
	return Point{X: x, Y: y}
}

// Add returns the sum of two points.
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Sub returns the difference of two points.
func (p Point) Sub(other Point) Point {
	return Point{X: p.X - other.X, Y: p.Y - other.Y}
}

// Mul returns the point scaled by a factor.
func (p Point) Mul(factor float64) Point {
	return Point{X: p.X * factor, Y: p.Y * factor}
}

// Abs returns the Euclidean norm (distance from origin).
func (p Point) Abs() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Transform applies a matrix transformation to the point.
func (p Point) Transform(m Matrix) Point {
	return Point{
		X: p.X*m.A + p.Y*m.C + m.E,
		Y: p.X*m.B + p.Y*m.D + m.F,
	}
}

// IsZero returns true if the point is at the origin.
func (p Point) IsZero() bool {
	return p.X == 0 && p.Y == 0
}

// String returns a string representation.
func (p Point) String() string {
	return fmt.Sprintf("Point(%g, %g)", p.X, p.Y)
}

// Rect represents a rectangle defined by two corner points.
// Corresponds to PyMuPDF's fitz.Rect.
type Rect struct {
	X0, Y0, X1, Y1 float64
}

// NewRect creates a new Rect.
func NewRect(x0, y0, x1, y1 float64) Rect {
	return Rect{X0: x0, Y0: y0, X1: x1, Y1: y1}
}

// RectFromPoints creates a Rect from two Points.
func RectFromPoints(topLeft, bottomRight Point) Rect {
	return Rect{X0: topLeft.X, Y0: topLeft.Y, X1: bottomRight.X, Y1: bottomRight.Y}
}

// Width returns the width of the rectangle.
func (r Rect) Width() float64 {
	return math.Abs(r.X1 - r.X0)
}

// Height returns the height of the rectangle.
func (r Rect) Height() float64 {
	return math.Abs(r.Y1 - r.Y0)
}

// IsEmpty returns true if the rectangle has zero or negative area.
func (r Rect) IsEmpty() bool {
	return r.X0 >= r.X1 || r.Y0 >= r.Y1
}

// IsInfinite returns true if the rectangle is the infinite rectangle.
func (r Rect) IsInfinite() bool {
	return r.X0 > r.X1-1 && r.Y0 > r.Y1-1
}

// Contains returns true if the point is inside the rectangle.
func (r Rect) Contains(p Point) bool {
	return p.X >= r.X0 && p.X <= r.X1 && p.Y >= r.Y0 && p.Y <= r.Y1
}

// ContainsRect returns true if other is fully inside this rectangle.
func (r Rect) ContainsRect(other Rect) bool {
	return other.X0 >= r.X0 && other.Y0 >= r.Y0 && other.X1 <= r.X1 && other.Y1 <= r.Y1
}

// Intersects returns true if the two rectangles overlap.
func (r Rect) Intersects(other Rect) bool {
	if r.IsEmpty() || other.IsEmpty() {
		return false
	}
	return r.X0 < other.X1 && r.X1 > other.X0 && r.Y0 < other.Y1 && r.Y1 > other.Y0
}

// Intersect returns the intersection of two rectangles.
func (r Rect) Intersect(other Rect) Rect {
	return Rect{
		X0: math.Max(r.X0, other.X0),
		Y0: math.Max(r.Y0, other.Y0),
		X1: math.Min(r.X1, other.X1),
		Y1: math.Min(r.Y1, other.Y1),
	}
}

// Union returns the smallest rectangle containing both rectangles.
func (r Rect) Union(other Rect) Rect {
	if r.IsEmpty() {
		return other
	}
	if other.IsEmpty() {
		return r
	}
	return Rect{
		X0: math.Min(r.X0, other.X0),
		Y0: math.Min(r.Y0, other.Y0),
		X1: math.Max(r.X1, other.X1),
		Y1: math.Max(r.Y1, other.Y1),
	}
}

// IncludePoint returns the smallest rectangle containing the rect and the point.
func (r Rect) IncludePoint(p Point) Rect {
	return Rect{
		X0: math.Min(r.X0, p.X),
		Y0: math.Min(r.Y0, p.Y),
		X1: math.Max(r.X1, p.X),
		Y1: math.Max(r.Y1, p.Y),
	}
}

// Transform applies a matrix transformation to the rectangle.
func (r Rect) Transform(m Matrix) Rect {
	if r.IsEmpty() {
		return r
	}
	p1 := Point{r.X0, r.Y0}.Transform(m)
	p2 := Point{r.X1, r.Y0}.Transform(m)
	p3 := Point{r.X1, r.Y1}.Transform(m)
	p4 := Point{r.X0, r.Y1}.Transform(m)
	return Rect{
		X0: math.Min(math.Min(p1.X, p2.X), math.Min(p3.X, p4.X)),
		Y0: math.Min(math.Min(p1.Y, p2.Y), math.Min(p3.Y, p4.Y)),
		X1: math.Max(math.Max(p1.X, p2.X), math.Max(p3.X, p4.X)),
		Y1: math.Max(math.Max(p1.Y, p2.Y), math.Max(p3.Y, p4.Y)),
	}
}

// Normalize ensures x0 <= x1 and y0 <= y1.
func (r Rect) Normalize() Rect {
	return Rect{
		X0: math.Min(r.X0, r.X1),
		Y0: math.Min(r.Y0, r.Y1),
		X1: math.Max(r.X0, r.X1),
		Y1: math.Max(r.Y0, r.Y1),
	}
}

// TopLeft returns the top-left corner point.
func (r Rect) TopLeft() Point {
	return Point{r.X0, r.Y0}
}

// TopRight returns the top-right corner point.
func (r Rect) TopRight() Point {
	return Point{r.X1, r.Y0}
}

// BottomLeft returns the bottom-left corner point.
func (r Rect) BottomLeft() Point {
	return Point{r.X0, r.Y1}
}

// BottomRight returns the bottom-right corner point.
func (r Rect) BottomRight() Point {
	return Point{r.X1, r.Y1}
}

// Quad returns the rectangle as a Quad.
func (r Rect) Quad() Quad {
	return Quad{
		UL: r.TopLeft(),
		UR: r.TopRight(),
		LL: r.BottomLeft(),
		LR: r.BottomRight(),
	}
}

// IRect returns the integer version of this rectangle (rounding outward).
func (r Rect) IRect() IRect {
	return IRect{
		X0: int(math.Floor(r.X0)),
		Y0: int(math.Floor(r.Y0)),
		X1: int(math.Ceil(r.X1)),
		Y1: int(math.Ceil(r.Y1)),
	}
}

// String returns a string representation.
func (r Rect) String() string {
	return fmt.Sprintf("Rect(%g, %g, %g, %g)", r.X0, r.Y0, r.X1, r.Y1)
}

// IRect represents an integer rectangle.
// Corresponds to PyMuPDF's fitz.IRect.
type IRect struct {
	X0, Y0, X1, Y1 int
}

// NewIRect creates a new IRect.
func NewIRect(x0, y0, x1, y1 int) IRect {
	return IRect{X0: x0, Y0: y0, X1: x1, Y1: y1}
}

// Width returns the width.
func (r IRect) Width() int {
	return r.X1 - r.X0
}

// Height returns the height.
func (r IRect) Height() int {
	return r.Y1 - r.Y0
}

// IsEmpty returns true if the rectangle has zero or negative area.
func (r IRect) IsEmpty() bool {
	return r.X0 >= r.X1 || r.Y0 >= r.Y1
}

// Rect returns the float version of this rectangle.
func (r IRect) Rect() Rect {
	return Rect{
		X0: float64(r.X0),
		Y0: float64(r.Y0),
		X1: float64(r.X1),
		Y1: float64(r.Y1),
	}
}

// String returns a string representation.
func (r IRect) String() string {
	return fmt.Sprintf("IRect(%d, %d, %d, %d)", r.X0, r.Y0, r.X1, r.Y1)
}

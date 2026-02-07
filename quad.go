package gomupdf

import "fmt"

// Quad represents a quadrilateral defined by four corner points.
// Corresponds to PyMuPDF's fitz.Quad.
type Quad struct {
	UL, UR, LL, LR Point // upper-left, upper-right, lower-left, lower-right
}

// NewQuad creates a new Quad from four points.
func NewQuad(ul, ur, ll, lr Point) Quad {
	return Quad{UL: ul, UR: ur, LL: ll, LR: lr}
}

// QuadFromRect creates a Quad from a Rect.
func QuadFromRect(r Rect) Quad {
	return r.Quad()
}

// Rect returns the smallest enclosing rectangle.
func (q Quad) Rect() Rect {
	r := Rect{X0: q.UL.X, Y0: q.UL.Y, X1: q.UL.X, Y1: q.UL.Y}
	r = r.IncludePoint(q.UR)
	r = r.IncludePoint(q.LL)
	r = r.IncludePoint(q.LR)
	return r
}

// IsEmpty returns true if the quad has zero area.
func (q Quad) IsEmpty() bool {
	return q.Rect().IsEmpty()
}

// IsRectangular returns true if the quad is a rectangle.
func (q Quad) IsRectangular() bool {
	r := q.Rect()
	return q.UL == r.TopLeft() && q.UR == r.TopRight() &&
		q.LL == r.BottomLeft() && q.LR == r.BottomRight()
}

// IsConvex returns true if the quad is convex.
func (q Quad) IsConvex() bool {
	return crossProduct(q.UL, q.UR, q.LR) >= 0 &&
		crossProduct(q.UR, q.LR, q.LL) >= 0 &&
		crossProduct(q.LR, q.LL, q.UL) >= 0 &&
		crossProduct(q.LL, q.UL, q.UR) >= 0
}

// Transform applies a matrix transformation to the quad.
func (q Quad) Transform(m Matrix) Quad {
	return Quad{
		UL: q.UL.Transform(m),
		UR: q.UR.Transform(m),
		LL: q.LL.Transform(m),
		LR: q.LR.Transform(m),
	}
}

// String returns a string representation.
func (q Quad) String() string {
	return fmt.Sprintf("Quad(%s, %s, %s, %s)", q.UL, q.UR, q.LL, q.LR)
}

func crossProduct(a, b, c Point) float64 {
	return (b.X-a.X)*(c.Y-a.Y) - (b.Y-a.Y)*(c.X-a.X)
}

package gomupdf

import (
	"fmt"
	"math"
)

// Matrix represents a 3x3 transformation matrix (using 6 values).
// Corresponds to PyMuPDF's fitz.Matrix.
//
//	| A B 0 |
//	| C D 0 |
//	| E F 1 |
type Matrix struct {
	A, B, C, D, E, F float64
}

// Identity is the identity matrix.
var Identity = Matrix{A: 1, B: 0, C: 0, D: 1, E: 0, F: 0}

// NewMatrix creates a new Matrix from 6 values.
func NewMatrix(a, b, c, d, e, f float64) Matrix {
	return Matrix{A: a, B: b, C: c, D: d, E: e, F: f}
}

// ScaleMatrix creates a scaling matrix.
func ScaleMatrix(sx, sy float64) Matrix {
	return Matrix{A: sx, B: 0, C: 0, D: sy, E: 0, F: 0}
}

// TranslateMatrix creates a translation matrix.
func TranslateMatrix(tx, ty float64) Matrix {
	return Matrix{A: 1, B: 0, C: 0, D: 1, E: tx, F: ty}
}

// RotateMatrix creates a rotation matrix (angle in degrees).
func RotateMatrix(deg float64) Matrix {
	rad := deg * math.Pi / 180.0
	s := math.Sin(rad)
	c := math.Cos(rad)
	return Matrix{A: c, B: s, C: -s, D: c, E: 0, F: 0}
}

// ShearMatrix creates a shear matrix.
func ShearMatrix(sx, sy float64) Matrix {
	return Matrix{A: 1, B: sy, C: sx, D: 1, E: 0, F: 0}
}

// Concat returns the product of two matrices (this * other).
func (m Matrix) Concat(other Matrix) Matrix {
	return Matrix{
		A: m.A*other.A + m.B*other.C,
		B: m.A*other.B + m.B*other.D,
		C: m.C*other.A + m.D*other.C,
		D: m.C*other.B + m.D*other.D,
		E: m.E*other.A + m.F*other.C + other.E,
		F: m.E*other.B + m.F*other.D + other.F,
	}
}

// PreScale returns the matrix pre-multiplied by a scaling matrix.
func (m Matrix) PreScale(sx, sy float64) Matrix {
	return ScaleMatrix(sx, sy).Concat(m)
}

// PreTranslate returns the matrix pre-multiplied by a translation matrix.
func (m Matrix) PreTranslate(tx, ty float64) Matrix {
	return TranslateMatrix(tx, ty).Concat(m)
}

// PreRotate returns the matrix pre-multiplied by a rotation matrix.
func (m Matrix) PreRotate(deg float64) Matrix {
	return RotateMatrix(deg).Concat(m)
}

// Invert returns the inverse of the matrix, or Identity if not invertible.
func (m Matrix) Invert() (Matrix, bool) {
	det := m.A*m.D - m.B*m.C
	if det == 0 {
		return Identity, false
	}
	invDet := 1.0 / det
	return Matrix{
		A: m.D * invDet,
		B: -m.B * invDet,
		C: -m.C * invDet,
		D: m.A * invDet,
		E: (m.C*m.F - m.D*m.E) * invDet,
		F: (m.B*m.E - m.A*m.F) * invDet,
	}, true
}

// IsRectilinear returns true if the matrix maps rectangles to rectangles.
func (m Matrix) IsRectilinear() bool {
	return (m.B == 0 && m.C == 0) || (m.A == 0 && m.D == 0)
}

// String returns a string representation.
func (m Matrix) String() string {
	return fmt.Sprintf("Matrix(%g, %g, %g, %g, %g, %g)", m.A, m.B, m.C, m.D, m.E, m.F)
}

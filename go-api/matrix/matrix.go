package matrix

import (
	"math"
)

func Rotate90Clockwise(matrix [][]float64) [][]float64 {
	if len(matrix) == 0 {
		return [][]float64{}
	}
	rows := len(matrix)
	cols := len(matrix[0])

	rotated := make([][]float64, cols)
	for i := range rotated {
		rotated[i] = make([]float64, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			rotated[j][rows-1-i] = matrix[i][j]
		}
	}
	return rotated
}

func FactorizeQR(matrix [][]float64) ([][]float64, [][]float64) {
	if len(matrix) == 0 {
		return [][]float64{}, [][]float64{}
	}

	m := len(matrix)
	n := len(matrix[0])

	q := make([][]float64, m)
	for i := range q {
		q[i] = make([]float64, n)
	}

	r := make([][]float64, n)
	for i := range r {
		r[i] = make([]float64, n)
	}

	for j := 0; j < n; j++ {

		v := make([]float64, m)
		for i := 0; i < m; i++ {
			v[i] = matrix[i][j]
		}

		for i := 0; i < j; i++ {

			dot := 0.0
			for k := 0; k < m; k++ {
				dot += q[k][i] * matrix[k][j]
			}
			r[i][j] = dot

			for k := 0; k < m; k++ {
				v[k] -= r[i][j] * q[k][i]
			}
		}

		norm := 0.0
		for k := 0; k < m; k++ {
			norm += v[k] * v[k]
		}
		norm = math.Sqrt(norm)
		r[j][j] = norm

		if norm > 1e-9 {
			for k := 0; k < m; k++ {
				q[k][j] = v[k] / norm
			}
		}
	}

	return q, r
}

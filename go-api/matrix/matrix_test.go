package matrix

import (
	"reflect"
	"testing"
)

func TestRotate90Clockwise(t *testing.T) {
	t.Run("Debe retornar vacio si la matriz esta vacia", func(t *testing.T) {
		var empty [][]float64
		res := Rotate90Clockwise(empty)
		if len(res) != 0 {
			t.Errorf("Se esperaba matriz vacia, pero se obtuvo %v", res)
		}
	})

	t.Run("Debe rotar la matriz 90 grados en sentido horario correctamente", func(t *testing.T) {
		input := [][]float64{
			{1, 2},
			{3, 4},
		}
		expected := [][]float64{
			{3, 1},
			{4, 2},
		}
		res := Rotate90Clockwise(input)
		if !reflect.DeepEqual(res, expected) {
			t.Errorf("Esperado %v, obtenido %v", expected, res)
		}
	})
}

func TestFactorizeQR(t *testing.T) {
	t.Run("Debe retornar vacio si la matriz esta vacia", func(t *testing.T) {
		var empty [][]float64
		q, r := FactorizeQR(empty)
		if len(q) != 0 || len(r) != 0 {
			t.Errorf("Se esperaban matrices Q y R vacias")
		}
	})

	t.Run("Debe factorizar correctamente una matriz identidad", func(t *testing.T) {
		// La matriz identidad factorizada en QR da como resultado la misma matriz identidad
		input := [][]float64{
			{1, 0},
			{0, 1},
		}
		q, r := FactorizeQR(input)
		if !reflect.DeepEqual(q, input) || !reflect.DeepEqual(r, input) {
			t.Errorf("Fallo en la factorizacion QR. Q: %v, R: %v", q, r)
		}
	})

	t.Run("Debe manejar matrices donde la norma es muy pequena (menor a 1e-9)", func(t *testing.T) {
		// Una matriz llena de ceros obliga al código a saltarse el "if norm > 1e-9"
		input := [][]float64{
			{0, 0},
			{0, 0},
		}
		q, _ := FactorizeQR(input)
		// Verificamos que no modificó Q (se queda en 0 por defecto)
		if q[0][0] != 0 {
			t.Errorf("Se esperaba 0 debido a la norma pequena, se obtuvo %v", q[0][0])
		}
	})
}

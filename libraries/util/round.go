package util

// 根据精度,四舍五入
// 原理是放大 精度**10 倍,小数和0.50000000001比,四舍五入
// WARN 如果放大后变为Inf或NaN,则不处理

import (
	"math"
)

// magnify precision**10, round and return an int64
// return 0 when encounter Inf or NaN
func RoundInt64(val float64, precision int) int64 {
	var t float64
	f := math.Pow10(precision)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return 0
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}

	return int64(t)
}

func Round(val float64, precision int) float64 {
	var t float64
	f := math.Pow10(precision)
	x := val * f
	if math.IsInf(x, 0) || math.IsNaN(x) {
		return val
	}
	if x >= 0.0 {
		t = math.Ceil(x)
		if (t - x) > 0.50000000001 {
			t -= 1.0
		}
	} else {
		t = math.Ceil(-x)
		if (t + x) > 0.50000000001 {
			t -= 1.0
		}
		t = -t
	}
	x = t / f

	if !math.IsInf(x, 0) {
		return x
	}

	return t
}

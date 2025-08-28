package utils

import (
	"math"
	"testing"
)

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		decimals int
		expected float64
	}{
		{"Round 1.235 to 2 decimals", 1.235, 2, 1.24},
		{"Round 1.234 to 2 decimals", 1.234, 2, 1.23},
		{"Round 1.5 to 0 decimals", 1.5, 0, 2},
		{"Round 1.4 to 0 decimals", 1.4, 0, 1},
		{"Round negative decimals", 1.234, -1, 1},
		{"Round zero", 0, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.value, tt.decimals)
			if result != tt.expected {
				t.Errorf("Round(%f, %d) = %f; want %f", tt.value, tt.decimals, result, tt.expected)
			}
		})
	}
}

func TestRoundWithMode(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		decimals int
		mode     RoundingMode
		expected float64
	}{
		{"RoundHalfUp 1.235", 1.235, 2, RoundHalfUp, 1.24},
		{"RoundHalfDown 1.235", 1.235, 2, RoundHalfDown, 1.23},
		{"RoundHalfEven 1.235", 1.235, 2, RoundHalfEven, 1.24},
		{"RoundHalfEven 1.225", 1.225, 2, RoundHalfEven, 1.22},
		{"RoundUp 1.231", 1.231, 2, RoundUp, 1.24},
		{"RoundDown 1.239", 1.239, 2, RoundDown, 1.23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RoundWithMode(tt.value, tt.decimals, tt.mode)
			if result != tt.expected {
				t.Errorf("RoundWithMode(%f, %d, %v) = %f; want %f", tt.value, tt.decimals, tt.mode, result, tt.expected)
			}
		})
	}
}

func TestRoundToCurrency(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{1.235, 1.24},
		{1.234, 1.23},
		{10.999, 11.00},
	}

	for _, tt := range tests {
		result := RoundToCurrency(tt.value)
		if result != tt.expected {
			t.Errorf("RoundToCurrency(%f) = %f; want %f", tt.value, result, tt.expected)
		}
	}
}

func TestRoundToPercent(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{1.23456, 1.2346},
		{1.23454, 1.2345},
	}

	for _, tt := range tests {
		result := RoundToPercent(tt.value)
		if result != tt.expected {
			t.Errorf("RoundToPercent(%f) = %f; want %f", tt.value, result, tt.expected)
		}
	}
}

func TestMinMax(t *testing.T) {
	tests := []struct {
		a, b     float64
		expMin   float64
		expMax   float64
	}{
		{1.5, 2.5, 1.5, 2.5},
		{-1.5, 1.5, -1.5, 1.5},
		{0, 0, 0, 0},
	}

	for _, tt := range tests {
		if result := Min(tt.a, tt.b); result != tt.expMin {
			t.Errorf("Min(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expMin)
		}
		if result := Max(tt.a, tt.b); result != tt.expMax {
			t.Errorf("Max(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expMax)
		}
	}
}

func TestMinMaxInt(t *testing.T) {
	tests := []struct {
		a, b     int
		expMin   int
		expMax   int
	}{
		{1, 2, 1, 2},
		{-1, 1, -1, 1},
		{0, 0, 0, 0},
	}

	for _, tt := range tests {
		if result := MinInt(tt.a, tt.b); result != tt.expMin {
			t.Errorf("MinInt(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expMin)
		}
		if result := MaxInt(tt.a, tt.b); result != tt.expMax {
			t.Errorf("MaxInt(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expMax)
		}
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		value, min, max float64
		expected        float64
	}{
		{5, 1, 10, 5},
		{0, 1, 10, 1},
		{15, 1, 10, 10},
		{-5, -10, -1, -5},
	}

	for _, tt := range tests {
		result := Clamp(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("Clamp(%f, %f, %f) = %f; want %f", tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestClampInt(t *testing.T) {
	tests := []struct {
		value, min, max int
		expected        int
	}{
		{5, 1, 10, 5},
		{0, 1, 10, 1},
		{15, 1, 10, 10},
	}

	for _, tt := range tests {
		result := ClampInt(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("ClampInt(%d, %d, %d) = %d; want %d", tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		value    float64
		expected float64
	}{
		{5.5, 5.5},
		{-5.5, 5.5},
		{0, 0},
	}

	for _, tt := range tests {
		result := Abs(tt.value)
		if result != tt.expected {
			t.Errorf("Abs(%f) = %f; want %f", tt.value, result, tt.expected)
		}
	}
}

func TestAbsInt(t *testing.T) {
	tests := []struct {
		value    int
		expected int
	}{
		{5, 5},
		{-5, 5},
		{0, 0},
	}

	for _, tt := range tests {
		result := AbsInt(tt.value)
		if result != tt.expected {
			t.Errorf("AbsInt(%d) = %d; want %d", tt.value, result, tt.expected)
		}
	}
}

func TestPercentage(t *testing.T) {
	tests := []struct {
		value, total float64
		expected     float64
	}{
		{25, 100, 25},
		{50, 200, 25},
		{0, 100, 0},
		{100, 0, 0}, // Division by zero
	}

	for _, tt := range tests {
		result := Percentage(tt.value, tt.total)
		if result != tt.expected {
			t.Errorf("Percentage(%f, %f) = %f; want %f", tt.value, tt.total, result, tt.expected)
		}
	}
}

func TestPercentageOf(t *testing.T) {
	tests := []struct {
		percentage, total float64
		expected          float64
	}{
		{25, 100, 25},
		{50, 200, 100},
		{0, 100, 0},
	}

	for _, tt := range tests {
		result := PercentageOf(tt.percentage, tt.total)
		if result != tt.expected {
			t.Errorf("PercentageOf(%f, %f) = %f; want %f", tt.percentage, tt.total, result, tt.expected)
		}
	}
}

func TestPercentageChange(t *testing.T) {
	tests := []struct {
		oldValue, newValue float64
		expected           float64
	}{
		{100, 120, 20},
		{100, 80, -20},
		{0, 100, 100},
		{0, 0, 0},
	}

	for _, tt := range tests {
		result := PercentageChange(tt.oldValue, tt.newValue)
		if result != tt.expected {
			t.Errorf("PercentageChange(%f, %f) = %f; want %f", tt.oldValue, tt.newValue, result, tt.expected)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 15},
		{[]float64{-1, 1}, 0},
		{[]float64{}, 0},
		{[]float64{5.5}, 5.5},
	}

	for _, tt := range tests {
		result := Sum(tt.values)
		if result != tt.expected {
			t.Errorf("Sum(%v) = %f; want %f", tt.values, result, tt.expected)
		}
	}
}

func TestSumInt(t *testing.T) {
	tests := []struct {
		values   []int
		expected int
	}{
		{[]int{1, 2, 3, 4, 5}, 15},
		{[]int{-1, 1}, 0},
		{[]int{}, 0},
		{[]int{5}, 5},
	}

	for _, tt := range tests {
		result := SumInt(tt.values)
		if result != tt.expected {
			t.Errorf("SumInt(%v) = %d; want %d", tt.values, result, tt.expected)
		}
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 3},
		{[]float64{-1, 1}, 0},
		{[]float64{}, 0},
		{[]float64{5.5}, 5.5},
	}

	for _, tt := range tests {
		result := Average(tt.values)
		if result != tt.expected {
			t.Errorf("Average(%v) = %f; want %f", tt.values, result, tt.expected)
		}
	}
}

func TestAverageInt(t *testing.T) {
	tests := []struct {
		values   []int
		expected float64
	}{
		{[]int{1, 2, 3, 4, 5}, 3},
		{[]int{-1, 1}, 0},
		{[]int{}, 0},
		{[]int{5}, 5},
	}

	for _, tt := range tests {
		result := AverageInt(tt.values)
		if result != tt.expected {
			t.Errorf("AverageInt(%v) = %f; want %f", tt.values, result, tt.expected)
		}
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 3},
		{[]float64{1, 2, 3, 4}, 2.5},
		{[]float64{5, 1, 3}, 3},
		{[]float64{}, 0},
		{[]float64{5.5}, 5.5},
	}

	for _, tt := range tests {
		result := Median(tt.values)
		if result != tt.expected {
			t.Errorf("Median(%v) = %f; want %f", tt.values, result, tt.expected)
		}
	}
}

func TestStandardDeviation(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		// For [2, 4, 4, 4, 5, 5, 7, 9]: std dev = sqrt(variance) = sqrt(32/7) ≈ 2.138090
		{[]float64{2, 4, 4, 4, 5, 5, 7, 9}, 2.138090},
		{[]float64{1}, 0},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := StandardDeviation(tt.values)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("StandardDeviation(%v) = %f; want %f", tt.values, result, tt.expected)
		}
	}
}

func TestVariance(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		// For [2, 4, 4, 4, 5, 5, 7, 9]: mean = 5, variance = 32/7 ≈ 4.571429
		{[]float64{2, 4, 4, 4, 5, 5, 7, 9}, 4.571429},
		{[]float64{1}, 0},
		{[]float64{}, 0},
	}

	for _, tt := range tests {
		result := Variance(tt.values)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Variance(%v) = %f; want %f", tt.values, result, tt.expected)
		}
	}
}

func TestWeightedAverage(t *testing.T) {
	tests := []struct {
		values   []float64
		weights  []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, []float64{1, 1, 1}, 2},
		{[]float64{1, 2, 3}, []float64{3, 2, 1}, 1.67},
		{[]float64{}, []float64{}, 0},
		{[]float64{1, 2}, []float64{1}, 0}, // Mismatched lengths
		{[]float64{1, 2}, []float64{0, 0}, 0}, // Zero weights
	}

	for _, tt := range tests {
		result := WeightedAverage(tt.values, tt.weights)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("WeightedAverage(%v, %v) = %f; want %f", tt.values, tt.weights, result, tt.expected)
		}
	}
}

func TestLinearInterpolation(t *testing.T) {
	tests := []struct {
		x, x1, y1, x2, y2 float64
		expected          float64
	}{
		{1.5, 1, 1, 2, 3, 2},
		{0, 0, 0, 1, 1, 0},
		{0.5, 0, 0, 1, 2, 1},
		{1, 1, 5, 1, 10, 5}, // Same x values
	}

	for _, tt := range tests {
		result := LinearInterpolation(tt.x, tt.x1, tt.y1, tt.x2, tt.y2)
		if result != tt.expected {
			t.Errorf("LinearInterpolation(%f, %f, %f, %f, %f) = %f; want %f", tt.x, tt.x1, tt.y1, tt.x2, tt.y2, result, tt.expected)
		}
	}
}

func TestCompoundInterest(t *testing.T) {
	tests := []struct {
		principal, rate float64
		periods         int
		expected        float64
	}{
		{1000, 0.05, 1, 1050},
		{1000, 0.05, 2, 1102.5},
		{1000, 0, 5, 1000},
	}

	for _, tt := range tests {
		result := CompoundInterest(tt.principal, tt.rate, tt.periods)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("CompoundInterest(%f, %f, %d) = %f; want %f", tt.principal, tt.rate, tt.periods, result, tt.expected)
		}
	}
}

func TestPresentValue(t *testing.T) {
	tests := []struct {
		futureValue, discountRate float64
		periods                   int
		expected                  float64
	}{
		{1050, 0.05, 1, 1000},
		{1102.5, 0.05, 2, 1000},
	}

	for _, tt := range tests {
		result := PresentValue(tt.futureValue, tt.discountRate, tt.periods)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("PresentValue(%f, %f, %d) = %f; want %f", tt.futureValue, tt.discountRate, tt.periods, result, tt.expected)
		}
	}
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		a, b, tolerance float64
		expected        bool
	}{
		{1.0, 1.0, 0.01, true},
		{1.0, 1.005, 0.01, true},
		{1.0, 1.02, 0.01, false},
	}

	for _, tt := range tests {
		result := IsEqual(tt.a, tt.b, tt.tolerance)
		if result != tt.expected {
			t.Errorf("IsEqual(%f, %f, %f) = %t; want %t", tt.a, tt.b, tt.tolerance, result, tt.expected)
		}
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		value    float64
		expected bool
	}{
		{0.0, true},
		{1e-10, true},
		{0.1, false},
	}

	for _, tt := range tests {
		result := IsZero(tt.value)
		if result != tt.expected {
			t.Errorf("IsZero(%f) = %t; want %t", tt.value, result, tt.expected)
		}
	}
}

func TestSafeDivide(t *testing.T) {
	tests := []struct {
		numerator, denominator float64
		expected               float64
	}{
		{10, 2, 5},
		{10, 0, 0},
		{0, 5, 0},
	}

	for _, tt := range tests {
		result := SafeDivide(tt.numerator, tt.denominator)
		if result != tt.expected {
			t.Errorf("SafeDivide(%f, %f) = %f; want %f", tt.numerator, tt.denominator, result, tt.expected)
		}
	}
}

func TestSafeDivideInt(t *testing.T) {
	tests := []struct {
		numerator, denominator int
		expected               float64
	}{
		{10, 2, 5},
		{10, 0, 0},
		{0, 5, 0},
	}

	for _, tt := range tests {
		result := SafeDivideInt(tt.numerator, tt.denominator)
		if result != tt.expected {
			t.Errorf("SafeDivideInt(%d, %d) = %f; want %f", tt.numerator, tt.denominator, result, tt.expected)
		}
	}
}

func TestGCD(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{12, 8, 4},
		{17, 13, 1},
		{0, 5, 5},
		{-12, 8, 4},
	}

	for _, tt := range tests {
		result := GCD(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("GCD(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestLCM(t *testing.T) {
	tests := []struct {
		a, b     int
		expected int
	}{
		{12, 8, 24},
		{17, 13, 221},
		{0, 5, 0},
		{5, 0, 0},
	}

	for _, tt := range tests {
		result := LCM(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("LCM(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestFactorial(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 1},
		{1, 1},
		{5, 120},
		{-1, 0},
	}

	for _, tt := range tests {
		result := Factorial(tt.n)
		if result != tt.expected {
			t.Errorf("Factorial(%d) = %d; want %d", tt.n, result, tt.expected)
		}
	}
}

func TestFibonacci(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{10, 55},
		{-1, 0},
	}

	for _, tt := range tests {
		result := Fibonacci(tt.n)
		if result != tt.expected {
			t.Errorf("Fibonacci(%d) = %d; want %d", tt.n, result, tt.expected)
		}
	}
}

func TestIsPrime(t *testing.T) {
	tests := []struct {
		n        int
		expected bool
	}{
		{2, true},
		{3, true},
		{4, false},
		{17, true},
		{25, false},
		{1, false},
		{0, false},
		{-5, false},
	}

	for _, tt := range tests {
		result := IsPrime(tt.n)
		if result != tt.expected {
			t.Errorf("IsPrime(%d) = %t; want %t", tt.n, result, tt.expected)
		}
	}
}

func TestRandomFloat(t *testing.T) {
	min, max := 1.0, 10.0
	result := RandomFloat(min, max)
	if result < min || result > max {
		t.Errorf("RandomFloat(%f, %f) = %f; should be between %f and %f", min, max, result, min, max)
	}
}

func TestRandomInt(t *testing.T) {
	min, max := 1, 10
	result := RandomInt(min, max)
	if result < min || result > max {
		t.Errorf("RandomInt(%d, %d) = %d; should be between %d and %d", min, max, result, min, max)
	}
}

func TestRandomIntWithSeed(t *testing.T) {
	min, max := 1, 10
	seed := int64(12345)
	result1 := RandomIntWithSeed(min, max, seed)
	result2 := RandomIntWithSeed(min, max, seed)
	
	if result1 != result2 {
		t.Errorf("RandomIntWithSeed with same seed should produce same result: %d != %d", result1, result2)
	}
	
	if result1 < min || result1 > max {
		t.Errorf("RandomIntWithSeed(%d, %d, %d) = %d; should be between %d and %d", min, max, seed, result1, min, max)
	}
}

func TestNormalizeToRange(t *testing.T) {
	tests := []struct {
		value, oldMin, oldMax, newMin, newMax float64
		expected                              float64
	}{
		{5, 0, 10, 0, 100, 50},
		{0, 0, 10, 0, 100, 0},
		{10, 0, 10, 0, 100, 100},
		{5, 5, 5, 0, 100, 0}, // Same old min/max
	}

	for _, tt := range tests {
		result := NormalizeToRange(tt.value, tt.oldMin, tt.oldMax, tt.newMin, tt.newMax)
		if result != tt.expected {
			t.Errorf("NormalizeToRange(%f, %f, %f, %f, %f) = %f; want %f", tt.value, tt.oldMin, tt.oldMax, tt.newMin, tt.newMax, result, tt.expected)
		}
	}
}

func TestInRange(t *testing.T) {
	tests := []struct {
		value, min, max float64
		expected        bool
	}{
		{5, 1, 10, true},
		{0, 1, 10, false},
		{15, 1, 10, false},
		{1, 1, 10, true},
		{10, 1, 10, true},
	}

	for _, tt := range tests {
		result := InRange(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("InRange(%f, %f, %f) = %t; want %t", tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestInRangeInt(t *testing.T) {
	tests := []struct {
		value, min, max int
		expected        bool
	}{
		{5, 1, 10, true},
		{0, 1, 10, false},
		{15, 1, 10, false},
		{1, 1, 10, true},
		{10, 1, 10, true},
	}

	for _, tt := range tests {
		result := InRangeInt(tt.value, tt.min, tt.max)
		if result != tt.expected {
			t.Errorf("InRangeInt(%d, %d, %d) = %t; want %t", tt.value, tt.min, tt.max, result, tt.expected)
		}
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		expected       float64
	}{
		{0, 0, 3, 4, 5},
		{0, 0, 0, 0, 0},
		{1, 1, 4, 5, 5},
	}

	for _, tt := range tests {
		result := Distance(tt.x1, tt.y1, tt.x2, tt.y2)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Distance(%f, %f, %f, %f) = %f; want %f", tt.x1, tt.y1, tt.x2, tt.y2, result, tt.expected)
		}
	}
}

func TestManhattanDistance(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		expected       float64
	}{
		{0, 0, 3, 4, 7},
		{0, 0, 0, 0, 0},
		{1, 1, 4, 5, 7},
	}

	for _, tt := range tests {
		result := ManhattanDistance(tt.x1, tt.y1, tt.x2, tt.y2)
		if result != tt.expected {
			t.Errorf("ManhattanDistance(%f, %f, %f, %f) = %f; want %f", tt.x1, tt.y1, tt.x2, tt.y2, result, tt.expected)
		}
	}
}

func TestDegreeToRadian(t *testing.T) {
	tests := []struct {
		degrees  float64
		expected float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{360, 2 * math.Pi},
	}

	for _, tt := range tests {
		result := DegreeToRadian(tt.degrees)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("DegreeToRadian(%f) = %f; want %f", tt.degrees, result, tt.expected)
		}
	}
}

func TestRadianToDegree(t *testing.T) {
	tests := []struct {
		radians  float64
		expected float64
	}{
		{0, 0},
		{math.Pi / 2, 90},
		{math.Pi, 180},
		{2 * math.Pi, 360},
	}

	for _, tt := range tests {
		result := RadianToDegree(tt.radians)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("RadianToDegree(%f) = %f; want %f", tt.radians, result, tt.expected)
		}
	}
}

func TestSigmoid(t *testing.T) {
	tests := []struct {
		x        float64
		expected float64
	}{
		{0, 0.5},
		{1, 0.731},
		{-1, 0.269},
	}

	for _, tt := range tests {
		result := Sigmoid(tt.x)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Sigmoid(%f) = %f; want %f", tt.x, result, tt.expected)
		}
	}
}

func TestMovingAverage(t *testing.T) {
	tests := []struct {
		values   []float64
		window   int
		expected []float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 3, []float64{2, 3, 4}},
		{[]float64{1, 2}, 3, []float64{}},
		{[]float64{1, 2, 3}, 0, []float64{}},
	}

	for _, tt := range tests {
		result := MovingAverage(tt.values, tt.window)
		if len(result) != len(tt.expected) {
			t.Errorf("MovingAverage length mismatch: got %d, want %d", len(result), len(tt.expected))
			continue
		}
		for i, v := range result {
			if math.Abs(v-tt.expected[i]) > 0.01 {
				t.Errorf("MovingAverage[%d] = %f; want %f", i, v, tt.expected[i])
			}
		}
	}
}

func TestExponentialMovingAverage(t *testing.T) {
	tests := []struct {
		values   []float64
		alpha    float64
		expected []float64
	}{
		{[]float64{1, 2, 3}, 0.5, []float64{1, 1.5, 2.25}},
		{[]float64{}, 0.5, []float64{}},
		{[]float64{1, 2, 3}, 0, []float64{}}, // Invalid alpha
		{[]float64{1, 2, 3}, 1.5, []float64{}}, // Invalid alpha
	}

	for _, tt := range tests {
		result := ExponentialMovingAverage(tt.values, tt.alpha)
		if len(result) != len(tt.expected) {
			t.Errorf("ExponentialMovingAverage length mismatch: got %d, want %d", len(result), len(tt.expected))
			continue
		}
		for i, v := range result {
			if math.Abs(v-tt.expected[i]) > 0.01 {
				t.Errorf("ExponentialMovingAverage[%d] = %f; want %f", i, v, tt.expected[i])
			}
		}
	}
}

func TestCorrelation(t *testing.T) {
	tests := []struct {
		x, y     []float64
		expected float64
	}{
		{[]float64{1, 2, 3}, []float64{1, 2, 3}, 1},
		{[]float64{1, 2, 3}, []float64{3, 2, 1}, -1},
		{[]float64{}, []float64{}, 0},
		{[]float64{1, 2}, []float64{1}, 0}, // Mismatched lengths
	}

	for _, tt := range tests {
		result := Correlation(tt.x, tt.y)
		if math.Abs(result-tt.expected) > 0.01 {
			t.Errorf("Correlation(%v, %v) = %f; want %f", tt.x, tt.y, result, tt.expected)
		}
	}
}

func TestLinearRegression(t *testing.T) {
	tests := []struct {
		x, y              []float64
		expectedSlope     float64
		expectedIntercept float64
	}{
		{[]float64{1, 2, 3}, []float64{2, 4, 6}, 2, 0},
		{[]float64{}, []float64{}, 0, 0},
		{[]float64{1, 2}, []float64{1}, 0, 0}, // Mismatched lengths
	}

	for _, tt := range tests {
		slope, intercept := LinearRegression(tt.x, tt.y)
		if math.Abs(slope-tt.expectedSlope) > 0.01 {
			t.Errorf("LinearRegression slope = %f; want %f", slope, tt.expectedSlope)
		}
		if math.Abs(intercept-tt.expectedIntercept) > 0.01 {
			t.Errorf("LinearRegression intercept = %f; want %f", intercept, tt.expectedIntercept)
		}
	}
}

// Test for ScaleToRange function (currently 0% coverage)
func TestScaleToRange(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min, max float64
		expected float64
	}{
		{"Value within range", 5.0, 0.0, 10.0, 5.0},
		{"Value below range", -5.0, 0.0, 10.0, 0.0},
		{"Value above range", 15.0, 0.0, 10.0, 10.0},
		{"Value at min", 0.0, 0.0, 10.0, 0.0},
		{"Value at max", 10.0, 0.0, 10.0, 10.0},
		{"Negative range", -5.0, -10.0, -1.0, -5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ScaleToRange(tt.value, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("ScaleToRange(%f, %f, %f) = %f; want %f", tt.value, tt.min, tt.max, result, tt.expected)
			}
		})
	}
}

// Test for Logistic function (currently 0% coverage)
func TestLogistic(t *testing.T) {
	tests := []struct {
		name     string
		x, k, x0, l float64
		expected float64
	}{
		{"Basic logistic", 5.0, 1.0, 0.0, 10.0, 9.933},
		{"At inflection point", 0.0, 1.0, 0.0, 10.0, 5.0},
		{"Negative x", -5.0, 1.0, 0.0, 10.0, 0.067},
		{"Different carrying capacity", 2.0, 0.5, 1.0, 100.0, 62.246},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Logistic(tt.x, tt.k, tt.x0, tt.l)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("Logistic(%f, %f, %f, %f) = %f; want %f", tt.x, tt.k, tt.x0, tt.l, result, tt.expected)
			}
		})
	}
}

// Test for ExponentialDecay function (currently 0% coverage)
func TestExponentialDecay(t *testing.T) {
	tests := []struct {
		name     string
		initial, rate, time float64
		expected float64
	}{
		{"Basic decay", 100.0, 0.1, 1.0, 90.484},
		{"No time", 100.0, 0.1, 0.0, 100.0},
		{"High decay rate", 1000.0, 0.5, 2.0, 367.879},
		{"Low decay rate", 100.0, 0.01, 10.0, 90.484},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExponentialDecay(tt.initial, tt.rate, tt.time)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("ExponentialDecay(%f, %f, %f) = %f; want %f", tt.initial, tt.rate, tt.time, result, tt.expected)
			}
		})
	}
}

// Test for ExponentialGrowth function (currently 0% coverage)
func TestExponentialGrowth(t *testing.T) {
	tests := []struct {
		name     string
		initial, rate, time float64
		expected float64
	}{
		{"Basic growth", 100.0, 0.1, 1.0, 110.517},
		{"No time", 100.0, 0.1, 0.0, 100.0},
		{"High growth rate", 100.0, 0.5, 2.0, 271.828},
		{"Low growth rate", 1000.0, 0.01, 10.0, 1105.171},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExponentialGrowth(tt.initial, tt.rate, tt.time)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("ExponentialGrowth(%f, %f, %f) = %f; want %f", tt.initial, tt.rate, tt.time, result, tt.expected)
			}
		})
	}
}

// Additional test for MaxInt to improve coverage from 66.7%
func TestMaxIntEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"First larger", 10, 5, 10},
		{"Second larger", 3, 8, 8},
		{"Equal values", 7, 7, 7},
		{"Negative values", -5, -10, -5},
		{"Mixed signs", -3, 2, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxInt(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("MaxInt(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
// Package utils provides mathematical utility functions for the ecommerce engine.
// It includes comprehensive mathematical operations such as rounding, statistical calculations,
// financial computations, geometric functions, and various mathematical utilities commonly
// needed in ecommerce applications like pricing calculations, tax computations, and data analysis.
//
// The package is designed to handle floating-point precision issues and provides safe
// mathematical operations with proper error handling for edge cases like division by zero.
//
// Example usage:
//
//	price := 19.999
//	roundedPrice := utils.RoundToCurrency(price) // 20.00
//
//	values := []float64{10.5, 15.2, 8.7, 12.1}
//	avg := utils.Average(values) // 11.625
//
//	discount := utils.Percentage(5.0, 25.0) // 20.0%
package utils

import (
	"math"
	"math/rand"
	"time"
)

// RoundingMode represents different rounding modes for mathematical operations.
// These modes determine how values are rounded when they fall exactly between
// two possible rounded values (e.g., 1.5 can round to either 1 or 2).
type RoundingMode int

const (
	// RoundHalfUp rounds 0.5 up to the next integer (default behavior).
	// Example: 1.5 → 2, -1.5 → -1
	RoundHalfUp RoundingMode = iota
	
	// RoundHalfDown rounds 0.5 down to the previous integer.
	// Example: 1.5 → 1, -1.5 → -2
	RoundHalfDown
	
	// RoundHalfEven rounds 0.5 to the nearest even number (banker's rounding).
	// This method reduces bias in repeated rounding operations.
	// Example: 1.5 → 2, 2.5 → 2, 3.5 → 4
	RoundHalfEven
	
	// RoundUp always rounds up to the next integer (ceiling function).
	// Example: 1.1 → 2, 1.9 → 2, -1.1 → -1
	RoundUp
	
	// RoundDown always rounds down to the previous integer (floor function).
	// Example: 1.1 → 1, 1.9 → 1, -1.1 → -2
	RoundDown
)

// Round rounds a float64 value to the specified number of decimal places using
// the default RoundHalfUp mode.
//
// Parameters:
//   - value: The floating-point number to round
//   - decimals: Number of decimal places to round to (negative values are treated as 0)
//
// Returns:
//   - The rounded value as float64
//
// Example:
//	result := Round(3.14159, 2) // 3.14
//	result := Round(1.235, 2)   // 1.24
//	result := Round(1.5, 0)     // 2
func Round(value float64, decimals int) float64 {
	return RoundWithMode(value, decimals, RoundHalfUp)
}

// RoundWithMode rounds a float64 value using the specified rounding mode.
// This function provides fine-grained control over rounding behavior, which is
// particularly important for financial calculations where rounding consistency matters.
//
// Parameters:
//   - value: The floating-point number to round
//   - decimals: Number of decimal places to round to (negative values are treated as 0)
//   - mode: The rounding mode to use (RoundHalfUp, RoundHalfDown, etc.)
//
// Returns:
//   - The rounded value as float64
//
// Example:
//	result := RoundWithMode(1.235, 2, RoundHalfUp)   // 1.24
//	result := RoundWithMode(1.235, 2, RoundHalfDown) // 1.23
//	result := RoundWithMode(1.235, 2, RoundHalfEven) // 1.24
func RoundWithMode(value float64, decimals int, mode RoundingMode) float64 {
	if decimals < 0 {
		decimals = 0
	}

	multiplier := math.Pow(10, float64(decimals))
	scaled := value * multiplier

	switch mode {
	case RoundHalfUp:
		return math.Floor(scaled+0.5) / multiplier
	case RoundHalfDown:
		// For RoundHalfDown, we round 0.5 down
		fracPart := scaled - math.Floor(scaled)
		epsilon := 1e-10
		if math.Abs(fracPart-0.5) < epsilon {
			return math.Floor(scaled) / multiplier
		}
		return math.Round(scaled) / multiplier
	case RoundHalfEven:
		return roundHalfEven(scaled) / multiplier
	case RoundUp:
		return math.Ceil(scaled) / multiplier
	case RoundDown:
		return math.Floor(scaled) / multiplier
	default:
		return math.Floor(scaled+0.5) / multiplier
	}
}

// roundHalfEven implements banker's rounding (round half to even).
// This is an internal helper function that implements the IEEE 754 standard
// for rounding, which helps reduce bias in statistical calculations.
//
// The function handles negative values by recursively calling itself with
// the absolute value and then negating the result.
//
// Parameters:
//   - value: The floating-point number to round
//
// Returns:
//   - The rounded value using banker's rounding
func roundHalfEven(value float64) float64 {
	if value < 0 {
		return -roundHalfEven(-value)
	}

	intPart := math.Floor(value)
	fracPart := value - intPart

	// Use a small epsilon to handle floating point precision issues
	epsilon := 1e-10
	
	if fracPart < 0.5-epsilon {
		return intPart
	} else if fracPart > 0.5+epsilon {
		return intPart + 1
	} else {
		// fracPart is approximately 0.5, round to even
		// Check if the integer part is even
		if int(intPart)%2 == 0 {
			return intPart // Even, round down
		} else {
			return intPart + 1 // Odd, round up
		}
	}
}

// RoundToCurrency rounds a value to standard currency precision (2 decimal places).
// This is a convenience function commonly used in ecommerce applications for
// price calculations, tax computations, and financial transactions.
//
// Parameters:
//   - value: The monetary value to round
//
// Returns:
//   - The value rounded to 2 decimal places
//
// Example:
//	price := 19.999
//	rounded := RoundToCurrency(price) // 20.00
//	tax := RoundToCurrency(15.678)    // 15.68
func RoundToCurrency(value float64) float64 {
	return Round(value, 2)
}

// RoundToPercent rounds a value to percentage precision (4 decimal places).
// This function is useful for precise percentage calculations in discount
// computations, tax rates, and statistical analysis.
//
// Parameters:
//   - value: The percentage value to round
//
// Returns:
//   - The value rounded to 4 decimal places
//
// Example:
//	rate := 15.678912
//	rounded := RoundToPercent(rate) // 15.6789
//	discount := RoundToPercent(7.12345) // 7.1235
func RoundToPercent(value float64) float64 {
	return Round(value, 4)
}

// Min returns the minimum of two float64 values.
// This function is useful for finding the smaller value in comparisons,
// such as determining the lowest price or minimum quantity.
//
// Parameters:
//   - a: First floating-point value
//   - b: Second floating-point value
//
// Returns:
//   - The smaller of the two values
//
// Example:
//	minPrice := Min(19.99, 24.99) // 19.99
//	minValue := Min(-5.5, 3.2)    // -5.5
func Min(a, b float64) float64 {
	return math.Min(a, b)
}

// Max returns the maximum of two float64 values.
// This function is useful for finding the larger value in comparisons,
// such as determining the highest price or maximum quantity.
//
// Parameters:
//   - a: First floating-point value
//   - b: Second floating-point value
//
// Returns:
//   - The larger of the two values
//
// Example:
//	maxPrice := Max(19.99, 24.99) // 24.99
//	maxValue := Max(-5.5, 3.2)    // 3.2
func Max(a, b float64) float64 {
	return math.Max(a, b)
}

// MinInt returns the minimum of two integer values.
// This function provides integer-specific minimum comparison without
// the overhead of floating-point operations.
//
// Parameters:
//   - a: First integer value
//   - b: Second integer value
//
// Returns:
//   - The smaller of the two integers
//
// Example:
//	minQty := MinInt(5, 10)   // 5
//	minVal := MinInt(-3, 7)   // -3
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt returns the maximum of two integer values.
// This function provides integer-specific maximum comparison without
// the overhead of floating-point operations.
//
// Parameters:
//   - a: First integer value
//   - b: Second integer value
//
// Returns:
//   - The larger of the two integers
//
// Example:
//	maxQty := MaxInt(5, 10)   // 10
//	maxVal := MaxInt(-3, 7)   // 7
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp constrains a floating-point value between specified minimum and maximum bounds.
// This function is essential for ensuring values stay within acceptable ranges,
// such as limiting discount percentages or constraining price adjustments.
//
// Parameters:
//   - value: The value to constrain
//   - min: The minimum allowed value
//   - max: The maximum allowed value
//
// Returns:
//   - The value if it's within bounds, otherwise the nearest boundary
//
// Example:
//	discount := Clamp(150.0, 0.0, 100.0) // 100.0 (clamped to max)
//	price := Clamp(15.99, 10.0, 50.0)    // 15.99 (within bounds)
//	negative := Clamp(-5.0, 0.0, 100.0)  // 0.0 (clamped to min)
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt constrains an integer value between specified minimum and maximum bounds.
// This function provides integer-specific clamping without floating-point overhead,
// useful for quantities, counts, and other discrete values.
//
// Parameters:
//   - value: The integer value to constrain
//   - min: The minimum allowed value
//   - max: The maximum allowed value
//
// Returns:
//   - The value if it's within bounds, otherwise the nearest boundary
//
// Example:
//	quantity := ClampInt(150, 1, 100)  // 100 (clamped to max)
//	count := ClampInt(25, 10, 50)      // 25 (within bounds)
//	negQty := ClampInt(-5, 0, 100)     // 0 (clamped to min)
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Abs returns the absolute value of a floating-point number.
// This function removes the sign from a number, making negative values positive
// while leaving positive values unchanged. Useful for distance calculations,
// error measurements, and ensuring non-negative results.
//
// Parameters:
//   - value: The floating-point number to get absolute value of
//
// Returns:
//   - The absolute value (always non-negative)
//
// Example:
//	diff := Abs(-15.5)     // 15.5
//	distance := Abs(10.0)  // 10.0
//	error := Abs(-0.001)   // 0.001
func Abs(value float64) float64 {
	return math.Abs(value)
}

// AbsInt returns the absolute value of an integer.
// This function provides integer-specific absolute value calculation
// without floating-point overhead, useful for counts and discrete values.
//
// Parameters:
//   - value: The integer to get absolute value of
//
// Returns:
//   - The absolute value (always non-negative)
//
// Example:
//	diff := AbsInt(-25)    // 25
//	count := AbsInt(10)    // 10
//	offset := AbsInt(-1)   // 1
func AbsInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// Percentage calculates what percentage one value represents of a total.
// This function is essential for ecommerce applications to calculate discounts,
// tax rates, commission percentages, and statistical analysis.
//
// Parameters:
//   - value: The partial value
//   - total: The total value (if zero, returns 0 to avoid division by zero)
//
// Returns:
//   - The percentage as a float64 (e.g., 25.0 for 25%)
//
// Example:
//	discountPct := Percentage(5.0, 25.0)    // 20.0 (5 is 20% of 25)
//	taxRate := Percentage(2.5, 50.0)        // 5.0 (2.5 is 5% of 50)
//	commission := Percentage(150, 1000)     // 15.0 (150 is 15% of 1000)
func Percentage(value, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

// PercentageOf calculates the actual value that represents a given percentage of a total.
// This function is the inverse of Percentage and is commonly used to calculate
// discount amounts, tax amounts, and commission values.
//
// Parameters:
//   - percentage: The percentage value (e.g., 15.0 for 15%)
//   - total: The total value to calculate percentage of
//
// Returns:
//   - The calculated value representing the percentage of total
//
// Example:
//	discountAmount := PercentageOf(20.0, 100.0)  // 20.0 (20% of 100)
//	taxAmount := PercentageOf(8.5, 250.0)       // 21.25 (8.5% of 250)
//	commission := PercentageOf(5.0, 1000.0)     // 50.0 (5% of 1000)
func PercentageOf(percentage, total float64) float64 {
	return (percentage / 100) * total
}

// PercentageChange calculates the percentage change between two values.
// This function is useful for analyzing price changes, growth rates,
// performance metrics, and trend analysis in ecommerce applications.
//
// Parameters:
//   - oldValue: The original value
//   - newValue: The new value to compare against
//
// Returns:
//   - The percentage change (positive for increase, negative for decrease)
//   - Returns 100 if oldValue is 0 and newValue is non-zero
//   - Returns 0 if both values are 0
//
// Example:
//	priceChange := PercentageChange(100.0, 120.0)  // 20.0 (20% increase)
//	salesChange := PercentageChange(200.0, 150.0)  // -25.0 (25% decrease)
//	growth := PercentageChange(0.0, 50.0)          // 100.0 (from zero)
func PercentageChange(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		if newValue == 0 {
			return 0
		}
		return 100 // or could return infinity
	}
	return ((newValue - oldValue) / oldValue) * 100
}

// Sum calculates the sum of a slice of float64 values
func Sum(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum
}

// SumInt calculates the sum of a slice of int values
func SumInt(values []int) int {
	sum := 0
	for _, value := range values {
		sum += value
	}
	return sum
}

// Average calculates the average of a slice of float64 values
func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Sum(values) / float64(len(values))
}

// AverageInt calculates the average of a slice of int values
func AverageInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(SumInt(values)) / float64(len(values))
}

// Median calculates the median of a slice of float64 values
func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	// Create a copy and sort it
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sortFloat64Slice(sorted)

	n := len(sorted)
	if n%2 == 0 {
		// Even number of elements
		return (sorted[n/2-1] + sorted[n/2]) / 2
	} else {
		// Odd number of elements
		return sorted[n/2]
	}
}

// sortFloat64Slice sorts a slice of float64 values in ascending order
func sortFloat64Slice(values []float64) {
	for i := 0; i < len(values)-1; i++ {
		for j := 0; j < len(values)-i-1; j++ {
			if values[j] > values[j+1] {
				values[j], values[j+1] = values[j+1], values[j]
			}
		}
	}
}

// StandardDeviation calculates the standard deviation of a slice of float64 values
func StandardDeviation(values []float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	avg := Average(values)
	variance := 0.0

	for _, value := range values {
		variance += math.Pow(value-avg, 2)
	}

	variance /= float64(len(values) - 1) // Sample standard deviation
	return math.Sqrt(variance)
}

// Variance calculates the variance of a slice of float64 values
func Variance(values []float64) float64 {
	if len(values) <= 1 {
		return 0
	}

	avg := Average(values)
	variance := 0.0

	for _, value := range values {
		variance += math.Pow(value-avg, 2)
	}

	return variance / float64(len(values)-1)
}

// WeightedAverage calculates weighted average of values
func WeightedAverage(values, weights []float64) float64 {
	if len(values) != len(weights) || len(values) == 0 {
		return 0
	}

	weightedSum := 0.0
	totalWeight := 0.0

	for i, value := range values {
		weightedSum += value * weights[i]
		totalWeight += weights[i]
	}

	if totalWeight == 0 {
		return 0
	}

	return weightedSum / totalWeight
}

// LinearInterpolation performs linear interpolation between two points
func LinearInterpolation(x, x1, y1, x2, y2 float64) float64 {
	if x2 == x1 {
		return y1
	}
	return y1 + (y2-y1)*(x-x1)/(x2-x1)
}

// CompoundInterest calculates compound interest
func CompoundInterest(principal, rate float64, periods int) float64 {
	return principal * math.Pow(1+rate, float64(periods))
}

// PresentValue calculates present value given future value and discount rate
func PresentValue(futureValue, discountRate float64, periods int) float64 {
	return futureValue / math.Pow(1+discountRate, float64(periods))
}

// IsEqual checks if two float64 values are equal within a tolerance
func IsEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// IsZero checks if a float64 value is effectively zero
func IsZero(value float64) bool {
	return IsEqual(value, 0, 1e-9)
}

// SafeDivide performs division with zero check
func SafeDivide(numerator, denominator float64) float64 {
	if IsZero(denominator) {
		return 0
	}
	return numerator / denominator
}

// SafeDivideInt performs integer division with zero check
func SafeDivideInt(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

// GCD calculates the greatest common divisor of two integers
func GCD(a, b int) int {
	a = AbsInt(a)
	b = AbsInt(b)

	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM calculates the least common multiple of two integers
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return AbsInt(a*b) / GCD(a, b)
}

// Factorial calculates the factorial of a non-negative integer
func Factorial(n int) int {
	if n < 0 {
		return 0
	}
	if n <= 1 {
		return 1
	}

	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}
	return result
}

// Fibonacci calculates the nth Fibonacci number
func Fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}

// IsPrime checks if a number is prime
func IsPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	i := 5
	for i*i <= n {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
		i += 6
	}
	return true
}

// RandomFloat generates a random float64 between min and max
func RandomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// RandomInt generates a random int between min and max (inclusive)
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// RandomIntWithSeed generates a random int with a specific seed
func RandomIntWithSeed(min, max int, seed int64) int {
	source := rand.NewSource(seed)
	r := rand.New(source)
	return r.Intn(max-min+1) + min
}

// NormalizeToRange normalizes a value from one range to another
func NormalizeToRange(value, oldMin, oldMax, newMin, newMax float64) float64 {
	if oldMax == oldMin {
		return newMin
	}
	return newMin + (value-oldMin)*(newMax-newMin)/(oldMax-oldMin)
}

// ScaleToRange scales a value to fit within a specific range.
// This function is a convenience wrapper around Clamp that ensures
// a value stays within specified bounds, useful for constraining
// values like prices, ratings, or percentages in ecommerce applications.
//
// Parameters:
//   - value: The value to scale/constrain
//   - min: The minimum allowed value
//   - max: The maximum allowed value
//
// Returns:
//   - The value if it's within bounds, otherwise the nearest boundary
//
// Example:
//	discount := ScaleToRange(150.0, 0.0, 100.0) // 100.0 (clamped to max)
//	rating := ScaleToRange(4.5, 1.0, 5.0)       // 4.5 (within bounds)
//	price := ScaleToRange(-10.0, 0.0, 1000.0)   // 0.0 (clamped to min)
func ScaleToRange(value, min, max float64) float64 {
	return Clamp(value, min, max)
}

// InRange checks if a value is within a specified range (inclusive)
func InRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// InRangeInt checks if an int value is within a specified range (inclusive)
func InRangeInt(value, min, max int) bool {
	return value >= min && value <= max
}

// Distance calculates the Euclidean distance between two 2D points
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// ManhattanDistance calculates the Manhattan distance between two 2D points
func ManhattanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Abs(x2-x1) + math.Abs(y2-y1)
}

// DegreeToRadian converts degrees to radians
func DegreeToRadian(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadianToDegree converts radians to degrees
func RadianToDegree(radians float64) float64 {
	return radians * 180 / math.Pi
}

// Sigmoid calculates the sigmoid activation function.
// This function maps any real number to a value between 0 and 1,
// creating an S-shaped curve. Useful for machine learning applications,
// probability calculations, and smooth transitions in ecommerce analytics.
//
// Parameters:
//   - x: The input value
//
// Returns:
//   - A value between 0 and 1 following the sigmoid curve
//
// Example:
//	// Convert score to probability-like value
//	prob := Sigmoid(2.0)     // ~0.88
//	// Smooth transition function
//	weight := Sigmoid(-1.0)  // ~0.27
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// Logistic calculates the logistic function with custom parameters.
// This function provides a generalized logistic curve that can be customized
// for specific growth patterns, capacity limits, and inflection points.
// Useful for modeling customer adoption, sales growth, and market saturation.
//
// Parameters:
//   - x: The input variable (often time)
//   - k: The growth rate
//   - x0: The x-value of the inflection point
//   - l: The maximum value (carrying capacity)
//
// Returns:
//   - The logistic function value
//
// Example:
//	// Model customer adoption over time
//	adoption := Logistic(6.0, 0.5, 5.0, 1000.0)  // customers at month 6
//	// Model sales growth with market cap
//	sales := Logistic(12.0, 0.3, 10.0, 50000.0)  // sales at month 12
func Logistic(x, k, x0, l float64) float64 {
	return l / (1 + math.Exp(-k*(x-x0)))
}

// ExponentialDecay calculates exponential decay over time.
// This function models how values decrease exponentially, useful for
// calculating depreciation, customer churn rates, inventory spoilage,
// and time-based discounts in ecommerce applications.
//
// Parameters:
//   - initial: The initial value
//   - rate: The decay rate (positive value)
//   - time: The time elapsed
//
// Returns:
//   - The decayed value after the specified time
//
// Example:
//	// Calculate product value after depreciation
//	value := ExponentialDecay(1000.0, 0.1, 2.0)  // ~818.73 after 2 years
//	// Model customer retention
//	retained := ExponentialDecay(1000.0, 0.05, 12.0)  // customers after 12 months
func ExponentialDecay(initial, rate, time float64) float64 {
	return initial * math.Exp(-rate*time)
}

// ExponentialGrowth calculates exponential growth over time.
// This function models how values increase exponentially, useful for
// calculating compound interest, viral growth, customer acquisition,
// and revenue projections in ecommerce applications.
//
// Parameters:
//   - initial: The initial value
//   - rate: The growth rate (positive value)
//   - time: The time elapsed
//
// Returns:
//   - The grown value after the specified time
//
// Example:
//	// Calculate investment growth
//	value := ExponentialGrowth(1000.0, 0.08, 5.0)  // ~1491.82 after 5 years
//	// Model user base growth
//	users := ExponentialGrowth(100.0, 0.15, 3.0)   // users after 3 periods
func ExponentialGrowth(initial, rate, time float64) float64 {
	return initial * math.Exp(rate*time)
}

// MovingAverage calculates simple moving average
func MovingAverage(values []float64, window int) []float64 {
	if len(values) < window || window <= 0 {
		return []float64{}
	}

	result := make([]float64, len(values)-window+1)
	for i := 0; i <= len(values)-window; i++ {
		sum := 0.0
		for j := i; j < i+window; j++ {
			sum += values[j]
		}
		result[i] = sum / float64(window)
	}
	return result
}

// ExponentialMovingAverage calculates exponential moving average
func ExponentialMovingAverage(values []float64, alpha float64) []float64 {
	if len(values) == 0 || alpha <= 0 || alpha > 1 {
		return []float64{}
	}

	result := make([]float64, len(values))
	result[0] = values[0]

	for i := 1; i < len(values); i++ {
		result[i] = alpha*values[i] + (1-alpha)*result[i-1]
	}
	return result
}

// Correlation calculates Pearson correlation coefficient
func Correlation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0
	}

	n := float64(len(x))
	sumX := Sum(x)
	sumY := Sum(y)
	sumXY := 0.0
	sumX2 := 0.0
	sumY2 := 0.0

	for i := 0; i < len(x); i++ {
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
		sumY2 += y[i] * y[i]
	}

	numerator := n*sumXY - sumX*sumY
	denominator := math.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))

	if IsZero(denominator) {
		return 0
	}

	return numerator / denominator
}

// LinearRegression calculates linear regression coefficients (slope, intercept)
func LinearRegression(x, y []float64) (slope, intercept float64) {
	if len(x) != len(y) || len(x) == 0 {
		return 0, 0
	}

	n := float64(len(x))
	sumX := Sum(x)
	sumY := Sum(y)
	sumXY := 0.0
	sumX2 := 0.0

	for i := 0; i < len(x); i++ {
		sumXY += x[i] * y[i]
		sumX2 += x[i] * x[i]
	}

	denominator := n*sumX2 - sumX*sumX
	if IsZero(denominator) {
		return 0, Average(y)
	}

	slope = (n*sumXY - sumX*sumY) / denominator
	intercept = (sumY - slope*sumX) / n

	return slope, intercept
}
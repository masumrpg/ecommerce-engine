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

// Sum calculates the sum of a slice of float64 values.
// This function is essential for aggregating numerical data such as
// order totals, tax amounts, shipping costs, and statistical calculations.
//
// Parameters:
//   - values: Slice of floating-point values to sum
//
// Returns:
//   - The sum of all values in the slice (0.0 for empty slice)
//
// Example:
//	prices := []float64{19.99, 24.99, 15.50}
//	total := Sum(prices) // 60.48
//	taxes := []float64{1.60, 2.00, 1.24}
//	totalTax := Sum(taxes) // 4.84
func Sum(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum
}

// SumInt calculates the sum of a slice of integer values.
// This function is useful for aggregating discrete quantities such as
// item counts, inventory levels, order quantities, and other countable metrics.
//
// Parameters:
//   - values: Slice of integer values to sum
//
// Returns:
//   - The sum of all values in the slice (0 for empty slice)
//
// Example:
//	quantities := []int{5, 10, 3, 7}
//	totalQty := SumInt(quantities) // 25
//	counts := []int{100, 250, 75}
//	totalCount := SumInt(counts) // 425
func SumInt(values []int) int {
	sum := 0
	for _, value := range values {
		sum += value
	}
	return sum
}

// Average calculates the arithmetic mean of a slice of float64 values.
// This function is fundamental for statistical analysis, performance metrics,
// price averaging, and data analysis in ecommerce applications.
//
// Parameters:
//   - values: Slice of floating-point values to average
//
// Returns:
//   - The arithmetic mean of all values (0.0 for empty slice)
//
// Example:
//	prices := []float64{10.0, 20.0, 30.0}
//	avgPrice := Average(prices) // 20.0
//	ratings := []float64{4.5, 3.8, 4.2, 4.9}
//	avgRating := Average(ratings) // 4.35
func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	return Sum(values) / float64(len(values))
}

// AverageInt calculates the arithmetic mean of a slice of integer values.
// This function converts integer values to floating-point for precise averaging,
// useful for analyzing discrete metrics like quantities, counts, and ratings.
//
// Parameters:
//   - values: Slice of integer values to average
//
// Returns:
//   - The arithmetic mean as a float64 (0.0 for empty slice)
//
// Example:
//	quantities := []int{5, 10, 15}
//	avgQty := AverageInt(quantities) // 10.0
//	scores := []int{85, 92, 78, 88}
//	avgScore := AverageInt(scores) // 85.75
func AverageInt(values []int) float64 {
	if len(values) == 0 {
		return 0
	}
	return float64(SumInt(values)) / float64(len(values))
}

// Median calculates the median (middle value) of a slice of float64 values.
// The median is the value that separates the higher half from the lower half
// of a data set. It's less affected by outliers than the mean, making it
// useful for price analysis, performance metrics, and statistical reporting.
//
// Parameters:
//   - values: Slice of floating-point values to find median of
//
// Returns:
//   - The median value (0.0 for empty slice)
//   - For even number of elements, returns average of two middle values
//   - For odd number of elements, returns the exact middle value
//
// Example:
//	prices := []float64{10.0, 15.0, 20.0, 25.0, 30.0}
//	medianPrice := Median(prices) // 20.0
//	scores := []float64{85.5, 92.0, 78.5, 88.0}
//	medianScore := Median(scores) // 86.75 (average of 85.5 and 88.0)
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

// StandardDeviation calculates the sample standard deviation of a slice of float64 values.
// Standard deviation measures the amount of variation or dispersion in a dataset.
// It's useful for analyzing price volatility, performance consistency, quality metrics,
// and risk assessment in ecommerce applications.
//
// Parameters:
//   - values: Slice of floating-point values to calculate standard deviation for
//
// Returns:
//   - The sample standard deviation (0.0 for empty slice or single value)
//
// Example:
//	prices := []float64{10.0, 12.0, 14.0, 16.0, 18.0}
//	stdDev := StandardDeviation(prices) // ~3.16
//	responseTime := []float64{120.5, 135.2, 118.9, 142.1}
//	variability := StandardDeviation(responseTime) // measures consistency
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

// Variance calculates the sample variance of a slice of float64 values.
// Variance measures how far the values are spread from the mean.
// It's the square of standard deviation and is useful for statistical analysis,
// risk assessment, and quality control in ecommerce applications.
//
// Parameters:
//   - values: Slice of floating-point values to calculate variance for
//
// Returns:
//   - The sample variance (0.0 for empty slice or single value)
//
// Example:
//	prices := []float64{10.0, 12.0, 14.0, 16.0, 18.0}
//	variance := Variance(prices) // ~10.0
//	salesData := []float64{1000.0, 1200.0, 950.0, 1100.0}
//	salesVariance := Variance(salesData) // measures sales consistency
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

// WeightedAverage calculates the weighted average of values using corresponding weights.
// This function gives more importance to certain values based on their weights,
// useful for calculating weighted prices, priority-based scoring, and importance-adjusted metrics.
//
// Parameters:
//   - values: Slice of floating-point values to average
//   - weights: Slice of weights corresponding to each value (must be same length as values)
//
// Returns:
//   - The weighted average (0.0 for empty slices, mismatched lengths, or zero total weight)
//
// Example:
//	prices := []float64{10.0, 20.0, 30.0}
//	quantities := []float64{5.0, 2.0, 1.0}
//	avgPrice := WeightedAverage(prices, quantities) // ~13.75
//	scores := []float64{85.0, 92.0, 78.0}
//	importance := []float64{0.5, 0.3, 0.2}
//	weightedScore := WeightedAverage(scores, importance) // ~85.9
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

// LinearInterpolation performs linear interpolation between two points.
// This function calculates a value along a straight line between two known points,
// useful for price scaling, progressive discounts, shipping rate calculations,
// and smooth transitions in ecommerce applications.
//
// Parameters:
//   - x: The x-coordinate where you want to find the y-value
//   - x1, y1: Coordinates of the first known point
//   - x2, y2: Coordinates of the second known point
//
// Returns:
//   - The interpolated y-value at position x
//   - Returns y1 if x1 equals x2 (vertical line)
//
// Example:
//	// Calculate shipping cost between weight ranges
//	cost := LinearInterpolation(7.5, 5.0, 10.0, 10.0, 15.0) // 12.5
//	// Progressive discount based on quantity
//	discount := LinearInterpolation(15.0, 10.0, 5.0, 20.0, 15.0) // 10.0%
func LinearInterpolation(x, x1, y1, x2, y2 float64) float64 {
	if x2 == x1 {
		return y1
	}
	return y1 + (y2-y1)*(x-x1)/(x2-x1)
}

// CompoundInterest calculates the future value using compound interest formula.
// This function computes how an investment grows over time with compound interest,
// useful for financial projections, loyalty point calculations, subscription growth,
// and investment analysis in ecommerce applications.
//
// Parameters:
//   - principal: The initial amount (starting value)
//   - rate: The interest rate per period (as decimal, e.g., 0.05 for 5%)
//   - periods: The number of compounding periods
//
// Returns:
//   - The future value after compound interest
//
// Example:
//	// Calculate investment growth
//	futureValue := CompoundInterest(1000.0, 0.08, 5) // ~1469.33 after 5 years at 8%
//	// Loyalty points growth
//	points := CompoundInterest(100.0, 0.02, 12) // points after 12 months at 2%
func CompoundInterest(principal, rate float64, periods int) float64 {
	return principal * math.Pow(1+rate, float64(periods))
}

// PresentValue calculates the present value of a future amount using discount rate.
// This function determines what a future amount is worth in today's terms,
// useful for financial analysis, investment evaluation, subscription pricing,
// and time-value calculations in ecommerce applications.
//
// Parameters:
//   - futureValue: The amount expected in the future
//   - discountRate: The discount rate per period (as decimal, e.g., 0.05 for 5%)
//   - periods: The number of periods until the future value is received
//
// Returns:
//   - The present value of the future amount
//
// Example:
//	// Calculate present value of future payment
//	presentVal := PresentValue(1000.0, 0.08, 3) // ~793.83 (today's value)
//	// Subscription value analysis
//	currentWorth := PresentValue(1200.0, 0.05, 2) // current worth of future revenue
func PresentValue(futureValue, discountRate float64, periods int) float64 {
	return futureValue / math.Pow(1+discountRate, float64(periods))
}

// IsEqual checks if two float64 values are equal within a specified tolerance.
// This function handles floating-point precision issues by comparing values
// within an acceptable margin of error, essential for reliable financial
// calculations and price comparisons in ecommerce applications.
//
// Parameters:
//   - a: First floating-point value to compare
//   - b: Second floating-point value to compare
//   - tolerance: Maximum acceptable difference between the values
//
// Returns:
//   - true if the absolute difference is within tolerance, false otherwise
//
// Example:
//	// Compare calculated prices with precision tolerance
//	equal := IsEqual(19.999, 20.0, 0.01) // true (within 1 cent)
//	// Validate tax calculations
//	valid := IsEqual(calculatedTax, expectedTax, 0.001) // precise comparison
func IsEqual(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}

// IsZero checks if a float64 value is effectively zero within a small tolerance.
// This function handles floating-point precision issues when checking for zero values,
// essential for reliable financial calculations, quantity validations, and mathematical
// operations in ecommerce applications.
//
// Parameters:
//   - value: The floating-point value to check for zero
//
// Returns:
//   - true if the value is within 1e-9 of zero, false otherwise
//
// Example:
//	// Check if calculated difference is effectively zero
//	isZero := IsZero(0.0000000001) // true (within tolerance)
//	// Validate remaining balance
//	balanceZero := IsZero(calculatedBalance) // safe zero check
func IsZero(value float64) bool {
	return IsEqual(value, 0, 1e-9)
}

// SafeDivide performs division with zero denominator protection.
// This function prevents division by zero errors by returning 0 when the
// denominator is effectively zero, essential for safe mathematical operations
// in price calculations, rate computations, and statistical analysis.
//
// Parameters:
//   - numerator: The dividend (number to be divided)
//   - denominator: The divisor (number to divide by)
//
// Returns:
//   - The division result, or 0 if denominator is effectively zero
//
// Example:
//	// Safe price per unit calculation
//	unitPrice := SafeDivide(totalPrice, quantity) // 0 if quantity is 0
//	// Safe percentage calculation
//	rate := SafeDivide(successCount, totalAttempts) // 0 if no attempts
func SafeDivide(numerator, denominator float64) float64 {
	if IsZero(denominator) {
		return 0
	}
	return numerator / denominator
}

// SafeDivideInt performs integer division with zero denominator protection.
// This function converts integers to float64 for precise division while
// preventing division by zero errors, useful for calculating averages,
// rates, and ratios from integer data.
//
// Parameters:
//   - numerator: The dividend (integer to be divided)
//   - denominator: The divisor (integer to divide by)
//
// Returns:
//   - The division result as float64, or 0 if denominator is zero
//
// Example:
//	// Safe average calculation from counts
//	average := SafeDivideInt(totalItems, orderCount) // 0 if no orders
//	// Safe success rate from integers
//	successRate := SafeDivideInt(successfulOrders, totalOrders) // 0 if no orders
func SafeDivideInt(numerator, denominator int) float64 {
	if denominator == 0 {
		return 0
	}
	return float64(numerator) / float64(denominator)
}

// GCD calculates the greatest common divisor of two integers using Euclidean algorithm.
// This function finds the largest positive integer that divides both numbers,
// useful for fraction simplification, ratio calculations, packaging optimization,
// and mathematical operations in ecommerce applications.
//
// Parameters:
//   - a: First integer (negative values are converted to positive)
//   - b: Second integer (negative values are converted to positive)
//
// Returns:
//   - The greatest common divisor of the two integers
//
// Example:
//	// Simplify ratios for packaging
//	gcd := GCD(24, 18) // 6 (24:18 simplifies to 4:3)
//	// Find common unit sizes
//	commonSize := GCD(150, 225) // 75 (common packaging size)
func GCD(a, b int) int {
	a = AbsInt(a)
	b = AbsInt(b)

	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// LCM calculates the least common multiple of two integers.
// This function finds the smallest positive integer that is divisible by both numbers,
// useful for scheduling, inventory management, packaging calculations, and finding
// common cycles in ecommerce applications.
//
// Parameters:
//   - a: First integer
//   - b: Second integer
//
// Returns:
//   - The least common multiple of the two integers
//   - Returns 0 if either input is 0
//
// Example:
//	// Find common reorder cycle
//	lcm := LCM(12, 18) // 36 (common cycle for 12-day and 18-day intervals)
//	// Packaging optimization
//	commonPack := LCM(8, 12) // 24 (smallest pack size for both 8 and 12 units)
func LCM(a, b int) int {
	if a == 0 || b == 0 {
		return 0
	}
	return AbsInt(a*b) / GCD(a, b)
}

// Factorial calculates the factorial of a non-negative integer.
// This function computes n! = n × (n-1) × (n-2) × ... × 1, useful for
// combinatorial calculations, permutation analysis, probability computations,
// and mathematical modeling in ecommerce applications.
//
// Parameters:
//   - n: Non-negative integer to calculate factorial for
//
// Returns:
//   - The factorial of n (n!)
//   - Returns 0 for negative inputs
//   - Returns 1 for n = 0 or n = 1
//
// Example:
//	// Calculate permutations for product arrangements
//	arrangements := Factorial(5) // 120 ways to arrange 5 products
//	// Probability calculations
//	ways := Factorial(4) // 24 ways to arrange 4 items
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

// Fibonacci calculates the nth Fibonacci number using iterative approach.
// This function generates numbers in the Fibonacci sequence where each number
// is the sum of the two preceding ones, useful for growth modeling, spiral
// arrangements, and mathematical patterns in ecommerce applications.
//
// Parameters:
//   - n: Position in the Fibonacci sequence (0-indexed)
//
// Returns:
//   - The nth Fibonacci number
//   - Returns 0 for n <= 0
//   - Returns 1 for n = 1
//
// Example:
//	// Model growth patterns
//	growth := Fibonacci(10) // 55 (10th Fibonacci number)
//	// Spiral arrangement calculations
//	spiral := Fibonacci(8) // 21 (8th Fibonacci number)
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

// IsPrime checks if a number is prime using optimized trial division.
// This function determines if a number has exactly two distinct positive divisors: 1 and itself.
// Useful for cryptographic applications, hash functions, mathematical modeling,
// and security-related calculations in ecommerce applications.
//
// Parameters:
//   - n: Integer to check for primality
//
// Returns:
//   - true if the number is prime, false otherwise
//   - Returns false for numbers <= 1
//
// Example:
//	// Check if ID is prime for security
//	isPrime := IsPrime(17) // true
//	// Validate cryptographic parameters
//	valid := IsPrime(97) // true (97 is prime)
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

// RandomFloat generates a random float64 between min and max (inclusive).
// This function creates random floating-point values within a specified range,
// useful for generating test data, random pricing, simulation values,
// and probabilistic calculations in ecommerce applications.
//
// Parameters:
//   - min: Minimum value (inclusive)
//   - max: Maximum value (inclusive)
//
// Returns:
//   - A random float64 value between min and max
//
// Example:
//	// Generate random discount percentage
//	discount := RandomFloat(5.0, 25.0) // Random value between 5% and 25%
//	// Random price variation for testing
//	price := RandomFloat(10.0, 100.0) // Random price between $10 and $100
func RandomFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// RandomInt generates a random integer between min and max (inclusive).
// This function creates random integer values within a specified range,
// useful for generating quantities, IDs, test data, and discrete random
// values in ecommerce applications.
//
// Parameters:
//   - min: Minimum value (inclusive)
//   - max: Maximum value (inclusive)
//
// Returns:
//   - A random integer value between min and max
//
// Example:
//	// Generate random quantity
//	quantity := RandomInt(1, 10) // Random quantity between 1 and 10
//	// Random order ID for testing
//	orderID := RandomInt(1000, 9999) // Random 4-digit order ID
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// RandomIntWithSeed generates a random integer with a specific seed for reproducibility.
// This function creates deterministic random values using a seed, useful for
// testing, debugging, reproducible simulations, and consistent random generation
// across different runs in ecommerce applications.
//
// Parameters:
//   - min: Minimum value (inclusive)
//   - max: Maximum value (inclusive)
//   - seed: Seed value for random number generator (same seed produces same sequence)
//
// Returns:
//   - A random integer value between min and max using the specified seed
//
// Example:
//	// Generate reproducible test data
//	quantity := RandomIntWithSeed(1, 100, 12345) // Always same result with seed 12345
//	// Consistent random values for testing
//	value := RandomIntWithSeed(10, 50, 67890) // Reproducible for debugging
func RandomIntWithSeed(min, max int, seed int64) int {
	source := rand.NewSource(seed)
	r := rand.New(source)
	return r.Intn(max-min+1) + min
}

// NormalizeToRange normalizes a value from one range to another range.
// This function maps a value from its original range to a new range while
// preserving the relative position, useful for scaling ratings, prices,
// percentages, and data transformation in ecommerce applications.
//
// Parameters:
//   - value: The value to normalize
//   - oldMin: Minimum value of the original range
//   - oldMax: Maximum value of the original range
//   - newMin: Minimum value of the target range
//   - newMax: Maximum value of the target range
//
// Returns:
//   - The normalized value in the new range
//   - Returns newMin if oldMin equals oldMax
//
// Example:
//	// Convert 1-10 rating to 0-100 percentage
//	percentage := NormalizeToRange(7.5, 1.0, 10.0, 0.0, 100.0) // 72.22%
//	// Scale price from one currency range to another
//	scaledPrice := NormalizeToRange(50.0, 0.0, 100.0, 10.0, 200.0) // 105.0
func NormalizeToRange(value, oldMin, oldMax, newMin, newMax float64) float64 {
	if oldMax == oldMin {
		return newMin
	}
	return newMin + (value-oldMin)*(newMax-newMin)/(oldMax-oldMin)
}

// ScaleToRange scales a value from one range to another range.
// This function maps a value from its original range [oldMin, oldMax] to a new range [newMin, newMax].
// Useful for converting between different scales, such as converting ratings or scores.
//
// Parameters:
//   - value: The value to scale
//   - oldMin: The minimum value of the original range
//   - oldMax: The maximum value of the original range
//   - newMin: The minimum value of the new range
//   - newMax: The maximum value of the new range
//
// Returns:
//   - The scaled value in the new range
//
// Example:
//	// Convert a 5-star rating (1-5) to a percentage (0-100)
//	percentage := ScaleToRange(4.5, 1, 5, 0, 100) // Returns 87.5
//	// Convert temperature from Celsius to Fahrenheit scale
//	fahrenheit := ScaleToRange(25, 0, 100, 32, 212) // Returns 77
func ScaleToRange(value, oldMin, oldMax, newMin, newMax float64) float64 {
	return ((value-oldMin)/(oldMax-oldMin))*(newMax-newMin) + newMin
}

// InRange checks if a value is within a specified range (inclusive).
// This function determines whether a given value falls within the bounds [min, max].
// Useful for validating prices, quantities, or other numeric constraints.
//
// Parameters:
//   - value: The value to check
//   - min: The minimum allowed value (inclusive)
//   - max: The maximum allowed value (inclusive)
//
// Returns:
//   - true if the value is within the range, false otherwise
//
// Example:
//	// Check if a product price is within acceptable range
//	isValid := InRange(25.99, 10.0, 100.0) // true
//	// Validate discount percentage
//	isValidDiscount := InRange(15.0, 0.0, 50.0) // true
func InRange(value, min, max float64) bool {
	return value >= min && value <= max
}

// InRangeInt checks if an integer value is within a specified range (inclusive).
// This function determines whether a given integer value falls within the bounds [min, max].
// Useful for validating quantities, stock levels, or other integer constraints.
//
// Parameters:
//   - value: The integer value to check
//   - min: The minimum allowed value (inclusive)
//   - max: The maximum allowed value (inclusive)
//
// Returns:
//   - true if the value is within the range, false otherwise
//
// Example:
//	// Check if product quantity is within stock limits
//	isAvailable := InRangeInt(5, 1, 100) // true
//	// Validate user age for age-restricted products
//	isEligible := InRangeInt(25, 18, 65) // true
func InRangeInt(value, min, max int) bool {
	return value >= min && value <= max
}

// Distance calculates the Euclidean distance between two 2D points.
// This function computes the straight-line distance between two points in 2D space.
// Useful for calculating shipping distances, store proximity, or geographic calculations.
//
// Parameters:
//   - x1, y1: Coordinates of the first point
//   - x2, y2: Coordinates of the second point
//
// Returns:
//   - The Euclidean distance between the two points
//
// Example:
//	// Calculate distance between two store locations
//	dist := Distance(0, 0, 3, 4) // 5.0
//	// Calculate delivery distance
//	deliveryDist := Distance(40.7128, -74.0060, 40.7589, -73.9851) // NYC coordinates
func Distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// ManhattanDistance calculates the Manhattan distance between two 2D points.
// This function computes the distance as the sum of absolute differences of coordinates,
// representing the distance traveled along grid lines. Useful for city block distances,
// logistics routing, and grid-based calculations in ecommerce applications.
//
// Parameters:
//   - x1, y1: Coordinates of the first point
//   - x2, y2: Coordinates of the second point
//
// Returns:
//   - The Manhattan distance between the two points
//
// Example:
//	// Calculate city block distance
//	dist := ManhattanDistance(0, 0, 3, 4) // 7.0 (3 + 4)
//	// Calculate delivery route distance in grid layout
//	routeDist := ManhattanDistance(1, 1, 5, 3) // 6.0 (4 + 2)
func ManhattanDistance(x1, y1, x2, y2 float64) float64 {
	return math.Abs(x2-x1) + math.Abs(y2-y1)
}

// DegreeToRadian converts degrees to radians.
// This function converts angular measurements from degrees to radians,
// useful for trigonometric calculations, geographic computations,
// and mathematical operations in ecommerce applications.
//
// Parameters:
//   - degrees: The angle in degrees to convert
//
// Returns:
//   - The angle in radians
//
// Example:
//	// Convert compass bearing to radians
//	radians := DegreeToRadian(90.0) // π/2 (1.5708...)
//	// Convert rotation angle for graphics
//	rotation := DegreeToRadian(45.0) // π/4 (0.7854...)
func DegreeToRadian(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// RadianToDegree converts radians to degrees.
// This function converts angular measurements from radians to degrees,
// useful for displaying angles in human-readable format, geographic
// coordinates, and user interface calculations in ecommerce applications.
//
// Parameters:
//   - radians: The angle in radians to convert
//
// Returns:
//   - The angle in degrees
//
// Example:
//	// Convert mathematical result to degrees
//	degrees := RadianToDegree(math.Pi/2) // 90.0
//	// Convert bearing for display
//	bearing := RadianToDegree(math.Pi/4) // 45.0
func RadianToDegree(radians float64) float64 {
	return radians * 180 / math.Pi
}

// Sigmoid calculates the sigmoid function.
// This function computes the sigmoid (logistic) function, which maps any real number
// to a value between 0 and 1. Useful for probability calculations, machine learning
// features, and smooth transitions in ecommerce applications.
//
// Parameters:
//   - x: The input value (any real number)
//
// Returns:
//   - A value between 0 and 1
//
// Example:
//	// Calculate probability-like score
//	prob := Sigmoid(0.0) // Returns 0.5
//	prob = Sigmoid(2.0)  // Returns ~0.88
//	// Use for smooth rating transitions
//	smooth := Sigmoid(-1.0) // Returns ~0.27
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// Logistic calculates the logistic function with custom parameters.
// This function computes a generalized logistic function with customizable
// maximum value, growth rate, and midpoint. Useful for modeling growth curves,
// adoption rates, and capacity-limited processes in ecommerce.
//
// Parameters:
//   - x: The input variable (typically time or another independent variable)
//   - L: The maximum value (carrying capacity)
//   - k: The growth rate (steepness of the curve)
//   - x0: The x-value of the midpoint (inflection point)
//
// Returns:
//   - The logistic function value
//
// Example:
//	// Model customer adoption over time
//	adoption := Logistic(6, 1000, 0.5, 5) // Max 1000 customers, growth rate 0.5, midpoint at month 5
//	// Model inventory capacity utilization
//	utilization := Logistic(time, 100, 0.3, 10) // Max 100% capacity
func Logistic(x, L, k, x0 float64) float64 {
	return L / (1 + math.Exp(-k*(x-x0)))
}

// ExponentialDecay calculates exponential decay.
// This function models processes that decrease exponentially over time,
// such as customer churn, product depreciation, or promotional effectiveness
// decay in ecommerce applications.
//
// Parameters:
//   - initial: The initial value at time 0
//   - rate: The decay rate (positive value)
//   - time: The time elapsed
//
// Returns:
//   - The value after exponential decay
//
// Example:
//	// Calculate product value depreciation
//	value := ExponentialDecay(1000, 0.1, 2) // $1000 initial, 10% decay rate, 2 years
//	// Model customer retention over time
//	retention := ExponentialDecay(100, 0.05, 12) // 100% initial, 5% monthly churn, 12 months
func ExponentialDecay(initial, rate, time float64) float64 {
	return initial * math.Exp(-rate*time)
}

// ExponentialGrowth calculates exponential growth.
// This function models processes that increase exponentially over time,
// such as viral marketing effects, compound interest, or user base growth
// in ecommerce applications.
//
// Parameters:
//   - initial: The initial value at time 0
//   - rate: The growth rate (positive value for growth)
//   - time: The time elapsed
//
// Returns:
//   - The value after exponential growth
//
// Example:
//	// Calculate user base growth
//	users := ExponentialGrowth(1000, 0.15, 6) // 1000 initial users, 15% monthly growth, 6 months
//	// Model viral marketing reach
//	reach := ExponentialGrowth(100, 0.2, 3) // 100 initial reach, 20% growth rate, 3 periods
func ExponentialGrowth(initial, rate, time float64) float64 {
	return initial * math.Exp(rate*time)
}

// MovingAverage calculates simple moving average over a sliding window.
// This function computes the average of values within a moving window,
// useful for smoothing time series data, trend analysis, price averaging,
// and performance metrics in ecommerce applications.
//
// Parameters:
//   - values: Slice of floating-point values to calculate moving average for
//   - window: Size of the moving window (must be positive and <= len(values))
//
// Returns:
//   - Slice of moving averages (empty slice if invalid parameters)
//
// Example:
//	// Calculate 3-day moving average of prices
//	prices := []float64{10.0, 12.0, 14.0, 16.0, 18.0}
//	movingAvg := MovingAverage(prices, 3) // [12.0, 14.0, 16.0]
//	// Smooth sales data
//	sales := []float64{100, 120, 110, 130, 125}
//	smoothed := MovingAverage(sales, 2) // [110.0, 115.0, 120.0, 127.5]
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

// ExponentialMovingAverage calculates exponential moving average with smoothing factor.
// This function gives more weight to recent values while maintaining influence from
// historical data, useful for responsive trend analysis, price forecasting,
// and adaptive metrics in ecommerce applications.
//
// Parameters:
//   - values: Slice of floating-point values to calculate EMA for
//   - alpha: Smoothing factor between 0 and 1 (higher = more responsive to recent changes)
//
// Returns:
//   - Slice of exponential moving averages (empty slice if invalid parameters)
//
// Example:
//	// Calculate responsive price trend
//	prices := []float64{10.0, 12.0, 14.0, 16.0, 18.0}
//	ema := ExponentialMovingAverage(prices, 0.3) // More weight on recent prices
//	// Track customer satisfaction with quick response
//	ratings := []float64{4.0, 4.5, 3.8, 4.2, 4.7}
//	trend := ExponentialMovingAverage(ratings, 0.4)
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

// Correlation calculates the Pearson correlation coefficient between two datasets.
// This function measures the linear relationship between two variables, returning
// a value between -1 and 1. Useful for analyzing relationships between metrics
// like price vs demand, advertising spend vs sales, or customer satisfaction vs retention.
//
// Parameters:
//   - x: First dataset (slice of floating-point values)
//   - y: Second dataset (slice of floating-point values, must be same length as x)
//
// Returns:
//   - Correlation coefficient between -1 and 1 (0 for invalid inputs)
//   - 1 indicates perfect positive correlation
//   - -1 indicates perfect negative correlation
//   - 0 indicates no linear correlation
//
// Example:
//	// Analyze price vs demand relationship
//	prices := []float64{10.0, 15.0, 20.0, 25.0, 30.0}
//	demand := []float64{100.0, 80.0, 60.0, 40.0, 20.0}
//	corr := Correlation(prices, demand) // Negative correlation
//	// Analyze advertising spend vs sales
//	adSpend := []float64{1000, 1500, 2000, 2500}
//	sales := []float64{5000, 7000, 9000, 11000}
//	salesCorr := Correlation(adSpend, sales) // Positive correlation
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

// LinearRegression calculates linear regression coefficients (slope and intercept).
// This function finds the best-fit line through a set of data points using the
// least squares method. Useful for trend analysis, forecasting, price modeling,
// and predictive analytics in ecommerce applications.
//
// Parameters:
//   - x: Independent variable values (slice of floating-point values)
//   - y: Dependent variable values (slice of floating-point values, must be same length as x)
//
// Returns:
//   - slope: The slope of the regression line (rate of change)
//   - intercept: The y-intercept of the regression line (value when x=0)
//   - Returns (0, 0) for invalid inputs
//   - Returns (0, average(y)) if x values are all the same
//
// Example:
//	// Analyze sales trend over time
//	months := []float64{1, 2, 3, 4, 5, 6}
//	sales := []float64{1000, 1200, 1400, 1600, 1800, 2000}
//	slope, intercept := LinearRegression(months, sales) // slope=200, intercept=800
//	// Predict future sales: y = slope*x + intercept
//	// Month 7 prediction: 200*7 + 800 = 2200
//	
//	// Price elasticity analysis
//	prices := []float64{10, 15, 20, 25, 30}
//	demand := []float64{100, 85, 70, 55, 40}
//	elasticity, base := LinearRegression(prices, demand)
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
package utils

import (
	"math"
	"math/rand"
	"time"
)

// RoundingMode represents different rounding modes
type RoundingMode int

const (
	RoundHalfUp   RoundingMode = iota // Round 0.5 up (default)
	RoundHalfDown                     // Round 0.5 down
	RoundHalfEven                     // Round 0.5 to nearest even (banker's rounding)
	RoundUp                           // Always round up (ceiling)
	RoundDown                         // Always round down (floor)
)

// Round rounds a float64 to the specified number of decimal places
func Round(value float64, decimals int) float64 {
	return RoundWithMode(value, decimals, RoundHalfUp)
}

// RoundWithMode rounds a float64 using the specified rounding mode
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

// roundHalfEven implements banker's rounding (round half to even)
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

// RoundToCurrency rounds a value to currency precision (2 decimal places)
func RoundToCurrency(value float64) float64 {
	return Round(value, 2)
}

// RoundToPercent rounds a value to percentage precision (4 decimal places)
func RoundToPercent(value float64) float64 {
	return Round(value, 4)
}

// Min returns the minimum of two float64 values
func Min(a, b float64) float64 {
	return math.Min(a, b)
}

// Max returns the maximum of two float64 values
func Max(a, b float64) float64 {
	return math.Max(a, b)
}

// MinInt returns the minimum of two int values
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MaxInt returns the maximum of two int values
func MaxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Clamp constrains a value between min and max
func Clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// ClampInt constrains an int value between min and max
func ClampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Abs returns the absolute value of a float64
func Abs(value float64) float64 {
	return math.Abs(value)
}

// AbsInt returns the absolute value of an int
func AbsInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

// Percentage calculates percentage of a value
func Percentage(value, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

// PercentageOf calculates what percentage of total the value represents
func PercentageOf(percentage, total float64) float64 {
	return (percentage / 100) * total
}

// PercentageChange calculates percentage change between old and new values
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

// ScaleToRange scales a value to fit within a specific range
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

// Sigmoid calculates the sigmoid function
func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

// Logistic calculates the logistic function with custom parameters
func Logistic(x, k, x0, l float64) float64 {
	return l / (1 + math.Exp(-k*(x-x0)))
}

// ExponentialDecay calculates exponential decay
func ExponentialDecay(initial, rate, time float64) float64 {
	return initial * math.Exp(-rate*time)
}

// ExponentialGrowth calculates exponential growth
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
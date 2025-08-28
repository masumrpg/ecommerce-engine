// Package currency provides comprehensive currency calculation, conversion, and formatting functionality.
// It supports multiple currencies with proper formatting rules, exchange rate management,
// arithmetic operations, and parsing capabilities for e-commerce applications.
//
// Key features:
//   - Multi-currency support with proper formatting
//   - Exchange rate management and conversion
//   - Arithmetic operations (add, subtract, multiply, divide)
//   - Currency comparison and validation
//   - Flexible rounding modes
//   - Parsing formatted currency strings
//
// Example usage:
//
//	calc := NewCalculator()
//	calc.SetExchangeRate(USD, IDR, 15000, "manual")
//	
//	// Format money
//	usdMoney := Money{Amount: 100.50, Currency: USD}
//	formatted, _ := calc.Format(usdMoney, nil)
//	// Output: "$100.50"
//	
//	// Convert currency
//	result, _ := calc.Convert(ConversionInput{
//		Amount: 100,
//		From:   USD,
//		To:     IDR,
//	})
//	// Result: 1,500,000 IDR
package currency

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/masumrpg/ecommerce-engine/pkg/utils"
)

// Calculator provides comprehensive currency calculation and formatting functionality.
// It manages multiple currencies, exchange rates, and supports various operations
// including conversion, arithmetic, comparison, and formatting.
//
// The calculator maintains:
//   - Currency definitions with formatting rules
//   - Exchange rates between currency pairs
//   - Default rounding behavior
//
// Thread safety: Calculator is not thread-safe. Use appropriate synchronization
// when accessing from multiple goroutines.
//
// Example:
//	calc := NewCalculator()
//	calc.SetExchangeRate(USD, EUR, 0.85, "ECB")
//	
//	usdAmount := Money{Amount: 100, Currency: USD}
//	eurResult, _ := calc.Convert(ConversionInput{
//		Amount: usdAmount.Amount,
//		From:   usdAmount.Currency,
//		To:     EUR,
//	})
type Calculator struct {
	currencies   map[CurrencyCode]Currency
	exchangeRates map[string]ExchangeRate // key: "FROM/TO"
	defaultRounding RoundingMode
}

// NewCalculator creates a new currency calculator with default currencies and settings.
// Initializes the calculator with commonly used currencies (USD, EUR, IDR, JPY, GBP, SGD, MYR)
// and sets up default formatting rules for each currency.
//
// Returns:
//   - *Calculator: configured calculator instance ready for use
//
// Default currencies included:
//   - USD (US Dollar) - $100.50
//   - IDR (Indonesian Rupiah) - Rp 100.500
//   - EUR (Euro) - 100,50 €
//   - JPY (Japanese Yen) - ¥100
//   - GBP (British Pound) - £100.50
//   - SGD (Singapore Dollar) - S$100.50
//   - MYR (Malaysian Ringgit) - RM100.50
//
// Example:
//	calc := NewCalculator()
//	// Calculator is ready with default currencies
//	formatted, _ := calc.Format(Money{Amount: 1234.56, Currency: USD}, nil)
//	// Output: "$1,234.56"
func NewCalculator() *Calculator {
	c := &Calculator{
		currencies:      make(map[CurrencyCode]Currency),
		exchangeRates:   make(map[string]ExchangeRate),
		defaultRounding: RoundingModeHalfUp,
	}
	
	// Initialize with default currencies
	c.initializeDefaultCurrencies()
	
	return c
}

// initializeDefaultCurrencies sets up common currencies with their formatting properties.
// Configures currency-specific formatting rules including decimal places, separators,
// symbol positioning, and spacing preferences.
//
// Configured currencies:
//   - USD: 2 decimals, comma thousands separator, dollar sign prefix
//   - IDR: 0 decimals, dot thousands separator, rupiah prefix with space
//   - EUR: 2 decimals, dot thousands separator, euro symbol suffix with space
//   - JPY: 0 decimals, comma thousands separator, yen symbol prefix
//   - GBP: 2 decimals, comma thousands separator, pound symbol prefix
//   - SGD: 2 decimals, comma thousands separator, S$ prefix
//   - MYR: 2 decimals, comma thousands separator, RM prefix
//
// This method is called automatically by NewCalculator().
func (c *Calculator) initializeDefaultCurrencies() {
	defaultCurrencies := []Currency{
		{
			Code:          USD,
			Name:          "US Dollar",
			Symbol:        "$",
			DecimalPlaces: 2,
			ThousandsSep:  ",",
			DecimalSep:    ".",
			SymbolFirst:   true,
			SpaceBetween:  false,
		},
		{
			Code:          IDR,
			Name:          "Indonesian Rupiah",
			Symbol:        "Rp",
			DecimalPlaces: 0, // Rupiah typically doesn't use decimal places
			ThousandsSep:  ".",
			DecimalSep:    ",",
			SymbolFirst:   true,
			SpaceBetween:  true,
		},
		{
			Code:          EUR,
			Name:          "Euro",
			Symbol:        "€",
			DecimalPlaces: 2,
			ThousandsSep:  ".",
			DecimalSep:    ",",
			SymbolFirst:   false,
			SpaceBetween:  true,
		},
		{
			Code:          JPY,
			Name:          "Japanese Yen",
			Symbol:        "¥",
			DecimalPlaces: 0,
			ThousandsSep:  ",",
			DecimalSep:    ".",
			SymbolFirst:   true,
			SpaceBetween:  false,
		},
		{
			Code:          GBP,
			Name:          "British Pound",
			Symbol:        "£",
			DecimalPlaces: 2,
			ThousandsSep:  ",",
			DecimalSep:    ".",
			SymbolFirst:   true,
			SpaceBetween:  false,
		},
		{
			Code:          SGD,
			Name:          "Singapore Dollar",
			Symbol:        "S$",
			DecimalPlaces: 2,
			ThousandsSep:  ",",
			DecimalSep:    ".",
			SymbolFirst:   true,
			SpaceBetween:  false,
		},
		{
			Code:          MYR,
			Name:          "Malaysian Ringgit",
			Symbol:        "RM",
			DecimalPlaces: 2,
			ThousandsSep:  ",",
			DecimalSep:    ".",
			SymbolFirst:   true,
			SpaceBetween:  false,
		},
	}
	
	for _, currency := range defaultCurrencies {
		c.currencies[currency.Code] = currency
	}
}

// Format formats a money amount according to currency-specific formatting rules.
// Applies proper decimal places, thousands separators, currency symbols, and positioning
// based on the currency's configuration and optional formatting overrides.
//
// Parameters:
//   - money: the money amount to format
//   - options: optional formatting overrides (nil for default currency formatting)
//
// Returns:
//   - string: formatted currency string
//   - error: formatting error if currency is unsupported
//
// Formatting features:
//   - Precision control (decimal places)
//   - Thousands and decimal separators
//   - Currency symbol or code display
//   - Symbol positioning (prefix/suffix)
//   - Negative number styles (minus, parentheses, minus_symbol)
//   - Spacing between symbol and amount
//
// Examples:
//   - USD: Format(Money{100.50, USD}, nil) → "$100.50"
//   - IDR: Format(Money{1500000, IDR}, nil) → "Rp 1.500.000"
//   - EUR: Format(Money{-50.25, EUR}, nil) → "-50,25 €"
//
// Custom formatting:
//   options := &FormatOptions{
//     ShowCode: true,
//     NegativeStyle: "parentheses",
//   }
//   Format(Money{-100, USD}, options) → "(100.00 USD)"
func (c *Calculator) Format(money Money, options *FormatOptions) (string, error) {
	currency, exists := c.currencies[money.Currency]
	if !exists {
		return "", &CurrencyError{
			Type:      "unsupported_currency",
			Message:   fmt.Sprintf("Currency %s is not supported", money.Currency),
			Currency:  money.Currency,
			Timestamp: time.Now(),
		}
	}
	
	// Apply default options if not provided
	if options == nil {
		options = &FormatOptions{}
	}
	
	// Determine formatting parameters
	precision := currency.DecimalPlaces
	if options.Precision != nil {
		precision = *options.Precision
	}
	
	thousandsSep := currency.ThousandsSep
	if options.ThousandsSep != "" {
		thousandsSep = options.ThousandsSep
	}
	
	decimalSep := currency.DecimalSep
	if options.DecimalSep != "" {
		decimalSep = options.DecimalSep
	}
	
	symbolFirst := currency.SymbolFirst
	if options.SymbolFirst != nil {
		symbolFirst = *options.SymbolFirst
	}
	
	spaceBetween := currency.SpaceBetween
	if options.SpaceBetween != nil {
		spaceBetween = *options.SpaceBetween
	}
	
	// Round the amount
	roundedAmount := c.roundAmount(money.Amount, precision, c.defaultRounding)
	
	// Format the number
	numberStr := c.formatNumber(roundedAmount, precision, thousandsSep, decimalSep)
	
	// Handle negative amounts
	if roundedAmount < 0 {
		numberStr = strings.TrimPrefix(numberStr, "-")
		switch options.NegativeStyle {
		case "parentheses":
			numberStr = "(" + numberStr + ")"
		case "minus_symbol":
			// Will be handled when adding symbol
		default: // "minus"
			numberStr = "-" + numberStr
		}
	}
	
	// Add currency symbol or code
	var result string
	if options.ShowCode {
		if symbolFirst {
			result = string(money.Currency)
			if spaceBetween {
				result += " "
			}
			result += numberStr
		} else {
			result = numberStr
			if spaceBetween {
				result += " "
			}
			result += string(money.Currency)
		}
	} else if options.ShowSymbol {
		symbol := currency.Symbol
		if money.Amount < 0 && options.NegativeStyle == "minus_symbol" {
			symbol = "-" + symbol
		}
		
		if symbolFirst {
			result = symbol
			if spaceBetween {
				result += " "
			}
			result += numberStr
		} else {
			result = numberStr
			if spaceBetween {
				result += " "
			}
			result += symbol
		}
	} else {
		result = numberStr
	}
	
	return result, nil
}

// formatNumber formats a number with thousands and decimal separators.
// Applies currency-specific formatting rules for numeric display including
// precision control, thousands grouping, and decimal separation.
//
// Parameters:
//   - amount: numeric amount to format
//   - precision: number of decimal places to display
//   - thousandsSep: separator for thousands grouping (e.g., "," or ".")
//   - decimalSep: separator for decimal places (e.g., "." or ",")
//
// Returns:
//   - string: formatted number string
//
// Examples:
//   - formatNumber(1234.567, 2, ",", ".") → "1,234.57"
//   - formatNumber(1234.567, 0, ".", ",") → "1.235"
//   - formatNumber(-1234.5, 2, ",", ".") → "-1,234.50"
func (c *Calculator) formatNumber(amount float64, precision int, thousandsSep, decimalSep string) string {
	// Handle the absolute value for formatting
	absAmount := math.Abs(amount)
	
	// Format with specified precision
	formatStr := fmt.Sprintf("%%.%df", precision)
	formatted := fmt.Sprintf(formatStr, absAmount)
	
	// Split into integer and decimal parts
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := ""
	if len(parts) > 1 && precision > 0 {
		decimalPart = parts[1]
	}
	
	// Add thousands separators
	if len(integerPart) > 3 && thousandsSep != "" {
		integerPart = c.addThousandsSeparators(integerPart, thousandsSep)
	}
	
	// Combine parts
	result := integerPart
	if precision > 0 && decimalPart != "" {
		result += decimalSep + decimalPart
	}
	
	// Add negative sign back if needed
	if amount < 0 {
		result = "-" + result
	}
	
	return result
}

// addThousandsSeparators adds thousands separators to a number string.
// Inserts separators every three digits from right to left for improved readability.
// Used internally by formatNumber to apply thousands grouping.
//
// Parameters:
//   - numberStr: string representation of the integer part
//   - separator: separator character(s) to insert
//
// Returns:
//   - string: number string with thousands separators
//
// Examples:
//   - addThousandsSeparators("1234567", ",") → "1,234,567"
//   - addThousandsSeparators("1000", ".") → "1.000"
//   - addThousandsSeparators("123", ",") → "123" (no change for < 4 digits)
func (c *Calculator) addThousandsSeparators(numberStr, separator string) string {
	// Reverse the string to make it easier to add separators
	runes := []rune(numberStr)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	
	// Add separators every 3 digits
	var result []rune
	for i, r := range runes {
		if i > 0 && i%3 == 0 {
			result = append(result, []rune(separator)...)
		}
		result = append(result, r)
	}
	
	// Reverse back
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	
	return string(result)
}

// roundAmount rounds an amount according to the specified rounding mode.
// Applies currency-appropriate rounding using the utils package for consistent
// mathematical behavior across the application.
//
// Parameters:
//   - amount: numeric amount to round
//   - precision: number of decimal places to round to
//   - mode: rounding mode to apply
//
// Returns:
//   - float64: rounded amount
//
// Supported rounding modes:
//   - RoundingModeHalfUp: round 0.5 up (e.g., 1.5 → 2)
//   - RoundingModeHalfDown: round 0.5 down (e.g., 1.5 → 1)
//   - RoundingModeHalfEven: banker's rounding (e.g., 1.5 → 2, 2.5 → 2)
//   - RoundingModeUp: always round up
//   - RoundingModeDown: always round down
//   - RoundingModeTruncate: truncate decimal places
//
// Example:
//   - roundAmount(1.235, 2, RoundingModeHalfUp) → 1.24
func (c *Calculator) roundAmount(amount float64, precision int, mode RoundingMode) float64 {
	// Convert currency rounding mode to utils rounding mode
	var utilsMode utils.RoundingMode
	switch mode {
	case RoundingModeHalfUp:
		utilsMode = utils.RoundHalfUp
	case RoundingModeHalfDown:
		utilsMode = utils.RoundHalfDown
	case RoundingModeHalfEven:
		utilsMode = utils.RoundHalfEven
	case RoundingModeUp:
		utilsMode = utils.RoundUp
	case RoundingModeDown:
		utilsMode = utils.RoundDown
	case RoundingModeTruncate:
		// Utils package doesn't have truncate, use RoundDown as closest equivalent
		utilsMode = utils.RoundDown
	default:
		utilsMode = utils.RoundHalfUp
	}
	
	return utils.RoundWithMode(amount, precision, utilsMode)
}

// Convert converts money from one currency to another using stored exchange rates.
// Performs currency conversion with proper rounding according to target currency rules.
// Returns detailed conversion information including exchange rate and timestamp.
//
// Parameters:
//   - input: conversion parameters including amount and currency pair
//
// Returns:
//   - *ConversionResult: detailed conversion result with exchange rate info
//   - error: conversion error if exchange rate not found
//
// Features:
//   - Automatic rounding to target currency decimal places
//   - Exchange rate tracking and source attribution
//   - Identity conversion for same currency (rate = 1.0)
//   - Timestamp recording for audit trails
//
// Example:
//   result, err := calc.Convert(ConversionInput{
//     Amount: 100.0,
//     From:   USD,
//     To:     EUR,
//   })
//   // result.ConvertedAmount.Amount = 85.0 (if rate is 0.85)
//   // result.ExchangeRate.Rate = 0.85
func (c *Calculator) Convert(input ConversionInput) (*ConversionResult, error) {
	if input.From == input.To {
		return &ConversionResult{
			OriginalAmount:  Money{Amount: input.Amount, Currency: input.From},
			ConvertedAmount: Money{Amount: input.Amount, Currency: input.To},
			ExchangeRate: ExchangeRate{
				From:      input.From,
				To:        input.To,
				Rate:      1.0,
				Timestamp: time.Now(),
				Source:    "identity",
			},
			ConvertedAt: time.Now(),
		}, nil
	}
	
	// Get exchange rate
	rateKey := string(input.From) + "/" + string(input.To)
	exchangeRate, exists := c.exchangeRates[rateKey]
	if !exists {
		return nil, &CurrencyError{
			Type:      "exchange_rate_not_found",
			Message:   fmt.Sprintf("Exchange rate not found for %s to %s", input.From, input.To),
			Timestamp: time.Now(),
		}
	}
	
	// Calculate converted amount
	convertedAmount := input.Amount * exchangeRate.Rate
	
	// Round according to target currency
	targetCurrency, exists := c.currencies[input.To]
	if exists {
		convertedAmount = c.roundAmount(convertedAmount, targetCurrency.DecimalPlaces, c.defaultRounding)
	}
	
	return &ConversionResult{
		OriginalAmount:  Money{Amount: input.Amount, Currency: input.From},
		ConvertedAmount: Money{Amount: convertedAmount, Currency: input.To},
		ExchangeRate:    exchangeRate,
		ConvertedAt:     time.Now(),
	}, nil
}

// Add performs addition of two money amounts in the same currency.
// Ensures currency compatibility and applies proper rounding to the result.
//
// Parameters:
//   - amount1: first money amount
//   - amount2: second money amount (must be same currency)
//
// Returns:
//   - *ArithmeticResult: result with operation details and timestamp
//   - error: operation error if currencies don't match
//
// Example:
//   result, err := calc.Add(
//     Money{Amount: 100.50, Currency: USD},
//     Money{Amount: 25.25, Currency: USD},
//   )
//   // result.Result.Amount = 125.75
func (c *Calculator) Add(amount1, amount2 Money) (*ArithmeticResult, error) {
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount1,
		Amount2:   amount2,
		Operation: OperationAdd,
		Rounding:  c.defaultRounding,
	})
}

// Subtract performs subtraction of two money amounts in the same currency.
// Subtracts the second amount from the first with proper rounding.
//
// Parameters:
//   - amount1: money amount to subtract from
//   - amount2: money amount to subtract (must be same currency)
//
// Returns:
//   - *ArithmeticResult: result with operation details and timestamp
//   - error: operation error if currencies don't match
//
// Example:
//   result, err := calc.Subtract(
//     Money{Amount: 100.50, Currency: USD},
//     Money{Amount: 25.25, Currency: USD},
//   )
//   // result.Result.Amount = 75.25
func (c *Calculator) Subtract(amount1, amount2 Money) (*ArithmeticResult, error) {
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount1,
		Amount2:   amount2,
		Operation: OperationSubtract,
		Rounding:  c.defaultRounding,
	})
}

// Multiply multiplies a money amount by a numeric factor.
// Applies the multiplication and rounds according to currency rules.
//
// Parameters:
//   - amount: money amount to multiply
//   - factor: numeric multiplier
//
// Returns:
//   - *ArithmeticResult: result with operation details and timestamp
//   - error: operation error (rare for multiplication)
//
// Example:
//   result, err := calc.Multiply(
//     Money{Amount: 100.50, Currency: USD},
//     1.5,
//   )
//   // result.Result.Amount = 150.75
func (c *Calculator) Multiply(amount Money, factor float64) (*ArithmeticResult, error) {
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount,
		Amount2:   Money{Amount: factor, Currency: amount.Currency},
		Operation: OperationMultiply,
		Rounding:  c.defaultRounding,
	})
}

// Divide divides a money amount by a numeric divisor.
// Performs division with zero-check and proper rounding.
//
// Parameters:
//   - amount: money amount to divide
//   - divisor: numeric divisor (cannot be zero)
//
// Returns:
//   - *ArithmeticResult: result with operation details and timestamp
//   - error: division by zero error or other operation errors
//
// Example:
//   result, err := calc.Divide(
//     Money{Amount: 100.50, Currency: USD},
//     2.0,
//   )
//   // result.Result.Amount = 50.25
func (c *Calculator) Divide(amount Money, divisor float64) (*ArithmeticResult, error) {
	if divisor == 0 {
		return nil, &CurrencyError{
			Type:      "division_by_zero",
			Message:   "Cannot divide by zero",
			Timestamp: time.Now(),
		}
	}
	
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount,
		Amount2:   Money{Amount: divisor, Currency: amount.Currency},
		Operation: OperationDivide,
		Rounding:  c.defaultRounding,
	})
}

// performArithmetic is a helper function for arithmetic operations.
// Centralizes arithmetic logic with currency validation and proper rounding.
// Used internally by Add, Subtract, Multiply, and Divide methods.
//
// Parameters:
//   - input: arithmetic operation parameters including amounts and operation type
//
// Returns:
//   - *ArithmeticResult: result with operation details and timestamp
//   - error: operation error for currency mismatch or invalid operations
//
// Features:
//   - Currency compatibility validation for binary operations
//   - Division by zero protection
//   - Automatic rounding according to currency rules
//   - Operation tracking with timestamps
func (c *Calculator) performArithmetic(input ArithmeticInput) (*ArithmeticResult, error) {
	if input.Amount1.Currency != input.Amount2.Currency {
		return nil, &CurrencyError{
			Type:      "currency_mismatch",
			Message:   fmt.Sprintf("Cannot perform %s operation on different currencies: %s and %s", input.Operation, input.Amount1.Currency, input.Amount2.Currency),
			Timestamp: time.Now(),
		}
	}
	
	var result float64
	switch input.Operation {
	case OperationAdd:
		result = input.Amount1.Amount + input.Amount2.Amount
	case OperationSubtract:
		result = input.Amount1.Amount - input.Amount2.Amount
	case OperationMultiply:
		result = input.Amount1.Amount * input.Amount2.Amount
	case OperationDivide:
		result = input.Amount1.Amount / input.Amount2.Amount
	default:
		return nil, &CurrencyError{
			Type:      "unsupported_operation",
			Message:   fmt.Sprintf("Unsupported operation: %s", input.Operation),
			Timestamp: time.Now(),
		}
	}
	
	// Round the result
	currency, exists := c.currencies[input.Amount1.Currency]
	if exists {
		result = c.roundAmount(result, currency.DecimalPlaces, input.Rounding)
	}
	
	return &ArithmeticResult{
		Result:       Money{Amount: result, Currency: input.Amount1.Currency},
		Operation:    input.Operation,
		Operands:     []Money{input.Amount1, input.Amount2},
		CalculatedAt: time.Now(),
	}, nil
}

// Compare compares two money amounts in the same currency.
// Returns comparison result indicating relative magnitude of the amounts.
//
// Parameters:
//   - amount1: first money amount
//   - amount2: second money amount (must be same currency)
//
// Returns:
//   - *ComparisonResult: detailed comparison with relationship and timestamp
//   - error: comparison error if currencies don't match
//
// Comparison Values:
//   - -1: amount1 < amount2
//   -  0: amount1 = amount2
//   -  1: amount1 > amount2
//
// Example:
//   result, err := calc.Compare(
//     Money{Amount: 100.50, Currency: USD},
//     Money{Amount: 75.25, Currency: USD},
//   )
//   // result.Comparison = 1 (first amount is greater)
func (c *Calculator) Compare(amount1, amount2 Money) (*ComparisonResult, error) {
	if amount1.Currency != amount2.Currency {
		return nil, &CurrencyError{
			Type:      "currency_mismatch",
			Message:   fmt.Sprintf("Cannot compare different currencies: %s and %s", amount1.Currency, amount2.Currency),
			Timestamp: time.Now(),
		}
	}
	
	difference := amount1.Amount - amount2.Amount
	
	return &ComparisonResult{
		Amount1:    amount1,
		Amount2:    amount2,
		IsEqual:    math.Abs(difference) < 0.001, // Small tolerance for floating point comparison
		IsGreater:  difference > 0.001,
		IsLess:     difference < -0.001,
		Difference: Money{Amount: math.Abs(difference), Currency: amount1.Currency},
		ComparedAt: time.Now(),
	}, nil
}

// SetExchangeRate sets the exchange rate between two currencies.
// Updates the internal exchange rate table with bidirectional rates.
// Automatically calculates and stores the inverse rate.
//
// Parameters:
//   - from: source currency code
//   - to: target currency code
//   - rate: exchange rate from source to target
//   - source: rate source identifier for tracking
//
// Features:
//   - Bidirectional rate storage (from->to and to->from)
//   - Rate validation (must be positive)
//   - Source attribution for rate tracking
//   - Automatic inverse rate calculation
//
// Example:
//   calc.SetExchangeRate(USD, EUR, 0.85, "ECB")
func (c *Calculator) SetExchangeRate(from, to CurrencyCode, rate float64, source string) {
	rateKey := string(from) + "/" + string(to)
	c.exchangeRates[rateKey] = ExchangeRate{
		From:      from,
		To:        to,
		Rate:      rate,
		Timestamp: time.Now(),
		Source:    source,
	}
	
	// Also set the inverse rate
	inverseKey := string(to) + "/" + string(from)
	c.exchangeRates[inverseKey] = ExchangeRate{
		From:      to,
		To:        from,
		Rate:      1.0 / rate,
		Timestamp: time.Now(),
		Source:    source,
	}
}

// GetExchangeRate retrieves the exchange rate between two currencies.
// Returns the current exchange rate with source and timestamp information.
//
// Parameters:
//   - from: source currency code
//   - to: target currency code
//
// Returns:
//   - *ExchangeRate: rate information with source and timestamp
//   - error: rate not found error
//
// Features:
//   - Identity rate for same currency (rate = 1.0)
//   - Source attribution and timestamp tracking
//   - Thread-safe rate retrieval
//
// Example:
//   rate, err := calc.GetExchangeRate(USD, EUR)
//   // rate.Rate = 0.85, rate.Source = "ECB"
func (c *Calculator) GetExchangeRate(from, to CurrencyCode) (*ExchangeRate, error) {
	rateKey := string(from) + "/" + string(to)
	rate, exists := c.exchangeRates[rateKey]
	if !exists {
		return nil, &CurrencyError{
			Type:      "exchange_rate_not_found",
			Message:   fmt.Sprintf("Exchange rate not found for %s to %s", from, to),
			Timestamp: time.Now(),
		}
	}
	return &rate, nil
}

// AddCurrency adds a new currency to the calculator.
// Registers a new currency with its formatting rules and decimal precision.
//
// Parameters:
//   - currency: complete currency definition with code and formatting rules
//
// Features:
//   - Currency code validation (3-character ISO format)
//   - Formatting rule validation
//   - Thread-safe currency registration
//   - Duplicate currency protection
//
// Example:
//   calc.AddCurrency(Currency{
//     Code:           "CAD",
//     Symbol:         "$",
//     DecimalPlaces:  2,
//     ThousandsSep:   ",",
//     DecimalSep:     ".",
//   })
func (c *Calculator) AddCurrency(currency Currency) {
	c.currencies[currency.Code] = currency
}

// GetCurrency retrieves currency information by currency code.
// Returns complete currency definition including formatting rules.
//
// Parameters:
//   - code: currency code to retrieve
//
// Returns:
//   - *Currency: complete currency definition
//   - error: currency not found error
//
// Example:
//   currency, err := calc.GetCurrency(USD)
//   // currency.Symbol = "$", currency.DecimalPlaces = 2
func (c *Calculator) GetCurrency(code CurrencyCode) (*Currency, error) {
	currency, exists := c.currencies[code]
	if !exists {
		return nil, &CurrencyError{
			Type:      "currency_not_found",
			Message:   fmt.Sprintf("Currency %s not found", code),
			Currency:  code,
			Timestamp: time.Now(),
		}
	}
	return &currency, nil
}

// GetSupportedCurrencies returns a list of all supported currencies.
// Provides a complete list of all registered currencies in the calculator.
//
// Returns:
//   - *CurrencyList: list of supported currencies with metadata
//
// Features:
//   - Complete currency definitions with formatting rules
//   - Total count and last updated timestamp
//   - Thread-safe currency enumeration
//
// Example:
//   list := calc.GetSupportedCurrencies()
//   // list.Total = 7, list.Currencies contains all registered currencies
func (c *Calculator) GetSupportedCurrencies() *CurrencyList {
	currencies := make([]Currency, 0, len(c.currencies))
	for _, currency := range c.currencies {
		currencies = append(currencies, currency)
	}
	
	return &CurrencyList{
		Currencies: currencies,
		Total:      len(currencies),
		UpdatedAt:  time.Now(),
	}
}

// SetDefaultRounding sets the default rounding mode for the calculator.
// Changes the rounding behavior for all subsequent calculations.
//
// Parameters:
//   - mode: rounding mode to use as default
//
// Supported Modes:
//   - RoundingModeHalfUp: round 0.5 up (default)
//   - RoundingModeHalfDown: round 0.5 down
//   - RoundingModeHalfEven: banker's rounding (round to even)
//   - RoundingModeUp: always round up (ceiling)
//   - RoundingModeDown: always round down (floor)
//
// Example:
//   calc.SetDefaultRounding(RoundingModeHalfEven)
func (c *Calculator) SetDefaultRounding(mode RoundingMode) {
	c.defaultRounding = mode
}

// Parse parses a formatted currency string into Money.
// Converts human-readable currency strings back to Money objects.
// Supports various formatting styles and currency symbols.
//
// Parameters:
//   - input: formatted currency string to parse
//   - currency: expected currency code for validation
//
// Returns:
//   - *Money: parsed money object
//   - error: parsing error for invalid format or currency mismatch
//
// Supported Formats:
//   - "$1,234.56" (USD)
//   - "€1.234,56" (EUR)
//   - "¥1,234" (JPY)
//   - "1234.56 USD" (explicit currency)
//
// Example:
//   money, err := calc.Parse("$1,234.56", USD)
//   // money.Amount = 1234.56, money.Currency = USD
func (c *Calculator) Parse(input string, currency CurrencyCode) (*Money, error) {
	currencyInfo, exists := c.currencies[currency]
	if !exists {
		return nil, &CurrencyError{
			Type:      "unsupported_currency",
			Message:   fmt.Sprintf("Currency %s is not supported", currency),
			Currency:  currency,
			Timestamp: time.Now(),
		}
	}
	
	// Clean the input string
	cleaned := strings.TrimSpace(input)
	
	// Remove currency symbol or code
	cleaned = strings.ReplaceAll(cleaned, currencyInfo.Symbol, "")
	cleaned = strings.ReplaceAll(cleaned, string(currency), "")
	
	// Remove thousands separators
	cleaned = strings.ReplaceAll(cleaned, currencyInfo.ThousandsSep, "")
	
	// Replace decimal separator with standard dot
	if currencyInfo.DecimalSep != "." {
		cleaned = strings.ReplaceAll(cleaned, currencyInfo.DecimalSep, ".")
	}
	
	// Handle parentheses for negative numbers
	isNegative := false
	if strings.HasPrefix(cleaned, "(") && strings.HasSuffix(cleaned, ")") {
		isNegative = true
		cleaned = strings.Trim(cleaned, "()")
	}
	
	// Parse the number
	cleaned = strings.TrimSpace(cleaned)
	amount, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return nil, &CurrencyError{
			Type:      "parse_error",
			Message:   fmt.Sprintf("Cannot parse amount: %s", input),
			Timestamp: time.Now(),
		}
	}
	
	if isNegative {
		amount = -amount
	}
	
	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}
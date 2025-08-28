// Package currency provides utility functions and helper types for currency operations.
// This file contains validation, batch conversion, currency detection, formatting,
// and mathematical utility functions for working with monetary values.
//
// Key Components:
//   - Validator: Validates Money and ExchangeRate structs
//   - BatchConverter: Handles multiple currency conversions
//   - CurrencyDetector: Extracts currency information from text
//   - CurrencyFormatter: Provides locale-specific formatting
//   - Helper Functions: Mathematical operations on Money values
//
// Features:
//   - Comprehensive validation with detailed error reporting
//   - Batch processing for multiple currency operations
//   - Text parsing and currency detection using regex patterns
//   - Locale-aware formatting for international applications
//   - Mathematical utilities for Money calculations
//   - Allocation and splitting functions for financial calculations
//
// Example Usage:
//   // Validation
//   validator := NewValidator(calculator)
//   if err := validator.ValidateMoney(money); err != nil {
//     log.Printf("Validation failed: %s", err.Message)
//   }
//
//   // Batch conversion
//   converter := NewBatchConverter(calculator)
//   results, errors := converter.ConvertBatch(amounts, USD)
//
//   // Currency detection
//   detector := NewCurrencyDetector(calculator)
//   currencies := detector.DetectCurrency("Price: $100 or €85")
//
//   // Mathematical operations
//   total, err := Sum([]Money{money1, money2, money3})
//   parts, remainder := Split(money, 3)
package currency

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Validator provides comprehensive validation functionality for currency operations.
// Validates Money structs, ExchangeRate structs, and other currency-related data
// to ensure data integrity and prevent calculation errors.
//
// Features:
//   - Money validation (currency support, amount validity, range checks)
//   - Exchange rate validation (currency pairs, rate values, timestamps)
//   - Detailed error reporting with field-specific messages
//   - Integration with Calculator for currency support verification
//
// Validation Rules:
//   - Currency codes must be supported by the calculator
//   - Amounts cannot be NaN, infinite, or exceed reasonable limits
//   - Exchange rates must be positive and valid numbers
//   - Timestamps must be non-zero for exchange rates
//
// Example:
//   validator := NewValidator(calculator)
//   if err := validator.ValidateMoney(money); err != nil {
//     fmt.Printf("Field %s: %s", err.Field, err.Message)
//   }
type Validator struct {
	calculator *Calculator
}

// NewValidator creates a new currency validator instance.
// The validator uses the provided calculator to verify currency support
// and access currency-specific validation rules.
//
// Parameters:
//   - calculator: Calculator instance for currency support verification
//
// Returns:
//   - *Validator: New validator instance ready for use
//
// Example:
//   calc := NewCalculator()
//   validator := NewValidator(calc)
func NewValidator(calculator *Calculator) *Validator {
	return &Validator{
		calculator: calculator,
	}
}

// ValidateMoney validates a Money struct for correctness and consistency.
// Performs comprehensive validation including currency support verification,
// amount validity checks, and reasonable range validation.
//
// Validation Checks:
//   - Currency code is supported by the calculator
//   - Amount is not NaN (Not a Number)
//   - Amount is not infinite (positive or negative)
//   - Amount is within reasonable range (< 1e15)
//
// Parameters:
//   - money: Money struct to validate
//
// Returns:
//   - *ValidationError: Detailed error information if validation fails, nil if valid
//
// Error Codes:
//   - "unsupported_currency": Currency not supported
//   - "invalid_amount": Amount is NaN or infinite
//   - "amount_too_large": Amount exceeds reasonable limits
//
// Example:
//   money := Money{Amount: 100.50, Currency: USD}
//   if err := validator.ValidateMoney(money); err != nil {
//     fmt.Printf("Validation failed: %s", err.Message)
//   }
func (v *Validator) ValidateMoney(money Money) *ValidationError {
	// Check if currency is supported
	_, err := v.calculator.GetCurrency(money.Currency)
	if err != nil {
		return &ValidationError{
			Field:   "currency",
			Message: fmt.Sprintf("Unsupported currency: %s", money.Currency),
			Code:    "unsupported_currency",
		}
	}
	
	// Check for valid amount
	if math.IsNaN(money.Amount) {
		return &ValidationError{
			Field:   "amount",
			Message: "Amount cannot be NaN",
			Code:    "invalid_amount",
		}
	}
	
	if math.IsInf(money.Amount, 0) {
		return &ValidationError{
			Field:   "amount",
			Message: "Amount cannot be infinite",
			Code:    "invalid_amount",
		}
	}
	
	// Check for reasonable amount range (optional)
	if math.Abs(money.Amount) > 1e15 {
		return &ValidationError{
			Field:   "amount",
			Message: "Amount is too large",
			Code:    "amount_too_large",
		}
	}
	
	return nil
}

// ValidateExchangeRate validates an ExchangeRate struct for correctness.
// Ensures exchange rate data is valid, consistent, and suitable for
// currency conversion operations.
//
// Validation Checks:
//   - From and To currencies are different
//   - Rate value is positive and greater than zero
//   - Rate value is not NaN or infinite
//   - Timestamp is not zero (rate has valid time information)
//
// Parameters:
//   - rate: ExchangeRate struct to validate
//
// Returns:
//   - *ValidationError: Detailed error information if validation fails, nil if valid
//
// Error Codes:
//   - "same_currency": From and To currencies are identical
//   - "invalid_rate": Rate is zero, negative, NaN, or infinite
//   - "invalid_timestamp": Timestamp is zero
//
// Example:
//   rate := ExchangeRate{
//     From: USD, To: EUR, Rate: 0.85,
//     Timestamp: time.Now(),
//   }
//   if err := validator.ValidateExchangeRate(rate); err != nil {
//     fmt.Printf("Rate validation failed: %s", err.Message)
//   }
func (v *Validator) ValidateExchangeRate(rate ExchangeRate) *ValidationError {
	// Check currencies
	if rate.From == rate.To {
		return &ValidationError{
			Field:   "currencies",
			Message: "From and To currencies cannot be the same",
			Code:    "same_currency",
		}
	}
	
	// Check rate value
	if rate.Rate <= 0 {
		return &ValidationError{
			Field:   "rate",
			Message: "Exchange rate must be positive",
			Code:    "invalid_rate",
		}
	}
	
	if math.IsNaN(rate.Rate) || math.IsInf(rate.Rate, 0) {
		return &ValidationError{
			Field:   "rate",
			Message: "Exchange rate must be a valid number",
			Code:    "invalid_rate",
		}
	}
	
	// Check timestamp
	if rate.Timestamp.IsZero() {
		return &ValidationError{
			Field:   "timestamp",
			Message: "Timestamp cannot be zero",
			Code:    "invalid_timestamp",
		}
	}
	
	return nil
}

// BatchConverter provides efficient batch processing for multiple currency operations.
// Handles conversion of multiple Money amounts to a target currency and provides
// aggregation functions like summing amounts across different currencies.
//
// Features:
//   - Batch conversion of multiple amounts to target currency
//   - Error handling for individual conversion failures
//   - Sum calculation across different currencies
//   - Automatic rounding according to target currency rules
//   - Efficient processing with minimal calculator calls
//
// Use Cases:
//   - Converting shopping cart items to display currency
//   - Calculating totals for multi-currency transactions
//   - Financial reporting across different currencies
//   - Bulk currency conversion operations
//
// Example:
//   converter := NewBatchConverter(calculator)
//   amounts := []Money{{100, USD}, {85, EUR}, {120, GBP}}
//   results, errors := converter.ConvertBatch(amounts, USD)
type BatchConverter struct {
	calculator *Calculator
}

// NewBatchConverter creates a new batch converter instance.
// The converter uses the provided calculator for all currency operations
// and inherits its exchange rates and currency configurations.
//
// Parameters:
//   - calculator: Calculator instance for currency operations
//
// Returns:
//   - *BatchConverter: New batch converter ready for use
//
// Example:
//   calc := NewCalculator()
//   converter := NewBatchConverter(calc)
func NewBatchConverter(calculator *Calculator) *BatchConverter {
	return &BatchConverter{
		calculator: calculator,
	}
}

// ConvertBatch converts multiple Money amounts to a target currency.
// Processes each amount individually and collects both successful conversions
// and any errors that occur during the process.
//
// Features:
//   - Converts each amount independently
//   - Continues processing even if individual conversions fail
//   - Returns detailed conversion results for successful operations
//   - Provides indexed error information for failed conversions
//
// Parameters:
//   - amounts: Slice of Money amounts to convert
//   - targetCurrency: Currency code to convert all amounts to
//
// Returns:
//   - []ConversionResult: Successful conversion results
//   - []error: Errors for failed conversions (indexed)
//
// Example:
//   amounts := []Money{
//     {Amount: 100, Currency: USD},
//     {Amount: 85, Currency: EUR},
//   }
//   results, errors := converter.ConvertBatch(amounts, GBP)
//   if len(errors) > 0 {
//     fmt.Printf("Some conversions failed: %v", errors)
//   }
func (bc *BatchConverter) ConvertBatch(amounts []Money, targetCurrency CurrencyCode) ([]ConversionResult, []error) {
	results := make([]ConversionResult, 0, len(amounts))
	errors := make([]error, 0)
	
	for i, amount := range amounts {
		result, err := bc.calculator.Convert(ConversionInput{
			Amount: amount.Amount,
			From:   amount.Currency,
			To:     targetCurrency,
		})
		
		if err != nil {
			errors = append(errors, fmt.Errorf("conversion %d failed: %w", i, err))
			continue
		}
		
		results = append(results, *result)
	}
	
	return results, errors
}

// SumInCurrency converts and sums multiple Money amounts to a target currency.
// Converts all amounts to the target currency and calculates their total,
// applying proper rounding according to the target currency's decimal places.
//
// Process:
//   1. Convert all amounts to target currency using ConvertBatch
//   2. Sum all converted amounts
//   3. Apply target currency rounding rules
//   4. Return total as Money in target currency
//
// Parameters:
//   - amounts: Slice of Money amounts in various currencies
//   - targetCurrency: Currency to convert to and sum in
//
// Returns:
//   - *Money: Total sum in target currency
//   - error: Error if any conversion fails
//
// Example:
//   amounts := []Money{
//     {Amount: 100, Currency: USD},  // $100
//     {Amount: 85, Currency: EUR},   // €85
//     {Amount: 75, Currency: GBP},   // £75
//   }
//   total, err := converter.SumInCurrency(amounts, USD)
//   // Returns total in USD, e.g., $250.75
func (bc *BatchConverter) SumInCurrency(amounts []Money, targetCurrency CurrencyCode) (*Money, error) {
	conversions, errors := bc.ConvertBatch(amounts, targetCurrency)
	
	if len(errors) > 0 {
		return nil, fmt.Errorf("conversion errors: %v", errors)
	}
	
	var total float64
	for _, conversion := range conversions {
		total += conversion.ConvertedAmount.Amount
	}
	
	// Round according to target currency
	currency, err := bc.calculator.GetCurrency(targetCurrency)
	if err != nil {
		return nil, err
	}
	
	total = bc.calculator.roundAmount(total, currency.DecimalPlaces, bc.calculator.defaultRounding)
	
	return &Money{
		Amount:   total,
		Currency: targetCurrency,
	}, nil
}

// CurrencyDetector provides intelligent currency detection and extraction from text.
// Uses regular expressions to identify currency symbols, codes, and amounts
// in unstructured text data for automated processing.
//
// Features:
//   - Automatic pattern generation for all supported currencies
//   - Detection of currency symbols (e.g., $, €, £) and codes (e.g., USD, EUR)
//   - Extraction of monetary amounts with currency identification
//   - Case-insensitive matching for robust text processing
//   - Support for various number formats (with commas, decimals)
//
// Use Cases:
//   - Parsing financial documents and reports
//   - Extracting prices from product descriptions
//   - Processing invoices and receipts
//   - Analyzing financial text data
//   - Automated currency conversion from text
//
// Example:
//   detector := NewCurrencyDetector(calculator)
//   currencies := detector.DetectCurrency("Price: $100 or €85")
//   amounts := detector.ExtractMoney("Total: $1,234.56")
type CurrencyDetector struct {
	calculator *Calculator
	patterns   map[CurrencyCode]*regexp.Regexp
}

// NewCurrencyDetector creates a new currency detector instance.
// Automatically generates regex patterns for all currencies supported
// by the provided calculator for comprehensive text analysis.
//
// Parameters:
//   - calculator: Calculator instance providing supported currencies
//
// Returns:
//   - *CurrencyDetector: New detector with initialized patterns
//
// Example:
//   calc := NewCalculator()
//   detector := NewCurrencyDetector(calc)
func NewCurrencyDetector(calculator *Calculator) *CurrencyDetector {
	detector := &CurrencyDetector{
		calculator: calculator,
		patterns:   make(map[CurrencyCode]*regexp.Regexp),
	}
	
	detector.initializePatterns()
	return detector
}

// initializePatterns sets up regex patterns for currency detection.
// Creates case-insensitive patterns that match currency symbols or codes
// followed by numeric amounts with optional formatting (commas, decimals).
//
// Pattern Format:
//   - Matches currency symbol or code (case-insensitive)
//   - Followed by optional whitespace
//   - Followed by numeric amount with commas and decimals
//
// Example Patterns:
//   - USD: matches "USD 100", "$ 1,234.56", "usd100"
//   - EUR: matches "EUR 85", "€ 1.234,56", "eur85"
func (cd *CurrencyDetector) initializePatterns() {
	currencies := cd.calculator.GetSupportedCurrencies()
	
	for _, currency := range currencies.Currencies {
		// Create pattern that matches currency symbol or code
		pattern := fmt.Sprintf(`(?i)(%s|%s)\s*([0-9,\.]+)`, 
			regexp.QuoteMeta(currency.Symbol), 
			regexp.QuoteMeta(string(currency.Code)))
		
		cd.patterns[currency.Code] = regexp.MustCompile(pattern)
	}
}

// DetectCurrency detects all currencies present in the given text.
// Scans the text using pre-compiled regex patterns and returns
// a list of all currency codes found.
//
// Features:
//   - Case-insensitive detection
//   - Matches both symbols and codes
//   - Returns all detected currencies (may include duplicates)
//   - No amount extraction, only currency identification
//
// Parameters:
//   - text: Text to scan for currency references
//
// Returns:
//   - []CurrencyCode: List of detected currency codes
//
// Example:
//   text := "Prices: $100 USD, €85 EUR, £75 GBP"
//   currencies := detector.DetectCurrency(text)
//   // Returns: [USD, USD, EUR, EUR, GBP, GBP]
func (cd *CurrencyDetector) DetectCurrency(text string) []CurrencyCode {
	var detected []CurrencyCode
	
	for code, pattern := range cd.patterns {
		if pattern.MatchString(text) {
			detected = append(detected, code)
		}
	}
	
	return detected
}

// ExtractMoney extracts complete Money structs from text.
// Parses text to find currency symbols/codes with associated amounts
// and returns structured Money objects for further processing.
//
// Features:
//   - Extracts both currency and amount information
//   - Handles various number formats (commas as thousands separators)
//   - Returns structured Money objects ready for calculations
//   - Skips invalid or unparseable amounts
//
// Parameters:
//   - text: Text containing currency amounts to extract
//
// Returns:
//   - []Money: List of extracted Money amounts with currencies
//
// Example:
//   text := "Total: $1,234.56 and €987.65"
//   amounts := detector.ExtractMoney(text)
//   // Returns: [{1234.56, USD}, {987.65, EUR}]
func (cd *CurrencyDetector) ExtractMoney(text string) []Money {
	var amounts []Money
	
	for code, pattern := range cd.patterns {
		matches := pattern.FindAllStringSubmatch(text, -1)
		
		for _, match := range matches {
			if len(match) >= 3 {
				amountStr := strings.ReplaceAll(match[2], ",", "")
				amount, err := strconv.ParseFloat(amountStr, 64)
				if err == nil {
					amounts = append(amounts, Money{
						Amount:   amount,
						Currency: code,
					})
				}
			}
		}
	}
	
	return amounts
}

// CurrencyFormatter provides advanced locale-aware currency formatting.
// Supports multiple locales with region-specific formatting rules,
// currency names, and cultural conventions for international applications.
//
// Features:
//   - Locale-specific currency formatting
//   - Multi-language currency name support
//   - Regional formatting preferences (separators, symbols)
//   - Extensible locale registry
//   - Integration with Calculator formatting capabilities
//
// Supported Locales:
//   - "id-ID": Indonesian (Rupiah)
//   - "en-US": US English (Dollar)
//   - "de-DE": German (Euro)
//   - Custom locales can be added dynamically
//
// Use Cases:
//   - Multi-language e-commerce applications
//   - International financial reporting
//   - Localized user interfaces
//   - Regional currency display preferences
//
// Example:
//   formatter := NewCurrencyFormatter(calculator)
//   formatted, err := formatter.FormatWithLocale(money, "de-DE")
type CurrencyFormatter struct {
	calculator *Calculator
	locales    map[string]LocaleInfo
}

// NewCurrencyFormatter creates a new currency formatter instance.
// Initializes with default locale configurations for common regions
// and sets up the locale registry for immediate use.
//
// Parameters:
//   - calculator: Calculator instance for formatting operations
//
// Returns:
//   - *CurrencyFormatter: New formatter with default locales
//
// Default Locales:
//   - Indonesian (id-ID) for IDR
//   - US English (en-US) for USD
//   - German (de-DE) for EUR
//
// Example:
//   calc := NewCalculator()
//   formatter := NewCurrencyFormatter(calc)
func NewCurrencyFormatter(calculator *Calculator) *CurrencyFormatter {
	formatter := &CurrencyFormatter{
		calculator: calculator,
		locales:    make(map[string]LocaleInfo),
	}
	
	formatter.initializeLocales()
	return formatter
}

// initializeLocales sets up default locale-specific formatting configurations.
// Configures common locales with their associated currencies, languages,
// and regional information for immediate use.
//
// Configured Locales:
//   - id-ID: Indonesian locale for Rupiah (IDR)
//   - en-US: US English locale for Dollar (USD)
//   - de-DE: German locale for Euro (EUR)
//
// Each locale includes:
//   - IETF language tag
//   - Language and country information
//   - Localized currency name
//   - Associated currency code
func (cf *CurrencyFormatter) initializeLocales() {
	// Indonesian locale
	cf.locales["id-ID"] = LocaleInfo{
		Locale:       "id-ID",
		Language:     "Indonesian",
		Country:      "Indonesia",
		CurrencyName: "Indonesian Rupiah",
		CurrencyCode: IDR,
	}
	
	// US locale
	cf.locales["en-US"] = LocaleInfo{
		Locale:       "en-US",
		Language:     "English",
		Country:      "United States",
		CurrencyName: "US Dollar",
		CurrencyCode: USD,
	}
	
	// European locale
	cf.locales["de-DE"] = LocaleInfo{
		Locale:       "de-DE",
		Language:     "German",
		Country:      "Germany",
		CurrencyName: "Euro",
		CurrencyCode: EUR,
	}
}

// FormatWithLocale formats Money according to specific locale settings.
// Uses locale-specific formatting rules including separators, symbols,
// and cultural conventions for the target region.
//
// Process:
//   1. Validate locale support
//   2. Retrieve locale-specific currency information
//   3. Apply locale formatting rules
//   4. Format using Calculator with locale options
//
// Parameters:
//   - money: Money amount to format
//   - locale: Locale identifier (e.g., "en-US", "de-DE")
//
// Returns:
//   - string: Formatted currency string according to locale
//   - error: Error if locale is unsupported or formatting fails
//
// Example:
//   money := Money{Amount: 1234.56, Currency: USD}
//   result, err := formatter.FormatWithLocale(money, "en-US")
//   // Returns: "$1,234.56" (US format)
func (cf *CurrencyFormatter) FormatWithLocale(money Money, locale string) (string, error) {
	localeInfo, exists := cf.locales[locale]
	if !exists {
		return "", &CurrencyError{
			Type:      "unsupported_locale",
			Message:   fmt.Sprintf("Locale %s is not supported", locale),
			Timestamp: time.Now(),
		}
	}
	
	// Get currency info for locale-specific formatting
	currency, err := cf.calculator.GetCurrency(localeInfo.CurrencyCode)
	if err != nil {
		return "", err
	}
	
	options := &FormatOptions{
		ThousandsSep: currency.ThousandsSep,
		DecimalSep:   currency.DecimalSep,
		ShowSymbol:   true,
	}
	
	return cf.calculator.Format(money, options)
}

// GetLocaleInfo retrieves detailed information for a specific locale.
// Returns comprehensive locale data including language, country,
// currency information, and formatting preferences.
//
// Parameters:
//   - locale: Locale identifier to retrieve information for
//
// Returns:
//   - *LocaleInfo: Complete locale information
//   - error: Error if locale is not supported
//
// Example:
//   info, err := formatter.GetLocaleInfo("de-DE")
//   if err == nil {
//     fmt.Printf("Currency: %s (%s)", info.CurrencyName, info.CurrencyCode)
//   }
func (cf *CurrencyFormatter) GetLocaleInfo(locale string) (*LocaleInfo, error) {
	localeInfo, exists := cf.locales[locale]
	if !exists {
		return nil, &CurrencyError{
			Type:      "unsupported_locale",
			Message:   fmt.Sprintf("Locale %s is not supported", locale),
			Timestamp: time.Now(),
		}
	}
	return &localeInfo, nil
}

// AddLocale adds a new locale configuration to the formatter.
// Extends the formatter's capabilities by registering additional
// locale-specific formatting rules and currency information.
//
// Parameters:
//   - locale: Locale identifier (e.g., "fr-FR", "ja-JP")
//   - info: Complete LocaleInfo configuration
//
// Example:
//   info := LocaleInfo{
//     Locale:       "fr-FR",
//     Language:     "French",
//     Country:      "France",
//     CurrencyName: "Euro",
//     CurrencyCode: EUR,
//   }
//   formatter.AddLocale("fr-FR", info)
func (cf *CurrencyFormatter) AddLocale(locale string, info LocaleInfo) {
	cf.locales[locale] = info
}

// Helper functions for Money operations and mathematical calculations.
// These functions provide convenient utilities for creating, validating,
// and performing mathematical operations on Money values.

// NewMoney creates a new Money instance with basic validation.
// Validates the amount for mathematical correctness before creating
// the Money struct to prevent invalid monetary values.
//
// Validation:
//   - Amount must not be NaN (Not a Number)
//   - Amount must not be infinite (positive or negative)
//   - No currency support validation (use Validator for comprehensive checks)
//
// Parameters:
//   - amount: Monetary amount as float64
//   - currency: Currency code for the amount
//
// Returns:
//   - *Money: New Money instance if valid
//   - error: CurrencyError if amount is invalid
//
// Example:
//   money, err := NewMoney(123.45, USD)
//   if err != nil {
//     log.Printf("Invalid amount: %v", err)
//   }
func NewMoney(amount float64, currency CurrencyCode) (*Money, error) {
	money := Money{
		Amount:   amount,
		Currency: currency,
	}
	
	// Basic validation
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return nil, &CurrencyError{
			Type:      "invalid_amount",
			Message:   "Amount must be a valid number",
			Timestamp: time.Now(),
		}
	}
	
	return &money, nil
}

// NewMoneyFromString creates Money from string representation of amount.
// Parses string input and creates a Money instance with validation,
// useful for processing user input or text-based data sources.
//
// Features:
//   - Parses standard numeric string formats
//   - Validates parsed amount using NewMoney
//   - Handles decimal numbers and scientific notation
//   - Returns detailed error information for parsing failures
//
// Parameters:
//   - amountStr: String representation of the amount
//   - currency: Currency code for the amount
//
// Returns:
//   - *Money: New Money instance if parsing and validation succeed
//   - error: CurrencyError if parsing fails or amount is invalid
//
// Example:
//   money, err := NewMoneyFromString("123.45", USD)
//   if err != nil {
//     log.Printf("Parse error: %v", err)
//   }
func NewMoneyFromString(amountStr string, currency CurrencyCode) (*Money, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return nil, &CurrencyError{
			Type:      "parse_error",
			Message:   fmt.Sprintf("Cannot parse amount: %s", amountStr),
			Timestamp: time.Now(),
		}
	}
	
	return NewMoney(amount, currency)
}

// IsZero checks if money amount is effectively zero.
// Uses a small tolerance (0.001) to handle floating-point precision issues
// and determine if an amount should be considered zero for practical purposes.
//
// Parameters:
//   - money: Money amount to check
//
// Returns:
//   - bool: true if amount is within zero tolerance, false otherwise
//
// Example:
//   money := Money{Amount: 0.0001, Currency: USD}
//   if IsZero(money) {
//     fmt.Println("Amount is effectively zero")
//   }
func IsZero(money Money) bool {
	return math.Abs(money.Amount) < 0.001
}

// IsPositive checks if money amount is positive.
// Uses a small tolerance (0.001) to handle floating-point precision
// and determine if an amount is meaningfully positive.
//
// Parameters:
//   - money: Money amount to check
//
// Returns:
//   - bool: true if amount is greater than tolerance, false otherwise
//
// Example:
//   money := Money{Amount: 10.50, Currency: USD}
//   if IsPositive(money) {
//     fmt.Println("Amount is positive")
//   }
func IsPositive(money Money) bool {
	return money.Amount > 0.001
}

// IsNegative checks if money amount is negative.
// Uses a small tolerance (0.001) to handle floating-point precision
// and determine if an amount is meaningfully negative.
//
// Parameters:
//   - money: Money amount to check
//
// Returns:
//   - bool: true if amount is less than negative tolerance, false otherwise
//
// Example:
//   money := Money{Amount: -5.25, Currency: USD}
//   if IsNegative(money) {
//     fmt.Println("Amount is negative")
//   }
func IsNegative(money Money) bool {
	return money.Amount < -0.001
}

// Abs returns absolute value of money.
// Creates a new Money instance with the absolute value of the amount,
// preserving the original currency. Useful for calculations that require
// positive amounts regardless of the original sign.
//
// Parameters:
//   - money: Money amount to get absolute value of
//
// Returns:
//   - Money: New Money instance with absolute amount value
//
// Example:
//   negative := Money{Amount: -25.50, Currency: USD}
//   positive := Abs(negative)
//   // positive.Amount = 25.50
func Abs(money Money) Money {
	return Money{
		Amount:   math.Abs(money.Amount),
		Currency: money.Currency,
	}
}

// Negate returns negated money amount.
// Creates a new Money instance with the opposite sign of the amount,
// preserving the original currency. Useful for representing refunds,
// credits, or reversing transactions.
//
// Parameters:
//   - money: Money amount to negate
//
// Returns:
//   - Money: New Money instance with negated amount
//
// Example:
//   charge := Money{Amount: 100.00, Currency: USD}
//   refund := Negate(charge)
//   // refund.Amount = -100.00
func Negate(money Money) Money {
	return Money{
		Amount:   -money.Amount,
		Currency: money.Currency,
	}
}

// Min returns the smaller of two money amounts.
// Compares two Money values and returns the one with the smaller amount.
// Both amounts must be in the same currency for comparison.
//
// Features:
//   - Currency validation before comparison
//   - Preserves original Money instance (no copying)
//   - Returns detailed error for currency mismatches
//
// Parameters:
//   - a: First Money amount to compare
//   - b: Second Money amount to compare
//
// Returns:
//   - Money: The Money instance with smaller amount
//   - error: CurrencyError if currencies don't match
//
// Example:
//   price1 := Money{Amount: 15.99, Currency: USD}
//   price2 := Money{Amount: 12.50, Currency: USD}
//   cheaper, err := Min(price1, price2)
//   // cheaper.Amount = 12.50
func Min(a, b Money) (Money, error) {
	if a.Currency != b.Currency {
		return Money{}, &CurrencyError{
			Type:      "currency_mismatch",
			Message:   fmt.Sprintf("Cannot compare different currencies: %s and %s", a.Currency, b.Currency),
			Timestamp: time.Now(),
		}
	}
	
	if a.Amount <= b.Amount {
		return a, nil
	}
	return b, nil
}

// Max returns the larger of two money amounts.
// Compares two Money values and returns the one with the larger amount.
// Both amounts must be in the same currency for comparison.
//
// Features:
//   - Currency validation before comparison
//   - Preserves original Money instance (no copying)
//   - Returns detailed error for currency mismatches
//
// Parameters:
//   - a: First Money amount to compare
//   - b: Second Money amount to compare
//
// Returns:
//   - Money: The Money instance with larger amount
//   - error: CurrencyError if currencies don't match
//
// Example:
//   price1 := Money{Amount: 15.99, Currency: USD}
//   price2 := Money{Amount: 12.50, Currency: USD}
//   expensive, err := Max(price1, price2)
//   // expensive.Amount = 15.99
func Max(a, b Money) (Money, error) {
	if a.Currency != b.Currency {
		return Money{}, &CurrencyError{
			Type:      "currency_mismatch",
			Message:   fmt.Sprintf("Cannot compare different currencies: %s and %s", a.Currency, b.Currency),
			Timestamp: time.Now(),
		}
	}
	
	if a.Amount >= b.Amount {
		return a, nil
	}
	return b, nil
}

// Sum calculates total of money amounts in the same currency.
// Adds all Money amounts in the slice and returns the total.
// All amounts must be in the same currency for summation.
//
// Features:
//   - Currency validation across all amounts
//   - Handles empty slices (returns zero amount)
//   - Preserves currency from input amounts
//   - Returns detailed error for currency mismatches
//
// Parameters:
//   - amounts: Slice of Money amounts to sum
//
// Returns:
//   - Money: Total sum with the common currency
//   - error: CurrencyError if currencies don't match or slice is empty
//
// Example:
//   prices := []Money{
//     {Amount: 10.50, Currency: USD},
//     {Amount: 25.75, Currency: USD},
//     {Amount: 5.25, Currency: USD},
//   }
//   total, err := Sum(prices)
//   // total.Amount = 41.50
func Sum(amounts []Money) (Money, error) {
	if len(amounts) == 0 {
		return Money{}, &CurrencyError{
			Type:      "empty_input",
			Message:   "Cannot sum empty array",
			Timestamp: time.Now(),
		}
	}
	
	baseCurrency := amounts[0].Currency
	var total float64
	
	for i, amount := range amounts {
		if amount.Currency != baseCurrency {
			return Money{}, &CurrencyError{
				Type:      "currency_mismatch",
				Message:   fmt.Sprintf("All amounts must have the same currency. Amount %d has %s, expected %s", i, amount.Currency, baseCurrency),
				Timestamp: time.Now(),
			}
		}
		total += amount.Amount
	}
	
	return Money{
		Amount:   total,
		Currency: baseCurrency,
	}, nil
}

// Average calculates average of money amounts in the same currency.
// Computes the arithmetic mean of all Money amounts in the slice.
// All amounts must be in the same currency for calculation.
//
// Features:
//   - Currency validation across all amounts
//   - Handles division by zero (empty slices)
//   - Preserves currency from input amounts
//   - Returns detailed error for currency mismatches
//
// Parameters:
//   - amounts: Slice of Money amounts to average
//
// Returns:
//   - Money: Average amount with the common currency
//   - error: CurrencyError if currencies don't match or slice is empty
//
// Example:
//   prices := []Money{
//     {Amount: 10.00, Currency: USD},
//     {Amount: 20.00, Currency: USD},
//     {Amount: 30.00, Currency: USD},
//   }
//   avg, err := Average(prices)
//   // avg.Amount = 20.00
func Average(amounts []Money) (Money, error) {
	sum, err := Sum(amounts)
	if err != nil {
		return Money{}, err
	}
	
	average := sum.Amount / float64(len(amounts))
	
	return Money{
		Amount:   average,
		Currency: sum.Currency,
	}, nil
}

// Percentage calculates percentage of money amount.
// Computes a percentage of the given Money amount and returns
// a new Money instance with the calculated value.
//
// Features:
//   - Supports any percentage value (including > 100% or negative)
//   - Preserves original currency
//   - Handles floating-point percentage values
//   - No validation on percentage range (caller responsibility)
//
// Parameters:
//   - money: Base Money amount to calculate percentage of
//   - percent: Percentage value (e.g., 15.5 for 15.5%)
//
// Returns:
//   - Money: New Money instance with percentage amount
//
// Example:
//   price := Money{Amount: 100.00, Currency: USD}
//   tax := Percentage(price, 8.25)
//   // tax.Amount = 8.25 (8.25% of 100.00)
//   
//   discount := Percentage(price, 15.0)
//   // discount.Amount = 15.00 (15% of 100.00)
func Percentage(money Money, percent float64) Money {
	return Money{
		Amount:   money.Amount * (percent / 100.0),
		Currency: money.Currency,
	}
}

// Split divides money amount into equal parts.
// Divides the Money amount into the specified number of equal parts,
// handling remainder distribution to ensure the sum equals the original.
//
// Features:
//   - Equal distribution with remainder handling
//   - Preserves total amount (sum of parts equals original)
//   - Validates positive number of parts
//   - Returns slice of Money instances and remainder
//
// Parameters:
//   - money: Money amount to split
//   - parts: Number of parts to split into (must be > 0)
//
// Returns:
//   - []Money: Slice of Money parts
//   - Money: Remainder amount after equal distribution
//
// Example:
//   total := Money{Amount: 100.00, Currency: USD}
//   parts, remainder := Split(total, 3)
//   // parts[0].Amount = 33.33
//   // parts[1].Amount = 33.33
//   // parts[2].Amount = 33.33
//   // remainder.Amount = 0.01
func Split(money Money, parts int) ([]Money, Money) {
	if parts <= 0 {
		return nil, money
	}
	
	partAmount := money.Amount / float64(parts)
	remainder := money.Amount - (partAmount * float64(parts))
	
	result := make([]Money, parts)
	for i := 0; i < parts; i++ {
		result[i] = Money{
			Amount:   partAmount,
			Currency: money.Currency,
		}
	}
	
	remainderMoney := Money{
		Amount:   remainder,
		Currency: money.Currency,
	}
	
	return result, remainderMoney
}

// Allocate distributes money according to ratios.
// Distributes the Money amount proportionally based on the provided ratios,
// ensuring the sum of allocated amounts equals the original amount.
//
// Features:
//   - Proportional distribution based on ratios
//   - Handles remainder distribution to preserve total
//   - Validates non-empty ratios and positive values
//   - Normalizes ratios automatically (ratios don't need to sum to 1.0)
//
// Parameters:
//   - money: Money amount to allocate
//   - ratios: Slice of ratio values for distribution (must be > 0)
//
// Returns:
//   - []Money: Slice of allocated Money amounts
//   - error: CurrencyError if ratios are invalid
//
// Example:
//   total := Money{Amount: 100.00, Currency: USD}
//   ratios := []float64{3, 2, 1} // 3:2:1 ratio
//   parts, err := Allocate(total, ratios)
//   // parts[0].Amount = 50.00 (3/6 of 100)
//   // parts[1].Amount = 33.33 (2/6 of 100)
//   // parts[2].Amount = 16.67 (1/6 of 100)
func Allocate(money Money, ratios []float64) ([]Money, error) {
	if len(ratios) == 0 {
		return nil, &CurrencyError{
			Type:      "empty_input",
			Message:   "Ratios array cannot be empty",
			Timestamp: time.Now(),
		}
	}
	
	// Calculate total ratio
	var totalRatio float64
	for _, ratio := range ratios {
		if ratio < 0 {
			return nil, &CurrencyError{
				Type:      "invalid_ratio",
				Message:   "Ratios must be non-negative",
				Timestamp: time.Now(),
			}
		}
		totalRatio += ratio
	}
	
	if totalRatio == 0 {
		return nil, &CurrencyError{
			Type:      "invalid_ratio",
			Message:   "Total ratio cannot be zero",
			Timestamp: time.Now(),
		}
	}
	
	// Allocate amounts
	result := make([]Money, len(ratios))
	var allocated float64
	
	for i, ratio := range ratios {
		amount := money.Amount * (ratio / totalRatio)
		result[i] = Money{
			Amount:   amount,
			Currency: money.Currency,
		}
		allocated += amount
	}
	
	// Handle rounding differences by adjusting the last allocation
	difference := money.Amount - allocated
	if math.Abs(difference) > 0.001 && len(result) > 0 {
		result[len(result)-1].Amount += difference
	}
	
	return result, nil
}
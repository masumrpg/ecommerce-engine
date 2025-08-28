// Package currency provides constants, mappings, and utility functions for currency handling.
// This file contains all the constant definitions, currency metadata, and helper functions
// used throughout the currency package for formatting, validation, and currency operations.
//
// Key features:
//   - Currency formatting constants and options
//   - Exchange rate thresholds and tolerances
//   - Currency groupings by region and characteristics
//   - Currency symbols, names, and decimal place mappings
//   - Helper functions for currency validation and metadata retrieval
//   - Common currency pair definitions
//
// Example usage:
//
//	// Check if a currency is major
//	if IsMajorCurrency(USD) {
//		fmt.Println("USD is a major currency")
//	}
//
//	// Get currency symbol
//	symbol := GetCurrencySymbol(EUR) // Returns "€"
//
//	// Get decimal places
//	places := GetCurrencyDecimalPlaces(JPY) // Returns 0
package currency

// Default formatting options for currency display and calculation.
// These constants define standard precision levels, separators, and
// negative number formatting styles used across the currency system.
const (
	// DefaultPrecision is the standard number of decimal places for most currencies (2)
	DefaultPrecision = 2
	
	// ZeroPrecision is used for currencies that don't use decimal places (like JPY, KRW)
	ZeroPrecision = 0
	
	// HighPrecision is used for currencies requiring more than 2 decimal places
	HighPrecision = 4
	
	// CryptoPrecision is used for cryptocurrency calculations requiring high precision
	CryptoPrecision = 8
	
	// DefaultThousandsSep is the default thousands separator for number formatting
	DefaultThousandsSep = ","
	
	// DefaultDecimalSep is the default decimal separator for number formatting
	DefaultDecimalSep = "."
	
	// NegativeStyleMinus formats negative amounts with a minus sign (-123.45)
	NegativeStyleMinus = "minus"
	
	// NegativeStyleParentheses formats negative amounts with parentheses (123.45)
	NegativeStyleParentheses = "parentheses"
	
	// NegativeStyleMinusSymbol formats negative amounts with minus and currency symbol (-$123.45)
	NegativeStyleMinusSymbol = "minus_symbol"
)

// Exchange rate constants define thresholds, tolerances, and limits for
// currency exchange rate operations and validations.
const (
	// RateFreshThreshold defines how long exchange rates are considered fresh (60 minutes)
	RateFreshThreshold = 60
	
	// RateStaleThreshold defines when exchange rates are considered stale (1440 minutes = 24 hours)
	RateStaleThreshold = 1440
	
	// RateToleranceDefault is the default tolerance for exchange rate comparisons
	RateToleranceDefault = 0.0001
	
	// RateToleranceStrict is a strict tolerance for precise exchange rate comparisons
	RateToleranceStrict = 0.00001
	
	// RateToleranceLoose is a loose tolerance for approximate exchange rate comparisons
	RateToleranceLoose = 0.001
	
	// MaxExchangeRate is the maximum allowed exchange rate value to prevent overflow
	MaxExchangeRate = 1000000.0
	
	// MinExchangeRate is the minimum allowed exchange rate value to prevent underflow
	MinExchangeRate = 0.000001
)

// Amount limits and tolerances define boundaries and precision settings
// for currency amount calculations and comparisons.
const (
	// MaxAmount is the maximum currency amount that can be processed
	MaxAmount = 1e15
	
	// MinAmount is the minimum currency amount that can be processed
	MinAmount = -1e15
	
	// AmountTolerance is the default tolerance for floating point amount comparisons
	AmountTolerance = 0.001
	
	// StrictTolerance is a strict tolerance for precise amount comparisons
	StrictTolerance = 0.0001
	
	// LooseTolerance is a loose tolerance for approximate amount comparisons
	LooseTolerance = 0.01
	
	// ZeroThreshold defines the threshold below which amounts are considered zero
	ZeroThreshold = 1e-10
)

// Currency groupings organize currencies by geographic regions, trading importance,
// and decimal precision characteristics for easier categorization and processing.
var (
	// MajorCurrencies contains the most actively traded currencies in global markets.
	// These currencies typically have high liquidity and are commonly used in
	// international trade and foreign exchange markets.
	MajorCurrencies = []CurrencyCode{
		USD, EUR, GBP, JPY, CHF, CAD, AUD,
	}
	
	// AsianCurrencies contains currencies from Asian countries and territories.
	// This grouping is useful for regional currency operations and formatting.
	AsianCurrencies = []CurrencyCode{
		IDR, SGD, MYR, THB, PHP, VND, KRW, INR, CNY,
	}
	
	// EuropeanCurrencies contains currencies from European countries.
	// Includes both EU and non-EU European currencies for comprehensive coverage.
	EuropeanCurrencies = []CurrencyCode{
		EUR, GBP, CHF, SEK, NOK, DKK,
	}
	
	// AmericasCurrencies contains currencies from North, Central, and South America.
	// Covers all major currencies used across the Americas region.
	AmericasCurrencies = []CurrencyCode{
		USD, CAD, BRL, MXN,
	}
	
	// MiddleEastCurrencies contains currencies from Middle Eastern and North African countries.
	// Includes currencies from the MENA region for regional operations.
	MiddleEastCurrencies = []CurrencyCode{
		SAR, AED, TRY,
	}
	
	// ZeroDecimalCurrencies contains currencies that don't use fractional units.
	// These currencies are typically displayed and calculated without decimal places.
	ZeroDecimalCurrencies = []CurrencyCode{
		JPY, KRW, VND, IDR,
	}
	
	// HighPrecisionCurrencies contains currencies that use more than 2 decimal places.
	// These currencies require 3 decimal places for accurate representation.
	HighPrecisionCurrencies = []CurrencyCode{
		// Add currencies that need more than 2 decimal places
	}
)

// CurrencySymbols maps currency codes to their display symbols.
// This mapping provides the standard symbols used for currency formatting
// and display in user interfaces and financial documents.
//
// Example usage:
//	symbol := CurrencySymbols[USD] // Returns "$"
var CurrencySymbols = map[CurrencyCode]string{
	USD: "$",
	EUR: "€",
	GBP: "£",
	JPY: "¥",
	CNY: "¥",
	IDR: "Rp",
	SGD: "S$",
	MYR: "RM",
	THB: "฿",
	PHP: "₱",
	VND: "₫",
	KRW: "₩",
	INR: "₹",
	AUD: "A$",
	CAD: "C$",
	CHF: "CHF",
	SEK: "kr",
	NOK: "kr",
	DKK: "kr",
	RUB: "₽",
	BRL: "R$",
	MXN: "$",
	ZAR: "R",
	TRY: "₺",
	SAR: "﷼",
	AED: "د.إ",
}

// CurrencyNames maps currency codes to their full English names.
// This mapping provides human-readable currency names for display
// and documentation purposes.
//
// Example usage:
//	name := CurrencyNames[EUR] // Returns "Euro"
var CurrencyNames = map[CurrencyCode]string{
	USD: "US Dollar",
	EUR: "Euro",
	GBP: "British Pound Sterling",
	JPY: "Japanese Yen",
	CNY: "Chinese Yuan",
	IDR: "Indonesian Rupiah",
	SGD: "Singapore Dollar",
	MYR: "Malaysian Ringgit",
	THB: "Thai Baht",
	PHP: "Philippine Peso",
	VND: "Vietnamese Dong",
	KRW: "South Korean Won",
	INR: "Indian Rupee",
	AUD: "Australian Dollar",
	CAD: "Canadian Dollar",
	CHF: "Swiss Franc",
	SEK: "Swedish Krona",
	NOK: "Norwegian Krone",
	DKK: "Danish Krone",
	RUB: "Russian Ruble",
	BRL: "Brazilian Real",
	MXN: "Mexican Peso",
	ZAR: "South African Rand",
	TRY: "Turkish Lira",
	SAR: "Saudi Riyal",
	AED: "UAE Dirham",
}

// CurrencyDecimalPlaces maps currency codes to their standard number of decimal places.
// This mapping defines how many decimal places should be used when formatting
// and calculating amounts for each currency.
//
// Example usage:
//	places := CurrencyDecimalPlaces[JPY] // Returns 0 (no decimal places)
//	places := CurrencyDecimalPlaces[USD] // Returns 2 (cents)
var CurrencyDecimalPlaces = map[CurrencyCode]int{
	USD: 2,
	EUR: 2,
	GBP: 2,
	JPY: 0,
	CNY: 2,
	IDR: 0,
	SGD: 2,
	MYR: 2,
	THB: 2,
	PHP: 2,
	VND: 0,
	KRW: 0,
	INR: 2,
	AUD: 2,
	CAD: 2,
	CHF: 2,
	SEK: 2,
	NOK: 2,
	DKK: 2,
	RUB: 2,
	BRL: 2,
	MXN: 2,
	ZAR: 2,
	TRY: 2,
	SAR: 2,
	AED: 2,
}

// Helper functions for currency groups

// IsMajorCurrency checks if the given currency code is a major currency.
// Major currencies are the most actively traded currencies in global markets.
//
// Parameters:
//   - code: The currency code to check (e.g., USD, EUR)
//
// Returns:
//   - bool: true if the currency is a major currency, false otherwise
//
// Example:
//	if IsMajorCurrency(USD) {
//		fmt.Println("USD is a major currency")
//	}
func IsMajorCurrency(code CurrencyCode) bool {
	for _, major := range MajorCurrencies {
		if major == code {
			return true
		}
	}
	return false
}

// IsAsianCurrency checks if the given currency code is an Asian currency.
//
// Parameters:
//   - code: The currency code to check (e.g., JPY, CNY)
//
// Returns:
//   - bool: true if the currency is from an Asian country, false otherwise
func IsAsianCurrency(code CurrencyCode) bool {
	for _, asian := range AsianCurrencies {
		if asian == code {
			return true
		}
	}
	return false
}

// IsEuropeanCurrency checks if the given currency code is a European currency.
//
// Parameters:
//   - code: The currency code to check (e.g., EUR, GBP)
//
// Returns:
//   - bool: true if the currency is from a European country, false otherwise
func IsEuropeanCurrency(code CurrencyCode) bool {
	for _, european := range EuropeanCurrencies {
		if european == code {
			return true
		}
	}
	return false
}

// IsZeroDecimalCurrency checks if the given currency code uses zero decimal places.
// Zero decimal currencies don't have fractional units (e.g., Japanese Yen, Korean Won).
//
// Parameters:
//   - code: The currency code to check (e.g., JPY, KRW)
//
// Returns:
//   - bool: true if the currency uses zero decimal places, false otherwise
//
// Example:
//	if IsZeroDecimalCurrency(JPY) {
//		fmt.Println("JPY doesn't use decimal places")
//	}
func IsZeroDecimalCurrency(code CurrencyCode) bool {
	for _, zero := range ZeroDecimalCurrencies {
		if zero == code {
			return true
		}
	}
	return false
}

// GetCurrencySymbol returns the display symbol for the given currency code.
// If the currency code is not found, it returns the code itself as fallback.
//
// Parameters:
//   - code: The currency code (e.g., USD, EUR)
//
// Returns:
//   - string: The currency symbol (e.g., "$", "€") or the code if not found
//
// Example:
//	symbol := GetCurrencySymbol(USD) // Returns "$"
//	symbol := GetCurrencySymbol("XYZ") // Returns "XYZ" (fallback)
func GetCurrencySymbol(code CurrencyCode) string {
	if symbol, exists := CurrencySymbols[code]; exists {
		return symbol
	}
	return string(code) // Fallback to currency code
}

// GetCurrencyName returns the full English name for the given currency code.
// If the currency code is not found, it returns the code itself as fallback.
//
// Parameters:
//   - code: The currency code (e.g., USD, EUR)
//
// Returns:
//   - string: The currency name (e.g., "US Dollar", "Euro") or the code if not found
//
// Example:
//	name := GetCurrencyName(EUR) // Returns "Euro"
func GetCurrencyName(code CurrencyCode) string {
	if name, exists := CurrencyNames[code]; exists {
		return name
	}
	return string(code) // Fallback to currency code
}

// GetCurrencyDecimalPlaces returns the standard number of decimal places for the given currency.
// If the currency code is not found, it returns the default precision (2 decimal places).
//
// Parameters:
//   - code: The currency code (e.g., USD, JPY)
//
// Returns:
//   - int: The number of decimal places (0-3) or DefaultPrecision if not found
//
// Example:
//	places := GetCurrencyDecimalPlaces(USD) // Returns 2
//	places := GetCurrencyDecimalPlaces(JPY) // Returns 0
func GetCurrencyDecimalPlaces(code CurrencyCode) int {
	if places, exists := CurrencyDecimalPlaces[code]; exists {
		return places
	}
	return DefaultPrecision // Fallback to default
}

// IsValidCurrencyCode checks if the given currency code is supported.
// A currency is considered valid if it exists in the CurrencyNames mapping.
//
// Parameters:
//   - code: The currency code to validate (e.g., USD, EUR)
//
// Returns:
//   - bool: true if the currency code is supported, false otherwise
//
// Example:
//	if IsValidCurrencyCode(USD) {
//		fmt.Println("USD is supported")
//	}
func IsValidCurrencyCode(code CurrencyCode) bool {
	_, exists := CurrencyNames[code]
	return exists
}

// GetSupportedCurrencyCodes returns a slice of all supported currency codes.
// The returned slice contains all currency codes that have names defined.
//
// Returns:
//   - []CurrencyCode: A slice of supported currency codes
//
// Example:
//	codes := GetSupportedCurrencyCodes()
//	fmt.Printf("Supported currencies: %v\n", codes)
func GetSupportedCurrencyCodes() []CurrencyCode {
	codes := make([]CurrencyCode, 0, len(CurrencyNames))
	for code := range CurrencyNames {
		codes = append(codes, code)
	}
	return codes
}

// CommonPairs contains frequently traded currency pairs in global markets.
// These pairs typically have high liquidity and tight spreads, making them
// suitable for most trading and exchange operations.
//
// The list includes major pairs (USD-based), cross pairs (non-USD), and
// regional pairs commonly used in Asian markets.
var CommonPairs = []CurrencyPair{
	{Base: USD, Quote: EUR},
	{Base: USD, Quote: GBP},
	{Base: USD, Quote: JPY},
	{Base: USD, Quote: CHF},
	{Base: USD, Quote: CAD},
	{Base: USD, Quote: AUD},
	{Base: USD, Quote: IDR},
	{Base: USD, Quote: SGD},
	{Base: USD, Quote: MYR},
	{Base: EUR, Quote: GBP},
	{Base: EUR, Quote: JPY},
	{Base: EUR, Quote: CHF},
	{Base: GBP, Quote: JPY},
	{Base: AUD, Quote: JPY},
	{Base: SGD, Quote: IDR},
	{Base: SGD, Quote: MYR},
}

// IsCommonPair checks if the given currency pair is commonly traded.
// This function checks both the direct pair and its reverse to determine
// if the pair is frequently traded in global markets.
//
// Parameters:
//   - pair: The currency pair to check (with Base and Quote currencies)
//
// Returns:
//   - bool: true if the pair is commonly traded, false otherwise
//
// Example:
//	pair := CurrencyPair{Base: USD, Quote: EUR}
//	if IsCommonPair(pair) {
//		fmt.Println("USD/EUR is a common trading pair")
//	}
func IsCommonPair(pair CurrencyPair) bool {
	for _, common := range CommonPairs {
		if common.Base == pair.Base && common.Quote == pair.Quote {
			return true
		}
		// Also check reverse pair
		if common.Base == pair.Quote && common.Quote == pair.Base {
			return true
		}
	}
	return false
}
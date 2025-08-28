// Package currency provides comprehensive currency handling capabilities for e-commerce applications.
// It supports multi-currency operations, exchange rate management, formatting, arithmetic operations,
// and currency conversion with proper rounding and validation.
//
// Key Features:
//   - ISO 4217 currency code support with 25+ predefined currencies
//   - Multi-currency arithmetic operations (add, subtract, multiply, divide)
//   - Currency conversion with exchange rate management
//   - Flexible formatting with locale-specific rules
//   - Multiple rounding modes for precise calculations
//   - Currency comparison and validation
//   - Exchange rate provider abstraction
//   - Thread-safe operations with proper error handling
//
// Basic Usage:
//   // Create money amounts
//   usd := Money{Amount: 100.50, Currency: USD}
//   eur := Money{Amount: 85.25, Currency: EUR}
//
//   // Currency conversion
//   calc := NewCalculator()
//   calc.SetExchangeRate(USD, EUR, 0.85, "ECB")
//   result, err := calc.Convert(ConversionInput{
//     Amount: 100.0,
//     From:   USD,
//     To:     EUR,
//   })
//
//   // Arithmetic operations
//   sum, err := calc.Add(usd1, usd2)
//   product, err := calc.Multiply(usd, 1.5)
//
//   // Currency formatting
//   formatted := calc.Format(usd, FormatOptions{
//     ShowSymbol: true,
//     Precision:  &[]int{2}[0],
//   })
//   // Output: "$100.50"
package currency

import "time"

// CurrencyCode represents ISO 4217 standard three-letter currency codes.
// Used throughout the system to identify currencies in a standardized format.
// All currency codes follow the ISO 4217 international standard.
//
// Example:
//   var code CurrencyCode = USD
//   fmt.Println(string(code)) // "USD"
type CurrencyCode string

// Predefined ISO 4217 currency codes for major world currencies.
// These constants provide type-safe currency identification and cover
// the most commonly used currencies in global e-commerce.
//
// Regional Coverage:
//   - North America: USD, CAD, MXN
//   - Europe: EUR, GBP, CHF, SEK, NOK, DKK, RUB, TRY
//   - Asia-Pacific: JPY, CNY, SGD, MYR, THB, PHP, VND, KRW, INR, IDR, AUD
//   - Middle East & Africa: SAR, AED, ZAR
//   - South America: BRL
const (
	USD CurrencyCode = "USD" // US Dollar - Primary global reserve currency
	EUR CurrencyCode = "EUR" // Euro - European Union currency
	GBP CurrencyCode = "GBP" // British Pound Sterling
	JPY CurrencyCode = "JPY" // Japanese Yen - No decimal places
	CNY CurrencyCode = "CNY" // Chinese Yuan Renminbi
	IDR CurrencyCode = "IDR" // Indonesian Rupiah - No decimal places
	SGD CurrencyCode = "SGD" // Singapore Dollar
	MYR CurrencyCode = "MYR" // Malaysian Ringgit
	THB CurrencyCode = "THB" // Thai Baht
	PHP CurrencyCode = "PHP" // Philippine Peso
	VND CurrencyCode = "VND" // Vietnamese Dong - No decimal places
	KRW CurrencyCode = "KRW" // South Korean Won - No decimal places
	INR CurrencyCode = "INR" // Indian Rupee
	AUD CurrencyCode = "AUD" // Australian Dollar
	CAD CurrencyCode = "CAD" // Canadian Dollar
	CHF CurrencyCode = "CHF" // Swiss Franc
	SEK CurrencyCode = "SEK" // Swedish Krona
	NOK CurrencyCode = "NOK" // Norwegian Krone
	DKK CurrencyCode = "DKK" // Danish Krone
	RUB CurrencyCode = "RUB" // Russian Ruble
	BRL CurrencyCode = "BRL" // Brazilian Real
	MXN CurrencyCode = "MXN" // Mexican Peso
	ZAR CurrencyCode = "ZAR" // South African Rand
	TRY CurrencyCode = "TRY" // Turkish Lira
	SAR CurrencyCode = "SAR" // Saudi Riyal
	AED CurrencyCode = "AED" // UAE Dirham
)

// Currency represents a complete currency definition with formatting rules.
// Contains all necessary information for currency display, calculation, and formatting
// according to locale-specific conventions and international standards.
//
// Fields:
//   - Code: ISO 4217 three-letter currency code
//   - Name: Full currency name (e.g., "US Dollar")
//   - Symbol: Currency symbol (e.g., "$", "€", "¥")
//   - DecimalPlaces: Number of decimal places for the currency (0-4)
//   - ThousandsSep: Thousands separator character (",", ".", " ")
//   - DecimalSep: Decimal separator character (".", ",")
//   - SymbolFirst: Whether symbol appears before amount
//   - SpaceBetween: Whether space appears between symbol and amount
//
// Example:
//   usd := Currency{
//     Code:          USD,
//     Name:          "US Dollar",
//     Symbol:        "$",
//     DecimalPlaces: 2,
//     ThousandsSep:  ",",
//     DecimalSep:    ".",
//     SymbolFirst:   true,
//     SpaceBetween:  false,
//   }
type Currency struct {
	Code         CurrencyCode `json:"code"`
	Name         string       `json:"name"`
	Symbol       string       `json:"symbol"`
	DecimalPlaces int         `json:"decimal_places"`
	ThousandsSep string       `json:"thousands_separator"`
	DecimalSep   string       `json:"decimal_separator"`
	SymbolFirst  bool         `json:"symbol_first"`
	SpaceBetween bool         `json:"space_between"`
}

// Money represents a monetary amount in a specific currency.
// The fundamental type for all currency operations, ensuring type safety
// and preventing accidental mixing of different currencies.
//
// Fields:
//   - Amount: Numeric value of the money (supports negative values)
//   - Currency: ISO 4217 currency code identifying the currency
//
// Usage Guidelines:
//   - Always specify both amount and currency
//   - Use consistent decimal precision for calculations
//   - Validate currency compatibility before operations
//
// Example:
//   price := Money{Amount: 99.99, Currency: USD}
//   discount := Money{Amount: -10.00, Currency: USD}
//   total := Money{Amount: 1234.56, Currency: EUR}
type Money struct {
	Amount   float64      `json:"amount"`
	Currency CurrencyCode `json:"currency"`
}

// ExchangeRate represents the exchange rate between two currencies.
// Contains rate information with metadata for tracking and auditing purposes.
// Supports bidirectional conversion and rate source attribution.
//
// Fields:
//   - From: Source currency code
//   - To: Target currency code
//   - Rate: Exchange rate multiplier (From * Rate = To)
//   - Timestamp: When the rate was set or last updated
//   - Source: Rate provider or source identifier
//
// Rate Calculation:
//   - 1 unit of From currency = Rate units of To currency
//   - Example: USD/EUR rate of 0.85 means 1 USD = 0.85 EUR
//
// Example:
//   rate := ExchangeRate{
//     From:      USD,
//     To:        EUR,
//     Rate:      0.8542,
//     Timestamp: time.Now(),
//     Source:    "ECB",
//   }
type ExchangeRate struct {
	From      CurrencyCode `json:"from"`
	To        CurrencyCode `json:"to"`
	Rate      float64      `json:"rate"`
	Timestamp time.Time    `json:"timestamp"`
	Source    string       `json:"source"`
}

// ConversionInput represents input parameters for currency conversion.
// Specifies the conversion requirements including amount, currency pair,
// and optional historical rate date for time-specific conversions.
//
// Fields:
//   - Amount: Numeric amount to convert
//   - From: Source currency code
//   - To: Target currency code
//   - RateDate: Optional specific date for historical rates
//
// Usage:
//   - For current rates: leave RateDate as nil
//   - For historical rates: specify exact date
//   - Amount can be negative for refunds or adjustments
//
// Example:
//   input := ConversionInput{
//     Amount: 100.50,
//     From:   USD,
//     To:     EUR,
//     // RateDate: nil for current rate
//   }
type ConversionInput struct {
	Amount   float64      `json:"amount"`
	From     CurrencyCode `json:"from"`
	To       CurrencyCode `json:"to"`
	RateDate *time.Time   `json:"rate_date,omitempty"`
}

// ConversionResult represents the complete result of currency conversion.
// Provides detailed information about the conversion including original amount,
// converted amount, exchange rate used, and conversion timestamp.
//
// Fields:
//   - OriginalAmount: Input money amount before conversion
//   - ConvertedAmount: Output money amount after conversion
//   - ExchangeRate: Exchange rate used for the conversion
//   - ConvertedAt: Timestamp when conversion was performed
//
// Features:
//   - Full audit trail with timestamps
//   - Exchange rate transparency
//   - Proper rounding according to target currency
//   - Immutable result for record keeping
//
// Example:
//   result := ConversionResult{
//     OriginalAmount:  Money{Amount: 100.00, Currency: USD},
//     ConvertedAmount: Money{Amount: 85.42, Currency: EUR},
//     ExchangeRate:    ExchangeRate{From: USD, To: EUR, Rate: 0.8542},
//     ConvertedAt:     time.Now(),
//   }
type ConversionResult struct {
	OriginalAmount Money        `json:"original_amount"`
	ConvertedAmount Money       `json:"converted_amount"`
	ExchangeRate   ExchangeRate `json:"exchange_rate"`
	ConvertedAt    time.Time    `json:"converted_at"`
}

// FormatOptions represents customizable options for currency formatting.
// Allows fine-grained control over currency display appearance,
// overriding default currency formatting rules when specified.
//
// Fields:
//   - ShowSymbol: Whether to display currency symbol (e.g., "$")
//   - ShowCode: Whether to display currency code (e.g., "USD")
//   - Precision: Override decimal places (nil uses currency default)
//   - ThousandsSep: Override thousands separator
//   - DecimalSep: Override decimal separator
//   - SymbolFirst: Override symbol position (nil uses currency default)
//   - SpaceBetween: Override spacing (nil uses currency default)
//   - NegativeStyle: How to display negative amounts
//
// Negative Styles:
//   - "parentheses": ($100.00)
//   - "minus": -$100.00
//   - "minus_symbol": -$100.00
//
// Example:
//   opts := FormatOptions{
//     ShowSymbol:    true,
//     ShowCode:      false,
//     Precision:     &[]int{2}[0],
//     NegativeStyle: "parentheses",
//   }
type FormatOptions struct {
	ShowSymbol    bool   `json:"show_symbol"`
	ShowCode      bool   `json:"show_code"`
	Precision     *int   `json:"precision,omitempty"`
	ThousandsSep  string `json:"thousands_separator,omitempty"`
	DecimalSep    string `json:"decimal_separator,omitempty"`
	SymbolFirst   *bool  `json:"symbol_first,omitempty"`
	SpaceBetween  *bool  `json:"space_between,omitempty"`
	NegativeStyle string `json:"negative_style,omitempty"` // "parentheses", "minus", "minus_symbol"
}

// RoundingMode represents different rounding strategies for currency calculations.
// Provides precise control over how fractional currency amounts are rounded
// to match currency-specific decimal place requirements.
//
// Available Modes:
//   - RoundingModeHalfUp: Round 0.5 up (most common, default)
//   - RoundingModeHalfDown: Round 0.5 down
//   - RoundingModeHalfEven: Banker's rounding (round to nearest even)
//   - RoundingModeUp: Always round up (ceiling)
//   - RoundingModeDown: Always round down (floor)
//   - RoundingModeTruncate: Simply truncate decimal places
//
// Examples (rounding 1.235 to 2 decimal places):
//   - HalfUp: 1.24
//   - HalfDown: 1.23
//   - HalfEven: 1.24 (rounds to even)
//   - Up: 1.24
//   - Down: 1.23
//   - Truncate: 1.23
type RoundingMode string

const (
	RoundingModeHalfUp   RoundingMode = "half_up"   // Round 0.5 up (default)
	RoundingModeHalfDown RoundingMode = "half_down" // Round 0.5 down
	RoundingModeHalfEven RoundingMode = "half_even" // Banker's rounding
	RoundingModeUp       RoundingMode = "up"        // Always round up
	RoundingModeDown     RoundingMode = "down"      // Always round down
	RoundingModeTruncate RoundingMode = "truncate"  // Truncate decimals
)

// ArithmeticOperation represents different arithmetic operations for money calculations.
// Defines the type of mathematical operation to perform on monetary amounts
// with proper currency validation and rounding.
//
// Supported Operations:
//   - OperationAdd: Addition of two money amounts (same currency)
//   - OperationSubtract: Subtraction of two money amounts (same currency)
//   - OperationMultiply: Multiplication of money by numeric factor
//   - OperationDivide: Division of money by numeric divisor
//
// Currency Requirements:
//   - Add/Subtract: Both amounts must have same currency
//   - Multiply/Divide: Only first amount needs currency, factor is numeric
type ArithmeticOperation string

const (
	OperationAdd      ArithmeticOperation = "add"      // Addition: money + money
	OperationSubtract ArithmeticOperation = "subtract" // Subtraction: money - money
	OperationMultiply ArithmeticOperation = "multiply" // Multiplication: money * factor
	OperationDivide   ArithmeticOperation = "divide"   // Division: money / divisor
)

// ArithmeticInput represents input parameters for arithmetic operations.
// Specifies the operands, operation type, and optional rounding mode
// for performing calculations on monetary amounts.
//
// Fields:
//   - Amount1: First monetary amount (primary operand)
//   - Amount2: Second monetary amount (secondary operand for binary ops)
//   - Operation: Type of arithmetic operation to perform
//   - Rounding: Optional rounding mode (uses default if omitted)
//
// Operation-Specific Usage:
//   - Add/Subtract: Both Amount1 and Amount2 must have same currency
//   - Multiply: Amount1 is money, Amount2.Amount is the multiplier
//   - Divide: Amount1 is money, Amount2.Amount is the divisor
//
// Example:
//   input := ArithmeticInput{
//     Amount1:   Money{Amount: 100.50, Currency: USD},
//     Amount2:   Money{Amount: 25.25, Currency: USD},
//     Operation: OperationAdd,
//     Rounding:  RoundingModeHalfUp,
//   }
type ArithmeticInput struct {
	Amount1   Money               `json:"amount1"`
	Amount2   Money               `json:"amount2"`
	Operation ArithmeticOperation `json:"operation"`
	Rounding  RoundingMode        `json:"rounding,omitempty"`
}

// ArithmeticResult represents the complete result of arithmetic operations.
// Provides detailed information about the calculation including result,
// operation performed, operands used, and calculation timestamp.
//
// Fields:
//   - Result: Final monetary amount after operation
//   - Operation: Type of arithmetic operation performed
//   - Operands: List of input amounts used in calculation
//   - CalculatedAt: Timestamp when calculation was performed
//
// Features:
//   - Complete audit trail for calculations
//   - Proper rounding according to result currency
//   - Immutable result for record keeping
//   - Operation transparency for debugging
//
// Example:
//   result := ArithmeticResult{
//     Result:       Money{Amount: 125.75, Currency: USD},
//     Operation:    OperationAdd,
//     Operands:     []Money{{100.50, USD}, {25.25, USD}},
//     CalculatedAt: time.Now(),
//   }
type ArithmeticResult struct {
	Result      Money               `json:"result"`
	Operation   ArithmeticOperation `json:"operation"`
	Operands    []Money             `json:"operands"`
	CalculatedAt time.Time          `json:"calculated_at"`
}

// ComparisonResult represents the result of comparing two monetary amounts.
// Provides comprehensive comparison information including relationship flags,
// difference calculation, and comparison timestamp.
//
// Fields:
//   - Amount1: First money amount in comparison
//   - Amount2: Second money amount in comparison
//   - IsEqual: Whether amounts are exactly equal
//   - IsGreater: Whether Amount1 > Amount2
//   - IsLess: Whether Amount1 < Amount2
//   - Difference: Absolute difference between amounts
//   - ComparedAt: Timestamp when comparison was performed
//
// Comparison Logic:
//   - Requires same currency for both amounts
//   - Uses precise decimal comparison
//   - Difference is always positive (absolute value)
//
// Example:
//   result := ComparisonResult{
//     Amount1:    Money{Amount: 100.50, Currency: USD},
//     Amount2:    Money{Amount: 75.25, Currency: USD},
//     IsEqual:    false,
//     IsGreater:  true,
//     IsLess:     false,
//     Difference: Money{Amount: 25.25, Currency: USD},
//     ComparedAt: time.Now(),
//   }
type ComparisonResult struct {
	Amount1     Money     `json:"amount1"`
	Amount2     Money     `json:"amount2"`
	IsEqual     bool      `json:"is_equal"`
	IsGreater   bool      `json:"is_greater"`
	IsLess      bool      `json:"is_less"`
	Difference  Money     `json:"difference"`
	ComparedAt  time.Time `json:"compared_at"`
}

// CurrencyList represents a collection of supported currencies.
// Provides a complete inventory of all registered currencies in the system
// with metadata about the collection.
//
// Fields:
//   - Currencies: Array of all supported currency definitions
//   - Total: Count of currencies in the collection
//   - UpdatedAt: Timestamp of last collection update
//
// Usage:
//   - Enumerate all available currencies
//   - Display currency selection options
//   - Validate currency support
//   - Track currency registry changes
//
// Example:
//   list := CurrencyList{
//     Currencies: []Currency{usd, eur, gbp},
//     Total:      3,
//     UpdatedAt:   time.Now(),
//   }
type CurrencyList struct {
	Currencies []Currency `json:"currencies"`
	Total      int        `json:"total"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ValidationError represents field-specific validation errors.
// Provides detailed information about validation failures for
// currency-related data with structured error reporting.
//
// Fields:
//   - Field: Name of the field that failed validation
//   - Message: Human-readable error description
//   - Code: Machine-readable error code for programmatic handling
//
// Common Error Codes:
//   - "required": Field is required but missing
//   - "invalid_format": Field format is incorrect
//   - "out_of_range": Field value is outside valid range
//   - "unsupported": Field value is not supported
//
// Example:
//   err := ValidationError{
//     Field:   "currency_code",
//     Message: "Currency code must be 3 characters",
//     Code:    "invalid_format",
//   }
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// CurrencyError represents comprehensive currency operation errors.
// Implements the error interface and provides detailed error information
// including validation details and contextual information.
//
// Fields:
//   - Type: Category of error ("validation", "conversion", "calculation")
//   - Message: Primary error message
//   - Currency: Related currency code (if applicable)
//   - Validations: Detailed validation errors
//   - Timestamp: When the error occurred
//
// Error Types:
//   - "validation": Data validation failures
//   - "conversion": Currency conversion errors
//   - "calculation": Arithmetic operation errors
//   - "format": Currency formatting errors
//   - "rate": Exchange rate errors
//
// Example:
//   err := &CurrencyError{
//     Type:      "conversion",
//     Message:   "Exchange rate not found",
//     Currency:  USD,
//     Timestamp: time.Now(),
//   }
type CurrencyError struct {
	Type        string             `json:"type"`
	Message     string             `json:"message"`
	Currency    CurrencyCode       `json:"currency,omitempty"`
	Validations []ValidationError  `json:"validations,omitempty"`
	Timestamp   time.Time          `json:"timestamp"`
}

// Error implements the error interface for CurrencyError.
// Returns the primary error message for standard error handling.
func (e *CurrencyError) Error() string {
	return e.Message
}

// LocaleInfo represents locale-specific currency information.
// Provides localization context for currency display and formatting
// according to regional preferences and standards.
//
// Fields:
//   - Locale: IETF language tag (e.g., "en-US", "de-DE")
//   - Language: ISO 639 language code (e.g., "en", "de")
//   - Country: ISO 3166 country code (e.g., "US", "DE")
//   - CurrencyName: Localized currency name
//   - CurrencyCode: ISO 4217 currency code
//
// Usage:
//   - Localized currency display
//   - Regional formatting preferences
//   - Multi-language support
//   - Cultural currency conventions
//
// Example:
//   info := LocaleInfo{
//     Locale:       "en-US",
//     Language:     "en",
//     Country:      "US",
//     CurrencyName: "US Dollar",
//     CurrencyCode: USD,
//   }
type LocaleInfo struct {
	Locale       string `json:"locale"`
	Language     string `json:"language"`
	Country      string `json:"country"`
	CurrencyName string `json:"currency_name"`
	CurrencyCode CurrencyCode `json:"currency_code"`
}

// CurrencyPair represents a currency pair for exchange rate operations.
// Defines the relationship between two currencies for conversion and trading.
// Follows standard financial market conventions for currency pair notation.
//
// Fields:
//   - Base: Base currency (the currency being converted from)
//   - Quote: Quote currency (the currency being converted to)
//
// Convention:
//   - Base/Quote format (e.g., USD/EUR means USD to EUR)
//   - Rate represents how many Quote units equal 1 Base unit
//   - Example: USD/EUR rate of 0.85 means 1 USD = 0.85 EUR
//
// Example:
//   pair := CurrencyPair{Base: USD, Quote: EUR}
//   fmt.Println(pair.String()) // "USD/EUR"
type CurrencyPair struct {
	Base  CurrencyCode `json:"base"`
	Quote CurrencyCode `json:"quote"`
}

// String returns the standard string representation of a currency pair.
// Uses the format "BASE/QUOTE" following financial market conventions.
//
// Returns:
//   - string: Currency pair in "BASE/QUOTE" format
//
// Example:
//   pair := CurrencyPair{Base: USD, Quote: EUR}
//   result := pair.String() // "USD/EUR"
func (cp CurrencyPair) String() string {
	return string(cp.Base) + "/" + string(cp.Quote)
}

// RateProvider represents different sources of exchange rate data.
// Categorizes exchange rate providers by their characteristics and reliability
// for proper rate source attribution and validation.
//
// Provider Types:
//   - ProviderManual: Manually entered rates
//   - ProviderFixed: Fixed rates for testing or specific business rules
//   - ProviderAPI: Real-time rates from external APIs
//   - ProviderCentralBank: Official rates from central banks
//
// Reliability Order (highest to lowest):
//   1. ProviderCentralBank: Official government rates
//   2. ProviderAPI: Real-time market rates
//   3. ProviderFixed: Business-defined rates
//   4. ProviderManual: User-entered rates
type RateProvider string

const (
	ProviderManual      RateProvider = "manual"       // Manually entered rates
	ProviderFixed       RateProvider = "fixed"        // Fixed business rates
	ProviderAPI         RateProvider = "api"          // External API rates
	ProviderCentralBank RateProvider = "central_bank" // Official central bank rates
)

// RateSource represents detailed information about exchange rate sources.
// Provides comprehensive metadata about rate providers including
// authentication, update frequency, and reliability metrics.
//
// Fields:
//   - Provider: Type of rate provider
//   - Name: Human-readable source name
//   - URL: Source API endpoint or website (optional)
//   - APIKey: Authentication key for API access (optional)
//   - UpdateFreq: How often rates are updated (e.g., "hourly", "daily")
//   - Reliability: Source reliability score (0.0 to 1.0)
//
// Reliability Scale:
//   - 0.9-1.0: Central banks, official sources
//   - 0.7-0.9: Major financial data providers
//   - 0.5-0.7: Commercial APIs, market data
//   - 0.0-0.5: Manual entry, fixed rates
//
// Example:
//   source := RateSource{
//     Provider:    ProviderAPI,
//     Name:        "European Central Bank",
//     URL:         "https://api.ecb.europa.eu/rates",
//     UpdateFreq:  "daily",
//     Reliability: 0.95,
//   }
type RateSource struct {
	Provider    RateProvider `json:"provider"`
	Name        string       `json:"name"`
	URL         string       `json:"url,omitempty"`
	APIKey      string       `json:"api_key,omitempty"`
	UpdateFreq  string       `json:"update_frequency"`
	Reliability float64      `json:"reliability"` // 0.0 to 1.0
}
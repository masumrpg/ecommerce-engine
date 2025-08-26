package currency

import "time"

// CurrencyCode represents ISO 4217 currency codes
type CurrencyCode string

// Common currency codes
const (
	USD CurrencyCode = "USD" // US Dollar
	EUR CurrencyCode = "EUR" // Euro
	GBP CurrencyCode = "GBP" // British Pound
	JPY CurrencyCode = "JPY" // Japanese Yen
	CNY CurrencyCode = "CNY" // Chinese Yuan
	IDR CurrencyCode = "IDR" // Indonesian Rupiah
	SGD CurrencyCode = "SGD" // Singapore Dollar
	MYR CurrencyCode = "MYR" // Malaysian Ringgit
	THB CurrencyCode = "THB" // Thai Baht
	PHP CurrencyCode = "PHP" // Philippine Peso
	VND CurrencyCode = "VND" // Vietnamese Dong
	KRW CurrencyCode = "KRW" // South Korean Won
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

// Currency represents a currency with its properties
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

// Money represents an amount of money in a specific currency
type Money struct {
	Amount   float64      `json:"amount"`
	Currency CurrencyCode `json:"currency"`
}

// ExchangeRate represents the exchange rate between two currencies
type ExchangeRate struct {
	From      CurrencyCode `json:"from"`
	To        CurrencyCode `json:"to"`
	Rate      float64      `json:"rate"`
	Timestamp time.Time    `json:"timestamp"`
	Source    string       `json:"source"`
}

// ConversionInput represents input for currency conversion
type ConversionInput struct {
	Amount   float64      `json:"amount"`
	From     CurrencyCode `json:"from"`
	To       CurrencyCode `json:"to"`
	RateDate *time.Time   `json:"rate_date,omitempty"`
}

// ConversionResult represents the result of currency conversion
type ConversionResult struct {
	OriginalAmount Money        `json:"original_amount"`
	ConvertedAmount Money       `json:"converted_amount"`
	ExchangeRate   ExchangeRate `json:"exchange_rate"`
	ConvertedAt    time.Time    `json:"converted_at"`
}

// FormatOptions represents options for formatting currency
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

// RoundingMode represents different rounding modes for currency
type RoundingMode string

const (
	RoundingModeHalfUp   RoundingMode = "half_up"
	RoundingModeHalfDown RoundingMode = "half_down"
	RoundingModeHalfEven RoundingMode = "half_even"
	RoundingModeUp       RoundingMode = "up"
	RoundingModeDown     RoundingMode = "down"
	RoundingModeTruncate RoundingMode = "truncate"
)

// ArithmeticOperation represents different arithmetic operations
type ArithmeticOperation string

const (
	OperationAdd      ArithmeticOperation = "add"
	OperationSubtract ArithmeticOperation = "subtract"
	OperationMultiply ArithmeticOperation = "multiply"
	OperationDivide   ArithmeticOperation = "divide"
)

// ArithmeticInput represents input for arithmetic operations
type ArithmeticInput struct {
	Amount1   Money               `json:"amount1"`
	Amount2   Money               `json:"amount2"`
	Operation ArithmeticOperation `json:"operation"`
	Rounding  RoundingMode        `json:"rounding,omitempty"`
}

// ArithmeticResult represents the result of arithmetic operations
type ArithmeticResult struct {
	Result      Money               `json:"result"`
	Operation   ArithmeticOperation `json:"operation"`
	Operands    []Money             `json:"operands"`
	CalculatedAt time.Time          `json:"calculated_at"`
}

// ComparisonResult represents the result of comparing two money amounts
type ComparisonResult struct {
	Amount1     Money     `json:"amount1"`
	Amount2     Money     `json:"amount2"`
	IsEqual     bool      `json:"is_equal"`
	IsGreater   bool      `json:"is_greater"`
	IsLess      bool      `json:"is_less"`
	Difference  Money     `json:"difference"`
	ComparedAt  time.Time `json:"compared_at"`
}

// CurrencyList represents a list of supported currencies
type CurrencyList struct {
	Currencies []Currency `json:"currencies"`
	Total      int        `json:"total"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ValidationError represents currency validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// CurrencyError represents currency-specific errors
type CurrencyError struct {
	Type        string             `json:"type"`
	Message     string             `json:"message"`
	Currency    CurrencyCode       `json:"currency,omitempty"`
	Validations []ValidationError  `json:"validations,omitempty"`
	Timestamp   time.Time          `json:"timestamp"`
}

// Error implements the error interface
func (e *CurrencyError) Error() string {
	return e.Message
}

// LocaleInfo represents locale-specific currency information
type LocaleInfo struct {
	Locale       string `json:"locale"`
	Language     string `json:"language"`
	Country      string `json:"country"`
	CurrencyName string `json:"currency_name"`
	CurrencyCode CurrencyCode `json:"currency_code"`
}

// CurrencyPair represents a currency pair for exchange rates
type CurrencyPair struct {
	Base  CurrencyCode `json:"base"`
	Quote CurrencyCode `json:"quote"`
}

// String returns the string representation of a currency pair
func (cp CurrencyPair) String() string {
	return string(cp.Base) + "/" + string(cp.Quote)
}

// RateProvider represents different exchange rate providers
type RateProvider string

const (
	ProviderManual    RateProvider = "manual"
	ProviderFixed     RateProvider = "fixed"
	ProviderAPI       RateProvider = "api"
	ProviderCentralBank RateProvider = "central_bank"
)

// RateSource represents the source of exchange rates
type RateSource struct {
	Provider    RateProvider `json:"provider"`
	Name        string       `json:"name"`
	URL         string       `json:"url,omitempty"`
	APIKey      string       `json:"api_key,omitempty"`
	UpdateFreq  string       `json:"update_frequency"`
	Reliability float64      `json:"reliability"` // 0.0 to 1.0
}
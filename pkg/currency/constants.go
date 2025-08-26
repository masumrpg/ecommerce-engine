package currency

// Default formatting options
const (
	// Default precision for different currency types
	DefaultPrecision     = 2
	ZeroPrecision       = 0
	HighPrecision       = 4
	CryptoPrecision     = 8
	
	// Default separators
	DefaultThousandsSep = ","
	DefaultDecimalSep   = "."
	
	// Common formatting styles
	NegativeStyleMinus       = "minus"
	NegativeStyleParentheses = "parentheses"
	NegativeStyleMinusSymbol = "minus_symbol"
)

// Exchange rate constants
const (
	// Rate freshness thresholds (in minutes)
	RateFreshThreshold    = 60  // 1 hour
	RateStaleThreshold    = 1440 // 24 hours
	
	// Rate tolerance for comparisons
	RateToleranceDefault  = 0.0001
	RateToleranceStrict   = 0.00001
	RateToleranceLoose    = 0.001
	
	// Maximum allowed exchange rate values
	MaxExchangeRate = 1000000.0
	MinExchangeRate = 0.000001
)

// Amount limits and tolerances
const (
	// Maximum amount values
	MaxAmount = 1e15
	MinAmount = -1e15
	
	// Tolerance for floating point comparisons
	AmountTolerance = 0.001
	StrictTolerance = 0.0001
	LooseTolerance  = 0.01
	
	// Zero threshold
	ZeroThreshold = 1e-10
)

// Currency groups for easier management
var (
	// Major currencies (most traded)
	MajorCurrencies = []CurrencyCode{
		USD, EUR, GBP, JPY, CHF, CAD, AUD,
	}
	
	// Asian currencies
	AsianCurrencies = []CurrencyCode{
		IDR, SGD, MYR, THB, PHP, VND, KRW, INR, CNY,
	}
	
	// European currencies
	EuropeanCurrencies = []CurrencyCode{
		EUR, GBP, CHF, SEK, NOK, DKK,
	}
	
	// Americas currencies
	AmericasCurrencies = []CurrencyCode{
		USD, CAD, BRL, MXN,
	}
	
	// Middle East currencies
	MiddleEastCurrencies = []CurrencyCode{
		SAR, AED, TRY,
	}
	
	// Currencies with zero decimal places
	ZeroDecimalCurrencies = []CurrencyCode{
		JPY, KRW, VND, IDR,
	}
	
	// Currencies with high precision (more than 2 decimal places)
	HighPrecisionCurrencies = []CurrencyCode{
		// Add currencies that need more than 2 decimal places
	}
)

// Common currency symbols mapping
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

// Currency names mapping
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

// Default decimal places for currencies
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

// IsMajorCurrency checks if a currency is a major currency
func IsMajorCurrency(code CurrencyCode) bool {
	for _, major := range MajorCurrencies {
		if major == code {
			return true
		}
	}
	return false
}

// IsAsianCurrency checks if a currency is an Asian currency
func IsAsianCurrency(code CurrencyCode) bool {
	for _, asian := range AsianCurrencies {
		if asian == code {
			return true
		}
	}
	return false
}

// IsEuropeanCurrency checks if a currency is a European currency
func IsEuropeanCurrency(code CurrencyCode) bool {
	for _, european := range EuropeanCurrencies {
		if european == code {
			return true
		}
	}
	return false
}

// IsZeroDecimalCurrency checks if a currency uses zero decimal places
func IsZeroDecimalCurrency(code CurrencyCode) bool {
	for _, zero := range ZeroDecimalCurrencies {
		if zero == code {
			return true
		}
	}
	return false
}

// GetCurrencySymbol returns the symbol for a currency code
func GetCurrencySymbol(code CurrencyCode) string {
	if symbol, exists := CurrencySymbols[code]; exists {
		return symbol
	}
	return string(code) // Fallback to currency code
}

// GetCurrencyName returns the name for a currency code
func GetCurrencyName(code CurrencyCode) string {
	if name, exists := CurrencyNames[code]; exists {
		return name
	}
	return string(code) // Fallback to currency code
}

// GetCurrencyDecimalPlaces returns the default decimal places for a currency
func GetCurrencyDecimalPlaces(code CurrencyCode) int {
	if places, exists := CurrencyDecimalPlaces[code]; exists {
		return places
	}
	return DefaultPrecision // Fallback to default
}

// IsValidCurrencyCode checks if a currency code is valid (exists in our mappings)
func IsValidCurrencyCode(code CurrencyCode) bool {
	_, exists := CurrencyNames[code]
	return exists
}

// GetSupportedCurrencyCodes returns all supported currency codes
func GetSupportedCurrencyCodes() []CurrencyCode {
	codes := make([]CurrencyCode, 0, len(CurrencyNames))
	for code := range CurrencyNames {
		codes = append(codes, code)
	}
	return codes
}

// Common exchange rate pairs
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

// IsCommonPair checks if a currency pair is commonly traded
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
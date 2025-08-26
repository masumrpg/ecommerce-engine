package currency

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Validator provides currency validation functionality
type Validator struct {
	calculator *Calculator
}

// NewValidator creates a new currency validator
func NewValidator(calculator *Calculator) *Validator {
	return &Validator{
		calculator: calculator,
	}
}

// ValidateMoney validates a Money struct
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

// ValidateExchangeRate validates an exchange rate
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

// BatchConverter provides batch conversion functionality
type BatchConverter struct {
	calculator *Calculator
}

// NewBatchConverter creates a new batch converter
func NewBatchConverter(calculator *Calculator) *BatchConverter {
	return &BatchConverter{
		calculator: calculator,
	}
}

// ConvertBatch converts multiple amounts to a target currency
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

// SumInCurrency sums multiple amounts in different currencies to a target currency
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

// CurrencyDetector provides currency detection from text
type CurrencyDetector struct {
	calculator *Calculator
	patterns   map[CurrencyCode]*regexp.Regexp
}

// NewCurrencyDetector creates a new currency detector
func NewCurrencyDetector(calculator *Calculator) *CurrencyDetector {
	detector := &CurrencyDetector{
		calculator: calculator,
		patterns:   make(map[CurrencyCode]*regexp.Regexp),
	}
	
	detector.initializePatterns()
	return detector
}

// initializePatterns sets up regex patterns for currency detection
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

// DetectCurrency detects currency from text
func (cd *CurrencyDetector) DetectCurrency(text string) []CurrencyCode {
	var detected []CurrencyCode
	
	for code, pattern := range cd.patterns {
		if pattern.MatchString(text) {
			detected = append(detected, code)
		}
	}
	
	return detected
}

// ExtractMoney extracts money amounts from text
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

// CurrencyFormatter provides advanced formatting options
type CurrencyFormatter struct {
	calculator *Calculator
	locales    map[string]LocaleInfo
}

// NewCurrencyFormatter creates a new currency formatter
func NewCurrencyFormatter(calculator *Calculator) *CurrencyFormatter {
	formatter := &CurrencyFormatter{
		calculator: calculator,
		locales:    make(map[string]LocaleInfo),
	}
	
	formatter.initializeLocales()
	return formatter
}

// initializeLocales sets up locale-specific formatting
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

// FormatWithLocale formats money according to locale settings
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

// GetLocaleInfo returns locale information
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

// AddLocale adds a new locale
func (cf *CurrencyFormatter) AddLocale(locale string, info LocaleInfo) {
	cf.locales[locale] = info
}

// Helper functions

// NewMoney creates a new Money instance with validation
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

// NewMoneyFromString creates Money from string representation
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

// IsZero checks if money amount is zero
func IsZero(money Money) bool {
	return math.Abs(money.Amount) < 0.001
}

// IsPositive checks if money amount is positive
func IsPositive(money Money) bool {
	return money.Amount > 0.001
}

// IsNegative checks if money amount is negative
func IsNegative(money Money) bool {
	return money.Amount < -0.001
}

// Abs returns absolute value of money
func Abs(money Money) Money {
	return Money{
		Amount:   math.Abs(money.Amount),
		Currency: money.Currency,
	}
}

// Negate returns negated money amount
func Negate(money Money) Money {
	return Money{
		Amount:   -money.Amount,
		Currency: money.Currency,
	}
}

// Min returns the smaller of two money amounts (same currency)
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

// Max returns the larger of two money amounts (same currency)
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

// Sum calculates the sum of multiple money amounts (same currency)
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

// Average calculates the average of multiple money amounts (same currency)
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

// Percentage calculates percentage of money amount
func Percentage(money Money, percent float64) Money {
	return Money{
		Amount:   money.Amount * (percent / 100.0),
		Currency: money.Currency,
	}
}

// Split splits money amount into equal parts
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

// Allocate allocates money according to ratios
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
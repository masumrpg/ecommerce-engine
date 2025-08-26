package currency

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/masumrpg/ecommerce-engine/pkg/utils"
)

// Calculator provides currency calculation and formatting functionality
type Calculator struct {
	currencies   map[CurrencyCode]Currency
	exchangeRates map[string]ExchangeRate // key: "FROM/TO"
	defaultRounding RoundingMode
}

// NewCalculator creates a new currency calculator with default currencies
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

// initializeDefaultCurrencies sets up common currencies with their properties
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

// Format formats a money amount according to currency rules
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

// formatNumber formats a number with thousands and decimal separators
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

// addThousandsSeparators adds thousands separators to a number string
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

// roundAmount rounds an amount according to the specified rounding mode
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

// Convert converts money from one currency to another
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

// Add performs addition of two money amounts
func (c *Calculator) Add(amount1, amount2 Money) (*ArithmeticResult, error) {
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount1,
		Amount2:   amount2,
		Operation: OperationAdd,
		Rounding:  c.defaultRounding,
	})
}

// Subtract performs subtraction of two money amounts
func (c *Calculator) Subtract(amount1, amount2 Money) (*ArithmeticResult, error) {
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount1,
		Amount2:   amount2,
		Operation: OperationSubtract,
		Rounding:  c.defaultRounding,
	})
}

// Multiply multiplies a money amount by a factor
func (c *Calculator) Multiply(amount Money, factor float64) (*ArithmeticResult, error) {
	return c.performArithmetic(ArithmeticInput{
		Amount1:   amount,
		Amount2:   Money{Amount: factor, Currency: amount.Currency},
		Operation: OperationMultiply,
		Rounding:  c.defaultRounding,
	})
}

// Divide divides a money amount by a divisor
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

// performArithmetic performs arithmetic operations on money amounts
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

// Compare compares two money amounts
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

// SetExchangeRate sets an exchange rate between two currencies
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

// GetExchangeRate gets the exchange rate between two currencies
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

// AddCurrency adds a new currency to the calculator
func (c *Calculator) AddCurrency(currency Currency) {
	c.currencies[currency.Code] = currency
}

// GetCurrency gets currency information
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

// GetSupportedCurrencies returns a list of all supported currencies
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

// SetDefaultRounding sets the default rounding mode
func (c *Calculator) SetDefaultRounding(mode RoundingMode) {
	c.defaultRounding = mode
}

// Parse parses a formatted currency string back to Money
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
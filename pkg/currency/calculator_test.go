package currency

import (
	"testing"
)

func TestNewCalculator(t *testing.T) {
	calc := NewCalculator()
	
	if calc == nil {
		t.Fatal("NewCalculator should not return nil")
	}
	
	// Test that default currencies are loaded
	currencies := calc.GetSupportedCurrencies()
	if currencies.Total == 0 {
		t.Error("Calculator should have default currencies loaded")
	}
	
	// Test that IDR is included
	idr, err := calc.GetCurrency(IDR)
	if err != nil {
		t.Errorf("IDR should be supported by default: %v", err)
	}
	if idr.Symbol != "Rp" {
		t.Errorf("IDR symbol should be 'Rp', got '%s'", idr.Symbol)
	}
	if idr.DecimalPlaces != 0 {
		t.Errorf("IDR should have 0 decimal places, got %d", idr.DecimalPlaces)
	}
}

func TestFormat(t *testing.T) {
	calc := NewCalculator()
	
	tests := []struct {
		name     string
		money    Money
		options  *FormatOptions
		expected string
		wantErr  bool
	}{
		{
			name:     "USD with default options",
			money:    Money{Amount: 1234.56, Currency: USD},
			options:  &FormatOptions{ShowSymbol: true},
			expected: "$1,234.56",
			wantErr:  false,
		},
		{
			name:     "IDR with default options",
			money:    Money{Amount: 15000, Currency: IDR},
			options:  &FormatOptions{ShowSymbol: true},
			expected: "Rp 15.000",
			wantErr:  false,
		},
		{
			name:     "EUR with custom options",
			money:    Money{Amount: 999.99, Currency: EUR},
			options:  &FormatOptions{ShowSymbol: true, SymbolFirst: &[]bool{false}[0]},
			expected: "999,99 €",
			wantErr:  false,
		},
		{
			name:     "Negative amount with parentheses",
			money:    Money{Amount: -100.50, Currency: USD},
			options:  &FormatOptions{ShowSymbol: true, NegativeStyle: "parentheses"},
			expected: "$(100.50)",
			wantErr:  false,
		},
		{
			name:     "Show currency code instead of symbol",
			money:    Money{Amount: 500, Currency: USD},
			options:  &FormatOptions{ShowCode: true},
			expected: "USD500.00",
			wantErr:  false,
		},
		{
			name:    "Unsupported currency",
			money:   Money{Amount: 100, Currency: "XXX"},
			options: &FormatOptions{ShowSymbol: true},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Format(tt.money, tt.options)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestConvert(t *testing.T) {
	calc := NewCalculator()
	
	// Set up exchange rates
	calc.SetExchangeRate(USD, IDR, 15000, "test")
	calc.SetExchangeRate(USD, EUR, 0.85, "test")
	
	tests := []struct {
		name     string
		input    ConversionInput
		expected float64
		wantErr  bool
	}{
		{
			name:     "USD to IDR",
			input:    ConversionInput{Amount: 100, From: USD, To: IDR},
			expected: 1500000,
			wantErr:  false,
		},
		{
			name:     "USD to EUR",
			input:    ConversionInput{Amount: 100, From: USD, To: EUR},
			expected: 85,
			wantErr:  false,
		},
		{
			name:     "Same currency conversion",
			input:    ConversionInput{Amount: 100, From: USD, To: USD},
			expected: 100,
			wantErr:  false,
		},
		{
			name:    "No exchange rate available",
			input:   ConversionInput{Amount: 100, From: USD, To: JPY},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Convert(tt.input)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.ConvertedAmount.Amount != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result.ConvertedAmount.Amount)
			}
			
			if result.ConvertedAmount.Currency != tt.input.To {
				t.Errorf("Expected currency %s, got %s", tt.input.To, result.ConvertedAmount.Currency)
			}
		})
	}
}

func TestArithmeticOperations(t *testing.T) {
	calc := NewCalculator()
	
	tests := []struct {
		name      string
		amount1   Money
		amount2   Money
		operation string
		expected  float64
		wantErr   bool
	}{
		{
			name:      "Add USD amounts",
			amount1:   Money{Amount: 100.50, Currency: USD},
			amount2:   Money{Amount: 50.25, Currency: USD},
			operation: "add",
			expected:  150.75,
			wantErr:   false,
		},
		{
			name:      "Subtract USD amounts",
			amount1:   Money{Amount: 100.50, Currency: USD},
			amount2:   Money{Amount: 50.25, Currency: USD},
			operation: "subtract",
			expected:  50.25,
			wantErr:   false,
		},
		{
			name:      "Multiply USD amount",
			amount1:   Money{Amount: 100, Currency: USD},
			amount2:   Money{Amount: 2.5, Currency: USD},
			operation: "multiply",
			expected:  250,
			wantErr:   false,
		},
		{
			name:      "Divide USD amount",
			amount1:   Money{Amount: 100, Currency: USD},
			amount2:   Money{Amount: 4, Currency: USD},
			operation: "divide",
			expected:  25,
			wantErr:   false,
		},
		{
			name:      "Different currencies",
			amount1:   Money{Amount: 100, Currency: USD},
			amount2:   Money{Amount: 100, Currency: EUR},
			operation: "add",
			wantErr:   true,
		},
		{
			name:      "Divide by zero",
			amount1:   Money{Amount: 100, Currency: USD},
			amount2:   Money{Amount: 0, Currency: USD},
			operation: "divide",
			wantErr:   true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result *ArithmeticResult
			var err error
			
			switch tt.operation {
			case "add":
				result, err = calc.Add(tt.amount1, tt.amount2)
			case "subtract":
				result, err = calc.Subtract(tt.amount1, tt.amount2)
			case "multiply":
				result, err = calc.Multiply(tt.amount1, tt.amount2.Amount)
			case "divide":
				result, err = calc.Divide(tt.amount1, tt.amount2.Amount)
			}
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.Result.Amount != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result.Result.Amount)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	calc := NewCalculator()
	
	tests := []struct {
		name      string
		amount1   Money
		amount2   Money
		isEqual   bool
		isGreater bool
		isLess    bool
		wantErr   bool
	}{
		{
			name:      "Equal amounts",
			amount1:   Money{Amount: 100, Currency: USD},
			amount2:   Money{Amount: 100, Currency: USD},
			isEqual:   true,
			isGreater: false,
			isLess:    false,
			wantErr:   false,
		},
		{
			name:      "First amount greater",
			amount1:   Money{Amount: 150, Currency: USD},
			amount2:   Money{Amount: 100, Currency: USD},
			isEqual:   false,
			isGreater: true,
			isLess:    false,
			wantErr:   false,
		},
		{
			name:      "First amount less",
			amount1:   Money{Amount: 75, Currency: USD},
			amount2:   Money{Amount: 100, Currency: USD},
			isEqual:   false,
			isGreater: false,
			isLess:    true,
			wantErr:   false,
		},
		{
			name:    "Different currencies",
			amount1: Money{Amount: 100, Currency: USD},
			amount2: Money{Amount: 100, Currency: EUR},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Compare(tt.amount1, tt.amount2)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.IsEqual != tt.isEqual {
				t.Errorf("Expected IsEqual %v, got %v", tt.isEqual, result.IsEqual)
			}
			
			if result.IsGreater != tt.isGreater {
				t.Errorf("Expected IsGreater %v, got %v", tt.isGreater, result.IsGreater)
			}
			
			if result.IsLess != tt.isLess {
				t.Errorf("Expected IsLess %v, got %v", tt.isLess, result.IsLess)
			}
		})
	}
}

func TestExchangeRateManagement(t *testing.T) {
	calc := NewCalculator()
	
	// Test setting exchange rate
	calc.SetExchangeRate(USD, IDR, 15000, "test")
	
	// Test getting exchange rate
	rate, err := calc.GetExchangeRate(USD, IDR)
	if err != nil {
		t.Errorf("Unexpected error getting exchange rate: %v", err)
	}
	
	if rate.Rate != 15000 {
		t.Errorf("Expected rate 15000, got %f", rate.Rate)
	}
	
	if rate.From != USD || rate.To != IDR {
		t.Errorf("Expected USD to IDR, got %s to %s", rate.From, rate.To)
	}
	
	// Test inverse rate is automatically created
	inverseRate, err := calc.GetExchangeRate(IDR, USD)
	if err != nil {
		t.Errorf("Unexpected error getting inverse exchange rate: %v", err)
	}
	
	expectedInverse := 1.0 / 15000
	if inverseRate.Rate != expectedInverse {
		t.Errorf("Expected inverse rate %f, got %f", expectedInverse, inverseRate.Rate)
	}
	
	// Test getting non-existent rate
	_, err = calc.GetExchangeRate(USD, JPY)
	if err == nil {
		t.Error("Expected error for non-existent exchange rate")
	}
}

func TestParse(t *testing.T) {
	calc := NewCalculator()
	
	tests := []struct {
		name     string
		input    string
		currency CurrencyCode
		expected float64
		wantErr  bool
	}{
		{
			name:     "USD with symbol",
			input:    "$1,234.56",
			currency: USD,
			expected: 1234.56,
			wantErr:  false,
		},
		{
			name:     "IDR with symbol",
			input:    "Rp 15.000",
			currency: IDR,
			expected: 15000,
			wantErr:  false,
		},
		{
			name:     "EUR with symbol",
			input:    "999,99 €",
			currency: EUR,
			expected: 999.99,
			wantErr:  false,
		},
		{
			name:     "Negative with parentheses",
			input:    "$(100.50)",
			currency: USD,
			expected: -100.50,
			wantErr:  false,
		},
		{
			name:     "Plain number",
			input:    "500.25",
			currency: USD,
			expected: 500.25,
			wantErr:  false,
		},
		{
			name:    "Invalid input",
			input:   "not a number",
			currency: USD,
			wantErr: true,
		},
		{
			name:    "Unsupported currency",
			input:   "100",
			currency: "XXX",
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calc.Parse(tt.input, tt.currency)
			
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			if result.Amount != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result.Amount)
			}
			
			if result.Currency != tt.currency {
				t.Errorf("Expected currency %s, got %s", tt.currency, result.Currency)
			}
		})
	}
}

func TestRoundingModes(t *testing.T) {
	calc := NewCalculator()
	
	tests := []struct {
		name      string
		amount    float64
		precision int
		mode      RoundingMode
		expected  float64
	}{
		{
			name:      "Half up rounding",
			amount:    1.235,
			precision: 2,
			mode:      RoundingModeHalfUp,
			expected:  1.24,
		},
		{
			name:      "Half down rounding",
			amount:    1.235,
			precision: 2,
			mode:      RoundingModeHalfDown,
			expected:  1.23,
		},
		{
			name:      "Up rounding",
			amount:    1.231,
			precision: 2,
			mode:      RoundingModeUp,
			expected:  1.24,
		},
		{
			name:      "Down rounding",
			amount:    1.239,
			precision: 2,
			mode:      RoundingModeDown,
			expected:  1.23,
		},
		{
			name:      "Truncate rounding",
			amount:    1.999,
			precision: 2,
			mode:      RoundingModeTruncate,
			expected:  1.99,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.roundAmount(tt.amount, tt.precision, tt.mode)
			
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestCurrencyManagement(t *testing.T) {
	calc := NewCalculator()
	
	// Test adding custom currency
	customCurrency := Currency{
		Code:          "TEST",
		Name:          "Test Currency",
		Symbol:        "T$",
		DecimalPlaces: 3,
		ThousandsSep:  ",",
		DecimalSep:    ".",
		SymbolFirst:   true,
		SpaceBetween:  false,
	}
	
	calc.AddCurrency(customCurrency)
	
	// Test getting the custom currency
	retrieved, err := calc.GetCurrency("TEST")
	if err != nil {
		t.Errorf("Unexpected error getting custom currency: %v", err)
	}
	
	if retrieved.Symbol != "T$" {
		t.Errorf("Expected symbol 'T$', got '%s'", retrieved.Symbol)
	}
	
	if retrieved.DecimalPlaces != 3 {
		t.Errorf("Expected 3 decimal places, got %d", retrieved.DecimalPlaces)
	}
	
	// Test formatting with custom currency
	money := Money{Amount: 1234.5678, Currency: "TEST"}
	formatted, err := calc.Format(money, &FormatOptions{ShowSymbol: true})
	if err != nil {
		t.Errorf("Unexpected error formatting custom currency: %v", err)
	}
	
	expected := "T$1,234.568" // Should round to 3 decimal places
	if formatted != expected {
		t.Errorf("Expected '%s', got '%s'", expected, formatted)
	}
}

func TestDefaultRounding(t *testing.T) {
	calc := NewCalculator()
	
	// Test default rounding mode
	originalMode := calc.defaultRounding
	if originalMode != RoundingModeHalfUp {
		t.Errorf("Expected default rounding mode to be HalfUp, got %s", originalMode)
	}
	
	// Test setting new default rounding mode
	calc.SetDefaultRounding(RoundingModeDown)
	if calc.defaultRounding != RoundingModeDown {
		t.Errorf("Expected rounding mode to be Down, got %s", calc.defaultRounding)
	}
	
	// Test that new rounding mode is used
	money := Money{Amount: 1.999, Currency: USD}
	result, err := calc.Multiply(money, 1.0) // This should trigger rounding
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	
	// With RoundingModeDown, 1.999 should round down to 1.99
	if result.Result.Amount != 1.99 {
		t.Errorf("Expected 1.99, got %f", result.Result.Amount)
	}
}

func BenchmarkFormat(b *testing.B) {
	calc := NewCalculator()
	money := Money{Amount: 1234.56, Currency: USD}
	options := &FormatOptions{ShowSymbol: true}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.Format(money, options)
	}
}

func BenchmarkConvert(b *testing.B) {
	calc := NewCalculator()
	calc.SetExchangeRate(USD, IDR, 15000, "test")
	input := ConversionInput{Amount: 100, From: USD, To: IDR}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.Convert(input)
	}
}

func BenchmarkArithmetic(b *testing.B) {
	calc := NewCalculator()
	amount1 := Money{Amount: 100.50, Currency: USD}
	amount2 := Money{Amount: 50.25, Currency: USD}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.Add(amount1, amount2)
	}
}
package currency

import (
	"math"
	"testing"
	"time"
)

func TestValidator(t *testing.T) {
	calc := NewCalculator()
	validator := NewValidator(calc)
	
	t.Run("ValidateMoney", func(t *testing.T) {
		tests := []struct {
			name    string
			money   Money
			wantErr bool
		}{
			{
				name:    "Valid USD money",
				money:   Money{Amount: 100.50, Currency: USD},
				wantErr: false,
			},
			{
				name:    "Valid IDR money",
				money:   Money{Amount: 15000, Currency: IDR},
				wantErr: false,
			},
			{
				name:    "Unsupported currency",
				money:   Money{Amount: 100, Currency: "XXX"},
				wantErr: true,
			},
			{
				name:    "Zero amount",
				money:   Money{Amount: 0, Currency: USD},
				wantErr: false,
			},
			{
				name:    "Negative amount",
				money:   Money{Amount: -100, Currency: USD},
				wantErr: false,
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateMoney(tt.money)
				
				if tt.wantErr {
					if err == nil {
						t.Error("Expected validation error but got none")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected validation error: %v", err)
					}
				}
			})
		}
	})
	
	t.Run("ValidateExchangeRate", func(t *testing.T) {
		tests := []struct {
			name    string
			rate    ExchangeRate
			wantErr bool
		}{
			{
				name: "Valid exchange rate",
				rate: ExchangeRate{
					From: USD,
					To:   IDR,
					Rate: 15000,
					Timestamp: time.Now(),
					Source: "test",
				},
				wantErr: false,
			},
			{
				name: "Same currency",
				rate: ExchangeRate{
					From: USD,
					To:   USD,
					Rate: 1.0,
				},
				wantErr: true,
			},
			{
				name: "Zero rate",
				rate: ExchangeRate{
					From: USD,
					To:   EUR,
					Rate: 0,
				},
				wantErr: true,
			},
			{
				name: "Negative rate",
				rate: ExchangeRate{
					From: USD,
					To:   EUR,
					Rate: -0.85,
				},
				wantErr: true,
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validator.ValidateExchangeRate(tt.rate)
				
				if tt.wantErr {
					if err == nil {
						t.Error("Expected validation error but got none")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected validation error: %v", err)
					}
				}
			})
		}
	})
}

func TestBatchConverter(t *testing.T) {
	calc := NewCalculator()
	calc.SetExchangeRate(USD, IDR, 15000, "test")
	calc.SetExchangeRate(USD, EUR, 0.85, "test")
	
	batchConverter := NewBatchConverter(calc)
	
	t.Run("ConvertBatch", func(t *testing.T) {
		amounts := []Money{
			{Amount: 100, Currency: USD},
			{Amount: 200, Currency: USD},
			{Amount: 50, Currency: USD},
		}
		
		results, errors := batchConverter.ConvertBatch(amounts, IDR)
		
		if len(errors) > 0 {
			t.Errorf("Unexpected errors: %v", errors)
		}
		
		if len(results) != 3 {
			t.Errorf("Expected 3 results, got %d", len(results))
		}
		
		// Check first conversion
		if results[0].ConvertedAmount.Amount != 1500000 {
			t.Errorf("Expected 1500000, got %f", results[0].ConvertedAmount.Amount)
		}
		
		if results[0].ConvertedAmount.Currency != IDR {
			t.Errorf("Expected IDR, got %s", results[0].ConvertedAmount.Currency)
		}
	})
	
	t.Run("SumInCurrency", func(t *testing.T) {
		amounts := []Money{
			{Amount: 100, Currency: USD},
			{Amount: 200, Currency: USD},
			{Amount: 50, Currency: USD},
		}
		
		result, err := batchConverter.SumInCurrency(amounts, IDR)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		// 100 + 200 + 50 = 350 USD = 350 * 15000 = 5,250,000 IDR
		expected := 350.0 * 15000
		if result.Amount != expected {
			t.Errorf("Expected %f, got %f", expected, result.Amount)
		}
		
		if result.Currency != IDR {
			t.Errorf("Expected IDR, got %s", result.Currency)
		}
	})
	
	t.Run("ConvertBatch with errors", func(t *testing.T) {
		amounts := []Money{
			{Amount: 100, Currency: USD},
			{Amount: 200, Currency: JPY}, // No exchange rate for JPY
		}
		
		results, errors := batchConverter.ConvertBatch(amounts, IDR)
		
		if len(errors) == 0 {
			t.Error("Expected conversion errors but got none")
		}
		
		// Should have one successful conversion
		if len(results) != 1 {
			t.Errorf("Expected 1 successful result, got %d", len(results))
		}
	})
}

func TestCurrencyDetector(t *testing.T) {
	calc := NewCalculator()
	detector := NewCurrencyDetector(calc)
	
	t.Run("DetectCurrency", func(t *testing.T) {
		tests := []struct {
			name     string
			text     string
			expected []CurrencyCode
		}{
			{
				name:     "USD symbol",
				text:     "The price is $100",
				expected: []CurrencyCode{USD},
			},
			{
				name:     "IDR symbol",
				text:     "Harga Rp 15000",
				expected: []CurrencyCode{IDR},
			},
			{
				name:     "EUR symbol",
				text:     "Cost is €85.50",
				expected: []CurrencyCode{EUR},
			},
			{
				name:     "Currency code",
				text:     "Amount: USD 100",
				expected: []CurrencyCode{USD},
			},
			{
				name:     "No currency",
				text:     "Just a number 100",
				expected: []CurrencyCode{},
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := detector.DetectCurrency(tt.text)
				
				if len(result) != len(tt.expected) {
					t.Errorf("Expected %d currencies, got %d", len(tt.expected), len(result))
					return
				}
				
				for _, expected := range tt.expected {
				found := false
				for _, detected := range result {
					if detected == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected currency %s not found in result %v", expected, result)
				}
			}
			})
		}
	})
	
	t.Run("ExtractMoney", func(t *testing.T) {
		tests := []struct {
			name     string
			text     string
			expected []Money
		}{
			{
				name: "USD amount",
				text: "The price is $100.50",
				expected: []Money{
					{Amount: 100.50, Currency: USD},
				},
			},
			{
				name: "IDR amount",
				text: "Harga Rp 15,000",
				expected: []Money{
					{Amount: 15000, Currency: IDR},
				},
			},
			{
				name: "Multiple amounts",
				text: "USD 100 and EUR 85",
				expected: []Money{
					{Amount: 100, Currency: USD},
					{Amount: 85, Currency: EUR},
				},
			},
			{
				name:     "No amounts",
				text:     "Just text without money",
				expected: []Money{},
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := detector.ExtractMoney(tt.text)
				
				if len(result) != len(tt.expected) {
					t.Errorf("Expected %d money amounts, got %d", len(tt.expected), len(result))
					return
				}
				
				for _, expected := range tt.expected {
				found := false
				for _, extracted := range result {
					if extracted.Amount == expected.Amount && extracted.Currency == expected.Currency {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected money %+v not found in result %+v", expected, result)
				}
			}
			})
		}
	})
}

func TestCurrencyFormatter(t *testing.T) {
	calc := NewCalculator()
	formatter := NewCurrencyFormatter(calc)
	
	t.Run("FormatWithLocale", func(t *testing.T) {
		tests := []struct {
			name     string
			money    Money
			locale   string
			wantErr  bool
		}{
			{
				name:    "Indonesian locale with IDR",
				money:   Money{Amount: 15000, Currency: IDR},
				locale:  "id-ID",
				wantErr: false,
			},
			{
				name:    "US locale with USD",
				money:   Money{Amount: 100.50, Currency: USD},
				locale:  "en-US",
				wantErr: false,
			},
			{
				name:    "German locale with EUR",
				money:   Money{Amount: 85.50, Currency: EUR},
				locale:  "de-DE",
				wantErr: false,
			},
			{
				name:    "Unsupported locale",
				money:   Money{Amount: 100, Currency: USD},
				locale:  "xx-XX",
				wantErr: true,
			},
		}
		
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := formatter.FormatWithLocale(tt.money, tt.locale)
				
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
				
				if result == "" {
					t.Error("Expected non-empty result")
				}
			})
		}
	})
	
	t.Run("GetLocaleInfo", func(t *testing.T) {
		localeInfo, err := formatter.GetLocaleInfo("id-ID")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if localeInfo.Language != "Indonesian" {
			t.Errorf("Expected Indonesian, got %s", localeInfo.Language)
		}
		
		if localeInfo.CurrencyCode != IDR {
			t.Errorf("Expected IDR, got %s", localeInfo.CurrencyCode)
		}
		
		// Test unsupported locale
		_, err = formatter.GetLocaleInfo("xx-XX")
		if err == nil {
			t.Error("Expected error for unsupported locale")
		}
	})
	
	t.Run("AddLocale", func(t *testing.T) {
		customLocale := LocaleInfo{
			Locale:       "test-TEST",
			Language:     "Test Language",
			Country:      "Test Country",
			CurrencyName: "Test Currency",
			CurrencyCode: USD,
		}
		
		formatter.AddLocale("test-TEST", customLocale)
		
		// Test that the locale was added
		retrieved, err := formatter.GetLocaleInfo("test-TEST")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if retrieved.Language != "Test Language" {
			t.Errorf("Expected 'Test Language', got %s", retrieved.Language)
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("NewMoney", func(t *testing.T) {
		money, err := NewMoney(100.50, USD)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if money.Amount != 100.50 {
			t.Errorf("Expected 100.50, got %f", money.Amount)
		}
		
		if money.Currency != USD {
			t.Errorf("Expected USD, got %s", money.Currency)
		}
		
		// Test invalid amount (NaN)
		_, err = NewMoney(math.NaN(), USD)
		if err == nil {
			t.Error("Expected error for NaN amount")
		}
	})
	
	t.Run("NewMoneyFromString", func(t *testing.T) {
		money, err := NewMoneyFromString("100.50", USD)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if money.Amount != 100.50 {
			t.Errorf("Expected 100.50, got %f", money.Amount)
		}
		
		// Test invalid string
		_, err = NewMoneyFromString("not a number", USD)
		if err == nil {
			t.Error("Expected error for invalid string")
		}
	})
	
	t.Run("IsZero", func(t *testing.T) {
		if !IsZero(Money{Amount: 0, Currency: USD}) {
			t.Error("Expected true for zero amount")
		}
		
		if IsZero(Money{Amount: 0.1, Currency: USD}) {
			t.Error("Expected false for non-zero amount")
		}
		
		// Test very small amount (within tolerance)
		if !IsZero(Money{Amount: 0.0001, Currency: USD}) {
			t.Error("Expected true for very small amount")
		}
	})
	
	t.Run("IsPositive", func(t *testing.T) {
		if !IsPositive(Money{Amount: 100, Currency: USD}) {
			t.Error("Expected true for positive amount")
		}
		
		if IsPositive(Money{Amount: -100, Currency: USD}) {
			t.Error("Expected false for negative amount")
		}
		
		if IsPositive(Money{Amount: 0, Currency: USD}) {
			t.Error("Expected false for zero amount")
		}
	})
	
	t.Run("IsNegative", func(t *testing.T) {
		if !IsNegative(Money{Amount: -100, Currency: USD}) {
			t.Error("Expected true for negative amount")
		}
		
		if IsNegative(Money{Amount: 100, Currency: USD}) {
			t.Error("Expected false for positive amount")
		}
		
		if IsNegative(Money{Amount: 0, Currency: USD}) {
			t.Error("Expected false for zero amount")
		}
	})
	
	t.Run("Abs", func(t *testing.T) {
		result := Abs(Money{Amount: -100.50, Currency: USD})
		if result.Amount != 100.50 {
			t.Errorf("Expected 100.50, got %f", result.Amount)
		}
		
		result = Abs(Money{Amount: 100.50, Currency: USD})
		if result.Amount != 100.50 {
			t.Errorf("Expected 100.50, got %f", result.Amount)
		}
	})
	
	t.Run("Negate", func(t *testing.T) {
		result := Negate(Money{Amount: 100.50, Currency: USD})
		if result.Amount != -100.50 {
			t.Errorf("Expected -100.50, got %f", result.Amount)
		}
		
		result = Negate(Money{Amount: -100.50, Currency: USD})
		if result.Amount != 100.50 {
			t.Errorf("Expected 100.50, got %f", result.Amount)
		}
	})
	
	t.Run("Min", func(t *testing.T) {
		a := Money{Amount: 100, Currency: USD}
		b := Money{Amount: 150, Currency: USD}
		
		result, err := Min(a, b)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if result.Amount != 100 {
			t.Errorf("Expected 100, got %f", result.Amount)
		}
		
		// Test different currencies
		c := Money{Amount: 100, Currency: EUR}
		_, err = Min(a, c)
		if err == nil {
			t.Error("Expected error for different currencies")
		}
	})
	
	t.Run("Max", func(t *testing.T) {
		a := Money{Amount: 100, Currency: USD}
		b := Money{Amount: 150, Currency: USD}
		
		result, err := Max(a, b)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if result.Amount != 150 {
			t.Errorf("Expected 150, got %f", result.Amount)
		}
	})
	
	t.Run("Sum", func(t *testing.T) {
		amounts := []Money{
			{Amount: 100, Currency: USD},
			{Amount: 200, Currency: USD},
			{Amount: 50, Currency: USD},
		}
		
		result, err := Sum(amounts)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if result.Amount != 350 {
			t.Errorf("Expected 350, got %f", result.Amount)
		}
		
		if result.Currency != USD {
			t.Errorf("Expected USD, got %s", result.Currency)
		}
		
		// Test empty array
		_, err = Sum([]Money{})
		if err == nil {
			t.Error("Expected error for empty array")
		}
		
		// Test mixed currencies
		mixedAmounts := []Money{
			{Amount: 100, Currency: USD},
			{Amount: 200, Currency: EUR},
		}
		_, err = Sum(mixedAmounts)
		if err == nil {
			t.Error("Expected error for mixed currencies")
		}
	})
	
	t.Run("Average", func(t *testing.T) {
		amounts := []Money{
			{Amount: 100, Currency: USD},
			{Amount: 200, Currency: USD},
			{Amount: 300, Currency: USD},
		}
		
		result, err := Average(amounts)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if result.Amount != 200 {
			t.Errorf("Expected 200, got %f", result.Amount)
		}
	})
	
	t.Run("Percentage", func(t *testing.T) {
		money := Money{Amount: 1000, Currency: USD}
		result := Percentage(money, 15) // 15%
		
		if result.Amount != 150 {
			t.Errorf("Expected 150, got %f", result.Amount)
		}
		
		if result.Currency != USD {
			t.Errorf("Expected USD, got %s", result.Currency)
		}
	})
	
	t.Run("Split", func(t *testing.T) {
		money := Money{Amount: 100, Currency: USD}
		parts, remainder := Split(money, 3)
		
		if len(parts) != 3 {
			t.Errorf("Expected 3 parts, got %d", len(parts))
		}
		
		// Each part should be approximately 33.33
		expectedPart := 100.0 / 3.0
		for i, part := range parts {
			if part.Amount != expectedPart {
				t.Errorf("Part %d: expected %f, got %f", i, expectedPart, part.Amount)
			}
		}
		
		// Check remainder
		if remainder.Currency != USD {
			t.Errorf("Expected USD remainder, got %s", remainder.Currency)
		}
	})
	
	t.Run("Allocate", func(t *testing.T) {
		money := Money{Amount: 1000, Currency: USD}
		ratios := []float64{0.5, 0.3, 0.2} // 50%, 30%, 20%
		
		result, err := Allocate(money, ratios)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		
		if len(result) != 3 {
			t.Errorf("Expected 3 allocations, got %d", len(result))
		}
		
		// Check allocations
		expected := []float64{500, 300, 200}
		for i, allocation := range result {
			if allocation.Amount != expected[i] {
				t.Errorf("Allocation %d: expected %f, got %f", i, expected[i], allocation.Amount)
			}
		}
		
		// Test empty ratios
		_, err = Allocate(money, []float64{})
		if err == nil {
			t.Error("Expected error for empty ratios")
		}
		
		// Test negative ratio
		_, err = Allocate(money, []float64{0.5, -0.3})
		if err == nil {
			t.Error("Expected error for negative ratio")
		}
	})
}

func BenchmarkValidateMoney(b *testing.B) {
	calc := NewCalculator()
	validator := NewValidator(calc)
	money := Money{Amount: 100.50, Currency: USD}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.ValidateMoney(money)
	}
}

func BenchmarkBatchConvert(b *testing.B) {
	calc := NewCalculator()
	calc.SetExchangeRate(USD, IDR, 15000, "test")
	batchConverter := NewBatchConverter(calc)
	
	amounts := []Money{
		{Amount: 100, Currency: USD},
		{Amount: 200, Currency: USD},
		{Amount: 50, Currency: USD},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = batchConverter.ConvertBatch(amounts, IDR)
	}
}

func BenchmarkSum(b *testing.B) {
	amounts := []Money{
		{Amount: 100, Currency: USD},
		{Amount: 200, Currency: USD},
		{Amount: 50, Currency: USD},
		{Amount: 75, Currency: USD},
		{Amount: 125, Currency: USD},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Sum(amounts)
	}
}
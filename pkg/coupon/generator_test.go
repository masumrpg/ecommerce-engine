package coupon

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateCode(t *testing.T) {
	t.Run("BasicGeneration", func(t *testing.T) {
		config := GeneratorConfig{
			Length: 8,
		}
		
		code, err := GenerateCode(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if len(code) != 8 {
			t.Errorf("Expected code length 8, got %d", len(code))
		}
	})
	
	t.Run("WithPrefix", func(t *testing.T) {
		config := GeneratorConfig{
			Length: 6,
			Prefix: "SAVE",
		}
		
		code, err := GenerateCode(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.HasPrefix(code, "SAVE") {
			t.Errorf("Expected code to start with 'SAVE', got %s", code)
		}
	})
	
	t.Run("WithSuffix", func(t *testing.T) {
		config := GeneratorConfig{
			Length: 6,
			Suffix: "OFF",
		}
		
		code, err := GenerateCode(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.HasSuffix(code, "OFF") {
			t.Errorf("Expected code to end with 'OFF', got %s", code)
		}
	})
	
	t.Run("WithPrefixPattern", func(t *testing.T) {
		config := GeneratorConfig{
			Pattern: "PREFIX-XXXXXX",
			Prefix:  "SALE",
			Length:  6,
		}
		
		code, err := GenerateCode(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.HasPrefix(code, "SALE-") {
			t.Errorf("Expected code to start with 'SALE-', got %s", code)
		}
	})
	
	t.Run("WithExcludeChars", func(t *testing.T) {
		config := GeneratorConfig{
			Length:       20,
			ExcludeChars: "0O1I",
		}
		
		code, err := GenerateCode(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// Check that excluded characters are not present
		excludedChars := []string{"0", "O", "1", "I"}
		for _, char := range excludedChars {
			if strings.Contains(code, char) {
				t.Errorf("Code contains excluded character '%s': %s", char, code)
			}
		}
	})
	
	t.Run("ZeroLength", func(t *testing.T) {
		config := GeneratorConfig{
			Length: 0,
		}
		
		code, err := GenerateCode(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		// With zero length, should use default length of 8
		if len(code) != 8 {
			t.Errorf("Expected default length 8, got %d", len(code))
		}
	})
}

func TestGenerateCodes(t *testing.T) {
	t.Run("MultipleCodesGeneration", func(t *testing.T) {
		config := GeneratorConfig{
			Length: 6,
			Count:  5,
		}
		
		codes, err := GenerateCodes(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if len(codes) != 5 {
			t.Errorf("Expected 5 codes, got %d", len(codes))
		}
		
		// Check all codes are unique
		uniqueMap := make(map[string]bool)
		for _, code := range codes {
			if uniqueMap[code] {
				t.Errorf("Duplicate code found: %s", code)
			}
			uniqueMap[code] = true
			
			if len(code) != 6 {
				t.Errorf("Expected code length 6, got %d for code %s", len(code), code)
			}
		}
	})
	
	t.Run("ZeroCount", func(t *testing.T) {
		config := GeneratorConfig{
			Length: 8,
			Count:  0,
		}
		
		codes, err := GenerateCodes(config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if len(codes) != 0 {
			t.Errorf("Expected 0 codes, got %d", len(codes))
		}
	})
}

func TestGenerateSeasonalCode(t *testing.T) {
	t.Run("ChristmasCode", func(t *testing.T) {
		config := GeneratorConfig{
			ExcludeChars: "0O1I",
		}
		
		code, err := GenerateSeasonalCode("christmas", 2024, config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.HasPrefix(code, "XMAS2024") {
			t.Errorf("Expected code to start with 'XMAS2024', got %s", code)
		}
	})
	
	t.Run("UnknownSeason", func(t *testing.T) {
		config := GeneratorConfig{}
		
		code, err := GenerateSeasonalCode("unknown", 2024, config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.HasPrefix(code, "SPECIAL2024") {
			t.Errorf("Expected code to start with 'SPECIAL2024', got %s", code)
		}
	})
}

func TestGenerateFlashSaleCode(t *testing.T) {
	t.Run("FlashSale50", func(t *testing.T) {
		config := GeneratorConfig{}
		
		code, err := GenerateFlashSaleCode(50, config)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if !strings.Contains(code, "50") {
			t.Errorf("Expected code to contain '50', got %s", code)
		}
		
		// Should start with one of the flash prefixes
		flashPrefixes := []string{"FLASH", "QUICK", "RUSH", "SPEED", "FAST"}
		hasValidPrefix := false
		for _, prefix := range flashPrefixes {
			if strings.HasPrefix(code, prefix) {
				hasValidPrefix = true
				break
			}
		}
		
		if !hasValidPrefix {
			t.Errorf("Expected code to start with flash prefix, got %s", code)
		}
	})
}

func TestValidateCodeFormat(t *testing.T) {
	t.Run("ValidCode", func(t *testing.T) {
		config := GeneratorConfig{
			Prefix: "SAVE",
			Suffix: "OFF",
		}
		
		valid := ValidateCodeFormat("SAVE123OFF", config)
		if !valid {
			t.Error("Expected code to be valid")
		}
	})
	
	t.Run("EmptyCode", func(t *testing.T) {
		config := GeneratorConfig{}
		
		valid := ValidateCodeFormat("", config)
		if valid {
			t.Error("Expected empty code to be invalid")
		}
	})
	
	t.Run("WithExcludedChars", func(t *testing.T) {
		config := GeneratorConfig{
			ExcludeChars: "0O1I",
		}
		
		valid := ValidateCodeFormat("SAVE0", config)
		if valid {
			t.Error("Expected code with excluded chars to be invalid")
		}
	})
	
	t.Run("WrongPrefix", func(t *testing.T) {
		config := GeneratorConfig{
			Prefix: "SAVE",
		}
		
		valid := ValidateCodeFormat("DEAL123", config)
		if valid {
			t.Error("Expected code with wrong prefix to be invalid")
		}
	})
}

func TestGenerateExpiryDate(t *testing.T) {
	t.Run("OneWeekExpiry", func(t *testing.T) {
		duration := 7 * 24 * time.Hour
		expiryDate := GenerateExpiryDate(duration)
		
		now := time.Now()
		expected := now.Add(duration)
		
		// Allow for small time differences (1 second)
		if expiryDate.Sub(expected).Abs() > time.Second {
			t.Errorf("Expected expiry date around %v, got %v", expected, expiryDate)
		}
	})
}

func TestGenerateBulkCodes(t *testing.T) {
	t.Run("MultiplePatternsGeneration", func(t *testing.T) {
		configs := []GeneratorConfig{
			{
				Length: 6,
				Count:  2,
				Pattern: "XXXXXXXX",
			},
			{
				Length: 4,
				Count:  3,
				Prefix: "SAVE",
			},
		}
		
		results, err := GenerateBulkCodes(configs)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		
		if len(results) != 2 {
			t.Errorf("Expected 2 pattern results, got %d", len(results))
		}
		
		// Check first pattern
		if codes, exists := results["XXXXXXXX"]; exists {
			if len(codes) != 2 {
				t.Errorf("Expected 2 codes for first pattern, got %d", len(codes))
			}
		} else {
			t.Error("Expected XXXXXXXX pattern in results")
		}
		
		// Check second pattern
		if codes, exists := results["pattern_1"]; exists {
			if len(codes) != 3 {
				t.Errorf("Expected 3 codes for second pattern, got %d", len(codes))
			}
			for _, code := range codes {
				if !strings.HasPrefix(code, "SAVE") {
					t.Errorf("Expected code to start with 'SAVE', got %s", code)
				}
			}
		} else {
			t.Error("Expected pattern_1 in results")
		}
	})
}

func BenchmarkGenerateCode(b *testing.B) {
	config := GeneratorConfig{
		Length: 8,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateCode(config)
	}
}

func BenchmarkGenerateCodes(b *testing.B) {
	config := GeneratorConfig{
		Length: 8,
		Count:  10,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = GenerateCodes(config)
	}
}

func BenchmarkValidateCodeFormat(b *testing.B) {
	config := GeneratorConfig{
		Prefix: "SAVE",
	}
	code := "SAVE123"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCodeFormat(code, config)
	}
}
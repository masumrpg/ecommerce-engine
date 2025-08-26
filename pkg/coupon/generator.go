package coupon

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// GenerateCode generates a single coupon code based on the configuration
func GenerateCode(config GeneratorConfig) (string, error) {
	if config.Length <= 0 {
		config.Length = 8 // Default length
	}

	if config.ExcludeChars == "" {
		config.ExcludeChars = "0O1I" // Default excluded characters
	}

	switch config.Pattern {
	case "PREFIX-XXXXXX":
		return generatePrefixPattern(config)
	case "XXXXXXXX":
		return generateRandomPattern(config)
	case "WORD-NUMBER":
		return generateWordNumberPattern(config)
	default:
		return generateRandomPattern(config)
	}
}

// GenerateCodes generates multiple unique coupon codes
func GenerateCodes(config GeneratorConfig) ([]string, error) {
	if config.Count == 0 {
		return []string{}, nil
	}
	
	if config.Count < 0 {
		config.Count = 1
	}

	codes := make([]string, 0, config.Count)
	uniqueCheck := make(map[string]bool)

	maxAttempts := config.Count * 10 // Prevent infinite loop
	attempts := 0

	for len(codes) < config.Count && attempts < maxAttempts {
		code, err := GenerateCode(config)
		if err != nil {
			return nil, err
		}

		if !uniqueCheck[code] {
			codes = append(codes, code)
			uniqueCheck[code] = true
		}

		attempts++
	}

	if len(codes) < config.Count {
		return codes, fmt.Errorf("could only generate %d unique codes out of %d requested", len(codes), config.Count)
	}

	return codes, nil
}

// generatePrefixPattern generates code with prefix pattern (e.g., SAVE-ABC123)
func generatePrefixPattern(config GeneratorConfig) (string, error) {
	prefix := config.Prefix
	if prefix == "" {
		prefix = "COUPON"
	}

	randomPart, err := generateRandomString(config.Length, config.ExcludeChars)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%s-%s", prefix, randomPart)

	if config.Suffix != "" {
		code += "-" + config.Suffix
	}

	return strings.ToUpper(code), nil
}

// generateRandomPattern generates purely random code
func generateRandomPattern(config GeneratorConfig) (string, error) {
	randomPart, err := generateRandomString(config.Length, config.ExcludeChars)
	if err != nil {
		return "", err
	}

	code := config.Prefix + randomPart + config.Suffix
	return strings.ToUpper(code), nil
}

// generateWordNumberPattern generates word-number pattern (e.g., SAVE2024)
func generateWordNumberPattern(config GeneratorConfig) (string, error) {
	words := []string{"SAVE", "DEAL", "OFFER", "SALE", "BONUS", "GIFT", "SPECIAL", "MEGA", "SUPER", "BEST"}

	// Select random word
	wordIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
	if err != nil {
		return "", err
	}
	word := words[wordIndex.Int64()]

	// Generate random number
	numberLength := config.Length
	if numberLength <= 0 {
		numberLength = 4
	}

	number, err := generateRandomNumber(numberLength)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%s%s%s", config.Prefix, word, number)

	if config.Suffix != "" {
		code += config.Suffix
	}

	return strings.ToUpper(code), nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int, excludeChars string) (string, error) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Remove excluded characters
	for _, char := range excludeChars {
		charset = strings.ReplaceAll(charset, string(char), "")
	}

	if len(charset) == 0 {
		return "", fmt.Errorf("no valid characters available after exclusions")
	}

	result := make([]byte, length)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}

	return string(result), nil
}

// generateRandomNumber generates a random number string of specified length
func generateRandomNumber(length int) (string, error) {
	charset := "0123456789"
	result := make([]byte, length)

	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[index.Int64()]
	}

	return string(result), nil
}

// GenerateExpiryDate generates expiry date based on duration
func GenerateExpiryDate(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

// GenerateSeasonalCode generates seasonal coupon codes
func GenerateSeasonalCode(season string, year int, config GeneratorConfig) (string, error) {
	seasonPrefixes := map[string]string{
		"spring": "SPRING",
		"summer": "SUMMER",
		"autumn": "AUTUMN",
		"winter": "WINTER",
		"newyear": "NY",
		"valentine": "LOVE",
		"easter": "EASTER",
		"halloween": "SPOOKY",
		"christmas": "XMAS",
	}

	prefix, exists := seasonPrefixes[strings.ToLower(season)]
	if !exists {
		prefix = "SPECIAL"
	}

	randomPart, err := generateRandomString(4, config.ExcludeChars)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%s%d%s", prefix, year, randomPart)
	return strings.ToUpper(code), nil
}

// GenerateFlashSaleCode generates flash sale specific codes
func GenerateFlashSaleCode(discountPercent int, config GeneratorConfig) (string, error) {
	prefixes := []string{"FLASH", "QUICK", "RUSH", "SPEED", "FAST"}

	// Select random prefix
	prefixIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(prefixes))))
	if err != nil {
		return "", err
	}
	prefix := prefixes[prefixIndex.Int64()]

	randomPart, err := generateRandomString(3, config.ExcludeChars)
	if err != nil {
		return "", err
	}

	code := fmt.Sprintf("%s%d%s", prefix, discountPercent, randomPart)
	return strings.ToUpper(code), nil
}

// GenerateBulkCodes generates codes in bulk with different patterns
func GenerateBulkCodes(configs []GeneratorConfig) (map[string][]string, error) {
	results := make(map[string][]string)

	for i, config := range configs {
		codes, err := GenerateCodes(config)
		if err != nil {
			return nil, fmt.Errorf("failed to generate codes for config %d: %w", i, err)
		}

		patternName := fmt.Sprintf("pattern_%d", i)
		if config.Pattern != "" {
			patternName = config.Pattern
		}

		results[patternName] = codes
	}

	return results, nil
}

// ValidateCodeFormat validates if a code follows expected format
func ValidateCodeFormat(code string, config GeneratorConfig) bool {
	if len(code) == 0 {
		return false
	}

	// Check for excluded characters
	for _, char := range config.ExcludeChars {
		if strings.Contains(code, string(char)) {
			return false
		}
	}

	// Check prefix/suffix if specified
	if config.Prefix != "" && !strings.HasPrefix(code, strings.ToUpper(config.Prefix)) {
		return false
	}

	if config.Suffix != "" && !strings.HasSuffix(code, strings.ToUpper(config.Suffix)) {
		return false
	}

	return true
}
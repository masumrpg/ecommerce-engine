package coupon

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
)

// GenerateCode generates a single coupon code based on the provided configuration.
// It supports multiple patterns including prefix-based, random, and word-number combinations.
// The function applies length constraints, character exclusions, and formatting rules.
//
// Parameters:
//   - config: GeneratorConfig containing pattern, length, prefix, suffix, and exclusion rules
//
// Returns:
//   - string: generated coupon code in uppercase format
//   - error: nil on success, error if generation fails
//
// Supported patterns:
//   - "PREFIX-XXXXXX": generates codes like "SAVE-ABC123"
//   - "XXXXXXXX": generates purely random codes
//   - "WORD-NUMBER": generates codes like "DEAL2024"
//
// Example:
//
//	config := GeneratorConfig{
//		Pattern: "PREFIX-XXXXXX",
//		Length: 6,
//		Prefix: "SAVE",
//		ExcludeChars: "0O1I",
//	}
//	code, err := GenerateCode(config)
//	// Result: "SAVE-ABC123"
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

// GenerateCodes generates multiple unique coupon codes using the specified configuration.
// It ensures all generated codes are unique within the batch and prevents infinite loops
// by limiting generation attempts. If uniqueness cannot be achieved, it returns partial results.
//
// Parameters:
//   - config: GeneratorConfig with Count field specifying number of codes to generate
//
// Returns:
//   - []string: slice of unique coupon codes
//   - error: nil on success, error if generation fails or partial generation occurs
//
// Behavior:
//   - Returns empty slice if Count is 0
//   - Sets Count to 1 if negative
//   - Limits attempts to Count × 10 to prevent infinite loops
//   - Returns partial results with error if full uniqueness cannot be achieved
//
// Example:
//
//	config := GeneratorConfig{Count: 100, Length: 8, Pattern: "XXXXXXXX"}
//	codes, err := GenerateCodes(config)
//	// Result: ["ABC12345", "DEF67890", ...]
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

// generatePrefixPattern generates coupon codes with prefix pattern format.
// Creates codes in the format "PREFIX-RANDOM" with optional suffix.
// Uses "COUPON" as default prefix if none specified.
//
// Parameters:
//   - config: GeneratorConfig containing prefix, suffix, length, and exclusion settings
//
// Returns:
//   - string: formatted coupon code like "SAVE-ABC123" or "COUPON-XYZ789-SPECIAL"
//   - error: nil on success, error if random string generation fails
//
// Format: PREFIX-RANDOM or PREFIX-RANDOM-SUFFIX
//
// Example:
//   Input: {Prefix: "SAVE", Length: 6, Suffix: "2024"}
//   Output: "SAVE-ABC123-2024"
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

// generateRandomPattern generates purely random coupon codes without separators.
// Concatenates prefix, random string, and suffix directly without hyphens.
//
// Parameters:
//   - config: GeneratorConfig containing prefix, suffix, length, and exclusion settings
//
// Returns:
//   - string: concatenated coupon code like "PREFIXABC123SUFFIX"
//   - error: nil on success, error if random string generation fails
//
// Format: PREFIXRANDOMSUFFIX (no separators)
//
// Example:
//   Input: {Prefix: "DEAL", Length: 4, Suffix: "END"}
//   Output: "DEALX7Y9END"
func generateRandomPattern(config GeneratorConfig) (string, error) {
	randomPart, err := generateRandomString(config.Length, config.ExcludeChars)
	if err != nil {
		return "", err
	}

	code := config.Prefix + randomPart + config.Suffix
	return strings.ToUpper(code), nil
}

// generateWordNumberPattern generates coupon codes using word-number pattern.
// Combines a randomly selected promotional word with a random number sequence.
// Supports configurable prefix and suffix additions.
//
// Parameters:
//   - config: GeneratorConfig containing prefix, suffix, length (for number part), and exclusions
//
// Returns:
//   - string: formatted coupon code like "PREFIXDEAL1234SUFFIX"
//   - error: nil on success, error if random generation fails
//
// Word pool: SAVE, DEAL, OFFER, SALE, BONUS, GIFT, SPECIAL, MEGA, SUPER, BEST
// Number length defaults to 4 if config.Length <= 0
//
// Example:
//   Input: {Prefix: "MEGA", Length: 3, Suffix: "END"}
//   Output: "MEGABONUS123END"
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

// generateRandomString generates a cryptographically secure random string of specified length.
// Uses alphanumeric charset (A-Z, 0-9) with configurable character exclusions.
// Commonly excludes visually similar characters like 0O1I to improve readability.
//
// Parameters:
//   - length: desired length of the generated string
//   - excludeChars: string containing characters to exclude from generation
//
// Returns:
//   - string: random string using allowed characters
//   - error: nil on success, error if no valid characters remain or crypto fails
//
// Default charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//
// Example:
//   generateRandomString(6, "0O1I") → "ABC2EF" (excludes confusing chars)
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

// generateRandomNumber generates a cryptographically secure random numeric string.
// Uses only digits 0-9 to create number-only sequences for coupon codes.
//
// Parameters:
//   - length: desired length of the numeric string
//
// Returns:
//   - string: random numeric string
//   - error: nil on success, error if crypto random generation fails
//
// Example:
//   generateRandomNumber(4) → "7392"
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

// GenerateExpiryDate generates an expiry date by adding the specified duration to the current time.
// This utility function is commonly used when creating coupons with time-based validity.
//
// Parameters:
//   - duration: time.Duration to add to current time (e.g., 24*time.Hour for 1 day)
//
// Returns:
//   - time.Time: expiry date calculated from current time + duration
//
// Example:
//
//	// Create coupon valid for 30 days
//	expiry := GenerateExpiryDate(30 * 24 * time.Hour)
//	// Create coupon valid for 2 hours
//	expiry := GenerateExpiryDate(2 * time.Hour)
func GenerateExpiryDate(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

// GenerateSeasonalCode generates themed coupon codes for seasonal promotions.
// Uses predefined seasonal prefixes combined with year and random components.
// Falls back to "SPECIAL" prefix for unrecognized seasons.
//
// Parameters:
//   - season: season name (spring, summer, autumn, winter, newyear, valentine, easter, halloween, christmas)
//   - year: year to include in the code
//   - config: GeneratorConfig for character exclusions and other settings
//
// Returns:
//   - string: seasonal coupon code like "XMAS2024ABC1"
//   - error: nil on success, error if random generation fails
//
// Supported seasons and their prefixes:
//   spring→SPRING, summer→SUMMER, autumn→AUTUMN, winter→WINTER,
//   newyear→NY, valentine→LOVE, easter→EASTER, halloween→SPOOKY, christmas→XMAS
//
// Example:
//
//	code, err := GenerateSeasonalCode("christmas", 2024, config)
//	// Result: "XMAS2024ABC1"
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

// GenerateFlashSaleCode generates coupon codes specifically for flash sale promotions.
// Combines urgency-themed prefixes with discount percentage and random components.
//
// Parameters:
//   - discountPercent: discount percentage to include in the code
//   - config: GeneratorConfig for character exclusions and other settings
//
// Returns:
//   - string: flash sale coupon code like "FLASH25ABC"
//   - error: nil on success, error if random generation fails
//
// Prefix pool: FLASH, QUICK, RUSH, SPEED, FAST
// Random part length: 3 characters
//
// Example:
//
//	code, err := GenerateFlashSaleCode(25, config)
//	// Result: "RUSH25XYZ"
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

// GenerateBulkCodes generates multiple batches of coupon codes using different configurations.
// Each configuration produces a separate batch of codes, allowing for diverse coupon campaigns
// with different patterns, lengths, and formatting rules in a single operation.
//
// Parameters:
//   - configs: slice of GeneratorConfig, each defining a different code generation pattern
//
// Returns:
//   - map[string][]string: map where keys are pattern names and values are code slices
//   - error: nil on success, error if any configuration fails
//
// Map keys:
//   - Uses config.Pattern if specified, otherwise "pattern_N" where N is the index
//   - Allows easy identification of which codes belong to which pattern
//
// Example:
//
//	configs := []GeneratorConfig{
//		{Pattern: "SUMMER", Count: 50, Length: 6},
//		{Pattern: "FLASH", Count: 100, Length: 8},
//	}
//	results, err := GenerateBulkCodes(configs)
//	// Result: {"SUMMER": ["SUMMER123", ...], "FLASH": ["FLASH456", ...]}
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

// ValidateCodeFormat validates whether a coupon code conforms to the specified format rules.
// Performs validation checks including character exclusions, prefix/suffix requirements,
// and basic format compliance. Used for validating user-entered codes or generated codes.
//
// Parameters:
//   - code: coupon code string to validate
//   - config: GeneratorConfig containing validation rules (exclusions, prefix, suffix)
//
// Returns:
//   - bool: true if code passes all validation checks, false otherwise
//
// Validation checks:
//   - Code is not empty
//   - Code does not contain any excluded characters
//   - Code starts with required prefix (if specified)
//   - Code ends with required suffix (if specified)
//
// Example:
//
//	config := GeneratorConfig{Prefix: "SAVE", ExcludeChars: "0O1I"}
//	isValid := ValidateCodeFormat("SAVE-ABC123", config)
//	// Returns: true (valid format)
//	isValid := ValidateCodeFormat("SAVE-0BC123", config)
//	// Returns: false (contains excluded char '0')
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
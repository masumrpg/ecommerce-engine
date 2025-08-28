// Package utils provides utility functions for e-commerce applications.
// This package includes various generators for IDs, tokens, codes, and other
// commonly needed identifiers in e-commerce systems.
//
// The package offers the following main functionalities:
//   - ID generation (UUID, sequential, timestamp-based, hash-based)
//   - Coupon code generation with customizable patterns
//   - Secure password generation with various options
//   - Token generation (API keys, OTP, verification codes)
//   - Reference number generation (orders, invoices, transactions)
//   - Barcode generation (EAN-13, UPC, SKU)
//   - URL slug generation
//   - Color code generation
//   - Cryptographic utilities (nonce, salt, checksum)
package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"strconv"
	"strings"
	"time"
)

// IDGenerator provides various ID generation methods with support for
// sequential numbering and custom prefixes. It maintains an internal
// counter that starts from the current nanosecond timestamp to ensure
// uniqueness across different instances.
//
// Example usage:
//
//	gen := NewIDGenerator("USER")
//	id1 := gen.GenerateSequentialID() // Returns "USER-1234567890123456789"
//	id2 := gen.GenerateSequentialID() // Returns "USER-1234567890123456790"
type IDGenerator struct {
	prefix  string // Optional prefix for generated IDs
	counter int64  // Internal counter for sequential ID generation
}

// NewIDGenerator creates a new ID generator with an optional prefix.
// The internal counter is initialized with the current nanosecond timestamp
// to ensure uniqueness across different generator instances.
//
// Parameters:
//   - prefix: Optional prefix to prepend to generated IDs. Can be empty.
//
// Returns:
//   - *IDGenerator: A new ID generator instance.
//
// Example:
//
//	gen := NewIDGenerator("ORDER")
//	// Creates a generator that will produce IDs like "ORDER-1234567890123456789"
func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{
		prefix:  prefix,
		counter: time.Now().UnixNano(),
	}
}

// GenerateUUID generates a UUID-like string in the format 8-4-4-4-12.
// This function uses cryptographically secure random number generation
// when possible, falling back to time-based generation if crypto/rand fails.
//
// The generated UUID follows the standard format but is not RFC 4122 compliant.
// For production systems requiring strict UUID compliance, consider using
// a dedicated UUID library.
//
// Returns:
//   - string: A UUID-like string (e.g., "550e8400-e29b-41d4-a716-446655440000")
//
// Example:
//
//	uuid := GenerateUUID()
//	fmt.Println(uuid) // Output: "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		// Fallback to time-based generation
		return fmt.Sprintf("%d-%d-%d-%d-%d",
			time.Now().UnixNano(),
			mathRand.Int63n(10000),
			mathRand.Int63n(10000),
			mathRand.Int63n(10000),
			mathRand.Int63n(100000000))
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

// GenerateShortID generates a short alphanumeric ID of the specified length.
// The generated ID contains a mix of uppercase letters, lowercase letters,
// and digits, providing good entropy for short identifiers.
//
// This function uses cryptographically secure random number generation
// to ensure unpredictability of the generated IDs.
//
// Parameters:
//   - length: The desired length of the generated ID. Must be positive.
//
// Returns:
//   - string: An alphanumeric ID of the specified length.
//
// Example:
//
//	id := GenerateShortID(8)
//	fmt.Println(id) // Output: "aB3dE7gH" (example)
func GenerateShortID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// GenerateNumericID generates a numeric ID of the specified length.
// The generated ID consists only of digits (0-9) and does not start with zero
// to ensure the full length is maintained when used as a number.
//
// If length is less than or equal to 0, it defaults to 8 digits.
//
// Parameters:
//   - length: The desired length of the numeric ID. Defaults to 8 if <= 0.
//
// Returns:
//   - string: A numeric ID of the specified length (e.g., "12345678").
//
// Example:
//
//	id := GenerateNumericID(6)
//	fmt.Println(id) // Output: "123456" (example, will not start with 0)
func GenerateNumericID(length int) string {
	if length <= 0 {
		length = 8
	}

	min := int64(1)
	max := int64(1)
	for i := 0; i < length; i++ {
		max *= 10
	}
	max--
	for i := 1; i < length; i++ {
		min *= 10
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
	return fmt.Sprintf("%0*d", length, n.Int64()+min)
}

// GenerateSequentialID generates a sequential ID with the configured prefix.
// Each call increments the internal counter, ensuring unique sequential IDs.
// The counter is thread-safe for single-instance usage but not across
// multiple goroutines.
//
// Returns:
//   - string: A sequential ID in the format "prefix-counter" or just "counter" if no prefix.
//
// Example:
//
//	gen := NewIDGenerator("ORDER")
//	id1 := gen.GenerateSequentialID() // "ORDER-1234567890123456790"
//	id2 := gen.GenerateSequentialID() // "ORDER-1234567890123456791"
func (g *IDGenerator) GenerateSequentialID() string {
	g.counter++
	if g.prefix != "" {
		return fmt.Sprintf("%s-%d", g.prefix, g.counter)
	}
	return fmt.Sprintf("%d", g.counter)
}

// GenerateTimestampID generates an ID based on the current nanosecond timestamp.
// This provides a time-ordered ID that can be useful for sorting and
// approximate chronological ordering.
//
// Note: IDs generated in rapid succession may have the same timestamp
// on some systems, so this should not be relied upon for uniqueness
// in high-concurrency scenarios.
//
// Returns:
//   - string: A timestamp-based ID (nanoseconds since Unix epoch).
//
// Example:
//
//	id := GenerateTimestampID()
//	fmt.Println(id) // Output: "1640995200123456789"
func GenerateTimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GenerateTimestampIDWithPrefix generates a timestamp-based ID with an optional prefix.
// Combines the benefits of timestamp ordering with a readable prefix for
// categorization and identification.
//
// Parameters:
//   - prefix: Optional prefix to prepend. If empty, returns same as GenerateTimestampID().
//
// Returns:
//   - string: A prefixed timestamp ID in format "prefix-timestamp" or just timestamp.
//
// Example:
//
//	id := GenerateTimestampIDWithPrefix("TXN")
//	fmt.Println(id) // Output: "TXN-1640995200123456789"
func GenerateTimestampIDWithPrefix(prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}
	return GenerateTimestampID()
}

// GenerateHashID generates a deterministic hash-based ID from the input string.
// Uses SHA-256 hashing to create a consistent 16-character hexadecimal ID
// that will always be the same for the same input.
//
// This is useful for creating reproducible IDs based on content or
// for deduplication purposes.
//
// Parameters:
//   - input: The string to hash. Can be any string value.
//
// Returns:
//   - string: A 16-character hexadecimal hash ID.
//
// Example:
//
//	id := GenerateHashID("user@example.com")
//	fmt.Println(id) // Output: "a1b2c3d4e5f67890" (consistent for same input)
func GenerateHashID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])[:16] // Return first 16 characters
}

// GenerateCustomID generates a custom ID with flexible formatting options.
// Allows combining prefix, timestamp, and random components based on requirements.
// Components are joined with hyphens when multiple parts are included.
//
// Format variations:
//   - prefix only: "PREFIX"
//   - prefix + timestamp: "PREFIX-1640995200"
//   - prefix + random: "PREFIX-aB3dE7"
//   - prefix + timestamp + random: "PREFIX-1640995200-aB3dE7"
//   - timestamp + random: "1640995200-aB3dE7"
//
// Parameters:
//   - prefix: Optional prefix string. Can be empty.
//   - includeTimestamp: Whether to include Unix timestamp (seconds).
//   - includeRandom: Whether to include a 6-character random component.
//
// Returns:
//   - string: A custom formatted ID, or 8-character random string if no components specified.
//
// Example:
//
//	id := GenerateCustomID("USER", true, true)
//	fmt.Println(id) // Output: "USER-1640995200-aB3dE7"
func GenerateCustomID(prefix string, includeTimestamp, includeRandom bool) string {
	parts := []string{}

	if prefix != "" {
		parts = append(parts, prefix)
	}

	if includeTimestamp {
		parts = append(parts, strconv.FormatInt(time.Now().Unix(), 10))
	}

	if includeRandom {
		parts = append(parts, GenerateShortID(6))
	}

	if len(parts) == 0 {
		return GenerateShortID(8)
	}

	return strings.Join(parts, "-")
}

// CouponCodeGenerator provides customizable coupon code generation with
// support for character set filtering and pattern-based generation.
// By default, it excludes visually similar characters to reduce user confusion.
//
// Example usage:
//
//	gen := NewCouponCodeGenerator(8)
//	code := gen.GenerateCouponCode() // Returns "ABCD2345" (example)
//
//	// Custom pattern
//	patternCode := gen.GenerateCouponCodeWithPattern("SAVE-XXX") // Returns "SAVE-ABC" (example)
type CouponCodeGenerator struct {
	length   int      // Length of generated coupon codes
	charset  string   // Character set to use for generation
	excluded []string // Characters to exclude from generation
}

// NewCouponCodeGenerator creates a new coupon code generator with the specified length.
// The generator is initialized with a default character set of uppercase letters
// and digits, excluding visually similar characters (0, O, 1, I, L) to prevent
// user confusion when entering codes.
//
// Parameters:
//   - length: The length of coupon codes to generate. Must be positive.
//
// Returns:
//   - *CouponCodeGenerator: A new coupon code generator instance.
//
// Example:
//
//	gen := NewCouponCodeGenerator(8)
//	// Creates a generator for 8-character coupon codes
func NewCouponCodeGenerator(length int) *CouponCodeGenerator {
	return &CouponCodeGenerator{
		length:  length,
		charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		excluded: []string{"0", "O", "1", "I", "L"}, // Exclude confusing characters
	}
}

// SetCharset sets a custom character set for coupon code generation.
// This allows you to define exactly which characters can appear in
// generated coupon codes.
//
// Parameters:
//   - charset: String containing all allowed characters for generation.
//
// Example:
//
//	gen.SetCharset("ABCDEF123456") // Only use these specific characters
func (g *CouponCodeGenerator) SetCharset(charset string) {
	g.charset = charset
}

// SetExcludedChars sets characters to exclude from the current character set.
// This is useful for removing problematic characters while keeping the
// rest of the character set intact.
//
// Parameters:
//   - excluded: Slice of strings representing characters to exclude.
//
// Example:
//
//	gen.SetExcludedChars([]string{"0", "O", "Q"}) // Remove confusing characters
func (g *CouponCodeGenerator) SetExcludedChars(excluded []string) {
	g.excluded = excluded
}

// GenerateCouponCode generates a random coupon code using the configured
// character set and length. Excluded characters are automatically filtered
// out from the generation process.
//
// Returns:
//   - string: A randomly generated coupon code of the configured length.
//
// Example:
//
//	gen := NewCouponCodeGenerator(6)
//	code := gen.GenerateCouponCode() // Returns "ABC123" (example)
func (g *CouponCodeGenerator) GenerateCouponCode() string {
	charset := g.getFilteredCharset()
	code := make([]byte, g.length)

	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}

	return string(code)
}

// GenerateCouponCodeWithPattern generates a coupon code following a specific pattern.
// The pattern uses 'X' as placeholders that get replaced with random characters
// from the configured character set. All other characters in the pattern are preserved.
//
// Parameters:
//   - pattern: Pattern string where 'X' represents positions for random characters.
//
// Returns:
//   - string: A coupon code following the specified pattern.
//
// Example:
//
//	gen := NewCouponCodeGenerator(8)
//	code := gen.GenerateCouponCodeWithPattern("SAVE-XXX") // Returns "SAVE-A2B" (example)
//	code2 := gen.GenerateCouponCodeWithPattern("XX-XX-XX") // Returns "AB-C3-D4" (example)
func (g *CouponCodeGenerator) GenerateCouponCodeWithPattern(pattern string) string {
	charset := g.getFilteredCharset()
	result := []rune(pattern)

	for i, char := range result {
		if char == 'X' {
			n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			result[i] = rune(charset[n.Int64()])
		}
	}

	return string(result)
}

// GenerateBatchCouponCodes generates multiple unique coupon codes in a single operation.
// This method ensures all generated codes are unique within the batch by using
// a map to track duplicates and regenerating when necessary.
//
// Parameters:
//   - count: Number of unique coupon codes to generate.
//
// Returns:
//   - []string: Slice of unique coupon codes.
//
// Example:
//
//	gen := NewCouponCodeGenerator(6)
//	codes := gen.GenerateBatchCouponCodes(5)
//	// Returns ["ABC123", "DEF456", "GHI789", "JKL234", "MNP567"] (example)
func (g *CouponCodeGenerator) GenerateBatchCouponCodes(count int) []string {
	codes := make([]string, 0, count)
	generated := make(map[string]bool)

	for len(codes) < count {
		code := g.GenerateCouponCode()
		if !generated[code] {
			codes = append(codes, code)
			generated[code] = true
		}
	}

	return codes
}

// getFilteredCharset returns the character set with all excluded characters removed.
// This is an internal helper method used by the generation functions to ensure
// excluded characters don't appear in generated codes.
//
// Returns:
//   - string: Filtered character set with excluded characters removed.
func (g *CouponCodeGenerator) getFilteredCharset() string {
	charset := g.charset
	for _, excluded := range g.excluded {
		charset = strings.ReplaceAll(charset, excluded, "")
	}
	return charset
}

// PasswordGenerator provides secure password generation with customizable
// character sets and complexity requirements. It supports uppercase letters,
// lowercase letters, digits, and symbols, with options to exclude ambiguous
// characters for better usability.
//
// Example usage:
//
//	gen := NewPasswordGenerator(12)
//	password := gen.GeneratePassword() // Returns "AbC3dEf7HiJk" (example)
//
//	// Custom configuration
//	gen.SetOptions(true, true, true, true, true)
//	password = gen.GeneratePassword() // Returns "Ab@3dE#7Hi!k" (example)
type PasswordGenerator struct {
	length           int  // Length of generated passwords
	includeUppercase bool // Include uppercase letters (A-Z)
	includeLowercase bool // Include lowercase letters (a-z)
	includeNumbers   bool // Include digits (0-9)
	includeSymbols   bool // Include symbols (!@#$%^&*)
	excludeAmbiguous bool // Exclude ambiguous characters (0,O,1,l,I)
}

// NewPasswordGenerator creates a new password generator with the specified length.
// The generator is initialized with secure defaults: uppercase and lowercase letters,
// digits enabled, symbols disabled, and ambiguous characters excluded for better
// user experience.
//
// Parameters:
//   - length: The length of passwords to generate. Should be at least 8 for security.
//
// Returns:
//   - *PasswordGenerator: A new password generator instance with secure defaults.
//
// Example:
//
//	gen := NewPasswordGenerator(12)
//	// Creates a generator for 12-character passwords with letters and digits
func NewPasswordGenerator(length int) *PasswordGenerator {
	return &PasswordGenerator{
		length:           length,
		includeUppercase: true,
		includeLowercase: true,
		includeNumbers:   true,
		includeSymbols:   false,
		excludeAmbiguous: true,
	}
}

// SetOptions configures all password generation options in a single call.
// This is a convenience method for setting multiple options at once.
//
// Parameters:
//   - uppercase: true to include uppercase letters (A-Z)
//   - lowercase: true to include lowercase letters (a-z)
//   - numbers: true to include digits (0-9)
//   - symbols: true to include symbols (!@#$%^&*)
//   - excludeAmbiguous: true to exclude visually similar characters (0,O,1,l,I)
//
// Example:
//
//	gen.SetOptions(true, true, true, true, false)
//	// Enable all character types including ambiguous characters
func (g *PasswordGenerator) SetOptions(uppercase, lowercase, numbers, symbols, excludeAmbiguous bool) {
	g.includeUppercase = uppercase
	g.includeLowercase = lowercase
	g.includeNumbers = numbers
	g.includeSymbols = symbols
	g.excludeAmbiguous = excludeAmbiguous
}

// GeneratePassword generates a secure password using the configured character sets
// and length. The password uses cryptographically secure random number generation
// to ensure unpredictability.
//
// Returns:
//   - string: A randomly generated password meeting the configured requirements,
//             or empty string if no character sets are enabled.
//
// Example:
//
//	gen := NewPasswordGenerator(12)
//	password := gen.GeneratePassword() // Returns "AbC3dEf7HiJk" (example)
//
//	gen.SetOptions(true, true, true, true, true)
//	password = gen.GeneratePassword() // Returns "Ab@3dE#7Hi!k" (example)
func (g *PasswordGenerator) GeneratePassword() string {
	charset := g.buildCharset()
	if len(charset) == 0 {
		return ""
	}

	password := make([]byte, g.length)
	for i := range password {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[n.Int64()]
	}

	return string(password)
}

// buildCharset builds the character set based on the current configuration options.
// This is an internal helper method that combines all enabled character types
// while respecting the excludeAmbiguous setting.
//
// Returns:
//   - string: Combined character set containing all enabled characters.
func (g *PasswordGenerator) buildCharset() string {
	var charset strings.Builder

	if g.includeUppercase {
		if g.excludeAmbiguous {
			charset.WriteString("ABCDEFGHJKLMNPQRSTUVWXYZ") // Exclude I, O
		} else {
			charset.WriteString("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		}
	}

	if g.includeLowercase {
		if g.excludeAmbiguous {
			charset.WriteString("abcdefghjkmnpqrstuvwxyz") // Exclude i, l, o
		} else {
			charset.WriteString("abcdefghijklmnopqrstuvwxyz")
		}
	}

	if g.includeNumbers {
		if g.excludeAmbiguous {
			charset.WriteString("23456789") // Exclude 0, 1
		} else {
			charset.WriteString("0123456789")
		}
	}

	if g.includeSymbols {
		if g.excludeAmbiguous {
			charset.WriteString("!@#$%^&*()_+-=") // Exclude similar looking symbols
		} else {
			charset.WriteString("!@#$%^&*()_+-=[]{}|;:,.<>?")
		}
	}

	return charset.String()
}

// TokenGenerator provides secure token generation for various purposes including
// API keys, session tokens, and authentication tokens. It supports multiple
// output formats (hex, base64, alphanumeric) to meet different requirements.
//
// Example usage:
//
//	gen := NewTokenGenerator()
//	token := gen.GenerateSecureToken(32) // Returns "a1b2c3d4e5f6..." (64 chars hex)
//	apiKey := gen.GenerateAPIKey("api") // Returns "api_a1b2c3d4..."
type TokenGenerator struct{}

// NewTokenGenerator creates a new token generator instance.
// The generator provides various methods for creating secure tokens
// suitable for different authentication and identification purposes.
//
// Returns:
//   - *TokenGenerator: A new token generator instance.
//
// Example:
//
//	gen := NewTokenGenerator()
//	// Use gen.GenerateSecureToken(), gen.GenerateAPIKey(), etc.
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// GenerateSecureToken generates a cryptographically secure token using
// crypto/rand for maximum security. The token is hex-encoded and suitable
// for security-sensitive applications such as API keys and session tokens.
//
// Parameters:
//   - length: Desired length of the final token string.
//
// Returns:
//   - string: A securely generated hex token of the specified length.
//
// Example:
//
//	gen := NewTokenGenerator()
//	token := gen.GenerateSecureToken(32) // Returns 32-character hex token
func (g *TokenGenerator) GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to less secure method
		return GenerateShortID(length)
	}
	return hex.EncodeToString(bytes)[:length]
}

// GenerateAPIKey generates an API key with an optional prefix.
// The API key consists of a prefix (if provided) followed by an underscore
// and a 32-character secure token. This format is commonly used for
// API authentication in web services.
//
// Parameters:
//   - prefix: Optional prefix to identify the API key type or service.
//
// Returns:
//   - string: API key in format "prefix_token" or just "token" if no prefix.
//
// Example:
//
//	gen := NewTokenGenerator()
//	apiKey := gen.GenerateAPIKey("api") // Returns "api_a1b2c3d4..."
//	apiKey = gen.GenerateAPIKey("")    // Returns "a1b2c3d4..." (no prefix)
func (g *TokenGenerator) GenerateAPIKey(prefix string) string {
	token := g.GenerateSecureToken(32)
	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, token)
	}
	return token
}

// GenerateOTP generates a one-time password (OTP) consisting of numeric digits.
// OTPs are commonly used for two-factor authentication and account verification.
// The default length is 6 digits if an invalid length is provided.
//
// Parameters:
//   - length: Number of digits in the OTP. Defaults to 6 if <= 0.
//
// Returns:
//   - string: Numeric OTP of the specified length.
//
// Example:
//
//	gen := NewTokenGenerator()
//	otp := gen.GenerateOTP(6) // Returns "123456" (example)
//	otp = gen.GenerateOTP(4)  // Returns "7890" (example)
func (g *TokenGenerator) GenerateOTP(length int) string {
	if length <= 0 {
		length = 6
	}
	return GenerateNumericID(length)
}

// GenerateVerificationCode generates a 6-digit numeric verification code.
// These codes are commonly used for email verification, phone verification,
// and other account security processes. The code is always 6 digits long.
//
// Returns:
//   - string: A 6-digit numeric verification code.
//
// Example:
//
//	gen := NewTokenGenerator()
//	code := gen.GenerateVerificationCode() // Returns "123456" (example)
func (g *TokenGenerator) GenerateVerificationCode() string {
	return GenerateNumericID(6)
}

// ReferenceGenerator provides reference number generation for orders, invoices,
// transactions, and other business documents. It supports customizable prefixes,
// suffixes, and lengths to create systematic and human-readable reference numbers.
//
// Example usage:
//
//	gen := NewReferenceGenerator("ORD", "SHOP", 8)
//	ref := gen.GenerateOrderReference() // Returns "ORD-20240115-12345678-SHOP"
//
//	gen = NewReferenceGenerator("INV", "", 6)
//	ref = gen.GenerateInvoiceReference() // Returns "INV-2024-01-123456"
type ReferenceGenerator struct {
	prefix string // Prefix to identify reference type
	suffix string // Optional suffix for additional identification
	length int    // Length of the random numeric component
}

// NewReferenceGenerator creates a new reference generator with the specified
// configuration. Reference numbers are commonly used for tracking business
// documents and maintaining systematic record keeping.
//
// Parameters:
//   - prefix: String prefix to identify the reference type (e.g., "ORD", "INV", "TXN").
//   - suffix: Optional suffix for additional identification. Can be empty.
//   - length: Length of the random numeric component in the reference.
//
// Returns:
//   - *ReferenceGenerator: A new reference generator instance.
//
// Example:
//
//	gen := NewReferenceGenerator("ORD", "SHOP", 8) // Order references
//	gen = NewReferenceGenerator("INV", "", 6)       // Invoice references
//	gen = NewReferenceGenerator("TXN", "PAY", 10)   // Transaction references
func NewReferenceGenerator(prefix, suffix string, length int) *ReferenceGenerator {
	return &ReferenceGenerator{
		prefix: prefix + "-",
		suffix: "-" + suffix,
		length: length,
	}
}

// GenerateOrderReference generates an order reference number with date and
// random numeric components. The format includes the current date (YYYYMMDD)
// for easy chronological sorting and identification.
//
// Returns:
//   - string: Order reference in format "PREFIX-YYYYMMDD-NUMBERS-SUFFIX" or
//             "PREFIX-YYYYMMDD-NUMBERS" if no suffix is configured.
//
// Example:
//
//	gen := NewReferenceGenerator("ORD", "SHOP", 8)
//	ref := gen.GenerateOrderReference() // Returns "ORD-20240115-12345678-SHOP"
//
//	gen = NewReferenceGenerator("ORD", "", 6)
//	ref = gen.GenerateOrderReference() // Returns "ORD-20240115-123456"
func (g *ReferenceGenerator) GenerateOrderReference() string {
	timestamp := time.Now().Format("20060102")
	random := GenerateNumericID(g.length)

	parts := []string{}
	if g.prefix != "" {
		parts = append(parts, g.prefix)
	}
	parts = append(parts, timestamp, random)
	if g.suffix != "" {
		parts = append(parts, g.suffix)
	}

	return strings.Join(parts, "")
}

// GenerateInvoiceReference generates an invoice reference number with year-month
// and random numeric components. The format uses YYYY-MM for monthly grouping
// of invoices, which is common in accounting systems.
//
// Returns:
//   - string: Invoice reference in format "PREFIX-YYYY-MM-NUMBERS-SUFFIX" or
//             "PREFIX-YYYY-MM-NUMBERS" if no suffix is configured.
//
// Example:
//
//	gen := NewReferenceGenerator("INV", "STORE", 6)
//	ref := gen.GenerateInvoiceReference() // Returns "INV-2024-01-123456-STORE"
//
//	gen = NewReferenceGenerator("INV", "", 8)
//	ref = gen.GenerateInvoiceReference() // Returns "INV-2024-01-12345678"
func (g *ReferenceGenerator) GenerateInvoiceReference() string {
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	sequence := GenerateNumericID(g.length)

	parts := []string{}
	if g.prefix != "" {
		parts = append(parts, g.prefix)
	}
	parts = append(parts, year, month, sequence)
	if g.suffix != "" {
		parts = append(parts, g.suffix)
	}

	return strings.Join(parts, "")
}

// GenerateTransactionReference generates a transaction reference number with
// Unix timestamp and random alphanumeric components. The result is converted
// to uppercase for consistency and readability.
//
// Returns:
//   - string: Transaction reference in uppercase format "PREFIX-TIMESTAMP-RANDOM-SUFFIX"
//             or "PREFIX-TIMESTAMP-RANDOM" if no suffix is configured.
//
// Example:
//
//	gen := NewReferenceGenerator("TXN", "PAY", 8)
//	ref := gen.GenerateTransactionReference() // Returns "TXN-1640995200-ABC123DE-PAY"
//
//	gen = NewReferenceGenerator("TXN", "", 6)
//	ref = gen.GenerateTransactionReference() // Returns "TXN-1640995200-ABC123"
func (g *ReferenceGenerator) GenerateTransactionReference() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	random := GenerateShortID(8)

	parts := []string{}
	if g.prefix != "" {
		parts = append(parts, g.prefix)
	}
	parts = append(parts, timestamp, random)
	if g.suffix != "" {
		parts = append(parts, g.suffix)
	}

	return strings.ToUpper(strings.Join(parts, ""))
}

// BarcodeGenerator provides barcode number generation for product identification
// and inventory management. It supports common barcode formats including EAN13,
// UPC, and SKU codes with proper check digit calculation.
//
// Example usage:
//
//	gen := NewBarcodeGenerator()
//	ean13 := gen.GenerateEAN13() // Returns "1234567890123" (with check digit)
//	upc := gen.GenerateUPC()     // Returns "123456789012" (with check digit)
//	sku := gen.GenerateSKU("ELEC", "MOB") // Returns "ELE-MOB-123456"
type BarcodeGenerator struct{}

// NewBarcodeGenerator creates a new barcode generator instance.
// The generator can create valid barcode numbers with proper check digits
// for product identification and retail systems.
//
// Returns:
//   - *BarcodeGenerator: A new barcode generator instance.
//
// Example:
//
//	gen := NewBarcodeGenerator()
//	// Use gen.GenerateEAN13(), gen.GenerateUPC(), gen.GenerateSKU(), etc.
func NewBarcodeGenerator() *BarcodeGenerator {
	return &BarcodeGenerator{}
}

// GenerateEAN13 generates a valid EAN-13 barcode with proper check digit calculation.
// EAN-13 is the European standard for product identification and consists of
// 13 digits including a calculated check digit for validation.
//
// Returns:
//   - string: A 13-digit EAN-13 barcode with valid check digit.
//
// Example:
//
//	gen := NewBarcodeGenerator()
//	barcode := gen.GenerateEAN13() // Returns "1234567890123" (example with valid check digit)
func (g *BarcodeGenerator) GenerateEAN13() string {
	// Generate 12 digits, the 13th would be check digit
	code := GenerateNumericID(12)
	checkDigit := g.calculateEAN13CheckDigit(code)
	return code + strconv.Itoa(checkDigit)
}

// GenerateUPC generates a valid UPC (Universal Product Code) barcode with
// proper check digit calculation. UPC is the North American standard for
// product identification and consists of 12 digits.
//
// Returns:
//   - string: A 12-digit UPC barcode with valid check digit.
//
// Example:
//
//	gen := NewBarcodeGenerator()
//	barcode := gen.GenerateUPC() // Returns "123456789012" (example with valid check digit)
func (g *BarcodeGenerator) GenerateUPC() string {
	// Generate 11 digits, the 12th would be check digit
	code := GenerateNumericID(11)
	checkDigit := g.calculateUPCCheckDigit(code)
	return code + strconv.Itoa(checkDigit)
}

// GenerateSKU generates a Stock Keeping Unit (SKU) code for inventory management.
// The SKU combines category and subcategory abbreviations with a random numeric
// component to create unique product identifiers.
//
// Parameters:
//   - category: Product category name (first 3 characters used).
//   - subcategory: Product subcategory name (first 3 characters used).
//
// Returns:
//   - string: SKU in format "CAT-SUB-123456" or variations based on provided parameters.
//
// Example:
//
//	gen := NewBarcodeGenerator()
//	sku := gen.GenerateSKU("Electronics", "Mobile") // Returns "ELE-MOB-123456"
//	sku = gen.GenerateSKU("Books", "")              // Returns "BOO-123456"
//	sku = gen.GenerateSKU("", "")                   // Returns "123456"
func (g *BarcodeGenerator) GenerateSKU(category, subcategory string) string {
	parts := []string{}
	if category != "" {
		parts = append(parts, strings.ToUpper(category[:MinInt(3, len(category))]))
	}
	if subcategory != "" {
		parts = append(parts, strings.ToUpper(subcategory[:MinInt(3, len(subcategory))]))
	}
	parts = append(parts, GenerateNumericID(6))
	return strings.Join(parts, "-")
}

// calculateEAN13CheckDigit calculates the check digit for EAN-13 barcodes
// using the standard algorithm. This is an internal helper method that
// implements the EAN-13 check digit calculation formula.
//
// Parameters:
//   - code: String containing the first 12 digits of the EAN-13 barcode.
//
// Returns:
//   - int: The calculated check digit (0-9).
func (g *BarcodeGenerator) calculateEAN13CheckDigit(code string) int {
	sum := 0
	for i, digit := range code {
		num, _ := strconv.Atoi(string(digit))
		if i%2 == 0 {
			sum += num
		} else {
			sum += num * 3
		}
	}
	return (10 - (sum % 10)) % 10
}

// calculateUPCCheckDigit calculates the check digit for UPC barcodes
// using the standard algorithm. This is an internal helper method that
// implements the UPC check digit calculation formula.
//
// Parameters:
//   - code: String containing the first 11 digits of the UPC barcode.
//
// Returns:
//   - int: The calculated check digit (0-9).
func (g *BarcodeGenerator) calculateUPCCheckDigit(code string) int {
	sum := 0
	for i, digit := range code {
		num, _ := strconv.Atoi(string(digit))
		if i%2 == 0 {
			sum += num * 3
		} else {
			sum += num
		}
	}
	return (10 - (sum % 10)) % 10
}

// SlugGenerator provides URL slug generation for creating SEO-friendly URLs
// from text input. It handles text normalization, special character removal,
// and length constraints to create clean, web-safe URL segments.
//
// Example usage:
//
//	gen := NewSlugGenerator()
//	slug := gen.GenerateSlug("Hello World! This is a Test") // Returns "hello-world-this-is-a-test"
//	slug = gen.GenerateUniqueSlug("Product Name", existingSlugs) // Returns unique slug
type SlugGenerator struct{}

// NewSlugGenerator creates a new slug generator instance.
// Slugs are commonly used for creating SEO-friendly URLs and file names.
//
// Returns:
//   - *SlugGenerator: A new slug generator instance.
//
// Example:
//
//	gen := NewSlugGenerator()  // Standard web slugs
//	// Use gen.GenerateSlug(), gen.GenerateUniqueSlug(), etc.
func NewSlugGenerator() *SlugGenerator {
	return &SlugGenerator{}
}

// GenerateSlug generates a URL-friendly slug from the input text.
// The process includes converting to lowercase, removing special characters,
// replacing spaces with hyphens, and normalizing the result for web use.
//
// Parameters:
//   - text: Input text to convert to a slug.
//
// Returns:
//   - string: URL-friendly slug with special characters removed and spaces replaced.
//
// Example:
//
//	gen := NewSlugGenerator()
//	slug := gen.GenerateSlug("Hello World! This is a Test") // Returns "hello-world-this-is-a-test"
//	slug = gen.GenerateSlug("Product #123 (New)")          // Returns "product-123-new"
func (g *SlugGenerator) GenerateSlug(text string) string {
	// Convert to lowercase
	slug := strings.ToLower(text)

	// Replace spaces and special characters with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove special characters (keep only alphanumeric and hyphens)
	var result strings.Builder
	for _, char := range slug {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-' {
			result.WriteRune(char)
		}
	}

	slug = result.String()

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// GenerateUniqueSlug generates a unique slug by appending a number if needed.
// This method ensures the generated slug doesn't conflict with existing slugs
// by checking against a provided list and incrementing a counter until unique.
//
// Parameters:
//   - text: Input text to convert to a slug.
//   - existingSlugs: List of existing slugs to check for conflicts.
//
// Returns:
//   - string: Unique URL-friendly slug, potentially with numeric suffix.
//
// Example:
//
//	gen := NewSlugGenerator()
//	existing := []string{"hello-world", "hello-world-1"}
//	slug := gen.GenerateUniqueSlug("Hello World", existing) // Returns "hello-world-2"
func (g *SlugGenerator) GenerateUniqueSlug(text string, existingSlugs []string) string {
	baseSlug := g.GenerateSlug(text)
	slug := baseSlug
	counter := 1

	// Check if slug exists and increment counter until unique
	for g.slugExists(slug, existingSlugs) {
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}

	return slug
}

// slugExists checks if a slug exists in the provided list
func (g *SlugGenerator) slugExists(slug string, existingSlugs []string) bool {
	for _, existing := range existingSlugs {
		if existing == slug {
			return true
		}
	}
	return false
}

// ColorGenerator provides color code generation for UI themes, product variants,
// and design systems. It supports multiple color formats including hex, RGB,
// and specialized color types like pastel and dark colors.
//
// Example usage:
//
//	gen := NewColorGenerator()
//	hex := gen.GenerateHexColor()     // Returns "#A1B2C3"
//	r, g, b := gen.GenerateRGBColor() // Returns RGB values
//	pastel := gen.GeneratePastelColor() // Returns pastel hex color
type ColorGenerator struct{}

// NewColorGenerator creates a new color generator instance.
// The generator provides methods for creating colors in various formats
// suitable for web development, design systems, and product customization.
//
// Returns:
//   - *ColorGenerator: A new color generator instance.
//
// Example:
//
//	gen := NewColorGenerator()
//	// Use gen.GenerateHexColor(), gen.GenerateRGBColor(), etc.
func NewColorGenerator() *ColorGenerator {
	return &ColorGenerator{}
}

// GenerateHexColor generates a random hexadecimal color code.
// The color is suitable for web development and CSS styling.
//
// Returns:
//   - string: Hex color code in format "#RRGGBB".
//
// Example:
//
//	gen := NewColorGenerator()
//	color := gen.GenerateHexColor() // Returns "#A1B2C3" (example)
func (g *ColorGenerator) GenerateHexColor() string {
	r := RandomInt(0, 255)
	green := RandomInt(0, 255)
	b := RandomInt(0, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, green, b)
}

// GenerateRGBColor generates random RGB color values.
// Returns individual red, green, and blue components as integers
// suitable for programmatic color manipulation.
//
// Returns:
//   - r: Red component (0-255)
//   - green: Green component (0-255)
//   - b: Blue component (0-255)
//
// Example:
//
//	gen := NewColorGenerator()
//	r, g, b := gen.GenerateRGBColor() // Returns (161, 178, 195) (example)
func (g *ColorGenerator) GenerateRGBColor() (r, green, b int) {
	return RandomInt(0, 255), RandomInt(0, 255), RandomInt(0, 255)
}

// GeneratePastelColor generates a pastel color with soft, light tones.
// Pastel colors have higher RGB values (127-255) creating gentle,
// muted colors suitable for backgrounds and subtle UI elements.
//
// Returns:
//   - string: Hex color code for a pastel color.
//
// Example:
//
//	gen := NewColorGenerator()
//	color := gen.GeneratePastelColor() // Returns "#E1C2F3" (example)
func (g *ColorGenerator) GeneratePastelColor() string {
	r := RandomInt(127, 255)
	green := RandomInt(127, 255)
	b := RandomInt(127, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, green, b)
}

// GenerateDarkColor generates a dark color with low brightness values.
// Dark colors have lower RGB values (0-127) creating deep, rich colors
// suitable for text, borders, and high-contrast UI elements.
//
// Returns:
//   - string: Hex color code for a dark color.
//
// Example:
//
//	gen := NewColorGenerator()
//	color := gen.GenerateDarkColor() // Returns "#2A1B3C" (example)
func (g *ColorGenerator) GenerateDarkColor() string {
	r := RandomInt(0, 127)
	green := RandomInt(0, 127)
	b := RandomInt(0, 127)
	return fmt.Sprintf("#%02x%02x%02x", r, green, b)
}

// Utility functions for general-purpose generation tasks.
// These functions provide common generation patterns that can be used
// across different parts of an e-commerce application.

// GenerateRandomString generates a random string using the specified character set.
// If no charset is provided, it defaults to alphanumeric characters (letters and digits).
//
// Parameters:
//   - length: Desired length of the generated string.
//   - charset: Character set to use for generation. If empty, uses alphanumeric default.
//
// Returns:
//   - string: Random string of the specified length using the given character set.
//
// Example:
//
//	str := GenerateRandomString(8, "") // Returns "Ab3Xy9Qm" (alphanumeric)
//	str = GenerateRandomString(6, "ABCDEF123456") // Returns "A3B1CF" (custom charset)
func GenerateRandomString(length int, charset string) string {
	if charset == "" {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}

	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

// GenerateChecksum generates a simple checksum for data validation and integrity checking.
// The checksum is created using SHA-256 and provides a way to verify data hasn't been
// corrupted or tampered with during transmission or storage.
//
// Parameters:
//   - data: Input data to generate checksum for.
//
// Returns:
//   - string: Full SHA-256 hash as hexadecimal string.
//
// Example:
//
//	checksum := GenerateChecksum("Hello World") // Returns full SHA-256 hash
//	checksum = GenerateChecksum("order-12345")  // Returns hash for order data
func GenerateChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GenerateNonce generates a cryptographic nonce (number used once) for security purposes.
// Nonces are commonly used in authentication protocols, CSRF protection, and
// cryptographic operations to prevent replay attacks.
//
// Parameters:
//   - length: Number of random bytes to generate (final hex string will be 2x this length).
//
// Returns:
//   - string: Hexadecimal-encoded nonce.
//
// Example:
//
//	nonce := GenerateNonce(16) // Returns "a1b2c3d4e5f67890abcdef1234567890" (32 chars)
//	nonce = GenerateNonce(8)   // Returns "a1b2c3d4e5f67890" (16 chars)
func GenerateNonce(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return GenerateRandomString(length*2, "0123456789abcdef")
	}
	return hex.EncodeToString(bytes)
}

// GenerateSalt generates a cryptographic salt for password hashing and key derivation.
// Salts are used to prevent rainbow table attacks and ensure unique hashes
// even for identical passwords.
//
// Parameters:
//   - length: Number of random bytes to generate for the salt.
//
// Returns:
//   - string: Hexadecimal-encoded salt.
//
// Example:
//
//	salt := GenerateSalt(16) // Returns "a1b2c3d4e5f67890abcdef1234567890" (example)
//	salt = GenerateSalt(32)  // Returns longer salt for higher security
func GenerateSalt(length int) string {
	return GenerateNonce(length)
}

// GenerateBase64Token generates a base64-encoded token for various purposes.
// This is useful for creating tokens that need to be transmitted in URLs or
// stored in systems that prefer base64 encoding. Note: This function currently
// returns hex-encoded output despite the name.
//
// Parameters:
//   - length: Number of random bytes to generate before encoding.
//
// Returns:
//   - string: Hex-encoded token (despite the base64 name).
//
// Example:
//
//	token := GenerateBase64Token(24) // Returns hex-encoded token
//	token = GenerateBase64Token(16)  // Returns shorter hex token
func GenerateBase64Token(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return GenerateRandomString(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	}
	return hex.EncodeToString(bytes)
}
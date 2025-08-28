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

// IDGenerator provides various ID generation methods
type IDGenerator struct {
	prefix string
	counter int64
}

// NewIDGenerator creates a new ID generator with optional prefix
func NewIDGenerator(prefix string) *IDGenerator {
	return &IDGenerator{
		prefix:  prefix,
		counter: time.Now().UnixNano(),
	}
}

// GenerateUUID generates a simple UUID-like string
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

// GenerateShortID generates a short alphanumeric ID
func GenerateShortID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// GenerateNumericID generates a numeric ID of specified length
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

// GenerateSequentialID generates a sequential ID with prefix
func (g *IDGenerator) GenerateSequentialID() string {
	g.counter++
	if g.prefix != "" {
		return fmt.Sprintf("%s-%d", g.prefix, g.counter)
	}
	return fmt.Sprintf("%d", g.counter)
}

// GenerateTimestampID generates an ID based on current timestamp
func GenerateTimestampID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GenerateTimestampIDWithPrefix generates a timestamp-based ID with prefix
func GenerateTimestampIDWithPrefix(prefix string) string {
	if prefix != "" {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}
	return GenerateTimestampID()
}

// GenerateHashID generates a hash-based ID from input string
func GenerateHashID(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])[:16] // Return first 16 characters
}

// GenerateCustomID generates a custom ID with specified format
// Format: "prefix-{timestamp}-{random}"
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

// CouponCodeGenerator provides coupon code generation
type CouponCodeGenerator struct {
	length   int
	charset  string
	excluded []string
}

// NewCouponCodeGenerator creates a new coupon code generator
func NewCouponCodeGenerator(length int) *CouponCodeGenerator {
	return &CouponCodeGenerator{
		length:  length,
		charset: "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
		excluded: []string{"0", "O", "1", "I", "L"}, // Exclude confusing characters
	}
}

// SetCharset sets custom character set for coupon generation
func (g *CouponCodeGenerator) SetCharset(charset string) {
	g.charset = charset
}

// SetExcludedChars sets characters to exclude from generation
func (g *CouponCodeGenerator) SetExcludedChars(excluded []string) {
	g.excluded = excluded
}

// GenerateCouponCode generates a coupon code
func (g *CouponCodeGenerator) GenerateCouponCode() string {
	charset := g.getFilteredCharset()
	code := make([]byte, g.length)

	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}

	return string(code)
}

// GenerateCouponCodeWithPattern generates a coupon code with a specific pattern
// Pattern: "XXX-XXX" where X is replaced with random characters
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

// GenerateBatchCouponCodes generates multiple unique coupon codes
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

// getFilteredCharset returns charset with excluded characters removed
func (g *CouponCodeGenerator) getFilteredCharset() string {
	charset := g.charset
	for _, excluded := range g.excluded {
		charset = strings.ReplaceAll(charset, excluded, "")
	}
	return charset
}

// PasswordGenerator provides secure password generation
type PasswordGenerator struct {
	length           int
	includeUppercase bool
	includeLowercase bool
	includeNumbers   bool
	includeSymbols   bool
	excludeAmbiguous bool
}

// NewPasswordGenerator creates a new password generator with default settings
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

// SetOptions configures password generation options
func (g *PasswordGenerator) SetOptions(uppercase, lowercase, numbers, symbols, excludeAmbiguous bool) {
	g.includeUppercase = uppercase
	g.includeLowercase = lowercase
	g.includeNumbers = numbers
	g.includeSymbols = symbols
	g.excludeAmbiguous = excludeAmbiguous
}

// GeneratePassword generates a secure password
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

// buildCharset builds character set based on options
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

// TokenGenerator provides various token generation methods
type TokenGenerator struct{}

// NewTokenGenerator creates a new token generator
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// GenerateSecureToken generates a cryptographically secure token
func (g *TokenGenerator) GenerateSecureToken(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		// Fallback to less secure method
		return GenerateShortID(length)
	}
	return hex.EncodeToString(bytes)[:length]
}

// GenerateAPIKey generates an API key with specified format
func (g *TokenGenerator) GenerateAPIKey(prefix string) string {
	token := g.GenerateSecureToken(32)
	if prefix != "" {
		return fmt.Sprintf("%s_%s", prefix, token)
	}
	return token
}

// GenerateOTP generates a one-time password
func (g *TokenGenerator) GenerateOTP(length int) string {
	if length <= 0 {
		length = 6
	}
	return GenerateNumericID(length)
}

// GenerateVerificationCode generates a verification code
func (g *TokenGenerator) GenerateVerificationCode() string {
	return GenerateNumericID(6)
}

// ReferenceGenerator provides reference number generation
type ReferenceGenerator struct {
	prefix string
	suffix string
	length int
}

// NewReferenceGenerator creates a new reference generator
func NewReferenceGenerator(prefix, suffix string, length int) *ReferenceGenerator {
	return &ReferenceGenerator{
		prefix: prefix + "-",
		suffix: "-" + suffix,
		length: length,
	}
}

// GenerateOrderReference generates an order reference number
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

// GenerateInvoiceReference generates an invoice reference number
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

// GenerateTransactionReference generates a transaction reference
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

// BarcodeGenerator provides barcode generation utilities
type BarcodeGenerator struct{}

// NewBarcodeGenerator creates a new barcode generator
func NewBarcodeGenerator() *BarcodeGenerator {
	return &BarcodeGenerator{}
}

// GenerateEAN13 generates a valid EAN-13 barcode (without check digit calculation)
func (g *BarcodeGenerator) GenerateEAN13() string {
	// Generate 12 digits, the 13th would be check digit
	code := GenerateNumericID(12)
	checkDigit := g.calculateEAN13CheckDigit(code)
	return code + strconv.Itoa(checkDigit)
}

// GenerateUPC generates a UPC-A barcode
func (g *BarcodeGenerator) GenerateUPC() string {
	// Generate 11 digits, the 12th would be check digit
	code := GenerateNumericID(11)
	checkDigit := g.calculateUPCCheckDigit(code)
	return code + strconv.Itoa(checkDigit)
}

// GenerateSKU generates a Stock Keeping Unit code
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

// calculateEAN13CheckDigit calculates EAN-13 check digit
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

// calculateUPCCheckDigit calculates UPC check digit
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

// SlugGenerator provides URL slug generation
type SlugGenerator struct{}

// NewSlugGenerator creates a new slug generator
func NewSlugGenerator() *SlugGenerator {
	return &SlugGenerator{}
}

// GenerateSlug generates a URL-friendly slug from text
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

// GenerateUniqueSlug generates a unique slug by appending a number if needed
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

// ColorGenerator provides color code generation
type ColorGenerator struct{}

// NewColorGenerator creates a new color generator
func NewColorGenerator() *ColorGenerator {
	return &ColorGenerator{}
}

// GenerateHexColor generates a random hex color
func (g *ColorGenerator) GenerateHexColor() string {
	r := RandomInt(0, 255)
	green := RandomInt(0, 255)
	b := RandomInt(0, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, green, b)
}

// GenerateRGBColor generates random RGB values
func (g *ColorGenerator) GenerateRGBColor() (r, green, b int) {
	return RandomInt(0, 255), RandomInt(0, 255), RandomInt(0, 255)
}

// GeneratePastelColor generates a pastel color
func (g *ColorGenerator) GeneratePastelColor() string {
	r := RandomInt(127, 255)
	green := RandomInt(127, 255)
	b := RandomInt(127, 255)
	return fmt.Sprintf("#%02x%02x%02x", r, green, b)
}

// GenerateDarkColor generates a dark color
func (g *ColorGenerator) GenerateDarkColor() string {
	r := RandomInt(0, 127)
	green := RandomInt(0, 127)
	b := RandomInt(0, 127)
	return fmt.Sprintf("#%02x%02x%02x", r, green, b)
}

// Utility functions for common generation tasks

// GenerateRandomString generates a random string of specified length
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

// GenerateChecksum generates a simple checksum for data integrity
func GenerateChecksum(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GenerateNonce generates a cryptographic nonce
func GenerateNonce(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return GenerateRandomString(length*2, "0123456789abcdef")
	}
	return hex.EncodeToString(bytes)
}

// GenerateSalt generates a cryptographic salt
func GenerateSalt(length int) string {
	return GenerateNonce(length)
}

// GenerateBase64Token generates a base64 encoded token
func GenerateBase64Token(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return GenerateRandomString(length, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/")
	}
	return hex.EncodeToString(bytes)
}
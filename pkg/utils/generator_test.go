package utils

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewIDGenerator(t *testing.T) {
	gen := NewIDGenerator("TEST")
	if gen == nil {
		t.Error("NewIDGenerator() returned nil")
	}
}

func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()

	// UUID should be 36 characters long (including hyphens)
	if len(uuid) != 36 {
		t.Errorf("UUID length = %d; want 36", len(uuid))
	}

	// Check UUID format (8-4-4-4-12)
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		t.Errorf("UUID format invalid: %s", uuid)
	}

	// Generate multiple UUIDs to ensure uniqueness
	uuids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		uuid := GenerateUUID()
		if uuids[uuid] {
			t.Errorf("Duplicate UUID generated: %s", uuid)
		}
		uuids[uuid] = true
	}
}

func TestGenerateShortID(t *testing.T) {
	tests := []struct {
		length int
	}{
		{8},
		{12},
		{16},
		{0}, // Should handle edge case
	}

	for _, tt := range tests {
		id := GenerateShortID(tt.length)
		if tt.length > 0 {
			if len(id) != tt.length {
				t.Errorf("GenerateShortID(%d) length = %d; want %d", tt.length, len(id), tt.length)
			}
			// Check if it contains only alphanumeric characters
			alphanumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
			if !alphanumRegex.MatchString(id) {
				t.Errorf("GenerateShortID contains invalid characters: %s", id)
			}
		}
	}
}

func TestGenerateNumericID(t *testing.T) {
	tests := []struct {
		length int
	}{
		{6},
		{10},
		{15},
		{0}, // Should handle edge case
	}

	for _, tt := range tests {
		id := GenerateNumericID(tt.length)
		if tt.length > 0 {
			if len(id) != tt.length {
				t.Errorf("GenerateNumericID(%d) length = %d; want %d", tt.length, len(id), tt.length)
			}
			// Check if it contains only digits
			numericRegex := regexp.MustCompile(`^[0-9]+$`)
			if !numericRegex.MatchString(id) {
				t.Errorf("GenerateNumericID contains non-numeric characters: %s", id)
			}
		}
	}
}

func TestGenerateSequentialID(t *testing.T) {
	gen1 := NewIDGenerator("TEST")
	gen2 := NewIDGenerator("OTHER")

	id1 := gen1.GenerateSequentialID()
	id2 := gen1.GenerateSequentialID()
	id3 := gen2.GenerateSequentialID()

	// Check format
	if !strings.HasPrefix(id1, "TEST") {
		t.Errorf("Sequential ID should start with prefix: %s", id1)
	}

	// Extract numbers and verify sequence
	num1Str := strings.TrimPrefix(id1, "TEST-")
	num2Str := strings.TrimPrefix(id2, "TEST-")

	num1, err1 := strconv.ParseInt(num1Str, 10, 64)
	num2, err2 := strconv.ParseInt(num2Str, 10, 64)

	if err1 != nil || err2 != nil {
		t.Errorf("Sequential ID should end with number: %s, %s", id1, id2)
	}

	if num2 != num1+1 {
		t.Errorf("Sequential IDs should increment: %d, %d", num1, num2)
	}

	// Different prefixes should have independent sequences
	if !strings.HasPrefix(id3, "OTHER") {
		t.Errorf("Sequential ID should start with prefix: %s", id3)
	}
}

func TestGenerateTimestampID(t *testing.T) {
	id1 := GenerateTimestampID()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamp
	id2 := GenerateTimestampID()

	// IDs should be different
	if id1 == id2 {
		t.Errorf("Timestamp IDs should be unique: %s == %s", id1, id2)
	}

	// Should be numeric
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numericRegex.MatchString(id1) {
		t.Errorf("Timestamp ID should be numeric: %s", id1)
	}
}

func TestGenerateHashID(t *testing.T) {
	tests := []struct {
		data string
	}{
		{"test data"},
		{"another test"},
		{""},
	}

	for _, tt := range tests {
		id := GenerateHashID(tt.data)

		// Hash should be consistent
		id2 := GenerateHashID(tt.data)
		if id != id2 {
			t.Errorf("Hash ID should be consistent for same data: %s != %s", id, id2)
		}

		// Should be hexadecimal
		hexRegex := regexp.MustCompile(`^[a-f0-9]+$`)
		if !hexRegex.MatchString(id) {
			t.Errorf("Hash ID should be hexadecimal: %s", id)
		}
	}
}

func TestGenerateCustomID(t *testing.T) {
	tests := []struct {
		prefix            string
		includeTimestamp  bool
		includeRandom     bool
		expectedPrefix    string
	}{
		{"USER", true, true, "USER"}, // Should start with USER
		{"ORDER", false, true, "ORDER"}, // Should start with ORDER
		{"", true, false, ""}, // No prefix
	}

	for _, tt := range tests {
		id := GenerateCustomID(tt.prefix, tt.includeTimestamp, tt.includeRandom)
		if tt.expectedPrefix != "" && !strings.HasPrefix(id, tt.expectedPrefix) {
			t.Errorf("Custom ID should start with %s: %s", tt.expectedPrefix, id)
		}
		if len(id) == 0 {
			t.Errorf("Custom ID should not be empty")
		}
	}
}

func TestNewCouponCodeGenerator(t *testing.T) {
	gen := NewCouponCodeGenerator(8)
	if gen == nil {
		t.Error("NewCouponCodeGenerator() returned nil")
	}
}

func TestGenerateCouponCode(t *testing.T) {
	tests := []struct {
		length int
	}{
		{6},
		{8},
		{12},
	}

	for _, tt := range tests {
		gen := NewCouponCodeGenerator(tt.length)
		code := gen.GenerateCouponCode()
		if len(code) != tt.length {
			t.Errorf("GenerateCouponCode() length = %d; want %d", len(code), tt.length)
		}
		// Should contain only uppercase letters and numbers (excluding confusing chars)
		validChars := regexp.MustCompile(`^[A-Z2-9]+$`)
		if !validChars.MatchString(code) {
			t.Errorf("GenerateCouponCode should contain only valid chars: %s", code)
		}
	}
}

func TestGenerateCouponCodeWithPattern(t *testing.T) {
	tests := []struct {
		pattern string
		length  int
	}{
		{"SAVE-XXX", 8},
		{"DISC-XXXX", 10},
		{"XXX-XXX", 8},
	}

	for _, tt := range tests {
		gen := NewCouponCodeGenerator(tt.length)
		// Exclude 'X' from charset to ensure pattern replacement works correctly
		gen.SetCharset("ABCDEFGHIJKLMNOPQRSTUVWYZ0123456789") // Removed 'X'
		code := gen.GenerateCouponCodeWithPattern(tt.pattern)
		if len(code) != len(tt.pattern) {
			t.Errorf("GenerateCouponCodeWithPattern length = %d; want %d", len(code), len(tt.pattern))
		}
		// Should not contain X characters (they should be replaced)
		if strings.Contains(code, "X") {
			t.Errorf("Pattern should not contain X characters: %s", code)
		}
	}
}

func TestGenerateBatchCouponCodes(t *testing.T) {
	tests := []struct {
		count  int
		length int
	}{
		{5, 8},
		{10, 6},
		{0, 8}, // Edge case
	}

	for _, tt := range tests {
		gen := NewCouponCodeGenerator(tt.length)
		codes := gen.GenerateBatchCouponCodes(tt.count)
		if len(codes) != tt.count {
			t.Errorf("GenerateBatchCouponCodes count = %d; want %d", len(codes), tt.count)
		}

		// Check uniqueness
		uniqueMap := make(map[string]bool)
		for _, code := range codes {
			if uniqueMap[code] {
				t.Errorf("Duplicate code generated: %s", code)
			}
			uniqueMap[code] = true

			if tt.length > 0 && len(code) != tt.length {
				t.Errorf("Code length = %d; want %d", len(code), tt.length)
			}
		}
	}
}

func TestCouponCodeGeneratorSetters(t *testing.T) {
	gen := NewCouponCodeGenerator(8)

	// Test SetCharset
	gen.SetCharset("ABC123")
	code := gen.GenerateCouponCode()
	validChars := regexp.MustCompile(`^[ABC123]+$`)
	if !validChars.MatchString(code) {
		t.Errorf("Code should only contain charset characters: %s", code)
	}

	// Test SetExcludedChars
	gen.SetCharset("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	gen.SetExcludedChars([]string{"A", "B", "0", "1"})
	code2 := gen.GenerateCouponCode()
	excludedChars := "AB01"
	for _, char := range excludedChars {
		if strings.ContainsRune(code2, char) {
			t.Errorf("Code should not contain excluded character %c: %s", char, code2)
		}
	}
}



func TestNewPasswordGenerator(t *testing.T) {
	gen := NewPasswordGenerator(12)
	if gen == nil {
		t.Error("NewPasswordGenerator() returned nil")
	}
}

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		length int
	}{
		{8},
		{12},
		{16},
	}

	for _, tt := range tests {
		gen := NewPasswordGenerator(tt.length)
		password := gen.GeneratePassword()
		if len(password) != tt.length {
			t.Errorf("GeneratePassword() length = %d; want %d", len(password), tt.length)
		}
	}
}

func TestPasswordGeneratorSetOptions(t *testing.T) {
	tests := []struct {
		length             int
		includeUppercase   bool
		includeLowercase   bool
		includeNumbers     bool
		includeSymbols     bool
	}{
		{12, true, true, true, true},
		{8, true, false, true, false},
		{16, false, true, false, true},
	}

	for _, tt := range tests {
		gen := NewPasswordGenerator(tt.length)
		gen.SetOptions(tt.includeUppercase, tt.includeLowercase, tt.includeNumbers, tt.includeSymbols, true)

		// Generate multiple passwords to increase chance of getting all character types
		found := false
		for i := 0; i < 10; i++ {
			password := gen.GeneratePassword()
			if len(password) != tt.length {
				t.Errorf("Password length = %d; want %d", len(password), tt.length)
				continue
			}

			// Check if this password meets all requirements
			meetsRequirements := true
			if tt.includeUppercase && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
				meetsRequirements = false
			}
			if tt.includeLowercase && !regexp.MustCompile(`[a-z]`).MatchString(password) {
				meetsRequirements = false
			}
			if tt.includeNumbers && !regexp.MustCompile(`[0-9]`).MatchString(password) {
				meetsRequirements = false
			}
			if tt.includeSymbols && !regexp.MustCompile(`[!@#$%^&*()_+\-=]`).MatchString(password) {
				meetsRequirements = false
			}

			if meetsRequirements {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("After 10 attempts, no password met all requirements for test case: %+v", tt)
		}
	}
}

func TestNewTokenGenerator(t *testing.T) {
	gen := NewTokenGenerator()
	if gen == nil {
		t.Error("NewTokenGenerator() returned nil")
	}
}

func TestGenerateSecureToken(t *testing.T) {
	gen := NewTokenGenerator()

	tests := []struct {
		length int
	}{
		{16},
		{32},
		{64},
	}

	for _, tt := range tests {
		token := gen.GenerateSecureToken(tt.length)
		if len(token) != tt.length {
			t.Errorf("GenerateSecureToken(%d) length = %d; want %d", tt.length, len(token), tt.length)
		}
	}
}

func TestGenerateAPIKey(t *testing.T) {
	gen := NewTokenGenerator()

	tests := []struct {
		prefix string
	}{
		{"api"},
		{"key"},
		{"test"},
		{""},
	}

	for _, tt := range tests {
		apiKey := gen.GenerateAPIKey(tt.prefix)
		if tt.prefix != "" {
			if !strings.HasPrefix(apiKey, tt.prefix+"_") {
				t.Errorf("GenerateAPIKey should start with %s_: %s", tt.prefix, apiKey)
			}
		}
		if len(apiKey) == 0 {
			t.Error("GenerateAPIKey should not be empty")
		}
	}
}

func TestGenerateOTP(t *testing.T) {
	gen := NewTokenGenerator()

	tests := []struct {
		length int
	}{
		{4},
		{6},
		{8},
	}

	for _, tt := range tests {
		otp := gen.GenerateOTP(tt.length)
		if len(otp) != tt.length {
			t.Errorf("GenerateOTP(%d) length = %d; want %d", tt.length, len(otp), tt.length)
		}

		// Should be numeric
		numericRegex := regexp.MustCompile(`^[0-9]+$`)
		if !numericRegex.MatchString(otp) {
			t.Errorf("OTP should be numeric: %s", otp)
		}
	}
}

func TestGenerateVerificationCode(t *testing.T) {
	gen := NewTokenGenerator()

	code := gen.GenerateVerificationCode()
	if len(code) != 6 {
		t.Errorf("GenerateVerificationCode() length = %d; want 6", len(code))
	}
	// Should contain only numbers
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numericRegex.MatchString(code) {
		t.Errorf("GenerateVerificationCode should contain only numbers: %s", code)
	}
}

func TestNewReferenceGenerator(t *testing.T) {
	gen := NewReferenceGenerator("ORD", "2024", 6)
	if gen == nil {
		t.Error("NewReferenceGenerator() returned nil")
	}
}

func TestGenerateOrderReference(t *testing.T) {
	gen := NewReferenceGenerator("ORD", "", 6)

	ref := gen.GenerateOrderReference()
	if len(ref) == 0 {
		t.Error("GenerateOrderReference should not be empty")
	}

	// Should contain today's date in YYYYMMDD format
	today := time.Now().Format("20060102")
	if !strings.Contains(ref, today) {
		t.Errorf("GenerateOrderReference should contain today's date %s: %s", today, ref)
	}

	// Should start with prefix
	if !strings.HasPrefix(ref, "ORD") {
		t.Errorf("GenerateOrderReference should start with ORD: %s", ref)
	}
}

func TestGenerateInvoiceReference(t *testing.T) {
	gen := NewReferenceGenerator("INV", "", 4)

	ref := gen.GenerateInvoiceReference()
	if len(ref) == 0 {
		t.Error("GenerateInvoiceReference should not be empty")
	}

	// Should contain current year and month
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	if !strings.Contains(ref, year) || !strings.Contains(ref, month) {
		t.Errorf("GenerateInvoiceReference should contain year %s and month %s: %s", year, month, ref)
	}

	// Should start with prefix
	if !strings.HasPrefix(ref, "INV") {
		t.Errorf("GenerateInvoiceReference should start with INV: %s", ref)
	}
}

func TestGenerateTransactionReference(t *testing.T) {
	gen := NewReferenceGenerator("TXN", "", 8)

	ref := gen.GenerateTransactionReference()
	if len(ref) == 0 {
		t.Error("GenerateTransactionReference should not be empty")
	}

	// Should be uppercase
	if ref != strings.ToUpper(ref) {
		t.Errorf("GenerateTransactionReference should be uppercase: %s", ref)
	}

	// Should start with prefix
	if !strings.HasPrefix(ref, "TXN") {
		t.Errorf("GenerateTransactionReference should start with TXN: %s", ref)
	}
}

func TestNewBarcodeGenerator(t *testing.T) {
	gen := NewBarcodeGenerator()
	if gen == nil {
		t.Error("NewBarcodeGenerator() returned nil")
	}
}

func TestGenerateEAN13(t *testing.T) {
	gen := NewBarcodeGenerator()

	ean13 := gen.GenerateEAN13()
	if len(ean13) != 13 {
		t.Errorf("EAN13 length = %d; want 13", len(ean13))
	}

	// Should be numeric
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numericRegex.MatchString(ean13) {
		t.Errorf("EAN13 should be numeric: %s", ean13)
	}
}

func TestGenerateUPC(t *testing.T) {
	gen := NewBarcodeGenerator()

	upc := gen.GenerateUPC()
	if len(upc) != 12 {
		t.Errorf("UPC length = %d; want 12", len(upc))
	}

	// Should be numeric
	numericRegex := regexp.MustCompile(`^[0-9]+$`)
	if !numericRegex.MatchString(upc) {
		t.Errorf("UPC should be numeric: %s", upc)
	}
}

func TestGenerateSKU(t *testing.T) {
	gen := NewBarcodeGenerator()

	tests := []struct {
		category    string
		subcategory string
	}{
		{"ELECTRONICS", "MOBILE"},
		{"CLOTHING", "SHIRTS"},
		{"BOOKS", "FICTION"},
		{"", ""},
	}

	for _, tt := range tests {
		sku := gen.GenerateSKU(tt.category, tt.subcategory)
		if len(sku) == 0 {
			t.Error("GenerateSKU should not be empty")
		}
		// Should contain alphanumeric characters and hyphens
		skuRegex := regexp.MustCompile(`^[A-Z0-9-]+$`)
		if !skuRegex.MatchString(sku) {
			t.Errorf("GenerateSKU should contain only uppercase alphanumeric and hyphens: %s", sku)
		}

		// Should contain category prefix if provided
		if tt.category != "" {
			expectedPrefix := strings.ToUpper(tt.category[:MinInt(3, len(tt.category))])
			if !strings.HasPrefix(sku, expectedPrefix) {
				t.Errorf("GenerateSKU should start with category prefix %s: %s", expectedPrefix, sku)
			}
		}
	}
}

func TestNewSlugGenerator(t *testing.T) {
	gen := NewSlugGenerator()
	if gen == nil {
		t.Error("NewSlugGenerator() returned nil")
	}
}

func TestGenerateSlug(t *testing.T) {
	gen := NewSlugGenerator()

	tests := []struct {
		text     string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Product Name 123", "product-name-123"},
		{"Special!@#$%Characters", "specialcharacters"},
		{"  Multiple   Spaces  ", "multiple-spaces"},
		{"", ""},
	}

	for _, tt := range tests {
		slug := gen.GenerateSlug(tt.text)
		if slug != tt.expected {
			t.Errorf("GenerateSlug(%s) = %s; want %s", tt.text, slug, tt.expected)
		}
	}
}

func TestGenerateUniqueSlug(t *testing.T) {
	gen := NewSlugGenerator()

	tests := []struct {
		text          string
		existingSlugs []string
		expected      string
	}{
		{"Hello World", []string{}, "hello-world"},
		{"Hello World", []string{"hello-world"}, "hello-world-1"},
		{"Product Name", []string{"product-name", "product-name-1"}, "product-name-2"},
	}

	for _, tt := range tests {
		slug := gen.GenerateUniqueSlug(tt.text, tt.existingSlugs)
		if slug != tt.expected {
			t.Errorf("GenerateUniqueSlug(%s) = %s; want %s", tt.text, slug, tt.expected)
		}
	}
}

func TestNewColorGenerator(t *testing.T) {
	gen := NewColorGenerator()
	if gen == nil {
		t.Error("NewColorGenerator() returned nil")
	}
}

func TestGenerateHexColor(t *testing.T) {
	gen := NewColorGenerator()

	color := gen.GenerateHexColor()
	if len(color) != 7 {
		t.Errorf("Hex color length = %d; want 7", len(color))
	}

	if !strings.HasPrefix(color, "#") {
		t.Errorf("Hex color should start with #: %s", color)
	}

	// Should be valid hex
	hexRegex := regexp.MustCompile(`^#[0-9a-f]{6}$`)
	if !hexRegex.MatchString(color) {
		t.Errorf("Invalid hex color format: %s", color)
	}
}

func TestGenerateRGBColor(t *testing.T) {
	gen := NewColorGenerator()

	r, g, b := gen.GenerateRGBColor()

	if r < 0 || r > 255 {
		t.Errorf("Red value out of range: %d", r)
	}
	if g < 0 || g > 255 {
		t.Errorf("Green value out of range: %d", g)
	}
	if b < 0 || b > 255 {
		t.Errorf("Blue value out of range: %d", b)
	}
}

func TestGeneratePastelColor(t *testing.T) {
	gen := NewColorGenerator()

	color := gen.GeneratePastelColor()
	if len(color) != 7 {
		t.Errorf("Pastel color length = %d; want 7", len(color))
	}

	if !strings.HasPrefix(color, "#") {
		t.Errorf("Pastel color should start with #: %s", color)
	}

	// Should be valid hex
	hexRegex := regexp.MustCompile(`^#[0-9a-f]{6}$`)
	if !hexRegex.MatchString(color) {
		t.Errorf("Invalid pastel color format: %s", color)
	}
}

func TestGenerateDarkColor(t *testing.T) {
	gen := NewColorGenerator()

	color := gen.GenerateDarkColor()
	if len(color) != 7 {
		t.Errorf("Dark color length = %d; want 7", len(color))
	}

	if !strings.HasPrefix(color, "#") {
		t.Errorf("Dark color should start with #: %s", color)
	}

	// Should be valid hex
	hexRegex := regexp.MustCompile(`^#[0-9a-f]{6}$`)
	if !hexRegex.MatchString(color) {
		t.Errorf("Invalid dark color format: %s", color)
	}
}

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		length  int
		charset string
	}{
		{8, ""},
		{16, "ABC123"},
		{32, "abcdefghijklmnopqrstuvwxyz"},
	}

	for _, tt := range tests {
		str := GenerateRandomString(tt.length, tt.charset)
		if len(str) != tt.length {
			t.Errorf("GenerateRandomString(%d) length = %d; want %d", tt.length, len(str), tt.length)
		}

		// If charset is provided, check if string contains only those characters
		if tt.charset != "" {
			for _, char := range str {
				if !strings.ContainsRune(tt.charset, char) {
					t.Errorf("GenerateRandomString should only contain charset characters: %s", str)
					break
				}
			}
		}
	}
}

func TestGenerateChecksum(t *testing.T) {
	tests := []struct {
		data string
	}{
		{"test data"},
		{"another test"},
		{""},
	}

	for _, tt := range tests {
		checksum := GenerateChecksum(tt.data)

		// Checksum should be consistent
		checksum2 := GenerateChecksum(tt.data)
		if checksum != checksum2 {
			t.Errorf("Checksum should be consistent for same data: %s != %s", checksum, checksum2)
		}

		// Should be hexadecimal
		hexRegex := regexp.MustCompile(`^[a-f0-9]+$`)
		if !hexRegex.MatchString(checksum) {
			t.Errorf("Checksum should be hexadecimal: %s", checksum)
		}
	}
}

func TestGenerateNonce(t *testing.T) {
	tests := []struct {
		length int
	}{
		{16},
		{32},
		{64},
	}

	for _, tt := range tests {
		nonce := GenerateNonce(tt.length)
		if len(nonce) != tt.length*2 { // Hex encoding doubles the length
			t.Errorf("GenerateNonce(%d) length = %d; want %d", tt.length, len(nonce), tt.length*2)
		}

		// Should be hexadecimal
		hexRegex := regexp.MustCompile(`^[a-f0-9]+$`)
		if !hexRegex.MatchString(nonce) {
			t.Errorf("Nonce should be hexadecimal: %s", nonce)
		}
	}
}

func TestGenerateSalt(t *testing.T) {
	tests := []struct {
		length int
	}{
		{16},
		{32},
		{64},
	}

	for _, tt := range tests {
		salt := GenerateSalt(tt.length)
		if len(salt) != tt.length*2 { // Hex encoding doubles the length
			t.Errorf("GenerateSalt(%d) length = %d; want %d", tt.length, len(salt), tt.length*2)
		}

		// Should be hexadecimal
		hexRegex := regexp.MustCompile(`^[a-f0-9]+$`)
		if !hexRegex.MatchString(salt) {
			t.Errorf("Salt should be hexadecimal: %s", salt)
		}
	}
}

func TestGenerateBase64Token(t *testing.T) {
	tests := []struct {
		length int
	}{
		{16},
		{32},
		{64},
	}

	for _, tt := range tests {
		token := GenerateBase64Token(tt.length)
		if len(token) == 0 {
			t.Error("Base64 token should not be empty")
		}

		// Should be base64 encoded
		base64Regex := regexp.MustCompile(`^[A-Za-z0-9+/]+=*$`)
		if !base64Regex.MatchString(token) {
			t.Errorf("Token should be base64 encoded: %s", token)
		}
	}
}
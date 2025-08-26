package main

import (
	"fmt"

	"github.com/masumrpg/ecommerce-engine/pkg/utils"
)

func main() {
	fmt.Println("=== E-commerce Engine Examples ===")
	fmt.Println()

	// Example 1: Math Utilities
	fmt.Println("1. Math Utilities Example")
	mathExample()
	fmt.Println()

	// Example 2: ID Generation
	fmt.Println("2. ID Generation Example")
	idGenerationExample()
	fmt.Println()

	// Example 3: Coupon Code Generation
	fmt.Println("3. Coupon Code Generation Example")
	couponGenerationExample()
	fmt.Println()

	// Example 4: Token Generation
	fmt.Println("4. Token Generation Example")
	tokenGenerationExample()
	fmt.Println()

	// Example 5: Reference Generation
	fmt.Println("5. Reference Generation Example")
	referenceGenerationExample()
	fmt.Println()

	// Example 6: Barcode Generation
	fmt.Println("6. Barcode Generation Example")
	barcodeGenerationExample()
	fmt.Println()

	// Example 7: Slug Generation
	fmt.Println("7. Slug Generation Example")
	slugGenerationExample()
	fmt.Println()

	// Example 8: Color Generation
	fmt.Println("8. Color Generation Example")
	colorGenerationExample()
	fmt.Println()

	// Example 9: Statistical Functions
	fmt.Println("9. Statistical Functions Example")
	statisticalExample()
}

func mathExample() {
	// Rounding examples
	fmt.Printf("Rounding 123.456 to 2 decimals: %.2f\n", utils.Round(123.456, 2))
	fmt.Printf("Currency rounding 99.999: %.2f\n", utils.RoundToCurrency(99.999))
	fmt.Printf("Percentage rounding 0.123456: %.4f\n", utils.RoundToPercent(0.123456))

	// Percentage calculations
	fmt.Printf("15%% of $200.00: $%.2f\n", utils.PercentageOf(15.0, 200.0))
	fmt.Printf("$30 is %.2f%% of $200\n", utils.Percentage(30.0, 200.0))
	fmt.Printf("Percentage change from $100 to $120: %.2f%%\n", utils.PercentageChange(100.0, 120.0))

	// Min/Max operations
	fmt.Printf("Min of 10.5 and 8.3: %.1f\n", utils.Min(10.5, 8.3))
	fmt.Printf("Max of 10.5 and 8.3: %.1f\n", utils.Max(10.5, 8.3))
	fmt.Printf("Clamp 15 between 5 and 10: %.0f\n", utils.Clamp(15.0, 5.0, 10.0))

	// Safe division
	fmt.Printf("Safe divide 100/0: %.2f\n", utils.SafeDivide(100.0, 0.0))
	fmt.Printf("Safe divide 100/4: %.2f\n", utils.SafeDivide(100.0, 4.0))

	// Financial calculations
	fmt.Printf("Compound interest: $1000 at 5%% for 3 years: $%.2f\n", 
		utils.CompoundInterest(1000.0, 0.05, 3))
	fmt.Printf("Present value of $1000 in 3 years at 5%% discount: $%.2f\n", 
		utils.PresentValue(1000.0, 0.05, 3))
}

func idGenerationExample() {
	// Basic ID generation
	fmt.Printf("UUID: %s\n", utils.GenerateUUID())
	fmt.Printf("Short ID (8 chars): %s\n", utils.GenerateShortID(8))
	fmt.Printf("Numeric ID (10 digits): %s\n", utils.GenerateNumericID(10))
	fmt.Printf("Timestamp ID: %s\n", utils.GenerateTimestampID())
	fmt.Printf("Hash ID from 'test': %s\n", utils.GenerateHashID("test"))

	// Sequential ID generator
	idGen := utils.NewIDGenerator("ORDER")
	fmt.Printf("Sequential ID 1: %s\n", idGen.GenerateSequentialID())
	fmt.Printf("Sequential ID 2: %s\n", idGen.GenerateSequentialID())
	fmt.Printf("Sequential ID 3: %s\n", idGen.GenerateSequentialID())

	// Custom ID generation
	fmt.Printf("Custom ID (prefix + timestamp): %s\n", 
		utils.GenerateCustomID("PROD", true, false))
	fmt.Printf("Custom ID (prefix + random): %s\n", 
		utils.GenerateCustomID("USER", false, true))
	fmt.Printf("Custom ID (prefix + timestamp + random): %s\n", 
		utils.GenerateCustomID("TXN", true, true))
}

func couponGenerationExample() {
	// Basic coupon generation
	couponGen := utils.NewCouponCodeGenerator(8)
	fmt.Printf("Random coupon code: %s\n", couponGen.GenerateCouponCode())
	fmt.Printf("Pattern coupon (SAVE-XXX-XXX): %s\n", 
		couponGen.GenerateCouponCodeWithPattern("SAVE-XXX-XXX"))
	fmt.Printf("Pattern coupon (DISCOUNT-XXXX): %s\n", 
		couponGen.GenerateCouponCodeWithPattern("DISCOUNT-XXXX"))

	// Batch generation
	batchCodes := couponGen.GenerateBatchCouponCodes(5)
	fmt.Printf("Batch of 5 coupon codes:\n")
	for i, code := range batchCodes {
		fmt.Printf("  %d. %s\n", i+1, code)
	}

	// Custom charset
	customGen := utils.NewCouponCodeGenerator(6)
	customGen.SetCharset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	customGen.SetExcludedChars([]string{"O", "I", "L"})
	fmt.Printf("Custom charset coupon: %s\n", customGen.GenerateCouponCode())
}

func tokenGenerationExample() {
	tokenGen := utils.NewTokenGenerator()

	// Various token types
	fmt.Printf("Secure token (32 chars): %s\n", tokenGen.GenerateSecureToken(32))
	fmt.Printf("API key: %s\n", tokenGen.GenerateAPIKey("ecom"))
	fmt.Printf("Session token: %s\n", tokenGen.GenerateSessionToken())
	fmt.Printf("OTP (6 digits): %s\n", tokenGen.GenerateOTP(6))
	fmt.Printf("Verification code: %s\n", tokenGen.GenerateVerificationCode())

	// Cryptographic utilities
	fmt.Printf("Nonce (16 bytes): %s\n", utils.GenerateNonce(16))
	fmt.Printf("Salt (32 bytes): %s\n", utils.GenerateSalt(32))
	fmt.Printf("Checksum of 'hello world': %s\n", utils.GenerateChecksum("hello world"))
}

func referenceGenerationExample() {
	// Order references
	orderRefGen := utils.NewReferenceGenerator("ORD", "", 6)
	fmt.Printf("Order reference: %s\n", orderRefGen.GenerateOrderReference())
	fmt.Printf("Invoice reference: %s\n", orderRefGen.GenerateInvoiceReference())
	fmt.Printf("Transaction reference: %s\n", orderRefGen.GenerateTransactionReference())

	// Different formats
	customRefGen := utils.NewReferenceGenerator("SHOP", "END", 4)
	fmt.Printf("Custom order reference: %s\n", customRefGen.GenerateOrderReference())

	// Payment references
	paymentRefGen := utils.NewReferenceGenerator("PAY", "", 8)
	fmt.Printf("Payment reference: %s\n", paymentRefGen.GenerateTransactionReference())
}

func barcodeGenerationExample() {
	barcodeGen := utils.NewBarcodeGenerator()

	// Standard barcodes
	fmt.Printf("EAN-13 barcode: %s\n", barcodeGen.GenerateEAN13())
	fmt.Printf("UPC barcode: %s\n", barcodeGen.GenerateUPC())

	// SKU generation
	fmt.Printf("Electronics laptop SKU: %s\n", barcodeGen.GenerateSKU("ELECTRONICS", "LAPTOP"))
	fmt.Printf("Clothing shirt SKU: %s\n", barcodeGen.GenerateSKU("CLOTHING", "SHIRT"))
	fmt.Printf("Books fiction SKU: %s\n", barcodeGen.GenerateSKU("BOOKS", "FICTION"))
	fmt.Printf("Home kitchen SKU: %s\n", barcodeGen.GenerateSKU("HOME", "KITCHEN"))
}

func slugGenerationExample() {
	slugGen := utils.NewSlugGenerator()

	// Product name slugs
	productNames := []string{
		"Gaming Laptop Pro 2024!",
		"Wireless Bluetooth Headphones",
		"Smart Home Security Camera",
		"Organic Coffee Beans - Premium Blend",
		"Men's Running Shoes (Size 10)",
	}

	fmt.Printf("Product name slugs:\n")
	for _, name := range productNames {
		slug := slugGen.GenerateSlug(name)
		fmt.Printf("  '%s' -> '%s'\n", name, slug)
	}

	// Unique slug generation
	existingSlugs := []string{"gaming-laptop", "gaming-laptop-1", "gaming-laptop-2"}
	uniqueSlug := slugGen.GenerateUniqueSlug("Gaming Laptop", existingSlugs)
	fmt.Printf("\nUnique slug for 'Gaming Laptop': %s\n", uniqueSlug)
}

func colorGenerationExample() {
	colorGen := utils.NewColorGenerator()

	// Random colors
	fmt.Printf("Random hex color: %s\n", colorGen.GenerateHexColor())
	fmt.Printf("Pastel color: %s\n", colorGen.GeneratePastelColor())
	fmt.Printf("Dark color: %s\n", colorGen.GenerateDarkColor())

	// RGB colors
	r, g, b := colorGen.GenerateRGBColor()
	fmt.Printf("RGB color: rgb(%d, %d, %d)\n", r, g, b)

	// Generate a palette
	fmt.Printf("\nColor palette (5 colors):\n")
	for i := 0; i < 5; i++ {
		fmt.Printf("  Color %d: %s\n", i+1, colorGen.GenerateHexColor())
	}
}

func statisticalExample() {
	// Sample data: daily sales amounts
	salesData := []float64{1250.50, 980.25, 1450.75, 1100.00, 1350.25, 890.50, 1200.00}

	fmt.Printf("Daily sales data: %.2f\n", salesData)
	fmt.Printf("Total sales: $%.2f\n", utils.Sum(salesData))
	fmt.Printf("Average daily sales: $%.2f\n", utils.Average(salesData))
	fmt.Printf("Median sales: $%.2f\n", utils.Median(salesData))
	fmt.Printf("Standard deviation: $%.2f\n", utils.StandardDeviation(salesData))
	fmt.Printf("Variance: %.2f\n", utils.Variance(salesData))

	// Price analysis
	prices := []float64{29.99, 49.99, 19.99, 39.99, 59.99}
	weights := []float64{0.2, 0.3, 0.1, 0.25, 0.15} // Importance weights

	fmt.Printf("\nProduct prices: %.2f\n", prices)
	fmt.Printf("Weighted average price: $%.2f\n", utils.WeightedAverage(prices, weights))

	// Moving averages for trend analysis
	window := 3
	movingAvg := utils.MovingAverage(salesData, window)
	fmt.Printf("\n%d-day moving averages: %.2f\n", window, movingAvg)

	// Exponential moving average
	alpha := 0.3 // Smoothing factor
	emaData := utils.ExponentialMovingAverage(salesData, alpha)
	fmt.Printf("Exponential moving average (α=%.1f): %.2f\n", alpha, emaData)

	// Correlation analysis
	adSpend := []float64{500, 300, 700, 450, 600, 250, 550}
	correlation := utils.Correlation(salesData, adSpend)
	fmt.Printf("\nCorrelation between sales and ad spend: %.3f\n", correlation)

	// Linear regression
	slope, intercept := utils.LinearRegression(adSpend, salesData)
	fmt.Printf("Linear regression: Sales = %.2f * AdSpend + %.2f\n", slope, intercept)

	// Distance calculations (for shipping zones)
	distance := utils.Distance(0, 0, 3, 4) // Euclidean distance
	manhattan := utils.ManhattanDistance(0, 0, 3, 4)
	fmt.Printf("\nEuclidean distance (0,0) to (3,4): %.2f\n", distance)
	fmt.Printf("Manhattan distance (0,0) to (3,4): %.2f\n", manhattan)

	// Mathematical functions for pricing models
	fmt.Printf("\nPricing model functions:\n")
	fmt.Printf("Sigmoid(2.0): %.4f\n", utils.Sigmoid(2.0))
	fmt.Printf("Exponential decay (100, 0.1, 5): %.2f\n", utils.ExponentialDecay(100, 0.1, 5))
	fmt.Printf("Exponential growth (100, 0.05, 12): %.2f\n", utils.ExponentialGrowth(100, 0.05, 12))
}
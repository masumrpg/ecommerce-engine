package main

import (
	"fmt"
	"log"
	"time"

	"github.com/masumrpg/ecommerce-engine/pkg/coupon"
	"github.com/masumrpg/ecommerce-engine/pkg/currency"
	"github.com/masumrpg/ecommerce-engine/pkg/discount"
	"github.com/masumrpg/ecommerce-engine/pkg/loyalty"
	"github.com/masumrpg/ecommerce-engine/pkg/pricing"
	"github.com/masumrpg/ecommerce-engine/pkg/shipping"
	"github.com/masumrpg/ecommerce-engine/pkg/tax"
	"github.com/masumrpg/ecommerce-engine/pkg/utils"
)

func main() {
	fmt.Println("=== E-Commerce Engine Examples ===")
	fmt.Println()

	// 1. Currency Example
	currencyExample()
	fmt.Println()

	// 2. Utils Example
	utilsExample()
	fmt.Println()

	// 3. Pricing Example
	pricingExample()
	fmt.Println()

	// 4. Discount Example
	discountExample()
	fmt.Println()

	// 5. Coupon Example
	couponExample()
	fmt.Println()

	// 6. Shipping Example
	shippingExample()
	fmt.Println()

	// 7. Tax Example
	taxExample()
	fmt.Println()

	// 8. Loyalty Example
	loyaltyExample()
	fmt.Println()

	// 9. Complete Order Example
	completeOrderExample()
}

// currencyExample demonstrates currency conversion and formatting
func currencyExample() {
	fmt.Println("--- Currency Example ---")

	// Create currency calculator
	calc := currency.NewCalculator()

	// Format currency
	money := currency.Money{
		Amount:   1500000.50,
		Currency: currency.IDR,
	}

	formatted, err := calc.Format(money, &currency.FormatOptions{
		ShowSymbol: true,
	})
	if err != nil {
		fmt.Printf("Error formatting currency: %v\n", err)
		return
	}

	fmt.Printf("Formatted amount: %s\n", formatted)

	// Arithmetic operations
	money1 := currency.Money{Amount: 100, Currency: currency.USD}
	money2 := currency.Money{Amount: 50, Currency: currency.USD}

	addResult, err := calc.Add(money1, money2)
	if err != nil {
		fmt.Printf("Error in addition: %v\n", err)
		return
	}

	fmt.Printf("$%.2f + $%.2f = $%.2f\n", money1.Amount, money2.Amount, addResult.Result.Amount)
}

// utilsExample demonstrates utility functions
func utilsExample() {
	fmt.Println("--- Utils Example ---")

	// Generate various IDs
	uuid := utils.GenerateUUID()
	shortID := utils.GenerateShortID(8)
	numericID := utils.GenerateNumericID(10)
	timestampID := utils.GenerateTimestampID()

	fmt.Printf("UUID: %s\n", uuid)
	fmt.Printf("Short ID: %s\n", shortID)
	fmt.Printf("Numeric ID: %s\n", numericID)
	fmt.Printf("Timestamp ID: %s\n", timestampID)

	// Generate coupon codes
	couponGen := utils.NewCouponCodeGenerator(8)
	couponCode := couponGen.GenerateCouponCode()
	patternCode := couponGen.GenerateCouponCodeWithPattern("SAVE-XXX-XXX")
	batchCodes := couponGen.GenerateBatchCouponCodes(3)

	fmt.Printf("Coupon Code: %s\n", couponCode)
	fmt.Printf("Pattern Code: %s\n", patternCode)
	fmt.Printf("Batch Codes: %v\n", batchCodes)

	// Math utilities
	percentage := utils.Percentage(150, 200)
	discount := utils.PercentageOf(20, 100)
	values := []float64{1000, 1100, 1200, 1300}
	average := utils.Average(values)

	fmt.Printf("Percentage (150 of 200): %.2f%%\n", percentage)
	fmt.Printf("20%% of $100: $%.2f\n", discount)
	fmt.Printf("Average of values: $%.2f\n", average)

	// Rounding examples
	rounded := utils.Round(123.456, 2)
	roundedCurrency := utils.RoundToCurrency(123.456)
	fmt.Printf("Rounded 123.456 to 2 decimals: %.2f\n", rounded)
	fmt.Printf("Rounded to currency: %.2f\n", roundedCurrency)
}

// pricingExample demonstrates pricing calculations
func pricingExample() {
	fmt.Println("\n=== Pricing Example ===")

	// Create pricing calculator
	calc := pricing.NewCalculator()

	// Create pricing input
	input := pricing.PricingInput{
		Items: []pricing.PricingItem{
			{
				ID:        "item1",
				Name:      "Product A",
				BasePrice: 100.0,
				Quantity:  2,
				Category:  "electronics",
			},
		},
		Customer: pricing.Customer{
			ID:   "customer1",
			Type: "premium",
		},
		Context: pricing.PricingContext{
			Channel:   "online",
			Region:    "US",
			Currency:  "USD",
			Timestamp: time.Now(),
		},
	}

	// Calculate pricing
	result, err := calc.Calculate(input)
	if err != nil {
		fmt.Printf("Error calculating pricing: %v\n", err)
		return
	}

	fmt.Printf("Subtotal: %.2f\n", result.Subtotal)
	fmt.Printf("Total Discount: %.2f\n", result.TotalDiscount)
	fmt.Printf("Grand Total: %.2f\n", result.GrandTotal)
}

// discountExample demonstrates discount calculations
func discountExample() {
	fmt.Println("\n=== Discount Example ===")

	// Create discount items
	items := []discount.DiscountItem{
		{
			ID:       "item1",
			Price:    100.0,
			Quantity: 5,
			Category: "electronics",
		},
	}

	// Create bulk discount rule
	bulkRule := discount.BulkDiscountRule{
		MinQuantity:          3,
		DiscountType:         "percentage",
		DiscountValue:        10.0,
		ApplicableCategories: []string{"electronics"},
	}

	// Create discount input
	input := discount.DiscountCalculationInput{
		Items:                     items,
		BulkRules:                 []discount.BulkDiscountRule{bulkRule},
		AllowStacking:             true,
		MaxStackedDiscountPercent: 50.0,
	}

	// Calculate discounts
	result := discount.Calculate(input)

	fmt.Printf("Original Amount: %.2f\n", result.OriginalAmount)
	fmt.Printf("Total Discount: %.2f\n", result.TotalDiscount)
	fmt.Printf("Final Amount: %.2f\n", result.FinalAmount)
}

// couponExample demonstrates coupon generation and validation
func couponExample() {
	fmt.Println("\n=== Coupon Example ===")

	// Generate coupon codes
	config := coupon.GeneratorConfig{
		Length:    8,
		Count:     5,
		Prefix:    "SAVE",
		Pattern:   "PREFIX-XXXXXX",
	}

	codes, err := coupon.GenerateCodes(config)
	if err != nil {
		fmt.Printf("Error generating codes: %v\n", err)
		return
	}
	fmt.Printf("Generated coupon codes: %v\n", codes)

	// Create a coupon
	couponData := coupon.Coupon{
		Code:         codes[0],
		Type:         coupon.CouponTypePercentage,
		Value:        15.0,
		MinOrder:     50.0,
		MaxDiscount:  25.0,
		ValidFrom:    time.Now(),
		ValidUntil:   time.Now().AddDate(0, 1, 0),
		IsActive:     true,
	}

	// Calculate coupon discount
	input := coupon.CalculationInput{
		Coupon:      couponData,
		OrderAmount: 100.0,
		UserID:      "user123",
		Items: []coupon.Item{
			{
				ID:       "item1",
				Price:    100.0,
				Quantity: 1,
				Category: "electronics",
			},
		},
		Usage: coupon.CouponUsage{
			CouponCode: codes[0],
			UserID:     "user123",
			UsageCount: 0,
			TotalUsage: 0,
		},
	}

	result := coupon.Calculate(input)

	fmt.Printf("Coupon Code: %s\n", couponData.Code)
	fmt.Printf("Discount Amount: %.2f\n", result.DiscountAmount)
	fmt.Printf("Is Valid: %t\n", result.IsValid)
	if result.ErrorMessage != "" {
		fmt.Printf("Error: %s\n", result.ErrorMessage)
	}
}

// shippingExample demonstrates shipping cost calculation
func shippingExample() {
	fmt.Println("--- Shipping Example ---")

	// Create shipping input
	input := shipping.ShippingCalculationInput{
		Items: []shipping.ShippingItem{
			{
				ID:       "item-001",
				Name:     "Laptop",
				Quantity: 1,
				Weight:   shipping.Weight{Value: 2.5, Unit: shipping.WeightUnitKG},
				Dimensions: shipping.Dimensions{
					Length: 35.0,
					Width:  25.0,
					Height: 3.0,
					Unit:   shipping.DimensionUnitCM,
				},
				Value:    1500.0,
				Category: "electronics",
			},
		},
		Origin: shipping.Address{
			Street1:    "123 Warehouse St",
			City:       "Jakarta",
			State:      "DKI Jakarta",
			PostalCode: "12345",
			Country:    "Indonesia",
		},
		Destination: shipping.Address{
			Street1:    "456 Customer Ave",
			City:       "Surabaya",
			State:      "East Java",
			PostalCode: "67890",
			Country:    "Indonesia",
		},
		ShippingRules: []shipping.ShippingRule{
			{
				ID:         "standard-rule",
				Name:       "Standard Shipping",
				Method:     shipping.ShippingMethodStandard,
				Zone:       shipping.ShippingZoneNational,
				BaseCost:   10.0,
				WeightRate: 2.0,
				IsActive:   true,
				ValidFrom:  time.Now().AddDate(0, -1, 0),
				ValidUntil: time.Now().AddDate(1, 0, 0),
			},
		},
	}

	// Calculate shipping
	result := shipping.Calculate(input)

	fmt.Printf("Shipping Result:\n")
	fmt.Printf("  Is Valid: %t\n", result.IsValid)
	fmt.Printf("  Zone: %s\n", result.Zone)
	fmt.Printf("  Total Weight: %.2f kg\n", result.TotalWeight.Value)
	fmt.Printf("  Total Value: $%.2f\n", result.TotalValue)

	if len(result.Options) > 0 {
		fmt.Printf("  Available Options:\n")
		for _, option := range result.Options {
			fmt.Printf("    - %s: $%.2f (%d days)\n", option.ServiceName, option.Cost, option.EstimatedDays)
		}
	}

	if result.ErrorMessage != "" {
		fmt.Printf("  Error: %s\n", result.ErrorMessage)
	}
}

// taxExample demonstrates tax calculation
func taxExample() {
	fmt.Println("--- Tax Example ---")

	// Create tax input
	input := tax.TaxCalculationInput{
		Items: []tax.TaxableItem{
			{
				ID:          "item-001",
				Name:        "Laptop",
				Category:    "electronics",
				Quantity:    1,
				UnitPrice:   1000.0,
				TotalAmount: 1000.0,
			},
		},
		Customer: tax.Customer{
			ID:   "customer-001",
			Type: "individual",
		},
		BillingAddress: tax.Address{
			Street1:    "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "US",
		},
		ShippingAddress: tax.Address{
			Street1:    "123 Main St",
			City:       "New York",
			State:      "NY",
			PostalCode: "10001",
			Country:    "US",
		},
		TransactionDate: time.Now(),
		TransactionType: "sale",
		Currency:        "USD",
		TaxRules: []tax.TaxRule{
			{
				ID:           "ny-sales-tax",
				Name:         "NY Sales Tax",
				Type:         tax.TaxTypeSales,
				Jurisdiction: tax.JurisdictionState,
				Method:       tax.TaxMethodPercentage,
				Rate:         8.25,
				ApplicableStates: []string{"NY"},
				IsActive:     true,
				ValidFrom:    time.Now().AddDate(-1, 0, 0),
				ValidUntil:   time.Now().AddDate(1, 0, 0),
			},
		},
	}

	// Calculate tax
	result := tax.Calculate(input)

	fmt.Printf("Tax Result:\n")
	fmt.Printf("  Subtotal: $%.2f\n", result.Subtotal)
	fmt.Printf("  Total Tax: $%.2f\n", result.TotalTax)
	fmt.Printf("  Grand Total: $%.2f\n", result.GrandTotal)
	fmt.Printf("  Currency: %s\n", result.Currency)
	fmt.Printf("  Effective Rate: %.2f%%\n", result.EffectiveRate)
	fmt.Printf("  Tax Breakdown: %d items\n", len(result.TaxBreakdown))

	for _, breakdown := range result.TaxBreakdown {
		fmt.Printf("  - %s: $%.2f tax on $%.2f\n", breakdown.ItemName, breakdown.TotalTax, breakdown.TaxableAmount)
	}
}

// loyaltyExample demonstrates loyalty points calculation
func loyaltyExample() {
	fmt.Println("--- Loyalty Example ---")

	// Create loyalty configuration
	config := &loyalty.LoyaltyConfiguration{
		BasePointsRate: 1.0,
		TierBenefits: map[loyalty.LoyaltyTier]loyalty.TierBenefit{
			loyalty.TierGold: {
				Tier:               loyalty.TierGold,
				PointsMultiplier:   1.5,
				BonusPointsPercent: 10.0,
				RedemptionBonus:    0.2,
				BirthdayBonus:      100,
				MaxPointsExpiry:    24,
			},
		},
		DefaultRules: []loyalty.LoyaltyRule{
			{
				ID:        "base-points",
				Name:      "Base Points Rule",
				Type:        "earn",
				IsActive:  true,
				ValidFrom: time.Now().AddDate(0, -1, 0),
				ValidUntil: time.Now().AddDate(1, 0, 0),
				Priority:  1,
				Actions: []loyalty.LoyaltyAction{
					{
						Type:       "earn_points",
						Value:      1.0,
						PointsType: loyalty.PointsTypeBase,
					},
				},
			},
		},
	}

	// Create loyalty calculator
	calc := loyalty.NewCalculator(config)

	// Create loyalty input
	input := loyalty.PointsCalculationInput{
		Customer: loyalty.Customer{
			ID:            "customer-001",
			Tier:          loyalty.TierGold,
			CurrentPoints: 1000,
			AnnualSpend:   5000.0,
			TotalSpend:    15000.0,
			JoinDate:      time.Now().AddDate(-1, 0, 0),
			IsActive:      true,
		},
		OrderAmount: 250.0,
		Timestamp:   time.Now(),
		OrderID:     "order-001",
		Channel:     "online",
	}

	// Calculate loyalty points
	result, err := calc.Calculate(input)
	if err != nil {
		log.Printf("Error calculating loyalty points: %v", err)
		return
	}

	fmt.Printf("Loyalty Points Calculation:\n")
	fmt.Printf("  Base Points: %d\n", result.BasePoints)
	fmt.Printf("  Bonus Points: %d\n", result.BonusPoints)
	fmt.Printf("  Total Points Earned: %d\n", result.TotalPoints)
	fmt.Printf("  New Balance: %d\n", result.NewBalance)
	fmt.Printf("  Customer Tier: %s\n", input.Customer.Tier)

	if len(result.AppliedRules) > 0 {
		fmt.Printf("  Applied Rules:\n")
		for _, rule := range result.AppliedRules {
			fmt.Printf("    - %s\n", rule.Name)
		}
	}

	if len(result.PointsBreakdown) > 0 {
		fmt.Printf("  Points Breakdown:\n")
		for _, breakdown := range result.PointsBreakdown {
			fmt.Printf("    - %s: %d points\n", breakdown.Description, breakdown.Points)
		}
	}
}

// completeOrderExample demonstrates a complete order flow using multiple packages
func completeOrderExample() {
	fmt.Println("--- Complete Order Example ---")

	// Order details
	orderItems := []struct {
		ID       string
		Name     string
		Price    float64
		Quantity int
		Category string
		Weight   float64
	}{
		{"laptop-001", "Gaming Laptop", 1200.0, 1, "electronics", 2.5},
		{"mouse-001", "Gaming Mouse", 45.0, 1, "accessories", 0.2},
		{"keyboard-001", "Mechanical Keyboard", 85.0, 1, "accessories", 0.8},
	}

	customerInfo := struct {
		ID     string
		Email  string
		Tier   string
		Points int
	}{"customer-001", "john@example.com", "silver", 750}

	fmt.Printf("Customer: %s (%s tier, %d points)\n", customerInfo.Email, customerInfo.Tier, customerInfo.Points)
	fmt.Println("Order Items:")
	subtotal := 0.0
	for _, item := range orderItems {
		itemTotal := item.Price * float64(item.Quantity)
		subtotal += itemTotal
		fmt.Printf("  - %s: $%.2f x %d = $%.2f\n", item.Name, item.Price, item.Quantity, itemTotal)
	}
	fmt.Printf("Subtotal: $%.2f\n\n", subtotal)

	// 1. Apply discount
	discountInput := discount.DiscountCalculationInput{
			Items: []discount.DiscountItem{
				{
					ID:       "laptop-001",
					Price:    1200.0,
					Quantity: 1,
					Category: "electronics",
				},
				{
					ID:       "mouse-001",
					Price:    45.0,
					Quantity: 1,
					Category: "accessories",
				},
				{
					ID:       "keyboard-001",
					Price:    85.0,
					Quantity: 1,
					Category: "accessories",
				},
			},
			BulkRules: []discount.BulkDiscountRule{
				{
					MinQuantity:          2,
					DiscountType:         "percentage",
					DiscountValue:        10.0,
					ApplicableCategories: []string{"electronics", "accessories"},
				},
			},
			AllowStacking: false,
		}

	discountResult := discount.Calculate(discountInput)
	fmt.Printf("Discount Applied: $%.2f (%.2f%%)\n", discountResult.TotalDiscount, discountResult.SavingsPercent)
	afterDiscount := discountResult.FinalAmount

	// 2. Apply coupon
	couponData := coupon.Coupon{
			Code:          "WELCOME10",
			Type:          coupon.CouponTypePercentage,
			Value:         10.0,
			MinOrder:      100.0,
			MaxDiscount:   50.0,
			IsActive:      true,
			ValidFrom:     time.Now().AddDate(0, -1, 0),
			ValidUntil:    time.Now().AddDate(0, 1, 0),
		}

	couponCalcInput := coupon.CalculationInput{
			Coupon:      couponData,
			OrderAmount: afterDiscount,
			UserID:      customerInfo.ID,
			Items: []coupon.Item{
				{
					ID:       "laptop-001",
					Price:    1200.0,
					Quantity: 1,
					Category: "electronics",
				},
				{
					ID:       "mouse-001",
					Price:    45.0,
					Quantity: 1,
					Category: "accessories",
				},
				{
					ID:       "keyboard-001",
					Price:    85.0,
					Quantity: 1,
					Category: "accessories",
				},
			},
			Usage: coupon.CouponUsage{
				CouponCode: "WELCOME10",
				UserID:     customerInfo.ID,
				UsageCount: 0,
				TotalUsage: 0,
			},
		}

	couponResult := coupon.Calculate(couponCalcInput)
	fmt.Printf("Coupon Applied: $%.2f\n", couponResult.DiscountAmount)
	afterCoupon := afterDiscount - couponResult.DiscountAmount

	// 3. Calculate shipping
	shippingInput := shipping.ShippingCalculationInput{
		Items: []shipping.ShippingItem{
			{
				ID:       "laptop-001",
				Quantity: 1,
				Weight:   shipping.Weight{Value: 2.5, Unit: shipping.WeightUnitKG},
				Value:    1200.0,
			},
			{
				ID:       "mouse-001",
				Quantity: 1,
				Weight:   shipping.Weight{Value: 0.2, Unit: shipping.WeightUnitKG},
				Value:    45.0,
			},
			{
				ID:       "keyboard-001",
				Quantity: 1,
				Weight:   shipping.Weight{Value: 0.8, Unit: shipping.WeightUnitKG},
				Value:    85.0,
			},
		},
		Origin: shipping.Address{
			City:    "New York",
			State:   "NY",
			Country: "US",
		},
		Destination: shipping.Address{
			City:    "Los Angeles",
			State:   "CA",
			Country: "US",
		},
		ShippingRules: []shipping.ShippingRule{
			{
				ID:         "standard-shipping",
				Name:       "Standard Shipping",
				Method:     shipping.ShippingMethodStandard,
				Zone:       shipping.ShippingZoneNational,
				BaseCost:   15.0,
				WeightRate: 2.0,
				IsActive:   true,
				ValidFrom:  time.Now().AddDate(0, -1, 0),
				ValidUntil: time.Now().AddDate(1, 0, 0),
			},
		},
	}

	shippingResult := shipping.Calculate(shippingInput)
	shippingCost := 0.0
	if len(shippingResult.Options) > 0 {
		shippingCost = shippingResult.Options[0].Cost // Use first option
		fmt.Printf("Shipping: $%.2f (%s)\n", shippingCost, shippingResult.Options[0].ServiceName)
	}

	// 4. Calculate tax
	taxInput := tax.TaxCalculationInput{
		Items: []tax.TaxableItem{
			{
				ID:          "laptop-001",
				Name:        "Gaming Laptop",
				Category:    "electronics",
				Quantity:    1,
				UnitPrice:   1200.0,
				TotalAmount: 1200.0,
			},
			{
				ID:          "mouse-001",
				Name:        "Gaming Mouse",
				Category:    "accessories",
				Quantity:    1,
				UnitPrice:   45.0,
				TotalAmount: 45.0,
			},
			{
				ID:          "keyboard-001",
				Name:        "Mechanical Keyboard",
				Category:    "accessories",
				Quantity:    1,
				UnitPrice:   85.0,
				TotalAmount: 85.0,
			},
		},
		Customer: tax.Customer{ID: "customer-001", Type: "individual"},
		BillingAddress: tax.Address{
			City:    "Los Angeles",
			State:   "CA",
			Country: "US",
		},
		ShippingAddress: tax.Address{
			City:    "Los Angeles",
			State:   "CA",
			Country: "US",
		},
		TransactionDate: time.Now(),
		Currency:        "USD",
		TaxRules: []tax.TaxRule{
			{
				ID:           "ca-sales-tax",
				Name:         "CA Sales Tax",
				Type:         tax.TaxTypeSales,
				Jurisdiction: tax.JurisdictionState,
				Method:       tax.TaxMethodPercentage,
				Rate:         7.25,
				ApplicableStates: []string{"CA"},
				IsActive:     true,
				ValidFrom:    time.Now().AddDate(-1, 0, 0),
				ValidUntil:   time.Now().AddDate(1, 0, 0),
			},
		},
	}

	taxResult := tax.Calculate(taxInput)
	fmt.Printf("Tax: $%.2f\n", taxResult.TotalTax)

	// 5. Calculate loyalty points
	loyaltyConfig := &loyalty.LoyaltyConfiguration{
		BasePointsRate: 1.0,
		TierBenefits: map[loyalty.LoyaltyTier]loyalty.TierBenefit{
			loyalty.TierSilver: {
				Tier:             loyalty.TierSilver,
				PointsMultiplier: 1.2,
				RedemptionBonus:  0.1,
			},
		},
		DefaultRules: []loyalty.LoyaltyRule{
			{
				ID:        "base-points",
				Name:      "Base Points",
				Type:        "earn",
				IsActive:  true,
				ValidFrom: time.Now().AddDate(0, -1, 0),
				ValidUntil: time.Now().AddDate(1, 0, 0),
			},
		},
	}

	loyaltyCalc := loyalty.NewCalculator(loyaltyConfig)
	loyaltyInput := loyalty.PointsCalculationInput{
		Customer: loyalty.Customer{
			ID:            customerInfo.ID,
			Email:         customerInfo.Email,
			Tier:          loyalty.TierSilver,
			CurrentPoints: customerInfo.Points,
			JoinDate:      time.Now().AddDate(-1, 0, 0),
			IsActive:      true,
		},
		OrderAmount: afterCoupon,
		Timestamp:   time.Now(),
		OrderID:     utils.GenerateShortID(8),
		Channel:     "online",
	}

	loyaltyResult, err := loyaltyCalc.Calculate(loyaltyInput)
	if err != nil {
		log.Printf("Error calculating loyalty points: %v", err)
	} else {
		fmt.Printf("Loyalty Points Earned: %d\n", loyaltyResult.TotalPoints)
	}

	// 6. Final totals
	finalTotal := afterCoupon + shippingCost + taxResult.TotalTax

	fmt.Println("\n--- Order Summary ---")
	fmt.Printf("Subtotal: $%.2f\n", subtotal)
	fmt.Printf("Discount: -$%.2f\n", discountResult.TotalDiscount)
	fmt.Printf("Coupon: -$%.2f\n", couponResult.DiscountAmount)
	fmt.Printf("Shipping: $%.2f\n", shippingCost)
	fmt.Printf("Tax: $%.2f\n", taxResult.TotalTax)
	fmt.Printf("Final Total: $%.2f\n", finalTotal)
	if loyaltyResult != nil {
		fmt.Printf("Points Earned: %d\n", loyaltyResult.TotalPoints)
	}
	fmt.Printf("Total Savings: $%.2f\n", discountResult.TotalDiscount+couponResult.DiscountAmount)

	// 7. Currency conversion example
	currencyCalc := currency.NewCalculator()
	currencyCalc.SetExchangeRate(currency.USD, currency.IDR, 15000.0, "Bank Indonesia")

	conversionInput := currency.ConversionInput{
		Amount: finalTotal,
		From:   currency.USD,
		To:     currency.IDR,
	}

	idrResult, err := currencyCalc.Convert(conversionInput)
	if err != nil {
		log.Printf("Error converting to IDR: %v", err)
	} else {
		formattedIDR, _ := currencyCalc.Format(idrResult.ConvertedAmount, &currency.FormatOptions{ShowSymbol: true})
		fmt.Printf("Total in IDR: %s\n", formattedIDR)
	}
}
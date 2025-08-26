package pricing

import (
	"testing"
	"time"
)

func TestNewCalculator(t *testing.T) {
	calc := NewCalculator()

	if calc == nil {
		t.Fatal("Expected calculator to be created")
	}
	if calc.rules == nil {
		t.Error("Expected rules to be initialized")
	}
	if calc.bundles == nil {
		t.Error("Expected bundles to be initialized")
	}
	if calc.tierPricing == nil {
		t.Error("Expected tierPricing to be initialized")
	}
	if calc.dynamicConfigs == nil {
		t.Error("Expected dynamicConfigs to be initialized")
	}
	if calc.marketData == nil {
		t.Error("Expected marketData to be initialized")
	}
	if calc.analytics == nil {
		t.Error("Expected analytics to be initialized")
	}
}

func TestCalculate(t *testing.T) {
	calc := NewCalculator()

	// Add test pricing rule
	calc.AddRule(PricingRule{
		ID:          "test-rule-1",
		Name:        "Test Discount",
		Type:        PricingTypePromo,
		Strategy:    StrategyFixed,
		IsActive:    true,
		Priority:    1,
		Description: "Test discount rule",
		Adjustments: []PriceAdjustment{
			{
				Type:  "percentage",
				Value: 10.0,
			},
		},
	})

	// Test basic pricing calculation
	input := PricingInput{
		Items: []PricingItem{
			{
				ID:        "item1",
				BasePrice: 100.0,
				Quantity:  2,
				Category:  "electronics",
			},
		},
		Customer: Customer{
			ID:   "customer1",
			Type: "regular",
			Tier: "bronze",
		},
		Context: PricingContext{
			Timestamp: time.Now(),
			Channel:   "online",
		},
	}

	result, err := calc.Calculate(input)

	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}
	if len(result.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(result.Items))
		return
	}
	if result.Items[0].ItemID != "item1" {
		t.Errorf("Expected item ID 'item1', got '%s'", result.Items[0].ItemID)
	}
	if result.Items[0].OriginalPrice != 100.0 {
		t.Errorf("Expected original price 100.0, got %f", result.Items[0].OriginalPrice)
	}
	if result.Items[0].Quantity != 2 {
		t.Errorf("Expected quantity 2, got %d", result.Items[0].Quantity)
	}

	// Test empty items
	emptyInput := PricingInput{
		Items: []PricingItem{},
	}

	_, err = calc.Calculate(emptyInput)
	if err == nil {
		t.Error("Expected error for empty items but got none")
	}

	// Test negative price
	negativeInput := PricingInput{
		Items: []PricingItem{
			{
				ID:        "item1",
				BasePrice: -100.0,
				Quantity:  1,
			},
		},
	}

	_, err = calc.Calculate(negativeInput)
	if err == nil {
		t.Error("Expected error for negative price but got none")
	}
}

func TestCalculateItemPricing(t *testing.T) {
	calc := NewCalculator()

	// Add test pricing rule
	calc.AddRule(PricingRule{
		ID:       "discount-rule",
		Name:     "Discount Rule",
		Type:     PricingTypePromo,
		Strategy: StrategyFixed,
		IsActive: true,
		Priority: 1,
		Adjustments: []PriceAdjustment{
			{
				Type:  "percentage",
				Value: 15.0,
			},
		},
	})

	item := PricingItem{
		ID:        "test-item",
		BasePrice: 200.0,
		Quantity:  1,
		Category:  "books",
	}

	customer := Customer{
		ID:   "customer1",
		Type: "premium",
		Tier: "gold",
	}

	context := PricingContext{
		Timestamp: time.Now(),
		Channel:   "online",
	}

	pricedItem, err := calc.calculateItemPricing(item, customer, context, []PricingRule{}, []TierPricing{}, PricingOptions{})

	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}
	if pricedItem == nil {
		t.Fatal("Expected pricedItem to not be nil")
	}
	if pricedItem.ItemID != "test-item" {
		t.Errorf("Expected item ID 'test-item', got '%s'", pricedItem.ItemID)
	}
	if pricedItem.OriginalPrice != 200.0 {
		t.Errorf("Expected original price 200.0, got %f", pricedItem.OriginalPrice)
	}
	if pricedItem.Quantity != 1 {
		t.Errorf("Expected quantity 1, got %d", pricedItem.Quantity)
	}
}

func TestCalculateDynamicPricing(t *testing.T) {
	calc := NewCalculator()

	// Add dynamic pricing config
	calc.AddDynamicConfig(DynamicPricingConfig{
		ID:       "dynamic-1",
		IsActive: true,
		Factors: []PricingFactor{
			{
				Type:   "demand",
				Weight: 0.3,
			},
			{
				Type:   "inventory",
				Weight: 0.2,
			},
		},
	})

	// Add market data
	calc.UpdateMarketData("item1", MarketData{
		ItemID:          "item1",
		AveragePrice:    100.0,
		MinPrice:        90.0,
		MaxPrice:        110.0,
		DemandLevel:     "high",
		TrendDirection:  "up",
		LastUpdated:     time.Now(),
	})

	item := PricingItem{
		ID:             "item1",
		BasePrice:      100.0,
		Quantity:       1,
		InventoryLevel: 50,
	}

	context := PricingContext{
		Timestamp: time.Now(),
		Channel:   "online",
	}

	adjustedPrice := calc.calculateDynamicPricing(item, context)

	// Dynamic pricing may return 0 if no config applies
	if adjustedPrice < 0 {
		t.Error("Expected adjusted price to be >= 0")
	}
}

func TestApplyAdjustment(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name       string
		price      float64
		adjustment PriceAdjustment
		expected   float64
	}{
		{
			name:  "percentage discount",
			price: 100.0,
			adjustment: PriceAdjustment{
				Type:  "percentage",
				Value: 10.0,
			},
			expected: 90.0,
		},
		{
			name:  "fixed discount",
			price: 100.0,
			adjustment: PriceAdjustment{
				Type:  "fixed",
				Value: 15.0,
			},
			expected: 85.0,
		},
		{
			name:  "markup",
			price: 100.0,
			adjustment: PriceAdjustment{
				Type:  "markup",
				Value: 20.0,
			},
			expected: 120.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.applyAdjustment(tt.price, tt.adjustment)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestValidateInput(t *testing.T) {
	calc := NewCalculator()

	// Test valid input
	validInput := PricingInput{
		Items: []PricingItem{
			{
				ID:        "item1",
				BasePrice: 100.0,
				Quantity:  1,
			},
		},
	}

	err := calc.validateInput(validInput)
	if err != nil {
		t.Errorf("Expected no error for valid input, got: %v", err)
	}

	// Test empty items
	emptyInput := PricingInput{
		Items: []PricingItem{},
	}

	err = calc.validateInput(emptyInput)
	if err == nil {
		t.Error("Expected error for empty items")
	}

	// Test negative price
	negativeInput := PricingInput{
		Items: []PricingItem{
			{
				ID:        "item1",
				BasePrice: -100.0,
				Quantity:  1,
			},
		},
	}

	err = calc.validateInput(negativeInput)
	if err == nil {
		t.Error("Expected error for negative price")
	}

	// Test zero quantity
	zeroQuantityInput := PricingInput{
		Items: []PricingItem{
			{
				ID:        "item1",
				BasePrice: 100.0,
				Quantity:  0,
			},
		},
	}

	err = calc.validateInput(zeroQuantityInput)
	if err == nil {
		t.Error("Expected error for zero quantity")
	}
}

func TestCalculatorConfiguration(t *testing.T) {
	calc := NewCalculator()

	// Test AddRule
	rule := PricingRule{
		ID:       "test-rule",
		Name:     "Test Rule",
		IsActive: true,
	}
	calc.AddRule(rule)
	if len(calc.rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(calc.rules))
	}
	if calc.rules[0].ID != "test-rule" {
		t.Errorf("Expected rule ID 'test-rule', got '%s'", calc.rules[0].ID)
	}

	// Test AddBundle
	bundle := Bundle{
		ID:       "test-bundle",
		Name:     "Test Bundle",
		IsActive: true,
	}
	calc.AddBundle(bundle)
	if len(calc.bundles) != 1 {
		t.Errorf("Expected 1 bundle, got %d", len(calc.bundles))
	}
	if calc.bundles[0].ID != "test-bundle" {
		t.Errorf("Expected bundle ID 'test-bundle', got '%s'", calc.bundles[0].ID)
	}

	// Test AddTierPricing
	tier := TierPricing{
		ID:       "test-tier",
		Name:     "Test Tier",
		IsActive: true,
	}
	calc.AddTierPricing(tier)
	if len(calc.tierPricing) != 1 {
		t.Errorf("Expected 1 tier pricing, got %d", len(calc.tierPricing))
	}
	if calc.tierPricing[0].ID != "test-tier" {
		t.Errorf("Expected tier ID 'test-tier', got '%s'", calc.tierPricing[0].ID)
	}

	// Test AddDynamicConfig
	config := DynamicPricingConfig{
		ID:       "test-config",
		IsActive: true,
	}
	calc.AddDynamicConfig(config)
	if len(calc.dynamicConfigs) != 1 {
		t.Errorf("Expected 1 dynamic config, got %d", len(calc.dynamicConfigs))
	}
	if calc.dynamicConfigs[0].ID != "test-config" {
		t.Errorf("Expected config ID 'test-config', got '%s'", calc.dynamicConfigs[0].ID)
	}

	// Test UpdateMarketData
	marketData := MarketData{
		ItemID:          "item1",
		AveragePrice:    100.0,
		DemandLevel:     "high",
		TrendDirection:  "up",
		LastUpdated:     time.Now(),
	}
	calc.UpdateMarketData("item1", marketData)
	if _, exists := calc.marketData["item1"]; !exists {
		t.Error("Expected market data to be stored for item1")
	}
	if calc.marketData["item1"].AveragePrice != 100.0 {
		t.Errorf("Expected average price 100.0, got %f", calc.marketData["item1"].AveragePrice)
	}

	// Test UpdateAnalytics
	analytics := PricingAnalytics{
		ItemID:         "item1",
		AveragePrice:   100.0,
		ConversionRate: 0.15,
		Margin:         0.25,
	}
	calc.UpdateAnalytics("item1", analytics)
	if _, exists := calc.analytics["item1"]; !exists {
		t.Error("Expected analytics to be stored for item1")
	}
	if calc.analytics["item1"].ConversionRate != 0.15 {
		t.Errorf("Expected conversion rate 0.15, got %f", calc.analytics["item1"].ConversionRate)
	}
}

// Benchmarks

func BenchmarkCalculate(b *testing.B) {
	calc := NewCalculator()

	// Add some test data
	calc.AddRule(PricingRule{
		ID:       "rule1",
		Name:     "Test Rule",
		Type:     PricingTypePromo,
		Strategy: StrategyFixed,
		IsActive: true,
		Priority: 1,
		Adjustments: []PriceAdjustment{
			{Type: "percentage", Value: 10.0},
		},
	})

	input := PricingInput{
		Items: []PricingItem{
			{
				ID:        "item1",
				BasePrice: 100.0,
				Quantity:  1,
				Category:  "electronics",
			},
			{
				ID:        "item2",
				BasePrice: 50.0,
				Quantity:  2,
				Category:  "books",
			},
		},
		Customer: Customer{
			ID:   "customer1",
			Type: "regular",
			Tier: "bronze",
		},
		Context: PricingContext{
			Timestamp: time.Now(),
			Channel:   "online",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = calc.Calculate(input)
	}
}
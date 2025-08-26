package discount

import (
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	t.Run("EmptyItems", func(t *testing.T) {
		input := DiscountCalculationInput{
			Items: []DiscountItem{},
		}
		
		result := Calculate(input)
		
		if result.IsValid {
			t.Error("Expected invalid result for empty items")
		}
		
		if result.ErrorMessage == "" {
			t.Error("Expected error message for empty items")
		}
	})
	
	t.Run("BasicCalculation", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 2, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 1, Category: "books"},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		expectedOriginal := 250.0 // (100*2) + (50*1)
		if result.OriginalAmount != expectedOriginal {
			t.Errorf("Expected original amount %f, got %f", expectedOriginal, result.OriginalAmount)
		}
		
		if result.FinalAmount != result.OriginalAmount-result.TotalDiscount {
			t.Error("Final amount calculation is incorrect")
		}
	})
	
	t.Run("BulkDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		}
		
		bulkRules := []BulkDiscountRule{
			{
				MinQuantity: 3,
				DiscountType: "percentage",
				DiscountValue: 10, // 10% discount
				ApplicableCategories: []string{"electronics"},
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			BulkRules: bulkRules,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		expectedDiscount := 50.0 // 10% of 500
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
	
	t.Run("TierPricing", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 10, Category: "electronics"},
		}
		
		tierRules := []TierPricingRule{
			{
				MinQuantity: 10,
				PricePerItem: 90, // Reduced price per item
				Category: "electronics",
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			TierRules: tierRules,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		expectedDiscount := 100.0 // (100-90) * 10
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
	
	t.Run("BundleDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "laptop", Price: 1000, Quantity: 1, Category: "electronics"},
			{ID: "mouse", Price: 50, Quantity: 1, Category: "accessories"},
		}
		
		bundleRules := []BundleDiscountRule{
			{
				ID: "laptop_bundle",
				Name: "Laptop Bundle",
				RequiredProducts: []string{"laptop", "mouse"},
				MinItems: 2,
				DiscountType: "percentage",
				DiscountValue: 5, // 5% bundle discount
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			BundleRules: bundleRules,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		expectedDiscount := 52.5 // 5% of 1050
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
	
	t.Run("LoyaltyDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 2, Category: "electronics"},
		}
		
		customer := Customer{
			ID: "customer1",
			LoyaltyTier: "gold",
		}
		
		loyaltyRules := []LoyaltyDiscountRule{
			{
				Tier: "gold",
				DiscountPercent: 15,
				MinOrderAmount: 100,
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			Customer: customer,
			LoyaltyRules: loyaltyRules,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		expectedDiscount := 30.0 // 15% of 200
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
	
	t.Run("CategoryDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 2, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 1, Category: "books"},
		}
		
		now := time.Now()
		categoryRules := []CategoryDiscountRule{
			{
				Category: "electronics",
				DiscountPercent: 20,
				MinQuantity: 1,
				ValidFrom: now.Add(-time.Hour),
				ValidUntil: now.Add(time.Hour),
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			CategoryRules: categoryRules,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		expectedDiscount := 40.0 // 20% of 200 (electronics items only)
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
	
	t.Run("ProgressiveDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 15, Category: "electronics"},
		}
		
		progressiveRules := []ProgressiveDiscountRule{
			{
				QuantityStep: 10, // Every 10 items
				DiscountPercent: 5, // Additional 5% discount
				MaxDiscount: 20, // Maximum 20% total discount
				Category: "electronics",
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			ProgressiveRules: progressiveRules,
			AllowStacking: false,
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		// 15 items = 1 step of 10, so 5% discount
		expectedDiscount := 75.0 // 5% of 1500
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
	
	t.Run("StackedDiscounts", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		}
		
		customer := Customer{
			ID: "customer1",
			LoyaltyTier: "silver",
		}
		
		bulkRules := []BulkDiscountRule{
			{
				MinQuantity: 3,
				DiscountType: "percentage",
				DiscountValue: 10,
				ApplicableCategories: []string{"electronics"},
			},
		}
		
		loyaltyRules := []LoyaltyDiscountRule{
			{
				Tier: "silver",
				DiscountPercent: 5,
				MinOrderAmount: 100,
			},
		}
		
		input := DiscountCalculationInput{
			Items: items,
			Customer: customer,
			BulkRules: bulkRules,
			LoyaltyRules: loyaltyRules,
			AllowStacking: true,
			MaxStackedDiscountPercent: 20, // Maximum 20% total discount
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		// Should apply both bulk (10%) and loyalty (5%) discounts
		// But limited by max stacked discount of 20%
		expectedMaxDiscount := 100.0 // 20% of 500
		if result.TotalDiscount > expectedMaxDiscount {
			t.Errorf("Discount exceeds maximum allowed: got %f, max %f", result.TotalDiscount, expectedMaxDiscount)
		}
		
		if len(result.AppliedDiscounts) < 2 {
			t.Error("Expected multiple discounts to be applied")
		}
	})
}

func TestCalculateBestDiscount(t *testing.T) {
	t.Run("MultipleInputs", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		}
		
		// Input 1: Bulk discount only
		input1 := DiscountCalculationInput{
			Items: items,
			BulkRules: []BulkDiscountRule{
				{
					MinQuantity: 3,
					DiscountType: "percentage",
					DiscountValue: 10,
				},
			},
			AllowStacking: false,
		}
		
		// Input 2: Loyalty discount only
		input2 := DiscountCalculationInput{
			Items: items,
			Customer: Customer{LoyaltyTier: "gold"},
			LoyaltyRules: []LoyaltyDiscountRule{
				{
					Tier: "gold",
					DiscountPercent: 15,
					MinOrderAmount: 100,
				},
			},
			AllowStacking: false,
		}
		
		inputs := []DiscountCalculationInput{input1, input2}
		result := CalculateBestDiscount(inputs)
		
		if !result.IsValid {
			t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
		}
		
		// Should choose loyalty discount (15% = 75) over bulk discount (10% = 50)
		expectedDiscount := 75.0
		if result.TotalDiscount != expectedDiscount {
			t.Errorf("Expected best discount %f, got %f", expectedDiscount, result.TotalDiscount)
		}
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("GetApplicableItems", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 1, Category: "books"},
			{ID: "item3", Price: 75, Quantity: 1, Category: "electronics"},
		}
		
		// Test category filter
		applicable := getApplicableItems(items, []string{"electronics"}, nil)
		if len(applicable) != 2 {
			t.Errorf("Expected 2 electronics items, got %d", len(applicable))
		}
		
		// Test product filter
		applicable = getApplicableItems(items, nil, []string{"item2"})
		if len(applicable) != 1 {
			t.Errorf("Expected 1 specific item, got %d", len(applicable))
		}
		
		// Test no filters (should return all)
		applicable = getApplicableItems(items, nil, nil)
		if len(applicable) != 3 {
			t.Errorf("Expected all 3 items, got %d", len(applicable))
		}
	})
	
	t.Run("GetItemsByCategory", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 1, Category: "books"},
			{ID: "item3", Price: 75, Quantity: 1, Category: "electronics"},
		}
		
		electronics := getItemsByCategory(items, "electronics")
		if len(electronics) != 2 {
			t.Errorf("Expected 2 electronics items, got %d", len(electronics))
		}
		
		books := getItemsByCategory(items, "books")
		if len(books) != 1 {
			t.Errorf("Expected 1 book item, got %d", len(books))
		}
	})
	
	t.Run("GetTotalQuantity", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 3, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 2, Category: "books"},
		}
		
		total := getTotalQuantity(items)
		expected := 5
		if total != expected {
			t.Errorf("Expected total quantity %d, got %d", expected, total)
		}
	})
	
	t.Run("CalculateItemsAmount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 2, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 3, Category: "books"},
		}
		
		total := calculateItemsAmount(items)
		expected := 350.0 // (100*2) + (50*3)
		if total != expected {
			t.Errorf("Expected total amount %f, got %f", expected, total)
		}
	})
	
	t.Run("CalculateBulkDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		}
		
		// Test percentage discount
		rule := BulkDiscountRule{
			DiscountType: "percentage",
			DiscountValue: 10,
		}
		discount := calculateBulkDiscount(items, rule)
		expected := 50.0 // 10% of 500
		if discount != expected {
			t.Errorf("Expected percentage discount %f, got %f", expected, discount)
		}
		
		// Test fixed amount discount
		rule.DiscountType = "fixed_amount"
		rule.DiscountValue = 75
		discount = calculateBulkDiscount(items, rule)
		expected = 75.0
		if discount != expected {
			t.Errorf("Expected fixed amount discount %f, got %f", expected, discount)
		}
		
		// Test fixed price discount
		rule.DiscountType = "fixed_price"
		rule.DiscountValue = 80 // 80 per item
		discount = calculateBulkDiscount(items, rule)
		expected = 100.0 // 500 - (80*5)
		if discount != expected {
			t.Errorf("Expected fixed price discount %f, got %f", expected, discount)
		}
	})
	
	t.Run("FindBundleMatches", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "laptop", Price: 1000, Quantity: 1, Category: "electronics"},
			{ID: "mouse", Price: 50, Quantity: 1, Category: "accessories"},
			{ID: "keyboard", Price: 100, Quantity: 1, Category: "accessories"},
		}
		
		// Test required products bundle
		rule := BundleDiscountRule{
			ID: "laptop_bundle",
			RequiredProducts: []string{"laptop", "mouse"},
			MinItems: 2,
		}
		
		matches := findBundleMatches(items, rule)
		if len(matches) != 1 {
			t.Errorf("Expected 1 bundle match, got %d", len(matches))
		}
		
		if len(matches[0].MatchedItems) != 2 {
			t.Errorf("Expected 2 matched items, got %d", len(matches[0].MatchedItems))
		}
		
		// Test required categories bundle
		rule = BundleDiscountRule{
			ID: "category_bundle",
			RequiredCategories: []string{"electronics", "accessories"},
			MinItems: 2,
		}
		
		matches = findBundleMatches(items, rule)
		if len(matches) != 1 {
			t.Errorf("Expected 1 category bundle match, got %d", len(matches))
		}
	})
	
	t.Run("CalculateBundleDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
			{ID: "item2", Price: 50, Quantity: 1, Category: "accessories"},
		}
		
		// Test percentage bundle discount
		rule := BundleDiscountRule{
			DiscountType: "percentage",
			DiscountValue: 10,
		}
		discount := calculateBundleDiscount(items, rule)
		expected := 15.0 // 10% of 150
		if discount != expected {
			t.Errorf("Expected bundle percentage discount %f, got %f", expected, discount)
		}
		
		// Test combo price bundle discount
		rule.DiscountType = "combo_price"
		rule.DiscountValue = 120 // Combo price
		discount = calculateBundleDiscount(items, rule)
		expected = 30.0 // 150 - 120
		if discount != expected {
			t.Errorf("Expected combo price discount %f, got %f", expected, discount)
		}
	})
}

func BenchmarkCalculate(t *testing.B) {
	items := []DiscountItem{
		{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		{ID: "item2", Price: 50, Quantity: 3, Category: "books"},
		{ID: "item3", Price: 75, Quantity: 2, Category: "electronics"},
	}
	
	input := DiscountCalculationInput{
		Items: items,
		Customer: Customer{LoyaltyTier: "gold"},
		BulkRules: []BulkDiscountRule{
			{
				MinQuantity: 3,
				DiscountType: "percentage",
				DiscountValue: 10,
			},
		},
		LoyaltyRules: []LoyaltyDiscountRule{
			{
				Tier: "gold",
				DiscountPercent: 15,
				MinOrderAmount: 100,
			},
		},
		AllowStacking: true,
	}
	
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		Calculate(input)
	}
}

func BenchmarkCalculateBestDiscount(t *testing.B) {
	items := []DiscountItem{
		{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
	}
	
	inputs := []DiscountCalculationInput{
		{
			Items: items,
			BulkRules: []BulkDiscountRule{
				{MinQuantity: 3, DiscountType: "percentage", DiscountValue: 10},
			},
		},
		{
			Items: items,
			Customer: Customer{LoyaltyTier: "gold"},
			LoyaltyRules: []LoyaltyDiscountRule{
				{Tier: "gold", DiscountPercent: 15, MinOrderAmount: 100},
			},
		},
	}
	
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		CalculateBestDiscount(inputs)
	}
}
package discount

import (
	"testing"
	"time"
)

func TestNewDiscountValidator(t *testing.T) {
	validator := NewDiscountValidator()
	
	if validator == nil {
		t.Fatal("NewDiscountValidator should not return nil")
	}
	
	if validator.MaxStackedDiscountPercent != 50.0 {
		t.Errorf("Expected default max stacked discount 50.0, got %f", validator.MaxStackedDiscountPercent)
	}
	
	if validator.MaxSingleDiscountPercent != 30.0 {
		t.Errorf("Expected default max single discount 30.0, got %f", validator.MaxSingleDiscountPercent)
	}
	
	if len(validator.AllowedCombinations) == 0 {
		t.Error("Expected default allowed combinations to be set")
	}
}

func TestValidateDiscountApplication(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
		}
		
		discount := DiscountApplication{
			Type: DiscountTypeBulk,
			DiscountAmount: 10,
			AppliedItems: items,
		}
		
		err := validator.ValidateDiscountApplication(discount, items, Customer{})
		if err != nil {
			t.Errorf("Expected valid discount, got error: %v", err)
		}
	})
	
	t.Run("NegativeDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
		}
		
		discount := DiscountApplication{
			Type: DiscountTypeBulk,
			DiscountAmount: -10,
			AppliedItems: items,
		}
		
		err := validator.ValidateDiscountApplication(discount, items, Customer{})
		if err == nil {
			t.Error("Expected error for negative discount amount")
		}
	})
	
	t.Run("NoAppliedItems", func(t *testing.T) {
		discount := DiscountApplication{
			Type: DiscountTypeBulk,
			DiscountAmount: 10,
			AppliedItems: []DiscountItem{},
		}
		
		err := validator.ValidateDiscountApplication(discount, []DiscountItem{}, Customer{})
		if err == nil {
			t.Error("Expected error for no applied items")
		}
	})
	
	t.Run("DiscountExceedsItemValue", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
		}
		
		discount := DiscountApplication{
			Type: DiscountTypeBulk,
			DiscountAmount: 150, // More than item value
			AppliedItems: items,
		}
		
		err := validator.ValidateDiscountApplication(discount, items, Customer{})
		if err == nil {
			t.Error("Expected error for discount exceeding item value")
		}
	})
	
	t.Run("ExceedsMaxSingleDiscountPercent", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
		}
		
		discount := DiscountApplication{
			Type: DiscountTypeBulk,
			DiscountAmount: 40, // 40% of 100, exceeds default 30% limit
			AppliedItems: items,
		}
		
		err := validator.ValidateDiscountApplication(discount, items, Customer{})
		if err == nil {
			t.Error("Expected error for exceeding max single discount percent")
		}
	})
}

func TestValidateStackedDiscounts(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidStackedDiscounts", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
		}
		
		discounts := []DiscountApplication{
			{
				Type: DiscountTypeBulk,
				DiscountAmount: 10,
				AppliedItems: items,
			},
			{
				Type: DiscountTypeLoyalty,
				DiscountAmount: 5,
				AppliedItems: items,
			},
		}
		
		err := validator.ValidateStackedDiscounts(discounts, 100.0)
		if err != nil {
			t.Errorf("Expected valid stacked discounts, got error: %v", err)
		}
	})
	
	t.Run("SingleDiscount", func(t *testing.T) {
		discounts := []DiscountApplication{
			{
				Type: DiscountTypeBulk,
				DiscountAmount: 10,
			},
		}
		
		err := validator.ValidateStackedDiscounts(discounts, 100.0)
		if err != nil {
			t.Errorf("Expected no error for single discount, got: %v", err)
		}
	})
	
	t.Run("TotalDiscountExceedsOriginal", func(t *testing.T) {
		discounts := []DiscountApplication{
			{
				Type: DiscountTypeBulk,
				DiscountAmount: 60,
			},
			{
				Type: DiscountTypeLoyalty,
				DiscountAmount: 50,
			},
		}
		
		err := validator.ValidateStackedDiscounts(discounts, 100.0)
		if err == nil {
			t.Error("Expected error for total discount exceeding original amount")
		}
	})
	
	t.Run("ExceedsMaxStackedPercent", func(t *testing.T) {
		discounts := []DiscountApplication{
			{
				Type: DiscountTypeBulk,
				DiscountAmount: 30,
			},
			{
				Type: DiscountTypeLoyalty,
				DiscountAmount: 25, // Total 55%, exceeds default 50% limit
			},
		}
		
		err := validator.ValidateStackedDiscounts(discounts, 100.0)
		if err == nil {
			t.Error("Expected error for exceeding max stacked discount percent")
		}
	})
	
	t.Run("InvalidCombination", func(t *testing.T) {
		discounts := []DiscountApplication{
			{
				Type: DiscountTypeProgressive,
				DiscountAmount: 10,
			},
			{
				Type: DiscountTypeBulk,
				DiscountAmount: 5,
			},
		}
		
		err := validator.ValidateStackedDiscounts(discounts, 100.0)
		if err == nil {
			t.Error("Expected error for invalid discount combination")
		}
	})
}

func TestValidateBulkDiscount(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidBulkDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		}
		
		rule := BulkDiscountRule{
			MinQuantity: 3,
			MaxQuantity: 10,
			ApplicableCategories: []string{"electronics"},
		}
		
		err := validator.ValidateBulkDiscount(rule, items)
		if err != nil {
			t.Errorf("Expected valid bulk discount, got error: %v", err)
		}
	})
	
	t.Run("InsufficientQuantity", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 2, Category: "electronics"},
		}
		
		rule := BulkDiscountRule{
			MinQuantity: 5,
			ApplicableCategories: []string{"electronics"},
		}
		
		err := validator.ValidateBulkDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for insufficient quantity")
		}
	})
	
	t.Run("ExceedsMaxQuantity", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 15, Category: "electronics"},
		}
		
		rule := BulkDiscountRule{
			MinQuantity: 3,
			MaxQuantity: 10,
			ApplicableCategories: []string{"electronics"},
		}
		
		err := validator.ValidateBulkDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for exceeding max quantity")
		}
	})
}

func TestValidateTierPricing(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidTierPricing", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 10, Category: "electronics"},
		}
		
		rule := TierPricingRule{
			MinQuantity: 5,
			Category: "electronics",
		}
		
		err := validator.ValidateTierPricing(rule, items)
		if err != nil {
			t.Errorf("Expected valid tier pricing, got error: %v", err)
		}
	})
	
	t.Run("InsufficientQuantityForTier", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 3, Category: "electronics"},
		}
		
		rule := TierPricingRule{
			MinQuantity: 5,
			Category: "electronics",
		}
		
		err := validator.ValidateTierPricing(rule, items)
		if err == nil {
			t.Error("Expected error for insufficient quantity for tier pricing")
		}
	})
}

func TestValidateBundleDiscount(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidBundleDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "laptop", Price: 1000, Quantity: 1, Category: "electronics"},
			{ID: "mouse", Price: 50, Quantity: 1, Category: "accessories"},
		}
		
		rule := BundleDiscountRule{
			RequiredProducts: []string{"laptop", "mouse"},
			MinItems: 2,
		}
		
		err := validator.ValidateBundleDiscount(rule, items)
		if err != nil {
			t.Errorf("Expected valid bundle discount, got error: %v", err)
		}
	})
	
	t.Run("BundleRequirementsNotMet", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "laptop", Price: 1000, Quantity: 1, Category: "electronics"},
		}
		
		rule := BundleDiscountRule{
			RequiredProducts: []string{"laptop", "mouse"},
			MinItems: 2,
		}
		
		err := validator.ValidateBundleDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for bundle requirements not met")
		}
	})
}

func TestValidateLoyaltyDiscount(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidLoyaltyDiscount", func(t *testing.T) {
		customer := Customer{
			LoyaltyTier: "gold",
			TotalPurchases: 1000,
		}
		
		rule := LoyaltyDiscountRule{
			Tier: "gold",
			MinOrderAmount: 500,
		}
		
		err := validator.ValidateLoyaltyDiscount(rule, customer)
		if err != nil {
			t.Errorf("Expected valid loyalty discount, got error: %v", err)
		}
	})
	
	t.Run("WrongLoyaltyTier", func(t *testing.T) {
		customer := Customer{
			LoyaltyTier: "silver",
			TotalPurchases: 1000,
		}
		
		rule := LoyaltyDiscountRule{
			Tier: "gold",
			MinOrderAmount: 500,
		}
		
		err := validator.ValidateLoyaltyDiscount(rule, customer)
		if err == nil {
			t.Error("Expected error for wrong loyalty tier")
		}
	})
	
	t.Run("InsufficientPurchases", func(t *testing.T) {
		customer := Customer{
			LoyaltyTier: "gold",
			TotalPurchases: 300,
		}
		
		rule := LoyaltyDiscountRule{
			Tier: "gold",
			MinOrderAmount: 500,
		}
		
		err := validator.ValidateLoyaltyDiscount(rule, customer)
		if err == nil {
			t.Error("Expected error for insufficient purchases")
		}
	})
}

func TestValidateCategoryDiscount(t *testing.T) {
	validator := NewDiscountValidator()
	now := time.Now()
	
	t.Run("ValidCategoryDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 3, Category: "electronics"},
		}
		
		rule := CategoryDiscountRule{
			Category: "electronics",
			MinQuantity: 2,
			ValidFrom: now.Add(-time.Hour),
			ValidUntil: now.Add(time.Hour),
		}
		
		err := validator.ValidateCategoryDiscount(rule, items)
		if err != nil {
			t.Errorf("Expected valid category discount, got error: %v", err)
		}
	})
	
	t.Run("DiscountNotYetValid", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 3, Category: "electronics"},
		}
		
		rule := CategoryDiscountRule{
			Category: "electronics",
			ValidFrom: now.Add(time.Hour),
			ValidUntil: now.Add(2 * time.Hour),
		}
		
		err := validator.ValidateCategoryDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for discount not yet valid")
		}
	})
	
	t.Run("DiscountExpired", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 3, Category: "electronics"},
		}
		
		rule := CategoryDiscountRule{
			Category: "electronics",
			ValidFrom: now.Add(-2 * time.Hour),
			ValidUntil: now.Add(-time.Hour),
		}
		
		err := validator.ValidateCategoryDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for expired discount")
		}
	})
	
	t.Run("NoCategoryItems", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 3, Category: "books"},
		}
		
		rule := CategoryDiscountRule{
			Category: "electronics",
			ValidFrom: now.Add(-time.Hour),
			ValidUntil: now.Add(time.Hour),
		}
		
		err := validator.ValidateCategoryDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for no category items")
		}
	})
	
	t.Run("InsufficientCategoryQuantity", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
		}
		
		rule := CategoryDiscountRule{
			Category: "electronics",
			MinQuantity: 3,
			ValidFrom: now.Add(-time.Hour),
			ValidUntil: now.Add(time.Hour),
		}
		
		err := validator.ValidateCategoryDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for insufficient category quantity")
		}
	})
}

func TestValidateProgressiveDiscount(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidProgressiveDiscount", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 15, Category: "electronics"},
		}
		
		rule := ProgressiveDiscountRule{
			QuantityStep: 10,
			Category: "electronics",
		}
		
		err := validator.ValidateProgressiveDiscount(rule, items)
		if err != nil {
			t.Errorf("Expected valid progressive discount, got error: %v", err)
		}
	})
	
	t.Run("InsufficientQuantityForStep", func(t *testing.T) {
		items := []DiscountItem{
			{ID: "item1", Price: 100, Quantity: 5, Category: "electronics"},
		}
		
		rule := ProgressiveDiscountRule{
			QuantityStep: 10,
			Category: "electronics",
		}
		
		err := validator.ValidateProgressiveDiscount(rule, items)
		if err == nil {
			t.Error("Expected error for insufficient quantity for progressive step")
		}
	})
}

func TestValidateCustomerEligibility(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidCustomer", func(t *testing.T) {
		customer := Customer{
			ID: "customer1",
			LoyaltyTier: "gold",
		}
		
		err := validator.ValidateCustomerEligibility(customer, DiscountTypeLoyalty)
		if err != nil {
			t.Errorf("Expected valid customer, got error: %v", err)
		}
	})
	
	t.Run("MissingCustomerID", func(t *testing.T) {
		customer := Customer{
			LoyaltyTier: "gold",
		}
		
		err := validator.ValidateCustomerEligibility(customer, DiscountTypeLoyalty)
		if err == nil {
			t.Error("Expected error for missing customer ID")
		}
	})
	
	t.Run("MissingLoyaltyTierForLoyaltyDiscount", func(t *testing.T) {
		customer := Customer{
			ID: "customer1",
		}
		
		err := validator.ValidateCustomerEligibility(customer, DiscountTypeLoyalty)
		if err == nil {
			t.Error("Expected error for missing loyalty tier")
		}
	})
}

func TestValidateDiscountLimits(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("WithinUsageLimit", func(t *testing.T) {
		err := validator.ValidateDiscountLimits("rule1", Customer{}, 3, 5)
		if err != nil {
			t.Errorf("Expected valid usage, got error: %v", err)
		}
	})
	
	t.Run("ExceedsUsageLimit", func(t *testing.T) {
		err := validator.ValidateDiscountLimits("rule1", Customer{}, 5, 5)
		if err == nil {
			t.Error("Expected error for exceeding usage limit")
		}
	})
	
	t.Run("NoUsageLimit", func(t *testing.T) {
		err := validator.ValidateDiscountLimits("rule1", Customer{}, 100, 0)
		if err != nil {
			t.Errorf("Expected no error for unlimited usage, got: %v", err)
		}
	})
}

func TestValidateTimeConstraints(t *testing.T) {
	validator := NewDiscountValidator()
	now := time.Now()
	
	t.Run("ValidTimeRange", func(t *testing.T) {
		err := validator.ValidateTimeConstraints(
			now.Add(-time.Hour),
			now.Add(time.Hour),
		)
		if err != nil {
			t.Errorf("Expected valid time range, got error: %v", err)
		}
	})
	
	t.Run("NotYetValid", func(t *testing.T) {
		err := validator.ValidateTimeConstraints(
			now.Add(time.Hour),
			now.Add(2*time.Hour),
		)
		if err == nil {
			t.Error("Expected error for not yet valid time range")
		}
	})
	
	t.Run("Expired", func(t *testing.T) {
		err := validator.ValidateTimeConstraints(
			now.Add(-2*time.Hour),
			now.Add(-time.Hour),
		)
		if err == nil {
			t.Error("Expected error for expired time range")
		}
	})
}

func TestValidatorConfiguration(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("SetMaxStackedDiscountPercent", func(t *testing.T) {
		validator.SetMaxStackedDiscountPercent(60.0)
		if validator.MaxStackedDiscountPercent != 60.0 {
			t.Errorf("Expected max stacked discount 60.0, got %f", validator.MaxStackedDiscountPercent)
		}
	})
	
	t.Run("SetMaxSingleDiscountPercent", func(t *testing.T) {
		validator.SetMaxSingleDiscountPercent(40.0)
		if validator.MaxSingleDiscountPercent != 40.0 {
			t.Errorf("Expected max single discount 40.0, got %f", validator.MaxSingleDiscountPercent)
		}
	})
	
	t.Run("AddAllowedCombination", func(t *testing.T) {
		validator.AddAllowedCombination(DiscountTypeProgressive, DiscountTypeBulk)
		
		if !validator.canCombineDiscounts(DiscountTypeProgressive, DiscountTypeBulk) {
			t.Error("Expected progressive and bulk discounts to be combinable after adding")
		}
	})
	
	t.Run("RemoveAllowedCombination", func(t *testing.T) {
		validator.RemoveAllowedCombination(DiscountTypeBulk, DiscountTypeLoyalty)
		
		if validator.canCombineDiscounts(DiscountTypeBulk, DiscountTypeLoyalty) {
			t.Error("Expected bulk and loyalty discounts to not be combinable after removal")
		}
	})
}

func TestValidateDiscountResult(t *testing.T) {
	validator := NewDiscountValidator()
	
	t.Run("ValidResult", func(t *testing.T) {
		result := DiscountCalculationResult{
			OriginalAmount: 100.0,
			TotalDiscount: 20.0,
			FinalAmount: 80.0,
			SavingsPercent: 20.0,
			IsValid: true,
			AppliedDiscounts: []DiscountApplication{
				{
					Type: DiscountTypeBulk,
					DiscountAmount: 20.0,
					AppliedItems: []DiscountItem{
						{ID: "item1", Price: 100, Quantity: 1},
					},
				},
			},
		}
		
		err := validator.ValidateDiscountResult(result)
		if err != nil {
			t.Errorf("Expected valid result, got error: %v", err)
		}
	})
	
	t.Run("NegativeTotalDiscount", func(t *testing.T) {
		result := DiscountCalculationResult{
			TotalDiscount: -10.0,
		}
		
		err := validator.ValidateDiscountResult(result)
		if err == nil {
			t.Error("Expected error for negative total discount")
		}
	})
	
	t.Run("NegativeFinalAmount", func(t *testing.T) {
		result := DiscountCalculationResult{
			FinalAmount: -10.0,
		}
		
		err := validator.ValidateDiscountResult(result)
		if err == nil {
			t.Error("Expected error for negative final amount")
		}
	})
	
	t.Run("DiscountExceedsOriginal", func(t *testing.T) {
		result := DiscountCalculationResult{
			OriginalAmount: 100.0,
			TotalDiscount: 150.0,
		}
		
		err := validator.ValidateDiscountResult(result)
		if err == nil {
			t.Error("Expected error for discount exceeding original amount")
		}
	})
	
	t.Run("SavingsPercentExceeds100", func(t *testing.T) {
		result := DiscountCalculationResult{
			OriginalAmount: 100.0,
			SavingsPercent: 150.0,
		}
		
		err := validator.ValidateDiscountResult(result)
		if err == nil {
			t.Error("Expected error for savings percent exceeding 100%")
		}
	})
}

func BenchmarkValidateDiscountApplication(t *testing.B) {
	validator := NewDiscountValidator()
	items := []DiscountItem{
		{ID: "item1", Price: 100, Quantity: 1, Category: "electronics"},
	}
	
	discount := DiscountApplication{
		Type: DiscountTypeBulk,
		DiscountAmount: 10,
		AppliedItems: items,
	}
	
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		validator.ValidateDiscountApplication(discount, items, Customer{})
	}
}

func BenchmarkValidateStackedDiscounts(t *testing.B) {
	validator := NewDiscountValidator()
	discounts := []DiscountApplication{
		{
			Type: DiscountTypeBulk,
			DiscountAmount: 10,
		},
		{
			Type: DiscountTypeLoyalty,
			DiscountAmount: 5,
		},
	}
	
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		validator.ValidateStackedDiscounts(discounts, 100.0)
	}
}
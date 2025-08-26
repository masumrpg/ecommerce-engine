package discount

import (
	"errors"
	"fmt"
	"time"
)

// DiscountValidator handles validation of discount applications
type DiscountValidator struct {
	MaxStackedDiscountPercent float64
	MaxSingleDiscountPercent  float64
	AllowedCombinations       map[DiscountType][]DiscountType
}

// NewDiscountValidator creates a new discount validator
func NewDiscountValidator() *DiscountValidator {
	return &DiscountValidator{
		MaxStackedDiscountPercent: 50.0, // Default max 50% total discount
		MaxSingleDiscountPercent:  30.0, // Default max 30% single discount
		AllowedCombinations: map[DiscountType][]DiscountType{
			DiscountTypeBulk:        {DiscountTypeLoyalty, DiscountTypeCategory},
			DiscountTypeBundle:      {DiscountTypeLoyalty},
			DiscountTypeLoyalty:     {DiscountTypeBulk, DiscountTypeBundle, DiscountTypeCategory, DiscountTypeTier},
			DiscountTypeCategory:    {DiscountTypeBulk, DiscountTypeLoyalty},
			DiscountTypeTier:        {DiscountTypeLoyalty},
			DiscountTypeProgressive: {},
		},
	}
}

// ValidateDiscountApplication validates if a discount can be applied
func (dv *DiscountValidator) ValidateDiscountApplication(discount DiscountApplication, items []DiscountItem, customer Customer) error {
	// Validate discount amount
	if discount.DiscountAmount < 0 {
		return errors.New("discount amount cannot be negative")
	}

	// Validate applied items
	if len(discount.AppliedItems) == 0 {
		return errors.New("discount must be applied to at least one item")
	}

	// Validate discount doesn't exceed item value
	itemsTotal := calculateItemsAmount(discount.AppliedItems)
	if discount.DiscountAmount > itemsTotal {
		return errors.New("discount amount cannot exceed item total")
	}

	// Validate single discount percentage limit
	if itemsTotal > 0 {
		discountPercent := (discount.DiscountAmount / itemsTotal) * 100
		if discountPercent > dv.MaxSingleDiscountPercent {
			return fmt.Errorf("single discount percentage (%.2f%%) exceeds maximum allowed (%.2f%%)",
				discountPercent, dv.MaxSingleDiscountPercent)
		}
	}

	return nil
}

// ValidateStackedDiscounts validates if multiple discounts can be stacked
func (dv *DiscountValidator) ValidateStackedDiscounts(discounts []DiscountApplication, originalAmount float64) error {
	if len(discounts) <= 1 {
		return nil // No stacking with single discount
	}

	// Calculate total discount amount
	totalDiscount := 0.0
	for _, discount := range discounts {
		totalDiscount += discount.DiscountAmount
	}

	// Validate total discount doesn't exceed original amount
	if totalDiscount > originalAmount {
		return errors.New("total discount cannot exceed original amount")
	}

	// Validate stacked discount percentage limit
	if originalAmount > 0 {
		totalDiscountPercent := (totalDiscount / originalAmount) * 100
		if totalDiscountPercent > dv.MaxStackedDiscountPercent {
			return fmt.Errorf("stacked discount percentage (%.2f%%) exceeds maximum allowed (%.2f%%)",
				totalDiscountPercent, dv.MaxStackedDiscountPercent)
		}
	}

	// Validate discount type combinations
	for i, discount1 := range discounts {
		for j, discount2 := range discounts {
			if i != j {
				if !dv.canCombineDiscounts(discount1.Type, discount2.Type) {
					return fmt.Errorf("cannot combine %s discount with %s discount",
						discount1.Type, discount2.Type)
				}
			}
		}
	}

	return nil
}

// ValidateBulkDiscount validates bulk discount rules
func (dv *DiscountValidator) ValidateBulkDiscount(rule BulkDiscountRule, items []DiscountItem) error {
	applicableItems := getApplicableItems(items, rule.ApplicableCategories, rule.ApplicableProducts)
	totalQuantity := getTotalQuantity(applicableItems)

	if totalQuantity < rule.MinQuantity {
		return fmt.Errorf("minimum quantity requirement not met: need %d, have %d",
			rule.MinQuantity, totalQuantity)
	}

	if rule.MaxQuantity > 0 && totalQuantity > rule.MaxQuantity {
		return fmt.Errorf("maximum quantity exceeded: limit %d, have %d",
			rule.MaxQuantity, totalQuantity)
	}

	return nil
}

// ValidateTierPricing validates tier pricing rules
func (dv *DiscountValidator) ValidateTierPricing(rule TierPricingRule, items []DiscountItem) error {
	var applicableItems []DiscountItem
	if rule.Category != "" {
		applicableItems = getItemsByCategory(items, rule.Category)
	} else {
		applicableItems = items
	}
	totalQuantity := getTotalQuantity(applicableItems)

	if totalQuantity < rule.MinQuantity {
		return fmt.Errorf("minimum quantity requirement not met for tier pricing: need %d, have %d",
			rule.MinQuantity, totalQuantity)
	}

	return nil
}

// ValidateBundleDiscount validates bundle discount rules
func (dv *DiscountValidator) ValidateBundleDiscount(rule BundleDiscountRule, items []DiscountItem) error {
	// Check if bundle requirements are met
	bundleMatches := findBundleMatches(items, rule)
	if len(bundleMatches) == 0 {
		return errors.New("bundle requirements not met")
	}

	// Validate minimum items requirement
	for _, match := range bundleMatches {
		if len(match.MatchedItems) < rule.MinItems {
			return fmt.Errorf("bundle minimum items requirement not met: need %d, have %d",
				rule.MinItems, len(match.MatchedItems))
		}
	}

	return nil
}

// ValidateLoyaltyDiscount validates loyalty discount rules
func (dv *DiscountValidator) ValidateLoyaltyDiscount(rule LoyaltyDiscountRule, customer Customer) error {
	if customer.LoyaltyTier != rule.Tier {
		return fmt.Errorf("customer loyalty tier '%s' does not match required tier '%s'",
			customer.LoyaltyTier, rule.Tier)
	}

	if rule.MinOrderAmount > 0 && customer.TotalPurchases < rule.MinOrderAmount {
		return fmt.Errorf("minimum order amount requirement not met: need %.2f, have %.2f",
			rule.MinOrderAmount, customer.TotalPurchases)
	}

	return nil
}

// ValidateCategoryDiscount validates category discount rules
func (dv *DiscountValidator) ValidateCategoryDiscount(rule CategoryDiscountRule, items []DiscountItem) error {
	now := time.Now()

	// Check if discount is currently valid
	if now.Before(rule.ValidFrom) {
		return fmt.Errorf("category discount not yet valid: starts %s", rule.ValidFrom.Format("2006-01-02"))
	}

	if now.After(rule.ValidUntil) {
		return fmt.Errorf("category discount expired: ended %s", rule.ValidUntil.Format("2006-01-02"))
	}

	// Check if any items match the category
	categoryItems := getItemsByCategory(items, rule.Category)
	if len(categoryItems) == 0 {
		return fmt.Errorf("no items found in category '%s'", rule.Category)
	}

	// Check minimum quantity requirement
	if rule.MinQuantity > 0 {
		totalQuantity := getTotalQuantity(categoryItems)
		if totalQuantity < rule.MinQuantity {
			return fmt.Errorf("minimum quantity requirement not met for category '%s': need %d, have %d",
				rule.Category, rule.MinQuantity, totalQuantity)
		}
	}

	return nil
}

// ValidateProgressiveDiscount validates progressive discount rules
func (dv *DiscountValidator) ValidateProgressiveDiscount(rule ProgressiveDiscountRule, items []DiscountItem) error {
	var applicableItems []DiscountItem
	if rule.Category != "" {
		applicableItems = getItemsByCategory(items, rule.Category)
	} else {
		applicableItems = items
	}

	totalQuantity := getTotalQuantity(applicableItems)
	if totalQuantity < rule.QuantityStep {
		return fmt.Errorf("minimum quantity requirement not met for progressive discount: need %d, have %d",
			rule.QuantityStep, totalQuantity)
	}

	return nil
}

// ValidateCustomerEligibility validates if customer is eligible for discounts
func (dv *DiscountValidator) ValidateCustomerEligibility(customer Customer, discountType DiscountType) error {
	// Check if customer ID is provided
	if customer.ID == "" {
		return errors.New("customer ID is required")
	}

	// Check loyalty tier requirements for loyalty discounts
	if discountType == DiscountTypeLoyalty && customer.LoyaltyTier == "" {
		return errors.New("customer must have a loyalty tier for loyalty discounts")
	}

	// Additional validation can be added here based on business requirements

	return nil
}

// ValidateDiscountLimits validates discount usage limits
func (dv *DiscountValidator) ValidateDiscountLimits(ruleID string, customer Customer, usageCount int, maxUsage int) error {
	if maxUsage > 0 && usageCount >= maxUsage {
		return fmt.Errorf("discount usage limit reached for rule '%s': %d/%d", ruleID, usageCount, maxUsage)
	}

	return nil
}

// ValidateTimeConstraints validates time-based constraints
func (dv *DiscountValidator) ValidateTimeConstraints(validFrom, validUntil time.Time) error {
	now := time.Now()

	if now.Before(validFrom) {
		return fmt.Errorf("discount not yet valid: starts %s", validFrom.Format("2006-01-02 15:04:05"))
	}

	if now.After(validUntil) {
		return fmt.Errorf("discount expired: ended %s", validUntil.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// Helper methods

// canCombineDiscounts checks if two discount types can be combined
func (dv *DiscountValidator) canCombineDiscounts(type1, type2 DiscountType) bool {
	allowedTypes, exists := dv.AllowedCombinations[type1]
	if !exists {
		return false
	}

	for _, allowedType := range allowedTypes {
		if allowedType == type2 {
			return true
		}
	}

	return false
}

// SetMaxStackedDiscountPercent sets the maximum allowed stacked discount percentage
func (dv *DiscountValidator) SetMaxStackedDiscountPercent(percent float64) {
	dv.MaxStackedDiscountPercent = percent
}

// SetMaxSingleDiscountPercent sets the maximum allowed single discount percentage
func (dv *DiscountValidator) SetMaxSingleDiscountPercent(percent float64) {
	dv.MaxSingleDiscountPercent = percent
}

// AddAllowedCombination adds an allowed discount type combination
func (dv *DiscountValidator) AddAllowedCombination(baseType DiscountType, allowedType DiscountType) {
	if dv.AllowedCombinations[baseType] == nil {
		dv.AllowedCombinations[baseType] = []DiscountType{}
	}
	dv.AllowedCombinations[baseType] = append(dv.AllowedCombinations[baseType], allowedType)
}

// RemoveAllowedCombination removes an allowed discount type combination
func (dv *DiscountValidator) RemoveAllowedCombination(baseType DiscountType, disallowedType DiscountType) {
	allowedTypes, exists := dv.AllowedCombinations[baseType]
	if !exists {
		return
	}

	for i, allowedType := range allowedTypes {
		if allowedType == disallowedType {
			dv.AllowedCombinations[baseType] = append(allowedTypes[:i], allowedTypes[i+1:]...)
			break
		}
	}
}

// ValidateDiscountResult validates the final discount calculation result
func (dv *DiscountValidator) ValidateDiscountResult(result DiscountCalculationResult) error {
	if result.TotalDiscount < 0 {
		return errors.New("total discount cannot be negative")
	}

	if result.FinalAmount < 0 {
		return errors.New("final amount cannot be negative")
	}

	if result.TotalDiscount > result.OriginalAmount {
		return errors.New("total discount cannot exceed original amount")
	}

	if result.OriginalAmount > 0 && result.SavingsPercent > 100 {
		return errors.New("savings percentage cannot exceed 100%")
	}

	// Validate individual discount applications
	for _, discount := range result.AppliedDiscounts {
		if err := dv.ValidateDiscountApplication(discount, nil, Customer{}); err != nil {
			return fmt.Errorf("invalid discount application '%s': %v", discount.Name, err)
		}
	}

	return nil
}
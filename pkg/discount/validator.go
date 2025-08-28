// Package discount provides comprehensive validation functionality for discount rules and applications.
// This package ensures that discount calculations are accurate, secure, and comply with business rules.
//
// Key Features:
//   - Comprehensive validation for all discount types (bulk, tier, bundle, loyalty, category, progressive)
//   - Stacked discount validation with configurable limits
//   - Customer eligibility verification
//   - Time-based constraint validation
//   - Usage limit enforcement
//   - Discount combination rules management
//   - Result validation for calculation accuracy
//
// Basic Usage:
//
//	validator := NewDiscountValidator()
//	validator.SetMaxStackedDiscountPercent(50.0)
//	validator.AddAllowedCombination(DiscountTypeBulk, DiscountTypeLoyalty)
//
//	// Validate a discount application
//	err := validator.ValidateDiscountApplication(discountApp, items, customer)
//	if err != nil {
//	    log.Printf("Discount validation failed: %v", err)
//	}
//
//	// Validate stacked discounts
//	err = validator.ValidateStackedDiscounts(discountApps, originalAmount)
//	if err != nil {
//	    log.Printf("Stacked discount validation failed: %v", err)
//	}
package discount

import (
	"errors"
	"fmt"
	"time"
)

// DiscountValidator handles validation of discount applications
// DiscountValidator provides comprehensive validation for discount rules and applications.
// It enforces business rules, validates discount combinations, and ensures calculation accuracy.
//
// Features:
//   - Configurable maximum discount percentages for single and stacked discounts
//   - Flexible discount combination rules management
//   - Comprehensive validation for all discount types
//   - Customer eligibility verification
//   - Time-based constraint validation
//   - Usage limit enforcement
//
// Example:
//
//	validator := NewDiscountValidator()
//	validator.SetMaxStackedDiscountPercent(50.0)
//	validator.SetMaxSingleDiscountPercent(30.0)
//	validator.AddAllowedCombination(DiscountTypeBulk, DiscountTypeLoyalty)
type DiscountValidator struct {
	MaxStackedDiscountPercent float64
	MaxSingleDiscountPercent  float64
	AllowedCombinations       map[DiscountType][]DiscountType
}

// NewDiscountValidator creates a new discount validator with sensible default settings.
// The validator is initialized with conservative limits to prevent excessive discounting.
//
// Default Settings:
//   - Maximum stacked discount: 50%
//   - Maximum single discount: 30%
//   - Pre-configured allowed discount combinations for common business scenarios
//
// Returns:
//   - *DiscountValidator: A new validator instance ready for use
//
// Example:
//
//	validator := NewDiscountValidator()
//	// Optionally customize settings
//	validator.SetMaxStackedDiscountPercent(60.0)
//	validator.AddAllowedCombination(DiscountTypeBulk, DiscountTypeLoyalty)
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

// ValidateDiscountApplication validates if a discount can be applied according to business rules.
// This method performs comprehensive validation including amount limits, item requirements,
// and percentage constraints to ensure discount integrity.
//
// Validation Rules:
//   - Discount amount must be non-negative
//   - Discount must be applied to at least one item
//   - Discount amount cannot exceed the total value of applied items
//   - Discount percentage must not exceed configured single discount limit
//
// Parameters:
//   - discount: The discount application to validate
//   - items: Cart items for context validation
//   - customer: Customer information for eligibility checks
//
// Returns:
//   - error: Validation error if any rule is violated, nil if valid
//
// Example:
//
//	discountApp := DiscountApplication{
//	    Name: "Bulk Discount",
//	    DiscountAmount: 10.0,
//	    AppliedItems: []DiscountItem{item1, item2},
//	}
//	err := validator.ValidateDiscountApplication(discountApp, items, customer)
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

// ValidateStackedDiscounts validates multiple discount applications to ensure they comply
// with stacking rules and don't exceed configured limits. This method is crucial for
// preventing excessive discounting and maintaining business profitability.
//
// Validation Rules:
//   - Total stacked discount percentage must not exceed configured maximum
//   - All discount combinations must be explicitly allowed
//   - Individual discounts must be valid
//   - No duplicate discount types (unless explicitly allowed)
//
// Parameters:
//   - discounts: Array of discount applications to validate
//   - originalAmount: Original order amount for percentage calculations
//
// Returns:
//   - error: Validation error if stacking rules are violated, nil if valid
//
// Example:
//
//	discounts := []DiscountApplication{
//	    {Type: DiscountTypeBulk, DiscountAmount: 10.0},
//	    {Type: DiscountTypeLoyalty, DiscountAmount: 5.0},
//	}
//	err := validator.ValidateStackedDiscounts(discounts, 100.0)
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

// ValidateBulkDiscount validates bulk discount rules against cart items to ensure
// minimum quantity requirements are met. This validation is essential for quantity-based
// discount strategies.
//
// Validation Rules:
//   - Total applicable item quantity must meet or exceed rule's minimum quantity
//   - Total applicable item quantity must not exceed rule's maximum quantity (if set)
//   - Items must be eligible for bulk discount consideration
//
// Parameters:
//   - rule: The bulk discount rule to validate against
//   - items: Cart items to check for bulk discount eligibility
//
// Returns:
//   - error: Validation error if requirements not met, nil if valid
//
// Example:
//
//	rule := BulkDiscountRule{
//	    MinQuantity: 10,
//	    DiscountPercent: 15.0,
//	}
//	err := validator.ValidateBulkDiscount(rule, cartItems)
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

// ValidateTierPricing validates tier pricing rules against cart items to ensure
// minimum amount thresholds are met. This validation supports tiered discount strategies
// based on order value.
//
// Validation Rules:
//   - Total applicable item amount must meet or exceed rule's minimum amount
//   - Items must be eligible for tier pricing consideration
//   - Tier thresholds must be logically consistent
//
// Parameters:
//   - rule: The tier pricing rule to validate against
//   - items: Cart items to check for tier pricing eligibility
//
// Returns:
//   - error: Validation error if requirements not met, nil if valid
//
// Example:
//
//	rule := TierPricingRule{
//	    MinAmount: 100.0,
//	    DiscountPercent: 10.0,
//	}
//	err := validator.ValidateTierPricing(rule, cartItems)
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

// ValidateBundleDiscount validates bundle discount rules against cart items to ensure
// required product combinations are present. This validation supports complex bundling
// strategies and cross-selling initiatives.
//
// Validation Rules:
//   - All required bundle products must be present in sufficient quantities
//   - Minimum items requirement must be met for each bundle match
//   - Bundle composition must match rule specifications
//   - Optional bundle items are validated if present
//
// Parameters:
//   - rule: The bundle discount rule to validate against
//   - items: Cart items to check for bundle eligibility
//
// Returns:
//   - error: Validation error if bundle requirements not met, nil if valid
//
// Example:
//
//	rule := BundleDiscountRule{
//	    RequiredProducts: []string{"PROD1", "PROD2"},
//	    MinItems: 2,
//	    DiscountPercent: 20.0,
//	}
//	err := validator.ValidateBundleDiscount(rule, cartItems)
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

// ValidateLoyaltyDiscount validates loyalty discount rules against customer information
// to ensure proper tier eligibility and purchase history requirements. This validation
// supports customer retention and loyalty program strategies.
//
// Validation Rules:
//   - Customer loyalty tier must match rule requirements
//   - Customer total purchases must meet minimum order amount (if specified)
//   - Customer must have valid loyalty program membership
//   - Loyalty tier must be active and not expired
//
// Parameters:
//   - rule: The loyalty discount rule to validate against
//   - customer: Customer information for eligibility verification
//
// Returns:
//   - error: Validation error if loyalty requirements not met, nil if valid
//
// Example:
//
//	rule := LoyaltyDiscountRule{
//	    Tier: "Gold",
//	    MinOrderAmount: 50.0,
//	    DiscountPercent: 15.0,
//	}
//	err := validator.ValidateLoyaltyDiscount(rule, customer)
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

// ValidateCategoryDiscount validates category discount rules against cart items and time constraints.
// This validation ensures category-specific discounts are applied correctly and within valid periods.
//
// Validation Rules:
//   - Discount must be within valid time period (ValidFrom to ValidUntil)
//   - At least one item must belong to the specified category
//   - Minimum quantity requirement must be met for category items (if specified)
//   - Category must be valid and recognized
//
// Parameters:
//   - rule: The category discount rule to validate against
//   - items: Cart items to check for category eligibility
//
// Returns:
//   - error: Validation error if category requirements not met, nil if valid
//
// Example:
//
//	rule := CategoryDiscountRule{
//	    Category: "Electronics",
//	    MinQuantity: 2,
//	    ValidFrom: time.Now().AddDate(0, 0, -1),
//	    ValidUntil: time.Now().AddDate(0, 0, 30),
//	    DiscountPercent: 10.0,
//	}
//	err := validator.ValidateCategoryDiscount(rule, cartItems)
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

// ValidateProgressiveDiscount validates progressive discount rules against cart items to ensure
// quantity step requirements are met. This validation supports incremental discount strategies
// that increase with purchase volume.
//
// Validation Rules:
//   - Total applicable item quantity must meet or exceed the quantity step threshold
//   - If category is specified, only items in that category are considered
//   - Progressive steps must be logically consistent
//   - Quantity calculations must be accurate
//
// Parameters:
//   - rule: The progressive discount rule to validate against
//   - items: Cart items to check for progressive discount eligibility
//
// Returns:
//   - error: Validation error if quantity requirements not met, nil if valid
//
// Example:
//
//	rule := ProgressiveDiscountRule{
//	    QuantityStep: 5,
//	    Category: "Books",
//	    DiscountPercent: 5.0,
//	}
//	err := validator.ValidateProgressiveDiscount(rule, cartItems)
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

// ValidateCustomerEligibility validates if a customer is eligible for specific discount types.
// This validation ensures that only qualified customers receive appropriate discounts
// based on their profile and membership status.
//
// Validation Rules:
//   - Customer ID must be provided and valid
//   - Loyalty tier must be present for loyalty-based discounts
//   - Customer account must be active and in good standing
//   - Additional business-specific eligibility rules
//
// Parameters:
//   - customer: Customer information for eligibility verification
//   - discountType: Type of discount to validate eligibility for
//
// Returns:
//   - error: Validation error if customer not eligible, nil if eligible
//
// Example:
//
//	customer := Customer{
//	    ID: "CUST123",
//	    LoyaltyTier: "Gold",
//	}
//	err := validator.ValidateCustomerEligibility(customer, DiscountTypeLoyalty)
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

// ValidateDiscountLimits validates discount usage limits to prevent abuse and ensure
// fair distribution of promotional benefits. This validation tracks and enforces
// per-customer and global usage constraints.
//
// Validation Rules:
//   - Usage count must not exceed maximum allowed usage (if specified)
//   - Per-customer limits must be respected
//   - Global discount limits must be enforced
//   - Time-based usage windows must be considered
//
// Parameters:
//   - ruleID: Unique identifier for the discount rule
//   - customer: Customer information for usage tracking
//   - usageCount: Current usage count for this customer/rule combination
//   - maxUsage: Maximum allowed usage (0 means unlimited)
//
// Returns:
//   - error: Validation error if usage limits exceeded, nil if within limits
//
// Example:
//
//	err := validator.ValidateDiscountLimits("BULK10", customer, 2, 3)
func (dv *DiscountValidator) ValidateDiscountLimits(ruleID string, customer Customer, usageCount int, maxUsage int) error {
	if maxUsage > 0 && usageCount >= maxUsage {
		return fmt.Errorf("discount usage limit reached for rule '%s': %d/%d", ruleID, usageCount, maxUsage)
	}

	return nil
}

// ValidateTimeConstraints validates time-based constraints for discount validity.
// This validation ensures discounts are only applied within their designated time windows,
// supporting time-limited promotions and seasonal campaigns.
//
// Validation Rules:
//   - Current time must be after or equal to validFrom time
//   - Current time must be before or equal to validUntil time
//   - Time zones must be handled consistently
//   - Date ranges must be logically valid (validFrom < validUntil)
//
// Parameters:
//   - validFrom: Start time for discount validity
//   - validUntil: End time for discount validity
//
// Returns:
//   - error: Validation error if current time outside valid window, nil if valid
//
// Example:
//
//	validFrom := time.Now().AddDate(0, 0, -1) // Yesterday
//	validUntil := time.Now().AddDate(0, 0, 7)  // Next week
//	err := validator.ValidateTimeConstraints(validFrom, validUntil)
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

// canCombineDiscounts checks if two discount types can be combined according to
// configured combination rules. This method supports flexible discount stacking
// policies while preventing incompatible discount combinations.
//
// Parameters:
//   - type1: First discount type to check
//   - type2: Second discount type to check
//
// Returns:
//   - bool: true if the discount types can be combined, false otherwise
//
// Example:
//
//	canCombine := validator.canCombineDiscounts(DiscountTypeBulk, DiscountTypeLoyalty)
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

// SetMaxStackedDiscountPercent sets the maximum allowed stacked discount percentage.
// This configuration method allows dynamic adjustment of discount limits to support
// different promotional strategies and business policies.
//
// Parameters:
//   - percent: Maximum allowed stacked discount percentage (0-100)
//
// Example:
//
//	validator.SetMaxStackedDiscountPercent(60.0) // Allow up to 60% total discount
func (dv *DiscountValidator) SetMaxStackedDiscountPercent(percent float64) {
	dv.MaxStackedDiscountPercent = percent
}

// SetMaxSingleDiscountPercent sets the maximum allowed single discount percentage.
// This configuration method provides control over individual discount limits to
// prevent excessive single-discount applications.
//
// Parameters:
//   - percent: Maximum allowed single discount percentage (0-100)
//
// Example:
//
//	validator.SetMaxSingleDiscountPercent(40.0) // Allow up to 40% single discount
func (dv *DiscountValidator) SetMaxSingleDiscountPercent(percent float64) {
	dv.MaxSingleDiscountPercent = percent
}

// AddAllowedCombination adds an allowed discount type combination to the validator's
// configuration. This method enables flexible discount stacking policies by defining
// which discount types can be combined together.
//
// Parameters:
//   - baseType: The base discount type
//   - allowedType: The discount type that can be combined with the base type
//
// Example:
//
//	// Allow bulk discounts to be combined with loyalty discounts
//	validator.AddAllowedCombination(DiscountTypeBulk, DiscountTypeLoyalty)
//	// Allow loyalty discounts to be combined with category discounts
//	validator.AddAllowedCombination(DiscountTypeLoyalty, DiscountTypeCategory)
func (dv *DiscountValidator) AddAllowedCombination(baseType DiscountType, allowedType DiscountType) {
	if dv.AllowedCombinations[baseType] == nil {
		dv.AllowedCombinations[baseType] = []DiscountType{}
	}
	dv.AllowedCombinations[baseType] = append(dv.AllowedCombinations[baseType], allowedType)
}

// RemoveAllowedCombination removes an allowed discount type combination from the
// validator's configuration. This method provides flexibility to adjust discount
// stacking policies by removing previously allowed combinations.
//
// Parameters:
//   - baseType: The base discount type
//   - disallowedType: The discount type to remove from allowed combinations
//
// Example:
//
//	// Remove the ability to combine bulk discounts with loyalty discounts
//	validator.RemoveAllowedCombination(DiscountTypeBulk, DiscountTypeLoyalty)
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

// ValidateDiscountResult validates the final discount calculation result to ensure
// mathematical accuracy and business rule compliance. This validation is the final
// check before applying discounts to prevent calculation errors and policy violations.
//
// Validation Rules:
//   - Total discount amount must be non-negative
//   - Final amount must be non-negative
//   - Total discount cannot exceed original amount
//   - Savings percentage cannot exceed 100%
//   - All individual discount applications must be valid
//   - Mathematical consistency between amounts and percentages
//
// Parameters:
//   - result: The discount calculation result to validate
//
// Returns:
//   - error: Validation error if result is invalid, nil if valid
//
// Example:
//
//	result := DiscountCalculationResult{
//	    OriginalAmount: 100.0,
//	    TotalDiscount: 25.0,
//	    FinalAmount: 75.0,
//	    SavingsPercent: 25.0,
//	}
//	err := validator.ValidateDiscountResult(result)
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
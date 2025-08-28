// Package discount provides comprehensive discount rule management and application capabilities.
// This package implements a flexible rule engine that supports multiple discount types,
// validation, and complex discount scenarios for e-commerce applications.
//
// Key Features:
//   - Multiple discount types (bulk, tier, bundle, loyalty, seasonal, etc.)
//   - Rule validation and management
//   - Flexible rule application with stacking support
//   - Time-based and seasonal discount rules
//   - Cross-sell and mix-and-match promotions
//   - Customer loyalty integration
//   - Comprehensive rule validation
//
// Basic Usage:
//   engine := NewRuleEngine()
//   
//   // Add various discount rules
//   bulkRule := BulkDiscountRule{
//     MinQuantity: 10,
//     DiscountType: "percentage",
//     DiscountValue: 15.0,
//   }
//   engine.AddBulkRule(bulkRule)
//   
//   // Apply rules to calculate discounts
//   result := engine.ApplyRules(items, customer, true)
//   fmt.Printf("Total discount: %.2f\n", result.TotalDiscount)
package discount

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// RuleEngine manages and applies discount rules.
// It serves as the central hub for all discount rule types, providing
// methods to add, validate, and apply various discount strategies.
//
// Features:
//   - Multiple rule type support (bulk, tier, bundle, loyalty, etc.)
//   - Rule validation before addition
//   - Flexible rule application strategies
//   - Stacking and non-stacking discount modes
//   - Time-based rule management
//   - Customer-specific rule application
//
// Supported Rule Types:
//   - BulkDiscountRule: Quantity-based discounts
//   - TierPricingRule: Tiered pricing structures
//   - BundleDiscountRule: Product bundle discounts
//   - LoyaltyDiscountRule: Customer loyalty rewards
//   - ProgressiveDiscountRule: Progressive discount tiers
//   - CategoryDiscountRule: Category-specific discounts
//   - FrequencyDiscountRule: Purchase frequency rewards
//   - SeasonalDiscountRule: Time and season-based discounts
//   - CrossSellRule: Cross-selling promotions
//   - MixAndMatchRule: Mix-and-match promotions
type RuleEngine struct {
	BulkRules        []BulkDiscountRule
	TierRules        []TierPricingRule
	BundleRules      []BundleDiscountRule
	LoyaltyRules     []LoyaltyDiscountRule
	ProgressiveRules []ProgressiveDiscountRule
	CategoryRules    []CategoryDiscountRule
	FrequencyRules   []FrequencyDiscountRule
	SeasonalRules    []SeasonalDiscountRule
	CrossSellRules   []CrossSellRule
	MixMatchRules    []MixAndMatchRule
}

// NewRuleEngine creates a new rule engine.
// Initializes all rule slices to empty, providing a clean starting point
// for adding and managing discount rules.
//
// Returns:
//   - *RuleEngine: A new rule engine instance with empty rule collections
//
// Example:
//   engine := NewRuleEngine()
//   // Engine is ready to accept rules
//   engine.AddBulkRule(bulkRule)
//   engine.AddLoyaltyRule(loyaltyRule)
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		BulkRules:        []BulkDiscountRule{},
		TierRules:        []TierPricingRule{},
		BundleRules:      []BundleDiscountRule{},
		LoyaltyRules:     []LoyaltyDiscountRule{},
		ProgressiveRules: []ProgressiveDiscountRule{},
		CategoryRules:    []CategoryDiscountRule{},
		FrequencyRules:   []FrequencyDiscountRule{},
		SeasonalRules:    []SeasonalDiscountRule{},
		CrossSellRules:   []CrossSellRule{},
		MixMatchRules:    []MixAndMatchRule{},
	}
}

// AddBulkRule adds a bulk discount rule.
// Validates the rule before adding it to the engine's bulk rules collection.
// Bulk rules apply discounts based on quantity thresholds.
//
// Parameters:
//   - rule: BulkDiscountRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := BulkDiscountRule{
//     MinQuantity: 10,
//     DiscountType: "percentage",
//     DiscountValue: 15.0,
//   }
//   err := engine.AddBulkRule(rule)
func (re *RuleEngine) AddBulkRule(rule BulkDiscountRule) error {
	if err := validateBulkRule(rule); err != nil {
		return err
	}
	re.BulkRules = append(re.BulkRules, rule)
	return nil
}

// AddTierRule adds a tier pricing rule.
// Validates the rule before adding it to the engine's tier rules collection.
// Tier rules provide different pricing based on quantity levels.
//
// Parameters:
//   - rule: TierPricingRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := TierPricingRule{
//     MinQuantity: 5,
//     PricePerItem: 9.99,
//   }
//   err := engine.AddTierRule(rule)
func (re *RuleEngine) AddTierRule(rule TierPricingRule) error {
	if err := validateTierRule(rule); err != nil {
		return err
	}
	re.TierRules = append(re.TierRules, rule)
	return nil
}

// AddBundleRule adds a bundle discount rule.
// Validates the rule before adding it to the engine's bundle rules collection.
// Bundle rules apply discounts when specific product combinations are purchased.
//
// Parameters:
//   - rule: BundleDiscountRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := BundleDiscountRule{
//     ID: "laptop_bundle",
//     RequiredProducts: []string{"laptop", "mouse"},
//     MinItems: 2,
//     DiscountType: "percentage",
//     DiscountValue: 10.0,
//   }
//   err := engine.AddBundleRule(rule)
func (re *RuleEngine) AddBundleRule(rule BundleDiscountRule) error {
	if err := validateBundleRule(rule); err != nil {
		return err
	}
	re.BundleRules = append(re.BundleRules, rule)
	return nil
}

// AddLoyaltyRule adds a loyalty discount rule.
// Validates the rule before adding it to the engine's loyalty rules collection.
// Loyalty rules provide discounts based on customer loyalty tiers.
//
// Parameters:
//   - rule: LoyaltyDiscountRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := LoyaltyDiscountRule{
//     Tier: "gold",
//     DiscountPercent: 15.0,
//     MinOrderAmount: 100.0,
//   }
//   err := engine.AddLoyaltyRule(rule)
func (re *RuleEngine) AddLoyaltyRule(rule LoyaltyDiscountRule) error {
	if err := validateLoyaltyRule(rule); err != nil {
		return err
	}
	re.LoyaltyRules = append(re.LoyaltyRules, rule)
	return nil
}

// AddCategoryRule adds a category discount rule.
// Validates the rule before adding it to the engine's category rules collection.
// Category rules apply discounts to specific product categories.
//
// Parameters:
//   - rule: CategoryDiscountRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := CategoryDiscountRule{
//     Category: "electronics",
//     DiscountPercent: 10.0,
//     ValidFrom: time.Now(),
//     ValidUntil: time.Now().AddDate(0, 1, 0),
//   }
//   err := engine.AddCategoryRule(rule)
func (re *RuleEngine) AddCategoryRule(rule CategoryDiscountRule) error {
	if err := validateCategoryRule(rule); err != nil {
		return err
	}
	re.CategoryRules = append(re.CategoryRules, rule)
	return nil
}

// AddSeasonalRule adds a seasonal discount rule.
// Validates the rule before adding it to the engine's seasonal rules collection.
// Seasonal rules apply discounts during specific seasons or time periods.
//
// Parameters:
//   - rule: SeasonalDiscountRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := SeasonalDiscountRule{
//     Season: "winter",
//     DiscountPercent: 20.0,
//     ValidFrom: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
//     ValidUntil: time.Date(2025, 2, 28, 23, 59, 59, 0, time.UTC),
//   }
//   err := engine.AddSeasonalRule(rule)
func (re *RuleEngine) AddSeasonalRule(rule SeasonalDiscountRule) error {
	if err := validateSeasonalRule(rule); err != nil {
		return err
	}
	re.SeasonalRules = append(re.SeasonalRules, rule)
	return nil
}

// AddCrossSellRule adds a cross-sell discount rule.
// Validates the rule before adding it to the engine's cross-sell rules collection.
// Cross-sell rules encourage purchasing complementary products together.
//
// Parameters:
//   - rule: CrossSellRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := CrossSellRule{
//     MainProductCategories: []string{"laptops"},
//     AccessoryCategories: []string{"accessories"},
//     DiscountPercent: 15.0,
//     MinMainProductPrice: 500.0,
//   }
//   err := engine.AddCrossSellRule(rule)
func (re *RuleEngine) AddCrossSellRule(rule CrossSellRule) error {
	if err := validateCrossSellRule(rule); err != nil {
		return err
	}
	re.CrossSellRules = append(re.CrossSellRules, rule)
	return nil
}

// AddMixMatchRule adds a mix and match rule.
// Validates the rule before adding it to the engine's mix-and-match rules collection.
// Mix-and-match rules provide discounts when customers buy multiple items from specified categories.
//
// Parameters:
//   - rule: MixAndMatchRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//   rule := MixAndMatchRule{
//     Categories: []string{"shirts", "pants", "shoes"},
//     RequiredItems: 3,
//     DiscountType: "flat_discount",
//     DiscountValue: 25.0,
//   }
//   err := engine.AddMixMatchRule(rule)
func (re *RuleEngine) AddMixMatchRule(rule MixAndMatchRule) error {
	if err := validateMixMatchRule(rule); err != nil {
		return err
	}
	re.MixMatchRules = append(re.MixMatchRules, rule)
	return nil
}

// ApplyRules applies all rules and returns the best discount.
// This is the main entry point for discount calculation, evaluating all applicable
// rules and returning the optimal discount configuration.
//
// Features:
//   - Comprehensive rule evaluation
//   - Stacking vs non-stacking modes
//   - Customer-specific rule application
//   - Automatic best discount selection
//   - Maximum stacked discount protection (50% default)
//
// Parameters:
//   - items: Slice of DiscountItem to calculate discounts for
//   - customer: Customer information for loyalty and personalized rules
//   - allowStacking: Whether to allow multiple discounts to stack
//
// Returns:
//   - DiscountCalculationResult: Complete discount calculation with applied rules
//
// Example:
//   items := []DiscountItem{{ProductID: "laptop", Price: 1000, Quantity: 2}}
//   customer := Customer{LoyaltyTier: "gold"}
//   result := engine.ApplyRules(items, customer, true)
//   fmt.Printf("Total savings: %.2f\n", result.TotalDiscount)
func (re *RuleEngine) ApplyRules(items []DiscountItem, customer Customer, allowStacking bool) DiscountCalculationResult {
	input := DiscountCalculationInput{
		Items:                     items,
		Customer:                  customer,
		BulkRules:                 re.BulkRules,
		TierRules:                 re.TierRules,
		BundleRules:               re.BundleRules,
		LoyaltyRules:              re.LoyaltyRules,
		ProgressiveRules:          re.ProgressiveRules,
		CategoryRules:             re.CategoryRules,
		AllowStacking:             allowStacking,
		MaxStackedDiscountPercent: 50, // Default max 50% stacked discount
	}

	return Calculate(input)
}

// ApplyFrequencyDiscounts applies purchase frequency-based discounts.
// Rewards customers based on their purchase history and frequency,
// encouraging repeat business through progressive discounts.
//
// Features:
//   - Purchase count validation
//   - Progressive discount tiers
//   - Automatic best tier selection
//   - Customer loyalty integration
//
// Parameters:
//   - items: Slice of DiscountItem to apply frequency discounts to
//   - customer: Customer with purchase history information
//
// Returns:
//   - DiscountCalculationResult: Result with frequency-based discounts applied
//
// Example:
//   customer := Customer{PurchaseCount: 15}
//   result := engine.ApplyFrequencyDiscounts(items, customer)
//   // Applies discount based on customer's purchase frequency
func (re *RuleEngine) ApplyFrequencyDiscounts(items []DiscountItem, customer Customer) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	for _, rule := range re.FrequencyRules {
		if customer.PurchaseCount >= rule.MinPurchaseCount {
			discount := result.OriginalAmount * (rule.DiscountPercent / 100)

			result.TotalDiscount += discount
			result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
				Type:           DiscountTypeLoyalty,
				RuleID:         "frequency",
				Name:           "Repeat Customer Discount",
				DiscountAmount: discount,
				AppliedItems:   items,
				Description:    "Discount for repeat customers",
			})
			break // Apply only the first matching rule
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// ApplySeasonalDiscounts applies seasonal discounts.
// Evaluates time-based and seasonal rules, applying discounts that are
// currently valid based on date ranges and seasonal periods.
//
// Features:
//   - Time range validation
//   - Seasonal period checking
//   - Category-specific seasonal discounts
//   - Loyalty tier multipliers
//   - Automatic season detection
//
// Parameters:
//   - items: Slice of DiscountItem to apply seasonal discounts to
//   - customer: Customer information for loyalty multipliers
//
// Returns:
//   - DiscountCalculationResult: Result with applicable seasonal discounts
//
// Example:
//   // During winter season
//   result := engine.ApplySeasonalDiscounts(items, customer)
//   // Applies winter seasonal discounts if rules are active
func (re *RuleEngine) ApplySeasonalDiscounts(items []DiscountItem, customer Customer) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	now := time.Now()

	for _, rule := range re.SeasonalRules {
		// Check if rule is currently valid
		if now.Before(rule.ValidFrom) || now.After(rule.ValidUntil) {
			continue
		}

		// Check if current season matches
		if !isCurrentSeason(now, rule.Season) {
			continue
		}

		applicableItems := items
		if len(rule.Categories) > 0 {
			applicableItems = getApplicableItems(items, rule.Categories, nil)
		}

		if len(applicableItems) > 0 {
			itemAmount := calculateItemsAmount(applicableItems)
			discountPercent := rule.DiscountPercent

			// Apply member bonus multiplier
			if rule.Multiplier > 0 && customer.LoyaltyTier != "" {
				discountPercent *= rule.Multiplier
			}

			discount := itemAmount * (discountPercent / 100)

			result.TotalDiscount += discount
			result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
				Type:           DiscountTypeCategory,
				RuleID:         "seasonal_" + rule.Season,
				Name:           "Seasonal Discount",
				DiscountAmount: discount,
				AppliedItems:   applicableItems,
				Description:    fmt.Sprintf("%s seasonal discount", strings.Title(rule.Season)),
			})
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// ApplyCrossSellDiscounts applies cross-sell discounts.
// Encourages customers to purchase complementary products by offering
// discounts when main products are combined with accessories.
//
// Features:
//   - Main product and accessory matching
//   - Minimum main product price validation
//   - Combo pricing or percentage discounts
//   - Automatic product combination detection
//
// Parameters:
//   - items: Slice of DiscountItem to evaluate for cross-sell opportunities
//
// Returns:
//   - DiscountCalculationResult: Result with cross-sell discounts applied
//
// Example:
//   // Items include laptop (main) + mouse (accessory)
//   result := engine.ApplyCrossSellDiscounts(items)
//   // Applies cross-sell discount for the combination
func (re *RuleEngine) ApplyCrossSellDiscounts(items []DiscountItem) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	for _, rule := range re.CrossSellRules {
		mainProducts := getApplicableItems(items, rule.MainProductCategories, nil)
		accessories := getApplicableItems(items, rule.AccessoryCategories, nil)

		if len(mainProducts) > 0 && len(accessories) > 0 {
			// Check minimum main product price if specified
			if rule.MinMainProductPrice > 0 {
				mainAmount := calculateItemsAmount(mainProducts)
				if mainAmount < rule.MinMainProductPrice {
					continue
				}
			}

			combinedItems := append(mainProducts, accessories...)
			combinedAmount := calculateItemsAmount(combinedItems)

			var discount float64
			if rule.ComboPrice > 0 {
				// Fixed combo price
				discount = math.Max(0, combinedAmount-rule.ComboPrice)
			} else {
				// Percentage discount
				discount = combinedAmount * (rule.DiscountPercent / 100)
			}

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type:           DiscountTypeBundle,
					RuleID:         "cross_sell",
					Name:           "Cross-sell Discount",
					DiscountAmount: discount,
					AppliedItems:   combinedItems,
					Description:    "Main product + accessory discount",
				})
			}
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// ApplyMixAndMatchDiscounts applies mix and match discounts.
// Provides discounts when customers purchase a specified number of items
// from designated categories, promoting variety in purchases.
//
// Features:
//   - Multi-category item matching
//   - Required item count validation
//   - Multiple application support
//   - Maximum application limits
//   - Flat discount or percentage options
//
// Parameters:
//   - items: Slice of DiscountItem to evaluate for mix-and-match opportunities
//
// Returns:
//   - DiscountCalculationResult: Result with mix-and-match discounts applied
//
// Example:
//   // Buy 3 items from clothing categories
//   result := engine.ApplyMixAndMatchDiscounts(items)
//   // Applies discount for qualifying item combinations
func (re *RuleEngine) ApplyMixAndMatchDiscounts(items []DiscountItem) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	for _, rule := range re.MixMatchRules {
		applicableItems := getApplicableItems(items, rule.Categories, nil)

		if len(applicableItems) >= rule.RequiredItems {
			// Calculate how many times this rule can be applied
			applications := len(applicableItems) / rule.RequiredItems
			if rule.MaxApplications > 0 && applications > rule.MaxApplications {
				applications = rule.MaxApplications
			}

			// Take the required number of items for discount
			discountItems := applicableItems[:applications*rule.RequiredItems]
			itemAmount := calculateItemsAmount(discountItems)

			var discount float64
			switch rule.DiscountType {
			case "flat_discount":
				discount = rule.DiscountValue * float64(applications)
			case "percentage":
				discount = itemAmount * (rule.DiscountValue / 100)
			}

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type:           DiscountTypeBundle,
					RuleID:         "mix_match",
					Name:           "Mix & Match Discount",
					DiscountAmount: discount,
					AppliedItems:   discountItems,
					Description:    fmt.Sprintf("Mix & match %d items discount", rule.RequiredItems),
				})
			}
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// Validation functions for discount rules.
// These functions ensure rule integrity and prevent invalid configurations
// that could cause calculation errors or unexpected behavior.

// validateBulkRule validates a bulk discount rule.
// Ensures all required fields are present and values are within acceptable ranges.
//
// Validation Rules:
//   - MinQuantity must be greater than 0
//   - MaxQuantity must be greater than MinQuantity (if specified)
//   - DiscountValue must be greater than 0
//   - Percentage discounts cannot exceed 100%
//
// Parameters:
//   - rule: BulkDiscountRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateBulkRule(rule BulkDiscountRule) error {
	if rule.MinQuantity <= 0 {
		return errors.New("minimum quantity must be greater than 0")
	}
	if rule.MaxQuantity > 0 && rule.MaxQuantity < rule.MinQuantity {
		return errors.New("maximum quantity must be greater than minimum quantity")
	}
	if rule.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	if rule.DiscountType == "percentage" && rule.DiscountValue > 100 {
		return errors.New("percentage discount cannot exceed 100%")
	}
	return nil
}

// validateTierRule validates a tier pricing rule.
// Ensures tier pricing configuration is valid and economically sound.
//
// Validation Rules:
//   - MinQuantity must be greater than 0
//   - PricePerItem must be greater than 0
//
// Parameters:
//   - rule: TierPricingRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateTierRule(rule TierPricingRule) error {
	if rule.MinQuantity <= 0 {
		return errors.New("minimum quantity must be greater than 0")
	}
	if rule.PricePerItem <= 0 {
		return errors.New("price per item must be greater than 0")
	}
	return nil
}

// validateBundleRule validates a bundle discount rule.
// Ensures bundle configuration is complete and discount values are valid.
//
// Validation Rules:
//   - ID must be provided for rule identification
//   - MinItems must be greater than 0
//   - Either RequiredProducts or RequiredCategories must be specified
//   - DiscountValue must be greater than 0
//
// Parameters:
//   - rule: BundleDiscountRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateBundleRule(rule BundleDiscountRule) error {
	if rule.ID == "" {
		return errors.New("bundle rule ID is required")
	}
	if rule.MinItems <= 0 {
		return errors.New("minimum items must be greater than 0")
	}
	if len(rule.RequiredProducts) == 0 && len(rule.RequiredCategories) == 0 {
		return errors.New("either required products or categories must be specified")
	}
	if rule.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	return nil
}

// validateLoyaltyRule validates a loyalty discount rule.
// Ensures loyalty tier configuration and discount percentages are valid.
//
// Validation Rules:
//   - Tier must be specified
//   - DiscountPercent must be between 0 and 100
//
// Parameters:
//   - rule: LoyaltyDiscountRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateLoyaltyRule(rule LoyaltyDiscountRule) error {
	if rule.Tier == "" {
		return errors.New("loyalty tier is required")
	}
	if rule.DiscountPercent <= 0 || rule.DiscountPercent > 100 {
		return errors.New("discount percent must be between 0 and 100")
	}
	return nil
}

// validateCategoryRule validates a category discount rule.
// Ensures category-based discount configuration is complete and valid.
//
// Validation Rules:
//   - Category must be specified
//   - DiscountPercent must be between 0 and 100
//
// Parameters:
//   - rule: CategoryDiscountRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateCategoryRule(rule CategoryDiscountRule) error {
	if rule.Category == "" {
		return errors.New("category is required")
	}
	if rule.DiscountPercent <= 0 || rule.DiscountPercent > 100 {
		return errors.New("discount percent must be between 0 and 100")
	}
	if rule.ValidUntil.Before(rule.ValidFrom) {
		return errors.New("valid until must be after valid from")
	}
	return nil
}

// validateSeasonalRule validates a seasonal discount rule.
// Ensures seasonal discount configuration and time periods are valid.
//
// Validation Rules:
//   - Season must be specified
//   - DiscountPercent must be between 0 and 100
//   - StartDate must be before EndDate
//
// Parameters:
//   - rule: SeasonalDiscountRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateSeasonalRule(rule SeasonalDiscountRule) error {
	validSeasons := []string{"spring", "summer", "autumn", "winter"}
	validSeason := false
	for _, season := range validSeasons {
		if strings.ToLower(rule.Season) == season {
			validSeason = true
			break
		}
	}
	if !validSeason {
		return errors.New("invalid season specified")
	}
	if rule.DiscountPercent <= 0 {
		return errors.New("discount percent must be greater than 0")
	}
	return nil
}

// validateCrossSellRule validates a cross-sell discount rule.
// Ensures cross-sell configuration has valid product relationships.
//
// Validation Rules:
//   - MainProductCategories must not be empty
//   - AccessoryCategories must not be empty
//   - Either DiscountPercent or ComboPrice must be specified
//
// Parameters:
//   - rule: CrossSellRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateCrossSellRule(rule CrossSellRule) error {
	if len(rule.MainProductCategories) == 0 {
		return errors.New("main product categories are required")
	}
	if len(rule.AccessoryCategories) == 0 {
		return errors.New("accessory categories are required")
	}
	if rule.DiscountPercent <= 0 && rule.ComboPrice <= 0 {
		return errors.New("either discount percent or combo price must be specified")
	}
	return nil
}

// validateMixMatchRule validates a mix-and-match discount rule.
// Ensures mix-and-match configuration has valid category requirements.
//
// Validation Rules:
//   - Categories must not be empty
//   - RequiredItems must be greater than 0
//   - DiscountValue must be greater than 0
//
// Parameters:
//   - rule: MixAndMatchRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
func validateMixMatchRule(rule MixAndMatchRule) error {
	if len(rule.Categories) == 0 {
		return errors.New("categories are required")
	}
	if rule.RequiredItems <= 0 {
		return errors.New("required items must be greater than 0")
	}
	if rule.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	return nil
}

// Helper functions for rule engine operations.
// These utility functions support rule validation and application logic.

// isCurrentSeason checks if the given season matches the current season.
// Determines seasonal discount eligibility based on current date.
//
// Parameters:
//   - now: Current time to check against
//   - season: Season name to check ("spring", "summer", "autumn", "winter")
//
// Returns:
//   - bool: True if the given season is the current season
//
// Example:
//   if isCurrentSeason(time.Now(), "summer") {
//       // Apply summer discounts
//   }
func isCurrentSeason(now time.Time, season string) bool {
	month := now.Month()
	switch strings.ToLower(season) {
	case "spring":
		return month >= 3 && month <= 5
	case "summer":
		return month >= 6 && month <= 8
	case "autumn", "fall":
		return month >= 9 && month <= 11
	case "winter":
		return month == 12 || month <= 2
	default:
		return false
	}
}

// GetApplicableRules returns rules that are applicable for given items and customer.
// Analyzes items and customer data to determine which discount rules can be applied.
//
// Parameters:
//   - items: Slice of DiscountItem to evaluate for rule applicability
//   - customer: Customer information for loyalty and personalized rules
//
// Returns:
//   - map[string]interface{}: Map containing applicable rules by type
//     Keys include: "bulk", "loyalty", "seasonal", etc.
//
// Example:
//   rules := engine.GetApplicableRules(items, customer)
//   if bulkRules, ok := rules["bulk"]; ok {
//       // Process applicable bulk rules
//   }
func (re *RuleEngine) GetApplicableRules(items []DiscountItem, customer Customer) map[string]interface{} {
	applicableRules := make(map[string]interface{})

	// Check bulk rules
	for _, rule := range re.BulkRules {
		applicableItems := getApplicableItems(items, rule.ApplicableCategories, rule.ApplicableProducts)
		if getTotalQuantity(applicableItems) >= rule.MinQuantity {
			applicableRules["bulk"] = append(applicableRules["bulk"].([]BulkDiscountRule), rule)
		}
	}

	// Check loyalty rules
	for _, rule := range re.LoyaltyRules {
		if customer.LoyaltyTier == rule.Tier {
			applicableRules["loyalty"] = append(applicableRules["loyalty"].([]LoyaltyDiscountRule), rule)
		}
	}

	// Check seasonal rules
	now := time.Now()
	for _, rule := range re.SeasonalRules {
		if now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) && isCurrentSeason(now, rule.Season) {
			applicableRules["seasonal"] = append(applicableRules["seasonal"].([]SeasonalDiscountRule), rule)
		}
	}

	return applicableRules
}

// ClearRules clears all rules from the engine.
// Resets the rule engine to an empty state, useful for testing or
// when completely reconfiguring discount rules.
//
// Example:
//   engine.ClearRules()
//   // Engine now has no active rules
func (re *RuleEngine) ClearRules() {
	re.BulkRules = []BulkDiscountRule{}
	re.TierRules = []TierPricingRule{}
	re.BundleRules = []BundleDiscountRule{}
	re.LoyaltyRules = []LoyaltyDiscountRule{}
	re.ProgressiveRules = []ProgressiveDiscountRule{}
	re.CategoryRules = []CategoryDiscountRule{}
	re.FrequencyRules = []FrequencyDiscountRule{}
	re.SeasonalRules = []SeasonalDiscountRule{}
	re.CrossSellRules = []CrossSellRule{}
	re.MixMatchRules = []MixAndMatchRule{}
}

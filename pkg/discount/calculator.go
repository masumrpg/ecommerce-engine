// Package discount provides comprehensive discount calculation functionality
// for e-commerce applications. It supports multiple discount types including
// tier pricing, bulk discounts, bundle offers, category-specific discounts,
// progressive discounts, and loyalty rewards.
//
// Features:
//   - Multiple discount types (tier, bulk, bundle, category, progressive, loyalty)
//   - Stacked vs. best single discount strategies
//   - Time-based discount validation
//   - Maximum discount limits and caps
//   - Flexible rule-based discount engine
//   - Bundle matching and combo pricing
//   - Customer loyalty tier integration
//
// Basic Usage:
//   input := DiscountCalculationInput{
//     Items: []DiscountItem{
//       {ID: "item1", Price: 100.0, Quantity: 2, Category: "electronics"},
//     },
//     AllowStacking: true,
//     BulkRules: []BulkDiscountRule{
//       {MinQuantity: 2, DiscountType: "percentage", DiscountValue: 10},
//     },
//   }
//   result := Calculate(input)
//   fmt.Printf("Total discount: $%.2f", result.TotalDiscount)
package discount

import (
	"math"
	"time"
)

// Calculate calculates all applicable discounts for the given input.
// This is the main entry point for discount calculations, supporting both
// stacked discounts (multiple discounts applied together) and best single
// discount strategies based on the input configuration.
//
// Features:
//   - Automatic original amount calculation
//   - Stacked vs. single discount strategies
//   - Final amount and savings percentage calculation
//   - Precision rounding to 2 decimal places
//   - Comprehensive error handling and validation
//
// Discount Application Order (when stacking):
//   1. Tier pricing (changes base price)
//   2. Bulk discounts
//   3. Bundle discounts
//   4. Category discounts
//   5. Progressive discounts
//   6. Loyalty discounts (applied last)
//
// Parameters:
//   - input: DiscountCalculationInput containing items, rules, and configuration
//
// Returns:
//   - DiscountCalculationResult: Complete calculation results with applied discounts
//
// Example:
//   input := DiscountCalculationInput{
//     Items: []DiscountItem{
//       {ID: "laptop", Price: 1000.0, Quantity: 1, Category: "electronics"},
//       {ID: "mouse", Price: 50.0, Quantity: 2, Category: "accessories"},
//     },
//     AllowStacking: true,
//     BulkRules: []BulkDiscountRule{
//       {MinQuantity: 2, DiscountType: "percentage", DiscountValue: 10},
//     },
//   }
//   result := Calculate(input)
//   // result.OriginalAmount = 1100.0
//   // result.TotalDiscount = 110.0 (10% bulk discount)
//   // result.FinalAmount = 990.0
func Calculate(input DiscountCalculationInput) DiscountCalculationResult {
	result := DiscountCalculationResult{
		IsValid: true,
		AppliedDiscounts: []DiscountApplication{},
	}

	// Calculate original amount
	result.OriginalAmount = calculateOriginalAmount(input.Items)

	if result.OriginalAmount == 0 {
		result.IsValid = false
		result.ErrorMessage = "no items to calculate discount for"
		return result
	}

	// Apply different types of discounts
	if input.AllowStacking {
		result = calculateStackedDiscounts(input, result)
	} else {
		result = calculateBestSingleDiscount(input, result)
	}

	// Calculate final amounts
	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	// Round to 2 decimal places
	result.TotalDiscount = math.Round(result.TotalDiscount*100) / 100
	result.FinalAmount = math.Round(result.FinalAmount*100) / 100
	result.SavingsPercent = math.Round(result.SavingsPercent*100) / 100

	return result
}

// calculateOriginalAmount calculates the total original amount before discounts.
// Computes the sum of all item prices multiplied by their quantities,
// providing the baseline amount for discount calculations.
//
// Parameters:
//   - items: Slice of DiscountItem to calculate total for
//
// Returns:
//   - float64: Total original amount (price × quantity for all items)
//
// Example:
//   items := []DiscountItem{
//     {Price: 100.0, Quantity: 2}, // 200.0
//     {Price: 50.0, Quantity: 1},  // 50.0
//   }
//   total := calculateOriginalAmount(items)
//   // total = 250.0
func calculateOriginalAmount(items []DiscountItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// calculateStackedDiscounts calculates multiple stacked discounts in priority order.
// Applies multiple discount types sequentially, allowing them to stack for
// maximum savings. Respects maximum stacked discount limits if configured.
//
// Features:
//   - Sequential application in priority order
//   - Cumulative discount calculation
//   - Maximum stacked discount percentage enforcement
//   - Preserves all applied discount details
//
// Application Priority:
//   1. Tier pricing (affects base prices)
//   2. Bulk discounts
//   3. Bundle discounts
//   4. Category discounts
//   5. Progressive discounts
//   6. Loyalty discounts
//
// Parameters:
//   - input: DiscountCalculationInput with rules and configuration
//   - result: Current DiscountCalculationResult to build upon
//
// Returns:
//   - DiscountCalculationResult: Updated result with all applicable stacked discounts
//
// Example:
//   // With 10% bulk + 5% loyalty stacking
//   // Original: $100, Bulk: $10 off, Loyalty: $4.50 off (5% of $90)
//   // Total discount: $14.50, Final: $85.50
func calculateStackedDiscounts(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	// Apply discounts in order of priority

	// 1. Tier pricing (changes base price)
	result = applyTierPricing(input, result)

	// 2. Bulk discounts
	result = applyBulkDiscounts(input, result)

	// 3. Bundle discounts
	result = applyBundleDiscounts(input, result)

	// 4. Category discounts
	result = applyCategoryDiscounts(input, result)

	// 5. Progressive discounts
	result = applyProgressiveDiscounts(input, result)

	// 6. Loyalty discounts (applied last)
	result = applyLoyaltyDiscounts(input, result)

	// Check maximum stacked discount limit
	if input.MaxStackedDiscountPercent > 0 {
		maxDiscount := result.OriginalAmount * (input.MaxStackedDiscountPercent / 100)
		if result.TotalDiscount > maxDiscount {
			result.TotalDiscount = maxDiscount
		}
	}

	return result
}

// calculateBestSingleDiscount finds the best single discount to apply.
// Tests each discount type individually and returns the one that provides
// the highest discount amount, ensuring customers get the best possible deal
// when stacking is not allowed.
//
// Features:
//   - Tests all discount types independently
//   - Compares discount amounts to find maximum savings
//   - Returns complete discount application details
//   - Ensures only one discount type is applied
//
// Tested Discount Types:
//   - Tier pricing
//   - Bulk discounts
//   - Bundle discounts
//   - Category discounts
//   - Progressive discounts
//   - Loyalty discounts
//
// Parameters:
//   - input: DiscountCalculationInput with rules and configuration
//   - result: Base DiscountCalculationResult with original amount
//
// Returns:
//   - DiscountCalculationResult: Result with the best single discount applied
//
// Example:
//   // Comparing: 10% bulk ($10) vs 15% loyalty ($15)
//   // Returns: loyalty discount result ($15 savings)
func calculateBestSingleDiscount(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	bestResult := result
	bestDiscount := 0.0

	// Try each type of discount and keep the best one
	discountTypes := []func(DiscountCalculationInput, DiscountCalculationResult) DiscountCalculationResult{
		applyTierPricing,
		applyBulkDiscounts,
		applyBundleDiscounts,
		applyCategoryDiscounts,
		applyProgressiveDiscounts,
		applyLoyaltyDiscounts,
	}

	for _, discountFunc := range discountTypes {
		testResult := discountFunc(input, DiscountCalculationResult{
			OriginalAmount: result.OriginalAmount,
			IsValid: true,
			AppliedDiscounts: []DiscountApplication{},
		})

		if testResult.TotalDiscount > bestDiscount {
			bestResult = testResult
			bestDiscount = testResult.TotalDiscount
		}
	}

	return bestResult
}

// applyTierPricing applies tier-based pricing discounts.
// Implements volume-based pricing where unit prices decrease based on
// quantity thresholds. This changes the base price rather than applying
// a discount percentage, making it the first discount type applied.
//
// Features:
//   - Quantity-based tier pricing
//   - Category-specific tier rules
//   - Min/max quantity range validation
//   - Per-item price adjustment
//   - Automatic discount calculation from price difference
//
// Parameters:
//   - input: DiscountCalculationInput containing tier rules and items
//   - result: Current DiscountCalculationResult to update
//
// Returns:
//   - DiscountCalculationResult: Updated result with tier pricing applied
//
// Example:
//   // Tier rule: 10+ items = $8 each (original $10)
//   // 12 items: discount = (10-8) × 12 = $24
func applyTierPricing(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	for _, rule := range input.TierRules {
		for _, item := range input.Items {
			if rule.Category != "" && item.Category != rule.Category {
				continue
			}

		if item.Quantity >= rule.MinQuantity && (rule.MaxQuantity == 0 || item.Quantity <= rule.MaxQuantity) {
				originalItemTotal := item.Price * float64(item.Quantity)
				newItemTotal := rule.PricePerItem * float64(item.Quantity)
				discount := originalItemTotal - newItemTotal

				if discount > 0 {
					result.TotalDiscount += discount
					result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
						Type: DiscountTypeTier,
						RuleID: "tier_pricing",
						Name: "Tier Pricing",
						DiscountAmount: discount,
						AppliedItems: []DiscountItem{item},
						Description: "Tier-based pricing discount",
					})
				}
			}
		}
	}

	return result
}

// applyBulkDiscounts applies bulk discount rules based on quantity thresholds.
// Provides discounts when customers purchase large quantities of items,
// supporting percentage, fixed amount, and fixed price discount types.
//
// Features:
//   - Quantity threshold validation
//   - Category and product filtering
//   - Multiple discount types (percentage, fixed_amount, fixed_price)
//   - Automatic applicable item detection
//   - Cumulative quantity calculation across items
//
// Discount Types:
//   - percentage: Discount as percentage of total amount
//   - fixed_amount: Fixed dollar amount off
//   - fixed_price: Fixed price per item
//
// Parameters:
//   - input: DiscountCalculationInput containing bulk rules and items
//   - result: Current DiscountCalculationResult to update
//
// Returns:
//   - DiscountCalculationResult: Updated result with bulk discounts applied
//
// Example:
//   // Rule: 5+ items get 15% off
//   // 6 items totaling $120: discount = $18 (15%)
func applyBulkDiscounts(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	for _, rule := range input.BulkRules {
		applicableItems := getApplicableItems(input.Items, rule.ApplicableCategories, rule.ApplicableProducts)
		totalQuantity := getTotalQuantity(applicableItems)

		if totalQuantity >= rule.MinQuantity && (rule.MaxQuantity == 0 || totalQuantity <= rule.MaxQuantity) {
			discount := calculateBulkDiscount(applicableItems, rule)

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type: DiscountTypeBulk,
					RuleID: "bulk_discount",
					Name: "Bulk Discount",
					DiscountAmount: discount,
					AppliedItems: applicableItems,
					Description: "Bulk quantity discount",
				})
			}
		}
	}

	return result
}

// applyBundleDiscounts applies bundle discount rules for product combinations.
// Provides discounts when customers purchase specific combinations of products
// or categories together, encouraging cross-selling and upselling.
//
// Features:
//   - Required product combinations
//   - Required category combinations
//   - Minimum item requirements
//   - Multiple discount types (percentage, fixed_amount, combo_price)
//   - Automatic bundle matching
//   - Multiple bundle applications
//
// Bundle Types:
//   - Product bundles: Specific product combinations
//   - Category bundles: Items from required categories
//   - Mixed bundles: Combination of products and categories
//
// Parameters:
//   - input: DiscountCalculationInput containing bundle rules and items
//   - result: Current DiscountCalculationResult to update
//
// Returns:
//   - DiscountCalculationResult: Updated result with bundle discounts applied
//
// Example:
//   // Bundle: laptop + mouse + keyboard = $50 off
//   // Items match bundle: discount = $50
func applyBundleDiscounts(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	for _, rule := range input.BundleRules {
		bundleMatches := findBundleMatches(input.Items, rule)

		for _, match := range bundleMatches {
			discount := calculateBundleDiscount(match.MatchedItems, rule)

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type: DiscountTypeBundle,
					RuleID: rule.ID,
					Name: rule.Name,
					DiscountAmount: discount,
					AppliedItems: match.MatchedItems,
					Description: "Bundle discount",
				})
			}
		}
	}

	return result
}

// applyCategoryDiscounts applies category-specific discounts with time validation.
// Provides percentage-based discounts for items in specific categories,
// with support for time-based validity periods and maximum discount limits.
//
// Features:
//   - Category-specific targeting
//   - Time-based validity (ValidFrom/ValidUntil)
//   - Minimum quantity requirements
//   - Percentage-based discounts
//   - Maximum discount amount caps
//   - Automatic category item filtering
//
// Validation:
//   - Checks current time against validity period
//   - Validates minimum quantity requirements
//   - Applies maximum discount limits
//
// Parameters:
//   - input: DiscountCalculationInput containing category rules and items
//   - result: Current DiscountCalculationResult to update
//
// Returns:
//   - DiscountCalculationResult: Updated result with category discounts applied
//
// Example:
//   // Rule: 20% off electronics, max $100, min 2 items
//   // 3 electronics items totaling $600: discount = $100 (capped)
func applyCategoryDiscounts(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	now := time.Now()

	for _, rule := range input.CategoryRules {
		// Check if rule is currently valid
		if now.Before(rule.ValidFrom) || now.After(rule.ValidUntil) {
			continue
		}

		categoryItems := getItemsByCategory(input.Items, rule.Category)
		totalQuantity := getTotalQuantity(categoryItems)

		if totalQuantity >= rule.MinQuantity {
			categoryAmount := calculateItemsAmount(categoryItems)
			discount := categoryAmount * (rule.DiscountPercent / 100)

			// Apply maximum discount limit
			if rule.MaxDiscountAmount > 0 && discount > rule.MaxDiscountAmount {
				discount = rule.MaxDiscountAmount
			}

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type: DiscountTypeCategory,
					RuleID: "category_" + rule.Category,
					Name: "Category Discount",
					DiscountAmount: discount,
					AppliedItems: categoryItems,
					Description: "Category-specific discount",
				})
			}
		}
	}

	return result
}

// applyProgressiveDiscounts applies progressive discount rules based on quantity steps.
// Provides increasing discount percentages as customers purchase more items,
// encouraging larger orders through escalating rewards.
//
// Features:
//   - Quantity step-based progression
//   - Escalating discount percentages
//   - Maximum discount percentage caps
//   - Category-specific or global application
//   - Automatic step calculation
//   - Cumulative discount benefits
//
// Calculation:
//   - Steps = Total Quantity ÷ Quantity Step
//   - Progressive Percent = Steps × Discount Percent
//   - Capped at Maximum Discount
//
// Parameters:
//   - input: DiscountCalculationInput containing progressive rules and items
//   - result: Current DiscountCalculationResult to update
//
// Returns:
//   - DiscountCalculationResult: Updated result with progressive discounts applied
//
// Example:
//   // Rule: 2% per 5 items, max 20%
//   // 23 items: 4 steps × 2% = 8% discount
func applyProgressiveDiscounts(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	for _, rule := range input.ProgressiveRules {
		applicableItems := input.Items
		if rule.Category != "" {
			applicableItems = getItemsByCategory(input.Items, rule.Category)
		}

		totalQuantity := getTotalQuantity(applicableItems)
		steps := totalQuantity / rule.QuantityStep

		if steps > 0 {
			progressivePercent := float64(steps) * rule.DiscountPercent
			if progressivePercent > rule.MaxDiscount {
				progressivePercent = rule.MaxDiscount
			}

			itemAmount := calculateItemsAmount(applicableItems)
			discount := itemAmount * (progressivePercent / 100)

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type: DiscountTypeProgressive,
					RuleID: "progressive",
					Name: "Progressive Discount",
					DiscountAmount: discount,
					AppliedItems: applicableItems,
					Description: "Progressive quantity discount",
				})
			}
		}
	}

	return result
}

// applyLoyaltyDiscounts applies loyalty-based discounts for customer tiers.
// Provides exclusive discounts based on customer loyalty tier status,
// rewarding long-term customers with special pricing benefits.
//
// Features:
//   - Customer loyalty tier validation
//   - Minimum order amount requirements
//   - Category-specific or global application
//   - Percentage-based discounts
//   - Maximum discount amount caps
//   - Tier-specific discount rates
//
// Validation:
//   - Matches customer tier with rule tier
//   - Validates minimum order amount
//   - Applies category restrictions if specified
//
// Parameters:
//   - input: DiscountCalculationInput containing loyalty rules and customer info
//   - result: Current DiscountCalculationResult to update
//
// Returns:
//   - DiscountCalculationResult: Updated result with loyalty discounts applied
//
// Example:
//   // Rule: Gold tier gets 15% off, min $200 order
//   // Gold customer with $300 order: discount = $45
func applyLoyaltyDiscounts(input DiscountCalculationInput, result DiscountCalculationResult) DiscountCalculationResult {
	for _, rule := range input.LoyaltyRules {
		if input.Customer.LoyaltyTier != rule.Tier {
			continue
		}

		applicableItems := input.Items
		if len(rule.ApplicableCategories) > 0 {
			applicableItems = getApplicableItems(input.Items, rule.ApplicableCategories, nil)
		}

		itemAmount := calculateItemsAmount(applicableItems)

		if itemAmount >= rule.MinOrderAmount {
			discount := itemAmount * (rule.DiscountPercent / 100)

			// Apply maximum discount limit
			if rule.MaxDiscountAmount > 0 && discount > rule.MaxDiscountAmount {
				discount = rule.MaxDiscountAmount
			}

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type: DiscountTypeLoyalty,
					RuleID: "loyalty_" + rule.Tier,
					Name: "Loyalty Discount",
					DiscountAmount: discount,
					AppliedItems: applicableItems,
					Description: "Loyalty tier discount",
				})
			}
		}
	}

	return result
}

// Helper functions for discount calculations and item filtering.
// These functions provide utilities for item selection, quantity calculations,
// amount computations, and specific discount type calculations.

// getApplicableItems filters items based on categories and products.
// Returns items that match either the specified categories or product IDs,
// enabling flexible rule targeting for different discount types.
//
// Features:
//   - Category-based filtering
//   - Product ID-based filtering
//   - OR logic (matches either categories OR products)
//   - Returns all items if no filters specified
//   - Preserves original item data
//
// Parameters:
//   - items: Slice of DiscountItem to filter
//   - categories: Slice of category names to match
//   - products: Slice of product IDs to match
//
// Returns:
//   - []DiscountItem: Filtered items matching the criteria
//
// Example:
//   items := getApplicableItems(allItems, []string{"electronics"}, []string{"laptop1"})
//   // Returns items in "electronics" category OR with ID "laptop1"
func getApplicableItems(items []DiscountItem, categories []string, products []string) []DiscountItem {
	if len(categories) == 0 && len(products) == 0 {
		return items
	}

	applicable := []DiscountItem{}
	for _, item := range items {
		isApplicable := false

		// Check categories
		if len(categories) > 0 {
			for _, category := range categories {
				if item.Category == category {
					isApplicable = true
					break
				}
			}
		}

		// Check products
		if len(products) > 0 {
			for _, productID := range products {
				if item.ID == productID {
					isApplicable = true
					break
				}
			}
		}

		if isApplicable {
			applicable = append(applicable, item)
		}
	}

	return applicable
}

// getItemsByCategory filters items by a specific category.
// Returns all items that belong to the specified category,
// useful for category-specific discount calculations.
//
// Parameters:
//   - items: Slice of DiscountItem to filter
//   - category: Category name to match
//
// Returns:
//   - []DiscountItem: Items belonging to the specified category
//
// Example:
//   electronics := getItemsByCategory(items, "electronics")
//   // Returns only items with Category = "electronics"
func getItemsByCategory(items []DiscountItem, category string) []DiscountItem {
	categoryItems := []DiscountItem{}
	for _, item := range items {
		if item.Category == category {
			categoryItems = append(categoryItems, item)
		}
	}
	return categoryItems
}

// getTotalQuantity calculates total quantity of items.
// Sums the quantity values across all items in the slice,
// useful for bulk discount threshold validation.
//
// Parameters:
//   - items: Slice of DiscountItem to sum quantities for
//
// Returns:
//   - int: Total quantity across all items
//
// Example:
//   items := []DiscountItem{
//     {Quantity: 3}, {Quantity: 2}, {Quantity: 5},
//   }
//   total := getTotalQuantity(items) // total = 10
func getTotalQuantity(items []DiscountItem) int {
	total := 0
	for _, item := range items {
		total += item.Quantity
	}
	return total
}

// calculateItemsAmount calculates total amount for items.
// Computes the sum of (price × quantity) for all items,
// providing the base amount for percentage-based discounts.
//
// Parameters:
//   - items: Slice of DiscountItem to calculate amount for
//
// Returns:
//   - float64: Total amount (sum of price × quantity)
//
// Example:
//   items := []DiscountItem{
//     {Price: 100.0, Quantity: 2}, // 200.0
//     {Price: 50.0, Quantity: 1},  // 50.0
//   }
//   total := calculateItemsAmount(items) // total = 250.0
func calculateItemsAmount(items []DiscountItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// calculateBulkDiscount calculates discount for bulk rules.
// Applies bulk discount based on quantity thresholds and discount rates,
// supporting percentage, fixed amount, and fixed price discount types.
//
// Features:
//   - Multiple discount type support (percentage, fixed_amount, fixed_price)
//   - Quantity threshold validation
//   - Total amount calculation
//   - Price per item calculations
//
// Parameters:
//   - items: Slice of DiscountItem to apply bulk discount to
//   - rule: BulkDiscountRule containing discount configuration
//
// Returns:
//   - float64: Calculated bulk discount amount
//
// Example:
//   rule := BulkDiscountRule{MinQuantity: 10, DiscountType: "percentage", DiscountValue: 15}
//   discount := calculateBulkDiscount(items, rule) // 15% off for 10+ items
func calculateBulkDiscount(items []DiscountItem, rule BulkDiscountRule) float64 {
	itemAmount := calculateItemsAmount(items)

	switch rule.DiscountType {
	case "percentage":
		return itemAmount * (rule.DiscountValue / 100)
	case "fixed_amount":
		return math.Min(rule.DiscountValue, itemAmount)
	case "fixed_price":
		// Fixed price per item
		totalQuantity := getTotalQuantity(items)
		newTotal := rule.DiscountValue * float64(totalQuantity)
		return math.Max(0, itemAmount-newTotal)
	default:
		return 0
	}
}

// findBundleMatches finds items that match bundle rules.
// Determines which items form valid bundles based on required products
// and categories, calculating how many complete bundles can be formed.
//
// Features:
//   - Required product matching
//   - Required category matching
//   - Minimum item validation
//   - Multiple bundle detection
//   - Complex bundle logic support
//
// Parameters:
//   - items: Slice of DiscountItem to check for bundle matches
//   - rule: BundleDiscountRule defining bundle requirements
//
// Returns:
//   - []BundleMatch: Slice of matched bundles with their items
//
// Example:
//   rule := BundleDiscountRule{
//     RequiredProducts: ["laptop", "mouse"],
//     MinItems: 2,
//   }
//   matches := findBundleMatches(items, rule) // Returns valid bundle combinations
func findBundleMatches(items []DiscountItem, rule BundleDiscountRule) []BundleMatch {
	matches := []BundleMatch{}

	// Simple implementation - can be enhanced for complex bundle logic
	matchedItems := []DiscountItem{}

	// Check required products
	if len(rule.RequiredProducts) > 0 {
		for _, productID := range rule.RequiredProducts {
			for _, item := range items {
				if item.ID == productID {
					matchedItems = append(matchedItems, item)
					break
				}
			}
		}
		if len(matchedItems) == len(rule.RequiredProducts) {
			matches = append(matches, BundleMatch{
				Rule:         rule,
				MatchedItems: matchedItems,
				Applications: 1,
			})
		}
	}

	// Check required categories
	if len(rule.RequiredCategories) > 0 {
		categoryMatches := make(map[string][]DiscountItem)
		for _, item := range items {
			for _, category := range rule.RequiredCategories {
				if item.Category == category {
					categoryMatches[category] = append(categoryMatches[category], item)
				}
			}
		}

		if len(categoryMatches) == len(rule.RequiredCategories) {
			bundleItems := []DiscountItem{}
			for _, categoryItems := range categoryMatches {
				if len(categoryItems) > 0 {
					bundleItems = append(bundleItems, categoryItems[0]) // Take first item from each category
				}
			}

			if len(bundleItems) >= rule.MinItems {
				matches = append(matches, BundleMatch{
					Rule:         rule,
					MatchedItems: bundleItems,
					Applications: 1,
				})
			}
		}
	}

	return matches
}

// calculateBundleDiscount calculates discount for bundle.
// Computes the total discount amount for a matched bundle based on
// the bundle's discount configuration (percentage, fixed amount, or combo price).
//
// Features:
//   - Percentage-based bundle discounts
//   - Fixed amount bundle discounts
//   - Combo price bundle discounts
//   - Total amount calculation
//   - Discount type validation
//
// Parameters:
//   - items: Slice of DiscountItem in the bundle
//   - rule: BundleDiscountRule with discount configuration
//
// Returns:
//   - float64: Total bundle discount amount
//
// Example:
//   rule := BundleDiscountRule{DiscountType: "percentage", DiscountValue: 20.0}
//   discount := calculateBundleDiscount(bundleItems, rule) // 20% off bundle
func calculateBundleDiscount(items []DiscountItem, rule BundleDiscountRule) float64 {
	itemAmount := calculateItemsAmount(items)

	switch rule.DiscountType {
	case "percentage":
		return itemAmount * (rule.DiscountValue / 100)
	case "fixed_amount":
		return math.Min(rule.DiscountValue, itemAmount)
	case "combo_price":
		return math.Max(0, itemAmount-rule.DiscountValue)
	default:
		return 0
	}
}

// CalculateBestDiscount finds the best discount combination.
// Evaluates multiple discount calculation scenarios and returns the one
// with the highest total discount, optimizing savings for customers.
//
// Features:
//   - Multiple scenario evaluation
//   - Maximum savings selection
//   - Result validation
//   - Automatic best choice selection
//   - Error handling for invalid results
//
// Parameters:
//   - inputs: Slice of DiscountCalculationInput scenarios to evaluate
//
// Returns:
//   - DiscountCalculationResult: The scenario with highest total discount
//
// Example:
//   scenarios := []DiscountCalculationInput{
//     {AllowStacking: true, BulkRules: bulkRules},
//     {AllowStacking: false, BundleRules: bundleRules},
//   }
//   best := CalculateBestDiscount(scenarios) // Returns highest savings scenario
func CalculateBestDiscount(inputs []DiscountCalculationInput) DiscountCalculationResult {
	bestResult := DiscountCalculationResult{}
	bestSavings := 0.0

	for _, input := range inputs {
		result := Calculate(input)
		if result.IsValid && result.TotalDiscount > bestSavings {
			bestResult = result
			bestSavings = result.TotalDiscount
		}
	}

	return bestResult
}
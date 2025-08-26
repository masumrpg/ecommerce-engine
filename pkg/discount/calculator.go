package discount

import (
	"math"
	"time"
)

// Calculate calculates all applicable discounts for the given input
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

// calculateOriginalAmount calculates the total original amount
func calculateOriginalAmount(items []DiscountItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// calculateStackedDiscounts calculates multiple stacked discounts
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

// calculateBestSingleDiscount finds the best single discount to apply
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

// applyTierPricing applies tier-based pricing
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

// applyBulkDiscounts applies bulk discount rules
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

// applyBundleDiscounts applies bundle discount rules
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

// applyCategoryDiscounts applies category-specific discounts
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

// applyProgressiveDiscounts applies progressive discount rules
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

// applyLoyaltyDiscounts applies loyalty-based discounts
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

// Helper functions

// getApplicableItems filters items based on categories and products
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

// getItemsByCategory filters items by category
func getItemsByCategory(items []DiscountItem, category string) []DiscountItem {
	categoryItems := []DiscountItem{}
	for _, item := range items {
		if item.Category == category {
			categoryItems = append(categoryItems, item)
		}
	}
	return categoryItems
}

// getTotalQuantity calculates total quantity of items
func getTotalQuantity(items []DiscountItem) int {
	total := 0
	for _, item := range items {
		total += item.Quantity
	}
	return total
}

// calculateItemsAmount calculates total amount for items
func calculateItemsAmount(items []DiscountItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

// calculateBulkDiscount calculates discount for bulk rules
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

// findBundleMatches finds items that match bundle rules
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

// calculateBundleDiscount calculates discount for bundle
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

// CalculateBestDiscount finds the best discount combination
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
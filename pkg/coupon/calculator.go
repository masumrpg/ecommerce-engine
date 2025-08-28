// Package coupon provides comprehensive coupon calculation functionality for e-commerce applications.
// It supports various coupon types including percentage discounts, fixed amount discounts,
// buy-X-get-Y promotions, and free shipping offers.
//
// The package handles coupon validation, discount calculations, and provides utilities
// for finding the best applicable coupon from multiple options.
//
// Example usage:
//
//	input := CalculationInput{
//		Coupon: Coupon{
//			Type:  CouponTypePercentage,
//			Value: 10.0, // 10% discount
//			IsActive: true,
//		},
//		OrderAmount: 100.0,
//		Items: []Item{{Price: 50.0, Quantity: 2}},
//	}
//	result := Calculate(input)
//	if result.IsValid {
//		fmt.Printf("Discount: $%.2f", result.DiscountAmount)
//	}
package coupon

import (
	"errors"
	"math"
	"time"
)

// Calculate calculates the discount amount for a given coupon and order.
// It validates the coupon first, then applies the appropriate calculation
// based on the coupon type (percentage, fixed amount, buy-X-get-Y, or free shipping).
//
// Parameters:
//   - input: CalculationInput containing coupon details, order information, and items
//
// Returns:
//   - CalculationResult with discount amount, validity status, and applied items
//
// Example:
//
//	input := CalculationInput{
//		Coupon: Coupon{Type: CouponTypePercentage, Value: 15.0},
//		OrderAmount: 200.0,
//		Items: []Item{{Price: 100.0, Quantity: 2}},
//	}
//	result := Calculate(input)
//	if result.IsValid {
//		fmt.Printf("You saved: $%.2f", result.DiscountAmount)
//	}
func Calculate(input CalculationInput) CalculationResult {
	result := CalculationResult{
		IsValid: false,
	}

	// Validate coupon first
	if validationErr := validateCoupon(input); validationErr != nil {
		result.ErrorMessage = validationErr.Error()
		return result
	}

	// Calculate discount based on coupon type
	switch input.Coupon.Type {
	case CouponTypePercentage:
		return calculatePercentageDiscount(input)
	case CouponTypeFixedAmount:
		return calculateFixedAmountDiscount(input)
	case CouponTypeBuyXGetY:
		return calculateBuyXGetYDiscount(input)
	case CouponTypeFreeShipping:
		return calculateFreeShippingDiscount(input)
	default:
		result.ErrorMessage = "unsupported coupon type"
		return result
	}
}

// calculatePercentageDiscount calculates percentage-based discount for the given coupon.
// It applies the percentage discount to applicable items and respects the maximum discount limit.
// The discount amount is rounded to 2 decimal places for currency precision.
//
// Parameters:
//   - input: CalculationInput containing coupon and order details
//
// Returns:
//   - CalculationResult with calculated percentage discount amount
//
// Example:
//   For a 20% coupon on $100 order: discount = $20.00
func calculatePercentageDiscount(input CalculationInput) CalculationResult {
	result := CalculationResult{IsValid: true}

	applicableAmount := getApplicableAmount(input)
	discountAmount := applicableAmount * (input.Coupon.Value / 100)

	// Apply maximum discount limit
	if input.Coupon.MaxDiscount > 0 && discountAmount > input.Coupon.MaxDiscount {
		discountAmount = input.Coupon.MaxDiscount
	}

	result.DiscountAmount = math.Round(discountAmount*100) / 100
	result.AppliedItems = getApplicableItems(input)
	return result
}

// calculateFixedAmountDiscount calculates fixed amount discount for the given coupon.
// It applies a fixed discount amount but ensures it doesn't exceed the applicable order amount.
// The discount is capped at the total value of applicable items.
//
// Parameters:
//   - input: CalculationInput containing coupon and order details
//
// Returns:
//   - CalculationResult with calculated fixed amount discount
//
// Example:
//   For a $15 fixed discount on $100 order: discount = $15.00
//   For a $15 fixed discount on $10 order: discount = $10.00 (capped)
func calculateFixedAmountDiscount(input CalculationInput) CalculationResult {
	result := CalculationResult{IsValid: true}

	applicableAmount := getApplicableAmount(input)
	discountAmount := input.Coupon.Value

	// Don't exceed the applicable amount
	if discountAmount > applicableAmount {
		discountAmount = applicableAmount
	}

	result.DiscountAmount = math.Round(discountAmount*100) / 100
	result.AppliedItems = getApplicableItems(input)
	return result
}

// calculateBuyXGetYDiscount calculates buy-X-get-Y promotional discount.
// It determines how many free items the customer gets based on their purchase quantity,
// then applies the discount to the cheapest qualifying items.
//
// Parameters:
//   - input: CalculationInput containing coupon with BuyX and GetY values
//
// Returns:
//   - CalculationResult with discount applied to cheapest qualifying items
//
// Example:
//   Buy 2 Get 1 Free: customer buys 4 items, gets 2 items free (cheapest ones)
func calculateBuyXGetYDiscount(input CalculationInput) CalculationResult {
	result := CalculationResult{IsValid: true}

	applicableItems := getApplicableItems(input)
	totalQuantity := 0
	for _, item := range applicableItems {
		totalQuantity += item.Quantity
	}

	// Calculate how many free items user gets
	freeItems := (totalQuantity / input.Coupon.BuyX) * input.Coupon.GetY

	// Find cheapest items to discount
	cheapestItems := findCheapestItems(applicableItems, freeItems)
	discountAmount := 0.0
	for _, item := range cheapestItems {
		discountAmount += item.Price * float64(item.Quantity)
	}

	result.DiscountAmount = math.Round(discountAmount*100) / 100
	result.AppliedItems = cheapestItems
	return result
}

// calculateFreeShippingDiscount calculates free shipping discount for the given coupon.
// This function primarily validates that the coupon is applicable for free shipping.
// The actual shipping discount calculation is typically handled by the shipping calculator.
//
// Parameters:
//   - input: CalculationInput containing coupon and order details
//
// Returns:
//   - CalculationResult with validation status (discount amount is 0.0)
//
// Note:
//   The actual shipping cost reduction is calculated by the shipping module
func calculateFreeShippingDiscount(input CalculationInput) CalculationResult {
	result := CalculationResult{IsValid: true}

	// Free shipping discount amount is typically handled by shipping calculator
	// This just validates the coupon is applicable
	result.DiscountAmount = 0.0 // Actual shipping discount calculated elsewhere
	result.AppliedItems = getApplicableItems(input)
	return result
}

// validateCoupon validates if a coupon can be applied to the given order.
// It performs comprehensive validation including activity status, date validity,
// minimum order requirements, usage limits, and item applicability.
//
// Parameters:
//   - input: CalculationInput containing coupon and order details
//
// Returns:
//   - error: nil if valid, otherwise an error describing the validation failure
//
// Validation checks:
//   - Coupon is active
//   - Current date is within validity period
//   - Order meets minimum amount requirement
//   - Usage limits are not exceeded
//   - At least one applicable item exists
func validateCoupon(input CalculationInput) error {
	coupon := input.Coupon

	// Check if coupon is active
	if !coupon.IsActive {
		return errors.New("coupon is not active")
	}

	// Check date validity
	now := time.Now()
	if now.Before(coupon.ValidFrom) {
		return errors.New("coupon is not yet valid")
	}
	if now.After(coupon.ValidUntil) {
		return errors.New("coupon has expired")
	}

	// Check minimum order amount
	if input.OrderAmount < coupon.MinOrder {
		return errors.New("order amount does not meet minimum requirement")
	}

	// Check usage limits
	if coupon.MaxUsage > 0 && input.Usage.TotalUsage >= coupon.MaxUsage {
		return errors.New("coupon usage limit exceeded")
	}

	if coupon.MaxUsagePerUser > 0 && input.Usage.UsageCount >= coupon.MaxUsagePerUser {
		return errors.New("user usage limit exceeded")
	}

	// Check if there are applicable items
	if len(getApplicableItems(input)) == 0 {
		return errors.New("no applicable items found")
	}

	return nil
}

// getApplicableItems returns items that the coupon can be applied to based on
// the coupon's category and product restrictions. If no restrictions are specified,
// all items are considered applicable.
//
// Parameters:
//   - input: CalculationInput containing coupon restrictions and order items
//
// Returns:
//   - []Item: slice of items that match the coupon's applicability criteria
//
// Logic:
//   - If no categories/products specified: all items are applicable
//   - Otherwise: items must match specified categories or product IDs
func getApplicableItems(input CalculationInput) []Item {
	coupon := input.Coupon
	applicableItems := []Item{}

	// If no specific categories or products, apply to all
	if len(coupon.ApplicableCategories) == 0 && len(coupon.ApplicableProducts) == 0 {
		return input.Items
	}

	for _, item := range input.Items {
		isApplicable := false

		// Check categories
		if len(coupon.ApplicableCategories) > 0 {
			for _, category := range coupon.ApplicableCategories {
				if item.Category == category {
					isApplicable = true
					break
				}
			}
		}

		// Check products
		if len(coupon.ApplicableProducts) > 0 {
			for _, productID := range coupon.ApplicableProducts {
				if item.ID == productID {
					isApplicable = true
					break
				}
			}
		}

		if isApplicable {
			applicableItems = append(applicableItems, item)
		}
	}

	return applicableItems
}

// getApplicableAmount calculates the total monetary amount for items that are
// applicable to the coupon. This is used as the base amount for percentage
// and fixed amount discount calculations.
//
// Parameters:
//   - input: CalculationInput containing coupon and order details
//
// Returns:
//   - float64: total amount of applicable items (price × quantity)
//
// Example:
//   Items: [{Price: 10.0, Quantity: 2}, {Price: 15.0, Quantity: 1}]
//   Result: 35.0 (10×2 + 15×1)
func getApplicableAmount(input CalculationInput) float64 {
	applicableItems := getApplicableItems(input)
	total := 0.0

	for _, item := range applicableItems {
		total += item.Price * float64(item.Quantity)
	}

	return total
}

// findCheapestItems finds the cheapest items up to the specified quantity.
// This function is used for buy-X-get-Y promotions to determine which items
// should receive the discount (typically the cheapest qualifying items).
//
// Parameters:
//   - items: slice of items to search through
//   - quantity: maximum number of items to return
//
// Returns:
//   - []Item: slice of cheapest items, each with quantity 1
//
// Algorithm:
//   1. Expands items by quantity (creates individual item entries)
//   2. Sorts by price in ascending order
//   3. Returns the cheapest items up to the specified quantity
func findCheapestItems(items []Item, quantity int) []Item {
	if quantity <= 0 {
		return []Item{}
	}

	// Create a list of individual items (expand quantities)
	individualItems := []Item{}
	for _, item := range items {
		for i := 0; i < item.Quantity; i++ {
			individualItems = append(individualItems, Item{
				ID:       item.ID,
				Price:    item.Price,
				Quantity: 1,
				Category: item.Category,
			})
		}
	}

	// Sort by price (ascending)
	for i := 0; i < len(individualItems)-1; i++ {
		for j := i + 1; j < len(individualItems); j++ {
			if individualItems[i].Price > individualItems[j].Price {
				individualItems[i], individualItems[j] = individualItems[j], individualItems[i]
			}
		}
	}

	// Take the cheapest items up to quantity
	result := []Item{}
	for i := 0; i < quantity && i < len(individualItems); i++ {
		result = append(result, individualItems[i])
	}

	return result
}

// CalculateMultiple calculates discounts for multiple coupons and returns the best one.
// This function is useful when customers have multiple applicable coupons and you want
// to automatically apply the one that provides the maximum discount.
//
// Parameters:
//   - coupons: slice of coupons to evaluate
//   - orderAmount: total order amount
//   - userID: user identifier for usage validation
//   - items: order items
//   - usages: slice of usage data corresponding to each coupon
//
// Returns:
//   - CalculationResult: result of the best coupon (highest discount amount)
//
// Logic:
//   - Evaluates each coupon independently
//   - Returns the result with the highest valid discount amount
//   - Returns invalid result if no coupons are applicable
func CalculateMultiple(coupons []Coupon, orderAmount float64, userID string, items []Item, usages []CouponUsage) CalculationResult {
	bestResult := CalculationResult{IsValid: false}
	bestDiscount := 0.0

	for i, coupon := range coupons {
		usage := CouponUsage{}
		if i < len(usages) {
			usage = usages[i]
		}

		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: orderAmount,
			UserID:      userID,
			Items:       items,
			Usage:       usage,
		}

		result := Calculate(input)
		if result.IsValid && result.DiscountAmount > bestDiscount {
			bestResult = result
			bestDiscount = result.DiscountAmount
		}
	}

	return bestResult
}
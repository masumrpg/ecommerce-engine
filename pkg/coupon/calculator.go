package coupon

import (
	"errors"
	"math"
	"time"
)

// Calculate calculates the discount amount for a given coupon and order
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

// calculatePercentageDiscount calculates percentage-based discount
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

// calculateFixedAmountDiscount calculates fixed amount discount
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

// calculateBuyXGetYDiscount calculates buy X get Y discount
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

// calculateFreeShippingDiscount calculates free shipping discount
func calculateFreeShippingDiscount(input CalculationInput) CalculationResult {
	result := CalculationResult{IsValid: true}

	// Free shipping discount amount is typically handled by shipping calculator
	// This just validates the coupon is applicable
	result.DiscountAmount = 0.0 // Actual shipping discount calculated elsewhere
	result.AppliedItems = getApplicableItems(input)
	return result
}

// validateCoupon validates if coupon can be applied
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

// getApplicableItems returns items that the coupon can be applied to
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

// getApplicableAmount calculates total amount for applicable items
func getApplicableAmount(input CalculationInput) float64 {
	applicableItems := getApplicableItems(input)
	total := 0.0

	for _, item := range applicableItems {
		total += item.Price * float64(item.Quantity)
	}

	return total
}

// findCheapestItems finds the cheapest items up to the specified quantity
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

// CalculateMultiple calculates discounts for multiple coupons and returns the best one
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
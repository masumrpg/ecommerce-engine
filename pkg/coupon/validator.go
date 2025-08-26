package coupon

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ValidateCouponRules validates coupon against multiple rules
func ValidateCouponRules(coupon Coupon, rules []ValidationRule, input CalculationInput, userEligibility UserEligibility) error {
	for _, rule := range rules {
		if err := validateSingleRule(coupon, rule, input, userEligibility); err != nil {
			return err
		}
	}
	return nil
}

// validateSingleRule validates a single validation rule
func validateSingleRule(coupon Coupon, rule ValidationRule, input CalculationInput, userEligibility UserEligibility) error {
	switch rule.Type {
	case "user_based":
		return validateUserBasedRule(rule, input, userEligibility)
	case "order_based":
		return validateOrderBasedRule(rule, input)
	case "time_based":
		return validateTimeBasedRule(rule, coupon)
	case "usage_based":
		return validateUsageBasedRule(rule, coupon, input)
	default:
		return fmt.Errorf("unknown rule type: %s", rule.Type)
	}
}

// validateUserBasedRule validates user-based rules
func validateUserBasedRule(rule ValidationRule, input CalculationInput, userEligibility UserEligibility) error {
	switch rule.Condition {
	case "first_purchase":
		if required, ok := rule.Value.(bool); ok && required {
			if !userEligibility.IsFirstPurchase {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "loyalty_tier":
		if requiredTier, ok := rule.Value.(string); ok {
			if !isValidLoyaltyTier(userEligibility.LoyaltyTier, requiredTier) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "birthday":
		if required, ok := rule.Value.(bool); ok && required {
			if !userEligibility.IsBirthday {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "member_since":
		if daysRequired, ok := rule.Value.(float64); ok {
			daysSinceMember := time.Since(userEligibility.MemberSince).Hours() / 24
			if daysSinceMember < daysRequired {
				return errors.New(rule.ErrorMessage)
			}
		}

	default:
		return fmt.Errorf("unknown user-based condition: %s", rule.Condition)
	}

	return nil
}

// validateOrderBasedRule validates order-based rules
func validateOrderBasedRule(rule ValidationRule, input CalculationInput) error {
	switch rule.Condition {
	case "minimum_amount":
		if minAmount, ok := rule.Value.(float64); ok {
			if input.OrderAmount < minAmount {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "specific_products":
		if productIDs, ok := rule.Value.([]string); ok {
			if !hasSpecificProducts(input.Items, productIDs) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "specific_categories":
		if categories, ok := rule.Value.([]string); ok {
			if !hasSpecificCategories(input.Items, categories) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "minimum_quantity":
		if minQty, ok := rule.Value.(float64); ok {
			totalQty := getTotalQuantity(input.Items)
			if float64(totalQty) < minQty {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "exclude_sale_items":
		if exclude, ok := rule.Value.(bool); ok && exclude {
			if hasSaleItems(input.Items) {
				return errors.New(rule.ErrorMessage)
			}
		}

	default:
		return fmt.Errorf("unknown order-based condition: %s", rule.Condition)
	}

	return nil
}

// validateTimeBasedRule validates time-based rules
func validateTimeBasedRule(rule ValidationRule, coupon Coupon) error {
	now := time.Now()

	switch rule.Condition {
	case "flash_sale":
		if duration, ok := rule.Value.(float64); ok {
			flashSaleEnd := coupon.ValidFrom.Add(time.Duration(duration) * time.Minute)
			if now.After(flashSaleEnd) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "seasonal":
		if season, ok := rule.Value.(string); ok {
			if !isValidSeason(now, season) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "recurring":
		if pattern, ok := rule.Value.(string); ok {
			if !isValidRecurringTime(now, pattern) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "time_window":
		if timeWindow, ok := rule.Value.(map[string]interface{}); ok {
			if !isWithinTimeWindow(now, timeWindow) {
				return errors.New(rule.ErrorMessage)
			}
		}

	default:
		return fmt.Errorf("unknown time-based condition: %s", rule.Condition)
	}

	return nil
}

// validateUsageBasedRule validates usage-based rules
func validateUsageBasedRule(rule ValidationRule, coupon Coupon, input CalculationInput) error {
	switch rule.Condition {
	case "single_use":
		if singleUse, ok := rule.Value.(bool); ok && singleUse {
			if input.Usage.UsageCount > 0 {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "limited_per_user":
		if maxPerUser, ok := rule.Value.(float64); ok {
			if float64(input.Usage.UsageCount) >= maxPerUser {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "total_usage_cap":
		if totalCap, ok := rule.Value.(float64); ok {
			if float64(input.Usage.TotalUsage) >= totalCap {
				return errors.New(rule.ErrorMessage)
			}
		}

	default:
		return fmt.Errorf("unknown usage-based condition: %s", rule.Condition)
	}

	return nil
}

// Helper functions

// isValidLoyaltyTier checks if user's loyalty tier meets requirement
func isValidLoyaltyTier(userTier, requiredTier string) bool {
	tierLevels := map[string]int{
		"bronze":   1,
		"silver":   2,
		"gold":     3,
		"platinum": 4,
	}

	userLevel, userExists := tierLevels[strings.ToLower(userTier)]
	requiredLevel, reqExists := tierLevels[strings.ToLower(requiredTier)]

	if !userExists || !reqExists {
		return false
	}

	return userLevel >= requiredLevel
}

// hasSpecificProducts checks if order contains specific products
func hasSpecificProducts(items []Item, productIDs []string) bool {
	for _, item := range items {
		for _, productID := range productIDs {
			if item.ID == productID {
				return true
			}
		}
	}
	return false
}

// hasSpecificCategories checks if order contains items from specific categories
func hasSpecificCategories(items []Item, categories []string) bool {
	for _, item := range items {
		for _, category := range categories {
			if strings.EqualFold(item.Category, category) {
				return true
			}
		}
	}
	return false
}

// getTotalQuantity calculates total quantity of all items
func getTotalQuantity(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Quantity
	}
	return total
}

// hasSaleItems checks if any items are on sale (assuming sale items have specific category)
func hasSaleItems(items []Item) bool {
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Category), "sale") {
			return true
		}
	}
	return false
}

// isValidSeason checks if current time is within specified season
func isValidSeason(now time.Time, season string) bool {
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

// isValidRecurringTime checks if current time matches recurring pattern
func isValidRecurringTime(now time.Time, pattern string) bool {
	switch strings.ToLower(pattern) {
	case "weekend":
		weekday := now.Weekday()
		return weekday == time.Saturday || weekday == time.Sunday
	case "weekday":
		weekday := now.Weekday()
		return weekday >= time.Monday && weekday <= time.Friday
	case "monthly_first_week":
		return now.Day() <= 7
	case "monthly_last_week":
		nextMonth := now.AddDate(0, 1, 0)
		lastDayOfMonth := nextMonth.AddDate(0, 0, -1).Day()
		return now.Day() > lastDayOfMonth-7
	default:
		return false
	}
}

// isWithinTimeWindow checks if current time is within specified time window
func isWithinTimeWindow(now time.Time, timeWindow map[string]interface{}) bool {
	startHour, startOk := timeWindow["start_hour"].(float64)
	endHour, endOk := timeWindow["end_hour"].(float64)

	if !startOk || !endOk {
		return false
	}

	currentHour := float64(now.Hour()) + float64(now.Minute())/60.0

	if startHour <= endHour {
		return currentHour >= startHour && currentHour <= endHour
	} else {
		// Overnight window (e.g., 22:00 to 06:00)
		return currentHour >= startHour || currentHour <= endHour
	}
}

// ValidateCouponStacking validates if multiple coupons can be stacked
func ValidateCouponStacking(coupons []Coupon, stackingRules map[string]interface{}) error {
	maxStackable, ok := stackingRules["max_stackable"].(float64)
	if ok && float64(len(coupons)) > maxStackable {
		return fmt.Errorf("cannot stack more than %.0f coupons", maxStackable)
	}

	allowSameType, ok := stackingRules["allow_same_type"].(bool)
	if ok && !allowSameType {
		typeCount := make(map[CouponType]int)
		for _, coupon := range coupons {
			typeCount[coupon.Type]++
			if typeCount[coupon.Type] > 1 {
				return errors.New("cannot stack coupons of the same type")
			}
		}
	}

	excludedCombinations, ok := stackingRules["excluded_combinations"].([][]string)
	if ok {
		for _, combination := range excludedCombinations {
			if hasAllCodes(coupons, combination) {
				return errors.New("this combination of coupons cannot be used together")
			}
		}
	}

	return nil
}

// hasAllCodes checks if all specified codes are present in the coupon list
func hasAllCodes(coupons []Coupon, codes []string) bool {
	couponCodes := make(map[string]bool)
	for _, coupon := range coupons {
		couponCodes[coupon.Code] = true
	}

	for _, code := range codes {
		if !couponCodes[code] {
			return false
		}
	}

	return true
}

// ValidateBusinessRules validates business-specific rules
func ValidateBusinessRules(coupon Coupon, input CalculationInput, businessRules map[string]interface{}) error {
	// Validate minimum margin
	if minMargin, ok := businessRules["minimum_margin_percent"].(float64); ok {
		discountPercent := (input.Coupon.Value / input.OrderAmount) * 100
		if discountPercent > (100 - minMargin) {
			return fmt.Errorf("discount exceeds maximum allowed margin")
		}
	}

	// Validate blacklisted users
	if blacklistedUsers, ok := businessRules["blacklisted_users"].([]string); ok {
		for _, userID := range blacklistedUsers {
			if input.UserID == userID {
				return errors.New("user is not eligible for coupons")
			}
		}
	}

	// Validate geographic restrictions
	if allowedRegions, ok := businessRules["allowed_regions"].([]string); ok {
		userRegion, regionOk := businessRules["user_region"].(string)
		if regionOk {
			allowed := false
			for _, region := range allowedRegions {
				if strings.EqualFold(userRegion, region) {
					allowed = true
					break
				}
			}
			if !allowed {
				return errors.New("coupon not available in your region")
			}
		}
	}

	return nil
}
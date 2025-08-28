package coupon

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ValidateCouponRules validates a coupon against multiple validation rules.
// Processes each rule sequentially and returns the first validation error encountered.
// Used to enforce complex business logic and eligibility criteria for coupon usage.
//
// Parameters:
//   - coupon: the coupon to validate
//   - rules: slice of ValidationRule to check against
//   - input: calculation input containing order and user details
//   - userEligibility: user-specific eligibility criteria
//
// Returns:
//   - error: nil if all rules pass, first validation error otherwise
//
// Rule processing:
//   - Rules are evaluated in order
//   - First failing rule stops validation and returns its error
//   - All rules must pass for successful validation
//
// Example:
//
//	rules := []ValidationRule{
//		{Type: "user_based", Condition: "loyalty_tier", Value: "gold"},
//		{Type: "order_based", Condition: "minimum_amount", Value: 100.0},
//	}
//	err := ValidateCouponRules(coupon, rules, input, eligibility)
//	if err != nil {
//		// Handle validation failure
//	}
func ValidateCouponRules(coupon Coupon, rules []ValidationRule, input CalculationInput, userEligibility UserEligibility) error {
	for _, rule := range rules {
		if err := validateSingleRule(coupon, rule, input, userEligibility); err != nil {
			return err
		}
	}
	return nil
}

// validateSingleRule validates a single validation rule against the provided coupon and context.
// Routes validation to the appropriate handler based on rule type.
// Internal function used by ValidateCouponRules for processing individual rules.
//
// Parameters:
//   - coupon: the coupon being validated
//   - rule: the specific validation rule to check
//   - input: calculation input with order and user details
//   - userEligibility: user eligibility criteria
//
// Returns:
//   - error: nil if rule passes, validation error if rule fails
//
// Supported rule types:
//   - "user_based": validates user profile and eligibility
//   - "order_based": validates order content and amounts
//   - "time_based": validates temporal conditions
//   - "usage_based": validates usage patterns and limits
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

// validateUserBasedRule validates rules related to user profile and eligibility.
// Checks user-specific conditions such as loyalty tier, purchase history, and demographics.
// Used to enforce user-targeted coupon restrictions and personalized offers.
//
// Parameters:
//   - rule: validation rule with user-based conditions
//   - input: calculation input containing user ID and context
//   - userEligibility: user profile and eligibility data
//
// Returns:
//   - error: nil if user meets criteria, validation error otherwise
//
// Supported conditions:
//   - "first_purchase": requires user to be making their first purchase
//   - "loyalty_tier": requires minimum loyalty tier level
//   - "birthday": requires current date to be user's birthday period
//   - "member_since": requires minimum membership duration
//   - "minimum_purchase_history": requires minimum number of past purchases
//   - "user_id_pattern": validates user ID format (email, numeric, uuid, etc.)
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

	case "minimum_purchase_history":
		if minPurchases, ok := rule.Value.(float64); ok {
			if float64(input.Usage.UsageCount) < minPurchases {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "user_id_pattern":
		if pattern, ok := rule.Value.(string); ok {
			if !isValidUserIDPattern(input.UserID, pattern) {
				return errors.New(rule.ErrorMessage)
			}
		}

	default:
		return fmt.Errorf("unknown user-based condition: %s", rule.Condition)
	}

	return nil
}

// validateOrderBasedRule validates rules related to order content and requirements.
// Checks order-specific conditions such as minimum amounts, product requirements, and quantities.
// Used to enforce order-level restrictions and ensure coupons apply to appropriate purchases.
//
// Parameters:
//   - rule: validation rule with order-based conditions
//   - input: calculation input containing order details and items
//
// Returns:
//   - error: nil if order meets criteria, validation error otherwise
//
// Supported conditions:
//   - "minimum_amount": requires minimum order total
//   - "specific_products": requires specific product IDs in the order
//   - "specific_categories": requires items from specific categories
//   - "minimum_quantity": requires minimum total item quantity
//   - "exclude_sale_items": prevents usage on sale/discounted items
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

// validateTimeBasedRule validates rules related to temporal conditions and timing.
// Checks time-specific conditions such as flash sales, seasonal availability, and recurring patterns.
// Used to enforce time-limited promotions and schedule-based coupon availability.
//
// Parameters:
//   - rule: validation rule with time-based conditions
//   - coupon: coupon entity containing validity periods
//
// Returns:
//   - error: nil if timing conditions are met, validation error otherwise
//
// Supported conditions:
//   - "flash_sale": validates coupon is within flash sale duration
//   - "seasonal": validates current season matches required season
//   - "recurring": validates recurring time patterns (weekend, weekday, etc.)
//   - "time_window": validates current time is within specified hours
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

// validateUsageBasedRule validates rules related to coupon usage patterns and limits.
// Checks usage-specific conditions such as single-use restrictions, per-user limits, and total caps.
// Used to enforce usage restrictions and prevent coupon abuse.
//
// Parameters:
//   - rule: validation rule with usage-based conditions
//   - coupon: coupon entity containing usage limits
//   - input: calculation input containing current usage statistics
//
// Returns:
//   - error: nil if usage conditions are met, validation error otherwise
//
// Supported conditions:
//   - "single_use": ensures coupon hasn't been used before by this user
//   - "limited_per_user": enforces per-user usage limits
//   - "total_usage_cap": enforces global usage limits across all users
//   - "coupon_expiry_buffer": ensures coupon won't expire too soon
//   - "coupon_value_threshold": validates minimum coupon value requirements
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

	case "coupon_expiry_buffer":
		if bufferHours, ok := rule.Value.(float64); ok {
			bufferTime := time.Duration(bufferHours) * time.Hour
			if time.Now().Add(bufferTime).After(coupon.ValidUntil) {
				return errors.New(rule.ErrorMessage)
			}
		}

	case "coupon_value_threshold":
		if minValue, ok := rule.Value.(float64); ok {
			if coupon.Value < minValue {
				return errors.New(rule.ErrorMessage)
			}
		}

	default:
		return fmt.Errorf("unknown usage-based condition: %s", rule.Condition)
	}

	return nil
}

// Helper functions

// isValidLoyaltyTier checks if a user's loyalty tier meets the minimum requirement.
// Compares tier levels using a hierarchical system where higher tiers satisfy lower requirements.
// Used for validating loyalty-based coupon eligibility.
//
// Parameters:
//   - userTier: user's current loyalty tier ("bronze", "silver", "gold", "platinum")
//   - requiredTier: minimum required tier for coupon eligibility
//
// Returns:
//   - bool: true if user tier meets or exceeds requirement, false otherwise
//
// Tier hierarchy (ascending):
//   1. bronze
//   2. silver
//   3. gold
//   4. platinum
//
// Example:
//   - User "gold" meets requirement "silver" (returns true)
//   - User "bronze" does not meet requirement "gold" (returns false)
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

// hasSpecificProducts checks if an order contains any of the specified product IDs.
// Used for validating product-specific coupon restrictions and requirements.
// Returns true if at least one required product is found in the order.
//
// Parameters:
//   - items: slice of order items to check
//   - productIDs: slice of required product IDs
//
// Returns:
//   - bool: true if any required product is found, false otherwise
//
// Example:
//   - Order contains ["LAPTOP001", "MOUSE002"]
//   - Required products ["LAPTOP001", "TABLET003"]
//   - Returns true (LAPTOP001 matches)
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

// hasSpecificCategories checks if an order contains items from any of the specified categories.
// Used for validating category-specific coupon restrictions and requirements.
// Performs case-insensitive category matching for flexibility.
//
// Parameters:
//   - items: slice of order items to check
//   - categories: slice of required category names
//
// Returns:
//   - bool: true if any required category is found, false otherwise
//
// Example:
//   - Order contains items with categories ["electronics", "books"]
//   - Required categories ["Electronics", "clothing"]
//   - Returns true (case-insensitive match on "electronics")
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

// getTotalQuantity calculates the total quantity of all items in an order.
// Sums up the quantity field from all order items.
// Used for validating minimum quantity requirements for coupon eligibility.
//
// Parameters:
//   - items: slice of order items
//
// Returns:
//   - int: total quantity across all items
//
// Example:
//   - Items: [{Quantity: 2}, {Quantity: 3}, {Quantity: 1}]
//   - Returns: 6
func getTotalQuantity(items []Item) int {
	total := 0
	for _, item := range items {
		total += item.Quantity
	}
	return total
}

// hasSaleItems checks if an order contains any items that are currently on sale.
// Used for validating coupon restrictions that exclude sale/discounted items.
// Identifies sale items by checking for "sale" keyword in item category.
//
// Parameters:
//   - items: slice of order items to check
//
// Returns:
//   - bool: true if any sale items are found, false otherwise
//
// Example:
//   - Item with Category: "electronics-sale" (on sale)
//   - Item with Category: "books" (not on sale)
//   - Returns true if any item category contains "sale"
func hasSaleItems(items []Item) bool {
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Category), "sale") {
			return true
		}
	}
	return false
}

// isValidSeason checks if the current season matches the required season for coupon validity.
// Determines season based on current month using meteorological seasons.
// Used for validating seasonal coupon restrictions and availability.
//
// Parameters:
//   - now: current time to check
//   - season: required season name ("spring", "summer", "autumn", "winter")
//
// Returns:
//   - bool: true if current season matches requirement, false otherwise
//
// Season mapping:
//   - Spring: March, April, May (months 3-5)
//   - Summer: June, July, August (months 6-8)
//   - Autumn: September, October, November (months 9-11)
//   - Winter: December, January, February (months 12, 1-2)
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

// isValidRecurringTime checks if the current time matches a recurring time pattern.
// Supports various recurring patterns for time-based coupon restrictions.
// Used for validating coupons that are only valid during specific recurring periods.
//
// Parameters:
//   - now: current time to check
//   - pattern: recurring time pattern to match
//
// Returns:
//   - bool: true if current time matches pattern, false otherwise
//
// Supported patterns:
//   - "weekend": Saturday and Sunday
//   - "weekday": Monday through Friday
//   - "monthly_first_week": first 7 days of the month
//   - "monthly_last_week": last 7 days of the month
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

// isWithinTimeWindow checks if the current time falls within a specified time window.
// Validates time-based coupon restrictions using hour-based windows.
// Supports both same-day and overnight time windows.
//
// Parameters:
//   - now: current time to check
//   - timeWindow: map containing "start_hour" and "end_hour" as float64 values
//
// Returns:
//   - bool: true if current time is within window, false otherwise
//
// Examples:
//   - Window 9-17: valid from 9 AM to 5 PM
//   - Window 22-6: valid from 10 PM to 6 AM (overnight)
//   - Current time 14:30 with window 9-17: returns true
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

// ValidateCouponStacking validates if multiple coupons can be used together in a single order.
// Checks stacking rules and restrictions to prevent incompatible coupon combinations.
// Used for enforcing business rules around coupon combination policies.
//
// Parameters:
//   - coupons: slice of coupons to validate for stacking compatibility
//   - stackingRules: map containing stacking rules and restrictions
//
// Returns:
//   - error: nil if coupons can be stacked, validation error otherwise
//
// Validation checks:
//   - Maximum number of stackable coupons
//   - Coupon type compatibility (e.g., only one percentage discount)
//   - Mutual exclusivity rules
//   - Total discount limits
//
// Example:
//   - Two percentage coupons: may be restricted
//   - Percentage + free shipping: typically allowed
//   - Coupons with "no_stacking" flag: rejected
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

// hasAllCodes checks if all specified coupon codes are present in the provided coupon list.
// Used for validating excluded combinations and ensuring all required codes are present.
// Supports validation of coupon combination restrictions.
//
// Parameters:
//   - coupons: slice of coupons to check
//   - codes: slice of coupon codes to look for
//
// Returns:
//   - bool: true if all specified codes are found, false otherwise
//
// Example:
//   - Codes: ["SAVE10", "FREESHIP"]
//   - Coupons: [{Code: "SAVE10"}, {Code: "FREESHIP"}, {Code: "EXTRA5"}]
//   - Returns: true (all specified codes present)
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

// isValidUserIDPattern checks if a user ID matches a specified pattern for coupon eligibility.
// Supports various pattern matching for targeted user promotions and validation.
// Used for validating user-specific or group-specific coupon restrictions.
//
// Parameters:
//   - userID: user identifier to validate
//   - pattern: required pattern that user ID must match
//
// Returns:
//   - bool: true if user ID matches pattern, false otherwise
//
// Supported patterns:
//   - "email": validates email format (contains @ and .)
//   - "numeric": validates numeric-only user IDs
//   - "alphanumeric": validates alphanumeric characters only
//   - "uuid": validates UUID format (8-4-4-4-12)
//   - Default: allows any pattern if not specified
func isValidUserIDPattern(userID, pattern string) bool {
	switch strings.ToLower(pattern) {
	case "email":
		return strings.Contains(userID, "@") && strings.Contains(userID, ".")
	case "numeric":
		for _, char := range userID {
			if char < '0' || char > '9' {
				return false
			}
		}
		return len(userID) > 0
	case "alphanumeric":
		for _, char := range userID {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
				return false
			}
		}
		return len(userID) > 0
	case "uuid":
		// Simple UUID pattern check (8-4-4-4-12)
		parts := strings.Split(userID, "-")
		if len(parts) != 5 {
			return false
		}
		expectedLengths := []int{8, 4, 4, 4, 12}
		for i, part := range parts {
			if len(part) != expectedLengths[i] {
				return false
			}
			for _, char := range part {
				if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
					return false
				}
			}
		}
		return true
	default:
		return true // Allow any pattern if not specified
	}
}

// ValidateBusinessRules validates custom business-specific rules for coupon usage.
// Applies organization-specific validation logic beyond standard coupon rules.
// Used for enforcing complex business policies and custom restrictions.
//
// Parameters:
//   - coupon: coupon entity to validate
//   - input: calculation input containing order and user details
//   - businessRules: map containing custom business rules and restrictions
//
// Returns:
//   - error: nil if all business rules pass, validation error otherwise
//
// Supported business rules:
//   - "minimum_margin_percent": ensures minimum profit margin
//   - "blacklisted_users": prevents usage by specific users
//   - "allowed_regions": restricts usage to specific geographic regions
//   - Custom validation logic for organization-specific requirements
//
// Note: This function serves as an extension point for organization-specific
// validation logic that cannot be covered by standard coupon rules.
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
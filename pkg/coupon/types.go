// Package coupon provides comprehensive coupon management functionality for e-commerce applications.
// It includes types, constants, and structures for creating, validating, and calculating coupon discounts
// with support for various coupon types including percentage, fixed amount, buy-X-get-Y, and free shipping.
//
// The package supports advanced features such as:
//   - Multiple coupon types with flexible configuration
//   - Usage tracking and limits per user and globally
//   - Category and product-specific applicability
//   - Time-based validity periods
//   - Complex validation rules and eligibility criteria
//   - Bulk coupon code generation with customizable patterns
//
// Example usage:
//
//	coupon := &Coupon{
//		Code: "SAVE20",
//		Type: CouponTypePercentage,
//		Value: 20.0,
//		MinOrder: 100.0,
//		MaxDiscount: 50.0,
//	}
//	result := Calculate(input)
package coupon

import (
	"time"
)

// CouponType represents the type of discount a coupon provides.
// Each type has different calculation logic and applicable scenarios.
// Used to determine how the coupon value should be interpreted and applied.
type CouponType string

const (
	// CouponTypePercentage applies a percentage discount to the order.
	// Value represents the percentage (0-100) to discount from the applicable amount.
	// Example: Value=20 means 20% discount
	CouponTypePercentage CouponType = "percentage"

	// CouponTypeFixedAmount applies a fixed monetary discount to the order.
	// Value represents the exact amount to subtract from the order total.
	// Example: Value=10.50 means $10.50 discount
	CouponTypeFixedAmount CouponType = "fixed_amount"

	// CouponTypeBuyXGetY provides free items when purchasing a certain quantity.
	// Requires BuyX and GetY fields to be set. Value is typically 0.
	// Example: BuyX=2, GetY=1 means "buy 2, get 1 free"
	CouponTypeBuyXGetY CouponType = "buy_x_get_y"

	// CouponTypeFreeShipping removes shipping costs from the order.
	// Value is typically 0 as it affects shipping calculation, not item prices.
	// Applied when order meets minimum requirements
	CouponTypeFreeShipping CouponType = "free_shipping"
)

// Coupon represents a complete coupon entity with all its properties and constraints.
// Contains all necessary information for coupon validation, calculation, and tracking.
// Supports multiple discount types and advanced features like usage limits and applicability rules.
//
// Field descriptions:
//   - Code: unique identifier for the coupon (e.g., "SAVE20", "WELCOME10")
//   - Type: determines how the discount is calculated (percentage, fixed, etc.)
//   - Value: discount amount - percentage (0-100) or fixed monetary amount
//   - MinOrder: minimum order amount required to use this coupon
//   - MaxDiscount: maximum discount amount (prevents excessive discounts on percentage coupons)
//   - MaxUsage: total number of times this coupon can be used across all users
//   - MaxUsagePerUser: maximum times a single user can use this coupon
//   - ValidFrom/ValidUntil: time window when the coupon is active
//   - IsActive: manual toggle to enable/disable the coupon
//   - BuyX/GetY: for buy-X-get-Y promotions (e.g., buy 2 get 1 free)
//   - ApplicableCategories/Products: restrict coupon to specific items
//
// Example:
//
//	coupon := Coupon{
//		Code: "SUMMER25",
//		Type: CouponTypePercentage,
//		Value: 25.0,
//		MinOrder: 50.0,
//		MaxDiscount: 100.0,
//		MaxUsage: 1000,
//		MaxUsagePerUser: 1,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(0, 3, 0),
//		IsActive: true,
//	}
type Coupon struct {
	Code           string     `json:"code"`
	Type           CouponType `json:"type"`
	Value          float64    `json:"value"`          // Percentage (0-100) or fixed amount
	MinOrder       float64    `json:"min_order"`      // Minimum order amount
	MaxDiscount    float64    `json:"max_discount"`   // Maximum discount amount (for percentage)
	MaxUsage       int        `json:"max_usage"`      // Maximum total usage
	MaxUsagePerUser int       `json:"max_usage_per_user"` // Maximum usage per user
	ValidFrom      time.Time  `json:"valid_from"`
	ValidUntil     time.Time  `json:"valid_until"`
	IsActive       bool       `json:"is_active"`
	BuyX           int        `json:"buy_x,omitempty"`  // For buy_x_get_y type
	GetY           int        `json:"get_y,omitempty"`  // For buy_x_get_y type
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
	ApplicableProducts   []string `json:"applicable_products,omitempty"`
}

// CouponUsage represents tracking information for coupon usage by users.
// Maintains counters for both individual user usage and total coupon usage
// to enforce usage limits and prevent abuse. Essential for coupon validation.
//
// Field descriptions:
//   - CouponCode: the coupon code this usage record tracks
//   - UserID: unique identifier of the user who used the coupon
//   - UsageCount: number of times this specific user has used this coupon
//   - TotalUsage: total number of times this coupon has been used by all users
//
// Usage patterns:
//   - Check UsageCount against Coupon.MaxUsagePerUser before allowing usage
//   - Check TotalUsage against Coupon.MaxUsage for global limits
//   - Increment counters after successful coupon application
//
// Example:
//
//	usage := CouponUsage{
//		CouponCode: "SAVE20",
//		UserID: "user123",
//		UsageCount: 1,
//		TotalUsage: 45,
//	}
type CouponUsage struct {
	CouponCode string `json:"coupon_code"`
	UserID     string `json:"user_id"`
	UsageCount int    `json:"usage_count"`
	TotalUsage int    `json:"total_usage"`
}

// CalculationInput represents all required input data for coupon discount calculation.
// Contains the coupon to apply, order details, user information, and current usage statistics.
// Used as input to the Calculate function to determine discount amount and validity.
//
// Field descriptions:
//   - Coupon: the coupon entity to be applied
//   - OrderAmount: total order amount before any discounts
//   - UserID: identifier of the user attempting to use the coupon
//   - Items: list of items in the order (for category/product-specific coupons)
//   - Usage: current usage statistics for validation
//
// Validation flow:
//   1. Check coupon validity (active, time window)
//   2. Verify user eligibility and usage limits
//   3. Validate minimum order requirements
//   4. Calculate applicable discount amount
//
// Example:
//
//	input := CalculationInput{
//		Coupon: coupon,
//		OrderAmount: 150.00,
//		UserID: "user123",
//		Items: []Item{{ID: "item1", Price: 50.00, Quantity: 3}},
//		Usage: usage,
//	}
type CalculationInput struct {
	Coupon      Coupon  `json:"coupon"`
	OrderAmount float64 `json:"order_amount"`
	UserID      string  `json:"user_id"`
	Items       []Item  `json:"items"`
	Usage       CouponUsage `json:"usage"`
}

// Item represents a single item in an order with pricing and categorization information.
// Used for calculating discounts on specific products or categories, and for buy-X-get-Y promotions.
// Essential for coupons that have product or category restrictions.
//
// Field descriptions:
//   - ID: unique identifier for the product/item
//   - Price: unit price of the item (before any discounts)
//   - Quantity: number of units of this item in the order
//   - Category: product category for category-specific coupon validation
//
// Usage in calculations:
//   - Total item value = Price × Quantity
//   - Category matching for applicable coupon validation
//   - Product ID matching for product-specific coupons
//   - Quantity consideration for buy-X-get-Y promotions
//
// Example:
//
//	item := Item{
//		ID: "LAPTOP001",
//		Price: 999.99,
//		Quantity: 1,
//		Category: "electronics",
//	}
type Item struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
}

// CalculationResult represents the outcome of a coupon discount calculation.
// Contains the calculated discount amount, validation status, and detailed information
// about the calculation process. Used to apply discounts and provide user feedback.
//
// Field descriptions:
//   - DiscountAmount: calculated discount amount to be applied to the order
//   - IsValid: whether the coupon is valid and can be applied
//   - ErrorMessage: detailed error description if IsValid is false
//   - AppliedItems: specific items the coupon discount was applied to
//
// Result interpretation:
//   - If IsValid=true: apply DiscountAmount to the order
//   - If IsValid=false: show ErrorMessage to user, no discount applied
//   - AppliedItems helps track which items received the discount
//
// Example successful result:
//
//	result := CalculationResult{
//		DiscountAmount: 25.00,
//		IsValid: true,
//		AppliedItems: []Item{{ID: "item1", Price: 100.00, Quantity: 1}},
//	}
//
// Example error result:
//
//	result := CalculationResult{
//		DiscountAmount: 0.00,
//		IsValid: false,
//		ErrorMessage: "Coupon has expired",
//	}
type CalculationResult struct {
	DiscountAmount float64 `json:"discount_amount"`
	IsValid        bool    `json:"is_valid"`
	ErrorMessage   string  `json:"error_message,omitempty"`
	AppliedItems   []Item  `json:"applied_items,omitempty"` // Items the coupon was applied to
}

// GeneratorConfig represents configuration parameters for automated coupon code generation.
// Defines patterns, formatting rules, and constraints for generating unique coupon codes.
// Used by the generator functions to create codes that meet specific business requirements.
//
// Field descriptions:
//   - Pattern: code generation pattern ("PREFIX-XXXXXX", "XXXXXXXX", "WORD-NUMBER")
//   - Length: length of the random/variable part of the code
//   - Prefix: fixed text to prepend to generated codes
//   - Suffix: fixed text to append to generated codes
//   - ExcludeChars: characters to avoid in generated codes (default: "0O1I" for clarity)
//   - Count: number of codes to generate in batch operations
//
// Pattern types:
//   - "PREFIX-XXXXXX": prefix + separator + random characters
//   - "XXXXXXXX": pure random characters
//   - "WORD-NUMBER": word + separator + numbers
//   - Custom patterns with X (random char) and N (random number) placeholders
//
// Example:
//
//	config := GeneratorConfig{
//		Pattern: "SAVE-XXXXXX",
//		Length: 6,
//		Prefix: "SAVE",
//		ExcludeChars: "0O1I",
//		Count: 100,
//	}
//	// Generates: SAVE-ABC123, SAVE-DEF456, etc.
type GeneratorConfig struct {
	Pattern    string `json:"pattern"`    // e.g., "PREFIX-XXXXXX", "XXXXXXXX", "WORD-NUMBER"
	Length     int    `json:"length"`     // Length of random part
	Prefix     string `json:"prefix"`     // Prefix for the code
	Suffix     string `json:"suffix"`     // Suffix for the code
	ExcludeChars string `json:"exclude_chars"` // Characters to exclude (default: "0O1I")
	Count      int    `json:"count"`      // Number of codes to generate
}

// ValidationRule represents a single validation constraint for coupon usage.
// Defines specific conditions that must be met for a coupon to be considered valid.
// Multiple rules can be combined to create complex validation logic.
//
// Field descriptions:
//   - Type: category of validation ("user_based", "order_based", "time_based", "usage_based")
//   - Condition: specific condition within the type ("min_age", "first_purchase", etc.)
//   - Value: the value to compare against (age, amount, date, etc.)
//   - ErrorMessage: user-friendly message to display when validation fails
//
// Validation types:
//   - "user_based": user profile validations (age, membership, loyalty tier)
//   - "order_based": order content validations (amount, items, categories)
//   - "time_based": temporal validations (day of week, hour, season)
//   - "usage_based": usage pattern validations (frequency, history)
//
// Example:
//
//	rule := ValidationRule{
//		Type: "user_based",
//		Condition: "loyalty_tier",
//		Value: "gold",
//		ErrorMessage: "This coupon is only available for Gold members",
//	}
type ValidationRule struct {
	Type        string      `json:"type"`        // "user_based", "order_based", "time_based", "usage_based"
	Condition   string      `json:"condition"`   // Specific condition
	Value       any `json:"value"`       // Value to check against
	ErrorMessage string     `json:"error_message"`
}

// UserEligibility represents user-specific criteria for coupon eligibility validation.
// Contains user profile information and status flags used to determine if a user
// qualifies for specific coupons with user-based restrictions.
//
// Field descriptions:
//   - IsFirstPurchase: true if this is the user's first purchase (for welcome coupons)
//   - LoyaltyTier: user's current loyalty program level ("bronze", "silver", "gold", "platinum")
//   - IsBirthday: true if current date is within the user's birthday period
//   - MemberSince: date when the user joined/registered (for tenure-based coupons)
//
// Usage patterns:
//   - First-time buyer coupons check IsFirstPurchase
//   - VIP coupons check LoyaltyTier against required levels
//   - Birthday promotions check IsBirthday flag
//   - Long-term member rewards check MemberSince duration
//
// Example:
//
//	eligibility := UserEligibility{
//		IsFirstPurchase: false,
//		LoyaltyTier: "gold",
//		IsBirthday: true,
//		MemberSince: time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
//	}
type UserEligibility struct {
	IsFirstPurchase bool   `json:"is_first_purchase"`
	LoyaltyTier     string `json:"loyalty_tier"`
	IsBirthday      bool   `json:"is_birthday"`
	MemberSince     time.Time `json:"member_since"`
}
package coupon

import (
	"time"
)

// CouponType represents the type of coupon
type CouponType string

const (
	CouponTypePercentage  CouponType = "percentage"
	CouponTypeFixedAmount CouponType = "fixed_amount"
	CouponTypeBuyXGetY    CouponType = "buy_x_get_y"
	CouponTypeFreeShipping CouponType = "free_shipping"
)

// Coupon represents a coupon with all its properties
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

// CouponUsage represents coupon usage tracking
type CouponUsage struct {
	CouponCode string `json:"coupon_code"`
	UserID     string `json:"user_id"`
	UsageCount int    `json:"usage_count"`
	TotalUsage int    `json:"total_usage"`
}

// CalculationInput represents input for coupon calculation
type CalculationInput struct {
	Coupon      Coupon  `json:"coupon"`
	OrderAmount float64 `json:"order_amount"`
	UserID      string  `json:"user_id"`
	Items       []Item  `json:"items"`
	Usage       CouponUsage `json:"usage"`
}

// Item represents an order item
type Item struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
}

// CalculationResult represents the result of coupon calculation
type CalculationResult struct {
	DiscountAmount float64 `json:"discount_amount"`
	IsValid        bool    `json:"is_valid"`
	ErrorMessage   string  `json:"error_message,omitempty"`
	AppliedItems   []Item  `json:"applied_items,omitempty"` // Items the coupon was applied to
}

// GeneratorConfig represents configuration for coupon code generation
type GeneratorConfig struct {
	Pattern    string `json:"pattern"`    // e.g., "PREFIX-XXXXXX", "XXXXXXXX", "WORD-NUMBER"
	Length     int    `json:"length"`     // Length of random part
	Prefix     string `json:"prefix"`     // Prefix for the code
	Suffix     string `json:"suffix"`     // Suffix for the code
	ExcludeChars string `json:"exclude_chars"` // Characters to exclude (default: "0O1I")
	Count      int    `json:"count"`      // Number of codes to generate
}

// ValidationRule represents a validation rule for coupons
type ValidationRule struct {
	Type        string      `json:"type"`        // "user_based", "order_based", "time_based", "usage_based"
	Condition   string      `json:"condition"`   // Specific condition
	Value       interface{} `json:"value"`       // Value to check against
	ErrorMessage string     `json:"error_message"`
}

// UserEligibility represents user eligibility criteria
type UserEligibility struct {
	IsFirstPurchase bool   `json:"is_first_purchase"`
	LoyaltyTier     string `json:"loyalty_tier"`
	IsBirthday      bool   `json:"is_birthday"`
	MemberSince     time.Time `json:"member_since"`
}
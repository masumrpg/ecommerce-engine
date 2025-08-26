package discount

import "time"

// DiscountType represents the type of discount
type DiscountType string

const (
	DiscountTypeBulk     DiscountType = "bulk"
	DiscountTypeBundle   DiscountType = "bundle"
	DiscountTypeLoyalty  DiscountType = "loyalty"
	DiscountTypeCategory DiscountType = "category"
	DiscountTypeTier     DiscountType = "tier"
	DiscountTypeProgressive DiscountType = "progressive"
)

// BulkDiscountRule represents bulk discount configuration
type BulkDiscountRule struct {
	MinQuantity    int     `json:"min_quantity"`
	MaxQuantity    int     `json:"max_quantity,omitempty"` // 0 means no max
	DiscountType   string  `json:"discount_type"`   // "percentage" or "fixed_amount" or "fixed_price"
	DiscountValue  float64 `json:"discount_value"`
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
	ApplicableProducts   []string `json:"applicable_products,omitempty"`
}

// TierPricingRule represents tier-based pricing
type TierPricingRule struct {
	MinQuantity int     `json:"min_quantity"`
	MaxQuantity int     `json:"max_quantity,omitempty"`
	PricePerItem float64 `json:"price_per_item"`
	Category    string  `json:"category,omitempty"`
}

// BundleDiscountRule represents bundle discount configuration
type BundleDiscountRule struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	RequiredProducts []string `json:"required_products,omitempty"`
	RequiredCategories []string `json:"required_categories,omitempty"`
	MinItems        int      `json:"min_items"`
	DiscountType    string   `json:"discount_type"` // "percentage", "fixed_amount", "combo_price"
	DiscountValue   float64  `json:"discount_value"`
	MaxApplications int      `json:"max_applications,omitempty"` // How many times this bundle can be applied
}

// LoyaltyDiscountRule represents loyalty-based discount
type LoyaltyDiscountRule struct {
	Tier            string  `json:"tier"`            // "bronze", "silver", "gold", "platinum"
	DiscountPercent float64 `json:"discount_percent"`
	MinOrderAmount  float64 `json:"min_order_amount,omitempty"`
	MaxDiscountAmount float64 `json:"max_discount_amount,omitempty"`
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
}

// ProgressiveDiscountRule represents progressive discount (e.g., every 10 items = extra 2%)
type ProgressiveDiscountRule struct {
	QuantityStep    int     `json:"quantity_step"`    // Every X items
	DiscountPercent float64 `json:"discount_percent"` // Additional discount percent
	MaxDiscount     float64 `json:"max_discount"`     // Maximum total discount
	Category        string  `json:"category,omitempty"`
}

// CategoryDiscountRule represents category-specific discount
type CategoryDiscountRule struct {
	Category        string  `json:"category"`
	DiscountPercent float64 `json:"discount_percent"`
	MinQuantity     int     `json:"min_quantity,omitempty"`
	MaxDiscountAmount float64 `json:"max_discount_amount,omitempty"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
}

// DiscountItem represents an item for discount calculation
type DiscountItem struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
	Weight   float64 `json:"weight,omitempty"`
	IsSale   bool    `json:"is_sale,omitempty"`
}

// Customer represents customer information for discount calculation
type Customer struct {
	ID              string    `json:"id"`
	LoyaltyTier     string    `json:"loyalty_tier"`
	MemberSince     time.Time `json:"member_since"`
	TotalPurchases  float64   `json:"total_purchases"`
	PurchaseCount   int       `json:"purchase_count"`
	IsRepeatCustomer bool     `json:"is_repeat_customer"`
}

// DiscountCalculationInput represents input for discount calculation
type DiscountCalculationInput struct {
	Items                   []DiscountItem          `json:"items"`
	Customer                Customer                `json:"customer"`
	BulkRules              []BulkDiscountRule      `json:"bulk_rules,omitempty"`
	TierRules              []TierPricingRule       `json:"tier_rules,omitempty"`
	BundleRules            []BundleDiscountRule    `json:"bundle_rules,omitempty"`
	LoyaltyRules           []LoyaltyDiscountRule   `json:"loyalty_rules,omitempty"`
	ProgressiveRules       []ProgressiveDiscountRule `json:"progressive_rules,omitempty"`
	CategoryRules          []CategoryDiscountRule  `json:"category_rules,omitempty"`
	AllowStacking          bool                    `json:"allow_stacking"`
	MaxStackedDiscountPercent float64             `json:"max_stacked_discount_percent,omitempty"`
}

// DiscountApplication represents a single discount application
type DiscountApplication struct {
	Type           DiscountType `json:"type"`
	RuleID         string       `json:"rule_id"`
	Name           string       `json:"name"`
	DiscountAmount float64      `json:"discount_amount"`
	AppliedItems   []DiscountItem `json:"applied_items"`
	Description    string       `json:"description"`
}

// DiscountCalculationResult represents the result of discount calculation
type DiscountCalculationResult struct {
	OriginalAmount    float64               `json:"original_amount"`
	TotalDiscount     float64               `json:"total_discount"`
	FinalAmount       float64               `json:"final_amount"`
	AppliedDiscounts  []DiscountApplication `json:"applied_discounts"`
	SavingsPercent    float64               `json:"savings_percent"`
	IsValid           bool                  `json:"is_valid"`
	ErrorMessage      string                `json:"error_message,omitempty"`
}

// BundleMatch represents a matched bundle
type BundleMatch struct {
	Rule         BundleDiscountRule `json:"rule"`
	MatchedItems []DiscountItem     `json:"matched_items"`
	Applications int                `json:"applications"`
}

// FrequencyDiscountRule represents purchase frequency-based discount
type FrequencyDiscountRule struct {
	MinPurchaseCount int     `json:"min_purchase_count"`
	DiscountPercent  float64 `json:"discount_percent"`
	ValidDays        int     `json:"valid_days"` // Valid for X days after qualifying
}

// SeasonalDiscountRule represents seasonal discount
type SeasonalDiscountRule struct {
	Season          string    `json:"season"` // "spring", "summer", "autumn", "winter"
	DiscountPercent float64   `json:"discount_percent"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
	Categories      []string  `json:"categories,omitempty"`
	Multiplier      float64   `json:"multiplier,omitempty"` // For member bonus events
}

// CrossSellRule represents cross-sell discount (main product + accessory)
type CrossSellRule struct {
	MainProductCategories []string `json:"main_product_categories"`
	AccessoryCategories   []string `json:"accessory_categories"`
	DiscountPercent       float64  `json:"discount_percent"`
	ComboPrice            float64  `json:"combo_price,omitempty"`
	MinMainProductPrice   float64  `json:"min_main_product_price,omitempty"`
}

// MixAndMatchRule represents mix and match discount
type MixAndMatchRule struct {
	Categories      []string `json:"categories"`
	RequiredItems   int      `json:"required_items"`
	DiscountType    string   `json:"discount_type"` // "flat_discount", "percentage"
	DiscountValue   float64  `json:"discount_value"`
	MaxApplications int      `json:"max_applications,omitempty"`
}
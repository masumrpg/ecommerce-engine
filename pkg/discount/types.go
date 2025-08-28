// Package discount provides comprehensive types and structures for implementing
// flexible discount systems in e-commerce applications.
//
// This package defines all the core types used throughout the discount system,
// including various discount rule types, calculation inputs and results, and
// customer information structures.
//
// Key Features:
//   - Multiple discount types (bulk, bundle, loyalty, category, tier, progressive)
//   - Flexible rule configuration with validation
//   - Comprehensive calculation input/output structures
//   - Customer loyalty tier support
//   - Time-based and seasonal discount rules
//   - Cross-sell and mix-and-match promotions
//   - Bundle matching and application tracking
//
// Basic Usage:
//   // Define discount items
//   items := []DiscountItem{
//       {ID: "item1", Price: 100.0, Quantity: 2, Category: "electronics"},
//   }
//
//   // Create customer information
//   customer := Customer{
//       ID: "customer1",
//       LoyaltyTier: "gold",
//       TotalPurchases: 5000.0,
//   }
//
//   // Configure discount calculation input
//   input := DiscountCalculationInput{
//       Items: items,
//       Customer: customer,
//       AllowStacking: true,
//   }
package discount

import "time"

// DiscountType represents the type of discount applied to items.
// Used to categorize and identify different discount mechanisms
// for reporting and calculation purposes.
type DiscountType string

const (
	// DiscountTypeBulk represents quantity-based discounts
	// Applied when purchasing large quantities of items
	DiscountTypeBulk DiscountType = "bulk"

	// DiscountTypeBundle represents bundle discounts
	// Applied when purchasing specific product combinations
	DiscountTypeBundle DiscountType = "bundle"

	// DiscountTypeLoyalty represents customer loyalty discounts
	// Applied based on customer loyalty tier and history
	DiscountTypeLoyalty DiscountType = "loyalty"

	// DiscountTypeCategory represents category-specific discounts
	// Applied to items within specific product categories
	DiscountTypeCategory DiscountType = "category"

	// DiscountTypeTier represents tier-based pricing discounts
	// Applied based on quantity tiers with different price points
	DiscountTypeTier DiscountType = "tier"

	// DiscountTypeProgressive represents progressive discounts
	// Applied with increasing discount rates based on quantity
	DiscountTypeProgressive DiscountType = "progressive"
)

// BulkDiscountRule represents bulk discount configuration.
// Defines quantity-based discounts that apply when customers purchase
// large quantities of items, encouraging bulk purchases.
//
// Features:
//   - Minimum and maximum quantity thresholds
//   - Multiple discount types (percentage, fixed amount, fixed price)
//   - Category and product-specific targeting
//   - Flexible quantity range configuration
//
// Example:
//   rule := BulkDiscountRule{
//       MinQuantity: 10,
//       DiscountType: "percentage",
//       DiscountValue: 15.0, // 15% off
//       ApplicableCategories: []string{"electronics"},
//   }
type BulkDiscountRule struct {
	MinQuantity    int     `json:"min_quantity"`
	MaxQuantity    int     `json:"max_quantity,omitempty"` // 0 means no max
	DiscountType   string  `json:"discount_type"`   // "percentage" or "fixed_amount" or "fixed_price"
	DiscountValue  float64 `json:"discount_value"`
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
	ApplicableProducts   []string `json:"applicable_products,omitempty"`
}

// TierPricingRule represents tier-based pricing configuration.
// Defines different price points based on quantity tiers,
// allowing for volume-based pricing strategies.
//
// Features:
//   - Quantity-based tier thresholds
//   - Fixed price per item in each tier
//   - Category-specific tier pricing
//   - Scalable pricing structure
//
// Example:
//   rule := TierPricingRule{
//       MinQuantity: 50,
//       MaxQuantity: 99,
//       PricePerItem: 8.50, // Reduced price for this tier
//       Category: "office-supplies",
//   }
type TierPricingRule struct {
	MinQuantity int     `json:"min_quantity"`
	MaxQuantity int     `json:"max_quantity,omitempty"`
	PricePerItem float64 `json:"price_per_item"`
	Category    string  `json:"category,omitempty"`
}

// BundleDiscountRule represents bundle discount configuration.
// Defines discounts for purchasing specific product combinations,
// encouraging customers to buy complementary items together.
//
// Features:
//   - Product or category-based bundle requirements
//   - Minimum item count validation
//   - Multiple discount types (percentage, fixed amount, combo price)
//   - Maximum application limits
//   - Flexible bundle composition
//
// Example:
//   rule := BundleDiscountRule{
//       ID: "laptop-bundle",
//       Name: "Laptop Starter Bundle",
//       RequiredCategories: []string{"laptops", "accessories"},
//       MinItems: 2,
//       DiscountType: "percentage",
//       DiscountValue: 10.0,
//   }
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

// LoyaltyDiscountRule represents loyalty-based discount configuration.
// Provides discounts based on customer loyalty tiers and purchase history,
// rewarding long-term customers with better pricing.
//
// Features:
//   - Tier-based discount percentages
//   - Minimum order amount requirements
//   - Maximum discount amount caps
//   - Category-specific loyalty benefits
//   - Flexible tier system support
//
// Example:
//   rule := LoyaltyDiscountRule{
//       Tier: "gold",
//       DiscountPercent: 12.0,
//       MinOrderAmount: 100.0,
//       MaxDiscountAmount: 50.0,
//   }
type LoyaltyDiscountRule struct {
	Tier            string  `json:"tier"`            // "bronze", "silver", "gold", "platinum"
	DiscountPercent float64 `json:"discount_percent"`
	MinOrderAmount  float64 `json:"min_order_amount,omitempty"`
	MaxDiscountAmount float64 `json:"max_discount_amount,omitempty"`
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
}

// ProgressiveDiscountRule represents progressive discount configuration.
// Provides increasing discount rates based on quantity milestones,
// encouraging larger purchases with escalating benefits.
//
// Features:
//   - Quantity step-based progression
//   - Cumulative discount percentage increases
//   - Maximum discount caps
//   - Category-specific progressive discounts
//
// Example:
//   rule := ProgressiveDiscountRule{
//       QuantityStep: 10,    // Every 10 items
//       DiscountPercent: 2.0, // Additional 2% discount
//       MaxDiscount: 20.0,   // Maximum 20% total
//       Category: "books",
//   }
type ProgressiveDiscountRule struct {
	QuantityStep    int     `json:"quantity_step"`    // Every X items
	DiscountPercent float64 `json:"discount_percent"` // Additional discount percent
	MaxDiscount     float64 `json:"max_discount"`     // Maximum total discount
	Category        string  `json:"category,omitempty"`
}

// CategoryDiscountRule represents category-specific discount configuration.
// Provides targeted discounts for specific product categories,
// enabling category-based promotional campaigns.
//
// Features:
//   - Category-targeted discounts
//   - Minimum quantity requirements
//   - Maximum discount amount limits
//   - Time-based validity periods
//   - Flexible category targeting
//
// Example:
//   rule := CategoryDiscountRule{
//       Category: "electronics",
//       DiscountPercent: 15.0,
//       MinQuantity: 2,
//       ValidFrom: time.Now(),
//       ValidUntil: time.Now().AddDate(0, 1, 0),
//   }
type CategoryDiscountRule struct {
	Category        string  `json:"category"`
	DiscountPercent float64 `json:"discount_percent"`
	MinQuantity     int     `json:"min_quantity,omitempty"`
	MaxDiscountAmount float64 `json:"max_discount_amount,omitempty"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
}

// DiscountItem represents an item for discount calculation.
// Contains all necessary information about a product item
// required for discount calculations and rule applications.
//
// Features:
//   - Unique item identification
//   - Price and quantity information
//   - Category classification
//   - Optional weight and sale status
//   - Flexible item attributes
//
// Example:
//   item := DiscountItem{
//       ID: "laptop-001",
//       Price: 999.99,
//       Quantity: 2,
//       Category: "electronics",
//       Weight: 2.5,
//   }
type DiscountItem struct {
	ID       string  `json:"id"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
	Weight   float64 `json:"weight,omitempty"`
	IsSale   bool    `json:"is_sale,omitempty"`
}

// Customer represents customer information for discount calculation.
// Contains customer data necessary for applying loyalty discounts,
// frequency-based promotions, and personalized pricing.
//
// Features:
//   - Unique customer identification
//   - Loyalty tier classification
//   - Purchase history tracking
//   - Membership duration information
//   - Repeat customer identification
//
// Example:
//   customer := Customer{
//       ID: "customer-123",
//       LoyaltyTier: "gold",
//       MemberSince: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
//       TotalPurchases: 5000.0,
//       PurchaseCount: 25,
//       IsRepeatCustomer: true,
//   }
type Customer struct {
	ID              string    `json:"id"`
	LoyaltyTier     string    `json:"loyalty_tier"`
	MemberSince     time.Time `json:"member_since"`
	TotalPurchases  float64   `json:"total_purchases"`
	PurchaseCount   int       `json:"purchase_count"`
	IsRepeatCustomer bool     `json:"is_repeat_customer"`
}

// DiscountCalculationInput represents input for discount calculation.
// Contains all necessary data for performing comprehensive discount calculations,
// including items, customer information, and applicable discount rules.
//
// Features:
//   - Complete item and customer data
//   - Multiple discount rule types
//   - Stacking configuration options
//   - Maximum discount limits
//   - Flexible rule combinations
//
// Example:
//   input := DiscountCalculationInput{
//       Items: []DiscountItem{{ID: "item1", Price: 100.0, Quantity: 2}},
//       Customer: Customer{ID: "customer1", LoyaltyTier: "gold"},
//       AllowStacking: true,
//       MaxStackedDiscountPercent: 50.0,
//   }
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

// DiscountApplication represents a single discount application.
// Records details about an individual discount that was applied
// during the calculation process for tracking and reporting.
//
// Features:
//   - Discount type and rule identification
//   - Applied discount amount tracking
//   - Item-specific application details
//   - Human-readable descriptions
//   - Comprehensive audit trail
//
// Example:
//   application := DiscountApplication{
//       Type: DiscountTypeBulk,
//       RuleID: "bulk-electronics-10",
//       Name: "Electronics Bulk Discount",
//       DiscountAmount: 50.0,
//       Description: "15% off for purchasing 10+ electronics items",
//   }
type DiscountApplication struct {
	Type           DiscountType `json:"type"`
	RuleID         string       `json:"rule_id"`
	Name           string       `json:"name"`
	DiscountAmount float64      `json:"discount_amount"`
	AppliedItems   []DiscountItem `json:"applied_items"`
	Description    string       `json:"description"`
}

// DiscountCalculationResult represents the result of discount calculation.
// Contains comprehensive information about the discount calculation outcome,
// including original amounts, applied discounts, and final pricing.
//
// Features:
//   - Original and final amount tracking
//   - Total discount calculation
//   - Applied discount details
//   - Savings percentage calculation
//   - Validation status and error handling
//
// Example:
//   result := DiscountCalculationResult{
//       OriginalAmount: 200.0,
//       TotalDiscount: 30.0,
//       FinalAmount: 170.0,
//       SavingsPercent: 15.0,
//       IsValid: true,
//   }
type DiscountCalculationResult struct {
	OriginalAmount    float64               `json:"original_amount"`
	TotalDiscount     float64               `json:"total_discount"`
	FinalAmount       float64               `json:"final_amount"`
	AppliedDiscounts  []DiscountApplication `json:"applied_discounts"`
	SavingsPercent    float64               `json:"savings_percent"`
	IsValid           bool                  `json:"is_valid"`
	ErrorMessage      string                `json:"error_message,omitempty"`
}

// BundleMatch represents a matched bundle configuration.
// Tracks successful bundle matches during discount calculation,
// including the matched items and application count.
//
// Features:
//   - Bundle rule reference
//   - Matched item tracking
//   - Application count monitoring
//   - Bundle optimization support
//
// Example:
//   match := BundleMatch{
//       Rule: bundleRule,
//       MatchedItems: []DiscountItem{laptop, mouse, keyboard},
//       Applications: 1,
//   }
type BundleMatch struct {
	Rule         BundleDiscountRule `json:"rule"`
	MatchedItems []DiscountItem     `json:"matched_items"`
	Applications int                `json:"applications"`
}

// FrequencyDiscountRule represents purchase frequency-based discount configuration.
// Provides discounts based on customer purchase frequency and patterns,
// encouraging repeat purchases and customer retention.
//
// Features:
//   - Minimum purchase count requirements
//   - Frequency-based discount percentages
//   - Time-limited validity periods
//   - Customer behavior tracking
//
// Example:
//   rule := FrequencyDiscountRule{
//       MinPurchaseCount: 5,  // After 5 purchases
//       DiscountPercent: 10.0, // 10% discount
//       ValidDays: 30,        // Valid for 30 days
//   }
type FrequencyDiscountRule struct {
	MinPurchaseCount int     `json:"min_purchase_count"`
	DiscountPercent  float64 `json:"discount_percent"`
	ValidDays        int     `json:"valid_days"` // Valid for X days after qualifying
}

// SeasonalDiscountRule represents seasonal discount configuration.
// Provides time-based discounts for specific seasons or periods,
// enabling seasonal promotional campaigns and holiday sales.
//
// Features:
//   - Season-based discount targeting
//   - Time period validation
//   - Category-specific seasonal discounts
//   - Member bonus multipliers
//   - Flexible seasonal definitions
//
// Example:
//   rule := SeasonalDiscountRule{
//       Season: "winter",
//       DiscountPercent: 20.0,
//       ValidFrom: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
//       ValidUntil: time.Date(2025, 2, 28, 23, 59, 59, 0, time.UTC),
//       Categories: []string{"clothing", "accessories"},
//   }
type SeasonalDiscountRule struct {
	Season          string    `json:"season"` // "spring", "summer", "autumn", "winter"
	DiscountPercent float64   `json:"discount_percent"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
	Categories      []string  `json:"categories,omitempty"`
	Multiplier      float64   `json:"multiplier,omitempty"` // For member bonus events
}

// CrossSellRule represents cross-sell discount configuration.
// Encourages customers to purchase complementary products by offering
// discounts when main products are combined with accessories.
//
// Features:
//   - Main product and accessory category matching
//   - Percentage or combo price discounts
//   - Minimum main product price requirements
//   - Automatic product combination detection
//   - Flexible category relationships
//
// Example:
//   rule := CrossSellRule{
//       MainProductCategories: []string{"laptops"},
//       AccessoryCategories: []string{"mice", "keyboards"},
//       DiscountPercent: 15.0,
//       MinMainProductPrice: 500.0,
//   }
type CrossSellRule struct {
	MainProductCategories []string `json:"main_product_categories"`
	AccessoryCategories   []string `json:"accessory_categories"`
	DiscountPercent       float64  `json:"discount_percent"`
	ComboPrice            float64  `json:"combo_price,omitempty"`
	MinMainProductPrice   float64  `json:"min_main_product_price,omitempty"`
}

// MixAndMatchRule represents mix and match discount configuration.
// Provides discounts when customers purchase a specified number of items
// from designated categories, promoting variety in purchases.
//
// Features:
//   - Multi-category item matching
//   - Required item count validation
//   - Flat discount or percentage options
//   - Maximum application limits
//   - Flexible category combinations
//
// Example:
//   rule := MixAndMatchRule{
//       Categories: []string{"shirts", "pants", "shoes"},
//       RequiredItems: 3,
//       DiscountType: "flat_discount",
//       DiscountValue: 25.0,
//       MaxApplications: 2,
//   }
type MixAndMatchRule struct {
	Categories      []string `json:"categories"`
	RequiredItems   int      `json:"required_items"`
	DiscountType    string   `json:"discount_type"` // "flat_discount", "percentage"
	DiscountValue   float64  `json:"discount_value"`
	MaxApplications int      `json:"max_applications,omitempty"`
}
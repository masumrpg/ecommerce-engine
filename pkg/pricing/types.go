// Package pricing provides comprehensive types and structures for flexible pricing calculations.
// This package defines all the core data types used throughout the pricing system, including
// pricing strategies, rules, bundles, tiers, and dynamic pricing configurations.
//
// Key Features:
//   - Multiple pricing strategies (fixed, dynamic, tiered, volume, bundle, etc.)
//   - Flexible pricing rules with conditions and adjustments
//   - Bundle pricing with various bundle types
//   - Tier-based volume pricing
//   - Dynamic pricing with market factors
//   - Comprehensive pricing analytics
//   - Customer segmentation support
//   - Multi-channel and multi-region pricing
//
// Basic Usage:
//
//	// Define a pricing item
//	item := pricing.PricingItem{
//		ID: "item-001",
//		Name: "Premium Widget",
//		Category: "electronics",
//		Quantity: 5,
//		BasePrice: 99.99,
//	}
//	
//	// Define customer information
//	customer := pricing.Customer{
//		ID: "customer-123",
//		Type: "premium",
//		Segment: "vip",
//		TotalSpent: 5000.00,
//	}
//	
//	// Create pricing context
//	context := pricing.PricingContext{
//		Channel: "online",
//		Region: "US",
//		Currency: "USD",
//		Timestamp: time.Now(),
//	}
//	
//	// Define pricing input
//	input := pricing.PricingInput{
//		Items: []pricing.PricingItem{item},
//		Customer: customer,
//		Context: context,
//	}
package pricing

import (
	"time"
)

// PricingStrategy represents different pricing strategies available in the system.
// Each strategy defines a different approach to calculating prices based on various factors.
type PricingStrategy string

const (
	StrategyFixed        PricingStrategy = "fixed"        // Fixed pricing
	StrategyDynamic      PricingStrategy = "dynamic"      // Dynamic pricing
	StrategyTiered       PricingStrategy = "tiered"       // Tiered pricing
	StrategyVolume       PricingStrategy = "volume"       // Volume-based pricing
	StrategyBundle       PricingStrategy = "bundle"       // Bundle pricing
	StrategySubscription PricingStrategy = "subscription" // Subscription pricing
	StrategyAuction      PricingStrategy = "auction"      // Auction pricing
	StrategyNegotiated   PricingStrategy = "negotiated"   // Negotiated pricing
	StrategyMarketBased  PricingStrategy = "market_based" // Market-based pricing
	StrategyCompetitive  PricingStrategy = "competitive"  // Competitive pricing
)

// PricingType represents the type of pricing calculation being performed.
// This determines the context and purpose of the price calculation.
type PricingType string

const (
	PricingTypeBase      PricingType = "base"      // Base price
	PricingTypePromo     PricingType = "promo"     // Promotional price
	PricingTypeClearance PricingType = "clearance" // Clearance price
	PricingTypeMSRP      PricingType = "msrp"      // Manufacturer's Suggested Retail Price
	PricingTypeWholesale PricingType = "wholesale" // Wholesale price
	PricingTypeRetail    PricingType = "retail"    // Retail price
	PricingTypeContract  PricingType = "contract"  // Contract price
)

// BundleType represents different types of product bundles available.
// Each bundle type has different rules for item selection and pricing.
type BundleType string

const (
	BundleTypeFixed      BundleType = "fixed"      // Fixed bundle with set items
	BundleTypeMixMatch   BundleType = "mix_match"  // Mix and match bundle
	BundleTypeFrequency  BundleType = "frequency"  // Frequency-based bundle
	BundleTypeCrossSell  BundleType = "cross_sell" // Cross-sell bundle
	BundleTypeUpSell     BundleType = "up_sell"    // Up-sell bundle
	BundleTypeCombo      BundleType = "combo"      // Combo bundle
	BundleTypeKit        BundleType = "kit"        // Kit bundle
	BundleTypeSubscription BundleType = "subscription" // Subscription bundle
)

// PricingRule represents a comprehensive pricing rule that can be applied to items.
// Rules define conditions under which specific price adjustments should be made,
// supporting complex business logic for dynamic pricing strategies.
//
// Key features:
//   - Conditional logic with multiple operators
//   - Multiple adjustment types (percentage, fixed, markup, markdown)
//   - Customer segmentation and channel targeting
//   - Priority-based rule application
//   - Time-based validity periods
//   - Item inclusion/exclusion lists
//
// Example:
//
//	// Volume discount rule
//	rule := PricingRule{
//		ID: "volume-discount-10",
//		Name: "10+ Items Volume Discount",
//		Strategy: StrategyVolume,
//		Type: PricingTypePromo,
//		Priority: 100,
//		Conditions: []PricingCondition{
//			{Type: "quantity", Operator: ">=", Value: 10},
//		},
//		Adjustments: []PriceAdjustment{
//			{Type: "percentage", Value: 10.0},
//		},
//		CustomerSegments: []string{"premium", "vip"},
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(0, 3, 0), // 3 months
//	}
type PricingRule struct {
	ID               string          `json:"id"`
	Name             string          `json:"name"`
	Description      string          `json:"description,omitempty"`
	Strategy         PricingStrategy `json:"strategy"`
	Type             PricingType     `json:"type"`
	Priority         int             `json:"priority"`
	IsActive         bool            `json:"is_active"`
	ValidFrom        time.Time       `json:"valid_from"`
	ValidUntil       time.Time       `json:"valid_until"`
	Conditions       []PricingCondition `json:"conditions,omitempty"`
	Adjustments      []PriceAdjustment  `json:"adjustments,omitempty"`
	ApplicableItems  []string        `json:"applicable_items,omitempty"`
	ExcludedItems    []string        `json:"excluded_items,omitempty"`
	CustomerSegments []string        `json:"customer_segments,omitempty"`
	Channels         []string        `json:"channels,omitempty"`
	Regions          []string        `json:"regions,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// PricingCondition represents a condition that must be met for a pricing rule to apply.
// Conditions support various comparison operators and can be combined with logical operators.
//
// Supported condition types:
//   - "quantity": Item quantity comparison
//   - "amount": Price amount comparison
//   - "customer_type": Customer type matching
//   - "time": Time-based conditions
//   - "inventory": Inventory level conditions
//   - "category": Product category matching
//   - "brand": Product brand matching
//
// Supported operators:
//   - ">", "<", ">=", "<=": Numeric comparisons
//   - "=", "!=": Equality comparisons
//   - "in": Value in list
//   - "between": Value between two values
//
// Example:
//
//	// Quantity condition
//	condition := PricingCondition{
//		Type: "quantity",
//		Operator: ">=",
//		Value: 10,
//		Logic: "AND",
//	}
type PricingCondition struct {
	Type     string      `json:"type"`     // "quantity", "amount", "customer_type", "time", "inventory"
	Operator string      `json:"operator"` // ">", "<", ">=", "<=", "=", "!=", "in", "between"
	Value    interface{} `json:"value"`    // Condition value
	Logic    string      `json:"logic,omitempty"` // "AND", "OR"
}

// PriceAdjustment represents a price adjustment that can be applied to an item.
// Supports various adjustment types with optional price limits and rounding.
//
// Supported adjustment types:
//   - "percentage": Percentage-based adjustment (positive or negative)
//   - "fixed": Fixed amount adjustment (positive or negative)
//   - "markup": Markup percentage (always positive)
//   - "markdown": Markdown percentage (always negative)
//
// Example:
//
//	// 15% discount with price limits
//	adjustment := PriceAdjustment{
//		Type: "percentage",
//		Value: -15.0, // 15% discount
//		MinPrice: 10.00, // Don't go below $10
//		MaxPrice: 500.00, // Don't go above $500
//		RoundTo: 0.99, // Round to .99 endings
//		Description: "15% Volume Discount",
//	}
type PriceAdjustment struct {
	Type        string  `json:"type"`        // "percentage", "fixed", "markup", "markdown"
	Value       float64 `json:"value"`       // Adjustment value
	MinPrice    float64 `json:"min_price,omitempty"`    // Minimum price limit
	MaxPrice    float64 `json:"max_price,omitempty"`    // Maximum price limit
	RoundTo     float64 `json:"round_to,omitempty"`     // Round to nearest value
	Description string  `json:"description,omitempty"`
}

// TierPricing represents a tiered pricing structure for volume-based discounts.
// Enables different pricing based on quantity thresholds, encouraging bulk purchases.
//
// Example:
//
//	// Bulk pricing tiers
//	tierPricing := TierPricing{
//		ID: "bulk-discount",
//		Name: "Bulk Purchase Discount",
//		Tiers: []PriceTier{
//			{MinQuantity: 1, MaxQuantity: 9, Price: 10.00},
//			{MinQuantity: 10, MaxQuantity: 49, Discount: 5.0}, // 5% off
//			{MinQuantity: 50, MaxQuantity: 99, Discount: 10.0}, // 10% off
//			{MinQuantity: 100, Discount: 15.0}, // 15% off for 100+
//		},
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(1, 0, 0), // 1 year
//	}
type TierPricing struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Tiers       []PriceTier `json:"tiers"`
	IsActive    bool        `json:"is_active"`
	ValidFrom   time.Time   `json:"valid_from"`
	ValidUntil  time.Time   `json:"valid_until"`
}

// PriceTier represents a single tier in a tiered pricing structure.
// Defines quantity ranges and associated pricing or discounts.
//
// Either Price, Discount, or FixedPrice should be specified:
//   - Price: Specific price for this tier
//   - Discount: Percentage discount from base price
//   - FixedPrice: Fixed price regardless of base price
//
// Example:
//
//	// Tier with percentage discount
//	tier := PriceTier{
//		MinQuantity: 10,
//		MaxQuantity: 49,
//		Discount: 10.0, // 10% discount
//	}
//	
//	// Tier with fixed price
//	tier2 := PriceTier{
//		MinQuantity: 50,
//		FixedPrice: 8.50, // Fixed price per unit
//	}
type PriceTier struct {
	MinQuantity int     `json:"min_quantity"`
	MaxQuantity int     `json:"max_quantity,omitempty"` // 0 means no upper limit
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount,omitempty"`    // Discount percentage
	FixedPrice  float64 `json:"fixed_price,omitempty"` // Fixed price for this tier
}

// Bundle represents a product bundle configuration for cross-sell and upsell opportunities.
// Bundles group related items together with special pricing to increase average order value.
//
// Supported bundle types:
//   - "fixed": Fixed bundle with predetermined items
//   - "mix_match": Customer can choose from available items
//   - "frequency": Recurring bundle for subscriptions
//   - "cross_sell": Related items bundle
//   - "up_sell": Higher-value alternative bundle
//
// Example:
//
//	// Laptop accessories bundle
//	bundle := Bundle{
//		ID: "laptop-bundle",
//		Name: "Laptop Essentials Bundle",
//		Type: BundleTypeFixed,
//		Items: []BundleItem{
//			{ItemID: "laptop-001", Quantity: 1, IsRequired: true},
//			{ItemID: "mouse-001", Quantity: 1, IsOptional: true},
//			{ItemID: "keyboard-001", Quantity: 1, IsOptional: true},
//		},
//		Pricing: BundlePricing{
//			Type: "percentage",
//			Value: 15.0, // 15% off bundle
//		},
//		MinItems: 2, // At least laptop + one accessory
//		IsActive: true,
//	}
type Bundle struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Type         BundleType  `json:"type"`
	Items        []BundleItem `json:"items"`
	Pricing      BundlePricing `json:"pricing"`
	MinItems     int         `json:"min_items,omitempty"`     // Minimum items required
	MaxItems     int         `json:"max_items,omitempty"`     // Maximum items allowed
	IsActive     bool        `json:"is_active"`
	ValidFrom    time.Time   `json:"valid_from"`
	ValidUntil   time.Time   `json:"valid_until"`
	Conditions   []PricingCondition `json:"conditions,omitempty"`
	Tags         []string    `json:"tags,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// BundleItem represents an individual item within a product bundle.
// Defines the item details, quantity, pricing, and whether it's required or optional.
//
// Example:
//
//	// Required bundle item
//	item := BundleItem{
//		ItemID: "laptop-001",
//		Name: "Premium Laptop",
//		Quantity: 1,
//		IsRequired: true,
//		BasePrice: 999.99,
//		BundlePrice: 899.99, // Special bundle price
//		Category: "computers",
//	}
//	
//	// Optional bundle item
//	optionalItem := BundleItem{
//		ItemID: "mouse-001",
//		Name: "Wireless Mouse",
//		Quantity: 1,
//		IsOptional: true,
//		BasePrice: 49.99,
//		Discount: 20.0, // 20% off when in bundle
//	}
type BundleItem struct {
	ItemID       string  `json:"item_id"`
	Name         string  `json:"name"`
	Quantity     int     `json:"quantity"`
	IsRequired   bool    `json:"is_required"`
	IsOptional   bool    `json:"is_optional"`
	BasePrice    float64 `json:"base_price"`
	BundlePrice  float64 `json:"bundle_price,omitempty"`
	Discount     float64 `json:"discount,omitempty"`
	Category     string  `json:"category,omitempty"`
	Subcategory  string  `json:"subcategory,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

// BundlePricing represents the pricing configuration for a product bundle.
// Defines how the bundle price is calculated and any savings offered.
//
// Supported pricing types:
//   - "fixed": Fixed bundle price regardless of individual item prices
//   - "percentage": Percentage discount off total individual prices
//   - "tiered": Different pricing based on quantity or value tiers
//   - "dynamic": Dynamic pricing based on market conditions
//
// Example:
//
//	// Percentage-based bundle pricing
//	pricing := BundlePricing{
//		Type: "percentage",
//		Value: 20.0, // 20% off individual prices
//		MinPrice: 100.00, // Minimum bundle price
//		MaxPrice: 1000.00, // Maximum bundle price
//		SavingsType: "percentage",
//		SavingsValue: 20.0,
//	}
//	
//	// Fixed bundle pricing
//	fixedPricing := BundlePricing{
//		Type: "fixed",
//		Value: 299.99, // Fixed bundle price
//		BasePrice: 350.00, // Original total price
//		SavingsType: "amount",
//		SavingsValue: 50.01, // $50.01 savings
//	}
type BundlePricing struct {
	Type         string  `json:"type"`         // "fixed", "percentage", "tiered", "dynamic"
	Value        float64 `json:"value"`        // Price or discount value
	MinPrice     float64 `json:"min_price,omitempty"`     // Minimum bundle price
	MaxPrice     float64 `json:"max_price,omitempty"`     // Maximum bundle price
	BasePrice    float64 `json:"base_price,omitempty"`    // Base bundle price
	SavingsType  string  `json:"savings_type,omitempty"`  // "amount", "percentage"
	SavingsValue float64 `json:"savings_value,omitempty"` // Savings amount or percentage
}

// PricingItem represents an item that needs pricing calculation.
// Contains all necessary information for applying pricing rules, bundles, and tiers.
//
// Example:
//
//	// Standard product item
//	item := PricingItem{
//		ID: "widget-001",
//		Name: "Premium Widget",
//		SKU: "WDG-001-PRM",
//		Category: "electronics",
//		Subcategory: "widgets",
//		Brand: "TechCorp",
//		Quantity: 5,
//		BasePrice: 99.99,
//		CostPrice: 60.00,
//		MSRP: 129.99,
//		Weight: 1.5,
//		Dimensions: Dimensions{
//			Length: 10.0,
//			Width: 5.0,
//			Height: 3.0,
//			Unit: "cm",
//		},
//		InventoryLevel: 150,
//		Tags: []string{"premium", "bestseller"},
//	}
type PricingItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	SKU          string  `json:"sku,omitempty"`
	Category     string  `json:"category"`
	Subcategory  string  `json:"subcategory,omitempty"`
	Brand        string  `json:"brand,omitempty"`
	Quantity     int     `json:"quantity"`
	BasePrice    float64 `json:"base_price"`
	CostPrice    float64 `json:"cost_price,omitempty"`
	MSRP         float64 `json:"msrp,omitempty"`
	Weight       float64 `json:"weight,omitempty"`
	Dimensions   Dimensions `json:"dimensions,omitempty"`
	InventoryLevel int   `json:"inventory_level,omitempty"`
	IsDigital    bool    `json:"is_digital,omitempty"`
	IsSubscription bool  `json:"is_subscription,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

// Dimensions represents the physical dimensions of an item.
// Used for shipping calculations and bundle compatibility checks.
//
// Example:
//
//	// Item dimensions in centimeters
//	dims := Dimensions{
//		Length: 25.4,
//		Width: 15.2,
//		Height: 8.9,
//		Unit: "cm",
//	}
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"` // "cm", "in", "mm"
}

// Customer represents customer information used for pricing calculations.
// Contains customer details, preferences, and history for personalized pricing.
//
// Example:
//
//	// Premium customer with loyalty tier
//	customer := Customer{
//		ID: "cust-12345",
//		Type: "premium",
//		Tier: "gold",
//		Segment: "enterprise",
//		Location: "US-CA",
//		Currency: "USD",
//		TotalSpent: 15000.00,
//		OrderCount: 45,
//		AverageOrderValue: 333.33,
//		LoyaltyPoints: 2500,
//		Tags: []string{"vip", "bulk-buyer"},
//		Preferences: map[string]interface{}{
//			"preferred_brands": []string{"TechCorp", "InnovateCo"},
//			"price_sensitivity": "low",
//		},
//	}
type Customer struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`         // "individual", "business", "vip", "wholesale"
	Segment      string    `json:"segment,omitempty"`      // Customer segment
	Tier         string    `json:"tier,omitempty"`         // Customer tier
	LoyaltyLevel string    `json:"loyalty_level,omitempty"` // Loyalty level
	JoinDate     time.Time `json:"join_date,omitempty"`
	TotalSpent   float64   `json:"total_spent,omitempty"`
	OrderCount   int       `json:"order_count,omitempty"`
	Region       string    `json:"region,omitempty"`
	Channel      string    `json:"channel,omitempty"`      // "online", "retail", "mobile", "api"
	Attributes   map[string]interface{} `json:"attributes,omitempty"`
}

// PricingContext represents the contextual information for pricing calculations.
// Includes timing, location, channel, and market conditions that affect pricing.
//
// Example:
//
//	// E-commerce pricing context
//	ctx := PricingContext{
//		Channel: "online",
//		Location: "US-NY",
//		Currency: "USD",
//		Timezone: "America/New_York",
//		Season: "holiday",
//		Campaign: "black-friday-2024",
//		UserAgent: "mobile-app",
//		Referrer: "google-ads",
//		SessionID: "sess-abc123",
//		MarketConditions: map[string]interface{}{
//			"demand_level": "high",
//			"inventory_pressure": "low",
//			"competitor_activity": "aggressive",
//		},
//	}
type PricingContext struct {
	Channel      string    `json:"channel"`      // Sales channel
	Region       string    `json:"region"`       // Geographic region
	Currency     string    `json:"currency"`     // Currency code
	ExchangeRate float64   `json:"exchange_rate,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
	Season       string    `json:"season,omitempty"`       // "spring", "summer", "fall", "winter"
	Event        string    `json:"event,omitempty"`        // Special events
	InventoryData map[string]int `json:"inventory_data,omitempty"`
	MarketData   map[string]interface{} `json:"market_data,omitempty"`
	CompetitorData map[string]interface{} `json:"competitor_data,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// PricingInput represents the complete input for a pricing calculation request.
// Combines items, customer, context, and options for comprehensive pricing.
//
// Example:
//
//	// Complete pricing input
//	input := PricingInput{
//		Items: []PricingItem{
//			{ID: "widget-001", Quantity: 2, BasePrice: 99.99},
//			{ID: "gadget-002", Quantity: 1, BasePrice: 149.99},
//		},
//		Customer: &Customer{
//			ID: "cust-12345",
//			Type: "premium",
//			Tier: "gold",
//		},
//		Context: &PricingContext{
//			Channel: "online",
//			Location: "US-CA",
//			Currency: "USD",
//		},
//		Options: &PricingOptions{
//			IncludeBundles: true,
//			IncludeTiers: true,
//			IncludeRecommendations: true,
//		},
//	}
type PricingInput struct {
	Items       []PricingItem   `json:"items"`
	Customer    Customer        `json:"customer"`
	Context     PricingContext  `json:"context"`
	Rules       []PricingRule   `json:"rules,omitempty"`
	Bundles     []Bundle        `json:"bundles,omitempty"`
	TierPricing []TierPricing   `json:"tier_pricing,omitempty"`
	Options     PricingOptions  `json:"options,omitempty"`
}

// PricingOptions represents configuration options for pricing calculations.
// Controls which pricing features and calculations to include in the result.
//
// Example:
//
//	// Comprehensive pricing options
//	options := PricingOptions{
//		IncludeTiers: true,
//		IncludeBundles: true,
//		IncludeDynamic: true,
//		IncludeRecommendations: true,
//		MaxRecommendations: 5,
//		RoundingMode: "round", // "floor", "ceil", "round"
//		Currency: "USD",
//		Precision: 2,
//		DebugMode: false,
//	}
//	
//	// Minimal pricing options
//	minimalOptions := PricingOptions{
//		IncludeTiers: false,
//		IncludeBundles: false,
//		IncludeDynamic: false,
//		RoundingMode: "round",
//		Precision: 2,
//	}
type PricingOptions struct {
	IncludeTax       bool    `json:"include_tax,omitempty"`
	IncludeShipping  bool    `json:"include_shipping,omitempty"`
	ApplyDiscounts   bool    `json:"apply_discounts,omitempty"`
	ApplyPromotions  bool    `json:"apply_promotions,omitempty"`
	RoundingMode     string  `json:"rounding_mode,omitempty"`     // "round", "floor", "ceil"
	RoundingPrecision int    `json:"rounding_precision,omitempty"` // Decimal places
	MaxDiscount      float64 `json:"max_discount,omitempty"`      // Maximum discount percentage
	MinMargin        float64 `json:"min_margin,omitempty"`        // Minimum profit margin
	CalculateBundle  bool    `json:"calculate_bundle,omitempty"`
	CalculateTiers   bool    `json:"calculate_tiers,omitempty"`
}

// PricedItem represents the pricing result for an individual item.
// Contains the final price, applied rules, discounts, and calculation details.
//
// Example:
//
//	// Priced item with discounts applied
//	pricedItem := PricedItem{
//		ItemID: "widget-001",
//		Name: "Premium Widget",
//		Quantity: 3,
//		BasePrice: 99.99,
//		FinalPrice: 89.99, // After 10% discount
//		OriginalPrice: 99.99,
//		Savings: 10.00,
//		SavingsPercent: 10.0,
//		TotalPrice: 269.97, // 3 × 89.99
//		AppliedRules: []AppliedPricingRule{
//			{RuleID: "bulk-discount", Adjustment: -10.0},
//		},
//		TierInfo: &TierInfo{
//			TierName: "bulk",
//			MinQuantity: 3,
//			TierDiscount: 10.0,
//		},
//	}
type PricedItem struct {
	ItemID        string            `json:"item_id"`
	Name          string            `json:"name"`
	Quantity      int               `json:"quantity"`
	BasePrice     float64           `json:"base_price"`
	FinalPrice    float64           `json:"final_price"`
	UnitPrice     float64           `json:"unit_price"`
	TotalPrice    float64           `json:"total_price"`
	OriginalPrice float64           `json:"original_price,omitempty"`
	Savings       float64           `json:"savings,omitempty"`
	SavingsPercent float64          `json:"savings_percent,omitempty"`
	AppliedRules  []AppliedPricingRule `json:"applied_rules,omitempty"`
	TierInfo      *TierInfo         `json:"tier_info,omitempty"`
	BundleInfo    *BundleInfo       `json:"bundle_info,omitempty"`
	Margin        float64           `json:"margin,omitempty"`
	Markup        float64           `json:"markup,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// AppliedPricingRule represents a pricing rule that was successfully applied to an item.
// Tracks which rule was used, the adjustment made, and the reasoning behind the application.
//
// Example:
//
//	// Volume discount rule applied
//	appliedRule := AppliedPricingRule{
//		RuleID: "volume-discount-10",
//		RuleName: "10% Volume Discount",
//		RuleType: "volume",
//		Adjustment: -9.99, // $9.99 discount
//		AdjustmentType: "percentage",
//		AdjustmentValue: 10.0, // 10% off
//		OriginalPrice: 99.99,
//		FinalPrice: 89.99,
//		Reason: "Customer purchased 5+ items qualifying for volume discount",
//		Priority: 1,
//		ConditionsMet: []string{"quantity >= 5", "category = electronics"},
//	}
type AppliedPricingRule struct {
	RuleID      string  `json:"rule_id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Adjustment  float64 `json:"adjustment"`
	Description string  `json:"description,omitempty"`
	Priority    int     `json:"priority"`
}

// TierInfo represents tier-based pricing information for an item.
// Contains details about which pricing tier was applied and the associated benefits.
//
// Example:
//
//	// Bulk pricing tier
//	tierInfo := TierInfo{
//		TierName: "Bulk Tier",
//		TierLevel: 2,
//		MinQuantity: 10,
//		MaxQuantity: 49,
//		TierDiscount: 15.0, // 15% discount
//		TierPrice: 84.99, // Price per unit in this tier
//		OriginalPrice: 99.99,
//		SavingsPerUnit: 15.00,
//		TotalSavings: 150.00, // For 10 units
//		NextTierAt: 50, // Next tier starts at 50 units
//		NextTierDiscount: 20.0, // 20% discount at next tier
//	}
type TierInfo struct {
	TierID       string  `json:"tier_id"`
	TierName     string  `json:"tier_name"`
	MinQuantity  int     `json:"min_quantity"`
	MaxQuantity  int     `json:"max_quantity,omitempty"`
	TierPrice    float64 `json:"tier_price"`
	TierDiscount float64 `json:"tier_discount,omitempty"`
}

// BundleInfo represents information about a bundle that was applied or is available.
// Contains bundle details, pricing, and savings information.
//
// Example:
//
//	// Applied laptop bundle
//	bundleInfo := BundleInfo{
//		BundleID: "laptop-essentials",
//		BundleName: "Laptop Essentials Bundle",
//		BundleType: "fixed",
//		Items: []string{"laptop-001", "mouse-001", "keyboard-001"},
//		OriginalPrice: 1299.97,
//		BundlePrice: 1099.99,
//		Savings: 199.98,
//		SavingsPercent: 15.38,
//		IsApplied: true,
//		QualificationsMet: []string{"all required items present", "minimum quantity met"},
//	}
type BundleInfo struct {
	BundleID     string  `json:"bundle_id"`
	BundleName   string  `json:"bundle_name"`
	BundleType   string  `json:"bundle_type"`
	BundlePrice  float64 `json:"bundle_price"`
	BundleSavings float64 `json:"bundle_savings"`
	ItemsInBundle []string `json:"items_in_bundle"`
}

// PricingResult represents the complete result of a pricing calculation.
// Contains all priced items, totals, applied rules, bundles, and recommendations.
//
// Example:
//
//	// Complete pricing result
//	result := PricingResult{
//		Items: []PricedItem{
//			{ItemID: "widget-001", FinalPrice: 89.99, TotalPrice: 179.98},
//			{ItemID: "gadget-002", FinalPrice: 134.99, TotalPrice: 134.99},
//		},
//		Subtotal: 314.97,
//		TotalDiscount: 25.01,
//		TotalSavings: 25.01,
//		GrandTotal: 314.97,
//		Currency: "USD",
//		AppliedBundles: []BundleInfo{
//			{BundleID: "tech-combo", Savings: 15.00},
//		},
//		Recommendations: []PricingRecommendation{
//			{Type: "bundle", Description: "Add mouse for extra 5% off"},
//		},
//		CalculationTime: time.Now(),
//	}
type PricingResult struct {
	Items           []PricedItem      `json:"items"`
	Subtotal        float64           `json:"subtotal"`
	TotalSavings    float64           `json:"total_savings"`
	TotalDiscount   float64           `json:"total_discount"`
	GrandTotal      float64           `json:"grand_total"`
	Currency        string            `json:"currency"`
	AppliedBundles  []BundleInfo      `json:"applied_bundles,omitempty"`
	AppliedTiers    []TierInfo        `json:"applied_tiers,omitempty"`
	Recommendations []PricingRecommendation `json:"recommendations,omitempty"`
	CalculationTime time.Time         `json:"calculation_time"`
	IsValid         bool              `json:"is_valid"`
	Errors          []string          `json:"errors,omitempty"`
	Warnings        []string          `json:"warnings,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// PricingRecommendation represents a pricing or product recommendation.
// Suggests actions to customers for better pricing or additional value.
//
// Supported recommendation types:
//   - "bundle": Suggests adding items to qualify for bundle pricing
//   - "tier": Suggests increasing quantity to reach better pricing tier
//   - "upsell": Suggests higher-value alternatives
//   - "cross_sell": Suggests complementary products
//   - "loyalty": Suggests actions to earn loyalty benefits
//
// Example:
//
//	// Bundle recommendation
//	recommendation := PricingRecommendation{
//		Type: "bundle",
//		Title: "Complete Your Tech Setup",
//		Description: "Add a wireless mouse to save an additional 10%",
//		ItemIDs: []string{"mouse-wireless-001"},
//		PotentialSavings: 15.99,
//		SavingsPercent: 10.0,
//		Priority: 1,
//		Reason: "Bundle discount available",
//		ActionRequired: "Add mouse-wireless-001 to cart",
//	}
type PricingRecommendation struct {
	Type        string  `json:"type"`        // "bundle", "tier", "upsell", "cross_sell"
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Savings     float64 `json:"savings"`
	ItemIDs     []string `json:"item_ids,omitempty"`
	BundleID    string  `json:"bundle_id,omitempty"`
	Priority    int     `json:"priority"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// DynamicPricingConfig represents configuration for dynamic pricing algorithms.
// Defines how prices should adjust based on real-time market conditions and factors.
//
// Example:
//
//	// Market-responsive dynamic pricing
//	config := DynamicPricingConfig{
//		ID: "electronics-dynamic",
//		Name: "Electronics Dynamic Pricing",
//		IsActive: true,
//		Strategy: "market_responsive",
//		Factors: []PricingFactor{
//			{Type: "demand", Weight: 0.4, MinImpact: -20.0, MaxImpact: 30.0},
//			{Type: "inventory", Weight: 0.3, MinImpact: -15.0, MaxImpact: 25.0},
//			{Type: "competition", Weight: 0.3, MinImpact: -10.0, MaxImpact: 15.0},
//		},
//		Rules: []DynamicPricingRule{
//			{Condition: "inventory_level < 10", Action: "increase_price", Value: 10.0},
//		},
//		MinPriceRatio: 0.7, // Never go below 70% of base price
//		MaxPriceRatio: 1.5, // Never go above 150% of base price
//	}
type DynamicPricingConfig struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Strategy          string            `json:"strategy"` // "demand_based", "inventory_based", "competitor_based", "time_based"
	UpdateFrequency   time.Duration     `json:"update_frequency"`
	MinPriceChange    float64           `json:"min_price_change"`    // Minimum price change percentage
	MaxPriceChange    float64           `json:"max_price_change"`    // Maximum price change percentage
	PriceFloor        float64           `json:"price_floor"`         // Minimum allowed price
	PriceCeiling      float64           `json:"price_ceiling"`       // Maximum allowed price
	Factors           []PricingFactor   `json:"factors"`
	Rules             []DynamicPricingRule `json:"rules"`
	IsActive          bool              `json:"is_active"`
	LastUpdated       time.Time         `json:"last_updated"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// PricingFactor represents an individual factor that influences dynamic pricing.
// Defines how external conditions affect price adjustments with weights and limits.
//
// Supported factor types:
//   - "demand": Customer demand levels
//   - "inventory": Stock levels and availability
//   - "competition": Competitor pricing data
//   - "time": Time-based factors (peak hours, seasons)
//   - "weather": Weather conditions impact
//   - "events": Special events or holidays
//
// Example:
//
//	// Demand-based pricing factor
//	demandFactor := PricingFactor{
//		Type: "demand",
//		Name: "Customer Demand Level",
//		Weight: 0.4, // 40% influence on final price
//		MinImpact: -20.0, // Can reduce price by up to 20%
//		MaxImpact: 30.0, // Can increase price by up to 30%
//		DataSource: "analytics_api",
//		UpdateFrequency: "hourly",
//		IsActive: true,
//	}
type PricingFactor struct {
	Type        string  `json:"type"`        // "demand", "inventory", "competition", "time", "weather", "events"
	Weight      float64 `json:"weight"`      // Factor weight (0-1)
	Threshold   float64 `json:"threshold,omitempty"`   // Threshold value
	Impact      float64 `json:"impact"`      // Price impact percentage
	IsActive    bool    `json:"is_active"`
	Description string  `json:"description,omitempty"`
}

// DynamicPricingRule represents a rule for dynamic pricing adjustments.
// Defines conditions and actions for automatic price changes based on market data.
//
// Example:
//
//	// Low inventory price increase rule
//	rule := DynamicPricingRule{
//		ID: "low-inventory-boost",
//		Name: "Low Inventory Price Boost",
//		Condition: "inventory_level < 5 AND demand_score > 0.8",
//		Action: "increase_price",
//		Value: 15.0, // Increase by 15%
//		MaxAdjustment: 25.0, // Never increase more than 25%
//		Priority: 1,
//		IsActive: true,
//		Description: "Increase price when inventory is low and demand is high",
//	}
type DynamicPricingRule struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Conditions  []PricingCondition  `json:"conditions"`
	Adjustments []PriceAdjustment   `json:"adjustments"`
	Priority    int                 `json:"priority"`
	IsActive    bool                `json:"is_active"`
	ValidFrom   time.Time           `json:"valid_from"`
	ValidUntil  time.Time           `json:"valid_until"`
}

// MarketData represents real-time market data used for pricing decisions.
// Contains competitor pricing, demand trends, and market conditions.
//
// Example:
//
//	// Market data for a product
//	marketData := MarketData{
//		ItemID: "widget-001",
//		AveragePrice: 95.50,
//		MinPrice: 79.99,
//		MaxPrice: 119.99,
//		CompetitorPrices: map[string]float64{
//			"competitor_a": 89.99,
//			"competitor_b": 94.50,
//			"competitor_c": 99.99,
//		},
//		DemandLevel: "high",
//		TrendDirection: "up",
//		LastUpdated: time.Now(),
//		Source: "market_intelligence_api",
//	}
type MarketData struct {
	ItemID          string            `json:"item_id"`
	AveragePrice    float64           `json:"average_price"`
	MinPrice        float64           `json:"min_price"`
	MaxPrice        float64           `json:"max_price"`
	CompetitorPrices map[string]float64 `json:"competitor_prices"`
	DemandLevel     string            `json:"demand_level"` // "low", "medium", "high"
	TrendDirection  string            `json:"trend_direction"` // "up", "down", "stable"
	LastUpdated     time.Time         `json:"last_updated"`
	Source          string            `json:"source"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// PricingAnalytics represents comprehensive analytics data for pricing performance.
// Tracks pricing effectiveness, sales impact, and optimization opportunities.
//
// Example:
//
//	// Monthly pricing analytics
//	analytics := PricingAnalytics{
//		ItemID: "widget-001",
//		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//		PeriodEnd: time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
//		AveragePrice: 92.50,
//		PriceChanges: 8,
//		SalesVolume: 1250,
//		Revenue: 115625.00,
//		Margin: 35.2,
//		ConversionRate: 12.5,
//		PriceElasticity: -1.2, // 1% price increase = 1.2% demand decrease
//		OptimalPrice: 89.99,
//		RecommendedPrice: 91.50,
//	}
type PricingAnalytics struct {
	ItemID           string    `json:"item_id"`
	PeriodStart      time.Time `json:"period_start"`
	PeriodEnd        time.Time `json:"period_end"`
	AveragePrice     float64   `json:"average_price"`
	PriceChanges     int       `json:"price_changes"`
	SalesVolume      int       `json:"sales_volume"`
	Revenue          float64   `json:"revenue"`
	Margin           float64   `json:"margin"`
	ConversionRate   float64   `json:"conversion_rate"`
	PriceElasticity  float64   `json:"price_elasticity"`
	OptimalPrice     float64   `json:"optimal_price,omitempty"`
	RecommendedPrice float64   `json:"recommended_price,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}
package pricing

import (
	"time"
)

// PricingStrategy represents different pricing strategies
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

// PricingType represents the type of pricing calculation
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

// BundleType represents different bundle types
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

// PricingRule represents a pricing rule
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

// PricingCondition represents conditions for pricing rules
type PricingCondition struct {
	Type     string      `json:"type"`     // "quantity", "amount", "customer_type", "time", "inventory"
	Operator string      `json:"operator"` // ">", "<", ">=", "<=", "=", "!=", "in", "between"
	Value    interface{} `json:"value"`    // Condition value
	Logic    string      `json:"logic,omitempty"` // "AND", "OR"
}

// PriceAdjustment represents price adjustments
type PriceAdjustment struct {
	Type        string  `json:"type"`        // "percentage", "fixed", "markup", "markdown"
	Value       float64 `json:"value"`       // Adjustment value
	MinPrice    float64 `json:"min_price,omitempty"`    // Minimum price limit
	MaxPrice    float64 `json:"max_price,omitempty"`    // Maximum price limit
	RoundTo     float64 `json:"round_to,omitempty"`     // Round to nearest value
	Description string  `json:"description,omitempty"`
}

// TierPricing represents tiered pricing structure
type TierPricing struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Tiers       []PriceTier `json:"tiers"`
	IsActive    bool        `json:"is_active"`
	ValidFrom   time.Time   `json:"valid_from"`
	ValidUntil  time.Time   `json:"valid_until"`
}

// PriceTier represents a single pricing tier
type PriceTier struct {
	MinQuantity int     `json:"min_quantity"`
	MaxQuantity int     `json:"max_quantity,omitempty"` // 0 means no upper limit
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount,omitempty"`    // Discount percentage
	FixedPrice  float64 `json:"fixed_price,omitempty"` // Fixed price for this tier
}

// Bundle represents a product bundle
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

// BundleItem represents an item in a bundle
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

// BundlePricing represents bundle pricing configuration
type BundlePricing struct {
	Type         string  `json:"type"`         // "fixed", "percentage", "tiered", "dynamic"
	Value        float64 `json:"value"`        // Price or discount value
	MinPrice     float64 `json:"min_price,omitempty"`     // Minimum bundle price
	MaxPrice     float64 `json:"max_price,omitempty"`     // Maximum bundle price
	BasePrice    float64 `json:"base_price,omitempty"`    // Base bundle price
	SavingsType  string  `json:"savings_type,omitempty"`  // "amount", "percentage"
	SavingsValue float64 `json:"savings_value,omitempty"` // Savings amount or percentage
}

// PricingItem represents an item for pricing calculation
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

// Dimensions represents item dimensions
type Dimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Unit   string  `json:"unit"` // "cm", "in", "mm"
}

// Customer represents customer information for pricing
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

// PricingContext represents the context for pricing calculation
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

// PricingInput represents input for pricing calculation
type PricingInput struct {
	Items       []PricingItem   `json:"items"`
	Customer    Customer        `json:"customer"`
	Context     PricingContext  `json:"context"`
	Rules       []PricingRule   `json:"rules,omitempty"`
	Bundles     []Bundle        `json:"bundles,omitempty"`
	TierPricing []TierPricing   `json:"tier_pricing,omitempty"`
	Options     PricingOptions  `json:"options,omitempty"`
}

// PricingOptions represents options for pricing calculation
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

// PricedItem represents an item with calculated pricing
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

// AppliedPricingRule represents a pricing rule that was applied
type AppliedPricingRule struct {
	RuleID      string  `json:"rule_id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Adjustment  float64 `json:"adjustment"`
	Description string  `json:"description,omitempty"`
	Priority    int     `json:"priority"`
}

// TierInfo represents tier pricing information
type TierInfo struct {
	TierID       string  `json:"tier_id"`
	TierName     string  `json:"tier_name"`
	MinQuantity  int     `json:"min_quantity"`
	MaxQuantity  int     `json:"max_quantity,omitempty"`
	TierPrice    float64 `json:"tier_price"`
	TierDiscount float64 `json:"tier_discount,omitempty"`
}

// BundleInfo represents bundle information
type BundleInfo struct {
	BundleID     string  `json:"bundle_id"`
	BundleName   string  `json:"bundle_name"`
	BundleType   string  `json:"bundle_type"`
	BundlePrice  float64 `json:"bundle_price"`
	BundleSavings float64 `json:"bundle_savings"`
	ItemsInBundle []string `json:"items_in_bundle"`
}

// PricingResult represents the result of pricing calculation
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

// PricingRecommendation represents pricing recommendations
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

// DynamicPricingConfig represents configuration for dynamic pricing
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

// PricingFactor represents factors that influence dynamic pricing
type PricingFactor struct {
	Type        string  `json:"type"`        // "demand", "inventory", "competition", "time", "weather", "events"
	Weight      float64 `json:"weight"`      // Factor weight (0-1)
	Threshold   float64 `json:"threshold,omitempty"`   // Threshold value
	Impact      float64 `json:"impact"`      // Price impact percentage
	IsActive    bool    `json:"is_active"`
	Description string  `json:"description,omitempty"`
}

// DynamicPricingRule represents rules for dynamic pricing
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

// MarketData represents market data for pricing decisions
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

// PricingAnalytics represents pricing analytics data
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
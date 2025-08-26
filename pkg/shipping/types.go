package shipping

import "time"

// ShippingMethod represents different shipping methods
type ShippingMethod string

const (
	ShippingMethodStandard   ShippingMethod = "standard"
	ShippingMethodExpress    ShippingMethod = "express"
	ShippingMethodOvernight  ShippingMethod = "overnight"
	ShippingMethodSameDay    ShippingMethod = "same_day"
	ShippingMethodPickup     ShippingMethod = "pickup"
	ShippingMethodFree       ShippingMethod = "free"
)

// ShippingZone represents shipping zones
type ShippingZone string

const (
	ShippingZoneLocal        ShippingZone = "local"
	ShippingZoneRegional     ShippingZone = "regional"
	ShippingZoneNational     ShippingZone = "national"
	ShippingZoneInternational ShippingZone = "international"
)

// WeightUnit represents weight measurement units
type WeightUnit string

const (
	WeightUnitKG WeightUnit = "kg"
	WeightUnitLB WeightUnit = "lb"
	WeightUnitG  WeightUnit = "g"
	WeightUnitOZ WeightUnit = "oz"
)

// DimensionUnit represents dimension measurement units
type DimensionUnit string

const (
	DimensionUnitCM DimensionUnit = "cm"
	DimensionUnitIN DimensionUnit = "in"
	DimensionUnitM  DimensionUnit = "m"
	DimensionUnitFT DimensionUnit = "ft"
)

// Address represents shipping address
type Address struct {
	Street1    string `json:"street1"`
	Street2    string `json:"street2,omitempty"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

// Dimensions represents package dimensions
type Dimensions struct {
	Length float64       `json:"length"`
	Width  float64       `json:"width"`
	Height float64       `json:"height"`
	Unit   DimensionUnit `json:"unit"`
}

// Weight represents package weight
type Weight struct {
	Value float64    `json:"value"`
	Unit  WeightUnit `json:"unit"`
}

// ShippingItem represents an item for shipping calculation
type ShippingItem struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Quantity    int        `json:"quantity"`
	Weight      Weight     `json:"weight"`
	Dimensions  Dimensions `json:"dimensions"`
	Value       float64    `json:"value"`
	Category    string     `json:"category"`
	IsFragile   bool       `json:"is_fragile,omitempty"`
	IsHazardous bool       `json:"is_hazardous,omitempty"`
	RequiresColdChain bool `json:"requires_cold_chain,omitempty"`
}

// Package represents a shipping package
type Package struct {
	ID         string        `json:"id"`
	Items      []ShippingItem `json:"items"`
	Weight     Weight        `json:"weight"`
	Dimensions Dimensions    `json:"dimensions"`
	Value      float64       `json:"value"`
	IsFragile  bool          `json:"is_fragile"`
	IsHazardous bool         `json:"is_hazardous"`
}

// ShippingRule represents shipping cost calculation rules
type ShippingRule struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Method            ShippingMethod `json:"method"`
	Zone              ShippingZone   `json:"zone"`
	MinWeight         Weight         `json:"min_weight,omitempty"`
	MaxWeight         Weight         `json:"max_weight,omitempty"`
	MinValue          float64        `json:"min_value,omitempty"`
	MaxValue          float64        `json:"max_value,omitempty"`
	BaseCost          float64        `json:"base_cost"`
	WeightRate        float64        `json:"weight_rate,omitempty"`        // Cost per weight unit
	ValueRate         float64        `json:"value_rate,omitempty"`         // Percentage of item value
	DimensionalRate   float64        `json:"dimensional_rate,omitempty"`   // Cost per dimensional weight
	FlatRate          float64        `json:"flat_rate,omitempty"`          // Fixed rate regardless of weight/value
	FreeShippingThreshold float64    `json:"free_shipping_threshold,omitempty"`
	Surcharges        []Surcharge    `json:"surcharges,omitempty"`
	IsActive          bool           `json:"is_active"`
	ValidFrom         time.Time      `json:"valid_from"`
	ValidUntil        time.Time      `json:"valid_until"`
	ApplicableCountries []string     `json:"applicable_countries,omitempty"`
	ApplicableStates    []string     `json:"applicable_states,omitempty"`
	ApplicableCategories []string    `json:"applicable_categories,omitempty"`
}

// Surcharge represents additional shipping charges
type Surcharge struct {
	Type        string  `json:"type"`        // "fragile", "hazardous", "oversized", "remote_area", "fuel"
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	IsPercentage bool   `json:"is_percentage"`
	Condition   string  `json:"condition,omitempty"` // Condition for applying surcharge
}

// ZoneRule represents zone-based shipping rules
type ZoneRule struct {
	Zone           ShippingZone `json:"zone"`
	Countries      []string     `json:"countries,omitempty"`
	States         []string     `json:"states,omitempty"`
	PostalCodes    []string     `json:"postal_codes,omitempty"`
	PostalCodeRanges []PostalCodeRange `json:"postal_code_ranges,omitempty"`
	DistanceKm     float64      `json:"distance_km,omitempty"`
}

// PostalCodeRange represents a range of postal codes
type PostalCodeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// CarrierRule represents carrier-specific shipping rules
type CarrierRule struct {
	CarrierID      string         `json:"carrier_id"`
	CarrierName    string         `json:"carrier_name"`
	Method         ShippingMethod `json:"method"`
	ServiceCode    string         `json:"service_code"`
	BaseCost       float64        `json:"base_cost"`
	WeightRate     float64        `json:"weight_rate"`
	ZoneRates      map[ShippingZone]float64 `json:"zone_rates"`
	MaxWeight      Weight         `json:"max_weight"`
	MaxDimensions  Dimensions     `json:"max_dimensions"`
	DeliveryDays   int            `json:"delivery_days"`
	TrackingIncluded bool         `json:"tracking_included"`
	InsuranceIncluded bool        `json:"insurance_included"`
	SignatureRequired bool        `json:"signature_required"`
}

// ShippingCalculationInput represents input for shipping calculation
type ShippingCalculationInput struct {
	Items           []ShippingItem `json:"items"`
	Packages        []Package      `json:"packages,omitempty"`
	Origin          Address        `json:"origin"`
	Destination     Address        `json:"destination"`
	ShippingRules   []ShippingRule `json:"shipping_rules"`
	ZoneRules       []ZoneRule     `json:"zone_rules,omitempty"`
	CarrierRules    []CarrierRule  `json:"carrier_rules,omitempty"`
	RequestedMethod ShippingMethod `json:"requested_method,omitempty"`
	InsuranceValue  float64        `json:"insurance_value,omitempty"`
	DeliveryDate    time.Time      `json:"delivery_date,omitempty"`
	IsPriority      bool           `json:"is_priority,omitempty"`
}

// ShippingOption represents a shipping option with cost and details
type ShippingOption struct {
	ID              string         `json:"id"`
	Method          ShippingMethod `json:"method"`
	CarrierID       string         `json:"carrier_id,omitempty"`
	CarrierName     string         `json:"carrier_name,omitempty"`
	ServiceName     string         `json:"service_name"`
	Cost            float64        `json:"cost"`
	BaseCost        float64        `json:"base_cost"`
	Surcharges      []AppliedSurcharge `json:"surcharges,omitempty"`
	EstimatedDays   int            `json:"estimated_days"`
	DeliveryDate    time.Time      `json:"delivery_date,omitempty"`
	TrackingIncluded bool          `json:"tracking_included"`
	InsuranceIncluded bool         `json:"insurance_included"`
	SignatureRequired bool         `json:"signature_required"`
	Zone            ShippingZone   `json:"zone"`
	Description     string         `json:"description"`
	Restrictions    []string       `json:"restrictions,omitempty"`
}

// AppliedSurcharge represents a surcharge that was applied
type AppliedSurcharge struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// ShippingCalculationResult represents the result of shipping calculation
type ShippingCalculationResult struct {
	Options         []ShippingOption `json:"options"`
	RecommendedOption *ShippingOption `json:"recommended_option,omitempty"`
	CheapestOption  *ShippingOption  `json:"cheapest_option,omitempty"`
	FastestOption   *ShippingOption  `json:"fastest_option,omitempty"`
	TotalWeight     Weight           `json:"total_weight"`
	TotalValue      float64          `json:"total_value"`
	Zone            ShippingZone     `json:"zone"`
	Distance        float64          `json:"distance,omitempty"`
	IsValid         bool             `json:"is_valid"`
	ErrorMessage    string           `json:"error_message,omitempty"`
	Warnings        []string         `json:"warnings,omitempty"`
}

// DeliveryTimeRule represents delivery time calculation rules
type DeliveryTimeRule struct {
	Method        ShippingMethod `json:"method"`
	Zone          ShippingZone   `json:"zone"`
	BaseDays      int            `json:"base_days"`
	WeightDelayDays int          `json:"weight_delay_days,omitempty"` // Additional days for heavy packages
	WeightThreshold Weight       `json:"weight_threshold,omitempty"`
	DistanceDelayDays int        `json:"distance_delay_days,omitempty"` // Additional days for long distances
	DistanceThreshold float64    `json:"distance_threshold,omitempty"`
	HolidayDelay  int            `json:"holiday_delay,omitempty"`
	WeekendDelay  int            `json:"weekend_delay,omitempty"`
}

// ShippingRestriction represents shipping restrictions
type ShippingRestriction struct {
	Type        string   `json:"type"`        // "item_category", "destination", "weight", "value", "dimensions"
	Condition   string   `json:"condition"`   // The restriction condition
	Message     string   `json:"message"`     // User-friendly restriction message
	Methods     []ShippingMethod `json:"methods,omitempty"` // Restricted methods
	Countries   []string `json:"countries,omitempty"`   // Restricted countries
	Categories  []string `json:"categories,omitempty"`  // Restricted item categories
}

// FreeShippingRule represents free shipping qualification rules
type FreeShippingRule struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	MinOrderValue   float64        `json:"min_order_value,omitempty"`
	MinWeight       Weight         `json:"min_weight,omitempty"`
	ApplicableZones []ShippingZone `json:"applicable_zones,omitempty"`
	ApplicableCategories []string  `json:"applicable_categories,omitempty"`
	ExcludedCategories []string    `json:"excluded_categories,omitempty"`
	MembershipRequired bool        `json:"membership_required,omitempty"`
	PromotionCode   string         `json:"promotion_code,omitempty"`
	ValidFrom       time.Time      `json:"valid_from"`
	ValidUntil      time.Time      `json:"valid_until"`
	IsActive        bool           `json:"is_active"`
}

// PackagingRule represents packaging optimization rules
type PackagingRule struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	MaxWeight       Weight     `json:"max_weight"`
	MaxDimensions   Dimensions `json:"max_dimensions"`
	PackagingCost   float64    `json:"packaging_cost"`
	MaterialType    string     `json:"material_type"` // "box", "envelope", "tube", "custom"
	IsDefault       bool       `json:"is_default"`
	FragileSupport  bool       `json:"fragile_support"`
	HazardousSupport bool      `json:"hazardous_support"`
}
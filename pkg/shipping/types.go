// Package shipping provides comprehensive types and data structures for shipping cost calculation,
// delivery time estimation, and shipping rule management in e-commerce applications.
//
// This package defines all the core types used throughout the shipping module, including:
//   - Shipping methods and zones
//   - Weight and dimension units with conversion support
//   - Address and package representations
//   - Shipping rules and carrier configurations
//   - Calculation inputs and results
//   - Delivery time and restriction rules
//   - Free shipping and packaging rules
//
// Basic Usage:
//
//	package main
//
//	import (
//		"fmt"
//		"time"
//		"your-project/pkg/shipping"
//	)
//
//	func main() {
//		// Create shipping items
//		items := []shipping.ShippingItem{
//			{
//				ID:       "item1",
//				Name:     "Laptop",
//				Quantity: 1,
//				Weight:   shipping.Weight{Value: 2.5, Unit: shipping.WeightUnitKG},
//				Value:    999.99,
//			},
//		}
//
//		// Define addresses
//		origin := shipping.Address{
//			Street1:    "123 Warehouse St",
//			City:       "New York",
//			State:      "NY",
//			PostalCode: "10001",
//			Country:    "US",
//		}
//
//		destination := shipping.Address{
//			Street1:    "456 Customer Ave",
//			City:       "Los Angeles",
//			State:      "CA",
//			PostalCode: "90210",
//			Country:    "US",
//		}
//
//		// Create calculation input
//		input := shipping.ShippingCalculationInput{
//			Items:       items,
//			Origin:      origin,
//			Destination: destination,
//		}
//
//		fmt.Printf("Shipping calculation input created for %d items\n", len(input.Items))
//	}
//
package shipping

import "time"

// ShippingMethod represents different shipping methods available for delivery.
// Each method has different cost structures, delivery times, and service levels.
//
// Example usage:
//
//	method := shipping.ShippingMethodExpress
//	if method == shipping.ShippingMethodOvernight {
//		fmt.Println("Fast delivery selected")
//	}
type ShippingMethod string

const (
	// ShippingMethodStandard represents standard ground shipping (3-7 business days)
	ShippingMethodStandard   ShippingMethod = "standard"
	// ShippingMethodExpress represents expedited shipping (1-3 business days)
	ShippingMethodExpress    ShippingMethod = "express"
	// ShippingMethodOvernight represents next-day delivery
	ShippingMethodOvernight  ShippingMethod = "overnight"
	// ShippingMethodSameDay represents same-day delivery (where available)
	ShippingMethodSameDay    ShippingMethod = "same_day"
	// ShippingMethodPickup represents customer pickup at store/warehouse
	ShippingMethodPickup     ShippingMethod = "pickup"
	// ShippingMethodFree represents free shipping (usually standard speed)
	ShippingMethodFree       ShippingMethod = "free"
)

// ShippingZone represents geographical shipping zones used for cost calculation.
// Zones are typically based on distance from origin and determine shipping rates.
//
// Example usage:
//
//	zone := shipping.ShippingZoneNational
//	if zone == shipping.ShippingZoneInternational {
//		fmt.Println("International shipping required")
//	}
type ShippingZone string

const (
	// ShippingZoneLocal represents local delivery within the same city/metro area
	ShippingZoneLocal        ShippingZone = "local"
	// ShippingZoneRegional represents regional delivery within the same state/province
	ShippingZoneRegional     ShippingZone = "regional"
	// ShippingZoneNational represents national delivery within the same country
	ShippingZoneNational     ShippingZone = "national"
	// ShippingZoneInternational represents international delivery to other countries
	ShippingZoneInternational ShippingZone = "international"
)

// WeightUnit represents weight measurement units used for shipping calculations.
// The system supports automatic conversion between different weight units.
//
// Example usage:
//
//	weight := shipping.Weight{
//		Value: 2.5,
//		Unit:  shipping.WeightUnitKG,
//	}
type WeightUnit string

const (
	// WeightUnitKG represents kilograms (metric system)
	WeightUnitKG WeightUnit = "kg"
	// WeightUnitLB represents pounds (imperial system)
	WeightUnitLB WeightUnit = "lb"
	// WeightUnitG represents grams (metric system)
	WeightUnitG  WeightUnit = "g"
	// WeightUnitOZ represents ounces (imperial system)
	WeightUnitOZ WeightUnit = "oz"
)

// DimensionUnit represents dimension measurement units used for package dimensions.
// The system supports automatic conversion between different dimension units.
//
// Example usage:
//
//	dimensions := shipping.Dimensions{
//		Length: 30,
//		Width:  20,
//		Height: 15,
//		Unit:   shipping.DimensionUnitCM,
//	}
type DimensionUnit string

const (
	// DimensionUnitCM represents centimeters (metric system)
	DimensionUnitCM DimensionUnit = "cm"
	// DimensionUnitIN represents inches (imperial system)
	DimensionUnitIN DimensionUnit = "in"
	// DimensionUnitM represents meters (metric system)
	DimensionUnitM  DimensionUnit = "m"
	// DimensionUnitFT represents feet (imperial system)
	DimensionUnitFT DimensionUnit = "ft"
)

// Address represents a shipping address for origin or destination.
// It includes geographical coordinates for distance calculations and zone determination.
//
// Example usage:
//
//	address := shipping.Address{
//		Street1:    "123 Main St",
//		Street2:    "Apt 4B",
//		City:       "New York",
//		State:      "NY",
//		PostalCode: "10001",
//		Country:    "US",
//		Latitude:   40.7128,
//		Longitude:  -74.0060,
//	}
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

// Dimensions represents the physical dimensions of a package.
// Used for dimensional weight calculations and carrier restrictions.
//
// Example usage:
//
//	dimensions := shipping.Dimensions{
//		Length: 30.0,
//		Width:  20.0,
//		Height: 15.0,
//		Unit:   shipping.DimensionUnitCM,
//	}
//
//	// Calculate volume
//	volume := dimensions.Length * dimensions.Width * dimensions.Height
type Dimensions struct {
	Length float64       `json:"length"`
	Width  float64       `json:"width"`
	Height float64       `json:"height"`
	Unit   DimensionUnit `json:"unit"`
}

// Weight represents the weight of a package or item.
// Supports different weight units with automatic conversion capabilities.
//
// Example usage:
//
//	weight := shipping.Weight{
//		Value: 2.5,
//		Unit:  shipping.WeightUnitKG,
//	}
//
//	// Check if weight exceeds limit
//	if weight.Value > 50.0 && weight.Unit == shipping.WeightUnitKG {
//		fmt.Println("Heavy package")
//	}
type Weight struct {
	Value float64    `json:"value"`
	Unit  WeightUnit `json:"unit"`
}

// ShippingItem represents an individual item that needs to be shipped.
// Contains all necessary information for shipping cost calculation and restrictions.
//
// Example usage:
//
//	item := shipping.ShippingItem{
//		ID:          "item123",
//		Name:        "Laptop Computer",
//		Quantity:    1,
//		Weight:      shipping.Weight{Value: 2.5, Unit: shipping.WeightUnitKG},
//		Dimensions:  shipping.Dimensions{Length: 35, Width: 25, Height: 3, Unit: shipping.DimensionUnitCM},
//		Value:       999.99,
//		Category:    "electronics",
//		IsFragile:   true,
//	}
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

// Package represents a shipping package containing one or more items.
// Used for optimized packaging and consolidated shipping calculations.
//
// Example usage:
//
//	package := shipping.Package{
//		ID:         "pkg001",
//		Items:      []shipping.ShippingItem{item1, item2},
//		Weight:     shipping.Weight{Value: 5.0, Unit: shipping.WeightUnitKG},
//		Dimensions: shipping.Dimensions{Length: 40, Width: 30, Height: 20, Unit: shipping.DimensionUnitCM},
//		Value:      1999.98,
//		IsFragile:  true,
//	}
type Package struct {
	ID         string        `json:"id"`
	Items      []ShippingItem `json:"items"`
	Weight     Weight        `json:"weight"`
	Dimensions Dimensions    `json:"dimensions"`
	Value      float64       `json:"value"`
	IsFragile  bool          `json:"is_fragile"`
	IsHazardous bool         `json:"is_hazardous"`
}

// ShippingRule represents a comprehensive shipping cost calculation rule.
// Defines how shipping costs are calculated based on various factors like weight, value, and destination.
//
// The rule supports multiple pricing models:
//   - Base cost + weight-based pricing
//   - Value-based percentage pricing
//   - Dimensional weight pricing
//   - Flat rate pricing
//   - Free shipping thresholds
//
// Example usage:
//
//	rule := shipping.ShippingRule{
//		ID:                    "rule_express_national",
//		Name:                  "Express National Shipping",
//		Method:                shipping.ShippingMethodExpress,
//		Zone:                  shipping.ShippingZoneNational,
//		BaseCost:              15.00,
//		WeightRate:            2.50, // $2.50 per kg
//		FreeShippingThreshold: 100.00,
//		IsActive:              true,
//		ValidFrom:             time.Now(),
//		ValidUntil:            time.Now().AddDate(1, 0, 0),
//	}
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

// Surcharge represents additional charges that can be applied to shipping costs.
// Surcharges can be fixed amounts or percentage-based and are applied based on specific conditions.
//
// Common surcharge types:
//   - "fragile": For fragile items requiring special handling
//   - "hazardous": For hazardous materials
//   - "oversized": For packages exceeding standard dimensions
//   - "remote_area": For delivery to remote locations
//   - "fuel": For fuel cost adjustments
//
// Example usage:
//
//	surcharge := shipping.Surcharge{
//		Type:         "fragile",
//		Name:         "Fragile Item Handling",
//		Amount:       5.00,
//		IsPercentage: false,
//		Condition:    "item.IsFragile == true",
//	}
type Surcharge struct {
	Type        string  `json:"type"`        // "fragile", "hazardous", "oversized", "remote_area", "fuel"
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	IsPercentage bool   `json:"is_percentage"`
	Condition   string  `json:"condition,omitempty"` // Condition for applying surcharge
}

// ZoneRule represents geographical zone definitions for shipping calculations.
// Defines which locations belong to specific shipping zones based on various criteria.
//
// Zone determination can be based on:
//   - Country codes
//   - State/province codes
//   - Postal code ranges
//   - Distance from origin
//
// Example usage:
//
//	zoneRule := shipping.ZoneRule{
//		Zone:      shipping.ShippingZoneRegional,
//		Countries: []string{"US", "CA"},
//		States:    []string{"NY", "NJ", "CT"},
//		PostalCodeRanges: []shipping.PostalCodeRange{
//			{Start: "10000", End: "19999"},
//		},
//		DistanceKm: 500.0,
//	}
type ZoneRule struct {
	Zone           ShippingZone `json:"zone"`
	Countries      []string     `json:"countries,omitempty"`
	States         []string     `json:"states,omitempty"`
	PostalCodes    []string     `json:"postal_codes,omitempty"`
	PostalCodeRanges []PostalCodeRange `json:"postal_code_ranges,omitempty"`
	DistanceKm     float64      `json:"distance_km,omitempty"`
}

// PostalCodeRange represents a range of postal codes for zone determination.
// Used to define geographical boundaries based on postal code ranges.
//
// Example usage:
//
//	range := shipping.PostalCodeRange{
//		Start: "10000",
//		End:   "19999", // Covers 10000-19999 postal codes
//	}
type PostalCodeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// CarrierRule represents shipping rules specific to a particular carrier (e.g., FedEx, UPS, DHL).
// Contains carrier-specific pricing, limitations, and service details.
//
// Example usage:
//
//	carrierRule := shipping.CarrierRule{
//		CarrierID:         "fedex",
//		CarrierName:       "FedEx",
//		Method:            shipping.ShippingMethodExpress,
//		ServiceCode:       "FEDEX_2_DAY",
//		BaseCost:          12.00,
//		WeightRate:        1.50,
//		ZoneRates: map[shipping.ShippingZone]float64{
//			shipping.ShippingZoneLocal:    0.50,
//			shipping.ShippingZoneNational: 2.00,
//		},
//		MaxWeight:         shipping.Weight{Value: 70, Unit: shipping.WeightUnitKG},
//		DeliveryDays:      2,
//		TrackingIncluded:  true,
//	}
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

// ShippingCalculationInput represents all the input data required for shipping cost calculation.
// Contains items, addresses, rules, and optional preferences for the calculation.
//
// Example usage:
//
//	input := shipping.ShippingCalculationInput{
//		Items: []shipping.ShippingItem{item1, item2},
//		Origin: shipping.Address{
//			Street1: "123 Warehouse St",
//			City:    "New York",
//			State:   "NY",
//			Country: "US",
//		},
//		Destination: shipping.Address{
//			Street1: "456 Customer Ave",
//			City:    "Los Angeles",
//			State:   "CA",
//			Country: "US",
//		},
//		ShippingRules:   rules,
//		RequestedMethod: shipping.ShippingMethodExpress,
//		InsuranceValue:  1000.00,
//	}
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

// ShippingOption represents a calculated shipping option with cost and service details.
// Contains all information needed for the customer to make a shipping choice.
//
// Example usage:
//
//	option := shipping.ShippingOption{
//		ID:                "opt_express_001",
//		Method:            shipping.ShippingMethodExpress,
//		CarrierName:       "FedEx",
//		ServiceName:       "FedEx 2Day",
//		Cost:              25.50,
//		BaseCost:          15.00,
//		EstimatedDays:     2,
//		TrackingIncluded:  true,
//		Zone:              shipping.ShippingZoneNational,
//		Description:       "Express delivery in 2 business days",
//	}
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

// AppliedSurcharge represents a surcharge that was actually applied to a shipping calculation.
// Contains the details of the surcharge and its calculated amount.
//
// Example usage:
//
//	appliedSurcharge := shipping.AppliedSurcharge{
//		Type:        "fragile",
//		Name:        "Fragile Item Handling",
//		Amount:      5.00,
//		Description: "Additional $5.00 for fragile item protection",
//	}
type AppliedSurcharge struct {
	Type        string  `json:"type"`
	Name        string  `json:"name"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// ShippingCalculationResult represents the complete result of a shipping cost calculation.
// Contains all available shipping options and recommendations for the customer.
//
// Example usage:
//
//	result := shipping.ShippingCalculationResult{
//		Options: []shipping.ShippingOption{option1, option2, option3},
//		RecommendedOption: &option2,
//		CheapestOption:   &option1,
//		FastestOption:    &option3,
//		TotalWeight:      shipping.Weight{Value: 5.0, Unit: shipping.WeightUnitKG},
//		TotalValue:       1999.98,
//		Zone:             shipping.ShippingZoneNational,
//		Distance:         2500.0,
//		IsValid:          true,
//	}
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

// DeliveryTimeRule represents rules for calculating delivery time estimates.
// Defines base delivery times and additional delays based on various factors.
//
// Example usage:
//
//	deliveryRule := shipping.DeliveryTimeRule{
//		Method:            shipping.ShippingMethodExpress,
//		Zone:              shipping.ShippingZoneNational,
//		BaseDays:          2,
//		WeightDelayDays:   1,
//		WeightThreshold:   shipping.Weight{Value: 20, Unit: shipping.WeightUnitKG},
//		DistanceDelayDays: 1,
//		DistanceThreshold: 1000.0,
//		HolidayDelay:      1,
//		WeekendDelay:      0,
//	}
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

// ShippingRestriction represents restrictions that prevent or limit shipping to certain destinations or for certain items.
// Used to enforce business rules, legal requirements, or carrier limitations.
//
// Common restriction types:
//   - "item_category": Restrictions based on item category (e.g., hazardous materials)
//   - "destination": Geographic restrictions
//   - "weight": Weight-based restrictions
//   - "value": Value-based restrictions
//   - "dimensions": Size-based restrictions
//
// Example usage:
//
//	restriction := shipping.ShippingRestriction{
//		Type:      "item_category",
//		Condition: "category == 'hazardous'",
//		Message:   "Hazardous materials cannot be shipped via air transport",
//		Methods:   []shipping.ShippingMethod{shipping.ShippingMethodOvernight, shipping.ShippingMethodSameDay},
//		Countries: []string{"CA", "MX"},
//		Categories: []string{"hazardous", "flammable"},
//	}
type ShippingRestriction struct {
	Type        string   `json:"type"`        // "item_category", "destination", "weight", "value", "dimensions"
	Condition   string   `json:"condition"`   // The restriction condition
	Message     string   `json:"message"`     // User-friendly restriction message
	Methods     []ShippingMethod `json:"methods,omitempty"` // Restricted methods
	Countries   []string `json:"countries,omitempty"`   // Restricted countries
	Categories  []string `json:"categories,omitempty"`  // Restricted item categories
}

// FreeShippingRule represents rules that determine when free shipping is offered.
// Defines conditions that must be met for customers to qualify for free shipping.
//
// Example usage:
//
//	freeShippingRule := shipping.FreeShippingRule{
//		ID:                   "free_shipping_promo",
//		Name:                 "Free Shipping on Orders Over $100",
//		MinOrderValue:        100.00,
//		ApplicableZones:      []shipping.ShippingZone{shipping.ShippingZoneNational},
//		ApplicableCategories: []string{"electronics", "books"},
//		ExcludedCategories:   []string{"hazardous"},
//		MembershipRequired:   false,
//		ValidFrom:            time.Now(),
//		ValidUntil:           time.Now().AddDate(0, 3, 0),
//		IsActive:             true,
//	}
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

// PackagingRule represents rules for package optimization and material selection.
// Defines packaging constraints, costs, and capabilities for different package types.
//
// Example usage:
//
//	packagingRule := shipping.PackagingRule{
//		ID:               "standard_box_medium",
//		Name:             "Medium Standard Box",
//		MaxWeight:        shipping.Weight{Value: 10, Unit: shipping.WeightUnitKG},
//		MaxDimensions:    shipping.Dimensions{Length: 40, Width: 30, Height: 20, Unit: shipping.DimensionUnitCM},
//		PackagingCost:    2.50,
//		MaterialType:     "box",
//		IsDefault:        true,
//		FragileSupport:   true,
//		HazardousSupport: false,
//	}
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
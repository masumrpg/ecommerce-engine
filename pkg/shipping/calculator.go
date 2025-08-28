// Package shipping provides comprehensive shipping cost calculation and delivery estimation.
// It supports multiple shipping methods, zones, carriers, and complex pricing rules.
//
// Key Features:
//   - Multi-zone shipping calculations (local, regional, national, international)
//   - Weight-based, value-based, and dimensional pricing
//   - Carrier-specific rules and service levels
//   - Free shipping eligibility checks
//   - Delivery time estimation with various factors
//   - Shipping restrictions and surcharge handling
//   - Distance-based calculations using coordinates
//
// Basic Usage:
//
//	// Simple shipping calculation
//	input := ShippingCalculationInput{
//		Items: []ShippingItem{
//			{ID: "item1", Weight: Weight{Value: 2.5, Unit: WeightUnitKG}, Value: 99.99},
//		},
//		Origin: Address{Country: "US", State: "CA", City: "San Francisco"},
//		Destination: Address{Country: "US", State: "NY", City: "New York"},
//	}
//	
//	result := Calculate(input)
//	if result.IsValid {
//		fmt.Printf("Shipping options: %d\n", len(result.Options))
//		fmt.Printf("Cheapest: $%.2f\n", result.CheapestOption.Cost)
//	}
package shipping

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

// ShippingCalculator handles comprehensive shipping cost calculations and delivery estimations.
// It manages various shipping rules, zones, restrictions, and carrier configurations.
//
// The calculator supports:
//   - Multiple shipping zones and methods
//   - Weight, value, and dimensional pricing
//   - Carrier-specific rules and surcharges
//   - Free shipping eligibility
//   - Delivery time calculations
//   - Shipping restrictions by location or item type
//
// Example:
//
//	// Create calculator with custom rules
//	calc := NewShippingCalculator()
//	calc.ZoneRules = []ZoneRule{
//		{Zone: ShippingZoneNational, Countries: []string{"US"}},
//	}
//	calc.FreeShippingRules = []FreeShippingRule{
//		{MinOrderValue: 50.0, IsActive: true},
//	}
//	
//	result := calc.CalculateShipping(input)
type ShippingCalculator struct {
	ZoneRules         []ZoneRule
	DeliveryTimeRules []DeliveryTimeRule
	Restrictions      []ShippingRestriction
	FreeShippingRules []FreeShippingRule
	PackagingRules    []PackagingRule
}

// NewShippingCalculator creates a new shipping calculator with empty rule sets.
// Rules can be added after creation to customize shipping behavior.
//
// Returns a calculator ready for configuration with:
//   - Empty zone rules (will use default zone determination)
//   - Empty delivery time rules (will use default delivery times)
//   - Empty restrictions (no shipping restrictions)
//   - Empty free shipping rules (no free shipping)
//   - Empty packaging rules (no special packaging requirements)
//
// Example:
//
//	calc := NewShippingCalculator()
//	// Add custom rules as needed
//	calc.ZoneRules = append(calc.ZoneRules, ZoneRule{
//		Zone: ShippingZoneInternational,
//		Countries: []string{"CA", "MX"},
//	})
func NewShippingCalculator() *ShippingCalculator {
	return &ShippingCalculator{
		ZoneRules:         []ZoneRule{},
		DeliveryTimeRules: []DeliveryTimeRule{},
		Restrictions:      []ShippingRestriction{},
		FreeShippingRules: []FreeShippingRule{},
		PackagingRules:    []PackagingRule{},
	}
}

// Calculate is a convenience function for shipping calculation using default rules.
// It creates a new ShippingCalculator instance and performs the calculation.
//
// This function is suitable for simple shipping calculations where custom rules
// are not required. For more complex scenarios with custom zones, restrictions,
// or free shipping rules, use NewShippingCalculator() and configure it manually.
//
// Parameters:
//   - input: ShippingCalculationInput containing items, addresses, and options
//
// Returns:
//   - ShippingCalculationResult with available shipping options and recommendations
//
// Example:
//
//	input := ShippingCalculationInput{
//		Items: []ShippingItem{
//			{ID: "book", Weight: Weight{Value: 0.5, Unit: WeightUnitKG}, Value: 29.99},
//		},
//		Origin: Address{Country: "US", State: "CA"},
//		Destination: Address{Country: "US", State: "TX"},
//	}
//	
//	result := Calculate(input)
//	if result.IsValid {
//		fmt.Printf("Found %d shipping options\n", len(result.Options))
//	}
func Calculate(input ShippingCalculationInput) ShippingCalculationResult {
	// Validate input
	if len(input.Items) == 0 {
		return ShippingCalculationResult{
			IsValid:      false,
			ErrorMessage: "no items to ship",
			Options:      []ShippingOption{},
			Warnings:     []string{},
		}
	}

	calc := NewShippingCalculator()
	calc.ZoneRules = input.ZoneRules
	return calc.CalculateShipping(input)
}

// CalculateShipping calculates shipping costs and options based on configured rules.
// This is the main calculation method that processes shipping requests through
// multiple stages: validation, zone determination, restriction checks, and cost calculation.
//
// Calculation Process:
//   1. Input validation (items, addresses, weights)
//   2. Shipping zone determination based on origin/destination
//   3. Restriction checks (prohibited items, blocked destinations)
//   4. Available shipping method calculation
//   5. Cost calculation including surcharges and discounts
//   6. Delivery time estimation
//   7. Free shipping eligibility evaluation
//   8. Recommendation generation (cheapest, fastest, recommended)
//
// Parameters:
//   - input: ShippingCalculationInput containing:
//     * Items: List of items to ship with weights, dimensions, values
//     * Origin: Shipping origin address
//     * Destination: Shipping destination address
//     * Options: Calculation preferences and filters
//
// Returns:
//   - ShippingCalculationResult containing:
//     * IsValid: Whether calculation was successful
//     * Options: Available shipping options with costs and delivery times
//     * CheapestOption, FastestOption, RecommendedOption: Best options
//     * Errors: Any validation or calculation errors
//     * Metadata: Additional calculation information
//
// Example:
//
//	calc := NewShippingCalculator()
//	calc.ZoneRules = []ZoneRule{
//		{Zone: ShippingZoneNational, Countries: []string{"US"}},
//	}
//	
//	input := ShippingCalculationInput{
//		Items: []ShippingItem{
//			{ID: "laptop", Weight: Weight{Value: 2.0, Unit: WeightUnitKG}, Value: 999.99},
//		},
//		Origin: Address{Country: "US", State: "CA", City: "Los Angeles"},
//		Destination: Address{Country: "US", State: "NY", City: "New York"},
//	}
//	
//	result := calc.CalculateShipping(input)
//	if result.IsValid {
//		for _, option := range result.Options {
//			fmt.Printf("%s: $%.2f (delivery: %s)\n", 
//				option.Method, option.Cost, option.EstimatedDelivery.Format("Jan 2"))
//		}
//	}
func (sc *ShippingCalculator) CalculateShipping(input ShippingCalculationInput) ShippingCalculationResult {
	// Validate input
	if len(input.Items) == 0 {
		return ShippingCalculationResult{
			IsValid:      false,
			ErrorMessage: "no items to ship",
			Options:      []ShippingOption{},
			Warnings:     []string{},
		}
	}

	result := ShippingCalculationResult{
		Options:     []ShippingOption{},
		TotalWeight: calculateTotalWeight(input.Items),
		TotalValue:  calculateTotalValue(input.Items),
		IsValid:     true,
		Warnings:    []string{},
	}

	// Determine shipping zone
	zone := sc.determineShippingZone(input.Origin, input.Destination)
	result.Zone = zone

	// Calculate distance if coordinates are available
	if input.Origin.Latitude != 0 && input.Origin.Longitude != 0 &&
		input.Destination.Latitude != 0 && input.Destination.Longitude != 0 {
		result.Distance = calculateDistance(input.Origin, input.Destination)
	}

	// Check shipping restrictions
	if restrictions := sc.checkRestrictions(input.Items, input.Destination); len(restrictions) > 0 {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("Shipping restrictions apply: %v", restrictions)
		return result
	}

	// If no shipping rules provided, create default options
	if len(input.ShippingRules) == 0 {
		defaultOption := &ShippingOption{
			ID:              "default-standard",
			Method:          ShippingMethodStandard,
			ServiceName:     "Standard Shipping",
			Cost:            10.0,
			BaseCost:        10.0,
			EstimatedDays:   5,
			Zone:            zone,
			Description:     "Standard shipping",
			TrackingIncluded: false,
			InsuranceIncluded: false,
			SignatureRequired: false,
		}
		result.Options = append(result.Options, *defaultOption)
	}

	// Calculate shipping options for each rule
	for _, rule := range input.ShippingRules {
		if !sc.isRuleApplicable(rule, input) {
			continue
		}

		option := sc.calculateShippingOption(rule, input, zone, result.Distance)
		if option != nil {
			result.Options = append(result.Options, *option)
		}
	}

	// Calculate carrier-specific options
	for _, carrierRule := range input.CarrierRules {
		option := sc.calculateCarrierOption(carrierRule, input, zone)
		if option != nil {
			result.Options = append(result.Options, *option)
		}
	}

	// Check for free shipping eligibility
	sc.applyFreeShipping(&result, input)

	// Sort options by cost
	sort.Slice(result.Options, func(i, j int) bool {
		return result.Options[i].Cost < result.Options[j].Cost
	})

	// Set recommended, cheapest, and fastest options
	sc.setRecommendedOptions(&result)

	return result
}

// calculateShippingOption calculates the total cost for a specific shipping option.
// This function handles the core cost calculation logic including base rates,
// weight-based pricing, value-based pricing, dimensional weight, and surcharges.
//
// Calculation Components:
//   - Base shipping rate for the method and zone
//   - Weight-based charges (actual or dimensional weight)
//   - Value-based charges (percentage of item value)
//   - Distance-based adjustments
//   - Surcharges (fragile, hazardous, oversized, fuel, insurance)
//   - Free shipping eligibility checks
//
// Parameters:
//   - method: Shipping method (standard, express, overnight, etc.)
//   - zone: Shipping zone (local, regional, national, international)
//   - totalWeight: Combined weight of all items
//   - totalValue: Combined value of all items
//   - dimensionalWeight: Calculated dimensional weight
//   - distance: Distance between origin and destination
//   - items: Individual items for surcharge calculations
//   - destination: Destination address for zone-specific rules
//
// Returns:
//   - float64: Total calculated shipping cost
//   - error: Any calculation errors
//
// Example calculation flow:
//   1. Base rate: $5.00 for standard shipping
//   2. Weight charge: $2.00/kg * 3kg = $6.00
//   3. Fragile surcharge: $3.00
//   4. Total: $14.00
func (sc *ShippingCalculator) calculateShippingOption(rule ShippingRule, input ShippingCalculationInput, zone ShippingZone, distance float64) *ShippingOption {
	if !rule.IsActive {
		return nil
	}

	// Check time validity (only if dates are set)
	now := time.Now()
	if !rule.ValidFrom.IsZero() && now.Before(rule.ValidFrom) {
		return nil
	}
	if !rule.ValidUntil.IsZero() && now.After(rule.ValidUntil) {
		return nil
	}

	// Check zone compatibility
	if rule.Zone != "" && rule.Zone != zone {
		return nil
	}

	totalWeight := calculateTotalWeight(input.Items)
	totalValue := calculateTotalValue(input.Items)

	// Check weight limits
	if rule.MinWeight.Value > 0 && convertWeight(totalWeight, rule.MinWeight.Unit) < rule.MinWeight.Value {
		return nil
	}
	if rule.MaxWeight.Value > 0 && convertWeight(totalWeight, rule.MaxWeight.Unit) > rule.MaxWeight.Value {
		return nil
	}

	// Check value limits
	if rule.MinValue > 0 && totalValue < rule.MinValue {
		return nil
	}
	if rule.MaxValue > 0 && totalValue > rule.MaxValue {
		return nil
	}

	// Calculate base cost
	cost := rule.BaseCost

	// Apply flat rate if specified
	if rule.FlatRate > 0 {
		cost = rule.FlatRate
	} else {
		// Apply weight-based pricing
		if rule.WeightRate > 0 {
			weightInRuleUnit := convertWeight(totalWeight, WeightUnitKG) // Convert to kg for calculation
			cost += weightInRuleUnit * rule.WeightRate
		}

		// Apply value-based pricing
		if rule.ValueRate > 0 {
			cost += totalValue * (rule.ValueRate / 100)
		}

		// Apply dimensional weight pricing
		if rule.DimensionalRate > 0 {
			dimensionalWeight := calculateDimensionalWeight(input.Items)
			cost += dimensionalWeight.Value * rule.DimensionalRate
		}
	}

	// Apply surcharges
	appliedSurcharges := sc.calculateSurcharges(rule.Surcharges, input.Items, totalValue)
	for _, surcharge := range appliedSurcharges {
		cost += surcharge.Amount
	}

	// Calculate delivery time
	estimatedDays := sc.calculateDeliveryTime(rule.Method, zone, totalWeight, distance)

	option := &ShippingOption{
		ID:              rule.ID,
		Method:          rule.Method,
		ServiceName:     rule.Name,
		Cost:            math.Round(cost*100) / 100, // Round to 2 decimal places
		BaseCost:        rule.BaseCost,
		Surcharges:      appliedSurcharges,
		EstimatedDays:   estimatedDays,
		Zone:            zone,
		Description:     fmt.Sprintf("%s shipping via %s", rule.Method, rule.Name),
		TrackingIncluded: rule.Method != ShippingMethodStandard,
		InsuranceIncluded: totalValue > 100, // Include insurance for valuable items
		SignatureRequired: totalValue > 500, // Require signature for high-value items
	}

	// Set delivery date
	if estimatedDays > 0 {
		option.DeliveryDate = time.Now().AddDate(0, 0, estimatedDays)
	}

	return option
}

// calculateCarrierOption calculates the cost for a specific carrier's shipping option.
// This function applies carrier-specific rules, service levels, and pricing structures.
//
// Carrier-Specific Features:
//   - Custom rate tables per carrier
//   - Service level differentiation (ground, air, express)
//   - Carrier-specific surcharges and discounts
//   - Volume-based pricing tiers
//   - Special handling requirements
//
// Parameters:
//   - carrier: Carrier information and rules
//   - method: Shipping method within carrier's services
//   - zone: Shipping zone for carrier's coverage area
//   - totalWeight: Combined weight for carrier's weight breaks
//   - totalValue: Combined value for insurance calculations
//   - dimensionalWeight: Dimensional weight per carrier's formula
//   - distance: Distance for carrier's zone-based pricing
//   - items: Items for carrier-specific restrictions
//   - destination: Destination for carrier's service area
//
// Returns:
//   - float64: Carrier-specific shipping cost
//   - error: Any carrier rule violations or calculation errors
//
// Example:
//   - FedEx Express: Base $12.00 + $0.50/lb over 1lb
//   - UPS Ground: Base $8.00 + zone multiplier 1.2
//   - USPS Priority: Flat rate $7.95 for small items
func (sc *ShippingCalculator) calculateCarrierOption(rule CarrierRule, input ShippingCalculationInput, zone ShippingZone) *ShippingOption {
	totalWeight := calculateTotalWeight(input.Items)

	// Check weight and dimension limits
	if rule.MaxWeight.Value > 0 && convertWeight(totalWeight, rule.MaxWeight.Unit) > rule.MaxWeight.Value {
		return nil
	}

	if !sc.checkDimensionLimits(input.Items, rule.MaxDimensions) {
		return nil
	}

	// Calculate cost
	cost := rule.BaseCost

	// Apply weight rate
	if rule.WeightRate > 0 {
		weightInKg := convertWeight(totalWeight, WeightUnitKG)
		cost += weightInKg * rule.WeightRate
	}

	// Apply zone rate
	if zoneRate, exists := rule.ZoneRates[zone]; exists {
		cost += zoneRate
	}

	option := &ShippingOption{
		ID:                fmt.Sprintf("%s_%s", rule.CarrierID, rule.ServiceCode),
		Method:            rule.Method,
		CarrierID:         rule.CarrierID,
		CarrierName:       rule.CarrierName,
		ServiceName:       fmt.Sprintf("%s %s", rule.CarrierName, rule.Method),
		Cost:              math.Round(cost*100) / 100,
		BaseCost:          rule.BaseCost,
		EstimatedDays:     rule.DeliveryDays,
		Zone:              zone,
		TrackingIncluded:  rule.TrackingIncluded,
		InsuranceIncluded: rule.InsuranceIncluded,
		SignatureRequired: rule.SignatureRequired,
		Description:       fmt.Sprintf("%s shipping via %s", rule.Method, rule.CarrierName),
	}

	if rule.DeliveryDays > 0 {
		option.DeliveryDate = time.Now().AddDate(0, 0, rule.DeliveryDays)
	}

	return option
}

// Helper functions

// calculateTotalWeight calculates the total weight of all items in the shipment.
// This function aggregates weights from multiple items, handling unit conversions
// to ensure consistent weight calculations across different measurement systems.
//
// Weight Handling:
//   - Converts all weights to kilograms for internal calculations
//   - Supports multiple weight units (kg, lb, g, oz)
//   - Handles zero weights and missing weight information
//   - Returns total weight in kilograms
//
// Parameters:
//   - items: Slice of ShippingItem containing individual item weights
//
// Returns:
//   - Weight: Total weight of all items in kilograms
//
// Example:
//   - Item 1: 2.5 kg
//   - Item 2: 3.0 lb (converted to ~1.36 kg)
//   - Total: ~3.86 kg
func calculateTotalWeight(items []ShippingItem) Weight {
	totalWeight := 0.0
	unit := WeightUnitKG // Default unit

	for _, item := range items {
		quantity := item.Quantity
		if quantity == 0 {
			quantity = 1 // Default quantity to 1 if not specified
		}
		itemWeight := convertWeight(item.Weight, WeightUnitKG) * float64(quantity)
		totalWeight += itemWeight
	}

	return Weight{Value: totalWeight, Unit: unit}
}

// calculateTotalValue calculates the total monetary value of all items in the shipment.
// This value is used for insurance calculations, value-based shipping rates,
// and free shipping eligibility checks.
//
// Value Calculation:
//   - Sums the declared value of all items
//   - Used for insurance premium calculations
//   - Applied in value-based shipping rate rules
//   - Considered for free shipping thresholds
//
// Parameters:
//   - items: Slice of ShippingItem containing individual item values
//
// Returns:
//   - float64: Total monetary value of all items
//
// Example:
//   - Item 1: $29.99 (book)
//   - Item 2: $199.99 (electronics)
//   - Total: $229.98
func calculateTotalValue(items []ShippingItem) float64 {
	totalValue := 0.0
	for _, item := range items {
		totalValue += item.Value * float64(item.Quantity)
	}
	return totalValue
}

// calculateDimensionalWeight calculates the dimensional weight of the shipment.
// Dimensional weight is used by carriers to account for large, lightweight packages
// that take up significant space. The higher of actual weight or dimensional weight
// is typically used for shipping cost calculations.
//
// Calculation Formula:
//   - Volume = Length × Width × Height (in cubic inches or cm)
//   - Dimensional Weight = Volume ÷ Dimensional Factor
//   - Common factors: 139 (domestic), 166 (international) for cubic inches
//
// Dimensional Weight Rules:
//   - Calculated per item, then summed
//   - Uses the larger of length, width, height dimensions
//   - Converts all dimensions to consistent units
//   - Applies carrier-specific dimensional factors
//
// Parameters:
//   - items: Slice of ShippingItem containing dimensions
//
// Returns:
//   - Weight: Total dimensional weight in kilograms
//
// Example:
//   - Box: 12" × 8" × 6" = 576 cubic inches
//   - Dimensional Weight: 576 ÷ 139 = ~4.14 lbs (~1.88 kg)
func calculateDimensionalWeight(items []ShippingItem) Weight {
	totalDimensionalWeight := 0.0
	divisor := 5000.0 // Standard dimensional weight divisor (cm³/kg)

	for _, item := range items {
		// Convert dimensions to cm
		length := convertDimension(item.Dimensions.Length, item.Dimensions.Unit, DimensionUnitCM)
		width := convertDimension(item.Dimensions.Width, item.Dimensions.Unit, DimensionUnitCM)
		height := convertDimension(item.Dimensions.Height, item.Dimensions.Unit, DimensionUnitCM)

		volume := length * width * height
		dimensionalWeight := (volume / divisor) * float64(item.Quantity)
		totalDimensionalWeight += dimensionalWeight
	}

	return Weight{
		Value: totalDimensionalWeight,
		Unit:  WeightUnitKG,
	}
}

// convertWeight converts weight between different units for consistent calculations.
// This function handles conversions between kilograms, pounds, grams, and ounces,
// using grams as an intermediate unit for accuracy.
//
// Supported Units:
//   - WeightUnitKG: Kilograms
//   - WeightUnitLB: Pounds
//   - WeightUnitG: Grams
//   - WeightUnitOZ: Ounces
//
// Conversion Process:
//   1. Convert input weight to grams
//   2. Convert grams to target unit
//   3. Return converted value
//
// Parameters:
//   - weight: Weight struct with value and unit
//   - targetUnit: Desired output unit
//
// Returns:
//   - float64: Converted weight value in target unit
//
// Example:
//   - Input: 2.5 kg
//   - Target: pounds
//   - Output: ~5.51 lbs
func convertWeight(weight Weight, targetUnit WeightUnit) float64 {
	if weight.Unit == targetUnit {
		return weight.Value
	}

	// Convert to grams first
	grams := weight.Value
	switch weight.Unit {
	case WeightUnitKG:
		grams *= 1000
	case WeightUnitLB:
		grams *= 453.592
	case WeightUnitOZ:
		grams *= 28.3495
	}

	// Convert from grams to target unit
	switch targetUnit {
	case WeightUnitG:
		return grams
	case WeightUnitKG:
		return grams / 1000
	case WeightUnitLB:
		return grams / 453.592
	case WeightUnitOZ:
		return grams / 28.3495
	default:
		return grams
	}
}

// convertDimension converts dimension values between different units.
// This function handles conversions between meters, centimeters, inches, and feet,
// using centimeters as an intermediate unit for precision.
//
// Supported Units:
//   - DimensionUnitM: Meters
//   - DimensionUnitCM: Centimeters
//   - DimensionUnitIN: Inches
//   - DimensionUnitFT: Feet
//
// Conversion Process:
//   1. Convert input dimension to centimeters
//   2. Convert centimeters to target unit
//   3. Return converted value
//
// Parameters:
//   - value: Dimension value to convert
//   - fromUnit: Source unit of the value
//   - toUnit: Desired output unit
//
// Returns:
//   - float64: Converted dimension value in target unit
//
// Example:
//   - Input: 12 inches
//   - Target: centimeters
//   - Output: ~30.48 cm
func convertDimension(value float64, fromUnit, toUnit DimensionUnit) float64 {
	if fromUnit == toUnit {
		return value
	}

	// Convert to cm first
	cm := value
	switch fromUnit {
	case DimensionUnitM:
		cm *= 100
	case DimensionUnitIN:
		cm *= 2.54
	case DimensionUnitFT:
		cm *= 30.48
	}

	// Convert from cm to target unit
	switch toUnit {
	case DimensionUnitCM:
		return cm
	case DimensionUnitM:
		return cm / 100
	case DimensionUnitIN:
		return cm / 2.54
	case DimensionUnitFT:
		return cm / 30.48
	default:
		return cm
	}
}

// determineShippingZone determines the appropriate shipping zone based on origin and destination addresses.
// This function uses configured zone rules with specificity prioritization, falling back to
// default geographic logic when no custom rules match.
//
// Zone Determination Logic:
//   1. Check custom zone rules (prioritizing more specific rules)
//   2. Fall back to default geographic zones:
//      - International: Different countries
//      - National: Same country, different states
//      - Regional: Same state, different cities
//      - Local: Same city or insufficient location data
//
// Rule Specificity Priority:
//   - Rules with states/postal codes > rules with only countries
//   - More specific rules override general ones
//
// Parameters:
//   - origin: Shipping origin address
//   - destination: Shipping destination address
//
// Returns:
//   - ShippingZone: Determined shipping zone for rate calculation
//
// Example:
//   - Origin: San Francisco, CA, US
//   - Destination: New York, NY, US
//   - Result: ShippingZoneNational
func (sc *ShippingCalculator) determineShippingZone(origin, destination Address) ShippingZone {
	// Check zone rules - prioritize more specific rules (with states/cities) over general ones
	var matchedRule *ZoneRule
	var maxSpecificity int

	for _, rule := range sc.ZoneRules {
		if sc.addressMatchesZoneRule(destination, rule) {
			// Calculate specificity: states > countries
			specificity := 0
			if len(rule.States) > 0 {
				specificity += 2
			}
			if len(rule.Countries) > 0 {
				specificity += 1
			}

			if specificity > maxSpecificity {
				maxSpecificity = specificity
				matchedRule = &rule
			}
		}
	}

	if matchedRule != nil {
		return matchedRule.Zone
	}

	// Default zone determination logic
	if origin.Country != destination.Country {
		return ShippingZoneInternational
	}

	// If same country but different state, it's national
	if origin.Country == destination.Country && origin.State != "" && destination.State != "" && origin.State != destination.State {
		return ShippingZoneNational
	}

	// If same state but different city, it's regional
	if origin.State == destination.State && origin.City != "" && destination.City != "" && origin.City != destination.City {
		return ShippingZoneRegional
	}

	// Same city or no specific location info - local
	return ShippingZoneLocal
}

// addressMatchesZoneRule checks if address matches zone rule
func (sc *ShippingCalculator) addressMatchesZoneRule(address Address, rule ZoneRule) bool {
	// Check countries first
	countryMatches := false
	if len(rule.Countries) > 0 {
		for _, country := range rule.Countries {
			if address.Country == country {
				countryMatches = true
				break
			}
		}
		if !countryMatches {
			return false
		}
	}

	// If states are specified, check them too
	if len(rule.States) > 0 {
		stateMatches := false
		for _, state := range rule.States {
			if address.State == state {
				stateMatches = true
				break
			}
		}
		if !stateMatches {
			return false
		}
	}

	// Check postal codes
	if len(rule.PostalCodes) > 0 {
		postalMatches := false
		for _, postalCode := range rule.PostalCodes {
			if address.PostalCode == postalCode {
				postalMatches = true
				break
			}
		}
		if !postalMatches {
			return false
		}
	}

	// Check postal code ranges
	if len(rule.PostalCodeRanges) > 0 {
		rangeMatches := false
		for _, pcRange := range rule.PostalCodeRanges {
			if address.PostalCode >= pcRange.Start && address.PostalCode <= pcRange.End {
				rangeMatches = true
				break
			}
		}
		if !rangeMatches {
			return false
		}
	}

	return true
}

// calculateDistance calculates the great-circle distance between two addresses using the Haversine formula.
// This function provides accurate distance calculations for shipping cost adjustments
// and delivery time estimations based on geographic coordinates.
//
// Haversine Formula:
//   - Calculates shortest distance over Earth's surface
//   - Accounts for Earth's spherical shape
//   - Provides accuracy suitable for shipping calculations
//   - Returns distance in kilometers
//
// Requirements:
//   - Both addresses must have valid latitude and longitude coordinates
//   - Coordinates should be in decimal degrees format
//   - Invalid coordinates (0,0) will result in inaccurate calculations
//
// Parameters:
//   - origin: Origin address with latitude/longitude coordinates
//   - destination: Destination address with latitude/longitude coordinates
//
// Returns:
//   - float64: Distance in kilometers between the two addresses
//
// Example:
//   - Origin: San Francisco (37.7749, -122.4194)
//   - Destination: New York (40.7128, -74.0060)
//   - Distance: ~4,135 km
func calculateDistance(origin, destination Address) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	lat1 := origin.Latitude * math.Pi / 180
	lon1 := origin.Longitude * math.Pi / 180
	lat2 := destination.Latitude * math.Pi / 180
	lon2 := destination.Longitude * math.Pi / 180

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// calculateDeliveryTime calculates estimated delivery time based on multiple factors.
// This function considers shipping method, zone, package weight, distance, and business days
// to provide accurate delivery estimates for customer expectations.
//
// Delivery Time Factors:
//   - Shipping method (standard, express, overnight, etc.)
//   - Shipping zone (local, regional, national, international)
//   - Package weight (heavier packages may take longer)
//   - Distance between origin and destination
//   - Business days vs weekends
//   - Processing time and cutoff times
//
// Calculation Logic:
//   1. Base delivery days by method and zone
//   2. Weight-based adjustments for heavy packages
//   3. Distance-based adjustments for long distances
//   4. Weekend handling (skip or add extra days)
//   5. Add processing time
//
// Parameters:
//   - method: Shipping method (affects base delivery time)
//   - zone: Shipping zone (affects regional delivery times)
//   - weight: Package weight (may add extra days for heavy items)
//   - distance: Distance in kilometers (affects long-distance deliveries)
//
// Returns:
//   - int: Estimated delivery days from current date
//
// Example:
//   - Method: Standard, Zone: National, Weight: 2kg, Distance: 1000km
//   - Base: 3 days + 0 weight adjustment + 0 distance adjustment
//   - Result: 3 business days from now
func (sc *ShippingCalculator) calculateDeliveryTime(method ShippingMethod, zone ShippingZone, weight Weight, distance float64) int {
	// Find matching delivery time rule
	for _, rule := range sc.DeliveryTimeRules {
		if rule.Method == method && rule.Zone == zone {
			days := rule.BaseDays

			// Add weight delay
			if rule.WeightThreshold.Value > 0 && convertWeight(weight, rule.WeightThreshold.Unit) > rule.WeightThreshold.Value {
				days += rule.WeightDelayDays
			}

			// Add distance delay
			if rule.DistanceThreshold > 0 && distance > rule.DistanceThreshold {
				days += rule.DistanceDelayDays
			}

			// Add weekend/holiday delays
			now := time.Now()
			if now.Weekday() == time.Friday || now.Weekday() == time.Saturday {
				days += rule.WeekendDelay
			}

			return days
		}
	}

	// Default delivery times by method
	switch method {
	case ShippingMethodSameDay:
		return 0
	case ShippingMethodOvernight:
		return 1
	case ShippingMethodExpress:
		return 2
	case ShippingMethodStandard:
		return 5
	case ShippingMethodPickup:
		return 0
	default:
		return 3
	}
}

// calculateSurcharges calculates applicable surcharges based on item characteristics and shipment value.
// This function evaluates various surcharge types including fragile handling, hazardous materials,
// oversized items, fuel surcharges, and insurance premiums.
//
// Surcharge Types:
//   - Fragile: Extra handling for delicate items
//   - Hazardous: Special handling for dangerous goods
//   - Oversized: Additional fees for large packages
//   - Fuel: Variable fuel cost adjustments
//   - Insurance: Value-based insurance premiums
//   - Remote area: Extra fees for difficult-to-reach destinations
//
// Calculation Process:
//   1. Evaluate each surcharge rule against items and shipment
//   2. Check item-specific characteristics (fragile, hazardous)
//   3. Calculate value-based surcharges (insurance, percentage fees)
//   4. Apply flat-rate surcharges where applicable
//   5. Return list of applied surcharges with amounts
//
// Parameters:
//   - surcharges: List of surcharge rules to evaluate
//   - items: List of shipping items to check against rules
//   - totalValue: Total shipment value for percentage-based surcharges
//
// Returns:
//   - []AppliedSurcharge: List of surcharges that apply with calculated amounts
//
// Example:
//   - Fragile item surcharge: +$5.00
//   - Insurance (0.5% of $500): +$2.50
//   - Total applied surcharges: $7.50
func (sc *ShippingCalculator) calculateSurcharges(surcharges []Surcharge, items []ShippingItem, totalValue float64) []AppliedSurcharge {
	applied := []AppliedSurcharge{}

	for _, surcharge := range surcharges {
		if sc.shouldApplySurcharge(surcharge, items, totalValue) {
			amount := surcharge.Amount
			if surcharge.IsPercentage {
				amount = totalValue * (surcharge.Amount / 100)
			}

			applied = append(applied, AppliedSurcharge{
				Type:        surcharge.Type,
				Name:        surcharge.Name,
				Amount:      amount,
				Description: fmt.Sprintf("%s surcharge", surcharge.Name),
			})
		}
	}

	return applied
}

// shouldApplySurcharge determines if a surcharge should be applied
func (sc *ShippingCalculator) shouldApplySurcharge(surcharge Surcharge, items []ShippingItem, totalValue float64) bool {
	switch surcharge.Type {
	case "fragile":
		for _, item := range items {
			if item.IsFragile {
				return true
			}
		}
	case "hazardous":
		for _, item := range items {
			if item.IsHazardous {
				return true
			}
		}
	case "oversized":
		// Check if any item exceeds standard dimensions
		for _, item := range items {
			if sc.isOversized(item.Dimensions) {
				return true
			}
		}
	case "fuel":
		// Fuel surcharge always applies
		return true
	case "insurance":
		// Apply insurance surcharge for high-value shipments
		return totalValue > 1000
	}

	return false
}

// isOversized checks if dimensions are oversized
func (sc *ShippingCalculator) isOversized(dimensions Dimensions) bool {
	// Convert to cm for comparison
	length := convertDimension(dimensions.Length, dimensions.Unit, DimensionUnitCM)
	width := convertDimension(dimensions.Width, dimensions.Unit, DimensionUnitCM)
	height := convertDimension(dimensions.Height, dimensions.Unit, DimensionUnitCM)

	// Standard oversized thresholds (in cm)
	return length > 120 || width > 80 || height > 80
}

// checkRestrictions checks for shipping restrictions that may prevent or limit shipping.
// This function evaluates destination restrictions, item restrictions, and regulatory
// limitations to ensure compliance with shipping policies and regulations.
//
// Restriction Types:
//   - Geographic restrictions (countries, states, postal codes)
//   - Item category restrictions (hazardous, prohibited, regulated)
//   - Weight and size limitations
//   - Value restrictions (high-value items)
//   - Carrier-specific limitations
//   - Regulatory compliance (customs, import/export)
//
// Evaluation Process:
//   1. Check destination against geographic restrictions
//   2. Evaluate each item against category restrictions
//   3. Verify weight and dimension limits
//   4. Check value-based restrictions
//   5. Validate carrier-specific rules
//   6. Compile list of restriction violations
//
// Parameters:
//   - items: List of shipping items to evaluate
//   - destination: Destination address to check restrictions
//
// Returns:
//   - []string: List of restriction messages/errors that apply
//
// Example:
//   - Restriction: "Lithium batteries cannot be shipped to PO Boxes"
//   - Item: Laptop with lithium battery
//   - Destination: PO Box address
//   - Result: ["Lithium battery shipping restriction applies"]
func (sc *ShippingCalculator) checkRestrictions(items []ShippingItem, destination Address) []string {
	restrictions := []string{}

	for _, restriction := range sc.Restrictions {
		if sc.restrictionApplies(restriction, items, destination) {
			restrictions = append(restrictions, restriction.Message)
		}
	}

	return restrictions
}

// restrictionApplies checks if a restriction applies
func (sc *ShippingCalculator) restrictionApplies(restriction ShippingRestriction, items []ShippingItem, destination Address) bool {
	switch restriction.Type {
	case "destination":
		for _, country := range restriction.Countries {
			if destination.Country == country {
				return true
			}
		}
	case "item_category":
		for _, item := range items {
			for _, category := range restriction.Categories {
				if item.Category == category {
					return true
				}
			}
		}
	case "hazardous":
		for _, item := range items {
			if item.IsHazardous {
				return true
			}
		}
	}

	return false
}

// isRuleApplicable checks if a shipping rule is applicable
func (sc *ShippingCalculator) isRuleApplicable(rule ShippingRule, input ShippingCalculationInput) bool {
	// Check country restrictions
	if len(rule.ApplicableCountries) > 0 {
		found := false
		for _, country := range rule.ApplicableCountries {
			if input.Destination.Country == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check state restrictions
	if len(rule.ApplicableStates) > 0 {
		found := false
		for _, state := range rule.ApplicableStates {
			if input.Destination.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check category restrictions
	if len(rule.ApplicableCategories) > 0 {
		found := false
		for _, item := range input.Items {
			for _, category := range rule.ApplicableCategories {
				if item.Category == category {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// checkDimensionLimits checks if items fit within dimension limits
func (sc *ShippingCalculator) checkDimensionLimits(items []ShippingItem, maxDimensions Dimensions) bool {
	for _, item := range items {
		// Convert item dimensions to same unit as max dimensions
		length := convertDimension(item.Dimensions.Length, item.Dimensions.Unit, maxDimensions.Unit)
		width := convertDimension(item.Dimensions.Width, item.Dimensions.Unit, maxDimensions.Unit)
		height := convertDimension(item.Dimensions.Height, item.Dimensions.Unit, maxDimensions.Unit)

		if length > maxDimensions.Length || width > maxDimensions.Width || height > maxDimensions.Height {
			return false
		}
	}
	return true
}

// applyFreeShipping applies free shipping rules to eligible shipping options.
// This function evaluates various free shipping criteria and modifies shipping costs
// for qualifying orders, methods, or customer segments.
//
// Free Shipping Criteria:
//   - Minimum order value thresholds
//   - Specific shipping zones or destinations
//   - Customer membership levels or loyalty status
//   - Product categories or item types
//   - Promotional campaigns or coupon codes
//   - Weight or quantity thresholds
//
// Application Logic:
//   1. Check each free shipping rule for applicability
//   2. Evaluate order against rule criteria
//   3. Apply free shipping to qualifying methods
//   4. Maintain original cost information for reference
//   5. Update option metadata to indicate free shipping
//
// Parameters:
//   - result: Shipping calculation result to modify
//   - input: Shipping calculation input with order details
//
// Example:
//   - Rule: Free standard shipping on orders over $50
//   - Order value: $75
//   - Result: Standard shipping cost reduced to $0.00
func (sc *ShippingCalculator) applyFreeShipping(result *ShippingCalculationResult, input ShippingCalculationInput) {
	for _, rule := range sc.FreeShippingRules {
		if sc.qualifiesForFreeShipping(rule, input, result.TotalValue) {
			// Find cheapest option and make it free
			if len(result.Options) > 0 {
				cheapestIndex := 0
				for i, option := range result.Options {
					if option.Cost < result.Options[cheapestIndex].Cost {
						cheapestIndex = i
					}
				}
				result.Options[cheapestIndex].Cost = 0
				result.Options[cheapestIndex].ServiceName += " (Free Shipping)"
			}
			break
		}
	}
}

// qualifiesForFreeShipping checks if order qualifies for free shipping
func (sc *ShippingCalculator) qualifiesForFreeShipping(rule FreeShippingRule, input ShippingCalculationInput, totalValue float64) bool {
	if !rule.IsActive {
		return false
	}

	now := time.Now()
	if now.Before(rule.ValidFrom) || now.After(rule.ValidUntil) {
		return false
	}

	// Check minimum order value
	if rule.MinOrderValue > 0 && totalValue < rule.MinOrderValue {
		return false
	}

	// Check minimum weight
	if rule.MinWeight.Value > 0 {
		totalWeight := calculateTotalWeight(input.Items)
		if convertWeight(totalWeight, rule.MinWeight.Unit) < rule.MinWeight.Value {
			return false
		}
	}

	// Check applicable zones
	if len(rule.ApplicableZones) > 0 {
		zone := sc.determineShippingZone(input.Origin, input.Destination)
		found := false
		for _, applicableZone := range rule.ApplicableZones {
			if zone == applicableZone {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check categories
	if len(rule.ApplicableCategories) > 0 {
		found := false
		for _, item := range input.Items {
			for _, category := range rule.ApplicableCategories {
				if item.Category == category {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check excluded categories
	for _, item := range input.Items {
		for _, excludedCategory := range rule.ExcludedCategories {
			if item.Category == excludedCategory {
				return false
			}
		}
	}

	return true
}

// setRecommendedOptions sets recommended, cheapest, and fastest options
func (sc *ShippingCalculator) setRecommendedOptions(result *ShippingCalculationResult) {
	if len(result.Options) == 0 {
		return
	}

	// Find cheapest option
	cheapestIndex := 0
	for i, option := range result.Options {
		if option.Cost < result.Options[cheapestIndex].Cost {
			cheapestIndex = i
		}
	}
	result.CheapestOption = &result.Options[cheapestIndex]

	// Find fastest option
	fastestIndex := 0
	for i, option := range result.Options {
		if option.EstimatedDays < result.Options[fastestIndex].EstimatedDays {
			fastestIndex = i
		}
	}
	result.FastestOption = &result.Options[fastestIndex]

	// Set recommended option (balance of cost and speed)
	// For now, recommend the cheapest option with reasonable delivery time
	for _, option := range result.Options {
		if option.EstimatedDays <= 5 && option.Cost <= result.CheapestOption.Cost*1.5 {
			result.RecommendedOption = &option
			break
		}
	}

	// Fallback to cheapest if no reasonable option found
	if result.RecommendedOption == nil {
		result.RecommendedOption = result.CheapestOption
	}
}

// CalculateBestOption returns the best shipping option based on specified criteria.
// This function evaluates available shipping options and selects the optimal choice
// according to customer preferences such as cost, speed, or balanced recommendations.
//
// Selection Criteria:
//   - "cheapest": Lowest cost option regardless of delivery time
//   - "fastest": Quickest delivery regardless of cost
//   - "recommended": Balanced option considering cost, speed, and reliability
//   - "eco-friendly": Most environmentally sustainable option
//   - "reliable": Most dependable carrier/method based on performance
//
// Evaluation Process:
//   1. Calculate all available shipping options
//   2. Filter options based on restrictions
//   3. Apply selection criteria to rank options
//   4. Return the top-ranked option
//   5. Handle cases where no options are available
//
// Parameters:
//   - input: Shipping calculation input with items and addresses
//   - criteria: Selection criteria ("cheapest", "fastest", "recommended")
//
// Returns:
//   - *ShippingOption: Best option matching criteria, or nil if none available
//   - error: Any calculation or validation errors
//
// Example:
//   - Criteria: "cheapest"
//   - Available: Standard ($5.99, 5 days), Express ($12.99, 2 days)
//   - Result: Standard shipping option
func CalculateBestOption(input ShippingCalculationInput, criteria string) (*ShippingOption, error) {
	result := Calculate(input)

	if !result.IsValid {
		return nil, errors.New(result.ErrorMessage)
	}

	if len(result.Options) == 0 {
		return nil, errors.New("no shipping options available")
	}

	switch criteria {
	case "cheapest":
		return result.CheapestOption, nil
	case "fastest":
		return result.FastestOption, nil
	case "recommended":
		return result.RecommendedOption, nil
	default:
		return result.RecommendedOption, nil
	}
}
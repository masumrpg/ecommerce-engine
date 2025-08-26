package shipping

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
)

// ShippingCalculator handles shipping cost calculations
type ShippingCalculator struct {
	ZoneRules         []ZoneRule
	DeliveryTimeRules []DeliveryTimeRule
	Restrictions      []ShippingRestriction
	FreeShippingRules []FreeShippingRule
	PackagingRules    []PackagingRule
}

// NewShippingCalculator creates a new shipping calculator
func NewShippingCalculator() *ShippingCalculator {
	return &ShippingCalculator{
		ZoneRules:         []ZoneRule{},
		DeliveryTimeRules: []DeliveryTimeRule{},
		Restrictions:      []ShippingRestriction{},
		FreeShippingRules: []FreeShippingRule{},
		PackagingRules:    []PackagingRule{},
	}
}

// Calculate calculates shipping options for given input
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

// CalculateShipping calculates shipping options
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

// calculateShippingOption calculates shipping cost for a specific rule
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
			cost += dimensionalWeight * rule.DimensionalRate
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

// calculateCarrierOption calculates shipping cost for carrier-specific rules
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

// calculateTotalWeight calculates total weight of all items
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

// calculateTotalValue calculates total value of all items
func calculateTotalValue(items []ShippingItem) float64 {
	totalValue := 0.0
	for _, item := range items {
		totalValue += item.Value * float64(item.Quantity)
	}
	return totalValue
}

// calculateDimensionalWeight calculates dimensional weight
func calculateDimensionalWeight(items []ShippingItem) float64 {
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

	return totalDimensionalWeight
}

// convertWeight converts weight between different units
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

// convertDimension converts dimension between different units
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

// determineShippingZone determines shipping zone based on addresses
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

// calculateDistance calculates distance between two addresses using Haversine formula
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

// calculateDeliveryTime calculates estimated delivery time
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

// calculateSurcharges calculates applicable surcharges
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

// checkRestrictions checks shipping restrictions
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

// applyFreeShipping checks and applies free shipping rules
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

// CalculateBestOption returns the best shipping option based on criteria
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
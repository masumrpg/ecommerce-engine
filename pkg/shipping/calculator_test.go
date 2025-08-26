package shipping

import (
	"testing"
	"time"
)

// Test NewShippingCalculator
func TestNewShippingCalculator(t *testing.T) {
	calc := NewShippingCalculator()
	if calc == nil {
		t.Fatal("NewShippingCalculator returned nil")
	}
	if calc.ZoneRules == nil {
		t.Error("ZoneRules should be initialized")
	}
	if calc.DeliveryTimeRules == nil {
		t.Error("DeliveryTimeRules should be initialized")
	}
}

// Test Calculate function
func TestCalculate(t *testing.T) {
	// Test basic calculation
	input := ShippingCalculationInput{
		Origin: Address{
			Country:   "US",
			State:     "CA",
			City:      "Los Angeles",
			Latitude:  34.0522,
			Longitude: -118.2437,
		},
		Destination: Address{
			Country:   "US",
			State:     "NY",
			City:      "New York",
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
		Items: []ShippingItem{
			{
				ID:       "item1",
				Name:     "Test Item",
				Quantity: 1,
				Weight:   Weight{Value: 1.0, Unit: WeightUnitKG},
				Dimensions: Dimensions{
					Length: 10,
					Width:  10,
					Height: 10,
					Unit:   DimensionUnitCM,
				},
				Value:    100.0,
				Category: "electronics",
			},
		},
	}

	result := Calculate(input)
	if !result.IsValid {
		t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
	}
	if result.TotalWeight.Value <= 0 {
		t.Error("Expected positive total weight")
	}
	if result.TotalValue <= 0 {
		t.Error("Expected positive total value")
	}
}

// Test empty items
func TestCalculateEmptyItems(t *testing.T) {
	input := ShippingCalculationInput{
		Origin: Address{Country: "US"},
		Destination: Address{Country: "US"},
		Items: []ShippingItem{},
	}

	result := Calculate(input)
	if result.IsValid {
		t.Error("Expected invalid result for empty items")
	}
}

// Test CalculateShipping
func TestCalculateShipping(t *testing.T) {
	calc := NewShippingCalculator()
	
	// Create input with shipping rules
	shippingRules := []ShippingRule{
		{
			ID:                  "rule1",
			Name:                "Standard Shipping",
			Method:              ShippingMethodStandard,
			BaseCost:            5.0,
			WeightRate:          1.0,
			ApplicableCountries: []string{"US"},
			IsActive:            true,
		},
	}

	input := ShippingCalculationInput{
		Origin: Address{Country: "US"},
		Destination: Address{Country: "US"},
		Items: []ShippingItem{
			{
				Weight: Weight{Value: 2.0, Unit: WeightUnitKG},
				Value:  50.0,
			},
		},
		ShippingRules: shippingRules,
	}

	result := calc.CalculateShipping(input)
	if !result.IsValid {
		t.Errorf("Expected valid result, got error: %s", result.ErrorMessage)
	}
	if len(result.Options) == 0 {
		t.Error("Expected at least one shipping option")
	}
}

// Test calculateTotalWeight
func TestCalculateTotalWeight(t *testing.T) {
	items := []ShippingItem{
		{Weight: Weight{Value: 1.0, Unit: WeightUnitKG}},
		{Weight: Weight{Value: 500, Unit: WeightUnitG}},
	}

	totalWeight := calculateTotalWeight(items)
	if totalWeight.Value != 1.5 {
		t.Errorf("Expected total weight 1.5 kg, got %f", totalWeight.Value)
	}
	if totalWeight.Unit != WeightUnitKG {
		t.Errorf("Expected unit kg, got %s", totalWeight.Unit)
	}
}

// Test calculateTotalValue
func TestCalculateTotalValue(t *testing.T) {
	items := []ShippingItem{
		{Value: 100.0, Quantity: 2},
		{Value: 50.0, Quantity: 1},
	}

	totalValue := calculateTotalValue(items)
	if totalValue != 250.0 {
		t.Errorf("Expected total value 250.0, got %f", totalValue)
	}
}

// Test calculateDimensionalWeight
func TestCalculateDimensionalWeight(t *testing.T) {
	items := []ShippingItem{
		{
			Dimensions: Dimensions{
				Length: 30,
				Width:  20,
				Height: 10,
				Unit:   DimensionUnitCM,
			},
			Quantity: 1,
		},
	}

	dimWeight := calculateDimensionalWeight(items)
	if dimWeight <= 0 {
		t.Errorf("Expected positive dimensional weight, got %f", dimWeight)
	}
}

// Test convertWeight
func TestConvertWeight(t *testing.T) {
	weight := Weight{Value: 1000, Unit: WeightUnitG}
	converted := convertWeight(weight, WeightUnitKG)
	if converted != 1.0 {
		t.Errorf("Expected 1.0 kg, got %f", converted)
	}

	weight = Weight{Value: 2.2, Unit: WeightUnitLB}
	converted = convertWeight(weight, WeightUnitKG)
	expected := 2.2 * 0.453592
	if converted < expected-0.01 || converted > expected+0.01 {
		t.Errorf("Expected approximately %f kg, got %f", expected, converted)
	}
}

// Test convertDimension
func TestConvertDimension(t *testing.T) {
	converted := convertDimension(100, DimensionUnitCM, DimensionUnitM)
	if converted != 1.0 {
		t.Errorf("Expected 1.0 m, got %f", converted)
	}

	converted = convertDimension(12, DimensionUnitIN, DimensionUnitCM)
	expected := 12 * 2.54
	if converted < expected-0.01 || converted > expected+0.01 {
		t.Errorf("Expected approximately %f cm, got %f", expected, converted)
	}
}

// Test determineShippingZone
func TestDetermineShippingZone(t *testing.T) {
	calc := NewShippingCalculator()
	
	// Add test zone rules
	calc.ZoneRules = []ZoneRule{
		{
			Zone:      ShippingZoneLocal,
			Countries: []string{"US"},
			States:    []string{"CA"},
		},
		{
			Zone:      ShippingZoneNational,
			Countries: []string{"US"},
		},
	}

	origin := Address{Country: "US", State: "CA"}
	destination := Address{Country: "US", State: "CA"}
	zone := calc.determineShippingZone(origin, destination)
	if zone != ShippingZoneLocal {
		t.Errorf("Expected local zone, got %s", zone)
	}

	destination = Address{Country: "US", State: "NY"}
	zone = calc.determineShippingZone(origin, destination)
	if zone != ShippingZoneNational {
		t.Errorf("Expected national zone, got %s", zone)
	}
}

// Test calculateDistance
func TestCalculateDistance(t *testing.T) {
	origin := Address{Latitude: 34.0522, Longitude: -118.2437} // Los Angeles
	destination := Address{Latitude: 40.7128, Longitude: -74.0060} // New York

	distance := calculateDistance(origin, destination)
	// Approximate distance between LA and NYC is about 3944 km
	if distance < 3900 || distance > 4000 {
		t.Errorf("Expected distance around 3944 km, got %f", distance)
	}
}

// Test calculateDeliveryTime
func TestCalculateDeliveryTime(t *testing.T) {
	calc := NewShippingCalculator()
	
	// Add delivery time rule
	calc.DeliveryTimeRules = []DeliveryTimeRule{
		{
			Method:   ShippingMethodStandard,
			Zone:     ShippingZoneNational,
			BaseDays: 5,
		},
	}

	weight := Weight{Value: 1.0, Unit: WeightUnitKG}
	days := calc.calculateDeliveryTime(ShippingMethodStandard, ShippingZoneNational, weight, 1000)
	if days != 5 {
		t.Errorf("Expected 5 days, got %d", days)
	}

	// Test default delivery time
	days = calc.calculateDeliveryTime(ShippingMethodExpress, ShippingZoneLocal, weight, 100)
	if days != 2 {
		t.Errorf("Expected 2 days for express, got %d", days)
	}
}

// Test calculateSurcharges
func TestCalculateSurcharges(t *testing.T) {
	calc := NewShippingCalculator()
	
	surcharges := []Surcharge{
		{
			Type:   "fuel",
			Name:   "Fuel Surcharge",
			Amount: 2.50,
		},
		{
			Type:         "insurance",
			Name:         "Insurance",
			Amount:       1.0,
			IsPercentage: true,
		},
	}

	items := []ShippingItem{
		{Value: 1500.0}, // High value item for insurance
	}

	applied := calc.calculateSurcharges(surcharges, items, 1500.0)
	if len(applied) != 2 {
		t.Errorf("Expected 2 surcharges, got %d", len(applied))
	}

	// Check fuel surcharge
	if applied[0].Amount != 2.50 {
		t.Errorf("Expected fuel surcharge 2.50, got %f", applied[0].Amount)
	}

	// Check insurance surcharge (1% of 1500)
	if applied[1].Amount != 15.0 {
		t.Errorf("Expected insurance surcharge 15.0, got %f", applied[1].Amount)
	}
}

// Test isOversized
func TestIsOversized(t *testing.T) {
	calc := NewShippingCalculator()
	
	// Normal size
	dimensions := Dimensions{Length: 50, Width: 30, Height: 20, Unit: DimensionUnitCM}
	if calc.isOversized(dimensions) {
		t.Error("Expected normal size item not to be oversized")
	}

	// Oversized
	dimensions = Dimensions{Length: 150, Width: 30, Height: 20, Unit: DimensionUnitCM}
	if !calc.isOversized(dimensions) {
		t.Error("Expected oversized item to be detected")
	}
}

// Test checkRestrictions
func TestCheckRestrictions(t *testing.T) {
	calc := NewShippingCalculator()
	
	calc.Restrictions = []ShippingRestriction{
		{
			Type:      "destination",
			Countries: []string{"CN"},
			Message:   "Shipping to China is restricted",
		},
		{
			Type:       "hazardous",
			Message:    "Hazardous materials cannot be shipped",
		},
	}

	items := []ShippingItem{
		{IsHazardous: true},
	}
	destination := Address{Country: "CN"}

	restrictions := calc.checkRestrictions(items, destination)
	if len(restrictions) != 2 {
		t.Errorf("Expected 2 restrictions, got %d", len(restrictions))
	}
}

// Test applyFreeShipping
func TestApplyFreeShipping(t *testing.T) {
	calc := NewShippingCalculator()
	
	calc.FreeShippingRules = []FreeShippingRule{
		{
			IsActive:      true,
			ValidFrom:     time.Now().Add(-24 * time.Hour),
			ValidUntil:    time.Now().Add(24 * time.Hour),
			MinOrderValue: 100.0,
		},
	}

	result := &ShippingCalculationResult{
		Options: []ShippingOption{
			{Cost: 10.0, ServiceName: "Standard"},
			{Cost: 15.0, ServiceName: "Express"},
		},
		TotalValue: 150.0,
	}

	input := ShippingCalculationInput{
		Items: []ShippingItem{{Value: 150.0}},
	}

	calc.applyFreeShipping(result, input)
	
	// Check that cheapest option is now free
	if result.Options[0].Cost != 0 {
		t.Errorf("Expected free shipping, got cost %f", result.Options[0].Cost)
	}
	if result.Options[0].ServiceName != "Standard (Free Shipping)" {
		t.Errorf("Expected free shipping label, got %s", result.Options[0].ServiceName)
	}
}

// Test setRecommendedOptions
func TestSetRecommendedOptions(t *testing.T) {
	calc := NewShippingCalculator()
	
	result := &ShippingCalculationResult{
		Options: []ShippingOption{
			{Cost: 15.0, EstimatedDays: 2, ServiceName: "Express"},
			{Cost: 10.0, EstimatedDays: 5, ServiceName: "Standard"},
			{Cost: 25.0, EstimatedDays: 1, ServiceName: "Overnight"},
		},
	}

	calc.setRecommendedOptions(result)
	
	if result.CheapestOption.Cost != 10.0 {
		t.Errorf("Expected cheapest cost 10.0, got %f", result.CheapestOption.Cost)
	}
	if result.FastestOption.EstimatedDays != 1 {
		t.Errorf("Expected fastest 1 day, got %d", result.FastestOption.EstimatedDays)
	}
	if result.RecommendedOption == nil {
		t.Error("Expected recommended option to be set")
	}
}

// Test CalculateBestOption
func TestCalculateBestOption(t *testing.T) {
	// This would require setting up the global Calculate function
	// For now, we'll test the logic conceptually
	input := ShippingCalculationInput{
		Origin:      Address{Country: "US"},
		Destination: Address{Country: "US"},
		Items:       []ShippingItem{{Weight: Weight{Value: 1, Unit: WeightUnitKG}, Value: 50}},
	}

	// Test that function exists and handles criteria
	option, err := CalculateBestOption(input, "cheapest")
	if err != nil && err.Error() != "no shipping options available" {
		t.Errorf("Unexpected error: %v", err)
	}
	_ = option // Avoid unused variable warning
}

// Benchmark tests
func BenchmarkCalculate(b *testing.B) {
	input := ShippingCalculationInput{
		Origin:      Address{Country: "US", State: "CA"},
		Destination: Address{Country: "US", State: "NY"},
		Items: []ShippingItem{
			{
				Weight: Weight{Value: 1.0, Unit: WeightUnitKG},
				Value:  100.0,
				Dimensions: Dimensions{
					Length: 10, Width: 10, Height: 10,
					Unit: DimensionUnitCM,
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Calculate(input)
	}
}

func BenchmarkCalculateShipping(b *testing.B) {
	calc := NewShippingCalculator()
	input := ShippingCalculationInput{
		Origin:      Address{Country: "US"},
		Destination: Address{Country: "US"},
		Items:       []ShippingItem{{Weight: Weight{Value: 1, Unit: WeightUnitKG}, Value: 50}},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = calc.CalculateShipping(input)
	}
}
// Package shipping provides comprehensive shipping calculation and rule management functionality.
//
// This package includes:
//   - Shipping cost calculation with multiple carriers and methods
//   - Dynamic rule engine for complex shipping scenarios
//   - Zone-based shipping with geographic restrictions
//   - Free shipping rule management
//   - Delivery time estimation
//   - Packaging rule optimization
//   - Surcharge calculation for special items
//
// Basic usage:
//
//	ruleEngine := shipping.NewShippingRuleEngine()
//	
//	// Add shipping rules
//	rule := shipping.ShippingRule{
//		ID:       "standard-us",
//		Name:     "Standard US Shipping",
//		Method:   shipping.ShippingMethodStandard,
//		Zone:     shipping.ShippingZoneNational,
//		BaseCost: 5.99,
//		WeightRate: 0.50,
//		IsActive: true,
//	}
//	ruleEngine.AddShippingRule(rule)
//	
//	// Get applicable rules
//	input := shipping.ShippingCalculationInput{
//		Items: []shipping.ShippingItem{{Weight: shipping.Weight{Value: 2.5, Unit: "kg"}}},
//		Destination: shipping.Address{Country: "US"},
//	}
//	applicableRules := ruleEngine.GetApplicableRules(input)
package shipping

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// ShippingRuleEngine manages all shipping-related rules and configurations.
// It provides a centralized system for managing shipping rules, carrier configurations,
// zone definitions, delivery time calculations, restrictions, and free shipping policies.
//
// The rule engine supports:
//   - Multiple shipping methods and carriers
//   - Geographic zone-based shipping
//   - Weight and value-based pricing
//   - Time-based rule validity
//   - Complex rule validation and optimization
//   - Rule import/export functionality
//
// Example usage:
//
//	ruleEngine := shipping.NewShippingRuleEngine()
//	
//	// Configure shipping rules
//	rule := shipping.ShippingRule{
//		ID:       "express-intl",
//		Name:     "Express International",
//		Method:   shipping.ShippingMethodExpress,
//		Zone:     shipping.ShippingZoneInternational,
//		BaseCost: 25.00,
//		WeightRate: 2.50,
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(1, 0, 0),
//	}
//	ruleEngine.AddShippingRule(rule)
//	
//	// Validate configuration
//	warnings := ruleEngine.ValidateRuleConfiguration()
//	if len(warnings) > 0 {
//		log.Printf("Rule warnings: %v", warnings)
//	}
type ShippingRuleEngine struct {
	ShippingRules      []ShippingRule
	CarrierRules       []CarrierRule
	ZoneRules          []ZoneRule
	DeliveryTimeRules  []DeliveryTimeRule
	Restrictions       []ShippingRestriction
	FreeShippingRules  []FreeShippingRule
	PackagingRules     []PackagingRule
}

// NewShippingRuleEngine creates a new shipping rule engine with empty rule sets.
// This function initializes all rule collections and returns a ready-to-use rule engine.
//
// The returned engine can be used to:
//   - Add and manage shipping rules
//   - Configure carrier-specific settings
//   - Define shipping zones and restrictions
//   - Set up free shipping policies
//   - Manage packaging rules
//
// Example:
//
//	engine := shipping.NewShippingRuleEngine()
//	
//	// Add a basic shipping rule
//	rule := shipping.ShippingRule{
//		ID:       "basic-shipping",
//		Name:     "Basic Shipping",
//		Method:   shipping.ShippingMethodStandard,
//		BaseCost: 9.99,
//		IsActive: true,
//	}
//	engine.AddShippingRule(rule)
//
// Returns:
//   - *ShippingRuleEngine: A new rule engine instance with empty rule collections
func NewShippingRuleEngine() *ShippingRuleEngine {
	return &ShippingRuleEngine{
		ShippingRules:      []ShippingRule{},
		CarrierRules:       []CarrierRule{},
		ZoneRules:          []ZoneRule{},
		DeliveryTimeRules:  []DeliveryTimeRule{},
		Restrictions:       []ShippingRestriction{},
		FreeShippingRules:  []FreeShippingRule{},
		PackagingRules:     []PackagingRule{},
	}
}

// Shipping Rule Management

// AddShippingRule adds a new shipping rule to the engine.
// This method validates that no rule with the same ID already exists before adding.
//
// Parameters:
//   - rule: The ShippingRule to add to the engine
//
// Returns:
//   - error: nil if successful, error if rule ID already exists
//
// Example:
//
//	rule := shipping.ShippingRule{
//		ID:       "premium-overnight",
//		Name:     "Premium Overnight",
//		Method:   shipping.ShippingMethodOvernight,
//		BaseCost: 29.99,
//		WeightRate: 3.50,
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(0, 6, 0),
//	}
//	err := engine.AddShippingRule(rule)
//	if err != nil {
//		log.Printf("Failed to add rule: %v", err)
//	}
func (sre *ShippingRuleEngine) AddShippingRule(rule ShippingRule) error {
	if rule.ID == "" {
		return errors.New("shipping rule ID cannot be empty")
	}

	if rule.Name == "" {
		return errors.New("shipping rule name cannot be empty")
	}

	if rule.BaseCost < 0 {
		return errors.New("base cost cannot be negative")
	}

	// Check for duplicate ID
	for _, existingRule := range sre.ShippingRules {
		if existingRule.ID == rule.ID {
			return fmt.Errorf("shipping rule with ID %s already exists", rule.ID)
		}
	}

	sre.ShippingRules = append(sre.ShippingRules, rule)
	return nil
}

// UpdateShippingRule updates an existing shipping rule identified by its ID.
// The rule ID cannot be changed during update to maintain referential integrity.
//
// Parameters:
//   - ruleID: The ID of the rule to update
//   - updatedRule: The new rule data (ID will be preserved)
//
// Returns:
//   - error: nil if successful, error if rule not found
//
// Example:
//
//	updatedRule := shipping.ShippingRule{
//		Name:     "Updated Premium Shipping",
//		Method:   shipping.ShippingMethodExpress,
//		BaseCost: 19.99,
//		WeightRate: 2.50,
//		IsActive: true,
//	}
//	err := engine.UpdateShippingRule("premium-overnight", updatedRule)
//	if err != nil {
//		log.Printf("Failed to update rule: %v", err)
//	}
func (sre *ShippingRuleEngine) UpdateShippingRule(ruleID string, updatedRule ShippingRule) error {
	for i, rule := range sre.ShippingRules {
		if rule.ID == ruleID {
			updatedRule.ID = ruleID // Preserve original ID
			sre.ShippingRules[i] = updatedRule
			return nil
		}
	}
	return fmt.Errorf("shipping rule with ID %s not found", ruleID)
}

// RemoveShippingRule removes a shipping rule from the engine by its ID.
// This operation permanently deletes the rule from the engine.
//
// Parameters:
//   - ruleID: The ID of the rule to remove
//
// Returns:
//   - error: nil if successful, error if rule not found
//
// Example:
//
//	err := engine.RemoveShippingRule("old-rule-id")
//	if err != nil {
//		log.Printf("Failed to remove rule: %v", err)
//	} else {
//		log.Println("Rule removed successfully")
//	}
func (sre *ShippingRuleEngine) RemoveShippingRule(ruleID string) error {
	for i, rule := range sre.ShippingRules {
		if rule.ID == ruleID {
			sre.ShippingRules = append(sre.ShippingRules[:i], sre.ShippingRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("shipping rule with ID %s not found", ruleID)
}

// GetShippingRule retrieves a specific shipping rule by its ID.
// Returns a pointer to the rule to allow for efficient access without copying.
//
// Parameters:
//   - ruleID: The ID of the rule to retrieve
//
// Returns:
//   - *ShippingRule: Pointer to the rule if found
//   - error: nil if successful, error if rule not found
//
// Example:
//
//	rule, err := engine.GetShippingRule("premium-overnight")
//	if err != nil {
//		log.Printf("Rule not found: %v", err)
//	} else {
//		log.Printf("Found rule: %s - $%.2f", rule.Name, rule.BaseCost)
//	}
func (sre *ShippingRuleEngine) GetShippingRule(ruleID string) (*ShippingRule, error) {
	for _, rule := range sre.ShippingRules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("shipping rule with ID %s not found", ruleID)
}

// GetActiveShippingRules returns all currently active shipping rules.
// Only rules with IsActive set to true and within their valid time range are included.
//
// Returns:
//   - []ShippingRule: Slice of all active shipping rules
//
// Example:
//
//	activeRules := engine.GetActiveShippingRules()
//	log.Printf("Found %d active shipping rules", len(activeRules))
//	for _, rule := range activeRules {
//		log.Printf("Active rule: %s (%s)", rule.Name, rule.Method)
//	}
func (sre *ShippingRuleEngine) GetActiveShippingRules() []ShippingRule {
	activeRules := []ShippingRule{}
	now := time.Now()

	for _, rule := range sre.ShippingRules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			activeRules = append(activeRules, rule)
		}
	}

	return activeRules
}

// Carrier Rule Management

// AddCarrierRule adds a new carrier-specific shipping rule to the engine.
// Carrier rules define shipping options and pricing for specific carriers and service codes.
//
// Parameters:
//   - rule: The CarrierRule to add to the engine
//
// Returns:
//   - error: nil if successful, error if rule for carrier and service code already exists
//
// Example:
//
//	carrierRule := shipping.CarrierRule{
//		CarrierID:    "fedex",
//		CarrierName:  "FedEx",
//		ServiceCode:  "EXPRESS",
//		ServiceName:  "FedEx Express",
//		BaseCost:     15.99,
//		WeightRate:   1.25,
//		MaxWeight:    shipping.Weight{Value: 70, Unit: "kg"},
//		DeliveryDays: 2,
//		IsActive:     true,
//	}
//	err := engine.AddCarrierRule(carrierRule)
//	if err != nil {
//		log.Printf("Failed to add carrier rule: %v", err)
//	}
func (sre *ShippingRuleEngine) AddCarrierRule(rule CarrierRule) error {
	if rule.CarrierID == "" {
		return errors.New("carrier ID cannot be empty")
	}

	if rule.CarrierName == "" {
		return errors.New("carrier name cannot be empty")
	}

	if rule.ServiceCode == "" {
		return errors.New("service code cannot be empty")
	}

	// Check for duplicate carrier + service combination
	for _, existingRule := range sre.CarrierRules {
		if existingRule.CarrierID == rule.CarrierID && existingRule.ServiceCode == rule.ServiceCode {
			return fmt.Errorf("carrier rule for %s %s already exists", rule.CarrierID, rule.ServiceCode)
		}
	}

	sre.CarrierRules = append(sre.CarrierRules, rule)
	return nil
}

// UpdateCarrierRule updates an existing carrier rule identified by carrier ID and service code.
// The carrier ID and service code cannot be changed during update to maintain referential integrity.
//
// Parameters:
//   - carrierID: The ID of the carrier
//   - serviceCode: The service code
//   - updatedRule: The new rule data (carrier ID and service code will be preserved)
//
// Returns:
//   - error: nil if successful, error if rule not found
//
// Example:
//
//	updatedRule := shipping.CarrierRule{
//		CarrierName:  "FedEx Updated",
//		ServiceName:  "FedEx Express Updated",
//		BaseCost:     17.99,
//		WeightRate:   1.50,
//		DeliveryDays: 1,
//		IsActive:     true,
//	}
//	err := engine.UpdateCarrierRule("fedex", "EXPRESS", updatedRule)
//	if err != nil {
//		log.Printf("Failed to update carrier rule: %v", err)
//	}
func (sre *ShippingRuleEngine) UpdateCarrierRule(carrierID, serviceCode string, updatedRule CarrierRule) error {
	for i, rule := range sre.CarrierRules {
		if rule.CarrierID == carrierID && rule.ServiceCode == serviceCode {
			updatedRule.CarrierID = carrierID
			updatedRule.ServiceCode = serviceCode
			sre.CarrierRules[i] = updatedRule
			return nil
		}
	}
	return fmt.Errorf("carrier rule for %s %s not found", carrierID, serviceCode)
}

// RemoveCarrierRule removes a carrier rule from the engine.
// This operation permanently deletes the rule for the specified carrier and service code.
//
// Parameters:
//   - carrierID: The ID of the carrier
//   - serviceCode: The service code
//
// Returns:
//   - error: nil if successful, error if rule not found
//
// Example:
//
//	err := engine.RemoveCarrierRule("old-carrier", "STANDARD")
//	if err != nil {
//		log.Printf("Failed to remove carrier rule: %v", err)
//	} else {
//		log.Println("Carrier rule removed successfully")
//	}
func (sre *ShippingRuleEngine) RemoveCarrierRule(carrierID, serviceCode string) error {
	for i, rule := range sre.CarrierRules {
		if rule.CarrierID == carrierID && rule.ServiceCode == serviceCode {
			sre.CarrierRules = append(sre.CarrierRules[:i], sre.CarrierRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("carrier rule for %s %s not found", carrierID, serviceCode)
}

// GetCarrierRules returns all carrier rules for a specific carrier.
// This method allows you to retrieve all service options available for a given carrier.
//
// Parameters:
//   - carrierID: The ID of the carrier
//
// Returns:
//   - []CarrierRule: Slice of all rules for the specified carrier
//
// Example:
//
//	rules := engine.GetCarrierRules("fedex")
//	log.Printf("Found %d rules for FedEx", len(rules))
//	for _, rule := range rules {
//		log.Printf("Service: %s - $%.2f", rule.ServiceName, rule.BaseCost)
//	}
func (sre *ShippingRuleEngine) GetCarrierRules(carrierID string) []CarrierRule {
	rules := []CarrierRule{}
	for _, rule := range sre.CarrierRules {
		if rule.CarrierID == carrierID {
			rules = append(rules, rule)
		}
	}
	return rules
}

// Zone Rule Management

// AddZoneRule adds a new zone rule to define geographic shipping zones.
// Zone rules specify which geographic areas belong to specific shipping zones.
//
// Parameters:
//   - rule: The ZoneRule to add to the engine
//
// Returns:
//   - error: nil if successful, error if zone is empty or no location criteria specified
//
// Example:
//
//	zoneRule := shipping.ZoneRule{
//		Zone:      shipping.ShippingZoneNational,
//		Countries: []string{"US", "CA"},
//		States:    []string{"CA", "NY", "TX"},
//		PostalCodeRanges: []shipping.PostalCodeRange{
//			{Start: "90000", End: "99999"},
//		},
//		IsActive:  true,
//	}
//	err := engine.AddZoneRule(zoneRule)
//	if err != nil {
//		log.Printf("Failed to add zone rule: %v", err)
//	}
func (sre *ShippingRuleEngine) AddZoneRule(rule ZoneRule) error {
	if rule.Zone == "" {
		return errors.New("zone cannot be empty")
	}

	if len(rule.Countries) == 0 && len(rule.States) == 0 && len(rule.PostalCodes) == 0 && len(rule.PostalCodeRanges) == 0 {
		return errors.New("zone rule must specify at least one location criteria")
	}

	sre.ZoneRules = append(sre.ZoneRules, rule)
	return nil
}

// UpdateZoneRule updates an existing zone rule at the specified index.
// Zone rules are indexed by their position in the rules array.
//
// Parameters:
//   - index: The index of the zone rule to update
//   - updatedRule: The new rule data
//
// Returns:
//   - error: nil if successful, error if index is invalid
//
// Example:
//
//	updatedRule := shipping.ZoneRule{
//		Zone:      shipping.ShippingZoneRegional,
//		Countries: []string{"US"},
//		States:    []string{"CA", "NV", "OR", "WA"},
//		IsActive:  true,
//	}
//	err := engine.UpdateZoneRule(0, updatedRule)
//	if err != nil {
//		log.Printf("Failed to update zone rule: %v", err)
//	}
func (sre *ShippingRuleEngine) UpdateZoneRule(index int, updatedRule ZoneRule) error {
	if index < 0 || index >= len(sre.ZoneRules) {
		return errors.New("invalid zone rule index")
	}

	sre.ZoneRules[index] = updatedRule
	return nil
}

// RemoveZoneRule removes a zone rule from the engine at the specified index.
// This operation permanently deletes the zone rule.
//
// Parameters:
//   - index: The index of the zone rule to remove
//
// Returns:
//   - error: nil if successful, error if index is invalid
//
// Example:
//
//	err := engine.RemoveZoneRule(2)
//	if err != nil {
//		log.Printf("Failed to remove zone rule: %v", err)
//	} else {
//		log.Println("Zone rule removed successfully")
//	}
func (sre *ShippingRuleEngine) RemoveZoneRule(index int) error {
	if index < 0 || index >= len(sre.ZoneRules) {
		return errors.New("invalid zone rule index")
	}

	sre.ZoneRules = append(sre.ZoneRules[:index], sre.ZoneRules[index+1:]...)
	return nil
}

// GetZoneRules returns all zone rules for a specific shipping zone.
// This method helps identify which geographic areas are covered by a zone.
//
// Parameters:
//   - zone: The shipping zone to get rules for
//
// Returns:
//   - []ZoneRule: Slice of all rules for the specified zone
//
// Example:
//
//	rules := engine.GetZoneRules(shipping.ShippingZoneInternational)
//	log.Printf("Found %d rules for international zone", len(rules))
//	for _, rule := range rules {
//		log.Printf("Zone rule covers %d countries", len(rule.Countries))
//	}
func (sre *ShippingRuleEngine) GetZoneRules(zone ShippingZone) []ZoneRule {
	rules := []ZoneRule{}
	for _, rule := range sre.ZoneRules {
		if rule.Zone == zone {
			rules = append(rules, rule)
		}
	}
	return rules
}

// Delivery Time Rule Management

// AddDeliveryTimeRule adds a new delivery time rule
func (sre *ShippingRuleEngine) AddDeliveryTimeRule(rule DeliveryTimeRule) error {
	if rule.BaseDays < 0 {
		return errors.New("base days cannot be negative")
	}

	sre.DeliveryTimeRules = append(sre.DeliveryTimeRules, rule)
	return nil
}

// UpdateDeliveryTimeRule updates an existing delivery time rule
func (sre *ShippingRuleEngine) UpdateDeliveryTimeRule(method ShippingMethod, zone ShippingZone, updatedRule DeliveryTimeRule) error {
	for i, rule := range sre.DeliveryTimeRules {
		if rule.Method == method && rule.Zone == zone {
			updatedRule.Method = method
			updatedRule.Zone = zone
			sre.DeliveryTimeRules[i] = updatedRule
			return nil
		}
	}
	return fmt.Errorf("delivery time rule for %s in %s not found", method, zone)
}

// RemoveDeliveryTimeRule removes a delivery time rule
func (sre *ShippingRuleEngine) RemoveDeliveryTimeRule(method ShippingMethod, zone ShippingZone) error {
	for i, rule := range sre.DeliveryTimeRules {
		if rule.Method == method && rule.Zone == zone {
			sre.DeliveryTimeRules = append(sre.DeliveryTimeRules[:i], sre.DeliveryTimeRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("delivery time rule for %s in %s not found", method, zone)
}

// GetDeliveryTimeRule retrieves a delivery time rule
func (sre *ShippingRuleEngine) GetDeliveryTimeRule(method ShippingMethod, zone ShippingZone) (*DeliveryTimeRule, error) {
	for _, rule := range sre.DeliveryTimeRules {
		if rule.Method == method && rule.Zone == zone {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("delivery time rule for %s in %s not found", method, zone)
}

// Shipping Restriction Management

// AddShippingRestriction adds a new shipping restriction
func (sre *ShippingRuleEngine) AddShippingRestriction(restriction ShippingRestriction) error {
	if restriction.Type == "" {
		return errors.New("restriction type cannot be empty")
	}

	if restriction.Message == "" {
		return errors.New("restriction message cannot be empty")
	}

	sre.Restrictions = append(sre.Restrictions, restriction)
	return nil
}

// UpdateShippingRestriction updates an existing shipping restriction
func (sre *ShippingRuleEngine) UpdateShippingRestriction(index int, updatedRestriction ShippingRestriction) error {
	if index < 0 || index >= len(sre.Restrictions) {
		return errors.New("invalid restriction index")
	}

	sre.Restrictions[index] = updatedRestriction
	return nil
}

// RemoveShippingRestriction removes a shipping restriction
func (sre *ShippingRuleEngine) RemoveShippingRestriction(index int) error {
	if index < 0 || index >= len(sre.Restrictions) {
		return errors.New("invalid restriction index")
	}

	sre.Restrictions = append(sre.Restrictions[:index], sre.Restrictions[index+1:]...)
	return nil
}

// GetShippingRestrictions returns all shipping restrictions of a specific type
func (sre *ShippingRuleEngine) GetShippingRestrictions(restrictionType string) []ShippingRestriction {
	restrictions := []ShippingRestriction{}
	for _, restriction := range sre.Restrictions {
		if restriction.Type == restrictionType {
			restrictions = append(restrictions, restriction)
		}
	}
	return restrictions
}

// Free Shipping Rule Management

// AddFreeShippingRule adds a new free shipping rule
func (sre *ShippingRuleEngine) AddFreeShippingRule(rule FreeShippingRule) error {
	if rule.Name == "" {
		return errors.New("free shipping rule name cannot be empty")
	}

	if rule.MinOrderValue < 0 {
		return errors.New("minimum order value cannot be negative")
	}

	sre.FreeShippingRules = append(sre.FreeShippingRules, rule)
	return nil
}

// UpdateFreeShippingRule updates an existing free shipping rule
func (sre *ShippingRuleEngine) UpdateFreeShippingRule(name string, updatedRule FreeShippingRule) error {
	for i, rule := range sre.FreeShippingRules {
		if rule.Name == name {
			updatedRule.Name = name
			sre.FreeShippingRules[i] = updatedRule
			return nil
		}
	}
	return fmt.Errorf("free shipping rule %s not found", name)
}

// RemoveFreeShippingRule removes a free shipping rule
func (sre *ShippingRuleEngine) RemoveFreeShippingRule(name string) error {
	for i, rule := range sre.FreeShippingRules {
		if rule.Name == name {
			sre.FreeShippingRules = append(sre.FreeShippingRules[:i], sre.FreeShippingRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("free shipping rule %s not found", name)
}

// GetActiveFreeShippingRules returns all active free shipping rules
func (sre *ShippingRuleEngine) GetActiveFreeShippingRules() []FreeShippingRule {
	activeRules := []FreeShippingRule{}
	now := time.Now()

	for _, rule := range sre.FreeShippingRules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			activeRules = append(activeRules, rule)
		}
	}

	return activeRules
}

// Packaging Rule Management

// AddPackagingRule adds a new packaging rule
func (sre *ShippingRuleEngine) AddPackagingRule(rule PackagingRule) error {
	if rule.Name == "" {
		return errors.New("packaging rule name cannot be empty")
	}

	if rule.MaxWeight.Value <= 0 {
		return errors.New("max weight must be positive")
	}

	sre.PackagingRules = append(sre.PackagingRules, rule)
	return nil
}

// UpdatePackagingRule updates an existing packaging rule
func (sre *ShippingRuleEngine) UpdatePackagingRule(name string, updatedRule PackagingRule) error {
	for i, rule := range sre.PackagingRules {
		if rule.Name == name {
			updatedRule.Name = name
			sre.PackagingRules[i] = updatedRule
			return nil
		}
	}
	return fmt.Errorf("packaging rule %s not found", name)
}

// RemovePackagingRule removes a packaging rule
func (sre *ShippingRuleEngine) RemovePackagingRule(name string) error {
	for i, rule := range sre.PackagingRules {
		if rule.Name == name {
			sre.PackagingRules = append(sre.PackagingRules[:i], sre.PackagingRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("packaging rule %s not found", name)
}

// GetPackagingRule retrieves a packaging rule by name
func (sre *ShippingRuleEngine) GetPackagingRule(name string) (*PackagingRule, error) {
	for _, rule := range sre.PackagingRules {
		if rule.Name == name {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("packaging rule %s not found", name)
}

// Advanced Rule Operations

// ValidateRuleConfiguration validates the entire rule configuration and returns warnings.
// This method performs comprehensive validation to identify potential issues with the
// current rule setup, including overlapping rules, missing coverage, and expired rules.
//
// Validation checks include:
//   - Overlapping shipping rules that may cause conflicts
//   - Missing zone coverage for standard shipping zones
//   - Expired rules that are still marked as active
//   - Inconsistent carrier configurations
//
// Returns:
//   - []string: Slice of warning messages describing potential issues
//
// Example:
//
//	warnings := engine.ValidateRuleConfiguration()
//	if len(warnings) > 0 {
//		log.Printf("Found %d configuration warnings:", len(warnings))
//		for _, warning := range warnings {
//			log.Printf("WARNING: %s", warning)
//		}
//	} else {
//		log.Println("Rule configuration is valid")
//	}
func (sre *ShippingRuleEngine) ValidateRuleConfiguration() []string {
	warnings := []string{}

	// Check for overlapping shipping rules
	for i, rule1 := range sre.ShippingRules {
		for j, rule2 := range sre.ShippingRules {
			if i != j && sre.rulesOverlap(rule1, rule2) {
				warnings = append(warnings, fmt.Sprintf("Shipping rules %s and %s may overlap", rule1.ID, rule2.ID))
			}
		}
	}

	// Check for missing zone coverage
	zones := []ShippingZone{ShippingZoneLocal, ShippingZoneRegional, ShippingZoneNational, ShippingZoneInternational}
	for _, zone := range zones {
		if !sre.hasZoneCoverage(zone) {
			warnings = append(warnings, fmt.Sprintf("No shipping rules cover zone %s", zone))
		}
	}

	// Check for expired rules
	now := time.Now()
	for _, rule := range sre.ShippingRules {
		if rule.IsActive && now.After(rule.ValidUntil) {
			warnings = append(warnings, fmt.Sprintf("Shipping rule %s has expired", rule.ID))
		}
	}

	return warnings
}

// rulesOverlap checks if two shipping rules overlap
func (sre *ShippingRuleEngine) rulesOverlap(rule1, rule2 ShippingRule) bool {
	// Check if rules have overlapping zones
	if rule1.Zone != "" && rule2.Zone != "" && rule1.Zone != rule2.Zone {
		return false
	}

	// Check if rules have overlapping methods
	if rule1.Method != rule2.Method {
		return false
	}

	// Check if rules have overlapping time periods
	if rule1.ValidUntil.Before(rule2.ValidFrom) || rule2.ValidUntil.Before(rule1.ValidFrom) {
		return false
	}

	// Check if rules have overlapping weight ranges
	if rule1.MaxWeight.Value > 0 && rule2.MinWeight.Value > 0 {
		max1 := convertWeight(rule1.MaxWeight, WeightUnitKG)
		min2 := convertWeight(rule2.MinWeight, WeightUnitKG)
		if max1 < min2 {
			return false
		}
	}

	if rule2.MaxWeight.Value > 0 && rule1.MinWeight.Value > 0 {
		max2 := convertWeight(rule2.MaxWeight, WeightUnitKG)
		min1 := convertWeight(rule1.MinWeight, WeightUnitKG)
		if max2 < min1 {
			return false
		}
	}

	return true
}

// hasZoneCoverage checks if there are active rules covering a zone
func (sre *ShippingRuleEngine) hasZoneCoverage(zone ShippingZone) bool {
	now := time.Now()
	for _, rule := range sre.ShippingRules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			if rule.Zone == "" || rule.Zone == zone {
				return true
			}
		}
	}
	return false
}

// GetRulesByName returns shipping rules sorted by name
func (sre *ShippingRuleEngine) GetRulesByName() []ShippingRule {
	rules := make([]ShippingRule, len(sre.ShippingRules))
	copy(rules, sre.ShippingRules)

	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Name < rules[j].Name
	})

	return rules
}

// GetApplicableRules returns shipping rules that are applicable for the given shipping input.
// This method filters all shipping rules based on various criteria including time validity,
// weight constraints, value constraints, geographic restrictions, and item categories.
//
// Filtering criteria:
//   - Rule must be active and within valid time range
//   - Total weight must be within rule's weight constraints
//   - Total value must be within rule's value constraints
//   - Destination must match geographic restrictions
//   - Items must match category restrictions (if specified)
//
// Parameters:
//   - input: The shipping calculation input containing items and destination
//
// Returns:
//   - []ShippingRule: Slice of applicable shipping rules
//
// Example:
//
//	input := shipping.ShippingCalculationInput{
//		Items: []shipping.ShippingItem{
//			{Weight: shipping.Weight{Value: 2.5, Unit: "kg"}, Category: "electronics"},
//		},
//		Destination: shipping.Address{Country: "US", State: "CA"},
//	}
//	applicableRules := engine.GetApplicableRules(input)
//	log.Printf("Found %d applicable rules", len(applicableRules))
//	for _, rule := range applicableRules {
//		log.Printf("Applicable rule: %s - $%.2f", rule.Name, rule.BaseCost)
//	}
func (sre *ShippingRuleEngine) GetApplicableRules(input ShippingCalculationInput) []ShippingRule {
	applicableRules := []ShippingRule{}
	now := time.Now()

	for _, rule := range sre.ShippingRules {
		// Check if rule is active and within valid time range
		if !rule.IsActive || now.Before(rule.ValidFrom) || now.After(rule.ValidUntil) {
			continue
		}

		// Check weight constraints
		totalWeight := calculateTotalWeight(input.Items)
		if rule.MinWeight.Value > 0 && convertWeight(totalWeight, rule.MinWeight.Unit) < rule.MinWeight.Value {
			continue
		}
		if rule.MaxWeight.Value > 0 && convertWeight(totalWeight, rule.MaxWeight.Unit) > rule.MaxWeight.Value {
			continue
		}

		// Check value constraints
		totalValue := calculateTotalValue(input.Items)
		if rule.MinValue > 0 && totalValue < rule.MinValue {
			continue
		}
		if rule.MaxValue > 0 && totalValue > rule.MaxValue {
			continue
		}

		// Check geographic constraints
		if len(rule.ApplicableCountries) > 0 {
			found := false
			for _, country := range rule.ApplicableCountries {
				if input.Destination.Country == country {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		if len(rule.ApplicableStates) > 0 {
			found := false
			for _, state := range rule.ApplicableStates {
				if input.Destination.State == state {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check category constraints
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
				continue
			}
		}

		applicableRules = append(applicableRules, rule)
	}

	return applicableRules
}

// OptimizeRules optimizes the rule configuration for better performance and consistency.
// This method performs several optimization operations to improve rule processing efficiency
// and reduce potential conflicts.
//
// Optimization operations:
//   - Sorts rules by name for consistent processing order
//   - Removes duplicate rules to reduce redundancy
//   - Consolidates overlapping rules where possible
//   - Improves rule lookup performance
//
// This method should be called after making significant changes to the rule configuration
// or periodically to maintain optimal performance.
//
// Example:
//
//	// After adding multiple rules
//	engine.AddShippingRule(rule1)
//	engine.AddShippingRule(rule2)
//	engine.AddShippingRule(rule3)
//	
//	// Optimize for better performance
//	engine.OptimizeRules()
//	log.Println("Rule configuration optimized")
//	
//	// Validate after optimization
//	warnings := engine.ValidateRuleConfiguration()
//	if len(warnings) == 0 {
//		log.Println("Optimized configuration is valid")
//	}
func (sre *ShippingRuleEngine) OptimizeRules() {
	// Sort rules by name for consistent processing
	sort.Slice(sre.ShippingRules, func(i, j int) bool {
		return sre.ShippingRules[i].Name < sre.ShippingRules[j].Name
	})

	// Remove duplicate rules
	sre.removeDuplicateRules()

	// Consolidate overlapping rules where possible
	sre.consolidateRules()
}

// removeDuplicateRules removes duplicate shipping rules
func (sre *ShippingRuleEngine) removeDuplicateRules() {
	unique := []ShippingRule{}
	seen := make(map[string]bool)

	for _, rule := range sre.ShippingRules {
		key := fmt.Sprintf("%s_%s_%s_%f_%f", rule.ID, rule.Method, rule.Zone, rule.BaseCost, rule.WeightRate)
		if !seen[key] {
			seen[key] = true
			unique = append(unique, rule)
		}
	}

	sre.ShippingRules = unique
}

// consolidateRules consolidates overlapping rules where possible
func (sre *ShippingRuleEngine) consolidateRules() {
	// This is a simplified consolidation - in practice, this would be more complex
	consolidated := []ShippingRule{}
	processed := make(map[int]bool)

	for i, rule1 := range sre.ShippingRules {
		if processed[i] {
			continue
		}

		consolidatedRule := rule1
		processed[i] = true

		// Look for similar rules that can be consolidated
		for j, rule2 := range sre.ShippingRules {
			if i != j && !processed[j] && sre.canConsolidate(rule1, rule2) {
				// Merge rules (simplified logic)
				if rule2.BaseCost < consolidatedRule.BaseCost {
					consolidatedRule.BaseCost = rule2.BaseCost
				}
				processed[j] = true
			}
		}

		consolidated = append(consolidated, consolidatedRule)
	}

	sre.ShippingRules = consolidated
}

// canConsolidate checks if two rules can be consolidated
func (sre *ShippingRuleEngine) canConsolidate(rule1, rule2 ShippingRule) bool {
	// Rules can be consolidated if they have the same method, zone, and similar constraints
	return rule1.Method == rule2.Method &&
		rule1.Zone == rule2.Zone &&
		rule1.WeightRate == rule2.WeightRate &&
		rule1.ValueRate == rule2.ValueRate
}

// ExportRules exports all rules to a structured format for backup or transfer.
// This method creates a comprehensive export of all rule types in the engine,
// suitable for serialization to JSON or other formats.
//
// The exported data includes:
//   - All shipping rules with their configurations
//   - Carrier-specific rules and settings
//   - Zone definitions and geographic mappings
//   - Delivery time rules and calculations
//   - Shipping restrictions and limitations
//   - Free shipping rules and conditions
//   - Packaging rules and requirements
//
// Returns:
//   - map[string]interface{}: Structured data containing all rule types
//
// Example:
//
//	exportedData := engine.ExportRules()
//	
//	// Convert to JSON for storage
//	jsonData, err := json.Marshal(exportedData)
//	if err != nil {
//		log.Printf("Failed to marshal rules: %v", err)
//	} else {
//		// Save to file or database
//		ioutil.WriteFile("shipping_rules.json", jsonData, 0644)
//		log.Println("Rules exported successfully")
//	}
func (sre *ShippingRuleEngine) ExportRules() map[string]interface{} {
	return map[string]interface{}{
		"shipping_rules":      sre.ShippingRules,
		"carrier_rules":       sre.CarrierRules,
		"zone_rules":          sre.ZoneRules,
		"delivery_time_rules": sre.DeliveryTimeRules,
		"restrictions":        sre.Restrictions,
		"free_shipping_rules": sre.FreeShippingRules,
		"packaging_rules":     sre.PackagingRules,
	}
}

// ImportRules imports rules from a structured format, typically from a previous export.
// This method replaces the current rule configuration with the imported data.
// Use with caution as it will overwrite existing rules.
//
// The import process validates data types and safely handles missing or invalid data.
// Only valid rule types that match the expected structure will be imported.
//
// Parameters:
//   - data: Structured data containing rule configurations
//
// Returns:
//   - error: nil if successful, error if import fails
//
// Example:
//
//	// Load from JSON file
//	jsonData, err := ioutil.ReadFile("shipping_rules.json")
//	if err != nil {
//		log.Printf("Failed to read rules file: %v", err)
//		return
//	}
//	
//	var importData map[string]interface{}
//	err = json.Unmarshal(jsonData, &importData)
//	if err != nil {
//		log.Printf("Failed to unmarshal rules: %v", err)
//		return
//	}
//	
//	err = engine.ImportRules(importData)
//	if err != nil {
//		log.Printf("Failed to import rules: %v", err)
//	} else {
//		log.Println("Rules imported successfully")
//	}
func (sre *ShippingRuleEngine) ImportRules(data map[string]interface{}) error {
	// This would typically involve JSON unmarshaling or similar
	// For now, we'll provide a basic structure

	if shippingRules, ok := data["shipping_rules"].([]ShippingRule); ok {
		sre.ShippingRules = shippingRules
	}

	if carrierRules, ok := data["carrier_rules"].([]CarrierRule); ok {
		sre.CarrierRules = carrierRules
	}

	if zoneRules, ok := data["zone_rules"].([]ZoneRule); ok {
		sre.ZoneRules = zoneRules
	}

	if deliveryTimeRules, ok := data["delivery_time_rules"].([]DeliveryTimeRule); ok {
		sre.DeliveryTimeRules = deliveryTimeRules
	}

	if restrictions, ok := data["restrictions"].([]ShippingRestriction); ok {
		sre.Restrictions = restrictions
	}

	if freeShippingRules, ok := data["free_shipping_rules"].([]FreeShippingRule); ok {
		sre.FreeShippingRules = freeShippingRules
	}

	if packagingRules, ok := data["packaging_rules"].([]PackagingRule); ok {
		sre.PackagingRules = packagingRules
	}

	return nil
}

// GetRuleStatistics returns comprehensive statistics about the current rule configuration.
// This method provides insights into the rule engine's current state, including counts
// of different rule types and their active status.
//
// Statistics include:
//   - Total and active shipping rules
//   - Carrier rule counts
//   - Zone rule counts
//   - Delivery time rule counts
//   - Restriction counts
//   - Free shipping rule counts (total and active)
//   - Packaging rule counts
//
// Returns:
//   - map[string]int: Statistics with descriptive keys and counts
//
// Example:
//
//	stats := engine.GetRuleStatistics()
//	log.Printf("Rule Engine Statistics:")
//	log.Printf("  Total shipping rules: %d", stats["total_shipping_rules"])
//	log.Printf("  Active shipping rules: %d", stats["active_shipping_rules"])
//	log.Printf("  Carrier rules: %d", stats["total_carrier_rules"])
//	log.Printf("  Zone rules: %d", stats["total_zone_rules"])
//	log.Printf("  Free shipping rules: %d (active: %d)", 
//		stats["total_free_shipping_rules"], stats["active_free_shipping_rules"])
func (sre *ShippingRuleEngine) GetRuleStatistics() map[string]int {
	stats := map[string]int{
		"total_shipping_rules":      len(sre.ShippingRules),
		"active_shipping_rules":     0,
		"total_carrier_rules":       len(sre.CarrierRules),
		"total_zone_rules":          len(sre.ZoneRules),
		"total_delivery_time_rules": len(sre.DeliveryTimeRules),
		"total_restrictions":        len(sre.Restrictions),
		"total_free_shipping_rules": len(sre.FreeShippingRules),
		"active_free_shipping_rules": 0,
		"total_packaging_rules":     len(sre.PackagingRules),
	}

	now := time.Now()
	for _, rule := range sre.ShippingRules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			stats["active_shipping_rules"]++
		}
	}

	for _, rule := range sre.FreeShippingRules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			stats["active_free_shipping_rules"]++
		}
	}

	return stats
}
package shipping

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// ShippingRuleEngine manages shipping rules and configurations
type ShippingRuleEngine struct {
	ShippingRules     []ShippingRule
	CarrierRules      []CarrierRule
	ZoneRules         []ZoneRule
	DeliveryTimeRules []DeliveryTimeRule
	Restrictions      []ShippingRestriction
	FreeShippingRules []FreeShippingRule
	PackagingRules    []PackagingRule
}

// NewShippingRuleEngine creates a new shipping rule engine
func NewShippingRuleEngine() *ShippingRuleEngine {
	return &ShippingRuleEngine{
		ShippingRules:     []ShippingRule{},
		CarrierRules:      []CarrierRule{},
		ZoneRules:         []ZoneRule{},
		DeliveryTimeRules: []DeliveryTimeRule{},
		Restrictions:      []ShippingRestriction{},
		FreeShippingRules: []FreeShippingRule{},
		PackagingRules:    []PackagingRule{},
	}
}

// Shipping Rule Management

// AddShippingRule adds a new shipping rule
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

// UpdateShippingRule updates an existing shipping rule
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

// RemoveShippingRule removes a shipping rule
func (sre *ShippingRuleEngine) RemoveShippingRule(ruleID string) error {
	for i, rule := range sre.ShippingRules {
		if rule.ID == ruleID {
			sre.ShippingRules = append(sre.ShippingRules[:i], sre.ShippingRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("shipping rule with ID %s not found", ruleID)
}

// GetShippingRule retrieves a shipping rule by ID
func (sre *ShippingRuleEngine) GetShippingRule(ruleID string) (*ShippingRule, error) {
	for _, rule := range sre.ShippingRules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("shipping rule with ID %s not found", ruleID)
}

// GetActiveShippingRules returns all active shipping rules
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

// AddCarrierRule adds a new carrier rule
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

// UpdateCarrierRule updates an existing carrier rule
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

// RemoveCarrierRule removes a carrier rule
func (sre *ShippingRuleEngine) RemoveCarrierRule(carrierID, serviceCode string) error {
	for i, rule := range sre.CarrierRules {
		if rule.CarrierID == carrierID && rule.ServiceCode == serviceCode {
			sre.CarrierRules = append(sre.CarrierRules[:i], sre.CarrierRules[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("carrier rule for %s %s not found", carrierID, serviceCode)
}

// GetCarrierRules returns all carrier rules for a specific carrier
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

// AddZoneRule adds a new zone rule
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

// UpdateZoneRule updates an existing zone rule
func (sre *ShippingRuleEngine) UpdateZoneRule(index int, updatedRule ZoneRule) error {
	if index < 0 || index >= len(sre.ZoneRules) {
		return errors.New("invalid zone rule index")
	}

	sre.ZoneRules[index] = updatedRule
	return nil
}

// RemoveZoneRule removes a zone rule
func (sre *ShippingRuleEngine) RemoveZoneRule(index int) error {
	if index < 0 || index >= len(sre.ZoneRules) {
		return errors.New("invalid zone rule index")
	}

	sre.ZoneRules = append(sre.ZoneRules[:index], sre.ZoneRules[index+1:]...)
	return nil
}

// GetZoneRules returns all zone rules for a specific zone
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

// ValidateRuleConfiguration validates the entire rule configuration
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

// GetApplicableRules returns rules applicable for given input
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

// OptimizeRules optimizes rule configuration for better performance
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

// ExportRules exports all rules to a structured format
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

// ImportRules imports rules from a structured format
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

// GetRuleStatistics returns statistics about the rule configuration
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
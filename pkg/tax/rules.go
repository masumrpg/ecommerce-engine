// Package tax provides comprehensive tax rule management and validation capabilities.
// It includes a rule engine for managing tax rules, validation rules, and audit trails.
//
// Key Features:
//   - Tax rule management (CRUD operations)
//   - Rule validation and conflict detection
//   - Audit trail for compliance tracking
//   - Rule optimization and statistics
//   - Import/export functionality
//   - Default rule creation
//
// Basic Usage:
//
//	config := TaxConfiguration{
//		DefaultRate: 8.25,
//		RoundingMode: "round_half_up",
//	}
//	engine := NewTaxRuleEngine(config)
//
//	rule := TaxRule{
//		ID:           "sales_tax_ca",
//		Name:         "California Sales Tax",
//		Type:         TaxTypeSales,
//		Jurisdiction: JurisdictionState,
//		Rate:         8.25,
//		IsActive:     true,
//	}
//
//	err := engine.AddRule(rule)
//	if err != nil {
//		log.Fatal(err)
//	}
package tax

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// TaxRuleEngine manages tax rules and configurations for an e-commerce system.
// It provides comprehensive functionality for managing tax rules, validation rules,
// and maintaining an audit trail for compliance purposes.
//
// The engine supports:
//   - CRUD operations for tax rules
//   - Rule validation and conflict detection
//   - Audit trail management
//   - Rule optimization and statistics
//   - Import/export functionality
//
// Example:
//
//	config := TaxConfiguration{
//		DefaultRate: 8.25,
//		RoundingMode: "round_half_up",
//	}
//	engine := NewTaxRuleEngine(config)
//
//	// Add a new tax rule
//	rule := TaxRule{
//		ID:           "sales_tax_ca",
//		Name:         "California Sales Tax",
//		Type:         TaxTypeSales,
//		Jurisdiction: JurisdictionState,
//		Rate:         8.25,
//		IsActive:     true,
//	}
//	err := engine.AddRule(rule)
type TaxRuleEngine struct {
	// Rules contains all tax rules managed by this engine
	Rules []TaxRule
	// ValidationRules contains custom validation rules for tax rule validation
	ValidationRules []TaxValidationRule
	// Configuration contains the engine's configuration settings
	Configuration TaxConfiguration
	// AuditTrail contains a log of all operations performed on the engine
	AuditTrail []TaxAuditTrail
}

// NewTaxRuleEngine creates a new tax rule engine with the specified configuration.
// The engine is initialized with empty rule sets and audit trail.
//
// Parameters:
//   - config: TaxConfiguration containing engine settings like default rates and rounding modes
//
// Returns:
//   - *TaxRuleEngine: A new tax rule engine instance ready for use
//
// Example:
//
//	config := TaxConfiguration{
//		DefaultRate:  8.25,
//		RoundingMode: "round_half_up",
//		Precision:    2,
//	}
//	engine := NewTaxRuleEngine(config)
func NewTaxRuleEngine(config TaxConfiguration) *TaxRuleEngine {
	return &TaxRuleEngine{
		Rules:           []TaxRule{},
		ValidationRules: []TaxValidationRule{},
		Configuration:   config,
		AuditTrail:      []TaxAuditTrail{},
	}
}

// AddRule adds a new tax rule to the engine after validation and conflict checking.
// The rule is validated against the engine's validation rules and checked for conflicts
// with existing rules before being added.
//
// Parameters:
//   - rule: TaxRule to be added to the engine
//
// Returns:
//   - error: nil if successful, otherwise an error describing the validation failure or conflict
//
// The method performs the following checks:
//   - Rule validation (required fields, valid rates, date ranges)
//   - Conflict detection (duplicate IDs, overlapping jurisdiction/type/time)
//   - Custom validation rules execution
//
// Example:
//
//	rule := TaxRule{
//		ID:           "vat_uk",
//		Name:         "UK VAT",
//		Type:         TaxTypeVAT,
//		Jurisdiction: JurisdictionCountry,
//		Rate:         20.0,
//		IsActive:     true,
//		ValidFrom:    time.Now(),
//		ValidUntil:   time.Now().AddDate(1, 0, 0),
//	}
//	err := engine.AddRule(rule)
func (tre *TaxRuleEngine) AddRule(rule TaxRule) error {
	if err := tre.validateRule(rule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	// Check for conflicts
	if conflicts := tre.findRuleConflicts(rule); len(conflicts) > 0 {
		return fmt.Errorf("rule conflicts with existing rules: %v", conflicts)
	}

	tre.Rules = append(tre.Rules, rule)
	tre.logAuditTrail("ADD_RULE", fmt.Sprintf("Added rule: %s", rule.Name), rule.ID)
	return nil
}

// UpdateRule updates an existing tax rule with new values after validation and conflict checking.
// The updated rule is validated and checked for conflicts with other rules (excluding itself).
//
// Parameters:
//   - ruleID: string ID of the rule to update
//   - updatedRule: TaxRule containing the new rule data
//
// Returns:
//   - error: nil if successful, otherwise an error describing the failure
//
// Example:
//
//	updatedRule := TaxRule{
//		ID:           "vat_uk",
//		Name:         "UK VAT Updated",
//		Type:         TaxTypeVAT,
//		Rate:         21.0, // Updated rate
//		IsActive:     true,
//	}
//	err := engine.UpdateRule("vat_uk", updatedRule)
func (tre *TaxRuleEngine) UpdateRule(ruleID string, updatedRule TaxRule) error {
	for i, rule := range tre.Rules {
		if rule.ID == ruleID {
			if err := tre.validateRule(updatedRule); err != nil {
				return fmt.Errorf("invalid updated rule: %w", err)
			}

			// Check for conflicts (excluding current rule)
			tempRules := make([]TaxRule, 0, len(tre.Rules)-1)
			for j, r := range tre.Rules {
				if j != i {
					tempRules = append(tempRules, r)
				}
			}
			oldEngine := &TaxRuleEngine{Rules: tempRules}
			if conflicts := oldEngine.findRuleConflicts(updatedRule); len(conflicts) > 0 {
				return fmt.Errorf("updated rule conflicts with existing rules: %v", conflicts)
			}

			oldRule := tre.Rules[i]
			tre.Rules[i] = updatedRule
			tre.logAuditTrail("UPDATE_RULE", fmt.Sprintf("Updated rule: %s -> %s", oldRule.Name, updatedRule.Name), ruleID)
			return nil
		}
	}
	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// RemoveRule removes a tax rule from the engine by its ID.
// The operation is logged in the audit trail for compliance tracking.
//
// Parameters:
//   - ruleID: string ID of the rule to remove
//
// Returns:
//   - error: nil if successful, otherwise an error if the rule is not found
//
// Example:
//
//	err := engine.RemoveRule("vat_uk")
//	if err != nil {
//		log.Printf("Failed to remove rule: %v", err)
//	}
func (tre *TaxRuleEngine) RemoveRule(ruleID string) error {
	for i, rule := range tre.Rules {
		if rule.ID == ruleID {
			tre.Rules = append(tre.Rules[:i], tre.Rules[i+1:]...)
			tre.logAuditTrail("REMOVE_RULE", fmt.Sprintf("Removed rule: %s", rule.Name), ruleID)
			return nil
		}
	}
	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetRule retrieves a tax rule by its ID.
//
// Parameters:
//   - ruleID: string ID of the rule to retrieve
//
// Returns:
//   - *TaxRule: pointer to the found rule, nil if not found
//   - error: nil if successful, otherwise an error if the rule is not found
//
// Example:
//
//	rule, err := engine.GetRule("vat_uk")
//	if err != nil {
//		log.Printf("Rule not found: %v", err)
//		return
//	}
//	fmt.Printf("Rule rate: %.2f%%", rule.Rate)
func (tre *TaxRuleEngine) GetRule(ruleID string) (*TaxRule, error) {
	for _, rule := range tre.Rules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetRulesByJurisdiction returns all rules that apply to a specific jurisdiction.
//
// Parameters:
//   - jurisdiction: TaxJurisdiction to filter rules by (e.g., federal, state, county, city)
//
// Returns:
//   - []TaxRule: slice of rules matching the specified jurisdiction
//
// Example:
//
//	stateRules := engine.GetRulesByJurisdiction(JurisdictionState)
//	fmt.Printf("Found %d state-level tax rules", len(stateRules))
func (tre *TaxRuleEngine) GetRulesByJurisdiction(jurisdiction TaxJurisdiction) []TaxRule {
	rules := []TaxRule{}
	for _, rule := range tre.Rules {
		if rule.Jurisdiction == jurisdiction {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetRulesByType returns all rules that apply to a specific tax type.
//
// Parameters:
//   - taxType: TaxType to filter rules by (e.g., sales, VAT, luxury, excise)
//
// Returns:
//   - []TaxRule: slice of rules matching the specified tax type
//
// Example:
//
//	salesRules := engine.GetRulesByType(TaxTypeSales)
//	for _, rule := range salesRules {
//		fmt.Printf("Sales tax rule: %s (%.2f%%)", rule.Name, rule.Rate)
//	}
func (tre *TaxRuleEngine) GetRulesByType(taxType TaxType) []TaxRule {
	rules := []TaxRule{}
	for _, rule := range tre.Rules {
		if rule.Type == taxType {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetActiveRules returns all rules that are currently active and within their valid date range.
// A rule is considered active if:
//   - IsActive flag is true
//   - Current time is after ValidFrom date
//   - Current time is before ValidUntil date
//
// Returns:
//   - []TaxRule: slice of currently active rules
//
// Example:
//
//	activeRules := engine.GetActiveRules()
//	fmt.Printf("Currently %d active tax rules", len(activeRules))
func (tre *TaxRuleEngine) GetActiveRules() []TaxRule {
	rules := []TaxRule{}
	now := time.Now()
	for _, rule := range tre.Rules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetRulesByName returns all rules sorted alphabetically by name.
// This method creates a copy of the rules slice to avoid modifying the original order.
//
// Returns:
//   - []TaxRule: slice of all rules sorted by name in ascending order
//
// Example:
//
//	sortedRules := engine.GetRulesByName()
//	for _, rule := range sortedRules {
//		fmt.Printf("Rule: %s", rule.Name)
//	}
func (tre *TaxRuleEngine) GetRulesByName() []TaxRule {
	rules := make([]TaxRule, len(tre.Rules))
	copy(rules, tre.Rules)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Name < rules[j].Name
	})
	return rules
}

// ActivateRule activates a tax rule by setting its IsActive flag to true.
// The operation is logged in the audit trail.
//
// Parameters:
//   - ruleID: string ID of the rule to activate
//
// Returns:
//   - error: nil if successful, otherwise an error if the rule is not found
//
// Example:
//
//	err := engine.ActivateRule("seasonal_tax")
//	if err != nil {
//		log.Printf("Failed to activate rule: %v", err)
//	}
func (tre *TaxRuleEngine) ActivateRule(ruleID string) error {
	for i, rule := range tre.Rules {
		if rule.ID == ruleID {
			tre.Rules[i].IsActive = true
			tre.logAuditTrail("ACTIVATE_RULE", fmt.Sprintf("Activated rule: %s", rule.Name), ruleID)
			return nil
		}
	}
	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// DeactivateRule deactivates a tax rule by setting its IsActive flag to false.
// The operation is logged in the audit trail.
//
// Parameters:
//   - ruleID: string ID of the rule to deactivate
//
// Returns:
//   - error: nil if successful, otherwise an error if the rule is not found
//
// Example:
//
//	err := engine.DeactivateRule("temporary_tax")
//	if err != nil {
//		log.Printf("Failed to deactivate rule: %v", err)
//	}
func (tre *TaxRuleEngine) DeactivateRule(ruleID string) error {
	for i, rule := range tre.Rules {
		if rule.ID == ruleID {
			tre.Rules[i].IsActive = false
			tre.logAuditTrail("DEACTIVATE_RULE", fmt.Sprintf("Deactivated rule: %s", rule.Name), ruleID)
			return nil
		}
	}
	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// AddValidationRule adds a custom validation rule to the engine.
// Validation rules are used to enforce business constraints when adding or updating tax rules.
//
// Parameters:
//   - rule: TaxValidationRule to be added to the engine
//
// Example:
//
//	validationRule := TaxValidationRule{
//		ID:        "max_rate_check",
//		Name:      "Maximum Rate Validation",
//		Type:      "rate_limit",
//		Condition: "rate <= 25.0",
//		Message:   "Tax rate cannot exceed 25%",
//		Severity:  "error",
//		IsActive:  true,
//	}
//	engine.AddValidationRule(validationRule)
func (tre *TaxRuleEngine) AddValidationRule(rule TaxValidationRule) {
	tre.ValidationRules = append(tre.ValidationRules, rule)
	tre.logAuditTrail("ADD_VALIDATION_RULE", fmt.Sprintf("Added validation rule: %s", rule.Name), rule.ID)
}

// RemoveValidationRule removes a validation rule from the engine by its ID.
//
// Parameters:
//   - ruleID: string ID of the validation rule to remove
//
// Returns:
//   - error: nil if successful, otherwise an error if the validation rule is not found
//
// Example:
//
//	err := engine.RemoveValidationRule("max_rate_check")
//	if err != nil {
//		log.Printf("Failed to remove validation rule: %v", err)
//	}
func (tre *TaxRuleEngine) RemoveValidationRule(ruleID string) error {
	for i, rule := range tre.ValidationRules {
		if rule.ID == ruleID {
			tre.ValidationRules = append(tre.ValidationRules[:i], tre.ValidationRules[i+1:]...)
			tre.logAuditTrail("REMOVE_VALIDATION_RULE", fmt.Sprintf("Removed validation rule: %s", rule.Name), ruleID)
			return nil
		}
	}
	return fmt.Errorf("validation rule with ID %s not found", ruleID)
}

// ValidateRules validates all rules in the engine against the current validation rules.
// This method is useful for batch validation after importing rules or configuration changes.
//
// Returns:
//   - []error: slice of validation errors, empty if all rules are valid
//
// Example:
//
//	errors := engine.ValidateRules()
//	if len(errors) > 0 {
//		for _, err := range errors {
//			log.Printf("Validation error: %v", err)
//		}
//	} else {
//		log.Println("All rules are valid")
//	}
func (tre *TaxRuleEngine) ValidateRules() []error {
	errors := []error{}
	for _, rule := range tre.Rules {
		if err := tre.validateRule(rule); err != nil {
			errors = append(errors, fmt.Errorf("rule %s: %w", rule.ID, err))
		}
	}
	return errors
}

// OptimizeRules optimizes the order of rules for better performance.
// Currently sorts rules alphabetically by name for consistent ordering.
// Future implementations may include more sophisticated optimization strategies.
//
// Example:
//
//	engine.OptimizeRules()
//	log.Println("Rules have been optimized for performance")
func (tre *TaxRuleEngine) OptimizeRules() {
	// Sort by name for consistent ordering
	sort.Slice(tre.Rules, func(i, j int) bool {
		return tre.Rules[i].Name < tre.Rules[j].Name
	})
	tre.logAuditTrail("OPTIMIZE_RULES", "Optimized rule order", "")
}

// ExportRules exports all rules and configuration to a map suitable for serialization.
// The exported data includes rules, validation rules, configuration, and metadata.
//
// Returns:
//   - map[string]interface{}: map containing all engine data with export metadata
//
// The returned map contains:
//   - "rules": all tax rules
//   - "validation_rules": all validation rules
//   - "configuration": engine configuration
//   - "export_date": timestamp of export
//   - "version": export format version
//
// Example:
//
//	exportData := engine.ExportRules()
//	jsonData, err := json.Marshal(exportData)
//	if err == nil {
//		ioutil.WriteFile("tax_rules_backup.json", jsonData, 0644)
//	}
func (tre *TaxRuleEngine) ExportRules() map[string]interface{} {
	return map[string]interface{}{
		"rules":            tre.Rules,
		"validation_rules": tre.ValidationRules,
		"configuration":    tre.Configuration,
		"export_date":      time.Now(),
		"version":          "1.0",
	}
}

// ImportRules imports rules and configuration from a map (typically from ExportRules).
// This method replaces the current rules and configuration with the imported data.
//
// Parameters:
//   - data: map[string]interface{} containing exported rule data
//
// Returns:
//   - error: nil if successful, otherwise an error describing the import failure
//
// Example:
//
//	var importData map[string]interface{}
//	jsonData, _ := ioutil.ReadFile("tax_rules_backup.json")
//	json.Unmarshal(jsonData, &importData)
//	err := engine.ImportRules(importData)
//	if err != nil {
//		log.Printf("Import failed: %v", err)
//	}
func (tre *TaxRuleEngine) ImportRules(data map[string]interface{}) error {
	if rules, ok := data["rules"].([]TaxRule); ok {
		tre.Rules = rules
	}
	if validationRules, ok := data["validation_rules"].([]TaxValidationRule); ok {
		tre.ValidationRules = validationRules
	}
	if config, ok := data["configuration"].(TaxConfiguration); ok {
		tre.Configuration = config
	}
	tre.logAuditTrail("IMPORT_RULES", "Imported rules from external source", "")
	return nil
}

// GetStatistics returns comprehensive statistics about the rule engine.
// Provides insights into rule distribution, activity status, and engine usage.
//
// Returns:
//   - map[string]interface{}: map containing various statistics
//
// The returned statistics include:
//   - "total_rules": total number of rules
//   - "active_rules": number of currently active rules
//   - "inactive_rules": number of inactive rules
//   - "validation_rules": number of validation rules
//   - "audit_entries": number of audit trail entries
//   - "jurisdictions": breakdown by jurisdiction
//   - "tax_types": breakdown by tax type
//   - "methods": breakdown by calculation method
//
// Example:
//
//	stats := engine.GetStatistics()
//	fmt.Printf("Total rules: %d, Active: %d", 
//		stats["total_rules"], stats["active_rules"])
func (tre *TaxRuleEngine) GetStatistics() map[string]interface{} {
	stats := map[string]interface{}{
		"total_rules":      len(tre.Rules),
		"active_rules":     0,
		"inactive_rules":   0,
		"validation_rules": len(tre.ValidationRules),
		"audit_entries":    len(tre.AuditTrail),
		"jurisdictions":    make(map[TaxJurisdiction]int),
		"tax_types":        make(map[TaxType]int),
		"methods":          make(map[TaxCalculationMethod]int),
	}

	jurisdictions := make(map[TaxJurisdiction]int)
	taxTypes := make(map[TaxType]int)
	methods := make(map[TaxCalculationMethod]int)
	activeCount := 0
	inactiveCount := 0
	now := time.Now()

	for _, rule := range tre.Rules {
		if rule.IsActive && now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) {
			activeCount++
		} else {
			inactiveCount++
		}

		jurisdictions[rule.Jurisdiction]++
		taxTypes[rule.Type]++
		methods[rule.Method]++
	}

	stats["active_rules"] = activeCount
	stats["inactive_rules"] = inactiveCount
	stats["jurisdictions"] = jurisdictions
	stats["tax_types"] = taxTypes
	stats["methods"] = methods

	return stats
}

// GetAuditTrail returns the complete audit trail of all operations performed on the engine.
// The audit trail provides a chronological record of all rule modifications for compliance.
//
// Returns:
//   - []TaxAuditTrail: slice of audit entries in chronological order
//
// Example:
//
//	auditTrail := engine.GetAuditTrail()
//	for _, entry := range auditTrail {
//		fmt.Printf("%s: %s at %s", entry.Action, entry.Reason, entry.Timestamp)
//	}
func (tre *TaxRuleEngine) GetAuditTrail() []TaxAuditTrail {
	return tre.AuditTrail
}

// ClearAuditTrail clears the audit trail by removing all audit entries.
// Use with caution as this removes compliance tracking history.
//
// Example:
//
//	engine.ClearAuditTrail()
//	log.Println("Audit trail has been cleared")
func (tre *TaxRuleEngine) ClearAuditTrail() {
	tre.AuditTrail = []TaxAuditTrail{}
}

// validateRule validates a single tax rule against basic constraints and custom validation rules.
// This is an internal method used by AddRule and UpdateRule.
//
// Parameters:
//   - rule: TaxRule to validate
//
// Returns:
//   - error: nil if valid, otherwise an error describing the validation failure
//
// Validation checks include:
//   - Required fields (ID, Name)
//   - Rate constraints (non-negative, percentage limits)
//   - Date range validity
//   - Amount constraints
//   - Threshold validation
//   - Custom validation rules
func (tre *TaxRuleEngine) validateRule(rule TaxRule) error {
	if rule.ID == "" {
		return errors.New("rule ID is required")
	}
	if rule.Name == "" {
		return errors.New("rule name is required")
	}
	if rule.Rate < 0 {
		return errors.New("tax rate cannot be negative")
	}
	if rule.Rate > 100 && rule.Method == TaxMethodPercentage {
		return errors.New("percentage tax rate cannot exceed 100%")
	}
	if rule.ValidFrom.After(rule.ValidUntil) {
		return errors.New("valid from date must be before valid until date")
	}
	if rule.MinAmount < 0 {
		return errors.New("minimum amount cannot be negative")
	}
	if rule.MaxAmount > 0 && rule.MaxAmount < rule.MinAmount {
		return errors.New("maximum amount must be greater than minimum amount")
	}

	// Validate thresholds
	for i, threshold := range rule.Thresholds {
		if threshold.MinAmount < 0 {
			return fmt.Errorf("threshold %d: minimum amount cannot be negative", i)
		}
		if threshold.MaxAmount > 0 && threshold.MaxAmount < threshold.MinAmount {
			return fmt.Errorf("threshold %d: maximum amount must be greater than minimum amount", i)
		}
		if threshold.Rate < 0 {
			return fmt.Errorf("threshold %d: rate cannot be negative", i)
		}
		if threshold.FixedAmount < 0 {
			return fmt.Errorf("threshold %d: fixed amount cannot be negative", i)
		}
	}

	// Run custom validation rules
	for _, validationRule := range tre.ValidationRules {
		if err := tre.runValidationRule(validationRule, rule); err != nil {
			return fmt.Errorf("validation rule %s failed: %w", validationRule.Name, err)
		}
	}

	return nil
}

// findRuleConflicts identifies conflicts between a new rule and existing rules.
// This is an internal method used to prevent rule conflicts during addition/updates.
//
// Parameters:
//   - newRule: TaxRule to check for conflicts
//
// Returns:
//   - []string: slice of conflict descriptions, empty if no conflicts
//
// Conflict detection includes:
//   - Duplicate rule IDs
//   - Overlapping jurisdiction, type, time period, and geographic coverage
func (tre *TaxRuleEngine) findRuleConflicts(newRule TaxRule) []string {
	conflicts := []string{}

	for _, existingRule := range tre.Rules {
		// Check for duplicate IDs
		if existingRule.ID == newRule.ID {
			conflicts = append(conflicts, fmt.Sprintf("duplicate rule ID: %s", newRule.ID))
		}

		// Check for overlapping jurisdiction and type
		if existingRule.Jurisdiction == newRule.Jurisdiction &&
			existingRule.Type == newRule.Type &&
			tre.hasTimeOverlap(existingRule, newRule) &&
			tre.hasGeographicOverlap(existingRule, newRule) {
			conflicts = append(conflicts, fmt.Sprintf("overlapping rule: %s", existingRule.ID))
		}
	}

	return conflicts
}

// hasTimeOverlap checks if two rules have overlapping valid time periods.
// This is a helper method for conflict detection.
//
// Parameters:
//   - rule1, rule2: TaxRule instances to compare
//
// Returns:
//   - bool: true if the rules have overlapping time periods
func (tre *TaxRuleEngine) hasTimeOverlap(rule1, rule2 TaxRule) bool {
	return rule1.ValidFrom.Before(rule2.ValidUntil) && rule2.ValidFrom.Before(rule1.ValidUntil)
}

// hasGeographicOverlap checks if two rules have overlapping geographic coverage.
// This is a helper method for conflict detection.
//
// Parameters:
//   - rule1, rule2: TaxRule instances to compare
//
// Returns:
//   - bool: true if the rules have overlapping geographic coverage
//
// Rules overlap geographically if:
//   - Either rule has no country restrictions (applies globally)
//   - Both rules share at least one common country
func (tre *TaxRuleEngine) hasGeographicOverlap(rule1, rule2 TaxRule) bool {
	// If either rule has no geographic restrictions, they overlap
	if len(rule1.ApplicableCountries) == 0 || len(rule2.ApplicableCountries) == 0 {
		return true
	}

	// Check for common countries
	for _, country1 := range rule1.ApplicableCountries {
		for _, country2 := range rule2.ApplicableCountries {
			if country1 == country2 {
				return true
			}
		}
	}

	return false
}

// runValidationRule executes a custom validation rule against a tax rule.
// This is an internal method that implements the validation logic for different rule types.
//
// Parameters:
//   - validationRule: TaxValidationRule to execute
//   - rule: TaxRule to validate against
//
// Returns:
//   - error: nil if validation passes, otherwise an error describing the failure
//
// Supported validation types:
//   - "rate_limit": validates maximum tax rate
//   - "jurisdiction_limit": validates allowed jurisdictions
//   - "date_range": validates rule duration limits
func (tre *TaxRuleEngine) runValidationRule(validationRule TaxValidationRule, rule TaxRule) error {
	// This is a simplified implementation
	// In a real system, you might have a more sophisticated rule engine
	switch validationRule.Type {
	case "rate_limit":
		// Simple rate validation - max 50%
		if rule.Rate > 50.0 {
			return fmt.Errorf("rate %.2f exceeds maximum allowed rate 50.0", rule.Rate)
		}
	case "jurisdiction_limit":
		// Validate allowed jurisdictions
		allowedJurisdictions := []TaxJurisdiction{
			JurisdictionFederal,
			JurisdictionState,
			JurisdictionCounty,
			JurisdictionCity,
		}
		allowed := false
		for _, jurisdiction := range allowedJurisdictions {
			if rule.Jurisdiction == jurisdiction {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("jurisdiction %s is not allowed", rule.Jurisdiction)
		}
	case "date_range":
		// Validate rule duration - max 365 days
		duration := rule.ValidUntil.Sub(rule.ValidFrom)
		if int(duration.Hours()/24) > 365 {
			return fmt.Errorf("rule duration exceeds maximum allowed days: 365")
		}
	}

	return nil
}

// logAuditTrail records an operation in the audit trail for compliance tracking.
// This is an internal method used by all rule modification operations.
//
// Parameters:
//   - action: string describing the action performed (e.g., "ADD_RULE", "UPDATE_RULE")
//   - reason: string providing additional context about the action
//   - ruleID: string ID of the rule affected by the action
func (tre *TaxRuleEngine) logAuditTrail(action, reason, ruleID string) {
	auditEntry := TaxAuditTrail{
		ID:            fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		Action:        action,
		Reason:        reason,
		TransactionID: ruleID,
		Timestamp:     time.Now(),
		UserID:        "system", // In a real system, this would be the actual user ID
	}
	tre.AuditTrail = append(tre.AuditTrail, auditEntry)
}

// CreateDefaultRules creates a set of default tax rules for common scenarios.
// This function provides pre-configured rules that can be used as a starting point
// for new tax rule engines.
//
// Returns:
//   - []TaxRule: slice of default tax rules including sales tax and luxury tax
//
// The default rules include:
//   - Standard sales tax (8.25% for US states)
//   - Luxury tax (15% for high-value items)
//
// Example:
//
//	defaultRules := CreateDefaultRules()
//	for _, rule := range defaultRules {
//		engine.AddRule(rule)
//	}
func CreateDefaultRules() []TaxRule {
	return []TaxRule{
		{
			ID:           "default_sales_tax",
			Name:         "Default Sales Tax",
			Description:  "Standard sales tax rate",
			Type:         TaxTypeSales,
			Jurisdiction: JurisdictionState,
			Method:       TaxMethodPercentage,
			Rate:         8.25,
			IsActive:     true,
			ValidFrom:    time.Now(),
			ValidUntil:   time.Now().AddDate(1, 0, 0),
			ApplicableCountries: []string{"US"},
			ApplicableStates:    []string{"CA", "NY", "TX"},
			MinAmount:           0,
			MaxAmount:           0,
			Thresholds:          []TaxThreshold{},
			Conditions:          []TaxCondition{},
			ApplicableCategories: []string{},
			ExemptCategories:     []string{"food", "medicine"},
		},
		{
			ID:           "luxury_tax",
			Name:         "Luxury Tax",
			Description:  "Tax on luxury items",
			Type:         TaxTypeLuxury,
			Jurisdiction: JurisdictionFederal,
			Method:       TaxMethodPercentage,
			Rate:         15.0,
			IsActive:     true,
			ValidFrom:    time.Now(),
			ValidUntil:   time.Now().AddDate(1, 0, 0),
			ApplicableCountries: []string{"US"},
			MinAmount:           1000,
			MaxAmount:           0,
			Thresholds:          []TaxThreshold{},
			Conditions: []TaxCondition{
				{
					Type:     "amount",
					Operator: ">=",
					Value:    1000.0,
				},
			},
			ApplicableCategories: []string{"jewelry", "luxury_cars", "yachts"},
			ExemptCategories:     []string{},
		},
	}
}

// CreateDefaultValidationRules creates a set of default validation rules.
// These rules enforce common business constraints and can be used as a foundation
// for custom validation requirements.
//
// Returns:
//   - []TaxValidationRule: slice of default validation rules
//
// The default validation rules include:
//   - Maximum rate validation (50% limit)
//   - Jurisdiction validation (allowed jurisdiction types)
//
// Example:
//
//	defaultValidationRules := CreateDefaultValidationRules()
//	for _, rule := range defaultValidationRules {
//		engine.AddValidationRule(rule)
//	}
func CreateDefaultValidationRules() []TaxValidationRule {
	return []TaxValidationRule{
		{
			ID:        "max_rate_validation",
			Name:      "Maximum Rate Validation",
			Type:      "rate_limit",
			Condition: "rate <= 50.0",
			Message:   "Tax rate cannot exceed 50%",
			Severity:  "error",
			IsActive:  true,
			ValidFrom: time.Now(),
			ValidUntil: time.Now().AddDate(1, 0, 0),
		},
		{
			ID:        "jurisdiction_validation",
			Name:      "Jurisdiction Validation",
			Type:      "jurisdiction_limit",
			Condition: "jurisdiction in [federal, state, county, city]",
			Message:   "Invalid jurisdiction specified",
			Severity:  "error",
			IsActive:  true,
			ValidFrom: time.Now(),
			ValidUntil: time.Now().AddDate(1, 0, 0),
		},
	}
}
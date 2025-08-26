package tax

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// TaxRuleEngine manages tax rules and configurations
type TaxRuleEngine struct {
	Rules           []TaxRule
	ValidationRules []TaxValidationRule
	Configuration   TaxConfiguration
	AuditTrail      []TaxAuditTrail
}

// NewTaxRuleEngine creates a new tax rule engine
func NewTaxRuleEngine(config TaxConfiguration) *TaxRuleEngine {
	return &TaxRuleEngine{
		Rules:           []TaxRule{},
		ValidationRules: []TaxValidationRule{},
		Configuration:   config,
		AuditTrail:      []TaxAuditTrail{},
	}
}

// AddRule adds a new tax rule
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

// UpdateRule updates an existing tax rule
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

// RemoveRule removes a tax rule
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

// GetRule retrieves a tax rule by ID
func (tre *TaxRuleEngine) GetRule(ruleID string) (*TaxRule, error) {
	for _, rule := range tre.Rules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}
	return nil, fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetRulesByJurisdiction returns rules for a specific jurisdiction
func (tre *TaxRuleEngine) GetRulesByJurisdiction(jurisdiction TaxJurisdiction) []TaxRule {
	rules := []TaxRule{}
	for _, rule := range tre.Rules {
		if rule.Jurisdiction == jurisdiction {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetRulesByType returns rules for a specific tax type
func (tre *TaxRuleEngine) GetRulesByType(taxType TaxType) []TaxRule {
	rules := []TaxRule{}
	for _, rule := range tre.Rules {
		if rule.Type == taxType {
			rules = append(rules, rule)
		}
	}
	return rules
}

// GetActiveRules returns all active rules
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

// GetRulesByName returns rules sorted by name
func (tre *TaxRuleEngine) GetRulesByName() []TaxRule {
	rules := make([]TaxRule, len(tre.Rules))
	copy(rules, tre.Rules)
	sort.Slice(rules, func(i, j int) bool {
		return rules[i].Name < rules[j].Name
	})
	return rules
}

// ActivateRule activates a tax rule
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

// DeactivateRule deactivates a tax rule
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

// AddValidationRule adds a validation rule
func (tre *TaxRuleEngine) AddValidationRule(rule TaxValidationRule) {
	tre.ValidationRules = append(tre.ValidationRules, rule)
	tre.logAuditTrail("ADD_VALIDATION_RULE", fmt.Sprintf("Added validation rule: %s", rule.Name), rule.ID)
}

// RemoveValidationRule removes a validation rule
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

// ValidateRules validates all rules in the engine
func (tre *TaxRuleEngine) ValidateRules() []error {
	errors := []error{}
	for _, rule := range tre.Rules {
		if err := tre.validateRule(rule); err != nil {
			errors = append(errors, fmt.Errorf("rule %s: %w", rule.ID, err))
		}
	}
	return errors
}

// OptimizeRules optimizes rule order for better performance
func (tre *TaxRuleEngine) OptimizeRules() {
	// Sort by name for consistent ordering
	sort.Slice(tre.Rules, func(i, j int) bool {
		return tre.Rules[i].Name < tre.Rules[j].Name
	})
	tre.logAuditTrail("OPTIMIZE_RULES", "Optimized rule order", "")
}

// ExportRules exports rules to a map for serialization
func (tre *TaxRuleEngine) ExportRules() map[string]interface{} {
	return map[string]interface{}{
		"rules":            tre.Rules,
		"validation_rules": tre.ValidationRules,
		"configuration":    tre.Configuration,
		"export_date":      time.Now(),
		"version":          "1.0",
	}
}

// ImportRules imports rules from a map
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

// GetStatistics returns statistics about the rule engine
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

// GetAuditTrail returns the audit trail
func (tre *TaxRuleEngine) GetAuditTrail() []TaxAuditTrail {
	return tre.AuditTrail
}

// ClearAuditTrail clears the audit trail
func (tre *TaxRuleEngine) ClearAuditTrail() {
	tre.AuditTrail = []TaxAuditTrail{}
}

// validateRule validates a single tax rule
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

// findRuleConflicts finds conflicts with existing rules
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

// hasTimeOverlap checks if two rules have overlapping time periods
func (tre *TaxRuleEngine) hasTimeOverlap(rule1, rule2 TaxRule) bool {
	return rule1.ValidFrom.Before(rule2.ValidUntil) && rule2.ValidFrom.Before(rule1.ValidUntil)
}

// hasGeographicOverlap checks if two rules have overlapping geographic coverage
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

// runValidationRule runs a custom validation rule
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

// logAuditTrail logs an action to the audit trail
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

// CreateDefaultRules creates a set of default tax rules
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

// CreateDefaultValidationRules creates default validation rules
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
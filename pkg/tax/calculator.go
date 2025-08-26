package tax

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// TaxCalculator handles tax calculations
type TaxCalculator struct {
	Configuration TaxConfiguration
	Rules         []TaxRule
	ValidationRules []TaxValidationRule
}

// NewTaxCalculator creates a new tax calculator
func NewTaxCalculator(config TaxConfiguration) *TaxCalculator {
	return &TaxCalculator{
		Configuration: config,
		Rules:         config.DefaultRules,
		ValidationRules: []TaxValidationRule{},
	}
}

// Calculate calculates taxes for given input
func Calculate(input TaxCalculationInput) TaxCalculationResult {
	config := TaxConfiguration{
		DefaultCurrency:     "USD",
		RoundingMode:        "round",
		RoundingPrecision:   2,
		TaxInclusivePricing: false,
		CompoundTaxes:       false,
		TaxOnShipping:       true,
		TaxOnDiscounts:      true,
	}

	calc := NewTaxCalculator(config)
	if len(input.TaxRules) > 0 {
		calc.Rules = input.TaxRules
	}

	return calc.CalculateTax(input)
}

// CalculateTax calculates taxes for the given input
func (tc *TaxCalculator) CalculateTax(input TaxCalculationInput) TaxCalculationResult {
	result := TaxCalculationResult{
		AppliedTaxes:       []AppliedTax{},
		TaxBreakdown:       []TaxBreakdown{},
		JurisdictionTotals: make(map[TaxJurisdiction]float64),
		TaxTypeTotals:      make(map[TaxType]float64),
		Currency:           input.Currency,
		CalculationDate:    time.Now(),
		IsValid:            true,
		Errors:             []string{},
		Warnings:           []string{},
		Metadata:           make(map[string]interface{}),
	}

	if input.Currency == "" {
		result.Currency = tc.Configuration.DefaultCurrency
	}

	// Validate input
	if validationErrors := tc.validateInput(input); len(validationErrors) > 0 {
		result.IsValid = false
		result.Errors = validationErrors
		return result
	}

	// Calculate subtotal
	result.Subtotal = tc.calculateSubtotal(input.Items)

	// Apply discounts if configured
	if tc.Configuration.TaxOnDiscounts && input.DiscountAmount > 0 {
		result.Subtotal -= input.DiscountAmount
	}

	// Add shipping if configured to be taxed
	if tc.Configuration.TaxOnShipping && input.ShippingAmount > 0 {
		result.Subtotal += input.ShippingAmount
	}

	// Get applicable tax rules
	applicableRules := tc.getApplicableRules(input)

	// Sort rules by priority (higher priority first)
	sort.Slice(applicableRules, func(i, j int) bool {
		return applicableRules[i].Priority > applicableRules[j].Priority
	})

	// Calculate taxes for each item
	for _, item := range input.Items {
		breakdown := tc.calculateItemTax(item, applicableRules, input)
		result.TaxBreakdown = append(result.TaxBreakdown, breakdown)
		result.TotalTax += breakdown.TotalTax
		result.TaxableAmount += breakdown.TaxableAmount
		result.ExemptAmount += breakdown.ExemptAmount

		// Aggregate applied taxes
		for _, appliedTax := range breakdown.AppliedTaxes {
			tc.aggregateAppliedTax(&result, appliedTax)
		}
	}

	// Apply tax overrides
	tc.applyTaxOverrides(&result, input.Overrides)

	// Calculate totals
	result.GrandTotal = result.Subtotal + result.TotalTax

	// Calculate effective rate
	if result.Subtotal > 0 {
		result.EffectiveRate = (result.TotalTax / result.Subtotal) * 100
	}

	// Round amounts based on configuration
	tc.roundAmounts(&result)

	// Validate result
	if warnings := tc.validateResult(result); len(warnings) > 0 {
		result.Warnings = warnings
	}

	return result
}

// calculateSubtotal calculates the subtotal of all items
func (tc *TaxCalculator) calculateSubtotal(items []TaxableItem) float64 {
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.TotalAmount
	}
	return subtotal
}

// getApplicableRules returns tax rules applicable to the input
func (tc *TaxCalculator) getApplicableRules(input TaxCalculationInput) []TaxRule {
	applicableRules := []TaxRule{}
	now := time.Now()

	for _, rule := range tc.Rules {
		// Check if rule is active and within valid time range
		if !rule.IsActive || now.Before(rule.ValidFrom) || now.After(rule.ValidUntil) {
			continue
		}

		// Check geographic applicability
		if !tc.isGeographicallyApplicable(rule, input.BillingAddress, input.ShippingAddress) {
			continue
		}

		// Check conditions
		if !tc.evaluateConditions(rule.Conditions, input) {
			continue
		}

		applicableRules = append(applicableRules, rule)
	}

	return applicableRules
}

// isGeographicallyApplicable checks if a rule applies to the given addresses
func (tc *TaxCalculator) isGeographicallyApplicable(rule TaxRule, billingAddr, shippingAddr Address) bool {
	// Use shipping address for tax calculation (common practice)
	addr := shippingAddr
	if addr.Country == "" {
		addr = billingAddr
	}

	// Check countries
	if len(rule.ApplicableCountries) > 0 {
		found := false
		for _, country := range rule.ApplicableCountries {
			if addr.Country == country {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check states
	if len(rule.ApplicableStates) > 0 {
		found := false
		for _, state := range rule.ApplicableStates {
			if addr.State == state {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check cities
	if len(rule.ApplicableCities) > 0 {
		found := false
		for _, city := range rule.ApplicableCities {
			if addr.City == city {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check postal codes
	if len(rule.PostalCodes) > 0 {
		found := false
		for _, postalCode := range rule.PostalCodes {
			if addr.PostalCode == postalCode {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// evaluateConditions evaluates rule conditions
func (tc *TaxCalculator) evaluateConditions(conditions []TaxCondition, input TaxCalculationInput) bool {
	if len(conditions) == 0 {
		return true
	}

	results := []bool{}
	for _, condition := range conditions {
		result := tc.evaluateCondition(condition, input)
		results = append(results, result)
	}

	// Apply logic (simplified - assumes AND by default)
	for _, result := range results {
		if !result {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a single condition
func (tc *TaxCalculator) evaluateCondition(condition TaxCondition, input TaxCalculationInput) bool {
	switch condition.Type {
	case "amount":
		totalAmount := tc.calculateSubtotal(input.Items)
		return tc.compareValues(totalAmount, condition.Operator, condition.Value)
	case "quantity":
		totalQuantity := 0
		for _, item := range input.Items {
			totalQuantity += item.Quantity
		}
		return tc.compareValues(float64(totalQuantity), condition.Operator, condition.Value)
	case "weight":
		totalWeight := 0.0
		for _, item := range input.Items {
			totalWeight += item.Weight * float64(item.Quantity)
		}
		return tc.compareValues(totalWeight, condition.Operator, condition.Value)
	case "category":
		for _, item := range input.Items {
			if tc.compareValues(item.Category, condition.Operator, condition.Value) {
				return true
			}
		}
		return false
	case "customer_type":
		return tc.compareValues(input.Customer.Type, condition.Operator, condition.Value)
	default:
		return true
	}
}

// compareValues compares two values based on operator
func (tc *TaxCalculator) compareValues(actual interface{}, operator string, expected interface{}) bool {
	switch operator {
	case ">":
		actualFloat, _ := tc.toFloat64(actual)
		expectedFloat, _ := tc.toFloat64(expected)
		return actualFloat > expectedFloat
	case "<":
		actualFloat, _ := tc.toFloat64(actual)
		expectedFloat, _ := tc.toFloat64(expected)
		return actualFloat < expectedFloat
	case ">=":
		actualFloat, _ := tc.toFloat64(actual)
		expectedFloat, _ := tc.toFloat64(expected)
		return actualFloat >= expectedFloat
	case "<=":
		actualFloat, _ := tc.toFloat64(actual)
		expectedFloat, _ := tc.toFloat64(expected)
		return actualFloat <= expectedFloat
	case "=":
		return fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", expected)
	case "!=":
		return fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected)
	case "in":
		if expectedSlice, ok := expected.([]interface{}); ok {
			actualStr := fmt.Sprintf("%v", actual)
			for _, item := range expectedSlice {
				if actualStr == fmt.Sprintf("%v", item) {
					return true
				}
			}
		}
		return false
	case "not_in":
		if expectedSlice, ok := expected.([]interface{}); ok {
			actualStr := fmt.Sprintf("%v", actual)
			for _, item := range expectedSlice {
				if actualStr == fmt.Sprintf("%v", item) {
					return false
				}
			}
			return true
		}
		return true
	default:
		return true
	}
}

// toFloat64 converts interface{} to float64
func (tc *TaxCalculator) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// calculateItemTax calculates tax for a single item
func (tc *TaxCalculator) calculateItemTax(item TaxableItem, rules []TaxRule, input TaxCalculationInput) TaxBreakdown {
	breakdown := TaxBreakdown{
		ItemID:        item.ID,
		ItemName:      item.Name,
		ItemAmount:    item.TotalAmount,
		AppliedTaxes:  []AppliedTax{},
		TotalTax:      0,
		TaxableAmount: item.TotalAmount,
		ExemptAmount:  0,
	}

	// Check if item is exempt
	if item.IsExempt {
		breakdown.ExemptAmount = item.TotalAmount
		breakdown.TaxableAmount = 0
		breakdown.ExemptionReason = item.ExemptionReason
		return breakdown
	}

	// Check customer exemptions
	if tc.isCustomerExempt(input.Customer, item) {
		breakdown.ExemptAmount = item.TotalAmount
		breakdown.TaxableAmount = 0
		breakdown.ExemptionReason = "Customer exemption"
		return breakdown
	}

	// Apply applicable tax rules
	for _, rule := range rules {
		if tc.isRuleApplicableToItem(rule, item) {
			appliedTax := tc.calculateTaxForRule(rule, breakdown.TaxableAmount, item)
			if appliedTax.TaxAmount > 0 {
				breakdown.AppliedTaxes = append(breakdown.AppliedTaxes, appliedTax)
				breakdown.TotalTax += appliedTax.TaxAmount

				// For compound taxes, reduce taxable amount
				if tc.Configuration.CompoundTaxes {
					breakdown.TaxableAmount += appliedTax.TaxAmount
				}
			}
		}
	}

	return breakdown
}

// isCustomerExempt checks if customer is exempt from tax
func (tc *TaxCalculator) isCustomerExempt(customer Customer, item TaxableItem) bool {
	for _, exemption := range customer.Exemptions {
		if tc.isExemptionApplicable(exemption, item) {
			return true
		}
	}
	return false
}

// isExemptionApplicable checks if an exemption applies to an item
func (tc *TaxCalculator) isExemptionApplicable(exemption TaxExemption, item TaxableItem) bool {
	now := time.Now()
	if now.Before(exemption.ValidFrom) || now.After(exemption.ValidUntil) {
		return false
	}

	switch exemption.Type {
	case "item":
		// Check if exemption applies to this item category
		for _, condition := range exemption.Conditions {
			if condition.Type == "category" && condition.Value == item.Category {
				return true
			}
		}
	case "customer":
		// Customer-level exemption applies to all items
		return true
	}

	return false
}

// isRuleApplicableToItem checks if a tax rule applies to a specific item
func (tc *TaxCalculator) isRuleApplicableToItem(rule TaxRule, item TaxableItem) bool {
	// Check applicable categories
	if len(rule.ApplicableCategories) > 0 {
		found := false
		for _, category := range rule.ApplicableCategories {
			if item.Category == category || item.Subcategory == category {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check exempt categories
	for _, category := range rule.ExemptCategories {
		if item.Category == category || item.Subcategory == category {
			return false
		}
	}

	// Check amount thresholds
	if rule.MinAmount > 0 && item.TotalAmount < rule.MinAmount {
		return false
	}
	if rule.MaxAmount > 0 && item.TotalAmount > rule.MaxAmount {
		return false
	}

	return true
}

// calculateTaxForRule calculates tax amount for a specific rule
func (tc *TaxCalculator) calculateTaxForRule(rule TaxRule, taxableAmount float64, item TaxableItem) AppliedTax {
	appliedTax := AppliedTax{
		RuleID:        rule.ID,
		Name:          rule.Name,
		Type:          rule.Type,
		Jurisdiction:  rule.Jurisdiction,
		Method:        rule.Method,
		Rate:          rule.Rate,
		TaxableAmount: taxableAmount,
		TaxAmount:     0,
		Description:   rule.Description,
	}

	switch rule.Method {
	case TaxMethodPercentage:
		appliedTax.TaxAmount = taxableAmount * (rule.Rate / 100)
	case TaxMethodFixed:
		appliedTax.TaxAmount = rule.Rate * float64(item.Quantity)
	case TaxMethodTiered:
		appliedTax.TaxAmount = tc.calculateTieredTax(rule.Thresholds, taxableAmount)
	case TaxMethodProgressive:
		appliedTax.TaxAmount = tc.calculateProgressiveTax(rule.Thresholds, taxableAmount)
	case TaxMethodCompound:
		// Compound tax is calculated on the amount including previous taxes
		appliedTax.TaxAmount = taxableAmount * (rule.Rate / 100)
	default:
		appliedTax.TaxAmount = taxableAmount * (rule.Rate / 100)
	}

	return appliedTax
}

// calculateTieredTax calculates tax using tiered rates
func (tc *TaxCalculator) calculateTieredTax(thresholds []TaxThreshold, amount float64) float64 {
	if len(thresholds) == 0 {
		return 0
	}

	// Sort thresholds by minimum amount
	sort.Slice(thresholds, func(i, j int) bool {
		return thresholds[i].MinAmount < thresholds[j].MinAmount
	})

	// Find applicable tier
	for _, threshold := range thresholds {
		if amount >= threshold.MinAmount && (threshold.MaxAmount == 0 || amount <= threshold.MaxAmount) {
			if threshold.FixedAmount > 0 {
				return threshold.FixedAmount
			}
			return amount * (threshold.Rate / 100)
		}
	}

	return 0
}

// calculateProgressiveTax calculates tax using progressive rates
func (tc *TaxCalculator) calculateProgressiveTax(thresholds []TaxThreshold, amount float64) float64 {
	if len(thresholds) == 0 {
		return 0
	}

	// Sort thresholds by minimum amount
	sort.Slice(thresholds, func(i, j int) bool {
		return thresholds[i].MinAmount < thresholds[j].MinAmount
	})

	totalTax := 0.0
	remainingAmount := amount

	for _, threshold := range thresholds {
		if remainingAmount <= 0 {
			break
		}

		// Calculate taxable amount for this tier
		tierAmount := 0.0
		if threshold.MaxAmount > 0 {
			tierAmount = math.Min(remainingAmount, threshold.MaxAmount-threshold.MinAmount)
		} else {
			tierAmount = remainingAmount
		}

		if tierAmount > 0 {
			if threshold.FixedAmount > 0 {
				totalTax += threshold.FixedAmount
			} else {
				totalTax += tierAmount * (threshold.Rate / 100)
			}
			remainingAmount -= tierAmount
		}
	}

	return totalTax
}

// aggregateAppliedTax aggregates applied tax into result totals
func (tc *TaxCalculator) aggregateAppliedTax(result *TaxCalculationResult, appliedTax AppliedTax) {
	// Add to jurisdiction totals
	result.JurisdictionTotals[appliedTax.Jurisdiction] += appliedTax.TaxAmount

	// Add to tax type totals
	result.TaxTypeTotals[appliedTax.Type] += appliedTax.TaxAmount

	// Add to applied taxes list (avoid duplicates)
	found := false
	for i, existing := range result.AppliedTaxes {
		if existing.RuleID == appliedTax.RuleID {
			result.AppliedTaxes[i].TaxAmount += appliedTax.TaxAmount
			result.AppliedTaxes[i].TaxableAmount += appliedTax.TaxableAmount
			found = true
			break
		}
	}

	if !found {
		result.AppliedTaxes = append(result.AppliedTaxes, appliedTax)
	}
}

// applyTaxOverrides applies manual tax overrides
func (tc *TaxCalculator) applyTaxOverrides(result *TaxCalculationResult, overrides []TaxOverride) {
	for _, override := range overrides {
		switch override.Type {
		case "rate":
			// Override tax rate for specific tax type
			for i, appliedTax := range result.AppliedTaxes {
				if appliedTax.Type == override.TaxType {
					oldAmount := appliedTax.TaxAmount
					newAmount := appliedTax.TaxableAmount * (override.Value / 100)
					difference := newAmount - oldAmount

					result.AppliedTaxes[i].TaxAmount = newAmount
					result.AppliedTaxes[i].Rate = override.Value
					result.AppliedTaxes[i].IsOverridden = true
					result.AppliedTaxes[i].OverrideReason = override.Reason

					result.TotalTax += difference
					result.TaxTypeTotals[override.TaxType] += difference
				}
			}
		case "amount":
			// Override tax amount for specific tax type
			for i, appliedTax := range result.AppliedTaxes {
				if appliedTax.Type == override.TaxType {
					oldAmount := appliedTax.TaxAmount
					difference := override.Value - oldAmount

					result.AppliedTaxes[i].TaxAmount = override.Value
					result.AppliedTaxes[i].IsOverridden = true
					result.AppliedTaxes[i].OverrideReason = override.Reason

					result.TotalTax += difference
					result.TaxTypeTotals[override.TaxType] += difference
				}
			}
		case "exempt":
			// Exempt from specific tax type
			for i, appliedTax := range result.AppliedTaxes {
				if appliedTax.Type == override.TaxType {
					result.TotalTax -= appliedTax.TaxAmount
					result.TaxTypeTotals[override.TaxType] -= appliedTax.TaxAmount
					result.ExemptAmount += appliedTax.TaxableAmount

					result.AppliedTaxes[i].TaxAmount = 0
					result.AppliedTaxes[i].IsOverridden = true
					result.AppliedTaxes[i].OverrideReason = override.Reason
				}
			}
		}
	}
}

// roundAmounts rounds amounts based on configuration
func (tc *TaxCalculator) roundAmounts(result *TaxCalculationResult) {
	precision := tc.Configuration.RoundingPrecision
	multiplier := math.Pow(10, float64(precision))

	switch tc.Configuration.RoundingMode {
	case "round":
		result.TotalTax = math.Round(result.TotalTax*multiplier) / multiplier
		result.GrandTotal = math.Round(result.GrandTotal*multiplier) / multiplier
		result.Subtotal = math.Round(result.Subtotal*multiplier) / multiplier
	case "floor":
		result.TotalTax = math.Floor(result.TotalTax*multiplier) / multiplier
		result.GrandTotal = math.Floor(result.GrandTotal*multiplier) / multiplier
		result.Subtotal = math.Floor(result.Subtotal*multiplier) / multiplier
	case "ceil":
		result.TotalTax = math.Ceil(result.TotalTax*multiplier) / multiplier
		result.GrandTotal = math.Ceil(result.GrandTotal*multiplier) / multiplier
		result.Subtotal = math.Ceil(result.Subtotal*multiplier) / multiplier
	}

	// Round applied taxes
	for i := range result.AppliedTaxes {
		switch tc.Configuration.RoundingMode {
		case "round":
			result.AppliedTaxes[i].TaxAmount = math.Round(result.AppliedTaxes[i].TaxAmount*multiplier) / multiplier
		case "floor":
			result.AppliedTaxes[i].TaxAmount = math.Floor(result.AppliedTaxes[i].TaxAmount*multiplier) / multiplier
		case "ceil":
			result.AppliedTaxes[i].TaxAmount = math.Ceil(result.AppliedTaxes[i].TaxAmount*multiplier) / multiplier
		}
	}

	// Round tax breakdown
	for i := range result.TaxBreakdown {
		switch tc.Configuration.RoundingMode {
		case "round":
			result.TaxBreakdown[i].TotalTax = math.Round(result.TaxBreakdown[i].TotalTax*multiplier) / multiplier
		case "floor":
			result.TaxBreakdown[i].TotalTax = math.Floor(result.TaxBreakdown[i].TotalTax*multiplier) / multiplier
		case "ceil":
			result.TaxBreakdown[i].TotalTax = math.Ceil(result.TaxBreakdown[i].TotalTax*multiplier) / multiplier
		}
	}
}

// validateInput validates the tax calculation input
func (tc *TaxCalculator) validateInput(input TaxCalculationInput) []string {
	errors := []string{}

	if len(input.Items) == 0 {
		errors = append(errors, "no items provided for tax calculation")
	}

	for i, item := range input.Items {
		if item.ID == "" {
			errors = append(errors, fmt.Sprintf("item %d missing ID", i))
		}
		if item.TotalAmount < 0 {
			errors = append(errors, fmt.Sprintf("item %s has negative amount", item.ID))
		}
		if item.Quantity <= 0 {
			errors = append(errors, fmt.Sprintf("item %s has invalid quantity", item.ID))
		}
	}

	if input.BillingAddress.Country == "" && input.ShippingAddress.Country == "" {
		errors = append(errors, "no valid address provided for tax calculation")
	}

	if input.TransactionDate.IsZero() {
		errors = append(errors, "transaction date is required")
	}

	return errors
}

// validateResult validates the tax calculation result
func (tc *TaxCalculator) validateResult(result TaxCalculationResult) []string {
	warnings := []string{}

	// Check for unusually high tax rates
	if result.EffectiveRate > 50 {
		warnings = append(warnings, fmt.Sprintf("unusually high effective tax rate: %.2f%%", result.EffectiveRate))
	}

	// Check for negative tax amounts
	if result.TotalTax < 0 {
		warnings = append(warnings, "negative total tax amount calculated")
	}

	// Check for inconsistent totals
	expectedTotal := result.Subtotal + result.TotalTax
	if math.Abs(expectedTotal-result.GrandTotal) > 0.01 {
		warnings = append(warnings, "inconsistent total calculation detected")
	}

	return warnings
}

// CalculateTaxInclusive calculates tax-inclusive pricing
func (tc *TaxCalculator) CalculateTaxInclusive(input TaxCalculationInput) TaxCalculationResult {
	// First calculate tax normally
	result := tc.CalculateTax(input)

	if !tc.Configuration.TaxInclusivePricing {
		return result
	}

	// For tax-inclusive pricing, the total amount includes tax
	// We need to extract the tax from the total
	for i := range result.TaxBreakdown {
		breakdown := &result.TaxBreakdown[i]
		if breakdown.TotalTax > 0 {
			// Calculate tax-exclusive amount
			taxRate := breakdown.TotalTax / breakdown.ItemAmount
			taxExclusiveAmount := breakdown.ItemAmount / (1 + taxRate)
			taxAmount := breakdown.ItemAmount - taxExclusiveAmount

			breakdown.TaxableAmount = taxExclusiveAmount
			breakdown.TotalTax = taxAmount
		}
	}

	// Recalculate totals
	result.TotalTax = 0
	result.TaxableAmount = 0
	for _, breakdown := range result.TaxBreakdown {
		result.TotalTax += breakdown.TotalTax
		result.TaxableAmount += breakdown.TaxableAmount
	}

	result.Subtotal = result.TaxableAmount
	result.GrandTotal = result.Subtotal + result.TotalTax

	return result
}

// GetTaxSummary returns a summary of tax calculations
func GetTaxSummary(results []TaxCalculationResult) map[string]interface{} {
	summary := map[string]interface{}{
		"total_transactions": len(results),
		"total_subtotal":     0.0,
		"total_tax":          0.0,
		"total_grand_total":  0.0,
		"average_tax_rate":   0.0,
		"jurisdiction_totals": make(map[TaxJurisdiction]float64),
		"tax_type_totals":    make(map[TaxType]float64),
	}

	totalSubtotal := 0.0
	totalTax := 0.0
	totalGrandTotal := 0.0
	jurisdictionTotals := make(map[TaxJurisdiction]float64)
	taxTypeTotals := make(map[TaxType]float64)

	for _, result := range results {
		totalSubtotal += result.Subtotal
		totalTax += result.TotalTax
		totalGrandTotal += result.GrandTotal

		for jurisdiction, amount := range result.JurisdictionTotals {
			jurisdictionTotals[jurisdiction] += amount
		}

		for taxType, amount := range result.TaxTypeTotals {
			taxTypeTotals[taxType] += amount
		}
	}

	summary["total_subtotal"] = totalSubtotal
	summary["total_tax"] = totalTax
	summary["total_grand_total"] = totalGrandTotal
	summary["jurisdiction_totals"] = jurisdictionTotals
	summary["tax_type_totals"] = taxTypeTotals

	if totalSubtotal > 0 {
		summary["average_tax_rate"] = (totalTax / totalSubtotal) * 100
	}

	return summary
}

// CalculateBestTaxStrategy calculates the best tax strategy for given scenarios
func CalculateBestTaxStrategy(scenarios []TaxCalculationInput) (*TaxCalculationInput, error) {
	if len(scenarios) == 0 {
		return nil, errors.New("no scenarios provided")
	}

	bestScenario := scenarios[0]
	lowestTax := math.MaxFloat64

	for _, scenario := range scenarios {
		result := Calculate(scenario)
		if result.IsValid && result.TotalTax < lowestTax {
			lowestTax = result.TotalTax
			bestScenario = scenario
		}
	}

	return &bestScenario, nil
}
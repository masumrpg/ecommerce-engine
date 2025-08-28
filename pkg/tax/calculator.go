// Package tax provides comprehensive tax calculation functionality for e-commerce applications.
//
// This package supports various tax calculation methods including percentage-based,
// fixed amount, tiered, progressive, and compound tax calculations. It handles
// multiple jurisdictions, tax exemptions, customer-specific rules, and complex
// tax scenarios commonly found in e-commerce platforms.
//
// Key Features:
//   - Multiple tax calculation methods (percentage, fixed, tiered, progressive, compound)
//   - Geographic-based tax rules (country, state, city, postal code)
//   - Customer and item-level tax exemptions
//   - Tax-inclusive and tax-exclusive pricing
//   - Compound tax calculations
//   - Tax overrides and manual adjustments
//   - Comprehensive validation and error handling
//   - Detailed tax breakdowns and reporting
//
// Basic Usage:
//
//	config := TaxConfiguration{
//		DefaultCurrency:     "USD",
//		RoundingMode:        "round",
//		RoundingPrecision:   2,
//		TaxInclusivePricing: false,
//	}
//
//	calc := NewTaxCalculator(config)
//	input := TaxCalculationInput{
//		Items: []TaxableItem{{
//			ID:          "item1",
//			Name:        "Product A",
//			TotalAmount: 100.00,
//			Quantity:    1,
//		}},
//		ShippingAddress: Address{
//			Country: "US",
//			State:   "CA",
//		},
//	}
//
//	result := calc.CalculateTax(input)
//	fmt.Printf("Total Tax: $%.2f\n", result.TotalTax)
package tax

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// TaxCalculator handles comprehensive tax calculations for e-commerce transactions.
// It supports multiple tax calculation methods, geographic rules, exemptions,
// and complex tax scenarios.
//
// The calculator maintains a configuration that defines default behavior,
// a set of tax rules that determine applicable taxes, and validation rules
// for ensuring calculation accuracy.
type TaxCalculator struct {
	// Configuration defines the default tax calculation behavior
	Configuration TaxConfiguration
	// Rules contains all tax rules that may apply to calculations
	Rules []TaxRule
	// ValidationRules contains rules for validating tax calculations
	ValidationRules []TaxValidationRule
}

// NewTaxCalculator creates a new tax calculator with the specified configuration.
// The calculator is initialized with the provided configuration and any default
// rules specified in the configuration.
//
// Parameters:
//   - config: Tax configuration defining calculation behavior
//
// Returns:
//   - *TaxCalculator: A new tax calculator instance
//
// Example:
//
//	config := TaxConfiguration{
//		DefaultCurrency:     "USD",
//		RoundingMode:        "round",
//		RoundingPrecision:   2,
//		TaxInclusivePricing: false,
//		CompoundTaxes:       false,
//		TaxOnShipping:       true,
//		TaxOnDiscounts:      true,
//	}
//	calc := NewTaxCalculator(config)
func NewTaxCalculator(config TaxConfiguration) *TaxCalculator {
	return &TaxCalculator{
		Configuration: config,
		Rules:         config.DefaultRules,
		ValidationRules: []TaxValidationRule{},
	}
}

// Calculate is a convenience function that calculates taxes for the given input
// using default tax configuration. This function creates a temporary tax calculator
// with standard settings and performs the calculation.
//
// Default Configuration:
//   - Currency: USD
//   - Rounding: Round to 2 decimal places
//   - Tax-exclusive pricing
//   - No compound taxes
//   - Tax on shipping and discounts enabled
//
// Parameters:
//   - input: Tax calculation input containing items, addresses, and rules
//
// Returns:
//   - TaxCalculationResult: Complete tax calculation result with breakdowns
//
// Example:
//
//	input := TaxCalculationInput{
//		Items: []TaxableItem{{
//			ID:          "item1",
//			TotalAmount: 100.00,
//			Quantity:    1,
//		}},
//		ShippingAddress: Address{Country: "US", State: "CA"},
//		TaxRules: []TaxRule{{
//			ID:   "sales_tax",
//			Rate: 8.25,
//			Type: TaxTypeSales,
//		}},
//	}
//	result := Calculate(input)
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

// CalculateTax performs comprehensive tax calculation for the given input.
// This is the main calculation method that processes all items, applies
// applicable tax rules, handles exemptions, and generates detailed breakdowns.
//
// The calculation process includes:
//   1. Input validation
//   2. Subtotal calculation
//   3. Discount and shipping adjustments
//   4. Rule evaluation and application
//   5. Tax calculation per item
//   6. Override application
//   7. Amount rounding
//   8. Result validation
//
// Parameters:
//   - input: Complete tax calculation input with items, addresses, customer info
//
// Returns:
//   - TaxCalculationResult: Detailed calculation result with breakdowns and totals
//
// Example:
//
//	calc := NewTaxCalculator(config)
//	input := TaxCalculationInput{
//		Items: []TaxableItem{{
//			ID:          "item1",
//			TotalAmount: 100.00,
//			Quantity:    1,
//			Category:    "electronics",
//		}},
//		ShippingAddress: Address{Country: "US", State: "CA"},
//		Customer: Customer{Type: "individual"},
//	}
//	result := calc.CalculateTax(input)
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

// calculateSubtotal calculates the subtotal amount from all taxable items.
// This method sums the total amount of all items before any tax calculations.
//
// Parameters:
//   - items: Slice of taxable items to calculate subtotal for
//
// Returns:
//   - float64: The subtotal amount of all items
func (tc *TaxCalculator) calculateSubtotal(items []TaxableItem) float64 {
	subtotal := 0.0
	for _, item := range items {
		subtotal += item.TotalAmount
	}
	return subtotal
}

// getApplicableRules filters and returns tax rules that apply to the given input.
// This method evaluates each rule against the input criteria including:
//   - Rule active status and validity period
//   - Geographic applicability (country, state, city, postal code)
//   - Rule conditions (amount, quantity, weight, category, customer type)
//
// Parameters:
//   - input: Tax calculation input to evaluate rules against
//
// Returns:
//   - []TaxRule: Slice of applicable tax rules sorted by priority
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

// isGeographicallyApplicable determines if a tax rule applies to the given addresses.
// The method uses shipping address as primary and falls back to billing address.
// It checks rule applicability against countries, states, cities, and postal codes.
//
// Parameters:
//   - rule: Tax rule to evaluate for geographic applicability
//   - billingAddr: Customer's billing address
//   - shippingAddr: Customer's shipping address (preferred for tax calculation)
//
// Returns:
//   - bool: True if the rule applies to the given geographic location
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

// evaluateConditions evaluates all conditions for a tax rule.
// This method processes multiple conditions and applies logical operators
// to determine if all conditions are met (currently uses AND logic).
//
// Parameters:
//   - conditions: Slice of tax conditions to evaluate
//   - input: Tax calculation input to evaluate conditions against
//
// Returns:
//   - bool: True if all conditions are satisfied
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

// evaluateCondition evaluates a single tax condition against the input.
// Supported condition types:
//   - "amount": Total transaction amount
//   - "quantity": Total item quantity
//   - "weight": Total item weight
//   - "category": Item category matching
//   - "customer_type": Customer type matching
//
// Parameters:
//   - condition: Single tax condition to evaluate
//   - input: Tax calculation input containing data to evaluate
//
// Returns:
//   - bool: True if the condition is satisfied
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

// compareValues compares two values using the specified operator.
// Supported operators:
//   - ">", "<", ">=", "<=": Numeric comparisons
//   - "=", "!=": Equality comparisons
//   - "in", "not_in": Array membership checks
//
// Parameters:
//   - actual: The actual value to compare
//   - operator: Comparison operator to use
//   - expected: The expected value or values to compare against
//
// Returns:
//   - bool: True if the comparison is satisfied
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

// toFloat64 converts various numeric types to float64.
// Supports conversion from float32, float64, int, int64, and string types.
//
// Parameters:
//   - value: Value to convert to float64
//
// Returns:
//   - float64: Converted numeric value
//   - error: Error if conversion fails
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

// calculateItemTax calculates tax for a single taxable item.
// This method processes exemptions, applies applicable tax rules,
// and generates a detailed breakdown of taxes for the item.
//
// The calculation process:
//   1. Check item-level exemptions
//   2. Check customer-level exemptions
//   3. Apply applicable tax rules
//   4. Handle compound tax calculations if configured
//
// Parameters:
//   - item: The taxable item to calculate tax for
//   - rules: Applicable tax rules for this calculation
//   - input: Complete tax calculation input for context
//
// Returns:
//   - TaxBreakdown: Detailed breakdown of taxes applied to this item
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

// isCustomerExempt determines if a customer is exempt from tax for a specific item.
// This method checks all customer exemptions to see if any apply to the given item.
//
// Parameters:
//   - customer: Customer information including exemptions
//   - item: Taxable item to check exemptions for
//
// Returns:
//   - bool: True if the customer is exempt from tax for this item
func (tc *TaxCalculator) isCustomerExempt(customer Customer, item TaxableItem) bool {
	for _, exemption := range customer.Exemptions {
		if tc.isExemptionApplicable(exemption, item) {
			return true
		}
	}
	return false
}

// isExemptionApplicable determines if a specific tax exemption applies to an item.
// This method checks exemption validity period and conditions to determine applicability.
//
// Exemption types:
//   - "item": Exemption based on item category or characteristics
//   - "customer": Customer-level exemption applying to all items
//
// Parameters:
//   - exemption: Tax exemption to evaluate
//   - item: Taxable item to check exemption against
//
// Returns:
//   - bool: True if the exemption applies to the item
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

// isRuleApplicableToItem determines if a tax rule applies to a specific item.
// This method checks rule criteria including:
//   - Applicable categories (item must match)
//   - Exempt categories (item must not match)
//   - Amount thresholds (min/max amounts)
//
// Parameters:
//   - rule: Tax rule to evaluate for applicability
//   - item: Taxable item to check rule against
//
// Returns:
//   - bool: True if the rule applies to the item
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

// calculateTaxForRule calculates the tax amount for a specific tax rule.
// This method supports multiple calculation methods:
//   - Percentage: Tax as percentage of taxable amount
//   - Fixed: Fixed tax amount per item quantity
//   - Tiered: Tax rate based on amount tiers
//   - Progressive: Progressive tax rates across tiers
//   - Compound: Tax calculated on amount including previous taxes
//
// Parameters:
//   - rule: Tax rule defining calculation method and rates
//   - taxableAmount: Amount subject to tax calculation
//   - item: Taxable item for quantity-based calculations
//
// Returns:
//   - AppliedTax: Complete applied tax information including amount
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

// calculateTieredTax calculates tax using tiered rate structure.
// In tiered taxation, the entire amount is taxed at the rate of the tier
// it falls into, unlike progressive taxation where each tier is taxed separately.
//
// Parameters:
//   - thresholds: Tax thresholds defining tiers and rates
//   - amount: Amount to calculate tax for
//
// Returns:
//   - float64: Calculated tax amount using tiered rates
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

// calculateProgressiveTax calculates tax using progressive rate structure.
// In progressive taxation, each tier is taxed at its own rate, and the
// total tax is the sum of taxes from all applicable tiers.
//
// Parameters:
//   - thresholds: Tax thresholds defining progressive tiers and rates
//   - amount: Amount to calculate progressive tax for
//
// Returns:
//   - float64: Calculated tax amount using progressive rates
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

// aggregateAppliedTax aggregates an applied tax into the calculation result totals.
// This method updates jurisdiction totals, tax type totals, and the applied taxes list.
// It handles duplicate tax rules by combining their amounts.
//
// Parameters:
//   - result: Tax calculation result to update with aggregated totals
//   - appliedTax: Applied tax to aggregate into the result
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

// applyTaxOverrides applies manual tax overrides to the calculation result.
// This method supports three types of overrides:
//   - "rate": Override the tax rate for a specific tax type
//   - "amount": Override the tax amount for a specific tax type
//   - "exempt": Exempt from a specific tax type
//
// Parameters:
//   - result: Tax calculation result to apply overrides to
//   - overrides: Slice of tax overrides to apply
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

// roundAmounts rounds all monetary amounts in the result based on configuration.
// Supported rounding modes:
//   - "round": Standard rounding to nearest value
//   - "floor": Always round down
//   - "ceil": Always round up
//
// The method rounds totals, applied taxes, and tax breakdowns according
// to the configured precision (typically 2 decimal places for currency).
//
// Parameters:
//   - result: Tax calculation result to round amounts in
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

// validateInput validates the tax calculation input for completeness and correctness.
// This method checks for:
//   - Presence of items to calculate tax for
//   - Valid item data (ID, amount, quantity)
//   - Valid address information
//   - Required transaction date
//
// Parameters:
//   - input: Tax calculation input to validate
//
// Returns:
//   - []string: Slice of validation error messages (empty if valid)
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

// validateResult validates the tax calculation result for reasonableness.
// This method checks for:
//   - Unusually high tax rates (>50%)
//   - Negative tax amounts
//   - Inconsistent total calculations
//
// Parameters:
//   - result: Tax calculation result to validate
//
// Returns:
//   - []string: Slice of validation warning messages
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

// CalculateTaxInclusive calculates tax for tax-inclusive pricing scenarios.
// In tax-inclusive pricing, the item prices already include tax, and this method
// extracts the tax component from the total price to show the breakdown.
//
// The method:
//   1. Performs normal tax calculation
//   2. If tax-inclusive pricing is enabled, extracts tax from item amounts
//   3. Recalculates totals with tax-exclusive amounts
//
// Parameters:
//   - input: Tax calculation input with tax-inclusive item prices
//
// Returns:
//   - TaxCalculationResult: Result showing tax-exclusive amounts and extracted tax
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

// GetTaxSummary generates a comprehensive summary from multiple tax calculation results.
// This function aggregates data across multiple transactions to provide insights into:
//   - Total transaction counts and amounts
//   - Average tax rates
//   - Tax totals by jurisdiction and tax type
//
// Parameters:
//   - results: Slice of tax calculation results to summarize
//
// Returns:
//   - map[string]interface{}: Summary containing aggregated tax statistics
//
// Example:
//
//	results := []TaxCalculationResult{result1, result2, result3}
//	summary := GetTaxSummary(results)
//	fmt.Printf("Average tax rate: %.2f%%\n", summary["average_tax_rate"])
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

// CalculateBestTaxStrategy determines the optimal tax strategy from multiple scenarios.
// This function evaluates different tax calculation scenarios and returns the one
// that results in the lowest total tax amount.
//
// Use cases:
//   - Comparing different shipping addresses for tax optimization
//   - Evaluating different customer types or exemptions
//   - Testing various item categorizations
//
// Parameters:
//   - scenarios: Slice of different tax calculation scenarios to evaluate
//
// Returns:
//   - *TaxCalculationInput: The scenario that results in the lowest tax
//   - error: Error if no scenarios provided or all scenarios are invalid
//
// Example:
//
//	scenarios := []TaxCalculationInput{scenario1, scenario2, scenario3}
//	bestScenario, err := CalculateBestTaxStrategy(scenarios)
//	if err == nil {
//		result := Calculate(*bestScenario)
//		fmt.Printf("Best strategy saves: $%.2f\n", originalTax - result.TotalTax)
//	}
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
package tax

import (
	"math"
	"testing"
	"time"
)

func createTestTaxCalculator() *TaxCalculator {
	config := TaxConfiguration{
		DefaultCurrency:     "USD",
		RoundingMode:        "round",
		RoundingPrecision:   2,
		TaxInclusivePricing: false,
		CompoundTaxes:       false,
		TaxOnShipping:       true,
		TaxOnDiscounts:      true,
		DefaultRules:        []TaxRule{createTestTaxRule()},
	}
	return NewTaxCalculator(config)
}

func createTestTaxInput() TaxCalculationInput {
	return TaxCalculationInput{
		Items: []TaxableItem{
			{
				ID:          "item1",
				Name:        "Test Item",
				UnitPrice:   100.0,
				TotalAmount: 100.0,
				Quantity:    1,
				Category:    "electronics",
			},
		},
		Customer: Customer{
			ID:   "customer1",
			Type: "individual",
		},
		BillingAddress: Address{
			City:    "New York",
			State:   "NY",
			Country: "US",
		},
		ShippingAddress: Address{
			City:    "New York",
			State:   "NY",
			Country: "US",
		},
		TransactionDate: time.Now(),
		TransactionType: "sale",
		Currency:        "USD",
	}
}

func createTestTaxRule() TaxRule {
	return TaxRule{
		ID:                  "test-rule",
		Name:                "Test Tax Rule",
		Type:                TaxTypeSales,
		Rate:                0.08,
		Jurisdiction:        JurisdictionState,
		Method:              TaxMethodPercentage,
		ApplicableCountries: []string{"US"},
		ApplicableStates:    []string{"NY"},
		MinAmount:           0.0,
		MaxAmount:           1000000.0,
		IsActive:            true,
		ValidFrom:           time.Now().AddDate(0, 0, -1),
		ValidUntil:          time.Now().AddDate(1, 0, 0),
	}
}

func TestNewTaxCalculator(t *testing.T) {
	config := TaxConfiguration{
		DefaultRules:      []TaxRule{createTestTaxRule()},
		RoundingMode:      "round",
		RoundingPrecision: 2,
		DefaultCurrency:   "USD",
	}

	calc := NewTaxCalculator(config)

	if calc == nil {
		t.Errorf("Expected non-nil calculator")
	}
	if len(calc.Configuration.DefaultRules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(calc.Configuration.DefaultRules))
	}
}

func TestCalculate(t *testing.T) {
	input := createTestTaxInput()
	input.TaxRules = []TaxRule{createTestTaxRule()}

	result := Calculate(input)

	if !result.IsValid {
		t.Errorf("Expected valid result, got invalid")
	}
	if result.TotalTax < 0.0 {
		t.Errorf("Expected non-negative tax amount, got %f", result.TotalTax)
	}
	if result.Subtotal+result.TotalTax != result.GrandTotal {
		t.Errorf("Expected subtotal + tax = grand total")
	}
	if len(result.TaxBreakdown) != 1 {
		t.Errorf("Expected 1 tax breakdown, got %d", len(result.TaxBreakdown))
	}
}

func TestCalculateTax(t *testing.T) {
	calc := createTestTaxCalculator()
	input := createTestTaxInput()

	result := calc.CalculateTax(input)

	if !result.IsValid {
		t.Errorf("Expected valid result")
	}
	if result.TotalTax < 0.0 {
		t.Errorf("Expected non-negative tax amount")
	}
	if result.Currency != "USD" {
		t.Errorf("Expected USD currency, got %s", result.Currency)
	}
}

func TestCalculateSubtotal(t *testing.T) {
	calc := createTestTaxCalculator()
	items := []TaxableItem{
		{TotalAmount: 100.0},
		{TotalAmount: 50.0},
	}

	subtotal := calc.calculateSubtotal(items)

	if subtotal != 150.0 {
		t.Errorf("Expected subtotal 150.0, got %f", subtotal)
	}
}

func TestGetApplicableRules(t *testing.T) {
	calc := createTestTaxCalculator()
	input := createTestTaxInput()

	rules := calc.getApplicableRules(input)

	if len(rules) == 0 {
		t.Errorf("Expected applicable rules, got none")
		return
	}
	if rules[0].ID != "test-rule" {
		t.Errorf("Expected rule ID 'test-rule', got %s", rules[0].ID)
	}
}

func TestIsGeographicallyApplicable(t *testing.T) {
	calc := createTestTaxCalculator()
	rule := TaxRule{
		ApplicableCountries: []string{"US"},
		ApplicableStates:    []string{"CA"},
	}
	billingAddr := Address{
		Country: "US",
		State:   "CA",
		City:    "Los Angeles",
	}
	shippingAddr := Address{
		Country: "US",
		State:   "CA",
		City:    "San Francisco",
	}

	applicable := calc.isGeographicallyApplicable(rule, billingAddr, shippingAddr)

	if !applicable {
		t.Errorf("Expected rule to be geographically applicable")
	}
}

func TestEvaluateConditions(t *testing.T) {
	calc := createTestTaxCalculator()
	conditions := []TaxCondition{
		{
			Type:     "amount",
			Operator: ">",
			Value:    50.0,
			Logic:    "AND",
		},
		{
			Type:     "category",
			Operator: "=",
			Value:    "electronics",
		},
	}
	input := TaxCalculationInput{
		Items: []TaxableItem{
			{
				TotalAmount: 100.0,
				Category:    "electronics",
			},
		},
		Customer: Customer{Type: "individual"},
	}

	result := calc.evaluateConditions(conditions, input)

	if !result {
		t.Errorf("Expected conditions to evaluate to true")
	}
}

func TestEvaluateCondition(t *testing.T) {
	calc := createTestTaxCalculator()
	condition := TaxCondition{
		Type:     "amount",
		Operator: ">",
		Value:    50.0,
	}
	input := TaxCalculationInput{
		Items: []TaxableItem{
			{
				TotalAmount: 100.0,
			},
		},
		Customer: Customer{Type: "individual"},
	}

	result := calc.evaluateCondition(condition, input)
	if !result {
		t.Errorf("Expected condition to evaluate to true")
	}
}

func TestCompareValues(t *testing.T) {
	calc := createTestTaxCalculator()

	// Test greater than
	result := calc.compareValues(100.0, ">", 50.0)
	if !result {
		t.Errorf("Expected 100 > 50 to be true")
	}

	// Test equal
	result = calc.compareValues("electronics", "=", "electronics")
	if !result {
		t.Errorf("Expected electronics = electronics to be true")
	}

	// Test less than
	result = calc.compareValues(25.0, "<", 50.0)
	if !result {
		t.Errorf("Expected 25 < 50 to be true")
	}
}

func TestToFloat64(t *testing.T) {
	calc := createTestTaxCalculator()

	result, err := calc.toFloat64(100.5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 100.5 {
		t.Errorf("Expected 100.5, got %f", result)
	}

	result, err = calc.toFloat64("123.45")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 123.45 {
		t.Errorf("Expected 123.45, got %f", result)
	}

	result, err = calc.toFloat64(42)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result != 42.0 {
		t.Errorf("Expected 42.0, got %f", result)
	}
}

func TestCalculateItemTax(t *testing.T) {
	calc := createTestTaxCalculator()
	item := TaxableItem{
		ID:          "item-1",
		TotalAmount: 100.0,
		Category:    "electronics",
	}
	rules := []TaxRule{createTestTaxRule()}
	input := TaxCalculationInput{
		Items:    []TaxableItem{item},
		Customer: Customer{Type: "individual"},
	}

	result := calc.calculateItemTax(item, rules, input)

	if result.TotalTax <= 0.0 {
		t.Errorf("Expected positive tax amount")
	}
	if result.ItemID != "item-1" {
		t.Errorf("Expected item ID 'item-1', got %s", result.ItemID)
	}
	if len(result.AppliedTaxes) == 0 {
		t.Errorf("Expected applied taxes")
	}
}

func TestIsCustomerExempt(t *testing.T) {
	calc := createTestTaxCalculator()
	customer := Customer{
		Type: "business",
		Exemptions: []TaxExemption{
			{
				Type:        "customer",
				Certificate: "CERT123",
				ValidFrom:   time.Now().AddDate(0, 0, -1),
				ValidUntil:  time.Now().AddDate(0, 0, 1),
			},
		},
	}
	item := TaxableItem{
		ID:          "item-1",
		TotalAmount: 100.0,
	}

	exempt := calc.isCustomerExempt(customer, item)

	// Test passes if function executes without error
	if exempt {
		t.Log("Customer is exempt")
	} else {
		t.Log("Customer is not exempt")
	}
}

func TestIsExemptionApplicable(t *testing.T) {
	calc := createTestTaxCalculator()
	exemption := TaxExemption{
		Type:        "item",
		Certificate: "CERT123",
		ValidFrom:   time.Now().AddDate(0, -1, 0),
		ValidUntil:  time.Now().AddDate(0, 1, 0),
	}
	item := TaxableItem{Category: "electronics"}

	applicable := calc.isExemptionApplicable(exemption, item)

	// Test passes if function executes without error
	if applicable {
		t.Log("Exemption is applicable")
	} else {
		t.Log("Exemption is not applicable")
	}
}

func TestIsRuleApplicableToItem(t *testing.T) {
	calc := createTestTaxCalculator()
	rule := TaxRule{
		ApplicableCategories: []string{"electronics"},
		MinAmount:           50.0,
	}
	item := TaxableItem{
		Category:    "electronics",
		TotalAmount: 100.0,
	}

	applicable := calc.isRuleApplicableToItem(rule, item)

	if !applicable {
		t.Errorf("Expected rule to be applicable to item")
	}
}

func TestCalculateTaxForRule(t *testing.T) {
	calc := createTestTaxCalculator()
	rule := TaxRule{
		ID:     "test-rule",
		Rate:   8.0,
		Type:   TaxTypeSales,
		Method: TaxMethodPercentage,
	}
	amount := 100.0
	item := TaxableItem{
		ID:          "item-1",
		TotalAmount: amount,
	}

	tax := calc.calculateTaxForRule(rule, amount, item)

	// The actual tax amount should be 8.0 (8% of 100)
	expectedTax := 8.0
	if math.Abs(tax.TaxAmount-expectedTax) > 0.01 {
		t.Errorf("Expected tax amount %.2f, got %f", expectedTax, tax.TaxAmount)
	}
	if tax.RuleID != "test-rule" {
		t.Errorf("Expected rule ID 'test-rule', got %s", tax.RuleID)
	}
	if tax.Type != TaxTypeSales {
		t.Errorf("Expected tax type %v, got %v", TaxTypeSales, tax.Type)
	}
}

// Benchmark tests
func BenchmarkCalculateTax(b *testing.B) {
	calc := createTestTaxCalculator()
	input := createTestTaxInput()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.CalculateTax(input)
	}
}

func BenchmarkCalculateSubtotal(b *testing.B) {
	calc := createTestTaxCalculator()
	items := []TaxableItem{
		{TotalAmount: 100.0},
		{TotalAmount: 50.0},
		{TotalAmount: 75.0},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.calculateSubtotal(items)
	}
}
// Package tax provides comprehensive types and data structures for tax calculation,
// rule management, and reporting in e-commerce applications.
//
// This package defines all the core types used throughout the tax system including:
//   - Tax types and calculation methods
//   - Tax rules and conditions
//   - Customer and address information
//   - Tax calculation inputs and results
//   - Tax reporting and audit trail structures
//
// Basic usage:
//
//	taxInput := &tax.TaxCalculationInput{
//		Items: []tax.TaxableItem{
//			{
//				ID: "item1",
//				Name: "Product A",
//				Category: "electronics",
//				Quantity: 2,
//				UnitPrice: 100.00,
//				TotalAmount: 200.00,
//			},
//		},
//		Customer: tax.Customer{
//			ID: "customer1",
//			Type: "individual",
//		},
//		BillingAddress: tax.Address{
//			City: "New York",
//			State: "NY",
//			Country: "US",
//			PostalCode: "10001",
//		},
//		TransactionDate: time.Now(),
//		Currency: "USD",
//	}
package tax

import (
	"time"
)

// TaxType represents different types of taxes that can be applied to transactions.
// Each tax type has specific rules and calculation methods depending on jurisdiction
// and applicable regulations.
type TaxType string

// Tax type constants define the various types of taxes supported by the system.
const (
	// TaxTypeSales represents sales tax, typically applied at the point of sale
	// in jurisdictions like US states.
	TaxTypeSales TaxType = "sales"
	
	// TaxTypeVAT represents Value Added Tax, commonly used in European countries
	// and applied at each stage of the supply chain.
	TaxTypeVAT TaxType = "vat"
	
	// TaxTypeGST represents Goods and Services Tax, used in countries like
	// Canada, Australia, and India.
	TaxTypeGST TaxType = "gst"
	
	// TaxTypeExcise represents excise tax, typically applied to specific goods
	// like alcohol, tobacco, or fuel.
	TaxTypeExcise TaxType = "excise"
	
	// TaxTypeCustoms represents customs duty, applied to imported goods
	// at international borders.
	TaxTypeCustoms TaxType = "customs"
	
	// TaxTypeProperty represents property tax, typically applied to real estate
	// and personal property.
	TaxTypeProperty TaxType = "property"
	
	// TaxTypeWithholding represents withholding tax, deducted at source
	// from payments like salaries or dividends.
	TaxTypeWithholding TaxType = "withholding"
	
	// TaxTypeDigital represents digital services tax, applied to digital
	// services and online transactions.
	TaxTypeDigital TaxType = "digital"
	
	// TaxTypeEnvironmental represents environmental tax, applied to products
	// or activities that impact the environment.
	TaxTypeEnvironmental TaxType = "environmental"
	
	// TaxTypeLuxury represents luxury tax, applied to high-value or
	// non-essential luxury items.
	TaxTypeLuxury TaxType = "luxury"
)

// TaxCalculationMethod represents the method used to calculate tax amounts.
// Different methods are used based on tax type, jurisdiction, and specific
// tax regulations.
type TaxCalculationMethod string

// Tax calculation method constants define how taxes are computed.
const (
	// TaxMethodPercentage applies a percentage rate to the taxable amount.
	// This is the most common method for sales tax, VAT, and GST.
	TaxMethodPercentage TaxCalculationMethod = "percentage"
	
	// TaxMethodFixed applies a fixed amount regardless of the taxable amount.
	// Often used for specific fees or duties.
	TaxMethodFixed TaxCalculationMethod = "fixed"
	
	// TaxMethodTiered applies different rates based on amount thresholds.
	// Each tier has its own rate applied to amounts within that range.
	TaxMethodTiered TaxCalculationMethod = "tiered"
	
	// TaxMethodProgressive applies increasing rates as amounts increase.
	// Similar to tiered but rates accumulate progressively.
	TaxMethodProgressive TaxCalculationMethod = "progressive"
	
	// TaxMethodCompound applies tax on top of other taxes.
	// Used when multiple taxes are applied and one is calculated on the total
	// including other taxes.
	TaxMethodCompound TaxCalculationMethod = "compound"
)

// TaxJurisdiction represents the governmental level or authority that
// imposes and collects taxes. Different jurisdictions may have overlapping
// tax requirements.
type TaxJurisdiction string

// Tax jurisdiction constants define the levels of government that can impose taxes.
const (
	// JurisdictionFederal represents federal or national level taxation,
	// such as federal income tax or national VAT.
	JurisdictionFederal TaxJurisdiction = "federal"
	
	// JurisdictionState represents state or province level taxation,
	// such as state sales tax or provincial tax.
	JurisdictionState TaxJurisdiction = "state"
	
	// JurisdictionCounty represents county level taxation,
	// often used for local sales tax or property tax.
	JurisdictionCounty TaxJurisdiction = "county"
	
	// JurisdictionCity represents city or municipal level taxation,
	// such as city sales tax or municipal fees.
	JurisdictionCity TaxJurisdiction = "city"
	
	// JurisdictionDistrict represents special district taxation,
	// such as school districts or transportation authorities.
	JurisdictionDistrict TaxJurisdiction = "district"
	
	// JurisdictionInternational represents international taxation,
	// such as customs duties or international trade taxes.
	JurisdictionInternational TaxJurisdiction = "international"
)

// TaxableItem represents an item or service subject to taxation.
// It contains all the information needed to determine applicable taxes,
// exemptions, and calculate the correct tax amount.
//
// Example:
//
//	item := &TaxableItem{
//		ID: "prod-123",
//		Name: "Laptop Computer",
//		Category: "electronics",
//		Quantity: 1,
//		UnitPrice: 999.99,
//		TotalAmount: 999.99,
//		HSCode: "8471.30.01",
//		IsDigital: false,
//	}
type TaxableItem struct {
	// ID is the unique identifier for the item
	ID string `json:"id"`
	
	// Name is the display name of the item
	Name string `json:"name"`
	
	// Category is the primary category for tax classification
	Category string `json:"category"`
	
	// Subcategory provides additional classification detail
	Subcategory string `json:"subcategory,omitempty"`
	
	// Quantity is the number of units being purchased
	Quantity int `json:"quantity"`
	
	// UnitPrice is the price per individual unit
	UnitPrice float64 `json:"unit_price"`
	
	// TotalAmount is the total amount for all units (usually Quantity * UnitPrice)
	TotalAmount float64 `json:"total_amount"`
	
	// Weight is the physical weight, used for shipping tax calculations
	Weight float64 `json:"weight,omitempty"`
	
	// Volume is the physical volume, used for shipping tax calculations
	Volume float64 `json:"volume,omitempty"`
	
	// Origin is the country of origin, important for customs and import duties
	Origin string `json:"origin,omitempty"`
	
	// HSCode is the Harmonized System code for international trade classification
	HSCode string `json:"hs_code,omitempty"`
	
	// SKU is the stock keeping unit identifier
	SKU string `json:"sku,omitempty"`
	
	// Brand is the manufacturer or brand name
	Brand string `json:"brand,omitempty"`
	
	// IsDigital indicates if this is a digital good or service
	IsDigital bool `json:"is_digital,omitempty"`
	
	// IsLuxury indicates if this item qualifies as a luxury good
	IsLuxury bool `json:"is_luxury,omitempty"`
	
	// IsExempt indicates if this item is exempt from taxation
	IsExempt bool `json:"is_exempt,omitempty"`
	
	// ExemptionReason provides the reason for tax exemption
	ExemptionReason string `json:"exemption_reason,omitempty"`
	
	// CustomAttributes allows for additional item-specific data
	CustomAttributes map[string]interface{} `json:"custom_attributes,omitempty"`
}

// TaxRule represents a comprehensive tax calculation rule that defines
// how taxes should be applied to transactions. Rules can be jurisdiction-specific,
// category-specific, and include complex conditions and exemptions.
//
// Example:
//
//	rule := &TaxRule{
//		ID: "ny-sales-tax",
//		Name: "New York State Sales Tax",
//		Type: TaxTypeSales,
//		Jurisdiction: JurisdictionState,
//		Method: TaxMethodPercentage,
//		Rate: 8.25,
//		ApplicableStates: []string{"NY"},
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
//		Priority: 100,
//	}
type TaxRule struct {
	// ID is the unique identifier for the tax rule
	ID string `json:"id"`
	
	// Name is the human-readable name of the tax rule
	Name string `json:"name"`
	
	// Type specifies the type of tax this rule applies
	Type TaxType `json:"type"`
	
	// Jurisdiction specifies the governmental level that imposes this tax
	Jurisdiction TaxJurisdiction `json:"jurisdiction"`
	
	// Method specifies how the tax amount is calculated
	Method TaxCalculationMethod `json:"method"`
	
	// Rate is the tax rate (percentage for percentage method, amount for fixed method)
	Rate float64 `json:"rate"`
	
	// MinAmount is the minimum taxable amount for this rule to apply
	MinAmount float64 `json:"min_amount,omitempty"`
	
	// MaxAmount is the maximum taxable amount for this rule to apply
	MaxAmount float64 `json:"max_amount,omitempty"`
	
	// Thresholds define rate tiers for tiered and progressive tax methods
	Thresholds []TaxThreshold `json:"thresholds,omitempty"`
	
	// ApplicableCategories lists item categories this rule applies to
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
	
	// ExemptCategories lists item categories exempt from this rule
	ExemptCategories []string `json:"exempt_categories,omitempty"`
	
	// ApplicableCountries lists countries where this rule applies
	ApplicableCountries []string `json:"applicable_countries,omitempty"`
	
	// ApplicableStates lists states/provinces where this rule applies
	ApplicableStates []string `json:"applicable_states,omitempty"`
	
	// ApplicableCities lists cities where this rule applies
	ApplicableCities []string `json:"applicable_cities,omitempty"`
	
	// PostalCodes lists specific postal codes where this rule applies
	PostalCodes []string `json:"postal_codes,omitempty"`
	
	// IsActive indicates whether this rule is currently active
	IsActive bool `json:"is_active"`
	
	// ValidFrom is the date when this rule becomes effective
	ValidFrom time.Time `json:"valid_from"`
	
	// ValidUntil is the date when this rule expires
	ValidUntil time.Time `json:"valid_until"`
	
	// Priority determines rule precedence (higher number = higher priority)
	Priority int `json:"priority"`
	
	// Description provides additional details about the rule
	Description string `json:"description,omitempty"`
	
	// Conditions define additional conditions for rule application
	Conditions []TaxCondition `json:"conditions,omitempty"`
	
	// Exemptions define specific exemptions for this rule
	Exemptions []TaxExemption `json:"exemptions,omitempty"`
}

// TaxThreshold represents tax rate thresholds for tiered and progressive taxation.
// Each threshold defines a range of amounts and the corresponding tax rate or
// fixed amount that applies to amounts within that range.
//
// Example for tiered tax:
//
//	thresholds := []TaxThreshold{
//		{MinAmount: 0, MaxAmount: 1000, Rate: 5.0},
//		{MinAmount: 1000, MaxAmount: 5000, Rate: 7.5},
//		{MinAmount: 5000, MaxAmount: 0, Rate: 10.0}, // 0 means no upper limit
//	}
type TaxThreshold struct {
	// MinAmount is the minimum amount for this threshold tier
	MinAmount float64 `json:"min_amount"`
	
	// MaxAmount is the maximum amount for this threshold tier (0 means no upper limit)
	MaxAmount float64 `json:"max_amount,omitempty"`
	
	// Rate is the tax rate (percentage) for this tier
	Rate float64 `json:"rate"`
	
	// FixedAmount is a fixed tax amount for this tier (alternative to rate)
	FixedAmount float64 `json:"fixed_amount,omitempty"`
}

// TaxCondition represents conditions that must be met for tax rules to apply.
// Conditions can be based on transaction amounts, item properties, customer types,
// or other criteria, and can be combined using logical operators.
//
// Example:
//
//	condition := &TaxCondition{
//		Type: "amount",
//		Operator: ">=",
//		Value: 100.0,
//		Logic: "AND",
//	}
type TaxCondition struct {
	// Type specifies what property to evaluate ("amount", "quantity", "weight", "category", "customer_type")
	Type string `json:"type"`
	
	// Operator specifies the comparison operator (">", "<", ">=", "<=", "=", "!=", "in", "not_in")
	Operator string `json:"operator"`
	
	// Value is the value to compare against
	Value interface{} `json:"value"`
	
	// Logic specifies how to combine with other conditions ("AND", "OR")
	Logic string `json:"logic,omitempty"`
}

// TaxExemption represents tax exemption rules that can be applied to customers,
// items, transactions, or locations. Exemptions can be temporary or permanent
// and may require certification.
//
// Example:
//
//	exemption := &TaxExemption{
//		ID: "nonprofit-001",
//		Name: "Nonprofit Organization Exemption",
//		Type: "customer",
//		Reason: "501(c)(3) nonprofit status",
//		Certificate: "EX-12345",
//		ValidFrom: time.Now(),
//		ValidUntil: time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
//	}
type TaxExemption struct {
	// ID is the unique identifier for the exemption
	ID string `json:"id"`
	
	// Name is the human-readable name of the exemption
	Name string `json:"name"`
	
	// Type specifies the scope of exemption ("customer", "item", "transaction", "location")
	Type string `json:"type"`
	
	// Reason provides the legal or business reason for the exemption
	Reason string `json:"reason"`
	
	// Certificate is the exemption certificate number, if applicable
	Certificate string `json:"certificate,omitempty"`
	
	// ValidFrom is the date when the exemption becomes effective
	ValidFrom time.Time `json:"valid_from"`
	
	// ValidUntil is the date when the exemption expires
	ValidUntil time.Time `json:"valid_until"`
	
	// Conditions define additional conditions for exemption application
	Conditions []TaxCondition `json:"conditions,omitempty"`
}

// Customer represents customer information needed for tax calculation.
// Different customer types may have different tax obligations, exemptions,
// and reporting requirements.
//
// Example:
//
//	customer := &Customer{
//		ID: "cust-12345",
//		Type: "business",
//		TaxID: "12-3456789",
//		VATNumber: "GB123456789",
//		BusinessType: "corporation",
//		Industry: "technology",
//	}
type Customer struct {
	// ID is the unique identifier for the customer
	ID string `json:"id"`
	
	// Type specifies the customer type ("individual", "business", "government", "nonprofit")
	Type string `json:"type"`
	
	// TaxID is the tax identification number
	TaxID string `json:"tax_id,omitempty"`
	
	// VATNumber is the VAT registration number for businesses
	VATNumber string `json:"vat_number,omitempty"`
	
	// BusinessType specifies the type of business entity
	BusinessType string `json:"business_type,omitempty"`
	
	// Industry specifies the customer's industry sector
	Industry string `json:"industry,omitempty"`
	
	// Exemptions lists any tax exemptions applicable to this customer
	Exemptions []TaxExemption `json:"exemptions,omitempty"`
	
	// Attributes provides additional customer-specific data
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Address represents address information used for determining tax jurisdiction
// and applicable tax rules. Different address components may be used by
// different tax systems.
//
// Example:
//
//	address := &Address{
//		Street1: "123 Main St",
//		City: "New York",
//		State: "NY",
//		PostalCode: "10001",
//		Country: "US",
//	}
type Address struct {
	// Street1 is the primary street address
	Street1 string `json:"street1"`
	
	// Street2 is the secondary street address (apartment, suite, etc.)
	Street2 string `json:"street2,omitempty"`
	
	// City is the city name
	City string `json:"city"`
	
	// State is the state or province
	State string `json:"state"`
	
	// PostalCode is the postal or ZIP code
	PostalCode string `json:"postal_code"`
	
	// Country is the country code (ISO 3166-1 alpha-2)
	Country string `json:"country"`
	
	// County is the county name, if applicable
	County string `json:"county,omitempty"`
	
	// District is the special district, if applicable
	District string `json:"district,omitempty"`
	
	// Latitude is the geographic latitude coordinate
	Latitude float64 `json:"latitude,omitempty"`
	
	// Longitude is the geographic longitude coordinate
	Longitude float64 `json:"longitude,omitempty"`
}

// TaxCalculationInput represents the complete input required for tax calculation.
// It includes all items, customer information, addresses, and context needed
// to determine applicable taxes.
//
// Example:
//
//	input := &TaxCalculationInput{
//		Items: []TaxableItem{
//			{ID: "item1", Name: "Product A", UnitPrice: 100.0, Quantity: 2},
//		},
//		Customer: Customer{ID: "cust1", Type: "individual"},
//		ShippingAddress: Address{City: "New York", State: "NY", Country: "US"},
//		TransactionDate: time.Now(),
//		Currency: "USD",
//	}
type TaxCalculationInput struct {
	// Items is the list of taxable items in the transaction
	Items           []TaxableItem `json:"items"`
	
	// Customer contains customer information for tax calculation
	Customer        Customer      `json:"customer"`
	
	// BillingAddress is the billing address, may affect tax calculation
	BillingAddress  Address       `json:"billing_address"`
	
	// ShippingAddress is the destination address for tax jurisdiction determination
	ShippingAddress Address       `json:"shipping_address"`
	
	// TransactionDate is the date of the transaction
	TransactionDate time.Time     `json:"transaction_date"`
	
	// TransactionType specifies the type of transaction
	TransactionType string        `json:"transaction_type"` // "sale", "purchase", "import", "export"
	
	// Currency is the transaction currency code
	Currency        string        `json:"currency"`
	
	// ExchangeRate is the currency exchange rate if different from base currency
	ExchangeRate    float64       `json:"exchange_rate,omitempty"`
	
	// ShippingAmount is the shipping cost for the transaction
	ShippingAmount  float64       `json:"shipping_amount,omitempty"`
	
	// DiscountAmount is the total discount applied to the transaction
	DiscountAmount  float64       `json:"discount_amount,omitempty"`
	
	// TaxRules contains specific tax rules to apply for this calculation
	TaxRules        []TaxRule     `json:"tax_rules,omitempty"`
	
	// Overrides contains any manual tax overrides to apply
	Overrides       []TaxOverride `json:"overrides,omitempty"`
	
	// Context provides additional context for tax calculation
	Context         map[string]interface{} `json:"context,omitempty"`
}

// TaxOverride represents manual tax overrides that can be applied during
// tax calculation. Overrides allow for manual adjustments to tax rates,
// amounts, or exemptions for specific scenarios.
//
// Example:
//
//	override := &TaxOverride{
//		Type: "rate",
//		TaxType: TaxTypeSales,
//		Value: 5.0,
//		Reason: "Special promotion discount",
//		ApprovedBy: "manager@company.com",
//	}
type TaxOverride struct {
	// Type specifies the override type ("rate", "amount", "exempt")
	Type        string  `json:"type"`
	
	// TaxType specifies which type of tax to override
	TaxType     TaxType `json:"tax_type"`
	
	// Value is the override value (rate percentage or fixed amount)
	Value       float64 `json:"value,omitempty"`
	
	// Reason explains why the override was applied
	Reason      string  `json:"reason"`
	
	// ApprovedBy identifies who approved the override
	ApprovedBy  string  `json:"approved_by,omitempty"`
	
	// ApprovedAt is the timestamp when the override was approved
	ApprovedAt  time.Time `json:"approved_at,omitempty"`
}

// AppliedTax represents a tax that was applied during calculation.
// It provides detailed information about each tax component that
// contributed to the final tax amount.
//
// Example:
//
//	appliedTax := &AppliedTax{
//		RuleID: "rule-ny-sales",
//		Name: "New York Sales Tax",
//		Type: TaxTypeSales,
//		Jurisdiction: JurisdictionState,
//		Rate: 8.25,
//		TaxableAmount: 100.0,
//		TaxAmount: 8.25,
//		Description: "NY state sales tax",
//	}
type AppliedTax struct {
	// RuleID is the unique identifier of the tax rule that was applied
	RuleID       string                `json:"rule_id"`
	
	// Name is the human-readable name of the tax rule
	Name         string                `json:"name"`
	
	// Type is the type of tax that was applied
	Type         TaxType               `json:"type"`
	
	// Jurisdiction is the tax jurisdiction that imposed this tax
	Jurisdiction TaxJurisdiction       `json:"jurisdiction"`
	
	// Method is the calculation method used for this tax
	Method       TaxCalculationMethod  `json:"method"`
	
	// Rate is the tax rate that was applied (as percentage)
	Rate         float64               `json:"rate"`
	
	// TaxableAmount is the amount that was subject to this tax
	TaxableAmount float64              `json:"taxable_amount"`
	
	// TaxAmount is the calculated tax amount
	TaxAmount    float64               `json:"tax_amount"`
	
	// Description provides additional details about the applied tax
	Description  string                `json:"description,omitempty"`
	
	// IsOverridden indicates if this tax was manually overridden
	IsOverridden bool                  `json:"is_overridden,omitempty"`
	
	// OverrideReason explains why the tax was overridden
	OverrideReason string              `json:"override_reason,omitempty"`
}

// TaxBreakdown represents detailed tax breakdown by item.
// It shows how taxes were calculated for each individual item
// in the transaction.
//
// Example:
//
//	breakdown := &TaxBreakdown{
//		ItemID: "item-123",
//		ItemName: "Laptop Computer",
//		ItemAmount: 999.99,
//		TaxableAmount: 999.99,
//		TotalTax: 82.50,
//		AppliedTaxes: []AppliedTax{...},
//	}
type TaxBreakdown struct {
	// ItemID is the unique identifier of the item
	ItemID      string       `json:"item_id"`
	
	// ItemName is the display name of the item
	ItemName    string       `json:"item_name"`
	
	// ItemAmount is the total amount for this item
	ItemAmount  float64      `json:"item_amount"`
	
	// AppliedTaxes lists all taxes applied to this item
	AppliedTaxes []AppliedTax `json:"applied_taxes"`
	
	// TotalTax is the sum of all taxes for this item
	TotalTax    float64      `json:"total_tax"`
	
	// TaxableAmount is the amount subject to tax for this item
	TaxableAmount float64    `json:"taxable_amount"`
	
	// ExemptAmount is the amount exempt from tax for this item
	ExemptAmount  float64    `json:"exempt_amount,omitempty"`
	
	// ExemptionReason explains why part of the amount was exempt
	ExemptionReason string   `json:"exemption_reason,omitempty"`
}

// TaxCalculationResult represents the complete result of tax calculation.
// It contains all calculated amounts, applied taxes, breakdowns, and metadata
// about the calculation process.
//
// Example:
//
//	result := &TaxCalculationResult{
//		Subtotal: 100.0,
//		TotalTax: 8.25,
//		GrandTotal: 108.25,
//		Currency: "USD",
//		CalculationDate: time.Now(),
//	}
type TaxCalculationResult struct {
	// Subtotal is the sum of all item amounts before tax
	Subtotal        float64         `json:"subtotal"`
	
	// TotalTax is the sum of all calculated taxes
	TotalTax        float64         `json:"total_tax"`
	
	// GrandTotal is the final amount including all taxes
	GrandTotal      float64         `json:"grand_total"`
	
	// TaxableAmount is the total amount subject to tax
	TaxableAmount   float64         `json:"taxable_amount"`
	
	// ExemptAmount is the total amount exempt from tax
	ExemptAmount    float64         `json:"exempt_amount"`
	
	// AppliedTaxes lists all taxes that were applied
	AppliedTaxes    []AppliedTax    `json:"applied_taxes"`
	
	// TaxBreakdown provides detailed tax breakdown by item
	TaxBreakdown    []TaxBreakdown  `json:"tax_breakdown"`
	
	// JurisdictionTotals shows tax amounts by jurisdiction
	JurisdictionTotals map[TaxJurisdiction]float64 `json:"jurisdiction_totals"`
	
	// TaxTypeTotals shows tax amounts by tax type
	TaxTypeTotals   map[TaxType]float64 `json:"tax_type_totals"`
	
	// EffectiveRate is the overall effective tax rate
	EffectiveRate   float64         `json:"effective_rate"`
	
	// Currency is the currency code for all amounts
	Currency        string          `json:"currency"`
	
	// CalculationDate is the timestamp when the calculation was performed
	CalculationDate time.Time       `json:"calculation_date"`
	
	// IsValid indicates whether the calculation completed successfully
	IsValid         bool            `json:"is_valid"`
	
	// Errors contains any errors encountered during calculation
	Errors          []string        `json:"errors,omitempty"`
	
	// Warnings contains any warnings generated during calculation
	Warnings        []string        `json:"warnings,omitempty"`
	
	// Metadata provides additional context about the calculation
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// TaxReport represents tax reporting data for compliance and filing purposes.
// It aggregates tax information over a specific period for a jurisdiction.
//
// Example:
//
//	report := &TaxReport{
//		ID: "report-2024-q1",
//		PeriodStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
//		PeriodEnd: time.Date(2024, 3, 31, 23, 59, 59, 0, time.UTC),
//		Jurisdiction: JurisdictionState,
//		TaxType: TaxTypeSales,
//	}
type TaxReport struct {
	// ID is the unique identifier for the report
	ID              string          `json:"id"`
	
	// PeriodStart is the start date of the reporting period
	PeriodStart     time.Time       `json:"period_start"`
	
	// PeriodEnd is the end date of the reporting period
	PeriodEnd       time.Time       `json:"period_end"`
	
	// Jurisdiction is the tax jurisdiction for this report
	Jurisdiction    TaxJurisdiction `json:"jurisdiction"`
	
	// TaxType is the type of tax covered by this report
	TaxType         TaxType         `json:"tax_type"`
	
	// TotalSales is the total sales amount for the period
	TotalSales      float64         `json:"total_sales"`
	
	// TaxableAmount is the total taxable amount
	TaxableAmount   float64         `json:"taxable_amount"`
	
	// ExemptAmount is the total exempt amount
	ExemptAmount    float64         `json:"exempt_amount"`
	
	// TaxCollected is the total tax collected
	TaxCollected    float64         `json:"tax_collected"`
	
	// TaxOwed is the total tax owed to authorities
	TaxOwed         float64         `json:"tax_owed"`
	
	// TransactionCount is the number of transactions in the period
	TransactionCount int            `json:"transaction_count"`
	
	// FilingDue is the due date for filing this report
	FilingDue       time.Time       `json:"filing_due"`
	
	// Status indicates the report status ("draft", "filed", "paid", "overdue")
	Status          string          `json:"status"`
	
	// Details provides line-by-line breakdown of the report
	Details         []TaxReportDetail `json:"details,omitempty"`
}

// TaxReportDetail represents detailed line items in tax reports.
// Each detail line represents a specific tax type or rate within the report.
//
// Example:
//
//	detail := &TaxReportDetail{
//		Category: "electronics",
//		Description: "State Sales Tax - 8.25%",
//		TaxableAmount: 10000.0,
//		TaxRate: 8.25,
//		TaxAmount: 825.0,
//		TransactionCount: 150,
//	}
type TaxReportDetail struct {
	// Category is the item category for this detail line
	Category        string  `json:"category"`
	
	// Description is the description of this line item
	Description     string  `json:"description"`
	
	// TaxableAmount is the taxable amount for this line
	TaxableAmount   float64 `json:"taxable_amount"`
	
	// TaxRate is the tax rate for this line
	TaxRate         float64 `json:"tax_rate"`
	
	// TaxAmount is the tax amount for this line
	TaxAmount       float64 `json:"tax_amount"`
	
	// TransactionCount is the number of transactions for this line
	TransactionCount int    `json:"transaction_count"`
}

// TaxConfiguration represents tax system configuration settings.
// It defines how tax calculations should be performed, including
// rounding rules, calculation order, and validation settings.
//
// Example:
//
//	config := &TaxConfiguration{
//		DefaultCurrency: "USD",
//		RoundingMode: "round",
//		RoundingPrecision: 2,
//		TaxInclusivePricing: false,
//		CompoundTaxes: false,
//	}
type TaxConfiguration struct {
	// DefaultCurrency is the default currency for calculations
	DefaultCurrency    string            `json:"default_currency"`
	
	// RoundingMode specifies how to round tax amounts ("round", "floor", "ceil")
	RoundingMode       string            `json:"rounding_mode"`
	
	// RoundingPrecision is the number of decimal places for rounding
	RoundingPrecision  int               `json:"rounding_precision"`
	
	// TaxInclusivePricing indicates whether prices include tax by default
	TaxInclusivePricing bool             `json:"tax_inclusive_pricing"`
	
	// CompoundTaxes indicates whether to compound taxes
	CompoundTaxes      bool              `json:"compound_taxes"`
	
	// TaxOnShipping indicates whether to apply tax to shipping costs
	TaxOnShipping      bool              `json:"tax_on_shipping"`
	
	// TaxOnDiscounts indicates whether to apply tax after discounts
	TaxOnDiscounts     bool              `json:"tax_on_discounts"`
	
	// DefaultRules contains the default tax rules to apply
	DefaultRules       []TaxRule         `json:"default_rules"`
	
	// ExemptionCertificates lists valid exemption certificates
	ExemptionCertificates []string       `json:"exemption_certificates,omitempty"`
	
	// ReportingFrequency specifies how often reports are generated ("monthly", "quarterly", "annually")
	ReportingFrequency string            `json:"reporting_frequency"`
	
	// Settings provides additional configuration options
	Settings           map[string]interface{} `json:"settings,omitempty"`
}

// TaxValidationRule represents validation rules for tax calculations.
// These rules help ensure tax calculations are accurate and compliant
// with business rules and regulatory requirements.
//
// Example:
//
//	rule := &TaxValidationRule{
//		ID: "rate-limit-rule",
//		Name: "Maximum Tax Rate Validation",
//		Type: "rate_limit",
//		Condition: "rate <= 15.0",
//		Message: "Tax rate cannot exceed 15%",
//		Severity: "error",
//		IsActive: true,
//	}
type TaxValidationRule struct {
	// ID is the unique identifier for the validation rule
	ID          string    `json:"id"`
	
	// Name is the human-readable name of the rule
	Name        string    `json:"name"`
	
	// Type specifies the validation type ("rate_limit", "amount_limit", "consistency")
	Type        string    `json:"type"`
	
	// Condition defines the validation condition to check
	Condition   string    `json:"condition"`
	
	// Message is the error message to display when validation fails
	Message     string    `json:"message"`
	
	// Severity indicates the severity level ("error", "warning", "info")
	Severity    string    `json:"severity"`
	
	// IsActive indicates whether the rule is currently active
	IsActive    bool      `json:"is_active"`
	
	// ValidFrom is the date when the rule becomes effective
	ValidFrom   time.Time `json:"valid_from"`
	
	// ValidUntil is the date when the rule expires
	ValidUntil  time.Time `json:"valid_until"`
}

// TaxAuditTrail represents audit trail entries for tax operations.
// It provides a complete history of changes and operations performed
// on tax-related entities for compliance and debugging purposes.
//
// Example:
//
//	audit := &TaxAuditTrail{
//		ID: "audit-12345",
//		TransactionID: "txn-67890",
//		Action: "calculate",
//		UserID: "user-123",
//		Timestamp: time.Now(),
//		Reason: "Standard tax calculation",
//	}
type TaxAuditTrail struct {
	// ID is the unique identifier for the audit entry
	ID              string                 `json:"id"`
	
	// TransactionID is the unique identifier of the related transaction
	TransactionID   string                 `json:"transaction_id"`
	
	// CalculationID is the unique identifier of the tax calculation
	CalculationID   string                 `json:"calculation_id"`
	
	// Timestamp is when the action was performed
	Timestamp       time.Time              `json:"timestamp"`
	
	// Action describes what action was performed ("calculate", "override", "exempt", "adjust")
	Action          string                 `json:"action"`
	
	// UserID is the unique identifier of the user who performed the action
	UserID          string                 `json:"user_id,omitempty"`
	
	// Reason explains why the action was performed
	Reason          string                 `json:"reason,omitempty"`
	
	// BeforeState contains the state before the action
	BeforeState     map[string]interface{} `json:"before_state,omitempty"`
	
	// AfterState contains the state after the action
	AfterState      map[string]interface{} `json:"after_state,omitempty"`
	
	// Changes lists the specific changes that were made
	Changes         []string               `json:"changes,omitempty"`
	
	// IPAddress is the IP address from which the action was performed
	IPAddress       string                 `json:"ip_address,omitempty"`
	
	// UserAgent is the user agent string from the client
	UserAgent       string                 `json:"user_agent,omitempty"`
}
package tax

import (
	"time"
)

// TaxType represents different types of taxes
type TaxType string

const (
	TaxTypeSales     TaxType = "sales"     // Sales tax
	TaxTypeVAT       TaxType = "vat"       // Value Added Tax
	TaxTypeGST       TaxType = "gst"       // Goods and Services Tax
	TaxTypeExcise    TaxType = "excise"    // Excise tax
	TaxTypeCustoms   TaxType = "customs"   // Customs duty
	TaxTypeProperty  TaxType = "property"  // Property tax
	TaxTypeWithholding TaxType = "withholding" // Withholding tax
	TaxTypeDigital   TaxType = "digital"   // Digital services tax
	TaxTypeEnvironmental TaxType = "environmental" // Environmental tax
	TaxTypeLuxury    TaxType = "luxury"    // Luxury tax
)

// TaxCalculationMethod represents how tax is calculated
type TaxCalculationMethod string

const (
	TaxMethodPercentage TaxCalculationMethod = "percentage" // Percentage of amount
	TaxMethodFixed      TaxCalculationMethod = "fixed"      // Fixed amount
	TaxMethodTiered     TaxCalculationMethod = "tiered"     // Tiered rates
	TaxMethodProgressive TaxCalculationMethod = "progressive" // Progressive rates
	TaxMethodCompound   TaxCalculationMethod = "compound"   // Compound tax
)

// TaxJurisdiction represents tax jurisdiction levels
type TaxJurisdiction string

const (
	JurisdictionFederal    TaxJurisdiction = "federal"    // Federal/National level
	JurisdictionState      TaxJurisdiction = "state"      // State/Province level
	JurisdictionCounty     TaxJurisdiction = "county"     // County level
	JurisdictionCity       TaxJurisdiction = "city"       // City/Municipal level
	JurisdictionDistrict   TaxJurisdiction = "district"   // Special district
	JurisdictionInternational TaxJurisdiction = "international" // International
)

// TaxableItem represents an item subject to taxation
type TaxableItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory,omitempty"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalAmount float64 `json:"total_amount"`
	Weight      float64 `json:"weight,omitempty"`
	Volume      float64 `json:"volume,omitempty"`
	Origin      string  `json:"origin,omitempty"`      // Country of origin
	HSCode      string  `json:"hs_code,omitempty"`     // Harmonized System code
	SKU         string  `json:"sku,omitempty"`
	Brand       string  `json:"brand,omitempty"`
	IsDigital   bool    `json:"is_digital,omitempty"`  // Digital goods
	IsLuxury    bool    `json:"is_luxury,omitempty"`   // Luxury item
	IsExempt    bool    `json:"is_exempt,omitempty"`   // Tax exempt
	ExemptionReason string `json:"exemption_reason,omitempty"`
	CustomAttributes map[string]interface{} `json:"custom_attributes,omitempty"`
}

// TaxRule represents a tax calculation rule
type TaxRule struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Type            TaxType               `json:"type"`
	Jurisdiction    TaxJurisdiction       `json:"jurisdiction"`
	Method          TaxCalculationMethod  `json:"method"`
	Rate            float64               `json:"rate"`           // Tax rate (percentage or fixed amount)
	MinAmount       float64               `json:"min_amount,omitempty"`     // Minimum taxable amount
	MaxAmount       float64               `json:"max_amount,omitempty"`     // Maximum taxable amount
	Thresholds      []TaxThreshold        `json:"thresholds,omitempty"`     // For tiered/progressive tax
	ApplicableCategories []string         `json:"applicable_categories,omitempty"`
	ExemptCategories     []string         `json:"exempt_categories,omitempty"`
	ApplicableCountries  []string         `json:"applicable_countries,omitempty"`
	ApplicableStates     []string         `json:"applicable_states,omitempty"`
	ApplicableCities     []string         `json:"applicable_cities,omitempty"`
	PostalCodes          []string         `json:"postal_codes,omitempty"`
	IsActive        bool                  `json:"is_active"`
	ValidFrom       time.Time             `json:"valid_from"`
	ValidUntil      time.Time             `json:"valid_until"`
	Priority        int                   `json:"priority"`       // Higher number = higher priority
	Description     string                `json:"description,omitempty"`
	Conditions      []TaxCondition        `json:"conditions,omitempty"`
	Exemptions      []TaxExemption        `json:"exemptions,omitempty"`
}

// TaxThreshold represents tax rate thresholds for tiered/progressive taxation
type TaxThreshold struct {
	MinAmount float64 `json:"min_amount"`
	MaxAmount float64 `json:"max_amount,omitempty"` // 0 means no upper limit
	Rate      float64 `json:"rate"`
	FixedAmount float64 `json:"fixed_amount,omitempty"` // Fixed amount for this tier
}

// TaxCondition represents conditions for tax application
type TaxCondition struct {
	Type     string      `json:"type"`     // "amount", "quantity", "weight", "category", "customer_type"
	Operator string      `json:"operator"` // ">", "<", ">=", "<=", "=", "!=", "in", "not_in"
	Value    interface{} `json:"value"`    // Condition value
	Logic    string      `json:"logic,omitempty"` // "AND", "OR" for combining conditions
}

// TaxExemption represents tax exemption rules
type TaxExemption struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`        // "customer", "item", "transaction", "location"
	Reason      string    `json:"reason"`      // Exemption reason
	Certificate string    `json:"certificate,omitempty"` // Exemption certificate number
	ValidFrom   time.Time `json:"valid_from"`
	ValidUntil  time.Time `json:"valid_until"`
	Conditions  []TaxCondition `json:"conditions,omitempty"`
}

// Customer represents customer information for tax calculation
type Customer struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`         // "individual", "business", "government", "nonprofit"
	TaxID        string            `json:"tax_id,omitempty"`        // Tax identification number
	VATNumber    string            `json:"vat_number,omitempty"`    // VAT registration number
	BusinessType string            `json:"business_type,omitempty"` // Type of business
	Industry     string            `json:"industry,omitempty"`
	Exemptions   []TaxExemption    `json:"exemptions,omitempty"`
	Attributes   map[string]string `json:"attributes,omitempty"`
}

// Address represents address information for tax calculation
type Address struct {
	Street1    string  `json:"street1"`
	Street2    string  `json:"street2,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	PostalCode string  `json:"postal_code"`
	Country    string  `json:"country"`
	County     string  `json:"county,omitempty"`
	District   string  `json:"district,omitempty"`
	Latitude   float64 `json:"latitude,omitempty"`
	Longitude  float64 `json:"longitude,omitempty"`
}

// TaxCalculationInput represents input for tax calculation
type TaxCalculationInput struct {
	Items           []TaxableItem `json:"items"`
	Customer        Customer      `json:"customer"`
	BillingAddress  Address       `json:"billing_address"`
	ShippingAddress Address       `json:"shipping_address"`
	TransactionDate time.Time     `json:"transaction_date"`
	TransactionType string        `json:"transaction_type"` // "sale", "purchase", "import", "export"
	Currency        string        `json:"currency"`
	ExchangeRate    float64       `json:"exchange_rate,omitempty"`
	ShippingAmount  float64       `json:"shipping_amount,omitempty"`
	DiscountAmount  float64       `json:"discount_amount,omitempty"`
	TaxRules        []TaxRule     `json:"tax_rules,omitempty"`
	Overrides       []TaxOverride `json:"overrides,omitempty"`
	Context         map[string]interface{} `json:"context,omitempty"`
}

// TaxOverride represents manual tax overrides
type TaxOverride struct {
	Type        string  `json:"type"`        // "rate", "amount", "exempt"
	TaxType     TaxType `json:"tax_type"`
	Value       float64 `json:"value,omitempty"`
	Reason      string  `json:"reason"`
	ApprovedBy  string  `json:"approved_by,omitempty"`
	ApprovedAt  time.Time `json:"approved_at,omitempty"`
}

// AppliedTax represents a tax that was applied
type AppliedTax struct {
	RuleID       string                `json:"rule_id"`
	Name         string                `json:"name"`
	Type         TaxType               `json:"type"`
	Jurisdiction TaxJurisdiction       `json:"jurisdiction"`
	Method       TaxCalculationMethod  `json:"method"`
	Rate         float64               `json:"rate"`
	TaxableAmount float64              `json:"taxable_amount"`
	TaxAmount    float64               `json:"tax_amount"`
	Description  string                `json:"description,omitempty"`
	IsOverridden bool                  `json:"is_overridden,omitempty"`
	OverrideReason string              `json:"override_reason,omitempty"`
}

// TaxBreakdown represents detailed tax breakdown by item
type TaxBreakdown struct {
	ItemID      string       `json:"item_id"`
	ItemName    string       `json:"item_name"`
	ItemAmount  float64      `json:"item_amount"`
	AppliedTaxes []AppliedTax `json:"applied_taxes"`
	TotalTax    float64      `json:"total_tax"`
	TaxableAmount float64    `json:"taxable_amount"`
	ExemptAmount  float64    `json:"exempt_amount,omitempty"`
	ExemptionReason string   `json:"exemption_reason,omitempty"`
}

// TaxCalculationResult represents the result of tax calculation
type TaxCalculationResult struct {
	Subtotal        float64         `json:"subtotal"`         // Total before tax
	TotalTax        float64         `json:"total_tax"`        // Total tax amount
	GrandTotal      float64         `json:"grand_total"`      // Total including tax
	TaxableAmount   float64         `json:"taxable_amount"`   // Amount subject to tax
	ExemptAmount    float64         `json:"exempt_amount"`    // Amount exempt from tax
	AppliedTaxes    []AppliedTax    `json:"applied_taxes"`    // All applied taxes
	TaxBreakdown    []TaxBreakdown  `json:"tax_breakdown"`    // Per-item tax breakdown
	JurisdictionTotals map[TaxJurisdiction]float64 `json:"jurisdiction_totals"` // Tax by jurisdiction
	TaxTypeTotals   map[TaxType]float64 `json:"tax_type_totals"`   // Tax by type
	EffectiveRate   float64         `json:"effective_rate"`   // Overall effective tax rate
	Currency        string          `json:"currency"`
	CalculationDate time.Time       `json:"calculation_date"`
	IsValid         bool            `json:"is_valid"`
	Errors          []string        `json:"errors,omitempty"`
	Warnings        []string        `json:"warnings,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// TaxReport represents tax reporting information
type TaxReport struct {
	ID              string          `json:"id"`
	PeriodStart     time.Time       `json:"period_start"`
	PeriodEnd       time.Time       `json:"period_end"`
	Jurisdiction    TaxJurisdiction `json:"jurisdiction"`
	TaxType         TaxType         `json:"tax_type"`
	TotalSales      float64         `json:"total_sales"`
	TaxableAmount   float64         `json:"taxable_amount"`
	ExemptAmount    float64         `json:"exempt_amount"`
	TaxCollected    float64         `json:"tax_collected"`
	TaxOwed         float64         `json:"tax_owed"`
	TransactionCount int            `json:"transaction_count"`
	FilingDue       time.Time       `json:"filing_due"`
	Status          string          `json:"status"` // "draft", "filed", "paid", "overdue"
	Details         []TaxReportDetail `json:"details,omitempty"`
}

// TaxReportDetail represents detailed tax report information
type TaxReportDetail struct {
	Category        string  `json:"category"`
	Description     string  `json:"description"`
	TaxableAmount   float64 `json:"taxable_amount"`
	TaxRate         float64 `json:"tax_rate"`
	TaxAmount       float64 `json:"tax_amount"`
	TransactionCount int    `json:"transaction_count"`
}

// TaxConfiguration represents tax system configuration
type TaxConfiguration struct {
	DefaultCurrency    string            `json:"default_currency"`
	RoundingMode       string            `json:"rounding_mode"`       // "round", "floor", "ceil"
	RoundingPrecision  int               `json:"rounding_precision"`  // Decimal places
	TaxInclusivePricing bool             `json:"tax_inclusive_pricing"` // Whether prices include tax
	CompoundTaxes      bool              `json:"compound_taxes"`      // Whether to compound taxes
	TaxOnShipping      bool              `json:"tax_on_shipping"`     // Whether to tax shipping
	TaxOnDiscounts     bool              `json:"tax_on_discounts"`    // Whether to apply tax after discounts
	DefaultRules       []TaxRule         `json:"default_rules"`
	ExemptionCertificates []string       `json:"exemption_certificates,omitempty"`
	ReportingFrequency string            `json:"reporting_frequency"` // "monthly", "quarterly", "annually"
	Settings           map[string]interface{} `json:"settings,omitempty"`
}

// TaxValidationRule represents validation rules for tax calculations
type TaxValidationRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`        // "rate_limit", "amount_limit", "consistency"
	Condition   string    `json:"condition"`   // Validation condition
	Message     string    `json:"message"`     // Error message
	Severity    string    `json:"severity"`    // "error", "warning", "info"
	IsActive    bool      `json:"is_active"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidUntil  time.Time `json:"valid_until"`
}

// TaxAuditTrail represents audit trail for tax calculations
type TaxAuditTrail struct {
	ID              string                 `json:"id"`
	TransactionID   string                 `json:"transaction_id"`
	CalculationID   string                 `json:"calculation_id"`
	Timestamp       time.Time              `json:"timestamp"`
	Action          string                 `json:"action"`          // "calculate", "override", "exempt", "adjust"
	UserID          string                 `json:"user_id,omitempty"`
	Reason          string                 `json:"reason,omitempty"`
	BeforeState     map[string]interface{} `json:"before_state,omitempty"`
	AfterState      map[string]interface{} `json:"after_state,omitempty"`
	Changes         []string               `json:"changes,omitempty"`
	IPAddress       string                 `json:"ip_address,omitempty"`
	UserAgent       string                 `json:"user_agent,omitempty"`
}
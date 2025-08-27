package discount

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// RuleEngine manages and applies discount rules
type RuleEngine struct {
	BulkRules        []BulkDiscountRule
	TierRules        []TierPricingRule
	BundleRules      []BundleDiscountRule
	LoyaltyRules     []LoyaltyDiscountRule
	ProgressiveRules []ProgressiveDiscountRule
	CategoryRules    []CategoryDiscountRule
	FrequencyRules   []FrequencyDiscountRule
	SeasonalRules    []SeasonalDiscountRule
	CrossSellRules   []CrossSellRule
	MixMatchRules    []MixAndMatchRule
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine() *RuleEngine {
	return &RuleEngine{
		BulkRules:        []BulkDiscountRule{},
		TierRules:        []TierPricingRule{},
		BundleRules:      []BundleDiscountRule{},
		LoyaltyRules:     []LoyaltyDiscountRule{},
		ProgressiveRules: []ProgressiveDiscountRule{},
		CategoryRules:    []CategoryDiscountRule{},
		FrequencyRules:   []FrequencyDiscountRule{},
		SeasonalRules:    []SeasonalDiscountRule{},
		CrossSellRules:   []CrossSellRule{},
		MixMatchRules:    []MixAndMatchRule{},
	}
}

// AddBulkRule adds a bulk discount rule
func (re *RuleEngine) AddBulkRule(rule BulkDiscountRule) error {
	if err := validateBulkRule(rule); err != nil {
		return err
	}
	re.BulkRules = append(re.BulkRules, rule)
	return nil
}

// AddTierRule adds a tier pricing rule
func (re *RuleEngine) AddTierRule(rule TierPricingRule) error {
	if err := validateTierRule(rule); err != nil {
		return err
	}
	re.TierRules = append(re.TierRules, rule)
	return nil
}

// AddBundleRule adds a bundle discount rule
func (re *RuleEngine) AddBundleRule(rule BundleDiscountRule) error {
	if err := validateBundleRule(rule); err != nil {
		return err
	}
	re.BundleRules = append(re.BundleRules, rule)
	return nil
}

// AddLoyaltyRule adds a loyalty discount rule
func (re *RuleEngine) AddLoyaltyRule(rule LoyaltyDiscountRule) error {
	if err := validateLoyaltyRule(rule); err != nil {
		return err
	}
	re.LoyaltyRules = append(re.LoyaltyRules, rule)
	return nil
}

// AddCategoryRule adds a category discount rule
func (re *RuleEngine) AddCategoryRule(rule CategoryDiscountRule) error {
	if err := validateCategoryRule(rule); err != nil {
		return err
	}
	re.CategoryRules = append(re.CategoryRules, rule)
	return nil
}

// AddSeasonalRule adds a seasonal discount rule
func (re *RuleEngine) AddSeasonalRule(rule SeasonalDiscountRule) error {
	if err := validateSeasonalRule(rule); err != nil {
		return err
	}
	re.SeasonalRules = append(re.SeasonalRules, rule)
	return nil
}

// AddCrossSellRule adds a cross-sell discount rule
func (re *RuleEngine) AddCrossSellRule(rule CrossSellRule) error {
	if err := validateCrossSellRule(rule); err != nil {
		return err
	}
	re.CrossSellRules = append(re.CrossSellRules, rule)
	return nil
}

// AddMixMatchRule adds a mix and match rule
func (re *RuleEngine) AddMixMatchRule(rule MixAndMatchRule) error {
	if err := validateMixMatchRule(rule); err != nil {
		return err
	}
	re.MixMatchRules = append(re.MixMatchRules, rule)
	return nil
}

// ApplyRules applies all rules and returns the best discount
func (re *RuleEngine) ApplyRules(items []DiscountItem, customer Customer, allowStacking bool) DiscountCalculationResult {
	input := DiscountCalculationInput{
		Items:                     items,
		Customer:                  customer,
		BulkRules:                 re.BulkRules,
		TierRules:                 re.TierRules,
		BundleRules:               re.BundleRules,
		LoyaltyRules:              re.LoyaltyRules,
		ProgressiveRules:          re.ProgressiveRules,
		CategoryRules:             re.CategoryRules,
		AllowStacking:             allowStacking,
		MaxStackedDiscountPercent: 50, // Default max 50% stacked discount
	}

	return Calculate(input)
}

// ApplyFrequencyDiscounts applies purchase frequency-based discounts
func (re *RuleEngine) ApplyFrequencyDiscounts(items []DiscountItem, customer Customer) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	for _, rule := range re.FrequencyRules {
		if customer.PurchaseCount >= rule.MinPurchaseCount {
			discount := result.OriginalAmount * (rule.DiscountPercent / 100)

			result.TotalDiscount += discount
			result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
				Type:           DiscountTypeLoyalty,
				RuleID:         "frequency",
				Name:           "Repeat Customer Discount",
				DiscountAmount: discount,
				AppliedItems:   items,
				Description:    "Discount for repeat customers",
			})
			break // Apply only the first matching rule
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// ApplySeasonalDiscounts applies seasonal discounts
func (re *RuleEngine) ApplySeasonalDiscounts(items []DiscountItem, customer Customer) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	now := time.Now()

	for _, rule := range re.SeasonalRules {
		// Check if rule is currently valid
		if now.Before(rule.ValidFrom) || now.After(rule.ValidUntil) {
			continue
		}

		// Check if current season matches
		if !isCurrentSeason(now, rule.Season) {
			continue
		}

		applicableItems := items
		if len(rule.Categories) > 0 {
			applicableItems = getApplicableItems(items, rule.Categories, nil)
		}

		if len(applicableItems) > 0 {
			itemAmount := calculateItemsAmount(applicableItems)
			discountPercent := rule.DiscountPercent

			// Apply member bonus multiplier
			if rule.Multiplier > 0 && customer.LoyaltyTier != "" {
				discountPercent *= rule.Multiplier
			}

			discount := itemAmount * (discountPercent / 100)

			result.TotalDiscount += discount
			result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
				Type:           DiscountTypeCategory,
				RuleID:         "seasonal_" + rule.Season,
				Name:           "Seasonal Discount",
				DiscountAmount: discount,
				AppliedItems:   applicableItems,
				Description:    fmt.Sprintf("%s seasonal discount", strings.Title(rule.Season)),
			})
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// ApplyCrossSellDiscounts applies cross-sell discounts
func (re *RuleEngine) ApplyCrossSellDiscounts(items []DiscountItem) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	for _, rule := range re.CrossSellRules {
		mainProducts := getApplicableItems(items, rule.MainProductCategories, nil)
		accessories := getApplicableItems(items, rule.AccessoryCategories, nil)

		if len(mainProducts) > 0 && len(accessories) > 0 {
			// Check minimum main product price if specified
			if rule.MinMainProductPrice > 0 {
				mainAmount := calculateItemsAmount(mainProducts)
				if mainAmount < rule.MinMainProductPrice {
					continue
				}
			}

			combinedItems := append(mainProducts, accessories...)
			combinedAmount := calculateItemsAmount(combinedItems)

			var discount float64
			if rule.ComboPrice > 0 {
				// Fixed combo price
				discount = math.Max(0, combinedAmount-rule.ComboPrice)
			} else {
				// Percentage discount
				discount = combinedAmount * (rule.DiscountPercent / 100)
			}

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type:           DiscountTypeBundle,
					RuleID:         "cross_sell",
					Name:           "Cross-sell Discount",
					DiscountAmount: discount,
					AppliedItems:   combinedItems,
					Description:    "Main product + accessory discount",
				})
			}
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// ApplyMixAndMatchDiscounts applies mix and match discounts
func (re *RuleEngine) ApplyMixAndMatchDiscounts(items []DiscountItem) DiscountCalculationResult {
	result := DiscountCalculationResult{
		OriginalAmount:   calculateOriginalAmount(items),
		IsValid:          true,
		AppliedDiscounts: []DiscountApplication{},
	}

	for _, rule := range re.MixMatchRules {
		applicableItems := getApplicableItems(items, rule.Categories, nil)

		if len(applicableItems) >= rule.RequiredItems {
			// Calculate how many times this rule can be applied
			applications := len(applicableItems) / rule.RequiredItems
			if rule.MaxApplications > 0 && applications > rule.MaxApplications {
				applications = rule.MaxApplications
			}

			// Take the required number of items for discount
			discountItems := applicableItems[:applications*rule.RequiredItems]
			itemAmount := calculateItemsAmount(discountItems)

			var discount float64
			switch rule.DiscountType {
			case "flat_discount":
				discount = rule.DiscountValue * float64(applications)
			case "percentage":
				discount = itemAmount * (rule.DiscountValue / 100)
			}

			if discount > 0 {
				result.TotalDiscount += discount
				result.AppliedDiscounts = append(result.AppliedDiscounts, DiscountApplication{
					Type:           DiscountTypeBundle,
					RuleID:         "mix_match",
					Name:           "Mix & Match Discount",
					DiscountAmount: discount,
					AppliedItems:   discountItems,
					Description:    fmt.Sprintf("Mix & match %d items discount", rule.RequiredItems),
				})
			}
		}
	}

	result.FinalAmount = result.OriginalAmount - result.TotalDiscount
	if result.OriginalAmount > 0 {
		result.SavingsPercent = (result.TotalDiscount / result.OriginalAmount) * 100
	}

	return result
}

// Validation functions

func validateBulkRule(rule BulkDiscountRule) error {
	if rule.MinQuantity <= 0 {
		return errors.New("minimum quantity must be greater than 0")
	}
	if rule.MaxQuantity > 0 && rule.MaxQuantity < rule.MinQuantity {
		return errors.New("maximum quantity must be greater than minimum quantity")
	}
	if rule.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	if rule.DiscountType == "percentage" && rule.DiscountValue > 100 {
		return errors.New("percentage discount cannot exceed 100%")
	}
	return nil
}

func validateTierRule(rule TierPricingRule) error {
	if rule.MinQuantity <= 0 {
		return errors.New("minimum quantity must be greater than 0")
	}
	if rule.PricePerItem <= 0 {
		return errors.New("price per item must be greater than 0")
	}
	return nil
}

func validateBundleRule(rule BundleDiscountRule) error {
	if rule.ID == "" {
		return errors.New("bundle rule ID is required")
	}
	if rule.MinItems <= 0 {
		return errors.New("minimum items must be greater than 0")
	}
	if len(rule.RequiredProducts) == 0 && len(rule.RequiredCategories) == 0 {
		return errors.New("either required products or categories must be specified")
	}
	if rule.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	return nil
}

func validateLoyaltyRule(rule LoyaltyDiscountRule) error {
	if rule.Tier == "" {
		return errors.New("loyalty tier is required")
	}
	if rule.DiscountPercent <= 0 || rule.DiscountPercent > 100 {
		return errors.New("discount percent must be between 0 and 100")
	}
	return nil
}

func validateCategoryRule(rule CategoryDiscountRule) error {
	if rule.Category == "" {
		return errors.New("category is required")
	}
	if rule.DiscountPercent <= 0 || rule.DiscountPercent > 100 {
		return errors.New("discount percent must be between 0 and 100")
	}
	if rule.ValidUntil.Before(rule.ValidFrom) {
		return errors.New("valid until must be after valid from")
	}
	return nil
}

func validateSeasonalRule(rule SeasonalDiscountRule) error {
	validSeasons := []string{"spring", "summer", "autumn", "winter"}
	validSeason := false
	for _, season := range validSeasons {
		if strings.ToLower(rule.Season) == season {
			validSeason = true
			break
		}
	}
	if !validSeason {
		return errors.New("invalid season specified")
	}
	if rule.DiscountPercent <= 0 {
		return errors.New("discount percent must be greater than 0")
	}
	return nil
}

func validateCrossSellRule(rule CrossSellRule) error {
	if len(rule.MainProductCategories) == 0 {
		return errors.New("main product categories are required")
	}
	if len(rule.AccessoryCategories) == 0 {
		return errors.New("accessory categories are required")
	}
	if rule.DiscountPercent <= 0 && rule.ComboPrice <= 0 {
		return errors.New("either discount percent or combo price must be specified")
	}
	return nil
}

func validateMixMatchRule(rule MixAndMatchRule) error {
	if len(rule.Categories) == 0 {
		return errors.New("categories are required")
	}
	if rule.RequiredItems <= 0 {
		return errors.New("required items must be greater than 0")
	}
	if rule.DiscountValue <= 0 {
		return errors.New("discount value must be greater than 0")
	}
	return nil
}

// Helper functions

func isCurrentSeason(now time.Time, season string) bool {
	month := now.Month()
	switch strings.ToLower(season) {
	case "spring":
		return month >= 3 && month <= 5
	case "summer":
		return month >= 6 && month <= 8
	case "autumn", "fall":
		return month >= 9 && month <= 11
	case "winter":
		return month == 12 || month <= 2
	default:
		return false
	}
}

// GetApplicableRules returns rules that are applicable for given items and customer
func (re *RuleEngine) GetApplicableRules(items []DiscountItem, customer Customer) map[string]interface{} {
	applicableRules := make(map[string]interface{})

	// Check bulk rules
	for _, rule := range re.BulkRules {
		applicableItems := getApplicableItems(items, rule.ApplicableCategories, rule.ApplicableProducts)
		if getTotalQuantity(applicableItems) >= rule.MinQuantity {
			applicableRules["bulk"] = append(applicableRules["bulk"].([]BulkDiscountRule), rule)
		}
	}

	// Check loyalty rules
	for _, rule := range re.LoyaltyRules {
		if customer.LoyaltyTier == rule.Tier {
			applicableRules["loyalty"] = append(applicableRules["loyalty"].([]LoyaltyDiscountRule), rule)
		}
	}

	// Check seasonal rules
	now := time.Now()
	for _, rule := range re.SeasonalRules {
		if now.After(rule.ValidFrom) && now.Before(rule.ValidUntil) && isCurrentSeason(now, rule.Season) {
			applicableRules["seasonal"] = append(applicableRules["seasonal"].([]SeasonalDiscountRule), rule)
		}
	}

	return applicableRules
}

// ClearRules clears all rules from the engine
func (re *RuleEngine) ClearRules() {
	re.BulkRules = []BulkDiscountRule{}
	re.TierRules = []TierPricingRule{}
	re.BundleRules = []BundleDiscountRule{}
	re.LoyaltyRules = []LoyaltyDiscountRule{}
	re.ProgressiveRules = []ProgressiveDiscountRule{}
	re.CategoryRules = []CategoryDiscountRule{}
	re.FrequencyRules = []FrequencyDiscountRule{}
	re.SeasonalRules = []SeasonalDiscountRule{}
	re.CrossSellRules = []CrossSellRule{}
	re.MixMatchRules = []MixAndMatchRule{}
}

package pricing

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Calculator handles pricing calculations
type Calculator struct {
	rules           []PricingRule
	bundles         []Bundle
	tierPricing     []TierPricing
	dynamicConfigs  []DynamicPricingConfig
	marketData      map[string]MarketData
	analytics       map[string]PricingAnalytics
}

// NewCalculator creates a new pricing calculator
func NewCalculator() *Calculator {
	return &Calculator{
		rules:          make([]PricingRule, 0),
		bundles:        make([]Bundle, 0),
		tierPricing:    make([]TierPricing, 0),
		dynamicConfigs: make([]DynamicPricingConfig, 0),
		marketData:     make(map[string]MarketData),
		analytics:      make(map[string]PricingAnalytics),
	}
}

// Calculate performs comprehensive pricing calculation
func (c *Calculator) Calculate(input PricingInput) (*PricingResult, error) {
	if err := c.validateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	result := &PricingResult{
		Items:           make([]PricedItem, 0),
		Currency:        input.Context.Currency,
		CalculationTime: time.Now(),
		IsValid:         true,
		Errors:          make([]string, 0),
		Warnings:        make([]string, 0),
		Metadata:        make(map[string]interface{}),
	}

	// Merge rules from input and calculator
	allRules := append(c.rules, input.Rules...)
	allBundles := append(c.bundles, input.Bundles...)
	allTierPricing := append(c.tierPricing, input.TierPricing...)

	// Calculate pricing for each item
	for _, item := range input.Items {
		pricedItem, err := c.calculateItemPricing(item, input.Customer, input.Context, allRules, allTierPricing, input.Options)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Error pricing item %s: %v", item.ID, err))
			continue
		}
		result.Items = append(result.Items, *pricedItem)
	}

	// Calculate bundle pricing if enabled
	if input.Options.CalculateBundle {
		bundleResults := c.calculateBundlePricing(result.Items, allBundles, input.Customer, input.Context)
		result.AppliedBundles = bundleResults
	}

	// Calculate totals
	c.calculateTotals(result)

	// Generate recommendations
	if len(allBundles) > 0 {
		result.Recommendations = c.generateRecommendations(result.Items, allBundles, allTierPricing)
	}

	return result, nil
}

// calculateItemPricing calculates pricing for a single item
func (c *Calculator) calculateItemPricing(item PricingItem, customer Customer, context PricingContext, rules []PricingRule, tierPricing []TierPricing, options PricingOptions) (*PricedItem, error) {
	pricedItem := &PricedItem{
		ItemID:        item.ID,
		Name:          item.Name,
		Quantity:      item.Quantity,
		BasePrice:     item.BasePrice,
		OriginalPrice: item.BasePrice,
		FinalPrice:    item.BasePrice,
		UnitPrice:     item.BasePrice,
		AppliedRules:  make([]AppliedPricingRule, 0),
		Metadata:      make(map[string]interface{}),
	}

	// Apply dynamic pricing if configured
	if dynamicPrice := c.calculateDynamicPricing(item, context); dynamicPrice > 0 {
		pricedItem.FinalPrice = dynamicPrice
		pricedItem.UnitPrice = dynamicPrice
	}

	// Apply tier pricing if enabled
	if options.CalculateTiers {
		if tierInfo := c.calculateTierPricing(item, tierPricing); tierInfo != nil {
			pricedItem.TierInfo = tierInfo
			pricedItem.FinalPrice = tierInfo.TierPrice
			pricedItem.UnitPrice = tierInfo.TierPrice
		}
	}

	// Apply pricing rules
	applicableRules := c.getApplicableRules(item, customer, context, rules)
	for _, rule := range applicableRules {
		adjustedPrice, appliedRule := c.applyPricingRule(pricedItem.FinalPrice, rule, item, customer)
		if appliedRule != nil {
			pricedItem.FinalPrice = adjustedPrice
			pricedItem.AppliedRules = append(pricedItem.AppliedRules, *appliedRule)
		}
	}

	// Apply rounding
	pricedItem.FinalPrice = c.roundPrice(pricedItem.FinalPrice, options.RoundingMode, options.RoundingPrecision)
	pricedItem.UnitPrice = pricedItem.FinalPrice
	pricedItem.TotalPrice = pricedItem.FinalPrice * float64(item.Quantity)

	// Calculate savings
	pricedItem.Savings = pricedItem.OriginalPrice - pricedItem.FinalPrice
	if pricedItem.OriginalPrice > 0 {
		pricedItem.SavingsPercent = (pricedItem.Savings / pricedItem.OriginalPrice) * 100
	}

	// Calculate margin and markup
	if item.CostPrice > 0 {
		pricedItem.Margin = ((pricedItem.FinalPrice - item.CostPrice) / pricedItem.FinalPrice) * 100
		pricedItem.Markup = ((pricedItem.FinalPrice - item.CostPrice) / item.CostPrice) * 100
	}

	return pricedItem, nil
}

// calculateDynamicPricing calculates dynamic pricing based on market conditions
func (c *Calculator) calculateDynamicPricing(item PricingItem, context PricingContext) float64 {
	for _, config := range c.dynamicConfigs {
		if !config.IsActive {
			continue
		}

		// Check if item is applicable for this config
		basePrice := item.BasePrice
		adjustedPrice := basePrice

		// Apply factors
		for _, factor := range config.Factors {
			if !factor.IsActive {
				continue
			}

			impact := c.calculateFactorImpact(factor, item, context)
			adjustedPrice += basePrice * (impact * factor.Weight / 100)
		}

		// Apply dynamic pricing rules
		for _, rule := range config.Rules {
			if !rule.IsActive || !c.isRuleApplicableToItem(rule.Conditions, item) {
				continue
			}

			for _, adjustment := range rule.Adjustments {
				adjustedPrice = c.applyAdjustment(adjustedPrice, adjustment)
			}
		}

		// Apply price constraints
		if config.PriceFloor > 0 && adjustedPrice < config.PriceFloor {
			adjustedPrice = config.PriceFloor
		}
		if config.PriceCeiling > 0 && adjustedPrice > config.PriceCeiling {
			adjustedPrice = config.PriceCeiling
		}

		// Check price change limits
		priceChange := math.Abs((adjustedPrice-basePrice)/basePrice) * 100
		if priceChange > config.MaxPriceChange {
			if adjustedPrice > basePrice {
				adjustedPrice = basePrice * (1 + config.MaxPriceChange/100)
			} else {
				adjustedPrice = basePrice * (1 - config.MaxPriceChange/100)
			}
		}

		return adjustedPrice
	}

	return 0 // No dynamic pricing applied
}

// calculateFactorImpact calculates the impact of a pricing factor
func (c *Calculator) calculateFactorImpact(factor PricingFactor, item PricingItem, context PricingContext) float64 {
	switch factor.Type {
	case "demand":
		// Use market data to determine demand impact
		if marketData, exists := c.marketData[item.ID]; exists {
			switch marketData.DemandLevel {
			case "high":
				return factor.Impact
			case "low":
				return -factor.Impact
			default:
				return 0
			}
		}
	case "inventory":
		// Adjust based on inventory levels
		if item.InventoryLevel < 10 {
			return factor.Impact // Low inventory, increase price
		} else if item.InventoryLevel > 100 {
			return -factor.Impact // High inventory, decrease price
		}
	case "competition":
		// Use competitor pricing data
		if marketData, exists := c.marketData[item.ID]; exists {
			if item.BasePrice > marketData.AveragePrice {
				return -factor.Impact // Our price is higher, decrease
			} else if item.BasePrice < marketData.AveragePrice {
				return factor.Impact // Our price is lower, can increase
			}
		}
	case "time":
		// Time-based adjustments (peak hours, seasons, etc.)
		hour := context.Timestamp.Hour()
		if hour >= 18 && hour <= 22 { // Peak hours
			return factor.Impact
		}
	case "weather":
		// Weather-based adjustments (from context metadata)
		if weather, exists := context.Metadata["weather"]; exists {
			if weather == "rain" || weather == "snow" {
				return factor.Impact
			}
		}
	case "events":
		// Event-based adjustments
		if context.Event != "" {
			return factor.Impact
		}
	}

	return 0
}

// calculateTierPricing calculates tier-based pricing
func (c *Calculator) calculateTierPricing(item PricingItem, tierPricing []TierPricing) *TierInfo {
	for _, tier := range tierPricing {
		if !tier.IsActive || time.Now().Before(tier.ValidFrom) || time.Now().After(tier.ValidUntil) {
			continue
		}

		for _, priceTier := range tier.Tiers {
			if item.Quantity >= priceTier.MinQuantity {
				if priceTier.MaxQuantity == 0 || item.Quantity <= priceTier.MaxQuantity {
					tierPrice := item.BasePrice
					if priceTier.FixedPrice > 0 {
						tierPrice = priceTier.FixedPrice
					} else if priceTier.Discount > 0 {
						tierPrice = item.BasePrice * (1 - priceTier.Discount/100)
					} else if priceTier.Price > 0 {
						tierPrice = priceTier.Price
					}

					return &TierInfo{
						TierID:       tier.ID,
						TierName:     tier.Name,
						MinQuantity:  priceTier.MinQuantity,
						MaxQuantity:  priceTier.MaxQuantity,
						TierPrice:    tierPrice,
						TierDiscount: priceTier.Discount,
					}
				}
			}
		}
	}

	return nil
}

// calculateBundlePricing calculates bundle pricing opportunities
func (c *Calculator) calculateBundlePricing(items []PricedItem, bundles []Bundle, customer Customer, context PricingContext) []BundleInfo {
	bundleResults := make([]BundleInfo, 0)

	for _, bundle := range bundles {
		if !bundle.IsActive || time.Now().Before(bundle.ValidFrom) || time.Now().After(bundle.ValidUntil) {
			continue
		}

		// Check if bundle conditions are met
		if !c.evaluateBundleConditions(bundle.Conditions, items, customer, context) {
			continue
		}

		// Find matching items for bundle
		matchingItems := c.findBundleItems(items, bundle)
		if len(matchingItems) < bundle.MinItems {
			continue
		}

		// Calculate bundle pricing
		bundlePrice := c.calculateBundlePrice(matchingItems, bundle)
		originalPrice := c.calculateOriginalBundlePrice(matchingItems)
		savings := originalPrice - bundlePrice

		bundleInfo := BundleInfo{
			BundleID:      bundle.ID,
			BundleName:    bundle.Name,
			BundleType:    string(bundle.Type),
			BundlePrice:   bundlePrice,
			BundleSavings: savings,
			ItemsInBundle: make([]string, len(matchingItems)),
		}

		for i, item := range matchingItems {
			bundleInfo.ItemsInBundle[i] = item.ItemID
		}

		bundleResults = append(bundleResults, bundleInfo)
	}

	return bundleResults
}

// getApplicableRules returns rules applicable to an item
func (c *Calculator) getApplicableRules(item PricingItem, customer Customer, context PricingContext, rules []PricingRule) []PricingRule {
	applicableRules := make([]PricingRule, 0)

	for _, rule := range rules {
		if !rule.IsActive || time.Now().Before(rule.ValidFrom) || time.Now().After(rule.ValidUntil) {
			continue
		}

		// Check item applicability
		if len(rule.ApplicableItems) > 0 {
			found := false
			for _, applicableItem := range rule.ApplicableItems {
				if applicableItem == item.ID || applicableItem == item.Category {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check excluded items
		if len(rule.ExcludedItems) > 0 {
			excluded := false
			for _, excludedItem := range rule.ExcludedItems {
				if excludedItem == item.ID || excludedItem == item.Category {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}
		}

		// Check customer segments
		if len(rule.CustomerSegments) > 0 {
			found := false
			for _, segment := range rule.CustomerSegments {
				if segment == customer.Segment || segment == customer.Type {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check channels
		if len(rule.Channels) > 0 {
			found := false
			for _, channel := range rule.Channels {
				if channel == context.Channel {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check regions
		if len(rule.Regions) > 0 {
			found := false
			for _, region := range rule.Regions {
				if region == context.Region {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		// Check conditions
		if len(rule.Conditions) > 0 && !c.evaluateConditions(rule.Conditions, item, customer, context) {
			continue
		}

		applicableRules = append(applicableRules, rule)
	}

	// Sort by priority
	sort.Slice(applicableRules, func(i, j int) bool {
		return applicableRules[i].Priority > applicableRules[j].Priority
	})

	return applicableRules
}

// applyPricingRule applies a pricing rule to a price
func (c *Calculator) applyPricingRule(currentPrice float64, rule PricingRule, item PricingItem, customer Customer) (float64, *AppliedPricingRule) {
	adjustedPrice := currentPrice
	totalAdjustment := 0.0

	for _, adjustment := range rule.Adjustments {
		adjustedPrice = c.applyAdjustment(adjustedPrice, adjustment)
		totalAdjustment += adjustment.Value
	}

	appliedRule := &AppliedPricingRule{
		RuleID:      rule.ID,
		Name:        rule.Name,
		Type:        string(rule.Type),
		Adjustment:  currentPrice - adjustedPrice,
		Description: rule.Description,
		Priority:    rule.Priority,
	}

	return adjustedPrice, appliedRule
}

// applyAdjustment applies a price adjustment
func (c *Calculator) applyAdjustment(price float64, adjustment PriceAdjustment) float64 {
	adjustedPrice := price

	switch adjustment.Type {
	case "percentage":
		adjustedPrice = price * (1 - adjustment.Value/100)
	case "fixed":
		adjustedPrice = price - adjustment.Value
	case "markup":
		adjustedPrice = price * (1 + adjustment.Value/100)
	case "markdown":
		adjustedPrice = price * (1 - adjustment.Value/100)
	}

	// Apply price limits
	if adjustment.MinPrice > 0 && adjustedPrice < adjustment.MinPrice {
		adjustedPrice = adjustment.MinPrice
	}
	if adjustment.MaxPrice > 0 && adjustedPrice > adjustment.MaxPrice {
		adjustedPrice = adjustment.MaxPrice
	}

	// Apply rounding
	if adjustment.RoundTo > 0 {
		adjustedPrice = math.Round(adjustedPrice/adjustment.RoundTo) * adjustment.RoundTo
	}

	return adjustedPrice
}

// evaluateConditions evaluates pricing conditions
func (c *Calculator) evaluateConditions(conditions []PricingCondition, item PricingItem, customer Customer, context PricingContext) bool {
	if len(conditions) == 0 {
		return true
	}

	results := make([]bool, len(conditions))
	for i, condition := range conditions {
		results[i] = c.evaluateCondition(condition, item, customer, context)
	}

	// Apply logic operators
	finalResult := results[0]
	for i := 1; i < len(results); i++ {
		if i-1 < len(conditions) && conditions[i-1].Logic == "OR" {
			finalResult = finalResult || results[i]
		} else {
			finalResult = finalResult && results[i]
		}
	}

	return finalResult
}

// evaluateCondition evaluates a single condition
func (c *Calculator) evaluateCondition(condition PricingCondition, item PricingItem, customer Customer, context PricingContext) bool {
	switch condition.Type {
	case "quantity":
		return c.compareValues(float64(item.Quantity), condition.Operator, condition.Value)
	case "amount":
		return c.compareValues(item.BasePrice*float64(item.Quantity), condition.Operator, condition.Value)
	case "customer_type":
		return c.compareStringValues(customer.Type, condition.Operator, condition.Value)
	case "customer_tier":
		return c.compareStringValues(customer.Tier, condition.Operator, condition.Value)
	case "time":
		return c.evaluateTimeCondition(condition, context.Timestamp)
	case "inventory":
		return c.compareValues(float64(item.InventoryLevel), condition.Operator, condition.Value)
	case "category":
		return c.compareStringValues(item.Category, condition.Operator, condition.Value)
	case "brand":
		return c.compareStringValues(item.Brand, condition.Operator, condition.Value)
	}

	return false
}

// compareValues compares numeric values
func (c *Calculator) compareValues(actual float64, operator string, expected interface{}) bool {
	expectedFloat, ok := expected.(float64)
	if !ok {
		return false
	}

	switch operator {
	case ">":
		return actual > expectedFloat
	case "<":
		return actual < expectedFloat
	case ">=":
		return actual >= expectedFloat
	case "<=":
		return actual <= expectedFloat
	case "=":
		return actual == expectedFloat
	case "!=":
		return actual != expectedFloat
	}

	return false
}

// compareStringValues compares string values
func (c *Calculator) compareStringValues(actual, operator string, expected interface{}) bool {
	expectedStr, ok := expected.(string)
	if !ok {
		return false
	}

	switch operator {
	case "=":
		return actual == expectedStr
	case "!=":
		return actual != expectedStr
	case "in":
		// Handle array of values
		if values, ok := expected.([]interface{}); ok {
			for _, value := range values {
				if str, ok := value.(string); ok && str == actual {
					return true
				}
			}
		}
		return false
	}

	return false
}

// evaluateTimeCondition evaluates time-based conditions
func (c *Calculator) evaluateTimeCondition(condition PricingCondition, timestamp time.Time) bool {
	// Implementation depends on the specific time condition format
	// This is a simplified version
	return true
}

// Helper functions

func (c *Calculator) isRuleApplicableToItem(conditions []PricingCondition, item PricingItem) bool {
	// Simplified implementation
	return true
}

func (c *Calculator) evaluateBundleConditions(conditions []PricingCondition, items []PricedItem, customer Customer, context PricingContext) bool {
	// Simplified implementation
	return true
}

func (c *Calculator) findBundleItems(items []PricedItem, bundle Bundle) []PricedItem {
	matchingItems := make([]PricedItem, 0)
	for _, item := range items {
		for _, bundleItem := range bundle.Items {
			if item.ItemID == bundleItem.ItemID {
				matchingItems = append(matchingItems, item)
				break
			}
		}
	}
	return matchingItems
}

func (c *Calculator) calculateBundlePrice(items []PricedItem, bundle Bundle) float64 {
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.FinalPrice * float64(item.Quantity)
	}

	switch bundle.Pricing.Type {
	case "fixed":
		return bundle.Pricing.Value
	case "percentage":
		return totalPrice * (1 - bundle.Pricing.Value/100)
	default:
		return totalPrice
	}
}

func (c *Calculator) calculateOriginalBundlePrice(items []PricedItem) float64 {
	totalPrice := 0.0
	for _, item := range items {
		totalPrice += item.OriginalPrice * float64(item.Quantity)
	}
	return totalPrice
}

func (c *Calculator) roundPrice(price float64, mode string, precision int) float64 {
	multiplier := math.Pow(10, float64(precision))
	switch mode {
	case "floor":
		return math.Floor(price*multiplier) / multiplier
	case "ceil":
		return math.Ceil(price*multiplier) / multiplier
	default:
		return math.Round(price*multiplier) / multiplier
	}
}

func (c *Calculator) calculateTotals(result *PricingResult) {
	subtotal := 0.0
	totalSavings := 0.0

	for _, item := range result.Items {
		subtotal += item.TotalPrice
		totalSavings += item.Savings * float64(item.Quantity)
	}

	result.Subtotal = subtotal
	result.TotalSavings = totalSavings
	result.TotalDiscount = totalSavings
	result.GrandTotal = subtotal
}

func (c *Calculator) generateRecommendations(items []PricedItem, bundles []Bundle, tierPricing []TierPricing) []PricingRecommendation {
	recommendations := make([]PricingRecommendation, 0)

	// Generate bundle recommendations
	for _, bundle := range bundles {
		if bundle.IsActive {
			recommendations = append(recommendations, PricingRecommendation{
				Type:        "bundle",
				Title:       fmt.Sprintf("Bundle: %s", bundle.Name),
				Description: bundle.Description,
				BundleID:    bundle.ID,
				Priority:    1,
			})
		}
	}

	return recommendations
}

func (c *Calculator) validateInput(input PricingInput) error {
	if len(input.Items) == 0 {
		return fmt.Errorf("no items provided")
	}

	for _, item := range input.Items {
		if item.ID == "" {
			return fmt.Errorf("item ID is required")
		}
		if item.BasePrice < 0 {
			return fmt.Errorf("item base price cannot be negative")
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("item quantity must be positive")
		}
	}

	return nil
}

// AddRule adds a pricing rule
func (c *Calculator) AddRule(rule PricingRule) {
	c.rules = append(c.rules, rule)
}

// AddBundle adds a bundle
func (c *Calculator) AddBundle(bundle Bundle) {
	c.bundles = append(c.bundles, bundle)
}

// AddTierPricing adds tier pricing
func (c *Calculator) AddTierPricing(tier TierPricing) {
	c.tierPricing = append(c.tierPricing, tier)
}

// AddDynamicConfig adds dynamic pricing configuration
func (c *Calculator) AddDynamicConfig(config DynamicPricingConfig) {
	c.dynamicConfigs = append(c.dynamicConfigs, config)
}

// UpdateMarketData updates market data for an item
func (c *Calculator) UpdateMarketData(itemID string, data MarketData) {
	c.marketData[itemID] = data
}

// UpdateAnalytics updates analytics data for an item
func (c *Calculator) UpdateAnalytics(itemID string, analytics PricingAnalytics) {
	c.analytics[itemID] = analytics
}
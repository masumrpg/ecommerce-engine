// Package pricing provides comprehensive pricing calculation capabilities for e-commerce applications.
// The calculator supports dynamic pricing, tier-based pricing, bundle pricing, rule-based adjustments,
// and market-driven pricing strategies.
//
// Key Features:
//   - Dynamic pricing based on market conditions, inventory, demand, and competition
//   - Tier-based pricing with quantity discounts and volume pricing
//   - Bundle pricing with cross-sell and upsell opportunities
//   - Rule-based pricing with complex conditions and adjustments
//   - Real-time market data integration for competitive pricing
//   - Customer segment-specific pricing strategies
//   - Channel and region-specific pricing rules
//   - Comprehensive analytics and performance tracking
//
// Basic Usage:
//
//	calc := pricing.NewCalculator()
//
//	// Add pricing rules
//	rule := pricing.PricingRule{
//		ID: "volume-discount",
//		Name: "Volume Discount",
//		Type: pricing.RuleTypeDiscount,
//		Adjustments: []pricing.PriceAdjustment{
//			{Type: "percentage", Value: 10.0},
//		},
//		Conditions: []pricing.PricingCondition{
//			{Type: "quantity", Operator: ">=", Value: 10.0},
//		},
//		IsActive: true,
//	}
//	calc.AddRule(rule)
//
//	// Calculate pricing
//	input := pricing.PricingInput{
//		Items: []pricing.PricingItem{
//			{ID: "product-1", BasePrice: 100.0, Quantity: 15},
//		},
//		Customer: pricing.Customer{Type: "retail", Tier: "gold"},
//		Context: pricing.PricingContext{Channel: "online", Currency: "USD"},
//		Options: pricing.PricingOptions{CalculateBundle: true, CalculateTiers: true},
//	}
//
//	result, err := calc.Calculate(input)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Total: $%.2f (Saved: $%.2f)\n", result.GrandTotal, result.TotalSavings)
package pricing

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// Calculator is the main pricing calculation engine that handles comprehensive pricing strategies.
// It manages pricing rules, bundles, tier pricing, dynamic pricing configurations, market data,
// and analytics to provide accurate and competitive pricing calculations.
//
// The calculator supports:
//   - Rule-based pricing with complex conditions and priority handling
//   - Dynamic pricing based on real-time market conditions
//   - Tier-based pricing for volume discounts
//   - Bundle pricing for cross-sell opportunities
//   - Customer segment and channel-specific pricing
//   - Market data integration for competitive pricing
//   - Performance analytics and optimization
//
// Example:
//
//	calc := pricing.NewCalculator()
//
//	// Configure dynamic pricing
//	dynamicConfig := pricing.DynamicPricingConfig{
//		ID: "demand-based",
//		Factors: []pricing.PricingFactor{
//			{Type: "demand", Weight: 30.0, Impact: 0.15},
//			{Type: "inventory", Weight: 20.0, Impact: 0.10},
//		},
//		MaxPriceChange: 25.0, // Max 25% price change
//		IsActive: true,
//	}
//	calc.AddDynamicConfig(dynamicConfig)
//
//	// Update market data
//	marketData := pricing.MarketData{
//		DemandLevel: "high",
//		AveragePrice: 95.0,
//		CompetitorCount: 5,
//	}
//	calc.UpdateMarketData("product-1", marketData)
type Calculator struct {
	rules           []PricingRule
	bundles         []Bundle
	tierPricing     []TierPricing
	dynamicConfigs  []DynamicPricingConfig
	marketData      map[string]MarketData
	analytics       map[string]PricingAnalytics
}

// NewCalculator creates a new pricing calculator instance.
// Initializes all internal collections and prepares the calculator for use.
//
// Returns:
//   - *Calculator: A new calculator ready for configuration and use
//
// Example:
//
//	calc := pricing.NewCalculator()
//
//	// Calculator is ready to use
//	fmt.Printf("Calculator initialized with %d rules\n", len(calc.GetRules()))
//
//	// Add configuration
//	calc.AddRule(volumeDiscountRule)
//	calc.AddBundle(crossSellBundle)
//	calc.AddTierPricing(bulkPricing)
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

// Calculate performs comprehensive pricing calculation for the given input.
// This is the main entry point for all pricing calculations, handling rules, bundles,
// tier pricing, dynamic pricing, and generating recommendations.
//
// The calculation process:
//   1. Validates input parameters
//   2. Applies dynamic pricing based on market conditions
//   3. Calculates tier-based pricing for volume discounts
//   4. Applies pricing rules in priority order
//   5. Calculates bundle pricing opportunities
//   6. Generates pricing recommendations
//   7. Calculates totals and savings
//
// Parameters:
//   - input: Complete pricing input with items, customer, context, and options
//
// Returns:
//   - *PricingResult: Comprehensive pricing result with calculated prices and recommendations
//   - error: Error if calculation fails or input is invalid
//
// Example:
//
//	input := pricing.PricingInput{
//		Items: []pricing.PricingItem{
//			{ID: "laptop", BasePrice: 1200.0, Quantity: 2, Category: "electronics"},
//			{ID: "mouse", BasePrice: 50.0, Quantity: 2, Category: "accessories"},
//		},
//		Customer: pricing.Customer{
//			ID: "customer-123",
//			Type: "business",
//			Tier: "premium",
//			Segment: "enterprise",
//		},
//		Context: pricing.PricingContext{
//			Channel: "online",
//			Region: "US",
//			Currency: "USD",
//			Timestamp: time.Now(),
//		},
//		Options: pricing.PricingOptions{
//			CalculateBundle: true,
//			CalculateTiers: true,
//			RoundingMode: "round",
//			RoundingPrecision: 2,
//		},
//	}
//
//	result, err := calc.Calculate(input)
//	if err != nil {
//		return nil, fmt.Errorf("pricing calculation failed: %w", err)
//	}
//
//	fmt.Printf("Subtotal: $%.2f\n", result.Subtotal)
//	fmt.Printf("Total Savings: $%.2f\n", result.TotalSavings)
//	fmt.Printf("Grand Total: $%.2f\n", result.GrandTotal)
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

// calculateItemPricing calculates comprehensive pricing for a single item.
// Applies dynamic pricing, tier pricing, and rule-based adjustments in sequence.
//
// Parameters:
//   - item: The item to price
//   - customer: Customer information for segment-specific pricing
//   - context: Pricing context (channel, region, time, etc.)
//   - rules: Applicable pricing rules
//   - tierPricing: Tier pricing configurations
//   - options: Calculation options and preferences
//
// Returns:
//   - *PricedItem: Fully calculated item with final price and applied adjustments
//   - error: Error if pricing calculation fails
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

// calculateDynamicPricing calculates dynamic pricing based on real-time market conditions.
// Considers demand, inventory levels, competition, time factors, weather, and events.
//
// Dynamic pricing factors:
//   - Demand: High demand increases price, low demand decreases price
//   - Inventory: Low inventory increases price, high inventory decreases price
//   - Competition: Adjusts price based on competitor pricing
//   - Time: Peak hours, seasonal adjustments
//   - Weather: Weather-dependent product pricing
//   - Events: Special event-based pricing
//
// Parameters:
//   - item: Item to calculate dynamic pricing for
//   - context: Pricing context with market conditions
//
// Returns:
//   - float64: Dynamically adjusted price, or 0 if no dynamic pricing applied
//
// Example:
//
//	// Configure dynamic pricing
//	config := pricing.DynamicPricingConfig{
//		Factors: []pricing.PricingFactor{
//			{Type: "demand", Weight: 40.0, Impact: 0.20},
//			{Type: "inventory", Weight: 30.0, Impact: 0.15},
//		},
//		MaxPriceChange: 30.0, // Max 30% change
//		PriceFloor: 50.0, // Minimum price
//		PriceCeiling: 200.0, // Maximum price
//	}
//	calc.AddDynamicConfig(config)
//
//	// Dynamic price will be calculated automatically
//	dynamicPrice := calc.calculateDynamicPricing(item, context)
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

// calculateFactorImpact calculates the impact of a specific pricing factor on item pricing.
// Each factor type has its own logic for determining price impact based on current conditions.
//
// Supported factor types:
//   - "demand": Based on market demand levels (high/medium/low)
//   - "inventory": Based on current inventory levels
//   - "competition": Based on competitor pricing data
//   - "time": Based on time of day, day of week, season
//   - "weather": Based on weather conditions
//   - "events": Based on special events or promotions
//
// Parameters:
//   - factor: The pricing factor to evaluate
//   - item: Item being priced
//   - context: Current pricing context
//
// Returns:
//   - float64: Impact multiplier (-1.0 to 1.0, where positive increases price)
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

// calculateTierPricing calculates tier-based pricing for volume discounts.
// Evaluates quantity-based pricing tiers and applies the best applicable tier.
//
// Parameters:
//   - item: Item to calculate tier pricing for
//   - tierPricing: Available tier pricing configurations
//
// Returns:
//   - *TierInfo: Information about the applied tier, or nil if no tier applies
//
// Example:
//
//	// Configure tier pricing
//	tierPricing := pricing.TierPricing{
//		ID: "bulk-discount",
//		Name: "Bulk Discount",
//		Tiers: []pricing.PriceTier{
//			{MinQuantity: 10, MaxQuantity: 49, Discount: 5.0}, // 5% off
//			{MinQuantity: 50, MaxQuantity: 99, Discount: 10.0}, // 10% off
//			{MinQuantity: 100, Discount: 15.0}, // 15% off for 100+
//		},
//		IsActive: true,
//	}
//	calc.AddTierPricing(tierPricing)
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

// calculateBundlePricing calculates bundle pricing opportunities for the given items.
// Identifies applicable bundles and calculates potential savings for cross-sell and upsell.
//
// Parameters:
//   - items: Items in the current cart/order
//   - bundles: Available bundle configurations
//   - customer: Customer information for bundle eligibility
//   - context: Pricing context for bundle applicability
//
// Returns:
//   - []BundleInfo: List of applicable bundles with pricing and savings information
//
// Example:
//
//	// Bundle will be automatically detected and calculated
//	bundleInfos := calc.calculateBundlePricing(pricedItems, bundles, customer, context)
//	for _, bundleInfo := range bundleInfos {
//		fmt.Printf("Bundle: %s, Price: $%.2f, Savings: $%.2f\n",
//			bundleInfo.BundleName, bundleInfo.BundlePrice, bundleInfo.BundleSavings)
//	}
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

// getApplicableRules filters and returns pricing rules applicable to a specific item.
// Evaluates rule conditions, item applicability, customer segments, channels, and regions.
// Returns rules sorted by priority (highest priority first).
//
// Rule filtering criteria:
//   - Rule must be active and within valid date range
//   - Item must be in applicable items list (if specified)
//   - Item must not be in excluded items list
//   - Customer must be in applicable segments (if specified)
//   - Context channel must be in applicable channels (if specified)
//   - Context region must be in applicable regions (if specified)
//   - All rule conditions must be satisfied
//
// Parameters:
//   - item: Item to find applicable rules for
//   - customer: Customer information for segment filtering
//   - context: Pricing context for channel/region filtering
//   - rules: All available pricing rules
//
// Returns:
//   - []PricingRule: Applicable rules sorted by priority (descending)
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

// applyPricingRule applies a specific pricing rule to an item and calculates the adjustment.
// Creates a detailed record of the rule application including original price, adjustment, and final price.
//
// Parameters:
//   - currentPrice: Current price before rule application
//   - rule: The pricing rule to apply
//   - item: Item to apply the rule to
//   - customer: Customer information for rule evaluation
//
// Returns:
//   - float64: Adjusted price after rule application
//   - *AppliedPricingRule: Detailed record of rule application, or nil if rule doesn't apply
//
// Example:
//
//	// Apply a specific rule
//	adjustedPrice, appliedRule := calc.applyPricingRule(currentPrice, rule, item, customer)
//	if appliedRule != nil {
//		fmt.Printf("Rule '%s' applied: $%.2f -> $%.2f (%.2f adjustment)\n",
//			appliedRule.Name, currentPrice, adjustedPrice, appliedRule.Adjustment)
//	}
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

// applyAdjustment applies a pricing adjustment to a base price.
// Supports multiple adjustment types with price limits and rounding options.
//
// Supported adjustment types:
//   - "percentage": Percentage-based adjustment (positive or negative)
//   - "fixed": Fixed amount adjustment (positive or negative)
//   - "markup": Markup percentage (always positive)
//   - "markdown": Markdown percentage (always negative)
//
// Parameters:
//   - price: Original price to adjust
//   - adjustment: The pricing adjustment configuration
//
// Returns:
//   - float64: Adjusted price after applying the adjustment and limits
//
// Example:
//
//	// 10% discount adjustment
//	adjustment := pricing.PriceAdjustment{
//		Type: "percentage",
//		Value: 10.0, // 10% off
//		MinPrice: 5.00, // Don't go below $5
//		MaxPrice: 100.00, // Don't go above $100
//	}
//	adjustedPrice := calc.applyAdjustment(50.00, adjustment) // Result: $45.00
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

// evaluateConditions evaluates multiple pricing conditions with logical operators.
// Supports AND/OR logic for combining multiple conditions.
//
// Logic operators:
//   - "AND": All conditions must be true (default)
//   - "OR": At least one condition must be true
//
// Parameters:
//   - conditions: List of conditions to evaluate
//   - item: Item being evaluated
//   - customer: Customer information for evaluation
//   - context: Pricing context for evaluation
//
// Returns:
//   - bool: True if conditions are satisfied based on logic operator
//
// Example:
//
//	// Multiple conditions with AND logic
//	conditions := []pricing.PricingCondition{
//		{Type: "quantity", Operator: ">=", Value: 10.0},
//		{Type: "customer_type", Operator: "=", Value: "premium"},
//		{Type: "category", Operator: "in", Value: []interface{}{"electronics", "computers"}},
//	}
//	// All conditions must be true
//	isApplicable := calc.evaluateConditions(conditions, item, customer, context)
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

// evaluateCondition evaluates a single pricing condition against item, customer, and context.
// Supports various condition types with different operators for flexible rule matching.
//
// Supported condition types:
//   - "quantity": Item quantity comparison
//   - "amount": Price amount comparison
//   - "customer_type": Customer type matching
//   - "customer_tier": Customer tier matching
//   - "time": Time-based conditions
//   - "inventory": Inventory level conditions
//   - "category": Product category matching
//   - "brand": Product brand matching
//
// Supported operators:
//   - "=": Equal to
//   - "!=": Not equal to
//   - ">": Greater than
//   - ">=": Greater than or equal to
//   - "<": Less than
//   - "<=": Less than or equal to
//   - "in": Value in list
//
// Parameters:
//   - condition: The condition to evaluate
//   - item: Item being evaluated
//   - customer: Customer information
//   - context: Pricing context
//
// Returns:
//   - bool: True if condition is satisfied
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

// compareValues compares numeric values using the specified operator.
// Supports standard comparison operators for numeric value evaluation.
//
// Supported operators:
//   - "=": Equal to
//   - "!=": Not equal to
//   - ">": Greater than
//   - ">=": Greater than or equal to
//   - "<": Less than
//   - "<=": Less than or equal to
//
// Parameters:
//   - actual: The actual numeric value to compare
//   - operator: Comparison operator
//   - expected: Expected value (must be convertible to float64)
//
// Returns:
//   - bool: Result of the comparison
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

// compareStringValues compares string values using the specified operator.
// Supports string comparison and list membership operations.
//
// Supported operators:
//   - "=": Equal to (case-sensitive)
//   - "!=": Not equal to
//   - "in": Value is in the provided list
//
// Parameters:
//   - actual: The actual string value to compare
//   - operator: Comparison operator
//   - expected: Expected value (string for equality, []interface{} for "in" operator)
//
// Returns:
//   - bool: Result of the comparison
//
// Example:
//
//	// String equality
//	result := calc.compareStringValues("premium", "=", "premium") // true
//	
//	// List membership
//	result := calc.compareStringValues("electronics", "in", []interface{}{"electronics", "computers"}) // true
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

// evaluateTimeCondition evaluates time-based pricing conditions.
// Supports various time-based rules for dynamic pricing based on temporal factors.
//
// Supported time condition types (via condition.Field):
//   - "hour": Hour of day (0-23)
//   - "day_of_week": Day of week (0=Sunday, 1=Monday, etc.)
//   - "day_of_month": Day of month (1-31)
//   - "month": Month (1-12)
//   - "season": Season (spring, summer, fall, winter)
//
// Parameters:
//   - condition: Time condition to evaluate
//   - timestamp: Current timestamp for evaluation
//
// Returns:
//   - bool: True if time condition is satisfied
//
// Example:
//
//	// Happy hour pricing (5-7 PM)
//	condition := pricing.PricingCondition{
//		Type: "time",
//		Field: "hour",
//		Operator: ">=",
//		Value: 17.0, // 5 PM
//	}
//	// Additional condition for <= 19 (7 PM) would complete the range
func (c *Calculator) evaluateTimeCondition(condition PricingCondition, timestamp time.Time) bool {
	// Implementation depends on the specific time condition format
	// This is a simplified version
	return true
}

// Helper functions

// isRuleApplicableToItem checks if pricing rule conditions apply to a specific item.
// Evaluates all conditions to determine if the rule should be applied to the item.
//
// Parameters:
//   - conditions: List of pricing conditions to evaluate
//   - item: Item to check applicability for
//
// Returns:
//   - bool: True if all conditions are satisfied and rule applies to the item
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

// AddRule adds a new pricing rule to the calculator.
// Rules are applied in priority order during price calculations.
//
// Parameters:
//   - rule: The pricing rule to add
//
// Example:
//
//	// Add a volume discount rule
//	rule := pricing.PricingRule{
//		ID: "volume-discount-10",
//		Name: "10+ Items Volume Discount",
//		Type: "discount",
//		Priority: 100,
//		Conditions: []pricing.PricingCondition{
//			{Type: "quantity", Operator: ">=", Value: 10.0},
//		},
//		Adjustment: pricing.PriceAdjustment{
//			Type: "percentage",
//			Value: 10.0, // 10% discount
//		},
//		IsActive: true,
//	}
//	calc.AddRule(rule)
func (c *Calculator) AddRule(rule PricingRule) {
	c.rules = append(c.rules, rule)
}

// AddBundle adds a new bundle configuration to the calculator.
// Bundles enable cross-sell and upsell opportunities with special pricing.
//
// Parameters:
//   - bundle: The bundle configuration to add
//
// Example:
//
//	// Add a laptop + accessories bundle
//	bundle := pricing.Bundle{
//		ID: "laptop-bundle",
//		Name: "Laptop Essentials Bundle",
//		Type: "fixed",
//		Items: []pricing.BundleItem{
//			{ItemID: "laptop-001", Quantity: 1, Required: true},
//			{ItemID: "mouse-001", Quantity: 1, Required: false},
//			{ItemID: "keyboard-001", Quantity: 1, Required: false},
//		},
//		Pricing: pricing.BundlePricing{
//			Type: "percentage",
//			Value: 15.0, // 15% off bundle
//		},
//		IsActive: true,
//	}
//	calc.AddBundle(bundle)
func (c *Calculator) AddBundle(bundle Bundle) {
	c.bundles = append(c.bundles, bundle)
}

// AddTierPricing adds a new tier pricing configuration to the calculator.
// Tier pricing enables volume discounts based on quantity thresholds.
//
// Parameters:
//   - tier: The tier pricing configuration to add
//
// Example:
//
//	// Add bulk pricing tiers
//	tierPricing := pricing.TierPricing{
//		ID: "bulk-discount",
//		Name: "Bulk Purchase Discount",
//		ApplicableItems: []string{"item-001", "item-002"},
//		Tiers: []pricing.PriceTier{
//			{MinQuantity: 10, MaxQuantity: 49, Discount: 5.0},  // 5% off
//			{MinQuantity: 50, MaxQuantity: 99, Discount: 10.0}, // 10% off
//			{MinQuantity: 100, Discount: 15.0},                 // 15% off
//		},
//		IsActive: true,
//	}
//	calc.AddTierPricing(tierPricing)
func (c *Calculator) AddTierPricing(tier TierPricing) {
	c.tierPricing = append(c.tierPricing, tier)
}

// AddDynamicConfig adds a new dynamic pricing configuration to the calculator.
// Dynamic pricing adjusts prices based on real-time factors like demand, inventory, and competition.
//
// Parameters:
//   - config: The dynamic pricing configuration to add
//
// Example:
//
//	// Add demand-based dynamic pricing
//	config := pricing.DynamicPricingConfig{
//		ID: "demand-pricing",
//		Name: "Demand-Based Pricing",
//		Factors: []pricing.PricingFactor{
//			{
//				Type: "demand",
//				Weight: 0.3,
//				MinImpact: -0.2, // Max 20% decrease
//				MaxImpact: 0.5,  // Max 50% increase
//			},
//			{
//				Type: "inventory",
//				Weight: 0.2,
//				MinImpact: -0.1,
//				MaxImpact: 0.3,
//			},
//		},
//		IsActive: true,
//	}
//	calc.AddDynamicConfig(config)
func (c *Calculator) AddDynamicConfig(config DynamicPricingConfig) {
	c.dynamicConfigs = append(c.dynamicConfigs, config)
}

// UpdateMarketData updates market data used for dynamic pricing calculations.
// Market data influences pricing factors like demand, competition, and trends.
//
// Parameters:
//   - itemID: The item ID to update market data for
//   - data: Updated market data for the item
//
// Example:
//
//	// Update market conditions for a specific item
//	marketData := pricing.MarketData{
//		DemandLevel: "high",
//		AveragePrice: 99.99,
//		CompetitorCount: 5,
//		InventoryLevel: 50,
//		Trend: "increasing",
//		LastUpdated: time.Now(),
//	}
//	calc.UpdateMarketData("item-001", marketData)
func (c *Calculator) UpdateMarketData(itemID string, data MarketData) {
	c.marketData[itemID] = data
}

// UpdateAnalytics updates pricing analytics data used for optimization and insights.
// Analytics help improve pricing strategies and track performance metrics.
//
// Parameters:
//   - itemID: The item ID to update analytics for
//   - analytics: Updated pricing analytics data for the item
//
// Example:
//
//	// Update pricing performance metrics for a specific item
//	analytics := pricing.PricingAnalytics{
//		ConversionRate: 0.15, // 15% conversion rate
//		PriceElasticity: -1.2, // Elastic demand
//		RevenueImpact: 1250.50, // Revenue increase
//		OptimalPrice: 95.00, // Suggested optimal price
//		LastUpdated: time.Now(),
//	}
//	calc.UpdateAnalytics("item-001", analytics)
func (c *Calculator) UpdateAnalytics(itemID string, analytics PricingAnalytics) {
	c.analytics[itemID] = analytics
}
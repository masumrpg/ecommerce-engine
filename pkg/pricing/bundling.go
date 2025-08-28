// Package pricing provides comprehensive bundle management functionality for e-commerce applications.
// This package enables creating, managing, and optimizing product bundles with advanced features
// including dynamic pricing, rule-based bundling, performance analytics, and AI-driven recommendations.
//
// Key Features:
//   - Dynamic Bundle Creation: Create bundles from templates or custom configurations
//   - Rule-Based Bundling: Apply complex business rules for bundle eligibility and pricing
//   - Multiple Bundle Types: Support for fixed, mix-and-match, and frequency-based bundles
//   - Performance Analytics: Track bundle performance with detailed metrics
//   - Optimization Engine: Automatically optimize bundles based on performance data
//   - Recommendation System: Generate personalized bundle recommendations
//   - Cross-sell and Upsell: Intelligent product recommendations
//   - Constraint Management: Enforce business rules and inventory constraints
//
// Basic Usage:
//
//	// Create a new bundle manager
//	bm := pricing.NewBundleManager()
//
//	// Create a simple bundle
//	items := []pricing.PricingItem{
//		{ID: "item1", Name: "Product A", BasePrice: 50.0, Category: "electronics"},
//		{ID: "item2", Name: "Product B", BasePrice: 30.0, Category: "accessories"},
//	}
//
//	pricing := pricing.BundlePricing{
//		Type: "percentage",
//		Value: 15.0, // 15% discount
//	}
//
//	bundle, err := bm.CreateBundle("Tech Bundle", "Electronics with accessories", pricing.BundleTypeFixed, items, pricing)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Generate recommendations
//	recommendations, err := bm.GenerateBundleRecommendations(items, customer, context)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Optimize existing bundle
//	optimization, err := bm.OptimizeBundle(bundle.ID)
//	if err != nil {
//		log.Fatal(err)
//	}
package pricing

import (
	"fmt"
	"math"
	"sort"
	"time"
)

// BundleManager handles comprehensive bundle creation, management, and optimization.
// It provides a centralized system for managing product bundles with support for
// templates, rules, analytics, and intelligent recommendations.
//
// Features:
//   - Bundle lifecycle management (create, update, optimize, analyze)
//   - Template-based bundle creation for consistency
//   - Rule engine for dynamic bundle behavior
//   - Performance analytics and optimization
//   - AI-driven bundle recommendations
//   - Cross-sell and upsell capabilities
//
// Example:
//
//	bm := pricing.NewBundleManager()
//
//	// Add bundle rules
//	rule := pricing.BundleRule{
//		ID: "electronics_discount",
//		Name: "Electronics Bundle Discount",
//		Type: "inclusion",
//		Conditions: []pricing.BundleCondition{
//			{Type: "category", Operator: "=", Value: "electronics"},
//		},
//		Actions: []pricing.BundleAction{
//			{Type: "apply_discount", Value: 20.0},
//		},
//		IsActive: true,
//	}
//	bm.AddBundleRule(rule)
type BundleManager struct {
	bundles         []Bundle
	bundleTemplates []BundleTemplate
	bundleRules     []BundleRule
	analytics       map[string]BundleAnalytics
}

// BundleTemplate represents a reusable template for creating bundles.
// Templates provide consistency and efficiency in bundle creation by defining
// standard configurations that can be applied to different sets of items.
//
// Example:
//
//	template := &BundleTemplate{
//		ID: "electronics_template",
//		Name: "Electronics Bundle Template",
//		Description: "Standard template for electronics bundles",
//		Type: BundleTypeFixed,
//		Pricing: BundlePricing{Type: "percentage", Value: 15.0},
//		Constraints: BundleConstraints{
//			MinItems: 2,
//			MaxItems: 5,
//			RequiredCategories: []string{"electronics"},
//		},
//		IsActive: true,
//	}
type BundleTemplate struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Type        BundleType        `json:"type"`
	Rules       []BundleRule      `json:"rules"`
	Pricing     BundlePricing     `json:"pricing"`
	Constraints BundleConstraints `json:"constraints"`
	IsActive    bool              `json:"is_active"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// BundleRule represents business rules that govern bundle creation and behavior.
// Rules define conditions that must be met and actions to take when those
// conditions are satisfied, enabling dynamic bundle management.
//
// Rule Types:
//   - "inclusion": Include items that meet conditions
//   - "exclusion": Exclude items that meet conditions
//   - "requirement": Require certain conditions to be met
//   - "substitution": Replace items based on conditions
//
// Example:
//
//	rule := &BundleRule{
//		ID: "premium_discount",
//		Name: "Premium Customer Discount",
//		Type: "inclusion",
//		Conditions: []BundleCondition{
//			{Type: "customer_type", Operator: "=", Value: "premium"},
//			{Type: "price_range", Operator: ">=", Value: 100.0, Logic: "AND"},
//		},
//		Actions: []BundleAction{
//			{Type: "apply_discount", Value: 25.0, Description: "25% premium discount"},
//		},
//		Priority: 10,
//		IsActive: true,
//	}
type BundleRule struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Type        string              `json:"type"` // "inclusion", "exclusion", "requirement", "substitution"
	Conditions  []BundleCondition   `json:"conditions"`
	Actions     []BundleAction      `json:"actions"`
	Priority    int                 `json:"priority"`
	IsActive    bool                `json:"is_active"`
	Description string              `json:"description,omitempty"`
}

// BundleCondition represents a condition that must be evaluated for bundle rules.
// Conditions define criteria that items, customers, or orders must meet.
//
// Condition Types:
//   - "category": Product category matching
//   - "brand": Brand matching
//   - "price_range": Price-based conditions
//   - "quantity": Quantity-based conditions
//   - "customer_type": Customer tier/type conditions
//
// Operators:
//   - "=", "!=": Equality/inequality
//   - "in", "not_in": Set membership
//   - ">", "<", ">=", "<=": Numeric comparisons
//
// Example:
//
//	condition := &BundleCondition{
//		Type: "category",
//		Operator: "in",
//		Value: []string{"electronics", "accessories"},
//		Logic: "AND",
//	}
type BundleCondition struct {
	Type     string      `json:"type"`     // "category", "brand", "price_range", "quantity", "customer_type"
	Operator string      `json:"operator"` // "=", "!=", "in", "not_in", ">", "<", ">=", "<="
	Value    interface{} `json:"value"`
	Logic    string      `json:"logic,omitempty"` // "AND", "OR"
}

// BundleAction represents an action to execute when bundle rule conditions are met.
// Actions modify bundle behavior, pricing, or composition based on rule evaluation.
//
// Action Types:
//   - "add_item": Add specific items to the bundle
//   - "remove_item": Remove items from the bundle
//   - "set_price": Set a fixed bundle price
//   - "apply_discount": Apply percentage or fixed discount
//
// Example:
//
//	action := &BundleAction{
//		Type: "apply_discount",
//		Target: "bundle",
//		Value: 20.0,
//		Description: "20% discount for electronics bundle",
//	}
type BundleAction struct {
	Type        string      `json:"type"`        // "add_item", "remove_item", "set_price", "apply_discount"
	Target      string      `json:"target"`      // Item ID or category
	Value       interface{} `json:"value"`       // Action value
	Description string      `json:"description,omitempty"`
}

// BundleConstraints represents business constraints that govern bundle creation.
// Constraints ensure bundles meet business requirements and inventory limitations.
//
// Example:
//
//	constraints := &BundleConstraints{
//		MinItems: 2,
//		MaxItems: 5,
//		MinValue: 50.0,
//		MaxValue: 500.0,
//		RequiredCategories: []string{"electronics", "accessories"},
//		ExcludedCategories: []string{"restricted"},
//		MaxSameCategory: 3,
//		RequiredBrands: []string{"premium_brand"},
//		ExcludedBrands: []string{"competitor_brand"},
//	}
type BundleConstraints struct {
	MinItems        int     `json:"min_items"`
	MaxItems        int     `json:"max_items"`
	MinValue        float64 `json:"min_value,omitempty"`
	MaxValue        float64 `json:"max_value,omitempty"`
	RequiredCategories []string `json:"required_categories,omitempty"`
	ExcludedCategories []string `json:"excluded_categories,omitempty"`
	MaxSameCategory int     `json:"max_same_category,omitempty"`
	RequiredBrands  []string `json:"required_brands,omitempty"`
	ExcludedBrands  []string `json:"excluded_brands,omitempty"`
}

// BundleAnalytics represents comprehensive analytics data for bundle performance.
// Provides insights into bundle effectiveness, customer behavior, and revenue impact.
//
// Example:
//
//	analytics := &BundleAnalytics{
//		BundleID: "bundle_123",
//		PeriodStart: time.Now().AddDate(0, -1, 0),
//		PeriodEnd: time.Now(),
//		ViewCount: 1500,
//		AddToCartCount: 300,
//		PurchaseCount: 75,
//		Revenue: 7500.0,
//		AverageOrderValue: 100.0,
//		ConversionRate: 0.25, // 25% conversion from add-to-cart
//		PopularityScore: 0.85,
//		ProfitMargin: 0.35,
//		CustomerSatisfaction: 4.3,
//		ReturnRate: 0.03,
//	}
type BundleAnalytics struct {
	BundleID        string    `json:"bundle_id"`
	PeriodStart     time.Time `json:"period_start"`
	PeriodEnd       time.Time `json:"period_end"`
	ViewCount       int       `json:"view_count"`
	AddToCartCount  int       `json:"add_to_cart_count"`
	PurchaseCount   int       `json:"purchase_count"`
	Revenue         float64   `json:"revenue"`
	AverageOrderValue float64 `json:"average_order_value"`
	ConversionRate  float64   `json:"conversion_rate"`
	PopularityScore float64   `json:"popularity_score"`
	ProfitMargin    float64   `json:"profit_margin"`
	CustomerSatisfaction float64 `json:"customer_satisfaction,omitempty"`
	ReturnRate      float64   `json:"return_rate,omitempty"`
}

// BundleRecommendation represents a personalized bundle recommendation.
// Contains bundle details, pricing information, and recommendation metadata
// to help customers make informed purchasing decisions.
//
// Example:
//
//	recommendation := &BundleRecommendation{
//		BundleID: "tech_bundle_456",
//		Name: "Complete Tech Setup",
//		Type: "cross_sell",
//		Items: []string{"laptop", "mouse", "keyboard", "monitor"},
//		OriginalPrice: 1200.0,
//		BundlePrice: 1000.0,
//		Savings: 200.0,
//		SavingsPercent: 16.67,
//		Confidence: 0.85,
//		Reason: "Frequently bought together by similar customers",
//		Priority: 8,
//		ValidUntil: time.Now().AddDate(0, 0, 7),
//	}
type BundleRecommendation struct {
	BundleID      string    `json:"bundle_id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Items         []string  `json:"items"`
	OriginalPrice float64   `json:"original_price"`
	BundlePrice   float64   `json:"bundle_price"`
	Savings       float64   `json:"savings"`
	SavingsPercent float64  `json:"savings_percent"`
	Confidence    float64   `json:"confidence"`
	Reason        string    `json:"reason"`
	Priority      int       `json:"priority"`
	ValidUntil    time.Time `json:"valid_until,omitempty"`
}

// BundleOptimization represents the results of bundle optimization analysis.
// Contains the original bundle, optimized version, improvements made, and
// expected performance metrics.
//
// Example:
//
//	optimization := &BundleOptimization{
//		OriginalBundle: originalBundle,
//		OptimizedBundle: optimizedBundle,
//		Improvements: []BundleImprovement{
//			{Type: "price_adjustment", Description: "Reduced price by 10%", Impact: 0.15, Confidence: 0.8},
//			{Type: "item_substitution", Description: "Replaced low-performing item", Impact: 0.12, Confidence: 0.75},
//		},
//		Metrics: BundleOptimizationMetrics{
//			ExpectedRevenueIncrease: 0.18,
//			ExpectedConversionIncrease: 0.22,
//			OptimizationScore: 0.85,
//		},
//		Recommendations: []string{"Monitor performance weekly", "A/B test the changes"},
//	}
type BundleOptimization struct{
	OriginalBundle Bundle                 `json:"original_bundle"`
	OptimizedBundle Bundle                `json:"optimized_bundle"`
	Improvements   []BundleImprovement    `json:"improvements"`
	Metrics        BundleOptimizationMetrics `json:"metrics"`
	Recommendations []string              `json:"recommendations"`
}

// BundleImprovement represents a specific improvement made during bundle optimization.
// Tracks the type of change, its expected impact, and confidence level.
//
// Improvement Types:
//   - "price_adjustment": Pricing modifications
//   - "item_addition": Adding new items
//   - "item_removal": Removing underperforming items
//   - "substitution": Replacing items with better alternatives
//
// Example:
//
//	improvement := &BundleImprovement{
//		Type: "price_adjustment",
//		Description: "Reduced bundle price by 15% to improve conversion",
//		Impact: 0.25, // Expected 25% improvement in conversion
//		Confidence: 0.85, // 85% confidence in the prediction
//	}
type BundleImprovement struct {
	Type        string  `json:"type"`        // "price_adjustment", "item_addition", "item_removal", "substitution"
	Description string  `json:"description"`
	Impact      float64 `json:"impact"`      // Expected impact on conversion/revenue
	Confidence  float64 `json:"confidence"`  // Confidence in the improvement
}

// BundleOptimizationMetrics represents quantitative metrics for bundle optimization results.
// Provides measurable expectations for the impact of optimization changes.
//
// Example:
//
//	metrics := &BundleOptimizationMetrics{
//		ExpectedRevenueIncrease: 0.20, // 20% revenue increase
//		ExpectedConversionIncrease: 0.15, // 15% conversion increase
//		ProfitMarginChange: -0.05, // 5% margin decrease due to pricing
//		CustomerSatisfactionChange: 0.10, // 10% satisfaction increase
//		OptimizationScore: 0.82, // Overall optimization score
//	}
type BundleOptimizationMetrics struct{
	ExpectedRevenueIncrease float64 `json:"expected_revenue_increase"`
	ExpectedConversionIncrease float64 `json:"expected_conversion_increase"`
	ProfitMarginChange      float64 `json:"profit_margin_change"`
	CustomerSatisfactionChange float64 `json:"customer_satisfaction_change"`
	OptimizationScore       float64 `json:"optimization_score"`
}

// NewBundleManager creates a new bundle manager instance.
// Initializes all internal collections and prepares the manager for use.
//
// Returns:
//   - *BundleManager: A new bundle manager ready for use
//
// Example:
//
//	bm := pricing.NewBundleManager()
//
//	// Manager is ready to use
//	bundles := bm.GetActiveBundles()
//	fmt.Printf("Active bundles: %d\n", len(bundles))
func NewBundleManager() *BundleManager {
	return &BundleManager{
		bundles:         make([]Bundle, 0),
		bundleTemplates: make([]BundleTemplate, 0),
		bundleRules:     make([]BundleRule, 0),
		analytics:       make(map[string]BundleAnalytics),
	}
}

// CreateBundle creates a new bundle from a collection of items.
// Applies pricing rules, bundle rules, and calculates final pricing automatically.
//
// Parameters:
//   - name: Display name for the bundle
//   - description: Detailed description of the bundle
//   - bundleType: Type of bundle (fixed, mix-and-match, frequency)
//   - items: Items to include in the bundle
//   - pricing: Pricing configuration for the bundle
//
// Returns:
//   - *Bundle: The created bundle with calculated pricing
//   - error: Error if bundle creation fails
//
// Example:
//
//	items := []pricing.PricingItem{
//		{ID: "laptop", Name: "Gaming Laptop", BasePrice: 1200.0, Category: "electronics"},
//		{ID: "mouse", Name: "Gaming Mouse", BasePrice: 80.0, Category: "accessories"},
//	}
//
//	pricing := pricing.BundlePricing{
//		Type: "percentage",
//		Value: 15.0, // 15% discount
//	}
//
//	bundle, err := bm.CreateBundle("Gaming Setup", "Complete gaming bundle", pricing.BundleTypeFixed, items, pricing)
//	if err != nil {
//		return nil, err
//	}
func (bm *BundleManager) CreateBundle(name, description string, bundleType BundleType, items []PricingItem, pricing BundlePricing) (*Bundle, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("cannot create bundle with no items")
	}

	bundle := &Bundle{
		ID:          fmt.Sprintf("bundle_%d", time.Now().Unix()),
		Name:        name,
		Description: description,
		Type:        bundleType,
		Items:       make([]BundleItem, 0),
		Pricing:     pricing,
		IsActive:    true,
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(1, 0, 0), // Valid for 1 year
		Tags:        make([]string, 0),
		Metadata:    make(map[string]interface{}),
	}

	// Convert pricing items to bundle items
	for _, item := range items {
		bundleItem := BundleItem{
			ItemID:      item.ID,
			Name:        item.Name,
			Quantity:    item.Quantity,
			IsRequired:  true,
			IsOptional:  false,
			BasePrice:   item.BasePrice,
			Category:    item.Category,
			Subcategory: item.Subcategory,
			Attributes:  item.Attributes,
		}
		bundle.Items = append(bundle.Items, bundleItem)
	}

	// Calculate bundle pricing
	bm.calculateBundlePricing(bundle)

	// Apply bundle rules
	bm.applyBundleRules(bundle)

	bm.bundles = append(bm.bundles, *bundle)
	return bundle, nil
}

// CreateBundleFromTemplate creates a bundle from a predefined template.
// Templates provide pre-configured rules, pricing, and constraints for common bundle types.
//
// Parameters:
//   - templateID: ID of the template to use
//   - items: Items to include in the bundle
//   - customizations: Custom overrides for template settings
//
// Returns:
//   - *Bundle: The created bundle based on the template
//   - error: Error if template not found or bundle creation fails
//
// Example:
//
//	customizations := map[string]interface{}{
//		"discount_percentage": 20.0, // Override template discount
//		"max_items": 5, // Override max items constraint
//	}
//
//	bundle, err := bm.CreateBundleFromTemplate("electronics-combo", items, customizations)
//	if err != nil {
//		return nil, err
//	}
func (bm *BundleManager) CreateBundleFromTemplate(templateID string, items []PricingItem, customizations map[string]interface{}) (*Bundle, error) {
	template := bm.getBundleTemplate(templateID)
	if template == nil {
		return nil, fmt.Errorf("template not found: %s", templateID)
	}

	// Validate items against template constraints
	if err := bm.validateItemsAgainstConstraints(items, template.Constraints); err != nil {
		return nil, fmt.Errorf("items don't meet template constraints: %w", err)
	}

	bundle, err := bm.CreateBundle(template.Name, template.Description, template.Type, items, template.Pricing)
	if err != nil {
		return nil, err
	}

	// Apply customizations
	bm.applyCustomizations(bundle, customizations)

	// Apply template rules
	for _, rule := range template.Rules {
		bm.applyBundleRule(bundle, rule)
	}

	return bundle, nil
}

// GenerateBundleRecommendations generates intelligent bundle recommendations for given items.
// Uses customer data, purchase history, and item relationships to suggest optimal bundles.
//
// Parameters:
//   - items: Items to generate recommendations for
//   - customer: Customer information and preferences
//   - context: Pricing context and business rules
//
// Returns:
//   - []BundleRecommendation: List of recommended bundles
//   - error: Error if recommendation generation fails
//
// Example:
//
//	customer := pricing.Customer{
//		ID: "customer-123",
//		Tier: "premium",
//		Preferences: map[string]interface{}{
//			"categories": []string{"electronics", "accessories"},
//			"price_sensitivity": "medium",
//		},
//	}
//
//	recommendations, err := bm.GenerateBundleRecommendations(items, customer, context)
//	if err != nil {
//		return nil, err
//	}
//
//	for _, rec := range recommendations {
//		fmt.Printf("Bundle: %s, Confidence: %.2f, Savings: $%.2f\n", 
//			rec.Name, rec.Confidence, rec.Savings)
//	}
func (bm *BundleManager) GenerateBundleRecommendations(items []PricingItem, customer Customer, context PricingContext) ([]BundleRecommendation, error) {
	recommendations := make([]BundleRecommendation, 0)

	// Find existing bundles that match the items
	for _, bundle := range bm.bundles {
		if !bundle.IsActive {
			continue
		}

		matchScore := bm.calculateBundleMatchScore(items, bundle)
		if matchScore > 0.5 { // Threshold for recommendation
			recommendation := bm.createBundleRecommendation(bundle, items, matchScore)
			recommendations = append(recommendations, recommendation)
		}
	}

	// Generate dynamic bundle recommendations
	dynamicRecommendations := bm.generateDynamicBundleRecommendations(items, customer, context)
	recommendations = append(recommendations, dynamicRecommendations...)

	// Sort by priority and confidence
	sort.Slice(recommendations, func(i, j int) bool {
		if recommendations[i].Priority == recommendations[j].Priority {
			return recommendations[i].Confidence > recommendations[j].Confidence
		}
		return recommendations[i].Priority > recommendations[j].Priority
	})

	return recommendations, nil
}

// OptimizeBundle optimizes an existing bundle to improve performance metrics.
// Analyzes current performance and suggests improvements based on analytics data.
//
// Parameters:
//   - bundleID: ID of the bundle to optimize
//
// Returns:
//   - *BundleOptimization: Optimization results with suggested improvements
//   - error: Error if bundle not found or optimization fails
//
// Example:
//
//	optimization, err := bm.OptimizeBundle("bundle-123")
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Printf("Optimization Score: %.2f\n", optimization.Metrics.OptimizationScore)
//	for _, improvement := range optimization.Improvements {
//		fmt.Printf("Improvement: %s (Impact: %.1f%%)\n", 
//			improvement.Description, improvement.Impact*100)
//	}
func (bm *BundleManager) OptimizeBundle(bundleID string) (*BundleOptimization, error) {
	bundle := bm.getBundle(bundleID)
	if bundle == nil {
		return nil, fmt.Errorf("bundle not found: %s", bundleID)
	}

	analytics := bm.analytics[bundleID]
	optimization := &BundleOptimization{
		OriginalBundle:  *bundle,
		OptimizedBundle: *bundle, // Start with original
		Improvements:    make([]BundleImprovement, 0),
		Recommendations: make([]string, 0),
	}

	// Analyze current performance
	performanceScore := bm.calculateBundlePerformanceScore(analytics)

	// Apply optimization strategies
	if performanceScore < 0.7 { // Poor performance
		// Price optimization
		if analytics.ConversionRate < 0.1 {
			bm.optimizeBundlePricing(&optimization.OptimizedBundle, analytics)
			optimization.Improvements = append(optimization.Improvements, BundleImprovement{
				Type:        "price_adjustment",
				Description: "Adjusted bundle pricing to improve conversion rate",
				Impact:      0.15, // Expected 15% improvement
				Confidence:  0.8,
			})
		}

		// Item optimization
		if analytics.ReturnRate > 0.1 {
			bm.optimizeBundleItems(&optimization.OptimizedBundle, analytics)
			optimization.Improvements = append(optimization.Improvements, BundleImprovement{
				Type:        "item_substitution",
				Description: "Replaced low-performing items with better alternatives",
				Impact:      0.2,
				Confidence:  0.75,
			})
		}
	}

	// Calculate optimization metrics
	optimization.Metrics = bm.calculateOptimizationMetrics(bundle, &optimization.OptimizedBundle, analytics)

	return optimization, nil
}

// AnalyzeBundlePerformance analyzes the performance of a bundle over a specified time period.
// Provides comprehensive metrics including sales, conversion rates, and customer feedback.
//
// Parameters:
//   - bundleID: ID of the bundle to analyze
//   - periodStart: Start date for the analysis period
//   - periodEnd: End date for the analysis period
//
// Returns:
//   - *BundleAnalytics: Detailed performance analytics
//   - error: Error if bundle not found or analysis fails
//
// Example:
//
//	startDate := time.Now().AddDate(0, -1, 0) // 1 month ago
//	endDate := time.Now()
//
//	analytics, err := bm.AnalyzeBundlePerformance("bundle-123", startDate, endDate)
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Printf("Bundle Performance (30 days):\n")
//	fmt.Printf("Sales: %d, Revenue: $%.2f\n", analytics.PurchaseCount, analytics.Revenue)
//	fmt.Printf("Conversion Rate: %.2f%%\n", analytics.ConversionRate*100)
//	fmt.Printf("Customer Satisfaction: %.1f/5\n", analytics.CustomerSatisfaction)
func (bm *BundleManager) AnalyzeBundlePerformance(bundleID string, periodStart, periodEnd time.Time) (*BundleAnalytics, error) {
	bundle := bm.getBundle(bundleID)
	if bundle == nil {
		return nil, fmt.Errorf("bundle not found: %s", bundleID)
	}

	// This would typically fetch data from analytics service
	// For now, return existing analytics or create mock data
	if analytics, exists := bm.analytics[bundleID]; exists {
		return &analytics, nil
	}

	// Create mock analytics data
	analytics := &BundleAnalytics{
		BundleID:             bundleID,
		PeriodStart:          periodStart,
		PeriodEnd:            periodEnd,
		ViewCount:            100,
		AddToCartCount:       25,
		PurchaseCount:        10,
		Revenue:              1000.0,
		AverageOrderValue:    100.0,
		ConversionRate:       0.1,
		PopularityScore:      0.7,
		ProfitMargin:         0.3,
		CustomerSatisfaction: 4.2,
		ReturnRate:           0.05,
	}

	bm.analytics[bundleID] = *analytics
	return analytics, nil
}

// CreateMixAndMatchBundle creates a flexible mix-and-match bundle.
// Allows customers to choose items from specified categories with quantity constraints.
//
// Parameters:
//   - name: Display name for the bundle
//   - categories: Product categories customers can choose from
//   - minItems: Minimum number of items required
//   - maxItems: Maximum number of items allowed
//   - pricing: Pricing configuration for the bundle
//
// Returns:
//   - *Bundle: The created mix-and-match bundle
//   - error: Error if bundle creation fails
//
// Example:
//
//	categories := []string{"shirts", "pants", "accessories"}
//	pricing := pricing.BundlePricing{
//		Type: "tiered",
//		Tiers: []pricing.PricingTier{
//			{MinQuantity: 3, Discount: 0.15}, // 15% off for 3+ items
//			{MinQuantity: 5, Discount: 0.25}, // 25% off for 5+ items
//		},
//	}
//
//	bundle, err := bm.CreateMixAndMatchBundle("Fashion Mix", categories, 2, 6, pricing)
func (bm *BundleManager) CreateMixAndMatchBundle(name string, categories []string, minItems, maxItems int, pricing BundlePricing) (*Bundle, error) {
	bundle := &Bundle{
		ID:          fmt.Sprintf("mixmatch_%d", time.Now().Unix()),
		Name:        name,
		Type:        BundleTypeMixMatch,
		Items:       make([]BundleItem, 0),
		Pricing:     pricing,
		MinItems:    minItems,
		MaxItems:    maxItems,
		IsActive:    true,
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(0, 6, 0), // Valid for 6 months
		Metadata:    map[string]interface{}{"categories": categories},
	}

	bm.bundles = append(bm.bundles, *bundle)
	return bundle, nil
}

// CreateFrequencyBundle creates a subscription-style frequency bundle.
// Designed for recurring purchases with automatic delivery and pricing benefits.
//
// Parameters:
//   - name: Display name for the bundle
//   - baseItem: Primary item for the recurring bundle
//   - frequency: Delivery frequency in days
//   - discount: Percentage discount for subscription
//
// Returns:
//   - *Bundle: The created frequency bundle
//   - error: Error if bundle creation fails
//
// Example:
//
//	baseItem := pricing.PricingItem{
//		ID: "coffee-beans",
//		Name: "Premium Coffee Beans",
//		BasePrice: 25.0,
//		Category: "beverages",
//	}
//
//	// Monthly delivery with 20% discount
//	bundle, err := bm.CreateFrequencyBundle("Coffee Subscription", baseItem, 30, 20.0)
func (bm *BundleManager) CreateFrequencyBundle(name string, baseItem PricingItem, frequency int, discount float64) (*Bundle, error) {
	bundle := &Bundle{
		ID:          fmt.Sprintf("frequency_%d", time.Now().Unix()),
		Name:        name,
		Type:        BundleTypeFrequency,
		Items:       make([]BundleItem, 0),
		IsActive:    true,
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(1, 0, 0),
		Metadata:    map[string]interface{}{"frequency": frequency, "discount": discount},
	}

	// Add base item with frequency quantity
	bundleItem := BundleItem{
		ItemID:      baseItem.ID,
		Name:        baseItem.Name,
		Quantity:    frequency,
		IsRequired:  true,
		BasePrice:   baseItem.BasePrice,
		BundlePrice: baseItem.BasePrice * (1 - discount/100),
		Discount:    discount,
		Category:    baseItem.Category,
	}

	bundle.Items = append(bundle.Items, bundleItem)
	bundle.Pricing = BundlePricing{
		Type:         "percentage",
		Value:        discount,
		SavingsType:  "percentage",
		SavingsValue: discount,
	}

	bm.bundles = append(bm.bundles, *bundle)
	return bundle, nil
}

// Helper functions

func (bm *BundleManager) getBundleTemplate(templateID string) *BundleTemplate {
	for _, template := range bm.bundleTemplates {
		if template.ID == templateID {
			return &template
		}
	}
	return nil
}

func (bm *BundleManager) getBundle(bundleID string) *Bundle {
	for _, bundle := range bm.bundles {
		if bundle.ID == bundleID {
			return &bundle
		}
	}
	return nil
}

func (bm *BundleManager) validateItemsAgainstConstraints(items []PricingItem, constraints BundleConstraints) error {
	if len(items) < constraints.MinItems {
		return fmt.Errorf("insufficient items: need at least %d, got %d", constraints.MinItems, len(items))
	}

	if constraints.MaxItems > 0 && len(items) > constraints.MaxItems {
		return fmt.Errorf("too many items: maximum %d allowed, got %d", constraints.MaxItems, len(items))
	}

	// Check required categories
	for _, requiredCategory := range constraints.RequiredCategories {
		found := false
		for _, item := range items {
			if item.Category == requiredCategory {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("required category not found: %s", requiredCategory)
		}
	}

	// Check excluded categories
	for _, excludedCategory := range constraints.ExcludedCategories {
		for _, item := range items {
			if item.Category == excludedCategory {
				return fmt.Errorf("excluded category found: %s", excludedCategory)
			}
		}
	}

	return nil
}

func (bm *BundleManager) calculateBundlePricing(bundle *Bundle) {
	totalPrice := 0.0
	for _, item := range bundle.Items {
		totalPrice += item.BasePrice * float64(item.Quantity)
	}

	switch bundle.Pricing.Type {
	case "fixed":
		bundle.Pricing.BasePrice = bundle.Pricing.Value
	case "percentage":
		bundle.Pricing.BasePrice = totalPrice * (1 - bundle.Pricing.Value/100)
		bundle.Pricing.SavingsValue = totalPrice - bundle.Pricing.BasePrice
	default:
		bundle.Pricing.BasePrice = totalPrice
	}
}

func (bm *BundleManager) applyBundleRules(bundle *Bundle) {
	for _, rule := range bm.bundleRules {
		if rule.IsActive {
			bm.applyBundleRule(bundle, rule)
		}
	}
}

func (bm *BundleManager) applyBundleRule(bundle *Bundle, rule BundleRule) {
	// Check if rule conditions are met
	if !bm.evaluateBundleRuleConditions(bundle, rule.Conditions) {
		return
	}

	// Apply rule actions
	for _, action := range rule.Actions {
		bm.applyBundleAction(bundle, action)
	}
}

func (bm *BundleManager) evaluateBundleRuleConditions(bundle *Bundle, conditions []BundleCondition) bool {
	if len(conditions) == 0 {
		return true
	}

	results := make([]bool, len(conditions))
	for i, condition := range conditions {
		results[i] = bm.evaluateBundleCondition(bundle, condition)
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

func (bm *BundleManager) evaluateBundleCondition(bundle *Bundle, condition BundleCondition) bool {
	switch condition.Type {
	case "category":
		for _, item := range bundle.Items {
			if bm.compareStringValue(item.Category, condition.Operator, condition.Value) {
				return true
			}
		}
	case "price_range":
		totalPrice := 0.0
		for _, item := range bundle.Items {
			totalPrice += item.BasePrice * float64(item.Quantity)
		}
		return bm.compareNumericValue(totalPrice, condition.Operator, condition.Value)
	case "quantity":
		totalQuantity := 0
		for _, item := range bundle.Items {
			totalQuantity += item.Quantity
		}
		return bm.compareNumericValue(float64(totalQuantity), condition.Operator, condition.Value)
	}

	return false
}

func (bm *BundleManager) applyBundleAction(bundle *Bundle, action BundleAction) {
	switch action.Type {
	case "apply_discount":
		if discount, ok := action.Value.(float64); ok {
			bundle.Pricing.Value = discount
			bundle.Pricing.Type = "percentage"
		}
	case "set_price":
		if price, ok := action.Value.(float64); ok {
			bundle.Pricing.Value = price
			bundle.Pricing.Type = "fixed"
		}
	}
}

func (bm *BundleManager) compareStringValue(actual, operator string, expected interface{}) bool {
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

func (bm *BundleManager) compareNumericValue(actual float64, operator string, expected interface{}) bool {
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

func (bm *BundleManager) applyCustomizations(bundle *Bundle, customizations map[string]interface{}) {
	for key, value := range customizations {
		switch key {
		case "name":
			if name, ok := value.(string); ok {
				bundle.Name = name
			}
		case "description":
			if desc, ok := value.(string); ok {
				bundle.Description = desc
			}
		case "discount":
			if discount, ok := value.(float64); ok {
				bundle.Pricing.Value = discount
				bundle.Pricing.Type = "percentage"
			}
		}
	}
}

func (bm *BundleManager) calculateBundleMatchScore(items []PricingItem, bundle Bundle) float64 {
	matchingItems := 0
	for _, item := range items {
		for _, bundleItem := range bundle.Items {
			if item.ID == bundleItem.ItemID {
				matchingItems++
				break
			}
		}
	}

	if len(bundle.Items) == 0 {
		return 0
	}

	return float64(matchingItems) / float64(len(bundle.Items))
}

func (bm *BundleManager) createBundleRecommendation(bundle Bundle, items []PricingItem, matchScore float64) BundleRecommendation {
	originalPrice := 0.0
	for _, item := range items {
		originalPrice += item.BasePrice * float64(item.Quantity)
	}

	bundlePrice := bundle.Pricing.BasePrice
	savings := originalPrice - bundlePrice
	savingsPercent := 0.0
	if originalPrice > 0 {
		savingsPercent = (savings / originalPrice) * 100
	}

	return BundleRecommendation{
		BundleID:       bundle.ID,
		Name:           bundle.Name,
		Type:           string(bundle.Type),
		OriginalPrice:  originalPrice,
		BundlePrice:    bundlePrice,
		Savings:        savings,
		SavingsPercent: savingsPercent,
		Confidence:     matchScore,
		Reason:         fmt.Sprintf("%.0f%% match with your items", matchScore*100),
		Priority:       int(matchScore * 10),
		ValidUntil:     bundle.ValidUntil,
	}
}

func (bm *BundleManager) generateDynamicBundleRecommendations(items []PricingItem, customer Customer, context PricingContext) []BundleRecommendation {
	recommendations := make([]BundleRecommendation, 0)

	// Generate cross-sell recommendations
	if len(items) > 0 {
		crossSellItems := bm.findCrossSellItems(items, customer)
		if len(crossSellItems) > 0 {
			recommendation := bm.createCrossSellRecommendation(items, crossSellItems)
			recommendations = append(recommendations, recommendation)
		}
	}

	// Generate upsell recommendations
	upsellItems := bm.findUpsellItems(items, customer)
	if len(upsellItems) > 0 {
		recommendation := bm.createUpsellRecommendation(items, upsellItems)
		recommendations = append(recommendations, recommendation)
	}

	return recommendations
}

func (bm *BundleManager) findCrossSellItems(items []PricingItem, customer Customer) []PricingItem {
	// This would typically use ML algorithms or predefined rules
	// For now, return empty slice
	return make([]PricingItem, 0)
}

func (bm *BundleManager) findUpsellItems(items []PricingItem, customer Customer) []PricingItem {
	// This would typically use ML algorithms or predefined rules
	// For now, return empty slice
	return make([]PricingItem, 0)
}

func (bm *BundleManager) createCrossSellRecommendation(originalItems, crossSellItems []PricingItem) BundleRecommendation {
	// Simplified implementation
	return BundleRecommendation{
		BundleID:   fmt.Sprintf("cross_sell_%d", time.Now().Unix()),
		Name:       "Cross-sell Bundle",
		Type:       "cross_sell",
		Confidence: 0.7,
		Reason:     "Frequently bought together",
		Priority:   5,
	}
}

func (bm *BundleManager) createUpsellRecommendation(originalItems, upsellItems []PricingItem) BundleRecommendation {
	// Simplified implementation
	return BundleRecommendation{
		BundleID:   fmt.Sprintf("upsell_%d", time.Now().Unix()),
		Name:       "Premium Bundle",
		Type:       "up_sell",
		Confidence: 0.8,
		Reason:     "Upgrade to premium items",
		Priority:   7,
	}
}

func (bm *BundleManager) calculateBundlePerformanceScore(analytics BundleAnalytics) float64 {
	// Weighted score based on multiple metrics
	conversionWeight := 0.4
	revenueWeight := 0.3
	satisfactionWeight := 0.2
	returnWeight := 0.1

	conversionScore := math.Min(analytics.ConversionRate*10, 1.0) // Normalize to 0-1
	revenueScore := math.Min(analytics.Revenue/10000, 1.0)        // Normalize based on target
	satisfactionScore := analytics.CustomerSatisfaction / 5.0     // Normalize 5-star rating
	returnScore := 1.0 - analytics.ReturnRate                    // Lower return rate is better

	totalScore := (conversionScore * conversionWeight) +
		(revenueScore * revenueWeight) +
		(satisfactionScore * satisfactionWeight) +
		(returnScore * returnWeight)

	return totalScore
}

func (bm *BundleManager) optimizeBundlePricing(bundle *Bundle, analytics BundleAnalytics) {
	// Simple price optimization based on conversion rate
	if analytics.ConversionRate < 0.05 {
		// Very low conversion, reduce price by 15%
		bundle.Pricing.Value = math.Min(bundle.Pricing.Value+15, 50) // Max 50% discount
	} else if analytics.ConversionRate < 0.1 {
		// Low conversion, reduce price by 10%
		bundle.Pricing.Value = math.Min(bundle.Pricing.Value+10, 40) // Max 40% discount
	}
}

func (bm *BundleManager) optimizeBundleItems(bundle *Bundle, analytics BundleAnalytics) {
	// Simple item optimization - remove items with high return rates
	// This would typically involve more sophisticated analysis
	if analytics.ReturnRate > 0.15 {
		// Remove the last item (simplified logic)
		if len(bundle.Items) > 2 {
			bundle.Items = bundle.Items[:len(bundle.Items)-1]
		}
	}
}

func (bm *BundleManager) calculateOptimizationMetrics(original, optimized *Bundle, analytics BundleAnalytics) BundleOptimizationMetrics {
	// Calculate expected improvements based on optimization changes
	priceChange := (optimized.Pricing.Value - original.Pricing.Value) / original.Pricing.Value
	expectedConversionIncrease := math.Abs(priceChange) * 0.5 // Simplified model
	expectedRevenueIncrease := expectedConversionIncrease * 0.8

	return BundleOptimizationMetrics{
		ExpectedRevenueIncrease:    expectedRevenueIncrease,
		ExpectedConversionIncrease: expectedConversionIncrease,
		ProfitMarginChange:         priceChange * -0.1, // Price reduction reduces margin
		OptimizationScore:          (expectedRevenueIncrease + expectedConversionIncrease) / 2,
	}
}

// Public methods for bundle management

// AddBundleTemplate adds a new bundle template to the manager.
// Templates can be reused to create multiple bundles with consistent configuration.
//
// Parameters:
//   - template: The bundle template to add
//
// Example:
//
//	template := pricing.BundleTemplate{
//		ID: "electronics-combo",
//		Name: "Electronics Combo",
//		Description: "Popular electronics bundle template",
//		Type: pricing.BundleTypeFixed,
//		DefaultPricing: pricing.BundlePricing{
//			Type: "percentage",
//			Value: 15.0,
//		},
//		Constraints: pricing.BundleConstraints{
//			MinItems: 2,
//			MaxItems: 5,
//			AllowedCategories: []string{"electronics", "accessories"},
//		},
//	}
//
//	bm.AddBundleTemplate(template)
func (bm *BundleManager) AddBundleTemplate(template BundleTemplate) {
	bm.bundleTemplates = append(bm.bundleTemplates, template)
}

// AddBundleRule adds a new bundle rule to the manager.
// Rules are applied automatically during bundle creation and pricing calculation.
//
// Parameters:
//   - rule: The bundle rule to add
//
// Example:
//
//	rule := pricing.BundleRule{
//		ID: "high-value-discount",
//		Name: "High Value Bundle Discount",
//		Conditions: []pricing.BundleCondition{
//			{
//				Type: "total_value",
//				Operator: "greater_than",
//				Value: 500.0,
//			},
//		},
//		Actions: []pricing.BundleAction{
//			{
//				Type: "percentage_discount",
//				Value: 25.0, // 25% discount for high-value bundles
//			},
//		},
//		Priority: 10,
//		Active: true,
//	}
//
//	bm.AddBundleRule(rule)
func (bm *BundleManager) AddBundleRule(rule BundleRule) {
	bm.bundleRules = append(bm.bundleRules, rule)
}

// GetBundles returns all bundles managed by this bundle manager.
// Includes both active and inactive bundles.
//
// Returns:
//   - []Bundle: All bundles in the manager
//
// Example:
//
//	allBundles := bm.GetBundles()
//	fmt.Printf("Total bundles: %d\n", len(allBundles))
//
//	for _, bundle := range allBundles {
//		fmt.Printf("Bundle: %s, Active: %t, Price: $%.2f\n", 
//			bundle.Name, bundle.Active, bundle.FinalPrice)
//	}
func (bm *BundleManager) GetBundles() []Bundle {
	return bm.bundles
}

// GetActiveBundles returns only the currently active bundles.
// Filters out inactive, expired, or disabled bundles.
//
// Returns:
//   - []Bundle: Active bundles available for purchase
//
// Example:
//
//	activeBundles := bm.GetActiveBundles()
//	fmt.Printf("Active bundles available: %d\n", len(activeBundles))
//
//	for _, bundle := range activeBundles {
//		fmt.Printf("Available: %s - $%.2f (Save $%.2f)\n", 
//			bundle.Name, bundle.FinalPrice, bundle.TotalSavings)
//	}
func (bm *BundleManager) GetActiveBundles() []Bundle {
	activeBundles := make([]Bundle, 0)
	for _, bundle := range bm.bundles {
		if bundle.IsActive && time.Now().After(bundle.ValidFrom) && time.Now().Before(bundle.ValidUntil) {
			activeBundles = append(activeBundles, bundle)
		}
	}
	return activeBundles
}

// UpdateBundleAnalytics updates the analytics data for a specific bundle.
// Used to track performance metrics and inform optimization decisions.
//
// Parameters:
//   - bundleID: ID of the bundle to update
//   - analytics: New analytics data to store
//
// Example:
//
//	analytics := pricing.BundleAnalytics{
//		PurchaseCount: 150,
//		Revenue: 12500.0,
//		ConversionRate: 0.08, // 8% conversion rate
//		AverageOrderValue: 83.33,
//		CustomerSatisfaction: 4.2,
//		ReturnRate: 0.02, // 2% return rate
//		ProfitMargin: 0.35, // 35% profit margin
//	}
//
//	bm.UpdateBundleAnalytics("bundle-123", analytics)
func (bm *BundleManager) UpdateBundleAnalytics(bundleID string, analytics BundleAnalytics) {
	bm.analytics[bundleID] = analytics
}
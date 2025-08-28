// Package loyalty provides comprehensive types and structures for implementing
// a flexible loyalty program system. This package defines all the core data
// structures needed for customer loyalty management, points calculation,
// rewards redemption, tier management, and analytics.
//
// Key Features:
//   - Multi-tier loyalty system (Bronze, Silver, Gold, Platinum)
//   - Flexible points earning and redemption system
//   - Rule-based loyalty program configuration
//   - Comprehensive reward management
//   - Referral program support
//   - Review-based rewards
//   - Advanced analytics and recommendations
//   - Transaction tracking and audit trails
//
// Basic Usage:
//
//	// Create a customer
//	customer := &Customer{
//		ID: "customer123",
//		Email: "customer@example.com",
//		Tier: TierBronze,
//		CurrentPoints: 500,
//		JoinDate: time.Now(),
//	}
//
//	// Define calculation input
//	input := &PointsCalculationInput{
//		Customer: *customer,
//		OrderAmount: 100.0,
//		Timestamp: time.Now(),
//	}
//
//	// Create a reward
//	reward := &Reward{
//		ID: "discount10",
//		Name: "$10 Discount",
//		Type: RewardTypeDiscount,
//		PointsCost: 1000,
//		Value: 10.0,
//	}
package loyalty

import (
	"time"
)

// LoyaltyTier represents customer loyalty tiers.
// Defines the hierarchical levels in the loyalty program based on customer spending.
type LoyaltyTier string

const (
	TierBronze   LoyaltyTier = "bronze"   // Entry level: $0-$999 annual spend
	TierSilver   LoyaltyTier = "silver"   // Mid level: $1,000-$4,999 annual spend
	TierGold     LoyaltyTier = "gold"     // High level: $5,000-$14,999 annual spend
	TierPlatinum LoyaltyTier = "platinum" // Premium level: $15,000+ annual spend
)

// PointsType represents different types of points that can be earned.
// Used to categorize and track the source of loyalty points.
type PointsType string

const (
	PointsTypeBase     PointsType = "base"     // Base points earned from regular purchases
	PointsTypeBonus    PointsType = "bonus"    // Bonus points from promotions and special offers
	PointsTypeReview   PointsType = "review"   // Points earned from writing product reviews
	PointsTypeReferral PointsType = "referral" // Points earned from successful referrals
	PointsTypeBirthday PointsType = "birthday" // Special birthday bonus points
	PointsTypeEvent    PointsType = "event"    // Points from special events and campaigns
)

// RewardType represents different types of rewards available in the loyalty program.
// Used to categorize rewards and determine their behavior and application.
type RewardType string

const (
	RewardTypeDiscount   RewardType = "discount"   // Percentage or fixed amount discounts
	RewardTypeFreeItem   RewardType = "free_item"  // Free products or gifts
	RewardTypeShipping   RewardType = "shipping"   // Free or discounted shipping
	RewardTypeUpgrade    RewardType = "upgrade"    // Service or product upgrades
	RewardTypeExperience RewardType = "experience" // Special experiences or events
	RewardTypeCashback   RewardType = "cashback"   // Cash back rewards
)

// TransactionType represents different types of loyalty point transactions.
// Used to track and categorize all point movements in the system.
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "earn"   // Points earned from purchases or activities
	TransactionTypeRedeem TransactionType = "redeem" // Points redeemed for rewards
	TransactionTypeExpire TransactionType = "expire" // Points expired due to time limits
	TransactionTypeAdjust TransactionType = "adjust" // Manual adjustments by administrators
)

// Customer represents a loyalty program customer.
// Contains all customer-related data including tier status, points balance,
// spending history, and preferences.
//
// Example:
//
//	customer := &Customer{
//		ID: "cust_123",
//		Email: "john@example.com",
//		Tier: TierSilver,
//		CurrentPoints: 2500,
//		LifetimePoints: 5000,
//		AnnualSpend: 3500.00,
//		JoinDate: time.Now().AddDate(-1, 0, 0),
//		IsActive: true,
//	}
type Customer struct {
	ID                string      `json:"id"`
	Email             string      `json:"email,omitempty"`
	Tier              LoyaltyTier `json:"tier"`
	CurrentPoints     int         `json:"current_points"`
	LifetimePoints    int         `json:"lifetime_points"`
	AnnualSpend       float64     `json:"annual_spend"`
	TotalSpend        float64     `json:"total_spend"`
	JoinDate          time.Time   `json:"join_date"`
	LastActivity      time.Time   `json:"last_activity"`
	Birthday          time.Time   `json:"birthday,omitempty"`
	ReferralCode      string      `json:"referral_code,omitempty"`
	ReferredBy        string      `json:"referred_by,omitempty"`
	TierAchievedDate  time.Time   `json:"tier_achieved_date"`
	NextTierThreshold float64     `json:"next_tier_threshold,omitempty"`
	IsActive          bool        `json:"is_active"`
	Preferences       CustomerPreferences `json:"preferences,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// CustomerPreferences represents customer communication and program preferences.
// Used to customize the loyalty program experience for each customer.
//
// Example:
//
//	prefs := &CustomerPreferences{
//		EmailNotifications: true,
//		SMSNotifications: false,
//		PreferredCategories: []string{"electronics", "books"},
//		CommunicationLanguage: "en",
//	}
type CustomerPreferences struct {
	EmailNotifications bool     `json:"email_notifications"`
	SMSNotifications   bool     `json:"sms_notifications"`
	PreferredCategories []string `json:"preferred_categories,omitempty"`
	CommunicationLanguage string `json:"communication_language,omitempty"`
}

// LoyaltyRule represents rules for earning and redeeming points.
// Defines conditions and actions that determine how customers earn points
// and what benefits they receive.
//
// Rule Types:
//   - "earn": Rules for earning points from purchases
//   - "redeem": Rules for redeeming points for rewards
//   - "tier": Rules for tier upgrades and benefits
//   - "bonus": Special bonus point rules
//
// Example:
//
//	rule := &LoyaltyRule{
//		ID: "base_points",
//		Name: "Base Points Earning",
//		Type: "earn",
//		Conditions: []LoyaltyCondition{
//			{Type: "amount", Operator: ">", Value: 0},
//		},
//		Actions: []LoyaltyAction{
//			{Type: "earn_points", Value: 1.0, PointsType: PointsTypeBase},
//		},
//		Priority: 1,
//		IsActive: true,
//	}
type LoyaltyRule struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Description      string            `json:"description,omitempty"`
	Type             string            `json:"type"` // "earn", "redeem", "tier", "bonus"
	Conditions       []LoyaltyCondition `json:"conditions,omitempty"`
	Actions          []LoyaltyAction   `json:"actions"`
	Priority         int               `json:"priority"`
	IsActive         bool              `json:"is_active"`
	ValidFrom        time.Time         `json:"valid_from"`
	ValidUntil       time.Time         `json:"valid_until"`
	ApplicableTiers  []LoyaltyTier     `json:"applicable_tiers,omitempty"`
	MaxUsagePerUser  int               `json:"max_usage_per_user,omitempty"`
	TotalUsageLimit  int               `json:"total_usage_limit,omitempty"`
	CurrentUsage     int               `json:"current_usage"`
	Categories       []string          `json:"categories,omitempty"`
	PaymentMethods   []string          `json:"payment_methods,omitempty"`
	Channels         []string          `json:"channels,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// LoyaltyCondition represents conditions that must be met for a loyalty rule to apply.
// Supports various condition types and operators for flexible rule configuration.
//
// Condition Types:
//   - "amount": Order amount conditions
//   - "quantity": Item quantity conditions
//   - "category": Product category conditions
//   - "payment_method": Payment method conditions
//   - "time": Time-based conditions
//   - "tier": Customer tier conditions
//
// Operators:
//   - ">", "<", ">=", "<=": Numeric comparisons
//   - "=", "!=": Equality comparisons
//   - "in": Value in list
//   - "between": Value between two numbers
//
// Example:
//
//	condition := &LoyaltyCondition{
//		Type: "amount",
//		Operator: ">=",
//		Value: 100.0,
//		Logic: "AND",
//	}
type LoyaltyCondition struct {
	Type     string      `json:"type"`     // "amount", "quantity", "category", "payment_method", "time", "tier"
	Operator string      `json:"operator"` // ">", "<", ">=", "<=", "=", "!=", "in", "between"
	Value    interface{} `json:"value"`    // Condition value
	Logic    string      `json:"logic,omitempty"` // "AND", "OR"
}

// LoyaltyAction represents actions to execute when loyalty rule conditions are met.
// Defines what happens when a rule is triggered.
//
// Action Types:
//   - "earn_points": Award points based on value
//   - "multiply_points": Multiply earned points by value
//   - "bonus_points": Award fixed bonus points
//   - "tier_upgrade": Upgrade customer tier
//
// Example:
//
//	action := &LoyaltyAction{
//		Type: "earn_points",
//		Value: 2.0, // 2 points per dollar
//		PointsType: PointsTypeBase,
//		Description: "Base points earning",
//	}
type LoyaltyAction struct {
	Type        string      `json:"type"`        // "earn_points", "multiply_points", "bonus_points", "tier_upgrade"
	Value       interface{} `json:"value"`       // Action value
	PointsType  PointsType  `json:"points_type,omitempty"`
	Description string      `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PointsTransaction represents a points transaction record.
// Tracks all point movements including earning, redemption, expiry, and adjustments.
// Provides complete audit trail for loyalty program activities.
//
// Example:
//
//	transaction := &PointsTransaction{
//		ID: "txn_123",
//		CustomerID: "cust_456",
//		Type: TransactionTypeEarn,
//		PointsType: PointsTypeBase,
//		Amount: 100,
//		Balance: 1500,
//		OrderID: "order_789",
//		Description: "Points earned from purchase",
//		Timestamp: time.Now(),
//	}
type PointsTransaction struct {
	ID              string          `json:"id"`
	CustomerID      string          `json:"customer_id"`
	Type            TransactionType `json:"type"`
	PointsType      PointsType      `json:"points_type"`
	Amount          int             `json:"amount"`
	Balance         int             `json:"balance"`
	OrderID         string          `json:"order_id,omitempty"`
	RuleID          string          `json:"rule_id,omitempty"`
	RewardID        string          `json:"reward_id,omitempty"`
	Description     string          `json:"description"`
	Timestamp       time.Time       `json:"timestamp"`
	ExpiryDate      time.Time       `json:"expiry_date,omitempty"`
	IsExpired       bool            `json:"is_expired"`
	Source          string          `json:"source,omitempty"` // "purchase", "review", "referral", "manual"
	Channel         string          `json:"channel,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// Reward represents a loyalty reward that customers can redeem with points.
// Supports various reward types including discounts, free items, and experiences.
//
// Reward Types:
//   - Discount: Percentage or fixed amount discounts
//   - Free Item: Complimentary products or services
//   - Shipping: Free or discounted shipping
//   - Upgrade: Service or product upgrades
//   - Experience: Special events or experiences
//   - Cashback: Direct cash rewards
//
// Example:
//
//	reward := &Reward{
//		ID: "discount_10",
//		Name: "$10 Off Your Order",
//		Type: RewardTypeDiscount,
//		PointsCost: 1000,
//		Value: 10.0,
//		MinOrderAmount: 50.0,
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(0, 3, 0),
//	}
type Reward struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Description      string      `json:"description,omitempty"`
	Type             RewardType  `json:"type"`
	PointsCost       int         `json:"points_cost"`
	Value            float64     `json:"value"`            // Monetary value or discount amount
	DiscountPercent  float64     `json:"discount_percent,omitempty"`
	MaxDiscount      float64     `json:"max_discount,omitempty"`
	MinOrderAmount   float64     `json:"min_order_amount,omitempty"`
	IsActive         bool        `json:"is_active"`
	ValidFrom        time.Time   `json:"valid_from"`
	ValidUntil       time.Time   `json:"valid_until"`
	Stock            int         `json:"stock,omitempty"`            // Available quantity
	MaxPerCustomer   int         `json:"max_per_customer,omitempty"` // Max redemptions per customer
	RequiredTier     LoyaltyTier `json:"required_tier,omitempty"`
	Categories       []string    `json:"categories,omitempty"`       // Applicable categories
	ExcludedCategories []string  `json:"excluded_categories,omitempty"`
	TermsConditions  string      `json:"terms_conditions,omitempty"`
	ImageURL         string      `json:"image_url,omitempty"`
	Tags             []string    `json:"tags,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// TierBenefit represents benefits and privileges for each loyalty tier.
// Defines the advantages customers receive at different tier levels.
//
// Benefits include:
//   - Points multipliers for enhanced earning
//   - Bonus points percentages
//   - Redemption bonuses for better value
//   - Free shipping thresholds
//   - Exclusive access and support
//   - Special occasion bonuses
//
// Example:
//
//	benefit := &TierBenefit{
//		Tier: TierGold,
//		PointsMultiplier: 1.5,
//		BonusPointsPercent: 10.0,
//		RedemptionBonus: 1.2,
//		FreeShippingThreshold: 75.0,
//		EarlyAccess: true,
//		PrioritySupport: true,
//		BirthdayBonus: 500,
//		MaxPointsExpiry: 24,
//	}
type TierBenefit struct {
	Tier                LoyaltyTier `json:"tier"`
	PointsMultiplier    float64     `json:"points_multiplier"`    // Base points multiplier
	BonusPointsPercent  float64     `json:"bonus_points_percent"` // Additional bonus percentage
	RedemptionBonus     float64     `json:"redemption_bonus"`     // Extra value when redeeming (20% = 1.2x value)
	FreeShippingThreshold float64   `json:"free_shipping_threshold,omitempty"`
	EarlyAccess         bool        `json:"early_access"`         // Early access to sales
	PrioritySupport     bool        `json:"priority_support"`     // Priority customer support
	BirthdayBonus       int         `json:"birthday_bonus"`       // Birthday bonus points
	AnnualBonus         int         `json:"annual_bonus"`         // Annual tier bonus
	ExclusiveRewards    []string    `json:"exclusive_rewards,omitempty"` // Exclusive reward IDs
	MaxPointsExpiry     int         `json:"max_points_expiry"`    // Points expiry in months
	Description         string      `json:"description,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// PointsCalculationInput represents input data for calculating loyalty points.
// Contains all necessary information to determine how many points a customer
// should earn from a purchase or activity.
//
// Example:
//
//	input := &PointsCalculationInput{
//		Customer: customer,
//		OrderAmount: 150.75,
//		Items: []OrderItem{
//			{ID: "item1", Category: "electronics", Price: 100.0, Quantity: 1},
//			{ID: "item2", Category: "books", Price: 50.75, Quantity: 1},
//		},
//		PaymentMethod: "credit_card",
//		Channel: "online",
//		Timestamp: time.Now(),
//		IsFirstPurchase: false,
//	}
type PointsCalculationInput struct {
	Customer       Customer        `json:"customer"`
	OrderAmount    float64         `json:"order_amount"`
	Items          []OrderItem     `json:"items,omitempty"`
	PaymentMethod  string          `json:"payment_method,omitempty"`
	Channel        string          `json:"channel,omitempty"`
	Timestamp      time.Time       `json:"timestamp"`
	OrderID        string          `json:"order_id,omitempty"`
	IsFirstPurchase bool           `json:"is_first_purchase,omitempty"`
	ReferralCode   string          `json:"referral_code,omitempty"`
	SpecialEvent   string          `json:"special_event,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
}

// OrderItem represents an individual item in an order for loyalty calculation.
// Used to apply category-specific rules and calculate item-level benefits.
//
// Example:
//
//	item := &OrderItem{
//		ID: "prod_123",
//		Name: "Wireless Headphones",
//		Category: "electronics",
//		Subcategory: "audio",
//		Brand: "TechBrand",
//		Price: 99.99,
//		Quantity: 2,
//		TotalAmount: 199.98,
//		IsGift: false,
//	}
type OrderItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Subcategory string  `json:"subcategory,omitempty"`
	Brand       string  `json:"brand,omitempty"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	IsGift      bool    `json:"is_gift,omitempty"`
	Attributes  map[string]interface{} `json:"attributes,omitempty"`
}

// RedemptionInput represents input data for redeeming loyalty points.
// Contains information needed to process a reward redemption request.
//
// Example:
//
//	input := &RedemptionInput{
//		Customer: customer,
//		RewardID: "discount_10",
//		Quantity: 1,
//		OrderAmount: 75.50,
//		Channel: "mobile_app",
//		Timestamp: time.Now(),
//	}
type RedemptionInput struct {
	Customer    Customer `json:"customer"`
	RewardID    string   `json:"reward_id"`
	Quantity    int      `json:"quantity,omitempty"`
	OrderAmount float64  `json:"order_amount,omitempty"`
	OrderItems  []OrderItem `json:"order_items,omitempty"`
	Channel     string   `json:"channel,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PointsCalculationResult represents the result of a points calculation.
// Contains detailed breakdown of points earned, applied rules, tier information,
// and recommendations for the customer.
//
// Example:
//
//	result := &PointsCalculationResult{
//		CustomerID: "cust_123",
//		BasePoints: 150,
//		BonusPoints: 50,
//		TotalPoints: 200,
//		NewBalance: 2700,
//		IsValid: true,
//	}
type PointsCalculationResult struct {
	CustomerID        string              `json:"customer_id"`
	BasePoints        int                 `json:"base_points"`
	BonusPoints       int                 `json:"bonus_points"`
	TotalPoints       int                 `json:"total_points"`
	PointsBreakdown   []PointsBreakdown   `json:"points_breakdown"`
	AppliedRules      []AppliedLoyaltyRule `json:"applied_rules"`
	NewBalance        int                 `json:"new_balance"`
	TierInfo          TierInfo            `json:"tier_info"`
	ExpiryDate        time.Time           `json:"expiry_date,omitempty"`
	Transactions      []PointsTransaction `json:"transactions"`
	Recommendations   []LoyaltyRecommendation `json:"recommendations,omitempty"`
	IsValid           bool                `json:"is_valid"`
	Errors            []string            `json:"errors,omitempty"`
	Warnings          []string            `json:"warnings,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// PointsBreakdown represents a detailed breakdown of how points were calculated.
// Provides transparency into the points calculation process.
//
// Example:
//
//	breakdown := &PointsBreakdown{
//		Source: "base_purchase",
//		Description: "Base points from purchase",
//		Amount: 100.0,
//		Rate: 1.0,
//		Multiplier: 1.5,
//		Points: 150,
//		PointsType: PointsTypeBase,
//		RuleID: "base_points_rule",
//	}
type PointsBreakdown struct {
	Source      string     `json:"source"`      // "base", "category_bonus", "payment_bonus", "tier_bonus"
	Description string     `json:"description"`
	Amount      float64    `json:"amount"`      // Order amount for this breakdown
	Rate        float64    `json:"rate"`        // Points per currency unit
	Multiplier  float64    `json:"multiplier"`  // Applied multiplier
	Points      int        `json:"points"`      // Calculated points
	PointsType  PointsType `json:"points_type"`
	RuleID      string     `json:"rule_id,omitempty"`
}

// AppliedLoyaltyRule represents a loyalty rule that was applied during calculation.
// Tracks which rules were triggered and their impact.
//
// Example:
//
//	appliedRule := &AppliedLoyaltyRule{
//		RuleID: "tier_bonus",
//		Name: "Silver Tier Bonus",
//		Type: "tier_bonus",
//		Description: "Additional points for Silver tier customers",
//		PointsAwarded: 50,
//		Multiplier: 1.25,
//	}
type AppliedLoyaltyRule struct {
	RuleID      string `json:"rule_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	PointsAwarded int  `json:"points_awarded"`
	Multiplier  float64 `json:"multiplier,omitempty"`
	BonusAmount float64 `json:"bonus_amount,omitempty"`
}

// TierInfo represents customer tier information and progress toward next tier.
// Provides insights into tier status and advancement opportunities.
//
// Example:
//
//	tierInfo := &TierInfo{
//		CurrentTier: TierSilver,
//		NextTier: TierGold,
//		CurrentSpend: 3500.0,
//		NextTierThreshold: 5000.0,
//		SpendToNextTier: 1500.0,
//		ProgressPercent: 70.0,
//		TierAchievedDate: time.Now().AddDate(0, -6, 0),
//		IsUpgraded: false,
//	}
type TierInfo struct {
	CurrentTier       LoyaltyTier `json:"current_tier"`
	NextTier          LoyaltyTier `json:"next_tier,omitempty"`
	CurrentSpend      float64     `json:"current_spend"`
	NextTierThreshold float64     `json:"next_tier_threshold,omitempty"`
	SpendToNextTier   float64     `json:"spend_to_next_tier,omitempty"`
	ProgressPercent   float64     `json:"progress_percent"`
	TierAchievedDate  time.Time   `json:"tier_achieved_date"`
	TierExpiryDate    time.Time   `json:"tier_expiry_date,omitempty"`
	Benefits          TierBenefit `json:"benefits"`
	IsUpgraded        bool        `json:"is_upgraded"`
}

// RedemptionResult represents the result of a points redemption transaction.
// Contains details about the redemption including codes and new balance.
//
// Example:
//
//	result := &RedemptionResult{
//		CustomerID: "cust_123",
//		RewardID: "discount_10",
//		RewardName: "$10 Off Your Order",
//		PointsRedeemed: 1000,
//		DiscountAmount: 10.0,
//		NewBalance: 1500,
//		RedemptionCode: "SAVE10-ABC123",
//		ValidUntil: time.Now().AddDate(0, 1, 0),
//		IsSuccessful: true,
//	}
type RedemptionResult struct {
	CustomerID       string            `json:"customer_id"`
	RewardID         string            `json:"reward_id"`
	RewardName       string            `json:"reward_name"`
	PointsRedeemed   int               `json:"points_redeemed"`
	DiscountAmount   float64           `json:"discount_amount,omitempty"`
	NewBalance       int               `json:"new_balance"`
	RedemptionCode   string            `json:"redemption_code,omitempty"`
	ValidUntil       time.Time         `json:"valid_until,omitempty"`
	Transaction      PointsTransaction `json:"transaction"`
	IsSuccessful     bool              `json:"is_successful"`
	Errors           []string          `json:"errors,omitempty"`
	Warnings         []string          `json:"warnings,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// LoyaltyRecommendation represents a personalized recommendation for the customer.
// Provides actionable suggestions to maximize loyalty benefits.
//
// Example:
//
//	recommendation := &LoyaltyRecommendation{
//		Type: "tier_upgrade",
//		Title: "Reach Gold Tier",
//		Description: "Spend $500 more to unlock Gold tier benefits",
//		ActionText: "Make a purchase of $500 or more",
//		Value: 500.0,
//		Priority: 1,
//		ValidUntil: time.Now().AddDate(0, 1, 0),
//	}
type LoyaltyRecommendation struct {
	Type        string  `json:"type"`        // "tier_upgrade", "reward", "bonus_opportunity"
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ActionText  string  `json:"action_text,omitempty"`
	Value       float64 `json:"value,omitempty"`       // Potential points or savings
	RewardID    string  `json:"reward_id,omitempty"`
	Priority    int     `json:"priority"`
	ValidUntil  time.Time `json:"valid_until,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// LoyaltyAnalytics represents comprehensive analytics data for the loyalty program.
// Provides insights into customer behavior and program performance.
//
// Example:
//
//	analytics := &LoyaltyAnalytics{
//		CustomerID: "cust_123",
//		PeriodStart: time.Now().AddDate(0, -1, 0),
//		PeriodEnd: time.Now(),
//		TotalPointsEarned: 2500,
//		TotalPointsRedeemed: 1200,
//		TotalSpend: 3500.00,
//		OrderCount: 15,
//		AverageOrderValue: 233.33,
//		RedemptionRate: 48.0,
//		EngagementScore: 85.5,
//	}
type LoyaltyAnalytics struct {
	CustomerID          string    `json:"customer_id"`
	PeriodStart         time.Time `json:"period_start"`
	PeriodEnd           time.Time `json:"period_end"`
	TotalPointsEarned   int       `json:"total_points_earned"`
	TotalPointsRedeemed int       `json:"total_points_redeemed"`
	TotalSpend          float64   `json:"total_spend"`
	OrderCount          int       `json:"order_count"`
	AverageOrderValue   float64   `json:"average_order_value"`
	RedemptionRate      float64   `json:"redemption_rate"`      // Points redeemed / points earned
	EngagementScore     float64   `json:"engagement_score"`     // Overall engagement score
	TierUpgrades        int       `json:"tier_upgrades"`
	ReferralsCount      int       `json:"referrals_count"`
	ReviewsCount        int       `json:"reviews_count"`
	LastActivity        time.Time `json:"last_activity"`
	ChurnRisk           float64   `json:"churn_risk"`           // Risk of customer churning (0-1)
	LifetimeValue       float64   `json:"lifetime_value"`
	PredictedValue      float64   `json:"predicted_value"`      // Predicted future value
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// LoyaltyConfiguration represents the overall loyalty program configuration.
// Defines program settings, tiers, benefits, and operational parameters.
//
// Example:
//
//	config := &LoyaltyConfiguration{
//		ProgramName: "VIP Rewards",
//		BaseCurrency: "USD",
//		BasePointsRate: 1.0,
//		RedemptionRate: 0.01,
//		PointsExpiry: 12,
//		MinRedemption: 100,
//		MaxRedemptionPercent: 50.0,
//		IsActive: true,
//		CreatedAt: time.Now(),
//		UpdatedAt: time.Now(),
//	}
type LoyaltyConfiguration struct {
	ProgramName         string        `json:"program_name"`
	BaseCurrency        string        `json:"base_currency"`
	BasePointsRate      float64       `json:"base_points_rate"`      // Points per currency unit
	RedemptionRate      float64       `json:"redemption_rate"`       // Currency value per point
	PointsExpiry        int           `json:"points_expiry"`         // Expiry in months
	MinRedemption       int           `json:"min_redemption"`        // Minimum points for redemption
	MaxRedemptionPercent float64      `json:"max_redemption_percent"` // Max % of order that can be paid with points
	TierThresholds      map[LoyaltyTier]float64 `json:"tier_thresholds"`
	TierBenefits        map[LoyaltyTier]TierBenefit `json:"tier_benefits"`
	DefaultRules        []LoyaltyRule `json:"default_rules"`
	IsActive            bool          `json:"is_active"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// ReferralProgram represents referral program configuration and settings.
// Defines how customers can refer others and earn rewards.
//
// Example:
//
//	referral := &ReferralProgram{
//		ID: "ref_program_1",
//		Name: "Friend Referral Program",
//		ReferrerReward: 500,
//		RefereeReward: 250,
//		MinOrderAmount: 50.0,
//		MaxReferrals: 10,
//		ValidityPeriod: 30,
//		IsActive: true,
//		ValidFrom: time.Now(),
//		ValidUntil: time.Now().AddDate(1, 0, 0),
//	}
type ReferralProgram struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
	ReferrerReward      int       `json:"referrer_reward"`      // Points for referrer
	RefereeReward       int       `json:"referee_reward"`       // Points for referee
	MinOrderAmount      float64   `json:"min_order_amount,omitempty"` // Min order for referee to qualify
	MaxReferrals        int       `json:"max_referrals,omitempty"`    // Max referrals per customer
	ValidityPeriod      int       `json:"validity_period"`      // Validity in days
	IsActive            bool      `json:"is_active"`
	ValidFrom           time.Time `json:"valid_from"`
	ValidUntil          time.Time `json:"valid_until"`
	TermsConditions     string    `json:"terms_conditions,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

// ReviewReward represents rewards configuration for product reviews.
// Defines how customers earn points for writing product reviews.
//
// Example:
//
//	reviewReward := &ReviewReward{
//		ID: "review_reward_1",
//		Name: "Product Review Rewards",
//		BasePoints: 50,
//		PhotoBonus: 25,
//		VideoBonus: 50,
//		VerifiedBonus: 25,
//		MinRating: 1,
//		MinCharacters: 50,
//		MaxPerProduct: 1,
//		MaxPerMonth: 10,
//		IsActive: true,
//	}
type ReviewReward struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	BasePoints      int       `json:"base_points"`      // Base points for review
	PhotoBonus      int       `json:"photo_bonus"`      // Extra points for photo
	VideoBonus      int       `json:"video_bonus"`      // Extra points for video
	VerifiedBonus   int       `json:"verified_bonus"`   // Extra points for verified purchase
	MinRating       int       `json:"min_rating"`       // Minimum rating required
	MinCharacters   int       `json:"min_characters"`   // Minimum review length
	MaxPerProduct   int       `json:"max_per_product"`  // Max reviews per product per customer
	MaxPerMonth     int       `json:"max_per_month"`    // Max reviews per month per customer
	IsActive        bool      `json:"is_active"`
	ValidFrom       time.Time `json:"valid_from"`
	ValidUntil      time.Time `json:"valid_until"`
	ApplicableCategories []string `json:"applicable_categories,omitempty"`
	ExcludedCategories   []string `json:"excluded_categories,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}
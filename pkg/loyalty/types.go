package loyalty

import (
	"time"
)

// LoyaltyTier represents customer loyalty tiers
type LoyaltyTier string

const (
	TierBronze   LoyaltyTier = "bronze"   // 0-999k annual spend
	TierSilver   LoyaltyTier = "silver"   // 1M-4.99M annual spend
	TierGold     LoyaltyTier = "gold"     // 5M-14.99M annual spend
	TierPlatinum LoyaltyTier = "platinum" // 15M+ annual spend
)

// PointsType represents different types of points
type PointsType string

const (
	PointsTypeBase     PointsType = "base"     // Base points from purchases
	PointsTypeBonus    PointsType = "bonus"    // Bonus points from promotions
	PointsTypeReview   PointsType = "review"   // Points from reviews
	PointsTypeReferral PointsType = "referral" // Points from referrals
	PointsTypeBirthday PointsType = "birthday" // Birthday bonus points
	PointsTypeEvent    PointsType = "event"    // Special event points
)

// RewardType represents different types of rewards
type RewardType string

const (
	RewardTypeDiscount   RewardType = "discount"   // Discount rewards
	RewardTypeFreeItem   RewardType = "free_item"  // Free item rewards
	RewardTypeShipping   RewardType = "shipping"   // Free shipping
	RewardTypeUpgrade    RewardType = "upgrade"    // Service upgrades
	RewardTypeExperience RewardType = "experience" // Experience rewards
	RewardTypeCashback   RewardType = "cashback"   // Cashback rewards
)

// TransactionType represents loyalty transaction types
type TransactionType string

const (
	TransactionTypeEarn   TransactionType = "earn"   // Points earned
	TransactionTypeRedeem TransactionType = "redeem" // Points redeemed
	TransactionTypeExpire TransactionType = "expire" // Points expired
	TransactionTypeAdjust TransactionType = "adjust" // Manual adjustment
)

// Customer represents a loyalty program customer
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

// CustomerPreferences represents customer preferences
type CustomerPreferences struct {
	EmailNotifications bool     `json:"email_notifications"`
	SMSNotifications   bool     `json:"sms_notifications"`
	PreferredCategories []string `json:"preferred_categories,omitempty"`
	CommunicationLanguage string `json:"communication_language,omitempty"`
}

// LoyaltyRule represents rules for earning/redeeming points
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

// LoyaltyCondition represents conditions for loyalty rules
type LoyaltyCondition struct {
	Type     string      `json:"type"`     // "amount", "quantity", "category", "payment_method", "time", "tier"
	Operator string      `json:"operator"` // ">", "<", ">=", "<=", "=", "!=", "in", "between"
	Value    interface{} `json:"value"`    // Condition value
	Logic    string      `json:"logic,omitempty"` // "AND", "OR"
}

// LoyaltyAction represents actions to take when conditions are met
type LoyaltyAction struct {
	Type        string      `json:"type"`        // "earn_points", "multiply_points", "bonus_points", "tier_upgrade"
	Value       interface{} `json:"value"`       // Action value
	PointsType  PointsType  `json:"points_type,omitempty"`
	Description string      `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PointsTransaction represents a points transaction
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

// Reward represents a loyalty reward
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

// TierBenefit represents benefits for each loyalty tier
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

// PointsCalculationInput represents input for points calculation
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

// OrderItem represents an item in an order for loyalty calculation
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

// RedemptionInput represents input for points redemption
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

// PointsCalculationResult represents the result of points calculation
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

// PointsBreakdown represents breakdown of points calculation
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

// AppliedLoyaltyRule represents a loyalty rule that was applied
type AppliedLoyaltyRule struct {
	RuleID      string `json:"rule_id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	PointsAwarded int  `json:"points_awarded"`
	Multiplier  float64 `json:"multiplier,omitempty"`
	BonusAmount float64 `json:"bonus_amount,omitempty"`
}

// TierInfo represents tier information and progress
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

// RedemptionResult represents the result of points redemption
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

// LoyaltyRecommendation represents loyalty program recommendations
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

// LoyaltyAnalytics represents analytics data for loyalty program
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

// LoyaltyConfiguration represents loyalty program configuration
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

// ReferralProgram represents referral program configuration
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

// ReviewReward represents rewards for product reviews
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
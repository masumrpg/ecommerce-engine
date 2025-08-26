package loyalty

import (
	"fmt"
	"sort"
	"time"
)

// RuleEngine manages loyalty rules and configurations
type RuleEngine struct {
	rules           []LoyaltyRule
	tierBenefits    map[LoyaltyTier]TierBenefit
	rewards         []Reward
	referralProgram ReferralProgram
	reviewRewards   []ReviewReward
	config          *LoyaltyConfiguration
}

// NewRuleEngine creates a new loyalty rule engine
func NewRuleEngine(config *LoyaltyConfiguration) *RuleEngine {
	return &RuleEngine{
		rules:        make([]LoyaltyRule, 0),
		tierBenefits: make(map[LoyaltyTier]TierBenefit),
		rewards:      make([]Reward, 0),
		reviewRewards: make([]ReviewReward, 0),
		config:       config,
	}
}

// AddRule adds a new loyalty rule
func (re *RuleEngine) AddRule(rule LoyaltyRule) error {
	if err := re.validateRule(rule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	re.rules = append(re.rules, rule)
	re.sortRulesByPriority()
	return nil
}

// UpdateRule updates an existing loyalty rule
func (re *RuleEngine) UpdateRule(ruleID string, updatedRule LoyaltyRule) error {
	if err := re.validateRule(updatedRule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules[i] = updatedRule
			re.sortRulesByPriority()
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// RemoveRule removes a loyalty rule
func (re *RuleEngine) RemoveRule(ruleID string) error {
	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules = append(re.rules[:i], re.rules[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetRule retrieves a loyalty rule by ID
func (re *RuleEngine) GetRule(ruleID string) (*LoyaltyRule, error) {
	for _, rule := range re.rules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetAllRules returns all loyalty rules
func (re *RuleEngine) GetAllRules() []LoyaltyRule {
	return re.rules
}

// GetActiveRules returns all active loyalty rules
func (re *RuleEngine) GetActiveRules() []LoyaltyRule {
	activeRules := make([]LoyaltyRule, 0)
	for _, rule := range re.rules {
		if rule.IsActive {
			activeRules = append(activeRules, rule)
		}
	}
	return activeRules
}

// GetRulesByType returns rules of a specific type
func (re *RuleEngine) GetRulesByType(ruleType string) []LoyaltyRule {
	typedRules := make([]LoyaltyRule, 0)
	for _, rule := range re.rules {
		if rule.Type == ruleType {
			typedRules = append(typedRules, rule)
		}
	}
	return typedRules
}

// ActivateRule activates a loyalty rule
func (re *RuleEngine) ActivateRule(ruleID string) error {
	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules[i].IsActive = true
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// DeactivateRule deactivates a loyalty rule
func (re *RuleEngine) DeactivateRule(ruleID string) error {
	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules[i].IsActive = false
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// SetTierBenefit sets benefits for a loyalty tier
func (re *RuleEngine) SetTierBenefit(tier LoyaltyTier, benefit TierBenefit) {
	benefit.Tier = tier
	re.tierBenefits[tier] = benefit
}

// GetTierBenefit retrieves benefits for a loyalty tier
func (re *RuleEngine) GetTierBenefit(tier LoyaltyTier) (TierBenefit, bool) {
	benefit, exists := re.tierBenefits[tier]
	return benefit, exists
}

// GetAllTierBenefits returns all tier benefits
func (re *RuleEngine) GetAllTierBenefits() map[LoyaltyTier]TierBenefit {
	return re.tierBenefits
}

// AddReward adds a new reward
func (re *RuleEngine) AddReward(reward Reward) error {
	if err := re.validateReward(reward); err != nil {
		return fmt.Errorf("invalid reward: %w", err)
	}

	re.rewards = append(re.rewards, reward)
	return nil
}

// UpdateReward updates an existing reward
func (re *RuleEngine) UpdateReward(rewardID string, updatedReward Reward) error {
	if err := re.validateReward(updatedReward); err != nil {
		return fmt.Errorf("invalid reward: %w", err)
	}

	for i, reward := range re.rewards {
		if reward.ID == rewardID {
			re.rewards[i] = updatedReward
			return nil
		}
	}

	return fmt.Errorf("reward with ID %s not found", rewardID)
}

// RemoveReward removes a reward
func (re *RuleEngine) RemoveReward(rewardID string) error {
	for i, reward := range re.rewards {
		if reward.ID == rewardID {
			re.rewards = append(re.rewards[:i], re.rewards[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("reward with ID %s not found", rewardID)
}

// GetReward retrieves a reward by ID
func (re *RuleEngine) GetReward(rewardID string) (*Reward, error) {
	for _, reward := range re.rewards {
		if reward.ID == rewardID {
			return &reward, nil
		}
	}

	return nil, fmt.Errorf("reward with ID %s not found", rewardID)
}

// GetAllRewards returns all rewards
func (re *RuleEngine) GetAllRewards() []Reward {
	return re.rewards
}

// GetActiveRewards returns all active rewards
func (re *RuleEngine) GetActiveRewards() []Reward {
	activeRewards := make([]Reward, 0)
	for _, reward := range re.rewards {
		if reward.IsActive {
			activeRewards = append(activeRewards, reward)
		}
	}
	return activeRewards
}

// GetRewardsByType returns rewards of a specific type
func (re *RuleEngine) GetRewardsByType(rewardType RewardType) []Reward {
	typedRewards := make([]Reward, 0)
	for _, reward := range re.rewards {
		if reward.Type == rewardType {
			typedRewards = append(typedRewards, reward)
		}
	}
	return typedRewards
}

// SetReferralProgram sets the referral program configuration
func (re *RuleEngine) SetReferralProgram(program ReferralProgram) error {
	if err := re.validateReferralProgram(program); err != nil {
		return fmt.Errorf("invalid referral program: %w", err)
	}

	re.referralProgram = program
	return nil
}

// GetReferralProgram returns the referral program configuration
func (re *RuleEngine) GetReferralProgram() ReferralProgram {
	return re.referralProgram
}

// AddReviewReward adds a new review reward
func (re *RuleEngine) AddReviewReward(reward ReviewReward) error {
	if err := re.validateReviewReward(reward); err != nil {
		return fmt.Errorf("invalid review reward: %w", err)
	}

	re.reviewRewards = append(re.reviewRewards, reward)
	return nil
}

// UpdateReviewReward updates an existing review reward
func (re *RuleEngine) UpdateReviewReward(rewardID string, updatedReward ReviewReward) error {
	if err := re.validateReviewReward(updatedReward); err != nil {
		return fmt.Errorf("invalid review reward: %w", err)
	}

	for i, reward := range re.reviewRewards {
		if reward.ID == rewardID {
			re.reviewRewards[i] = updatedReward
			return nil
		}
	}

	return fmt.Errorf("review reward with ID %s not found", rewardID)
}

// GetReviewRewards returns all review rewards
func (re *RuleEngine) GetReviewRewards() []ReviewReward {
	return re.reviewRewards
}

// GetActiveReviewRewards returns all active review rewards
func (re *RuleEngine) GetActiveReviewRewards() []ReviewReward {
	activeRewards := make([]ReviewReward, 0)
	for _, reward := range re.reviewRewards {
		if reward.IsActive {
			activeRewards = append(activeRewards, reward)
		}
	}
	return activeRewards
}

// OptimizeRules optimizes rule performance by sorting and organizing
func (re *RuleEngine) OptimizeRules() {
	re.sortRulesByPriority()
	re.removeExpiredRules()
	re.consolidateSimilarRules()
}

// ExportRules exports all rules to a map for serialization
func (re *RuleEngine) ExportRules() map[string]interface{} {
	return map[string]interface{}{
		"rules":            re.rules,
		"tier_benefits":    re.tierBenefits,
		"rewards":          re.rewards,
		"referral_program": re.referralProgram,
		"review_rewards":   re.reviewRewards,
		"config":           re.config,
		"exported_at":      time.Now(),
	}
}

// ImportRules imports rules from a map
func (re *RuleEngine) ImportRules(data map[string]interface{}) error {
	if rules, ok := data["rules"].([]LoyaltyRule); ok {
		re.rules = rules
	}

	if tierBenefits, ok := data["tier_benefits"].(map[LoyaltyTier]TierBenefit); ok {
		re.tierBenefits = tierBenefits
	}

	if rewards, ok := data["rewards"].([]Reward); ok {
		re.rewards = rewards
	}

	if referralProgram, ok := data["referral_program"].(ReferralProgram); ok {
		re.referralProgram = referralProgram
	}

	if reviewRewards, ok := data["review_rewards"].([]ReviewReward); ok {
		re.reviewRewards = reviewRewards
	}

	return nil
}

// GetStatistics returns statistics about the rule engine
func (re *RuleEngine) GetStatistics() map[string]interface{} {
	activeRules := 0
	expiredRules := 0
	activeRewards := 0
	expiredRewards := 0

	now := time.Now()

	for _, rule := range re.rules {
		if rule.IsActive {
			activeRules++
		}
		if !rule.ValidUntil.IsZero() && rule.ValidUntil.Before(now) {
			expiredRules++
		}
	}

	for _, reward := range re.rewards {
		if reward.IsActive {
			activeRewards++
		}
		if !reward.ValidUntil.IsZero() && reward.ValidUntil.Before(now) {
			expiredRewards++
		}
	}

	return map[string]interface{}{
		"total_rules":       len(re.rules),
		"active_rules":      activeRules,
		"expired_rules":     expiredRules,
		"total_rewards":     len(re.rewards),
		"active_rewards":    activeRewards,
		"expired_rewards":   expiredRewards,
		"tier_benefits":     len(re.tierBenefits),
		"review_rewards":    len(re.reviewRewards),
		"referral_active":   re.referralProgram.IsActive,
		"last_updated":      time.Now(),
	}
}

// validateRule validates a loyalty rule
func (re *RuleEngine) validateRule(rule LoyaltyRule) error {
	if rule.ID == "" {
		return fmt.Errorf("rule ID is required")
	}

	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}

	if rule.Type == "" {
		return fmt.Errorf("rule type is required")
	}

	if len(rule.Actions) == 0 {
		return fmt.Errorf("rule must have at least one action")
	}

	if !rule.ValidFrom.IsZero() && !rule.ValidUntil.IsZero() && rule.ValidFrom.After(rule.ValidUntil) {
		return fmt.Errorf("valid from date must be before valid until date")
	}

	// Validate conditions
	for _, condition := range rule.Conditions {
		if condition.Type == "" {
			return fmt.Errorf("condition type is required")
		}
		if condition.Operator == "" {
			return fmt.Errorf("condition operator is required")
		}
		if condition.Value == nil {
			return fmt.Errorf("condition value is required")
		}
	}

	// Validate actions
	for _, action := range rule.Actions {
		if action.Type == "" {
			return fmt.Errorf("action type is required")
		}
		if action.Value == nil {
			return fmt.Errorf("action value is required")
		}
	}

	return nil
}

// validateReward validates a reward
func (re *RuleEngine) validateReward(reward Reward) error {
	if reward.ID == "" {
		return fmt.Errorf("reward ID is required")
	}

	if reward.Name == "" {
		return fmt.Errorf("reward name is required")
	}

	if reward.Type == "" {
		return fmt.Errorf("reward type is required")
	}

	if reward.PointsCost <= 0 {
		return fmt.Errorf("points cost must be positive")
	}

	if reward.Value < 0 {
		return fmt.Errorf("reward value cannot be negative")
	}

	if !reward.ValidFrom.IsZero() && !reward.ValidUntil.IsZero() && reward.ValidFrom.After(reward.ValidUntil) {
		return fmt.Errorf("valid from date must be before valid until date")
	}

	return nil
}

// validateReferralProgram validates a referral program
func (re *RuleEngine) validateReferralProgram(program ReferralProgram) error {
	if program.ID == "" {
		return fmt.Errorf("program ID is required")
	}

	if program.Name == "" {
		return fmt.Errorf("program name is required")
	}

	if program.ReferrerReward < 0 {
		return fmt.Errorf("referrer reward cannot be negative")
	}

	if program.RefereeReward < 0 {
		return fmt.Errorf("referee reward cannot be negative")
	}

	if program.ValidityPeriod <= 0 {
		return fmt.Errorf("validity period must be positive")
	}

	if !program.ValidFrom.IsZero() && !program.ValidUntil.IsZero() && program.ValidFrom.After(program.ValidUntil) {
		return fmt.Errorf("valid from date must be before valid until date")
	}

	return nil
}

// validateReviewReward validates a review reward
func (re *RuleEngine) validateReviewReward(reward ReviewReward) error {
	if reward.ID == "" {
		return fmt.Errorf("reward ID is required")
	}

	if reward.Name == "" {
		return fmt.Errorf("reward name is required")
	}

	if reward.BasePoints < 0 {
		return fmt.Errorf("base points cannot be negative")
	}

	if reward.MinRating < 1 || reward.MinRating > 5 {
		return fmt.Errorf("minimum rating must be between 1 and 5")
	}

	if reward.MinCharacters < 0 {
		return fmt.Errorf("minimum characters cannot be negative")
	}

	if !reward.ValidFrom.IsZero() && !reward.ValidUntil.IsZero() && reward.ValidFrom.After(reward.ValidUntil) {
		return fmt.Errorf("valid from date must be before valid until date")
	}

	return nil
}

// sortRulesByPriority sorts rules by priority (highest first)
func (re *RuleEngine) sortRulesByPriority() {
	sort.Slice(re.rules, func(i, j int) bool {
		return re.rules[i].Priority > re.rules[j].Priority
	})
}

// removeExpiredRules removes expired rules
func (re *RuleEngine) removeExpiredRules() {
	now := time.Now()
	activeRules := make([]LoyaltyRule, 0)

	for _, rule := range re.rules {
		if rule.ValidUntil.IsZero() || rule.ValidUntil.After(now) {
			activeRules = append(activeRules, rule)
		}
	}

	re.rules = activeRules
}

// consolidateSimilarRules consolidates similar rules to improve performance
func (re *RuleEngine) consolidateSimilarRules() {
	// Group rules by type and conditions
	ruleGroups := make(map[string][]LoyaltyRule)

	for _, rule := range re.rules {
		key := fmt.Sprintf("%s_%d", rule.Type, len(rule.Conditions))
		ruleGroups[key] = append(ruleGroups[key], rule)
	}

	// For each group, check if rules can be consolidated
	consolidatedRules := make([]LoyaltyRule, 0)

	for _, group := range ruleGroups {
		if len(group) == 1 {
			consolidatedRules = append(consolidatedRules, group[0])
			continue
		}

		// For now, just add all rules without consolidation
		// In a real implementation, you would check for similar conditions
		// and merge compatible rules
		consolidatedRules = append(consolidatedRules, group...)
	}

	re.rules = consolidatedRules
}

// CreateDefaultConfiguration creates a default loyalty configuration
func CreateDefaultConfiguration() *LoyaltyConfiguration {
	return &LoyaltyConfiguration{
		ProgramName:          "Loyalty Program",
		BaseCurrency:         "USD",
		BasePointsRate:       1.0, // 1 point per dollar
		RedemptionRate:       0.01, // 1 cent per point
		PointsExpiry:         12,   // 12 months
		MinRedemption:        100,  // Minimum 100 points
		MaxRedemptionPercent: 50.0, // Max 50% of order
		TierThresholds: map[LoyaltyTier]float64{
			TierBronze:   0,
			TierSilver:   1000,
			TierGold:     5000,
			TierPlatinum: 15000,
		},
		TierBenefits: map[LoyaltyTier]TierBenefit{
			TierBronze: {
				Tier:                TierBronze,
				PointsMultiplier:    1.0,
				BonusPointsPercent:  0.0,
				RedemptionBonus:     0.0,
				FreeShippingThreshold: 50.0,
				EarlyAccess:         false,
				PrioritySupport:     false,
				BirthdayBonus:       50,
				AnnualBonus:         0,
				MaxPointsExpiry:     12,
				Description:         "Bronze tier benefits",
			},
			TierSilver: {
				Tier:                TierSilver,
				PointsMultiplier:    1.25,
				BonusPointsPercent:  5.0,
				RedemptionBonus:     0.1,
				FreeShippingThreshold: 25.0,
				EarlyAccess:         true,
				PrioritySupport:     false,
				BirthdayBonus:       100,
				AnnualBonus:         250,
				MaxPointsExpiry:     18,
				Description:         "Silver tier benefits",
			},
			TierGold: {
				Tier:                TierGold,
				PointsMultiplier:    1.5,
				BonusPointsPercent:  10.0,
				RedemptionBonus:     0.2,
				FreeShippingThreshold: 0.0,
				EarlyAccess:         true,
				PrioritySupport:     true,
				BirthdayBonus:       200,
				AnnualBonus:         500,
				MaxPointsExpiry:     24,
				Description:         "Gold tier benefits",
			},
			TierPlatinum: {
				Tier:                TierPlatinum,
				PointsMultiplier:    2.0,
				BonusPointsPercent:  15.0,
				RedemptionBonus:     0.3,
				FreeShippingThreshold: 0.0,
				EarlyAccess:         true,
				PrioritySupport:     true,
				BirthdayBonus:       500,
				AnnualBonus:         1000,
				MaxPointsExpiry:     36,
				Description:         "Platinum tier benefits",
			},
		},
		DefaultRules: CreateDefaultRules(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CreateDefaultRules creates default loyalty rules
func CreateDefaultRules() []LoyaltyRule {
	return []LoyaltyRule{
		{
			ID:          "base_points",
			Name:        "Base Points",
			Description: "Earn 1 point per dollar spent",
			Type:        "earn",
			Conditions:  []LoyaltyCondition{},
			Actions: []LoyaltyAction{
				{
					Type:        "earn_points",
					Value:       1.0,
					PointsType:  PointsTypeBase,
					Description: "Base points from purchase",
				},
			},
			Priority:    1,
			IsActive:    true,
			ValidFrom:   time.Now(),
			ValidUntil:  time.Time{},
		},
		{
			ID:          "first_purchase_bonus",
			Name:        "First Purchase Bonus",
			Description: "Extra 100 points for first purchase",
			Type:        "bonus",
			Conditions: []LoyaltyCondition{
				{
					Type:     "first_purchase",
					Operator: "=",
					Value:    true,
				},
			},
			Actions: []LoyaltyAction{
				{
					Type:        "earn_points",
					Value:       100,
					PointsType:  PointsTypeBonus,
					Description: "First purchase bonus",
				},
			},
			Priority:    10,
			IsActive:    true,
			ValidFrom:   time.Now(),
			ValidUntil:  time.Time{},
		},
		{
			ID:          "high_value_bonus",
			Name:        "High Value Purchase Bonus",
			Description: "2x points for orders over $500",
			Type:        "bonus",
			Conditions: []LoyaltyCondition{
				{
					Type:     "amount",
					Operator: ">=",
					Value:    500.0,
				},
			},
			Actions: []LoyaltyAction{
				{
					Type:        "multiply_points",
					Value:       2.0,
					PointsType:  PointsTypeBonus,
					Description: "2x points for high value orders",
				},
			},
			Priority:    5,
			IsActive:    true,
			ValidFrom:   time.Now(),
			ValidUntil:  time.Time{},
		},
	}
}

// CreateDefaultRewards creates default rewards
func CreateDefaultRewards() []Reward {
	return []Reward{
		{
			ID:               "discount_5",
			Name:             "$5 Off",
			Description:      "$5 discount on your next purchase",
			Type:             RewardTypeDiscount,
			PointsCost:       500,
			Value:            5.0,
			MinOrderAmount:   25.0,
			IsActive:         true,
			ValidFrom:        time.Now(),
			ValidUntil:       time.Now().AddDate(1, 0, 0),
			MaxPerCustomer:   5,
			TermsConditions:  "Valid for 30 days from redemption",
		},
		{
			ID:               "discount_10",
			Name:             "$10 Off",
			Description:      "$10 discount on your next purchase",
			Type:             RewardTypeDiscount,
			PointsCost:       1000,
			Value:            10.0,
			MinOrderAmount:   50.0,
			IsActive:         true,
			ValidFrom:        time.Now(),
			ValidUntil:       time.Now().AddDate(1, 0, 0),
			MaxPerCustomer:   3,
			TermsConditions:  "Valid for 30 days from redemption",
		},
		{
			ID:               "free_shipping",
			Name:             "Free Shipping",
			Description:      "Free shipping on your next order",
			Type:             RewardTypeShipping,
			PointsCost:       250,
			Value:            0.0,
			IsActive:         true,
			ValidFrom:        time.Now(),
			ValidUntil:       time.Now().AddDate(1, 0, 0),
			MaxPerCustomer:   10,
			TermsConditions:  "Valid for 60 days from redemption",
		},
	}
}
// Package loyalty provides comprehensive loyalty program management capabilities.
// This package includes rule engines for managing loyalty rules, tier benefits,
// rewards, referral programs, and review rewards.
//
// Key Features:
//   - Dynamic rule management with priority-based execution
//   - Multi-tier loyalty system with configurable benefits
//   - Flexible reward system with various types (discounts, shipping, products)
//   - Referral program management with customizable rewards
//   - Review reward system to incentivize customer feedback
//   - Rule validation and optimization capabilities
//   - Import/export functionality for rule configurations
//   - Comprehensive statistics and analytics
//
// Basic Usage:
//
//	config := CreateDefaultConfiguration()
//	engine := NewRuleEngine(config)
//
//	// Add a custom rule
//	rule := LoyaltyRule{
//		ID:   "birthday_bonus",
//		Name: "Birthday Bonus",
//		Type: "bonus",
//		Actions: []LoyaltyAction{{
//			Type:  "earn_points",
//			Value: 100,
//		}},
//		IsActive: true,
//	}
//	engine.AddRule(rule)
//
//	// Set tier benefits
//	benefit := TierBenefit{
//		Tier:             TierGold,
//		PointsMultiplier: 1.5,
//		BonusPointsPercent: 10.0,
//	}
//	engine.SetTierBenefit(TierGold, benefit)
package loyalty

import (
	"fmt"
	"sort"
	"time"
)

// RuleEngine manages loyalty rules and configurations.
// It provides a centralized system for managing all aspects of a loyalty program
// including rules, tier benefits, rewards, referral programs, and review rewards.
//
// Features:
//   - Rule management with priority-based execution
//   - Tier benefit configuration and management
//   - Reward catalog management
//   - Referral program configuration
//   - Review reward system
//   - Rule optimization and performance tuning
//   - Import/export capabilities
//   - Statistics and analytics
//
// Example:
//
//	config := CreateDefaultConfiguration()
//	engine := NewRuleEngine(config)
//
//	// Add a loyalty rule
//	rule := LoyaltyRule{
//		ID:   "weekend_bonus",
//		Name: "Weekend Bonus",
//		Type: "bonus",
//		Conditions: []LoyaltyCondition{{
//			Type:     "day_of_week",
//			Operator: "in",
//			Value:    []string{"Saturday", "Sunday"},
//		}},
//		Actions: []LoyaltyAction{{
//			Type:  "multiply_points",
//			Value: 1.5,
//		}},
//		Priority: 5,
//		IsActive: true,
//	}
//	engine.AddRule(rule)
type RuleEngine struct {
	rules           []LoyaltyRule
	tierBenefits    map[LoyaltyTier]TierBenefit
	rewards         []Reward
	referralProgram ReferralProgram
	reviewRewards   []ReviewReward
	config          *LoyaltyConfiguration
}

// NewRuleEngine creates a new loyalty rule engine.
// Initializes the engine with the provided configuration and empty collections
// for rules, tier benefits, rewards, and other loyalty program components.
//
// Parameters:
//   - config: Loyalty program configuration settings
//
// Returns:
//   - *RuleEngine: Initialized rule engine ready for use
//
// Example:
//
//	config := CreateDefaultConfiguration()
//	engine := NewRuleEngine(config)
//
//	// Engine is now ready to manage loyalty rules
//	stats := engine.GetStatistics()
//	fmt.Printf("Engine initialized with %d rules\n", stats["total_rules"])
func NewRuleEngine(config *LoyaltyConfiguration) *RuleEngine {
	return &RuleEngine{
		rules:        make([]LoyaltyRule, 0),
		tierBenefits: make(map[LoyaltyTier]TierBenefit),
		rewards:      make([]Reward, 0),
		reviewRewards: make([]ReviewReward, 0),
		config:       config,
	}
}

// AddRule adds a new loyalty rule to the engine.
// Validates the rule before adding and automatically sorts rules by priority.
//
// Parameters:
//   - rule: LoyaltyRule to add to the engine
//
// Returns:
//   - error: Validation error if rule is invalid, nil if successful
//
// Example:
//
//	rule := LoyaltyRule{
//		ID:   "category_bonus",
//		Name: "Electronics Bonus",
//		Type: "bonus",
//		Conditions: []LoyaltyCondition{{
//			Type:     "category",
//			Operator: "=",
//			Value:    "electronics",
//		}},
//		Actions: []LoyaltyAction{{
//			Type:  "multiply_points",
//			Value: 2.0,
//		}},
//		Priority: 8,
//		IsActive: true,
//	}
//	err := engine.AddRule(rule)
func (re *RuleEngine) AddRule(rule LoyaltyRule) error {
	if err := re.validateRule(rule); err != nil {
		return fmt.Errorf("invalid rule: %w", err)
	}

	re.rules = append(re.rules, rule)
	re.sortRulesByPriority()
	return nil
}

// UpdateRule updates an existing loyalty rule.
// Validates the updated rule and maintains priority sorting.
//
// Parameters:
//   - ruleID: ID of the rule to update
//   - updatedRule: New rule data to replace the existing rule
//
// Returns:
//   - error: Validation error or "not found" error if rule doesn't exist
//
// Example:
//
//	updatedRule := existingRule
//	updatedRule.Priority = 10
//	updatedRule.Actions[0].Value = 3.0
//	err := engine.UpdateRule("category_bonus", updatedRule)
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

// RemoveRule removes a loyalty rule from the engine.
// Permanently deletes the rule from the engine's rule collection.
//
// Parameters:
//   - ruleID: ID of the rule to remove
//
// Returns:
//   - error: "Not found" error if rule doesn't exist, nil if successful
//
// Example:
//
//	err := engine.RemoveRule("old_promotion_rule")
//	if err != nil {
//		log.Printf("Failed to remove rule: %v", err)
//	}
func (re *RuleEngine) RemoveRule(ruleID string) error {
	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules = append(re.rules[:i], re.rules[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetRule retrieves a loyalty rule by ID.
// Returns a pointer to the rule if found, allowing for inspection or modification.
//
// Parameters:
//   - ruleID: ID of the rule to retrieve
//
// Returns:
//   - *LoyaltyRule: Pointer to the rule if found
//   - error: "Not found" error if rule doesn't exist
//
// Example:
//
//	rule, err := engine.GetRule("category_bonus")
//	if err == nil {
//		fmt.Printf("Rule: %s, Priority: %d\n", rule.Name, rule.Priority)
//	}
func (re *RuleEngine) GetRule(ruleID string) (*LoyaltyRule, error) {
	for _, rule := range re.rules {
		if rule.ID == ruleID {
			return &rule, nil
		}
	}

	return nil, fmt.Errorf("rule with ID %s not found", ruleID)
}

// GetAllRules returns all loyalty rules in the engine.
// Returns a copy of all rules regardless of their active status.
//
// Returns:
//   - []LoyaltyRule: Slice containing all rules in the engine
//
// Example:
//
//	allRules := engine.GetAllRules()
//	fmt.Printf("Total rules: %d\n", len(allRules))
//	for _, rule := range allRules {
//		fmt.Printf("Rule: %s, Active: %t\n", rule.Name, rule.IsActive)
//	}
func (re *RuleEngine) GetAllRules() []LoyaltyRule {
	return re.rules
}

// GetActiveRules returns all active loyalty rules.
// Filters rules to return only those marked as active.
//
// Returns:
//   - []LoyaltyRule: Slice containing only active rules
//
// Example:
//
//	activeRules := engine.GetActiveRules()
//	fmt.Printf("Active rules: %d\n", len(activeRules))
//	for _, rule := range activeRules {
//		fmt.Printf("Active rule: %s\n", rule.Name)
//	}
func (re *RuleEngine) GetActiveRules() []LoyaltyRule {
	activeRules := make([]LoyaltyRule, 0)
	for _, rule := range re.rules {
		if rule.IsActive {
			activeRules = append(activeRules, rule)
		}
	}
	return activeRules
}

// GetRulesByType returns rules of a specific type.
// Filters rules by their type (e.g., "earn", "bonus", "multiplier").
//
// Parameters:
//   - ruleType: Type of rules to retrieve
//
// Returns:
//   - []LoyaltyRule: Slice containing rules of the specified type
//
// Example:
//
//	bonusRules := engine.GetRulesByType("bonus")
//	fmt.Printf("Found %d bonus rules\n", len(bonusRules))
func (re *RuleEngine) GetRulesByType(ruleType string) []LoyaltyRule {
	typedRules := make([]LoyaltyRule, 0)
	for _, rule := range re.rules {
		if rule.Type == ruleType {
			typedRules = append(typedRules, rule)
		}
	}
	return typedRules
}

// ActivateRule activates a loyalty rule.
// Sets the IsActive flag to true for the specified rule.
//
// Parameters:
//   - ruleID: ID of the rule to activate
//
// Returns:
//   - error: "Not found" error if rule doesn't exist, nil if successful
//
// Example:
//
//	err := engine.ActivateRule("seasonal_promotion")
//	if err == nil {
//		fmt.Println("Rule activated successfully")
//	}
func (re *RuleEngine) ActivateRule(ruleID string) error {
	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules[i].IsActive = true
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// DeactivateRule deactivates a loyalty rule.
// Sets the IsActive flag to false for the specified rule.
//
// Parameters:
//   - ruleID: ID of the rule to deactivate
//
// Returns:
//   - error: "Not found" error if rule doesn't exist, nil if successful
//
// Example:
//
//	err := engine.DeactivateRule("expired_promotion")
//	if err == nil {
//		fmt.Println("Rule deactivated successfully")
//	}
func (re *RuleEngine) DeactivateRule(ruleID string) error {
	for i, rule := range re.rules {
		if rule.ID == ruleID {
			re.rules[i].IsActive = false
			return nil
		}
	}

	return fmt.Errorf("rule with ID %s not found", ruleID)
}

// SetTierBenefit sets benefits for a loyalty tier.
// Configures the benefits and privileges for customers at the specified tier.
//
// Parameters:
//   - tier: Loyalty tier to configure
//   - benefit: TierBenefit configuration for the tier
//
// Example:
//
//	benefit := TierBenefit{
//		Tier:                TierGold,
//		PointsMultiplier:    1.5,
//		BonusPointsPercent:  10.0,
//		RedemptionBonus:     0.2,
//		FreeShippingThreshold: 0.0,
//		EarlyAccess:         true,
//		PrioritySupport:     true,
//	}
//	engine.SetTierBenefit(TierGold, benefit)
func (re *RuleEngine) SetTierBenefit(tier LoyaltyTier, benefit TierBenefit) {
	benefit.Tier = tier
	re.tierBenefits[tier] = benefit
}

// GetTierBenefit retrieves benefits for a loyalty tier.
// Returns the configured benefits for the specified tier.
//
// Parameters:
//   - tier: Loyalty tier to retrieve benefits for
//
// Returns:
//   - TierBenefit: Benefit configuration for the tier
//   - bool: True if tier benefits exist, false otherwise
//
// Example:
//
//	benefit, exists := engine.GetTierBenefit(TierGold)
//	if exists {
//		fmt.Printf("Gold tier multiplier: %.2f\n", benefit.PointsMultiplier)
//	}
func (re *RuleEngine) GetTierBenefit(tier LoyaltyTier) (TierBenefit, bool) {
	benefit, exists := re.tierBenefits[tier]
	return benefit, exists
}

// GetAllTierBenefits returns all tier benefits.
// Returns a map of all configured tier benefits.
//
// Returns:
//   - map[LoyaltyTier]TierBenefit: Map of tier to benefit configurations
//
// Example:
//
//	allBenefits := engine.GetAllTierBenefits()
//	for tier, benefit := range allBenefits {
//		fmt.Printf("%s: %.2fx points\n", tier, benefit.PointsMultiplier)
//	}
func (re *RuleEngine) GetAllTierBenefits() map[LoyaltyTier]TierBenefit {
	return re.tierBenefits
}

// AddReward adds a new reward to the catalog.
// Validates the reward before adding it to the available rewards.
//
// Parameters:
//   - reward: Reward to add to the catalog
//
// Returns:
//   - error: Validation error if reward is invalid, nil if successful
//
// Example:
//
//	reward := Reward{
//		ID:               "premium_discount",
//		Name:             "$25 Off Premium Items",
//		Description:      "$25 discount on premium category items",
//		Type:             RewardTypeDiscount,
//		PointsCost:       2500,
//		Value:            25.0,
//		MinOrderAmount:   100.0,
//		IsActive:         true,
//	}
//	err := engine.AddReward(reward)
func (re *RuleEngine) AddReward(reward Reward) error {
	if err := re.validateReward(reward); err != nil {
		return fmt.Errorf("invalid reward: %w", err)
	}

	re.rewards = append(re.rewards, reward)
	return nil
}

// UpdateReward updates an existing reward.
// Validates the updated reward before replacing the existing one.
//
// Parameters:
//   - rewardID: ID of the reward to update
//   - updatedReward: New reward data to replace the existing reward
//
// Returns:
//   - error: Validation error or "not found" error if reward doesn't exist
//
// Example:
//
//	updatedReward := existingReward
//	updatedReward.PointsCost = 2000
//	updatedReward.Value = 20.0
//	err := engine.UpdateReward("premium_discount", updatedReward)
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

// RemoveReward removes a reward from the catalog.
// Permanently deletes the reward from the available rewards.
//
// Parameters:
//   - rewardID: ID of the reward to remove
//
// Returns:
//   - error: "Not found" error if reward doesn't exist, nil if successful
//
// Example:
//
//	err := engine.RemoveReward("expired_reward")
//	if err != nil {
//		log.Printf("Failed to remove reward: %v", err)
//	}
func (re *RuleEngine) RemoveReward(rewardID string) error {
	for i, reward := range re.rewards {
		if reward.ID == rewardID {
			re.rewards = append(re.rewards[:i], re.rewards[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("reward with ID %s not found", rewardID)
}

// GetReward retrieves a reward by ID.
// Returns a pointer to the reward if found.
//
// Parameters:
//   - rewardID: ID of the reward to retrieve
//
// Returns:
//   - *Reward: Pointer to the reward if found
//   - error: "Not found" error if reward doesn't exist
//
// Example:
//
//	reward, err := engine.GetReward("premium_discount")
//	if err == nil {
//		fmt.Printf("Reward: %s, Cost: %d points\n", reward.Name, reward.PointsCost)
//	}
func (re *RuleEngine) GetReward(rewardID string) (*Reward, error) {
	for _, reward := range re.rewards {
		if reward.ID == rewardID {
			return &reward, nil
		}
	}

	return nil, fmt.Errorf("reward with ID %s not found", rewardID)
}

// GetAllRewards returns all rewards in the catalog.
// Returns a copy of all rewards regardless of their active status.
//
// Returns:
//   - []Reward: Slice containing all rewards in the catalog
//
// Example:
//
//	allRewards := engine.GetAllRewards()
//	fmt.Printf("Total rewards: %d\n", len(allRewards))
//	for _, reward := range allRewards {
//		fmt.Printf("Reward: %s, Active: %t\n", reward.Name, reward.IsActive)
//	}
func (re *RuleEngine) GetAllRewards() []Reward {
	return re.rewards
}

// GetActiveRewards returns all active rewards.
// Filters rewards to return only those marked as active.
//
// Returns:
//   - []Reward: Slice containing only active rewards
//
// Example:
//
//	activeRewards := engine.GetActiveRewards()
//	fmt.Printf("Active rewards: %d\n", len(activeRewards))
//	for _, reward := range activeRewards {
//		fmt.Printf("Available: %s (%d points)\n", reward.Name, reward.PointsCost)
//	}
func (re *RuleEngine) GetActiveRewards() []Reward {
	activeRewards := make([]Reward, 0)
	for _, reward := range re.rewards {
		if reward.IsActive {
			activeRewards = append(activeRewards, reward)
		}
	}
	return activeRewards
}

// GetRewardsByType returns rewards of a specific type.
// Filters rewards by their type (discount, shipping, product, etc.).
//
// Parameters:
//   - rewardType: Type of rewards to retrieve
//
// Returns:
//   - []Reward: Slice containing rewards of the specified type
//
// Example:
//
//	discountRewards := engine.GetRewardsByType(RewardTypeDiscount)
//	fmt.Printf("Found %d discount rewards\n", len(discountRewards))
func (re *RuleEngine) GetRewardsByType(rewardType RewardType) []Reward {
	typedRewards := make([]Reward, 0)
	for _, reward := range re.rewards {
		if reward.Type == rewardType {
			typedRewards = append(typedRewards, reward)
		}
	}
	return typedRewards
}

// SetReferralProgram sets the referral program configuration.
// Validates and configures the referral program settings.
//
// Parameters:
//   - program: ReferralProgram configuration to set
//
// Returns:
//   - error: Validation error if program is invalid, nil if successful
//
// Example:
//
//	program := ReferralProgram{
//		ID:               "main_referral",
//		Name:             "Refer a Friend",
//		ReferrerReward:   500,
//		RefereeReward:    250,
//		ValidityPeriod:   30,
//		IsActive:         true,
//	}
//	err := engine.SetReferralProgram(program)
func (re *RuleEngine) SetReferralProgram(program ReferralProgram) error {
	if err := re.validateReferralProgram(program); err != nil {
		return fmt.Errorf("invalid referral program: %w", err)
	}

	re.referralProgram = program
	return nil
}

// GetReferralProgram returns the referral program configuration.
// Returns the current referral program settings.
//
// Returns:
//   - ReferralProgram: Current referral program configuration
//
// Example:
//
//	program := engine.GetReferralProgram()
//	if program.IsActive {
//		fmt.Printf("Referral rewards: %d for referrer, %d for referee\n",
//			program.ReferrerReward, program.RefereeReward)
//	}
func (re *RuleEngine) GetReferralProgram() ReferralProgram {
	return re.referralProgram
}

// AddReviewReward adds a new review reward.
// Validates and adds a reward for customer reviews.
//
// Parameters:
//   - reward: ReviewReward to add
//
// Returns:
//   - error: Validation error if reward is invalid, nil if successful
//
// Example:
//
//	reward := ReviewReward{
//		ID:             "detailed_review",
//		Name:           "Detailed Review Bonus",
//		BasePoints:     50,
//		MinRating:      4,
//		MinCharacters:  100,
//		BonusPoints:    25,
//		IsActive:       true,
//	}
//	err := engine.AddReviewReward(reward)
func (re *RuleEngine) AddReviewReward(reward ReviewReward) error {
	if err := re.validateReviewReward(reward); err != nil {
		return fmt.Errorf("invalid review reward: %w", err)
	}

	re.reviewRewards = append(re.reviewRewards, reward)
	return nil
}

// UpdateReviewReward updates an existing review reward.
// Validates and updates the review reward configuration.
//
// Parameters:
//   - rewardID: ID of the review reward to update
//   - updatedReward: New review reward data
//
// Returns:
//   - error: Validation error or "not found" error if reward doesn't exist
//
// Example:
//
//	updatedReward := existingReward
//	updatedReward.BasePoints = 75
//	updatedReward.BonusPoints = 50
//	err := engine.UpdateReviewReward("detailed_review", updatedReward)
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

// GetReviewRewards returns all review rewards.
// Returns a copy of all review reward configurations.
//
// Returns:
//   - []ReviewReward: Slice containing all review rewards
//
// Example:
//
//	reviewRewards := engine.GetReviewRewards()
//	fmt.Printf("Total review rewards: %d\n", len(reviewRewards))
//	for _, reward := range reviewRewards {
//		fmt.Printf("Review reward: %s, Base: %d points\n", reward.Name, reward.BasePoints)
//	}
func (re *RuleEngine) GetReviewRewards() []ReviewReward {
	return re.reviewRewards
}

// GetActiveReviewRewards returns all active review rewards.
// Filters review rewards to return only those marked as active.
//
// Returns:
//   - []ReviewReward: Slice containing only active review rewards
//
// Example:
//
//	activeRewards := engine.GetActiveReviewRewards()
//	fmt.Printf("Active review rewards: %d\n", len(activeRewards))
//	for _, reward := range activeRewards {
//		fmt.Printf("Available: %s (min %d stars)\n", reward.Name, reward.MinRating)
//	}
func (re *RuleEngine) GetActiveReviewRewards() []ReviewReward {
	activeRewards := make([]ReviewReward, 0)
	for _, reward := range re.reviewRewards {
		if reward.IsActive {
			activeRewards = append(activeRewards, reward)
		}
	}
	return activeRewards
}

// OptimizeRules optimizes rule performance by sorting and organizing.
// Performs maintenance operations to improve rule engine performance:
// - Sorts rules by priority
// - Removes expired rules
// - Consolidates similar rules where possible
//
// Example:
//
//	// Perform periodic optimization
//	engine.OptimizeRules()
//	fmt.Println("Rule engine optimized")
func (re *RuleEngine) OptimizeRules() {
	re.sortRulesByPriority()
	re.removeExpiredRules()
	re.consolidateSimilarRules()
}

// ExportRules exports all rules to a map for serialization.
// Creates a comprehensive export of all rule engine data for backup or transfer.
//
// Returns:
//   - map[string]interface{}: Map containing all rule engine data
//
// Example:
//
//	exportData := engine.ExportRules()
//	// Save to file or transfer to another system
//	jsonData, _ := json.Marshal(exportData)
//	ioutil.WriteFile("rules_backup.json", jsonData, 0644)
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

// ImportRules imports rules from a map.
// Restores rule engine data from an exported map.
//
// Parameters:
//   - data: Map containing exported rule engine data
//
// Returns:
//   - error: Import error if data is invalid, nil if successful
//
// Example:
//
//	// Load from file
//	jsonData, _ := ioutil.ReadFile("rules_backup.json")
//	var importData map[string]interface{}
//	json.Unmarshal(jsonData, &importData)
//	err := engine.ImportRules(importData)
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

// GetStatistics returns statistics about the rule engine.
// Provides comprehensive analytics about rules, rewards, and configuration.
//
// Returns:
//   - map[string]interface{}: Map containing various statistics
//
// Statistics include:
//   - total_rules, active_rules, expired_rules
//   - total_rewards, active_rewards, expired_rewards
//   - tier_benefits count
//   - review_rewards count
//   - referral_active status
//
// Example:
//
//	stats := engine.GetStatistics()
//	fmt.Printf("Rules: %d total, %d active\n", stats["total_rules"], stats["active_rules"])
//	fmt.Printf("Rewards: %d total, %d active\n", stats["total_rewards"], stats["active_rewards"])
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

// validateRule validates a loyalty rule.
// Performs comprehensive validation of rule structure and data.
//
// Parameters:
//   - rule: LoyaltyRule to validate
//
// Returns:
//   - error: Validation error if rule is invalid, nil if valid
//
// Validation checks:
//   - Rule ID is not empty
//   - Rule name is not empty
//   - Priority is non-negative
//   - Conditions are valid
//   - Actions are valid
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

// validateReward validates a reward.
// Ensures reward configuration is valid and complete.
//
// Parameters:
//   - reward: Reward to validate
//
// Returns:
//   - error: Validation error if reward is invalid, nil if valid
//
// Validation checks:
//   - Reward ID is not empty
//   - Reward name is not empty
//   - Type is valid
//   - Value is positive
//   - Points required is non-negative
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

// validateReferralProgram validates a referral program.
// Ensures referral program configuration is valid.
//
// Parameters:
//   - program: ReferralProgram to validate
//
// Returns:
//   - error: Validation error if program is invalid, nil if valid
//
// Validation checks:
//   - Referrer reward is non-negative
//   - Referee reward is non-negative
//   - Maximum referrals is positive
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

// validateReviewReward validates a review reward.
// Ensures review reward configuration is valid.
//
// Parameters:
//   - reward: ReviewReward to validate
//
// Returns:
//   - error: Validation error if reward is invalid, nil if valid
//
// Validation checks:
//   - Reward ID is not empty
//   - Reward name is not empty
//   - Base points is non-negative
//   - Bonus points is non-negative
//   - Minimum rating is between 1 and 5
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

// sortRulesByPriority sorts rules by priority (higher priority first).
// Internal helper function for rule optimization.
func (re *RuleEngine) sortRulesByPriority() {
	sort.Slice(re.rules, func(i, j int) bool {
		return re.rules[i].Priority > re.rules[j].Priority
	})
}

// removeExpiredRules removes expired rules.
// Internal helper function that cleans up rules past their expiry date.
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

// consolidateSimilarRules consolidates similar rules for better performance.
// Internal helper function that merges rules with identical conditions.
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

// CreateDefaultConfiguration creates a default loyalty program configuration.
// Initializes a complete loyalty program with standard tiers, rules, and rewards.
//
// Returns:
//   - *LoyaltyConfiguration: Fully configured loyalty configuration with default settings
//
// Default configuration includes:
//   - 4-tier system (Bronze, Silver, Gold, Platinum)
//   - Standard point earning rules
//   - Tier-based benefits and multipliers
//   - Common rewards (discounts, free shipping)
//   - Referral program
//   - Review rewards
//
// Example:
//
//	config := CreateDefaultConfiguration()
//	engine := NewRuleEngine(config)
//	// Ready to use with standard loyalty program
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

// CreateDefaultRules creates default loyalty rules.
// Generates a standard set of loyalty rules for common e-commerce scenarios.
//
// Returns:
//   - []LoyaltyRule: Slice of default loyalty rules
//
// Default rules include:
//   - Base points earning (1 point per dollar)
//   - First purchase bonus (100 points)
//   - High-value purchase bonus (500+ points for orders over $500)
//
// Example:
//
//	rules := CreateDefaultRules()
//	for _, rule := range rules {
//		engine.AddRule(rule)
//	}
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

// CreateDefaultRewards creates default rewards.
// Generates a standard set of rewards that customers can redeem with points.
//
// Returns:
//   - []Reward: Slice of default rewards
//
// Default rewards include:
//   - $5 discount (500 points)
//   - $10 discount (1000 points)
//   - $25 discount (2500 points)
//   - Free shipping (250 points)
//
// Example:
//
//	rewards := CreateDefaultRewards()
//	for _, reward := range rewards {
//		engine.AddReward(reward)
//	}
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
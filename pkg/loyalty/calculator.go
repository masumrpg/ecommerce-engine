package loyalty

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
)

// Calculator handles loyalty points calculations
type Calculator struct {
	config *LoyaltyConfiguration
	rules  []LoyaltyRule
	tierBenefits map[LoyaltyTier]TierBenefit
}

// NewCalculator creates a new loyalty calculator
func NewCalculator(config *LoyaltyConfiguration) *Calculator {
	return &Calculator{
		config: config,
		rules:  config.DefaultRules,
		tierBenefits: config.TierBenefits,
	}
}

// Calculate calculates loyalty points for a purchase
func (c *Calculator) Calculate(input PointsCalculationInput) (*PointsCalculationResult, error) {
	if err := c.validateInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	result := &PointsCalculationResult{
		CustomerID: input.Customer.ID,
		IsValid:    true,
	}

	// Calculate base points
	basePoints := c.calculateBasePoints(input)
	result.BasePoints = basePoints
	result.PointsBreakdown = append(result.PointsBreakdown, PointsBreakdown{
		Source:      "base",
		Description: "Base points from purchase",
		Amount:      input.OrderAmount,
		Rate:        c.config.BasePointsRate,
		Multiplier:  1.0,
		Points:      basePoints,
		PointsType:  PointsTypeBase,
	})

	// Apply tier benefits
	tierBenefit := c.getTierBenefit(input.Customer.Tier)
	tierMultiplier := tierBenefit.PointsMultiplier
	tierBonusPoints := int(math.Round(float64(basePoints) * (tierMultiplier - 1.0)))
	if tierBonusPoints > 0 {
		result.BonusPoints += tierBonusPoints
		result.PointsBreakdown = append(result.PointsBreakdown, PointsBreakdown{
			Source:      "tier_bonus",
			Description: fmt.Sprintf("%s tier multiplier", input.Customer.Tier),
			Amount:      input.OrderAmount,
			Rate:        c.config.BasePointsRate,
			Multiplier:  tierMultiplier,
			Points:      tierBonusPoints,
			PointsType:  PointsTypeBonus,
		})
	}

	// Apply loyalty rules
	applicableRules := c.getApplicableRules(input)
	for _, rule := range applicableRules {
		bonusPoints, breakdown, appliedRule := c.applyRule(rule, input, basePoints)
		if bonusPoints > 0 {
			result.BonusPoints += bonusPoints
			result.PointsBreakdown = append(result.PointsBreakdown, breakdown...)
			result.AppliedRules = append(result.AppliedRules, appliedRule)
		}
	}

	// Calculate total points
	result.TotalPoints = result.BasePoints + result.BonusPoints
	result.NewBalance = input.Customer.CurrentPoints + result.TotalPoints

	// Set expiry date
	result.ExpiryDate = c.calculateExpiryDate(input.Customer.Tier)

	// Update tier information
	result.TierInfo = c.calculateTierInfo(input.Customer, input.OrderAmount)

	// Create transactions
	result.Transactions = c.createTransactions(input, result)

	// Generate recommendations
	result.Recommendations = c.generateRecommendations(input.Customer, result)

	return result, nil
}

// RedeemPoints redeems points for rewards
func (c *Calculator) RedeemPoints(input RedemptionInput, reward Reward) (*RedemptionResult, error) {
	if err := c.validateRedemptionInput(input, reward); err != nil {
		return nil, fmt.Errorf("invalid redemption input: %w", err)
	}

	quantity := input.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	totalPointsCost := reward.PointsCost * quantity
	if input.Customer.CurrentPoints < totalPointsCost {
		return &RedemptionResult{
			CustomerID:   input.Customer.ID,
			RewardID:     input.RewardID,
			IsSuccessful: false,
			Errors:       []string{"Insufficient points balance"},
		}, nil
	}

	// Apply tier redemption bonus
	tierBenefit := c.getTierBenefit(input.Customer.Tier)
	discountAmount := reward.Value * float64(quantity)
	if tierBenefit.RedemptionBonus > 0 {
		discountAmount *= (1.0 + tierBenefit.RedemptionBonus)
	}

	// Create redemption transaction
	transaction := PointsTransaction{
		ID:          c.generateTransactionID(),
		CustomerID:  input.Customer.ID,
		Type:        TransactionTypeRedeem,
		PointsType:  PointsTypeBase,
		Amount:      -totalPointsCost,
		Balance:     input.Customer.CurrentPoints - totalPointsCost,
		RewardID:    input.RewardID,
		Description: fmt.Sprintf("Redeemed %s (x%d)", reward.Name, quantity),
		Timestamp:   input.Timestamp,
		Source:      "redemption",
		Channel:     input.Channel,
		Metadata:    input.Metadata,
	}

	result := &RedemptionResult{
		CustomerID:     input.Customer.ID,
		RewardID:       input.RewardID,
		RewardName:     reward.Name,
		PointsRedeemed: totalPointsCost,
		DiscountAmount: discountAmount,
		NewBalance:     transaction.Balance,
		RedemptionCode: c.generateRedemptionCode(),
		ValidUntil:     reward.ValidUntil,
		Transaction:    transaction,
		IsSuccessful:   true,
	}

	return result, nil
}

// CalculateReferralReward calculates referral rewards
func (c *Calculator) CalculateReferralReward(referrer Customer, referee Customer, program ReferralProgram, orderAmount float64) (*PointsCalculationResult, error) {
	if !program.IsActive {
		return nil, fmt.Errorf("referral program is not active")
	}

	if program.MinOrderAmount > 0 && orderAmount < program.MinOrderAmount {
		return nil, fmt.Errorf("order amount below minimum threshold")
	}

	result := &PointsCalculationResult{
		CustomerID: referrer.ID,
		IsValid:    true,
	}

	// Referrer reward
	result.BonusPoints = program.ReferrerReward
	result.TotalPoints = program.ReferrerReward
	result.NewBalance = referrer.CurrentPoints + result.TotalPoints

	result.PointsBreakdown = append(result.PointsBreakdown, PointsBreakdown{
		Source:      "referral",
		Description: fmt.Sprintf("Referral reward for %s", referee.Email),
		Amount:      orderAmount,
		Rate:        0,
		Multiplier:  1.0,
		Points:      program.ReferrerReward,
		PointsType:  PointsTypeReferral,
	})

	// Create transaction
	transaction := PointsTransaction{
		ID:          c.generateTransactionID(),
		CustomerID:  referrer.ID,
		Type:        TransactionTypeEarn,
		PointsType:  PointsTypeReferral,
		Amount:      program.ReferrerReward,
		Balance:     result.NewBalance,
		Description: fmt.Sprintf("Referral reward for %s", referee.Email),
		Timestamp:   time.Now(),
		Source:      "referral",
		Metadata:    map[string]interface{}{"referee_id": referee.ID, "program_id": program.ID},
	}

	result.Transactions = []PointsTransaction{transaction}

	return result, nil
}

// CalculateReviewReward calculates review rewards
func (c *Calculator) CalculateReviewReward(customer Customer, reward ReviewReward, hasPhoto, hasVideo, isVerified bool, rating int, reviewLength int) (*PointsCalculationResult, error) {
	if !reward.IsActive {
		return nil, fmt.Errorf("review reward is not active")
	}

	if rating < reward.MinRating {
		return nil, fmt.Errorf("rating below minimum threshold")
	}

	if reviewLength < reward.MinCharacters {
		return nil, fmt.Errorf("review too short")
	}

	result := &PointsCalculationResult{
		CustomerID: customer.ID,
		IsValid:    true,
	}

	// Base review points
	totalPoints := reward.BasePoints
	breakdown := []PointsBreakdown{
		{
			Source:      "review_base",
			Description: "Base review points",
			Amount:      0,
			Rate:        0,
			Multiplier:  1.0,
			Points:      reward.BasePoints,
			PointsType:  PointsTypeReview,
		},
	}

	// Photo bonus
	if hasPhoto && reward.PhotoBonus > 0 {
		totalPoints += reward.PhotoBonus
		breakdown = append(breakdown, PointsBreakdown{
			Source:      "review_photo",
			Description: "Photo bonus",
			Amount:      0,
			Rate:        0,
			Multiplier:  1.0,
			Points:      reward.PhotoBonus,
			PointsType:  PointsTypeReview,
		})
	}

	// Video bonus
	if hasVideo && reward.VideoBonus > 0 {
		totalPoints += reward.VideoBonus
		breakdown = append(breakdown, PointsBreakdown{
			Source:      "review_video",
			Description: "Video bonus",
			Amount:      0,
			Rate:        0,
			Multiplier:  1.0,
			Points:      reward.VideoBonus,
			PointsType:  PointsTypeReview,
		})
	}

	// Verified purchase bonus
	if isVerified && reward.VerifiedBonus > 0 {
		totalPoints += reward.VerifiedBonus
		breakdown = append(breakdown, PointsBreakdown{
			Source:      "review_verified",
			Description: "Verified purchase bonus",
			Amount:      0,
			Rate:        0,
			Multiplier:  1.0,
			Points:      reward.VerifiedBonus,
			PointsType:  PointsTypeReview,
		})
	}

	result.BonusPoints = totalPoints - reward.BasePoints
	result.BasePoints = reward.BasePoints
	result.TotalPoints = totalPoints
	result.NewBalance = customer.CurrentPoints + totalPoints
	result.PointsBreakdown = breakdown

	// Create transaction
	transaction := PointsTransaction{
		ID:          c.generateTransactionID(),
		CustomerID:  customer.ID,
		Type:        TransactionTypeEarn,
		PointsType:  PointsTypeReview,
		Amount:      totalPoints,
		Balance:     result.NewBalance,
		Description: "Review reward",
		Timestamp:   time.Now(),
		Source:      "review",
		Metadata:    map[string]interface{}{"reward_id": reward.ID, "rating": rating},
	}

	result.Transactions = []PointsTransaction{transaction}

	return result, nil
}

// GetAvailableRewards returns available rewards for a customer
func (c *Calculator) GetAvailableRewards(customer Customer, rewards []Reward) []Reward {
	available := make([]Reward, 0)

	for _, reward := range rewards {
		if c.isRewardAvailable(customer, reward) {
			available = append(available, reward)
		}
	}

	// Sort by points cost
	sort.Slice(available, func(i, j int) bool {
		return available[i].PointsCost < available[j].PointsCost
	})

	return available
}

// calculateBasePoints calculates base points from purchase amount
func (c *Calculator) calculateBasePoints(input PointsCalculationInput) int {
	return int(math.Floor(input.OrderAmount * c.config.BasePointsRate))
}

// getTierBenefit returns tier benefits for a given tier
func (c *Calculator) getTierBenefit(tier LoyaltyTier) TierBenefit {
	if benefit, exists := c.tierBenefits[tier]; exists {
		return benefit
	}
	return TierBenefit{PointsMultiplier: 1.0} // Default
}

// getApplicableRules returns rules applicable to the input
func (c *Calculator) getApplicableRules(input PointsCalculationInput) []LoyaltyRule {
	applicable := make([]LoyaltyRule, 0)

	for _, rule := range c.rules {
		if c.isRuleApplicable(rule, input) {
			applicable = append(applicable, rule)
		}
	}

	// Sort by priority
	sort.Slice(applicable, func(i, j int) bool {
		return applicable[i].Priority > applicable[j].Priority
	})

	return applicable
}

// isRuleApplicable checks if a rule is applicable to the input
func (c *Calculator) isRuleApplicable(rule LoyaltyRule, input PointsCalculationInput) bool {
	if !rule.IsActive {
		return false
	}

	// Check time validity
	if !rule.ValidFrom.IsZero() && input.Timestamp.Before(rule.ValidFrom) {
		return false
	}
	if !rule.ValidUntil.IsZero() && input.Timestamp.After(rule.ValidUntil) {
		return false
	}

	// Check tier applicability
	if len(rule.ApplicableTiers) > 0 {
		tierApplicable := false
		for _, tier := range rule.ApplicableTiers {
			if tier == input.Customer.Tier {
				tierApplicable = true
				break
			}
		}
		if !tierApplicable {
			return false
		}
	}

	// Check payment method
	if len(rule.PaymentMethods) > 0 && input.PaymentMethod != "" {
		paymentApplicable := false
		for _, method := range rule.PaymentMethods {
			if method == input.PaymentMethod {
				paymentApplicable = true
				break
			}
		}
		if !paymentApplicable {
			return false
		}
	}

	// Check channel
	if len(rule.Channels) > 0 && input.Channel != "" {
		channelApplicable := false
		for _, channel := range rule.Channels {
			if channel == input.Channel {
				channelApplicable = true
				break
			}
		}
		if !channelApplicable {
			return false
		}
	}

	// Evaluate conditions
	for _, condition := range rule.Conditions {
		if !c.evaluateCondition(condition, input) {
			return false
		}
	}

	return true
}

// evaluateCondition evaluates a loyalty condition
func (c *Calculator) evaluateCondition(condition LoyaltyCondition, input PointsCalculationInput) bool {
	switch condition.Type {
	case "amount":
		return c.compareValues(input.OrderAmount, condition.Operator, condition.Value)
	case "quantity":
		totalQuantity := 0
		for _, item := range input.Items {
			totalQuantity += item.Quantity
		}
		return c.compareValues(float64(totalQuantity), condition.Operator, condition.Value)
	case "category":
		for _, item := range input.Items {
			if c.compareStringValues(item.Category, condition.Operator, condition.Value) {
				return true
			}
		}
		return false
	case "payment_method":
		return c.compareStringValues(input.PaymentMethod, condition.Operator, condition.Value)
	case "first_purchase":
		expected, _ := condition.Value.(bool)
		return input.IsFirstPurchase == expected
	case "tier":
		return c.compareStringValues(string(input.Customer.Tier), condition.Operator, condition.Value)
	default:
		return true
	}
}

// compareValues compares numeric values based on operator
func (c *Calculator) compareValues(actual float64, operator string, expected interface{}) bool {
	expectedFloat, err := c.toFloat64(expected)
	if err != nil {
		return false
	}

	switch operator {
	case ">":
		return actual > expectedFloat
	case ">=":
		return actual >= expectedFloat
	case "<":
		return actual < expectedFloat
	case "<=":
		return actual <= expectedFloat
	case "=":
		return actual == expectedFloat
	case "!=":
		return actual != expectedFloat
	default:
		return false
	}
}

// compareStringValues compares string values based on operator
func (c *Calculator) compareStringValues(actual string, operator string, expected interface{}) bool {
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
		if expectedSlice, ok := expected.([]string); ok {
			for _, val := range expectedSlice {
				if actual == val {
					return true
				}
			}
		}
		return false
	default:
		return false
	}
}

// toFloat64 converts interface{} to float64
func (c *Calculator) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// applyRule applies a loyalty rule and returns bonus points
func (c *Calculator) applyRule(rule LoyaltyRule, input PointsCalculationInput, basePoints int) (int, []PointsBreakdown, AppliedLoyaltyRule) {
	bonusPoints := 0
	breakdown := make([]PointsBreakdown, 0)
	appliedRule := AppliedLoyaltyRule{
		RuleID:      rule.ID,
		Name:        rule.Name,
		Type:        rule.Type,
		Description: rule.Description,
	}

	for _, action := range rule.Actions {
		switch action.Type {
		case "earn_points":
			if points, ok := action.Value.(int); ok {
				bonusPoints += points
				appliedRule.PointsAwarded += points
				breakdown = append(breakdown, PointsBreakdown{
					Source:      "rule_bonus",
					Description: action.Description,
					Amount:      input.OrderAmount,
					Rate:        0,
					Multiplier:  1.0,
					Points:      points,
					PointsType:  action.PointsType,
					RuleID:      rule.ID,
				})
			}
		case "multiply_points":
			if multiplier, ok := action.Value.(float64); ok {
				additionalPoints := int(float64(basePoints) * (multiplier - 1.0))
				bonusPoints += additionalPoints
				appliedRule.PointsAwarded += additionalPoints
				appliedRule.Multiplier = multiplier
				breakdown = append(breakdown, PointsBreakdown{
					Source:      "rule_multiplier",
					Description: action.Description,
					Amount:      input.OrderAmount,
					Rate:        c.config.BasePointsRate,
					Multiplier:  multiplier,
					Points:      additionalPoints,
					PointsType:  action.PointsType,
					RuleID:      rule.ID,
				})
			}
		case "bonus_points":
			if rate, ok := action.Value.(float64); ok {
				additionalPoints := int(input.OrderAmount * rate)
				bonusPoints += additionalPoints
				appliedRule.PointsAwarded += additionalPoints
				breakdown = append(breakdown, PointsBreakdown{
					Source:      "rule_bonus_rate",
					Description: action.Description,
					Amount:      input.OrderAmount,
					Rate:        rate,
					Multiplier:  1.0,
					Points:      additionalPoints,
					PointsType:  action.PointsType,
					RuleID:      rule.ID,
				})
			}
		}
	}

	return bonusPoints, breakdown, appliedRule
}

// calculateExpiryDate calculates points expiry date based on tier
func (c *Calculator) calculateExpiryDate(tier LoyaltyTier) time.Time {
	tierBenefit := c.getTierBenefit(tier)
	months := tierBenefit.MaxPointsExpiry
	if months <= 0 {
		months = c.config.PointsExpiry
	}
	return time.Now().AddDate(0, months, 0)
}

// calculateTierInfo calculates tier information and progress
func (c *Calculator) calculateTierInfo(customer Customer, orderAmount float64) TierInfo {
	newSpend := customer.AnnualSpend + orderAmount
	currentTier := customer.Tier
	nextTier := c.getNextTier(currentTier)
	nextThreshold := c.getTierThreshold(nextTier)

	tierInfo := TierInfo{
		CurrentTier:       currentTier,
		NextTier:          nextTier,
		CurrentSpend:      newSpend,
		NextTierThreshold: nextThreshold,
		SpendToNextTier:   math.Max(0, nextThreshold-newSpend),
		ProgressPercent:   math.Min(100, (newSpend/nextThreshold)*100),
		TierAchievedDate:  customer.TierAchievedDate,
		Benefits:          c.getTierBenefit(currentTier),
		IsUpgraded:        false,
	}

	// Check for tier upgrade
	if newSpend >= nextThreshold && nextTier != currentTier {
		tierInfo.CurrentTier = nextTier
		tierInfo.IsUpgraded = true
		tierInfo.TierAchievedDate = time.Now()
		tierInfo.Benefits = c.getTierBenefit(nextTier)
		
		// Update next tier info
		newNextTier := c.getNextTier(nextTier)
		tierInfo.NextTier = newNextTier
		tierInfo.NextTierThreshold = c.getTierThreshold(newNextTier)
		tierInfo.SpendToNextTier = math.Max(0, tierInfo.NextTierThreshold-newSpend)
		tierInfo.ProgressPercent = math.Min(100, (newSpend/tierInfo.NextTierThreshold)*100)
	}

	return tierInfo
}

// getNextTier returns the next tier for a given tier
func (c *Calculator) getNextTier(current LoyaltyTier) LoyaltyTier {
	switch current {
	case TierBronze:
		return TierSilver
	case TierSilver:
		return TierGold
	case TierGold:
		return TierPlatinum
	default:
		return current // Already at highest tier
	}
}

// getTierThreshold returns the spending threshold for a tier
func (c *Calculator) getTierThreshold(tier LoyaltyTier) float64 {
	if threshold, exists := c.config.TierThresholds[tier]; exists {
		return threshold
	}
	return math.Inf(1) // No threshold (highest tier)
}

// createTransactions creates point transactions
func (c *Calculator) createTransactions(input PointsCalculationInput, result *PointsCalculationResult) []PointsTransaction {
	transactions := make([]PointsTransaction, 0)

	if result.TotalPoints > 0 {
		transaction := PointsTransaction{
			ID:          c.generateTransactionID(),
			CustomerID:  input.Customer.ID,
			Type:        TransactionTypeEarn,
			PointsType:  PointsTypeBase,
			Amount:      result.TotalPoints,
			Balance:     result.NewBalance,
			OrderID:     input.OrderID,
			Description: "Points earned from purchase",
			Timestamp:   input.Timestamp,
			ExpiryDate:  result.ExpiryDate,
			Source:      "purchase",
			Channel:     input.Channel,
			Metadata:    input.Metadata,
		}
		transactions = append(transactions, transaction)
	}

	return transactions
}

// generateRecommendations generates loyalty recommendations
func (c *Calculator) generateRecommendations(customer Customer, result *PointsCalculationResult) []LoyaltyRecommendation {
	recommendations := make([]LoyaltyRecommendation, 0)

	// Tier upgrade recommendation
	if result.TierInfo.SpendToNextTier > 0 && result.TierInfo.SpendToNextTier <= 1000 {
		recommendations = append(recommendations, LoyaltyRecommendation{
			Type:        "tier_upgrade",
			Title:       fmt.Sprintf("Upgrade to %s tier", result.TierInfo.NextTier),
			Description: fmt.Sprintf("Spend $%.2f more to reach %s tier", result.TierInfo.SpendToNextTier, result.TierInfo.NextTier),
			ActionText:  "Shop now",
			Value:       result.TierInfo.SpendToNextTier,
			Priority:    1,
		})
	}

	// Points redemption recommendation
	if customer.CurrentPoints >= c.config.MinRedemption {
		recommendations = append(recommendations, LoyaltyRecommendation{
			Type:        "reward",
			Title:       "Redeem your points",
			Description: fmt.Sprintf("You have %d points available for redemption", customer.CurrentPoints),
			ActionText:  "View rewards",
			Value:       float64(customer.CurrentPoints),
			Priority:    2,
		})
	}

	return recommendations
}

// isRewardAvailable checks if a reward is available for a customer
func (c *Calculator) isRewardAvailable(customer Customer, reward Reward) bool {
	if !reward.IsActive {
		return false
	}

	if customer.CurrentPoints < reward.PointsCost {
		return false
	}

	if reward.RequiredTier != "" && customer.Tier != reward.RequiredTier {
		return false
	}

	if !reward.ValidFrom.IsZero() && time.Now().Before(reward.ValidFrom) {
		return false
	}

	if !reward.ValidUntil.IsZero() && time.Now().After(reward.ValidUntil) {
		return false
	}

	if reward.Stock > 0 && reward.Stock <= 0 {
		return false
	}

	return true
}

// validateInput validates points calculation input
func (c *Calculator) validateInput(input PointsCalculationInput) error {
	if input.Customer.ID == "" {
		return fmt.Errorf("customer ID is required")
	}

	if input.OrderAmount < 0 {
		return fmt.Errorf("order amount cannot be negative")
	}

	if input.Timestamp.IsZero() {
		return fmt.Errorf("timestamp is required")
	}

	return nil
}

// validateRedemptionInput validates redemption input
func (c *Calculator) validateRedemptionInput(input RedemptionInput, reward Reward) error {
	if input.Customer.ID == "" {
		return fmt.Errorf("customer ID is required")
	}

	if input.RewardID == "" {
		return fmt.Errorf("reward ID is required")
	}

	if input.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	// Check basic reward availability (excluding points balance)
	if !reward.IsActive {
		return fmt.Errorf("reward is not active")
	}

	if reward.RequiredTier != "" && input.Customer.Tier != reward.RequiredTier {
		return fmt.Errorf("customer tier does not meet requirement")
	}

	if !reward.ValidFrom.IsZero() && time.Now().Before(reward.ValidFrom) {
		return fmt.Errorf("reward is not yet valid")
	}

	if !reward.ValidUntil.IsZero() && time.Now().After(reward.ValidUntil) {
		return fmt.Errorf("reward has expired")
	}

	if reward.Stock > 0 && reward.Stock <= 0 {
		return fmt.Errorf("reward is out of stock")
	}

	return nil
}

// generateTransactionID generates a unique transaction ID
func (c *Calculator) generateTransactionID() string {
	return fmt.Sprintf("txn_%d", time.Now().UnixNano())
}

// generateRedemptionCode generates a unique redemption code
func (c *Calculator) generateRedemptionCode() string {
	return fmt.Sprintf("RED%d", time.Now().UnixNano()%1000000)
}
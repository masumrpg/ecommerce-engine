package loyalty

import (
	"testing"
	"time"
)

func TestNewCalculator(t *testing.T) {
	config := &LoyaltyConfiguration{
		ProgramName:    "Test Program",
		BaseCurrency:   "USD",
		BasePointsRate: 1.0,
		RedemptionRate: 0.01,
		PointsExpiry:   12,
		MinRedemption:  100,
		TierThresholds: map[LoyaltyTier]float64{
			TierBronze:   0,
			TierSilver:   1000,
			TierGold:     5000,
			TierPlatinum: 15000,
		},
		TierBenefits: map[LoyaltyTier]TierBenefit{
			TierBronze: {
				Tier:             TierBronze,
				PointsMultiplier: 1.0,
				BonusPointsPercent: 0,
				RedemptionBonus:  0,
				MaxPointsExpiry:  12,
			},
			TierSilver: {
				Tier:             TierSilver,
				PointsMultiplier: 1.2,
				BonusPointsPercent: 5,
				RedemptionBonus:  0.1,
				MaxPointsExpiry:  18,
			},
			TierGold: {
				Tier:             TierGold,
				PointsMultiplier: 1.5,
				BonusPointsPercent: 10,
				RedemptionBonus:  0.2,
				MaxPointsExpiry:  24,
			},
			TierPlatinum: {
				Tier:             TierPlatinum,
				PointsMultiplier: 2.0,
				BonusPointsPercent: 15,
				RedemptionBonus:  0.3,
				MaxPointsExpiry:  36,
			},
		},
		IsActive: true,
	}
	
	calc := NewCalculator(config)
	
	if calc == nil {
		t.Fatal("NewCalculator should not return nil")
	}
	
	if calc.config != config {
		t.Error("Calculator config not set correctly")
	}
	
	if len(calc.tierBenefits) != len(config.TierBenefits) {
		t.Error("Tier benefits not set correctly")
	}
}

func TestCalculate(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	t.Run("BasicCalculation", func(t *testing.T) {
		customer := Customer{
			ID:            "customer1",
			Tier:          TierBronze,
			CurrentPoints: 100,
			AnnualSpend:   500,
		}
		
		input := PointsCalculationInput{
			Customer:    customer,
			OrderAmount: 100.0,
			Timestamp:   time.Now(),
			OrderID:     "order1",
		}
		
		result, err := calc.Calculate(input)
		if err != nil {
			t.Fatalf("Calculate failed: %v", err)
		}
		
		if !result.IsValid {
			t.Error("Result should be valid")
		}
		
		if result.CustomerID != customer.ID {
			t.Error("Customer ID mismatch")
		}
		
		expectedBasePoints := 100 // 100 * 1.0 rate
		if result.BasePoints != expectedBasePoints {
			t.Errorf("Expected base points %d, got %d", expectedBasePoints, result.BasePoints)
		}
		
		if result.NewBalance != customer.CurrentPoints+result.TotalPoints {
			t.Error("New balance calculation incorrect")
		}
		
		if len(result.PointsBreakdown) == 0 {
			t.Error("Points breakdown should not be empty")
		}
		
		if len(result.Transactions) == 0 {
			t.Error("Transactions should not be empty")
		}
	})
	
	t.Run("SilverTierMultiplier", func(t *testing.T) {
		customer := Customer{
			ID:            "customer2",
			Tier:          TierSilver,
			CurrentPoints: 200,
			AnnualSpend:   2000,
		}
		
		input := PointsCalculationInput{
			Customer:    customer,
			OrderAmount: 100.0,
			Timestamp:   time.Now(),
		}
		
		result, err := calc.Calculate(input)
		if err != nil {
			t.Fatalf("Calculate failed: %v", err)
		}
		
		expectedBasePoints := 100
		expectedBonusPoints := 20 // 100 * (1.2 - 1.0)
		
		if result.BasePoints != expectedBasePoints {
			t.Errorf("Expected base points %d, got %d", expectedBasePoints, result.BasePoints)
		}
		
		if result.BonusPoints < expectedBonusPoints {
			t.Errorf("Expected at least %d bonus points, got %d", expectedBonusPoints, result.BonusPoints)
		}
	})
	
	t.Run("TierUpgrade", func(t *testing.T) {
		customer := Customer{
			ID:            "customer3",
			Tier:          TierBronze,
			CurrentPoints: 100,
			AnnualSpend:   950, // Close to silver threshold
		}
		
		input := PointsCalculationInput{
			Customer:    customer,
			OrderAmount: 100.0, // This should trigger tier upgrade
			Timestamp:   time.Now(),
		}
		
		result, err := calc.Calculate(input)
		if err != nil {
			t.Fatalf("Calculate failed: %v", err)
		}
		
		if result.TierInfo.CurrentTier != TierSilver {
			t.Errorf("Expected tier upgrade to Silver, got %s", result.TierInfo.CurrentTier)
		}
		
		if !result.TierInfo.IsUpgraded {
			t.Error("IsUpgraded should be true")
		}
	})
	
	t.Run("InvalidInput", func(t *testing.T) {
		input := PointsCalculationInput{
			Customer: Customer{ID: ""}, // Invalid customer ID
			OrderAmount: 100.0,
			Timestamp:   time.Now(),
		}
		
		_, err := calc.Calculate(input)
		if err == nil {
			t.Error("Expected error for invalid input")
		}
	})
	
	t.Run("NegativeOrderAmount", func(t *testing.T) {
		customer := Customer{
			ID:            "customer4",
			Tier:          TierBronze,
			CurrentPoints: 100,
		}
		
		input := PointsCalculationInput{
			Customer:    customer,
			OrderAmount: -50.0, // Negative amount
			Timestamp:   time.Now(),
		}
		
		_, err := calc.Calculate(input)
		if err == nil {
			t.Error("Expected error for negative order amount")
		}
	})
}

func TestRedeemPoints(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	t.Run("SuccessfulRedemption", func(t *testing.T) {
		customer := Customer{
			ID:            "customer1",
			Tier:          TierBronze,
			CurrentPoints: 500,
		}
		
		reward := Reward{
			ID:         "reward1",
			Name:       "$5 Discount",
			Type:       RewardTypeDiscount,
			PointsCost: 100,
			Value:      5.0,
			IsActive:   true,
			ValidFrom:  time.Now().Add(-time.Hour),
			ValidUntil: time.Now().Add(time.Hour),
		}
		
		input := RedemptionInput{
			Customer:  customer,
			RewardID:  reward.ID,
			Quantity:  1,
			Timestamp: time.Now(),
		}
		
		result, err := calc.RedeemPoints(input, reward)
		if err != nil {
			t.Fatalf("RedeemPoints failed: %v", err)
		}
		
		if !result.IsSuccessful {
			t.Error("Redemption should be successful")
		}
		
		if result.PointsRedeemed != reward.PointsCost {
			t.Errorf("Expected %d points redeemed, got %d", reward.PointsCost, result.PointsRedeemed)
		}
		
		expectedBalance := customer.CurrentPoints - reward.PointsCost
		if result.NewBalance != expectedBalance {
			t.Errorf("Expected balance %d, got %d", expectedBalance, result.NewBalance)
		}
		
		if result.RedemptionCode == "" {
			t.Error("Redemption code should not be empty")
		}
	})
	
	t.Run("InsufficientPoints", func(t *testing.T) {
		customer := Customer{
			ID:            "customer2",
			Tier:          TierBronze,
			CurrentPoints: 50, // Not enough points
		}
		
		reward := Reward{
			ID:         "reward1",
			Name:       "$5 Discount",
			PointsCost: 100,
			Value:      5.0,
			IsActive:   true,
		}
		
		input := RedemptionInput{
			Customer:  customer,
			RewardID:  reward.ID,
			Quantity:  1,
			Timestamp: time.Now(),
		}
		
		result, err := calc.RedeemPoints(input, reward)
		if err != nil {
			t.Fatalf("RedeemPoints failed: %v", err)
		}
		
		if result.IsSuccessful {
			t.Error("Redemption should not be successful")
		}
		
		if len(result.Errors) == 0 {
			t.Error("Should have errors for insufficient points")
		}
	})
	
	t.Run("SilverTierRedemptionBonus", func(t *testing.T) {
		customer := Customer{
			ID:            "customer3",
			Tier:          TierSilver,
			CurrentPoints: 500,
		}
		
		reward := Reward{
			ID:         "reward1",
			Name:       "$10 Discount",
			PointsCost: 200,
			Value:      10.0,
			IsActive:   true,
		}
		
		input := RedemptionInput{
			Customer:  customer,
			RewardID:  reward.ID,
			Quantity:  1,
			Timestamp: time.Now(),
		}
		
		result, err := calc.RedeemPoints(input, reward)
		if err != nil {
			t.Fatalf("RedeemPoints failed: %v", err)
		}
		
		// Silver tier has 10% redemption bonus
		expectedDiscount := reward.Value * 1.1
		if result.DiscountAmount != expectedDiscount {
			t.Errorf("Expected discount amount %.2f, got %.2f", expectedDiscount, result.DiscountAmount)
		}
	})
	
	t.Run("MultipleQuantity", func(t *testing.T) {
		customer := Customer{
			ID:            "customer4",
			Tier:          TierBronze,
			CurrentPoints: 1000,
		}
		
		reward := Reward{
			ID:         "reward1",
			Name:       "$5 Discount",
			PointsCost: 100,
			Value:      5.0,
			IsActive:   true,
		}
		
		input := RedemptionInput{
			Customer:  customer,
			RewardID:  reward.ID,
			Quantity:  3,
			Timestamp: time.Now(),
		}
		
		result, err := calc.RedeemPoints(input, reward)
		if err != nil {
			t.Fatalf("RedeemPoints failed: %v", err)
		}
		
		expectedPointsRedeemed := reward.PointsCost * 3
		if result.PointsRedeemed != expectedPointsRedeemed {
			t.Errorf("Expected %d points redeemed, got %d", expectedPointsRedeemed, result.PointsRedeemed)
		}
		
		expectedDiscount := reward.Value * 3
		if result.DiscountAmount != expectedDiscount {
			t.Errorf("Expected discount amount %.2f, got %.2f", expectedDiscount, result.DiscountAmount)
		}
	})
}

func TestCalculateReferralReward(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	t.Run("ValidReferral", func(t *testing.T) {
		referrer := Customer{
			ID:            "referrer1",
			Email:         "referrer@test.com",
			CurrentPoints: 100,
		}
		
		referee := Customer{
			ID:    "referee1",
			Email: "referee@test.com",
		}
		
		program := ReferralProgram{
			ID:             "ref1",
			Name:           "Friend Referral",
			ReferrerReward: 50,
			RefereeReward:  25,
			MinOrderAmount: 50.0,
			IsActive:       true,
		}
		
		result, err := calc.CalculateReferralReward(referrer, referee, program, 100.0)
		if err != nil {
			t.Fatalf("CalculateReferralReward failed: %v", err)
		}
		
		if !result.IsValid {
			t.Error("Result should be valid")
		}
		
		if result.BonusPoints != program.ReferrerReward {
			t.Errorf("Expected %d bonus points, got %d", program.ReferrerReward, result.BonusPoints)
		}
		
		if result.TotalPoints != program.ReferrerReward {
			t.Errorf("Expected %d total points, got %d", program.ReferrerReward, result.TotalPoints)
		}
		
		expectedBalance := referrer.CurrentPoints + program.ReferrerReward
		if result.NewBalance != expectedBalance {
			t.Errorf("Expected balance %d, got %d", expectedBalance, result.NewBalance)
		}
		
		if len(result.Transactions) == 0 {
			t.Error("Should have transactions")
		}
	})
	
	t.Run("InactiveProgram", func(t *testing.T) {
		referrer := Customer{ID: "referrer1"}
		referee := Customer{ID: "referee1"}
		
		program := ReferralProgram{
			IsActive: false,
		}
		
		_, err := calc.CalculateReferralReward(referrer, referee, program, 100.0)
		if err == nil {
			t.Error("Expected error for inactive program")
		}
	})
	
	t.Run("BelowMinimumOrder", func(t *testing.T) {
		referrer := Customer{ID: "referrer1"}
		referee := Customer{ID: "referee1"}
		
		program := ReferralProgram{
			MinOrderAmount: 100.0,
			IsActive:       true,
		}
		
		_, err := calc.CalculateReferralReward(referrer, referee, program, 50.0)
		if err == nil {
			t.Error("Expected error for order below minimum")
		}
	})
}

func TestCalculateReviewReward(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	t.Run("BasicReview", func(t *testing.T) {
		customer := Customer{
			ID:            "customer1",
			CurrentPoints: 100,
		}
		
		reward := ReviewReward{
			ID:            "review1",
			Name:          "Product Review",
			BasePoints:    10,
			PhotoBonus:    5,
			VideoBonus:    10,
			VerifiedBonus: 5,
			MinRating:     3,
			MinCharacters: 50,
			IsActive:      true,
		}
		
		result, err := calc.CalculateReviewReward(customer, reward, false, false, false, 4, 100)
		if err != nil {
			t.Fatalf("CalculateReviewReward failed: %v", err)
		}
		
		if !result.IsValid {
			t.Error("Result should be valid")
		}
		
		if result.BasePoints != reward.BasePoints {
			t.Errorf("Expected %d base points, got %d", reward.BasePoints, result.BasePoints)
		}
		
		if result.BonusPoints != 0 {
			t.Errorf("Expected 0 bonus points, got %d", result.BonusPoints)
		}
		
		if result.TotalPoints != reward.BasePoints {
			t.Errorf("Expected %d total points, got %d", reward.BasePoints, result.TotalPoints)
		}
	})
	
	t.Run("ReviewWithAllBonuses", func(t *testing.T) {
		customer := Customer{
			ID:            "customer2",
			CurrentPoints: 100,
		}
		
		reward := ReviewReward{
			ID:            "review1",
			BasePoints:    10,
			PhotoBonus:    5,
			VideoBonus:    10,
			VerifiedBonus: 5,
			MinRating:     3,
			MinCharacters: 50,
			IsActive:      true,
		}
		
		result, err := calc.CalculateReviewReward(customer, reward, true, true, true, 5, 100)
		if err != nil {
			t.Fatalf("CalculateReviewReward failed: %v", err)
		}
		
		expectedBonusPoints := reward.PhotoBonus + reward.VideoBonus + reward.VerifiedBonus
		if result.BonusPoints != expectedBonusPoints {
			t.Errorf("Expected %d bonus points, got %d", expectedBonusPoints, result.BonusPoints)
		}
		
		expectedTotalPoints := reward.BasePoints + expectedBonusPoints
		if result.TotalPoints != expectedTotalPoints {
			t.Errorf("Expected %d total points, got %d", expectedTotalPoints, result.TotalPoints)
		}
	})
	
	t.Run("RatingTooLow", func(t *testing.T) {
		customer := Customer{ID: "customer3"}
		
		reward := ReviewReward{
			MinRating: 4,
			IsActive:  true,
		}
		
		_, err := calc.CalculateReviewReward(customer, reward, false, false, false, 3, 100)
		if err == nil {
			t.Error("Expected error for rating too low")
		}
	})
	
	t.Run("ReviewTooShort", func(t *testing.T) {
		customer := Customer{ID: "customer4"}
		
		reward := ReviewReward{
			MinRating:     3,
			MinCharacters: 100,
			IsActive:      true,
		}
		
		_, err := calc.CalculateReviewReward(customer, reward, false, false, false, 4, 50)
		if err == nil {
			t.Error("Expected error for review too short")
		}
	})
	
	t.Run("InactiveReward", func(t *testing.T) {
		customer := Customer{ID: "customer5"}
		
		reward := ReviewReward{
			IsActive: false,
		}
		
		_, err := calc.CalculateReviewReward(customer, reward, false, false, false, 4, 100)
		if err == nil {
			t.Error("Expected error for inactive reward")
		}
	})
}

func TestGetAvailableRewards(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	customer := Customer{
		ID:            "customer1",
		Tier:          TierSilver,
		CurrentPoints: 500,
	}
	
	rewards := []Reward{
		{
			ID:         "reward1",
			Name:       "$5 Discount",
			PointsCost: 100,
			IsActive:   true,
			ValidFrom:  time.Now().Add(-time.Hour),
			ValidUntil: time.Now().Add(time.Hour),
		},
		{
			ID:           "reward2",
			Name:         "Gold Exclusive",
			PointsCost:   200,
			RequiredTier: TierGold, // Customer is Silver
			IsActive:     true,
		},
		{
			ID:         "reward3",
			Name:       "Expensive Reward",
			PointsCost: 1000, // Customer doesn't have enough points
			IsActive:   true,
		},
		{
			ID:       "reward4",
			Name:     "Inactive Reward",
			IsActive: false,
		},
		{
			ID:         "reward5",
			Name:       "$10 Discount",
			PointsCost: 200,
			IsActive:   true,
			ValidFrom:  time.Now().Add(-time.Hour),
			ValidUntil: time.Now().Add(time.Hour),
		},
	}
	
	available := calc.GetAvailableRewards(customer, rewards)
	
	// Should only include reward1 and reward5
	if len(available) != 2 {
		t.Errorf("Expected 2 available rewards, got %d", len(available))
	}
	
	// Should be sorted by points cost (reward1 first)
	if available[0].ID != "reward1" {
		t.Errorf("Expected reward1 first, got %s", available[0].ID)
	}
	
	if available[1].ID != "reward5" {
		t.Errorf("Expected reward5 second, got %s", available[1].ID)
	}
}

func TestHelperFunctions(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	t.Run("CalculateBasePoints", func(t *testing.T) {
		input := PointsCalculationInput{
			OrderAmount: 123.45,
		}
		
		basePoints := calc.calculateBasePoints(input)
		expected := 123 // Floor of 123.45 * 1.0
		
		if basePoints != expected {
			t.Errorf("Expected %d base points, got %d", expected, basePoints)
		}
	})
	
	t.Run("GetTierBenefit", func(t *testing.T) {
		benefit := calc.getTierBenefit(TierGold)
		
		if benefit.Tier != TierGold {
			t.Errorf("Expected Gold tier benefit, got %s", benefit.Tier)
		}
		
		if benefit.PointsMultiplier != 1.5 {
			t.Errorf("Expected 1.5 multiplier, got %f", benefit.PointsMultiplier)
		}
		
		// Test unknown tier
		unknownBenefit := calc.getTierBenefit("unknown")
		if unknownBenefit.PointsMultiplier != 1.0 {
			t.Errorf("Expected default 1.0 multiplier for unknown tier, got %f", unknownBenefit.PointsMultiplier)
		}
	})
	
	t.Run("GetNextTier", func(t *testing.T) {
		if calc.getNextTier(TierBronze) != TierSilver {
			t.Error("Bronze should upgrade to Silver")
		}
		
		if calc.getNextTier(TierSilver) != TierGold {
			t.Error("Silver should upgrade to Gold")
		}
		
		if calc.getNextTier(TierGold) != TierPlatinum {
			t.Error("Gold should upgrade to Platinum")
		}
		
		if calc.getNextTier(TierPlatinum) != TierPlatinum {
			t.Error("Platinum should stay Platinum")
		}
	})
	
	t.Run("GetTierThreshold", func(t *testing.T) {
		if calc.getTierThreshold(TierSilver) != 1000 {
			t.Error("Silver threshold should be 1000")
		}
		
		if calc.getTierThreshold(TierGold) != 5000 {
			t.Error("Gold threshold should be 5000")
		}
		
		if calc.getTierThreshold(TierPlatinum) != 15000 {
			t.Error("Platinum threshold should be 15000")
		}
	})
	
	t.Run("CalculateExpiryDate", func(t *testing.T) {
		expiry := calc.calculateExpiryDate(TierSilver)
		
		// Silver tier has 18 months expiry
		expected := time.Now().AddDate(0, 18, 0)
		
		// Allow some tolerance for test execution time
		if expiry.Before(expected.Add(-time.Minute)) || expiry.After(expected.Add(time.Minute)) {
			t.Errorf("Expiry date not within expected range")
		}
	})
	
	t.Run("GenerateTransactionID", func(t *testing.T) {
		id1 := calc.generateTransactionID()
		id2 := calc.generateTransactionID()
		
		if id1 == id2 {
			t.Error("Transaction IDs should be unique")
		}
		
		if id1 == "" || id2 == "" {
			t.Error("Transaction IDs should not be empty")
		}
	})
	
	t.Run("GenerateRedemptionCode", func(t *testing.T) {
		code1 := calc.generateRedemptionCode()
		code2 := calc.generateRedemptionCode()
		
		if code1 == code2 {
			t.Error("Redemption codes should be unique")
		}
		
		if code1 == "" || code2 == "" {
			t.Error("Redemption codes should not be empty")
		}
	})
}

func TestValidation(t *testing.T) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	t.Run("ValidateInput", func(t *testing.T) {
		// Valid input
		validInput := PointsCalculationInput{
			Customer:    Customer{ID: "customer1"},
			OrderAmount: 100.0,
			Timestamp:   time.Now(),
		}
		
		if err := calc.validateInput(validInput); err != nil {
			t.Errorf("Valid input should not return error: %v", err)
		}
		
		// Invalid customer ID
		invalidInput := PointsCalculationInput{
			Customer:    Customer{ID: ""},
			OrderAmount: 100.0,
			Timestamp:   time.Now(),
		}
		
		if err := calc.validateInput(invalidInput); err == nil {
			t.Error("Invalid customer ID should return error")
		}
		
		// Negative order amount
		negativeInput := PointsCalculationInput{
			Customer:    Customer{ID: "customer1"},
			OrderAmount: -50.0,
			Timestamp:   time.Now(),
		}
		
		if err := calc.validateInput(negativeInput); err == nil {
			t.Error("Negative order amount should return error")
		}
		
		// Zero timestamp
		zeroTimeInput := PointsCalculationInput{
			Customer:    Customer{ID: "customer1"},
			OrderAmount: 100.0,
			Timestamp:   time.Time{},
		}
		
		if err := calc.validateInput(zeroTimeInput); err == nil {
			t.Error("Zero timestamp should return error")
		}
	})
	
	t.Run("ValidateRedemptionInput", func(t *testing.T) {
		reward := Reward{
			ID:         "reward1",
			PointsCost: 100,
			IsActive:   true,
			ValidFrom:  time.Now().Add(-time.Hour),
			ValidUntil: time.Now().Add(time.Hour),
		}
		
		// Valid input
		validInput := RedemptionInput{
			Customer: Customer{ID: "customer1", CurrentPoints: 500},
			RewardID: "reward1",
			Quantity: 1,
		}
		
		if err := calc.validateRedemptionInput(validInput, reward); err != nil {
			t.Errorf("Valid redemption input should not return error: %v", err)
		}
		
		// Invalid customer ID
		invalidInput := RedemptionInput{
			Customer: Customer{ID: ""},
			RewardID: "reward1",
			Quantity: 1,
		}
		
		if err := calc.validateRedemptionInput(invalidInput, reward); err == nil {
			t.Error("Invalid customer ID should return error")
		}
		
		// Invalid reward ID
		emptyRewardInput := RedemptionInput{
			Customer: Customer{ID: "customer1"},
			RewardID: "",
			Quantity: 1,
		}
		
		if err := calc.validateRedemptionInput(emptyRewardInput, reward); err == nil {
			t.Error("Empty reward ID should return error")
		}
		
		// Invalid quantity
		zeroQuantityInput := RedemptionInput{
			Customer: Customer{ID: "customer1"},
			RewardID: "reward1",
			Quantity: 0,
		}
		
		if err := calc.validateRedemptionInput(zeroQuantityInput, reward); err == nil {
			t.Error("Zero quantity should return error")
		}
	})
}

// Helper function to create test configuration
func getTestConfig() *LoyaltyConfiguration {
	return &LoyaltyConfiguration{
		ProgramName:    "Test Loyalty Program",
		BaseCurrency:   "USD",
		BasePointsRate: 1.0,
		RedemptionRate: 0.01,
		PointsExpiry:   12,
		MinRedemption:  100,
		MaxRedemptionPercent: 50.0,
		TierThresholds: map[LoyaltyTier]float64{
			TierBronze:   0,
			TierSilver:   1000,
			TierGold:     5000,
			TierPlatinum: 15000,
		},
		TierBenefits: map[LoyaltyTier]TierBenefit{
			TierBronze: {
				Tier:             TierBronze,
				PointsMultiplier: 1.0,
				BonusPointsPercent: 0,
				RedemptionBonus:  0,
				MaxPointsExpiry:  12,
			},
			TierSilver: {
				Tier:             TierSilver,
				PointsMultiplier: 1.2,
				BonusPointsPercent: 5,
				RedemptionBonus:  0.1,
				MaxPointsExpiry:  18,
			},
			TierGold: {
				Tier:             TierGold,
				PointsMultiplier: 1.5,
				BonusPointsPercent: 10,
				RedemptionBonus:  0.2,
				MaxPointsExpiry:  24,
			},
			TierPlatinum: {
				Tier:             TierPlatinum,
				PointsMultiplier: 2.0,
				BonusPointsPercent: 15,
				RedemptionBonus:  0.3,
				MaxPointsExpiry:  36,
			},
		},
		DefaultRules: []LoyaltyRule{},
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func BenchmarkCalculate(b *testing.B) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	customer := Customer{
		ID:            "customer1",
		Tier:          TierSilver,
		CurrentPoints: 500,
		AnnualSpend:   2000,
	}
	
	input := PointsCalculationInput{
		Customer:    customer,
		OrderAmount: 100.0,
		Timestamp:   time.Now(),
		OrderID:     "order1",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.Calculate(input)
	}
}

func BenchmarkRedeemPoints(b *testing.B) {
	config := getTestConfig()
	calc := NewCalculator(config)
	
	customer := Customer{
		ID:            "customer1",
		Tier:          TierSilver,
		CurrentPoints: 1000,
	}
	
	reward := Reward{
		ID:         "reward1",
		Name:       "$5 Discount",
		PointsCost: 100,
		Value:      5.0,
		IsActive:   true,
		ValidFrom:  time.Now().Add(-time.Hour),
		ValidUntil: time.Now().Add(time.Hour),
	}
	
	input := RedemptionInput{
		Customer:  customer,
		RewardID:  reward.ID,
		Quantity:  1,
		Timestamp: time.Now(),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		calc.RedeemPoints(input, reward)
	}
}
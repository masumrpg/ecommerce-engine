package coupon

import (
	"testing"
	"time"
)

// TestValidateCouponRules tests coupon validation with rules
func TestValidateCouponRules(t *testing.T) {
	tests := []struct {
		name           string
		coupon         Coupon
		rules          []ValidationRule
		input          CalculationInput
		userEligibility UserEligibility
		expectError    bool
	}{
		{
			name: "valid user-based rule - first purchase",
			coupon: Coupon{
				Code:       "FIRST10",
				Type:       CouponTypePercentage,
				Value:      10.0,
				ValidFrom:  time.Now().Add(-24 * time.Hour),
				ValidUntil: time.Now().Add(24 * time.Hour),
				MinOrder:   50.0,
				MaxUsage:   100,
				IsActive:   true,
			},
			rules: []ValidationRule{
				{
					Type:         "user_based",
					Condition:    "first_purchase",
					Value:        true,
					ErrorMessage: "This coupon is only for first-time customers",
				},
			},
			input: CalculationInput{
				OrderAmount: 100.0,
				UserID:      "user123",
			},
			userEligibility: UserEligibility{
				IsFirstPurchase: true,
				LoyaltyTier:     "bronze",
				MemberSince:     time.Now().Add(-30 * 24 * time.Hour),
			},
			expectError: false,
		},
		{
			name: "invalid user-based rule - not first purchase",
			coupon: Coupon{
				Code:       "FIRST10",
				Type:       CouponTypePercentage,
				Value:      10.0,
				ValidFrom:  time.Now().Add(-24 * time.Hour),
				ValidUntil: time.Now().Add(24 * time.Hour),
				MinOrder:   50.0,
				MaxUsage:   100,
				IsActive:   true,
			},
			rules: []ValidationRule{
				{
					Type:         "user_based",
					Condition:    "first_purchase",
					Value:        true,
					ErrorMessage: "This coupon is only for first-time customers",
				},
			},
			input: CalculationInput{
				OrderAmount: 100.0,
				UserID:      "user123",
			},
			userEligibility: UserEligibility{
				IsFirstPurchase: false,
				LoyaltyTier:     "bronze",
				MemberSince:     time.Now().Add(-30 * 24 * time.Hour),
			},
			expectError: true,
		},
		{
			name: "valid order-based rule - minimum amount",
			coupon: Coupon{
				Code:       "MIN100",
				Type:       CouponTypeFixedAmount,
				Value:      20.0,
				ValidFrom:  time.Now().Add(-24 * time.Hour),
				ValidUntil: time.Now().Add(24 * time.Hour),
				MinOrder:   50.0,
				MaxUsage:   100,
				IsActive:   true,
			},
			rules: []ValidationRule{
				{
					Type:         "order_based",
					Condition:    "minimum_amount",
					Value:        100.0,
					ErrorMessage: "Minimum order amount is $100",
				},
			},
			input: CalculationInput{
				OrderAmount: 150.0,
				UserID:      "user123",
			},
			userEligibility: UserEligibility{
				IsFirstPurchase: false,
				LoyaltyTier:     "silver",
				MemberSince:     time.Now().Add(-60 * 24 * time.Hour),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCouponRules(tt.coupon, tt.rules, tt.input, tt.userEligibility)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCouponRules() error = %v, expectError = %v", err, tt.expectError)
			}
		})
	}
}

// TestValidateCouponStacking tests coupon stacking validation
func TestValidateCouponStacking(t *testing.T) {
	tests := []struct {
		name         string
		coupons      []Coupon
		stackingRules map[string]interface{}
		expectError  bool
	}{
		{
			name: "valid stacking - under limit",
			coupons: []Coupon{
				{
					Code: "SAVE10",
					Type: CouponTypePercentage,
				},
				{
					Code: "FREESHIP",
					Type: CouponTypeFreeShipping,
				},
			},
			stackingRules: map[string]interface{}{
				"max_stackable": float64(3),
				"allow_same_type": true,
			},
			expectError: false,
		},
		{
			name: "invalid stacking - over limit",
			coupons: []Coupon{
				{Code: "SAVE10", Type: CouponTypePercentage},
				{Code: "SAVE20", Type: CouponTypePercentage},
				{Code: "FREESHIP", Type: CouponTypeFreeShipping},
			},
			stackingRules: map[string]interface{}{
				"max_stackable": float64(2),
			},
			expectError: true,
		},
		{
			name: "invalid stacking - same type not allowed",
			coupons: []Coupon{
				{Code: "SAVE10", Type: CouponTypePercentage},
				{Code: "SAVE20", Type: CouponTypePercentage},
			},
			stackingRules: map[string]interface{}{
				"allow_same_type": false,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCouponStacking(tt.coupons, tt.stackingRules)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCouponStacking() error = %v, expectError = %v", err, tt.expectError)
			}
		})
	}
}

// TestValidateBusinessRules tests business rules validation
func TestValidateBusinessRules(t *testing.T) {
	tests := []struct {
		name          string
		coupon        Coupon
		input         CalculationInput
		businessRules map[string]interface{}
		expectError   bool
	}{
		{
			name: "valid business rules",
			coupon: Coupon{
				Code:  "SAVE10",
				Type:  CouponTypePercentage,
				Value: 10.0,
			},
			input: CalculationInput{
				Coupon: Coupon{
					Value: 10.0,
				},
				OrderAmount: 100.0,
				UserID:      "user123",
			},
			businessRules: map[string]interface{}{
				"minimum_margin_percent": float64(20),
				"blacklisted_users":       []string{"baduser1", "baduser2"},
			},
			expectError: false,
		},
		{
			name: "invalid - blacklisted user",
			coupon: Coupon{
				Code:  "SAVE10",
				Type:  CouponTypePercentage,
				Value: 10.0,
			},
			input: CalculationInput{
				OrderAmount: 100.0,
				UserID:      "baduser1",
			},
			businessRules: map[string]interface{}{
				"blacklisted_users": []string{"baduser1", "baduser2"},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBusinessRules(tt.coupon, tt.input, tt.businessRules)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateBusinessRules() error = %v, expectError = %v", err, tt.expectError)
			}
		})
	}
}

// Benchmark tests
func BenchmarkValidateCouponRules(b *testing.B) {
	coupon := Coupon{
		Code:       "SAVE10",
		Type:       CouponTypePercentage,
		Value:      10.0,
		IsActive:   true,
		ValidFrom:  time.Now().Add(-24 * time.Hour),
		ValidUntil: time.Now().Add(24 * time.Hour),
		MinOrder:   0,
		MaxUsage:   100,
	}

	rules := []ValidationRule{
		{
			Type:         "user_based",
			Condition:    "first_purchase",
			Value:        true,
			ErrorMessage: "First purchase only",
		},
	}

	input := CalculationInput{
		OrderAmount: 100.0,
		UserID:      "user123",
	}

	userEligibility := UserEligibility{
		IsFirstPurchase: true,
		LoyaltyTier:     "bronze",
		MemberSince:     time.Now().Add(-30 * 24 * time.Hour),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCouponRules(coupon, rules, input, userEligibility)
	}
}

func BenchmarkValidateCouponStacking(b *testing.B) {
	coupons := []Coupon{
		{Code: "SAVE10", Type: CouponTypePercentage},
		{Code: "FREESHIP", Type: CouponTypeFreeShipping},
	}

	stackingRules := map[string]interface{}{
		"max_stackable": float64(3),
		"allow_same_type": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCouponStacking(coupons, stackingRules)
	}
}
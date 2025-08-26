package coupon

import (
	"testing"
	"time"
)

func TestCalculate(t *testing.T) {
	t.Run("PercentageDiscount", func(t *testing.T) {
		coupon := Coupon{
			Code:       "SAVE10",
			Type:       CouponTypePercentage,
			Value:      10.0,
			MinOrder:   50.0,
			ValidFrom:  time.Now().Add(-24 * time.Hour),
			ValidUntil: time.Now().Add(24 * time.Hour),
			IsActive:   true,
			MaxUsage:   100,
		}
		
		items := []Item{
			{ID: "item1", Price: 50.0, Quantity: 1, Category: "electronics"},
			{ID: "item2", Price: 50.0, Quantity: 1, Category: "electronics"},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 100.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Error("Expected coupon to be valid")
		}
		
		if result.DiscountAmount != 10.0 {
			t.Errorf("Expected discount amount 10.0, got %f", result.DiscountAmount)
		}
		
		if result.ErrorMessage != "" {
			t.Errorf("Expected no error message, got: %s", result.ErrorMessage)
		}
	})
	
	t.Run("FixedAmountDiscount", func(t *testing.T) {
		coupon := Coupon{
			Code:       "SAVE5",
			Type:       CouponTypeFixedAmount,
			Value:      5.0,
			MinOrder:   10.0,
			ValidFrom:  time.Now().Add(-24 * time.Hour),
			ValidUntil: time.Now().Add(24 * time.Hour),
			IsActive:   true,
			MaxUsage:   100,
		}
		
		items := []Item{
			{ID: "item1", Price: 20.0, Quantity: 1, Category: "books"},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 20.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Error("Expected coupon to be valid")
		}
		
		if result.DiscountAmount != 5.0 {
			t.Errorf("Expected discount amount 5.0, got %f", result.DiscountAmount)
		}
	})
	
	t.Run("FreeShippingDiscount", func(t *testing.T) {
		coupon := Coupon{
			Code:       "FREESHIP",
			Type:       CouponTypeFreeShipping,
			Value:      0.0,
			MinOrder:   25.0,
			ValidFrom:  time.Now().Add(-24 * time.Hour),
			ValidUntil: time.Now().Add(24 * time.Hour),
			IsActive:   true,
			MaxUsage:   100,
		}
		
		items := []Item{
			{ID: "item1", Price: 30.0, Quantity: 1, Category: "clothing"},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 30.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Error("Expected coupon to be valid")
		}
		
		// Free shipping typically has 0 discount amount
		if result.DiscountAmount != 0.0 {
			t.Errorf("Expected discount amount 0.0, got %f", result.DiscountAmount)
		}
	})
	
	t.Run("BuyXGetYDiscount", func(t *testing.T) {
		coupon := Coupon{
			Code:       "BUY2GET1",
			Type:       CouponTypeBuyXGetY,
			Value:      0.0,
			BuyX:       2,
			GetY:       1,
			ValidFrom:  time.Now().Add(-24 * time.Hour),
			ValidUntil: time.Now().Add(24 * time.Hour),
			IsActive:   true,
			MaxUsage:   100,
		}
		
		items := []Item{
			{ID: "item1", Price: 10.0, Quantity: 3, Category: "toys"},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 30.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if !result.IsValid {
			t.Error("Expected coupon to be valid")
		}
		
		// Should get 1 free item (buy 2 get 1)
		if result.DiscountAmount != 10.0 {
			t.Errorf("Expected discount amount 10.0, got %f", result.DiscountAmount)
		}
	})
	
	t.Run("InvalidCoupon - Inactive", func(t *testing.T) {
		coupon := Coupon{
			Code:       "INACTIVE",
			Type:       CouponTypePercentage,
			Value:      10.0,
			IsActive:   false,
		}
		
		items := []Item{
			{ID: "item1", Price: 100.0, Quantity: 1},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 100.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if result.IsValid {
			t.Error("Expected coupon to be invalid")
		}
		
		if result.ErrorMessage == "" {
			t.Error("Expected error message to be set")
		}
	})
	
	t.Run("InvalidCoupon - Expired", func(t *testing.T) {
		coupon := Coupon{
			Code:       "EXPIRED",
			Type:       CouponTypePercentage,
			Value:      10.0,
			ValidFrom:  time.Now().Add(-48 * time.Hour),
			ValidUntil: time.Now().Add(-24 * time.Hour),
			IsActive:   true,
		}
		
		items := []Item{
			{ID: "item1", Price: 100.0, Quantity: 1},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 100.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if result.IsValid {
			t.Error("Expected coupon to be invalid")
		}
		
		if result.ErrorMessage == "" {
			t.Error("Expected error message to be set")
		}
	})
	
	t.Run("InvalidCoupon - BelowMinOrder", func(t *testing.T) {
		coupon := Coupon{
			Code:       "MINORDER",
			Type:       CouponTypePercentage,
			Value:      10.0,
			MinOrder:   100.0,
			ValidFrom:  time.Now().Add(-24 * time.Hour),
			ValidUntil: time.Now().Add(24 * time.Hour),
			IsActive:   true,
		}
		
		items := []Item{
			{ID: "item1", Price: 50.0, Quantity: 1},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 50.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if result.IsValid {
			t.Error("Expected coupon to be invalid")
		}
		
		if result.ErrorMessage == "" {
			t.Error("Expected error message to be set")
		}
	})
	
	t.Run("InvalidCoupon - UsageLimitExceeded", func(t *testing.T) {
		coupon := Coupon{
			Code:       "LIMITEXCEEDED",
			Type:       CouponTypePercentage,
			Value:      10.0,
			ValidFrom:  time.Now().Add(-24 * time.Hour),
			ValidUntil: time.Now().Add(24 * time.Hour),
			IsActive:   true,
			MaxUsage:   10,
		}
		
		items := []Item{
			{ID: "item1", Price: 100.0, Quantity: 1},
		}
		
		input := CalculationInput{
			Coupon:      coupon,
			OrderAmount: 100.0,
			UserID:      "user123",
			Items:       items,
			Usage:       CouponUsage{TotalUsage: 10, UsageCount: 0},
		}
		
		result := Calculate(input)
		
		if result.IsValid {
			t.Error("Expected coupon to be invalid")
		}
		
		if result.ErrorMessage == "" {
			t.Error("Expected error message to be set")
		}
	})
}

func BenchmarkCalculate(b *testing.B) {
	coupon := Coupon{
		Code:       "BENCH",
		Type:       CouponTypePercentage,
		Value:      10.0,
		MinOrder:   50.0,
		ValidFrom:  time.Now().Add(-24 * time.Hour),
		ValidUntil: time.Now().Add(24 * time.Hour),
		IsActive:   true,
		MaxUsage:   1000000,
	}
	
	items := []Item{
		{ID: "item1", Price: 50.0, Quantity: 1, Category: "electronics"},
		{ID: "item2", Price: 50.0, Quantity: 1, Category: "electronics"},
	}
	
	input := CalculationInput{
		Coupon:      coupon,
		OrderAmount: 100.0,
		UserID:      "user123",
		Items:       items,
		Usage:       CouponUsage{TotalUsage: 0, UsageCount: 0},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Calculate(input)
	}
}
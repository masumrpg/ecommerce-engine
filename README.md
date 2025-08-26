# E-commerce Engine

A comprehensive, modular e-commerce calculation engine built in Go, designed with business logic purity and pure calculation/rule engine focus.

## Overview

This e-commerce engine provides a collection of independent, reusable packages for common e-commerce calculations including pricing, discounts, coupons, shipping, taxes, loyalty programs, and utility functions. Each package is designed to be stateless, pure, and focused on business logic without external dependencies.

## Features

- **🏷️ Pricing Engine**: Dynamic pricing, tier-based pricing, bundle pricing
- **💰 Discount System**: Rule-based discounts, promotional campaigns
- **🎫 Coupon Management**: Generation, validation, and redemption
- **📦 Shipping Calculator**: Multi-carrier support, zone-based pricing
- **💸 Tax Calculator**: Multi-jurisdiction tax calculations
- **⭐ Loyalty Program**: Points calculation, tier management, rewards
- **🔧 Utilities**: Math functions, ID generation, statistical analysis

## Architecture

### Design Principles

- **Business Logic Purity**: Pure calculation engines without side effects
- **Stateless Operations**: No persistent state, all data passed as parameters
- **Modular Design**: Independent packages that can be used separately
- **Type Safety**: Strong typing with comprehensive data structures
- **Performance**: Optimized for high-throughput calculations

### Package Structure

```
ecommerce-engine/
├── pkg/
│   ├── coupon/          # Coupon generation and validation
│   ├── discount/        # Discount calculation engine
│   ├── loyalty/         # Loyalty points and rewards
│   ├── pricing/         # Dynamic pricing engine
│   ├── shipping/        # Shipping cost calculations
│   ├── tax/            # Tax calculation engine
│   └── utils/          # Utility functions and helpers
├── examples/           # Usage examples
└── README.md
```

## Installation

```bash
go get github.com/masumrpg/ecommerce-engine
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/masumrpg/ecommerce-engine/pkg/utils"
)

func main() {
    // Generate a unique order ID
    orderID := utils.GenerateUUID()
    fmt.Printf("Order ID: %s\n", orderID)
    
    // Calculate percentage discount
    discount := utils.PercentageOf(15.0, 100.0) // 15% of $100
    fmt.Printf("Discount: $%.2f\n", discount)
    
    // Round to currency
    total := utils.RoundToCurrency(99.999)
    fmt.Printf("Total: $%.2f\n", total)
}
```

## Package Documentation

### Coupon Package

Manages coupon generation, validation, and redemption.

#### Key Features
- Multiple coupon types (percentage, fixed amount, free shipping)
- Usage limits and expiration dates
- Customer-specific restrictions
- Batch generation capabilities

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/coupon"

// Generate a coupon
generator := coupon.NewGenerator()
input := coupon.GenerationInput{
    Type:           coupon.CouponTypePercentage,
    Value:          15.0,
    MinOrderValue:  50.0,
    ExpirationDate: time.Now().AddDate(0, 1, 0),
}

result, err := generator.Generate(input)
if err == nil {
    fmt.Printf("Generated coupon: %s\n", result.Code)
}
```

### Discount Package

Rule-based discount calculation engine.

#### Key Features
- Percentage and fixed amount discounts
- Category-specific discounts
- Quantity-based discounts
- Time-based promotional rules

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/discount"

calc := discount.NewCalculator()
input := discount.CalculationInput{
    Items: []discount.DiscountableItem{
        {
            ID:       "item1",
            Price:    100.0,
            Quantity: 2,
            Category: "Electronics",
        },
    },
    CustomerID: "customer123",
    OrderValue: 200.0,
}

result, err := calc.Calculate(input)
```

### Loyalty Package

Comprehensive loyalty program management.

#### Key Features
- Points calculation based on purchase amount
- Tier-based benefits and multipliers
- Reward redemption system
- Referral program support
- Review reward system

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/loyalty"

calc := loyalty.NewCalculator()
input := loyalty.PointsCalculationInput{
    CustomerID: "customer123",
    OrderValue: 299.99,
    Items: []loyalty.OrderItem{
        {
            ID:       "item1",
            Price:    299.99,
            Quantity: 1,
            Category: "Electronics",
        },
    },
}

result, err := calc.CalculatePoints(input)
if err == nil {
    fmt.Printf("Points earned: %d\n", result.TotalPointsEarned)
}
```

### Pricing Package

Dynamic pricing engine with advanced features.

#### Key Features
- Base pricing calculations
- Tier-based pricing for different customer segments
- Bundle pricing and package deals
- Dynamic pricing based on demand/inventory
- Volume discounts

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/pricing"

calc := pricing.NewCalculator()
input := pricing.PricingInput{
    Items: []pricing.PricingItem{
        {
            ID:       "laptop",
            BasePrice: 999.99,
            Quantity: 1,
            Category: "Electronics",
        },
    },
    Customer: pricing.Customer{
        ID:   "customer123",
        Tier: "Premium",
    },
}

result, err := calc.Calculate(input)
```

### Shipping Package

Multi-carrier shipping cost calculator.

#### Key Features
- Weight and dimension-based calculations
- Zone-based pricing
- Multiple service types (standard, express, overnight)
- Insurance and tracking options
- Delivery time estimation

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/shipping"

calc := shipping.NewCalculator()
input := shipping.CalculationInput{
    Items: []shipping.ShippableItem{
        {
            Weight: 2.5,
            Dimensions: shipping.Dimensions{
                Length: 30.0,
                Width:  20.0,
                Height: 10.0,
            },
            Value: 299.99,
        },
    },
    Origin: shipping.Address{
        Country:    "US",
        State:      "CA",
        PostalCode: "90210",
    },
    Destination: shipping.Address{
        Country:    "US",
        State:      "NY",
        PostalCode: "10001",
    },
    ServiceType: shipping.ServiceStandard,
}

result, err := calc.Calculate(input)
```

### Tax Package

Multi-jurisdiction tax calculation engine.

#### Key Features
- Sales tax calculation
- VAT support
- Tax-exempt categories
- Business vs. consumer tax rates
- Multi-jurisdiction support

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/tax"

calc := tax.NewCalculator()
input := tax.TaxCalculationInput{
    Items: []tax.TaxableItem{
        {
            ID:       "item1",
            Price:    100.0,
            Quantity: 1,
            Category: "Electronics",
            Taxable:  true,
        },
    },
    Customer: tax.Customer{
        Address: tax.Address{
            Country:    "US",
            State:      "CA",
            PostalCode: "90210",
        },
    },
    Subtotal: 100.0,
}

result, err := calc.Calculate(input)
```

### Utils Package

Comprehensive utility functions for e-commerce operations.

#### Key Features

**Mathematical Functions:**
- Rounding with various modes
- Percentage calculations
- Statistical functions (average, median, standard deviation)
- Financial calculations (compound interest, present value)
- Distance calculations

**ID Generation:**
- UUID generation
- Sequential ID generation
- Custom ID patterns
- Timestamp-based IDs

**Code Generation:**
- Coupon code generation
- Token generation (API keys, OTP, session tokens)
- Reference number generation
- Barcode and SKU generation

**Text Processing:**
- URL-friendly slug generation
- Color code generation
- Random string generation

#### Example Usage

```go
import "github.com/masumrpg/ecommerce-engine/pkg/utils"

// Mathematical operations
total := utils.Round(123.456, 2)                    // 123.46
discount := utils.PercentageOf(15.0, 200.0)         // 30.00
avg := utils.Average([]float64{10, 20, 30})         // 20.00

// ID generation
uuid := utils.GenerateUUID()                        // "550e8400-e29b-41d4-a716-446655440000"
shortID := utils.GenerateShortID(8)                 // "A1B2C3D4"
timestampID := utils.GenerateTimestampID()          // "20240115123045"

// Code generation
couponGen := utils.NewCouponCodeGenerator(8)
couponCode := couponGen.GenerateCouponCode()        // "SAVE2024"

tokenGen := utils.NewTokenGenerator()
apiKey := tokenGen.GenerateAPIKey("ecom")           // "ecom_1234567890abcdef"
otp := tokenGen.GenerateOTP(6)                      // "123456"

// Reference generation
refGen := utils.NewReferenceGenerator("ORD", "", 6)
orderRef := refGen.GenerateOrderReference()         // "ORD123456"

// Barcode generation
barcodeGen := utils.NewBarcodeGenerator()
ean13 := barcodeGen.GenerateEAN13()                 // "1234567890123"
sku := barcodeGen.GenerateSKU("ELEC", "LAP")        // "ELEC-LAP-001"

// Slug generation
slugGen := utils.NewSlugGenerator()
slug := slugGen.GenerateSlug("Gaming Laptop Pro!")  // "gaming-laptop-pro"

// Color generation
colorGen := utils.NewColorGenerator()
hexColor := colorGen.GenerateHexColor()             // "#FF5733"
r, g, b := colorGen.GenerateRGBColor()              // 255, 87, 51
```

## Advanced Usage

### Statistical Analysis

The utils package provides comprehensive statistical functions for business analytics:

```go
// Sales data analysis
salesData := []float64{1250.50, 980.25, 1450.75, 1100.00, 1350.25}

totalSales := utils.Sum(salesData)                  // Total revenue
avgSales := utils.Average(salesData)                // Average daily sales
medianSales := utils.Median(salesData)              // Median sales
stdDev := utils.StandardDeviation(salesData)        // Sales volatility

// Moving averages for trend analysis
movingAvg := utils.MovingAverage(salesData, 3)      // 3-day moving average
emaData := utils.ExponentialMovingAverage(salesData, 0.3) // EMA with α=0.3

// Correlation analysis
adSpend := []float64{500, 300, 700, 450, 600}
correlation := utils.Correlation(salesData, adSpend) // Sales vs ad spend correlation

// Linear regression
slope, intercept := utils.LinearRegression(adSpend, salesData)
// Sales = slope * AdSpend + intercept
```

### Financial Calculations

```go
// Investment calculations
futureValue := utils.CompoundInterest(1000.0, 0.05, 3)  // $1000 at 5% for 3 years
presentValue := utils.PresentValue(1000.0, 0.05, 3)     // Present value of $1000

// Pricing models
sigmoidPrice := utils.Sigmoid(2.0)                      // Sigmoid pricing curve
decayPrice := utils.ExponentialDecay(100, 0.1, 5)       // Price decay over time
growthPrice := utils.ExponentialGrowth(100, 0.05, 12)   // Price growth projection
```

## Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Examples

See the `examples/` directory for comprehensive usage examples demonstrating all package features.

Run the examples:

```bash
go run examples/main.go
```

## Performance

The engine is optimized for high-throughput scenarios:

- **Stateless Design**: No memory overhead from persistent state
- **Pure Functions**: Optimizable by the Go compiler
- **Minimal Allocations**: Efficient memory usage
- **Concurrent Safe**: All functions are thread-safe

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
- Maintain business logic purity (no side effects)
- Add comprehensive tests for new features
- Update documentation for API changes
- Use meaningful variable and function names

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions, issues, or contributions, please:

- Open an issue on GitHub
- Submit a pull request
- Contact the maintainers

## Roadmap

- [ ] Add more shipping carriers
- [ ] Implement international tax calculations
- [ ] Add A/B testing framework for pricing
- [ ] Implement machine learning-based pricing recommendations
- [ ] Add GraphQL API layer
- [ ] Implement caching strategies
- [ ] Add monitoring and metrics

---

**Built with ❤️ for the e-commerce community**
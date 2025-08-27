# Go Documentation for ecommerce-engine

## Currency Package

package currency // import "github.com/masumrpg/ecommerce-engine/pkg/currency"

const DefaultPrecision = 2 ...
const RateFreshThreshold = 60 ...
const MaxAmount = 1e15 ...
var MajorCurrencies = []CurrencyCode{ ... } ...
var CommonPairs = []CurrencyPair{ ... }
var CurrencyDecimalPlaces = map[CurrencyCode]int{ ... }
var CurrencyNames = map[CurrencyCode]string{ ... }
var CurrencySymbols = map[CurrencyCode]string{ ... }
func GetCurrencyDecimalPlaces(code CurrencyCode) int
func GetCurrencyName(code CurrencyCode) string
func GetCurrencySymbol(code CurrencyCode) string
func IsAsianCurrency(code CurrencyCode) bool
func IsCommonPair(pair CurrencyPair) bool
func IsEuropeanCurrency(code CurrencyCode) bool
func IsMajorCurrency(code CurrencyCode) bool
func IsNegative(money Money) bool
func IsPositive(money Money) bool
func IsValidCurrencyCode(code CurrencyCode) bool
func IsZero(money Money) bool
func IsZeroDecimalCurrency(code CurrencyCode) bool
func Split(money Money, parts int) ([]Money, Money)
type ArithmeticInput struct{ ... }
type ArithmeticOperation string
    const OperationAdd ArithmeticOperation = "add" ...
type ArithmeticResult struct{ ... }
type BatchConverter struct{ ... }
    func NewBatchConverter(calculator *Calculator) *BatchConverter
type Calculator struct{ ... }
    func NewCalculator() *Calculator
type ComparisonResult struct{ ... }
type ConversionInput struct{ ... }
type ConversionResult struct{ ... }
type Currency struct{ ... }
type CurrencyCode string
    const USD CurrencyCode = "USD" ...
    func GetSupportedCurrencyCodes() []CurrencyCode
type CurrencyDetector struct{ ... }
    func NewCurrencyDetector(calculator *Calculator) *CurrencyDetector
type CurrencyError struct{ ... }
type CurrencyFormatter struct{ ... }
    func NewCurrencyFormatter(calculator *Calculator) *CurrencyFormatter
type CurrencyList struct{ ... }
type CurrencyPair struct{ ... }
type ExchangeRate struct{ ... }
type FormatOptions struct{ ... }
type LocaleInfo struct{ ... }
type Money struct{ ... }
    func Abs(money Money) Money
    func Allocate(money Money, ratios []float64) ([]Money, error)
    func Average(amounts []Money) (Money, error)
    func Max(a, b Money) (Money, error)
    func Min(a, b Money) (Money, error)
    func Negate(money Money) Money
    func NewMoney(amount float64, currency CurrencyCode) (*Money, error)
    func NewMoneyFromString(amountStr string, currency CurrencyCode) (*Money, error)
    func Percentage(money Money, percent float64) Money
    func Sum(amounts []Money) (Money, error)
type RateProvider string
    const ProviderManual RateProvider = "manual" ...
type RateSource struct{ ... }
type RoundingMode string
    const RoundingModeHalfUp RoundingMode = "half_up" ...
type ValidationError struct{ ... }
type Validator struct{ ... }
    func NewValidator(calculator *Calculator) *Validator

## Coupon Package

package coupon // import "github.com/masumrpg/ecommerce-engine/pkg/coupon"

func GenerateBulkCodes(configs []GeneratorConfig) (map[string][]string, error)
func GenerateCode(config GeneratorConfig) (string, error)
func GenerateCodes(config GeneratorConfig) ([]string, error)
func GenerateExpiryDate(duration time.Duration) time.Time
func GenerateFlashSaleCode(discountPercent int, config GeneratorConfig) (string, error)
func GenerateSeasonalCode(season string, year int, config GeneratorConfig) (string, error)
func ValidateBusinessRules(coupon Coupon, input CalculationInput, businessRules map[string]interface{}) error
func ValidateCodeFormat(code string, config GeneratorConfig) bool
func ValidateCouponRules(coupon Coupon, rules []ValidationRule, input CalculationInput, ...) error
func ValidateCouponStacking(coupons []Coupon, stackingRules map[string]interface{}) error
type CalculationInput struct{ ... }
type CalculationResult struct{ ... }
    func Calculate(input CalculationInput) CalculationResult
    func CalculateMultiple(coupons []Coupon, orderAmount float64, userID string, items []Item, ...) CalculationResult
type Coupon struct{ ... }
type CouponType string
    const CouponTypePercentage CouponType = "percentage" ...
type CouponUsage struct{ ... }
type GeneratorConfig struct{ ... }
type Item struct{ ... }
type UserEligibility struct{ ... }
type ValidationRule struct{ ... }

## Discount Package

package discount // import "github.com/masumrpg/ecommerce-engine/pkg/discount"

type BulkDiscountRule struct{ ... }
type BundleDiscountRule struct{ ... }
type BundleMatch struct{ ... }
type CategoryDiscountRule struct{ ... }
type CrossSellRule struct{ ... }
type Customer struct{ ... }
type DiscountApplication struct{ ... }
type DiscountCalculationInput struct{ ... }
type DiscountCalculationResult struct{ ... }
    func Calculate(input DiscountCalculationInput) DiscountCalculationResult
    func CalculateBestDiscount(inputs []DiscountCalculationInput) DiscountCalculationResult
type DiscountItem struct{ ... }
type DiscountType string
    const DiscountTypeBulk DiscountType = "bulk" ...
type DiscountValidator struct{ ... }
    func NewDiscountValidator() *DiscountValidator
type FrequencyDiscountRule struct{ ... }
type LoyaltyDiscountRule struct{ ... }
type MixAndMatchRule struct{ ... }
type ProgressiveDiscountRule struct{ ... }
type RuleEngine struct{ ... }
    func NewRuleEngine() *RuleEngine
type SeasonalDiscountRule struct{ ... }
type TierPricingRule struct{ ... }

## Shipping Package

package shipping // import "github.com/masumrpg/ecommerce-engine/pkg/shipping"

type Address struct{ ... }
type AppliedSurcharge struct{ ... }
type CarrierRule struct{ ... }
type DeliveryTimeRule struct{ ... }
type DimensionUnit string
    const DimensionUnitCM DimensionUnit = "cm" ...
type Dimensions struct{ ... }
type FreeShippingRule struct{ ... }
type Package struct{ ... }
type PackagingRule struct{ ... }
type PostalCodeRange struct{ ... }
type ShippingCalculationInput struct{ ... }
type ShippingCalculationResult struct{ ... }
    func Calculate(input ShippingCalculationInput) ShippingCalculationResult
type ShippingCalculator struct{ ... }
    func NewShippingCalculator() *ShippingCalculator
type ShippingItem struct{ ... }
type ShippingMethod string
    const ShippingMethodStandard ShippingMethod = "standard" ...
type ShippingOption struct{ ... }
    func CalculateBestOption(input ShippingCalculationInput, criteria string) (*ShippingOption, error)
type ShippingRestriction struct{ ... }
type ShippingRule struct{ ... }
type ShippingRuleEngine struct{ ... }
    func NewShippingRuleEngine() *ShippingRuleEngine
type ShippingZone string
    const ShippingZoneLocal ShippingZone = "local" ...
type Surcharge struct{ ... }
type Weight struct{ ... }
type WeightUnit string
    const WeightUnitKG WeightUnit = "kg" ...
type ZoneRule struct{ ... }

## Tax Package

package tax // import "github.com/masumrpg/ecommerce-engine/pkg/tax"

func GetTaxSummary(results []TaxCalculationResult) map[string]interface{}
type Address struct{ ... }
type AppliedTax struct{ ... }
type Customer struct{ ... }
type TaxAuditTrail struct{ ... }
type TaxBreakdown struct{ ... }
type TaxCalculationInput struct{ ... }
    func CalculateBestTaxStrategy(scenarios []TaxCalculationInput) (*TaxCalculationInput, error)
type TaxCalculationMethod string
    const TaxMethodPercentage TaxCalculationMethod = "percentage" ...
type TaxCalculationResult struct{ ... }
    func Calculate(input TaxCalculationInput) TaxCalculationResult
type TaxCalculator struct{ ... }
    func NewTaxCalculator(config TaxConfiguration) *TaxCalculator
type TaxCondition struct{ ... }
type TaxConfiguration struct{ ... }
type TaxExemption struct{ ... }
type TaxJurisdiction string
    const JurisdictionFederal TaxJurisdiction = "federal" ...
type TaxOverride struct{ ... }
type TaxReport struct{ ... }
type TaxReportDetail struct{ ... }
type TaxRule struct{ ... }
    func CreateDefaultRules() []TaxRule
type TaxRuleEngine struct{ ... }
    func NewTaxRuleEngine(config TaxConfiguration) *TaxRuleEngine
type TaxThreshold struct{ ... }
type TaxType string
    const TaxTypeSales TaxType = "sales" ...
type TaxValidationRule struct{ ... }
    func CreateDefaultValidationRules() []TaxValidationRule
type TaxableItem struct{ ... }

## Loyalty Package

package loyalty // import "github.com/masumrpg/ecommerce-engine/pkg/loyalty"

type AppliedLoyaltyRule struct{ ... }
type Calculator struct{ ... }
    func NewCalculator(config *LoyaltyConfiguration) *Calculator
type Customer struct{ ... }
type CustomerPreferences struct{ ... }
type LoyaltyAction struct{ ... }
type LoyaltyAnalytics struct{ ... }
type LoyaltyCondition struct{ ... }
type LoyaltyConfiguration struct{ ... }
    func CreateDefaultConfiguration() *LoyaltyConfiguration
type LoyaltyRecommendation struct{ ... }
type LoyaltyRule struct{ ... }
    func CreateDefaultRules() []LoyaltyRule
type LoyaltyTier string
    const TierBronze LoyaltyTier = "bronze" ...
type OrderItem struct{ ... }
type PointsBreakdown struct{ ... }
type PointsCalculationInput struct{ ... }
type PointsCalculationResult struct{ ... }
type PointsTransaction struct{ ... }
type PointsType string
    const PointsTypeBase PointsType = "base" ...
type RedemptionInput struct{ ... }
type RedemptionResult struct{ ... }
type ReferralProgram struct{ ... }
type ReviewReward struct{ ... }
type Reward struct{ ... }
    func CreateDefaultRewards() []Reward
type RewardType string
    const RewardTypeDiscount RewardType = "discount" ...
type RuleEngine struct{ ... }
    func NewRuleEngine(config *LoyaltyConfiguration) *RuleEngine
type TierBenefit struct{ ... }
type TierInfo struct{ ... }
type TransactionType string
    const TransactionTypeEarn TransactionType = "earn" ...

## Pricing Package

package pricing // import "github.com/masumrpg/ecommerce-engine/pkg/pricing"

type AppliedPricingRule struct{ ... }
type Bundle struct{ ... }
type BundleAction struct{ ... }
type BundleAnalytics struct{ ... }
type BundleCondition struct{ ... }
type BundleConstraints struct{ ... }
type BundleImprovement struct{ ... }
type BundleInfo struct{ ... }
type BundleItem struct{ ... }
type BundleManager struct{ ... }
    func NewBundleManager() *BundleManager
type BundleOptimization struct{ ... }
type BundleOptimizationMetrics struct{ ... }
type BundlePricing struct{ ... }
type BundleRecommendation struct{ ... }
type BundleRule struct{ ... }
type BundleTemplate struct{ ... }
type BundleType string
    const BundleTypeFixed BundleType = "fixed" ...
type Calculator struct{ ... }
    func NewCalculator() *Calculator
type Customer struct{ ... }
type Dimensions struct{ ... }
type DynamicPricingConfig struct{ ... }
type DynamicPricingRule struct{ ... }
type MarketData struct{ ... }
type PriceAdjustment struct{ ... }
type PriceTier struct{ ... }
type PricedItem struct{ ... }
type PricingAnalytics struct{ ... }
type PricingCondition struct{ ... }
type PricingContext struct{ ... }
type PricingFactor struct{ ... }
type PricingInput struct{ ... }
type PricingItem struct{ ... }
type PricingOptions struct{ ... }
type PricingRecommendation struct{ ... }
type PricingResult struct{ ... }
type PricingRule struct{ ... }
type PricingStrategy string
    const StrategyFixed PricingStrategy = "fixed" ...
type PricingType string
    const PricingTypeBase PricingType = "base" ...
type TierInfo struct{ ... }
type TierPricing struct{ ... }

## Utils Package

package utils // import "github.com/masumrpg/ecommerce-engine/pkg/utils"

func Abs(value float64) float64
func AbsInt(value int) int
func Average(values []float64) float64
func AverageInt(values []int) float64
func Clamp(value, min, max float64) float64
func ClampInt(value, min, max int) int
func CompoundInterest(principal, rate float64, periods int) float64
func Correlation(x, y []float64) float64
func DegreeToRadian(degrees float64) float64
func Distance(x1, y1, x2, y2 float64) float64
func ExponentialDecay(initial, rate, time float64) float64
func ExponentialGrowth(initial, rate, time float64) float64
func ExponentialMovingAverage(values []float64, alpha float64) []float64
func Factorial(n int) int
func Fibonacci(n int) int
func GCD(a, b int) int
func GenerateBase64Token(length int) string
func GenerateChecksum(data string) string
func GenerateCustomID(prefix string, includeTimestamp, includeRandom bool) string
func GenerateHashID(input string) string
func GenerateNonce(length int) string
func GenerateNumericID(length int) string
func GenerateRandomString(length int, charset string) string
func GenerateSalt(length int) string
func GenerateShortID(length int) string
func GenerateTimestampID() string
func GenerateTimestampIDWithPrefix(prefix string) string
func GenerateUUID() string
func InRange(value, min, max float64) bool
func InRangeInt(value, min, max int) bool
func IsEqual(a, b, tolerance float64) bool
func IsPrime(n int) bool
func IsZero(value float64) bool
func LCM(a, b int) int
func LinearInterpolation(x, x1, y1, x2, y2 float64) float64
func LinearRegression(x, y []float64) (slope, intercept float64)
func Logistic(x, k, x0, l float64) float64
func ManhattanDistance(x1, y1, x2, y2 float64) float64
func Max(a, b float64) float64
func MaxInt(a, b int) int
func Median(values []float64) float64
func Min(a, b float64) float64
func MinInt(a, b int) int
func MovingAverage(values []float64, window int) []float64
func NormalizeToRange(value, oldMin, oldMax, newMin, newMax float64) float64
func Percentage(value, total float64) float64
func PercentageChange(oldValue, newValue float64) float64
func PercentageOf(percentage, total float64) float64
func PresentValue(futureValue, discountRate float64, periods int) float64
func RadianToDegree(radians float64) float64
func RandomFloat(min, max float64) float64
func RandomInt(min, max int) int
func RandomIntWithSeed(min, max int, seed int64) int
func Round(value float64, decimals int) float64
func RoundToCurrency(value float64) float64
func RoundToPercent(value float64) float64
func RoundWithMode(value float64, decimals int, mode RoundingMode) float64
func SafeDivide(numerator, denominator float64) float64
func SafeDivideInt(numerator, denominator int) float64
func ScaleToRange(value, min, max float64) float64
func Sigmoid(x float64) float64
func StandardDeviation(values []float64) float64
func Sum(values []float64) float64
func SumInt(values []int) int
func Variance(values []float64) float64
func WeightedAverage(values, weights []float64) float64
type BarcodeGenerator struct{}
    func NewBarcodeGenerator() *BarcodeGenerator
type ColorGenerator struct{}
    func NewColorGenerator() *ColorGenerator
type CouponCodeGenerator struct{ ... }
    func NewCouponCodeGenerator(length int) *CouponCodeGenerator
type IDGenerator struct{ ... }
    func NewIDGenerator(prefix string) *IDGenerator
type PasswordGenerator struct{ ... }
    func NewPasswordGenerator(length int) *PasswordGenerator
type ReferenceGenerator struct{ ... }
    func NewReferenceGenerator(prefix, suffix string, length int) *ReferenceGenerator
type RoundingMode int
    const RoundHalfUp RoundingMode = iota ...
type SlugGenerator struct{}
    func NewSlugGenerator() *SlugGenerator
type TokenGenerator struct{}
    func NewTokenGenerator() *TokenGenerator

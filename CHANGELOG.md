# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of E-commerce Engine
- Currency package with multi-currency support
- Coupon package with validation and calculation
- Discount package with various discount types
- Shipping package with multiple calculation methods
- Tax package with configurable tax rules
- Loyalty package with points and tier management
- Pricing package with bundling support
- Utils package with common utilities
- Comprehensive test coverage
- Go documentation for all packages
- Example usage in examples/main.go

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A

## [1.0.0] - 2025-XX-XX

### Added
- Initial stable release
- Complete e-commerce calculation engine
- Multi-currency support with conversion
- Advanced coupon system with validation rules
- Flexible discount calculations
- Shipping cost calculations
- Tax calculations with multiple rules
- Loyalty points and tier system
- Product bundling and pricing
- Comprehensive documentation
- Unit tests for all packages
- MIT License
- Contributing guidelines

---

## How to Update This Changelog

When making changes to the project:

1. Add new entries under the "Unreleased" section
2. Use the following categories:
   - **Added** for new features
   - **Changed** for changes in existing functionality
   - **Deprecated** for soon-to-be removed features
   - **Removed** for now removed features
   - **Fixed** for any bug fixes
   - **Security** in case of vulnerabilities

3. When releasing a new version:
   - Move items from "Unreleased" to a new version section
   - Add the release date
   - Create a new empty "Unreleased" section

4. Follow semantic versioning:
   - **MAJOR** version for incompatible API changes
   - **MINOR** version for backwards-compatible functionality additions
   - **PATCH** version for backwards-compatible bug fixes
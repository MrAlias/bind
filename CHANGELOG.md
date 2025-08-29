# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-08-28

### Added

- Package-level documentation (`doc.go`) with comprehensive examples
- Comprehensive CI/CD workflows with multi-platform testing (Ubuntu, macOS, Windows)
- Security scanning with CodeQL analysis
- Automated dependency updates via Renovate configuration
- Contributing guidelines and community standards (`CONTRIBUTING.md`)
- Security policy and vulnerability reporting process (`SECURITY.md`)
- Issue templates for bug reports and feature requests
- Code ownership configuration (`.github/CODEOWNERS`)
- Golangci-lint configuration for code quality enforcement
- README improvements with badges for CI, coverage, and Go Report Card
- Benchmark testing in CI pipeline
- Coverage reporting with Codecov integration

## [0.2.0] - 2025-08-27

### Added

- `Meter` function to bind attributes to meters
- Test for unwrap of non-bound instruments

### Fixed

- Do not sort passed attribute slices to avoid side effects
- Early return optimization for no passed attributes

### Changed

- Added benchmarks for all measurement methods

## [0.1.1] - 2025-08-20

### Added

- Support for Go 1.23

## [0.1.0] - 2025-08-19

### Added

- Initial release
- Support for binding attributes to all OpenTelemetry synchronous metric instruments
- `Unwrap` function to retrieve original instruments and bound attributes
- Comprehensive test suite with benchmarks
- Example documentation
- Apache 2.0 license

[Unreleased]: https://github.com/MrAlias/bind/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/MrAlias/bind/releases/tag/v1.0.0
[0.2.0]: https://github.com/MrAlias/bind/releases/tag/v0.2.0
[0.1.1]: https://github.com/MrAlias/bind/releases/tag/v0.1.1
[0.1.0]: https://github.com/MrAlias/bind/releases/tag/v0.1.0

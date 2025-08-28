# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Package-level documentation
- Comprehensive CI/CD workflows
- Security scanning
- Automated dependency updates

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

[Unreleased]: https://github.com/MrAlias/bind/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/MrAlias/bind/releases/tag/v0.2.0
[0.1.1]: https://github.com/MrAlias/bind/releases/tag/v0.1.1
[0.1.0]: https://github.com/MrAlias/bind/releases/tag/v0.1.0

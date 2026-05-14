# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-01-24
### Added
- OAuth login support.
- `--web` flag to open the projects page in a browser.
- Dynamic versioning support for CLI builds.
- CI workflow for test/check.
- Improved UI/UX and display formatting.

### Changed
- Refactored token management and credential handling into dedicated packages.
- Extracted API and display logic into internal packages.
- Updated README and development guidelines.
- Added macOS entries to `.gitignore`.
- Release workflow enhancements (checksums, install instructions).

### Fixed
- Projects API parsing for `projects_users`.
- In-progress project validation status handling.
- Error handling when retrieving stored credentials.

### Security
- CSRF protection in OAuth flow.
- Masked tokens in display output.
- Sanitized sensitive data in API error messages.

## [0.1.0] - 2025-07-22
### Added
- Initial release with authentication (keyring storage).
- User profile and location information.
- Projects and corrections tracking.

[unreleased]: https://github.com/shiftwavedev/42-cli/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/shiftwavedev/42-cli/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/shiftwavedev/42-cli/releases/tag/v0.1.0

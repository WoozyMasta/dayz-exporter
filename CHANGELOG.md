# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][],
and this project adheres to [Semantic Versioning][].

<!--
## Unreleased

### Added
### Changed
### Removed
-->

## [0.3.0][] - 2025-01-16

### Added

* `a2s_info_ping_seconds` metric with A2S_INFO response time in seconds
* read/idle/write timeouts in http server listener
* extend logging
* `.golangci.yml` config and fix linting issues
* automatically turn on colors in log if terminal is used
* simple css style for index page
* 32x32 and 64x64 winres icons for cli

### Changed

* fixed the `game` label previously could use the server description,
  now only the `folder` from `a2s`
* update Grafana dashboard panels
* replaced `github.com/rumblefrog/go-a2s` with `github.com/woozymasta/a2s`
* replaced `github.com/sirupsen/logrus` with `github.com/rs/zerolog`
* internal dependencies for cli are moved to `internal/`
* logs are output to `stdout` rather than `stderr`

[0.3.0]: https://github.com/WoozyMasta/dayz-exporter/compare/v0.2.0...v0.3.0

## [0.2.0][] - 2024-12-22

### Added

* Basic authentication for all HTTP endpoints,
  with `/health*` optional for protection
* SBOM generation and cyclonedx-gomod dev tool dependency
* Added workflow action for check structures alignment

### Changed

* Grafana dashboard misspell
* Align all structs for less memory usage
* Workflow action for VirusTotal scan artifacts replaced with version
  that supports file masking in release

[0.2.0]: https://github.com/WoozyMasta/dayz-exporter/compare/v0.1.1...v0.2.0

## [0.1.1][] - 2024-12-07

### Added

* Windows manifest and icon for binary exe
* Scan release binaries on VirusTotal

### Changed

* Disabled UPX packer for Windows binaries to prevent false
  positives from some antivirus

[0.1.1]: https://github.com/WoozyMasta/dayz-exporter/compare/v0.1.0...v0.1.1

## [0.1.0][] - 2024-12-06

### Added

* First public release

[0.1.0]: https://github.com/WoozyMasta/dayz-exporter/tree/v0.1.0

<!--links-->
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

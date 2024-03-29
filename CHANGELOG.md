# Changelog


## v0.4.2

### Fixed
- Preserve spacing after TOC on overwrite 


## v0.4.1

### Fixed
- Properly handle parenthesis in anchor links


## v0.4.0

### Added
- Version (`-version`) CLI flag
- Config struct to store insert settings
- Flags/configuration to enable optional TOC headings (`-with-toc-heading`, `-toc-heading`)

### Changed
- `Insert` function signature now accepts `*Config` type


## v0.3.0

### Added
- Exclude headings with the `<!--mdtoc: ignore-->` comment

### Changed
- Simplified `Insert` function signature
- Updated TOC begin/end comments with proper formatting


## v0.2.1

### Fixed
- Ignore lines matching heading regex when inside code blocks


## v0.2.0

### Added
- Github action workflow for `go build`, `go test` commands

### Changed
- Removed unneccesary constant exports
- Modified exported function names for clarity
- Updated go library usage docs


## v0.1.2

### Fixed
- Properly handle braces in anchor links


## v0.1.1

### Fixed
- Properly handle backticks, single, and double quotes in anchor links


## v0.1.0

### Added
- Optional `-out` flag to redirect modified content


## v0.0.2

### Added
- Improved usage message formatting

### Fixed
- Properly handle repeated heading text


## v0.0.1
- Initial release

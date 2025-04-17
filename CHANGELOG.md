# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [3.2.0] - 2025-04-17

### Added

- Add `version` command

### Changed

- Use `~/.config/clockify2cats/config.yaml` instead of `~/.clockify2cats.yaml` as config file

## [3.1.0] - 2025-04-16

### Added

- Add "billable" flag to clockify entry
- Use billable flag to distribute the "(*)" entries to billable entries

## [3.0.1] - 2024-10-29

### Fixed

- Increase clockify api limit to 1000 entries

## [3.0.0] - 2024-10-18

### Added

- Accept comma separated list of CATS ids to split tracked times to multiple projects
- Add description delimiter config param to split clockify description for text, text 2 and text external columns, default is '#'

## [2.1.0] - 2024-02-26

### Added

- GitHub actions for CI

### Changed

- Refactor internal code to make it testable

### Fixed

- Fix duplicated console output
- Fix last week parameter if current week is the first of the new year

## [2.0.0] - 2024-02-24

### Changed

- Rebuild cli behavior to include specific subcommands
- Use [cobra](https://github.com/spf13/cobra) for the cli commands
- Use [viper](https://github.com/spf13/viper) for the configuration

# Copilot Instructions

**Go version:** 1.26

## Commands

```sh
# Run all tests (opens coverage HTML report)
make test

# Run a single test
go test ./internal/report/... -run TestReporter_Generate_withTwoTimeEntriesForTheSameProjectAndDescription -v

# Run a single test in cmd package
go test ./cmd/... -run TestGenerateCmd -v

# Format code
make fmt

# Run the tool locally
make run

# Build release snapshot
make release_local
```

Coverage excludes `internal/report/repository` (the Clockify HTTP client) — this is enforced via `.covignore`.

## Architecture

This is a Go CLI tool (Cobra + Viper) that fetches time entries from the Clockify API and formats them for pasting into SAP CATS.

**Data flow:**
1. `cmd/generate.go` parses CLI flags and calls `report.Reporter.Generate()`
2. `internal/report/repository.go` — `Repository.FetchClockifyData()` hits the Clockify REST API (`/time-entries?hydrated=1`)
3. `internal/report/report.go` — `Reporter.convertTimeEntries()` maps Clockify entries → `CatsEntity` structs, then `generateCatsReportData()` formats them as tab-separated output

**Key packages:**
- `cmd/` — Cobra commands: `root`, `generate`, `init`, `version`
- `internal/report/` — core logic: `Reporter`, `Repository`, models, date helpers

**Interfaces for testability:**
- `ReporterInterface` — implemented by `Reporter`, used in `cmd/generate.go` (injected via `newGenerateCmd`)
- `RepositoryInterface` — implemented by `Repository`, used by `Reporter`; tests use a `repositoryMock`

**Configuration** is stored in the OS user config dir (`$HOME/.config/clockify2cats/config.yaml` on Linux/Mac) via Viper. Keys: `workspace-id`, `user-id`, `api-key`, `description-delimiter`.

## `init` Command

Writes config to the OS user config dir via Viper. All three flags are required:

```sh
clockify2cats init \
  --workspace <WorkspaceID> \
  --user <UserID> \
  --api-key <API-KEY> \
  --description-delimiter "#"   # optional, defaults to "#"
```

Flag → Viper key mappings:
- `--workspace` → `workspace-id`
- `--user` → `user-id`
- `--api-key` → `api-key`
- `--description-delimiter` → `description-delimiter`

To get the workspace and user IDs from the API:
```sh
curl -H 'X-Api-Key: <API-KEY>' https://api.clockify.me/api/v1/user | jq '. | {id, defaultWorkspace}'
```

## `generate` Command

Generates a tab-separated CATS report for a given week and prints to stdout.

```sh
# One of these week selectors is required (mutually exclusive):
clockify2cats generate --current          # current ISO week
clockify2cats generate --last             # previous ISO week
clockify2cats generate --week 42          # specific week number (uses previous year if > current week)

# Optional flags:
clockify2cats generate --current --text             # include Text, Text 2, Text External columns
clockify2cats generate --current --copy             # also copy output to clipboard
clockify2cats generate --current --category "ID"   # override category column (default: "ID")
clockify2cats generate --current --month-boundary start      # only entries where start date is in a new month
clockify2cats generate --current --month-boundary end        # only entries where start date is in the current month
```

The `--month-boundary` flag handles weeks that span a month boundary: `start` keeps only entries in the new month, `end` keeps only entries in the current month.

## Key Conventions

**Clockify project naming** drives CATS ID extraction via regex `\((.*)\)`:
- `Project Name (CATS-123)` → single CATS ID
- `Project Name (CATS-1,CATS-2)` → time is split equally across IDs
- `Project Name (*)` → time is distributed proportionally across billable entries

**Description delimiter** (default `#`) splits descriptions into up to 3 CATS text fields (Text, Text 2, Text External). Only populated when `--text` flag is used.

**Time format:** `"2006-01-02T15:04:05.999Z"` (Go reference time). Clockify durations come in ISO 8601 (`PT1H30M`) and are lowercased/stripped of the `PT` prefix before parsing with `time.ParseDuration`.

**CATS output format:** tab-separated columns: `Rec. order`, `(empty)`, `Text`, `Text 2`, `Text External`, `Category`, Mon–Sun hours (formatted with German locale, e.g. `2,00`).

**Release process:** Update `version.txt`, then `make release`. The CI workflow auto-tags and runs GoReleaser when the version in `version.txt` is new.

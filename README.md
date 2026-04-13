[![CI](https://github.com/marvincaspar/clockify-cats-reporter/actions/workflows/ci.yml/badge.svg)](https://github.com/marvincaspar/clockify-cats-reporter/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/marvincaspar/clockify2cats/badge.svg?branch=main)](https://coveralls.io/github/marvincaspar/clockify2cats?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/marvincaspar/clockify2cats)](https://goreportcard.com/report/github.com/marvincaspar/clockify2cats)

# Clockify2Cats

`clockify2cats` is a tool which exports your time from clockify and prints a report into a format that you can copy and paste into SAP CATS (Cross-Application Time Sheet).

![Clockify2Cats usage](./clockify2cats.gif)

## Installation

### Homebrew

```sh
brew tap marvincaspar/tap
brew update
brew install clockify2cats
```

### Github

Download the latest binary for your system from the [GitHub release page](https://github.com/marvincaspar/clockify-cats-reporter/releases/latest/).

```sh
curl -o clockify2cats.tar.gz -L https://github.com/marvincaspar/clockify-cats-reporter/releases/latest/download/clockify2cats_$(uname -s)_$(uname -m).tar.gz
tar -xvzf clockify2cats.tar.gz
sudo mv clockify2cats /usr/local/bin
rm clockify2cats.tar.gz
```

## Usage

### 1. Configure

Run `init` once to store your Clockify credentials:

```sh
clockify2cats init \
  --workspace <WorkspaceID> \
  --user <UserID> \
  --api-key <API-KEY> \
  --description-delimiter "#"   # optional, defaults to "#"
```

Fetch your workspace and user IDs from the Clockify API:

```sh
curl -H 'X-Api-Key: <API-KEY>' https://api.clockify.me/api/v1/user \
  | jq '. | {id, defaultWorkspace}'
```

The configuration is stored in a platform-specific directory:

| OS | Path |
|----|------|
| Linux | `$XDG_CONFIG_HOME/clockify2cats/config.yaml` (or `$HOME/.config/clockify2cats/config.yaml`) |
| macOS | `$HOME/Library/Application Support/clockify2cats/config.yaml` |
| Windows | `%AppData%\clockify2cats\config.yaml` |

### 2. Generate a report

```sh
clockify2cats generate --current          # current ISO week
clockify2cats generate --last             # previous ISO week
clockify2cats generate --week <number>    # specific week number

# Optional flags:
#   -t, --text              include text columns (Text, Text 2, Text External)
#   -C, --copy              copy output to clipboard
#       --category string   override the category column (default "ID")
#   -m, --month-boundary end|start   filter a week that spans a month boundary
```

Example output:

```
$ clockify2cats generate --current

# Rec. order              Category  Mon   Tue   Wed   Thu   Fri   Sat   Sun
CATSID-1                    ID      8.06  4.68  1.26  2.34  7.62  0.00  0.00
CATSID-2                    ID      0.00  0.47  6.02  5.13  0.73  0.00  0.00
CATSID-3                    ID      0.00  2.93  0.00  0.53  0.00  0.00  0.00
```

> The comment line above is for illustration only — actual output is tab-separated with no header row.

Output columns (tab-separated): `Rec. order` · `Description` (empty) · `Text` · `Text 2` · `Text External` · `Category` · Mon–Sun hours

Use `--text` to populate the Text columns from your Clockify entry descriptions (see [Clockify setup](#clockify-setup)).  
Use `--month-boundary end` or `--month-boundary start` to split reporting for weeks that cross a month boundary.

## Clockify setup

### Project naming

Name your Clockify projects using the pattern `<ProjectName> (<CatsID>)`. The CATS ID is extracted from the parentheses.

| Project name | Behaviour |
|---|---|
| `My Project (CATSID-1)` | Maps all time to `CATSID-1` |
| `My Project (CATSID-1, CATSID-2)` | Splits time equally between `CATSID-1` and `CATSID-2` |
| `My Project (*)` | Distributes time proportionally across all other billable entries |

### Description delimiter

Use the description field in Clockify to populate the CATS text columns (only shown with `--text`). Fields are separated by the configured delimiter (default `#`):

| Description | Text | Text 2 | Text External |
|---|---|---|---|
| `Task description` | _(empty)_ | `Task description` | _(empty)_ |
| `Task # Detail` | `Task` | `Detail` | _(empty)_ |
| `Task # Detail # External` | `Task` | `Detail` | `External` |

### Proportional time distribution (`*`)

When a project is named `SharedProject (*)`, its recorded hours are distributed **proportionally** across all other entries marked as `billable=true`, weighted by hours already logged.

**Example:** 9 hours logged to `SharedProject (*)` with three other billable projects (3 h each):

| CATS ID | Hours |
|---|---|
| CATSID1 | 3.00 + 3.00 (shared) = **6.00** |
| CATSID2 | 3.00 + 3.00 (shared) = **6.00** |
| CATSID3 | 3.00 + 3.00 (shared) = **6.00** |
| CATSID4 _(not billable)_ | 0.00 |

## Release

```sh
brew install goreleaser
goreleaser init
```

Update the version in `version.txt` and run `make release`.
This will create a new git tag, build all binaries and publish it to GitHub.


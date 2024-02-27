[![CI](https://github.com/marvincaspar/clockify-cats-reporter/actions/workflows/ci.yml/badge.svg)](https://github.com/marvincaspar/clockify-cats-reporter/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/marvincaspar/clockify2cats/badge.svg?branch=main)](https://coveralls.io/github/marvincaspar/clockify2cats?branch=main)

# Clockify2Cats

`clockify2cats` is a tool which exports your time from clockify and prints a report into a format that you can capy and paste into SAP CATS (Cross-Application Time Sheet).

![Clockify2Cats usage](./clockify2cats.gif)

## Installation

Download the latest binary for your system from the [GitHub release page](https://github.com/marvincaspar/clockify-cats-reporter/releases/latest/).

```sh
curl -o clockify2cats.tar.gz -L https://github.com/marvincaspar/clockify-cats-reporter/releases/latest/download/clockify2cats_$(uname -s)_$(uname -m).tar.gz
tar -xvzf clockify2cats.tar.gz
sudo mv clockify2cats /usr/local/bin
rm clockify2cats.tar.gz
```

## Usage

First you need to setup your local configuration. Run `clockify2cats init --workspace <WorkspaceID> --user <UserID> --api-key <API-KEY>`. The configuration is stored in `~/.clockify2cats.yaml`.

Then you can use `clockify2cats` to generate a report. Run `clockify2cats generate [flags]`.

Flags:
```
      --category string   Category identifyer (default "ID")
  -C, --copy              Copy report to clipboard
  -c, --current           Current week
  -h, --help              help for generate
  -l, --last              Last week
  -t, --text              Print with text
  -w, --week int          Week number
```

Example: 
```
$ clockify2cats generate --current

CATSID-1                    ID      8.06            4.68            1.26            2.34            7.62            0.00            0.00
CATSID-2                    ID      0.00            0.47            6.02            5.13            0.73            0.00            0.00
CATSID-3                    ID      0.00            2.93            0.00            0.53            0.00            0.00            0.00
```

This report is build for the CATS columns:

- Rec. order
- Description - empty
- Text - use flag `-t` to use it
- Category - default `ID`, can be set with flag `--category <string>`
- Monday
- Tuesday
- Wensday
- Thursday
- Friday
- Saturday
- Sunday


## Clockify setup

Create projects and name it like `<ProjectName> (<CatsID>)`. Track your time for the projects. Use the clockify description field to add additional information for the CATS text field.

Generate an API for clockify. It can be fount in your profile settings.

Fetch your user id and your default workspace id from the api `curl -H 'X-Api-Key: <API-KEY>' https://api.clockify.me/api/v1/user | jq`.


## Release

```sh
brew install goreleaser
goreleaser init
```

Update the version in `version.txt` and run `make release`.
This will create a new git tag, build all binaries and publish it to github.


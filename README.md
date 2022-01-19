# Clockify2Cats

`clockify2cats` is a tool which exports your time from clockify and prints a report into a format that you can capy and paste into SAP CATS (Cross-Application Time Sheet).

![Clockify2Cats usage](./clockify2cats.gif)

## Usage

To generate the report run `clockify2cats` with the following arguments:

Usage of clockify2cats:
-  -w, --week week number for report (don't use in combination with start)
-  -s, --start Startdate for report YYYY-MM-DD (don't use in combination with week)
-  -C, --copy Copy report to clipboard
-  -c, --category Category identifyer (default: ID)
-  -t, --text Add Clockify description as text to report

Example: 
```
$ clockify2cats --week 2

CATSID-1                    ID      8.06            4.68            1.26            2.34            7.62            0.00            0.00
CATSID-2                    ID      0.00            0.47            6.02            5.13            0.73            0.00            0.00
CATSID-3                    ID      0.00            2.93            0.00            0.53            0.00            0.00            0.00
```

This report is build for the CATS columns:

- Rec. order
- Description - empty
- Text - use flag `-t` to use it
- Category - default `ID`, can be set with flag `-c`
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

Set your workspace id, user id and api key and set them as environment variable or create a `.env` file with the following variables:

```
CLOCKIFY_WORKSPACE_ID=xxxxxxxxxxxxxxxxxxxxxxxx
CLOCKIFY_USER_ID=xxxxxxxxxxxxxxxxxxxxxxxx
CLOCKIFY_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Release

```sh
brew install goreleaser
goreleaser init
```

Update the version in `version.txt` and run `make release`.
This will create a new git tag, build all binaries and publish it to github.


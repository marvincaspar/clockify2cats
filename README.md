# Clockify2Cats

This is a command line tool which exports your time from clockify and prints a report into a format that you can use for SAP CATS (Cross-Application Time Sheet).

This report is build for the following CATS columns:

- Rec. order
- Description - empty
- Text - currently empty
- Category - currently always ID
- Monday
- Tuesday
- Wensday
- Thursday
- Friday
- Saturday
- Sunday

## Clockify

Tracking time in clockify and follow one of the two conventions:

1. Create a project and name it like `<ProjectName> (<CatsID>)` (recommended)
2. Empty Project and use the description like `<CatsID> (<ProjectName>)`


Get your workspace id, user id and api key and set them as environment vaiable or create a `.env` file with the following variables:

```
CLOCKIFY_WORKSPACE_ID=xxxxxxxxxxxxxxxxxxxxxxxx
CLOCKIFY_USER_ID=xxxxxxxxxxxxxxxxxxxxxxxx
CLOCKIFY_API_KEY=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Generate report

To generate the report run the cli tool with one of the following arguments:

1. `--start <DATE>` where date is the first day of the week with the format `YYYY-MM-dd`
2. `--week <WEEK>` where week is the calendar week number

```

CATSID-1                    ID      8.06            4.68            1.26            2.34            7.62            0.00            0.00
CATSID-2                    ID      0.00            0.47            6.02            5.13            0.73            0.00            0.00
CATSID-3                    ID      0.00            2.93            0.00            0.53            0.00            0.00            0.00

```


## Release

```sh
brew install goreleaser
goreleaser init
```

Update the version in `version.txt` and run `make release`.
This will create a new git tag, build all binaries and publish it to github.


# Issue From Template

This action opens a new issue from an issue template. It parses the template's front matter and the body, then posts [an API request to open an issue](https://developer.github.com/v3/issues/#create-an-issue). Works best with a [scheduled workflow](https://developer.github.com/actions/managing-workflows/creating-and-cancelling-a-workflow/#scheduling-a-workflow) and the [Auto Closer](https://github.com/lowply/auto-closer) action.

## Environment variables

- `IFT_TEMPLATE_NAME` (*required*): The name of the issue template. For example, `report.md`. This action will look for the file in the `.github/ISSUE_TEMPLATE` directory.

## Available template variables

- `.Year`: Year of the week
- `.WeekStartDate`: Date of Monday of the week
- `.WeekEndDate`: Date of Sunday of the week
- `.WeekNumber`: ISO week number
- `.Dates`: Array of the dates of the week (Can be used as `{{ index .Dates 1 }}` in the template)

## Template example

```
---
name: Weekly Report
about: This is an example
title: 'Report on {{ .WeekStartDate }} (Week {{ .WeekNumber }}, {{ .Year }})'
labels: report
assignees: lowply
---

# This week's updates!

## {{ index .Dates 0 }} MON
## {{ index .Dates 1 }} TUE
## {{ index .Dates 2 }} WED
## {{ index .Dates 3 }} THU
## {{ index .Dates 4 }} FRI
## {{ index .Dates 5 }} SAT
## {{ index .Dates 6 }} SUN
```

## After post

Once the API call succeeds, the [API response](https://developer.github.com/v3/issues/#response-3) will be saved at `$HOME/resp.json` so you can use it in the next action.

## Running locally for development

This is designed to be used as a GitHub Action, but you can also just run it locally with the following env vars:

```
export GITHUB_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export GITHUB_REPOSITORY="owner/repository"
export GITHUB_WORKSPACE="/path/to/your/local/repository"
export IFT_TEMPLATE_NAME="issue.md"
go run src/main.go
```

# Issue From Template

This action opens a new issue from an issue template. It parses the template's front matter and the body, then posts [an API request to open an issue](https://developer.github.com/v3/issues/#create-an-issue). Works best with a [scheduled workflow](https://help.github.com/en/actions/automating-your-workflow-with-github-actions/events-that-trigger-workflows#scheduled-events-schedule) and the [Auto Closer](https://github.com/lowply/auto-closer) action.

## Environment variables

- `IFT_TEMPLATE_NAME` (*required*): The name of the issue template. For example, `report.md`. This action will look for the file in the `.github/ISSUE_TEMPLATE` directory.
- `ADD_DATES` (*optional*): Number of the dates to add. This is useful when you want to run this action to open an issue for the next week, not this week.

## Available template variables

- `.Year`: Year of the day when this action runs
- `.Month`: Month of the day when this action runs
- `.Day`: Day when this action runs
- `.WeekStartDate`: Date of Monday of the week (YYYY/MM/DD)
- `.WeekEndDate`: Date of Sunday of the week (YYYY/MM/DD)
- `.WeekNumber`: ISO week number
- `.WeekNumberYear`: Year of the Thursday of the week. Matches with [ISO week number](https://en.wikipedia.org/wiki/ISO_week_date#First_week)
- `.Dates`: Array of the dates of the week (Can be used as `{{ index .Dates 1 }}` in the template)

## Template example

```
---
name: Weekly Report
about: This is an example
title: 'Report for Week {{ .WeekNumber }}, {{ .WeekNumberYear }} (Week of {{ .WeekStartDate }})'
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

## Default comments

If the *.github/ift-comments.yaml* file exists, it also parses the content of the file and posts comments after creating the issue. This is useful for teams to have default comment to the issue. Here's an example of the comments in the YAML format:

```
- comment: |
    ## Sales
    Hello :wave: from the Sales team! Here's the [link](http://example.com) to the latest numbers.
- comment: |
    ## Support
    - Tickets from company A: [URL](http://example.com)
    - Tickets from company B: [URL](http://example.com)
    - High priority tickets: [URL](http://example.com)
- comment: |
    ## Workplace
    Hi everyone! Here's the latest news from us:
```

## Running locally for development

This is designed to be used as a GitHub Action, but you can also just run it locally with the following env vars:

```
export GITHUB_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export GITHUB_REPOSITORY="owner/repository"
export GITHUB_WORKSPACE="/path/to/your/local/repository"
export IFT_TEMPLATE_NAME="issue.md"
go run src/main.go
```

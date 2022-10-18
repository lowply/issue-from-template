# Issue From Template

This action opens a new issue from a markdown template file. It parses the template's front matter and the body, then posts [an API request to open an issue](https://docs.github.com/en/rest/issues/issues#create-an-issue). Works best with a [scheduled workflow](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#schedule) and the [Auto Closer](https://github.com/lowply/auto-closer) action.

Note that "template" here is just a markdown file, not [GitHub's issue template feature](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/about-issue-and-pull-request-templates) which uses yaml file.

## Environment variables

- `IFT_TEMPLATE_NAME` (*required*): The name of the issue template. For example, `report.md`. This action will look for the template file in the `.github` directory. **Note that it will not look for any sub directories** including the `.github/ISSUE_TEMPLATE` directory.
- `ADD_DATES` (*optional*): Number of the dates to add. This is useful when you want to run this action to open an issue for the next week, not this week.

## Available template variables

- `.Current`: The day when this action runs (time.Time)
- `.WeekStart`: Date of Monday of the week (time.Time)
- `.WeekEnd`: Date of Sunday of the week (time.Time)
- `.WeekNumber`: ISO week number (string)
- `.YearOfTheWeek`: Year of the Thursday of the week. Matches with [ISO week number](https://en.wikipedia.org/wiki/ISO_week_date#First_week) (string)
- `.Dates`: Array of the dates of the week (Can be used as `{{ index .Dates 1 }}` in the template, array of time.Time)

For variables that are in the `time.Time` type, you can pick your preferred format in the template e.g. `.Format "2006-01-02"`. See https://pkg.go.dev/time#Time.Format for more details.

### Date and time formatting layout

If you're not familiar with Go's time.Time layouts, there are other resources you can use e.g. [Date and time format in Go (Golang) cheatsheet](https://gosamples.dev/date-time-format-cheatsheet/) but in short, **Monday, Jan 2nd, 2006** is the day used to express any formatting. So for example, if you want to format your date in `YYYY-MM-DD`, the format would be `2006-01-02`.

## Template example

```
---
name: Weekly Report
about: This is an example
title: 'Report for Week {{ .WeekNumber }}, {{ .YearOfTheWeek }} (Week of {{ .WeekStartDate.Format "2006/01/02" }})'
labels: report
assignees: lowply
---

# This week's updates!

## {{ (index .Dates 0).Format "01/02 Mon" }}
## {{ (index .Dates 1).Format "01/02 Mon" }}
## {{ (index .Dates 2).Format "01/02 Mon" }}
## {{ (index .Dates 3).Format "01/02 Mon" }}
## {{ (index .Dates 4).Format "01/02 Mon" }}
## {{ (index .Dates 5).Format "01/02 Mon" }}
## {{ (index .Dates 6).Format "01/02 Mon" }}
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
cd src
export GITHUB_TOKEN="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
export GITHUB_REPOSITORY="owner/repository"
export GITHUB_WORKSPACE="/path/to/your/local/repository"
export IFT_TEMPLATE_NAME="issue.md"
go run .
```

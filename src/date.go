package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

type date struct {
	Year           string
	WeekStartDate  string
	WeekNumber     string
	WeekNumberYear string
	Dates          [7]string
}

func NewDate(t time.Time) *date {
	d := &date{}

	nc := &now.Config{
		WeekStartDay: time.Monday,
	}
	n := nc.With(t)

	d.Year = n.BeginningOfYear().Format("2006")
	d.WeekStartDate = n.BeginningOfWeek().Format("01/02")
	d.WeekStartDate = n.BeginningOfWeek().Format("2006/01/02")
	_, isoweek := n.Monday().ISOWeek()
	d.WeekNumber = fmt.Sprintf("%02d", isoweek)
	for j := range d.Dates {
		d.Dates[j] = n.Monday().AddDate(0, 0, j).Format("01/02")
	}
	// Thursday of the week, should be used with the week number
	// e.g. "2020 Week 01".
	// See https://en.wikipedia.org/wiki/ISO_week_date#First_week
	// for the ISO 8601 first week definition
	d.WeekNumberYear = n.BeginningOfWeek().AddDate(0, 0, 3).Format("2006")
	return d
}

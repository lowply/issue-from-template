package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

type data struct {
	Current       time.Time
	WeekStart     time.Time
	WeekEnd       time.Time
	WeekNumber    string
	YearOfTheWeek string
	Dates         [7]time.Time
}

func NewData(t time.Time) *data {
	d := &data{}

	nc := &now.Config{
		WeekStartDay: time.Monday,
	}
	n := nc.With(t)

	// https://github.com/jinzhu/now#mondaysunday
	d.Current = t
	d.WeekStart = n.Monday()
	d.WeekEnd = n.Sunday()

	_, isoweek := n.Monday().ISOWeek()
	d.WeekNumber = fmt.Sprintf("%02d", isoweek)

	// Thursday of the week, should be used with the week number
	// e.g. "2020 Week 01".
	// See https://en.wikipedia.org/wiki/ISO_week_date#First_week
	// for the ISO 8601 first week definition
	d.YearOfTheWeek = n.BeginningOfWeek().AddDate(0, 0, 3).Format("2006")

	for j := range d.Dates {
		d.Dates[j] = n.Monday().AddDate(0, 0, j)
	}
	return d
}

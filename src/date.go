package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

type date struct {
	Year          string
	WeekStartDate string
	WeekEndDate   string
	WeekNumber    string
	Dates         [7]string
}

func NewDate(t time.Time) *date {
	d := &date{}

	nc := &now.Config{
		WeekStartDay: time.Monday,
	}
	n := nc.With(t)

	d.Year = n.BeginningOfYear().Format("2006")
	d.WeekEndDate = n.EndOfSunday().Format("01/02")
	d.WeekStartDate = n.BeginningOfWeek().Format("01/02")
	_, isoweek := n.Monday().ISOWeek()
	d.WeekNumber = fmt.Sprintf("%02d", isoweek)
	for j := range d.Dates {
		d.Dates[j] = n.Monday().AddDate(0, 0, j).Format("01/02")
	}
	return d
}

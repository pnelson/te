package te

import (
	"sort"
	"time"
)

// Expression represents a temporal expression.
type Expression interface {
	IsActive(t time.Time) bool
	Next(t time.Time) time.Time
}

type hourExpr int

func (expr hourExpr) IsActive(t time.Time) bool {
	return t.Hour() == int(expr)
}

func (expr hourExpr) Next(t time.Time) time.Time {
	year, month, day := t.Date()
	loc := t.Location()
	next := time.Date(year, month, day, int(expr), 0, 0, 0, loc)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

type minuteExpr int

func (expr minuteExpr) IsActive(t time.Time) bool {
	return t.Minute() == int(expr)
}

func (expr minuteExpr) Next(t time.Time) time.Time {
	year, month, day := t.Date()
	hour := t.Hour()
	loc := t.Location()
	next := time.Date(year, month, day, hour, int(expr), 0, 0, loc)
	if t.Equal(next) || t.After(next) {
		next = next.Add(time.Hour)
	}
	return next
}

type secondExpr int

func (expr secondExpr) IsActive(t time.Time) bool {
	return t.Second() == int(expr)
}

func (expr secondExpr) Next(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, _ := t.Clock()
	loc := t.Location()
	next := time.Date(year, month, day, hour, min, int(expr), 0, loc)
	if t.Equal(next) || t.After(next) {
		next = next.Add(time.Minute)
	}
	return next
}

type dayExpr int

func (expr dayExpr) IsActive(t time.Time) bool {
	return t.Day() == int(expr)
}

func (expr dayExpr) Next(t time.Time) time.Time {
	loc := t.Location()
	next := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, loc)
	if int(expr) > 0 {
		next = next.AddDate(0, 0, int(expr)-1)
	} else {
		next = next.AddDate(0, 1, -1)
	}
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(0, 1, 0)
		if !expr.IsActive(next) {
			return expr.Next(next)
		}
	}
	return next
}

type weekdayExpr time.Weekday

func (expr weekdayExpr) IsActive(t time.Time) bool {
	return t.Weekday() == time.Weekday(expr)
}

func (expr weekdayExpr) Next(t time.Time) time.Time {
	loc := t.Location()
	year, month, day := t.Date()
	next := time.Date(year, month, day, 0, 0, 0, 0, loc)
	days := int(time.Weekday(expr) - t.Weekday())
	if days <= 0 {
		days += 7
	}
	return next.AddDate(0, 0, days)
}

type monthExpr time.Month

func (expr monthExpr) IsActive(t time.Time) bool {
	return t.Month() == time.Month(expr)
}

func (expr monthExpr) Next(t time.Time) time.Time {
	loc := t.Location()
	next := time.Date(t.Year(), time.Month(expr), 1, 0, 0, 0, 0, loc)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(1, 0, 0)
	}
	return next
}

type yearExpr int

func (expr yearExpr) IsActive(t time.Time) bool {
	return t.Year() == int(expr)
}

func (expr yearExpr) Next(t time.Time) time.Time {
	return time.Time{}
}

type dateRangeExpr struct {
	t1 time.Time
	t2 time.Time
}

func (expr dateRangeExpr) IsActive(t time.Time) bool {
	t = dateFrom(t, t)
	return isBetween(t, expr.t1, expr.t2)
}

func (expr dateRangeExpr) Next(t time.Time) time.Time {
	next := dateFrom(t, expr.t1)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(1, 0, 0)
	}
	return next
}

type timeRangeExpr struct {
	t1 time.Time
	t2 time.Time
}

func (expr timeRangeExpr) IsActive(t time.Time) bool {
	t1 := timeFrom(t, expr.t1)
	t2 := timeFrom(t, expr.t2)
	return isBetween(t, t1, t2)
}

func (expr timeRangeExpr) Next(t time.Time) time.Time {
	next := timeFrom(t, expr.t1)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

type unionExpr []Expression

func (expr unionExpr) IsActive(t time.Time) bool {
	for _, e := range expr {
		if e.IsActive(t) {
			return true
		}
	}
	return false
}

func (expr unionExpr) Next(t time.Time) time.Time {
	ts := make(byTime, len(expr))
	for i, e := range expr {
		ts[i] = e.Next(t)
	}
	sort.Sort(ts)
	return ts[0]
}

type intersectExpr []Expression

func (expr intersectExpr) IsActive(t time.Time) bool {
	for _, e := range expr {
		if !e.IsActive(t) {
			return false
		}
	}
	return true
}

func (expr intersectExpr) Next(t time.Time) time.Time {
	ts := make(byTime, 0)
	for _, e := range expr {
		next := e.Next(t)
		if next.IsZero() {
			if !e.IsActive(t) {
				return time.Time{}
			}
			continue
		}
		ts = append(ts, next)
	}
	sort.Sort(ts)
	t = ts[0]
	if expr.IsActive(t) {
		return t
	}
	return expr.Next(t)
}

type exceptExpr []Expression

func (expr exceptExpr) IsActive(t time.Time) bool {
	for _, e := range expr {
		if e.IsActive(t) {
			return false
		}
	}
	return true
}

func (expr exceptExpr) Next(t time.Time) time.Time {
	ts := make(byTime, len(expr))
	for i, e := range expr {
		ts[i] = e.Next(t)
	}
	sort.Sort(ts)
	return ts[0]
}

type nilExpr struct{}

func (expr nilExpr) IsActive(t time.Time) bool  { return false }
func (expr nilExpr) Next(t time.Time) time.Time { return time.Time{} }

type byTime []time.Time

func (ts byTime) Len() int           { return len(ts) }
func (ts byTime) Less(i, j int) bool { return ts[i].Before(ts[j]) }
func (ts byTime) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }

func dateFrom(t, date time.Time) time.Time {
	loc := date.Location()
	_, month, day := date.Date()
	return time.Date(t.Year(), month, day, 0, 0, 0, 0, loc)
}

func timeFrom(date, clock time.Time) time.Time {
	loc := clock.Location()
	year, month, day := date.Date()
	hour, min, sec := clock.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, loc)
}

func isBetween(t, from, to time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Equal(to) || t.Before(to))
}

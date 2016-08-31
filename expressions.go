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

type dayExpr struct {
	day int
	loc *time.Location
}

func (expr dayExpr) IsActive(t time.Time) bool {
	return t.In(expr.loc).Day() == expr.day
}

func (expr dayExpr) Next(t time.Time) time.Time {
	t = t.In(expr.loc)
	next := time.Date(t.Year(), t.Month(), expr.day, 0, 0, 0, 0, expr.loc)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(0, 1, 0)
	}
	return next
}

type timeExpr struct {
	from time.Time
	to   time.Time
}

func (expr timeExpr) IsActive(t time.Time) bool {
	t = timeOnly(t)
	return isBetween(t, expr.from, expr.to)
}

func (expr timeExpr) Next(t time.Time) time.Time {
	loc := expr.from.Location()
	year, month, day := t.Date()
	hour, min, sec := expr.from.Clock()
	next := time.Date(year, month, day, hour, min, sec, 0, loc)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(0, 0, 1)
	}
	return next
}

type weekdayExpr struct {
	weekday time.Weekday
	count   int
	loc     *time.Location
}

func (expr weekdayExpr) IsActive(t time.Time) bool {
	t = t.In(expr.loc)
	if t.Weekday() != expr.weekday {
		return false
	}
	if expr.count > 0 {
		return weekInMonth(t) == expr.count
	} else if expr.count < 0 {
		return weekInMonthFromEnd(t) == -expr.count
	}
	return true
}

func (expr weekdayExpr) Next(t time.Time) time.Time {
	t = t.In(expr.loc)
	t = t.Truncate(24 * time.Hour)
	t = expr.next(t)
	if expr.count > 0 {
		for weekInMonth(t) != expr.count {
			t = expr.next(t)
		}
	} else if expr.count < 0 {
		for weekInMonthFromEnd(t) != -expr.count {
			t = expr.next(t)
		}
	}
	return t
}

func (expr weekdayExpr) next(t time.Time) time.Time {
	days := int(expr.weekday - t.Weekday())
	if days <= 0 {
		days += 7
	}
	return t.AddDate(0, 0, days)
}

type monthExpr struct {
	month time.Month
	loc   *time.Location
}

func (expr monthExpr) IsActive(t time.Time) bool {
	return t.In(expr.loc).Month() == expr.month
}

func (expr monthExpr) Next(t time.Time) time.Time {
	t = t.In(expr.loc)
	next := time.Date(t.Year(), expr.month, t.Day(), 0, 0, 0, 0, expr.loc)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(1, 0, 0)
	}
	return next
}

type dateRangeExpr struct {
	from time.Time
	to   time.Time
}

func (expr dateRangeExpr) IsActive(t time.Time) bool {
	t = dateOnly(t)
	return isBetween(t, expr.from, expr.to)
}

func (expr dateRangeExpr) Next(t time.Time) time.Time {
	loc := expr.from.Location()
	_, month, day := expr.from.Date()
	next := time.Date(t.Year(), month, day, 0, 0, 0, 0, loc)
	if t.Equal(next) || t.After(next) {
		next = next.AddDate(1, 0, 0)
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
	ts := make(byTime, len(expr))
	for i, e := range expr {
		ts[i] = e.Next(t)
	}
	sort.Sort(ts)
	for _, next := range ts {
		if expr.IsActive(next) {
			return next
		}
	}
	t = ts[len(ts)-1]
	for i, e := range expr {
		ts[i] = e.Next(t)
	}
	sort.Sort(ts)
	for _, next := range ts {
		if expr.IsActive(next) {
			return next
		}
	}
	return time.Time{}
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
	return time.Time{}
}

type nilExpr struct{}

func (expr nilExpr) IsActive(t time.Time) bool  { return false }
func (expr nilExpr) Next(t time.Time) time.Time { return time.Time{} }

type byTime []time.Time

func (ts byTime) Len() int           { return len(ts) }
func (ts byTime) Less(i, j int) bool { return ts[i].Before(ts[j]) }
func (ts byTime) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }

func dateOnly(t time.Time) time.Time {
	loc := t.Location()
	_, month, day := t.Date()
	return time.Date(1, month, day, 0, 0, 0, 0, loc)
}

func timeOnly(t time.Time) time.Time {
	loc := t.Location()
	hour, min, sec := t.Clock()
	return time.Date(1, 1, 1, hour, min, sec, 0, loc)
}

func isBetween(t, from, to time.Time) bool {
	return (t.Equal(from) || t.After(from)) && (t.Equal(to) || t.Before(to))
}

func weekInMonth(t time.Time) int {
	day := t.Day()
	return weekInMonthFromDay(day)
}

func weekInMonthFromEnd(t time.Time) int {
	day := daysFromMonthEnd(t)
	return weekInMonthFromDay(day)
}

func weekInMonthFromDay(day int) int {
	return ((day - 1) / 7) + 1
}

func daysFromMonthEnd(t time.Time) int {
	day := t.Day()
	return t.AddDate(0, 1, -day).Day() - day
}

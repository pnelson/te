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
	from := timeFrom(t, expr.from)
	to := timeFrom(t, expr.to)
	return isBetween(t, from, to)
}

func (expr timeExpr) Next(t time.Time) time.Time {
	next := timeFrom(t, expr.from)
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
	year, month, day := t.In(expr.loc).Date()
	t = time.Date(year, month, day, 0, 0, 0, 0, expr.loc)
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
	t = dateFrom(t, t)
	return isBetween(t, expr.from, expr.to)
}

func (expr dateRangeExpr) Next(t time.Time) time.Time {
	next := dateFrom(t, expr.from)
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
	// Find each next occurrence from the given time.
	ts := make([]time.Time, len(expr))
	for i, e := range expr {
		next := e.Next(t)
		ts[i] = next
	}
	sort.Sort(sort.Reverse(byTime(ts)))
	// Choose the latest time to be the earliest possible intersection.
	t = ts[0]
	// Find the durations to the next occurrence from this new time.
	ds := make([]time.Duration, len(expr))
	for i, e := range expr {
		next := e.Next(t)
		ds[i] = next.Sub(t)
	}
	sort.Sort(sort.Reverse(byDuration(ds)))
	// Enumerate candidate intersection times by applying duration subsets.
	ts = []time.Time{t}
	for _, d := range ds[1:] {
		ss := make([]time.Time, len(ts))
		for i, t := range ts {
			ss[i] = t.Add(d)
		}
		ts = append(ts, ss...)
	}
	sort.Sort(byTime(ts))
	// Return the first active intersection time.
	for _, t := range ts {
		if expr.IsActive(t) {
			return t
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

type byDuration []time.Duration

func (ds byDuration) Len() int           { return len(ds) }
func (ds byDuration) Less(i, j int) bool { return ds[i] < ds[j] }
func (ds byDuration) Swap(i, j int)      { ds[i], ds[j] = ds[j], ds[i] }

func dateFrom(t, date time.Time) time.Time {
	loc := date.Location()
	_, month, day := date.Date()
	return time.Date(t.Year(), month, day, 0, 0, 0, 0, loc)
}

func timeFrom(date, clock time.Time) time.Time {
	loc := clock.Location()
	year, month, day := date.Date()
	hour, min, _ := clock.Clock()
	return time.Date(year, month, day, hour, min, 0, 0, loc)
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

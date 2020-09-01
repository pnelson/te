// Package te implements temporal expressions.
package te

import (
	"time"
)

// Hour returns a temporal expression for an hour.
// If hour is negative or greater than 23, the nil expression is returned.
func Hour(hour int) Expression {
	if hour < 0 || hour > 23 {
		return nilExpr{}
	}
	return hourExpr(hour)
}

// Minute returns a temporal expression for a minute.
// If min is negative or greater than 59, the nil expression is returned.
func Minute(min int) Expression {
	if min < 0 || min > 59 {
		return nilExpr{}
	}
	return minuteExpr(min)
}

// Second returns a temporal expression for a second.
// If sec is negative or greater than 59, the nil expression is returned.
func Second(sec int) Expression {
	if sec < 0 || sec > 59 {
		return nilExpr{}
	}
	return secondExpr(sec)
}

// Day returns a temporal expression for a day of the month.
// Months without the nth day are ignored. If n is -1, the expression
// represents the last day of the month. If n is greater than 31 or less
// than -1, the nil expression is returned.
func Day(n int) Expression {
	if n < -1 || n == 0 || n > 31 {
		return nilExpr{}
	}
	return dayExpr(n)
}

// Weekday returns a temporal expression for weekdays.
func Weekday(d time.Weekday) Expression {
	return weekdayExpr(d)
}

// Month returns a temporal expression for months of the year.
func Month(month time.Month) Expression {
	return monthExpr(month)
}

// Year returns a temporal expression for the given year.
func Year(year int) Expression {
	return yearExpr(year)
}

// Date returns a temporal expression for a date.
func Date(month time.Month, day int) Expression {
	me := Month(month)
	de := Day(day)
	return Intersect(me, de)
}

// Time returns a temporal expression for a time.
func Time(hour, min, sec int) Expression {
	he := Hour(hour)
	me := Minute(min)
	se := Second(sec)
	return Intersect(he, me, se)
}

// DateRange returns a temporal expression for an inclusive date range.
func DateRange(m1 time.Month, d1 int, m2 time.Month, d2 int) Expression {
	t1 := time.Date(1, m1, d1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(1, m2, d2, 0, 0, 0, 0, time.UTC)
	return dateRangeExpr{t1, t2}
}

// TimeRange returns a temporal expression for an inclusive time range.
func TimeRange(h1, m1, s1, h2, m2, s2 int) Expression {
	t1 := time.Date(1, 1, 1, h1, m1, s1, 0, time.UTC)
	t2 := time.Date(1, 1, 1, h2, m2, s2, 0, time.UTC)
	return timeRangeExpr{t1, t2}
}

// Union returns a temporal expression that represents the union
// of the provided expressions. This expression is active when
// any of the given expressions are active.
func Union(exprs ...Expression) Expression {
	if len(exprs) == 0 {
		return nilExpr{}
	}
	return unionExpr(exprs)
}

// Intersect returns a temporal expression that represents the
// intersection of the provided expressions. This expression is
// active when all of the given expressions are active.
func Intersect(exprs ...Expression) Expression {
	if len(exprs) == 0 {
		return nilExpr{}
	}
	return intersectExpr(exprs)
}

// Except returns a temporal expression that represents exceptions.
// This expression does not have a next active time. Compose with an
// intersection expression to represent the difference. This expression
// is active when none of the given expressions are active.
func Except(exprs ...Expression) Expression {
	if len(exprs) == 0 {
		return nilExpr{}
	}
	return exceptExpr(exprs)
}

// Iter returns a receive-only channel of next active times
// for the given expression. The channel is closed when the
// expression returns a zero time.
func Iter(expr Expression, t time.Time) <-chan time.Time {
	ch := make(chan time.Time)
	go func() {
		for {
			t = expr.Next(t)
			if t.IsZero() {
				close(ch)
				return
			}
			ch <- t
		}
	}()
	return ch
}

// Until returns the duration until the next occurrence of t.
func Until(expr Expression, t time.Time) time.Duration {
	next := expr.Next(t)
	return next.Sub(t)
}

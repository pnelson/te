// Package te implements temporal expressions.
package te

import "time"

// Day returns a temporal expression for a day of the month.
func Day(n int) Expression {
	return dayExpr(n)
}

// Weekday returns a temporal expression for weekdays.
// If n is zero, the expression represents every given weekday.
// If n is positive, the expression represents the nth given weekday.
// If n is negative, the expression represents the nth last given weekday.
func Weekday(w time.Weekday, n int) Expression {
	return weekdayExpr{w, n}
}

// Month returns a temporal expression for months of the year.
func Month(m time.Month) Expression {
	return monthExpr(m)
}

// DateRange returns a temporal expression for a date range.
// Only the month, day and location are considered.
func DateRange(t1, t2 time.Time) Expression {
	return dateRangeExpr{t1, t2}
}

// TimeRange returns a temporal expression for a time range.
// Only the time and location are considered.
func TimeRange(t1, t2 time.Time) Expression {
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

// Next returns the next active time for the given expression
// in the given time zone.
func Next(expr Expression, loc *time.Location) time.Time {
	now := time.Now()
	next := expr.Next(now)
	return next.In(loc)
}

// Until returns the duration until the next occurrence of t.
func Until(expr Expression, t time.Time) time.Duration {
	next := expr.Next(t)
	return next.Sub(t)
}

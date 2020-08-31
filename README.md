# te

Package te implements temporal expressions.

An expression is an implementation of the `Expression` interface. Two methods
must be implemented: `IsActive` which returns true if the expression is active
at the given time, and `Next` which returns the next active time.

A handful of expressions are included. They compute the next time within the
location of the given time. Daylight savings is respected.

To find the first day of the month:

```go
expr := te.Day(1)
next := time.Now()
for i := 0; i < 4; i++ {
  next = expr.Next(next)
  fmt.Println(next)
}
// 2020-09-01 00:00:00 -0400 EDT
// 2020-10-01 00:00:00 -0400 EDT
// 2020-11-01 00:00:00 -0400 EDT
// 2020-12-01 00:00:00 -0500 EST
```

To find every day at 4am:

```go
expr := te.Hour(4)
next := time.Now()
for i := 0; i < 4; i++ {
  next = expr.Next(next)
  fmt.Println(next)
}
2020-09-01 04:00:00 -0400 EDT
2020-09-02 04:00:00 -0400 EDT
2020-09-03 04:00:00 -0400 EDT
2020-09-04 04:00:00 -0400 EDT
```

Complex expressions can be composed from more primitive expressions. For
example, one can continuously generate `time.Time`s for every 3rd or 5th of
the month that falls on a Thursday or Friday except in February and December.

```go
expr := te.Intersect(
  te.Union(
    te.Day(3),
    te.Day(5),
  ),
  te.Union(
    te.Weekday(time.Thursday),
    te.Weekday(time.Friday),
  ),
  te.Except(
    te.Union(
      te.Month(time.February),
      te.Month(time.December),
    ),
  ),
)
next := time.Now()
for i := 0; i < 4; i++ {
  next = expr.Next(next)
  fmt.Println(next)
}
// 2020-09-03 00:00:00 -0400 EDT
// 2020-11-05 00:00:00 -0500 EST
// 2021-03-05 00:00:00 -0500 EST
// 2021-06-03 00:00:00 -0400 EDT
```

## Inspiration

This package is inspired by a paper on Recurring Events for Calendars
written by Martin Fowler. https://martinfowler.com/apsupp/recurring.pdf

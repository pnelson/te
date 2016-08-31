package te

import (
	"testing"
	"time"
)

func TestDay(t *testing.T) {
	tests := map[string]struct {
		day      int
		loc      *time.Location
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			day:      1,
			loc:      time.UTC,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			day:      2,
			loc:      time.UTC,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			day:      1,
			loc:      time.UTC,
			t:        time.Date(2016, 1, 1, 9, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
	}
	for name, tt := range tests {
		expr := Day(tt.day, tt.loc)
		isActive := expr.IsActive(tt.t)
		if tt.isActive != isActive {
			t.Errorf("%s\nhave isActive %v\nwant isActive %v", name, isActive, tt.isActive)
			continue
		}
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestTime(t *testing.T) {
	tests := map[string]struct {
		from     time.Time
		d        time.Duration
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			from:     time.Date(1, 1, 1, 6, 0, 0, 0, time.UTC),
			d:        time.Hour,
			t:        time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 6, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			from:     time.Date(1, 1, 1, 6, 0, 0, 0, time.UTC),
			d:        time.Hour,
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			from:     time.Date(1, 1, 1, 6, 0, 0, 0, time.UTC),
			d:        time.Hour,
			t:        time.Date(2016, 1, 1, 8, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 6, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Time(tt.from, tt.d)
		isActive := expr.IsActive(tt.t)
		if tt.isActive != isActive {
			t.Errorf("%s\nhave isActive %v\nwant isActive %v", name, isActive, tt.isActive)
			continue
		}
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestWeekday(t *testing.T) {
	tests := map[string]struct {
		weekday  time.Weekday
		count    int
		loc      *time.Location
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"monday": {
			weekday:  time.Monday,
			count:    0,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 8, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"tuesday": {
			weekday:  time.Tuesday,
			count:    0,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 2, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"wednesday local time zone": {
			weekday:  time.Wednesday,
			count:    0,
			loc:      time.Local,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.Local),
			next:     time.Date(2016, 8, 3, 0, 0, 0, 0, time.Local),
			isActive: false,
		},
		"wednesday mixed time zone": {
			weekday:  time.Wednesday,
			count:    0,
			loc:      time.Local,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 3, 0, 0, 0, 0, time.Local),
			isActive: false,
		},
		"sunday": {
			weekday:  time.Sunday,
			count:    0,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 7, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"first monday active": {
			weekday:  time.Monday,
			count:    1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 5, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"first friday upcoming this week": {
			weekday:  time.Friday,
			count:    1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 5, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"first sunday upcoming next week": {
			weekday:  time.Sunday,
			count:    1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 7, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"first sunday already passed": {
			weekday:  time.Sunday,
			count:    1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 8, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"second sunday upcoming": {
			weekday:  time.Sunday,
			count:    2,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 14, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"last sunday upcoming": {
			weekday:  time.Sunday,
			count:    -1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 28, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"last sunday active": {
			weekday:  time.Sunday,
			count:    -1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 28, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 25, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"last sunday already passed": {
			weekday:  time.Sunday,
			count:    -1,
			loc:      time.UTC,
			t:        time.Date(2016, 8, 29, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 25, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Weekday(tt.weekday, tt.count, tt.loc)
		isActive := expr.IsActive(tt.t)
		if tt.isActive != isActive {
			t.Errorf("%s\nhave isActive %v\nwant isActive %v", name, isActive, tt.isActive)
			continue
		}
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestMonth(t *testing.T) {
	tests := map[string]struct {
		month    time.Month
		loc      *time.Location
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			month:    time.February,
			loc:      time.UTC,
			t:        time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			month:    time.March,
			loc:      time.UTC,
			t:        time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 3, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			month:    time.January,
			loc:      time.UTC,
			t:        time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Month(tt.month, tt.loc)
		isActive := expr.IsActive(tt.t)
		if tt.isActive != isActive {
			t.Errorf("%s\nhave isActive %v\nwant isActive %v", name, isActive, tt.isActive)
			continue
		}
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestDateRange(t *testing.T) {
	tests := map[string]struct {
		from     time.Time
		to       time.Time
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			from:     time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			from:     time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			from:     time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
	}
	for name, tt := range tests {
		expr := DateRange(tt.from, tt.to)
		isActive := expr.IsActive(tt.t)
		if tt.isActive != isActive {
			t.Errorf("%s\nhave isActive %v\nwant isActive %v", name, isActive, tt.isActive)
			continue
		}
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestUnion(t *testing.T) {
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	expr := Union(Month(time.January, time.UTC), Day(4, time.UTC))
	isActive := expr.IsActive(now)
	if !isActive {
		t.Errorf("should be active")
	}
	next := expr.Next(now)
	want := time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("have %v\nwant %v", next, want)
	}
	next = expr.Next(want)
	want = time.Date(2016, 2, 4, 0, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("have %v\nwant %v", next, want)
	}
	outside := want.AddDate(0, 0, 1)
	isActive = expr.IsActive(outside)
	if isActive {
		t.Errorf("should not be active")
	}
}

func TestIntersection(t *testing.T) {
	now := time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC)
	expr := Intersect(Month(time.January, time.UTC), Day(4, time.UTC))
	if expr.IsActive(now) {
		t.Errorf("should not be active")
	}
	next := expr.Next(now)
	want := time.Date(2016, 1, 4, 0, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("have %v\nwant %v", next, want)
	}
	next = expr.Next(want)
	want = time.Date(2017, 1, 4, 0, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("have %v\nwant %v", next, want)
	}
	if !expr.IsActive(want) {
		t.Errorf("should be active")
	}
}

func TestExcept(t *testing.T) {
	now := time.Date(2016, 1, 3, 0, 0, 0, 0, time.UTC)
	expr := Intersect(Month(time.January, time.UTC), Except(Weekday(time.Sunday, 0, time.UTC), Day(4, time.UTC)))
	if expr.IsActive(now) {
		t.Errorf("should not be active (Sunday)")
	}
	now = now.AddDate(0, 0, 1)
	if expr.IsActive(now) {
		t.Errorf("should not be active (the 4th)")
	}
	for i := 0; i < 5; i++ {
		now = now.AddDate(0, 0, 1)
		if !expr.IsActive(now) {
			t.Errorf("should be active\ntime %v", now)
		}
	}
	now = now.AddDate(0, 0, 1)
	if expr.IsActive(now) {
		t.Errorf("should not be active (Sunday)")
	}
	for i := 0; i < 6; i++ {
		now = now.AddDate(0, 0, 1)
		if !expr.IsActive(now) {
			t.Errorf("should be active\ntime %v", now)
		}
	}
}

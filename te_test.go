package te

import (
	"testing"
	"time"
)

func TestHour(t *testing.T) {
	tests := map[string]struct {
		hour     int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			hour:     4,
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 4, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			hour:     4,
			t:        time.Date(2016, 1, 1, 3, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			hour:     4,
			t:        time.Date(2016, 1, 1, 5, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 4, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"hour negative": {
			hour:     -1,
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
		"hour greater than 23": {
			hour:     36,
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Hour(tt.hour)
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

func TestMinute(t *testing.T) {
	tests := map[string]struct {
		min      int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			min:      4,
			t:        time.Date(2016, 1, 1, 0, 4, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 1, 4, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			min:      4,
			t:        time.Date(2016, 1, 1, 0, 3, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 4, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			min:      4,
			t:        time.Date(2016, 1, 1, 0, 5, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 1, 4, 0, 0, time.UTC),
			isActive: false,
		},
		"min negative": {
			min:      -1,
			t:        time.Date(2016, 1, 1, 0, 4, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
		"min greater than 59": {
			min:      60,
			t:        time.Date(2016, 1, 1, 0, 4, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Minute(tt.min)
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

func TestSecond(t *testing.T) {
	tests := map[string]struct {
		sec      int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			sec:      4,
			t:        time.Date(2016, 1, 1, 0, 0, 4, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 1, 4, 0, time.UTC),
			isActive: true,
		},
		"before": {
			sec:      4,
			t:        time.Date(2016, 1, 1, 0, 0, 3, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 0, 4, 0, time.UTC),
			isActive: false,
		},
		"after": {
			sec:      4,
			t:        time.Date(2016, 1, 1, 0, 0, 5, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 1, 4, 0, time.UTC),
			isActive: false,
		},
		"sec negative": {
			sec:      -1,
			t:        time.Date(2016, 1, 1, 0, 0, 4, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
		"sec greater than 59": {
			sec:      60,
			t:        time.Date(2016, 1, 1, 0, 0, 4, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Second(tt.sec)
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

func TestDay(t *testing.T) {
	tests := map[string]struct {
		day      int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			day:      1,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			day:      2,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			day:      1,
			t:        time.Date(2016, 1, 1, 9, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
	}
	for name, tt := range tests {
		expr := Day(tt.day)
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
		d        time.Weekday
		n        int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"monday": {
			d:        time.Monday,
			n:        0,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 8, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"tuesday": {
			d:        time.Tuesday,
			n:        0,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 2, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"wednesday local time zone": {
			d:        time.Wednesday,
			n:        0,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.Local),
			next:     time.Date(2016, 8, 3, 0, 0, 0, 0, time.Local),
			isActive: false,
		},
		"sunday": {
			d:        time.Sunday,
			n:        0,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 7, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"first monday active": {
			d:        time.Monday,
			n:        1,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 5, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"first friday upcoming this week": {
			d:        time.Friday,
			n:        1,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 5, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"first sunday upcoming next week": {
			d:        time.Sunday,
			n:        1,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 7, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"first sunday already passed": {
			d:        time.Sunday,
			n:        1,
			t:        time.Date(2016, 8, 8, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"second sunday upcoming": {
			d:        time.Sunday,
			n:        2,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 14, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"last sunday upcoming": {
			d:        time.Sunday,
			n:        -1,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 28, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"last sunday active": {
			d:        time.Sunday,
			n:        -1,
			t:        time.Date(2016, 8, 28, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 25, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"last sunday already passed": {
			d:        time.Sunday,
			n:        -1,
			t:        time.Date(2016, 8, 29, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 9, 25, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Weekday(tt.d, tt.n)
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
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			month:    time.February,
			t:        time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			month:    time.March,
			t:        time.Date(2016, 2, 2, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 3, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			month:    time.January,
			t:        time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Month(tt.month)
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

func TestYear(t *testing.T) {
	tests := map[string]struct {
		t        time.Time
		isActive bool
	}{
		"equal": {
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			t:        time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			t:        time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	expr := Year(2016)
	for name, tt := range tests {
		isActive := expr.IsActive(tt.t)
		if tt.isActive != isActive {
			t.Errorf("%s\nhave isActive %v\nwant isActive %v", name, isActive, tt.isActive)
			continue
		}
		next := expr.Next(tt.t)
		if !next.IsZero() {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, time.Time{})
		}
	}
}

func TestDate(t *testing.T) {
	tests := map[string]struct {
		t    time.Time
		next time.Time
	}{
		"equal": {
			t:    time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			next: time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		"before": {
			t:    time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next: time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		"after": {
			t:    time.Date(2016, 2, 2, 0, 0, 0, 0, time.UTC),
			next: time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC),
		},
	}
	for name, tt := range tests {
		expr := Date(time.February, 1)
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestTime(t *testing.T) {
	tests := map[string]struct {
		t    time.Time
		next time.Time
	}{
		"equal": {
			t:    time.Date(2016, 2, 1, 15, 04, 05, 0, time.UTC),
			next: time.Date(2016, 2, 2, 15, 04, 05, 0, time.UTC),
		},
		"before": {
			t:    time.Date(2016, 2, 1, 15, 01, 00, 0, time.UTC),
			next: time.Date(2016, 2, 2, 15, 04, 05, 0, time.UTC),
		},
		"after": {
			t:    time.Date(2016, 2, 1, 15, 04, 06, 0, time.UTC),
			next: time.Date(2016, 2, 2, 15, 04, 05, 0, time.UTC),
		},
	}
	for name, tt := range tests {
		expr := Time(15, 04, 05)
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestDateRange(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			t1:       time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			t1:       time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			t1:       time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			t2:       time.Date(2016, 9, 4, 0, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
	}
	for name, tt := range tests {
		expr := DateRange(tt.t1, tt.t2)
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

func TestTimeRange(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			t1:       time.Date(1, 1, 1, 6, 0, 0, 0, time.UTC),
			t2:       time.Date(1, 1, 1, 7, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 6, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			t1:       time.Date(1, 1, 1, 6, 0, 0, 0, time.UTC),
			t2:       time.Date(1, 1, 1, 7, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			t1:       time.Date(1, 1, 1, 6, 0, 0, 0, time.UTC),
			t2:       time.Date(1, 1, 1, 7, 0, 0, 0, time.UTC),
			t:        time.Date(2016, 1, 1, 8, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 6, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := TimeRange(tt.t1, tt.t2)
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
	expr := Union(Month(time.January), Day(4))
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

func TestIntersect(t *testing.T) {
	now := time.Date(2015, 8, 1, 0, 0, 0, 0, time.UTC)
	expr := Intersect(Month(time.January), Day(4), Hour(9))
	if expr.IsActive(now) {
		t.Errorf("should not be active")
	}
	next := expr.Next(now)
	want := time.Date(2016, 1, 4, 9, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("next 1\nhave %v\nwant %v", next, want)
	}
	next = expr.Next(want)
	want = time.Date(2017, 1, 4, 9, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Errorf("next 2\nhave %v\nwant %v", next, want)
	}
	if !expr.IsActive(want) {
		t.Errorf("should be active")
	}
}

func TestExcept(t *testing.T) {
	now := time.Date(2016, 1, 3, 0, 0, 0, 0, time.UTC)
	expr := Intersect(Month(time.January), Except(Weekday(time.Sunday, 0), Day(4)))
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

func TestUntil(t *testing.T) {
	tests := []struct {
		expr Expression
		d    time.Duration
	}{
		{Day(4), 3 * 24 * time.Hour},
		{Weekday(time.Thursday, 0), 3 * 24 * time.Hour},
		{Month(time.September), 31 * 24 * time.Hour},
	}
	now := time.Date(2016, time.August, 1, 0, 0, 0, 0, time.UTC)
	for i, tt := range tests {
		d := Until(tt.expr, now)
		if d != tt.d {
			t.Errorf("Until %d.\nhave %v\nwant %v", i, d, tt.d)
		}
	}
}

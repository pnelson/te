package te

import (
	"fmt"
	"reflect"
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

func TestHourly(t *testing.T) {
	tests := map[string]struct {
		n        int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"zero": {
			n:        5,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 5, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"equal": {
			n:        5,
			t:        time.Date(2016, 1, 1, 5, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			n:        5,
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 5, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			n:        5,
			t:        time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"wrap": {
			n:        5,
			t:        time.Date(2016, 1, 1, 22, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"less than 1": {
			n:        0,
			t:        time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
		"greater than 12": {
			n:        13,
			t:        time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Hourly(tt.n)
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

func TestMinutely(t *testing.T) {
	tests := map[string]struct {
		n        int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"zero": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 15, 0, 0, time.UTC),
			isActive: true,
		},
		"equal": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 15, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 30, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 14, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 15, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 16, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 30, 0, 0, time.UTC),
			isActive: false,
		},
		"wrap": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 58, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 1, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"less than 1": {
			n:        0,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
		"greater than 30": {
			n:        31,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Minutely(tt.n)
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

func TestSecondly(t *testing.T) {
	tests := map[string]struct {
		n        int
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"zero": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 0, 15, 0, time.UTC),
			isActive: true,
		},
		"equal": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 0, 15, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 0, 30, 0, time.UTC),
			isActive: true,
		},
		"before": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 0, 14, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 0, 15, 0, time.UTC),
			isActive: false,
		},
		"after": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 0, 16, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 0, 30, 0, time.UTC),
			isActive: false,
		},
		"wrap": {
			n:        15,
			t:        time.Date(2016, 1, 1, 0, 0, 58, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 0, 1, 0, 0, time.UTC),
			isActive: false,
		},
		"less than 1": {
			n:        0,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
		"greater than 30": {
			n:        31,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Secondly(tt.n)
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
		"zero": {
			day:      0,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
		},
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
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"31st of the month": {
			day:      31,
			t:        time.Date(2016, 1, 31, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 3, 31, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"last day of the month": {
			day:      -1,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 31, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"nth last day of the month": {
			day:      -2,
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Time{},
			isActive: false,
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
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"Monday": {
			d:        time.Monday,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 8, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"Tuesday": {
			d:        time.Tuesday,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 2, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"Sunday": {
			d:        time.Sunday,
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 7, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := Weekday(tt.d)
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
	expr := Date(time.February, 1)
	for name, tt := range tests {
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
			next: time.Date(2016, 2, 1, 15, 04, 05, 0, time.UTC),
		},
		"after": {
			t:    time.Date(2016, 2, 1, 15, 04, 06, 0, time.UTC),
			next: time.Date(2016, 2, 2, 15, 04, 05, 0, time.UTC),
		},
	}
	expr := Time(15, 04, 05)
	for name, tt := range tests {
		next := expr.Next(tt.t)
		if !next.Equal(tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestDateRange(t *testing.T) {
	tests := map[string]struct {
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			t:        time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			t:        time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			t:        time.Date(2016, 8, 3, 0, 0, 0, 0, time.UTC),
			next:     time.Date(2017, 8, 1, 0, 0, 0, 0, time.UTC),
			isActive: true,
		},
	}
	for name, tt := range tests {
		expr := DateRange(time.August, 1, time.September, 4)
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
		t        time.Time
		next     time.Time
		isActive bool
	}{
		"equal": {
			t:        time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 6, 0, 0, 0, time.UTC),
			isActive: true,
		},
		"before": {
			t:        time.Date(2016, 1, 1, 4, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 1, 6, 0, 0, 0, time.UTC),
			isActive: false,
		},
		"after": {
			t:        time.Date(2016, 1, 1, 8, 0, 0, 0, time.UTC),
			next:     time.Date(2016, 1, 2, 6, 0, 0, 0, time.UTC),
			isActive: false,
		},
	}
	for name, tt := range tests {
		expr := TimeRange(6, 0, 0, 7, 0, 0)
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
	tests := map[string]struct {
		t    time.Time
		expr Expression
		next []time.Time
	}{
		"January 4th at 1am": {
			t:    time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			expr: Intersect(Month(time.January), Day(4), Hour(1)),
			next: []time.Time{
				time.Date(2016, 1, 4, 1, 0, 0, 0, time.UTC),
				time.Date(2017, 1, 4, 1, 0, 0, 0, time.UTC),
				time.Date(2018, 1, 4, 1, 0, 0, 0, time.UTC),
			},
		},
		"January 4th at 1am (nested)": {
			t:    time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			expr: Intersect(Date(time.January, 4), Hour(1)),
			next: []time.Time{
				time.Date(2016, 1, 4, 1, 0, 0, 0, time.UTC),
				time.Date(2017, 1, 4, 1, 0, 0, 0, time.UTC),
				time.Date(2018, 1, 4, 1, 0, 0, 0, time.UTC),
			},
		},
		"February 29th": {
			t:    time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			expr: Date(time.February, 29),
			next: []time.Time{
				time.Date(2016, 2, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
				time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			},
		},
		"1st and 5th of every month on Thursday and Friday except August and December in 2016": {
			t: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			expr: Intersect(
				Year(2016),
				Union(
					Day(1),
					Day(5),
				),
				Union(
					Weekday(time.Thursday),
					Weekday(time.Friday),
				),
				Except(
					Union(
						Month(time.August),
						Month(time.December),
					),
				),
			),
			next: []time.Time{
				time.Date(2016, 2, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2016, 4, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2016, 5, 5, 0, 0, 0, 0, time.UTC),
				time.Date(2016, 7, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
				time.Time{},
			},
		},
		"every Friday and Saturday between Jan 1 and Jan 14": {
			t: time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC),
			expr: Intersect(
				Union(
					Weekday(time.Friday),
					Weekday(time.Saturday),
				),
				DateRange(time.January, 1, time.January, 14),
			),
			next: []time.Time{
				time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
				time.Date(2016, 1, 8, 0, 0, 0, 0, time.UTC),
				time.Date(2016, 1, 9, 0, 0, 0, 0, time.UTC),
				time.Date(2017, 1, 6, 0, 0, 0, 0, time.UTC),
				time.Date(2017, 1, 7, 0, 0, 0, 0, time.UTC),
			},
		},
	}
	for name, tt := range tests {
		next := next(tt.expr, tt.t, len(tt.next))
		if !reflect.DeepEqual(next, tt.next) {
			t.Errorf("%s\nhave next %v\nwant next %v", name, next, tt.next)
		}
	}
}

func TestExcept(t *testing.T) {
	now := time.Date(2016, 1, 3, 0, 0, 0, 0, time.UTC)
	expr := Intersect(Month(time.January), Except(Weekday(time.Sunday), Day(4)))
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
		{Weekday(time.Thursday), 3 * 24 * time.Hour},
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

func TestString(t *testing.T) {
	tests := map[string]struct {
		expr Expression
		want string
	}{
		"hour": {
			Hour(15),
			"te.Hour(15)",
		},
		"hourly": {
			Hourly(2),
			"te.Hourly(2)",
		},
		"minute": {
			Minute(4),
			"te.Minute(4)",
		},
		"minutely": {
			Minutely(2),
			"te.Minutely(2)",
		},
		"second": {
			Second(5),
			"te.Second(5)",
		},
		"secondly": {
			Secondly(2),
			"te.Secondly(2)",
		},
		"day": {
			Day(2),
			"te.Day(2)",
		},
		"weekday": {
			Weekday(time.Monday),
			"te.Weekday(time.Monday)",
		},
		"month": {
			Month(time.January),
			"te.Month(time.January)",
		},
		"year": {
			Year(2016),
			"te.Year(2016)",
		},
		"date range": {
			DateRange(time.January, 2, time.February, 14),
			"te.DateRange(time.January, 2, time.February, 14)",
		},
		"time range": {
			TimeRange(6, 0, 0, 7, 0, 0),
			"te.TimeRange(6, 0, 0, 7, 0, 0)",
		},
		"union": {
			Union(Weekday(time.Tuesday), Weekday(time.Thursday)),
			"te.Union(te.Weekday(time.Tuesday), te.Weekday(time.Thursday))",
		},
		"intersect": {
			Intersect(Month(time.January), Weekday(time.Tuesday)),
			"te.Intersect(te.Month(time.January), te.Weekday(time.Tuesday))",
		},
		"except": {
			Except(Weekday(time.Sunday)),
			"te.Except(te.Weekday(time.Sunday))",
		},
	}
	for name, tt := range tests {
		e, ok := tt.expr.(fmt.Stringer)
		if !ok {
			t.Errorf("%s should be a fmt.Stringer", name)
			continue
		}
		have := e.String()
		if have != tt.want {
			t.Errorf("%s\nhave %s\nwant %s", name, have, tt.want)
		}
	}
}

// next returns the next N times for the given expression
// including zero times.
func next(expr Expression, t time.Time, n int) []time.Time {
	ts := make([]time.Time, n)
	for i := 0; i < n; i++ {
		t = expr.Next(t)
		ts[i] = t
	}
	return ts
}

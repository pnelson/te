package te

import (
	"reflect"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		want Expression
	}{
		{"daily", Hour(0)},
		{"midnight", Hour(0)},
		{"noon", Hour(12)},

		{"3pm", Hour(15)},
		{"3 PM", Hour(15)},
		{"15:00", Hour(15)},
		{"15:04", Intersect(Hour(15), Minute(4))},
		{"15:04:05", Intersect(Hour(15), Minute(4), Second(5))},

		{"at midnight", Hour(0)},
		{"at noon", Hour(12)},
		{"at 3pm", Hour(15)},

		{"daily at midnight", Hour(0)},
		{"daily at noon", Hour(12)},
		{"daily at 3pm", Hour(15)},

		{"every day", Hour(0)},
		{"every day at 3pm", Hour(15)},

		{"weekly", Weekday(time.Sunday)},
		{"weekly on Friday", Weekday(time.Friday)},
		{"every week", Weekday(time.Sunday)},
		{"every week on Friday", Weekday(time.Friday)},
		{"every Friday", Weekday(time.Friday)},
		{"Friday", Weekday(time.Friday)},

		{"monthly", Day(1)},
		{"every month", Day(1)},
		{"every 1st", Day(1)},
		{"every 4th", Day(4)},
		{"4th", Day(4)},

		{"yearly", Intersect(Month(time.January), Day(1))},
		{"annually", Intersect(Month(time.January), Day(1))},
		{"annually at 4am", Intersect(Intersect(Month(time.January), Day(1)), Hour(4))},
		{"every year", Intersect(Month(time.January), Day(1))},
		{"every November", Month(time.November)},
		{"November", Month(time.November)},

		{"hourly", Minute(0)},
		{"every hour", Minute(0)},

		{"every minute", Second(0)},
		{"every second", Secondly(1)},

		{"every 2 hours", Hourly(2)},
		{"every 15 minutes", Minutely(15)},
		{"every 30 seconds", Secondly(30)},

		{"3pm and 9pm", Union(Hour(15), Hour(21))},
		{"15:00 and 21:00", Union(Hour(15), Hour(21))},
		{"15:04:05 and 9pm", Union(Intersect(Hour(15), Minute(4), Second(5)), Hour(21))},

		{"at noon and 9pm", Union(Hour(12), Hour(21))},
		{"at noon and 6pm and 9pm", Union(Hour(12), Hour(18), Hour(21))},
		{"at noon, 6pm, and 9pm", Union(Hour(12), Hour(18), Hour(21))},
		{"at noon/6pm/9pm", Union(Hour(12), Hour(18), Hour(21))},

		{"Tuesday and Thursday", Union(Weekday(time.Tuesday), Weekday(time.Thursday))},
		{"Tue and Wed and Thu", Union(Weekday(time.Tuesday), Weekday(time.Wednesday), Weekday(time.Thursday))},
		{"Tue, Wed, and Thu", Union(Weekday(time.Tuesday), Weekday(time.Wednesday), Weekday(time.Thursday))},
		{"Tue/Wed/Thu", Union(Weekday(time.Tuesday), Weekday(time.Wednesday), Weekday(time.Thursday))},

		{"every 1st and 3rd", Union(Day(1), Day(3))},
		{"every 1st and 3rd and 9th", Union(Day(1), Day(3), Day(9))},
		{"every 1st, 3rd, and 9th", Union(Day(1), Day(3), Day(9))},
		{"every 1st/3rd/9th", Union(Day(1), Day(3), Day(9))},

		{"July and August", Union(Month(time.July), Month(time.August))},
		{"May and June and July", Union(Month(time.May), Month(time.June), Month(time.July))},
		{"May, June, and July", Union(Month(time.May), Month(time.June), Month(time.July))},
		{"May/Jun/Jul", Union(Month(time.May), Month(time.June), Month(time.July))},

		{"April 19th", Intersect(Month(time.April), Day(19))},
		{"every April 19th", Intersect(Month(time.April), Day(19))},
		{"every April 19th at 3pm", Intersect(Month(time.April), Day(19), Hour(15))},
		{"every 2 hours on Sunday", Intersect(Hourly(2), Weekday(time.Sunday))},
		{"Tue/Thu at 4am", Intersect(Union(Weekday(time.Tuesday), Weekday(time.Thursday)), Hour(4))},
	}
	for _, tt := range tests {
		have, err := Parse(tt.in, time.UTC)
		if err != nil {
			t.Fatalf("Parse(%q) %v", tt.in, err)
		} else if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("Parse(%q)\nhave %v\nwant %v", tt.in, have, tt.want)
		}
	}
}

func TestParseError(t *testing.T) {
	var tests = []string{
		"",
		"   ",
		"daily at",
		"daily at daily",
		"midnight at",
		"midnight at daily",
		"noon at",
		"noon at daily",
		"3:",
		"3:4",
		"15:04:",
		"every",
		"in noon",
		"at noon and",
	}
	for _, tt := range tests {
		have, err := Parse(tt, time.UTC)
		if err == nil {
			t.Errorf("Parse(%q)\nhave %v\nwant parse error", tt, have)
		}
	}
}

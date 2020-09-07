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

		{"yearly", Month(time.January)},
		{"annually", Month(time.January)},
		{"every year", Month(time.January)},
		{"every November", Month(time.November)},
		{"November", Month(time.November)},

		{"hourly", Minute(0)},
		{"every hour", Minute(0)},

		{"every 2 hours", Hourly(2)},
		{"every 15 minutes", Minutely(15)},
		{"every 30 seconds", Secondly(30)},
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
	}
	for _, tt := range tests {
		have, err := Parse(tt, time.UTC)
		if err == nil {
			t.Errorf("Parse(%q)\nhave %v\nwant parse error", tt, have)
		}
	}
}

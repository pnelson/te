package te

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		in   string
		want []token
	}{
		{
			"",
			[]token{},
		},
		{
			"   ",
			[]token{},
		},
		{
			"daily",
			[]token{
				{tokenDaily, "daily"},
			},
		},
		{
			"midnight",
			[]token{
				{tokenDaily, "midnight"},
			},
		},
		{
			"noon",
			[]token{
				{tokenNoon, "noon"},
			},
		},
		{
			"hourly",
			[]token{
				{tokenHourly, "hourly"},
			},
		},
		{
			"weekly",
			[]token{
				{tokenWeekly, "weekly"},
			},
		},
		{
			"monthly",
			[]token{
				{tokenMonthly, "monthly"},
			},
		},
		{
			"yearly",
			[]token{
				{tokenYearly, "yearly"},
			},
		},
		{
			"annually",
			[]token{
				{tokenYearly, "annually"},
			},
		},
		{
			"Sunday",
			[]token{
				{tokenWeekday, "sunday"},
			},
		},
		{
			"Wednesday",
			[]token{
				{tokenWeekday, "wednesday"},
			},
		},
		{
			"January",
			[]token{
				{tokenMonth, "january"},
			},
		},
		{
			"November",
			[]token{
				{tokenMonth, "november"},
			},
		},
		{
			"3",
			[]token{
				{tokenDigit, "3"},
			},
		},
		{
			"3am",
			[]token{
				{tokenDigit, "3"},
				{tokenTwelveHour, "am"},
			},
		},
		{
			"3pm",
			[]token{
				{tokenDigit, "3"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"3 PM",
			[]token{
				{tokenDigit, "3"},
				{tokenTwelveHour, "pm"},
			},
		},
		{
			"1st",
			[]token{
				{tokenDigit, "1"},
				{tokenOrdinal, "st"},
			},
		},
		{
			"2nd",
			[]token{
				{tokenDigit, "2"},
				{tokenOrdinal, "nd"},
			},
		},
		{
			"3rd",
			[]token{
				{tokenDigit, "3"},
				{tokenOrdinal, "rd"},
			},
		},
		{
			"4th",
			[]token{
				{tokenDigit, "4"},
				{tokenOrdinal, "th"},
			},
		},
		{
			"every last day of the month at",
			[]token{
				{tokenEvery, "every"},
				{tokenLast, "last"},
				{tokenUnitDay, "day"},
				{tokenOf, "of"},
				{tokenThe, "the"},
				{tokenUnitMonth, "month"},
				{tokenAt, "at"},
			},
		},
	}
	for _, tt := range tests {
		have, err := lex(tt.in)
		if err != nil {
			t.Fatalf("lex(%q) %v", tt.in, err)
		}
		if !reflect.DeepEqual(have, tt.want) {
			t.Errorf("lex(%q)\nhave %v\nwant %v", tt.in, have, tt.want)
		}
	}
}

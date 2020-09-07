package te

import "fmt"

type tokenType int

const (
	tokenAnd tokenType = iota
	tokenAt
	tokenColon
	tokenDaily
	tokenDigit
	tokenError
	tokenEvery
	tokenEOF
	tokenHourly
	tokenIn
	tokenLast
	tokenNoon
	tokenOf
	tokenOn
	tokenOrdinal
	tokenThe
	tokenTwelveHour
	tokenWeekday
	tokenWeekly
	tokenMonth
	tokenMonthly
	tokenUnitDay
	tokenUnitHour
	tokenUnitMinute
	tokenUnitMonth
	tokenUnitSecond
	tokenUnitWeek
	tokenUnitYear
	tokenYearly
)

const eof = rune(-1)

type token struct {
	typ tokenType
	val string
}

func (t token) String() string {
	switch t.typ {
	case tokenError:
		return t.val
	case tokenEOF:
		return "EOF"
	}
	return fmt.Sprintf("%q", t.val)
}

type parseError struct {
	token   token
	message string
}

func (e parseError) Error() string {
	return fmt.Sprintf("%s, token: %q", e.message, e.token.val)
}

func newParseError(t token, message string) parseError {
	return parseError{t, message}
}

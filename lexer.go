package te

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

type lexer struct {
	input  string
	i, j   int // position within input
	width  int // width of last rune
	state  stateFn
	tokens []token
}

func lex(s string) ([]token, error) {
	l := &lexer{
		input:  strings.ToLower(s),
		state:  readExpr,
		tokens: make([]token, 0),
	}
	for state := readExpr; state != nil; {
		state = state(l)
	}
	if len(l.tokens) > 0 {
		last := l.tokens[len(l.tokens)-1]
		if last.typ == tokenError {
			return nil, errors.New(last.val)
		}
	}
	return l.tokens, nil
}

func (l *lexer) emit(typ tokenType) {
	l.tokens = append(l.tokens, token{typ, l.value()})
	l.i = l.j
}

func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens = append(l.tokens, token{tokenError, fmt.Sprintf(format, args...)})
	return nil
}

func (l *lexer) ignore() {
	l.i = l.j
}

func (l *lexer) read() rune {
	if l.j >= len(l.input) {
		return eof
	}
	r, width := utf8.DecodeRuneInString(l.input[l.j:])
	l.j += width
	l.width = width
	return r
}

func (l *lexer) readFn(fn func(rune) bool) {
	for {
		r := l.read()
		if r == eof {
			break
		} else if !fn(r) {
			l.unread()
			break
		}
	}
}

func (l *lexer) readRune(ch rune) bool {
	r := l.read()
	return r == ch
}

func (l *lexer) peek() rune {
	r := l.read()
	l.unread()
	return r
}

func (l *lexer) unread() {
	l.j -= l.width
}

func (l *lexer) value() string {
	return l.input[l.i:l.j]
}

func readExpr(l *lexer) stateFn {
	l.readFn(unicode.IsSpace)
	l.ignore()
	r := l.peek()
	switch {
	case r == eof:
		return nil
	case unicode.IsLetter(r):
		return readLetter
	case unicode.IsDigit(r):
		return readDigit
	}
	return l.errorf("invalid character")
}

func readColon(l *lexer) stateFn {
	_ = l.read()
	l.emit(tokenColon)
	r := l.peek()
	if !unicode.IsDigit(r) {
		return l.errorf("colon must be followed by a digit")
	}
	return readDigit
}

func readDigit(l *lexer) stateFn {
	l.readFn(unicode.IsDigit)
	l.emit(tokenDigit)
	r := l.peek()
	switch r {
	case eof:
		return nil
	case 's', 'n', 'r', 't':
		return readOrdinal
	case 'a', 'p':
		return readTwelveHour
	case ':':
		return readColon
	}
	return readSpace
}

func readLetter(l *lexer) stateFn {
	l.readFn(unicode.IsLetter)
	val := l.value()
	switch val {
	case "daily", "midnight":
		l.emit(tokenDaily)
	case "noon":
		l.emit(tokenNoon)
	case "hourly":
		l.emit(tokenHourly)
	case "weekly":
		l.emit(tokenWeekly)
	case "monthly":
		l.emit(tokenMonthly)
	case "yearly", "annually":
		l.emit(tokenYearly)
	case "sun", "sunday":
		fallthrough
	case "mon", "monday":
		fallthrough
	case "tue", "tuesday":
		fallthrough
	case "wed", "wednesday":
		fallthrough
	case "thu", "thursday":
		fallthrough
	case "fri", "friday":
		fallthrough
	case "sat", "saturday":
		l.emit(tokenWeekday)
	case "jan", "january":
		fallthrough
	case "feb", "february":
		fallthrough
	case "mar", "march":
		fallthrough
	case "apr", "april":
		fallthrough
	case "may":
		fallthrough
	case "jun", "june":
		fallthrough
	case "jul", "july":
		fallthrough
	case "aug", "august":
		fallthrough
	case "sep", "september":
		fallthrough
	case "oct", "october":
		fallthrough
	case "nov", "november":
		fallthrough
	case "dec", "december":
		l.emit(tokenMonth)
	case "every":
		l.emit(tokenEvery)
	case "second", "seconds":
		l.emit(tokenUnitSecond)
	case "minute", "minutes":
		l.emit(tokenUnitMinute)
	case "hour", "hours":
		l.emit(tokenUnitHour)
	case "day", "days":
		l.emit(tokenUnitDay)
	case "week", "weeks":
		l.emit(tokenUnitWeek)
	case "month", "months":
		l.emit(tokenUnitMonth)
	case "year", "years":
		l.emit(tokenUnitYear)
	case "am", "pm":
		l.emit(tokenTwelveHour)
	case "at":
		l.emit(tokenAt)
	case "in":
		l.emit(tokenIn)
	case "of":
		l.emit(tokenOf)
	case "on":
		l.emit(tokenOn)
	case "and":
		l.emit(tokenAnd)
	case "the":
		l.emit(tokenThe)
	case "last":
		l.emit(tokenLast)
	default:
		return nil
	}
	return readSpace
}

func readOrdinal(l *lexer) stateFn {
	var ok bool
	r := l.read()
	switch r {
	case 's':
		ok = l.readRune('t') // 1st
	case 'n', 'r':
		ok = l.readRune('d') // 2nd, 3rd
	case 't':
		ok = l.readRune('h') // 4th-9th
	}
	if !ok {
		return l.errorf("invalid ordinal")
	}
	l.emit(tokenOrdinal)
	return readSpace
}

func readSpace(l *lexer) stateFn {
	r := l.peek()
	if r == eof {
		return nil
	}
	if !unicode.IsSpace(r) {
		return l.errorf("expected space")
	}
	l.readFn(unicode.IsSpace)
	l.ignore()
	return readExpr
}

func readTwelveHour(l *lexer) stateFn {
	_ = l.read()
	r := l.read()
	if r != 'm' {
		return l.errorf("expected twelve hour am/pm marker")
	}
	l.emit(tokenTwelveHour)
	return readSpace
}

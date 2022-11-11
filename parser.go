package te

import (
	"strconv"
	"strings"
	"time"
)

type parser struct {
	loc    *time.Location
	pos    int
	tokens []token
	exprs  []Expression
	join   bool
}

// Parse parses the provided string into an Expression.
func Parse(s string, loc *time.Location) (Expression, error) {
	s = strings.TrimSpace(s)
	tokens, err := lex(s)
	if err != nil {
		return nil, err
	}
	p := &parser{
		loc:    loc,
		tokens: tokens,
		exprs:  make([]Expression, 0),
	}
	err = p.parseExpr()
	if err != nil {
		return nilExpr{}, err
	}
	if len(p.exprs) == 1 {
		return p.exprs[0], nil
	}
	return Intersect(p.exprs...), nil
}

func (p *parser) add(expr Expression) error {
	if p.join {
		exprs := make([]Expression, 0)
		u, ok := p.exprs[len(p.exprs)-1].(unionExpr)
		if ok {
			for _, e := range u {
				exprs = append(exprs, e)
			}
			exprs = append(exprs, expr)
		} else {
			exprs = append(exprs, p.exprs[len(p.exprs)-1], expr)
		}
		p.exprs[len(p.exprs)-1] = Union(exprs...)
		p.join = false
	} else {
		p.exprs = append(p.exprs, expr)
	}
	return p.parseExpr()
}

func (p *parser) parseAt() error {
	t := p.next()
	switch t.typ {
	case tokenMidnight:
		return p.parseMidnight(t)
	case tokenDigit:
		return p.parseTime(t)
	case tokenNoon:
		return p.parseNoon()
	}
	return newParseError(t, "expected time or time constant")
}

func (p *parser) parseDaily() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.exprs = append(p.exprs, Daily())
		return nil
	case tokenAt:
		return p.parseAt()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseDigit(d token) error {
	t := p.next()
	switch t.typ {
	case tokenColon:
		m := p.next()
		return p.parseTwentyFourHour(d, m)
	case tokenOrdinal:
		return p.parseOrdinal(d)
	case tokenTwelveHour:
		return p.parseTwelveHour(d, t)
	case tokenUnitHour:
		return p.parseUnitHour(d)
	case tokenUnitMinute:
		return p.parseUnitMinute(d)
	case tokenUnitSecond:
		return p.parseUnitSecond(d)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseEvery() error {
	t := p.next()
	switch t.typ {
	case tokenDigit:
		return p.parseDigit(t)
	case tokenMonth:
		return p.parseMonth(t)
	case tokenUnitSecond:
		return p.parseSecondly()
	case tokenUnitMinute:
		return p.parseMinutely()
	case tokenUnitHour:
		return p.parseHourly()
	case tokenUnitDay:
		return p.parseDaily()
	case tokenUnitWeek:
		return p.parseWeekly()
	case tokenUnitMonth:
		return p.parseMonthly()
	case tokenUnitYear:
		return p.parseYearly()
	case tokenWeekday:
		return p.parseWeekday(t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseExpr() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		switch {
		case p.join:
			return newParseError(t, "incomplete expression")
		case len(p.exprs) == 0:
			return newParseError(t, "empty expression")
		default:
			return nil
		}
	case tokenAnd:
		p.join = true
		return p.parseExpr()
	case tokenAt:
		return p.parseAt()
	case tokenDaily:
		return p.parseDaily()
	case tokenDigit:
		return p.parseDigit(t)
	case tokenEvery:
		return p.parseEvery()
	case tokenHourly:
		return p.parseHourly()
	case tokenMidnight:
		return p.parseMidnight(t)
	case tokenMonth:
		return p.parseMonth(t)
	case tokenMonthly:
		return p.parseMonthly()
	case tokenNoon:
		return p.parseNoon()
	case tokenOn:
		return p.parseOn()
	case tokenWeekly:
		return p.parseWeekly()
	case tokenWeekday:
		return p.parseWeekday(t)
	case tokenYearly:
		return p.parseYearly()
	}
	return newParseError(t, "unexpected expression")
}

func (p *parser) parseHourly() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.exprs = append(p.exprs, Minute(0))
		return nil
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseMidnight(t token) error {
	if t.val != "midnight" {
		return newParseError(t, "expected at midnight")
	}
	expr := Hour(0)
	return p.add(expr)
}

func (p *parser) parseMinutely() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.exprs = append(p.exprs, Second(0))
		return nil
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseMonth(t token) error {
	var m time.Month
	switch t.val[:3] {
	case "jan":
		m = time.January
	case "feb":
		m = time.February
	case "mar":
		m = time.March
	case "apr":
		m = time.April
	case "may":
		m = time.May
	case "jun":
		m = time.June
	case "jul":
		m = time.July
	case "aug":
		m = time.August
	case "sep":
		m = time.September
	case "oct":
		m = time.October
	case "nov":
		m = time.November
	case "dec":
		m = time.December
	default:
		return newParseError(t, "invalid month")
	}
	expr := Month(m)
	return p.add(expr)
}

func (p *parser) parseMonthly() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.exprs = append(p.exprs, Day(1))
		return nil
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseNoon() error {
	expr := Hour(12)
	return p.add(expr)
}

func (p *parser) parseOn() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		return newParseError(t, "expected weekday")
	case tokenWeekday:
		return p.parseWeekday(t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseOrdinal(d token) error {
	n, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	expr := Day(n)
	return p.add(expr)
}

func (p *parser) parseSecondly() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.exprs = append(p.exprs, Secondly(1))
		return nil
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseTime(h token) error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		return newParseError(t, "expected time")
	case tokenColon:
		m := p.next()
		return p.parseTwentyFourHour(h, m)
	case tokenTwelveHour:
		return p.parseTwelveHour(h, t)
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseTwentyFourHour(h, m token) error {
	c := p.peek()
	if c.typ == tokenColon {
		p.next()
		return p.parseTwentyFourHourWithSeconds(h, m)
	}
	t, err := time.ParseInLocation("15:04", h.val+":"+m.val, p.loc)
	if err != nil {
		return err
	}
	hour, min, _ := t.Clock()
	expr := Hour(hour)
	if min != 0 {
		expr = Intersect(expr, Minute(min))
	}
	return p.add(expr)
}

func (p *parser) parseTwentyFourHourWithSeconds(h, m token) error {
	s := p.next()
	t, err := time.ParseInLocation("15:04:05", h.val+":"+m.val+":"+s.val, p.loc)
	if err != nil {
		return err
	}
	hour, min, sec := t.Clock()
	exprs := []Expression{Hour(hour)}
	if min != 0 || sec != 0 {
		exprs = append(exprs, Minute(min))
	}
	if sec != 0 {
		exprs = append(exprs, Second(sec))
	}
	expr := exprs[0]
	if len(exprs) > 1 {
		expr = Intersect(exprs...)
	}
	return p.add(expr)
}

func (p *parser) parseTwelveHour(h, ampm token) error {
	t, err := time.ParseInLocation("3pm", h.val+ampm.val, p.loc)
	if err != nil {
		return err
	}
	hour := t.Hour()
	expr := Hour(hour)
	return p.add(expr)
}

func (p *parser) parseUnitHour(d token) error {
	hour, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	expr := Hourly(hour)
	return p.add(expr)
}

func (p *parser) parseUnitMinute(d token) error {
	min, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	expr := Minutely(min)
	return p.add(expr)
}

func (p *parser) parseUnitSecond(d token) error {
	sec, err := strconv.Atoi(d.val)
	if err != nil {
		return err
	}
	expr := Secondly(sec)
	return p.add(expr)
}

func (p *parser) parseWeekday(t token) error {
	var d time.Weekday
	switch t.val[:3] {
	case "sun":
		d = time.Sunday
	case "mon":
		d = time.Monday
	case "tue":
		d = time.Tuesday
	case "wed":
		d = time.Wednesday
	case "thu":
		d = time.Thursday
	case "fri":
		d = time.Friday
	case "sat":
		d = time.Saturday
	default:
		return newParseError(t, "invalid weekday")
	}
	expr := Weekday(d)
	return p.add(expr)
}

func (p *parser) parseWeekly() error {
	t := p.next()
	switch t.typ {
	case tokenEOF:
		p.exprs = append(p.exprs, Weekday(time.Sunday))
		return nil
	case tokenOn:
		return p.parseOn()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) parseYearly() error {
	expr := Intersect(Month(time.January), Day(1))
	err := p.add(expr)
	if err != nil {
		return err
	}
	t := p.next()
	switch t.typ {
	case tokenEOF:
		return nil
	case tokenAt:
		return p.parseAt()
	}
	return newParseError(t, "unexpected token")
}

func (p *parser) peek() token {
	if p.pos >= len(p.tokens) {
		return token{tokenEOF, ""}
	}
	return p.tokens[p.pos]
}

func (p *parser) next() token {
	t := p.peek()
	p.pos++
	return t
}

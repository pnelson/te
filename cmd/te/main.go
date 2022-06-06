package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pnelson/te"
)

var (
	help = flag.Bool("h", false, "show this usage information")

	u = flag.Bool("u", false, "output as UTC")
	s = flag.Bool("s", false, "output as Unix time in seconds")
	f = flag.String("f", "Mon Jan 2 15:04 MST", "time format layout")
	l = flag.String("l", "Local", "expression timezone location")
	n = flag.Int("n", 1, "number of time generations")

	rfc3339 = flag.Bool("rfc-3339", false, "output as RFC 3339 format")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] EXPR\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	loc, err := time.LoadLocation(*l)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	args := flag.Args()
	expr := strings.Join(args, " ")
	e, err := te.Parse(expr, loc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
		return
	}
	if *rfc3339 {
		*f = time.RFC3339
	}
	next := time.Now()
	for i := 0; i < *n; i++ {
		next = e.Next(next)
		if *s {
			fmt.Println(next.Unix())
			continue
		}
		if *u {
			next = next.In(time.UTC)
		}
		fmt.Println(next.Format(*f))
	}
}

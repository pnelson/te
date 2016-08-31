// Example usage of the Schedule capabilities of package te.
// This example sends an Event every Mon/Wed/Fri/Sat at 06:00 local time.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/pnelson/te"
)

func main() {
	hours := te.Time(time.Date(1, 1, 1, 6, 0, 0, 0, time.Local), time.Hour)
	weekdays := te.Union(
		te.Weekday(time.Monday, 0, time.Local),
		te.Weekday(time.Wednesday, 0, time.Local),
		te.Weekday(time.Friday, 0, time.Local),
		te.Weekday(time.Saturday, 0, time.Local),
	)
	s := te.NewSchedule()
	s.Set("fitness", te.Intersect(weekdays, hours))
	quit := make(chan os.Signal, 1)
	done := make(chan struct{})
	signal.Notify(quit, os.Interrupt)
	go func() {
		for {
			select {
			case e := <-s.Events:
				fmt.Printf("%v (%s)\n", e.Time, e.Name)
			case <-quit:
				done <- struct{}{}
				return
			}
		}
	}()
	<-done
	s.Close()
}

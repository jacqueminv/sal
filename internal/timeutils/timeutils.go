package timeutils

import (
	"time"
)

const (
	fiftyOneWeeks = time.Hour * 24 * 7 * 51
)

func FiftyOneWeeksBeforeMondayThisWeek(t time.Time) time.Time {
	if t.Weekday() == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(t.Weekday())+1)
	}
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC).Add(-fiftyOneWeeks)
}

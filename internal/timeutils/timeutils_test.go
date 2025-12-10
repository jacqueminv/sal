package timeutils

import (
	"testing"
	"time"
)

func TestFiftyOneWeeksBeforeMondayThisWeek(t *testing.T) {
	var dayTests = []struct {
		given    time.Time
		expected time.Time
	}{
		{time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC), time.Date(2023, 12, 4, 0, 0, 0, 0, time.UTC)},
		{time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(2024, 12, 27, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
		{time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)},
		{time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC), time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)},
	}
	for _, test := range dayTests {
		result := FiftyOneWeeksBeforeMondayThisWeek(test.given)
		year, month, day := result.Date()
		expectedYear, expectedMonth, expectedDay := test.expected.Date()
		if year != expectedYear || month != expectedMonth || day != expectedDay {
			t.Errorf("Unexpected `fiftyOneWeeksBeforeMondayThisWeek` for %v should have been %v but is %v", test.given, test.expected, result)
		}
	}
}

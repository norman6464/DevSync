package domain

import "time"

// StartOfMonth は指定時刻が属する月の1日 00:00:00 を返す。
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// StartOfWeek は指定時刻が属する週の月曜日 00:00:00 を返す。
func StartOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysBack := int(weekday) - int(time.Monday)
	monday := t.AddDate(0, 0, -daysBack)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, t.Location())
}

// DaysAgo は指定時刻から n 日前の時刻を返す。
func DaysAgo(t time.Time, n int) time.Time {
	return t.AddDate(0, 0, -n)
}

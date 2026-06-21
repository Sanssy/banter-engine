package notify

import "time"

const (
	quietHourStart = 0
	quietHourEnd   = 8
)

// IsQuietHour returns true between 00h00 and 07h59 — no notifications sent.
func IsQuietHour(t time.Time) bool {
	h := t.Hour()
	return h >= quietHourStart && h < quietHourEnd
}

// IsNightSummaryHour returns true at 08h00 — the morning digest window.
func IsNightSummaryHour(t time.Time) bool {
	return t.Hour() == quietHourEnd
}

package logfile

import (
	"time"

	"im/internal/config"
)

// ShouldEmitHeading determines whether a new timestamp subheading should be
// written, and if so, what time to display. Returns (false, zero) when no
// heading should be emitted.
func ShouldEmitHeading(strategy config.TimestampRounding, now, last time.Time, hasLast bool) (bool, time.Time) {
	switch strategy {
	case config.TimestampRoundingRound10:
		return shouldEmitRound10(now, last, hasLast)
	default:
		return shouldEmitAdaptive(now, last, hasLast)
	}
}

func shouldEmitAdaptive(now, last time.Time, hasLast bool) (bool, time.Time) {
	if !hasLast {
		return true, now
	}

	gap := now.Sub(last)
	if gap >= 10*time.Minute {
		return true, now
	}

	// Within 10m — only emit on a round 10-minute boundary.
	if now.Minute()%10 == 0 {
		return true, roundDown10(now)
	}

	return false, time.Time{}
}

func shouldEmitRound10(now, last time.Time, hasLast bool) (bool, time.Time) {
	rounded := roundDown10(now)
	if hasLast && roundDown10(last).Equal(rounded) {
		return false, time.Time{}
	}
	return true, rounded
}

// roundDown10 truncates minutes to the nearest 10-minute boundary.
func roundDown10(t time.Time) time.Time {
	m := t.Minute()
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), m-m%10, 0, 0, t.Location())
}

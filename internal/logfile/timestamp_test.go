package logfile_test

import (
	"testing"
	"time"

	"im/internal/config"
	"im/internal/logfile"
)

func tm(h, m int) time.Time {
	return time.Date(2026, 1, 23, h, m, 0, 0, time.UTC)
}

func TestShouldEmitHeading(t *testing.T) {
	tests := []struct {
		name     string
		strategy config.TimestampRounding
		now      time.Time
		last     time.Time
		hasLast  bool
		wantEmit bool
		wantH    int
		wantM    int
	}{
		// --- Adaptive: no prior heading ---
		{
			name:     "adaptive/no prior heading uses exact time",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(10, 3),
			hasLast:  false,
			wantEmit: true,
			wantH:    10,
			wantM:    3,
		},
		// --- Adaptive: gap >= 10m ---
		{
			name:     "adaptive/gap >= 10m uses exact time",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(10, 47),
			last:     tm(10, 10),
			hasLast:  true,
			wantEmit: true,
			wantH:    10,
			wantM:    47,
		},
		// --- Adaptive: gap < 10m, not on round boundary ---
		{
			name:     "adaptive/gap < 10m not on boundary suppresses",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(10, 7),
			last:     tm(10, 3),
			hasLast:  true,
			wantEmit: false,
		},
		// --- Adaptive: gap < 10m, on round boundary ---
		{
			name:     "adaptive/gap < 10m on round boundary emits rounded",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(10, 10),
			last:     tm(10, 3),
			hasLast:  true,
			wantEmit: true,
			wantH:    10,
			wantM:    10,
		},
		// --- Adaptive: exactly 10m gap ---
		{
			name:     "adaptive/exactly 10m gap uses exact time",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(10, 13),
			last:     tm(10, 3),
			hasLast:  true,
			wantEmit: true,
			wantH:    10,
			wantM:    13,
		},
		// --- Adaptive: gap < 10m on :00 boundary ---
		{
			name:     "adaptive/gap < 10m on :00 emits rounded",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(11, 0),
			last:     tm(10, 55),
			hasLast:  true,
			wantEmit: true,
			wantH:    11,
			wantM:    0,
		},
		// --- Adaptive: gap < 10m on :20 boundary ---
		{
			name:     "adaptive/gap < 10m on :20 emits rounded",
			strategy: config.TimestampRoundingAdaptive,
			now:      tm(10, 20),
			last:     tm(10, 14),
			hasLast:  true,
			wantEmit: true,
			wantH:    10,
			wantM:    20,
		},
		// --- Round10: basic rounding ---
		{
			name:     "round10/rounds down to 10m boundary",
			strategy: config.TimestampRoundingRound10,
			now:      tm(10, 7),
			hasLast:  false,
			wantEmit: true,
			wantH:    10,
			wantM:    0,
		},
		{
			name:     "round10/exact boundary unchanged",
			strategy: config.TimestampRoundingRound10,
			now:      tm(10, 30),
			hasLast:  false,
			wantEmit: true,
			wantH:    10,
			wantM:    30,
		},
		// --- Round10: same rounded time as last suppresses ---
		{
			name:     "round10/same rounded time as last suppresses",
			strategy: config.TimestampRoundingRound10,
			now:      tm(10, 7),
			last:     tm(10, 3),
			hasLast:  true,
			wantEmit: false,
		},
		// --- Round10: different rounded time emits ---
		{
			name:     "round10/different rounded time emits",
			strategy: config.TimestampRoundingRound10,
			now:      tm(10, 13),
			last:     tm(10, 3),
			hasLast:  true,
			wantEmit: true,
			wantH:    10,
			wantM:    10,
		},
		// --- Round10: midnight ---
		{
			name:     "round10/midnight",
			strategy: config.TimestampRoundingRound10,
			now:      tm(0, 5),
			hasLast:  false,
			wantEmit: true,
			wantH:    0,
			wantM:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emit, display := logfile.ShouldEmitHeading(tt.strategy, tt.now, tt.last, tt.hasLast)
			if emit != tt.wantEmit {
				t.Fatalf("emit = %v, want %v", emit, tt.wantEmit)
			}
			if !emit {
				return
			}
			if display.Hour() != tt.wantH || display.Minute() != tt.wantM {
				t.Errorf("display = %02d:%02d, want %02d:%02d",
					display.Hour(), display.Minute(), tt.wantH, tt.wantM)
			}
		})
	}
}

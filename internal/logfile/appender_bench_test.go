package logfile_test

import (
	"testing"
	"time"

	"im/internal/config"
	"im/internal/logfile"
)

// BenchmarkAppend_Inline measures the full inline-mode write path (NF-002: < 100ms).
func BenchmarkAppend_Inline(b *testing.B) {
	dir := b.TempDir()
	cfg := config.Config{TimestampRounding: config.TimestampRoundingAdaptive}
	now := time.Date(2026, 3, 9, 14, 30, 0, 0, time.UTC)

	b.ResetTimer()
	for b.Loop() {
		a := logfile.NewAppender(func() time.Time { return now }, cfg)
		if err := a.Append(dir, "benchmark entry", false); err != nil {
			b.Fatal(err)
		}
	}

	// Fail if a single append takes >= 100ms.
	if b.Elapsed() > 0 && b.Elapsed()/time.Duration(b.N) >= 100*time.Millisecond {
		b.Fatalf("append too slow: %v per op (NF-002 requires < 100ms)", b.Elapsed()/time.Duration(b.N))
	}
}

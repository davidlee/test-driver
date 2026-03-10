package logfile_test

import (
	"testing"

	"im/internal/logfile"
)

func TestParseLastTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantOK  bool
		wantH   int
		wantM   int
	}{
		{
			name:    "no headings",
			content: "# Wednesday, 23 January 2026\n\nsome text\n",
			wantOK:  false,
		},
		{
			name:    "single heading",
			content: "# Wednesday, 23 January 2026\n\n## 10:03\n\nhello\n",
			wantOK:  true,
			wantH:   10,
			wantM:   3,
		},
		{
			name:    "multiple headings returns last",
			content: "## 10:03\n\nhello\n\n## 14:30\n\nworld\n",
			wantOK:  true,
			wantH:   14,
			wantM:   30,
		},
		{
			name:    "heading with trailing space",
			content: "## 09:15   \n\ntext\n",
			wantOK:  true,
			wantH:   9,
			wantM:   15,
		},
		{
			name:    "malformed heading ignored",
			content: "## 10:3\n\n## not-a-time\n\n## 25:00\n",
			wantOK:  false,
		},
		{
			name:    "heading not at line start ignored",
			content: "text ## 10:03\n",
			wantOK:  false,
		},
		{
			name:    "empty content",
			content: "",
			wantOK:  false,
		},
		{
			name:    "midnight",
			content: "## 00:00\n\nlate night\n",
			wantOK:  true,
			wantH:   0,
			wantM:   0,
		},
		{
			name:    "end of day",
			content: "## 23:59\n\nalmost midnight\n",
			wantOK:  true,
			wantH:   23,
			wantM:   59,
		},
		// 12h format cases.
		{
			name:    "12h afternoon",
			content: "## 2:30 PM\n\nhello\n",
			wantOK:  true,
			wantH:   14,
			wantM:   30,
		},
		{
			name:    "12h morning",
			content: "## 9:15 AM\n\nhello\n",
			wantOK:  true,
			wantH:   9,
			wantM:   15,
		},
		{
			name:    "12h noon",
			content: "## 12:00 PM\n\nlunch\n",
			wantOK:  true,
			wantH:   12,
			wantM:   0,
		},
		{
			name:    "12h midnight",
			content: "## 12:00 AM\n\nlate\n",
			wantOK:  true,
			wantH:   0,
			wantM:   0,
		},
		{
			name:    "12h mixed with 24h returns last",
			content: "## 10:03\n\nhello\n\n## 2:30 PM\n\nworld\n",
			wantOK:  true,
			wantH:   14,
			wantM:   30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := logfile.ParseLastTimestamp(tt.content)
			if ok != tt.wantOK {
				t.Fatalf("ParseLastTimestamp() ok = %v, want %v", ok, tt.wantOK)
			}
			if !ok {
				return
			}
			if got.Hour() != tt.wantH || got.Minute() != tt.wantM {
				t.Errorf("ParseLastTimestamp() = %02d:%02d, want %02d:%02d",
					got.Hour(), got.Minute(), tt.wantH, tt.wantM)
			}
		})
	}
}

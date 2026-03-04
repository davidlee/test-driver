package entry_test

import (
	"testing"

	"im/internal/entry"
)

func TestFormat(t *testing.T) {
	tests := []struct {
		name string
		body string
		task bool
		want string
	}{
		{
			name: "single word",
			body: "hello",
			task: false,
			want: "hello\n",
		},
		{
			name: "multiple words",
			body: "hello world",
			task: false,
			want: "hello world\n",
		},
		{
			name: "multi-line",
			body: "line one\nline two",
			task: false,
			want: "line one\nline two\n",
		},
		{
			name: "trailing newline preserved as single",
			body: "hello\n",
			task: false,
			want: "hello\n",
		},
		{
			name: "trailing whitespace trimmed",
			body: "hello   \n  ",
			task: false,
			want: "hello\n",
		},
		{
			name: "leading whitespace trimmed",
			body: "  \nhello",
			task: false,
			want: "hello\n",
		},
		{
			name: "empty string",
			body: "",
			task: false,
			want: "",
		},
		{
			name: "whitespace only",
			body: "   \n  \n  ",
			task: false,
			want: "",
		},
		{
			name: "task single line",
			body: "buy milk",
			task: true,
			want: "- [ ] buy milk\n",
		},
		{
			name: "task multi-line",
			body: "buy milk\nwalk dog",
			task: true,
			want: "- [ ] buy milk\n- [ ] walk dog\n",
		},
		{
			name: "task skips blank lines",
			body: "buy milk\n\nwalk dog",
			task: true,
			want: "- [ ] buy milk\n- [ ] walk dog\n",
		},
		{
			name: "task empty input",
			body: "",
			task: true,
			want: "",
		},
		{
			name: "task whitespace only",
			body: "  \n  ",
			task: true,
			want: "",
		},
		{
			name: "task trims lines before prefixing",
			body: "  buy milk  \n  walk dog  ",
			task: true,
			want: "- [ ] buy milk\n- [ ] walk dog\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entry.Format(tt.body, tt.task)
			if got != tt.want {
				t.Errorf("Format(%q, %v)\n got: %q\nwant: %q", tt.body, tt.task, got, tt.want)
			}
		})
	}
}

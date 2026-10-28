package ui

import (
	"testing"
	"time"
)

func TestBigClockIsAligned(t *testing.T) {
	rows := BigClock(25*time.Minute + 3*time.Second)
	if len(rows) != 5 {
		t.Fatalf("rows = %d, want 5", len(rows))
	}
	width := len([]rune(rows[0]))
	for i, row := range rows {
		if got := len([]rune(row)); got != width {
			t.Fatalf("row %d width = %d, want %d (%q)", i, got, width, row)
		}
	}
}

func TestClockText(t *testing.T) {
	cases := map[time.Duration]string{
		65 * time.Second:                        "01:05",
		25 * time.Minute:                        "25:00",
		time.Hour:                               "1:00:00",
		-5 * time.Second:                        "00:00",
		1500*time.Second + 400*time.Millisecond: "25:00",
	}
	for d, want := range cases {
		if got := ClockText(d); got != want {
			t.Errorf("ClockText(%s) = %q, want %q", d, got, want)
		}
	}
}

func TestProgressBarWidthIsStable(t *testing.T) {
	bar := ProgressBar(20, 0.5, "#7aa2f7", "#565f89")
	if got := len([]rune(stripANSI(bar))); got != 20 {
		t.Fatalf("bar width = %d, want 20", got)
	}
}

// stripANSI is enough for the single SGR sequence lipgloss emits here.
func stripANSI(s string) string {
	out := make([]rune, 0, len(s))
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		out = append(out, r)
	}
	return string(out)
}

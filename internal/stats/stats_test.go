package stats

import (
	"testing"
	"time"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/store"
)

func TestSummarizeLastSevenDays(t *testing.T) {
	now := time.Date(2026, 10, 28, 15, 0, 0, 0, time.Local)
	recs := []store.Record{
		{StartedAt: now.Add(-1 * time.Hour), FocusedSec: 1500, Completed: true},
		{StartedAt: now.Add(-24 * time.Hour), FocusedSec: 1500, Completed: true},
		{StartedAt: now.Add(-48 * time.Hour), FocusedSec: 600, Completed: true},
		{StartedAt: now.Add(-20 * 24 * time.Hour), FocusedSec: 1500, Completed: true},
		{StartedAt: now.Add(-3 * time.Hour), FocusedSec: 1500, Completed: false},
	}

	got := Summarize(recs, now, 2)
	if got.TodayFocusSec != 1500 || got.TodayBlocks != 1 {
		t.Fatalf("today = %d/%d, want 1500/1", got.TodayFocusSec, got.TodayBlocks)
	}
	if got.ActiveDays != 3 {
		t.Fatalf("active days = %d, want 3", got.ActiveDays)
	}
	if got.TotalFocusSec != 5100 || got.TotalBlocks != 4 {
		t.Fatalf("total = %d/%d, want 5100/4", got.TotalFocusSec, got.TotalBlocks)
	}
	if len(got.Days) != 7 {
		t.Fatalf("days = %d, want 7", len(got.Days))
	}
	if got.ParkingCount != 2 {
		t.Fatalf("parking = %d, want 2", got.ParkingCount)
	}
}

func TestFormatDuration(t *testing.T) {
	cases := map[int64]string{
		0:    "0m",
		59:   "1m",
		1500: "25m",
		3600: "1h",
		3900: "1h 05m",
	}
	for sec, want := range cases {
		if got := FormatDuration(sec); got != want {
			t.Errorf("FormatDuration(%d) = %q, want %q", sec, got, want)
		}
	}
}

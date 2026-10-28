// Package stats derives neutral, non-judgmental summaries from session
// history. There are deliberately no streaks, completion rates, or grades.
package stats

import (
	"fmt"
	"sort"
	"time"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/store"
)

type Day struct {
	Date     time.Time `json:"date"`
	FocusSec int64     `json:"focus_sec"`
	Blocks   int       `json:"blocks"`
}

type Summary struct {
	TodayFocusSec int64 `json:"today_focus_sec"`
	TodayBlocks   int   `json:"today_blocks"`
	TotalFocusSec int64 `json:"total_focus_sec"`
	TotalBlocks   int   `json:"total_blocks"`
	ActiveDays    int   `json:"active_days"`
	Days          []Day `json:"days"`
	ParkingCount  int   `json:"parking_count"`
}

// Summarize builds the last seven local calendar days, oldest first.
func Summarize(records []store.Record, now time.Time, parkingCount int) Summary {
	now = now.Local()
	start := startOfDay(now).AddDate(0, 0, -6)
	days := make([]Day, 7)
	index := make(map[string]int, 7)
	for i := range days {
		d := start.AddDate(0, 0, i)
		days[i] = Day{Date: d}
		index[d.Format("2006-01-02")] = i
	}

	var sum Summary
	for _, rec := range records {
		if !rec.Completed || rec.FocusedSec <= 0 {
			continue
		}
		sum.TotalFocusSec += rec.FocusedSec
		sum.TotalBlocks++
		key := rec.StartedAt.Local().Format("2006-01-02")
		i, ok := index[key]
		if !ok {
			continue
		}
		days[i].FocusSec += rec.FocusedSec
		days[i].Blocks++
	}

	sum.Days = days
	sum.ParkingCount = parkingCount
	for _, d := range days {
		if d.Blocks > 0 {
			sum.ActiveDays++
		}
	}
	today := days[len(days)-1]
	sum.TodayFocusSec = today.FocusSec
	sum.TodayBlocks = today.Blocks
	return sum
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// FormatDuration renders a compact human duration: 0m, 25m, 1h 05m.
func FormatDuration(sec int64) string {
	if sec <= 0 {
		return "0m"
	}
	mins := (sec + 30) / 60
	if mins < 60 {
		return fmt.Sprintf("%dm", mins)
	}
	h := mins / 60
	m := mins % 60
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %02dm", h, m)
}

// SortDays is handy for callers that build a day slice themselves.
func SortDays(days []Day) {
	sort.Slice(days, func(i, j int) bool { return days[i].Date.Before(days[j].Date) })
}

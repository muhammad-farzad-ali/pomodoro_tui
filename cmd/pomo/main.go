package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/app"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/stats"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/store"
)

var version = "0.1.0"

func main() {
	var (
		statsFlag  = flag.Bool("stats", false, "print today and last-7-days summary, then exit")
		jsonFlag   = flag.Bool("json", false, "with --stats, print machine-readable JSON")
		versionF   = flag.Bool("version", false, "print version and exit")
		duration   = flag.Int("duration", 0, "focus length in minutes for this run")
		breakMin   = flag.Int("break", 0, "short-break length in minutes for this run")
		activation = flag.Bool("activation", false, "open with a 5-minute activation sprint")
		noNotify   = flag.Bool("no-notify", false, "disable terminal bell and desktop notification")
	)
	flag.Parse()

	if *versionF {
		fmt.Printf("pomo %s\n", version)
		return
	}

	st := store.New(store.DefaultDir())

	if *statsFlag {
		printStats(st, *jsonFlag)
		return
	}
	if *jsonFlag {
		fmt.Fprintln(os.Stderr, "pomo: --json is only valid with --stats")
		os.Exit(2)
	}

	release, err := st.Lock()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pomo:", err)
		os.Exit(1)
	}
	defer release()

	model, err := app.New(st, app.Options{
		FocusMinutes: *duration,
		BreakMinutes: *breakMin,
		Activation:   *activation,
		NoNotify:     *noNotify,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "pomo:", err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		p.Send(app.QuitMsg{})
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "pomo:", err)
		os.Exit(1)
	}
}

func printStats(st *store.Store, asJSON bool) {
	records, err := st.Sessions()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pomo:", err)
		os.Exit(1)
	}
	state, err := st.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pomo:", err)
	}
	parked := 0
	for _, task := range state.Parking {
		if !task.Done {
			parked++
		}
	}
	summary := stats.Summarize(records, time.Now(), parked)

	if asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(summary); err != nil {
			fmt.Fprintln(os.Stderr, "pomo:", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("today  %s · %d blocks · %d parked\n", stats.FormatDuration(summary.TodayFocusSec), summary.TodayBlocks, summary.ParkingCount)
	fmt.Println("last 7 days")
	for _, day := range summary.Days {
		value := stats.FormatDuration(day.FocusSec)
		if day.FocusSec == 0 {
			value = "·"
		}
		fmt.Printf("  %s  %-6s  %d block(s)\n", day.Date.Format("Mon 02"), value, day.Blocks)
	}
	fmt.Printf("%d active days · no streaks, no grades\n", summary.ActiveDays)
}

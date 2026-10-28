package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/session"
)

func testState(now time.Time) State {
	st := DefaultState()
	st.FirstRunDone = true
	st.NextAction = "write the evidence doc"
	st.IfThen = "if I check email, then I park it"
	st.Parking = []Task{{ID: "1", Text: "look up a fact", CreatedAt: now}}
	st.Session = session.NewFocus(25*time.Minute, st.NextAction, st.IfThen, 0, false)
	st.Session.Begin(now)
	st.LastSaved = now
	return st
}

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	st := New(dir)
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	want := testState(now)

	if err := st.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.NextAction != want.NextAction {
		t.Fatalf("next action = %q, want %q", got.NextAction, want.NextAction)
	}
	if got.IfThen != want.IfThen {
		t.Fatalf("if-then = %q, want %q", got.IfThen, want.IfThen)
	}
	if len(got.Parking) != 1 || got.Parking[0].Text != "look up a fact" {
		t.Fatalf("parking = %#v", got.Parking)
	}
	if got.Session == nil || got.Session.PlannedSec != 1500 {
		t.Fatalf("session = %#v", got.Session)
	}

	info, err := os.Stat(filepath.Join(dir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("state mode = %o, want 600", got)
	}
}

func TestCorruptStateIsPreserved(t *testing.T) {
	dir := t.TempDir()
	st := New(dir)
	if err := os.WriteFile(filepath.Join(dir, "state.json"), []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := st.Load()
	if err == nil {
		t.Fatal("expected a corrupt-state error")
	}
	if !got.FirstRunDone {
		// Default state is a fresh start, but it must still be usable.
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "state.json.corrupt-*"))
	if len(matches) != 1 {
		t.Fatalf("backup files = %d, want 1", len(matches))
	}
}

func TestAppendAndReadHistory(t *testing.T) {
	dir := t.TempDir()
	st := New(dir)
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)

	rec := Record{
		StartedAt:  now,
		EndedAt:    now.Add(25 * time.Minute),
		PlannedSec: 1500,
		FocusedSec: 1500,
		Completed:  true,
		Intention:  "write",
	}
	if err := st.AppendSession(rec); err != nil {
		t.Fatal(err)
	}
	// A malformed line (e.g. crash during append) must not poison the rest.
	f, err := os.OpenFile(filepath.Join(dir, "sessions.jsonl"), os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString("{broken\n"); err != nil {
		t.Fatal(err)
	}
	f.Close()

	got, err := st.Sessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].FocusedSec != 1500 {
		t.Fatalf("records = %#v", got)
	}
}

func TestLockIsExclusive(t *testing.T) {
	dir := t.TempDir()
	st := New(dir)

	release, err := st.Lock()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := st.Lock(); err == nil {
		t.Fatal("second lock should fail")
	}
	release()

	release2, err := st.Lock()
	if err != nil {
		t.Fatalf("lock after release: %v", err)
	}
	release2()
}

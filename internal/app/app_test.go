package app

import (
	"strings"
	"testing"
	"time"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/session"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/store"
)

func testModel(t *testing.T, now time.Time) (*Model, *store.Store) {
	t.Helper()
	dir := t.TempDir()
	st := store.New(dir)
	state := store.DefaultState()
	state.FirstRunDone = true
	state.NextAction = "write the evidence doc"
	if err := st.Save(state); err != nil {
		t.Fatal(err)
	}
	m, err := New(st, Options{Now: func() time.Time { return now }, NoNotify: true})
	if err != nil {
		t.Fatal(err)
	}
	return m, st
}

func TestFirstRunStartsWithIntention(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	st := store.New(dir)
	if err := st.Save(store.DefaultState()); err != nil {
		t.Fatal(err)
	}
	m, err := New(st, Options{Now: func() time.Time { return now }, NoNotify: true})
	if err != nil {
		t.Fatal(err)
	}
	if m.modal != modalIntention {
		t.Fatalf("modal = %v, want intention", m.modal)
	}

	m.nextInput.SetValue("open the document")
	m.ifInput.SetValue("if I check email, then I park it")
	m.saveIntention(now)

	if m.modal != modalNone {
		t.Fatalf("modal = %v, want none", m.modal)
	}
	if m.sess == nil || m.sess.Status != session.StatusRunning {
		t.Fatalf("session = %#v, want running", m.sess)
	}
	if m.sess.Intention != "open the document" {
		t.Fatalf("intention = %q", m.sess.Intention)
	}

	state, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	if !state.FirstRunDone {
		t.Fatal("first_run_done not persisted")
	}
	if state.NextAction != "open the document" {
		t.Fatalf("next action = %q", state.NextAction)
	}
}

func TestCompletionLogsOnceAndSchedulesBreak(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, st := testModel(t, now)
	m.sess.RemainingSec = 1
	m.sess.Deadline = now.Add(time.Second)

	m.advance(now.Add(2 * time.Second))

	if m.sess.Status != session.StatusFinished {
		t.Fatalf("status = %s, want finished", m.sess.Status)
	}
	records, err := st.Sessions()
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatalf("records = %d, want 1", len(records))
	}
	if records[0].FocusedSec != 1500 {
		t.Fatalf("focused = %d, want 1500", records[0].FocusedSec)
	}
	if m.state.Cycle != 1 {
		t.Fatalf("cycle = %d, want 1", m.state.Cycle)
	}
	if m.autoBreakAt.IsZero() {
		t.Fatal("auto-break was not scheduled")
	}

	// Calling finish again must not duplicate history.
	m.handleFinished(now.Add(2*time.Second), true)
	records, _ = st.Sessions()
	if len(records) != 1 {
		t.Fatalf("records after duplicate finish = %d, want 1", len(records))
	}

	m.advance(m.autoBreakAt.Add(time.Second))
	if m.sess == nil || !m.sess.IsBreak() {
		t.Fatalf("session = %#v, want break", m.sess)
	}
	if m.sess.Status != session.StatusRunning {
		t.Fatalf("break status = %s, want running", m.sess.Status)
	}
}

func TestParkingCaptureAndPromote(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, _ := testModel(t, now)

	m.parkingInput.SetValue("check email")
	m.addParking(now)
	if len(m.state.Parking) != 1 {
		t.Fatalf("parking = %d, want 1", len(m.state.Parking))
	}

	m.parkingSel = 0
	task := m.state.Parking[0]
	m.state.NextAction = task.Text
	m.sess.Intention = task.Text
	m.removeParking(0)
	if len(m.state.Parking) != 0 {
		t.Fatalf("parking after promote = %d, want 0", len(m.state.Parking))
	}
	if m.state.NextAction != "check email" {
		t.Fatalf("next action = %q", m.state.NextAction)
	}
}

func TestExtendAfterCompletionStartsNewBlock(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, _ := testModel(t, now)
	m.sess.RemainingSec = 1
	m.sess.Deadline = now.Add(time.Second)
	m.advance(now.Add(2 * time.Second))
	if m.sess.Status != session.StatusFinished || !m.sess.Logged {
		t.Fatalf("session = %#v, want finished and logged", m.sess)
	}

	m.extend(now.Add(3 * time.Second))
	if m.sess.Status != session.StatusRunning {
		t.Fatalf("status = %s, want running", m.sess.Status)
	}
	if m.sess.PlannedSec != 600 {
		t.Fatalf("planned = %d, want 600", m.sess.PlannedSec)
	}
	if m.sess.Logged {
		t.Fatal("new extension block should not be marked logged")
	}
}

func TestBreakCanAutoStartFocus(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, _ := testModel(t, now)
	m.cfg.AutoStartFocus = true
	m.sess = session.NewBreak(session.PhaseShortBreak, time.Minute, 0)
	m.sess.Begin(now)
	m.sess.RemainingSec = 1
	m.sess.Deadline = now.Add(time.Second)

	m.advance(now.Add(2 * time.Second))
	if m.sess.Status != session.StatusFinished {
		t.Fatalf("break status = %s, want finished", m.sess.Status)
	}
	if m.autoFocusAt.IsZero() {
		t.Fatal("auto-focus was not scheduled")
	}

	m.advance(m.autoFocusAt.Add(time.Second))
	if m.sess == nil || !m.sess.IsFocus() {
		t.Fatalf("session = %#v, want focus", m.sess)
	}
	if m.sess.Status != session.StatusRunning {
		t.Fatalf("focus status = %s, want running", m.sess.Status)
	}
}

func BenchmarkViewRunning(b *testing.B) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	dir := b.TempDir()
	st := store.New(dir)
	state := store.DefaultState()
	state.FirstRunDone = true
	state.NextAction = "write the evidence doc"
	if err := st.Save(state); err != nil {
		b.Fatal(err)
	}
	m, err := New(st, Options{Now: func() time.Time { return now }, NoNotify: true})
	if err != nil {
		b.Fatal(err)
	}
	m.width, m.height = 100, 30
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.View()
	}
}

func TestCLIOverrideDoesNotPersist(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	dir := t.TempDir()
	st := store.New(dir)
	state := store.DefaultState()
	state.FirstRunDone = true
	if err := st.Save(state); err != nil {
		t.Fatal(err)
	}

	m, err := New(st, Options{FocusMinutes: 50, NoNotify: true, Now: func() time.Time { return now }})
	if err != nil {
		t.Fatal(err)
	}
	if m.sess.PlannedSec != 3000 {
		t.Fatalf("session planned = %d, want 3000", m.sess.PlannedSec)
	}
	m.save(now)

	got, err := st.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Config.FocusMinutes != 25 {
		t.Fatalf("persisted focus = %d, want 25", got.Config.FocusMinutes)
	}
}

func TestPauseStopsTicking(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, _ := testModel(t, now)
	m.togglePause(now)
	if m.sess.Status != session.StatusPaused {
		t.Fatalf("status = %s, want paused", m.sess.Status)
	}
	if m.ticking {
		t.Fatal("still ticking after pause")
	}
	if cmd := m.startTicking(); cmd != nil {
		t.Fatal("startTicking returned a command while paused")
	}
}

func TestStaleTickIsIgnored(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, _ := testModel(t, now)
	m.tickSeq = 7
	m.ticking = true
	before := m.sess.RemainingSec

	if cmd := m.onTick(tickMsg{at: now.Add(time.Minute), seq: 6}); cmd != nil {
		t.Fatal("stale tick returned a command")
	}
	if m.sess.RemainingSec != before {
		t.Fatalf("stale tick advanced session: %d -> %d", before, m.sess.RemainingSec)
	}
	if !m.ticking {
		t.Fatal("stale tick cleared the active chain flag")
	}
}

func TestViewIsSafeAcrossWidths(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	m, _ := testModel(t, now)

	for _, w := range []int{24, 40, 60, 80, 120} {
		m.width = w
		m.height = 24
		out := m.View()
		if out == "" {
			t.Fatalf("width %d rendered empty", w)
		}
		if strings.Contains(out, "%!") {
			t.Fatalf("width %d has formatting error: %s", w, out)
		}
	}
}

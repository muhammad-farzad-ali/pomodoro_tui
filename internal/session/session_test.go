package session

import (
	"testing"
	"time"
)

func TestCountdownReachesFinished(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	s := NewFocus(25*time.Minute, "write", "", 0, false)
	s.Begin(now)

	if got := s.RemainingAt(now.Add(5 * time.Minute)); got != 20*time.Minute {
		t.Fatalf("remaining after 5m = %s, want 20m", got)
	}
	s.Tick(now.Add(25 * time.Minute))
	if s.Status != StatusFinished {
		t.Fatalf("status = %s, want finished", s.Status)
	}
	if s.RemainingSec != 0 {
		t.Fatalf("remaining = %d, want 0", s.RemainingSec)
	}
}

func TestPauseFreezesAndResumeRestoresDeadline(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	s := NewFocus(25*time.Minute, "", "", 0, false)
	s.Begin(now)
	s.Pause(now.Add(2 * time.Minute))

	if s.Status != StatusPaused {
		t.Fatalf("status = %s, want paused", s.Status)
	}
	if got := s.RemainingAt(now.Add(20 * time.Minute)); got != 23*time.Minute {
		t.Fatalf("paused remaining = %s, want 23m", got)
	}

	s.Resume(now.Add(20 * time.Minute))
	if got := s.RemainingAt(now.Add(21 * time.Minute)); got != 22*time.Minute {
		t.Fatalf("resumed remaining = %s, want 22m", got)
	}
}

func TestExtendRunningAndFinished(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	s := NewFocus(25*time.Minute, "", "", 0, false)
	s.Begin(now)
	s.Extend(now.Add(time.Minute), 10*time.Minute)
	if got := s.RemainingAt(now.Add(5 * time.Minute)); got != 30*time.Minute {
		t.Fatalf("running extended remaining = %s, want 30m", got)
	}
	if s.Total() != 35*time.Minute {
		t.Fatalf("total = %s, want 35m", s.Total())
	}

	s.Tick(now.Add(35 * time.Minute))
	if s.Status != StatusFinished {
		t.Fatalf("status = %s, want finished", s.Status)
	}

	s.Extend(now.Add(35*time.Minute), 10*time.Minute)
	if s.Status != StatusRunning {
		t.Fatalf("status after extend = %s, want running", s.Status)
	}
	if got := s.RemainingAt(now.Add(36 * time.Minute)); got != 9*time.Minute {
		t.Fatalf("extended finished remaining = %s, want 9m", got)
	}
	if s.Extensions != 2 {
		t.Fatalf("extensions = %d, want 2", s.Extensions)
	}
}

func TestRemainingIsDriftFreeAcrossSnapshot(t *testing.T) {
	now := time.Date(2026, 10, 28, 9, 0, 0, 0, time.UTC)
	s := NewFocus(25*time.Minute, "", "", 0, false)
	s.Begin(now)

	// Simulate a save/load cycle: the deadline, not the tick count, is truth.
	copy := *s
	if got := copy.RemainingAt(now.Add(7*time.Minute + 30*time.Second)); got != 17*time.Minute+30*time.Second {
		t.Fatalf("remaining after restore = %s, want 17m30s", got)
	}
}

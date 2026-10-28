// Package session implements the drift-free focus/break state machine.
//
// A running session stores an absolute deadline. Every render recomputes the
// remaining time from the wall clock, so sleep, suspend, missed ticks, and
// terminal resizes do not accumulate error.
package session

import (
	"math"
	"time"
)

type Phase string

const (
	PhaseFocus      Phase = "focus"
	PhaseShortBreak Phase = "short_break"
	PhaseLongBreak  Phase = "long_break"
)

type Status string

const (
	StatusIdle     Status = "idle"
	StatusRunning  Status = "running"
	StatusPaused   Status = "paused"
	StatusFinished Status = "finished"
)

// Session is both the live state and its on-disk representation. Seconds are
// used in JSON so the file stays readable.
type Session struct {
	Phase        Phase     `json:"phase"`
	Status       Status    `json:"status"`
	PlannedSec   int64     `json:"planned_sec"`
	RemainingSec int64     `json:"remaining_sec"`
	ExtraSec     int64     `json:"extra_sec"`
	Deadline     time.Time `json:"deadline,omitempty"`
	StartedAt    time.Time `json:"started_at,omitempty"`
	Extensions   int       `json:"extensions,omitempty"`
	Cycle        int       `json:"cycle,omitempty"`
	Intention    string    `json:"intention,omitempty"`
	IfThen       string    `json:"if_then,omitempty"`
	Activation   bool      `json:"activation,omitempty"`
	Logged       bool      `json:"logged,omitempty"`
}

func NewFocus(d time.Duration, intention, ifThen string, cycle int, activation bool) *Session {
	sec := int64(d / time.Second)
	return &Session{
		Phase:        PhaseFocus,
		Status:       StatusIdle,
		PlannedSec:   sec,
		RemainingSec: sec,
		Cycle:        cycle,
		Intention:    intention,
		IfThen:       ifThen,
		Activation:   activation,
	}
}

func NewBreak(phase Phase, d time.Duration, cycle int) *Session {
	sec := int64(d / time.Second)
	return &Session{
		Phase:        phase,
		Status:       StatusIdle,
		PlannedSec:   sec,
		RemainingSec: sec,
		Cycle:        cycle,
	}
}

func (s *Session) IsFocus() bool {
	return s != nil && s.Phase == PhaseFocus
}

func (s *Session) IsBreak() bool {
	return s != nil && (s.Phase == PhaseShortBreak || s.Phase == PhaseLongBreak)
}

func (s *Session) Planned() time.Duration {
	if s == nil {
		return 0
	}
	return time.Duration(s.PlannedSec) * time.Second
}

// Total includes time added with the extend key.
func (s *Session) Total() time.Duration {
	if s == nil {
		return 0
	}
	return time.Duration(s.PlannedSec+s.ExtraSec) * time.Second
}

// Begin starts a session from its current remaining time. For a fresh
// session that is the full planned duration.
func (s *Session) Begin(now time.Time) {
	if s == nil {
		return
	}
	if s.RemainingSec <= 0 {
		s.RemainingSec = s.PlannedSec
	}
	s.Status = StatusRunning
	s.StartedAt = now
	s.Deadline = now.Add(time.Duration(s.RemainingSec) * time.Second)
}

func (s *Session) Pause(now time.Time) {
	if s == nil || s.Status != StatusRunning {
		return
	}
	s.RemainingSec = ceilSec(s.RemainingAt(now))
	s.Status = StatusPaused
	s.Deadline = time.Time{}
}

func (s *Session) Resume(now time.Time) {
	if s == nil || s.Status != StatusPaused || s.RemainingSec <= 0 {
		return
	}
	s.Status = StatusRunning
	s.Deadline = now.Add(time.Duration(s.RemainingSec) * time.Second)
	if s.StartedAt.IsZero() {
		s.StartedAt = now
	}
}

// Extend adds time to a running, paused, or just-finished session. Extending
// a finished focus session restarts it from the added time.
func (s *Session) Extend(now time.Time, d time.Duration) {
	if s == nil || d <= 0 {
		return
	}
	sec := int64(d / time.Second)
	s.ExtraSec += sec
	s.Extensions++

	switch s.Status {
	case StatusRunning:
		s.Deadline = s.Deadline.Add(d)
		s.RemainingSec += sec
	case StatusFinished:
		s.RemainingSec = sec
		s.Status = StatusRunning
		s.Deadline = now.Add(d)
	case StatusPaused, StatusIdle:
		s.RemainingSec += sec
	}
}

// Tick updates the remaining seconds and marks a running session finished.
func (s *Session) Tick(now time.Time) {
	if s == nil || s.Status != StatusRunning {
		return
	}
	rem := s.RemainingAt(now)
	if rem <= 0 {
		s.RemainingSec = 0
		s.Status = StatusFinished
		s.Deadline = time.Time{}
		return
	}
	s.RemainingSec = ceilSec(rem)
}

// RemainingAt is the source of truth while running; otherwise it returns the
// last saved remaining value.
func (s *Session) RemainingAt(now time.Time) time.Duration {
	if s == nil {
		return 0
	}
	if s.Status == StatusRunning && !s.Deadline.IsZero() {
		rem := s.Deadline.Sub(now)
		if rem < 0 {
			return 0
		}
		return rem
	}
	if s.RemainingSec < 0 {
		return 0
	}
	return time.Duration(s.RemainingSec) * time.Second
}

// SetRemaining is a helper for restoring a session that was paused or
// adjusted outside the normal flow.
func (s *Session) SetRemaining(d time.Duration) {
	if s == nil {
		return
	}
	if d < 0 {
		d = 0
	}
	s.RemainingSec = int64(d / time.Second)
}

func ceilSec(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	return int64(math.Ceil(d.Seconds()))
}

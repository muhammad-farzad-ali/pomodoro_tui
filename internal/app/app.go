// Package app is the Bubble Tea model: phase transitions, keys, persistence,
// and the small amount of UX that the evidence actually supports.
package app

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/config"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/notify"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/session"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/stats"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/store"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/ui"
)

const autoBreakDelay = 8 * time.Second

type screen int

const (
	screenMain screen = iota
	screenStats
	screenHelp
)

type modal int

const (
	modalNone modal = iota
	modalIntention
	modalParkingInput
	modalParkingList
)

type tickMsg time.Time

type toastMsg struct{ id int }

// QuitMsg lets main ask the model to save before exiting on a signal.
type QuitMsg struct{}

type Options struct {
	FocusMinutes int
	BreakMinutes int
	Activation   bool
	NoNotify     bool
	Now          func() time.Time
}

type Model struct {
	cfg   config.Config
	st    *store.Store
	state store.State
	sess  *session.Session
	theme ui.Theme

	screen screen
	modal  modal

	firstRun   bool
	activation bool
	noNotify   bool

	width  int
	height int
	now    func() time.Time

	ticking     bool
	autoBreakAt time.Time

	nextInput    textinput.Model
	ifInput      textinput.Model
	parkingInput textinput.Model
	field        int

	parkingSel    int
	parkingOffset int

	toast   string
	toastID int

	summary       stats.Summary
	todayFocusSec int64
	todayBlocks   int
}

func New(st *store.Store, opts Options) (*Model, error) {
	state, loadErr := st.Load()
	cfg := state.Config
	cfg.Normalize()
	if opts.FocusMinutes > 0 {
		cfg.FocusMinutes = clampInt(opts.FocusMinutes, 1, 180)
	}
	if opts.BreakMinutes > 0 {
		cfg.ShortBreakMinutes = clampInt(opts.BreakMinutes, 1, 60)
	}

	nowFn := opts.Now
	if nowFn == nil {
		nowFn = time.Now
	}

	m := &Model{
		cfg:      cfg,
		st:       st,
		state:    state,
		theme:    ui.TokyoNight(),
		noNotify: opts.NoNotify,
		now:      nowFn,
		width:    80,
		height:   24,
	}
	m.nextInput = newInput(m.theme, "next physical action", 52)
	m.ifInput = newInput(m.theme, "if …, then I'll park it (optional)", 52)
	m.parkingInput = newInput(m.theme, "what pulled at you?", 52)

	if loadErr != nil {
		m.toast = loadErr.Error()
		m.toastID++
	}
	m.refreshToday(nowFn())

	if !state.FirstRunDone {
		m.screen = screenMain
		m.modal = modalIntention
		m.firstRun = true
		m.nextInput.SetValue(state.NextAction)
		m.ifInput.SetValue(state.IfThen)
		m.nextInput.Focus()
		return m, loadErr
	}

	m.restoreOrStart(nowFn(), opts.Activation)
	return m, loadErr
}

func newInput(theme ui.Theme, placeholder string, width int) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Prompt = "› "
	ti.CharLimit = 200
	ti.Width = width
	ti.PromptStyle = theme.Fg(theme.Magenta)
	ti.TextStyle = theme.Fg(theme.FG)
	ti.PlaceholderStyle = theme.Fg(theme.Muted)
	ti.Cursor.Style = theme.Fg(theme.Cyan)
	return ti
}

func (m *Model) Init() tea.Cmd {
	if m.modal == modalIntention || m.modal == modalParkingInput {
		return textinput.Blink
	}
	return m.startTicking()
}

func (m *Model) restoreOrStart(now time.Time, activation bool) {
	m.activation = activation
	st := &m.state

	if st.Session != nil {
		s := *st.Session
		m.sess = &s

		switch s.Status {
		case session.StatusRunning:
			rem := s.RemainingAt(now)
			if rem <= 0 {
				s.Tick(now)
				m.handleFinished(now, true)
				return
			}
			age := now.Sub(st.LastSaved)
			if st.LastSaved.IsZero() || (age >= 0 && age <= 2*time.Hour) {
				s.Status = session.StatusRunning
				s.Deadline = now.Add(rem)
				s.RemainingSec = int64(math.Ceil(rem.Seconds()))
				m.sess = &s
				m.state.Session = m.sess
				m.activation = s.Activation
				m.save(now)
				return
			}
			m.startFocus(now, m.focusDuration(activation), activation)
			return

		case session.StatusPaused:
			m.sess = &s
			m.state.Session = m.sess
			m.activation = s.Activation
			return

		case session.StatusFinished:
			m.sess = &s
			m.state.Session = m.sess
			m.activation = s.Activation
			if s.IsFocus() && now.Sub(st.LastSaved) < 5*time.Minute {
				if !s.Activation && m.cfg.AutoStartBreak {
					m.autoBreakAt = now.Add(autoBreakDelay)
				}
				return
			}
		}
	}

	m.startFocus(now, m.focusDuration(activation), activation)
}

func (m *Model) focusDuration(activation bool) time.Duration {
	if activation {
		return m.cfg.ActivationDuration()
	}
	return m.cfg.FocusDuration()
}

func (m *Model) startFocus(now time.Time, d time.Duration, activation bool) {
	m.sess = session.NewFocus(d, m.state.NextAction, m.state.IfThen, m.state.Cycle, activation)
	m.sess.Begin(now)
	m.state.Session = m.sess
	m.activation = activation
	m.screen = screenMain
	m.autoBreakAt = time.Time{}
	m.save(now)
}

func (m *Model) startBreak(now time.Time) {
	if m.sess == nil || !m.sess.IsFocus() {
		return
	}
	phase := session.PhaseShortBreak
	d := m.cfg.ShortBreakDuration()
	if m.cfg.LongBreakEvery > 0 && m.state.Cycle > 0 && m.state.Cycle%m.cfg.LongBreakEvery == 0 {
		phase = session.PhaseLongBreak
		d = m.cfg.LongBreakDuration()
	}
	m.sess = session.NewBreak(phase, d, m.state.Cycle)
	m.sess.Begin(now)
	m.state.Session = m.sess
	m.activation = false
	m.autoBreakAt = time.Time{}
	m.screen = screenMain
	m.save(now)
}

func (m *Model) startNextFocus(now time.Time) {
	m.startFocus(now, m.cfg.FocusDuration(), false)
}

func (m *Model) handleFinished(now time.Time, restored bool) {
	if m.sess == nil {
		return
	}
	m.sess.Tick(now)
	m.autoBreakAt = time.Time{}

	if m.sess.IsFocus() {
		if !m.sess.Logged {
			rec := store.Record{
				StartedAt:  m.sess.StartedAt,
				EndedAt:    now,
				PlannedSec: m.sess.PlannedSec,
				FocusedSec: m.sess.PlannedSec + m.sess.ExtraSec,
				ExtraSec:   m.sess.ExtraSec,
				Activation: m.sess.Activation,
				Intention:  m.sess.Intention,
				Completed:  true,
			}
			if err := m.st.AppendSession(rec); err != nil {
				m.setToast(err.Error())
			} else {
				m.sess.Logged = true
			}
			m.state.Cycle++
		}
		if !m.sess.Activation && m.cfg.AutoStartBreak {
			m.autoBreakAt = now.Add(autoBreakDelay)
		}
		m.activation = m.sess.Activation
		m.refreshToday(now)
		m.save(now)
		if !restored {
			body := fmt.Sprintf("%s banked. Next: a short break.", stats.FormatDuration(m.sess.PlannedSec+m.sess.ExtraSec))
			if m.sess.Activation {
				body = "Momentum started. Continue, or take a break."
			}
			m.notify("Focus complete", body)
		}
		return
	}

	m.save(now)
	if !restored {
		m.notify("Break over", "Ready when you are.")
	}
}

func (m *Model) notify(title, body string) {
	if m.noNotify {
		return
	}
	if m.cfg.Bell {
		notify.Bell()
	}
	if m.cfg.DesktopNotify {
		notify.Desktop(title, body)
	}
}

func (m *Model) save(now time.Time) {
	if m.sess != nil {
		if m.sess.Status == session.StatusRunning {
			rem := m.sess.RemainingAt(now)
			m.sess.RemainingSec = int64(math.Ceil(rem.Seconds()))
		}
		if m.sess.IsFocus() {
			m.sess.Intention = m.state.NextAction
			m.sess.IfThen = m.state.IfThen
		}
	}
	m.state.Session = m.sess
	m.state.Config = m.cfg
	m.state.LastSaved = now
	if err := m.st.Save(m.state); err != nil {
		m.setToast(err.Error())
	}
}

func (m *Model) setToast(s string) {
	m.toast = s
	m.toastID++
}

func (m *Model) startTicking() tea.Cmd {
	if m.ticking {
		return nil
	}
	if (m.sess == nil || m.sess.Status != session.StatusRunning) && m.autoBreakAt.IsZero() {
		return nil
	}
	m.ticking = true
	return tickCmd()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func (m *Model) continueTicking() tea.Cmd {
	m.ticking = false
	return m.startTicking()
}

func (m *Model) onTick(now time.Time) tea.Cmd {
	m.ticking = false
	if m.sess != nil && m.sess.Status == session.StatusRunning {
		m.sess.Tick(now)
		if m.sess.Status == session.StatusFinished {
			m.handleFinished(now, false)
		}
	}
	if !m.autoBreakAt.IsZero() && !now.Before(m.autoBreakAt) {
		m.startBreak(now)
	}
	return m.startTicking()
}

func (m *Model) openStats(now time.Time) {
	records, err := m.st.Sessions()
	m.summary = stats.Summarize(records, now, m.parkingCount())
	m.todayFocusSec = m.summary.TodayFocusSec
	m.todayBlocks = m.summary.TodayBlocks
	m.screen = screenStats
	if err != nil {
		m.setToast(err.Error())
	}
}

func (m *Model) refreshToday(now time.Time) {
	records, err := m.st.Sessions()
	if err != nil {
		return
	}
	summary := stats.Summarize(records, now, m.parkingCount())
	m.todayFocusSec = summary.TodayFocusSec
	m.todayBlocks = summary.TodayBlocks
}

func (m *Model) parkingCount() int {
	n := 0
	for _, t := range m.state.Parking {
		if !t.Done {
			n++
		}
	}
	return n
}

func (m *Model) openIntention() tea.Cmd {
	m.nextInput.SetValue(m.state.NextAction)
	m.ifInput.SetValue(m.state.IfThen)
	m.field = 0
	m.nextInput.Focus()
	m.ifInput.Blur()
	m.modal = modalIntention
	return textinput.Blink
}

func (m *Model) saveIntention(now time.Time) tea.Cmd {
	m.state.NextAction = trimSpace(m.nextInput.Value())
	m.state.IfThen = trimSpace(m.ifInput.Value())
	if m.sess != nil && m.sess.IsFocus() {
		m.sess.Intention = m.state.NextAction
		m.sess.IfThen = m.state.IfThen
	}
	m.save(now)
	m.modal = modalNone
	m.nextInput.Blur()
	m.ifInput.Blur()
	if m.firstRun {
		m.firstRun = false
		m.state.FirstRunDone = true
		m.startFocus(now, m.focusDuration(m.activation), m.activation)
		return m.startTicking()
	}
	return m.setToastCmd("Intention saved")
}

func (m *Model) openParkingInput() tea.Cmd {
	m.parkingInput.SetValue("")
	m.parkingInput.Focus()
	m.modal = modalParkingInput
	return textinput.Blink
}

func (m *Model) addParking(now time.Time) tea.Cmd {
	text := trimSpace(m.parkingInput.Value())
	if text == "" {
		m.modal = modalNone
		m.parkingInput.Blur()
		return nil
	}
	m.state.Parking = append(m.state.Parking, store.Task{
		ID:        fmt.Sprintf("%d", now.UnixNano()),
		Text:      text,
		CreatedAt: now,
	})
	m.save(now)
	m.modal = modalNone
	m.parkingInput.Blur()
	return m.setToastCmd("Parked. Back to it.")
}

func (m *Model) setToastCmd(s string) tea.Cmd {
	m.setToast(s)
	id := m.toastID
	return tea.Tick(2200*time.Millisecond, func(time.Time) tea.Msg { return toastMsg{id} })
}

func clampInt(v, min, max int) int {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

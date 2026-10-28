package app

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/session"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		w := clampInt(msg.Width-18, 20, 60)
		m.nextInput.Width = w
		m.ifInput.Width = w
		m.parkingInput.Width = w
		return m, nil

	case tickMsg:
		return m, m.onTick(time.Time(msg))

	case toastMsg:
		if msg.id == m.toastID {
			m.toast = ""
		}
		return m, nil

	case QuitMsg:
		m.save(m.now())
		return m, tea.Quit

	case tea.KeyMsg:
		if m.modal != modalNone {
			return m.updateModal(msg)
		}
		return m.updateMain(msg)
	}

	// Bubble Tea's text input uses internal cursor-blink messages. Pass
	// anything we did not handle to whichever field is focused.
	if m.modal == modalIntention {
		var cmd tea.Cmd
		if m.field == 0 {
			m.nextInput, cmd = m.nextInput.Update(msg)
		} else {
			m.ifInput, cmd = m.ifInput.Update(msg)
		}
		return m, cmd
	}
	if m.modal == modalParkingInput {
		var cmd tea.Cmd
		m.parkingInput, cmd = m.parkingInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *Model) updateMain(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	now := m.now()

	if m.screen == screenStats {
		switch msg.String() {
		case "q", "ctrl+c":
			m.save(now)
			return m, tea.Quit
		case "esc", "v":
			m.screen = screenMain
		case "?":
			m.screen = screenHelp
		}
		return m, nil
	}
	if m.screen == screenHelp {
		switch msg.String() {
		case "q", "ctrl+c":
			m.save(now)
			return m, tea.Quit
		case "esc", "?":
			m.screen = screenMain
		case "v":
			m.openStats(now)
		}
		return m, nil
	}

	switch msg.String() {
	case "q", "ctrl+c":
		m.save(now)
		return m, tea.Quit

	case "?":
		m.screen = screenHelp
		return m, nil

	case "v":
		m.openStats(now)
		return m, nil

	case "i":
		return m, m.openIntention()

	case "d":
		return m, m.openParkingInput()

	case "p":
		m.modal = modalParkingList
		m.clampParkingSel()
		return m, nil

	case " ":
		return m, m.togglePause(now)

	case "1", "2", "3":
		return m, m.startPreset(now, msg.String())

	case "e":
		return m, m.extend(now)

	case "s":
		return m, m.skipBreak(now)

	case "b":
		if m.sess != nil && m.sess.IsFocus() && m.sess.Status == session.StatusFinished {
			m.startBreak(now)
			return m, m.startTicking()
		}
		return m, nil

	case "n":
		if m.sess == nil || m.sess.Status != session.StatusRunning {
			m.startNextFocus(now)
			return m, m.startTicking()
		}
		return m, nil

	case "enter":
		return m, m.enter(now)
	}

	return m, nil
}

func (m *Model) updateModal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	now := m.now()

	switch m.modal {
	case modalIntention:
		switch msg.String() {
		case "ctrl+c":
			m.save(now)
			return m, tea.Quit
		case "esc":
			if m.firstRun {
				return m, m.finishFirstRun(now)
			}
			m.modal = modalNone
			m.nextInput.Blur()
			m.ifInput.Blur()
			return m, nil
		case "tab", "shift+tab":
			m.field = 1 - m.field
			if m.field == 0 {
				m.nextInput.Focus()
				m.ifInput.Blur()
			} else {
				m.ifInput.Focus()
				m.nextInput.Blur()
			}
			return m, textinput.Blink
		case "enter":
			return m, m.saveIntention(now)
		}
		var cmd tea.Cmd
		if m.field == 0 {
			m.nextInput, cmd = m.nextInput.Update(msg)
		} else {
			m.ifInput, cmd = m.ifInput.Update(msg)
		}
		return m, cmd

	case modalParkingInput:
		switch msg.String() {
		case "ctrl+c":
			m.save(now)
			return m, tea.Quit
		case "esc":
			m.modal = modalNone
			m.parkingInput.Blur()
			return m, nil
		case "enter":
			return m, m.addParking(now)
		}
		var cmd tea.Cmd
		m.parkingInput, cmd = m.parkingInput.Update(msg)
		return m, cmd

	case modalParkingList:
		switch msg.String() {
		case "ctrl+c":
			m.save(now)
			return m, tea.Quit
		case "esc", "p":
			m.modal = modalNone
		case "j", "down":
			if m.parkingSel < len(m.state.Parking)-1 {
				m.parkingSel++
			}
		case "k", "up":
			if m.parkingSel > 0 {
				m.parkingSel--
			}
		case "u":
			if m.parkingSel >= 0 && m.parkingSel < len(m.state.Parking) {
				m.removeParking(m.parkingSel)
				m.save(now)
				return m, m.setToastCmd("Cleared")
			}
		case "x", "delete":
			if m.parkingSel >= 0 && m.parkingSel < len(m.state.Parking) {
				m.removeParking(m.parkingSel)
				m.save(now)
				return m, m.setToastCmd("Removed")
			}
		case "enter":
			if m.parkingSel >= 0 && m.parkingSel < len(m.state.Parking) {
				task := m.state.Parking[m.parkingSel]
				m.state.NextAction = task.Text
				if m.sess != nil && m.sess.IsFocus() {
					m.sess.Intention = task.Text
				}
				m.removeParking(m.parkingSel)
				m.save(now)
				m.modal = modalNone
				return m, m.setToastCmd("Next action set")
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) finishFirstRun(now time.Time) tea.Cmd {
	m.state.NextAction = trimSpace(m.nextInput.Value())
	m.state.IfThen = trimSpace(m.ifInput.Value())
	m.state.FirstRunDone = true
	m.firstRun = false
	m.modal = modalNone
	m.nextInput.Blur()
	m.ifInput.Blur()
	m.startFocus(now, m.focusDuration(m.activation), m.activation)
	return m.startTicking()
}

func (m *Model) togglePause(now time.Time) tea.Cmd {
	if m.sess == nil {
		return nil
	}
	switch m.sess.Status {
	case session.StatusRunning:
		m.sess.Pause(now)
		m.save(now)
		m.ticking = false
		return nil
	case session.StatusPaused:
		m.sess.Resume(now)
		m.save(now)
		m.ticking = false
		return m.startTicking()
	case session.StatusFinished:
		if m.sess.IsBreak() {
			m.startNextFocus(now)
			return m.startTicking()
		}
		m.startBreak(now)
		return m.startTicking()
	}
	return nil
}

func (m *Model) startPreset(now time.Time, key string) tea.Cmd {
	if m.sess != nil && m.sess.Status == session.StatusRunning {
		return m.setToastCmd("Pause first to switch")
	}
	d := m.cfg.FocusDuration()
	activation := false
	switch key {
	case "1":
		d = m.cfg.ActivationDuration()
		activation = true
	case "3":
		d = 2 * m.cfg.FocusDuration()
	}
	m.startFocus(now, d, activation)
	m.ticking = false
	return m.startTicking()
}

func (m *Model) extend(now time.Time) tea.Cmd {
	if m.sess == nil {
		return nil
	}
	m.sess.Extend(now, m.cfg.ExtendDuration())
	m.autoBreakAt = time.Time{}
	m.save(now)
	if m.sess.Status == session.StatusRunning {
		m.ticking = false
		return m.startTicking()
	}
	return m.setToastCmd("Added " + itoa(m.cfg.ExtendMinutes) + "m")
}

func (m *Model) skipBreak(now time.Time) tea.Cmd {
	if m.sess == nil || !m.sess.IsBreak() {
		return m.setToastCmd("s skips breaks only")
	}
	m.sess.Status = session.StatusFinished
	m.sess.RemainingSec = 0
	m.sess.Deadline = time.Time{}
	m.autoBreakAt = time.Time{}
	m.save(now)
	m.ticking = false
	return nil
}

func (m *Model) enter(now time.Time) tea.Cmd {
	if m.sess == nil {
		m.startNextFocus(now)
		return m.startTicking()
	}
	switch m.sess.Status {
	case session.StatusPaused:
		return m.togglePause(now)
	case session.StatusFinished:
		if m.sess.IsBreak() {
			m.startNextFocus(now)
			return m.startTicking()
		}
		if m.sess.Activation {
			m.startNextFocus(now)
			return m.startTicking()
		}
		m.startBreak(now)
		return m.startTicking()
	}
	return nil
}

func (m *Model) removeParking(i int) {
	if i < 0 || i >= len(m.state.Parking) {
		return
	}
	m.state.Parking = append(m.state.Parking[:i], m.state.Parking[i+1:]...)
	m.clampParkingSel()
}

func (m *Model) clampParkingSel() {
	if len(m.state.Parking) == 0 {
		m.parkingSel = 0
		return
	}
	if m.parkingSel >= len(m.state.Parking) {
		m.parkingSel = len(m.state.Parking) - 1
	}
	if m.parkingSel < 0 {
		m.parkingSel = 0
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}

package app

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/session"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/stats"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/ui"
)

func (m *Model) View() string {
	if m.width < 20 || m.height < 8 {
		return "pomo needs a slightly larger terminal"
	}

	var content string
	switch m.modal {
	case modalIntention:
		content = m.viewIntention()
	case modalParkingInput:
		content = m.viewParkingInput()
	case modalParkingList:
		content = m.viewParking()
	default:
		switch m.screen {
		case screenStats:
			content = m.viewStats()
		case screenHelp:
			content = m.viewHelp()
		default:
			content = m.viewMain()
		}
	}
	return ui.Screen(m.center(content), m.width, m.height, m.theme.BG)
}

func (m *Model) center(s string) string {
	lines := strings.Split(s, "\n")
	style := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Background(m.theme.BG)
	for i, line := range lines {
		lines[i] = style.Render(line)
	}
	return strings.Join(lines, "\n")
}

func (m *Model) viewMain() string {
	now := m.now()
	panelWidth := clampInt(m.width-10, 26, 76)
	body := m.viewBody(now)
	panel := m.theme.Panel().Width(panelWidth).Align(lipgloss.Center).Render(body)
	return m.header() + "\n\n" + panel + "\n" + m.footer()
}

func (m *Model) header() string {
	left := m.theme.Fg(m.theme.Magenta).Bold(true).Render("pomo")
	sep := m.theme.Fg(m.theme.Muted).Render("  ·  ")
	phaseLabel := "setup"
	color := m.theme.Muted
	if m.sess != nil {
		switch {
		case m.sess.Status == session.StatusPaused:
			phaseLabel, color = "paused", m.theme.Yellow
		case m.sess.Status == session.StatusFinished && m.sess.IsBreak():
			phaseLabel, color = "break done", m.theme.Green
		case m.sess.Status == session.StatusFinished:
			phaseLabel, color = "block done", m.theme.Green
		case m.sess.Activation:
			phaseLabel, color = "activation", m.theme.Blue
		case m.sess.IsBreak():
			phaseLabel, color = "break", m.theme.Green
		default:
			phaseLabel, color = "focus", m.theme.Cyan
		}
	}
	phase := m.theme.Fg(color).Render(strings.ToUpper(phaseLabel))
	right := m.theme.Fg(m.theme.Muted).Render(
		fmt.Sprintf("today %s · %d blocks", stats.FormatDuration(m.todayFocusSec), m.todayBlocks),
	)
	leftWidth := lipgloss.Width(left) + lipgloss.Width(sep) + lipgloss.Width(phase)
	gap := m.width - leftWidth - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}
	return "  " + left + sep + phase + strings.Repeat(" ", gap) + right + "  "
}

func (m *Model) viewBody(now time.Time) string {
	if m.sess == nil {
		return lipgloss.JoinVertical(lipgloss.Center,
			m.theme.Fg(m.theme.FG).Bold(true).Render("Ready."),
			"",
			m.theme.Fg(m.theme.Muted).Render("Press enter to start a focus block."),
		)
	}

	switch {
	case m.sess.IsBreak() && m.sess.Status == session.StatusFinished:
		return m.viewBreakDone()
	case m.sess.IsFocus() && m.sess.Status == session.StatusFinished:
		return m.viewComplete()
	default:
		return m.viewRunning(now)
	}
}

func (m *Model) viewRunning(now time.Time) string {
	rem := m.sess.RemainingAt(now)
	if m.sess.Status != session.StatusRunning {
		rem = m.sess.RemainingAt(now)
	}
	total := m.sess.Total()
	if total <= 0 {
		total = time.Second
	}
	color := m.clockColor(rem, total)
	clock := m.renderClock(rem, color)

	status := "FOCUS"
	switch {
	case m.sess.Status == session.StatusPaused:
		status = "PAUSED · space resumes"
	case m.sess.IsBreak():
		status = "BREAK · s skips"
	case m.sess.Activation:
		status = "ACTIVATION · small on purpose"
	}
	statusLine := m.theme.Fg(color).Bold(true).Render(status)

	endLine := ""
	if m.sess.Status == session.StatusRunning {
		endLine = m.theme.Fg(m.theme.Muted).Render("ends at " + now.Add(rem).Format("15:04"))
	} else if m.sess.Status == session.StatusPaused {
		endLine = m.theme.Fg(m.theme.Yellow).Render("paused with " + ui.ClockText(rem) + " left")
	}

	barWidth := clampInt(m.width-24, 10, 52)
	ratio := 1 - float64(rem)/float64(total)
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	bar := ui.ProgressBar(barWidth, ratio, color, m.theme.Muted)

	intention := m.intentionLine(clampInt(m.width-18, 20, 64))
	ifThen := ""
	if t := trimSpace(m.state.IfThen); t != "" {
		ifThen = m.theme.Fg(m.theme.Muted).Render(ui.Truncate("if "+t, clampInt(m.width-18, 20, 64)))
	}

	meta := m.theme.Fg(m.theme.Muted).Render(
		fmt.Sprintf("today %s · %d blocks · %d parked",
			stats.FormatDuration(m.todayFocusSec), m.todayBlocks, m.parkingCount()),
	)

	parts := []string{clock, "", statusLine, endLine, "", bar, "", intention}
	if ifThen != "" {
		parts = append(parts, ifThen)
	}
	parts = append(parts, "", meta)
	if m.toast != "" {
		parts = append(parts, "", m.theme.Fg(m.theme.Yellow).Render(m.toast))
	}
	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

func (m *Model) viewComplete() string {
	kept := m.sess.PlannedSec + m.sess.ExtraSec
	title := m.theme.Fg(m.theme.Green).Bold(true).Render("BLOCK COMPLETE")
	banked := m.theme.Fg(m.theme.FG).Bold(true).Render("+" + stats.FormatDuration(kept) + " banked")
	if m.sess.Activation {
		title = m.theme.Fg(m.theme.Blue).Bold(true).Render("MOMENTUM")
		banked = m.theme.Fg(m.theme.FG).Bold(true).Render(stats.FormatDuration(kept) + " in. Keep going?")
	}

	next := ""
	switch {
	case m.sess.Activation:
		next = "enter continue " + stats.FormatDuration(int64(m.cfg.FocusDuration()/time.Second)) + " · b break · q stop"
	case !m.autoBreakAt.IsZero():
		left := int(m.autoBreakAt.Sub(m.now()).Seconds() + 0.5)
		if left < 0 {
			left = 0
		}
		next = fmt.Sprintf("break in %ds · b break · n next · e +%dm", left, m.cfg.ExtendMinutes)
	default:
		next = fmt.Sprintf("b break · n next focus · e +%dm", m.cfg.ExtendMinutes)
	}

	intention := m.intentionLine(clampInt(m.width-18, 20, 64))
	meta := m.theme.Fg(m.theme.Muted).Render(
		fmt.Sprintf("today %s · %d blocks", stats.FormatDuration(m.todayFocusSec), m.todayBlocks),
	)

	parts := []string{title, "", banked, "", m.theme.Fg(m.theme.Muted).Render(next), "", intention, meta}
	if m.toast != "" {
		parts = append(parts, "", m.theme.Fg(m.theme.Yellow).Render(m.toast))
	}
	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

func (m *Model) viewBreakDone() string {
	title := m.theme.Fg(m.theme.Green).Bold(true).Render("BREAK DONE")
	line := m.theme.Fg(m.theme.FG).Render("Ready when you are.")
	next := m.theme.Fg(m.theme.Muted).Render(
		fmt.Sprintf("enter start %s · i change next action · q quit",
			stats.FormatDuration(int64(m.cfg.FocusDuration()/time.Second))),
	)
	parts := []string{title, "", line, "", next}
	if m.toast != "" {
		parts = append(parts, "", m.theme.Fg(m.theme.Yellow).Render(m.toast))
	}
	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

func (m *Model) renderClock(rem time.Duration, color lipgloss.Color) string {
	if m.width >= 54 && ui.BigClockWidth(rem) <= m.width-14 {
		rows := ui.BigClock(rem)
		style := m.theme.Fg(color).Bold(true)
		for i := range rows {
			rows[i] = style.Render(rows[i])
		}
		return strings.Join(rows, "\n")
	}
	return m.theme.Fg(color).Bold(true).Render(ui.ClockText(rem))
}

func (m *Model) clockColor(rem, total time.Duration) lipgloss.Color {
	if m.sess.Status == session.StatusPaused {
		return m.theme.Yellow
	}
	if m.sess.IsBreak() {
		return m.theme.Green
	}
	if rem <= time.Minute {
		return m.theme.Red
	}
	if total > 0 && rem*5 <= total {
		return m.theme.Orange
	}
	if m.sess.Activation {
		return m.theme.Blue
	}
	return m.theme.Cyan
}

func (m *Model) intentionLine(width int) string {
	if t := trimSpace(m.state.NextAction); t != "" {
		return m.theme.Fg(m.theme.FG).Bold(true).Render(ui.Truncate("▸ "+t, width))
	}
	return m.theme.Fg(m.theme.Muted).Render("press i to name your next action")
}

func (m *Model) footer() string {
	if m.sess == nil {
		return m.theme.Fg(m.theme.Muted).Render(ui.Truncate("enter start · i next action · v stats · q quit", m.width-4))
	}
	parts := ""
	switch {
	case m.sess.IsBreak() && m.sess.Status == session.StatusFinished:
		parts = "enter start focus · q quit"
	case m.sess.IsBreak():
		parts = "space pause · s skip · d dump · p lot · q quit"
	case m.sess.Status == session.StatusFinished && m.sess.Activation:
		parts = "enter continue · b break · e +10 · q stop"
	case m.sess.Status == session.StatusFinished:
		parts = "b break · n next · e +10 · q quit"
	default:
		parts = "space pause · 1 start 5m · 2 focus · 3 deep · d dump · i intent · p lot · v stats · q quit"
	}
	return m.theme.Fg(m.theme.Muted).Render(ui.Truncate(parts, m.width-4))
}

func (m *Model) viewStats() string {
	panelWidth := clampInt(m.width-10, 34, 76)
	barWidth := clampInt(m.width-40, 8, 28)

	title := m.theme.Fg(m.theme.Magenta).Bold(true).Render("last 7 days")
	today := m.theme.Fg(m.theme.FG).Render(fmt.Sprintf(
		"today  %s · %d blocks · %d parked",
		stats.FormatDuration(m.summary.TodayFocusSec), m.summary.TodayBlocks, m.summary.ParkingCount,
	))
	lines := []string{title, "", today, ""}

	max := int64(1)
	for _, d := range m.summary.Days {
		if d.FocusSec > max {
			max = d.FocusSec
		}
	}
	for _, d := range m.summary.Days {
		n := 0
		if d.FocusSec > 0 {
			n = int(float64(barWidth)*float64(d.FocusSec)/float64(max) + 0.5)
			if n < 1 {
				n = 1
			}
		}
		bar := m.theme.Fg(m.theme.Blue).Render(strings.Repeat("█", n)) +
			m.theme.Fg(m.theme.Muted).Render(strings.Repeat("░", barWidth-n))
		label := d.Date.Format("Mon 02")
		value := stats.FormatDuration(d.FocusSec)
		if d.FocusSec == 0 {
			value = "·"
		}
		lines = append(lines, fmt.Sprintf("%s  %s  %-6s", label, bar, value))
	}

	lines = append(lines, "",
		m.theme.Fg(m.theme.Muted).Render(fmt.Sprintf("%d active days · no streaks, no grades", m.summary.ActiveDays)),
	)
	if m.summary.TotalBlocks == 0 {
		lines = append(lines, "", m.theme.Fg(m.theme.Yellow).Render("No completed blocks yet. Try 5 minutes."))
	}
	lines = append(lines, "", m.theme.Fg(m.theme.Muted).Render("esc back · q quit"))

	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	panel := m.theme.Panel().Width(panelWidth).Render(body)
	return m.header() + "\n\n" + panel + "\n" + m.footer()
}

func (m *Model) viewHelp() string {
	panelWidth := clampInt(m.width-12, 34, 68)
	title := m.theme.Fg(m.theme.Magenta).Bold(true).Render("keys")
	rows := [][2]string{
		{"space", "pause / resume · no ticks while paused"},
		{"1", "activation sprint (5m)"},
		{"2", "focus (25m)"},
		{"3", "deep focus (50m)"},
		{"e", "extend by 10m"},
		{"d", "park a distraction, return immediately"},
		{"i", "next action + optional if-then plan"},
		{"p", "parking lot"},
		{"s", "skip a break (never focus)"},
		{"v", "stats"},
		{"?", "this help"},
		{"q", "quit · state is saved"},
	}
	var lines []string
	lines = append(lines, title, "")
	for _, r := range rows {
		key := m.theme.Fg(m.theme.Cyan).Bold(true).Render(fmt.Sprintf("%-6s", r[0]))
		lines = append(lines, key+"  "+m.theme.Fg(m.theme.FG).Render(r[1]))
	}
	lines = append(lines, "",
		m.theme.Fg(m.theme.Muted).Render("Evidence and caveats: docs/EVIDENCE.md"),
		m.theme.Fg(m.theme.Muted).Render("No streaks. No grades. Just what you did."),
		"", m.theme.Fg(m.theme.Muted).Render("esc back · q quit"),
	)
	panel := m.theme.Panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return m.header() + "\n\n" + panel + "\n" + m.footer()
}

func (m *Model) viewIntention() string {
	panelWidth := clampInt(m.width-12, 32, 68)
	title := "Next action"
	sub := "One physical move, not a project. Tab adds an if-then plan."
	if m.firstRun {
		title = "What's the next physical action?"
		sub = "No backlog. No planning session. Just the first move."
	}

	nextLabel := m.theme.Fg(m.theme.Muted).Render("NEXT")
	if m.field == 0 {
		nextLabel = m.theme.Fg(m.theme.Cyan).Bold(true).Render("NEXT")
	}
	ifLabel := m.theme.Fg(m.theme.Muted).Render("IF-THEN (optional)")
	if m.field == 1 {
		ifLabel = m.theme.Fg(m.theme.Magenta).Bold(true).Render("IF-THEN (optional)")
	}

	hint := "enter save · esc cancel · tab switch"
	if m.firstRun {
		hint = "enter start · esc skip · tab switch"
	}

	lines := []string{
		m.theme.Fg(m.theme.FG).Bold(true).Render(title),
		m.theme.Fg(m.theme.Muted).Render(sub),
		"",
		nextLabel,
		m.nextInput.View(),
		"",
		ifLabel,
		m.ifInput.View(),
		"",
		m.theme.Fg(m.theme.Muted).Render(hint),
	}
	panel := m.theme.Panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return panel
}

func (m *Model) viewParkingInput() string {
	panelWidth := clampInt(m.width-12, 32, 68)
	lines := []string{
		m.theme.Fg(m.theme.Magenta).Bold(true).Render("Park it"),
		m.theme.Fg(m.theme.Muted).Render("Capture the thought now. You are not acting on it."),
		"",
		m.parkingInput.View(),
		"",
		m.theme.Fg(m.theme.Muted).Render("enter save · esc cancel"),
	}
	panel := m.theme.Panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return panel
}

func (m *Model) viewParking() string {
	panelWidth := clampInt(m.width-10, 34, 72)
	title := m.theme.Fg(m.theme.Magenta).Bold(true).Render(fmt.Sprintf("parked · %d", len(m.state.Parking)))
	lines := []string{title, ""}

	if len(m.state.Parking) == 0 {
		lines = append(lines,
			m.theme.Fg(m.theme.Muted).Render("Nothing parked."),
			m.theme.Fg(m.theme.Muted).Render("Press d while focusing to capture a distraction."),
		)
	} else {
		maxRows := clampInt(m.height-14, 4, 12)
		start := 0
		if m.parkingSel >= maxRows {
			start = m.parkingSel - maxRows + 1
		}
		end := start + maxRows
		if end > len(m.state.Parking) {
			end = len(m.state.Parking)
		}
		for i := start; i < end; i++ {
			marker := "  "
			style := m.theme.Fg(m.theme.FG)
			if i == m.parkingSel {
				marker = "▸ "
				style = m.theme.Fg(m.theme.Magenta).Bold(true)
			}
			text := ui.Truncate(m.state.Parking[i].Text, panelWidth-8)
			lines = append(lines, style.Render(marker+text))
		}
	}

	lines = append(lines, "",
		m.theme.Fg(m.theme.Muted).Render("u clear · x delete · enter make next action"),
		m.theme.Fg(m.theme.Muted).Render("esc back · review at the break, not now"),
	)
	panel := m.theme.Panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return panel
}

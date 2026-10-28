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
	panel := m.panel().Width(panelWidth).Align(lipgloss.Center).Render(body)
	return m.header() + "\n\n" + panel + "\n" + m.footer()
}

func (m *Model) compact() bool {
	return m.height < 26
}

func (m *Model) narrow() bool {
	return m.width < 52
}

func (m *Model) panel() lipgloss.Style {
	style := m.theme.Panel()
	if m.compact() {
		style = style.Padding(0, 2)
	}
	return style
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
	if m.narrow() {
		return "  " + left + sep + phase + "  "
	}
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
	if m.narrow() {
		switch {
		case m.sess.Status == session.StatusPaused:
			status = "PAUSED"
		case m.sess.IsBreak():
			status = "BREAK"
		case m.sess.Activation:
			status = "ACTIVATION"
		default:
			status = "FOCUS"
		}
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

	metaText := fmt.Sprintf("today %s · %d blocks", stats.FormatDuration(m.todayFocusSec), m.todayBlocks)
	if !m.narrow() {
		metaText += fmt.Sprintf(" · %d parked", m.parkingCount())
	}
	meta := m.theme.Fg(m.theme.Muted).Render(metaText)

	parts := []string{clock}
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, statusLine)
	if endLine != "" {
		parts = append(parts, endLine)
	}
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, bar)
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, intention)
	if ifThen != "" {
		parts = append(parts, ifThen)
	}
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, meta)
	if m.toast != "" {
		if !m.compact() {
			parts = append(parts, "")
		}
		parts = append(parts, m.theme.Fg(m.theme.Yellow).Render(m.toast))
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
		if m.narrow() {
			next = "enter continue · b break"
		} else {
			next = "enter continue " + stats.FormatDuration(int64(m.cfg.FocusDuration()/time.Second)) + " · b break · q stop"
		}
	case !m.autoBreakAt.IsZero():
		left := int(m.autoBreakAt.Sub(m.now()).Seconds() + 0.5)
		if left < 0 {
			left = 0
		}
		if m.narrow() {
			next = fmt.Sprintf("break in %ds · b · n · e", left)
		} else {
			next = fmt.Sprintf("break in %ds · b break · n next · e +%dm", left, m.cfg.ExtendMinutes)
		}
	default:
		if m.narrow() {
			next = fmt.Sprintf("b · n next · e +%dm", m.cfg.ExtendMinutes)
		} else {
			next = fmt.Sprintf("b break · n next focus · e +%dm", m.cfg.ExtendMinutes)
		}
	}

	intention := m.intentionLine(clampInt(m.width-18, 20, 64))
	meta := m.theme.Fg(m.theme.Muted).Render(
		fmt.Sprintf("today %s · %d blocks", stats.FormatDuration(m.todayFocusSec), m.todayBlocks),
	)
	if m.narrow() {
		meta = m.theme.Fg(m.theme.Muted).Render("today " + stats.FormatDuration(m.todayFocusSec))
	}

	parts := []string{title}
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, banked)
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, m.theme.Fg(m.theme.Muted).Render(next))
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, intention, meta)
	if m.toast != "" {
		if !m.compact() {
			parts = append(parts, "")
		}
		parts = append(parts, m.theme.Fg(m.theme.Yellow).Render(m.toast))
	}
	return lipgloss.JoinVertical(lipgloss.Center, parts...)
}

func (m *Model) viewBreakDone() string {
	title := m.theme.Fg(m.theme.Green).Bold(true).Render("BREAK DONE")
	line := m.theme.Fg(m.theme.FG).Render("Ready when you are.")
	nextText := fmt.Sprintf("enter start %s · i change next action · q quit",
		stats.FormatDuration(int64(m.cfg.FocusDuration()/time.Second)))
	if m.narrow() {
		nextText = "enter start · q quit"
	}
	next := m.theme.Fg(m.theme.Muted).Render(nextText)
	if !m.autoFocusAt.IsZero() {
		left := int(m.autoFocusAt.Sub(m.now()).Seconds() + 0.5)
		if left < 0 {
			left = 0
		}
		if m.narrow() {
			next = m.theme.Fg(m.theme.Muted).Render(fmt.Sprintf("focus in %ds · enter · q", left))
		} else {
			next = m.theme.Fg(m.theme.Muted).Render(
				fmt.Sprintf("starting focus in %ds · enter now · q quit", left),
			)
		}
	}
	parts := []string{title}
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, line)
	if !m.compact() {
		parts = append(parts, "")
	}
	parts = append(parts, next)
	if m.toast != "" {
		if !m.compact() {
			parts = append(parts, "")
		}
		parts = append(parts, m.theme.Fg(m.theme.Yellow).Render(m.toast))
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
		parts = fmt.Sprintf("enter continue · b break · e +%dm · q stop", m.cfg.ExtendMinutes)
	case m.sess.Status == session.StatusFinished:
		parts = "b break · n next · e +10 · q quit"
	case m.sess.Status == session.StatusPaused:
		parts = "space resume · e +10 · d dump · i intent · v stats · q quit"
	default:
		parts = "space pause · 1 start 5m · 2 focus · 3 deep · d dump · i intent · p lot · v stats · q quit"
	}
	return m.staticFooter(parts)
}

func (m *Model) staticFooter(parts string) string {
	return m.theme.Fg(m.theme.Muted).Render(ui.Truncate(parts, m.width-4))
}

func (m *Model) viewStats() string {
	panelWidth := clampInt(m.width-10, 34, 76)
	barWidth := clampInt(m.width-40, 8, 28)

	title := m.theme.Fg(m.theme.Magenta).Bold(true).Render("last 7 days")
	todayText := fmt.Sprintf("today  %s · %d blocks · %d parked",
		stats.FormatDuration(m.summary.TodayFocusSec), m.summary.TodayBlocks, m.summary.ParkingCount)
	if m.narrow() {
		todayText = fmt.Sprintf("today  %s · %d blocks",
			stats.FormatDuration(m.summary.TodayFocusSec), m.summary.TodayBlocks)
	}
	today := m.theme.Fg(m.theme.FG).Render(todayText)
	lines := []string{title}
	if !m.compact() {
		lines = append(lines, "")
	}
	lines = append(lines, today)
	if !m.compact() {
		lines = append(lines, "")
	}

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

	activeText := fmt.Sprintf("%d active days · no streaks, no grades", m.summary.ActiveDays)
	if m.narrow() {
		activeText = fmt.Sprintf("%d active days", m.summary.ActiveDays)
	}
	lines = append(lines, m.theme.Fg(m.theme.Muted).Render(activeText))
	if m.summary.TotalBlocks == 0 && !m.compact() {
		lines = append(lines, "", m.theme.Fg(m.theme.Yellow).Render("No completed blocks yet. Try 5 minutes."))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, lines...)
	panel := m.panel().Width(panelWidth).Render(body)
	return m.header() + "\n\n" + panel + "\n" + m.staticFooter("esc back · q quit")
}

func (m *Model) viewHelp() string {
	panelWidth := clampInt(m.width-12, 34, 68)
	title := m.theme.Fg(m.theme.Magenta).Bold(true).Render("keys")
	rows := [][2]string{
		{"space", "pause / resume"},
		{"1", "activation 5m"},
		{"2", "focus 25m"},
		{"3", "deep 50m"},
		{"e", "extend 10m"},
		{"d", "park thought"},
		{"i", "next action"},
		{"p", "parking lot"},
		{"s", "skip break"},
		{"v", "stats"},
		{"?", "help"},
		{"q", "quit + save"},
	}
	var lines []string
	lines = append(lines, title)
	if !m.compact() {
		lines = append(lines, "")
	}
	if m.height < 18 && len(rows) > 10 {
		rows = rows[:10]
	}
	for _, r := range rows {
		key := m.theme.Fg(m.theme.Cyan).Bold(true).Render(fmt.Sprintf("%-6s", r[0]))
		lines = append(lines, key+"  "+m.theme.Fg(m.theme.FG).Render(r[1]))
	}
	if !m.compact() {
		lines = append(lines, "",
			m.theme.Fg(m.theme.Muted).Render("Evidence and caveats: docs/EVIDENCE.md"),
			m.theme.Fg(m.theme.Muted).Render("No streaks. No grades. Just what you did."),
		)
	}
	panel := m.panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	if m.height < 18 {
		return panel + "\n" + m.staticFooter("esc back · q quit")
	}
	return m.header() + "\n\n" + panel + "\n" + m.staticFooter("esc back · q quit")
}

func (m *Model) viewIntention() string {
	panelWidth := clampInt(m.width-12, 32, 68)
	title := "Next action"
	sub := "One physical move, not a project. Tab adds an if-then plan."
	if m.firstRun {
		title = "What's the next physical action?"
		sub = "No backlog. No planning session. Just the first move."
	}
	if m.narrow() {
		title = "Next action"
		sub = "Just the first move."
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
	if m.narrow() {
		hint = "enter save · esc cancel"
		if m.firstRun {
			hint = "enter start · esc skip"
		}
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
	panel := m.panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return panel
}

func (m *Model) viewParkingInput() string {
	panelWidth := clampInt(m.width-12, 32, 68)
	sub := "Capture the thought now. You are not acting on it."
	if m.narrow() {
		sub = "Capture it. Don't act on it."
	}
	lines := []string{
		m.theme.Fg(m.theme.Magenta).Bold(true).Render("Park it"),
		m.theme.Fg(m.theme.Muted).Render(sub),
		"",
		m.parkingInput.View(),
		"",
		m.theme.Fg(m.theme.Muted).Render("enter save · esc cancel"),
	}
	panel := m.panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return panel
}

func (m *Model) viewParking() string {
	panelWidth := clampInt(m.width-10, 34, 72)
	title := m.theme.Fg(m.theme.Magenta).Bold(true).Render(fmt.Sprintf("parked · %d", len(m.state.Parking)))
	lines := []string{title, ""}

	if len(m.state.Parking) == 0 {
		empty := "Press d while focusing to capture a distraction."
		if m.narrow() {
			empty = "Press d to capture."
		}
		lines = append(lines,
			m.theme.Fg(m.theme.Muted).Render("Nothing parked."),
			m.theme.Fg(m.theme.Muted).Render(empty),
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

	lines = append(lines, "")
	if m.narrow() {
		lines = append(lines,
			m.theme.Fg(m.theme.Muted).Render("u clear · x delete"),
			m.theme.Fg(m.theme.Muted).Render("enter promote · esc back"),
		)
	} else {
		lines = append(lines,
			m.theme.Fg(m.theme.Muted).Render("u clear · x delete · enter make next action"),
			m.theme.Fg(m.theme.Muted).Render("esc back · review at the break, not now"),
		)
	}
	panel := m.panel().Width(panelWidth).Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	return panel
}

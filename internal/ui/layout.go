package ui

import (
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Screen paints content onto a fixed-size background and centers it
// vertically. Horizontal centering is expected to be done by the caller.
func Screen(content string, width, height int, bg lipgloss.Color) string {
	if width <= 0 || height <= 0 {
		return content
	}
	lines := strings.Split(content, "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	blank := lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", width))
	top := (height - len(lines)) / 2
	if top < 0 {
		top = 0
	}
	out := make([]string, height)
	for i := range out {
		if i < top || i >= top+len(lines) {
			out[i] = blank
			continue
		}
		out[i] = paintLine(lines[i-top], width, bg)
	}
	return strings.Join(out, "\n")
}

func paintLine(line string, width int, bg lipgloss.Color) string {
	w := lipgloss.Width(line)
	if w == 0 {
		return lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", width))
	}
	if w > width {
		// Keep ANSI sequences intact while clamping visible width.
		line = truncateVisible(line, width)
		w = lipgloss.Width(line)
	}
	if w < width {
		line += lipgloss.NewStyle().Background(bg).Render(strings.Repeat(" ", width-w))
	}
	return line
}

// ProgressBar draws a solid elapsed bar and a muted remainder.
func ProgressBar(width int, ratio float64, filled, empty lipgloss.Color) string {
	if width < 1 {
		return ""
	}
	if ratio < 0 {
		ratio = 0
	}
	if ratio > 1 {
		ratio = 1
	}
	n := int(math.Round(ratio * float64(width)))
	if n < 0 {
		n = 0
	}
	if n > width {
		n = width
	}
	f := lipgloss.NewStyle().Foreground(filled).Render(strings.Repeat("█", n))
	e := lipgloss.NewStyle().Foreground(empty).Render(strings.Repeat("░", width-n))
	return f + e
}

func Truncate(s string, width int) string {
	if width <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= width {
		return s
	}
	if width <= 1 {
		return "…"
	}
	return truncateVisible(s, width-1) + "…"
}

// truncateVisible respects ANSI escape sequences by leaning on runes and
// lipgloss.Width in small steps. The strings used here are short, so a linear
// walk is fine.
func truncateVisible(s string, width int) string {
	if width <= 0 {
		return ""
	}
	var b strings.Builder
	visible := 0
	inEscape := false
	truncated := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
		}
		if inEscape {
			b.WriteRune(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		w := lipgloss.Width(string(r))
		if visible+w > width {
			truncated = true
			break
		}
		b.WriteRune(r)
		visible += w
	}
	if truncated {
		b.WriteString("\x1b[0m")
	}
	return b.String()
}

func Center(width int, s string) string {
	if width <= 0 {
		return s
	}
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(s)
}

func Pad(width int, s string) string {
	if width <= 0 {
		return s
	}
	return lipgloss.NewStyle().Width(width).Render(s)
}

// Package ui contains the Tokyo Night palette, block clock, and small layout
// primitives. It has no knowledge of timers or storage.
package ui

import "github.com/charmbracelet/lipgloss"

// Tokyo Night (night variant) by Enkia. Fixed dark values: the app is designed
// for a dark terminal and does not try to guess a light theme.
type Theme struct {
	BG      lipgloss.Color
	BGDark  lipgloss.Color
	FG      lipgloss.Color
	Muted   lipgloss.Color
	Blue    lipgloss.Color
	Cyan    lipgloss.Color
	Green   lipgloss.Color
	Yellow  lipgloss.Color
	Orange  lipgloss.Color
	Red     lipgloss.Color
	Magenta lipgloss.Color
}

func TokyoNight() Theme {
	return Theme{
		BG:      lipgloss.Color("#1a1b26"),
		BGDark:  lipgloss.Color("#16161e"),
		FG:      lipgloss.Color("#c0caf5"),
		Muted:   lipgloss.Color("#565f89"),
		Blue:    lipgloss.Color("#7aa2f7"),
		Cyan:    lipgloss.Color("#7dcfff"),
		Green:   lipgloss.Color("#9ece6a"),
		Yellow:  lipgloss.Color("#e0af68"),
		Orange:  lipgloss.Color("#ff9e64"),
		Red:     lipgloss.Color("#f7768e"),
		Magenta: lipgloss.Color("#bb9af7"),
	}
}

func (t Theme) Base() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.FG).Background(t.BG)
}

func (t Theme) Fg(c lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(c).Background(t.BG)
}

func (t Theme) Panel() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.FG).
		Background(t.BG).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Muted).
		Padding(1, 2)
}

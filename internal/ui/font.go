package ui

import (
	"fmt"
	"strings"
	"time"
)

// A hand-rolled 5x5 block font. No Nerd Fonts, no external assets. Each glyph
// is five columns wide; the colon is one column.
var glyphs = map[rune][5]string{
	'0': {"█████", "█   █", "█   █", "█   █", "█████"},
	'1': {"   ██", "    █", "    █", "    █", "   ███"},
	'2': {"█████", "    █", "█████", "█    ", "█████"},
	'3': {"█████", "    █", "█████", "    █", "█████"},
	'4': {"█   █", "█   █", "█████", "    █", "    █"},
	'5': {"█████", "█    ", "█████", "    █", "█████"},
	'6': {"█████", "█    ", "█████", "█   █", "█████"},
	'7': {"█████", "    █", "   █ ", "  █  ", "  █  "},
	'8': {"█████", "█   █", "█████", "█   █", "█████"},
	'9': {"█████", "█   █", "█████", "    █", "█████"},
	':': {" ", "█", " ", "█", " "},
	' ': {"  ", "  ", "  ", "  ", "  "},
}

// ClockText renders a duration as MM:SS (or H:MM:SS above an hour).
func ClockText(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	total := int(d.Round(time.Second) / time.Second)
	h := total / 3600
	m := (total % 3600) / 60
	s := total % 60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// BigClock returns five rows of block digits. Callers should use it only when
// the terminal is wide enough; ClockText is the narrow fallback.
func BigClock(d time.Duration) []string {
	text := ClockText(d)
	rows := make([]string, 5)
	for _, r := range text {
		g, ok := glyphs[r]
		if !ok {
			continue
		}
		for i := 0; i < 5; i++ {
			if rows[i] != "" {
				rows[i] += " "
			}
			rows[i] += g[i]
		}
	}
	return rows
}

func BigClockWidth(d time.Duration) int {
	rows := BigClock(d)
	if len(rows) == 0 {
		return 0
	}
	return len([]rune(rows[0]))
}

// JoinBigClock keeps the rows aligned as one block, optionally styled.
func JoinBigClock(rows []string) string {
	return strings.Join(rows, "\n")
}

// Package config holds the small, user-editable settings for pomo.
package config

import (
	"time"
)

const (
	DefaultFocusMinutes      = 25
	DefaultShortBreakMinutes = 5
	DefaultLongBreakMinutes  = 15
	DefaultLongBreakEvery    = 4
	DefaultActivationMinutes = 5
	DefaultExtendMinutes     = 10
	DefaultTheme             = "tokyonight"
)

// Config is persisted inside state.json so there is no separate config file
// to find, parse, or forget. Values are deliberately plain minutes.
type Config struct {
	FocusMinutes      int    `json:"focus_minutes"`
	ShortBreakMinutes int    `json:"short_break_minutes"`
	LongBreakMinutes  int    `json:"long_break_minutes"`
	LongBreakEvery    int    `json:"long_break_every"`
	ActivationMinutes int    `json:"activation_minutes"`
	ExtendMinutes     int    `json:"extend_minutes"`
	AutoStartBreak    bool   `json:"auto_start_break"`
	AutoStartFocus    bool   `json:"auto_start_focus"`
	Bell              bool   `json:"bell"`
	DesktopNotify     bool   `json:"desktop_notify"`
	Theme             string `json:"theme"`
}

// Default returns the deliberately unsurprising starting point. The 25/5
// split is familiar, not an evidence-based optimum; every value is editable.
func Default() Config {
	return Config{
		FocusMinutes:      DefaultFocusMinutes,
		ShortBreakMinutes: DefaultShortBreakMinutes,
		LongBreakMinutes:  DefaultLongBreakMinutes,
		LongBreakEvery:    DefaultLongBreakEvery,
		ActivationMinutes: DefaultActivationMinutes,
		ExtendMinutes:     DefaultExtendMinutes,
		AutoStartBreak:    true,
		AutoStartFocus:    false,
		Bell:              true,
		DesktopNotify:     true,
		Theme:             DefaultTheme,
	}
}

// Normalize replaces zero/absurd values with safe defaults. It also keeps
// durations inside a range a terminal clock can display cleanly.
func (c *Config) Normalize() {
	c.FocusMinutes = clamp(c.FocusMinutes, 1, 180, DefaultFocusMinutes)
	c.ShortBreakMinutes = clamp(c.ShortBreakMinutes, 1, 60, DefaultShortBreakMinutes)
	c.LongBreakMinutes = clamp(c.LongBreakMinutes, 1, 120, DefaultLongBreakMinutes)
	c.LongBreakEvery = clamp(c.LongBreakEvery, 1, 12, DefaultLongBreakEvery)
	c.ActivationMinutes = clamp(c.ActivationMinutes, 1, 30, DefaultActivationMinutes)
	c.ExtendMinutes = clamp(c.ExtendMinutes, 1, 60, DefaultExtendMinutes)
	if c.Theme == "" {
		c.Theme = DefaultTheme
	}
}

func (c Config) FocusDuration() time.Duration {
	return time.Duration(c.FocusMinutes) * time.Minute
}

func (c Config) ShortBreakDuration() time.Duration {
	return time.Duration(c.ShortBreakMinutes) * time.Minute
}

func (c Config) LongBreakDuration() time.Duration {
	return time.Duration(c.LongBreakMinutes) * time.Minute
}

func (c Config) ActivationDuration() time.Duration {
	return time.Duration(c.ActivationMinutes) * time.Minute
}

func (c Config) ExtendDuration() time.Duration {
	return time.Duration(c.ExtendMinutes) * time.Minute
}

func clamp(v, min, max, fallback int) int {
	if v == 0 {
		return fallback
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

package config

import "testing"

func TestNormalizeFillsAndClamps(t *testing.T) {
	c := Config{}
	c.Normalize()

	if c.FocusMinutes != DefaultFocusMinutes {
		t.Fatalf("focus = %d, want %d", c.FocusMinutes, DefaultFocusMinutes)
	}
	if c.Theme != DefaultTheme {
		t.Fatalf("theme = %q, want %q", c.Theme, DefaultTheme)
	}

	c.FocusMinutes = 999
	c.ShortBreakMinutes = -3
	c.LongBreakEvery = 0
	c.Normalize()

	if c.FocusMinutes != 180 {
		t.Fatalf("focus clamp = %d, want 180", c.FocusMinutes)
	}
	if c.ShortBreakMinutes != 1 {
		t.Fatalf("break clamp = %d, want 1", c.ShortBreakMinutes)
	}
	if c.LongBreakEvery != DefaultLongBreakEvery {
		t.Fatalf("long break every = %d, want %d", c.LongBreakEvery, DefaultLongBreakEvery)
	}
}

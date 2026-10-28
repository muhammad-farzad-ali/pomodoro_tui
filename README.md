# pomo

A small, fast terminal focus timer built around the parts of ADHD support that
have actual evidence behind them: external time, a single next action, graded
activation, immediate feedback, and a place to park distractions.

It is **not** a task manager, a streak game, or a body-doubling platform. See
[`docs/EVIDENCE.md`](docs/EVIDENCE.md) for the research basis and the caveats.

## Status

Working core: deadline-based timer, 5-minute activation sprint, pause/extend,
parking lot, resume-after-quit, passive stats, and a single boundary alert.
Hardening and packaging remain.

## Why these mechanics

The short version, with the caveats in [`docs/EVIDENCE.md`](docs/EVIDENCE.md):

- **External time** (countdown, bar, “ends at”) compensates for documented
  timing and prospective-memory difficulties.
- **One next action** and a **5-minute activation sprint** borrow from
  behavioral activation and adult ADHD organization/planning CBT.
- **Immediate feedback** (`+25m banked`) fits the steeper delay discounting
  seen in adults with ADHD.
- **Parking distractions** is cognitive offloading, not a promise to cure
  distraction.
- **Breaks are optional** because the fixed-Pomodoro evidence is mixed;
  micro-breaks help fatigue more reliably than they help performance.
- **No streaks, points, or grades.** There is no adult-ADHD evidence that
  streak loss helps, and it adds the task aversiveness we are trying to reduce.

This is not medical advice and not a substitute for medication or therapy.

## Build

```bash
sudo apt install golang-go   # Go 1.26+
git clone git@github.com:muhammad-farzad-ali/pomodoro_tui.git
cd pomodoro_tui
make build
./pomo
```

Install to `$GOBIN`:

```bash
make install
# or
go install github.com/muhammad-farzad-ali/pomodoro_tui/cmd/pomo@latest
```

## Use

```bash
pomo                    # start immediately, resuming a recent session
pomo --activation       # open with a 5-minute activation sprint
pomo --duration 50      # focus length for this run
pomo --break 10         # short-break length for this run
pomo --stats            # plain-text summary, no TUI
pomo --stats --json     # same, machine-readable
pomo --no-notify        # silence bell and desktop notification
```

On first run you are asked for one next action. After that, `pomo` starts
immediately and keeps the last next action; press `i` to change it.

## Keys

| Key | Action |
|---|---|
| `space` | pause / resume (no CPU ticks while paused) |
| `1` | 5-minute activation sprint |
| `2` | 25-minute focus |
| `3` | 50-minute deep focus |
| `e` | extend by 10 minutes |
| `d` | park a distraction and return immediately |
| `i` | next action, plus optional `if …, then …` plan |
| `p` | parking lot |
| `s` | skip a break (never focus) |
| `v` | stats |
| `?` | help |
| `q` | quit and save |

## Data

`$XDG_DATA_HOME/pomodoro-tui/state.json` and `sessions.jsonl`
(fallback: `~/.local/share/pomodoro-tui/`). Override with `POMO_DATA_DIR`.
Files are written atomically; a corrupt state file is preserved as
`state.json.corrupt-<timestamp>`.

Config lives inside `state.json`. There is no settings screen in v1; edit the
JSON while pomo is not running, or use the CLI flags. The CLI flags apply to
the current run and do not overwrite the saved config.

```json
"config": {
  "focus_minutes": 25,
  "short_break_minutes": 5,
  "long_break_minutes": 15,
  "long_break_every": 4,
  "activation_minutes": 5,
  "extend_minutes": 10,
  "auto_start_break": true,
  "auto_start_focus": false,
  "bell": true,
  "desktop_notify": true,
  "theme": "tokyonight"
}
```

`auto_start_break` starts the break after the 8-second completion screen.
`auto_start_focus` starts the next focus 5 seconds after a break ends; it is off
by default because the transition back to work is the hard one.

## Design budget

- single static binary, no CGO, no daemon, no network
- one tick per second only while running
- target: <50 ms startup, <20 MB RSS, <10 MB stripped binary
- single-instance lock prevents two terminals from clobbering state

## Development

```bash
make test   # go test ./...
make vet    # go vet ./...
make fmt
make perf   # build + show size/startup
```

## License

MIT. Not medical advice; not a substitute for treatment.

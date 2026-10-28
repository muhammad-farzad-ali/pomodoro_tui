# pomo

A small, fast terminal focus timer built around the parts of ADHD support that
have actual evidence behind them: external time, a single next action, graded
activation, immediate feedback, and a place to park distractions.

It is **not** a task manager, a streak game, or a body-doubling platform. See
[`docs/EVIDENCE.md`](docs/EVIDENCE.md) for the research basis and the caveats.

## Status

Early build. The core loop (timer, activation sprint, parking lot, stats) works
and is being hardened.

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
JSON while pomo is not running, or use the CLI flags.

## Design budget

- single static binary, no CGO, no daemon, no network
- one tick per second only while running
- target: <50 ms startup, <20 MB RSS, <10 MB stripped binary
- single-instance lock prevents two terminals from clobbering state

## License

MIT. Not medical advice; not a substitute for treatment.

# pomo

A small, fast terminal focus timer built around the parts of ADHD support that
have actual evidence behind them: external time, a single next action, graded
activation, immediate feedback, and a place to park distractions.

It is **not** a task manager, a streak game, or a body-doubling platform. See
[`docs/EVIDENCE.md`](docs/EVIDENCE.md) for the research basis and the caveats.

## Status

Working core: deadline-based timer, 5-minute activation sprint, pause/extend,
parking lot, resume-after-quit, passive stats, and a single boundary alert.
Single binary; no settings screen yet.

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

```text
  pomo  ·  FOCUS                                           today 25m · 1 block

    ╭──────────────────────────────────────────────────────────────────────╮
    │                       █████ █████   █████ █████                      │
    │                           █ █     █ █   █ █   █                      │
    │                       █████ █████   █   █ █   █                      │
    │                       █         █ █ █   █ █   █                      │
    │                       █████ █████   █████ █████                      │
    │                                 FOCUS                                │
    │                             ends at 15:24                            │
    │         ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░         │
    │                     ▸ draft the evidence section                     │
    │                  ↳ if I check email, then I park it                  │
    │                    today 25m · 1 block · 0 parked                    │
    ╰──────────────────────────────────────────────────────────────────────╯
  space pause · 1 start 5m · 2 focus · 3 deep · d dump · i intent · p lot · v…
```

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

## Walkthrough: a real session

**The situation:** a report is due, and every few minutes you want to check
email.

### Start

```bash
./pomo
```

First run only, the app asks:

```text
What's the next physical action?
› Open report.docx and write the first three bullets
```

Press **Enter**. The 25-minute focus block starts.

That phrasing matters: “write the report” is a project; “open the file and
write three bullets” is a physical action you can start.

### While working

The screen shows the countdown, a shrinking bar, the end time, your next
action, and today’s total. You do not have to remember how much time is left.

Ten minutes in, the email urge hits. Press **`d`**:

```text
Park it
› check email
```

Press **Enter**. The timer never stops. A brief “Parked. Back to it.” appears,
and the parking count becomes 1. You did not open email, and the thought is
not lost.

### When someone interrupts

Press **`space`** to pause. The clock turns yellow and no timer ticks happen.
After the call, press **`space`** again; it resumes with the same remaining
time and recalculates the end time.

If the interruption is long, press **`q`**. It saves and exits. Running
`./pomo` again within two hours resumes the session; a paused session stays
paused until you press space.

### When you are in flow

Press **`e`**. It adds 10 minutes, moves the end time, and recalculates the
bar. If the block had already finished, `e` starts a **new 10-minute block**
rather than silently inflating the completed one. The stats stay honest.

### When time is up

```text
BLOCK COMPLETE
+25m banked
break in 7s · b break · n next · e +10m
```

- **Do nothing:** the 5-minute break starts after 8 seconds.
- **`b`:** start the break now.
- **`n`:** skip the break and start the next focus block.
- **`e`:** take another 10 minutes.

There is one bell and one optional desktop notification. No repeated “are you
working?” pings.

### At the break

Stand up, drink water, look outside. The app just holds the break; it does not
pretend a five-minute break is a productivity trick.

Press **`p`** to review the parking lot:

```text
parked · 1
▸ check email

u clear · x delete · enter make next action
```

- **`u`** clears it (handled, or no longer relevant).
- **`x`** deletes it.
- **Enter** promotes it to the next action, if it is genuinely next.
- **`s`** skips the rest of the break.

At the end of the break, press **Enter** to start the next 25 minutes.
Auto-start focus is off by default because the transition back to work is the
hard part.

### End of the day

Press **`q`**. It saves and exits.

The next morning, `./pomo` starts immediately with the last next action; press
**`i`** to change it. If the previous session was still running less than two
hours ago, it resumes instead of restarting.

### When starting feels impossible

Do not start with 25 minutes. Press **`1`**, or launch with:

```bash
./pomo --activation
```

That gives you a 5-minute activation sprint. When it ends:

- **Enter** → continue into a full focus block
- **`b`** → take a break
- **`q`** → stop

The point is to lower the start threshold, not to finish the task in five
minutes.

### Check the week

Inside the app, press **`v`**. Or without the TUI:

```bash
./pomo --stats
./pomo --stats --json
```

You get today’s focus minutes, 7-day bars, and active days. No streaks, no
grades, no “you failed” state.

### A few practical rules

- **One next action, not a backlog.** Use another tool for project plans; this
  one is for the next move.
- **Edit config while pomo is closed.** The running app saves state over
  `state.json`.
- **Two terminals are blocked by design.** A second instance prints “another
  pomo is already running” rather than corrupting state.
- **Silence alerts** with `--no-notify`; the timer still works.
- **The 25/5 split is familiar, not sacred.** Change it in `state.json`, or
  use `--duration 50 --break 10`.

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
- one tick per second only while running; renderer capped at 10 FPS
- 24 MiB soft memory limit; measured ~21 MB RSS and 3.9 MB stripped binary
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

# pomo — build todos

Evidence-grounded ADHD focus timer. See `docs/EVIDENCE.md` for the research basis and caveats.

## M0 — Toolchain and skeleton
- [x] Install Go 1.26
- [x] Initialize module and pin TUI dependencies
- [ ] Tokyo Night theme scaffold
- [ ] Project docs, license, Makefile, `.gitignore`
- [ ] Initial commit

## M1 — Timer shell
- [ ] Deadline-based countdown (no drift, survives sleep)
- [ ] Big clock, progress bar, "ends at" time
- [ ] Presets, pause/resume, extend
- [ ] Responsive layout and footer help
- [ ] Commit

## M2 — State and persistence
- [ ] Phase/status state machine
- [ ] Atomic `state.json` save/load, resume
- [ ] `sessions.jsonl` history
- [ ] Single-instance lock
- [ ] SIGINT/SIGTERM save
- [ ] Tests
- [ ] Commit

## M3 — Activation loop
- [ ] First-run next-action prompt, auto-start afterwards
- [ ] 5-minute activation sprint with continue option
- [ ] Optional if-then coping plan
- [ ] Immediate completion feedback
- [ ] Break flow, one boundary alert
- [ ] Commit

## M4 — Offloading and stats
- [ ] Parking lot capture, review, promote
- [ ] Passive 7-day stats
- [ ] `--stats` and `--stats --json`
- [ ] Commit

## M5 — Rigor
- [ ] `docs/EVIDENCE.md` with grades and references
- [ ] README, keybindings, screenshots
- [ ] Unit tests, `go vet`, perf check
- [ ] Release build and final commit

# pomo — build todos

Evidence-grounded ADHD focus timer. See `docs/EVIDENCE.md` for the research basis and caveats.

## M0 — Toolchain and skeleton
- [x] Install Go 1.26
- [x] Initialize module and pin TUI dependencies
- [x] Tokyo Night theme scaffold
- [x] Project docs, license, Makefile, `.gitignore`
- [x] Initial commit

## M1 — Timer shell
- [x] Deadline-based countdown (no drift, survives sleep)
- [x] Big clock, progress bar, "ends at" time
- [x] Presets, pause/resume, extend
- [x] Responsive layout and footer help
- [x] Commit

## M2 — State and persistence
- [x] Phase/status state machine
- [x] Atomic `state.json` save/load, resume
- [x] `sessions.jsonl` history
- [x] Single-instance lock
- [x] SIGINT/SIGTERM save
- [x] Tests
- [x] Commit

## M3 — Activation loop
- [x] First-run next-action prompt, auto-start afterwards
- [x] 5-minute activation sprint with continue option
- [x] Optional if-then coping plan
- [x] Immediate completion feedback
- [x] Break flow, one boundary alert
- [x] Commit

## M4 — Offloading and stats
- [x] Parking lot capture, review, promote
- [x] Passive 7-day stats
- [x] `--stats` and `--stats --json`
- [x] Commit

## M5 — Rigor
- [x] `docs/EVIDENCE.md` with grades and references
- [x] README, keybindings, config, caveats
- [x] Unit tests, `go vet`
- [ ] Performance check and final polish
- [ ] Final commit

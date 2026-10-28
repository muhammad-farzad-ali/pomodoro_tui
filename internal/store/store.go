// Package store persists state atomically and keeps an append-only session
// history. It intentionally uses stdlib JSON rather than a database: one user,
// one process, a few kilobytes, and no CGO.
package store

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/config"
	"github.com/muhammad-farzad-ali/pomodoro_tui/internal/session"
)

const Version = 1

var ErrLocked = errors.New("another pomo is already running")

type Task struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	Done      bool      `json:"done,omitempty"`
}

type State struct {
	Version      int              `json:"version"`
	Config       config.Config    `json:"config"`
	FirstRunDone bool             `json:"first_run_done"`
	Session      *session.Session `json:"session,omitempty"`
	Parking      []Task           `json:"parking,omitempty"`
	NextAction   string           `json:"next_action,omitempty"`
	IfThen       string           `json:"if_then,omitempty"`
	Cycle        int              `json:"cycle"`
	LastSaved    time.Time        `json:"last_saved"`
}

// Record is one completed focus block. Breaks are not logged.
type Record struct {
	StartedAt  time.Time `json:"started_at"`
	EndedAt    time.Time `json:"ended_at"`
	PlannedSec int64     `json:"planned_sec"`
	FocusedSec int64     `json:"focused_sec"`
	ExtraSec   int64     `json:"extra_sec,omitempty"`
	Activation bool      `json:"activation,omitempty"`
	Intention  string    `json:"intention,omitempty"`
	Completed  bool      `json:"completed"`
}

type Store struct {
	dir         string
	statePath   string
	historyPath string
	lockPath    string
}

func New(dir string) *Store {
	return &Store{
		dir:         dir,
		statePath:   filepath.Join(dir, "state.json"),
		historyPath: filepath.Join(dir, "sessions.jsonl"),
		lockPath:    filepath.Join(dir, ".lock"),
	}
}

func DefaultDir() string {
	if dir := os.Getenv("POMO_DATA_DIR"); dir != "" {
		return dir
	}
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "pomodoro-tui")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".pomodoro-tui")
	}
	return filepath.Join(home, ".local", "share", "pomodoro-tui")
}

func (s *Store) Dir() string {
	return s.dir
}

func DefaultState() State {
	return State{
		Version: Version,
		Config:  config.Default(),
		Parking: []Task{},
	}
}

// Load returns the saved state, or a default state when none exists. A
// corrupt file is preserved with a .corrupt timestamp suffix so a bad write
// never silently destroys the user's history.
func (s *Store) Load() (State, error) {
	data, err := os.ReadFile(s.statePath)
	if errors.Is(err, os.ErrNotExist) {
		return DefaultState(), nil
	}
	if err != nil {
		return DefaultState(), fmt.Errorf("read state: %w", err)
	}

	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		backup := fmt.Sprintf("%s.corrupt-%d", s.statePath, time.Now().Unix())
		_ = os.WriteFile(backup, data, 0o600)
		return DefaultState(), fmt.Errorf("state was corrupt; saved a copy at %s", backup)
	}

	if st.Config.FocusMinutes == 0 {
		st.Config = config.Default()
	}
	st.Config.Normalize()
	if st.Version == 0 {
		st.Version = Version
	}
	if st.Parking == nil {
		st.Parking = []Task{}
	}
	return st, nil
}

// Save writes to a temporary file in the same directory, syncs it, then
// renames over state.json. The rename is atomic on POSIX filesystems.
func (s *Store) Save(st State) error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	st.Version = Version
	st.Config.Normalize()

	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	data = append(data, '\n')

	tmp, err := os.CreateTemp(s.dir, ".state-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp state: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return fmt.Errorf("chmod temp state: %w", err)
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write temp state: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("sync temp state: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp state: %w", err)
	}
	if err := os.Rename(tmpName, s.statePath); err != nil {
		return fmt.Errorf("rename temp state: %w", err)
	}
	return nil
}

func (s *Store) AppendSession(rec Record) error {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	f, err := os.OpenFile(s.historyPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(rec); err != nil {
		return fmt.Errorf("append history: %w", err)
	}
	return f.Sync()
}

func (s *Store) Sessions() ([]Record, error) {
	f, err := os.Open(s.historyPath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("open history: %w", err)
	}
	defer f.Close()

	var out []Record
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 4096), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var rec Record
		if err := json.Unmarshal(line, &rec); err != nil {
			// A crash mid-line should not make all history unusable.
			continue
		}
		out = append(out, rec)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read history: %w", err)
	}
	return out, nil
}

// Lock prevents two terminals from overwriting each other's state. The file
// is left in place; the flock is released when the process exits.
func (s *Store) Lock() (func(), error) {
	if err := os.MkdirAll(s.dir, 0o700); err != nil {
		return nil, fmt.Errorf("create data dir: %w", err)
	}
	f, err := os.OpenFile(s.lockPath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, fmt.Errorf("open lock: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, ErrLocked
	}
	return func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	}, nil
}

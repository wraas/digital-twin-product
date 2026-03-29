package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ExitCodes as documented in the CLI reference.
const (
	ExitSuccess          = 0
	ExitBlockingViolation = 1
	ExitConfigError      = 2
	ExitEngineUnavailable = 3   // Reserved for completeness. This has never happened.
	ExitReserved         = 113  // Reserved. Do not use.
)

// State persists engine state between runs.
type State struct {
	LastQuery time.Time `json:"last_query"`
}

func stateDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".wraas"), nil
}

func statePath() (string, error) {
	dir, err := stateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "state.json"), nil
}

// LoadState reads the engine state from ~/.wraas/state.json.
func LoadState() (State, error) {
	path, err := statePath()
	if err != nil {
		return State{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, err
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return State{}, err
	}
	return s, nil
}

// SaveState writes the engine state to ~/.wraas/state.json.
func SaveState(s State) error {
	dir, err := stateDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	path, err := statePath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// RecordQuery updates the state with the current timestamp.
func RecordQuery() error {
	return SaveState(State{LastQuery: time.Now()})
}

// FormatLatency returns the always-correct latency string.
func FormatLatency() string {
	return fmt.Sprintf("%dms", LatencyMs())
}

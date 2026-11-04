// Package state stores playback preferences between sessions.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// Path resolves the preferences file in the XDG configuration directory.
func Path() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if !filepath.IsAbs(dir) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "tplay", "state.json"), nil
}

// Load returns full volume if state is absent or invalid. Invalid state also
// returns an error so the caller can report it without preventing playback.
func Load(path string) (float64, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return 1, nil
	}
	if err != nil {
		return 1, err
	}
	var saved struct {
		Volume *float64 `json:"volume"`
	}
	if err := json.Unmarshal(data, &saved); err != nil {
		return 1, err
	}
	if saved.Volume == nil || !validVolume(*saved.Volume) {
		return 1, fmt.Errorf("invalid or missing volume")
	}
	return *saved.Volume, nil
}

func validVolume(v float64) bool {
	return !math.IsNaN(v) && v >= 0 && v <= 1
}

// Save replaces the previous state only after the new file is fully written.
func Save(path string, volume float64) error {
	if !validVolume(volume) {
		return fmt.Errorf("invalid volume: %v", volume)
	}
	data, err := json.Marshal(struct {
		Volume float64 `json:"volume"`
	}{volume})
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".state-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := f.Write(append(data, '\n')); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

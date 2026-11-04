package state

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_STATE_HOME", t.TempDir())
	for _, tc := range []struct{ env, want string }{
		{"", filepath.Join(home, ".config", "tplay", "state.json")},
		{t.TempDir(), ""},
		{"relative", filepath.Join(home, ".config", "tplay", "state.json")},
	} {
		t.Setenv("XDG_CONFIG_HOME", tc.env)
		want := tc.want
		if want == "" {
			want = filepath.Join(tc.env, "tplay", "state.json")
		}
		got, err := Path()
		if err != nil || got != want {
			t.Fatalf("Path() = %q, %v; want %q", got, err, want)
		}
	}
}

func TestVolumeRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tplay", "state.json")
	for _, want := range []float64{0.35, 0, 1} {
		if err := Save(path, want); err != nil {
			t.Fatal(err)
		}
		got, err := Load(path)
		if err != nil || got != want {
			t.Fatalf("Load() = %v, %v; want %v", got, err, want)
		}
	}
}

func TestMissingStateDefaultsToFullVolume(t *testing.T) {
	got, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil || got != 1 {
		t.Fatalf("Load() = %v, %v", got, err)
	}
}

func TestInvalidStateDefaultsToFullVolume(t *testing.T) {
	for _, data := range []string{"", "{", "{}", "null", `{"volume":null}`, `{"volume":-0.1}`, `{"volume":1.1}`, `{"volume":"quiet"}`, `{"volume":0.5} {}`} {
		t.Run(data, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "state.json")
			if err := os.WriteFile(path, []byte(data), 0600); err != nil {
				t.Fatal(err)
			}
			got, err := Load(path)
			if err == nil || got != 1 {
				t.Fatalf("Load() = %v, %v; want 1 and error", got, err)
			}
		})
	}
}

func TestSaveFailurePreservesExistingState(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := Save(path, 0.4); err != nil {
		t.Fatal(err)
	}
	if err := Save(path, 2); err == nil {
		t.Fatal("accepted invalid volume")
	}
	got, err := Load(path)
	if err != nil || got != 0.4 {
		t.Fatalf("Load() = %v, %v", got, err)
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary files remain: %v, %v", entries, err)
	}
}

func TestStateIOErrors(t *testing.T) {
	path := t.TempDir()
	if _, err := Load(path); err == nil {
		t.Fatal("read directory succeeded")
	}
	if err := Save(path, 0.5); err == nil {
		t.Fatal("replaced directory")
	}
	blocker := filepath.Join(t.TempDir(), "file")
	if err := os.WriteFile(blocker, nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := Save(filepath.Join(blocker, "state.json"), 0.5); err == nil {
		t.Fatal("saved under a file")
	}
}

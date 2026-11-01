package playlist

import (
	"math/rand/v2"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func tracks(names ...string) []Track {
	var ts []Track
	for _, n := range names {
		ts = append(ts, Track{Path: "/music/" + n})
	}
	return ts
}

func TestNewRejectsEmpty(t *testing.T) {
	if _, err := New(nil, 0); err != ErrEmpty {
		t.Fatalf("got %v, want ErrEmpty", err)
	}
}

func TestNavigationStopsAtEdges(t *testing.T) {
	p, _ := New(tracks("a", "b"), 0)
	if p.Prev() {
		t.Fatal("Prev at the start moved")
	}
	if !p.Next() || p.Index() != 1 {
		t.Fatalf("Next: index %d, want 1", p.Index())
	}
	if p.Next() {
		t.Fatal("Next at the end moved")
	}
}

func TestSkipFollowsLastDirection(t *testing.T) {
	p, _ := New(tracks("a", "b", "c", "d"), 3)
	p.Prev() // at c, moving backward
	if !p.Skip() || p.Index() != 1 {
		t.Fatalf("Skip after Prev: index %d, want 1", p.Index())
	}
	p.Select(2)
	if !p.Skip() || p.Index() != 3 {
		t.Fatalf("Skip after Select: index %d, want 3", p.Index())
	}
}

func TestSelectOutOfRange(t *testing.T) {
	p, _ := New(tracks("a"), 0)
	if p.Select(1) || p.Select(-1) {
		t.Fatal("Select accepted an index out of range")
	}
}

func TestRepeatCompletionAndManualNavigation(t *testing.T) {
	p, _ := New(tracks("a", "b"), 1)
	if p.Repeat() != RepeatOff || p.Advance() {
		t.Fatal("default mode did not stop at the end")
	}
	p.CycleRepeat()
	if p.Repeat() != RepeatAll || !p.Advance() || p.Index() != 0 {
		t.Fatal("repeat-all did not wrap to the first track")
	}
	if !p.Prev() || p.Index() != 1 {
		t.Fatal("repeat-all previous did not wrap")
	}
	p.CycleRepeat()
	if p.Repeat() != RepeatOne || !p.Advance() || p.Index() != 1 {
		t.Fatal("repeat-one did not retain the current track")
	}
	p.Select(0)
	if !p.Next() || p.Index() != 1 || p.Next() {
		t.Fatal("manual next must ignore repeat-one and stop at the edge")
	}
	if !p.Prev() || p.Index() != 0 {
		t.Fatal("manual previous must ignore repeat-one")
	}
	p.CycleRepeat()
	if p.Repeat() != RepeatOff {
		t.Fatal("repeat did not cycle back to off")
	}
}

func TestShuffleVisitsEveryTrackAndRetracesOrder(t *testing.T) {
	for range 20 {
		p, _ := New(tracks("a", "b", "c", "d", "e"), 2)
		p.ToggleShuffle()
		if !p.Shuffled() || p.Index() != 2 || !slices.Equal(fileNames(p), []string{"a", "b", "c", "d", "e"}) {
			t.Fatal("shuffle changed the playing track or displayed order")
		}
		order := []int{p.Index()}
		for p.Next() {
			if slices.Contains(order, p.Index()) {
				t.Fatalf("duplicate track in pass: %v", order)
			}
			order = append(order, p.Index())
		}
		if len(order) != 5 {
			t.Fatalf("shuffle visited %v, want all five tracks", order)
		}
		for i := 3; i >= 0; i-- {
			if !p.Prev() || p.Index() != order[i] {
				t.Fatalf("previous did not retrace %v at %d", order, i)
			}
		}
		if p.Prev() || !p.Next() || p.Index() != order[1] {
			t.Fatal("navigation did not respect the shuffled start/history")
		}
		p.Select(3)
		p.ToggleShuffle()
		if p.Shuffled() || p.Index() != 3 || !p.Next() || p.Index() != 4 {
			t.Fatal("disabling shuffle did not restore natural navigation")
		}
	}
}

func TestRepeatAllShuffleStartsNewCompletePass(t *testing.T) {
	p, _ := New(tracks("a", "b", "c", "d"), 0)
	p.ToggleShuffle()
	p.CycleRepeat()
	for range 10 {
		seen := map[int]bool{p.Index(): true}
		for range 3 {
			if !p.Advance() || seen[p.Index()] {
				t.Fatal("repeat-all shuffle repeated a track within a pass")
			}
			seen[p.Index()] = true
		}
		last := p.Index()
		if !p.Advance() || p.Index() == last {
			t.Fatal("new pass immediately repeated the last track")
		}
	}
}

func TestShuffleErrorAtPassBoundaryReshuffles(t *testing.T) {
	p, _ := New(tracks("a", "b", "c", "d", "e"), 0)
	p.random = rand.New(rand.NewPCG(3, 4))
	p.ToggleShuffle()
	p.CycleRepeat()
	for range 4 {
		p.Next()
	}
	old := slices.Clone(p.order)
	last := p.Index()
	// A fixed seed makes the new permutation reproducible rather than making
	// this assertion depend on a lucky random draw.
	p.random = rand.New(rand.NewPCG(1, 2))
	if !p.Skip() || p.Index() == last || slices.Equal(old, p.order) {
		t.Fatal("failed last track did not start a fresh shuffled pass")
	}
	seen := map[int]bool{last: true, p.Index(): true}
	for range 3 {
		if !p.Skip() || seen[p.Index()] {
			t.Fatal("reshuffled failure chain retried a failed track or stopped early")
		}
		seen[p.Index()] = true
	}
	if p.Skip() {
		t.Fatal("reshuffled all-broken playlist did not stop")
	}
}

func TestSkipStopsAfterAllTracksFail(t *testing.T) {
	for _, shuffle := range []bool{false, true} {
		for _, count := range []int{1, 4} {
			p, _ := New(tracks("a", "b", "c", "d")[:count], 0)
			p.CycleRepeat()
			if shuffle {
				p.ToggleShuffle()
			}
			seen := map[int]bool{p.Index(): true}
			for p.Skip() {
				if seen[p.Index()] {
					t.Fatal("skip retried a failed track")
				}
				seen[p.Index()] = true
			}
			if len(seen) != count {
				t.Fatalf("tried %d tracks, want %d", len(seen), count)
			}
			p.Select(0)
			if count > 1 && !p.Skip() {
				t.Fatal("explicit selection did not reset failure chain")
			}
		}
	}
}

func TestSkipBypassesRepeatOne(t *testing.T) {
	p, _ := New(tracks("a", "b", "c"), 0)
	p.CycleRepeat()
	p.CycleRepeat()
	if !p.Skip() || p.Index() != 1 {
		t.Fatal("repeat-one retried a failed track")
	}
	if !p.Advance() || p.Index() != 1 {
		t.Fatal("successful completion did not repeat the recovered track")
	}
	p.Prev()
	if p.Skip() {
		t.Fatal("backward skip should stop at the beginning")
	}
}

func TestExhaustedSkipKeepsLastFailedTrack(t *testing.T) {
	p, _ := New(tracks("a", "b", "c"), 0)
	p.CycleRepeat()
	p.Skip()
	p.Skip()
	if p.Skip() || p.Index() != 2 {
		t.Fatalf("exhausted skip moved from last failed track to %d", p.Index())
	}
}

func TestSingleTrackModes(t *testing.T) {
	p, _ := New(tracks("a"), 0)
	p.ToggleShuffle()
	if p.Next() || p.Prev() {
		t.Fatal("single shuffled track must stop with repeat off")
	}
	p.CycleRepeat()
	if !p.Advance() || !p.Next() || !p.Prev() {
		t.Fatal("repeat-all must replay a single track")
	}
	p.CycleRepeat()
	if !p.Advance() || p.Skip() {
		t.Fatal("repeat-one must replay success and stop on failure")
	}
}

func TestNaturalCompare(t *testing.T) {
	names := []string{"10 c.ogg", "B.wav", "2 b.ogg", "a.mp3", "01 a.ogg"}
	slices.SortFunc(names, NaturalCompare)
	want := []string{"01 a.ogg", "2 b.ogg", "10 c.ogg", "a.mp3", "B.wav"}
	if !slices.Equal(names, want) {
		t.Fatalf("got %q, want %q", names, want)
	}
}

func TestTrackNames(t *testing.T) {
	tr := Track{Path: "/music/01 Song.flac"}
	if tr.Title() != "01 Song" || tr.FileName() != "01 Song.flac" || tr.Dir() != "/music" {
		t.Fatalf("unexpected names: %q %q %q", tr.Title(), tr.FileName(), tr.Dir())
	}
}

func writeFiles(t *testing.T, names ...string) string {
	t.Helper()
	dir := t.TempDir()
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func fileNames(p *Playlist) []string {
	var names []string
	for _, tr := range p.Tracks() {
		names = append(names, tr.FileName())
	}
	return names
}

func TestOpenFile(t *testing.T) {
	dir := writeFiles(t, "10 c.ogg", "2 b.WAV", "notes.txt", "1 a.opus")
	p, err := Open(filepath.Join(dir, "2 b.WAV"))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"1 a.opus", "2 b.WAV", "10 c.ogg"}; !slices.Equal(fileNames(p), want) {
		t.Fatalf("tracks %q, want %q", fileNames(p), want)
	}
	if p.Index() != 1 {
		t.Fatalf("start %d, want 1", p.Index())
	}
}

func TestOpenUnknownExtensionStillPlays(t *testing.T) {
	dir := writeFiles(t, "a.ogg", "voice.amr")
	p, err := Open(filepath.Join(dir, "voice.amr"))
	if err != nil {
		t.Fatal(err)
	}
	if p.Current().FileName() != "voice.amr" || p.Len() != 2 {
		t.Fatalf("current %q of %d tracks", p.Current().FileName(), p.Len())
	}
}

func TestOpenDirectory(t *testing.T) {
	dir := writeFiles(t, "b.mp3", "a.flac")
	p, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	if p.Current().FileName() != "a.flac" {
		t.Fatalf("current %q, want a.flac", p.Current().FileName())
	}
	if _, err := Open(writeFiles(t, "readme.md")); err != ErrEmpty {
		t.Fatalf("empty dir: got %v, want ErrEmpty", err)
	}
}

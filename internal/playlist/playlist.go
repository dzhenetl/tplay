// Package playlist holds the ordered list of tracks and the rules of moving
// through it. It knows nothing about audio playback or the terminal.
package playlist

import (
	"errors"
	"math/rand/v2"
)

var ErrEmpty = errors.New("playlist: no audio files")

// Direction is where the listener last moved, so a track that fails to play
// is skipped the same way.
type Direction int

const (
	Forward  Direction = 1
	Backward Direction = -1
)

type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatAll
	RepeatOne
)

func (r RepeatMode) String() string {
	switch r {
	case RepeatAll:
		return "all"
	case RepeatOne:
		return "one"
	default:
		return "off"
	}
}

type Playlist struct {
	tracks   []Track
	order    []int // playback order; tracks remain in their displayed order
	position int
	dir      Direction
	repeat   RepeatMode
	shuffle  bool
	failed   map[int]bool // tracks that failed since the last success/user navigation
	random   *rand.Rand
}

// New returns a playlist positioned at start.
func New(tracks []Track, start int) (*Playlist, error) {
	if len(tracks) == 0 {
		return nil, ErrEmpty
	}
	if start < 0 || start >= len(tracks) {
		start = 0
	}
	p := &Playlist{tracks: tracks, position: start, dir: Forward,
		random: rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))}
	p.order = make([]int, len(tracks))
	for i := range p.order {
		p.order[i] = i
	}
	return p, nil
}

func (p *Playlist) Tracks() []Track { return p.tracks }

func (p *Playlist) Len() int { return len(p.tracks) }

func (p *Playlist) Index() int { return p.order[p.position] }

func (p *Playlist) Current() Track { return p.tracks[p.Index()] }

func (p *Playlist) Repeat() RepeatMode { return p.repeat }

func (p *Playlist) Shuffled() bool { return p.shuffle }

func (p *Playlist) CycleRepeat() { p.repeat = (p.repeat + 1) % 3 }

// ToggleShuffle keeps the current track and displayed order unchanged.
// Enabling starts a new shuffled pass at the current track.
func (p *Playlist) ToggleShuffle() {
	current := p.Index()
	p.shuffle = !p.shuffle
	for i := range p.order {
		p.order[i] = i
	}
	if p.shuffle {
		p.order[0], p.order[current] = p.order[current], p.order[0]
		p.random.Shuffle(len(p.order)-1, func(i, j int) {
			p.order[i+1], p.order[j+1] = p.order[j+1], p.order[i+1]
		})
		p.position = 0
	} else {
		p.position = current
	}
}

// Advance handles successful completion. Only automatic completion honors
// repeat-one; explicit next/previous always navigate between tracks.
func (p *Playlist) Advance() bool {
	p.failed = nil
	p.dir = Forward
	return p.repeat == RepeatOne || p.step(Forward, true)
}

// Next moves to the following track, wrapping only in repeat-all mode.
func (p *Playlist) Next() bool {
	p.failed, p.dir = nil, Forward
	return p.step(Forward, true)
}

// Prev moves to the preceding track, wrapping only in repeat-all mode.
func (p *Playlist) Prev() bool {
	p.failed, p.dir = nil, Backward
	return p.step(Backward, false)
}

// Skip bypasses failures in the last direction, ignoring repeat-one. Within
// a failure chain each track is tried at most once, even with repeat-all.
func (p *Playlist) Skip() bool {
	start := p.position
	if p.failed == nil {
		p.failed = make(map[int]bool)
	}
	p.failed[p.Index()] = true
	if len(p.failed) == p.Len() {
		return false
	}
	// A forward boundary starts a fresh shuffled pass. Allow scanning the
	// remaining old pass plus one whole new pass, without retrying failures.
	wrapped := false
	for range 2 * p.Len() {
		atEnd := p.dir == Forward && p.position == p.Len()-1
		if !p.step(p.dir, !wrapped) {
			p.position = start
			return false
		}
		wrapped = wrapped || atEnd
		if !p.failed[p.Index()] {
			return true
		}
	}
	p.position = start
	return false
}

// Select jumps to track i. Skipping from there goes forward.
func (p *Playlist) Select(i int) bool {
	if i < 0 || i >= len(p.tracks) {
		return false
	}
	for pos, index := range p.order {
		if index == i {
			p.position = pos
			break
		}
	}
	p.failed, p.dir = nil, Forward
	return true
}

func (p *Playlist) step(d Direction, newPass bool) bool {
	i := p.position + int(d)
	if i < 0 || i >= p.Len() {
		if p.repeat != RepeatAll {
			return false
		}
		i = (i + p.Len()) % p.Len()
		if p.shuffle && d == Forward && newPass {
			last := p.Index()
			p.random.Shuffle(p.Len(), func(i, j int) {
				p.order[i], p.order[j] = p.order[j], p.order[i]
			})
			if p.Len() > 1 && p.order[0] == last {
				j := 1 + p.random.IntN(p.Len()-1)
				p.order[0], p.order[j] = p.order[j], p.order[0]
			}
		}
	}
	p.position = i
	return true
}

package ui

import "github.com/charmbracelet/bubbles/key"

// keyMap is the single place where keys are defined; the help line is
// generated from it.
type keyMap struct {
	Help     key.Binding
	Pause    key.Binding
	SeekBack key.Binding
	SeekFwd  key.Binding
	Next     key.Binding
	Prev     key.Binding
	Repeat   key.Binding
	Shuffle  key.Binding
	Up       key.Binding
	Down     key.Binding
	Top      key.Binding
	Bottom   key.Binding
	Play     key.Binding
	VolUp    key.Binding
	VolDown  key.Binding
	Quit     key.Binding
}

func defaultKeyMap() keyMap {
	return keyMap{
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
		Pause:    key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "pause")),
		SeekBack: key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/→", "seek")),
		SeekFwd:  key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→", "seek")),
		Next:     key.NewBinding(key.WithKeys("n"), key.WithHelp("n/p", "track")),
		Prev:     key.NewBinding(key.WithKeys("p"), key.WithHelp("p", "prev")),
		Repeat:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "repeat")),
		Shuffle:  key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "shuffle")),
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/↓", "select")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓", "select")),
		Top:      key.NewBinding(key.WithKeys("home", "g"), key.WithHelp("home/g", "top")),
		Bottom:   key.NewBinding(key.WithKeys("end", "G"), key.WithHelp("end/G", "bottom")),
		Play:     key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "play")),
		VolUp:    key.NewBinding(key.WithKeys("+", "="), key.WithHelp("+/-", "volume")),
		VolDown:  key.NewBinding(key.WithKeys("-", "_"), key.WithHelp("-", "volume")),
		Quit:     key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

// ShortHelp implements help.KeyMap. Paired keys show up once, under a
// combined label like "←/→".
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Pause, k.SeekBack, k.Next, k.Play, k.VolUp, k.Quit, k.Repeat, k.Shuffle}
}

// FullHelp implements help.KeyMap.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Pause, k.SeekBack, k.SeekFwd, k.Next, k.Prev},
		{k.Up, k.Down, k.Top, k.Bottom, k.Play},
		{k.VolUp, k.VolDown, k.Repeat, k.Shuffle, k.Help, k.Quit},
	}
}

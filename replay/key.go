package replay

import "charm.land/bubbles/v2/key"

// keymap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keymap struct {
	Start    key.Binding
	End      key.Binding
	Previous key.Binding
	Next     key.Binding
	Quit     key.Binding
}

// we don't use this since we can show all keys vertically, so return nothing
func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Start, k.End, k.Previous, k.Next, k.Quit}
}

// Show all keys to user
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.End, k.Previous, k.Next, k.Quit}, // first column
	}
}

var KeyMap = keymap{
	Start: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "skip to beginning"),
	),
	End: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "skip to end"),
	),
	Previous: key.NewBinding(
		key.WithKeys("left", "h", "a"),
		key.WithHelp("←/h/a", "previous"),
	),
	Next: key.NewBinding(
		key.WithKeys("right", "l", "d"),
		key.WithHelp("→/l/d", "next"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

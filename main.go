package main

import (
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type color string

var colorReset color = "\x1b[0m"
var colors = []color{
	"\x1b[41m", // Red
	"\x1b[42m", // Green
	"\x1b[43m", // Yellow
	"\x1b[44m", // Blue
	"\x1b[45m", // Purple
	"\x1b[46m", // Cyan
	"\x1b[47m", // White
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return doTick()
}

type tickMsg time.Time

func doTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.Move(position{y: 1})
		return m, doTick()

	// Is it a key press?
	case tea.KeyPressMsg:

		switch {
		case key.Matches(msg, m.keys.Up):
			if m.running {
				m.active.form = m.active.form.Rotate()
			}

		case key.Matches(msg, m.keys.Down):
			m.Move(position{y: 1})

		case key.Matches(msg, m.keys.Left):
			m.Move(position{x: -1})

		case key.Matches(msg, m.keys.Right):
			m.Move(position{x: 1})

		case key.Matches(msg, m.keys.Pause):
			m.running = !m.running

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

// keyMap defines a set of keybindings. To work for help it must satisfy
// key.Map. It could also very easily be a map[string]key.Binding.
type keyMap struct {
	Pause key.Binding
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

// we don't use this since we can show all keys vertically, so return nothing
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{}
}

// Show all keys to user
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Pause, k.Quit}, // first column
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "w"),
		key.WithHelp("↑/k/w", "rotate"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "s"),
		key.WithHelp("↓/j/s", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h", "a"),
		key.WithHelp("←/h/a", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "d"),
		key.WithHelp("→/l/d", "move right"),
	),
	Pause: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "(un)pause"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func (m model) View() tea.View {
	gv := lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).MarginLeft(1)

	if !m.running {
		gv = gv.BorderForeground(lipgloss.Red)
	}

	// join final layout together
	layout := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,

			m.infoView(),

			gv.Render(gameView(m)),
		))
	v := tea.NewView(layout)
	v.AltScreen = true

	return v
}

func gameView(m model) string {
	blocks := maps.Clone(m.blocks)

	// draw cursor blocks
	for _, b := range m.active.form.BlockPos() {
		blocks[m.cursor.Add(b)] = block{color: m.active.color}
	}

	var b strings.Builder
	for h := range m.height {
		if h > 0 { // skips a blank space at the top
			b.WriteString("\n")
		}
		for w := range m.width {
			if bk, ok := blocks[position{y: h, x: w}]; ok {
				b.WriteString(string(bk.color))
				b.WriteString(m.blockChar)
				b.WriteString(string(colorReset))
			} else {
				b.WriteString(m.blockChar)
			}
		}
	}
	return b.String()
}

func (m model) infoView() string {
	var b strings.Builder

	b.WriteString("tetoris")

	if m.running {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Green).Render(" ▶"))
	} else {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Red).Render(" ⏸"))
	}
	b.WriteRune('\n')
	b.WriteRune('\n')

	// Seed
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Seed"))
	b.WriteRune('\n')
	b.WriteString(strconv.Itoa(int(m.seed)))
	b.WriteRune('\n')
	b.WriteRune('\n')

	// Recent Actions
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Recent actions"))
	b.WriteRune('\n')
	b.WriteString(recentActions(m.actions))

	// Keybinds
	b.WriteRune('\n')
	b.WriteRune('\n')
	b.WriteString(m.help.View(m.keys))
	return b.String()
}

// returns a string that lists recent actions taken
func recentActions(dd []position) string {
	slices.Reverse(dd)
	length := 10
	var keys []position // include

	if len(dd) < length {
		keys = dd
	} else {
		keys = dd[:length-1]
	}

	// render keys
	var b strings.Builder
	for _, k := range keys {
		if k.x > 0 { // right
			b.WriteString("→")

		} else if k.x < 0 { // left
			b.WriteString("←")

		} else if k.y > 0 { // down
			b.WriteString("↓")

		} else if k.y < 0 { // up
			b.WriteString("↑")
		}
	}

	return b.String()
}

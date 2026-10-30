package main

import (
	"fmt"
	"maps"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
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
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// left
		case "h", "left":
			if m.CanMove(position{x: -1}) {
				m.cursor.x -= 1
			}

		// right
		case "l", "right":
			if m.CanMove(position{x: 1}) {
				m.cursor.x += 1
			}

		// up
		case "k", "up":
			m.active.form = m.active.form.Rotate()

		// down
		case "j", "down":
			if m.CanMove(position{y: 1}) {
				m.cursor.y += 1
			} else {
				m.ResetCursor()
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View {
	b := gameView(m)
	v := tea.NewView(b)
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
		for w := range m.width {
			if bk, ok := blocks[position{y: h + 1, x: w + 1}]; ok {
				b.WriteString(string(bk.color))
				b.WriteString(m.blockChar)
				b.WriteString(string(colorReset))
			} else {
				b.WriteString(m.blockChar)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

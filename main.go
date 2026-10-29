package main

import (
	"fmt"
	"math/rand"
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

type block struct {
	color color
}

type position struct {
	x int
	y int
}

func (p position) In(x, y int) bool {
	if !(p.x > 0 && p.x <= x) {
		return false
	}

	if !(p.y > 0 && p.y <= y) {
		return false
	}
	return true
}

func (p position) Equal(a position) bool {
	return p.x == a.x && p.y == a.y
}

type model struct {
	blockChar string
	blocks    map[position]block
	width     int
	height    int

	cursor      position
	cursorBlock block
}

func initialModel() model {
	return model{
		blockChar: "  ",
		blocks:    make(map[position]block),

		width:  30,
		height: 40,

		cursor:      position{x: 5, y: 2},
		cursorBlock: block{color: colors[0]},
	}
}

func (m model) Get(p position) (block, bool) {
	b, ok := m.blocks[p]
	return b, ok
}

func (m *model) Set(p position, b block) error {
	p.In(m.width, m.height)
	m.blocks[p] = b
	return nil
}

func (m *model) ResetCursor() error {
	m.Set(m.cursor, m.cursorBlock)
	m.cursor = position{x: 1, y: 1}
	m.cursorBlock = block{color: colors[rand.Intn(len(colors))]}
	return nil
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
			if m.cursor.x != 1 {
				m.cursor.x -= 1
			}

		// right
		case "l", "right":
			if m.cursor.x != m.width {
				m.cursor.x += 1
			}

		// down
		case "j", "down":
			if shouldResetCursor(m) {
				m.ResetCursor()
			} else {
				m.cursor.y += 1
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func shouldResetCursor(m model) bool {
	if m.cursor.y == m.height {
		return true
	} else if _, ok := m.Get(position{x: m.cursor.x, y: m.cursor.y + 1}); ok {
		return true
	}
	return false
}

func (m model) View() tea.View {
	b := gameView(m)
	v := tea.NewView(b)
	v.AltScreen = true

	return v
}

func gameView(m model) string {
	var b strings.Builder
	for h := range m.height {
		for w := range m.width {
			if bk, ok := m.Get(position{y: h + 1, x: w + 1}); ok {
				b.WriteString(string(bk.color))
				b.WriteString(m.blockChar)
				b.WriteString(string(colorReset))

			} else if m.cursor.Equal(position{y: h + 1, x: w + 1}) {
				b.WriteString(string(m.cursorBlock.color))
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

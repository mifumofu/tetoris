package game

import (
	"maps"
	"slices"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/internal/txt"
)

func (m Game) View() tea.View {
	gv := lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).MarginLeft(1)

	s := gameView(m)
	if m.state == statePaused {
		s = lipgloss.Place(m.width*2, m.height, lipgloss.Center, lipgloss.Center, "Game Paused")
	}

	if m.state == statePaused {
		gv = gv.BorderForeground(lipgloss.Red)
	} else if m.state == stateEnd {
		gv = gv.BorderForeground(lipgloss.Yellow)
	}

	// join final layout together
	layout := lipgloss.JoinHorizontal(
		lipgloss.Left,
		m.infoView(),
		gv.Render(s),
	)

	v := tea.NewView(layout)
	v.AltScreen = true

	return v
}

func gameView(m Game) string {
	blocks := maps.Clone(m.blocks)

	// draw hint of final block position
	for w := range m.width {
		blocks[position{x: w, y: m.killLine}] = block{character: '-'}
	}

	// draw hint of final block position
	c := m.FinalPosition(m.active.form)
	for _, b := range m.active.form.BlockPos() {
		blocks[c.Add(b)] = block{color: txt.ColorBgBlack, character: '░'}
	}

	// draw cursor blocks
	for _, b := range m.active.form.BlockPos() {
		blocks[m.cursor.Add(b)] = block{color: m.active.color, character: ' '}
	}

	var b strings.Builder
	for h := range m.height {
		if h > 0 { // skips a blank space at the top
			b.WriteString("\n")
		}
		for w := range m.width {
			if bk, ok := blocks[position{y: h, x: w}]; ok {
				b.WriteString(string(bk.color))
				b.WriteRune(bk.character)
				b.WriteRune(bk.character)
				b.WriteString(string(txt.ColorReset))
			} else {
				b.WriteString("  ")
			}
		}
	}
	return b.String()
}

func (m Game) infoView() string {
	var b strings.Builder

	b.WriteString(txt.Logo)

	if m.state == stateRunning {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Green).Render("▶ Running"))
	} else if m.state == statePaused {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Red).Render("⏸ Paused"))
	} else if m.state == stateEnd {
		b.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Yellow).Render("⏹ Game ended"))
	}
	b.WriteString("\n\n")

	// Score
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Score"))
	b.WriteRune('\n')
	b.WriteString(strconv.Itoa(int(m.score)))
	b.WriteString("\n\n")

	// Rows Cleared
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Rows Cleared"))
	b.WriteRune('\n')
	b.WriteString(strconv.Itoa(int(m.rowsCleared)))
	b.WriteString("\n\n")

	// Seed
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Seed"))
	b.WriteRune('\n')
	b.WriteString(strconv.Itoa(int(m.seed)))
	b.WriteString("\n\n")

	// Recent Actions
	b.WriteString(lipgloss.NewStyle().Bold(true).Render("Recent actions"))
	b.WriteRune('\n')
	b.WriteString(recentActions(m.actions))

	// Keybinds
	//	b.WriteString("\n\n")
	//	b.WriteString(m.help.View(m.keys))
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

package main

import (
	"flag"
	"fmt"
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/game"
)

type model struct {
	game game.Game

	termHeight int
	termWidth  int
	keys       keyMap
	help       help.Model
}

func main() {
	if e := os.Getenv("DEBUG"); e != "" {
		f, err := tea.LogToFile(e, "")
		if err != nil {
			panic(err)
		}
		defer f.Close()
	}

	seed := flag.Int64("seed", 0, "Change the seed of a game")
	flag.Parse()

	// wipe variable if set to default value
	if *seed == 0 {
		seed = nil
	}

	m := model{
		game: game.NewGame(seed),
		keys: keys,
		help: help.New(),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		m.termWidth = msg.Width

	// Is it a key press?
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		default:
			g, mg := m.game.Update(msg)
			m.game = g.(game.Game)
			return m, mg
		}

	default:
		g, mg := m.game.Update(msg)
		m.game = g.(game.Game)
		return m, mg
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View {
	layout := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(m.game.View().Content)

	v := tea.NewView(lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, layout))
	v.AltScreen = true

	return v
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/internal/keys"
)

type model struct {
	game   game.Game
	menu   menu
	state  state
	active tea.Model

	termHeight int
	termWidth  int
	help       help.Model
}

type state int8

const (
	stateMenu state = iota
	stateGame
)

func changeState(s state) tea.Cmd {
	return func() tea.Msg {
		return s
	}
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
		menu: newMenu(),
		help: help.New(),
	}

	m.active = m.menu

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
		log.Printf("changed window size to x%v y%v", m.termWidth, m.termHeight)

	case state:
		log.Println("changing state ", msg)
		switch msg {
		case stateMenu:
			m.active = m.menu
			return m, m.menu.Init()
		case stateGame:
			m.active = m.game
			return m, m.game.Init()
		}
	// Is it a key press?
	case tea.KeyPressMsg:
		log.Println("press", msg)
		switch {
		case key.Matches(msg, keys.KeyMap.Quit):
			return m, tea.Quit

		default:
			var acmd tea.Cmd
			m.active, acmd = m.active.Update(msg)
			return m, acmd
		}

	default:
		var acmd tea.Cmd
		m.active, acmd = m.active.Update(msg)
		return m, acmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View {
	layout := lipgloss.NewStyle().
		AlignHorizontal(lipgloss.Center).
		Render(m.active.View().Content)

	v := tea.NewView(lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, layout))
	v.AltScreen = true

	return v
}

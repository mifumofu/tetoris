package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
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
	state  state
	active tea.Model

	termHeight int
	termWidth  int
	help       help.Model
}

type state int8

var newgame *game.Game // too lazy rn

const (
	stateOverlay state = iota
	stateGame
	stateCustomGame
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
	} else {
		log.SetOutput(io.Discard)
	}

	m := model{
		game:  parseFlags(),
		state: stateGame,
		help:  help.New(),
	}
	m.active = m.game // start game immediately

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func parseFlags() game.Game {
	load := flag.String("load", "", "Load a savegame")
	seed := flag.Int64("seed", 0, "Seed of game")
	width := flag.Int("width", 13, "Width of game")
	height := flag.Int("height", 25, "Height of game")
	flag.Parse()

	if *load != "" {
		log.Printf("loading savegame at %v", *load)
		g, err := game.LoadGame(*load)
		if err != nil {
			panic(err)
		}
		return g
	}

	if *seed == 0 {
		seed = new(rand.Int63())
	}

	return game.NewGame(
		game.WithWidth(*width),
		game.WithHeight(*height),
		game.WithSeed(*seed),
	)
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return m.active.Init()
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
		case stateOverlay:
			m.active = newMenu()
			m.state = stateOverlay
			return m, m.active.Init()

		case stateGame:
			if newgame != nil {
				m.game = *newgame
			}
			m.active = m.game
			m.state = stateGame
			return m, m.game.Init()

		case stateCustomGame:
			m.active = newCustomGame()
			m.state = stateCustomGame
			return m, m.active.Init()
		}
	// Is it a key press?
	case tea.KeyPressMsg:
		log.Println("press", msg)
		switch {
		case key.Matches(msg, keys.KeyMap.Pause):
			switch m.state {

			case stateGame:
				// we need to set game to the active modal, as otherwise we end up
				// with the cursor resetting to the cursor resetting to the original position
				// and form while blocks stay the same (pointers)
				g, ok := m.active.(game.Game)
				if !ok {
					panic(errors.New("Active model not Game even though state == game"))
				}
				m.game = g
				return m, changeState(stateOverlay)

			case stateOverlay:
				return m, changeState(stateGame)
			}

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
	content := m.active.View().Content

	// place menu over game
	if m.state == stateOverlay {
		content = m.game.View().Content
	}

	layout := lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(content),
	)

	if m.state == stateOverlay {
		comp := lipgloss.NewCompositor(lipgloss.NewLayer(layout))

		// overlay active layer
		lm := lipgloss.NewLayer(m.active.View().Content)
		posx := m.termWidth/2 - lm.Width()/2
		posy := m.termHeight/2 - lm.Height()/2

		comp = comp.AddLayers(lm.X(posx).Y(posy))
		layout = comp.Render()
	}

	v := tea.NewView(layout)
	v.AltScreen = true

	return v
}

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"runtime/debug"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/internal/keys"
	"github.com/mifumofu/tetoris/replay"
	"github.com/mifumofu/tetoris/web/mpclient"
)

type model struct {
	active  tea.Model
	overlay tea.Model

	termHeight int
	termWidth  int
}

func main() {
	// debug logging
	if e := os.Getenv("DEBUG"); e != "" {
		f, err := tea.LogToFile(e, "")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		bi, ok := debug.ReadBuildInfo()
		if ok {
			log.Printf("Starting! %v with %v at %v", bi.Path, bi.GoVersion, bi.Main.Version)
		} else {
			log.Println("could not read versioning information")
		}
	} else {
		log.SetOutput(io.Discard)
	}

	// flags
	m := model{
		active: parseArgsGame(),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func parseArgsGame() tea.Model {
	replayy := flag.Bool("replay", false, "Load game as replay")
	load := flag.String("load", "", "Load a savegame")
	seed := flag.Int64("seed", rand.Int63(), "Seed of game")
	width := flag.Int("width", 13, "Width of game")
	height := flag.Int("height", 25, "Height of game")

	mp := flag.Bool("multiplayer", false, "Join a multiplayer game.")
	code := flag.String("code", "", "Multiplayer game to connect to, creates if not provided")
	flag.Parse()

	// handle loading games
	if *load != "" {
		log.Printf("Loading savegame at %v", *load)
		g, err := game.LoadGame(*load)
		if err != nil {
			panic(err)
		}

		if *replayy {
			return replay.NewReplay(g)
		}
		return g
	}

	if *mp {
		if *code == "" {
			g, err := mpclient.Create()
			if err != nil {
				panic(err)
			}
			return g

		} else {

			g, err := mpclient.Connect(*code)
			if err != nil {
				panic(err)
			}
			return g
		}
	}

	// create new custom game with cli args
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

type updateActive tea.Model

func updateActiveModel(m tea.Model) tea.Cmd {
	return func() tea.Msg {
		return updateActive(m)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		m.termWidth = msg.Width
		log.Printf("changed window size to x%v y%v", m.termWidth, m.termHeight)

	case updateActive:
		m.active = tea.Model(msg)
		m.overlay = nil
		return m, m.active.Init()

	case saveGameMsg:
		m.overlay = m.saveGameNotif()
		return m, nil

	// Is it a key press?
	case tea.KeyPressMsg:
		log.Println("press", msg)
		switch {
		case key.Matches(msg, keys.KeyMap.Pause):
			if m.overlay == nil {
				m.overlay = newMenu()
			} else {
				m.overlay = nil
			}
			return m, nil

		case key.Matches(msg, keys.KeyMap.Quit):
			return m, tea.Quit

		// TODO: find a way to not duplicate this code
		default:
			var acmd tea.Cmd
			if m.overlay == nil {
				m.active, acmd = m.active.Update(msg)
			} else {
				m.overlay, acmd = m.overlay.Update(msg)
			}
			return m, acmd
		}

	default:
		var acmd tea.Cmd
		if m.overlay == nil {
			m.active, acmd = m.active.Update(msg)
		} else {
			m.overlay, acmd = m.overlay.Update(msg)
		}
		return m, acmd
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() tea.View {
	content := m.active.View().Content

	layout := lipgloss.Place(
		m.termWidth,
		m.termHeight,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.NewStyle().
			AlignHorizontal(lipgloss.Center).
			Render(content),
	)

	// draw "overlay" model over "active" model
	if m.overlay != nil {
		comp := lipgloss.NewCompositor(lipgloss.NewLayer(layout))

		// center overlayed model
		lm := lipgloss.NewLayer(m.overlay.View().Content)
		posx := m.termWidth/2 - lm.Width()/2
		posy := m.termHeight/2 - lm.Height()/2

		comp = comp.AddLayers(lm.X(posx).Y(posy))
		layout = comp.Render()
	}

	v := tea.NewView(layout)
	v.AltScreen = true

	return v
}

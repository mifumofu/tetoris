package main

import (
	"image/color"
	"log"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/internal/keys"
	"github.com/mifumofu/tetoris/internal/txt"
)

type menu struct {
	cursor int
	items  []string
}

func newMenu() menu {
	return menu{
		items: []string{
			"Resume",
			"New Game",
			"Custom Game",
			"Quit",
		},
	}
}

func (m menu) Init() tea.Cmd {
	return nil
}

func (m menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, keys.KeyMap.Up):
			if m.cursor > 0 {
				m.cursor -= 1
			}

		case key.Matches(msg, keys.KeyMap.Down):
			if len(m.items)-1 > m.cursor {
				m.cursor += 1
			}

		case key.Matches(msg, keys.KeyMap.Select):
			log.Printf("got item '%v'", m.items[m.cursor])
			switch m.items[m.cursor] {
			case "Resume":
				return m, changeState(stateGame)

			case "New Game":
				newgame = new(game.NewGame())
				return m, changeState(stateGame)

			case "Custom Game":
				return m, changeState(stateCustomGame)

			case "Quit":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m menu) View() tea.View {
	var b strings.Builder

	for pos, i := range m.items {
		if pos == m.cursor {
			b.WriteString(string(txt.ColorAccent))
			b.WriteString("> ")
			b.WriteString(string(txt.ColorReset))
		} else {
			b.WriteString("  ")
		}

		b.WriteString(i)
		b.WriteRune('\n')
	}

	return tea.NewView(
		lipgloss.NewStyle().
			Background(color.Black).
			Margin(2).
			MarginBackground(color.Black).
			Padding(3).
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left, txt.Logo, b.String()),
			),
	)
}

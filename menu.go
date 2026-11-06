package main

import (
	"log"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
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
			"Play",
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
			if len(m.items)-1 >= m.cursor {
				m.cursor += 1
			}

		case key.Matches(msg, keys.KeyMap.Select):
			log.Printf("got item '%v'", m.items[m.cursor])
			switch m.items[m.cursor] {
			case "Play":
				return m, changeState(stateGame)
			case "Quit":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m menu) View() tea.View {
	var b strings.Builder

	b.WriteString(txt.Logo)

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

	return tea.NewView(b.String())
}

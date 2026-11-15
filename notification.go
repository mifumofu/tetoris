package main

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/internal/keys"
)

// Provides a notification model that displays a message in the
// overlay part of the main model, and disappears when dismissed
type notification string

func (m notification) Init() tea.Cmd {
	return nil
}

func (m notification) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	kMsg, ofType := msg.(tea.KeyPressMsg)
	// do nothing if event not keypress and Select not pressed
	if !ofType || !key.Matches(kMsg, keys.KeyMap.Select) {
		return m, nil
	}

	return nil, nil
}

func (m notification) View() tea.View {
	return tea.NewView(
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Width(35).
			Padding(1).
			PaddingLeft(2).
			PaddingRight(2).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Top,
					string(m+"\n"),

					lipgloss.NewStyle().
						Background(lipgloss.Yellow).
						Foreground(lipgloss.Black).
						Render(" ok "),
				)),
	)
}

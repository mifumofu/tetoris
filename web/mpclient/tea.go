package mpclient

import (
	tea "charm.land/bubbletea/v2"
	"github.com/mifumofu/tetoris/game"
)

func (m *MpGame) actionListener(a game.Action) {
	m.Send(a)
}

func (m MpGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case StateLobby:
		return m.LobbyUpdate(msg)
	case StatePlaying:
		a, b := m.game.Update(msg)

		c, ok := a.(game.Game)
		if !ok {
			panic("Game did not return a game model")
		}

		m.game = c
		return m, b
	}

	return m, nil
}

func (m MpGame) Init() tea.Cmd {
	return m.receive()
}

func (m MpGame) View() tea.View {
	switch m.state {
	case StateLobby:
		return m.LobbyView()
	case StatePlaying:
		return m.game.View()
	}

	panic("Invalid state for mpclient")
}

package mpclient

import (
	tea "charm.land/bubbletea/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/web/multiplayer"
)

func (m *MpGame) actionListener(a game.Action) {
	m.Send(a)
}

func (m MpGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case multiplayer.StateLobby:
		return m.LobbyUpdate(msg)

	case multiplayer.StatePlaying:
		return m.GameUpdate(msg)
	}

	return m, nil
}

func (m MpGame) Init() tea.Cmd {
	return m.receive()
}

func (m MpGame) View() tea.View {
	switch m.state {
	case multiplayer.StateLobby:
		return m.LobbyView()
	case multiplayer.StatePlaying:
		return m.GameView()
	}

	panic("Invalid state for mpclient")
}

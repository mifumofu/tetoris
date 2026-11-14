package mpclient

import (
	tea "charm.land/bubbletea/v2"
	"github.com/mifumofu/tetoris/game"
)

func (m *MpGame) actionListener(a game.Action) {
	m.Send(a)
}

func (m MpGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.game.Update(msg)
}

func (m MpGame) Init() tea.Cmd {
	return m.game.Init()
}

func (m MpGame) View() tea.View {
	return m.game.View()
}

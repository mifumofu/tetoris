package mpclient

import (
	"encoding/json"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/web/multiplayer"
)

func (m MpGame) LobbyUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case wsUpdate:
		var mpg multiplayer.GameMpUpdate
		err := json.Unmarshal(msg, &mpg)
		if err != nil {
			panic(err)
		}

		m.game = mpg.Game
		m.selfId = mpg.YourID

		m.players = make(map[multiplayer.PlayerID]game.Game)
		for _, i := range mpg.Players {
			m.players[i] = m.game
		}

		if mpg.Start {
			m, scmd := m.Start()
			return m, tea.Batch(scmd, m.receive())
		}

		return m, m.receive()
	}

	return m, nil
}

func (m MpGame) LobbyView() tea.View {
	var b strings.Builder

	b.WriteString("# Multiplayer game")
	fmt.Fprintf(&b, "\nJoin code: %v", m.id)
	b.WriteString("\n\n> Settings")

	fmt.Fprintf(&b, "\nSeed: %v", m.game.Seed)
	fmt.Fprintf(&b, "\nWidth: %v", m.game.Width)
	fmt.Fprintf(&b, "\nHeight: %v", m.game.Height)
	fmt.Fprintf(&b, "\nKill line: %v", m.game.KillLine)

	b.WriteString("\n\n> Players")
	fmt.Fprintf(&b, "\nYou are %v.", m.selfId)
	for x := range m.players {
		fmt.Fprintf(&b, "\n- %v", x.String())
	}
	fmt.Fprintf(&b, "\n\nWaiting for %v more players to join", 2-len(m.players))

	return tea.NewView(b.String())
}

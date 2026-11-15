package mpclient

import (
	"encoding/json"
	"fmt"
	"strings"
	"uuid"

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

		m.players = make(map[uuid.UUID]game.Game)
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
	b.WriteString("\n> Settings")

	fmt.Fprintf(&b, "\nSeed: %v", m.game.Seed)
	fmt.Fprintf(&b, "\nWidth: %v", m.game.Width)
	fmt.Fprintf(&b, "\nHeight: %v", m.game.Height)
	fmt.Fprintf(&b, "\nKill line: %v", m.game.KillLine)

	b.WriteString("\n\n> Players")
	for x := range m.players {
		fmt.Fprintf(&b, "\n- %v", x.String())
	}

	return tea.NewView(b.String())
}

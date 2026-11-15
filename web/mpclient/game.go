package mpclient

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/internal/txt"
	"github.com/mifumofu/tetoris/web/multiplayer"
)

func (m MpGame) Start() (MpGame, tea.Cmd) {
	m.state = multiplayer.StatePlaying

	m.game = game.NewGame(
		game.WithBlocks(m.game.Shapes),
		game.WithWidth(m.game.Width),
		game.WithHeight(m.game.Height),
		game.WithSeed(m.game.Seed),
		game.WithActionListener(m.actionListener),
	)

	for i := range m.players {
		m.players[i] = game.NewGame(
			game.WithBlocks(m.game.Shapes),
			game.WithWidth(m.game.Width),
			game.WithHeight(m.game.Height),
			game.WithSeed(m.game.Seed),
		)
	}

	return m, m.game.Init()
}

func (m MpGame) GameUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case wsUpdate:
		up := decodeAction(msg)
		g := m.players[up.PlayerID]
		g.Do(up.Action)
		m.players[up.PlayerID] = g

		log.Printf("received message! %v", up.Action.String())

		return m, m.receive()

	default:
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

func (m MpGame) GameView() tea.View {
	var games = []string{
		m.sideBar(),
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Render(m.game.GameView()),
	}
	for i, g := range m.players {
		if i == m.selfId {
			continue
		}
		games = append(games, i.String()+"\n"+g.GameView())
	}

	return tea.NewView(
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			games...,
		),
	)
}

func (m MpGame) sideBar() string {
	var b strings.Builder
	b.WriteString(txt.Logo)
	b.WriteString("\nMultiplayer        \nCode: ")
	b.WriteString(string(m.id))
	return b.String()
}

// encodes an action for sending over a websocket
// returns encoding in bytes
func decodeAction(input []byte) multiplayer.MpUpdate {
	var buf bytes.Buffer
	buf.Write(input)
	out := decodeActionReader(&buf)
	return out
}

// encodes an action for sending over a websocket
// uses an io.Writer
func decodeActionReader(r io.Reader) multiplayer.MpUpdate {
	dec := gob.NewDecoder(r)
	var out multiplayer.MpUpdate
	dec.Decode(&out)
	return out
}

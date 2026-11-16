package multiplayer

import (
	"encoding/json"
	"log/slog"

	"github.com/mifumofu/tetoris/game"
)

type GameMpUpdate struct {
	Players []PlayerID
	Game    game.Game
	Start   bool
	YourID  PlayerID
}

// start a game
func (g *MpGame) Start() {
	g.state = StatePlaying
	g.sendLobbyUpdate(true)
}

// updates lobby on state
func (g *MpGame) sendLobbyUpdate(start bool) {
	for p := range g.Players {
		b, err := json.Marshal(GameMpUpdate{
			Game: g.InitialGame,
			Players: func() []PlayerID {
				var ids []PlayerID
				for i := range g.Players {
					ids = append(ids, i.ID)
				}
				return ids
			}(),
			Start:  start,
			YourID: p.ID,
		})

		if err != nil {
			slog.Warn("failed to encode lobby update", "err", err)
			return
		}

		// send message
		p.send <- b
	}
}

// handles websocket messages when state == lobby
func handleMpLobbyMessages(msg []byte, c client, g *MpGame) bool {
	return true
}

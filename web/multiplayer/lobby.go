package multiplayer

import (
	"encoding/json"
	"log/slog"
	"uuid"

	"github.com/mifumofu/tetoris/game"
)

type GameMpUpdate struct {
	Players []uuid.UUID
	Game    game.Game
	Start   bool
	YourID  uuid.UUID
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
			Players: func() []uuid.UUID {
				var ids []uuid.UUID
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

// encodes a lobby update
func (g *MpGame) encodeLobbyUpdate(start bool) ([]byte, error) {
	return json.Marshal(GameMpUpdate{
		Game: g.InitialGame,
		Players: func() []uuid.UUID {
			var ids []uuid.UUID
			for i := range g.Players {
				ids = append(ids, i.ID)
			}
			return ids
		}(),
		Start: start,
	})
}

// handles websocket messages when state == lobby
func handleMpLobbyMessages(msg []byte, c client, g *MpGame) bool {
	return true
}

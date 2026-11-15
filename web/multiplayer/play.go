package multiplayer

import (
	"bytes"
	"encoding/gob"
	"io"
	"strconv"
	"uuid"

	"github.com/gorilla/websocket"
	"github.com/mifumofu/tetoris/game"
)

// sending data
type MpUpdate struct {
	PlayerID uuid.UUID
	Action   game.Action
}

// sends an MpUpdate to all players in a game
func (g *MpGame) sendAllUpdate(id uuid.UUID, action game.Action) {
	g.sendAll(websocket.BinaryMessage, encodeAction(id, action))
}

// encodes an action for sending over a websocket
// returns encoding in bytes
func encodeAction(id uuid.UUID, action game.Action) []byte {
	var buf bytes.Buffer
	encodeActionWriter(&buf, id, action)
	return buf.Bytes()
}

// encodes an action for sending over a websocket
// uses an io.Writer
func encodeActionWriter(w io.Writer, id uuid.UUID, action game.Action) {
	enc := gob.NewEncoder(w)
	enc.Encode(MpUpdate{id, action})
}

// handles websocket messages when state == playing
func handleMpPlayingMsgs(msg []byte, c client, g *MpGame) bool {
	a, err := parseMpPlayMsg(msg)
	if err != nil {
		c.logger.Warn("failed to parse message", "err", err)
		g.DisconnectClient(c)
		return false
	}

	c.logger.Debug("msg", "action", a.String())

	// Broadcast message to all clients
	g.sendAllUpdate(c.ID, a)

	return true
}

// clients should be sending the int of the game.Actions enum action
// they are performing in their game
func parseMpPlayMsg(b []byte) (game.Action, error) {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		return game.ActionDown, ErrInvalidMessage
	}

	a := game.Action(i)

	// check if valid
	// TODO: find a better way to see if action is valid
	if i > int(game.ActionInstadown) || i < 0 {
		return a, ErrInvalidMessage
	}

	return a, nil
}

package replay

import (
	"log"
	"slices"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/mifumofu/tetoris/game"
	"github.com/mifumofu/tetoris/internal/keys"
)

type replay struct {
	game    game.Game
	orig    game.Game // original game passed in
	actions []game.Action
	cursor  int
}

func NewReplay(g game.Game) replay {
	log.Println("actions in save", g.Actions)
	return replay{
		actions: slices.Clone(g.Actions),
		orig:    g,
		game: game.NewGame(
			game.WithBlocks(g.Shapes),
			game.WithWidth(g.Width),
			game.WithHeight(g.Height),
			game.WithSeed(g.Seed),
		),
	}
}

func (m *replay) reset() {
	m.game = game.NewGame(
		game.WithBlocks(m.orig.Shapes),
		game.WithWidth(m.orig.Width),
		game.WithHeight(m.orig.Height),
		game.WithSeed(m.orig.Seed),
	)

}

func (m replay) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m replay) View() tea.View {
	return tea.NewView(m.game.GameView())
}

func (m replay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyPressMsg:

		switch {
		case key.Matches(msg, keys.KeyMap.Up):
			m.reset()
			m.cursor = 0
			return m, nil

		case key.Matches(msg, keys.KeyMap.Down):
			m.reset()
			for _, x := range m.actions {
				m.game.Do(x)
			}
			m.cursor = len(m.actions) - 1
			return m, nil

		case key.Matches(msg, keys.KeyMap.Left):
			if m.cursor <= 0 {
				return m, nil
			}

			m.reset()
			m.cursor -= 1
			// right now, it currently takes around half a frame to replay
			// an entire game, so even though this is inefficient, it should
			// be good enough since we don't have a proper way to reverse actions
			// which is important mainly for ActionDown since that has the possibliity
			// to reset the cursor and because of the random number gen
			for x := range m.cursor {
				m.game.Do(m.actions[x])
			}
			return m, nil

		case key.Matches(msg, keys.KeyMap.Right):
			if m.cursor >= len(m.actions)-1 {
				return m, nil
			}
			m.game.Do(m.actions[m.cursor])
			m.cursor += 1
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

package game

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/mifumofu/tetoris/internal/keys"
)

//
// update
//

// handle given update events
func (m Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m.Move(position{y: 1})
		return m, m.doTick()

	// Is it a key press?
	case tea.KeyPressMsg:

		switch {
		case key.Matches(msg, keys.KeyMap.Up):
			if m.state == stateRunning {
				rot, canrotate := m.active.form.Rotate(), true
				// don't allow the player to rotate the block into a wall
				for _, b := range rot.BlockPos() {
					if b.Add(m.cursor).x >= m.Width || b.Add(m.cursor).y >= m.Height {
						canrotate = false
						break
					}
				}

				if canrotate {
					m.active.form = rot
				}
			}

		case key.Matches(msg, keys.KeyMap.Down):
			m.Move(position{y: 1})

		case key.Matches(msg, keys.KeyMap.Left):
			m.Move(position{x: -1})

		case key.Matches(msg, keys.KeyMap.Right):
			m.Move(position{x: 1})

		case key.Matches(msg, keys.KeyMap.Instadown):
			m.cursor = m.FinalPosition(m.active.form)
			m.Move(position{y: 1})

		case key.Matches(msg, keys.KeyMap.Pause):
			if m.state == stateRunning {
				m.state = statePaused
			} else if m.state == statePaused {
				m.state = stateRunning
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

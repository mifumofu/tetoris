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
		m.Do(ActionDown)
		return m, m.doTick()

	// Is it a key press?
	case tea.KeyPressMsg:

		switch {
		case key.Matches(msg, keys.KeyMap.Up):
			m.Do(ActionRotate)

		case key.Matches(msg, keys.KeyMap.Down):
			m.Do(ActionDown)

		case key.Matches(msg, keys.KeyMap.Left):
			m.Do(ActionLeft)

		case key.Matches(msg, keys.KeyMap.Right):
			m.Do(ActionRight)

		case key.Matches(msg, keys.KeyMap.Instadown):
			m.Do(ActionInstadown)
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

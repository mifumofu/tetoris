package game

import (
	"fmt"
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Game and helper functions
type Game struct {
	score   int
	state   state
	actions []position // list of all moves made
	blocks  map[position]block

	rand *rand.Rand
	seed int64

	killLine int // y of kill line
	width    int
	height   int

	cursor position
	active shape
}

func NewGame(seed *int64) Game {
	m := Game{
		state: stateRunning,
		seed:  rand.Int63(),

		blocks: make(map[position]block),

		killLine: 2,
		width:    13,
		height:   25,

		cursor: position{x: 5, y: 0},
	}

	// use set seed if requested
	if seed != nil {
		m.seed = *seed
	}

	m.rand = rand.New(rand.NewSource(m.seed))
	m.active = m.NewShape()
	return m
}

//
// init
//

func (m Game) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return doTick()
}

type tickMsg time.Time

func doTick() tea.Cmd {
	fmt.Println("tick")
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

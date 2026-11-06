package game

import (
	"log"
	"math/rand"
	"time"

	tea "charm.land/bubbletea/v2"
)

// Game and helper functions
type Game struct {
	score       int
	rowsCleared int
	state       state
	actions     []position // list of all moves made
	blocks      map[position]block

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
	return m.doTick()
}

type tickMsg time.Time

func (m Game) doTick() tea.Cmd {
	log.Println("got tick for falling block!")
	nextTick := max(time.Second-(time.Millisecond*20*time.Duration(m.rowsCleared)), 150*time.Millisecond)
	log.Printf("making block fall! next in %v", nextTick)
	return tea.Tick(nextTick, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

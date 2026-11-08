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
	Actions     []Action // list of all moves made
	blocks      map[position]block

	rand *rand.Rand
	Seed int64

	KillLine int // y of kill line
	Width    int
	Height   int
	Shapes   []ShapeForm

	cursor position
	active Shape
}

func NewGame(opts ...gameOption) Game {
	m := Game{
		state: stateRunning,
		Seed:  rand.Int63(),

		blocks: make(map[position]block),

		KillLine: 2,
		Width:    13,
		Height:   25,
		Shapes:   shapesForms,

		cursor: position{X: 5, Y: 0},
	}

	// apply options
	for _, opt := range opts {
		opt(&m)
	}

	m.rand = rand.New(rand.NewSource(m.Seed))
	m.active = m.NewShape()

	return m
}

//
// options
//

type gameOption func(*Game)

func WithSeed(w int64) gameOption {
	return func(s *Game) {
		s.Seed = w
		s.rand = rand.New(rand.NewSource(s.Seed))
	}
}

func WithWidth(w int) gameOption {
	return func(s *Game) {
		s.Width = w
	}
}

func WithHeight(h int) gameOption {
	return func(s *Game) {
		s.Height = h
	}
}

func WithBlocks(sf []ShapeForm) gameOption {
	return func(s *Game) {
		s.Shapes = sf
	}
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

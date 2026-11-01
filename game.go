package main

import (
	"math/rand"
	"slices"

	"charm.land/bubbles/v2/help"
)

// model and helper functions
type model struct {
	score   int
	state   state
	actions []position // list of all moves made
	rand    *rand.Rand
	seed    int64

	termHeight int
	termWidth  int

	keys keyMap
	help help.Model

	blocks   map[position]block
	killLine int // y of kill line
	width    int
	height   int

	cursor position
	active shape
}

func initialModel() model {
	m := model{
		state: statePaused,
		keys:  keys,
		help:  help.New(),
		seed:  rand.Int63(),

		blocks: make(map[position]block),

		killLine: 2,
		width:    13,
		height:   25,

		cursor: position{x: 5, y: 0},
	}
	m.rand = rand.New(rand.NewSource(m.seed))
	m.active = m.NewShape()
	m.help.ShowAll = true
	return m
}

type state int8

const (
	stateRunning state = iota
	statePaused
	stateEnd
)

func (m model) Get(p position) (block, bool) {
	// block does not exist if its out of bounds
	if p.x >= m.width || p.x < 0 || p.y >= m.height || p.y < 0 {
		return block{}, true
	}

	b, ok := m.blocks[p]
	if b.color == "" {
		return b, false
	}
	return b, ok
}

func (m *model) Set(p position, b block) {
	p.In(m.width, m.height)
	m.blocks[p] = b
}

func (m *model) ResetCursor() {
	m.SetShape(m.cursor, m.active)
	m.cursor = position{x: m.width/2 - 1, y: 0}
	m.active = m.NewShape()

	if m.FinalPosition(m.active.form).y <= m.killLine {
		m.state = stateEnd
	}

	m.score += 5
}

// handles moving a block, resets position if necessary
func (m *model) Move(d position) {
	if m.state != stateRunning {
		return
	}

	if !m.CanMove(d) {
		if d.y > 0 && d.x == 0 {
			cur := m.cursor
			form := m.active.form.BlockPos()
			m.ResetCursor()

			for _, b := range form {
				if m.CanClearRow(cur.Add(b).y) {
					m.ClearRow(cur.Add(b).y)
				}
			}
		}
		return
	}

	(*m).actions = append((*m).actions, d)
	(*m).cursor.y += d.y
	(*m).cursor.x += d.x
}

// checks if s can move the difference of p in m
// only supports one direction at a time
func (m model) CanMove(d position) bool {
	bpos := m.active.form.BlockPos()
	for _, b := range bpos {
		npos := m.cursor.Add(b).Add(d)
		if _, exists := m.Get(npos); exists {
			return false

		} else if npos.x < 0 || npos.x == m.width {
			return false

		} else if npos.y < 0 || npos.y == m.height {
			return false
		}
	}
	return true
}

// checks if a row can be cleared
func (m model) CanClearRow(y int) bool {
	for px := range m.width {
		pos := position{x: px, y: y}
		b, exists := m.Get(pos)
		_ = b
		if pos.x == 0 {
			continue
		}
		if !exists {
			return false
		}
	}
	return true
}

// clears a row, does not check if it is completely full
func (m model) ClearRow(y int) {
	rangeover := m.height - (m.height - y)
	for py := range rangeover {
		if rangeover-py == m.height {
			continue
		}
		for px := range m.width {
			pos := position{y: rangeover - py, x: px}
			block, exists := m.Get(position{x: pos.x, y: pos.y - 1})

			if exists {
				m.Set(pos, block)
			} else {
				delete(m.blocks, pos)
			}
		}
	}

	m.score += 50
}

// position
type position struct {
	x int
	y int
}

func (p position) Add(d position) position {
	return position{x: p.x + d.x, y: p.y + d.y}
}

func (p position) In(x, y int) bool {
	if !(p.x > 0 && p.x <= x) {
		return false
	}

	if !(p.y > 0 && p.y <= y) {
		return false
	}
	return true
}

func (p position) Equal(a position) bool {
	return p.x == a.x && p.y == a.y
}

// block
type block struct {
	color     color
	character rune
}

// shapes
type shape struct {
	form      shapeForm
	character rune
	color     color
}

type shapeForm [][]bool

// generates a new random shape
func (m model) NewShape() shape {
	return shape{
		color:     colors[m.rand.Intn(len(colors))],
		form:      randomShapeForm(shapesForms, m.rand),
		character: ' ',
	}
}

func randomShapeForm(sp []shapeForm, seed *rand.Rand) shapeForm {
	s := sp[seed.Intn(len(sp))]

	// randomly flip shape form
	if seed.Int()%2 == 0 {
		s = s.Flip()
	}

	// randomly flop shape form
	if seed.Int()%2 == 0 {
		s = s.Flop()
	}

	return s
}

// omg slice overload
var shapesForms = []shapeForm{
	{
		{true, true, true, true},
	},
	{
		{true, false, false},
		{true, true, true},
	},
	{
		{false, false, true},
		{true, true, true},
	},
	{
		{true, true},
		{true, true},
	},
	{
		{false, true, true},
		{true, true, false},
	},
	{
		{false, true, false},
		{true, true, true},
	},
}

// rotates a shape's form 90deg
func (s shapeForm) Rotate() shapeForm {
	var n shapeForm
	// transpose the slice, move rows to columns
	// we assume all rows are the same length
	for pa, a := range s {
		for pb, b := range a {
			if pa == 0 {
				n = append(n, []bool{b})
			} else {
				n[pb] = append(n[pb], b)
			}
		}
	}

	// rotate al
	for p := range n {
		slices.Reverse(n[p])
	}

	return n
}

// flips a shape's form horizontally
func (s shapeForm) Flip() shapeForm {
	for _, v := range s {
		slices.Reverse(v)
	}
	return s
}

// flips a shape's form vertically
func (s shapeForm) Flop() shapeForm {
	slices.Reverse(s)
	return s
}

// gets the position of all blocks in a shape, relative to
// the shape itself. for position on grid combine with p.Add
func (s shapeForm) BlockPos() []position {
	var pos []position
	for y, xs := range s {
		for x, exists := range xs {
			if exists == false {
				continue
			}

			pos = append(pos, position{y: y, x: x})
		}
	}
	return pos
}

// converts a shape into blocks on a grid
func (m model) SetShape(p position, s shape) {
	bpos := s.form.BlockPos()
	for _, b := range bpos {
		m.Set(b.Add(p), block{color: s.color, character: s.character})
	}
}

// the position the shape will end up at if the
// user doesn't input anything
func (m model) FinalPosition(s shapeForm) position {
	bpos := s.BlockPos()
	rangeover := m.height - m.cursor.y

	for y := range rangeover {
		var canuse bool = true
		for _, b := range bpos {
			p := b.Add(position{y: m.height - y, x: m.cursor.x})

			if _, exists := m.Get(p); exists {
				canuse = false
				break
			}
		}

		if canuse {
			return position{y: m.height - y, x: m.cursor.x}
		}
	}

	return m.cursor
}

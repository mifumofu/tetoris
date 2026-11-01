package main

import (
	"math/rand"
	"slices"
)

// model and helper functions
type model struct {
	blockChar string
	blocks    map[position]block
	width     int
	height    int

	cursor position
	active shape
}

func initialModel() model {
	return model{
		blockChar: "  ",
		blocks:    make(map[position]block),

		width:  30,
		height: 40,

		cursor: position{x: 5, y: 0},
		active: NewShape(),
	}
}

func (m model) Get(p position) (block, bool) {
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
	m.cursor = position{x: 0, y: 0}
	m.active = NewShape()
}

// handles moving a block, resets position if necessary
func (m *model) Move(d position) {
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
			checkb := rangeover - py
			fromBlock, _ := m.Get(position{y: checkb - 1, x: px})
			m.Set(position{y: checkb, x: px}, fromBlock)
		}
	}
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
	color color
}

// shapes
type shape struct {
	form  shapeForm
	color color
}

type shapeForm [][]bool

// generates a new random shape
func NewShape() shape {
	return shape{
		color: colors[rand.Intn(len(colors))],
		form:  randomShapeForm(shapesForms),
	}
}

func randomShapeForm(sp []shapeForm) shapeForm {
	s := sp[rand.Intn(len(sp))]

	// randomly flip shape form
	if rand.Int()%2 == 0 {
		s = s.Flip()
	}

	// randomly flop shape form
	if rand.Int()%2 == 0 {
		s = s.Flop()
	}

	return s
}

// omg slice overload
var shapesForms = []shapeForm{
	{
		{false, true},
		{false, true},
		{true, true},
	},
	{
		{true},
	},
	{
		{true, true, true},
	},
	{
		{true, true},
		{true, true},
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
		m.Set(b.Add(p), block{color: s.color})
	}
}

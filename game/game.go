package game

import (
	"github.com/mifumofu/tetoris/internal/txt"
)

type state int8

const (
	stateRunning state = iota
	statePaused
	stateEnd
)

func (m Game) Get(p position) (block, bool) {
	// block does not exist if its out of bounds
	if p.x >= m.Width || p.x < 0 || p.y >= m.Height || p.y < 0 {
		return block{}, true
	}

	b, ok := m.blocks[p]
	if b.color == "" {
		return b, false
	}
	return b, ok
}

func (m *Game) Set(p position, b block) {
	p.In(m.Width, m.Height)
	m.blocks[p] = b
}

// place down the active shape into permanent blocks, then
// resets the cursor with a new shape
func (m *Game) ResetCursor() {
	m.SetShape(m.cursor, m.active)
	m.cursor = position{x: m.Width/2 - 1, y: 0}
	m.active = m.NewShape()

	if m.FinalPosition(m.active.form).y <= m.killLine {
		m.state = stateEnd
	}

	m.score += 5
}

// handles moving a block, resets position if necessary
func (m *Game) Move(d position) {
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
func (m Game) CanMove(d position) bool {
	bpos := m.active.form.BlockPos()
	for _, b := range bpos {
		npos := m.cursor.Add(b).Add(d)
		if _, exists := m.Get(npos); exists {
			return false

		} else if npos.x < 0 || npos.x == m.Width {
			return false

		} else if npos.y < 0 || npos.y == m.Height {
			return false
		}
	}
	return true
}

// checks if a row can be cleared
func (m Game) CanClearRow(y int) bool {
	for px := range m.Width {
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
func (m *Game) ClearRow(y int) {
	rangeover := m.Height - (m.Height - y)
	for py := range rangeover {
		if rangeover-py == m.Height {
			continue
		}
		for px := range m.Width {
			pos := position{y: rangeover - py, x: px}
			block, exists := m.Get(position{x: pos.x, y: pos.y - 1})

			if exists {
				m.Set(pos, block)
			} else {
				delete(m.blocks, pos)
			}
		}
	}

	(*m).score += 50
	(*m).rowsCleared += 1
}

// position
type position struct {
	x int
	y int
}

func (p position) Add(d position) position {
	return position{x: p.x + d.x, y: p.y + d.y}
}

// finds if p is in the boundaries of X, Y
func (p position) In(x, y int) bool {
	if !(p.x > 0 && p.x <= x) {
		return false
	}

	if !(p.y > 0 && p.y <= y) {
		return false
	}
	return true
}

// checks if a position is equal to another
func (p position) Equal(a position) bool {
	return p.x == a.x && p.y == a.y
}

// block
type block struct {
	color     txt.Color
	character rune
}

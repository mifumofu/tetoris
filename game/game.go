package game

import (
	"os"
	"path"

	"github.com/mifumofu/tetoris/internal/txt"
)

type state int8

const (
	stateRunning state = iota
	stateEnd
)

func (m Game) Get(p position) (block, bool) {
	// block does not exist if its out of bounds
	if p.X >= m.Width || p.X < 0 || p.Y >= m.Height || p.Y < 0 {
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
	m.cursor = position{X: m.Width/2 - 1, Y: 0}
	m.active = m.NewShape()

	if m.FinalPosition(m.active.form).Y <= m.KillLine {
		m.state = stateEnd
	}

	m.score += 5
}

func (m *Game) Do(a Action) {
	if m.state != stateRunning {
		return
	}
	switch a {
	case ActionRotate:
		m.Rotate()

	case ActionDown:
		m.Move(position{Y: 1})

	case ActionLeft:
		m.Move(position{X: -1})

	case ActionRight:
		m.Move(position{X: 1})

	case ActionInstadown:
		m.cursor = m.FinalPosition(m.active.form)
		m.Move(position{Y: 1})

	}
	(*m).Actions = append(m.Actions, a)
}

// handles moving a block, resets position if necessary
func (m *Game) Move(d position) {
	if m.state != stateRunning {
		return
	}

	if !m.CanMove(d) {
		if d.Y > 0 && d.X == 0 {
			cur := m.cursor
			form := m.active.form.BlockPos()
			m.ResetCursor()

			for _, b := range form {
				if m.CanClearRow(cur.Add(b).Y) {
					m.ClearRow(cur.Add(b).Y)
				}
			}
		}
		return
	}

	(*m).cursor.Y += d.Y
	(*m).cursor.X += d.X
}

// checks if s can move the difference of p in m
// only supports one direction at a time
func (m Game) CanMove(d position) bool {
	bpos := m.active.form.BlockPos()
	for _, b := range bpos {
		npos := m.cursor.Add(b).Add(d)
		if _, exists := m.Get(npos); exists {
			return false

		} else if npos.X < 0 || npos.X == m.Width {
			return false

		} else if npos.Y < 0 || npos.Y == m.Height {
			return false
		}
	}
	return true
}

// checks if a row can be cleared
func (m Game) CanClearRow(y int) bool {
	for px := range m.Width {
		pos := position{X: px, Y: y}
		b, exists := m.Get(pos)
		_ = b
		if pos.X == 0 {
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
			pos := position{Y: rangeover - py, X: px}
			block, exists := m.Get(position{X: pos.X, Y: pos.Y - 1})

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
	X int
	Y int
}

func (p position) Add(d position) position {
	return position{X: p.X + d.X, Y: p.Y + d.Y}
}

// finds if p is in the boundaries of X, Y
func (p position) In(x, y int) bool {
	if !(p.X > 0 && p.X <= x) {
		return false
	}

	if !(p.Y > 0 && p.Y <= y) {
		return false
	}
	return true
}

// checks if a position is equal to another
func (p position) Equal(a position) bool {
	return p.X == a.X && p.Y == a.Y
}

// block
type block struct {
	color     txt.Color
	character rune
}

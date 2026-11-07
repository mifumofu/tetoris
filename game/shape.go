package game

import (
	"math/rand"
	"slices"

	"github.com/mifumofu/tetoris/internal/txt"
)

// shapes
type Shape struct {
	form      ShapeForm
	character rune
	color     txt.Color
}

type ShapeForm [][]bool

//
// new shapes
//

// generates a new random shape
func (m Game) NewShape() Shape {
	return Shape{
		color:     txt.Colors[m.rand.Intn(len(txt.Colors))],
		form:      randomShapeForm(m.Shapes, m.rand),
		character: ' ',
	}
}

func randomShapeForm(sp []ShapeForm, seed *rand.Rand) ShapeForm {
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
var shapesForms = []ShapeForm{
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

//
// manipulating form
//

// rotates a shape's form 90deg
func (s ShapeForm) Rotate() ShapeForm {
	var n ShapeForm
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
func (s ShapeForm) Flip() ShapeForm {
	for _, v := range s {
		slices.Reverse(v)
	}
	return s
}

// flips a shape's form vertically
func (s ShapeForm) Flop() ShapeForm {
	slices.Reverse(s)
	return s
}

//
// misc
//

// gets the position of all blocks in a shape, relative to
// the shape itself. for position on grid combine with p.Add
func (s ShapeForm) BlockPos() []position {
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
func (m Game) SetShape(p position, s Shape) {
	bpos := s.form.BlockPos()
	for _, b := range bpos {
		m.Set(b.Add(p), block{color: s.color, character: s.character})
	}
}

// the position the shape will end up at if the
// user doesn't input anything
func (m Game) FinalPosition(s ShapeForm) position {
	bpos := s.BlockPos()
	rangeover := m.Height - m.cursor.y - len(m.active.form) + 1
	var prev position

	for y := range rangeover {
		var canuse bool = false
		for _, b := range bpos {
			p := b.Add(position{y: m.cursor.y + y, x: m.cursor.x})

			if _, exists := m.Get(p); exists {
				canuse = true
				break
			}
		}

		if canuse {
			return prev
		}

		prev = position{y: m.cursor.y + y, x: m.cursor.x}
	}

	return prev
}

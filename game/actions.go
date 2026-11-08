package game

type Action int16

const (
	ActionRotate Action = iota
	ActionDown
	ActionLeft
	ActionRight
	ActionInstadown
)

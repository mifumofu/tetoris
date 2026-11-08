package game

type Action int16

const (
	ActionRotate Action = iota
	ActionDown
	ActionLeft
	ActionRight
	ActionInstadown
)

func (a Action) String() string {
	switch a {
	case ActionRotate:
		return "↷ Rotated"

	case ActionDown:
		return "↓ Down"

	case ActionLeft:
		return "← Move Left"

	case ActionRight:
		return "→ Move Right"

	case ActionInstadown:
		return "↡ Instadown"

	default:
		return ""
	}
}

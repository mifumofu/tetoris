package main

import (
	"strconv"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/mifumofu/tetoris/game"
)

type customGame struct {
	form *huh.Form
}

func newCustomGame() customGame {
	g := game.NewGame()

	return customGame{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Key("seed").
					Value(new(strconv.Itoa(int(g.Seed)))).
					Validate(func(s string) error {
						_, err := strconv.Atoi(s)
						return err
					}).
					Title("Seed"),

				huh.NewInput().
					Key("width").
					Value(new(strconv.Itoa(int(g.Width)))).
					Validate(func(s string) error {
						_, err := strconv.Atoi(s)
						return err
					}).
					Title("Width"),

				huh.NewInput().
					Key("height").
					Value(new(strconv.Itoa(int(g.Height)))).
					Validate(func(s string) error {
						_, err := strconv.Atoi(s)
						return err
					}).
					Title("Height"),
			),
		),
	}
}

func (m customGame) Init() tea.Cmd {
	return m.form.Init()
}

func (m customGame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	form, cmd := m.form.Update(msg)
	f, ok := form.(*huh.Form)
	if ok && f.State != huh.StateCompleted {
		return m, cmd
	}

	m.form = f

	// width
	seed, err := strconv.Atoi(f.GetString("seed"))
	if err != nil {
		panic(err)
	}

	// width
	w, err := strconv.Atoi(f.GetString("width"))
	if err != nil {
		panic(err)
	}

	// height
	h, err := strconv.Atoi(f.GetString("height"))
	if err != nil {
		panic(err)
	}

	newgame = new(game.NewGame(
		game.WithSeed(seed),
		game.WithWidth(w),
		game.WithHeight(h),
	))

	return m, changeState(stateGame)
}

func (m customGame) View() tea.View {
	return tea.NewView(m.form.View())
}

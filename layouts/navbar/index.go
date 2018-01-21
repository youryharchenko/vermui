package navbar

import (
	"github.com/verdverm/vermui/hoc/console"
	"github.com/verdverm/vermui/layouts/grid"
)

type NavBar struct {
	*grid.Grid
	console *console.DevConsoleWidget
	usererr *console.ErrorConsoleWidget
}

func New() *NavBar {
	rows := []*grid.Row{}

	ue := console.NewErrorConsoleWidget()
	ue.Init()
	ueRow := grid.NewRow(
		grid.NewCol(12, 0, ue),
	)
	rows = append(rows, ueRow)

	cw := console.NewDevConsoleWidget()
	cw.Init()
	cwRow := grid.NewRow(
		grid.NewCol(12, 0, cw),
	)
	rows = append(rows, cwRow)

	g := grid.NewGrid(rows...)

	nb := &NavBar{
		Grid: g,
		// console: cw,
		usererr: ue,
	}

	return nb
}

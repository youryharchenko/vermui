package main

import (
	"github.com/verdverm/vermui/hoc/cmdbox"
	"github.com/verdverm/vermui/hoc/console"
	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/layouts/grid"
	"github.com/verdverm/vermui/layouts/navbar"
	"github.com/verdverm/vermui/layouts/router"
	"github.com/verdverm/vermui/lib/mux"
	"github.com/verdverm/vermui/lib/render"
	"github.com/verdverm/vermui/widgets/text"
)

func buildLayout() layouts.Layout {
	n := text.NewPar(":PRESS Ctrl-c to quit demo")
	n.Height = 3
	n.TextFgColor = render.ColorWhite
	n.BorderLabel = " VermUI - Not Found "
	n.BorderFg = render.ColorCyan

	p := text.NewPar(":PRESS Ctrl-c to quit demo")
	p.Height = 3
	p.TextFgColor = render.ColorWhite
	p.BorderLabel = " VermUI - P "
	p.BorderFg = render.ColorCyan

	q := text.NewPar(":PRESS Ctrl-c to quit demo")
	q.Height = 3
	q.TextFgColor = render.ColorWhite
	q.BorderLabel = " VermUI - Q "
	q.BorderFg = render.ColorCyan

	layoutN := grid.NewGrid()
	layoutN.AddRows(
		grid.NewRow(
			grid.NewCol(6, 0, n),
		),
	)

	layoutP := grid.NewGrid()
	layoutP.AddRows(
		grid.NewRow(
			grid.NewCol(6, 0, p),
		),
	)

	layoutQ := grid.NewGrid()
	layoutQ.AddRows(
		grid.NewRow(
			grid.NewCol(6, 0, q),
		),
	)

	r := router.New()

	cb := cmdbox.New()
	cb.BorderLabel = " VermUI "

	nav := navbar.New()

	cw := console.NewDevConsoleWidget()
	cw.Init()

	g := grid.NewGrid()
	r1 := grid.NewRow(
		grid.NewCol(12, 0, cb),
	)
	navRow := grid.NewRow(
		grid.NewCol(12, 0, nav),
	)
	rg := grid.NewGrid(
		grid.NewRow(
			grid.NewCol(12, 0, cw),
		),
	)
	r2 := grid.NewRow(
		grid.NewCol(12, 0, rg),
	)
	r3 := grid.NewRow(
		grid.NewCol(12, 0, r),
	)
	g.AddRows(r1, navRow, r2, r3)

	r.SetNotFound(layoutN)
	r.AddRouteLayout("/p", layoutP)
	r.AddRouteHandlerFunc("/q/{cnt}", func(req *mux.Request) (layouts.Layout, error) {
		vars := mux.Vars(req)
		q.BorderLabel = " VermUI - Q @" + vars["cnt"] + " "
		return layoutQ, nil
	})

	return g
}

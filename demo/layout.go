package main

import (
	"github.com/verdverm/vermui/hoc/cmdbox"
	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/layouts/grid"
	"github.com/verdverm/vermui/layouts/navbar"
	"github.com/verdverm/vermui/layouts/router"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/mux"
	"github.com/verdverm/vermui/lib/render"
	"github.com/verdverm/vermui/widgets/text"
)

func buildLayout() layouts.Layout {
	// create out NotFound view
	n := text.NewPar(":PRESS Ctrl-c to quit demo")
	n.Height = 3
	n.TextFgColor = render.ColorWhite
	n.BorderLabel = " VermUI - Not Found "
	n.BorderFg = render.ColorCyan

	// Wrap the views in grids
	layoutN := grid.NewGrid()
	layoutN.AddRows(
		grid.NewRow(
			grid.NewCol(6, 3, n),
		),
	)

	// the Home View
	h := text.NewPar(":PRESS Ctrl-c to quit demo")
	h.Height = 3
	h.TextFgColor = render.ColorWhite
	h.BorderLabel = " VermUI - Home "
	h.BorderFg = render.ColorCyan

	layoutH := grid.NewGrid()
	layoutH.AddRows(
		grid.NewRow(
			grid.NewCol(6, 3, h),
		),
	)

	// the Echo View
	q := text.NewPar(":PRESS Ctrl-c to quit demo")
	q.Height = 3
	q.TextFgColor = render.ColorWhite
	q.BorderLabel = " VermUI - Echo "
	q.BorderFg = render.ColorCyan

	layoutQ := grid.NewGrid()
	layoutQ.AddRows(
		grid.NewRow(
			grid.NewCol(6, 3, q),
		),
	)

	// The layouts we just made are our main views
	// They will be added to the router below

	// Lets create some other goodies
	// and the main layout next

	// Setup a Command Box
	cbox := cmdbox.New()
	cbox.BorderLabel = " VermUI "
	cboxRow := grid.NewRow(
		grid.NewCol(12, 0, cbox),
	)

	// NavBar has the Console and UserError HOCs
	nav := navbar.New()

	// The navbar is actually hidden
	// C-l and C-e to see the console/user-error
	hiddenRow := grid.NewRow(
		grid.NewCol(12, 0, nav),
	)

	// Now to tie everything together

	// First,
	//   Create a new Router View
	//   This uses an internal router based on gorilla/mux
	rtr := router.New()

	// Set a NotFound View (aka 404 w/o the internet)
	rtr.SetNotFound(layoutN)

	// Add a layout directly
	rtr.AddRouteLayout("/home", layoutH)

	// Add a route with a Handler function
	rtr.AddRouteHandlerFunc("/echo/{what}", func(req *mux.Request) (layouts.Layout, error) {
		vars := mux.Vars(req)
		q.BorderLabel = " VermUI - Echo - '" + vars["what"] + "' "
		return layoutQ, nil
	})
	// There is also a Handler type for more complex requirements

	// Add a command to the commandbox
	cbox.AddCommandCallback("echo", func(args []string, context map[string]interface{}) {
		go events.SendCustomEvent("/router/dispatch", "/echo/"+args[0])
	})

	// Then,
	//   Grid to use the router view as our main view
	mainRow := grid.NewRow(
		grid.NewCol(12, 0, rtr),
	)

	// Finally,
	//   we will use a grid as our top-most layout
	g := grid.NewGrid()
	g.AddRows(cboxRow, hiddenRow, mainRow)
	return g
}

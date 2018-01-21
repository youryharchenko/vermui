package modules

import (
	"github.com/verdverm/vermui/hoc/cmdbox"
	"github.com/verdverm/vermui/layouts/grid"
	"github.com/verdverm/vermui/layouts/navbar"
	"github.com/verdverm/vermui/layouts/router"
)

type RootView struct {
	*grid.Grid

	cmdbox *cmdbox.CmdBoxWidget
	nav    *navbar.NavBar
	view   *router.Router
}

func NewRootView() *RootView {
	cb := cmdbox.New()
	cb.BorderLabel = " Edsger "

	nav := navbar.New()
	view := router.New()

	rows := []*grid.Row{}

	RV := &RootView{
		cmdbox: cb,
		nav:    nav,
		view:   view,
	}

	cbRow := grid.NewRow(
		grid.NewCol(6, 0, RV.cmdbox),
	)
	navRow := grid.NewRow(
		grid.NewCol(12, 0, RV.nav),
	)
	viewRow := grid.NewRow(
		grid.NewCol(12, 0, RV.view),
	)
	rows = append(rows, cbRow)
	rows = append(rows, navRow)
	rows = append(rows, viewRow)

	g := grid.NewGrid(rows...)
	RV.Grid = g

	RV.addRoutes()
	RV.addCommands()

	return RV
}

func (RV *RootView) Name() string {
	return "root-view"
}

func (RV *RootView) addRoutes() {
	for _, mod := range enabledModules {
		r, ok := mod.(router.Routable)
		if ok {
			for _, pair := range r.Routings() {
				RV.view.AddRoute(pair.Path, pair.Thing)
			}
		}
	}
}

func (RV *RootView) addCommands() {
	for _, mod := range enabledModules {
		cmd, ok := mod.(cmdbox.Command)
		if ok {
			RV.cmdbox.Add(cmd)
		}
	}
}

func (RV *RootView) removeCommands() {
	for _, mod := range enabledModules {
		cmd, ok := mod.(cmdbox.Command)
		if ok {
			RV.cmdbox.Remove(cmd)
		}
	}
}

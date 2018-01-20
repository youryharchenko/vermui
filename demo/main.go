package main

import (
	"fmt"

	"github.com/verdverm/vermui"

	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/mux"
	"github.com/verdverm/vermui/lib/render"

	"github.com/verdverm/vermui/layouts/grid"
	"github.com/verdverm/vermui/layouts/router"

	"github.com/verdverm/vermui/widgets/text"
)

func main() {
	fmt.Println("Welcome to VermUI")

	err := vermui.Init()
	if err != nil {
		panic(err)
	}

	n := text.NewPar(":PRESS Ctrl-c to quit demo")
	n.Height = 3
	n.Width = 50
	n.TextFgColor = render.ColorWhite
	n.BorderLabel = " VermUI - Not Found "
	n.BorderFg = render.ColorCyan

	p := text.NewPar(":PRESS Ctrl-c to quit demo")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = render.ColorWhite
	p.BorderLabel = " VermUI - P "
	p.BorderFg = render.ColorCyan

	q := text.NewPar(":PRESS Ctrl-c to quit demo")
	q.Height = 3
	q.Width = 50
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
	r.SetNotFound(layoutN)
	r.AddRouteLayout("/p", layoutP)
	r.AddRouteHandlerFunc("/q/{cnt}", func(req *mux.Request) (layouts.Layout, error) {
		vars := mux.Vars(req)
		q.BorderLabel = " VermUI - Q @" + vars["cnt"] + " "
		return layoutQ, nil
	})
	r.SetActive("/p")

	vermui.SetLayout(r)

	// lib.Render(p) // feel free to call Render, it's async and non-block

	// handle key C-c pressing
	vermui.AddGlobalHandler("/sys/kbd/C-c", func(events.Event) {
		// press q to quit
		vermui.Stop()
	})
	vermui.AddGlobalHandler("/sys/kbd/p", func(events.Event) {
		// press p to show p
		go events.SendCustomEvent("/router/dispatch", "/p")
	})
	cnt := 0
	vermui.AddGlobalHandler("/sys/kbd/q", func(events.Event) {
		// press q to show q
		cnt += 1
		path := fmt.Sprintf("/q/%d", cnt)
		go events.SendCustomEvent("/router/dispatch", path)
	})

	vermui.Start() // block until Stop is called
}

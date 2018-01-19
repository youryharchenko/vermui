package main

import (
	"fmt"

	"github.com/verdverm/vermui"

	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/render"

	"github.com/verdverm/vermui/layouts/grid"

	"github.com/verdverm/vermui/widgets/text"
)

func main() {
	fmt.Println("Welcome to VermUI")

	err := vermui.Init()
	if err != nil {
		panic(err)
	}

	p := text.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = render.ColorWhite
	p.BorderLabel = "VermUI"
	p.BorderFg = render.ColorCyan

	q := text.NewPar(":PRESS q TO QUIT DEMO")
	q.Height = 3
	q.Width = 50
	q.TextFgColor = render.ColorWhite
	q.BorderLabel = "VermUI - Q"
	q.BorderFg = render.ColorCyan

	layout := grid.NewGrid()
	layout.AddRows(
		grid.NewRow(
			grid.NewCol(6, 0, p),
		),
		grid.NewRow(
			grid.NewCol(6, 0, q),
		),
	)

	vermui.SetLayout(layout)

	// lib.Render(p) // feel free to call Render, it's async and non-block

	// handle key q pressing
	vermui.AddGlobalHandler("/sys/kbd/q", func(events.Event) {
		// press q to quit
		vermui.Stop()
	})

	vermui.Start() // block until Stop is called
}

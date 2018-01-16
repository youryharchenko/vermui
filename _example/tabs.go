// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"github.com/verdverm/vermui"
)

func main() {
	err := vermui.Init()
	if err != nil {
		panic(err)
	}
	defer vermui.Close()

	//vermui.UseTheme("helloworld")

	header := vermui.NewPar("Press q to quit, Press j or k to switch tabs")
	header.Height = 1
	header.Width = 50
	header.Border = false
	header.TextBgColor = vermui.ColorBlue

	tab1 := vermui.NewTab("pierwszy")
	par2 := vermui.NewPar("Press q to quit\nPress j or k to switch tabs\n")
	par2.Height = 5
	par2.Width = 37
	par2.Y = 0
	par2.BorderLabel = "Keys"
	par2.BorderFg = vermui.ColorYellow
	tab1.AddBlocks(par2)

	tab2 := vermui.NewTab("drugi")
	bc := vermui.NewBarChart()
	data := []int{3, 2, 5, 3, 9, 5, 3, 2, 5, 8, 3, 2, 4, 5, 3, 2, 5, 7, 5, 3, 2, 6, 7, 4, 6, 3, 6, 7, 8, 3, 6, 4, 5, 3, 2, 4, 6, 4, 8, 5, 9, 4, 3, 6, 5, 3, 6}
	bclabels := []string{"S0", "S1", "S2", "S3", "S4", "S5"}
	bc.BorderLabel = "Bar Chart"
	bc.Data = data
	bc.Width = 26
	bc.Height = 10
	bc.DataLabels = bclabels
	bc.TextColor = vermui.ColorGreen
	bc.BarColor = vermui.ColorRed
	bc.NumColor = vermui.ColorYellow
	tab2.AddBlocks(bc)

	tab3 := vermui.NewTab("trzeci")
	tab4 := vermui.NewTab("żółw")
	tab5 := vermui.NewTab("four")
	tab6 := vermui.NewTab("five")

	tabpane := vermui.NewTabpane()
	tabpane.Y = 1
	tabpane.Width = 30
	tabpane.Border = true

	tabpane.SetTabs(*tab1, *tab2, *tab3, *tab4, *tab5, *tab6)

	vermui.Render(header, tabpane)

	vermui.Handle("/sys/kbd/q", func(vermui.Event) {
		vermui.StopLoop()
	})

	vermui.Handle("/sys/kbd/j", func(vermui.Event) {
		tabpane.SetActiveLeft()
		vermui.Clear()
		vermui.Render(header, tabpane)
	})

	vermui.Handle("/sys/kbd/k", func(vermui.Event) {
		tabpane.SetActiveRight()
		vermui.Clear()
		vermui.Render(header, tabpane)
	})

	vermui.Loop()
}

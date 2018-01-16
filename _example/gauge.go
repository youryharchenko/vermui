// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import "github.com/verdverm/vermui"

func main() {
	err := vermui.Init()
	if err != nil {
		panic(err)
	}
	defer vermui.Close()

	//vermui.UseTheme("helloworld")

	g0 := vermui.NewGauge()
	g0.Percent = 40
	g0.Width = 50
	g0.Height = 3
	g0.BorderLabel = "Slim Gauge"
	g0.BarColor = vermui.ColorRed
	g0.BorderFg = vermui.ColorWhite
	g0.BorderLabelFg = vermui.ColorCyan

	gg := vermui.NewBlock()
	gg.Width = 50
	gg.Height = 5
	gg.Y = 12
	gg.BorderLabel = "TEST"
	gg.Align()

	g2 := vermui.NewGauge()
	g2.Percent = 60
	g2.Width = 50
	g2.Height = 3
	g2.PercentColor = vermui.ColorBlue
	g2.Y = 3
	g2.BorderLabel = "Slim Gauge"
	g2.BarColor = vermui.ColorYellow
	g2.BorderFg = vermui.ColorWhite

	g1 := vermui.NewGauge()
	g1.Percent = 30
	g1.Width = 50
	g1.Height = 5
	g1.Y = 6
	g1.BorderLabel = "Big Gauge"
	g1.PercentColor = vermui.ColorYellow
	g1.BarColor = vermui.ColorGreen
	g1.BorderFg = vermui.ColorWhite
	g1.BorderLabelFg = vermui.ColorMagenta

	g3 := vermui.NewGauge()
	g3.Percent = 50
	g3.Width = 50
	g3.Height = 3
	g3.Y = 11
	g3.BorderLabel = "Gauge with custom label"
	g3.Label = "{{percent}}% (100MBs free)"
	g3.LabelAlign = vermui.AlignRight

	g4 := vermui.NewGauge()
	g4.Percent = 50
	g4.Width = 50
	g4.Height = 3
	g4.Y = 14
	g4.BorderLabel = "Gauge"
	g4.Label = "Gauge with custom highlighted label"
	g4.PercentColor = vermui.ColorYellow
	g4.BarColor = vermui.ColorGreen
	g4.PercentColorHighlighted = vermui.ColorBlack

	vermui.Render(g0, g1, g2, g3, g4)

	vermui.Handle("/sys/kbd/q", func(vermui.Event) {
		vermui.StopLoop()
	})

	vermui.Loop()
}

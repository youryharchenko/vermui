// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"math"

	"github.com/verdverm/vermui"
)

func main() {
	err := vermui.Init()
	if err != nil {
		panic(err)
	}
	defer vermui.Close()

	//vermui.UseTheme("helloworld")

	sinps := (func() []float64 {
		n := 220
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5)
		}
		return ps
	})()

	lc0 := vermui.NewLineChart()
	lc0.BorderLabel = "braille-mode Line Chart"
	lc0.Data = sinps
	lc0.Width = 50
	lc0.Height = 12
	lc0.X = 0
	lc0.Y = 0
	lc0.AxesColor = vermui.ColorWhite
	lc0.LineColor = vermui.ColorGreen | vermui.AttrBold

	lc1 := vermui.NewLineChart()
	lc1.BorderLabel = "dot-mode Line Chart"
	lc1.Mode = "dot"
	lc1.Data = sinps
	lc1.Width = 26
	lc1.Height = 12
	lc1.X = 51
	lc1.DotStyle = '+'
	lc1.AxesColor = vermui.ColorWhite
	lc1.LineColor = vermui.ColorYellow | vermui.AttrBold

	lc2 := vermui.NewLineChart()
	lc2.BorderLabel = "dot-mode Line Chart"
	lc2.Mode = "dot"
	lc2.Data = sinps[4:]
	lc2.Width = 77
	lc2.Height = 16
	lc2.X = 0
	lc2.Y = 12
	lc2.AxesColor = vermui.ColorWhite
	lc2.LineColor = vermui.ColorCyan | vermui.AttrBold

	vermui.Render(lc0, lc1, lc2)
	vermui.Handle("/sys/kbd/q", func(vermui.Event) {
		vermui.StopLoop()
	})
	vermui.Loop()

}

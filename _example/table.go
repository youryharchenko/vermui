// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import "github.com/verdverm/vermui"

func main() {
	err := vermui.Init()
	if err != nil {
		panic(err)
	}
	defer vermui.Close()
	rows1 := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"你好吗", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "10", "11"},
	}

	table1 := vermui.NewTable()
	table1.Rows = rows1
	table1.FgColor = vermui.ColorWhite
	table1.BgColor = vermui.ColorDefault
	table1.Y = 0
	table1.X = 0
	table1.Width = 62
	table1.Height = 7

	vermui.Render(table1)

	rows2 := [][]string{
		[]string{"header1", "header2", "header3"},
		[]string{"Foundations", "Go-lang is so cool", "Im working on Ruby"},
		[]string{"2016", "11", "11"},
	}

	table2 := vermui.NewTable()
	table2.Rows = rows2
	table2.FgColor = vermui.ColorWhite
	table2.BgColor = vermui.ColorDefault
	table2.TextAlign = vermui.AlignCenter
	table2.Separator = false
	table2.Analysis()
	table2.SetSize()
	table2.BgColors[2] = vermui.ColorRed
	table2.Y = 10
	table2.X = 0
	table2.Border = true

	vermui.Render(table2)
	vermui.Handle("/sys/kbd/q", func(vermui.Event) {
		vermui.StopLoop()
	})
	vermui.Loop()
}

// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Copyright 2018 Tony Worm <verdverm@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package grid

import (
	"sync"

	"github.com/verdverm/vermui/lib/render"
)

// Grid implements 12 columns system.
// A simple example:
/*
   import ui "github.com/verdverm/vermui"
   // init and create widgets...

   // build
   ui.Body.AddRows(
       ui.NewRow(
           ui.NewCol(6, 0, widget0),
           ui.NewCol(6, 0, widget1)),
       ui.NewRow(
           ui.NewCol(3, 0, widget2),
           ui.NewCol(3, 0, widget30, widget31, widget32),
           ui.NewCol(6, 0, widget4)))

   // calculate layout
   ui.Body.Align()

   ui.Render(ui.Body)
*/
type Grid struct {
	sync.Mutex
	Rows    []*Row
	Width   int
	Height  int
	X       int
	Y       int
	BgColor render.Attribute
}

// NewGrid returns *Grid with given rows.
func NewGrid(rows ...*Row) *Grid {
	return &Grid{Rows: rows}
}

// AddRows appends given rows to Grid.
func (g *Grid) AddRows(rs ...*Row) {
	g.Rows = append(g.Rows, rs...)
}

// NewRow creates a new row out of given columns.
func NewRow(cols ...*Row) *Row {
	rs := &Row{Span: 12, Cols: cols}
	return rs
}

// NewCol accepts: widgets are LayoutBufferer or widgets is A NewRow.
// Note that if multiple widgets are provided, they will stack up in the col.
func NewCol(span, offset int, widgets ...GridBufferer) *Row {
	r := &Row{Span: span, Offset: offset}

	if widgets != nil && len(widgets) == 1 {

		wgt := widgets[0]
		// go events.SendCustomEvent("/console/debug", fmt.Sprintf("Col widget type: %q", reflect.TypeOf(wgt)))
		nw, isRow := wgt.(*Row)
		if isRow {
			r.Cols = nw.Cols
		} else {
			r.Widget = wgt
		}
		return r
	}

	r.Cols = []*Row{}
	ir := r
	for _, w := range widgets {
		nr := &Row{Span: 12, Widget: w}
		ir.Cols = []*Row{nr}
		ir = nr
	}

	return r
}

func (g *Grid) GetX() int {
	return g.X
}

func (g *Grid) SetX(x int) {
	//g.Lock()
	//defer g.Unlock()
	g.X = x
}

func (g *Grid) GetY() int {
	return g.Y
}

func (g *Grid) SetY(y int) {
	//g.Lock()
	//defer g.Unlock()
	g.Y = y
}

func (g *Grid) GetHeight() int {
	// g.Align()
	return g.Height
}

func (g *Grid) SetHeight(h int) {
	// g.Height = h
}

func (g *Grid) GetWidth() int {
	return g.Width
}

func (g *Grid) SetWidth(w int) {
	//g.Lock()
	//defer g.Unlock()
	g.Width = w
}

// Align calculate each rows' layout.
func (g *Grid) Align() {
	g.calcLayout()
	g.Lock()
	defer g.Unlock()
	for _, r := range g.Rows {
		r.Align()
	}
}

func (g *Grid) calcLayout() {

	g.Lock()
	defer g.Unlock()
	h := 0
	for _, r := range g.Rows {
		r.SetWidth(g.Width)
		r.SetX(g.X)
		r.SetY(g.Y + h)
		r.calcLayout()
		rh := r.GetHeight()
		h += rh
	}

	g.Height = h
}

// Buffer implements Bufferer interface.
func (g *Grid) Buffer() render.Buffer {
	g.Lock()
	defer g.Unlock()
	buf := render.NewBuffer()

	for _, r := range g.Rows {
		buf.Merge(r.Buffer())
	}
	return buf
}

func (g *Grid) Mount() error {
	for _, row := range g.Rows {
		row.Mount()
	}
	return nil
}
func (g *Grid) Unmount() error {
	for _, row := range g.Rows {
		row.Unmount()
	}
	return nil
}

func (g *Grid) Show() {}
func (g *Grid) Hide() {}

func (g *Grid) Focus()   {}
func (g *Grid) Unfocus() {}

// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Copyright 2018 Tony Worm <verdverm@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package grid

import (
	"sync"

	"github.com/verdverm/vermui/lib/render"
)

// GridBufferer introduces a Bufferer that can be manipulated by Grid.
type GridBufferer interface {
	render.Bufferer
	GetHeight() int
	SetWidth(int)
	SetX(int)
	SetY(int)
	Align()
	Mount() error
	Unmount() error
}

// Row builds a layout tree
type Row struct {
	sync.Mutex
	Cols   []*Row       //children
	Widget GridBufferer // root
	X      int
	Y      int
	Width  int
	Height int
	Span   int
	Offset int
}

func (r *Row) Align() {
	r.Lock()
	defer r.Unlock()

	if r.isRenderableLeaf() {
		/*
			switch t := r.Widget.(type) {
			case *Grid:
				go events.SendCustomEvent("/console/debug", fmt.Sprintf("row widget A-1 type: %q", t))
			default:
				if t != nil {
					go events.SendCustomEvent("/console/debug", fmt.Sprintf("row widget A-2 type: %q", t))
				}
			}
		*/

		r.Widget.Align()
	}

	if r.Widget != nil {
		/*
			switch t := r.Widget.(type) {
			case *Grid:
				go events.SendCustomEvent("/console/debug", fmt.Sprintf("row widget B-1 type: %q", t))
			default:
				if t != nil {
					go events.SendCustomEvent("/console/debug", fmt.Sprintf("row widget B-2 type: %q", t))
				}
			}
		*/

		r.Widget.Align()
	}
	if !r.isLeaf() {
		for _, c := range r.Cols {
			c.Align()
		}
	}
}

// calculate and set the underlying layout tree's x, y, height and width.
func (r *Row) calcLayout() {

	r.assignWidth(r.Width)
	r.Height = r.solveHeight()
	r.assignX(r.X)
	r.assignY(r.Y)

}

// tell if the node is leaf in the tree.
func (r *Row) isLeaf() bool {
	return r.Cols == nil || len(r.Cols) == 0
}

func (r *Row) isRenderableLeaf() bool {
	return r.isLeaf() && r.Widget != nil
}

// assign widgets' (and their parent rows') width recursively.
func (r *Row) assignWidth(w int) {
	r.SetWidth(w)

	/*
		r.Lock()
		defer r.Unlock()
	*/

	accW := 0                            // acc span and offset
	calcW := make([]int, len(r.Cols))    // calculated width
	calcOftX := make([]int, len(r.Cols)) // computed start position of x

	for i, c := range r.Cols {
		accW += c.Span + c.Offset
		cw := int(float64(c.Span*r.Width) / 12.0)

		if i >= 1 {
			calcOftX[i] = calcOftX[i-1] +
				calcW[i-1] +
				int(float64(r.Cols[i-1].Offset*r.Width)/12.0)
		}

		// use up the space if it is the last col
		if i == len(r.Cols)-1 && accW == 12 {
			cw = r.Width - calcOftX[i]
		}
		calcW[i] = cw
		r.Cols[i].assignWidth(cw)
	}
}

// bottom up calc and set rows' (and their widgets') height,
// return r's total height.
func (r *Row) solveHeight() int {
	if r.isRenderableLeaf() {
		r.Widget.Align()
		r.Height = r.Widget.GetHeight()
		return r.Height
	}

	maxh := 0
	if !r.isLeaf() {
		for _, c := range r.Cols {
			nh := c.solveHeight()
			// when embed rows in Cols, row widgets stack up
			if r.Widget != nil {
				nh += r.Widget.GetHeight()
			}
			if nh > maxh {
				maxh = nh
			}
		}
	}

	r.Height = maxh
	return maxh
}

// recursively assign x position for r tree.
func (r *Row) assignX(x int) {
	r.SetX(x)

	/*
		r.Lock()
		defer r.Unlock()
	*/

	if !r.isLeaf() {
		acc := 0
		for i, c := range r.Cols {
			if c.Offset != 0 {
				acc += int(float64(c.Offset*r.Width) / 12.0)
			}
			r.Cols[i].assignX(x + acc)
			acc += c.Width
		}
	}
}

// recursively assign y position to r.
func (r *Row) assignY(y int) {
	r.SetY(y)

	/*
		r.Lock()
		defer r.Unlock()
	*/

	if r.isLeaf() {
		return
	}

	for i := range r.Cols {
		acc := 0
		if r.Widget != nil {
			acc = r.Widget.GetHeight()
		}
		r.Cols[i].assignY(y + acc)
	}

}

// GetHeight implements GridBufferer interface.
func (r Row) GetHeight() int {
	//r.Lock()
	//defer r.Unlock()

	return r.Height
}

func (r *Row) GetX() int {
	r.Lock()
	defer r.Unlock()

	return r.X
}

// SetX implements GridBufferer interface.
func (r *Row) SetX(x int) {
	r.Lock()
	defer r.Unlock()

	r.X = x
	if r.Widget != nil {
		r.Widget.SetX(x)
	}
}
func (r *Row) GetY() int {
	r.Lock()
	defer r.Unlock()

	return r.Y
}

// SetY implements GridBufferer interface.
func (r *Row) SetY(y int) {
	r.Lock()
	defer r.Unlock()

	r.Y = y
	if r.Widget != nil {
		r.Widget.SetY(y)
	}
}

// SetWidth implements GridBufferer interface.
func (r *Row) SetWidth(w int) {
	r.Lock()
	defer r.Unlock()

	r.Width = w
	if r.Widget != nil {
		r.Widget.SetWidth(w)
	}
}

// Buffer implements Bufferer interface,
// recursively merge all widgets buffer
func (r *Row) Buffer() render.Buffer {
	r.Lock()
	defer r.Unlock()

	merged := render.NewBuffer()

	if r.isRenderableLeaf() {
		return r.Widget.Buffer()
	}

	// for those are not leaves but have a renderable widget
	if r.Widget != nil {
		merged.Merge(r.Widget.Buffer())
	}

	// collect buffer from children
	if !r.isLeaf() {
		for _, c := range r.Cols {
			merged.Merge(c.Buffer())
		}
	}

	return merged
}

func (r *Row) Mount() error {
	if r.isRenderableLeaf() {
		return r.Widget.Mount()
	}

	if r.Widget != nil {
		r.Widget.Mount()
	}
	if !r.isLeaf() {
		for _, c := range r.Cols {
			c.Mount()
		}
	}
	return nil
}

func (r *Row) Unmount() error {
	if r.isRenderableLeaf() {
		return r.Widget.Unmount()
	}

	if r.Widget != nil {
		r.Widget.Unmount()
	}
	if !r.isLeaf() {
		for _, c := range r.Cols {
			c.Unmount()
		}
	}
	return nil
}

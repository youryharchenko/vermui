package base

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/render"
)

type Layout struct {
	*render.Block

	Parent   layouts.Layout
	Children []layouts.Layout
}

func New() *Layout {
	l := &Layout{
		Block:    render.NewBlock(),
		Children: []layouts.Layout{},
	}

	return l
}

// Buffer implements Bufferer interface.
func (L *Layout) Buffer() render.Buffer {
	L.Align()
	buf := render.NewBuffer()
	for _, c := range L.Children {
		buf.Merge(c.Buffer())
	}
	return buf
}

func (L *Layout) SetX(x int) {
	L.Block.SetX(x)
	for _, c := range L.Children {
		c.SetX(x)
	}
}
func (L *Layout) SetY(y int) {
	L.Block.SetY(y)
	h := y
	for _, c := range L.Children {
		c.SetY(h)
		h += c.GetHeight()
	}
}

func (L *Layout) GetHeight() int {
	h := 0
	for _, c := range L.Children {
		h += c.GetHeight()
	}
	return h
}

func (L *Layout) SetWidth(w int) {
	L.Block.Width = 0
	for _, c := range L.Children {
		c.SetWidth(w)
	}
}

func (L *Layout) Align() {
	L.Block.Align()
	for _, c := range L.Children {
		c.Align()
	}
}

func (L *Layout) Mount() error {
	fmt.Println("layouts.Base.Mount")
	for _, c := range L.Children {
		err := c.Mount()
		if err != nil {
			errors.Wrapf(err, "in Layout.Mount()")
		}
	}
	return nil
}

func (L *Layout) Unmount() error {
	for _, c := range L.Children {
		err := c.Unmount()
		if err != nil {
			errors.Wrapf(err, "in Layout.Unmount()")
		}
	}
	return nil
}

func (L *Layout) Show() {
	for _, c := range L.Children {
		c.Show()
	}
}

func (L *Layout) Hide() {
	for _, c := range L.Children {
		c.Hide()
	}
}

func (L *Layout) Focus() {
}
func (L *Layout) Unfocus() {
}

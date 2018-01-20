package base

import (
	"github.com/pkg/errors"

	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/render"
)

type Layout struct {
	width    int
	Children []layouts.Layout
}

func New() *Layout {
	l := &Layout{
		Children: []layouts.Layout{},
	}

	return l
}

// Buffer implements Bufferer interface.
func (L *Layout) Buffer() render.Buffer {
	buf := render.NewBuffer()
	for _, c := range L.Children {
		buf.Merge(c.Buffer())
	}
	return buf
}

func (L *Layout) GetWidth() int {
	return L.width
}
func (L *Layout) SetWidth(w int) {
	L.width = w
}

func (L *Layout) Align() {
	return
	for _, c := range L.Children {
		c.Align()
	}
}

func (L *Layout) Mount() error {
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

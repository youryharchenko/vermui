package base

import (
	ui "github.com/verdverm/vermui"

	"github.com/verdverm/vermui/layouts"
)

type Layout struct {
	Children layouts.Layout
}

func New() *Layout {
	l := &Layout{
		rows: []*ui.Row{},
	}

	return l
}

// Buffer implements Bufferer interface.
func (L *Layout) Buffer() Buffer {
	buf := NewBuffer()

	for _, r := range L.rows {
		buf.Merge(r.Buffer())
	}
	return buf
}

func (L *Layout) Rows() []*ui.Row {
	rows := []*ui.Row{}

	if L.Header != nil {
		rows = append(rows, L.Header.Rows()...)
	}

	if L.Body != nil {
		rows = append(rows, L.Body.Rows()...)
	}

	if L.Footer != nil {
		rows = append(rows, L.Footer.Rows()...)
	}

	L.rows = rows
	return L.rows
}

func (L *Layout) Mount() error {
	if L.Header != nil {
		err := L.Header.Mount()
		if err != nil {
			return err
		}
	}

	if L.Body != nil {
		err := L.Body.Mount()
		if err != nil {
			return err
		}
	}

	if L.Footer != nil {
		err := L.Footer.Mount()
		if err != nil {
			return err
		}
	}

	L.Rows()

	return nil
}

func (L *Layout) Unmount() error {
	if L.Header != nil {
		err := L.Header.Unmount()
		if err != nil {
			return err
		}
	}

	if L.Body != nil {
		err := L.Body.Unmount()
		if err != nil {
			return err
		}
	}

	if L.Footer != nil {
		err := L.Footer.Unmount()
		if err != nil {
			return err
		}
	}

	L.rows = nil

	return nil
}

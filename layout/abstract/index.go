package abstract

import (
	ui "github.com/verdverm/vermui"
)

type Layout interface {
	Name() string

	Rows() []*ui.Row

	Mount() error
	Unmount() error

	Focus()
	Unfocus()
}

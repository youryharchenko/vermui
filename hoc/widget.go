package widgets

import (
	ui "github.com/verdverm/vermui"
)

type Widget interface {
	ui.GridBufferer

	Mount() error
	Unmount() error

	Show()
	Hide()

	Focus()
	Unfocus()
}

type DataWidget interface {
	Widget
	SetData(interface{})
}

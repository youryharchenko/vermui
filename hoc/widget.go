package widgets

import (
	ui "github.com/verdverm/termui"
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

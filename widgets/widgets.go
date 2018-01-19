package widgets

import (
	"sync"

	"github.com/verdverm/vermui/lib/render"
)

type Widget interface {
	render.Bufferer

	Mount() error
	Unmount() error

	Show()
	Hide()

	Focus()
	Unfocus()

	Props() sync.Map
	SetProps(sync.Map)
}

type DataWidget interface {
	Widget
	SetData(interface{})
}

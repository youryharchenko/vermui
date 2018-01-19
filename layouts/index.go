package layouts

import "github.com/verdverm/vermui/lib/render"

type Layout interface {
	render.Bufferer

	GetWidth() int
	SetWidth(int)
	Align()

	Mount() error
	Unmount() error

	Show()
	Hide()

	Focus()
	Unfocus()
}

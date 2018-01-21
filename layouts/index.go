package layouts

import "github.com/verdverm/vermui/lib/render"

type Layout interface {
	render.Bufferer

	GetHeight() int
	SetHeight(int)

	GetWidth() int
	SetWidth(int)

	GetX() int
	SetX(int)

	GetY() int
	SetY(int)

	Align()

	Mount() error
	Unmount() error

	Show()
	Hide()

	Focus()
	Unfocus()
}

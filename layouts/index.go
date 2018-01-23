package layouts

import "github.com/rivo/tview"

type Layout interface {
	tview.Primitive

	// Mount() error
	// Unmount() error
}

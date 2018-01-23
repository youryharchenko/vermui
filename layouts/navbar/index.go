package navbar

import (
	"github.com/rivo/tview"
	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/hoc/console"
	"github.com/verdverm/vermui/lib"
	"github.com/verdverm/vermui/lib/events"
)

type NavBar struct {
	*tview.Flex

	box *tview.Box

	console *console.DevConsoleWidget
	usererr *console.ErrorConsoleWidget
}

func New() *NavBar {
	box := tview.NewBox().SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle(" VermUI ")
	rhs := tview.NewBox().SetBorder(true).SetTitleAlign(tview.AlignLeft).SetTitle(" Home ")

	topRow := tview.NewFlex()
	topRow.AddItem(box, 0, 1, false)
	topRow.AddItem(rhs, 0, 1, false)

	ue := console.NewErrorConsoleWidget()
	ue.Init()
	ueRow := tview.NewFlex().AddItem(ue, 0, 1, false)

	cw := console.NewDevConsoleWidget()
	cw.Init()
	cwRow := tview.NewFlex().AddItem(cw, 0, 1, false)

	f := tview.NewFlex().SetDirection(tview.FlexRow)
	f.AddItem(topRow, 3, 1, false)

	nb := &NavBar{
		Flex: f,

		console: cw,
		usererr: ue,
	}

	showUE := false
	vermui.AddGlobalHandler("/sys/key/C-e", func(ev events.Event) {
		showUE = !showUE
		if showUE {
			f.AddItem(ueRow, 0, 5, false)
		} else {
			for idx, item := range f.GetItems() {
				if item.Item == ueRow {
					f.DelItem(idx)
					break
				}
			}
		}
		lib.Draw()
	})

	showCW := false
	vermui.AddGlobalHandler("/sys/key/C-l", func(ev events.Event) {
		showCW = !showCW
		if showCW {
			f.AddItem(cwRow, 0, 5, false)
		} else {
			for idx, item := range f.GetItems() {
				if item.Item == cwRow {
					f.DelItem(idx)
					break
				}
			}
		}
		lib.Draw()
	})

	return nb
}

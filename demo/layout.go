package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell"
	"github.com/verdverm/tview"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"

	"github.com/verdverm/vermui/hoc/console"
)

const corporate = `Leverage agile frameworks to provide a robust synopsis for high level overviews. Iterative approaches to corporate strategy foster collaborative thinking to further the overall value proposition. Organically grow the holistic world view of disruptive innovation via workplace diversity and empowerment.

Bring to the table win-win survival strategies to ensure proactive domination. At the end of the day, going forward, a new normal that has evolved from generation X is on the runway heading towards a streamlined cloud solution. User generated content in real-time will have multiple touchpoints for offshoring.

Capitalize on low hanging fruit to identify a ballpark value added activity to beta test. Override the digital divide with additional clickthroughs from DevOps. Nanotechnology immersion along the information highway will close the loop on focusing solely on the bottom line.

[yellow]Press Enter, then Tab/Backtab for word selections`

func buildLayout() tview.Primitive {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	topBar := tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	topBar.SetTitle("  [blue]VermUI[white]  ").SetTitleAlign(tview.AlignLeft).SetBorder(true)
	fmt.Fprint(topBar, "A React-like terminal UI framework.")
	layout.AddItem(topBar, 3, 1, false)

	ue := console.NewErrorConsoleWidget()
	ue.Init()
	ueRow := tview.NewFlex().AddItem(ue, 0, 1, false)

	cw := console.NewDevConsoleWidget()
	cw.Init()
	cwRow := tview.NewFlex().AddItem(cw, 0, 1, false)

	showUE := false
	vermui.AddGlobalHandler("/sys/key/C-e", func(ev events.Event) {
		showUE = !showUE
		if showUE {
			layout.InsItem(1, ueRow, 0, 5, false)
		} else {
			for idx, item := range layout.GetItems() {
				if item.Item == ueRow {
					layout.DelItem(idx)
					break
				}
			}
		}
		vermui.Draw()
	})

	showCW := false
	vermui.AddGlobalHandler("/sys/key/C-l", func(ev events.Event) {
		showCW = !showCW
		if showCW {
			layout.InsItem(1, cwRow, 0, 5, false)
		} else {
			for idx, item := range layout.GetItems() {
				if item.Item == cwRow {
					layout.DelItem(idx)
					break
				}
			}
		}
		vermui.Draw()
	})

	home := genTextView(corporate)
	help := tview.NewTextView()
	fmt.Fprint(help, "help, i need some help!")

	items := map[string]tview.Primitive{
		"home": home,
		"Help": help,
	}

	pages := tview.NewPages()
	for name, page := range items {
		key := name[:1]
		pg := name
		events.Handle("/sys/key/"+key, func(e events.Event) {
			pages.SwitchToPage(pg)
			vermui.Draw()
		})

		pages.AddPage(pg, page, true, false)
	}
	pagesRow := tview.NewFlex().AddItem(pages, 0, 1, true)
	pages.SwitchToPage("home")

	layout.AddItem(pagesRow, 0, 5, true)

	return layout
}

func genTextView(text string) tview.Primitive {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			vermui.Draw()
		})

	numSelections := 0
	go func() {
		for _, word := range strings.Split(text, " ") {
			if word == "the" {
				word = "[red]the[white]"

			}
			if word == "to" {
				word = fmt.Sprintf(`["%d"]to[""]`, numSelections)
				numSelections++

			}
			fmt.Fprintf(textView, "%s ", word)
			time.Sleep(200 * time.Millisecond)

		}

	}()
	textView.SetDoneFunc(func(key tcell.Key) {
		currentSelection := textView.GetHighlights()
		if key == tcell.KeyEnter {
			if len(currentSelection) > 0 {
				textView.Highlight()

			} else {
				textView.Highlight("0").ScrollToHighlight()

			}

		} else if len(currentSelection) > 0 {
			index, _ := strconv.Atoi(currentSelection[0])
			if key == tcell.KeyTab {
				index = (index + 1) % numSelections

			} else if key == tcell.KeyBacktab {
				index = (index - 1 + numSelections) % numSelections

			} else {
				return

			}
			textView.Highlight(strconv.Itoa(index)).ScrollToHighlight()

		}

	})
	textView.SetBorder(true)

	return textView
}

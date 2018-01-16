package tabview

import (
	ui "github.com/verdverm/vermui"
)

type TabView struct {
	Name string

	rows       []*ui.Row
	headerTabs []*ui.Row
	tabs       map[string]*Tab
	active     string
	activeTab  *Tab

	// for handlers
	dummy *ui.Block
}

type Tab struct {
	Name   string
	Hotkey string
	Widget hoc.Widget

	header *ui.Par
	row    *ui.Row
}

func NewTabView(name string, ts []*Tab) *TabView {

	TV := &TabView{}

	tabs := map[string]*Tab{}
	cols := []*ui.Row{}
	rows := []*ui.Row{}

	for _, t := range ts {
		tabs[t.Name] = t

		// tab header element
		p := ui.NewPar("")
		p.Border = false
		p.BorderLabel = t.Name
		p.Height = 1
		col := ui.NewCol(1, 0, p)
		cols = append(cols, col)

		t.header = p

		// tab widget element
		if tab.Rows {

		} else {
			tab.Widget.SetVisible(false)
			tab.Widget.Init()
			row := ui.NewRow(ui.NewCol(12, 0, tab.Widget))
			rows = append(rows, row)
		}

		localT := t
		ui.Handle("/sys/kbd/"+tab.Hotkey, func(ev ui.Event) {
			ui.SendCustomEvt("/tabs/activate", localT)
		})

	}

	header := ui.NewRow(cols...)

	viewRows := []*ui.Row{
		header,
	}
	viewRows = append(viewRows, rows...)

	TV := &TabView{
		Tabs:       tabs,
		Rows:       viewRows,
		HeaderTabs: cols,
	}

	TV.setActive(ts[0].Name)

	return t
}

func (TV *TabView) Mount() {

	TV.dummy.Handle("/"+TV.Name+"/activate", func(ev ui.Event) {
		active := ev.Data.(string)
		TV.setActive(active)
	})
}

func (T *TabView) setActive(active string) {
	// do nothing if already active
	if active == T.Active {
		return
	}
	T.Active = active

	// hide deactivating
	if T.ActiveTab != nil {
		T.ActiveTab.Widget.SetVisible(false)
		T.ActiveTab.Header.BorderLabelFg = ui.ColorGreen
	}

	// show
	tab := T.Tabs[active]
	tab.Widget.SetVisible(true)

	// Highlight Tab
	tab.Header.BorderLabelFg = ui.ColorYellow

	// finally, set the new current tab
	T.ActiveTab = tab
	go ui.SendCustomEvt("/sys/redraw", "tabs - "+active)
}

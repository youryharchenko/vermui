package router

import (
	"fmt"

	ui "github.com/verdverm/vermui"

	"github.com/verdverm/vermui/layout/abstract"
	"github.com/verdverm/vermui/layout/base"
)

type Route interface {
	abstract.Layout

	Id() string
	Route() string
	HotKey() string

	Subroutes() []Route

	AddSubroute() []Route
	RemoveSubroute() []Route
}

type BaseRoute struct {
	*base.Layout

	id        string
	name      string
	shortName string
	shortCode string
	hotKey    string
}

type Router struct {
	activeName   string
	activeLayout abstract.Layout

	routes []Route

	// for holding the event handles
	dummy *ui.Block
}

func New(name string) *Layout {
	lv := &Layout{
		name:           name,
		enabledLayouts: make(map[string]abstract.Layout),
		dummy:          ui.NewBlock(),
	}
	return lv
}

func (L *Layout) Name() string {
	return L.name
}

func (L *Layout) Rows() []*ui.Row {
	var rows []*ui.Row

	if L.activeLayout != nil {
		rows = L.activeLayout.Rows()
	}

	L.rows = rows

	return L.rows
}

func (L *Layout) Mount() error {
	L.dummy.Handle(fmt.Sprintf("/%s/activate", L.name), func(ev ui.Event) {
		active := ev.Data.(string)
		L.SetActive(active)
	})

	for _, layout := range L.enabledLayouts {
		if l, ok := layout.(SubLayout); ok {
			if key := l.HotKey(); key != "" {
				name := l.Name()
				// go ui.SendCustomEvt("/dev/messages", "adding hotkey "+key+" for "+name)
				L.dummy.Handle("/sys/kbd/"+key, func(e ui.Event) {
					L.SetActive(name)
				})
			}
		}
	}

	if L.activeLayout != nil {
		return L.activeLayout.Mount()
	}
	return nil
}

func (L *Layout) Unmount() error {
	L.dummy.RemoveHandle(fmt.Sprintf("/%s/activate", L.name))
	for _, layout := range L.enabledLayouts {
		if l, ok := layout.(SubLayout); ok {
			if key := l.HotKey(); key != "" {
				L.dummy.RemoveHandle("/sys/kbd/" + key)
			}
		}
	}
	if L.activeLayout != nil {
		err := L.activeLayout.Unmount()
		if err != nil {
			return err
		}
		L.activeLayout = nil
	}

	return nil
}

func (L *Layout) AddRoute(r Route) error {
	L.enabledLayouts[r.Route()] = layout

	return nil
}

func (L *Layout) AddSubLayouts(subs []abstract.Layout) error {
	for _, layout := range subs {
		L.enabledLayouts[layout.Route()] = layout
	}

	return nil
}

func (L *Layout) AddSubLayouts(subs []abstract.Layout) error {
	for _, layout := range subs {
		L.enabledLayouts[layout.Route()] = layout
	}

	return nil
}

func (L *Layout) SetActive(active string) {
	// go ui.SendCustomEvt("/dev/messages", "switcher receiving request for "+active)

	// do nothing if already active
	if active == L.activeName {
		return
	}

	// go ui.SendCustomEvt("/dev/messages", "switcher setting active to "+active)

	// make sure the layout exists
	layout, ok := L.enabledLayouts[active]
	if !ok {
		go ui.SendCustomEvt("/console/crit", fmt.Sprint("Error! Layout '%s' does not exist"))
		return
	}

	// mount new layout
	layout.Mount()

	// unmount deactivating
	if L.activeLayout != nil {
		L.activeLayout.Unmount()
	}

	// finally, set the active layout and redraw
	L.activeName = active
	L.activeLayout = layout
	L.Rows()
	go ui.SendCustomEvt("/sys/redraw", "switch - "+active)
}

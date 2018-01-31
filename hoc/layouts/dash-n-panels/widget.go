// DashAndPanels Layout is a Flex widget with
// hidable panels and a main content.
// Can be vert or horz oriendted throught the Flex widget.
// Panels can be jumped to with <hotkey> and hidden with <shift>-<hotkey>
// Recommend making the <hotkey> an: '<alt>-<key>' and hidden will be '<shift>-<alt>-<key>'
// Normal movement and interaction keys within the focussed panel.
//
// main (middle) panel, can be anything, including...
// - the router (when this is the root view)
// - another DashAndPanels
// - a pager, grid, or any other primitive
package dashnpanels

import (
	"sort"

	"github.com/verdverm/tview"
	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

type Panel struct {
	Name       string
	Item       tview.Primitive
	FixedSize  int
	Proportion int
	Focus      int
	Hidden     bool
	HotKey     string
}

type Layout struct {
	*tview.Flex

	// first (left/top) panels, can be almost anything and hidden.
	fPanels map[string]Panel

	// main (middle) panel, can be anything, I think.
	mPanel Panel

	// last (right/bottom) panels, can be almost anything and hidden.
	lPanels map[string]Panel
}

func New() *Layout {
	L := &Layout{
		Flex:    tview.NewFlex(),
		fPanels: map[string]Panel{},
		lPanels: map[string]Panel{},
	}

	return L
}

// AddFirstPanel adds a Panel to the left or top, depending on orientation.
func (L *Layout) AddFirstPanel(name string, item tview.Primitive, fixedSize, proportion, focus int, hidden bool, hotkey string) {
	panel := Panel{
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		Hidden:     hidden,
		HotKey:     hotkey,
	}

	L.fPanels[name] = panel
}

// AddLastPanel adds a Panel to the right or bottom, depending on orientation.
func (L *Layout) AddLastPanel(name string, item tview.Primitive, fixedSize, proportion, focus int, hidden bool, hotkey string) {
	panel := Panel{
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		Hidden:     hidden,
		HotKey:     hotkey,
	}

	L.lPanels[name] = panel
}

func (L *Layout) SetMainPanel(name string, item tview.Primitive, fixedSize, proportion, focus int, hotkey string) {
	panel := Panel{
		Name:       name,
		Item:       item,
		FixedSize:  fixedSize,
		Proportion: proportion,
		Focus:      focus,
		HotKey:     hotkey,
	}

	L.mPanel = panel
}

func (L *Layout) Mount(context map[string]interface{}) error {
	err := L.build()
	if err != nil {
		return err
	}

	items := L.GetItems()
	for _, item := range items {
		err := item.Item.Mount(context)
		if err != nil {
			return err
		}
	}

	// Setup hotkeys
	for _, panel := range L.fPanels {
		if panel.HotKey != "" {
			localPanel := panel
			vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.HotKey, func(e events.Event) {
				go events.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
				vermui.SetFocus(localPanel.Item)
			})
		}
	}
	if L.mPanel.HotKey != "" {
		localPanel := L.mPanel
		vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.HotKey, func(e events.Event) {
			go events.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
			vermui.SetFocus(localPanel.Item)
		})
	}
	for _, panel := range L.fPanels {
		if panel.HotKey != "" {
			localPanel := panel
			vermui.AddWidgetHandler(L, "/sys/key/"+localPanel.HotKey, func(e events.Event) {
				go events.SendCustomEvent("/console/trace", "Focus: "+localPanel.Name)
				vermui.SetFocus(localPanel.Item)
			})
		}
	}

	return nil
}

func (L *Layout) build() error {
	// get and order the fPanels
	fPs := []Panel{}
	for _, panel := range L.fPanels {
		if panel.Hidden {
			continue
		}
		fPs = append(fPs, panel)
	}
	sort.Slice(fPs, func(i, j int) bool {
		return fPs[i].Focus < fPs[j].Focus
	})

	// get and order the lPanels
	lPs := []Panel{}
	for _, panel := range L.lPanels {
		if panel.Hidden {
			continue
		}
		lPs = append(lPs, panel)
	}
	sort.Slice(lPs, func(i, j int) bool {
		return lPs[i].Focus < lPs[j].Focus
	})

	// Start a fresh Flex item
	orient := L.GetDirection()
	L.Flex = tview.NewFlex().SetDirection(orient)

	for _, p := range fPs {
		L.AddItem(p.Item, p.FixedSize, p.Proportion, p.Focus >= 0)
	}

	p := L.mPanel
	L.AddItem(p.Item, p.FixedSize, p.Proportion, p.Focus >= 0)

	for _, p := range lPs {
		L.AddItem(p.Item, p.FixedSize, p.Proportion, p.Focus >= 0)
	}

	return nil
}

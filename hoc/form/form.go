package form

import (
	"github.com/gdamore/tcell"
	"github.com/verdverm/tview"
	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

type Form interface {
	tview.Primitive

	Name() string

	GetValues() (values map[string]interface{})
	SetValues(values map[string]interface{})

	AddItem(name string, item FormItem, taborder, proportion int)
	GetItem(name string) FormItem
	GetItems(name string) []FormItem

	AddButton(name string, button FormButton, taborder, proportion int)
	GetButton(name string) FormButton
	GetButtons(name string) []FormButton
}

type FormItem interface {
	tview.Primitive

	Name() string

	GetValues() (values map[string]interface{})
	SetValues(values map[string]interface{})

	SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem
}

type FormButton interface {
	tview.Primitive

	Name() string

	OnSubmit()
	SetBlurFunc(handler func(key tcell.Key)) FormItem
}

type FormBase struct {
	*Flex

	name string

	items   []FormItem
	buttons []FormButton

	focusedElement int

	// An optional function which is called when the user hits Escape.
	cancel func()
}

func New(name string) *FormBase {
	L := &FormBase{
		Flex:    NewFlex(name),
		buttons: []FormButton{},
	}

	return L
}

func (f *FormBase) Focus(delegate func(p tview.Primitive)) {
	items := f.GetItems()

	if len(items)+len(f.buttons) == 0 {
		return
	}

	// Hand on the focus to one of our child elements.
	if f.focusedElement < 0 || f.focusedElement >= len(items)+len(f.buttons) {
		f.focusedElement = 0
	}
	handler := func(key tcell.Key) {
		switch key {
		case tcell.KeyTab, tcell.KeyEnter:
			f.focusedElement++
			f.Focus(delegate)
		case tcell.KeyBacktab:
			f.focusedElement--
			if f.focusedElement < 0 {
				f.focusedElement = len(items) + len(f.buttons) - 1
			}
			f.Focus(delegate)
		case tcell.KeyEscape:
			if f.cancel != nil {
				f.cancel()
			} else {
				f.focusedElement = 0
				f.Focus(delegate)
			}
		}
	}

	if f.focusedElement < len(items) {
		// We're selecting an item.
		item := items[f.focusedElement]
		item.SetFinishedFunc(handler)
		delegate(item)
	} else {
		// We're selecting a button.
		button := f.buttons[f.focusedElement-len(items)]
		button.SetBlurFunc(handler)
		delegate(button)
	}
}

/*
func (F *FormBase) Mount(context map[string]interface{}) error {
	go events.SendCustomEvent("/console/debug", "form mount")

	vermui.AddWidgetHandler(F, "/sys/key/<tab>", func(evt events.Event) {
		F.handleTab()
	})

	vermui.AddWidgetHandler(F, "/sys/key/<backtab>", func(evt events.Event) {
		F.handleBacktab()
	})

	vermui.AddWidgetHandler(F, "/sys/key/<enter>", func(evt events.Event) {
		F.handleEnter()
	})

	vermui.AddWidgetHandler(F, "/sys/key/<esc>", func(evt events.Event) {
		F.handleEscape()
	})

	return nil
}

func (F *FormBase) Unmount() error {

	vermui.RemoveWidgetHandler(F, "/sys/key/<tab>")
	vermui.RemoveWidgetHandler(F, "/sys/key/<backtab>")
	vermui.RemoveWidgetHandler(F, "/sys/key/<enter>")
	vermui.RemoveWidgetHandler(F, "/sys/key/<esc>")

	return nil
}
*/

func (F *FormBase) GetButton(name string) FormButton {
	for _, item := range F.buttons {
		if item.Name() == name {
			return item
		}
	}
	return nil
}

func (F *FormBase) GetButtons(name string) []FormButton {
	return F.buttons
}

func (F *FormBase) AddButton(button FormButton, fixedSize, proportion int) {
	F.buttons = append(F.buttons, button)
	F.Flex.Flex.AddItem(button, fixedSize, proportion, true)
}

func (F *FormBase) handleTab() {
	go events.SendCustomEvent("/console/debug", "form <tab>")

	items := F.GetItems()
	allLen := len(items) + len(F.buttons)

	idx := F.focusedElement
	idx++
	if idx >= allLen {
		idx = 0
	}
	F.focusedElement = idx

	if idx < len(items) {
		vermui.SetFocus(items[idx])
	} else {
		vermui.SetFocus(F.buttons[idx-len(items)])
	}

}
func (F *FormBase) handleBacktab() {
	go events.SendCustomEvent("/console/debug", "form <backtab>")

	items := F.GetItems()
	allLen := len(items) + len(F.buttons)

	idx := F.focusedElement
	idx--
	if idx < 0 {
		idx = allLen - 1
	}
	F.focusedElement = idx

	if idx < len(items) {
		vermui.SetFocus(items[idx])
	} else {
		vermui.SetFocus(F.buttons[idx-len(items)])
	}

}
func (F *FormBase) handleEnter() {
	go events.SendCustomEvent("/console/debug", "form <enter>")

	items := F.GetItems()
	if F.focusedElement < len(items) {
		F.handleTab()
		return
	}

	// otherwise, we've got a button

	idx := F.focusedElement - len(items)

	F.buttons[idx].OnSubmit()

}
func (F *FormBase) handleEscape() {
	go events.SendCustomEvent("/console/debug", "form <esc>")

	items := F.GetItems()
	if F.focusedElement < len(items) {
		item := items[F.focusedElement]
		item.Blur()
	}
	vermui.SetFocus(F)
}

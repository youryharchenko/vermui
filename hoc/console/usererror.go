package console

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/lib"
	"github.com/verdverm/vermui/lib/events"
)

type ErrorConsoleWidget struct {
	*tview.TextView
}

func NewErrorConsoleWidget() *ErrorConsoleWidget {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetChangedFunc(func() {
			lib.Draw()
		})

	textView.SetTitle(" errors ").
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	c := &ErrorConsoleWidget{
		TextView: textView,
	}

	return c
}

func (D *ErrorConsoleWidget) Init() {

	vermui.AddGlobalHandler("/user/error", func(ev events.Event) {
		text := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), ev.Data)
		fmt.Fprint(D, "%s", text)
	})
}

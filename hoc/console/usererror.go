package console

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/verdverm/tview"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

type ErrorConsoleWidget struct {
	*tview.TextView
}

func NewErrorConsoleWidget() *ErrorConsoleWidget {
	textView := tview.NewTextView().
		SetTextColor(tcell.ColorMaroon).
		SetScrollable(true).
		SetChangedFunc(func() {
			vermui.Draw()
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

	vermui.AddGlobalHandler("/user/error", func(evt events.Event) {
		str := evt.Data.(*events.EventCustom).Data()
		text := fmt.Sprintf("[%s] %v - %d\n", evt.When().Format("2006-01-02 15:04:05"), str, vermui.Application().Screen().Colors())
		fmt.Fprintf(D, "%s", text)
	})
}

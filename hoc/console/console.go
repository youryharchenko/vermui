package console

import (
	"fmt"
	"strings"

	"github.com/verdverm/tview"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

type DevConsoleWidget struct {
	*tview.TextView
}

func NewDevConsoleWidget() *DevConsoleWidget {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true).
		SetChangedFunc(func() {
			vermui.Draw()
		})

	textView.SetTitle(" console ").SetBorder(true)

	c := &DevConsoleWidget{
		TextView: textView,
	}

	return c
}

func (D *DevConsoleWidget) Init() {
	// capture all key strokes and print
	vermui.AddGlobalHandler("/console/key", func(ev events.Event) {
		str := ev.Data.(*events.EventCustom).Data()
		fmt.Fprintf(D, "[fuchsia]key %s[white]\n", str)
	})

	vermui.AddGlobalHandler("/sys/err", func(ev events.Event) {
		err := ev.Data.(*events.EventError)
		line := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), err)
		fmt.Fprintf(D, "[red]SYSTEM ERROR %v[white]\n", line)
	})

	vermui.AddGlobalHandler("/console", func(ev events.Event) {
		d := ev.Data
		switch t := ev.Data.(type) {
		case *events.EventCustom:
			d = t.Data()
		case *events.EventKey:
			d = t.KeyStr
		}
		line := fmt.Sprintf("[%s] %v", ev.When().Format("2006-01-02 15:04:05"), d)

		level := strings.TrimPrefix(ev.Path, "/console/")
		if level[:4] == "colo" {
			color := level[6:]
			line = fmt.Sprintf("[%s]%.5s  %s[white]", color, color, line)
		} else {
			switch level {
			case "crit":
				line = fmt.Sprintf("[red]CRIT   %s[white]", line)
			case "error":
				line = fmt.Sprintf("[orange]ERROR  %s[white]", line)
			case "warn":
				line = fmt.Sprintf("[yellow]WARN   %s[white]", line)
			case "info":
				line = fmt.Sprintf("INFO  %s", line)
			case "debug":
				line = fmt.Sprintf("[green]DEBUG  %s[white]", line)
			case "trace":
				line = fmt.Sprintf("[aqua]TRACE  %s[white]", line)
			}
		}

		fmt.Fprintln(D, line)
		D.ScrollToEnd()
	})
}

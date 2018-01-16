package console

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/verdverm/termui"
)

type ErrorConsoleWidget struct {
	Rows []*ui.Row

	content  *ui.Par
	messages []string
}

func NewErrorConsoleWidget() *ErrorConsoleWidget {
	content := ui.NewPar("")

	content.BorderLabel = " errors "
	content.Height = 0
	content.Border = false
	content.BorderFg = ui.ColorRed
	content.BorderLabelFg = ui.ColorWhite

	rows := []*ui.Row{
		ui.NewRow(
			ui.NewCol(12, 0, content),
		),
	}

	return &ErrorConsoleWidget{
		Rows:     rows,
		content:  content,
		messages: []string{},
	}
}

func (D *ErrorConsoleWidget) Init() {

	ui.Handle("/sys/kbd/C-e", func(ev ui.Event) {
		if D.content.Height > 0 {
			D.content.Height = 0
			D.content.Border = false
		} else {
			D.content.Height = 24
			D.content.Border = true
			D.UpdateText()
		}
	})

	ui.Handle("/user/error", func(ev ui.Event) {
		text := fmt.Sprintf("[%s] %v", time.Unix(ev.Time, 0).Format("2006-01-02 15:04:05"), ev.Data)
		lines := strings.Split(text, "\n")
		D.messages = append(D.messages, lines...)
		D.UpdateText()
	})
}

func (D *ErrorConsoleWidget) UpdateText() {
	if D.content.Height == 0 {
		return
	}

	H := D.content.Height - 2
	start := len(D.messages) - H

	for i := len(D.messages) - 1; i >= 0 && i > start; i -= 1 {
		line := D.messages[i]
		lCnt := (len(line) / D.content.Width)
		start += lCnt
	}
	if start < 0 {
		start = 0
	}

	content := D.messages[start:]
	D.content.Text = strings.Join(content, "\n")

	ui.Render(D.content)
}

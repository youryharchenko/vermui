package console

import (
	"fmt"
	"strings"
	"time"

	ui "github.com/verdverm/vermui"
)

type DevConsoleWidget struct {
	Rows []*ui.Row

	content  *ui.Par
	messages []string
}

func NewDevConsoleWidget() *DevConsoleWidget {
	content := ui.NewPar("")

	content.BorderLabel = " console "
	content.Height = 0
	content.Border = false
	content.BorderFg = ui.ColorGreen
	content.BorderLabelFg = ui.ColorWhite
	// content.TextFgColor = ui.ColorGreen
	// content.TextBgColor = ui.ColorBlack

	rows := []*ui.Row{
		ui.NewRow(
			ui.NewCol(12, 0, content),
		),
	}

	return &DevConsoleWidget{
		Rows:     rows,
		content:  content,
		messages: []string{},
	}
}

func (D *DevConsoleWidget) Init() {

	ui.Handle("/sys/kbd/C-l", func(ev ui.Event) {
		if D.content.Height > 0 {
			D.content.Height = 0
			D.content.Border = false
		} else {
			D.content.Height = 24
			D.content.Border = true
			D.UpdateText()
		}
	})

	ui.Handle("/console", func(ev ui.Event) {
		text := fmt.Sprintf("[%s] %v", time.Unix(ev.Time, 0).Format("2006-01-02 15:04:05"), ev.Data)

		level := strings.TrimPrefix(ev.Path, "/console/")
		switch level {
		case "crit":
			text = fmt.Sprintf("[[crit]  %s](fg-white,fg-bold,bg-red)", text)
		case "error":
			text = fmt.Sprintf("[[error] %s](fg-red)", text)
		case "warn":
			text = fmt.Sprintf("[[warn]  %s](fg-yellow)", text)
		case "info":
			text = fmt.Sprintf("[[info]  %s](fg-white)", text)
		case "debug":
			text = fmt.Sprintf("[[debug] %s](fg-green)", text)
		case "trace":
			text = fmt.Sprintf("[[trace] %s](fg-cyan)", text)
		}

		lines := strings.Split(text, "\n")
		D.messages = append(D.messages, lines...)
		D.UpdateText()
	})
}

func (D *DevConsoleWidget) UpdateText() {
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

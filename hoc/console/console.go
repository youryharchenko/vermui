package console

import (
	"fmt"
	"strings"
	"time"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/render"
	"github.com/verdverm/vermui/widgets/text"
)

type DevConsoleWidget struct {
	*text.Par
	messages []string
}

func NewDevConsoleWidget() *DevConsoleWidget {
	c := &DevConsoleWidget{
		Par:      text.NewPar(""),
		messages: []string{},
	}

	c.Height = 0
	c.Border = false
	c.BorderFg = render.ColorGreen
	c.BorderLabelFg = render.ColorWhite

	return c
}

func (D *DevConsoleWidget) Init() {

	vermui.AddGlobalHandler("/sys/kbd/C-l", func(ev events.Event) {
		D.Lock()
		if D.Height > 0 {
			D.BorderLabel = ""
			D.Height = 0
			D.Border = false
		} else {
			D.BorderLabel = " console "
			D.Height = 24
			D.Border = true
		}
		D.Unlock()
		D.UpdateText()
		go events.SendCustomEvent("/sys/redraw", "dev-console")
	})

	vermui.AddGlobalHandler("/console", func(ev events.Event) {
		text := fmt.Sprintf("[%s] %v", time.Unix(ev.Time, 0).Format("2006-01-02 15:04:05"), ev.Data)
		lines := strings.Split(text, "\n")

		level := strings.TrimPrefix(ev.Path, "/console/")

		D.Lock()
		for _, line := range lines {
			switch level {
			case "crit":
				line = fmt.Sprintf("[crit  %s](fg-white,fg-bold,bg-red)", line)
			case "error":
				line = fmt.Sprintf("[error %s](fg-red)", line)
			case "warn":
				line = fmt.Sprintf("[warn  %s](fg-yellow)", line)
			case "info":
				line = fmt.Sprintf("[info  %s](fg-white)", line)
			case "debug":
				line = fmt.Sprintf("[debug %s](fg-green)", line)
			case "trace":
				line = fmt.Sprintf("[trace %s](fg-cyan)", line)
			}

			D.messages = append(D.messages, line)
		}
		D.Unlock()

		D.UpdateText()
	})
}

func (D *DevConsoleWidget) UpdateText() {
	if D.Height == 0 {
		return
	}
	D.Lock()
	defer D.Unlock()

	H := D.Height - 2
	start := len(D.messages) - H

	for i := len(D.messages) - 1; i >= 0 && i > start; i -= 1 {
		line := D.messages[i]
		lCnt := (len(line) / D.Width)
		start += lCnt
	}
	if start < 0 {
		start = 0
	}

	content := D.messages[start:]
	D.Text = strings.Join(content, "\n")

	go vermui.Render(D)
}

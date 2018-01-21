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

type ErrorConsoleWidget struct {
	*text.Par
	messages []string
}

func NewErrorConsoleWidget() *ErrorConsoleWidget {
	c := &ErrorConsoleWidget{
		Par:      text.NewPar(""),
		messages: []string{},
	}

	c.Height = 0
	c.Border = false
	c.BorderFg = render.ColorRed
	c.BorderLabelFg = render.ColorWhite

	return c
}

func (D *ErrorConsoleWidget) Init() {

	vermui.AddGlobalHandler("/sys/kbd/C-e", func(ev events.Event) {
		//D.Lock()
		if D.Height > 0 {
			D.BorderLabel = ""
			D.Height = 0
			D.Border = false
		} else {
			D.BorderLabel = " errors "
			D.Height = 24
			D.Border = true
		}
		//D.Unlock()
		D.UpdateText()
		go events.SendCustomEvent("/sys/redraw", "dev-console")
	})

	vermui.AddGlobalHandler("/user/error", func(ev events.Event) {
		text := fmt.Sprintf("[%s] %v", time.Unix(ev.Time, 0).Format("2006-01-02 15:04:05"), ev.Data)
		lines := strings.Split(text, "\n")
		//D.Lock()
		D.messages = append(D.messages, lines...)
		//D.Unlock()
		D.UpdateText()
	})
}

func (D *ErrorConsoleWidget) UpdateText() {
	if D.Height == 0 {
		return
	}

	//D.Lock()
	//defer D.Unlock()

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

	go render.Render(D)

}

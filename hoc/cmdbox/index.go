package cmdbox

import (
	"fmt"
	"strings"

	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/render"
	"github.com/verdverm/vermui/widgets/text"
)

const emptyMsg = "press 'Ctrl-<space>' to enter a command or '/path/to/something' to navigate"
const inputRune = "█"

type Command interface {
	CommandName() string
	CommandUsage() string
	CommandHelp() string

	CommandCallback(args []string, context map[string]interface{})
}

type DefaultCommand struct {
	Name  string
	Usage string
	Help  string

	Callback func(args []string, context map[string]interface{})
}

func (DC *DefaultCommand) CommandName() string {
	return DC.Name
}

func (DC *DefaultCommand) CommandHelp() string {
	return DC.Help
}

func (DC *DefaultCommand) CommandUsage() string {
	return DC.Usage
}

func (DC *DefaultCommand) CommandCallback(args []string, context map[string]interface{}) {
	DC.Callback(args, context)
}

type CmdBoxWidget struct {
	*text.Par

	commands map[string]Command
}

func New() *CmdBoxWidget {
	cb := &CmdBoxWidget{
		Par:      text.NewPar(emptyMsg),
		commands: make(map[string]Command),
	}
	cb.Height = 3
	cb.PaddingLeft = 1
	cb.BorderFg = render.ColorBlue
	cb.BorderLabelFg = render.StringToAttribute("red,bold")

	return cb
}

func (CB *CmdBoxWidget) AddCommandCallback(command string, callback func([]string, map[string]interface{})) Command {
	c := &DefaultCommand{
		Name:     command,
		Usage:    command,
		Help:     "no help for " + command,
		Callback: callback,
	}
	CB.commands[c.CommandName()] = c
	return c
}

func (CB *CmdBoxWidget) AddCommand(command Command) {
	// go events.SendCustomEvent("/console/info", "adding command: "+command.CommandName())
	CB.commands[command.CommandName()] = command
}

func (CB *CmdBoxWidget) RemoveCommand(command Command) {
	delete(CB.commands, command.CommandName())
}

func (CB *CmdBoxWidget) Mount() error {
	// fmt.Println("cmdbox Mount")
	CB.AddHandler("/sys/kbd/C-<space>", func(e events.Event) {
		// fmt.Println("cmdbox - look at me!")
		// go events.SendCustomEvent("/console/trace", fmt.Sprint("cmdbox - look at me!", e.Path))
		CB.Focus()
	})

	CB.AddHandler("/user/error", func(e events.Event) {
		CB.Lock()
		CB.Text = fmt.Sprintf("[Error](bg-red,fg-bold): [%v](fg-yellow)", e.Data)
		CB.Unlock()

		render.Render(CB)
	})

	return nil
}
func (CB *CmdBoxWidget) Unmount() error {
	// fmt.Println("cmdbox - bye bye!")
	CB.RemoveHandler("/sys/kbd/C-<space>")
	CB.RemoveHandler("/user/error")

	return nil
}

func (CB *CmdBoxWidget) Focus() error {
	CB.Lock()
	CB.BorderFg = render.ColorRed
	CB.BorderLabelFg = render.StringToAttribute("bold,blue")
	CB.BorderLabelBg = render.ColorDefault
	CB.Text = inputRune
	CB.Unlock()

	render.Render(CB)
	CB.AddHandler("/sys/kbd", func(e events.Event) {
		key := e.Path[9:]
		CB.handleKey(key)
		render.Render(CB)
	})

	return nil
}

func (CB *CmdBoxWidget) Unfocus() error {
	CB.Lock()
	// fmt.Println("cmdbox - you lack focus!")
	CB.BorderFg = render.ColorBlue
	CB.BorderLabelFg = render.StringToAttribute("red,bold")
	CB.Text = emptyMsg
	CB.Unlock()

	go CB.RemoveHandler("/sys/kbd")

	render.Render(CB)
	return nil
}

func (CB *CmdBoxWidget) handleKey(key string) {
	// fmt.Println("cmdbox - key:", key)

	// handle first key after submit
	CB.Lock()
	if CB.Text == "Sent!" {
		CB.Text = ""
	}

	// strip the input rune
	CB.Text = strings.TrimSuffix(CB.Text, inputRune)
	L := len(CB.Text)
	CB.Unlock()

	// handle key
	switch key {
	case "<escape>":
		go CB.Unfocus()
		return

	case "<enter>":
		if L > 0 {
			CB.Lock()
			fields := strings.Fields(CB.Text)
			CB.Unlock()
			if len(fields) == 1 {
				CB.Submit(fields[0], nil)
			} else {
				CB.Submit(strings.ToLower(fields[0]), fields[1:])
			}
			go CB.Unfocus()
			return
		}

	case "<backspace>":
		CB.Lock()
		if L > 0 {
			CB.Text = CB.Text[:L-1] + inputRune
		} else {
			CB.Text = inputRune
		}
		CB.Unlock()

	case "<space>":
		CB.Lock()
		CB.Text += " " + inputRune
		CB.Unlock()

	default:
		// just one rune allowed, otherwise it's a special charactor
		if len(key) == 1 {
			CB.Text += key + inputRune
		}
	}

	render.Render(CB)
}

func (CB *CmdBoxWidget) Submit(command string, args []string) {
	if len(command) == 0 {
		return
	}
	command = strings.ToLower(command)
	if command[:1] == "/" {
		go events.SendCustomEvent("/router/dispatch", command)
		return
	}
	cmd, ok := CB.commands[command]
	if !ok {
		// render for the user
		go events.SendCustomEvent("/user/error", fmt.Sprintf("unknown command %q", command))
		// log to console
		go events.SendCustomEvent("/console/warn", fmt.Sprintf("unknown command %q", command))
		return
	}

	go cmd.CommandCallback(args, nil)
}

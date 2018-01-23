package cmdbox

import (
	"fmt"
	"strings"
	"sync"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/verdverm/vermui/events"
)

const emptyMsg = "press 'Ctrl-<space>' to enter a command or '/path/to/something' to navigate"

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
	events.HandledWidget
	sync.Mutex

	*tview.InputField

	commands map[string]Command
}

func New() *CmdBoxWidget {
	cb := &CmdBoxWidget{
		InputField: tview.NewInputField(),
		commands:   make(map[string]Command),
	}

	cb.
		SetTitle("  Edsger  ").
		SetTitleAlign(tview.AlignLeft).
		SetTitleColor(tcell.ColorRed).
		SetBorder(true).
		SetBorderColor(tcell.ColorBlue)

	return cb
}

func (CB *CmdBoxWidget) Id() string {
	return CB.InputField.Id()
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

func (CB *CmdBoxWidget) Init() {
	CB.Mount()
}

func (CB *CmdBoxWidget) Mount() error {
	CB.AddHandler("/sys/kbd/C-<space>", func(e events.Event) {
		// CB.Focus()
	})
	CB.AddHandler("/sys/kbd/C-space", func(e events.Event) {
		// CB.Focus()
	})

	CB.AddHandler("/user/error", func(e events.Event) {
		str := fmt.Sprintf("[Error](bg-red,fg-bold): [%v](fg-yellow)", e.Data)
		CB.Blur()
		CB.SetText(str)
	})

	return nil
}
func (CB *CmdBoxWidget) Unmount() error {
	// fmt.Println("cmdbox - bye bye!")
	CB.RemoveHandler("/sys/kbd/C-<space>")
	CB.RemoveHandler("/sys/kbd/C-space")
	CB.RemoveHandler("/user/error")

	return nil
}

/*
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
*/

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

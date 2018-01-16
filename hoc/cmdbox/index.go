package cmdbox

import (
	"fmt"
	"strings"

	ui "github.com/verdverm/vermui"
)

const emptyMsg = "press 'Ctrl-<space>' to enter a command"
const inputRune = "█"

type Command interface {
	CommandName() string
	CommandUsage() string
	CommandHelp() string

	CommandCallback(args []string, context map[string]interface{})
}

type CmdBoxWidget struct {
	*ui.Par

	commands map[string]Command
}

func New() *CmdBoxWidget {
	cb := &CmdBoxWidget{
		Par:      ui.NewPar(emptyMsg),
		commands: make(map[string]Command),
	}
	cb.Height = 3
	cb.PaddingLeft = 1
	cb.BorderFg = ui.ColorBlue
	cb.BorderLabelFg = ui.StringToAttribute("red,bold")
	return cb
}

func (CB *CmdBoxWidget) Add(command Command) {
	// go ui.SendCustomEvt("/console/info", "adding command: "+command.CommandName())
	CB.commands[command.CommandName()] = command
}

func (CB *CmdBoxWidget) Remove(command Command) {
	delete(CB.commands, command.CommandName())
}

func (CB *CmdBoxWidget) Mount() error {
	// fmt.Println("cmdbox Mount")
	CB.Handle("/sys/kbd/C-<space>", func(e ui.Event) {
		// fmt.Println("cmdbox - look at me!")
		// go ui.SendCustomEvt("/dev/messages", fmt.Sprint("cmdbox - look at me!", e.Path))
		CB.Focus()
	})

	return nil
}
func (CB *CmdBoxWidget) Unmount() error {
	// fmt.Println("cmdbox - bye bye!")
	CB.RemoveHandle("/sys/kbd/C-<space>")

	return nil
}

func (CB *CmdBoxWidget) Focus() error {
	// fmt.Println("cmdbox - pay attention!")
	CB.BorderFg = ui.ColorRed
	CB.BorderLabelFg = ui.StringToAttribute("bold,blue")
	CB.BorderLabelBg = ui.ColorDefault
	CB.Text = inputRune

	CB.Handle("/sys/kbd", func(e ui.Event) {
		// go ui.SendCustomEvt("/dev/messages", fmt.Sprint("cmdbox - keyboard!", e.Path))
		key := e.Path[9:]
		CB.handleKey(key)
	})

	CB.Handle("/user/error", func(e ui.Event) {
		CB.Text = fmt.Sprintf("[Error](bg-red,fg-bold): [%v](fg-yellow)", e.Data)
		// CB.BorderFg = ui.ColorWhite
		// CB.BorderLabelBg = ui.ColorWhite
	})

	return nil
}

func (CB *CmdBoxWidget) Unfocus() error {
	// fmt.Println("cmdbox - you lack focus!")
	CB.BorderFg = ui.ColorBlue
	CB.BorderLabelFg = ui.StringToAttribute("red,bold")
	CB.Text = emptyMsg

	CB.RemoveHandle("/sys/kbd")

	return nil
}

func (CB *CmdBoxWidget) handleKey(key string) {
	// fmt.Println("cmdbox - key:", key)

	// handle first key after submit
	if CB.Text == "Sent!" {
		CB.Text = ""
	}

	// strip the input rune
	CB.Text = strings.TrimSuffix(CB.Text, inputRune)
	L := len(CB.Text)

	// handle key
	switch key {
	case "<escape>":
		CB.Unfocus()

	case "<enter>":
		if L > 0 {
			fields := strings.Fields(CB.Text)
			if len(fields) == 1 {
				CB.Submit(fields[0], nil)
			} else {
				CB.Submit(strings.ToLower(fields[0]), fields[1:])
			}
			CB.Unfocus()
		}

	case "<backspace>":
		if L > 0 {
			CB.Text = CB.Text[:L-1] + inputRune
		} else {
			CB.Text = inputRune
		}

	case "<space>":
		CB.Text += " " + inputRune

	default:
		// just one rune allowed, otherwise it's a special charactor
		if len(key) == 1 {
			CB.Text += key + inputRune
		}
	}

	ui.Render(CB)
}

func (CB *CmdBoxWidget) Submit(command string, args []string) {
	command = strings.ToLower(command)
	cmd, ok := CB.commands[command]
	if !ok {
		// render for the user
		go ui.SendCustomEvt("/user/error", fmt.Sprintf("unknown command %q", command))
		// log to console
		go ui.SendCustomEvt("/console/warn", fmt.Sprintf("unknown command %q", command))
		return
	}

	go cmd.CommandCallback(args, nil)
}

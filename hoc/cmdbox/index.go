package cmdbox

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"github.com/verdverm/vermui"
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
	vermui.AddWidgetHandler(CB, "/sys/kbd/C-<space>", func(e events.Event) {
		// CB.Focus()
	})
	vermui.AddWidgetHandler(CB, "/sys/kbd/C-space", func(e events.Event) {
		// CB.Focus()
	})

	vermui.AddWidgetHandler(CB, "/user/error", func(e events.Event) {
		str := fmt.Sprintf("[Error](bg-red,fg-bold): [%v](fg-yellow)", e.Data)
		// CB.Blur()
		CB.SetText(str)
	})

	return nil
}
func (CB *CmdBoxWidget) Unmount() error {
	// fmt.Println("cmdbox - bye bye!")
	vermui.RemoveWidgetHandler(CB, "/sys/kbd/C-<space>")
	vermui.RemoveWidgetHandler(CB, "/sys/kbd/C-space")
	vermui.RemoveWidgetHandler(CB, "/user/error")

	return nil
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

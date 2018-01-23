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

	cb.InputField.
		SetFieldBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetLabel(" ")

	cb.Mount(nil)

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

func (CB *CmdBoxWidget) Mount(context map[string]interface{}) error {
	vermui.AddWidgetHandler(CB, "/sys/key/C-space", func(e events.Event) {
		CB.SetText("")
		CB.SetBorderColor(tcell.Color69)

		vermui.SetFocus(CB)
	})

	vermui.AddWidgetHandler(CB, "/user/error", func(e events.Event) {
		str := fmt.Sprintf("%v", e.Data.(*events.EventCustom).Data())
		CB.SetBorderColor(tcell.ColorRed)
		CB.SetFieldTextColor(tcell.ColorOrange)
		CB.SetText(str)

		vermui.Unfocus()
	})

	return nil
}
func (CB *CmdBoxWidget) Unmount() error {
	vermui.RemoveWidgetHandler(CB, "/sys/key/C-space")
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

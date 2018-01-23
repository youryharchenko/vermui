package vermui

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/rivo/tview"

	"github.com/verdverm/vermui/layouts"

	"github.com/verdverm/vermui/lib"
	"github.com/verdverm/vermui/lib/events"
)

func Init() error {
	err := lib.Init()
	if err != nil {
		return errors.Wrapf(err, "in vermui.Init()\n")
	}

	return nil
}

func Start() error {
	defer func() {
		err := recover()
		if err != nil {
			lib.Stop()
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
	}()
	err := lib.Start()
	if err != nil {
		return errors.Wrapf(err, "in vermui.Start()\n")
	}
	return nil
}

func Stop() error {
	err := lib.Stop()
	if err != nil {
		return errors.Wrapf(err, "in vermui.Stop()\n")
	}
	return nil
}

func GetLayout() layouts.Layout {
	return lib.GetRootLayout()
}

func SetLayout(l layouts.Layout) {
	lib.SetRootLayout(l)
}

func AddGlobalHandler(path string, handler func(events.Event)) {
	lib.AddGlobalHandler(path, handler)
}

func RemoveGlobalHandler(path string) {
	lib.RemoveGlobalHandler(path)
}

func ClearGlobalHandlers() {
	lib.ClearGlobalHandlers()
}

func Application() *tview.Application {
	return lib.Application()
}

func Draw() {
	lib.Draw()
}

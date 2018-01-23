package main

import (
	"fmt"
	"os"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/events"
)

func main() {
	fmt.Println("Welcome to VermUI")

	// initialize VermUI first
	err := vermui.Init()
	if err != nil {
		vermui.Stop()
		fmt.Printf("\n\n%v\n\n", err)
		os.Exit(1)
	}

	// build our routing and layouts, give them to VermUI
	layout := buildLayout()
	vermui.SetRootView(layout)

	// Handler: press <Ctrl>-c to quit
	vermui.AddGlobalHandler("/sys/kbd/C-c", func(events.Event) {
		vermui.Stop()
	})

	// Handler: press <home> to go to main screen
	vermui.AddGlobalHandler("/sys/kbd/<home>", func(events.Event) {
		go events.SendCustomEvent("/router/dispatch", "/home")
	})

	// go to the initial route/view
	// go events.SendCustomEvent("/router/dispatch", "/page-1")

	// block until Stop is called
	vermui.Start()
}

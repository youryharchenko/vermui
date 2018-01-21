package main

import (
	"fmt"
	"os"

	"github.com/verdverm/vermui"

	"github.com/verdverm/vermui/lib/events"
)

func main() {
	fmt.Println("Welcome to VermUI")

	err := vermui.Init()
	if err != nil {
		vermui.Stop()
		fmt.Printf("\n\n%v\n\n", err)
		os.Exit(1)
	}

	layout := buildLayout()
	vermui.SetLayout(layout)

	// lib.Render(p) // feel free to call Render, it's async and non-block

	// handle key C-c pressing
	vermui.AddGlobalHandler("/sys/kbd/C-c", func(events.Event) {
		// press q to quit
		vermui.Stop()
	})

	go events.SendCustomEvent("/router/dispatch", "/p")
	vermui.Start() // block until Stop is called
}

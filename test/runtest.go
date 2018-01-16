// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/debug"
)

func main() {
	// run as client
	if len(os.Args) > 1 {
		fmt.Print(debug.ConnectAndListen())
		return
	}

	// run as server
	go func() { panic(debug.ListenAndServe()) }()

	if err := vermui.Init(); err != nil {
		panic(err)
	}
	defer vermui.Close()

	//vermui.UseTheme("helloworld")
	b := vermui.NewBlock()
	b.Width = 20
	b.Height = 20
	b.Float = vermui.AlignCenter
	b.BorderLabel = "[HELLO](fg-red,bg-white) [WORLD](fg-blue,bg-green)"

	vermui.Render(b)

	vermui.Handle("/sys", func(e vermui.Event) {
		k, ok := e.Data.(vermui.EvtKbd)
		debug.Logf("->%v\n", e)
		if ok && k.KeyStr == "q" {
			vermui.StopLoop()
		}
	})

	vermui.Handle(("/usr"), func(e vermui.Event) {
		debug.Logf("->%v\n", e)
	})

	vermui.Handle("/timer/1s", func(e vermui.Event) {
		t := e.Data.(vermui.EvtTimer)
		vermui.SendCustomEvt("/usr/t", t.Count)

		if t.Count%2 == 0 {
			b.BorderLabel = "[HELLO](fg-red,bg-green) [WORLD](fg-blue,bg-white)"
		} else {
			b.BorderLabel = "[HELLO](fg-blue,bg-white) [WORLD](fg-red,bg-green)"
		}

		vermui.Render(b)

	})

	vermui.Loop()
}

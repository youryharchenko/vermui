package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/maruel/panicparse/stack"
	"github.com/rivo/tview"

	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/events"
)

var app *tview.Application
var rootLayout layouts.Layout

// Init initializes vermui library. This function should be called before any others.
// After initialization, the library must be finalized by 'Close' function.
func Init() error {
	/*
		err := render.Init()
		if err != nil {
			return err
		}
	*/

	app = tview.NewApplication()

	err := events.Init(app)
	if err != nil {
		return err
	}

	events.Handle("/", events.DefaultHandler)

	events.Handle("/sys/redraw", func(e events.Event) {
		app.Draw()
	})

	return nil
}

// blocking call
func Start() error {
	defer func() {
		e := recover()
		if e != nil {
			Stop()
			// Print a formatted panic output
			fmt.Fprintf(os.Stderr, "Captured a panic(value=%v) lib.Start()... Exit vermui and clean terminal...\nPrint stack trace:\n\n", e)
			//debug.PrintStack()
			gs, err := stack.ParseDump(bytes.NewReader(debug.Stack()), os.Stderr)
			if err != nil {
				debug.PrintStack()
				os.Exit(1)
			}
			p := &stack.Palette{}
			buckets := stack.SortBuckets(stack.Bucketize(gs, stack.AnyValue))
			srcLen, pkgLen := stack.CalcLengths(buckets, false)
			for _, bucket := range buckets {
				io.WriteString(os.Stdout, p.BucketHeader(&bucket, false, len(buckets) > 1))
				io.WriteString(os.Stdout, p.StackLines(&bucket.Signature, srcLen, pkgLen, false))
			}
			os.Exit(1)
		}
	}()

	app.SetFocus(rootLayout)
	app.SetRoot(rootLayout, true)

	go events.Start()
	go events.SendCustomEvent("/console/debug", "App Starting")

	// blocking
	return app.Run()
}

// Close finalizes vermui library,
// should be called after successful initialization when vermui's functionality isn't required anymore.
func Stop() error {
	app.Stop()
	err := events.Stop()
	return err
}

func Draw() {
	app.Draw()
}

func Application() *tview.Application {
	return app
}

func GetRootLayout() layouts.Layout {
	return rootLayout
}

func SetRootLayout(l layouts.Layout) {
	rootLayout = l
}

func SetFocus(p tview.Primitive) {
	app.SetFocus(p)
}

func AddGlobalHandler(path string, handler func(events.Event)) {
	events.Handle(path, handler)
}

func RemoveGlobalHandler(path string) {
	events.RemoveHandle(path)
}

func ClearGlobalHandlers() {
	events.ResetHandlers()
}

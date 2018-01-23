package vermui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/maruel/panicparse/stack"
	"github.com/rivo/tview"

	"github.com/verdverm/vermui/events"
)

var app *tview.Application
var rootView tview.Primitive

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

	events.AddGlobalHandler("/", events.DefaultHandler)

	events.AddGlobalHandler("/sys/redraw", func(e events.Event) {
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

	app.SetFocus(rootView)
	app.SetRoot(rootView, true)

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

func Application() *tview.Application {
	return app
}

func Draw() {
	app.Draw()
}

func GetRootView() tview.Primitive {
	return rootView
}

func SetRootView(v tview.Primitive) {
	rootView = v
}

func GetFocus() (p tview.Primitive) {
	return app.GetFocus()
}

func SetFocus(p tview.Primitive) {
	app.SetFocus(p)
	app.Draw()
}
func Unfocus() {
	cur := app.GetFocus()
	cur.Blur()
	app.Screen().HideCursor()

	app.SetFocus(rootView)
	app.Draw()
}

func AddGlobalHandler(path string, handler func(events.Event)) {
	events.AddGlobalHandler(path, handler)
}

func RemoveGlobalHandler(path string) {
	events.RemoveGlobalHandler(path)
}

func ClearGlobalHandlers() {
	events.ClearGlobalHandlers()
}

func AddWidgetHandler(widget tview.Primitive, path string, handler func(events.Event)) {
	events.AddWidgetHandler(widget, path, handler)
}

func RemoveWidgetHandler(widget tview.Primitive, path string) {
	events.RemoveWidgetHandler(widget, path)
}

func ClearWidgetHandlers(widget tview.Primitive) {
	events.ClearWidgetHandlers(widget)
}

package lib

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"

	"github.com/maruel/panicparse/stack"
	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/render"
)

var rootLayout layouts.Layout
var renderTimer *time.Ticker

// Init initializes vermui library. This function should be called before any others.
// After initialization, the library must be finalized by 'Close' function.
func Init() error {
	err := render.Init()
	if err != nil {
		return err
	}

	err = events.Init()
	if err != nil {
		return err
	}

	events.Handle("/", events.DefaultHandler)
	events.Handle("/sys/wnd/resize", func(e events.Event) {
		w := e.Data.(events.EvtWnd)
		rootLayout.SetWidth(w.Width)
		rootLayout.Align()
		render.Clear()
		render.Render(rootLayout)
	})

	events.Handle("/sys/redraw", func(e events.Event) {
		rootLayout.SetWidth(render.TermWidth())
		rootLayout.Align()
		render.Clear()
		render.Render(rootLayout)
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

	rootLayout.SetWidth(render.TermWidth())
	rootLayout.Mount()
	rootLayout.Align()
	render.Clear()
	render.Start()
	render.Render(rootLayout)

	/*
		go func() {
			renderTimer = time.NewTicker(time.Millisecond * 200)
			for range renderTimer.C {
				render.Clear()
				if rootLayout.GetWidth() != render.TermWidth() {
					rootLayout.SetWidth(render.TermWidth())
					rootLayout.Align()
				}
				render.Render(rootLayout)
			}
		}()
	*/

	// blocking
	return events.Start()
}

// Close finalizes vermui library,
// should be called after successful initialization when vermui's functionality isn't required anymore.
func Stop() error {
	if renderTimer != nil {
		renderTimer.Stop()

	}
	render.Stop()
	render.Clear()
	err := events.Stop()
	return err
}

func GetRootLayout() layouts.Layout {
	return rootLayout
}

func SetRootLayout(l layouts.Layout) {
	rootLayout = l
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

func Render(bs ...render.Bufferer) {
	render.Render(bs...)
}

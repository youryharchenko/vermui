package lib

import (
	"time"

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
		rootLayout.Align()
		render.Clear()
		render.Render(rootLayout)
	})

	return nil
}

// blocking call
func Start() error {
	render.Start()

	renderTimer = time.NewTicker(time.Millisecond * 2000)

	rootLayout.SetWidth(render.TermWidth())
	rootLayout.Align()
	rootLayout.Mount()
	render.Clear()
	render.Render(rootLayout)
	/*
		go func() {
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
	renderTimer.Stop()
	render.Clear()
	render.Stop()
	return events.Stop()
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

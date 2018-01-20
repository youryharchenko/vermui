package router

import (
	"github.com/pkg/errors"

	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/mux"
	"github.com/verdverm/vermui/lib/render"
)

type Router struct {
	width        int
	activeName   string
	activeLayout layouts.Layout

	// internal router
	iRouter *mux.Router

	// for holding the event handles
	dummy *render.Block
}

func New() *Router {
	lv := &Router{
		iRouter: mux.NewRouter(),
		dummy:   render.NewBlock(),
	}
	return lv
}

func (R *Router) GetWidth() int {
	if R.activeLayout != nil {
		return R.activeLayout.GetWidth()
	}
	return R.width
}

func (R *Router) SetWidth(w int) {
	R.width = w
	if R.activeLayout != nil {
		R.activeLayout.SetWidth(w)
	}
}

func (R *Router) Align() {
	if R.activeLayout != nil {
		// R.activeLayout.SetWidth(R.width)
		R.activeLayout.Align()
	}
}

func (R *Router) Buffer() render.Buffer {
	if R.activeLayout != nil {
		return R.activeLayout.Buffer()
	} else {
		return R.dummy.Buffer()
	}
}

func (R *Router) Mount() error {
	R.dummy.AddHandler("/router/dispatch", func(ev events.Event) {
		path := ev.Data.(string)
		// fmt.Println("\n\nDispatch", path)
		layout, err := R.iRouter.Dispatch(path)
		if err != nil {
			go events.SendCustomEvent("/console/error", errors.Wrap(err, "in dispatch handler"))
		}
		if layout != nil {
			R.setActive(layout)
		} else {
			go events.SendCustomEvent("/console/error", "nil layout in dispatch handler")
		}
	})

	if R.activeLayout != nil {
		return R.activeLayout.Mount()
	}
	return nil
}

func (R *Router) Unmount() error {
	R.dummy.RemoveHandler("/router/dispatch")

	if R.activeLayout != nil {
		err := R.activeLayout.Unmount()
		if err != nil {
			return err
		}
		R.activeLayout = nil
	}

	return nil
}

func (R *Router) SetNotFound(layout layouts.Layout) {
	handler := func(*mux.Request) (layouts.Layout, error) {
		return layout, nil
	}
	R.iRouter.NotFoundHandler = mux.NewDefaultHandler(handler)
}
func (R *Router) AddRouteLayout(path string, layout layouts.Layout) error {
	handler := func(*mux.Request) (layouts.Layout, error) {
		return layout, nil
	}
	R.iRouter.Handle(path, mux.NewDefaultHandler(handler))
	return nil
}

func (R *Router) AddRouteHandlerFunc(path string, handler mux.HandlerFunc) error {
	R.iRouter.Handle(path, mux.NewDefaultHandler(handler))
	return nil
}

func (R *Router) AddRouteHandler(path string, handler mux.Handler) error {
	R.iRouter.Handle(path, handler)
	return nil
}

func (R *Router) SetActive(path string) {
	layout, err := R.iRouter.Dispatch(path)
	if err != nil {
		go events.SendCustomEvent("/console/error", errors.Wrap(err, "in dispatch handler"))
	}
	if layout != nil {
		R.setActive(layout)
	} else {
		go events.SendCustomEvent("/console/error", "nil layout in dispatch handler")
	}
}

func (R *Router) setActive(layout layouts.Layout) {
	// mount new layout
	layout.SetWidth(R.width)
	layout.Mount()

	// unmount deactivating
	if R.activeLayout != nil {
		R.activeLayout.Unmount()
	}

	// finally, set the active layout and redraw
	R.activeLayout = layout
	go events.SendCustomEvent("/sys/redraw", "router - ")
}

func (R *Router) Show() {
	if R.activeLayout != nil {
		R.activeLayout.Show()
	}
}
func (R *Router) Hide() {
	if R.activeLayout != nil {
		R.activeLayout.Hide()
	}
}
func (R *Router) Focus() {
	if R.activeLayout != nil {
		R.activeLayout.Focus()
	}
}
func (R *Router) Unfocus() {
	if R.activeLayout != nil {
		R.activeLayout.Unfocus()
	}
}

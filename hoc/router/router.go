package router

import (
	"sync"

	"github.com/pkg/errors"
	"github.com/rivo/tview"

	"github.com/verdverm/vermui/events"
	"github.com/verdverm/vermui/mux"
)

type RoutePair struct {
	Path  string
	Thing interface{}
}

type Routable interface {
	Routings() []RoutePair
}

type Router struct {
	events.HandledWidget
	sync.Mutex

	*tview.Pages
	activeName   string
	activeLayout tview.Primitive

	// internal router
	iRouter *mux.Router
}

func New() *Router {
	lv := &Router{
		Pages:   tview.NewPages(),
		iRouter: mux.NewRouter(),
	}
	return lv
}

/*
func (R *Router) Mount() error {
	R.dummy.AddHandler("/router/dispatch", func(ev events.Event) {
		path := ev.Data.(string)
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
*/

func (R *Router) SetNotFound(layout tview.Primitive) {
	handler := func(*mux.Request) (tview.Primitive, error) {
		return layout, nil
	}
	R.iRouter.NotFoundHandler = mux.NewDefaultHandler(handler)
}

func (R *Router) AddRoute(path string, thing interface{}) error {

	switch t := thing.(type) {
	case tview.Primitive:
		R.AddRouteLayout(path, t)

	case mux.HandlerFunc:
		R.AddRouteHandlerFunc(path, t)

	case mux.Handler:
		R.AddRouteHandler(path, t)

	default:
		return errors.New("Unknown thing to be routed to...")
	}

	return nil
}

func (R *Router) AddRouteLayout(path string, layout tview.Primitive) error {
	handler := func(*mux.Request) (tview.Primitive, error) {
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

func (R *Router) setActive(layout tview.Primitive) {

	// mount new layout
	// layout.Mount()

	// lock R
	R.Lock()
	defer R.Unlock()

	// unmount deactivating
	if R.activeLayout != nil {
		// R.activeLayout.Unmount()
	}

	// finally, set the active layout and redraw
	R.activeLayout = layout

	go events.SendCustomEvent("/sys/redraw", "router - ")
}

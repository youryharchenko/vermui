package router

import (
	"github.com/pkg/errors"
	"github.com/rivo/tview"

	"github.com/verdverm/vermui"
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
	*tview.Pages

	// internal router
	iRouter *mux.Router
}

func New() *Router {
	r := &Router{
		Pages:   tview.NewPages(),
		iRouter: mux.NewRouter(),
	}
	vermui.AddWidgetHandler(r, "/router/dispatch", func(ev events.Event) {
		path := ev.Data.(*events.EventCustom).Data().(string)

		layout, err := r.iRouter.Dispatch(path)
		if err != nil {
			go events.SendCustomEvent("/console/error", errors.Wrap(err, "in dispatch handler"))
		}

		go events.SendCustomEvent("/console/debug", "router received path: "+path)
		return
		if layout != nil {
			r.setActive(layout)
		} else {
			go events.SendCustomEvent("/console/error", "nil layout in dispatch handler")
		}
	})

	return r
}

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
	R.AddPage(layout.Id(), layout, true, false)
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
	go events.SendCustomEvent("/console/error", "nil layout in dispatch handler")

	R.SetActive(layout.Id())

}

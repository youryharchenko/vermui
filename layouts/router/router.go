package router

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/verdverm/vermui/layouts"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/mux"
	"github.com/verdverm/vermui/lib/render"
)

type RoutePair struct {
	Path  string
	Thing interface{}
}

type Routable interface {
	Routings() []RoutePair
}

type Router struct {
	sync.Mutex
	x, y         int
	width        int
	height       int
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

func (R *Router) GetX() int {
	R.Lock()
	defer R.Unlock()

	if R.activeLayout != nil {
		return R.activeLayout.GetX()
	}
	return R.x
}
func (R *Router) SetX(x int) {
	R.Lock()
	defer R.Unlock()

	R.x = x
	if R.activeLayout != nil {
		R.activeLayout.SetX(x)
	}
}
func (R *Router) GetY() int {
	R.Lock()
	defer R.Unlock()

	if R.activeLayout != nil {
		return R.activeLayout.GetY()
	}
	return R.y
}
func (R *Router) SetY(y int) {
	R.Lock()
	defer R.Unlock()

	R.y = y
	if R.activeLayout != nil {
		R.activeLayout.SetY(y)
	}
}

func (R *Router) GetHeight() int {
	R.Lock()
	defer R.Unlock()

	if R.activeLayout != nil {
		h := R.activeLayout.GetHeight()
		return h
	}
	return R.height
}

func (R *Router) SetHeight(h int) {
	R.Lock()
	defer R.Unlock()

	R.height = h
}

func (R *Router) GetWidth() int {
	R.Lock()
	defer R.Unlock()

	if R.activeLayout != nil {
		return R.activeLayout.GetWidth()
	}
	return R.width
}

func (R *Router) SetWidth(w int) {
	R.Lock()
	defer R.Unlock()

	R.width = w
	if R.activeLayout != nil {
		R.activeLayout.SetWidth(w)
	}

}

func (R *Router) Align() {
	if R.GetWidth() == 0 {
		return
		panic(errors.New("Router has zero width"))
	}

	if R.activeLayout != nil {
		R.activeLayout.Align()

		// R.height = R.activeLayout.GetHeight()
	}
}

func (R *Router) Buffer() render.Buffer {
	//R.Lock()
	//defer R.Unlock()

	R.Align()
	if R.activeLayout != nil {
		return R.activeLayout.Buffer()
	} else {
		return render.NewBuffer()
	}
}

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

func (R *Router) SetNotFound(layout layouts.Layout) {
	handler := func(*mux.Request) (layouts.Layout, error) {
		return layout, nil
	}
	R.iRouter.NotFoundHandler = mux.NewDefaultHandler(handler)
}

func (R *Router) AddRoute(path string, thing interface{}) error {

	switch t := thing.(type) {
	case layouts.Layout:
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
	layout.Mount()
	layout.SetWidth(R.GetWidth())
	layout.SetX(R.GetX())
	layout.SetY(R.GetY())

	// unmount deactivating
	if R.activeLayout != nil {
		R.activeLayout.Unmount()
	}

	// finally, set the active layout and redraw
	R.Lock()
	R.activeLayout = layout
	R.Unlock()

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

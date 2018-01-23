package mux

import (
	"github.com/rivo/tview"
)

type HandlerFunc func(*Request) (tview.Primitive, error)

type Handler interface {
	Serve(*Request) (tview.Primitive, error)
}

type DefaultHandler struct {
	handler HandlerFunc
}

func NewDefaultHandler(handle HandlerFunc) *DefaultHandler {
	return &DefaultHandler{
		handler: handle,
	}
}

func (H *DefaultHandler) Serve(req *Request) (tview.Primitive, error) {
	return H.handler(req)
}

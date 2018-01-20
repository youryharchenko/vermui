package mux

import "github.com/verdverm/vermui/layouts"

type HandlerFunc func(*Request) (layouts.Layout, error)

type Handler interface {
	Serve(*Request) (layouts.Layout, error)
}

type DefaultHandler struct {
	handler HandlerFunc
}

func NewDefaultHandler(handle HandlerFunc) *DefaultHandler {
	return &DefaultHandler{
		handler: handle,
	}
}

func (H *DefaultHandler) Serve(req *Request) (layouts.Layout, error) {
	return H.handler(req)
}

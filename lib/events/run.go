package events

var defaultEventStream *EvtStream
var defaultWgtMgr WgtMgr

func Init() error {
	sysEvtChs = make([]chan Event, 0)
	go hookTermboxEvt()

	defaultEventStream = NewEventStream()
	defaultEventStream.Init()
	defaultEventStream.Merge("termbox", NewSysEvtCh())
	defaultEventStream.Merge("custom", usrEvtCh)

	defaultWgtMgr = NewWgtMgr()
	defaultEventStream.Hook(defaultWgtMgr.WgtHandlersHook())

	return nil
}

func Start() error {
	defaultEventStream.Loop()
	return nil
}

func Stop() error {
	defaultEventStream.StopLoop()
	return nil
}

func Merge(name string, ec chan Event) {
	defaultEventStream.Merge(name, ec)
}

func Handle(path string, handler func(Event)) {
	defaultEventStream.Handle(path, handler)
}

func RemoveHandle(path string) {
	defaultEventStream.RemoveHandle(path)
}

func ResetHandlers() {
	defaultEventStream.ResetHandlers()
}

func AddWgtHandler(wgt Widget, path string, handler func(Event)) {
	if _, ok := defaultWgtMgr[wgt.Id()]; !ok {
		defaultWgtMgr.AddWgt(wgt)
	}

	defaultWgtMgr.AddWgtHandler(wgt.Id(), path, handler)
}

func RemoveWgtHandler(wgt Widget, path string) {
	_, ok := defaultWgtMgr[wgt.Id()]
	if !ok {
		return
	}

	defaultWgtMgr.RmWgtHandler(wgt.Id(), path)
}

func ClearWgtHandlers(wgt Widget) {
	_, ok := defaultWgtMgr[wgt.Id()]
	if !ok {
		return
	}

	defaultWgtMgr.ClearWgtHandlers(wgt.Id())
}

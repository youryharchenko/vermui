package streamtable

import (
	"github.com/rivo/tview"

	"github.com/verdverm/vermui/lib/events"
)

type StreamTableSource func(chan string) chan interface{}
type StreamTableFormatter func(interface{}) [][]string

type StreamTable struct {
	*tview.Table

	TableHeader   [][]string
	DataSource    StreamTableSource
	DataFormatter StreamTableFormatter

	DataStreamer chan interface{}
	DataCommands chan string
	QuitChan     chan string
}

func NewStreamTable(header [][]string, source StreamTableSource, formatter StreamTableFormatter) *StreamTable {
	ST := &StreamTable{
		Table:         tview.NewTable(),
		TableHeader:   header,
		DataSource:    source,
		DataFormatter: formatter,
	}

	ST.QuitChan = make(chan string, 2)
	ST.DataCommands = make(chan string, 2)

	return ST
}

func (ST *StreamTable) Mount() error {
	return nil
}
func (ST *StreamTable) Unmount() error {
	return nil
}

func (ST *StreamTable) Show() {
	// already shown
	if ST.DataStreamer != nil {
		return
	}
	ST.DataStreamer = ST.DataSource(ST.DataCommands)
	first := true
	go func() {
		for {
			select {
			case data := <-ST.DataStreamer:
				ST.UpdateData(data)

			case <-ST.QuitChan:
				return
			}
			if first {
				events.SendCustomEvent("/sys/redraw", "stable")
				first = false
			}
		}
	}()
	go events.SendCustomEvent("/sys/redraw", "stable")
}
func (ST *StreamTable) Hide() {
	ST.DataCommands <- "quit"
	ST.QuitChan <- "quit"
	ST.DataStreamer = nil
	go events.SendCustomEvent("/sys/redraw", "stable")
}

func (ST *StreamTable) UpdateData(input interface{}) {

	data := ST.DataFormatter(input)

	rows := [][]string{}
	rows = append(rows, ST.TableHeader...)
	rows = append(rows, data...)

	for r := range rows {
		for c := range rows[r] {
			ST.Table.SetCell(r, c,
				tview.NewTableCell(rows[c][r]).
					SetAlign(tview.AlignCenter))
		}
	}
}

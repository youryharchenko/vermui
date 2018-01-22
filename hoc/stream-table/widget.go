package streamtable

import (
	"github.com/verdverm/vermui"
	"github.com/verdverm/vermui/layouts/align"
	"github.com/verdverm/vermui/lib/events"
	"github.com/verdverm/vermui/lib/render"
	"github.com/verdverm/vermui/widgets/table"
)

type StreamTableSource func(chan string) chan interface{}
type StreamTableFormatter func(interface{}) [][]string

type StreamTable struct {
	*table.Table

	TableHeader   [][]string
	DataSource    StreamTableSource
	DataFormatter StreamTableFormatter

	DataStreamer chan interface{}
	DataCommands chan string
	QuitChan     chan string
}

func NewStreamTable(header [][]string, source StreamTableSource, formatter StreamTableFormatter) *StreamTable {
	ST := &StreamTable{
		Table:         table.NewTable(),
		TableHeader:   header,
		DataSource:    source,
		DataFormatter: formatter,
	}

	ST.Table.FgColor = render.ColorWhite
	ST.Table.BgColor = render.ColorDefault
	ST.Table.TextAlign = align.AlignCenter
	ST.Table.Separator = false
	ST.Table.Border = false
	ST.Table.Height = 0

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
	ST.Table.Height = len(ST.Table.Rows) + 2
	ST.Table.Border = true
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
	ST.Table.Height = 0
	ST.Table.Border = false
	ST.Table.Rows = [][]string{}
	go events.SendCustomEvent("/sys/redraw", "stable")
}

func (ST *StreamTable) UpdateData(input interface{}) {

	data := ST.DataFormatter(input)

	rows := [][]string{}
	rows = append(rows, ST.TableHeader...)
	rows = append(rows, data...)

	oldRows := ST.Table.Rows
	ST.Table.Rows = rows
	if len(oldRows) != len(rows) {
		ST.Table.Height = len(ST.Table.Rows) + 2
		if ST.Table.Height > 20 {
			ST.Table.Height = 20
		}
		ST.Table.Analysis()
		ST.Table.SetSize()
	}
	vermui.Render(ST)

}

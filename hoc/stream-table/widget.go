package streamtable

import (
	ui "github.com/gizak/termui"
)

type StreamTableSource func(chan string) chan interface{}
type StreamTableFormatter func(interface{}) [][]string

type StreamTable struct {
	*ui.Table

	TableHeader   [][]string
	DataSource    StreamTableSource
	DataFormatter StreamTableFormatter

	DataStreamer chan interface{}
	DataCommands chan string
	QuitChan     chan string
}

func NewStreamTable(header [][]string, source StreamTableSource, formatter StreamTableFormatter) *StreamTable {
	ST := &StreamTable{
		Table:         ui.NewTable(),
		TableHeader:   header,
		DataSource:    source,
		DataFormatter: formatter,
	}

	ST.Table.FgColor = ui.ColorWhite
	ST.Table.BgColor = ui.ColorDefault
	ST.Table.TextAlign = ui.AlignCenter
	ST.Table.Separator = false
	ST.Table.Border = false
	ST.Table.Height = 0

	ST.QuitChan = make(chan string, 2)
	ST.DataCommands = make(chan string, 2)

	return ST
}

func (ST *StreamTable) Show() {
	// already shown
	if ST.DataStreamer != nil {
		return
	}
	ST.DataStreamer = ST.DataSource(ST.DataCommands)
	ST.Table.Height = len(ST.Table.Rows) + 2
	ST.Table.Border = true
	go func() {
		for {
			select {
			case data := <-ST.DataStreamer:
				ST.UpdateData(data)

			case <-ST.QuitChan:
				return
			}
		}
	}()
}
func (ST *StreamTable) Hide() {
	ST.DataCommands <- "quit"
	ST.QuitChan <- "quit"
	ST.DataStreamer = nil
	ST.Table.Height = 0
	ST.Table.Border = false
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
		ui.Render(ST.Table)
	}
}

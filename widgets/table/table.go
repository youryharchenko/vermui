// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package table

import (
	"strings"
	"sync"

	"github.com/verdverm/vermui/layouts/align"
	"github.com/verdverm/vermui/lib/render"
)

/* Table is like:

┌Awesome Table ────────────────────────────────────────────────┐
│  Col0          | Col1 | Col2 | Col3  | Col4  | Col5  | Col6  |
│──────────────────────────────────────────────────────────────│
│  Some Item #1  | AAA  | 123  | CCCCC | EEEEE | GGGGG | IIIII |
│──────────────────────────────────────────────────────────────│
│  Some Item #2  | BBB  | 456  | DDDDD | FFFFF | HHHHH | JJJJJ |
└──────────────────────────────────────────────────────────────┘

Datapoints are a two dimensional array of strings: [][]string

Example:
	data := [][]string{
		{"Col0", "Col1", "Col3", "Col4", "Col5", "Col6"},
		{"Some Item #1", "AAA", "123", "CCCCC", "EEEEE", "GGGGG", "IIIII"},
		{"Some Item #2", "BBB", "456", "DDDDD", "FFFFF", "HHHHH", "JJJJJ"},
	}

	table := vermui.NewTable()
	table.Rows = data  // type [][]string
	table.FgColor = vermui.ColorWhite
	table.BgColor = vermui.ColorDefault
	table.Height = 7
	table.Width = 62
	table.Y = 0
	table.X = 0
	table.Border = true
*/

// Table tracks all the attributes of a Table instance
type Table struct {
	render.Block
	Rows         [][]string
	CellWidth    []int
	FgColor      render.Attribute
	BgColor      render.Attribute
	FgColors     []render.Attribute
	BgColors     []render.Attribute
	CellFgColors [][]render.Attribute
	CellBgColors [][]render.Attribute
	Separator    bool
	TextAlign    align.Align
	sync.Mutex
}

// NewTable returns a new Table instance
func NewTable() *Table {
	table := &Table{Block: *render.NewBlock()}
	table.FgColor = render.ColorWhite
	table.BgColor = render.ColorDefault
	table.Separator = true
	table.Border = true
	return table
}

// CellsWidth calculates the width of a cell array and returns an int
func cellsWidth(cells []render.Cell) int {
	width := 0
	for _, c := range cells {
		width += c.Width()
	}
	return width
}

func (table *Table) analysis() [][]render.Cell {
	var rowCells [][]render.Cell
	length := len(table.Rows)
	if length < 1 {
		return rowCells
	}

	if len(table.FgColors) != length {
		table.FgColors = make([]render.Attribute, length)
		table.CellFgColors = make([][]render.Attribute, length)
		for y, row := range table.Rows {
			table.CellFgColors[y] = make([]render.Attribute, len(row))
		}
	}
	if len(table.BgColors) == 0 {
		table.BgColors = make([]render.Attribute, len(table.Rows))
		table.CellBgColors = make([][]render.Attribute, length)
		for y, row := range table.Rows {
			table.CellBgColors[y] = make([]render.Attribute, len(row))
		}
	}

	cellWidths := make([]int, len(table.Rows[0]))

	for y, row := range table.Rows {
		if table.FgColors[y] == 0 {
			table.FgColors[y] = table.FgColor
			for x := range row {
				table.CellFgColors[y][x] = table.FgColor
			}
		}
		if table.BgColors[y] == 0 {
			table.BgColors[y] = table.BgColor
			for x := range row {
				table.CellBgColors[y][x] = table.BgColor
			}
		}
		for x, str := range row {
			cells := render.DefaultTxBuilder.Build(str, table.CellFgColors[y][x], table.CellBgColors[y][x])
			cw := cellsWidth(cells)
			if cellWidths[x] < cw {
				cellWidths[x] = cw
			}
			rowCells = append(rowCells, cells)
		}
	}
	table.CellWidth = cellWidths
	return rowCells
}

// Analysis generates and returns an array of []Cell that represent all columns in the Table
func (table *Table) Analysis() [][]render.Cell {
	table.Lock()
	defer table.Unlock()
	return table.analysis()
}

// SetSize calculates the table size and sets the internal value
func (table *Table) SetSize() {
	table.Lock()
	defer table.Unlock()
	length := len(table.Rows)
	if table.Separator {
		table.Height = length*2 + 1
	} else {
		table.Height = length + 2
	}
	table.Width = 2
	if length != 0 {
		for _, cellWidth := range table.CellWidth {
			table.Width += cellWidth + 3
		}
	}
}

func (table *Table) calculatePosition(x int, y int, coordinateX *int, coordinateY *int, cellStart *int) {
	if table.Separator {
		*coordinateY = table.InnerArea().Min.Y + y*2
	} else {
		*coordinateY = table.InnerArea().Min.Y + y
	}
	if x == 0 {
		*cellStart = table.InnerArea().Min.X
	} else {
		*cellStart += table.CellWidth[x-1] + 3
	}

	switch table.TextAlign {
	case align.AlignRight:
		*coordinateX = *cellStart + (table.CellWidth[x] - len(table.Rows[y][x])) + 2
	case align.AlignCenter:
		*coordinateX = *cellStart + (table.CellWidth[x]-len(table.Rows[y][x]))/2 + 2
	default:
		*coordinateX = *cellStart + 2
	}
}

func (table *Table) calculateCell(x int, y int) (cx, cy int) {
	tMin := table.InnerArea().Min
	tMax := table.InnerArea().Max

	if x < tMin.X || x > tMax.X || y < tMin.Y || y > tMax.Y {
		return -1, -1
	}

	cy = y - tMin.Y
	if table.Separator {
		cy = cy / 2
	}

	cw := table.CellWidth
	cx = 0
	for cx < len(cw) && x-cw[cx] >= 0 {
		x -= cw[cx] + 2
		cx += 1
	}

	return cx, cy

	/*
		if x == 0 {
			*cellStart = table.InnerArea().Min.X
		} else {
			*cellStart += table.CellWidth[x-1] + 3
		}

		switch table.TextAlign {
		case AlignRight:
			*coordinateX = *cellStart + (table.CellWidth[x] - len(table.Rows[y][x])) + 2
		case AlignCenter:
			*coordinateX = *cellStart + (table.CellWidth[x]-len(table.Rows[y][x]))/2 + 2
		default:
			*coordinateX = *cellStart + 2
		}
	*/
}

// CalculatePosition ...
func (table *Table) CalculatePosition(x int, y int, coordinateX *int, coordinateY *int, cellStart *int) {
	table.Lock()
	defer table.Unlock()
	table.calculatePosition(x, y, coordinateX, coordinateY, cellStart)
}

// CalculateCell given position ...
func (table *Table) CalculateCell(x int, y int) (cx, cy int) {
	table.Lock()
	defer table.Unlock()
	return table.calculateCell(x, y)
}

// Buffer ...
func (table *Table) Buffer() render.Buffer {
	table.Lock()
	defer table.Unlock()
	buffer := table.Block.Buffer()
	rowCells := table.analysis()
	pointerX := table.InnerArea().Min.X + 2
	pointerY := table.InnerArea().Min.Y
	borderPointerX := table.InnerArea().Min.X
	for y, row := range table.Rows {
		for x := range row {
			table.calculatePosition(x, y, &pointerX, &pointerY, &borderPointerX)
			background := render.DefaultTxBuilder.Build(strings.Repeat(" ", table.CellWidth[x]+3), table.BgColors[y], table.BgColors[y])
			cells := rowCells[y*len(row)+x]
			for i, back := range background {
				buffer.Set(borderPointerX+i, pointerY, back)
			}

			coordinateX := pointerX
			for _, printer := range cells {
				buffer.Set(coordinateX, pointerY, printer)
				coordinateX += printer.Width()
			}

			if x != 0 {
				dividors := render.DefaultTxBuilder.Build("|", table.FgColors[y], table.BgColors[y])
				for _, dividor := range dividors {
					buffer.Set(borderPointerX, pointerY, dividor)
				}
			}
		}

		if table.Separator {
			border := render.DefaultTxBuilder.Build(strings.Repeat("─", table.Width-2), table.FgColor, table.BgColor)
			for i, cell := range border {
				buffer.Set(i+1, pointerY+1, cell)
			}
		}
	}

	return buffer
}

func (table *Table) SetRows(rows [][]string) {
	table.Lock()
	defer table.Unlock()

	oldlen := len(table.Rows)
	table.Rows = rows
	nrNewRows := len(table.Rows) - oldlen
	nrNewColors := len(table.Rows) - len(table.FgColors) /* FgColors and BgColors are in sync */

	/* if there is a positive delta between the current number of colors and then number we expect allocate them */
	/* we intentionally do not deallocate unnecessary colors. They are not used and we keep them "chached" */
	/* caching avoids reallocation and it is relatively unlikely that a table starts very big, decreades a lot and stays like that */
	if nrNewColors > 0 {
		newfgs := make([]render.Attribute, nrNewColors)
		newbgs := make([]render.Attribute, nrNewColors)

		table.FgColors = append(table.FgColors, newfgs...)
		table.BgColors = append(table.BgColors, newbgs...)
	}

	/* always reset the colors of additional rows */
	if nrNewRows > 0 {
		for i := 0; i < nrNewRows; i++ {
			table.FgColors[oldlen+i] = table.FgColor
			table.BgColors[oldlen+i] = table.BgColor
		}
	}
}

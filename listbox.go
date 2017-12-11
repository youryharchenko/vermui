// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.
//
// Portions copyright 2017 Patrick Devine <patrick@immense.ly>
// Portions copyright 2017 Philipp Resch <phil@2kd.de>

package termui

// List displays []Item as its items (items are pairs of text and values),
// it has a Overflow option (default is "hidden"), when set to "hidden",
// the item exceeding List's width is truncated, but when set to "wrap",
// the overflowed text breaks into next line.

// Item is the main struct for the listbox entries
type Item struct {
	ItemVal string
	Text    string
}

// ListBox is the main struct
type ListBox struct {
	Block
	Items       []Item
	ItemFgColor Attribute
	ItemBgColor Attribute
	Selected    int
	lowerBound  int
}

// NewListBox returns a new *ListBox with current theme.
func NewListBox() *ListBox {
	l := &ListBox{Block: *NewBlock()}
	l.ItemFgColor = ThemeAttr("list.item.fg")
	l.ItemBgColor = ThemeAttr("list.item.bg")
	l.Selected = 0
	l.lowerBound = 0
	return l
}

// Buffer implements Bufferer interface.
func (l *ListBox) Buffer() Buffer {
	buf := l.Block.Buffer()

	trimItems := l.GetItemsStrs()
	totalItems := len(l.GetItemsStrs())
	if len(trimItems) > l.innerArea.Dy() {
		trimItems = trimItems[l.lowerBound : l.innerArea.Dy()+l.lowerBound]
	}
	for i, v := range trimItems {
		var cs []Cell
		if i+l.lowerBound == l.Selected {
			cs = DTrimTxCls(DefaultTxBuilder.Build(v, l.ItemBgColor, l.ItemFgColor), l.innerArea.Dx())
		} else {
			cs = DTrimTxCls(DefaultTxBuilder.Build(v, l.ItemFgColor, l.ItemBgColor), l.innerArea.Dx())
		}
		j := 0
		for _, vv := range cs {
			w := vv.Width()
			buf.Set(l.innerArea.Min.X+j, l.innerArea.Min.Y+i, vv)
			j += w
		}
	}
	// display scroll arrows
	if l.lowerBound > 0 {
		buf.Set(l.innerArea.Dx(), 1, Cell{Ch: '^'})
	}
	if totalItems > l.lowerBound+l.innerArea.Dy() {
		buf.Set(l.innerArea.Dx(), l.innerArea.Dy(), Cell{Ch: 'v'})
	}
	return buf
}

func (l *ListBox) GetItemsStrs() []string {
	var strs []string
	for _, item := range l.Items {
		strs = append(strs, item.Text)
	}
	return strs
}

// Up moves the selection one up
func (l *ListBox) Up() {
	if l.Selected > 0 {
		l.Selected--
		if l.Selected < l.lowerBound {
			l.lowerBound--
		}
	}
}

// Down moves the selection one down
func (l *ListBox) Down() {
	if l.Selected < len(l.Items)-1 {
		l.Selected++
		if l.Selected >= l.innerArea.Dy()+l.lowerBound {
			l.lowerBound++
		}
	}
}

// Current gives the currently selected item
func (l *ListBox) Current() Item {
	// Failsafe
	if l.Selected > len(l.Items)-1 {
		l.Selected = len(l.Items) - 1
	}
	if l.Selected < 0 {
		l.Selected = 0
	}
	return l.Items[l.Selected]
}

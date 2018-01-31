package form

import (
	"github.com/gdamore/tcell"
	"github.com/verdverm/tview"
)

type Flex struct {
	*tview.Flex

	name  string
	items []FormItem
}

func NewFlex(name string) *Flex {
	F := &Flex{
		Flex:  tview.NewFlex(),
		name:  name,
		items: []FormItem{},
	}

	return F
}

func (F *Flex) Name() string {
	return F.name
}

func (F *Flex) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	return nil
}

func (F *Flex) GetValues() (values map[string]interface{}) {
	values = make(map[string]interface{})
	for _, item := range F.items {
		vals := item.GetValues()
		for field, value := range vals {
			values[field] = value
		}
	}
	return values
}

func (F *Flex) SetValues(values map[string]interface{}) {
	for _, item := range F.items {
		item.SetValues(values)
	}
}

func (F *Flex) GetItem(name string) FormItem {
	for _, item := range F.items {
		if item.Name() == name {
			return item
		}
	}

	return nil
}

func (F *Flex) GetItems() []FormItem {
	items := []FormItem{}
	for _, item := range F.items {
		switch typ := item.(type) {
		case *Flex:
			itms := typ.GetItems()
			items = append(items, itms...)
		default:
			items = append(items, item)
		}
	}

	return items
}

func (F *Flex) AddItem(item FormItem, fixedSize, proportion int) {
	F.items = append(F.items, item)

	F.Flex.AddItem(item, fixedSize, proportion, true)
}

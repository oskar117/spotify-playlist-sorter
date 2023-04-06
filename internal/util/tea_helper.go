package util

import "github.com/charmbracelet/bubbles/list"

func ConvertToListEntry[T list.Item](items []*T) []list.Item {
	listItems := make([]list.Item, len(items))
	for i, v := range items {
		listItems[i] = list.Item(*v)
	}
	return listItems
}

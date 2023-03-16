package tui

import "github.com/charmbracelet/bubbles/list"

type activeView int

const (
	loadingView activeView = iota
	playlistView
	sorterView
)

type model struct {
	activeView activeView

	playlistList list.Model
}

func New() *model {
	return nil	
}

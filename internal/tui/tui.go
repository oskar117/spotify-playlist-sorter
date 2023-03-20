package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/oskar117/spotify-playlist-sorter/internal/spotify"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/sorter"
)

type activeView int

const (
	loadingView activeView = iota
	playlistView
	sorterView
)

type model struct {
	client *spotify.SpotifyClient

	activeView activeView

	playlistList     list.Model
	selectedPlaylist sorter_model.Playlist

	sorterView sorter.Model
}

func convertPlaylistsToListEntry(playlists []*sorter_model.Playlist) []list.Item {
	listItems := make([]list.Item, len(playlists))
	for i, v := range playlists {
		listItems[i] = list.Item(*v)
	}
	return listItems
}

func New() *model {
	client := spotify.New()
	playlists := client.GetPlaylistsFirstPage()

	delegate := list.NewDefaultDelegate()
	list := list.New(convertPlaylistsToListEntry(playlists), delegate, 0, 0)
	list.Title = "Select playlist"
	list.SetShowHelp(false)

	return &model{
		playlistList: list,
		activeView:   playlistView,
		sorterView:   sorter.InitialModel(client),
		client:       client,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.playlistList.SetSize(h, v)
		m.sorterView, cmd = m.sorterView.Update(msg)
		return m, cmd
	}
	switch m.activeView {
	case sorterView:
		m.sorterView, cmd = m.sorterView.Update(msg)
	case playlistView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.selectedPlaylist = m.playlistList.SelectedItem().(sorter_model.Playlist)
				m.client.SetSelectedPlaylist(m.selectedPlaylist.ID)
				m.activeView = sorterView
				cmd = m.sorterView.FetchArtists()
				return m, cmd
			}
		}
		m.playlistList, cmd = m.playlistList.Update(msg)
	default:
		panic("Unknown view value!")
	}
	return m, cmd
}

func (m model) View() string {
	switch m.activeView {
	case loadingView:
		return "todo"
	case sorterView:
		return m.sorterView.View()
	case playlistView:
		return m.playlistList.View()
	default:
		panic("Unknown view value!")
	}
}

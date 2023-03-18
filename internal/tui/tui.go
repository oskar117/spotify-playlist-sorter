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
	activeView activeView

	playlistList list.Model
	sorterView   sorter.Model
}

func New() *model {
	client := spotify.New()
	playlistPage := client.GetPlaylistsFirstPage()

	var artists []*sorter_model.Artist

	for _, playlist := range playlistPage.Playlists {
		if playlist.Name == "asdf" {
			client.SetSelectedPlaylist(playlist.ID)
			artists = client.FetchArtists()
		}
	}
	return &model{
		activeView: sorterView,
		sorterView: sorter.InitialModel(artists, client),
	}
}

func (m model) Init() tea.Cmd {
	return m.sorterView.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.sorterView.Update(msg)
}

func (m model) View() string {
	return m.sorterView.View()
}

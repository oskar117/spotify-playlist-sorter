package tui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/oskar117/spotify-playlist-sorter/internal/spotify"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/command"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/sorter"
	"github.com/oskar117/spotify-playlist-sorter/internal/util"
)

var marginStyle = lipgloss.NewStyle().Margin(1)

type activeView int

const (
	playlistView activeView = iota
	sorterView
)

type model struct {
	client *spotify.SpotifyClient

	activeView activeView

	playlistList     list.Model
	selectedPlaylist sorter_model.Playlist

	sorterView sorter.Model

	spinner         spinner.Model
	loading         bool
	loadingMessages []string
}

func New() *model {
	client := spotify.New()
	delegate := list.NewDefaultDelegate()
	list := list.New(nil, delegate, 0, 0)
	list.Title = "Select playlist"
	list.SetShowHelp(false)

	spinnerObj := spinner.New()
	spinnerObj.Spinner = spinner.Dot

	return &model{
		playlistList: list,
		activeView:   playlistView,
		sorterView:   sorter.InitialModel(client),
		client:       client,
		spinner:      spinnerObj,
	}
}

func (m model) Init() tea.Cmd {
	return m.FetchPlaylists()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.playlistList.SetSize(h-2, v-2)
		m.sorterView, cmd = m.sorterView.Update(msg)
		return m, cmd
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case playlistsMsg:
		m.playlistList.SetItems(util.ConvertToListEntry(msg.playlists))
		m.loading = false
		m.loadingMessages = nil
	case command.LoadingMsg:
		m.loadingMessages = append(m.loadingMessages, msg.Message)
		m.loading = true
		return m, m.spinner.Tick
	case command.StopLoadingMsg:
		m.loadingMessages = nil
		m.loading = false
	case command.GoBackMessage:
		m.activeView = playlistView
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
	if m.loading {
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		return m, tea.Batch(cmd, spinnerCmd)
	}
	return m, cmd
}

func (m model) View() string {
	if m.loading {
		return m.renderLoadingView()
	}
	switch m.activeView {
	case sorterView:
		return m.sorterView.View()
	case playlistView:
		return marginStyle.Render(m.playlistList.View())
	default:
		panic("Unknown view value!")
	}
}

func (m model) renderLoadingView() string {
	spinner := m.spinner.View()
	messagesCopy := make([]string, len(m.loadingMessages))
	lastIndex := len(m.loadingMessages) - 1
	if lastIndex > 0 {
		for i, msg := range m.loadingMessages[:lastIndex] {
			messagesCopy[i] = strings.Repeat(" ", lipgloss.Width(spinner)) + msg
		}
	}
	messagesCopy[lastIndex] = spinner + m.loadingMessages[lastIndex]
	return marginStyle.Render(strings.Join(messagesCopy, "\n"))
}

type playlistsMsg struct {
	playlists []*sorter_model.Playlist
}

func (m model) FetchPlaylists() tea.Cmd {
	return tea.Batch(command.StartLoading("Fetching playlists..."), func() tea.Msg {
		time.Sleep(2 * time.Second)
		playlists := m.client.GetPlaylistsFirstPage()
		return playlistsMsg{playlists}
	})
}

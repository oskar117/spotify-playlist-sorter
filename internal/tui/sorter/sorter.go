package sorter

import (
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/oskar117/spotify-playlist-sorter/internal/spotify"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/command"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/songgroups"
	"github.com/oskar117/spotify-playlist-sorter/internal/util"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().
			Margin(1, 2)
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("283"))
	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

type activeFocus int

const (
	listFocus      activeFocus = iota
	songGroupFocus activeFocus = iota
)

type Model struct {
	artistsList         list.Model
	songGroups          songgroups.Model
	artists             []*sorter_model.Artist
	selected            string
	activeFocus         activeFocus
	songGroupsViewWidth int
	artistListViewWidth int
	client              *spotify.SpotifyClient
}

func InitialModel(client *spotify.SpotifyClient) Model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	list := list.New(nil, delegate, 0, 0)
	list.Title = "Spotify Playlist Sorter"
	list.SetShowHelp(false)
	viewport := songgroups.New(0, 0, client)
	return Model{
		artistsList: list,
		songGroups:  viewport,
		activeFocus: listFocus,
		client:      client,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		borderHeight := blurredBorderStyle.GetHorizontalFrameSize()
		m.artistsList.SetSize(h/2, v-borderHeight)
		m.songGroups.SetSize(h/2, v-borderHeight)
		m.artistListViewWidth = int(float64(h)*0.25) - 2*blurredBorderStyle.GetVerticalFrameSize()
		m.songGroupsViewWidth = h - m.artistListViewWidth - 2*blurredBorderStyle.GetVerticalFrameSize()
		if it := m.artistsList.SelectedItem(); it != nil {
			m.songGroups.ChangeArtist(it.(sorter_model.Artist))
		}
	case list.FilterMatchesMsg:
		m.artistsList, _ = m.artistsList.Update(msg)
		if it := m.artistsList.SelectedItem(); it != nil {
			m.songGroups.ChangeArtist(it.(sorter_model.Artist))
		}
		return m, nil
	case artistsMsg:
		m.artists = msg.artists
		filterCmd := m.artistsList.SetItems(util.ConvertToListEntry(m.artists))
		if it := m.artistsList.SelectedItem(); it != nil && m.artistsList.FilterState() == list.Unfiltered {
			m.songGroups.ChangeArtist(it.(sorter_model.Artist))
		}
		return m, tea.Batch(filterCmd, command.StopLoading())
	case command.SongGroupsUpdateMsg:
		return m, m.FetchArtists()
	}
	switch m.activeFocus {
	case listFocus:
		if m.artistsList.FilterState() != list.Filtering {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.String() == "enter" {
					m.selected = m.artistsList.SelectedItem().FilterValue()
					m.activeFocus = songGroupFocus
				} else if msg.String() == "esc" && m.artistsList.FilterState() == list.Unfiltered {
					return m, command.GoBackToPlaylistSelection()
				}
			}
		}
		m.artistsList, cmd = m.artistsList.Update(msg)
		cmds = append(cmds, cmd)
		if it := m.artistsList.SelectedItem(); it != nil {
			m.songGroups.ChangeArtist(it.(sorter_model.Artist))
		}
	case songGroupFocus:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				if m.songGroups.IsGroupSelected {
					m.selected = ""
					m.songGroups.Deselect()
				} else {
					m.activeFocus = listFocus
				}
				return m, nil
			}
		}
		m.songGroups, cmd = m.songGroups.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	artistsList := m.artistsList.View()
	songGroupsView := m.songGroups.View()
	switch m.activeFocus {
	case listFocus:
		artistsList = focusedBorderStyle.Width(m.artistListViewWidth).Render(artistsList)
		songGroupsView = blurredBorderStyle.Width(m.songGroupsViewWidth).Render(songGroupsView)
	case songGroupFocus:
		artistsList = blurredBorderStyle.Width(m.artistListViewWidth).Render(artistsList)
		songGroupsView = focusedBorderStyle.Width(m.songGroupsViewWidth).Render(songGroupsView)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, artistsList, songGroupsView)
}

type artistsMsg struct {
	artists []*sorter_model.Artist
}

func (m Model) FetchArtists() tea.Cmd {
	return tea.Batch(command.StartLoading("Fetching artists..."), func() tea.Msg {
		artists := m.client.FetchArtists()
		return artistsMsg{artists}
	})
}

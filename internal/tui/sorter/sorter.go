package sorter

import (
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/oskar117/spotify-playlist-sorter/internal/spotify"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/songgroups"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type activeFocus int

const (
	listFocus      activeFocus = iota
	songGroupFocus activeFocus = iota
)

type ViewArtist struct {
	Name string
	Desc string
}

var (
	docStyle = lipgloss.NewStyle().
			Margin(1, 2)
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("283"))
	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

func (i ViewArtist) FilterValue() string { return i.Name }
func (i ViewArtist) Description() string { return i.Desc }
func (i ViewArtist) Title() string       { return i.Name }

type Model struct {
	artistsList         list.Model
	songGroups          songgroups.Model
	artists             []*sorter_model.Artist
	selected            string
	activeFocus         activeFocus
	songGroupsViewWidth int
	artistListViewWidth int
}

func convertArtistsToListEntry(artists []*sorter_model.Artist) []list.Item {
	listItems := make([]list.Item, len(artists))
	for i, v := range artists {
		listItems[i] = list.Item(*v)
	}
	return listItems
}

func InitialModel(artists []*sorter_model.Artist, client *spotify.SpotifyClient) Model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	list := list.New(convertArtistsToListEntry(artists), delegate, 0, 0)
	list.Title = "Spotify Playlist Sorter"
	list.SetShowHelp(false)
	viewport := songgroups.New(0, 0, client)
	return Model{
		artistsList: list,
		songGroups:  viewport,
		artists:     artists,
		activeFocus: listFocus,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
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
	}
	switch m.activeFocus {
	case listFocus:
		if m.artistsList.FilterState() != list.Filtering {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				if msg.String() == "enter" {
					m.selected = m.artistsList.SelectedItem().FilterValue()
					m.activeFocus = songGroupFocus
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
				m.selected = ""
				m.songGroups.Deselect()
				m.activeFocus = listFocus
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

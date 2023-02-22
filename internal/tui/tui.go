package tui

import (
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"
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

type model struct {
	artistsList         list.Model
	songGroups          songgroups.Model
	artists             map[string]*sorter.Artist
	selected            string
	activeFocus         activeFocus
	songGroupsViewWidth int
	artistListViewWidth int
}

func InitialModel(artistNames []list.Item, artists map[string]*sorter.Artist) model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	list := list.New(artistNames, delegate, 0, 0)
	list.Title = "Spotify Playlist Sorter"
	list.SetShowHelp(false)
	viewport := songgroups.New(0, 0)
	return model{
		artistsList: list,
		songGroups:  viewport,
		artists:     artists,
		activeFocus: listFocus,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "esc" {
			m.selected = ""
			m.songGroups.Deselect()
			m.activeFocus = listFocus
			if m.artistsList.FilterState() == list.Unfiltered {
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		borderHeight := blurredBorderStyle.GetHorizontalFrameSize()
		m.artistsList.SetSize(h/2, v-borderHeight)
		m.songGroups.SetSize(h/2, v-borderHeight)
		m.artistListViewWidth = int(float64(h) * 0.25) - 2*blurredBorderStyle.GetVerticalFrameSize()
		m.songGroupsViewWidth = h - m.artistListViewWidth - 2*blurredBorderStyle.GetVerticalFrameSize()
		m.songGroups.ChangeArtist(*m.artists[m.artistsList.SelectedItem().FilterValue()])
	}
	switch m.activeFocus {
	case listFocus:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.selected = m.artistsList.SelectedItem().FilterValue()
				m.activeFocus = songGroupFocus
			}
		}
		m.artistsList, cmd = m.artistsList.Update(msg)

		cmds = append(cmds, cmd)
		m.songGroups.ChangeArtist(*m.artists[m.artistsList.SelectedItem().FilterValue()])
	case songGroupFocus:
		m.songGroups, cmd = m.songGroups.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
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

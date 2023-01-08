package tui

import (
	"fmt"
	"strings"

	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
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
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))
	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

func (i ViewArtist) FilterValue() string { return i.Name }
func (i ViewArtist) Description() string { return i.Desc }
func (i ViewArtist) Title() string       { return i.Name }

type model struct {
	artistsList list.Model
	songGroups  viewport.Model
	artists     map[string]*sorter.Artist
	selected    string
	activeFocus activeFocus
}

func InitialModel(artistNames []list.Item, artists map[string]*sorter.Artist) model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetSpacing(0)
	list := list.New(artistNames, delegate, 0, 0)
	list.Title = "Spotify Playlist Sorter"
	viewport := viewport.New(0, 0)
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
		if msg.String() == "enter" {
			m.selected = m.artistsList.SelectedItem().FilterValue()
			m.activeFocus = songGroupFocus
		}
		if msg.String() == "esc" {
			m.selected = ""
			m.activeFocus = listFocus
			if m.artistsList.FilterState() == list.Unfiltered {
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.artistsList.SetSize(h/2, v-10)
		m.songGroups = viewport.New(h, v-10)
		m.songGroups.SetContent(buildViewport(*m.artists[m.artistsList.SelectedItem().FilterValue()]))
	}
	switch m.activeFocus {
	case listFocus:
		m.artistsList, cmd = m.artistsList.Update(msg)
		cmds = append(cmds, cmd)
		m.songGroups.SetContent(buildViewport(*m.artists[m.artistsList.SelectedItem().FilterValue()]))
	case songGroupFocus:
		m.songGroups, cmd = m.songGroups.Update(msg)
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func buildViewport(choosen sorter.Artist) string {
	var builder strings.Builder
	for x, group := range choosen.SongGroups {
		builder.WriteString(fmt.Sprintln("Group", x, "first index", group.First, "last index", group.Last))
		for i, song := range group.SongTitles {
			builder.WriteString(fmt.Sprintln(i+group.First, song))
		}
	}
	return builder.String()
}

func (m model) View() string {
	artistsList := m.artistsList.View()
	songGroupsView := m.songGroups.View()
	switch m.activeFocus {
	case listFocus:
		artistsList = focusedBorderStyle.Margin(2).Render(artistsList)
	case songGroupFocus:
		songGroupsView = focusedBorderStyle.Margin(2).Render(songGroupsView)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, artistsList, songGroupsView)
}

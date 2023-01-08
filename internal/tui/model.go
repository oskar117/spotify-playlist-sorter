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

type ViewArtist struct {
	Name string
	Desc string
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (i ViewArtist) FilterValue() string { return i.Name }
func (i ViewArtist) Description() string { return i.Desc }
func (i ViewArtist) Title() string       { return i.Name }

type model struct {
	artistsList list.Model
	songGroups  viewport.Model
	artists     map[string]*sorter.Artist
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
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		var cmd tea.Cmd
		m.artistsList, cmd = m.artistsList.Update(msg)
		cmds = append(cmds, cmd)
		m.songGroups.SetContent(buildViewport(*m.artists[m.artistsList.SelectedItem().FilterValue()]))
	case tea.WindowSizeMsg:
		h, v := msg.Width, msg.Height
		m.artistsList.SetSize(h/2, v)
		m.songGroups = viewport.New(h/2, v)
		m.songGroups.SetContent(buildViewport(*m.artists[m.artistsList.SelectedItem().FilterValue()]))
	}
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
	return lipgloss.JoinHorizontal(lipgloss.Top, m.artistsList.View(), m.songGroups.View())
}

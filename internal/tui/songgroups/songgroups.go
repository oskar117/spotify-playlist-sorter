package songgroups

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/oskar117/spotify-playlist-sorter/internal/spotify"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui/command"
)

var (
	highlightedText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("12")).
			Bold(true)
	selectedText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("3")).
			Bold(true)
	neutralText = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0"))
)

type Model struct {
	viewport              viewport.Model
	artist                sorter_model.Artist
	highlightedGroupIndex int
	selectedGroupIndex    int
	IsGroupSelected       bool
	moveLocation          spotify.GroupLocation
	client                *spotify.SpotifyClient
}

func New(width, height int, client *spotify.SpotifyClient) Model {
	return Model{
		viewport: viewport.New(width, height),
		client:   client,
	}
}

func (m Model) Init() tea.Cmd {
	return m.viewport.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.viewport.KeyMap.Up):
			if !m.IsGroupSelected && m.highlightedGroupIndex > 0 {
				m.highlightedGroupIndex--
				group := len(m.artist.SongGroups[m.highlightedGroupIndex].SongTitles) - 1
				m.viewport.LineUp(group)
			} else if m.selectedGroupIndex == m.highlightedGroupIndex && m.selectedGroupIndex > 0 {
				m.selectedGroupIndex--
				m.moveLocation = spotify.Bottom
				group := len(m.artist.SongGroups[m.selectedGroupIndex].SongTitles) - 1
				m.viewport.LineUp(group)
			} else if m.selectedGroupIndex > 0 || m.moveLocation == spotify.Bottom {
				if m.moveLocation == spotify.Bottom {
					m.moveLocation = spotify.Top
					group := len(m.artist.SongGroups[m.selectedGroupIndex].SongTitles) - 1
					m.viewport.LineUp(group)
				} else {
					m.moveLocation = spotify.Bottom
					m.selectedGroupIndex--
				}
			}
		case key.Matches(msg, m.viewport.KeyMap.Down):
			if !m.IsGroupSelected && m.highlightedGroupIndex < len(m.artist.SongGroups)-1 {
				group := len(m.artist.SongGroups[m.highlightedGroupIndex].SongTitles) - 1
				m.viewport.LineDown(group)
				m.highlightedGroupIndex++
			} else if m.selectedGroupIndex == m.highlightedGroupIndex && m.selectedGroupIndex < len(m.artist.SongGroups)-1 {
				group := len(m.artist.SongGroups[m.selectedGroupIndex].SongTitles) - 1
				m.viewport.LineDown(group)
				m.selectedGroupIndex++
				m.moveLocation = spotify.Top
			} else if m.selectedGroupIndex < len(m.artist.SongGroups)-1 || m.moveLocation == spotify.Top {
				if m.moveLocation == spotify.Bottom {
					m.moveLocation = spotify.Top
					m.selectedGroupIndex++
				} else {
					m.moveLocation = spotify.Bottom
					group := len(m.artist.SongGroups[m.selectedGroupIndex].SongTitles) - 1
					m.viewport.LineDown(group)
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.IsGroupSelected {
				m.IsGroupSelected = false
				m.client.ReorderGroups(m.artist.SongGroups[m.highlightedGroupIndex], m.artist.SongGroups[m.selectedGroupIndex], m.moveLocation)
				return m, command.UpdateSongGroups()
			} else {
				m.IsGroupSelected = true
				m.selectedGroupIndex = m.highlightedGroupIndex
			}
		}
	}
	m.SetContent(m.buildContent())
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m *Model) SetContent(lines string) {
	m.viewport.SetContent(lines)
}

func (m *Model) SetSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height
}

func (m *Model) ChangeArtist(artist sorter_model.Artist) {
	m.artist = artist
	m.highlightedGroupIndex = 0
	m.viewport.GotoTop()
	m.SetContent(m.buildContent())
}

func (m Model) Width() int {
	return m.viewport.Width
}

func (m *Model) Deselect() {
	m.IsGroupSelected = false
	m.selectedGroupIndex = 0
	m.SetContent(m.buildContent())
}

func (m Model) buildContent() string {
	var builder strings.Builder
	groupModels := convertToModel(m.artist)
	if m.IsGroupSelected && m.selectedGroupIndex != m.highlightedGroupIndex {
		switch m.moveLocation {
		case spotify.Top:
			groupModels.mergeOnTop(m.highlightedGroupIndex, m.selectedGroupIndex)
		case spotify.Bottom:
			groupModels.mergeAtBottom(m.highlightedGroupIndex, m.selectedGroupIndex)
		}
	}
	for x, group := range groupModels {
		var localGroupBuilder strings.Builder
		localStyle := lipgloss.NewStyle().Inline(true)
		if m.IsGroupSelected && group.index == m.selectedGroupIndex {
			localStyle = localStyle.Inherit(selectedText)
		} else if !m.IsGroupSelected && group.index == m.highlightedGroupIndex {
			localStyle = localStyle.Inherit(highlightedText)
		}
		localGroupBuilder.WriteString(localStyle.Render(fmt.Sprintln("Group", x, "first index", group.first, "last index", group.last)))
		localGroupBuilder.WriteString("\n")
		for songIndex, song := range group.songs {
			localGroupBuilder.WriteString(localStyle.Render(fmt.Sprintln(songIndex+group.first, song.name)))
			localGroupBuilder.WriteString("\n")
		}
		localGroupBuilder.WriteString("\n")
		builder.WriteString(localGroupBuilder.String())
	}
	return builder.String()
}

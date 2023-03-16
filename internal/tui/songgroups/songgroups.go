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
	spotify_api "github.com/zmb3/spotify/v2"
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
	isSelected            bool
	moveLocation          spotify.GroupLocation
	playlistId			  spotify_api.ID
	client				  *spotify_api.Client
}

func New(width, height int, playlistId spotify_api.ID, client *spotify_api.Client) Model {
	return Model{
		viewport: viewport.New(width, height),
		playlistId: playlistId,
		client: client,
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
			if !m.isSelected && m.highlightedGroupIndex > 0 {
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
			if !m.isSelected && m.highlightedGroupIndex < len(m.artist.SongGroups)-1 {
				group := len(m.artist.SongGroups[m.highlightedGroupIndex].SongTitles) - 1
				m.viewport.LineDown(group)
				m.highlightedGroupIndex++
			} else if m.selectedGroupIndex == m.highlightedGroupIndex  && m.selectedGroupIndex < len(m.artist.SongGroups)-1 {
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
			if m.isSelected {
				m.isSelected = false
				// spotify.ReorderGroups(m.client, m.playlistId, m.artist.SongGroups[m.highlightedGroupIndex], m.artist.SongGroups[m.selectedGroupIndex], m.moveLocation)
			} else {
				m.isSelected = true
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
	m.isSelected = false
	m.selectedGroupIndex = 0
	m.SetContent(m.buildContent())
}

func (m Model) buildContent() string {
	var builder strings.Builder
	groupModels := convertToModel(m.artist)
	if m.isSelected && m.selectedGroupIndex != m.highlightedGroupIndex {
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
		if m.isSelected && group.index == m.selectedGroupIndex {
			localStyle = localStyle.Inherit(selectedText)
		} else if !m.isSelected && group.index == m.highlightedGroupIndex {
			localStyle = localStyle.Inherit(highlightedText)
		}
		localGroupBuilder.WriteString(localStyle.Render(fmt.Sprintln("Group", x, "first index", group.first, "last index", group.last)))
		localGroupBuilder.WriteString("\n")
		for songIndex, song := range group.songs {
			localGroupBuilder.WriteString(localStyle.Render(fmt.Sprintln(songIndex + group.first, song.name)))
			localGroupBuilder.WriteString("\n")
		}
		localGroupBuilder.WriteString("\n")
		builder.WriteString(localGroupBuilder.String())
	}
	return builder.String()
}

package songgroups

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"
)

type Model struct {
	viewport            viewport.Model
	artist              sorter.Artist
	selectedArtistIndex int
}

func New(width, height int) Model {
	return Model{
		viewport: viewport.New(width, height),
	}
}

func (m Model) Init() tea.Cmd {
	return m.viewport.Init()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m *Model) SetContent(lines string) {
	m.viewport.SetContent(lines)
	m.viewport.GotoTop()
}

func (m *Model) SetSize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height
}

func (m *Model) ChangeArtist(artist sorter.Artist) {
	m.artist = artist
	m.SetContent(buildViewport(artist))
}

func (m Model) Width() int {
	return m.viewport.Width
}

func buildViewport(choosen sorter.Artist) string {
	var builder strings.Builder
	for x, group := range choosen.SongGroups {
		builder.WriteString(fmt.Sprintln("Group", x, "first index", group.First, "last index", group.Last))
		for i, song := range group.SongTitles {
			builder.WriteString(fmt.Sprintln(i+group.First, song))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}

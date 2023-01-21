package songgroups

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"
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

type groupLocation int

const (
	top groupLocation = iota
	bottom
)

type Model struct {
	viewport              viewport.Model
	artist                sorter.Artist
	highlightedGroupIndex int
	selectedGroupIndex    int
	isSelected            bool
	moveLocation          groupLocation
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.viewport.KeyMap.Up):
			if m.highlightedGroupIndex > 0 {
				if !m.isSelected {
					m.highlightedGroupIndex--
				} else {
					if m.moveLocation == bottom || m.selectedGroupIndex == m.highlightedGroupIndex {
						m.moveLocation = top
						m.selectedGroupIndex--
					} else {
						m.moveLocation = bottom
					}
				}
			}
		case key.Matches(msg, m.viewport.KeyMap.Down):
			if m.highlightedGroupIndex < len(m.artist.SongGroups)-1 {
				if !m.isSelected {
					m.highlightedGroupIndex++
				} else {
					if m.moveLocation == bottom && m.selectedGroupIndex != m.highlightedGroupIndex {
						m.moveLocation = top
					} else {
						m.moveLocation = bottom
						m.selectedGroupIndex++
					}
				}
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.isSelected {
				m.isSelected = false
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

func (m *Model) ChangeArtist(artist sorter.Artist) {
	m.artist = artist
	m.highlightedGroupIndex = 0
	m.viewport.GotoTop()
	m.SetContent(m.buildContent())
}

func (m Model) Width() int {
	return m.viewport.Width
}

func (m Model) buildContent() string {
	var builder strings.Builder
	for x, group := range m.artist.SongGroups {
		var localGroupBuilder strings.Builder
		localStyle := lipgloss.NewStyle().Inline(true)
		if m.isSelected && x == m.selectedGroupIndex {
			localStyle = localStyle.Inherit(selectedText)
		} else if !m.isSelected && x == m.highlightedGroupIndex {
			localStyle = localStyle.Inherit(highlightedText)
		}
		localGroupBuilder.WriteString(localStyle.Render(fmt.Sprintln("Group", x, "first index", group.First, "last index", group.Last)))
		localGroupBuilder.WriteString("\n")
		for i, song := range group.SongTitles {
			localGroupBuilder.WriteString(localStyle.Render(fmt.Sprintln(i+group.First, song)))
			localGroupBuilder.WriteString("\n")
		}
		localGroupBuilder.WriteString("\n")
		builder.WriteString(localGroupBuilder.String())
	}
	return builder.String()
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

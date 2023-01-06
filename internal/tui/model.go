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
	choices     []string         // items on the to-do list
	cursor      int              // which to-do list item our cursor is pointing at
	selected    map[int]struct{} // which to-do items are selected
	artistsList list.Model
	songGroups  viewport.Model
	artists     map[string]*sorter.Artist
}

func InitialModel(artistNames []list.Item, artists map[string]*sorter.Artist) model {
	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	list := list.New(artistNames, delegate, 0, 0)
	list.Title = "Spotify Playlist Sorter"
	viewport := viewport.New(0, 0)
	viewport.SetContent("Test")
	// list.SetShowHelp(false)
	// list.SetShowStatusBar(false)
	// list.SetShowFilter(false)
	// list.SetFilteringEnabled(false)
	return model{
		// Our to-do list is a grocery list
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected:    make(map[int]struct{}),
		artistsList: list,
		songGroups:  viewport,
		artists:     artists,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {

	// Is it a key press?
	// case tea.KeyMsg:
	//
	// 	// Cool, what was the actual key pressed?
	// 	switch msg.String() {
	//
	// 	// These keys should exit the program.
	// 	case "ctrl+c", "q":
	// 		return m, tea.Quit
	//
	// 	// The "up" and "k" keys move the cursor up
	// 	case "up", "k":
	// 		if m.cursor > 0 {
	// 			m.cursor--
	// 		}
	//
	// 	// The "down" and "j" keys move the cursor down
	// 	case "down", "j":
	// 		if m.cursor < len(m.choices)-1 {
	// 			m.cursor++
	// 		}
	//
	// 	// The "enter" key and the spacebar (a literal space) toggle
	// 	// the selected state for the item that the cursor is pointing at.
	// 	case "enter", " ":
	// 		_, ok := m.selected[m.cursor]
	// 		if ok {
	// 			delete(m.selected, m.cursor)
	// 		} else {
	// 			m.selected[m.cursor] = struct{}{}
	// 		}
	// 	}
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
	// var cmd tea.Cmd
	// m.artistsList, cmd = m.artistsList.Update(msg)
	// cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	// return m, nil
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
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"
	return lipgloss.JoinHorizontal(lipgloss.Top, m.artistsList.View(), m.songGroups.View())
}

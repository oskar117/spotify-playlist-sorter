package tui

import tea "github.com/charmbracelet/bubbletea"

func (m model) testCmd() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

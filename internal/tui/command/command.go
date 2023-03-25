package command

import tea "github.com/charmbracelet/bubbletea"

type LoadingMsg struct {
	Message string
}

func StartLoading(message string) tea.Cmd {
	return func() tea.Msg {
		return LoadingMsg{message}
	}
}

type StopLoadingMsg struct {}

func StopLoading() tea.Cmd {
	return func() tea.Msg {
		return StopLoadingMsg{}
	}
}

type SongGroupsUpdateMsg struct {}

func UpdateSongGroups() tea.Cmd {
	return func() tea.Msg {
		return SongGroupsUpdateMsg{}
	}
}

type GoBackMessage struct {}

func GoBackToPlaylistSelection() tea.Cmd {
	return func() tea.Msg {
		return GoBackMessage{}
	}
}

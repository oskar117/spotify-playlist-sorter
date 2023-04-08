package program

import (
	"flag"
	"fmt"
	"io"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oskar117/spotify-playlist-sorter/internal/tui"
)

type Program struct {
	debugEnabled bool
}

func New() *Program {
	debugFlag := flag.Bool("debug", false, "Enable logging to debug.log file.")
	flag.Parse()
	return &Program{
		debugEnabled: *debugFlag,
	}
}

func (p *Program) Start() error {
	if p.debugEnabled {
		f, err := tea.LogToFile("debug.log", "debug")
		defer f.Close()
		if err != nil {
			return fmt.Errorf("Failed to enable log file: %w", err)
		}
	} else {
		log.SetOutput(io.Discard)
	}
	if _, err := tea.NewProgram(tui.New(), tea.WithAltScreen()).Run(); err != nil {
		return fmt.Errorf("Alas, there's been an error while starting ui: %w", err)
	}
	return nil
}

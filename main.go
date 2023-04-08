package main

import (
	"fmt"
	"os"

	"github.com/oskar117/spotify-playlist-sorter/internal/program"
)

func main() {
	if err := program.New().Start(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}


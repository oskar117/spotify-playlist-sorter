package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/oskar117/spotify-playlist-sorter/internal/tui/sorter"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/oskar117/spotify-playlist-sorter/internal/auth"
	loc_spotify "github.com/oskar117/spotify-playlist-sorter/internal/spotify"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zmb3/spotify/v2"
)

func main() {
	client := spotify.New(auth.GetHttpClient())

	// use the client to make calls that require authorization
	spotifyUser, err := client.CurrentUser(context.Background())
	if err != nil {
		auth.RemoveTokenFromKeyring()
		log.Fatal(err)
	}

	fmt.Println("You are logged in as:", spotifyUser.ID)
	playlistPage, _ := client.GetPlaylistsForUser(context.Background(), spotifyUser.ID)

	var playlistId spotify.ID
	var artists []*sorter_model.Artist

	for _, playlist := range playlistPage.Playlists {
		if playlist.Owner.ID == spotifyUser.ID && playlist.Name == "asdf" {
			playlistId = playlist.ID
			artists = loc_spotify.FetchArtists(client, playlistId)
		}
	}
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	p := tea.NewProgram(sorter.InitialModel(artists, playlistId, client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	os.Exit(6)
}


package spotify

import (
	"context"
	"fmt"
	"log"

	"github.com/oskar117/spotify-playlist-sorter/internal/auth"
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter_model"
	"github.com/zmb3/spotify/v2"
)

type GroupLocation int

const (
	Top GroupLocation = iota
	Bottom
)

type SpotifyClient struct {
	playlistId spotify.ID
	userId     string
	spotifyApi *spotify.Client
}

func New() *SpotifyClient {
	client := spotify.New(auth.GetHttpClient())
	spotifyUser, err := client.CurrentUser(context.Background())
	if err != nil {
		auth.RemoveTokenFromKeyring()
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", spotifyUser.ID)
	return &SpotifyClient{userId: spotifyUser.ID, spotifyApi: client}
}

func (client *SpotifyClient) SetSelectedPlaylist(playlistId spotify.ID) {
	client.playlistId = playlistId
}

func (client *SpotifyClient) GetPlaylistsFirstPage() []*sorter_model.Playlist {
	playlistPage, _ := client.spotifyApi.GetPlaylistsForUser(context.Background(), client.userId)
	playlists := make([]*sorter_model.Playlist, 0)
	for _, playlist := range playlistPage.Playlists {
		if playlist.Owner.ID != client.userId {
			continue
		}
		playlists = append(playlists, &sorter_model.Playlist{
			Name: playlist.Name, 
			Desc: playlist.Description, 
			ID: playlist.ID,
		})
	}
	return playlists
}

func (client *SpotifyClient) ReorderGroups(from, to *sorter_model.SongGroup, location GroupLocation) error {
	targetIndex := func() int {
		if location == Top {
			return to.First
		}
		return to.Last + 1
	}()
	options := spotify.PlaylistReorderOptions{RangeStart: from.First, RangeLength: len(from.SongTitles), InsertBefore: targetIndex}
	_, error := client.spotifyApi.ReorderPlaylistTracks(context.Background(), client.playlistId, options)
	return error
}

func (client *SpotifyClient) FetchArtists() []*sorter_model.Artist {
	artists := make([]*sorter_model.Artist, 0)
	firstItemsPage, _ := client.spotifyApi.GetPlaylistItems(context.Background(), client.playlistId)
	items := firstItemsPage.Items
	for firstItemsPage.Next != "" {
		client.spotifyApi.NextPage(context.Background(), firstItemsPage)
		items = append(items, firstItemsPage.Items...)
	}
	for index, item := range items {
		artistName := item.Track.Track.Artists[0].Name
		artistIndex := findArtistIndex(artists, artistName)
		if artistIndex < 0 {
			artists = append(artists, &sorter_model.Artist{Name: artistName, SongGroups: make([]*sorter_model.SongGroup, 0)})
			artistIndex = len(artists) - 1
		}
		artists[artistIndex].AddSong(item.Track.Track.Name, index)
	}
	return artists
}

func findArtistIndex(artists []*sorter_model.Artist, name string) int {
	for i, artist := range artists {
		if artist.Name == name {
			return i
		}
	}
	return -1
}

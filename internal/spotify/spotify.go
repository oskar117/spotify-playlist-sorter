package spotify

import (
	"context"

	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"
	"github.com/zmb3/spotify/v2"
)

type GroupLocation int

const (
	Top GroupLocation = iota
	Bottom
)

func ReorderGroups(client *spotify.Client, playlistId spotify.ID, from, to *sorter.SongGroup, location GroupLocation) error {
	targetIndex := func() int {
		if location == Top {
			return to.First
		}
		return to.Last
	}()
	options := spotify.PlaylistReorderOptions{RangeStart: from.First, RangeLength: len(from.SongTitles), InsertBefore: targetIndex + 1}
	_, error := client.ReorderPlaylistTracks(context.Background(), playlistId, options)
	return error
}

func FetchArtists(client *spotify.Client, playlistId spotify.ID) []*sorter.Artist {
	artists := make([]*sorter.Artist, 0)
	firstItemsPage, _ := client.GetPlaylistItems(context.Background(), playlistId)
	items := firstItemsPage.Items
	for firstItemsPage.Next != "" {
		client.NextPage(context.Background(), firstItemsPage)
		items = append(items, firstItemsPage.Items...)
	}
	for index, item := range items {
		artistName := item.Track.Track.Artists[0].Name
		artistIndex := findArtistIndex(artists, artistName)
		if artistIndex < 0 {
			artists = append(artists, &sorter.Artist{Name: artistName, SongGroups: make([]*sorter.SongGroup, 0)}) 
			artistIndex = len(artists) - 1
		}
		artists[artistIndex].AddSong(item.Track.Track.Name, index)
	}
	return artists
}

func findArtistIndex(artists []*sorter.Artist, name string) int {
	for i, artist := range artists {
		if artist.Name == name {
			return i
		}
	}
	return -1
}

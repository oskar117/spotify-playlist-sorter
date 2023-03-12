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
	options := spotify.PlaylistReorderOptions{RangeStart: from.First, RangeLength: len(from.SongTitles), InsertBefore: targetIndex+1}
	_, error := client.ReorderPlaylistTracks(context.Background(), playlistId, options)
	return error
}

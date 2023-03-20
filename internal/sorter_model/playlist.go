package sorter_model

import "github.com/zmb3/spotify/v2"

type Playlist struct {
	ID	 spotify.ID
	Name string
	Desc string
}

func (playlist Playlist) Title() string {
	return playlist.Name
}

func (playlist Playlist) Description() string {
	return playlist.Desc
}

func (playlist Playlist) FilterValue() string {
	return playlist.Name
}

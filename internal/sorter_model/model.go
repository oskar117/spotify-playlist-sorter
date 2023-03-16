package sorter_model

type SongGroup struct {
	First, Last int
	SongTitles  []string
}

type ViewArtist interface {
	FilterValue() string
	Description() string
	Title() string      
}

type Artist struct {
	Name       string
	SongGroups []*SongGroup
}

func (artist *Artist) AddSong(title string, index int) {
	if len(artist.SongGroups) > 0 {
		lastGroup := artist.SongGroups[len(artist.SongGroups)-1]
		if index-lastGroup.Last == 1 {
			lastGroup.Last++
			lastGroup.SongTitles = append(lastGroup.SongTitles, title)
			return
		}
	}
	artist.SongGroups = append(artist.SongGroups, &SongGroup{index, index, []string{title}})
}

func (artist Artist) Title() string {
	return artist.Name
}

func (artist Artist) Description() string {
	return ""
}

func (artist Artist) FilterValue() string {
	return artist.Name
}

func (group *SongGroup) instertAtEnd(songGroup SongGroup) {
	group.SongTitles = append(group.SongTitles, songGroup.SongTitles...)
}

func (group *SongGroup) instertAtBeginning(songGroup SongGroup) {
	group.SongTitles = append(songGroup.SongTitles, group.SongTitles...)
}

package sorter

type SongGroup struct {
	First, Last int
	SongTitles  []string
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

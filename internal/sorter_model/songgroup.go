package sorter_model

type SongGroup struct {
	First, Last int
	SongTitles  []string
}

func (group *SongGroup) instertAtEnd(songGroup SongGroup) {
	group.SongTitles = append(group.SongTitles, songGroup.SongTitles...)
}

func (group *SongGroup) instertAtBeginning(songGroup SongGroup) {
	group.SongTitles = append(songGroup.SongTitles, group.SongTitles...)
}

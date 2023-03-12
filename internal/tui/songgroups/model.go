package songgroups

import (
	"github.com/oskar117/spotify-playlist-sorter/internal/sorter"
)

type songGroups []songGroupModel

type songGroupModel struct {
	index, first, last int
	songs  []songModel
}

type songModel struct {
	index int
	name  string
}

func convertToModel(artist sorter.Artist) songGroups {
	result := make([]songGroupModel, 0)
	for index, group := range artist.SongGroups {
		groupResult := make([]songModel, 0)
		for songIndex, song := range group.SongTitles {
			groupResult = append(groupResult, songModel{songIndex + group.First, song})
		}
		result = append(result, songGroupModel{index, group.First, group.Last, groupResult})
	}
	return result
}

func (group *songGroups) mergeOnTop(from, to int) {
	sourceGroup := (*group)[from]
	targetGroup := &(*group)[to]
	targetGroup.songs = append(sourceGroup.songs, targetGroup.songs...)
	targetGroup.last += len(sourceGroup.songs)
	*group = append((*group)[:from], (*group)[from+1:]...)
}

func (group *songGroups) mergeAtBottom(from, to int) {
	sourceGroup := (*group)[from]
	targetGroup := &(*group)[to]
	for i, song := range sourceGroup.songs {
		targetGroup.songs = append(targetGroup.songs, songModel{targetGroup.last + i + 1, song.name})
	}
	targetGroup.last += len(sourceGroup.songs)
	*group = append((*group)[:from], (*group)[from+1:]...)
}

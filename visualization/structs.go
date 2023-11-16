package visualization

import (
	"embed"
	"groupie-tracker/get_info"
)

type DataForFilter struct {
	MinYear     int
	MaxYear     int
	TeamMembers int
	Location    string
}

type ViewData struct {
	Artists         []get_info.Artist
	LocCoords       map[string][][]float64
	LenArtists      int
	Name            string
	MaxMembers      int
	MaxMemberOption []int
	MinCreationDate int
	MaxCreationDate int
	MinYearAlbum    int
	MaxYearAlbum    int
	ListOfLocs      []string
	SelectedLoc     string
}

var (
	//go:embed htmlTemplates/*
	viewTemplate embed.FS
)

var (
	//go:embed mapTempl/*
	viewMaps embed.FS
)

type DataForFilterFunc struct {
	TeamNumFromValue []int
	MinAlbDate       int
	MaxAlbDate       int
	MaxCreationDate  int
	MinCreationDate  int
	SearchQuery      string
}
type MapStruct struct {
	ConcertData map[string][][]float64
	InitialLoc  [][]float64
}

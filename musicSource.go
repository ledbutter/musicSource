package musicSource

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type SearchCriteria struct {
	ArtistName string
	AlbumName  string
}

type Searcher interface {
	Search() string
}

type MusicBrainzSearcher struct {
	Criteria SearchCriteria
}

func (mbSearcher *MusicBrainzSearcher) Search() string {
	query := "http://musicbrainz.org/ws/2/artist?query=" + mbSearcher.Criteria.ArtistName + "&fmt=json"
	fmt.Println(query)
	resp, err := http.Get(query)
	if err != nil {
		return "Uh Oh, Get()"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Uh oh, ReadAllBytes()"
	}
	content := string(body)
	return content
}

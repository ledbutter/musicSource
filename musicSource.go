package musicSource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
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

type Artist struct {
	Id   string
	Name string
}

func (mbSearcher *MusicBrainzSearcher) Search() string {
	val := url.Values{}
	val.Add("query", mbSearcher.Criteria.ArtistName)
	val.Add("fmt", "json")
	query := "http://musicbrainz.org/ws/2/artist?" + val.Encode()
	resp, err := http.Get(query)
	if err != nil {
		return "Uh Oh, Get()"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "Uh oh, ReadAllBytes()"
	}

	var result struct {
		Artists [1]Artist
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "Uh oh, Unmarshal"
	}
	return result.Artists[0].Name
}

package musicSource

import (
	"encoding/json"
	//	"fmt"
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

	var topLevel interface{}
	err = json.Unmarshal(body, &topLevel)
	if err != nil {
		return "Uh oh, Unmarshal"
	} else {
		m := topLevel.(map[string]interface{}) //this seems really hacky/brittle, there has to be a better way?
		result := (m["artists"].([]interface{})[0]).(map[string]interface{})
		artist := new(Artist)
		artist.Id = result["id"].(string)
		artist.Name = result["name"].(string)
		return artist.Name
	}
}

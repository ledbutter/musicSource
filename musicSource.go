package musicSource

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// Represents the search criteria needed to search for an album.
type SearchCriteria struct {
	ArtistName string
	AlbumName  string
}

// Interface for searching.
// Only included so I can say I tried using an interface.
// Also allows for more simple substitution of the MusicBrainzSearcher type, in theory.
type Searcher interface {
	Search() (Album, error)
}

// Struct for searching musicbrainz for an album.
type MusicBrainzSearcher struct {
	Criteria SearchCriteria
}

// Represents a musical artist.
type Artist struct {
	Id   string
	Name string
}

// Represents an album.
type Album struct {
	Id     string
	Title  string
	Artist string
}

// Searches MusicBrainz for an artist/album combination.
func (mbSearcher *MusicBrainzSearcher) Search() (Album, error) {
	/*
		TODO: create sub-routine to handle executing query and extracting returnvalue
	*/
	album := Album{}
	val := url.Values{}
	val.Add("query", mbSearcher.Criteria.ArtistName)
	val.Add("fmt", "json")
	query := "http://musicbrainz.org/ws/2/artist?" + val.Encode()
	resp, err := http.Get(query)
	if err != nil {
		return album, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return album, err
	}

	var result struct {
		Artists [1]Artist
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return album, err
	}

	var artistId = result.Artists[0].Id
	val.Del("query")
	val.Del("fmt")
	//musicbrainz relies on the order of the parameters
	//if arid is first, then the results are all hosed up
	//furthermore, the url.Value type encodes sorted by key
	//so this is really hacky
	val.Add("release", mbSearcher.Criteria.AlbumName)

	albumQuery := val.Encode()

	val.Del("release")
	val.Add("arid", artistId)

	albumQuery += "&" + val.Encode()
	//for whatever reason, musicbrainz uses : instead of =
	//for query string parameter pairs
	albumQuery = strings.Replace(albumQuery, "=", ":", -1)

	query = "http://musicbrainz.org/ws/2/release/?query=" + albumQuery
	query += "&fmt=json"

	var resp2 *http.Response
	resp2, err = http.Get(query)
	if err != nil {
		return album, err
	}
	defer resp2.Body.Close()
	body2, err2 := ioutil.ReadAll(resp2.Body)
	if err2 != nil {
		return album, err
	}

	/*
		TODO: for whatever reason, the method I used to Unmarshal
		the artist is not working for album:
			If i try it, the Albums[0].Title is blank
	*/
	var topLevel interface{}
	err = json.Unmarshal(body2, &topLevel)
	m := topLevel.(map[string]interface{})
	albumHackResult := (m["releases"].([]interface{})[0]).(map[string]interface{})
	album.Title = albumHackResult["title"].(string)
	album.Id = albumHackResult["id"].(string)
	album.Artist = result.Artists[0].Name
	return album, nil
}

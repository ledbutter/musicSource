package musicSource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

type Album struct {
	Title string
}

func (mbSearcher *MusicBrainzSearcher) Search() string {
	/*
		TODO: create sub-routine to handle executing query and extracting returnvalue
	*/
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
		return "Uh Oh, Album Get()"
	}
	defer resp2.Body.Close()
	body2, err2 := ioutil.ReadAll(resp2.Body)
	if err2 != nil {
		return "Uh Oh, Album ReadAllBytes()"
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
	foo := albumHackResult["title"].(string)
	return foo
	/*var result2 struct {
		Albums [1]Album
	}

	err2 = json.Unmarshal(body2, &result2)
	if err2 != nil {
		return "Uh oh, Album Unmarshal"
	}

	fmt.Printf("%#v\n", result2.Albums[0])
	return result2.Albums[0].Title*/
}

package setlistfm

import (
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/setlist"
	"festwrap/internal/str"
)

type SetlistFMRepository struct {
	host                  string
	apiKey                string
	deserializer          serialization.Deserializer[setlistFMResponse]
	httpSender            httpsender.HTTPRequestSender
	maxPages              int
	nextPageSleepMs       int
	artistMaxEditDistance int
}

func NewSetlistFMSetlistRepository(apiKey string, httpSender httpsender.HTTPRequestSender) *SetlistFMRepository {
	deserializer := serialization.NewJsonDeserializer[setlistFMResponse]()
	return &SetlistFMRepository{
		host:                  "api.setlist.fm",
		apiKey:                apiKey,
		deserializer:          &deserializer,
		httpSender:            httpSender,
		maxPages:              1,
		nextPageSleepMs:       0,
		artistMaxEditDistance: 5,
	}
}

func (r *SetlistFMRepository) GetSetlist(artist string, minSongs int) (setlist.Setlist, error) {

	page := 1
	var resultSetlist setlist.Setlist
	var err error
	setlistFound := false

	for page <= r.maxPages {
		resultSetlist, err = r.getFirstSetlistFromPage(artist, page, minSongs)
		if err == nil {
			setlistFound = true
			break
		} else {
			page += 1
			// Sleep to avoid hitting Setlistfm rate limit
			time.Sleep(time.Duration(r.nextPageSleepMs) * time.Millisecond)
		}
	}

	if !setlistFound {
		return setlist.Setlist{}, fmt.Errorf("could not find setlist for artist %s", artist)
	}
	// Make sure we maintain the input artist name, as the API might return close matches
	resultSetlist.SetArtist(artist)
	return resultSetlist, nil
}

func (r *SetlistFMRepository) getFirstSetlistFromPage(artist string, page int, minSongs int) (setlist.Setlist, error) {
	httpOptions := r.createSetlistHttpOptions(artist, page)
	responseBody, err := r.httpSender.Send(httpOptions)
	if err != nil {
		return setlist.Setlist{}, err
	}

	var response setlistFMResponse
	err = r.deserializer.Deserialize(*responseBody, &response)
	if err != nil {
		return setlist.Setlist{}, err
	}
	result, err := r.findValidSetlist(artist, response, minSongs)
	if err != nil {
		return setlist.Setlist{}, err
	}

	return result, nil
}

func (r *SetlistFMRepository) findValidSetlist(
	artist string,
	response setlistFMResponse,
	minSongs int,
) (setlist.Setlist, error) {
	validSetlists := response.getSetlistsWithMinSongs(minSongs)
	if len(validSetlists) == 0 {
		return setlist.Setlist{}, fmt.Errorf(
			"could not find setlist for artist %s with minimum songs %d",
			artist,
			minSongs,
		)
	}

	foundSetlist, err := r.getMatchingArtistSetlist(artist, validSetlists)
	if err == nil {
		return foundSetlist, nil
	}

	foundSetlist, err = r.findClosestArtistMatch(artist, validSetlists)
	if err != nil {
		return setlist.Setlist{}, err
	}
	return foundSetlist, nil
}

func (r *SetlistFMRepository) getMatchingArtistSetlist(
	artist string,
	setlists []setlist.Setlist,
) (setlist.Setlist, error) {
	for _, setlist := range setlists {
		if setlist.GetArtist() == artist {
			return setlist, nil
		}
	}
	return setlist.Setlist{}, fmt.Errorf("could not find setlist for artist %s", artist)
}

func (r *SetlistFMRepository) findClosestArtistMatch(artist string, setlists []setlist.Setlist) (setlist.Setlist, error) {
	var distances []int
	for _, setlist := range setlists {
		distance := str.LevenshteinDistance{}.Compute(
			strings.ToLower(artist),
			strings.ToLower(setlist.GetArtist()),
		)
		if distance > r.artistMaxEditDistance {
			distance = math.MaxInt
		}
		distances = append(distances, distance)
	}

	minDistance := math.MaxInt
	minDistanceIndex := -1
	for i, distance := range distances {
		if distance < minDistance {
			minDistance = distance
			minDistanceIndex = i
		}
	}

	if minDistance == math.MaxInt {
		error := fmt.Errorf(
			"could not find setlist for artist %s and max edit distance of %d",
			artist,
			r.artistMaxEditDistance,
		)
		return setlist.Setlist{}, error
	} else {
		return setlists[minDistanceIndex], nil
	}
}

func (r *SetlistFMRepository) createSetlistHttpOptions(artist string, page int) httpsender.HTTPRequestOptions {
	url := r.getSetlistFullUrl(artist, page)
	httpOptions := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	httpOptions.SetHeaders(
		map[string]string{
			"x-api-key": r.apiKey,
			"Accept":    "application/json",
		},
	)
	return httpOptions
}

func (r *SetlistFMRepository) getSetlistFullUrl(artist string, page int) string {
	queryParams := url.Values{}
	queryParams.Set("artistName", artist)
	queryParams.Set("p", fmt.Sprint(page))
	setlistPath := "rest/1.0/search/setlists"
	return fmt.Sprintf("https://%s/%s?%s", r.host, setlistPath, queryParams.Encode())
}

func (r *SetlistFMRepository) SetMaxPages(maxPages int) {
	r.maxPages = maxPages
}

func (r *SetlistFMRepository) SetNextPageSleep(sleepMs int) {
	r.nextPageSleepMs = sleepMs
}

func (r *SetlistFMRepository) SetArtistMaxEditDistance(distance int) {
	r.artistMaxEditDistance = distance
}

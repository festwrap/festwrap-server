package setlistfm

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"festwrap/internal/setlist"
	"festwrap/internal/setlist/errors"
)

type SetlistFMSetlistRepositoryConfig struct {
	Client *http.Client
	Host   string
	ApiKey string
}

type SetlistFMRepository struct {
	config *SetlistFMSetlistRepositoryConfig
	parser setlist.SetlistParser
}

// TODO: abstract client
// TODO: check more pages until setlist found
// raise error if no 200
// hierarchy: HTTP multipage

func (s *SetlistFMRepository) GetSetlist(artist string) (*setlist.Setlist, error) {

	request, err := s.createSetlistRequest(artist)
	if err != nil {
		errorMsg := fmt.Sprintf("Cannot create request for %s: %s", s.config.Host, err.Error())
		return nil, errors.NewCannotRetrieveSetlistError(errorMsg)
	}

	response, err := s.config.Client.Do(request)
	if err != nil {
		errorMsg := fmt.Sprintf(
			"Cannot retrieve response for %s: %s", s.config.Host, err.Error(),
		)
		return nil, errors.NewCannotRetrieveSetlistError(errorMsg)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		errorMsg := fmt.Sprintf("Error reading response %s: ", err.Error())
		return nil, errors.NewCannotRetrieveSetlistError(errorMsg)
	}

	setlist, err := s.parser.Parse(body)
	if err != nil {
		return nil, err
	}
	if setlist == nil {
		// TODO: if no valid setlist found, we should check for the next page
		// TODO: probable a good idea to move the min songs filter to repository and keep parser simpler
		errorMsg := fmt.Sprintf("Could not find setlist for artist %s", artist)
		errors.NewCannotRetrieveSetlistError(errorMsg)
	}

	return setlist, nil
}

func (s *SetlistFMRepository) createSetlistRequest(artist string) (*http.Request, error) {
	request, err := http.NewRequest("GET", s.getSetlistFullUrl(artist), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("x-api-key", s.config.ApiKey)
	request.Header.Add("Accept", "application/json")
	return request, nil
}

func (s *SetlistFMRepository) getSetlistFullUrl(artist string) string {
	queryParams := url.Values{}
	queryParams.Set("artistName", artist)
	setlistPath := "rest/1.0/search/setlists"
	return fmt.Sprintf("https://%s/%s?%s", s.config.Host, setlistPath, queryParams.Encode())
}

func NewSetlistFMSetlistRepository(config *SetlistFMSetlistRepositoryConfig, parser setlist.SetlistParser) *SetlistFMRepository {
	return &SetlistFMRepository{config: config, parser: parser}
}

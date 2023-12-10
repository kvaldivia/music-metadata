package services

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/kvaldivia/music-metadata/internal/models"
)

type service struct {
	url        string
	token      string
	httpClient *http.Client
}

type spotifyArtist struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Popularity int    `json:"popularity"`
}

type spotifyExternalIDs struct {
	ISRC string `json:"isrc"`
	EAN  string `json:"ean"`
	UPC  string `json:"upc"`
}

type spotifyTrack struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Artists     []spotifyArtist    `json:"artists"`
	ExternalIDs spotifyExternalIDs `json:"external_ids"`
}

func NewSpotifyService(url string, token string) Service {
	return &service{url: url, token: token, httpClient: &http.Client{}}
}

func (s *service) Get(ctx context.Context, id string) (*models.Track, error) {
	var track models.Track
	var err error

	url := s.url + "/tracks/" + id
	println("Authorization", "Basic "+s.token)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Join(errors.New("unable to create request for track id: "+id), err)
	}

	req.Header.Add("Authorization", "Basic "+s.token)
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, errors.Join(errors.New("unable to fetch data from spotify api"), err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Join(errors.New("unable to read response body"), err)
	}

	err = json.Unmarshal(body, &track)

	if err != nil {
		return nil, errors.Join(errors.New("unable to unmarshal response body"), err)
	}

	return &track, nil
}

func (s *service) Search(ctx context.Context, q string) ([]*models.Track, error) {
	return nil, nil
}

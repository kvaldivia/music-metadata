package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/kvaldivia/music-metadata/internal/models"
	"github.com/kvaldivia/music-metadata/internal/tools/logger"
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

type spotifyImage struct {
	Url    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type spotifyAlbum struct {
	Images []spotifyImage `json:"images"`
}

type spotifyTrack struct {
	ID          string             `json:"id"`
	Album       spotifyAlbum       `json:"album"`
	Name        string             `json:"name"`
	Artists     []spotifyArtist    `json:"artists"`
	ExternalIDs spotifyExternalIDs `json:"external_ids"`
}

type spotifyAuthResponseBody struct {
	AccessToken string `json:"access_token"`
}

var l = logger.GetLogger()

func NewSpotifyService(serviceUrl string, clientId, clientSecret string) (Service, error) {
	client := http.Client{}
	var err error

	authUrl := "https://accounts.spotify.com/api/token"
	reqParams := url.Values{}
	reqParams.Add("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, authUrl, strings.NewReader(reqParams.Encode()))
	if err != nil {
		return nil, errors.Join(errors.New("could not prepare auth request for spotify api"), err)
	}
	authString := clientId + ":" + clientSecret
	encodedAuthStr := base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Add("Authorization", "Basic "+encodedAuthStr)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("could not complete authentication against spotify api"), err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var err error
		resBody, e := io.ReadAll(res.Body)
		if e != nil {
			err = errors.Join(errors.New("could not read response body"), e)
		} else {
			l.Info("Spotify authentication failed", "response", string(resBody))
		}
		reqBody, e := ioutil.ReadAll(req.Body)
		if e != nil {
			err = errors.Join(errors.New("could not reead request body"), err, e)
		} else {
			l.Info("unsuccessful request info", "body", string(reqBody), "auth string", authString, "request params string", reqParams.Encode())
		}

		return nil, errors.Join(errors.New("spotify auth failed with status code: "+fmt.Sprintf("%d", res.StatusCode)+""), err)
	}

	bodyStr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Join(errors.New("could not parse spotify auth response"), err)
	}

	var body spotifyAuthResponseBody
	err = json.Unmarshal(bodyStr, &body)

	if err != nil {
		l.Info("unmarshalling error", "bodyStr", string(bodyStr))
		return nil, errors.Join(errors.New("could not unmarshal auth response"), err)
	}

	return &service{url: serviceUrl, token: body.AccessToken, httpClient: &client}, nil
}

func (s *service) Get(ctx context.Context, id string) (*models.Track, *models.Artist, error) {
	var track models.Track
	var spotifyTrack spotifyTrack
	var err error

	url := s.url + "/tracks/" + id

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, errors.Join(errors.New("unable to create request for track id: "+id), err)
	}

	req.Header.Add("Authorization", "Bearer "+s.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)

	if err != nil {
		return nil, nil, errors.Join(errors.New("unable to fetch data from spotify api"), err)
	}

	defer resp.Body.Close()

	bodyStr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		l.Error("could not read response", "error", err)
		return nil, nil, errors.Join(errors.New("unable to read response body"), err)
	}

	err = json.Unmarshal(bodyStr, &spotifyTrack)

	if err != nil {
		l.Error("could not unmarshal track from response", "error", err)
		return nil, nil, errors.Join(errors.New("unable to unmarshal response body"), err)
	}

	track = models.Track{
		SpotifyID: spotifyTrack.ID,
		ISRC:      spotifyTrack.ExternalIDs.ISRC,
		ImageURI:  spotifyTrack.Album.Images[0].Url,
		Title:     spotifyTrack.Name,
	}

	artist := models.Artist{
		Tracks:    []models.Track{},
		Name:      spotifyTrack.Artists[0].Name,
		SpotifyID: spotifyTrack.ID,
	}

	return &track, &artist, nil
}

func (s *service) Search(ctx context.Context, q string) ([]*models.Track, error) {
	return nil, nil
}

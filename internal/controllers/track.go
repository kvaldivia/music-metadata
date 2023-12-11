package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kvaldivia/music-metadata/internal/models"
	"github.com/kvaldivia/music-metadata/internal/services"
	"github.com/kvaldivia/music-metadata/internal/store"
	"github.com/kvaldivia/music-metadata/internal/tools/logger"
)

var l = logger.GetLogger()

type tracksController struct {
	Store          store.Track
	SpotifyService services.Service
}

type getTrackByISRCParams struct {
	ISRC string `uri:"isrc" binding:"required"`
}

func NewTracksController(st *store.Track, ss *services.Service) tracksController {
	return tracksController{Store: *st, SpotifyService: *ss}
}

func (t *tracksController) GetTrackByISRC(c *gin.Context) {
	var params getTrackByISRCParams
	var err error

	c.Header("Content-Type", "application/json")

	if err = c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		c.Abort()
		return
	}

	if params.ISRC == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid isrc: " + params.ISRC + ". Should not be empty."})
		return
	}

	track, err := t.Store.Find(c.Request.Context(), params.ISRC)

	if err != nil {
		err = errors.Join(errors.New("Could not query for the ISRC: "+params.ISRC), err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if track != nil {
		c.JSON(http.StatusOK, track)
		return
	}

	// Out of scope cases fall here
	c.JSON(http.StatusInternalServerError, errors.New("could not process properly"))
	// Abort so other handlers don't waste time on this request
	c.Abort()
}

type addNewTrackRequestBody struct {
	ISRC      string `uri:"isrc"`
	SpotifyID string `uri:"spotifyId"`
}

func (t *tracksController) AddNewTrack(c *gin.Context) {
	var track *models.Track
	var artist *models.Artist
	var err error
	var requestBody addNewTrackRequestBody
	c.Header("Content-Type", "application/json")

	if err = c.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
		err = errors.Join(errors.New("could not process request body"), err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if requestBody.SpotifyID == "" {
		err = errors.New("invalid ID, should not be blank")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track, _ = t.Store.Find(c, requestBody.SpotifyID)

	if track.SpotifyID == requestBody.SpotifyID {
		c.JSON(http.StatusConflict, gin.H{"error": "A track with the same ID already exists."})
		return
	}

	defer c.Done()
	track, artist, err = t.SpotifyService.Get(c, requestBody.SpotifyID)
	if err != nil {
		err = errors.Join(errors.New("could not get data from API"), err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if track.SpotifyID == requestBody.SpotifyID {
		var elm interface{}
		if artist != nil {
			artist.Tracks = append(artist.Tracks, *track)
			elm = artist
		} else {
			elm = track
		}
		if err = t.Store.Create(c, elm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, &track)
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not process correctly"})
}

type allByArtistUriParams struct {
	ArtistID string `uri:"artistId" binding:"required"`
}

func (t *tracksController) AllByArtist(c *gin.Context) {
	var err error
	var params allByArtistUriParams

	err = c.ShouldBindUri(&params)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/json")

	if params.ArtistID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist id"})
		return
	}

	tracks, err := t.Store.AllByArtist(c, params.ArtistID)

	if err != nil {
		err = errors.Join(errors.New(""), err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracks)
}

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

type TracksController struct {
	Store          store.Track
	SpotifyService services.Service
}

type getTrackByISRCParams struct {
	ISRC string `uri:"isrc" binding:"required"`
}

func NewTracksController(st *store.Track, ss *services.Service) TracksController {
	return TracksController{Store: *st, SpotifyService: *ss}
}

// GetTrackByISRC godoc
//
//	@Summary    get matching track details, use isrc to match track
//	@Description	returns the details for an existing track that matches the provided isrc.
//	@Tags			v1
//	@Produce		json
//	@Param			isrc path string true	"used to look up a track record"
//	@Success		200		{object}	controllers.JSONResult{data=models.Track} "track details"
//	@Failure		400		{object}	controllers.JSONResult "bad request"
//	@Failure		404		{object}	controllers.JSONResult "not found"
//	@Failure		500		{object}	controllers.JSONResult "internal error"
//	@Router			/v1/tracks/{isrc} [get]
func (t *TracksController) GetTrackByISRC(c *gin.Context) {
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

// AddNewTrack godoc
//
// @Summary    request to add a new spotify track
//
// @Description	receives a spotify track id and fetches the track's info from spotify api
// @Tags  v1
// @Accept  json
// @Produce		json
// @Param			isrc path string true	"used to look up car record"
// @Success		201 {object}	controllers.JSONResult{data=models.Track} "track details"
// @Failure		400 {object}	controllers.JSONResult "bad request"
// @Failure		409 {object}	controllers.JSONResult "conflict"
// @Failure		500 {object}	controllers.JSONResult "internal error"
// @Router			/v1/tracks [post]
func (t *TracksController) AddNewTrack(c *gin.Context) {
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
		c.JSON(http.StatusCreated, &track)
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not process correctly"})
}

type allByArtistUriParams struct {
	ArtistID string `uri:"artistId" binding:"required"`
}

// AllByArtist godoc
//
// @Summary		list of existing tracks that belong to an artist
// @Description	provides a list of the existing track records that belong to a single existing artist
// @Tags			v1
// @Accept			json
// @Produce	  json
// @Param artistId path string true "used to look for an artist record"
// @Success		200		{object}	controllers.JSONResult{data=[]models.Track}	"tracks list"
// @Failure		400		{object}	controllers.JSONResult "bad request"
// @Failure		404		{object}	controllers.JSONResult "not found"
// @Failure		500		{object}	controllers.JSONResult "internal error"
// @Router			/v1/artists/{artistId}/tracks [get]
func (t *TracksController) AllByArtist(c *gin.Context) {
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

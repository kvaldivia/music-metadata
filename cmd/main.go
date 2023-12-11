package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kvaldivia/music-metadata/internal/controllers"
	"github.com/kvaldivia/music-metadata/internal/services"
	"github.com/kvaldivia/music-metadata/internal/store"
	"github.com/kvaldivia/music-metadata/internal/tools/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var l = logger.GetLogger()

func setupDb() (*gorm.DB, error) {
	dsn := "host=db user=music-metadata password=v9qsJRuw6e dbname=music-metadata sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
	} else {
		log.Println(db.Config)
	}
	return db, err
}

// @title MusicMetadata API
// @version 1.0
// @description This is a music metadata api that enriches track details using spotify API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email k_valdivia@gmx.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host www.music-metadata.io
// @BasePath /v1
func main() {
	db, err := setupDb()
	if err != nil {
		l.Error("could not start db", "error", err.Error())
		return
	}
	r := gin.Default()
	trackStore := store.NewTrackStore(db)
	spotifyService, err := services.NewSpotifyService("https://api.spotify.com/v1", os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	if err != nil {
		l.Error("could not start spotify service", "error", err.Error())
		return
	}
	trackController := controllers.NewTracksController(&trackStore, &spotifyService)

	// Get track by ISRC
	r.GET("/v1/tracks/:isrc", func(c *gin.Context) {
		trackController.GetTrackByISRC(c)
	})

	// Get user value
	r.GET("/v1/artists/:artistId/tracks", func(c *gin.Context) {
		trackController.AllByArtist(c)
	})

	// TODO(kvaldivia): implement auth
	//authorized := r.Group("/v1/tracks", gin.BasicAuth(gin.Accounts{
	//	"foo":  "bar", // user:foo password:bar
	//	"manu": "123", // user:manu password:123
	//}))

	//authorized.POST("", func(c *gin.Context) {
	//	var _ string = c.MustGet(gin.AuthUserKey).(string)
	//	controllers.AddNewTrack(c, trackStore)
	//})
	r.POST("/v1/tracks", func(c *gin.Context) {
		trackController.AddNewTrack(c)
	})

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

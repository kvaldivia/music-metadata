package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kvaldivia/music-metadata/internal/controllers"
	"github.com/kvaldivia/music-metadata/internal/services"
	"github.com/kvaldivia/music-metadata/internal/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db = make(map[string]string)

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

func main() {
	db, err := setupDb()
	if err != nil {
		log.Println(err)
		return
	}
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	trackStore := store.NewTrackStore(db)
	spotifyService := services.NewSpotifyService("https://api.spotify.com/v1", "")
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

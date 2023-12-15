
#define populate_db_with_tracks
#SPOTIFY_TRACK_IDS="7ouMYWpwJ422jRcDASZB7P" "4VqPOruhp5EdPBeR92t6lQ" "2takcwOaAZWiXQijPHIx7B"
#for track_id in ${echo $SPOTIFY_TRACK_IDS}; do
#	echo "Reaching here"
#	curl -d '{"spotifyId":"$track_id"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/tracks
#done
#endef

populate-db:
	curl -d '{"spotifyId":"11dFghVXANMlKmJXsNCbNl"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/tracks
	curl -d '{"spotifyId":"7ouMYWpwJ422jRcDASZB7P"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/tracks
	curl -d '{"spotifyId":"4VqPOruhp5EdPBeR92t6lQ"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/tracks
	curl -d '{"spotifyId":"2takcwOaAZWiXQijPHIx7B"}' -H "Content-Type: application/json" -X POST http://localhost:8080/api/v1/tracks

fetch-songs-by-muse:
	curl http://localhost:8080/api/v1/artists/12Chz98pHFMPJEknJQMWvI/tracks

fetch-songs-by-crj:
	curl http://localhost:8080/api/v1/artists/6sFIWsNpZYqfjUpaCgueju/tracks

fetch-track-info:
	curl http://localhost:8080/api/v1/tracks/11dFghVXANMlKmJXsNCbNl

fetch-unexistent-track-info:
	curl http://localhost:8080/api/v1/tracks/11dFghVXANMlKmJXsASDFASDF

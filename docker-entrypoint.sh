echo "Downloading go module deps"
go mod download
echo "Running migrations"
atlas migrate apply --env $ENV
if [ $ENV = "dev" ];
then
  echo "Running on dev mode"
  echo "spotify client id: " $SPOTIFY_CLIENT_ID
  gow run cmd/main.go
else
  echo "Running on prod mode"
  echo "Building go app"
  go build -o bin/main cmd/main.go
  echo "Running go module"
  ./bin/main
fi

echo "Downloading go module deps"
go mod download
echo "Building go module"
go build -o .
echo "Running migrations"
atlas migrate apply --env $ENV
echo "Running go module"
./music-metadata

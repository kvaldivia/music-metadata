echo "Downloading go module deps"
go mod download
echo "Building go app"
go build -o bin/main cmd/main.go
echo "Running migrations"
atlas migrate apply --env $ENV
echo "Running go module"
./bin/main

echo "Downloading go module deps"
go mod download
echo "Building go module"
go build -o .
echo "Running go module"
ls -la .
./music-metadata

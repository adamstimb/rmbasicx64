# The script builds the 64-bit MacOS executable

echo "Building RM BASICx64 for 64-bit MacOS: build/rmbasicx64"
export GOOS=darwin
export GOARCH=amd64
go build -o ../build/rmbasicx64 ../cmd/rmbasicx64/main.go
cp ../examples/*.BAS ../build/workspace/
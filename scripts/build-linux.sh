# The script builds the 64-bit Linux executable

echo "Building RM BASICx64 for 64-bit Linux: build/rmbasicx64"
export GOOS=linux
export GOARCH=amd64
go build -o ../build/rmbasicx64 ../cmd/rmbasicx64/main.go
cp ../examples/*.BAS ../build/workspace/
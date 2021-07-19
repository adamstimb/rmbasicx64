# The script builds the 64-bit Windows executable

echo "Building RM BASICx64 for 64-bit Windows: build/rmbasicx64.exe"
export GOOS=windows
export GOARCH=amd64
go build -o ../build/rmbasicx64.exe ../cmd/rmbasicx64/main.go
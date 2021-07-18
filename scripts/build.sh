# The script builds the 64-bit Windows executable

echo "Building RM BASICx64 for 64-bit Windows: build/windows/rmbasic.exe"
export GOOS=windows
export GOARCH=amd64
go build -o ../build/rmbasic.exe ../cmd/rmbasicx64/main.go
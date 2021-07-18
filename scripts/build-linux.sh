# The script builds the 64-bit Windows executable

echo "Building RM BASICx64 for 64-bit Linux: build/windows/rmbasic"
export GOOS=linux
export GOARCH=amd64
go build -o ../build/rmbasic ../cmd/rmbasicx64/main.go
# The script builds the executable for whichever platform the script is running

echo "Building RM BASICx64 for this platform: build/rmbasicx64"
go build -o ../build/rmbasicx64 ../cmd/rmbasicx64/main.go
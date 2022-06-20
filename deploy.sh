mkdir exec

echo "Deleting old executables"
rm ./exec/*

go mod tidy

echo "Compiling for MacOs AMD64"
GOOS=darwin GOARCH=amd64 go build  -o S3-JSON-VISUALIZER-mamd64 .
mv S3-JSON-VISUALIZER-mamd64 ./exec/S3-JSON-VISUALIZER-mamd64

echo "Compiling for LINUX"
GOOS=linux GOARCH=amd64 go build  -o S3-JSON-VISUALIZER-linux .
mv S3-JSON-VISUALIZER-linux ./exec/S3-JSON-VISUALIZER-linux

echo "Compiling for MacOs ARM1"
GOOS=darwin GOARCH=arm64 go build -o S3-JSON-VISUALIZER-marm1 .
mv S3-JSON-VISUALIZER-marm1  ./exec/S3-JSON-VISUALIZER-marm1

echo "Compiling for Windows"
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++  go build -ldflags="-H windowsgui" -o S3-JSON-VISUALIZER.exe
mv S3-JSON-VISUALIZER.exe ./exec/S3-JSON-VISUALIZER.exe

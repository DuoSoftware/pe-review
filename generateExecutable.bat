echo "Generating Executable"
echo "OS Code": %1
set GOOS=%1
echo "System Architecture": %2
set GOARCH=%2
cd %3
echo "Building Executable"
go build -v %4
echo "Executable generation complete!"
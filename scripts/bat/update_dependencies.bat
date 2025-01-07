@echo off

IF NOT EXIST "go.mod" (
    echo "go.mod not found! Please run this script from the root of your Go project."
    exit /b 1
)

echo Updating dependencies to their latest versions...
go get -u ./...

go mod tidy

echo Dependencies have been updated to the latest versions.
pause
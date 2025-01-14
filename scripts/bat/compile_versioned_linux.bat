@echo off
for /f "tokens=*" %%i in ('git describe --tags --always --abbrev=7') do set VERSION=v%%i

cd ..\..\src

set GOOS=linux
set GOARCH=amd64
go build -ldflags "-X version.Version=%VERSION%" -o bin\MonCaveau-%VERSION%

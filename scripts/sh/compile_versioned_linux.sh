#!/bin/bash

VERSION="v$(git describe --tags --always --abbrev=7)"

cd ../../src/

env GOOS=linux GOARCH=amd64 go build -ldflags "-X version.Version=$VERSION" -o bin/MonCaveau-$VERSION
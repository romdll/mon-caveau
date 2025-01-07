#!/bin/bash

if [ ! -f "go.mod" ]; then
  echo "go.mod not found! Please run this script from the root of your Go project."
  exit 1
fi

echo "Updating dependencies to their latest versions..."
go get -u ./...

go mod tidy

echo "Dependencies have been updated to the latest versions."
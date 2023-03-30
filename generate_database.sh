#!/usr/bin/env bash

echo "building the database..."
pushd ./assets/database/ || exit
go build
./database
popd || exit

echo "copying the database..."
cp ./assets/database/database_keywords.json .

echo "running go generate..."
go generate

echo "running all tests to see how it went..."
go test -v ./...

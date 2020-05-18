#!/usr/bin/env bash

echo "building the database..."
pushd ./assets/database/ || exit
go run build_database.go
popd || exit

echo "copying the database..."
cp ./assets/database/database_keywords.json .

echo "running go generate..."
go generate

echo "running all tests to see how it went..."
go test -v -run TestFullDatabase ./processor

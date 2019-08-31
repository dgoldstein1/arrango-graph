#!/bin/bash

# start arango in background
docker-compose up -d

while true; do

# for less verbose outout
export GIN_MODE=test

inotifywait -e modify,create,delete -r ./ && \
	clear
	go fmt ./... \
		&& go build -o build/destrib-graph \
		&& go test ./...
done

#!/bin/bash

# start arango in background
docker-compose up -d

export GRAPH_DB_NAME="test" # name of database in arango
export GRAPH_DB_COLLECTION_NAME="testingcollection" # collection name within arango db name
export GRAPH_DB_ARANGO_ENDPOINTS="http://localhost:8529" #list of arango db endpoints

while true; do

# for less verbose outout
export GIN_MODE=test

inotifywait -e modify,create,delete -r ./ && \
	clear
	go fmt ./... \
		&& go build -o build/destrib-graph \
		&& go test ./...
done

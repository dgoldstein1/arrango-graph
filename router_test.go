package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	// setup environment
	dbName := "graph-testing-router"
	os.Setenv("GRAPH_DB_NAME", dbName)
	os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-router-wikipedia")
	os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
	os.Setenv("GRAPH_DB_NAME", "wikipedia-graph")
	// run test
	router, err := SetupRouter("./api/*")
	assert.NotNil(t, router)
	assert.NotNil(t, err)
}

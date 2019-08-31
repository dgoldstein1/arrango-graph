package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"testing"
)

// export GRAPH_DB_NAME="arango_graphs" # name of database in arango
// export GRAPH_DB_COLLECTION_NAME="wikipedia" # collection name within arango db name
// export GRAPH_DB_ARANGO_ENDPOINTS="http://localhost:9520" #list of arango db endpoints
// export GRAPH_DB_NAME="wikipedia-graph" # name of graph within collection
func TestConnectToDB(t *testing.T) {
	// mock out log.Fatalf
	origLogFatalf := logFatalf
	defer func() { logFatalf = origLogFatalf }()
	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	// positive test
	os.Setenv("GRAPH_DB_NAME", "graph-testing")
	os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
	os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
	os.Setenv("GRAPH_DB_NAME", "wikipedia-graph")
	g := ConnectToDB()
	assert.NotNil(t, g)
	require.Equal(t, []string{}, errors)
	// try creating same graph again, should not fail
	os.Setenv("GRAPH_DB_NAME", "wikipedia-graph-"+string(rand.Int()))
	g = ConnectToDB()
	assert.NotNil(t, g)
	assert.Equal(t, []string{}, errors)
	// remove graph

}

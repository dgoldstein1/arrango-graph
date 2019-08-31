package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// export GRAPH_DB_NAME="arango_graphs" # name of database in arango
// export GRAPH_DB_COLLECTION_NAME="wikipedia" # collection name within arango db name
// export GRAPH_DB_ARANGO_ENDPOINTS="http://localhost:9520" #list of arango db endpoints
func TestConnectToDB(t *testing.T) {
	// positive test
	os.Setenv("GRAPH_DB_NAME", "graph-testing")
	os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
	os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
	err, g := ConnectToDB()
	assert.Nil(t, err)
	assert.NotNil(t, g)
	// try creating same graph again, should not fail
	err, g = ConnectToDB()
	assert.Nil(t, err)
	assert.NotNil(t, g)

}

package main

import (
	"os"
	"testing"
)

// export GRAPH_DB_NAME="arango_graphs" # name of database in arango
// export GRAPH_DB_COLLECTION_NAME="wikipedia" # collection name within arango db name
// export GRAPH_DB_arango_ENDPOINTS="http://localhost:9520" #list of arango db endpoints
func TestConnectToDB(t *testing.T) {
	// positive test
	os.Setenv("GRAPH_DB_NAME", "graph-testing")
	os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
	os.Setenv("GRAPH_DB_arango_ENDPOINTS", "http://localhost:9520")
}

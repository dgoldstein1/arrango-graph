package main

import (
	"os"
	"testing"
)

// export GRAPH_DB_NAME="arrango_graphs" # name of database in arrango
// export GRAPH_DB_COLLECTION_NAME="wikipedia" # collection name within arrango db name
// export GRAPH_DB_ARRANGO_ENDPOINTS="http://localhost:9520" #list of arrango db endpoints
func TestConnectToDB(t *testing.T) {
	// positive test
	os.Setenv("GRAPH_DB_NAME", "graph-testing")
	os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
	os.Setenv("GRAPH_DB_ARRANGO_ENDPOINTS", "http://localhost:9520")
}

package main

import (
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	graphsToDelete := []driver.Graph{}
	t.Run("connects to db that doesnt already exist", func(t *testing.T) {
		dbName := "graph-testing-db"
		os.Setenv("GRAPH_DB_NAME", dbName)
		os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-testing-wikipedia")
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
		os.Setenv("GRAPH_DB_NAME", "wikipedia-graph")
		g, nodes, edges := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		graphsToDelete = append(graphsToDelete, g)
	})
	t.Run("connnects to DB that already exists", func(t *testing.T) {
		g, nodes, edges := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
	})
	t.Run("connects to same DB with new graph name", func(t *testing.T) {
		dbName2 := "graph-testing-2"
		os.Setenv("GRAPH_DB_NAME", dbName2)
		g, nodes, edges := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		graphsToDelete = append(graphsToDelete, g)
	})
	t.Run("bad url endpoints", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8000")
		g, nodes, edges := ConnectToDB()
		assert.Nil(t, nodes)
		assert.Nil(t, edges)
		assert.Nil(t, g)
		assert.Equal(t, []string{"Could not establish connection to DB [Could not check if databse exists create database at [http://localhost:8000]: Get http://localhost:8000/_db/graph-testing-2/_api/database/current: dial tcp 127.0.0.1:8000: connect: connection refused]"}, errors)
		errors = []string{}
	})
	t.Run("bad db name", func(t *testing.T) {
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
		os.Setenv("GRAPH_DB_NAME", "sldjf093ur2n093r2039d[2e9ufsdf - -CC]")
		g, nodes, edges := ConnectToDB()
		assert.Nil(t, nodes)
		assert.Nil(t, edges)
		assert.Nil(t, g)
		assert.Equal(t, []string{"Could not establish connection to DB [Failed to initialize database: database name invalid]"}, errors)
		errors = []string{}
	})

	// remove created graphs
	for _, g := range graphsToDelete {
		require.Nil(t, g.Remove(nil))
	}
}

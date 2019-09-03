package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

// these need to be set before running tests:
//
// export GRAPH_DB_NAME="arango_graphs" # name of database in arango
// export GRAPH_DB_ARANGO_ENDPOINTS="http://localhost:8529" #list of arango db endpoints
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
	t.Run("connects to db that doesnt already exist and connects to graph that does exist", func(t *testing.T) {
		errors = []string{}
		g, nodes, edges, _ := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		// connect to graph we just created
		g, nodes, edges, _ = ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		require.Nil(t, g.Remove(nil))
	})
	t.Run("connects to same DB with new graph name", func(t *testing.T) {
		temp := os.Getenv("GRAPH_DB_NAME")
		defer os.Setenv("GRAPH_DB_NAME", temp)
		os.Setenv("GRAPH_DB_NAME", "testing-graph-2")
		errors = []string{}
		g, nodes, edges, _ := ConnectToDB()
		assert.NotNil(t, g)
		assert.NotNil(t, nodes)
		assert.NotNil(t, edges)
		require.Equal(t, []string{}, errors)
		require.Nil(t, g.Remove(nil))
	})
	t.Run("bad url endpoints", func(t *testing.T) {
		errors = []string{}
		temp := os.Getenv("GRAPH_DB_ARANGO_ENDPOINTS")
		defer os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", temp)
		os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8000")
		g, nodes, edges, _ := ConnectToDB()
		assert.Nil(t, nodes)
		assert.Nil(t, edges)
		assert.Nil(t, g)
		assert.Equal(t, 1, len(errors))
		errors = []string{}
	})
	t.Run("bad db name", func(t *testing.T) {
		errors = []string{}
		temp := os.Getenv("GRAPH_DB_NAME")
		defer os.Setenv("GRAPH_DB_NAME", temp)
		os.Setenv("GRAPH_DB_NAME", "sldjf093ur2n093r2039d[2e9ufsdf - -CC]")
		g, nodes, edges, _ := ConnectToDB()
		assert.Nil(t, nodes)
		assert.Nil(t, edges)
		assert.Nil(t, g)
		assert.Equal(t, []string{"Could not establish connection to DB [Failed to initialize database: database name invalid]"}, errors)
		errors = []string{}
	})
}

func TestAddEdgesDB(t *testing.T) {
	// mock out log.Fatalf
	logErrOriginal := logErr
	defer func() { logErr = logErrOriginal }()
	errors := []string{}
	logErr = func(format string, args ...interface{}) {
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}
	temp := os.Getenv("GRAPH_DB_NAME")
	defer os.Setenv("GRAPH_DB_NAME", temp)
	os.Setenv("GRAPH_DB_NAME", "testing-add-edges-to-graph")
	g, nodes, edges, _ := ConnectToDB()
	assert.NotNil(t, g)
	assert.NotNil(t, nodes)
	assert.NotNil(t, edges)
	require.Equal(t, []string{}, errors)
	defer require.Nil(t, g.Remove(nil))

	type Test struct {
		Before                   func()
		Name                     string
		Node                     string
		Neighbors                []string
		ExpectedError            error
		ExpectedNodesAddedLength int
		ExpectedErrorsLogged     []string
	}

	testTable := []Test{
		Test{
			Before: func() {
				g, nodes, edges, _ = ConnectToDB()
				nodes.RemoveDocuments(nil, []string{"new-node-2", "new-node-3"})
				edges.RemoveDocument(nil, "new-node-2TOnew-node-3")
			},
			Name:                     "addes all new edges",
			Node:                     "new-node-1",
			Neighbors:                []string{"new-node-2", "new-node-3"},
			ExpectedError:            nil,
			ExpectedNodesAddedLength: 2,
			ExpectedErrorsLogged:     []string{},
		},
		Test{
			Before: func() {
				g, nodes, edges, _ = ConnectToDB()
				nodes.RemoveDocuments(nil, []string{"new-node-2", "new-node-4"})
				edges.RemoveDocument(nil, "new-node-2TOnew-node-3")
			},
			Name:                     "only returns new nodes",
			Node:                     "new-node-1",
			Neighbors:                []string{"new-node-3", "new-node-4"},
			ExpectedError:            nil,
			ExpectedNodesAddedLength: 1,
			ExpectedErrorsLogged:     []string{},
		},
		Test{
			Before:                   func() {},
			Name:                     "bad node name",
			Node:                     "new-node-1-OSF#OK2O$ kCADK c/// adcaKf@",
			Neighbors:                []string{"new-node-3", "new-node-4"},
			ExpectedError:            nil,
			ExpectedNodesAddedLength: 0,
			ExpectedErrorsLogged:     []string{"Error adding nodes to graph [ illegal document key]: %!s(MISSING)", "Error adding edges to node [ document not found]: %!s(MISSING)", "Error adding edges to node [ document not found]: %!s(MISSING)"},
		},
	}

	s := Server{
		Nodes: nodes,
		Edges: edges,
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			errors = []string{}
			test.Before()
			e, nAdded := AddEdges(test.Node, test.Neighbors, s)
			assert.Equal(t, test.ExpectedError, e)
			assert.Equal(t, test.ExpectedErrorsLogged, errors)
			assert.Equal(t, test.ExpectedNodesAddedLength, len(nAdded))
			if test.ExpectedNodesAddedLength != len(nAdded) {
				fmt.Printf("nAdded : %v \n", nAdded)
			}
		})
	}

}

func TestGetEdgesFromDB(t *testing.T) {
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
	temp := os.Getenv("GRAPH_DB_NAME")
	defer os.Setenv("GRAPH_DB_NAME", temp)
	os.Setenv("GRAPH_DB_NAME", "testing-get-edges-from-graph")
	g, nodes, edges, db := ConnectToDB()
	assert.NotNil(t, g)
	assert.NotNil(t, nodes)
	assert.NotNil(t, edges)
	require.Equal(t, []string{}, errors)
	// defer require.Nil(t, g.Remove(nil))

	type Test struct {
		Before            func()
		Name              string
		Node              string
		ExpectedError     error
		ExpectedNeighbors []string
	}

	testTable := []Test{
		Test{
			Before: func() {
				nodes.CreateDocuments(nil, []Node{Node{"test1"}, Node{"test2"}, Node{"test3"}})
				meta, errors, err := edges.CreateDocuments(nil, []Edge{
					Edge{
						From: VERTICIES_COLLECTION_NAME + "/test1",
						To:   VERTICIES_COLLECTION_NAME + "/test2",
						Key:  "test1TOtest2",
					},
					Edge{
						From: VERTICIES_COLLECTION_NAME + "/test1",
						To:   VERTICIES_COLLECTION_NAME + "/test3",
						Key:  "test1TOtest3",
					},
				})
				if err != nil {
					fmt.Println(err.Error())
				}
				for i, e := range errors {
					if !strings.Contains(e.Error(), "conflicting key") {
						fmt.Printf("Could not add edge to graph %s: %v", meta[i].ID, e)
					}
				}

			},
			Name:              "gets edges of node in graph",
			Node:              "test1",
			ExpectedError:     nil,
			ExpectedNeighbors: []string{"test2", "test3"},
		},
		Test{
			Before:            func() {},
			Name:              "node doesn't exist",
			Node:              "sdf s/d/fas/ dfa/s####",
			ExpectedError:     nil,
			ExpectedNeighbors: []string(nil),
		},
	}

	s := Server{
		Nodes: nodes,
		Edges: edges,
		DB:    db,
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			test.Before()
			e, neighbors := GetEdges(test.Node, s)
			assert.Equal(t, test.ExpectedError, e)
			assert.Equal(t, test.ExpectedNeighbors, neighbors)
		})
	}

}

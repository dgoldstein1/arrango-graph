package main

import (
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

// connects to arango db using env vars
func ConnectToDB() (error, driver.Graph) {

	// Create an HTTP connection to the database
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		logFatalf("Failed to create HTTP connection: %v", err)
	}
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})

	// Create database
	db, err := c.CreateDatabase(nil, "my_graph_db", nil)
	if err != nil {
		logFatalf("Failed to create database: %v", err)
	}

	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = "myEdgeCollection"
	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{"myCollection1", "myCollection2"}

	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{"myCollection1", "myCollection3"}

	// A graph can contain additional vertex collections, defined in the set of orphan collections
	var options driver.CreateGraphOptions
	options.OrphanVertexCollections = []string{"myCollection4", "myCollection5"}
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	// now it's possible to create a graph
	graph, err := db.CreateGraph(nil, "myGraph", &options)
	if err != nil {
		logFatalf("Failed to create graph: %v", err)
	}
	return nil, graph
}

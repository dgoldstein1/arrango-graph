package main

import (
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"os"
	"strings"
)

// connects to arango db using env vars
func ConnectToDB() driver.Graph {
	db := establishConnectionToDb()
	options := configureGraph()
	// check if graph already exists
	var graph driver.Graph
	var err error
	if exists, _ := db.GraphExists(nil, os.Getenv("GRAPH_DB_NAME")); exists {
		// graph already exists, read current
		graph, err = db.Graph(nil, os.Getenv("GRAPH_DB_NAME"))
	} else {
		// graph does not exist, create new
		graph, err = db.CreateGraph(nil, os.Getenv("GRAPH_DB_NAME"), &options)
	}
	if err != nil {
		logFatalf("Failed to create graph: %v", err)
	}
	return graph
}

// establishes connection to DB. Exists on error
func establishConnectionToDb() driver.Database {
	// Create an HTTP connection to the database
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: strings.Split(os.Getenv("GRAPH_DB_ARANGO_ENDPOINTS"), "|"),
	})
	if err != nil {
		logFatalf("Failed to create HTTP connection: %v", err)
	}
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	// try fetching database
	exists, err := c.DatabaseExists(nil, os.Getenv("GRAPH_DB_NAME"))
	if err != nil {
		logFatalf("Could not check if databse exists create database: %v", err)
	}
	// retrieve db normally
	var db driver.Database
	if exists {
		db, err = c.Database(nil, os.Getenv("GRAPH_DB_NAME"))
	} else { // create new database
		db, err = c.CreateDatabase(nil, os.Getenv("GRAPH_DB_NAME"), nil)
	}
	if err != nil {
		logFatalf("Failed to initialize database: %v", err)
	}
	return db
}

func configureGraph() driver.CreateGraphOptions {
	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = os.Getenv("GRAPH_DB_COLLECTION_NAME")
	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{os.Getenv("GRAPH_DB_COLLECTION_NAME")}
	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{"GRAPH_DB_COLLECTION_NAME"}
	// A graph can contain additional vertex collections, defined in the set of orphan collections
	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}
	return options
}

package main

import (
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"os"
	"strings"
)

// connects to arango db using env vars
func ConnectToDB() (g driver.Graph, nodes driver.Collection, edges driver.Collection) {
	err, db := establishConnectionToDb()
	if err != nil {
		logFatalf("Could not establish connection to DB %v", err)
		return g, nodes, edges
	}
	options := configureGraph()
	// check if graph already exists
	if exists, _ := db.GraphExists(nil, os.Getenv("GRAPH_DB_NAME")); exists {
		// graph already exists, read current
		g, err = db.Graph(nil, os.Getenv("GRAPH_DB_NAME"))
	} else {
		// graph does not exist, create new
		g, err = db.CreateGraph(nil, os.Getenv("GRAPH_DB_NAME"), &options)
	}
	if err != nil {
		logFatalf("Failed to create graph: %v", err)
	}

	// initialize node and edge collections
	nodes, _ = g.VertexCollection(nil, os.Getenv("GRAPH_DB_COLLECTION_NAME"))
	edges, _, err = g.EdgeCollection(nil, "edges")
	return g, nodes, edges
}

// establishes connection to DB. Exists on error
func establishConnectionToDb() (error, driver.Database) {
	// Create an HTTP connection to the database
	urls := strings.Split(os.Getenv("GRAPH_DB_ARANGO_ENDPOINTS"), "|")
	conn, _ := http.NewConnection(http.ConnectionConfig{
		Endpoints: urls,
	})
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	// try fetching database
	exists, err := c.DatabaseExists(nil, os.Getenv("GRAPH_DB_NAME"))
	if err != nil {
		return fmt.Errorf("Could not check if databse exists create database at %v: %v", urls, err), nil
	}
	// retrieve db normally
	var db driver.Database
	if exists {
		db, err = c.Database(nil, os.Getenv("GRAPH_DB_NAME"))
	} else { // create new database
		db, err = c.CreateDatabase(nil, os.Getenv("GRAPH_DB_NAME"), nil)
	}
	if err != nil {
		return fmt.Errorf("Failed to initialize database: %v", err), nil
	}
	return nil, db
}

func configureGraph() driver.CreateGraphOptions {
	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = "edges"
	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{os.Getenv("GRAPH_DB_COLLECTION_NAME")}
	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{os.Getenv("GRAPH_DB_COLLECTION_NAME")}
	// A graph can contain additional vertex collections, defined in the set of orphan collections
	var options driver.CreateGraphOptions
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}
	return options
}

func GetEdges(node string) (err error, neighbors []string) {
	return err, neighbors
}

func AddEdges(node string, neighbors []string) (e error) {
	return e
}

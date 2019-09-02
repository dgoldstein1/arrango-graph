package main

import (
	"context"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"os"
	"strings"
)

var VERTICIES_COLLECTION_NAME = os.Getenv("GRAPH_DB_COLLECTION_NAME") + "verticies"
var EDGES_COLLECTION_NAME = os.Getenv("GRAPH_DB_COLLECTION_NAME") + "edges"
var NOT_FOUND_ERROR = "not found"

// connects to arango db using env vars
func ConnectToDB() (g driver.Graph, nodes driver.Collection, edges driver.Collection) {
	err, db := establishConnectionToDb()
	if err != nil {
		logFatalf("Could not establish connection to DB %v", err)
		return g, nodes, edges
	}
	options, nodes := configureGraph(db)
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
	// create edges if doesn't exist
	if edges, _, err = g.EdgeCollection(nil, EDGES_COLLECTION_NAME); err != nil {
		logFatalf("Error fetching edge collection", err)
	}
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

// configures collections and graph options in graph
// panics on failure
func configureGraph(db driver.Database) (options driver.CreateGraphOptions, nodes driver.Collection) {
	// create collections if they don't already exist
	ctx := context.Background()
	found, err := db.CollectionExists(ctx, VERTICIES_COLLECTION_NAME)
	if err != nil {
		logFatalf("Could not check if verticies collection exists: %v", err)
	}
	if !found {
		if nodes, err = db.CreateCollection(ctx, VERTICIES_COLLECTION_NAME, &driver.CreateCollectionOptions{}); err != nil {
			logFatalf("Could not create verticies collection: %v", err)
		}
	} else {
		// read in current collection
		if nodes, err = db.Collection(ctx, VERTICIES_COLLECTION_NAME); err != nil {
			logFatalf("Error reading existing verticies collection: %v", err)
		}
	}

	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = EDGES_COLLECTION_NAME
	// define a set of collections where an edge is going out...
	edgeDefinition.From = []string{VERTICIES_COLLECTION_NAME}
	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{VERTICIES_COLLECTION_NAME}
	// A graph can contain additional vertex collections, defined in the set of orphan collections
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}
	return options, nodes
}

func GetEdges(node string, s Server) (err error, neighbors []string) {
	return err, neighbors
}

func AddEdges(
	node string,
	neighbors []string,
	s Server,
) (
	e error,
	nodesAdded []string,
) {
	// create new nodes and edges
	nodes := []Node{Node{node}}
	edges := []Edge{}
	for _, n := range neighbors {
		nodes = append(nodes, Node{n})
		edges = append(edges, Edge{
			From: VERTICIES_COLLECTION_NAME + "/" + node,
			To:   VERTICIES_COLLECTION_NAME + "/" + n,
		})
	}
	// add all nodes to vertext collection
	newNodes, _, err := s.Nodes.CreateDocuments(nil, nodes)
	if err != nil {
		logErr("Could not create nodes: %v", err)
		return err, neighbors
	}
	metaslice, errors, err := s.Edges.CreateDocuments(nil, edges)
	if err != nil {
		logErr("Could not create edges: %v", err)
		return err, neighbors
	}
	for i, e := range errors {
		if err != nil {
			logErr("Error adding edges to node %v: %v", metaslice[i], e)
		}
	}
	// get current node
	// TODO: compare with neighbors..
	// add nodes back into []string{}
	temp := make(map[string]bool)
	for _, n := range newNodes {
		if n.Key != "" && !temp[n.Key] {
			neighbors = append(neighbors, n.Key)
			temp[n.Key] = true
		}
	}
	return e, neighbors
}

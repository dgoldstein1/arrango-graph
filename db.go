package main

import (
	"context"
	"fmt"
	driver "github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"os"
	"strings"
)

var VERTICIES_COLLECTION_NAME = "verticies"
var EDGES_COLLECTION_NAME = "edges"

// connects to arango db using env vars
func ConnectToDB() (g driver.Graph, nodes driver.Collection, edges driver.Collection, db driver.Database) {
	err, db := establishConnectionToDb()
	if err != nil {
		logFatalf("Could not establish connection to DB %v", err)
		return g, nodes, edges, db
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
		logFatalf("Error fetching edge collection: %v", err)
	}
	return g, nodes, edges, db
}

// establishes connection to DB. Exists on error
func establishConnectionToDb() (error, driver.Database) {
	// Create an HTTP connection to the database
	urls := []string{os.Getenv("GRAPH_DB_ARANGO_ENDPOINT")}
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: urls,
	})
	if err != nil {
		return fmt.Errorf("Could not create connection to DB: %v", err), nil
	}
	// Create a client
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		return fmt.Errorf("Could not create client driver: %v", err), nil
	}
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

// get destination of edges  which have FROM: "node"
func GetEdges(node string, s Server) (err error, neighbors []string) {
	// search edges which are leaving this node, going to another
	// Return the node they're going to
	query := "FOR edge IN edges FILTER edge._from == @node RETURN edge._to"
	bindVars := map[string]interface{}{
		"node": VERTICIES_COLLECTION_NAME + "/" + node,
	}
	cursor, err := s.DB.Query(context.Background(), query, bindVars)
	if err != nil {
		logErr("Error launching query %s: %v", query, err)
		return err, neighbors
	}
	defer cursor.Close()
	// read out results
	for cursor.HasMore() {
		var n string
		if _, err := cursor.ReadDocument(nil, &n); err != nil {
			logErr(err.Error())
		} else {
			nodeName := strings.TrimPrefix(n, VERTICIES_COLLECTION_NAME+"/")
			neighbors = append(neighbors, nodeName)
		}

	}
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
			Key:  VERTICIES_COLLECTION_NAME + "-" + node + "TO" + n,
		})
	}
	// add all nodes to vertext collection
	newNodes, errors, err := s.Nodes.CreateDocuments(nil, nodes)
	if err != nil {
		logErr("Could not create nodes: %v", err)
		return err, []string{}
	}
	for i, e := range errors {
		if e != nil {
			// check that is not a conflict error
			if !strings.Contains(e.Error(), "conflicting key") {
				logErr("Error adding nodes to graph %s: %s", newNodes[i].Key, e.Error())
			}
			// remove from list
			newNodes[i].Key = ""
		}
	}
	metaslice, errors, err := s.Edges.CreateDocuments(nil, edges)
	if err != nil {
		logErr("Could not create edges: %v", err)
		return err, []string{}
	}
	for i, e := range errors {
		if e != nil && !strings.Contains(e.Error(), "conflicting key") {
			logErr("Error adding edges to node %s: %s", metaslice[i].Key, e.Error())
		}
	}
	// add nodes back into []string{}
	neighbors = []string{}
	for _, n := range newNodes {
		if n.Key != "" && n.Key != node {
			neighbors = append(neighbors, n.Key)
		}
	}
	return e, neighbors
}

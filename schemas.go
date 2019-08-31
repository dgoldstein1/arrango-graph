package main

import (
	driver "github.com/arangodb/go-driver"
)

// server environment
type Server struct {
	G              driver.Graph
	Nodes          driver.Collection
	Edges          driver.Collection
	GetEdgesFromDB func(string) (error, []string)
	AddEdgesToDB   func(string, []string) (error, []string)
}

type Error struct {
	Code  int
	Error string
}

// util struct for GetEntries
type RetrievalError struct {
	Error    string // error on lookup
	NotFound bool   // is the error that it wasn't found?
}

/////////
// API //
/////////

type AddEdgesRequest struct {
	Neighbors []string `json:"neighbors"`
}

type AddEdgesResponse struct {
	NeighborsAdded []string `json:"neighborsAdded"`
}

////////////////////
// DB definitions //
////////////////////

// node in graph
type Node struct {
	Key string `json:"_key"`
}

// edge between nodes
type Edge struct {
	From string `json:"_from"`
	To   string `json:"_to"`
}

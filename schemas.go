package main

import (
	driver "github.com/arangodb/go-driver"
)

// server environment
type Server struct {
	G driver.Graph
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

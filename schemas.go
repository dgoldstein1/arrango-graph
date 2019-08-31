package main

// server environment
type Server struct {
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

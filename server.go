package main

import (
	"github.com/gin-gonic/gin"
)

// adds edges to "?node=NAME", creates nodes if they don't exist
func (s *Server) AddEdges(c *gin.Context) {}

// retrives edges to the edge "node=NAME" in uri
func (s *Server) GetEdges(c *gin.Context) {}

// finds shortest, or tied, path between two nodes
func (s *Server) ShortestPath(c *gin.Context) {}

// exports the graph to big JSON array
func (s *Server) Export(c *gin.Context) {}

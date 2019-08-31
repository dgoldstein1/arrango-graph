package main

import (
	"github.com/gin-gonic/gin"
)

// adds edges to "?node=NAME", creates nodes if they don't exist
func (s *Server) AddEdges(c *gin.Context) {}

// retrives edges to the edge "node=NAME" in uri
func (s *Server) GetEdges(c *gin.Context) {
	if c.Query("node") == "" {
		c.JSON(400, Error{400, "'node' is a required parameter"})
	}
	err, edges := s.GetEdgesFromDB(c.Query("node"))
	if err != nil {
		c.JSON(500, Error{500, err.Error()})
	}
	c.JSON(200, edges)
}

// finds shortest, or tied, path between two nodes
func (s *Server) ShortestPath(c *gin.Context) {}

// exports the graph to big JSON array
func (s *Server) Export(c *gin.Context) {}

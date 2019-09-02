package main

import (
	"github.com/gin-gonic/gin"
)

// adds edges to "?node=NAME", creates nodes if they don't exist
func (s *Server) AddEdges(c *gin.Context) {
	if c.Query("node") == "" {
		c.JSON(400, Error{400, "'node' is a required parameter"})
		return
	}
	request := AddEdgesRequest{}
	if err := c.BindJSON(&request); err != nil || len(request.Neighbors) == 0 {
		c.JSON(400, Error{400, "Bad request"})
		return
	}
	err, newNodes := s.AddEdgesToDB(c.Query("node"), request.Neighbors, *s)
	if err != nil {
		c.JSON(500, Error{500, err.Error()})
		return
	}
	c.JSON(200, AddEdgesResponse{newNodes})
}

// retrives edges to the edge "node=NAME" in uri
func (s *Server) GetEdges(c *gin.Context) {
	if c.Query("node") == "" {
		c.JSON(400, Error{400, "'node' is a required parameter"})
		return
	}
	err, edges := s.GetEdgesFromDB(c.Query("node"), *s)
	if err != nil {
		c.JSON(500, Error{500, err.Error()})
		return
	}
	c.JSON(200, edges)
}

// finds shortest, or tied, path between two nodes
func (s *Server) ShortestPath(c *gin.Context) {}

// exports the graph to big JSON array
func (s *Server) Export(c *gin.Context) {}

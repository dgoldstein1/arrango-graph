package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zsais/go-gin-prometheus"
	"net/http"
)

var (
	numberOfNodes = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "number_of_nodes",
			Help:      "Number of items in the collection GRAPH_DB_COLLECTION_NAME",
		})

	numberOfEdges = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "golang",
			Name:      "number_of_edges",
			Help:      "number of items in the 'edges' collection",
		})
)

// entrypoint
func SetupRouter(docs string) (*gin.Engine, *Server) {
	// try to connect to db
	logMsg("Connecting to DB")
	g, nodes, edges, db := ConnectToDB()
	logMsg("Done.")
	// create server object
	s := Server{g, nodes, edges, db, GetEdges, AddEdges}
	// define endpoints
	router := gin.Default()
	router.Use(gin.Logger())
	// set base page as readme html
	router.LoadHTMLGlob(docs)
	router.Static("/static", "static")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// metrics
	p := ginprometheus.NewPrometheus("gin")
	prometheus.MustRegister(numberOfNodes)
	prometheus.MustRegister(numberOfEdges)
	p.Use(router)
	// define endpoints
	router.POST("/edges", s.AddEdges)
	router.GET("/edges", s.GetEdges)
	router.GET("/shortestPath", s.ShortestPath)
	router.GET("/export", s.Export)
	// return server
	return router, &s
}

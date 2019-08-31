package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createMockRouter(s Server) *gin.Engine {
	router := gin.Default()
	// define endpoints
	router.POST("/edges", s.AddEdges)
	router.GET("/edges", s.GetEdges)
	router.GET("/shortestPath", s.ShortestPath)
	router.GET("/export", s.Export)

	return router
}

func TestAddEdges(t *testing.T) {}

func TestGetEdges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type Test struct {
		Name             string
		GetEdgesFromDB   func(string) (error, []string)
		Method           string
		Path             string
		Body             string
		ExpectedCode     int
		ExpectedResponse string
	}

	testTable := []Test{
		Test{
			Name: "positive test",
			GetEdgesFromDB: func(s string) (error, []string) {
				return nil, []string{"/wiki/test1", "/wiki/test2"}
			},
			Method:           "GET",
			Path:             "/edges?node=test1",
			ExpectedCode:     200,
			ExpectedResponse: `["test1", "test2"]`,
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create server object
			s := Server{nil, nil, nil, test.GetEdgesFromDB}
			router := createMockRouter(s)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, nil)
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code)

		})
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestAddEdges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type Test struct {
		Name             string
		AddEdgesToDB     func(string, []string) (error, []string)
		Method           string
		Path             string
		Body             []byte
		ExpectedCode     int
		ExpectedResponse AddEdgesResponse
		Error            Error
	}

	testTable := []Test{
		Test{
			Name: "positive test",
			AddEdgesToDB: func(node string, neighbors []string) (error, []string) {
				return nil, []string{"test1", "test2"}
			},
			Method:       "POST",
			Path:         "/edges?node=test",
			Body:         []byte(`{"neighbors" : ["test1", "test2"]}`),
			ExpectedCode: 200,
			ExpectedResponse: AddEdgesResponse{
				NeighborsAdded: []string{"test1", "test2"},
			},
			Error: Error{},
		},
		Test{
			Name: "no node passed",
			AddEdgesToDB: func(node string, neighbors []string) (error, []string) {
				return nil, []string{"test1", "test2"}
			},
			Method:           "POST",
			Path:             "/edges",
			Body:             []byte(`{"neighbors" : ["test1", "test2"]}`),
			ExpectedCode:     400,
			ExpectedResponse: AddEdgesResponse{},
			Error:            Error{400, "'node' is a required parameter"},
		},
		Test{
			Name: "bad request object",
			AddEdgesToDB: func(node string, neighbors []string) (error, []string) {
				return nil, []string{"test1", "test2"}
			},
			Method:           "POST",
			Path:             "/edges?node=test",
			Body:             []byte(`{"nesdfsd sd sd rs" : ["test1", "test2"]}`),
			ExpectedCode:     400,
			ExpectedResponse: AddEdgesResponse{},
			Error:            Error{400, "Bad request"},
		},
		Test{
			Name: "error adding edges",
			AddEdgesToDB: func(node string, neighbors []string) (error, []string) {
				return fmt.Errorf("node test was not found"), []string{}
			},
			Method:           "POST",
			Path:             "/edges?node=test",
			Body:             []byte(`{"neighbors" : ["test1", "test2"]}`),
			ExpectedCode:     500,
			ExpectedResponse: AddEdgesResponse{},
			Error:            Error{500, "node test was not found"},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create server object
			s := Server{
				AddEdgesToDB: test.AddEdgesToDB,
			}
			router := createMockRouter(s)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, bytes.NewBuffer(test.Body))
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code)
			body := []byte(w.Body.String())
			if test.ExpectedCode == 200 {
				resp := AddEdgesResponse{}
				err := json.Unmarshal(body, &resp)
				require.Nil(t, err)
				assert.Equal(t, test.ExpectedResponse, resp)
			} else {
				resp := Error{}
				err := json.Unmarshal(body, &resp)
				require.Nil(t, err)
				assert.Equal(t, test.Error, resp)
			}
		})
	}
}

func TestGetEdges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type Test struct {
		Name             string
		GetEdgesFromDB   func(string) (error, []string)
		Method           string
		Path             string
		ExpectedCode     int
		ExpectedResponse []string
		Error            Error
	}

	testTable := []Test{
		Test{
			Name: "positive test",
			GetEdgesFromDB: func(s string) (error, []string) {
				return nil, []string{"test1", "test2"}
			},
			Method:           "GET",
			Path:             "/edges?node=test1",
			ExpectedCode:     200,
			ExpectedResponse: []string{"test1", "test2"},
			Error:            Error{},
		},
		Test{
			Name: "no node passed",
			GetEdgesFromDB: func(s string) (error, []string) {
				return nil, []string{"test1", "test2"}
			},
			Method:           "GET",
			Path:             "/edges",
			ExpectedCode:     400,
			ExpectedResponse: []string{},
			Error:            Error{400, "'node' is a required parameter"},
		},
		Test{
			Name: "node does not exist",
			GetEdgesFromDB: func(s string) (error, []string) {
				return fmt.Errorf("node %s does not exist", s), []string{}
			},
			Method:           "GET",
			Path:             "/edges?node=6",
			ExpectedCode:     500,
			ExpectedResponse: []string{},
			Error:            Error{500, "node 6 does not exist"},
		},
	}

	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create server object
			s := Server{
				GetEdgesFromDB: test.GetEdgesFromDB,
			}
			router := createMockRouter(s)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.Method, test.Path, nil)
			req.Header.Add("Content-Type", "application/json")
			router.ServeHTTP(w, req)
			assert.Equal(t, test.ExpectedCode, w.Code)
			body := []byte(w.Body.String())
			if test.ExpectedCode == 200 {
				resp := []string{}
				err := json.Unmarshal(body, &resp)
				require.Nil(t, err)
				assert.Equal(t, test.ExpectedResponse, resp)
			} else {
				resp := Error{}
				err := json.Unmarshal(body, &resp)
				require.Nil(t, err)
				assert.Equal(t, test.Error, resp)
			}
		})
	}
}

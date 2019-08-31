package main

import (
	"encoding/json"
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

func TestAddEdges(t *testing.T) {}

func TestGetEdges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type Test struct {
		Name             string
		GetEdgesFromDB   func(string) (error, []string)
		Method           string
		Path             string
		ExpectedCode     int
		ExpectedResponse []string
		Error            string
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
			Error:            "",
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
				assert.Equal(t, test.ExpectedResponse, resp)
			}
		})
	}
}

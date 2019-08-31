package main

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"testing"
)

func TestAddEdges(t *testing.T) {}

func TestGetEdges(t *testing.T) {
	type Test struct {
		Name           string
		GetEdgesFromDB func(string) (error, []string)
		Params         []gin.Param
	}

	testTable := []Test{
		Test{
			Name: "positive test",
			GetEdgesFromDB: func(s string) (error, []string) {
				return nil, []string{"/wiki/test1", "/wiki/test2"}
			},
			Params: []gin.Param{gin.Param{Key: "node", Value: "/wiki/test"}},
		},
	}
	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create server object
			s := Server{nil, nil, nil, test.GetEdgesFromDB}
			// create context
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = test.Params
			s.GetEdges(c)
		})
	}
}

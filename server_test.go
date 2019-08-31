package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestAddEdges(t *testing.T) {}

func TestGetEdges(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type Test struct {
		Name             string
		GetEdgesFromDB   func(string) (error, []string)
		Params           []gin.Param
		ExpectedCode     int
		ExpectedResponse string
	}

	testTable := []Test{
		Test{
			Name: "positive test",
			GetEdgesFromDB: func(s string) (error, []string) {
				return nil, []string{"/wiki/test1", "/wiki/test2"}
			},
			Params:           []gin.Param{gin.Param{Key: "node", Value: "/wiki/test"}},
			ExpectedCode:     200,
			ExpectedResponse: `["/wiki/test1", "/wiki/test2"]`,
		},
	}
	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {
			// create server object
			s := Server{nil, nil, nil, test.GetEdgesFromDB}
			// create context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Params = test.Params

			s.GetEdges(c)
			assert.Equal(t, test.ExpectedCode, w.Code)
			b, _ := ioutil.ReadAll(w.Body)
			assert.Equal(t, test.ExpectedResponse, string(b))
		})
	}
}

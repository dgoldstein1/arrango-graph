package main

import (
	"fmt"
	// "github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	originLogPrintf := logMsg
	defer func() { logMsg = originLogPrintf }()
	logs := []string{}
	logMsg = func(format string, args ...interface{}) {
		if len(args) > 0 {
			logs = append(logs, fmt.Sprintf(format, args))
		} else {
			logs = append(logs, format)
		}
	}
	// setup environment
	// dbName := "graph-testing-router"
	// os.Setenv("GRAPH_DB_NAME", dbName)
	// os.Setenv("GRAPH_DB_COLLECTION_NAME", "graph-router-wikipedia")
	// os.Setenv("GRAPH_DB_ARANGO_ENDPOINTS", "http://localhost:8529")
	// os.Setenv("GRAPH_DB_NAME", "wikipedia-graph")
	// run test
	// router, err := SetupRouter("./api/*")
	// assert.NotNil(t, router)
	// assert.NotNil(t, err)
	// fmt.Println(logs)
}

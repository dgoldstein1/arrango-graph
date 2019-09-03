package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
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
	logErr(os.Getenv("GRAPH_DB_ARANGO_ENDPOINT"))
	router, err := SetupRouter("./api/*")
	assert.NotNil(t, router)
	assert.NotNil(t, err)
}

package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetupRouter(t *testing.T) {
	router, err := SetupRouter("./api/*")
	assert.NotNil(t, router)
	assert.NotNil(t, err)
}

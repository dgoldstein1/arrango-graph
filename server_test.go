package main

import (
	"testing"
)

func TestAddEdges(t *testing.T) {}

func TestGetEdges(t *testing.T) {
	type Test struct {
		Name     string
		Node     string
		GetEdges func(string) (error, []string)
	}

	testTable := []Test{
		Test{
			Name: "positive test",
			Node: "/wiki/test",
			GetEdges: func(s string) (error, []string) {
				return nil, []string{"/wiki/test1", "/wiki/test2"}
			},
		},
	}
	for _, test := range testTable {
		t.Run(test.Name, func(t *testing.T) {

		})
	}
}

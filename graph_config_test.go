package graph_engine_go

import (
	"fmt"
	"testing"
)

func TestGraphConfig(t *testing.T) {
	var jsonPath string = "./config/test/main_graph.json"
	graphConfig, err := LoadGraphConfig(jsonPath)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(*graphConfig)
}

package graph_engine_go

import "testing"

func TestGraphContextBuildGraph(t *testing.T) {
	var jsonPath string = "./config/test/graph_test.json"
	graphConfig, _ := LoadGraphConfig(jsonPath)
	var ctx GraphContext
	ctx.BuildGraph(graphConfig)
}

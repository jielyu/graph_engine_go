package graph_engine_go

import "fmt"

// import "fmt"

type GraphContext struct {
	graphConfig *GraphConfig
	graphNodes  map[string]*GraphOperator
}

func (ctx *GraphContext) BuildGraph(graphConfig *GraphConfig) error {
	ctx.graphConfig = graphConfig
	// 创建所有节点
	for _, nodeconfig := range graphConfig.Nodes {
		fmt.Println(nodeconfig.Name)
	}
	return nil
}

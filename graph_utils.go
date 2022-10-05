package graph_engine_go

import "fmt"

// 用于SetUp函数中创建GraphData
func GetGraphDataByName(ctx *GraphContext, config *GraphNodeConfig, name string) *GraphData {
	// 检查发布数据名称是否存在
	if _, ok := config.Emitters[name]; !ok {
		panic(fmt.Errorf("not found Emitter '%s' in node '%s'", name, config.Name))
	}
	return ctx.NameEmit(config.Emitters[name])
}

// 用于SetUp函数中创建GraphDep
func GetGraphDependByName(ctx *GraphContext, config *GraphNodeConfig, name string) *GraphDep {
	// 检查发布数据名称是否存在
	if _, ok := config.Depends[name]; !ok {
		panic(fmt.Errorf("not found Depend '%s' in node '%s'", name, config.Name))
	}
	return ctx.NameDepend(config.Depends[name])
}

package graph_engine_go

// 所有节点都需要实现的基础接口
type GraphOperator interface {
	// 设置并记录节点的config信息
	SetConfig(config *GraphNodeConfig) error
	// 建立图结构，主要用于创建依赖和发布的数据结构，只在创建时执行一次
	SetUp(ctx *GraphContext) error
	// 初始化操作，一般用于设置初始化参数，只在创建时执行一次
	Initailize(ctx *GraphContext) error
	// 执行任务，业务逻辑主要封装在这个函数，每次都会执行
	Process(ctx *GraphContext) error
}

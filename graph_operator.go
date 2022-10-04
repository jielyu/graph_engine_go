package graph_engine_go

type GraphOperator interface {
	SetConfig(config *GraphNodeConfig) error
	SetUp(ctx *GraphContext) error
	Initailize(ctx *GraphContext) error
	Process(ctx *GraphContext) error
}

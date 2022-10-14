package nodes

import (
	"fmt"

	ge "github.com/jielyu/graph_engine_go"
)

type AddOp struct {
	Config *ge.GraphNodeConfig
	a      *ge.GraphDep
	b      *ge.GraphDep
	c      *ge.GraphData
}

func (g *AddOp) SetConfig(config *ge.GraphNodeConfig) error {
	g.Config = config
	return nil
}

func (g *AddOp) SetUp(ctx *ge.GraphContext) error {
	g.a = ge.GetGraphDependByName(ctx, g.Config, "A")
	g.b = ge.GetGraphDependByName(ctx, g.Config, "B")
	g.c = ge.GetGraphDataByName(ctx, g.Config, "C")
	fmt.Printf("compute: setup '%s' node '%s' succ.%p\r\n", g.Config.NodeType, g.Config.Name, g.Config)
	return nil
}

func (g *AddOp) Initailize(ctx *ge.GraphContext) error {
	return nil
}

func (g *AddOp) Process(ctx *ge.GraphContext) error {
	a := ge.Dep[int](g.a)
	b := ge.Dep[int](g.b)
	c := ge.Emit[int](g.c)
	*c = *a + *b
	return nil
}

func init() {
	ge.Register[AddOp]()
}

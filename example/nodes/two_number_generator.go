package nodes

import (
	"fmt"

	ge "github.com/jielyu/graph_engine_go"
)

type TwoNumbersGeneratorOp struct {
	Config *ge.GraphNodeConfig
	a      *ge.GraphData
	b      *ge.GraphData
}

func (g *TwoNumbersGeneratorOp) SetConfig(config *ge.GraphNodeConfig) error {
	g.Config = config
	return nil
}

func (g *TwoNumbersGeneratorOp) SetUp(ctx *ge.GraphContext) error {
	g.a = ge.GetGraphDataByName(ctx, g.Config, "A")
	g.b = ge.GetGraphDataByName(ctx, g.Config, "B")
	fmt.Printf("input: setup '%s' node '%s' succ.%p\r\n", g.Config.NodeType, g.Config.Name, g.Config)
	return nil
}

func (g *TwoNumbersGeneratorOp) Initailize(ctx *ge.GraphContext) error {
	return nil
}

func (g *TwoNumbersGeneratorOp) Process(ctx *ge.GraphContext) error {
	a := ge.Emit[int](g.a)
	b := ge.Emit[int](g.b)
	*a = 2
	*b = 3
	return nil
}

func init() {
	ge.Register[TwoNumbersGeneratorOp]()
}

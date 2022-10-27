package nodes

import (
	"fmt"
	"strconv"

	ge "github.com/jielyu/graph_engine_go"
)

type TwoNumbersGeneratorOp struct {
	Config *ge.GraphNodeConfig
	a      *ge.GraphData
	b      *ge.GraphData
	aVal   int
	bVal   int
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
	// 从配置中解析a参数
	aVal, ok := g.Config.Config["A"]
	if ok {
		g.aVal, _ = strconv.Atoi(aVal)
	} else {
		g.aVal = 2
	}
	// 从配置中解析b参数
	bVal, ok := g.Config.Config["B"]
	if ok {
		g.bVal, _ = strconv.Atoi(bVal)
	} else {
		g.bVal = 3
	}
	fmt.Printf("node[%s], a=%d, b=%d\r\n", g.Config.Name, g.aVal, g.bVal)
	return nil
}

func (g *TwoNumbersGeneratorOp) Process(ctx *ge.GraphContext) error {
	a := ge.Emit[int](g.a)
	b := ge.Emit[int](g.b)
	*a = g.aVal
	*b = g.bVal
	return nil
}

func init() {
	ge.Register[TwoNumbersGeneratorOp]()
}

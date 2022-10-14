package nodes

import (
	"fmt"

	ge "github.com/jielyu/graph_engine_go"
)

type PrinterOp struct {
	Config *ge.GraphNodeConfig
	a      *ge.GraphDep
}

func (g *PrinterOp) SetConfig(config *ge.GraphNodeConfig) error {
	g.Config = config
	return nil
}

func (g *PrinterOp) SetUp(ctx *ge.GraphContext) error {
	g.a = ge.GetGraphDependByName(ctx, g.Config, "A")
	fmt.Printf("output: setup '%s' node '%s' succ.%p\r\n", g.Config.NodeType, g.Config.Name, g.Config)
	return nil
}

func (g *PrinterOp) Initailize(ctx *ge.GraphContext) error {

	return nil
}

func (g *PrinterOp) Process(ctx *ge.GraphContext) error {
	a := ge.Dep[int](g.a)
	// time.Sleep(time.Duration(2) * time.Microsecond)
	fmt.Printf("result=%d\r\n", *a)
	return nil
}

func init() {
	ge.Register[PrinterOp]()
}

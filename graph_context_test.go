package graph_engine_go

import (
	"fmt"
	"testing"
)

type TwoNumbersGeneratorOp struct {
	Config *GraphNodeConfig
	a      *GraphData
	b      *GraphData
}

func (g *TwoNumbersGeneratorOp) SetConfig(config *GraphNodeConfig) error {
	g.Config = config
	return nil
}

func (g *TwoNumbersGeneratorOp) SetUp(ctx *GraphContext) error {
	g.a = GetGraphDataByName(ctx, g.Config, "A")
	g.b = GetGraphDataByName(ctx, g.Config, "B")
	fmt.Printf("input: setup '%s' node '%s' succ.%p\r\n", g.Config.NodeType, g.Config.Name, g.Config)
	return nil
}

func (g *TwoNumbersGeneratorOp) Initailize(ctx *GraphContext) error {
	return nil
}

func (g *TwoNumbersGeneratorOp) Process(ctx *GraphContext) error {
	a := Emit[int](g.a)
	b := Emit[int](g.b)
	*a = 2
	*b = 3
	return nil
}

type AddOp struct {
	Config *GraphNodeConfig
	a      *GraphDep
	b      *GraphDep
	c      *GraphData
}

func (g *AddOp) SetConfig(config *GraphNodeConfig) error {
	g.Config = config
	return nil
}

func (g *AddOp) SetUp(ctx *GraphContext) error {
	g.a = GetGraphDependByName(ctx, g.Config, "A")
	g.b = GetGraphDependByName(ctx, g.Config, "B")
	g.c = GetGraphDataByName(ctx, g.Config, "C")
	fmt.Printf("compute: setup '%s' node '%s' succ.%p\r\n", g.Config.NodeType, g.Config.Name, g.Config)
	return nil
}

func (g *AddOp) Initailize(ctx *GraphContext) error {
	return nil
}

func (g *AddOp) Process(ctx *GraphContext) error {
	a := Dep[int](g.a)
	b := Dep[int](g.b)
	c := Emit[int](g.c)
	*c = *a + *b
	return nil
}

type PrinterOp struct {
	Config *GraphNodeConfig
	a      *GraphDep
}

func (g *PrinterOp) SetConfig(config *GraphNodeConfig) error {
	g.Config = config
	return nil
}

func (g *PrinterOp) SetUp(ctx *GraphContext) error {
	g.a = GetGraphDependByName(ctx, g.Config, "A")
	fmt.Printf("output: setup '%s' node '%s' succ.%p\r\n", g.Config.NodeType, g.Config.Name, g.Config)
	return nil
}

func (g *PrinterOp) Initailize(ctx *GraphContext) error {

	return nil
}

func (g *PrinterOp) Process(ctx *GraphContext) error {
	a := Dep[int](g.a)
	fmt.Printf("result=%d\r\n", *a)
	return nil
}

func init() {
	RegisterClass("TwoNumbersGeneratorOp", func() GraphOperator { return new(TwoNumbersGeneratorOp) })
	RegisterClass("AddOp", func() GraphOperator { return new(AddOp) })
	RegisterClass("PrinterOp", func() GraphOperator { return new(PrinterOp) })
}

func TestGraphContextBuildGraph(t *testing.T) {
	var jsonPath string = "./config/test/main_graph.json"
	graphConfig, err := LoadGraphConfig(jsonPath)
	if err != nil {
		panic((err.Error()))
	}
	var ctx = NewGraphContext()
	ctx.Build(graphConfig)
	ctx.Process()
}

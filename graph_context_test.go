package graph_engine_go

import "testing"

type TwoNumbersGeneratorOp struct {
	Config *GraphNodeConfig
}

func (g *TwoNumbersGeneratorOp) SetConfig(config *GraphNodeConfig) error {

	return nil
}

func (g *TwoNumbersGeneratorOp) SetUp(ctx *GraphContext) error {
	return nil
}

func (g *TwoNumbersGeneratorOp) Initailize(ctx *GraphContext) error {
	return nil
}

func (g *TwoNumbersGeneratorOp) Process(ctx *GraphContext) error {
	return nil
}

type AddOp struct {
	Config *GraphNodeConfig
}

func (g *AddOp) SetConfig(config *GraphNodeConfig) error {

	return nil
}

func (g *AddOp) SetUp(ctx *GraphContext) error {
	return nil
}

func (g *AddOp) Initailize(ctx *GraphContext) error {
	return nil
}

func (g *AddOp) Process(ctx *GraphContext) error {
	return nil
}

type PrinterOp struct {
	Config *GraphNodeConfig
}

func (g *PrinterOp) SetConfig(config *GraphNodeConfig) error {

	return nil
}

func (g *PrinterOp) SetUp(ctx *GraphContext) error {
	return nil
}

func (g *PrinterOp) Initailize(ctx *GraphContext) error {
	return nil
}

func (g *PrinterOp) Process(ctx *GraphContext) error {
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
	ctx.BuildGraph(graphConfig)
}

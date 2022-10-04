package graph_engine_go

import (
	"fmt"
	"testing"
)

type GraphOpTest struct {
	Name string
}

func (g *GraphOpTest) SetConfig(config *GraphNodeConfig) error {

	return nil
}

func (g *GraphOpTest) SetUp(ctx *GraphContext) error {
	return nil
}

func (g *GraphOpTest) Initailize(ctx *GraphContext) error {
	return nil
}

func (g *GraphOpTest) Process(ctx *GraphContext) error {
	fmt.Printf("obj.Name=%s\r\n", g.Name)
	return nil
}

func newGraphOpTest() GraphOperator {
	op := new(GraphOpTest)
	op.Name = "Hello"
	return op
}

func init() {
	RegisterClass("GraphOpTest", newGraphOpTest)
}

func TestCreateNode(t *testing.T) {
	obj, _ := CreateInstance("GraphOpTest")
	var ctx GraphContext
	obj.Process(&ctx)
}

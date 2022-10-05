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

func init() {
	Register[GraphOpTest]()
}

func TestCreateNode(t *testing.T) {
	obj, _ := CreateInstance("GraphOpTest")
	obj2, _ := CreateInstance("GraphOpTest")
	obj.(*GraphOpTest).Name = "hello"
	obj2.(*GraphOpTest).Name = "world"
	var ctx GraphContext
	obj.Process(&ctx)
	obj2.Process(&ctx)
}

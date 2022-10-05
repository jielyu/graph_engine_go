package graph_engine_go

import (
	"fmt"
	"testing"
)

func TestGraphDataEmitDependInt(t *testing.T) {
	graphData := new(GraphData)
	graphData.Name = "var1"
	v := Emit[int](graphData)
	*v = 12
	fmt.Println("GraphData:", *v)

	graphDep := new(GraphDep)
	graphDep.RefGraphData = graphData

	depV := Dep[int](graphDep)
	fmt.Println("GraphDep:", *depV)
}

type DataType struct {
	Name string
	Age  int
}

func TestGraphDataEmitDependStruct(t *testing.T) {
	graphData := new(GraphData)
	graphData.Name = "var1"
	v := Emit[DataType](graphData)
	v.Name = "hello"
	v.Age = 12
	fmt.Println("GraphData:", *v)

	graphDep := new(GraphDep)
	graphDep.RefGraphData = graphData

	depV := Dep[DataType](graphDep)
	fmt.Println("GraphDep:", *depV)
}

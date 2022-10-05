package graph_engine_go

import (
	"fmt"
	"reflect"
)

/***
用于管理节点发布的数据
*/
type GraphData struct {
	Name     string
	Active   bool
	TypeName string
	Data     interface{}
}

func Emit[T any](gdata *GraphData) *T {
	// 无法直接从类型获取类型名，只能先创建一个临时对象，再用反射机制
	// 后续有更好的方式再优化
	var tmp T
	typeOfT := reflect.TypeOf(tmp)
	typeTName := typeOfT.Name()
	if gdata.Data == nil {
		var d *T = new(T)
		gdata.TypeName = typeTName
		gdata.Data = d
	}
	if typeTName != gdata.TypeName {
		panic(fmt.Errorf("current emitter type '%s' different from history type '%s'",
			typeTName, gdata.TypeName))
	}
	return (gdata.Data).(*T)
}

func EmitDep[T any](gdata *GraphData, gdep *GraphDep) *T {
	gdata.Data = gdep.RefGraphData.Data
	return gdata.Data.(*T)
}

/***
用于管理节点依赖的数据
*/
type GraphDep struct {
	Mutable      bool
	RefGraphData *GraphData
}

func Dep[T any](gdep *GraphDep) *T {
	var tmp T
	typeOfT := reflect.TypeOf(tmp)
	typeTName := typeOfT.Name()
	if typeTName != gdep.RefGraphData.TypeName {
		panic(fmt.Errorf("depend type '%s' different from emitter type '%s'",
			typeTName, gdep.RefGraphData.TypeName))
	}
	return gdep.RefGraphData.Data.(*T)
}

func (gdep *GraphDep) SetMutable(m bool) {
	gdep.Mutable = m
}

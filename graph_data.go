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
	// 获取类型名，nil不会占用存储空间，避免内存碎片化
	t := reflect.TypeOf((*T)(nil)).Elem()
	typeTName := t.Name()
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
	t := reflect.TypeOf((*T)(nil)).Elem()
	typeTName := t.Name()
	if typeTName != gdep.RefGraphData.TypeName {
		panic(fmt.Errorf("depend type '%s' different from emitter type '%s'",
			typeTName, gdep.RefGraphData.TypeName))
	}
	return gdep.RefGraphData.Data.(*T)
}

func (gdep *GraphDep) SetMutable(m bool) {
	gdep.Mutable = m
}

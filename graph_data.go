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
		if len(typeTName) == 0 {
			panic(fmt.Errorf("not allow unknown type GraphData"))
		}
		var d *T = new(T)
		gdata.TypeName = typeTName
		gdata.Data = d
	}
	if typeTName != gdata.TypeName {
		panic(fmt.Errorf("current emitter type '%s' different from history type '%s'",
			typeTName, gdata.TypeName))
	}
	gdata.Active = true
	return (gdata.Data).(*T)
}

// 直接发布依赖的数据，用于inplace形式的修改，节省内存
// 模版不能用于类型方法，因此只能单独提供函数
func EmitDep[T any](gdata *GraphData, gdep *GraphDep) *T {
	// 必须声明为mutable才能原处修改
	if !gdep.Mutable {
		panic(fmt.Errorf("not allow emit immutable GraphDep[%s]", gdep.Name))
	}
	gdata.TypeName = gdep.RefGraphData.TypeName
	gdata.Data = gdep.RefGraphData.Data
	gdata.Active = true
	return gdata.Data.(*T)
}

/***
用于管理节点依赖的数据
*/
type GraphDep struct {
	Name         string
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
	if !gdep.RefGraphData.Active {
		panic(fmt.Errorf("GraphDep[%s] get inactive GraphData[%s], please check emitter", gdep.Name, gdep.RefGraphData.Name))
	}
	return gdep.RefGraphData.Data.(*T)
}

func (gdep *GraphDep) SetMutable() {
	gdep.Mutable = true
}

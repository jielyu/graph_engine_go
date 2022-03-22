package types

/***
用于管理节点发布的数据
*/
type GraphData struct {
	Name     string
	Active   bool
	TypeName string
	Data     *interface{}
}

func (gdata *GraphData) Emit() *interface{} {
	return gdata.Data
}

func (gdata *GraphData) EmitDep(gdep *GraphDep) *interface{} {
	gdata.Data = gdep.RefGraphData.Data
	return gdata.Data
}

/***
用于管理节点依赖的数据
*/
type GraphDep struct {
	Mutable      bool
	RefGraphData *GraphData
}

func (gdep *GraphDep) Dep() *interface{} {
	return gdep.RefGraphData.Data
}

func (gdep *GraphDep) SetMutable(m bool) {
	gdep.Mutable = m
}

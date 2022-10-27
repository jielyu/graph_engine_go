package graph_engine_go

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	mapset "github.com/deckarep/golang-set"
)

type GraphContext struct {
	graphConfig    *GraphConfig             // 图结构配置信息
	graphNodes     map[string]GraphOperator // 存储所有节点对象
	nodeEmitMap    map[string]mapset.Set    // 存储每个节点发布的数据名字
	nodeDepMap     map[string]mapset.Set    // 存储每个节点依赖数据的名字
	inputNodes     []string                 // 存储所有需要运行的输入节点的名字
	computeNodes   []string                 // 存储所有需要运行的计算节点的名字
	outputNodes    []string                 // 存储所有需要运行的输出节点名字
	allNodes       []string                 // 存储所有需要运行的节点名字
	allGraphData   map[string]*GraphData    // 存储所有发布的数据
	allGraphDep    map[string]*GraphDep     // 存储所有依赖的数据结构
	dataByNodeMap  map[string]string        // 用于存储发布数据与发布节点之间的映射
	dataForNodeMap map[string][]string      // 用于存储发布的数据与依赖该数据的节点之间的映射
	InputData      interface{}
	OutputData     interface{}
	Busy           bool   // ctx是否处于繁忙的标志
	Id             int    // ctx在缓冲池中的Id
	ReqId          uint64 // 每次请求尽量使用不同的ReqId，便于日志查询
}

// 创建GraphContext对象
func NewGraphContext() *GraphContext {
	ctx := new(GraphContext)
	ctx.graphNodes = make(map[string]GraphOperator, 16)
	ctx.allGraphData = make(map[string]*GraphData, 32)
	ctx.allGraphDep = make(map[string]*GraphDep, 32)
	ctx.dataByNodeMap = make(map[string]string, 32)
	ctx.dataForNodeMap = make(map[string][]string, 32)
	return ctx
}

// 用于创建发布数据的对象
func (ctx *GraphContext) NameEmit(name string) *GraphData {
	// 如果发布的数据已经存在，则直接返回。只可能是NameDepend()中发布的，不会重复
	if _, ok := ctx.allGraphData[name]; ok {
		return ctx.allGraphData[name]
	}
	// 创建发布数据结构
	graphData := new(GraphData)
	graphData.Name = name
	graphData.Data = nil
	ctx.allGraphData[name] = graphData
	// fmt.Printf("create GraphData '%s' succ\r\n", name)
	return graphData
}

// 用于创建依赖数据的对象
func (ctx *GraphContext) NameDepend(name string) *GraphDep {
	// 依赖数据已经存在的情况下，直接返回
	if _, ok := ctx.allGraphDep[name]; ok {
		return ctx.allGraphDep[name]
	}
	// 依赖尚未发布的数据，先操作发布，创建时已经避免重复，这里不需检查
	if _, ok := ctx.allGraphData[name]; !ok {
		ctx.NameEmit(name)
	}
	// 创建依赖数据结构
	graphDep := new(GraphDep)
	graphDep.Name = name
	graphDep.RefGraphData = ctx.allGraphData[name]
	ctx.allGraphDep[name] = graphDep
	return graphDep
}

// 创建图结构
func (ctx *GraphContext) Build(graphConfig *GraphConfig) error {
	// 创建节点
	err := ctx.create(graphConfig)
	if err != nil {
		return err
	}
	// 创建数据
	err = ctx.setup()
	if err != nil {
		return err
	}
	// 初始化
	err = ctx.initailize()
	if err != nil {
		return err
	}
	return nil
}

// 运行图结构
func (ctx *GraphContext) Process() error {
	// 记录节点依赖情况
	nodeState := make(map[string]*int32)
	for _, name := range ctx.allNodes {
		var stat int32 = int32(ctx.nodeDepMap[name].Cardinality())
		nodeState[name] = &stat
	}
	// 创建Group同步
	var wg sync.WaitGroup
	wg.Add(len(nodeState))
	// 使用切片记录所有就绪的节点
	readyNode := make([]string, 0, len(nodeState))
	for len(nodeState) > 0 {
		// 查找准备就绪的节点
		for k, v := range nodeState {
			if *v == 0 {
				delete(nodeState, k)
				readyNode = append(readyNode, k)
			}
		}
		// 运行准备就绪的节点
		for _, name := range readyNode {
			go func(nname string) {
				// fmt.Println("process node ", nname)
				op := ctx.graphNodes[nname]
				op.Process(ctx)
				// 更新所有依赖该节点的发布数据的节点状态
				for v := range ctx.nodeEmitMap[nname].Iter() {
					emitName := v.(string)
					for _, nodeName := range ctx.dataForNodeMap[emitName] {
						atomic.AddInt32(nodeState[nodeName], -1)
					}
				}
				wg.Done()
			}(name)
		}
		// 清空就绪列表
		readyNode = readyNode[:0]
		// 睡眠1ms，避免占用太多计算资源
		time.Sleep(time.Duration(1) * time.Microsecond)
	}
	// 避免节点还没运行完就返回
	wg.Wait()
	return nil
}

// 构建各个节点
func (ctx *GraphContext) create(graphConfig *GraphConfig) error {
	ctx.graphConfig = graphConfig
	// 创建所有节点
	nodeEmitMap := make(map[string]mapset.Set, 16)
	nodeDepMap := make(map[string]mapset.Set, 16)
	for _, nodeconfig := range graphConfig.Nodes {
		nodeName := nodeconfig.Name
		typeName := nodeconfig.NodeType
		// 创建对象
		node, err := CreateInstance(typeName)
		if err != nil {
			return err
		}
		// 设置节点配置
		// fmt.Println(nodeName, nodeconfig)
		var configCopy GraphNodeConfig = nodeconfig
		node.SetConfig(&configCopy)
		// 记录节点
		if _, ok := ctx.graphNodes[nodeName]; ok {
			return fmt.Errorf("not allow duplicated node name: %s", nodeName)
		}
		ctx.graphNodes[nodeName] = node
		// 设置节点发布数据的信息
		nodeEmitMap[nodeName] = mapset.NewSet()
		for _, emitName := range nodeconfig.Emitters {
			// 避免同一个节点发布两份一样的数据
			if nodeEmitMap[nodeName].Contains(emitName) {
				return fmt.Errorf("not allow duplicated emitter name '%s' in node '%s'", emitName, nodeName)
			}
			nodeEmitMap[nodeName].Add(emitName)
			// 记录发布数据的节点，同时避免不同节点重复发布数据
			if _, ok := ctx.dataByNodeMap[emitName]; ok {
				return fmt.Errorf("not allow duplicated emitter name '%s' in node '%s' and '%s'",
					emitName, nodeName, ctx.dataByNodeMap[emitName])
			}
			ctx.dataByNodeMap[emitName] = nodeName
		}
		// 设置节点依赖数据的信息
		nodeDepMap[nodeName] = mapset.NewSet()
		for _, depName := range nodeconfig.Depends {
			if nodeDepMap[nodeName].Contains(depName) {
				return fmt.Errorf("not allow duplicated depend name '%s' in node '%s'", depName, nodeName)
			}
			nodeDepMap[nodeName].Add(depName)
			// 记录依赖数据的所有节点名字
			ctx.dataForNodeMap[depName] = append(ctx.dataForNodeMap[depName], nodeName)
		}
		// 判断节点类型
		if len(nodeconfig.Emitters) == 0 {
			ctx.outputNodes = append(ctx.outputNodes, nodeName)
		} else if len(nodeconfig.Depends) == 0 {
			ctx.inputNodes = append(ctx.inputNodes, nodeName)
		} else {
			ctx.computeNodes = append(ctx.computeNodes, nodeName)
		}
	}
	ctx.nodeEmitMap = nodeEmitMap
	ctx.nodeDepMap = nodeDepMap
	fmt.Printf("num_inputs:%d, num_compute:%d, num_output:%d\r\n",
		len(ctx.inputNodes), len(ctx.computeNodes), len(ctx.outputNodes))
	// 找出所有对输出有贡献的节点
	var nodeFound []string
	validNodes := mapset.NewSet()
	for _, name := range ctx.outputNodes {
		nodeFound = append(nodeFound, name)
		validNodes.Add(name)
	}
	for idx := 0; idx < len(nodeFound); idx++ {
		name := nodeFound[idx]
		depends := nodeDepMap[name]
		for v := range depends.Iter() {
			// fmt.Printf("node %s, depend: %s\r\n", name, v)
			// 查找当前依赖数据的发布节点
			for emitName, c := range nodeEmitMap {
				if c.Contains(v) {
					if !validNodes.Contains(emitName) {
						validNodes.Add(emitName)
						nodeFound = append(nodeFound, emitName)
					}
					break
				}
			}
		}
	}
	ctx.allNodes = nodeFound
	fmt.Printf("found valid nodes: %v\r\n", ctx.allNodes)
	// 记录所有需要运行的节点
	ctx.inputNodes = ctx.inputNodes[:0]
	ctx.computeNodes = ctx.computeNodes[:0]
	ctx.outputNodes = ctx.outputNodes[:0]
	for _, name := range nodeFound {
		// 不依赖数据的节点为输入节点
		if nodeDepMap[name].Cardinality() == 0 {
			ctx.inputNodes = append(ctx.inputNodes, name)
		} else if nodeEmitMap[name].Cardinality() == 0 {
			ctx.outputNodes = append(ctx.outputNodes, name)
		} else {
			ctx.computeNodes = append(ctx.computeNodes, name)
		}
	}
	fmt.Printf("input:%v, compute:%v, output:%v\r\n",
		ctx.inputNodes, ctx.computeNodes, ctx.outputNodes)
	return nil
}

// 调用每一个node的SetUp函数
func (ctx *GraphContext) setup() error {
	for _, name := range ctx.allNodes {
		// fmt.Printf("start to setup node '%s'\r\n", name)
		err := ctx.graphNodes[name].SetUp(ctx)
		if err != nil {
			fmt.Printf("node[%s]:failed to run setup()", name)
			return err
		}
		// 获取当前节点应该发布和依赖的数据
		emitSet := mapset.NewSet()
		depSet := mapset.NewSet()
		for v := range ctx.nodeEmitMap[name].Iter() {
			emitName := v.(string)
			emitSet.Add(emitName)
		}
		for v := range ctx.nodeDepMap[name].Iter() {
			depName := v.(string)
			depSet.Add(depName)
		}
		// 检查节点发布数据是否足够
		nodeType := reflect.TypeOf(ctx.graphNodes[name]).Elem()
		nodeVal := reflect.ValueOf(ctx.graphNodes[name]).Elem()
		for i := 0; i < nodeType.NumField(); i++ {
			fieldType := nodeType.Field(i)
			if fieldType.Type.Kind().String() == "ptr" {
				elemName := fieldType.Type.Elem().Name()
				if elemName == "GraphData" {
					// 获取字段，并校验空指针
					nodeField := nodeVal.Field(i)
					if nodeField.IsNil() {
						panic(fmt.Errorf("node[%s]:not allow to make nil GraphData[%s] in setup()", name, fieldType.Name))
					}
					// 获取发布数据的名字
					emitName := nodeField.Elem().FieldByName("Name").String()
					if emitSet.Contains(emitName) {
						emitSet.Remove(emitName)
					} else {
						panic(fmt.Errorf("node[%s]:not allow to emit GraphData[%s] not found in config", name, emitName))
					}
				} else if elemName == "GraphDep" {
					// 获取字段，并校验空指针
					nodeField := nodeVal.Field(i)
					if nodeField.IsNil() {
						panic(fmt.Errorf("node[%s]:not allow to make nil GraphDep[%s] in setup()", name, fieldType.Name))
					}
					// 获取依赖数据的名字
					depName := nodeField.Elem().FieldByName("Name").String()
					if depSet.Contains(depName) {
						depSet.Remove(depName)
					} else {
						panic(fmt.Errorf("node[%s]:not allow to depend GraphDep[%s] not found in config", name, depName))
					}
				}
			}
		}
		// 检查发布数据是否发布完全
		if emitSet.Cardinality() > 0 {
			notEmitData := make([]string, 0, 4)
			for v := range emitSet.Iter() {
				notEmitData = append(notEmitData, v.(string))
			}
			panic(fmt.Errorf("node[%s]:not allow to miss GraphData[%s]", name, strings.Join(notEmitData[:], ",")))
		}
		// 检查依赖数据是否依赖完全
		if depSet.Cardinality() > 0 {
			notDepData := make([]string, 0, 4)
			for v := range depSet.Iter() {
				notDepData = append(notDepData, v.(string))
				panic(fmt.Errorf("node[%s]:not allow to miss GraphDep[%s]", name, strings.Join(notDepData[:], ",")))
			}
		}
	}
	// 检查mutable依赖是否唯一， 不唯一则报错
	for dataName, nodeNames := range ctx.dataForNodeMap {
		if ctx.allGraphDep[dataName].Mutable {
			if len(nodeNames) > 1 {
				nodeNameStr := strings.Join(nodeNames[:], ",")
				panic(fmt.Errorf("not allow >1 nodes depend mutable GraphDep[%s] in Nodes %s", dataName, nodeNameStr))
			}
		}
	}
	fmt.Printf("setup all nodes successfully.\r\n")
	return nil
}

// 调用每一个node的Initialize函数
func (ctx *GraphContext) initailize() error {
	for _, name := range ctx.allNodes {
		err := ctx.graphNodes[name].Initailize(ctx)
		if err != nil {
			fmt.Printf("failed to initialize node '%s'", name)
			return err
		}
	}
	fmt.Printf("initialize all nodes successfully.\r\n")
	return nil
}
